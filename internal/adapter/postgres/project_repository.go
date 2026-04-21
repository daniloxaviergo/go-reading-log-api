package postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-reading-log-api-next/internal/config"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/domain/models"
	"go-reading-log-api-next/internal/repository"
)

// parseLogDate attempts to parse a date string using multiple formats.
// Supported formats:
//   - RFC3339 (e.g., "2024-01-15T10:30:00Z")
//   - Date only (e.g., "2024-01-15")
//   - Standard datetime (e.g., "2024-01-15 10:30:00")
func parseLogDate(dateStr string) (*time.Time, bool) {
	formats := []string{
		time.RFC3339,          // 2006-01-02T15:04:05Z
		"2006-01-02",          // YYYY-MM-DD
		"2006-01-02 15:04:05", // Standard datetime
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return &t, true
		}
	}

	return nil, false
}

const defaultContextTimeout = 15 * time.Second

// ProjectRepositoryImpl implements ProjectRepository interface using PostgreSQL
type ProjectRepositoryImpl struct {
	pool *pgxpool.Pool
}

// NewProjectRepositoryImpl creates a new ProjectRepositoryImpl with the given connection pool
func NewProjectRepositoryImpl(pool *pgxpool.Pool) *ProjectRepositoryImpl {
	return &ProjectRepositoryImpl{pool: pool}
}

// Create inserts a new project into the database and returns the created project with ID
func (r *ProjectRepositoryImpl) Create(ctx context.Context, project *models.Project) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	query := `
		INSERT INTO projects (name, total_page, started_at, page, reinicia)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, started_at, created_at, updated_at
	`

	var createdAt, updatedAt time.Time
	var startedAt *time.Time
	var id int64

	err := r.pool.QueryRow(ctx, query,
		project.Name,
		project.TotalPage,
		project.StartedAt,
		project.Page,
		project.Reinicia,
	).Scan(&id, &startedAt, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Set the ID and timestamps on the project
	project.ID = id
	project.StartedAt = startedAt
	project.CreatedAt = &createdAt
	project.UpdatedAt = &updatedAt

	return project, nil
}

// GetByID retrieves a project by its ID
func (r *ProjectRepositoryImpl) GetByID(ctx context.Context, id int64) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	query := `
		SELECT id, name, total_page, started_at, page, reinicia
		FROM projects
		WHERE id = $1
	`

	var project models.Project
	var startedAt *time.Time

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&project.ID,
		&project.Name,
		&project.TotalPage,
		&startedAt,
		&project.Page,
		&project.Reinicia,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("project with ID %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get project by ID %d: %w", id, err)
	}

	// Set optional fields
	project.StartedAt = startedAt

	return &project, nil
}

// GetAll retrieves all projects
func (r *ProjectRepositoryImpl) GetAll(ctx context.Context) ([]*models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	query := `
		SELECT id, name, total_page, started_at, page, reinicia
		FROM projects
		ORDER BY id ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []*models.Project
	for rows.Next() {
		var project models.Project
		var startedAt *time.Time

		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.TotalPage,
			&startedAt,
			&project.Page,
			&project.Reinicia,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}

		project.StartedAt = startedAt

		projects = append(projects, &project)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating projects: %w", err)
	}

	return projects, nil
}

// GetWithLogs retrieves a project with its associated logs
func (r *ProjectRepositoryImpl) GetWithLogs(ctx context.Context, id int64) (*repository.ProjectWithLogs, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	// First, get the project (only base fields, no computed fields)
	projectQuery := `
		SELECT id, name, total_page, started_at, page, reinicia
		FROM projects
		WHERE id = $1
	`

	var domainProject models.Project
	var startedAt *time.Time

	err := r.pool.QueryRow(ctx, projectQuery, id).Scan(
		&domainProject.ID,
		&domainProject.Name,
		&domainProject.TotalPage,
		&startedAt,
		&domainProject.Page,
		&domainProject.Reinicia,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("project with ID %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get project by ID %d: %w", id, err)
	}

	domainProject.StartedAt = startedAt

	// Get logs for this project ordered by data DESC
	logs, err := r.getLogsByProjectID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for project %d: %w", id, err)
	}

	// Convert logs to DTOs
	logResponses := make([]*dto.LogResponse, len(logs))
	for i, log := range logs {
		// Parse the data string to time.Time for RFC3339 compliance
		var dataTime *time.Time
		if log.Data != nil && *log.Data != "" {
			dataTime, _ = parseLogDate(*log.Data)
		}
		logResponses[i] = &dto.LogResponse{
			ID:        log.ID,
			Data:      dataTime,
			StartPage: log.StartPage,
			EndPage:   log.EndPage,
			Note:      log.Note,
		}
	}

	// Calculate derived fields
	daysUnread := domainProject.CalculateDaysUnreading(logResponses)

	// Convert to DTO
	var project dto.ProjectResponse
	var startedAtStr *string
	if startedAt != nil {
		s := startedAt.Format(time.RFC3339)
		startedAtStr = &s
	}
	project = *dto.NewProjectResponse(
		domainProject.ID,
		domainProject.Name,
		startedAtStr,
		domainProject.TotalPage,
		domainProject.Page,
	)
	logsCount := domainProject.CalculateLogsCount(logResponses)
	project.LogsCount = logsCount
	project.Status = domainProject.CalculateStatus(logResponses, config.LoadConfig())
	project.DaysUnread = daysUnread
	project.Progress = domainProject.CalculateProgress()
	project.MedianDay = domainProject.CalculateMedianDay()
	finishedAtPtr := domainProject.CalculateFinishedAt(logResponses)
	project.FinishedAt = formatTimePtr(finishedAtPtr)

	return &repository.ProjectWithLogs{
		Project: &project,
		Logs:    logResponses,
	}, nil
}

// getLogsByProjectID retrieves logs for a specific project ID ordered by data DESC
func (r *ProjectRepositoryImpl) getLogsByProjectID(ctx context.Context, projectID int64) ([]*models.Log, error) {
	// Cast timestamp columns to text to avoid binary format scanning issues
	query := `
		SELECT id, project_id, data::text as data_text, start_page, end_page, wday, note, text, created_at, updated_at
		FROM logs
		WHERE project_id = $1
		ORDER BY data DESC
	`

	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs by project ID %d: %w", projectID, err)
	}
	defer rows.Close()

	var logs []*models.Log
	for rows.Next() {
		var log models.Log
		var data *string
		var createdAt, updatedAt time.Time

		// Scan data as string (matches VARCHAR column type)
		var dataStr *string
		err := rows.Scan(
			&log.ID,
			&log.ProjectID,
			&dataStr,
			&log.StartPage,
			&log.EndPage,
			&log.Wday,
			&log.Note,
			&log.Text,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log row: %w", err)
		}

		// Convert timestamp to formatted string for JSON compatibility
		if dataStr != nil && *dataStr != "" {
			data = dataStr
		}
		log.Data = data
		log.CreatedAt = &createdAt
		log.UpdatedAt = &updatedAt

		logs = append(logs, &log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}

// GetAllWithLogs retrieves all projects with their associated logs using a single LEFT OUTER JOIN query
// The query returns projects with their logs joined, ordered by logs.data DESC
// Projects without logs are included with NULL log fields (LEFT OUTER JOIN behavior)
// Note: This matches Rails API behavior which orders by logs.data DESC
func (r *ProjectRepositoryImpl) GetAllWithLogs(ctx context.Context) ([]*repository.ProjectWithLogs, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	// Single LEFT OUTER JOIN query to fetch projects with their logs
	// Orders by logs.data DESC to match Rails eager loading behavior
	// Cast timestamp columns to text to avoid binary format scanning issues
	query := `
		SELECT 
			p.id, p.name, p.total_page, p.started_at, p.page, p.reinicia,
			l.id as log_id, l.data::text as data_text, l.start_page, l.end_page, l.note, l.wday, l.text,
			l.created_at as log_created_at, l.updated_at as log_updated_at
		FROM projects p
		LEFT OUTER JOIN logs l ON p.id = l.project_id
		ORDER BY l.data DESC NULLS LAST
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects with logs: %w", err)
	}
	defer rows.Close()

	// Scan projects and group logs from joined result
	// Since JOIN creates one row per log, project data is duplicated
	// We need to group logs by project ID
	var domainProjects []*models.Project
	var logsByProject = make(map[int64][]*dto.LogResponse)
	var seenProjectIDs = make(map[int64]bool)

	for rows.Next() {
		var project models.Project
		var startedAt *time.Time
		var logID *int64
		var logCreatedAt, logUpdatedAt *time.Time
		var wday, startPage, endPage *int

		// Scan data as string (matches VARCHAR column type)
		var dataStr *string
		var note, text *string

		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.TotalPage,
			&startedAt,
			&project.Page,
			&project.Reinicia,
			&logID,
			&dataStr,
			&startPage,
			&endPage,
			&note,
			&wday,
			&text,
			&logCreatedAt,
			&logUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan joined row: %w", err)
		}

		project.StartedAt = startedAt

		// Track if we've seen this project before (project data is duplicated in JOIN result)
		if !seenProjectIDs[project.ID] {
			domainProjects = append(domainProjects, &project)
			seenProjectIDs[project.ID] = true
		}

		// If log_id is not NULL, we have a log entry to process
		if logID != nil && *logID != 0 {
			// Parse the data string to time.Time for RFC3339 compliance
			var dataTime *time.Time
			if dataStr != nil && *dataStr != "" {
				dataTime, _ = parseLogDate(*dataStr)
			}
			logResponse := &dto.LogResponse{
				ID:   *logID,
				Data: dataTime,
				StartPage: func() int {
					if startPage == nil {
						return 0
					}
					return *startPage
				}(),
				EndPage: func() int {
					if endPage == nil {
						return 0
					}
					return *endPage
				}(),
				Note: note,
			}
			logsByProject[project.ID] = append(logsByProject[project.ID], logResponse)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating joined rows: %w", err)
	}

	// Build result
	var result []*repository.ProjectWithLogs
	for _, project := range domainProjects {
		// Get logs for this project
		logsForProject := logsByProject[project.ID]

		// Calculate derived fields
		daysUnread := project.CalculateDaysUnreading(logsForProject)
		logsCount := project.CalculateLogsCount(logsForProject)

		// Create DTO with calculated fields
		var startedAtStr *string
		if project.StartedAt != nil {
			s := project.StartedAt.Format(time.RFC3339)
			startedAtStr = &s
		}

		projectResp := dto.NewProjectResponse(
			project.ID,
			project.Name,
			startedAtStr,
			project.TotalPage,
			project.Page,
		)
		projectResp.LogsCount = logsCount
		projectResp.Status = project.CalculateStatus(logsForProject, config.LoadConfig())
		projectResp.DaysUnread = daysUnread
		projectResp.Progress = project.CalculateProgress()
		projectResp.MedianDay = project.CalculateMedianDay()
		projectResp.FinishedAt = formatTimePtr(project.CalculateFinishedAt(logsForProject))

		pw := &repository.ProjectWithLogs{
			Project: projectResp,
			Logs:    logsForProject,
		}

		result = append(result, pw)
	}

	return result, nil
}

// scanProjects scans rows into ProjectResponse DTOs
func (r *ProjectRepositoryImpl) scanProjects(rows pgx.Rows) ([]*dto.ProjectResponse, error) {
	var projects []*dto.ProjectResponse

	for rows.Next() {
		var project dto.ProjectResponse
		var startedAt, finishedAt *time.Time
		var progress *float64
		var status *string
		var logsCount, daysUnread *int
		var medianDayStr *string

		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.TotalPage,
			&startedAt,
			&project.Page,
			nil, // reinicia
			&progress,
			&status,
			&logsCount,
			&daysUnread,
			&medianDayStr,
			&finishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}

		// Convert timestamps to strings for DTO
		if startedAt != nil {
			startedAtStr := startedAt.Format(time.RFC3339)
			project.StartedAt = &startedAtStr
		}
		if finishedAt != nil {
			finishedAtStr := finishedAt.Format(time.RFC3339)
			project.FinishedAt = &finishedAtStr
		}
		// Convert median_day string to float64 if available
		if medianDayStr != nil && *medianDayStr != "" {
			if medianDay, err := strconv.ParseFloat(*medianDayStr, 64); err == nil {
				project.MedianDay = &medianDay
			}
		}
		project.LogsCount = logsCount
		project.DaysUnread = daysUnread
		project.Status = status
		project.Progress = progress

		projects = append(projects, &project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating projects: %w", err)
	}

	return projects, nil
}

// getLogsByProjectIDs retrieves logs for multiple project IDs, ordered by data DESC
func (r *ProjectRepositoryImpl) getLogsByProjectIDs(ctx context.Context, projectIDs []int64) ([]*models.Log, error) {
	query := `
		SELECT id, project_id, data, start_page, end_page, wday, note, text, created_at, updated_at
		FROM logs
		WHERE project_id = ANY($1)
		ORDER BY data DESC
	`

	rows, err := r.pool.Query(ctx, query, projectIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	var logs []*models.Log
	for rows.Next() {
		var log models.Log
		var data, note, text *string
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&log.ID,
			&log.ProjectID,
			&data,
			&log.StartPage,
			&log.EndPage,
			&log.Wday,
			&note,
			&text,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log row: %w", err)
		}

		log.Data = data
		log.Note = note
		log.Text = text
		log.CreatedAt = &createdAt
		log.UpdatedAt = &updatedAt

		logs = append(logs, &log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}

// formatTimePtr converts a time.Time pointer to a string pointer for JSON serialization
func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}
