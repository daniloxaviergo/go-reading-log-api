package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/domain/models"
	"go-reading-log-api-next/internal/repository"
)

const defaultContextTimeout = 5 * time.Second

// ProjectRepositoryImpl implements ProjectRepository interface using PostgreSQL
type ProjectRepositoryImpl struct {
	pool *pgxpool.Pool
}

// NewProjectRepositoryImpl creates a new ProjectRepositoryImpl with the given connection pool
func NewProjectRepositoryImpl(pool *pgxpool.Pool) *ProjectRepositoryImpl {
	return &ProjectRepositoryImpl{pool: pool}
}

// GetByID retrieves a project by its ID
func (r *ProjectRepositoryImpl) GetByID(ctx context.Context, id int64) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	query := `
		SELECT id, name, total_page, started_at, page, reinicia, progress, status, logs_count, days_unread, median_day, finished_at
		FROM projects
		WHERE id = $1
	`

	var project models.Project
	var startedAt, finishedAt, medianDay *time.Time
	var progress, status *string
	var logsCount, daysUnread *int

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&project.ID,
		&project.Name,
		&project.TotalPage,
		&startedAt,
		&project.Page,
		&project.Reinicia,
		&progress,
		&status,
		&logsCount,
		&daysUnread,
		&medianDay,
		&finishedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("project with ID %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get project by ID %d: %w", id, err)
	}

	// Set optional fields
	project.StartedAt = startedAt
	project.FinishedAt = finishedAt
	project.MedianDay = medianDay
	project.LogsCount = logsCount
	project.DaysUnread = daysUnread
	project.Status = status

	return &project, nil
}

// GetAll retrieves all projects
func (r *ProjectRepositoryImpl) GetAll(ctx context.Context) ([]*models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	query := `
		SELECT id, name, total_page, started_at, page, reinicia, progress, status, logs_count, days_unread, median_day, finished_at
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
		var startedAt, finishedAt, medianDay *time.Time
		var progress, status *string
		var logsCount, daysUnread *int

		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.TotalPage,
			&startedAt,
			&project.Page,
			&project.Reinicia,
			&progress,
			&status,
			&logsCount,
			&daysUnread,
			&medianDay,
			&finishedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}

		project.StartedAt = startedAt
		project.FinishedAt = finishedAt
		project.MedianDay = medianDay
		project.LogsCount = logsCount
		project.DaysUnread = daysUnread
		project.Status = status

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

	// First, get the project
	projectQuery := `
		SELECT id, name, total_page, started_at, page, reinicia, progress, status, logs_count, days_unread, median_day, finished_at
		FROM projects
		WHERE id = $1
	`

	var project dto.ProjectResponse
	var startedAt, finishedAt, medianDay *time.Time
	var progress *float64
	var status *string
	var logsCount, daysUnread *int

	err := r.pool.QueryRow(ctx, projectQuery, id).Scan(
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
		&medianDay,
		&finishedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("project with ID %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get project by ID %d: %w", id, err)
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
	if medianDay != nil {
		medianDayStr := medianDay.Format(time.RFC3339)
		project.MedianDay = &medianDayStr
	}
	project.LogsCount = logsCount
	project.DaysUnread = daysUnread
	project.Status = status
	project.Progress = progress

	// Get logs for this project ordered by data DESC
	logs, err := r.getLogsByProjectID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for project %d: %w", id, err)
	}

	// Convert logs to DTOs
	logResponses := make([]*dto.LogResponse, len(logs))
	for i, log := range logs {
		logResponses[i] = &dto.LogResponse{
			ID:        log.ID,
			Data:      log.Data,
			StartPage: log.StartPage,
			EndPage:   log.EndPage,
			Note:      log.Note,
		}
	}

	return &repository.ProjectWithLogs{
		Project: &project,
		Logs:    logResponses,
	}, nil
}

// getLogsByProjectID retrieves logs for a specific project ID ordered by data DESC
func (r *ProjectRepositoryImpl) getLogsByProjectID(ctx context.Context, projectID int64) ([]*models.Log, error) {
	query := `
		SELECT id, project_id, data, start_page, end_page, wday, note, text, created_at, updated_at
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

// GetAllWithLogs retrieves all projects with their associated logs ordered by logs data DESC
func (r *ProjectRepositoryImpl) GetAllWithLogs(ctx context.Context) ([]*repository.ProjectWithLogs, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultContextTimeout)
	defer cancel()

	// First, get all projects
	projectQuery := `
		SELECT id, name, total_page, started_at, page, reinicia, progress, status, logs_count, days_unread, median_day, finished_at
		FROM projects
		ORDER BY id ASC
	`

	projectRows, err := r.pool.Query(ctx, projectQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer projectRows.Close()

	projects, err := r.scanProjects(projectRows)
	if err != nil {
		return nil, err
	}

	// Get all project IDs
	var projectIDs []int64
	for _, p := range projects {
		projectIDs = append(projectIDs, p.ID)
	}

	// Get all logs for these projects, ordered by data DESC
	logs, err := r.getLogsByProjectIDs(ctx, projectIDs)
	if err != nil {
		return nil, err
	}

	// Group logs by project ID
	logsByProject := make(map[int64][]*models.Log)
	for _, log := range logs {
		logsByProject[log.ProjectID] = append(logsByProject[log.ProjectID], log)
	}

	// Build result
	var result []*repository.ProjectWithLogs
	for _, project := range projects {
		pw := &repository.ProjectWithLogs{
			Project: project,
			Logs:    make([]*dto.LogResponse, 0),
		}

		// Convert logs to DTOs
		for _, log := range logsByProject[project.ID] {
			logResponse := &dto.LogResponse{
				ID:        log.ID,
				Data:      log.Data,
				StartPage: log.StartPage,
				EndPage:   log.EndPage,
				Note:      log.Note,
			}
			pw.Logs = append(pw.Logs, logResponse)
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
		var startedAt, finishedAt, medianDay *time.Time
		var progress *float64
		var status *string
		var logsCount, daysUnread *int

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
			&medianDay,
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
		if medianDay != nil {
			medianDayStr := medianDay.Format(time.RFC3339)
			project.MedianDay = &medianDayStr
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
