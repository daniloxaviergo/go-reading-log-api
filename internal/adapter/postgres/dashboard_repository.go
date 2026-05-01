package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
)

const dashboardContextTimeout = 15 * time.Second

// DashboardRepositoryImpl implements DashboardRepository interface using PostgreSQL
type DashboardRepositoryImpl struct {
	pool *pgxpool.Pool
}

// NewDashboardRepositoryImpl creates a new DashboardRepositoryImpl with the given connection pool
func NewDashboardRepositoryImpl(pool *pgxpool.Pool) *DashboardRepositoryImpl {
	return &DashboardRepositoryImpl{pool: pool}
}

// GetPool returns the underlying database connection pool
func (r *DashboardRepositoryImpl) GetPool() repository.PoolInterface {
	return r.pool
}

// GetDailyStats returns daily page statistics with weekday breakdown
// Uses COALESCE to handle NULL values from empty result sets
func (r *DashboardRepositoryImpl) GetDailyStats(ctx context.Context, date time.Time) (*dto.DailyStats, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	// Cast timestamp columns to text to avoid binary format scanning issues
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN start_page IS NOT NULL AND end_page IS NOT NULL THEN end_page - start_page ELSE 0 END), 0) as total_pages,
			COUNT(*) as log_count
		FROM logs
		WHERE data::date = $1
	`

	var stats dto.DailyStats
	err := r.pool.QueryRow(ctx, query, date).Scan(
		&stats.TotalPages,
		&stats.LogCount,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return dto.NewDailyStats(0, 0), nil // Return zero values instead of error
		}
		return nil, fmt.Errorf("failed to get daily stats: %w", err)
	}

	return &stats, nil
}

// GetProjectAggregates returns project-level sums and counts for all projects
// Joins projects with logs to calculate aggregates per project
func (r *DashboardRepositoryImpl) GetProjectAggregates(ctx context.Context) ([]*dto.ProjectAggregate, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	// Cast timestamp columns to text to avoid binary format scanning issues
	query := `
		SELECT 
			p.id as project_id,
			p.name as project_name,
			COALESCE(SUM(CASE WHEN l.start_page IS NOT NULL AND l.end_page IS NOT NULL THEN l.end_page ELSE 0 END), 0) as total_pages,
			COUNT(l.id) as log_count,
			COALESCE(p.total_page, 0) as total_page
		FROM projects p
		LEFT OUTER JOIN logs l ON p.id = l.project_id
		GROUP BY p.id, p.name, p.total_page
		ORDER BY p.id ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query project aggregates: %w", err)
	}
	defer rows.Close()

	var aggregates []*dto.ProjectAggregate
	for rows.Next() {
		var agg dto.ProjectAggregate
		err := rows.Scan(
			&agg.ProjectID,
			&agg.ProjectName,
			&agg.TotalPages,
			&agg.LogCount,
			&agg.TotalPage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project aggregate row: %w", err)
		}
		aggregates = append(aggregates, &agg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating project aggregates: %w", err)
	}

	return aggregates, nil
}

// GetFaultsByDateRange returns the count of faults within a date range
// Note: Currently counts all logs as 'faults' - adjust query based on actual fault criteria
func (r *DashboardRepositoryImpl) GetFaultsByDateRange(ctx context.Context, start, end time.Time) (*dto.FaultStats, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	// Cast timestamp columns to text to avoid binary format scanning issues
	query := `
		SELECT COUNT(*) as fault_count
		FROM logs
		WHERE data::date BETWEEN $1 AND $2
	`

	var stats dto.FaultStats
	err := r.pool.QueryRow(ctx, query, start, end).Scan(&stats.FaultCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get faults by date range: %w", err)
	}

	return &stats, nil
}

// GetWeekdayFaults returns fault distribution by weekday (0-6 = Sunday-Saturday)
// Ensures all 7 days are present in the result map with default value of 0
func (r *DashboardRepositoryImpl) GetWeekdayFaults(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	// Cast timestamp columns to text to avoid binary format scanning issues
	// data is VARCHAR, so we need to cast it to TIMESTAMP for EXTRACT
	query := `
		SELECT 
			EXTRACT(DOW FROM data::timestamp)::int as weekday,
			COUNT(*) as fault_count
		FROM logs
		WHERE data::date BETWEEN $1 AND $2
		GROUP BY EXTRACT(DOW FROM data::timestamp)
		ORDER BY weekday
	`

	rows, err := r.pool.Query(ctx, query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query weekday faults: %w", err)
	}
	defer rows.Close()

	result := make(map[int]int)
	for rows.Next() {
		var weekday int
		var count int
		if err := rows.Scan(&weekday, &count); err != nil {
			return nil, fmt.Errorf("failed to scan weekday fault: %w", err)
		}
		result[weekday] = count
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating weekday faults: %w", err)
	}

	// Ensure all 7 days are present (0-6) with default value of 0
	for i := 0; i < 7; i++ {
		if _, exists := result[i]; !exists {
			result[i] = 0
		}
	}

	return dto.NewWeekdayFaults(result), nil
}

// GetLogsByDateRange returns log entries within a date range with page calculations
func (r *DashboardRepositoryImpl) GetLogsByDateRange(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	query := `
		SELECT 
			l.id,
			l.project_id,
			l.data,
			l.start_page,
			l.end_page,
			p.name as project_name,
			p.total_page,
			p.page
		FROM logs l
		LEFT JOIN projects p ON l.project_id = p.id
		WHERE l.data::date BETWEEN $1 AND $2
		ORDER BY l.data DESC
	`

	rows, err := r.pool.Query(ctx, query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs by date range: %w", err)
	}
	defer rows.Close()

	var entries []*dto.LogEntry
	for rows.Next() {
		var entry dto.LogEntry
		project := &dto.Project{}

		err := rows.Scan(
			&entry.ID,
			&entry.ProjectID,
			&entry.Data,
			&entry.StartPage,
			&entry.EndPage,
			&project.Name,
			&project.TotalPage,
			&project.Page,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log entry: %w", err)
		}

		// Calculate read pages
		entry.ReadPages = entry.EndPage - entry.StartPage

		// Set project data
		entry.Project = project

		entries = append(entries, &entry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating log entries: %w", err)
	}

	return entries, nil
}

// GetProjectWeekdayMean calculates the mean pages for a project on a specific weekday
func (r *DashboardRepositoryImpl) GetProjectWeekdayMean(ctx context.Context, projectID int64, weekday int) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	// EXTRACT(DOW FROM data::timestamp) returns 0-6 (Sunday-Saturday)
	query := `
		SELECT COALESCE(AVG(CASE 
			WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
			THEN end_page - start_page 
			ELSE 0 
		END), 0)
		FROM logs
		WHERE project_id = $1
		AND EXTRACT(DOW FROM data::timestamp)::int = $2
	`

	var mean float64
	err := r.pool.QueryRow(ctx, query, projectID, weekday).Scan(&mean)
	if err != nil {
		return 0.0, fmt.Errorf("failed to get project weekday mean: %w", err)
	}

	return mean, nil
}

// CalculatePeriodPages calculates total pages within a date range
func (r *DashboardRepositoryImpl) CalculatePeriodPages(ctx context.Context, start, end time.Time) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	query := `
		SELECT COALESCE(SUM(CASE 
			WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
			THEN end_page - start_page 
			ELSE 0 
		END), 0)
		FROM logs
		WHERE data::date BETWEEN $1 AND $2
	`

	var totalPages int
	err := r.pool.QueryRow(ctx, query, start, end).Scan(&totalPages)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate period pages: %w", err)
	}

	return totalPages, nil
}

// GetProjectsWithLogs returns all projects with eager-loaded logs (first 4 per project)
// Uses a CTE (Common Table Expression) to efficiently fetch the first 4 logs per project
func (r *DashboardRepositoryImpl) GetProjectsWithLogs(ctx context.Context) ([]*dto.ProjectAggregateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	// First, get all projects with their aggregates
	query := `
		SELECT 
			p.id as project_id,
			p.name as project_name,
			COALESCE(SUM(CASE WHEN l.start_page IS NOT NULL AND l.end_page IS NOT NULL THEN l.end_page - l.start_page ELSE 0 END), 0) as total_pages,
			COUNT(l.id) as log_count
		FROM projects p
		LEFT OUTER JOIN logs l ON p.id = l.project_id
		GROUP BY p.id, p.name
		ORDER BY p.id ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query project aggregates: %w", err)
	}
	defer rows.Close()

	var aggregates []*dto.ProjectAggregateResponse
	for rows.Next() {
		var agg dto.ProjectAggregateResponse
		err := rows.Scan(
			&agg.ProjectID,
			&agg.ProjectName,
			&agg.TotalPages,
			&agg.LogCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project aggregate row: %w", err)
		}

		// Calculate progress for this project
		if agg.TotalPages > 0 && agg.ProjectID > 0 {
			// Get the total_page from projects table for progress calculation
			var totalPage int
			err := r.pool.QueryRow(ctx, "SELECT total_page FROM projects WHERE id = $1", agg.ProjectID).Scan(&totalPage)
			if err == nil && totalPage > 0 {
				agg.Progress = math.Round((float64(agg.TotalPages)/float64(totalPage))*100*1000) / 1000
			} else {
				agg.Progress = 0.0
			}
		} else {
			agg.Progress = 0.0
		}

		aggregates = append(aggregates, &agg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating project aggregates: %w", err)
	}

	return aggregates, nil
}

// GetProjectLogs returns logs for a specific project ordered by date DESC with limit
func (r *DashboardRepositoryImpl) GetProjectLogs(ctx context.Context, projectID int64, limit int) ([]*dto.LogEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	// Cast timestamp columns to text to avoid binary format scanning issues
	query := `
		SELECT 
			l.id,
			l.project_id,
			l.data,
			l.start_page,
			l.end_page,
			p.name as project_name,
			p.total_page,
			p.page
		FROM logs l
		LEFT JOIN projects p ON l.project_id = p.id
		WHERE l.project_id = $1
		ORDER BY l.data DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, projectID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query project logs: %w", err)
	}
	defer rows.Close()

	var entries []*dto.LogEntry
	for rows.Next() {
		var entry dto.LogEntry
		project := &dto.Project{}

		err := rows.Scan(
			&entry.ID,
			&entry.ProjectID,
			&entry.Data,
			&entry.StartPage,
			&entry.EndPage,
			&project.Name,
			&project.TotalPage,
			&project.Page,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log entry: %w", err)
		}

		// Calculate read pages
		entry.ReadPages = entry.EndPage - entry.StartPage

		// Set project data
		entry.Project = project

		entries = append(entries, &entry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating log entries: %w", err)
	}

	return entries, nil
}

// GetMaxByWeekday returns the maximum pages read in a single day for the target weekday
func (r *DashboardRepositoryImpl) GetMaxByWeekday(ctx context.Context, date time.Time) (*float64, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	// Get the weekday (0-6 = Sunday-Saturday) from the target date
	weekday := int(date.Weekday())

	// Query to find the maximum pages read on any log entry for this weekday
	query := `
		SELECT MAX(CASE 
			WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
			THEN end_page - start_page 
			ELSE 0 
		END)
		FROM logs
		WHERE EXTRACT(DOW FROM data::timestamp)::int = $1
	`

	var maxPages *float64
	err := r.pool.QueryRow(ctx, query, weekday).Scan(&maxPages)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Return nil for no data
		}
		return nil, fmt.Errorf("failed to get max by weekday: %w", err)
	}

	return maxPages, nil
}

// GetOverallMean calculates the overall mean across all weekdays for the target date
func (r *DashboardRepositoryImpl) GetOverallMean(ctx context.Context, date time.Time) (*float64, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	// Calculate the overall mean by averaging pages across all logs
	query := `
		SELECT COALESCE(AVG(CASE 
			WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
			THEN end_page - start_page 
			ELSE 0 
		END), 0)
		FROM logs
	`

	var mean float64
	err := r.pool.QueryRow(ctx, query).Scan(&mean)
	if err != nil {
		return nil, fmt.Errorf("failed to get overall mean: %w", err)
	}

	// Round to 3 decimals
	rounded := math.Round(mean*1000) / 1000
	return &rounded, nil
}

// GetPreviousPeriodMean returns the mean for the same weekday 7 days prior
func (r *DashboardRepositoryImpl) GetPreviousPeriodMean(ctx context.Context, date time.Time) (*float64, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	// Calculate the date 7 days prior
	prevDate := date.AddDate(0, 0, -7)
	weekday := int(prevDate.Weekday())

	// Query to calculate mean for the same weekday 7 days prior
	query := `
		SELECT COALESCE(AVG(CASE 
			WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
			THEN end_page - start_page 
			ELSE 0 
		END), 0)
		FROM logs
		WHERE EXTRACT(DOW FROM data::timestamp)::int = $1
	`

	var mean float64
	err := r.pool.QueryRow(ctx, query, weekday).Scan(&mean)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous period mean: %w", err)
	}

	// Round to 3 decimals
	rounded := math.Round(mean*1000) / 1000
	return &rounded, nil
}

// GetPreviousPeriodSpecMean returns the speculative mean for the same weekday 7 days prior
func (r *DashboardRepositoryImpl) GetPreviousPeriodSpecMean(ctx context.Context, date time.Time) (*float64, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	// Calculate the date 7 days prior
	prevDate := date.AddDate(0, 0, -7)
	weekday := int(prevDate.Weekday())

	// Query to calculate speculative mean (mean * 1.15) for the same weekday 7 days prior
	query := `
		SELECT COALESCE(AVG(CASE 
			WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
			THEN end_page - start_page 
			ELSE 0 
		END), 0) * 1.15
		FROM logs
		WHERE EXTRACT(DOW FROM data::timestamp)::int = $1
	`

	var specMean float64
	err := r.pool.QueryRow(ctx, query, weekday).Scan(&specMean)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous period spec mean: %w", err)
	}

	// Round to 3 decimals
	rounded := math.Round(specMean*1000) / 1000
	return &rounded, nil
}

// GetMeanByWeekday calculates the mean pages per 7-day interval for a specific weekday
// Algorithm (V1::MeanLog):
// 1. Filter logs by target weekday (DOW 0-6)
// 2. Calculate total_pages = sum(end_page - start_page) for all filtered logs
// 3. Find begin_data (first log timestamp) and log_data (most recent log timestamp)
// 4. Calculate count_reads = floor((log_data - begin_data) / 7 days)
// 5. Calculate mean_day = total_pages / count_reads, rounded to 3 decimals
// 6. Edge cases: return nil if no logs or count_reads is zero
func (r *DashboardRepositoryImpl) GetMeanByWeekday(ctx context.Context, weekday int) (*float64, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	// Query to get all logs for the target weekday, ordered chronologically
	query := `
		SELECT 
			COALESCE(SUM(CASE 
				WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
				THEN end_page - start_page 
				ELSE 0 
			END), 0) as total_pages,
			MIN(data::timestamp) as begin_data,
			MAX(data::timestamp) as log_data,
			COUNT(*) as log_count
		FROM logs
		WHERE EXTRACT(DOW FROM data::timestamp)::int = $1
	`

	var totalPages int
	var beginData, logData time.Time
	var logCount int

	err := r.pool.QueryRow(ctx, query, weekday).Scan(
		&totalPages,
		&beginData,
		&logData,
		&logCount,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Return nil for no data
		}
		return nil, fmt.Errorf("failed to get mean by weekday: %w", err)
	}

	// Edge case: no logs exist
	if logCount == 0 {
		return nil, nil
	}

	// Calculate 7-day intervals
	// Use Hours() to get precise difference, then divide by 24 to get days
	daysDiff := logData.Sub(beginData).Hours() / 24
	countReads := int(daysDiff / 7)

	// Edge case: zero intervals (logs within same 7-day period)
	if countReads == 0 {
		return nil, nil
	}

	// Calculate mean and round to 3 decimals
	mean := float64(totalPages) / float64(countReads)
	rounded := math.Round(mean*1000) / 1000
	return &rounded, nil
}

// GetRunningProjectsWithLogs returns all projects with eager-loaded logs using a single SQL query
// Uses CTE with ROW_NUMBER() window function to limit logs to first 4 per project
// Filters by running status (page != total_page) to match Rails only_status(:running)
// Orders projects by progress DESC, then by logs.data DESC to match Rails order_progress.order('logs.data DESC')
// Returns projects with first 4 logs per project ordered by data DESC
//
// Acceptance Criteria:
// #1 - SQL query joins projects with logs table
// #2 - Logs limited to first 4 per project ordered by data DESC
// #3 - Progress ordering implemented via SQL CASE statement
// #4 - NULL values handled with COALESCE
// #5 - Query returns all required project and log fields
// #6 - Filters by running status (page != total_page) to match Rails only_status(:running)
// #7 - Orders by logs.data DESC as secondary sort to match Rails order('logs.data DESC')
func (r *DashboardRepositoryImpl) GetRunningProjectsWithLogs(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
	ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
	defer cancel()

	// Single query using CTE with window function to eager-load first 4 logs per project
	// This eliminates the N+1 query problem from the previous implementation
	// Filters by page != total_page to match Rails only_status(:running) scope
	// Orders by progress DESC, then logs.data DESC to match Rails order_progress.order('logs.data DESC')
	query := `
		WITH log_ranked AS (
			SELECT 
				l.id,
				l.project_id,
				l.data,
				l.start_page,
				l.end_page,
				l.note,
				ROW_NUMBER() OVER (PARTITION BY l.project_id ORDER BY l.data DESC) as rn
			FROM logs l
		)
		SELECT 
			p.id as project_id,
			p.name as project_name,
			p.total_page as project_total_page,
			p.page as project_page,
			COALESCE(SUM(CASE WHEN l.start_page IS NOT NULL AND l.end_page IS NOT NULL 
				THEN l.end_page - l.start_page ELSE 0 END), 0) as total_pages,
			COUNT(l.id) as log_count,
			lr.id as log_id,
			lr.project_id as log_project_id,
			lr.data as log_data,
			lr.start_page as log_start_page,
			lr.end_page as log_end_page,
			lr.note as log_note
		FROM projects p
		LEFT JOIN log_ranked lr ON p.id = lr.project_id AND lr.rn <= 4
		LEFT JOIN logs l ON p.id = l.project_id
		WHERE p.page != p.total_page
		GROUP BY p.id, p.name, p.total_page, p.page, 
				 lr.id, lr.project_id, lr.data, lr.start_page, lr.end_page, lr.note
		ORDER BY 
			CASE 
				WHEN p.total_page = 0 THEN 0 
				ELSE p.page::float / p.total_page::float 
			END DESC,
			MAX(lr.data) DESC,
			p.id ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query running projects with logs: %w", err)
	}
	defer rows.Close()

	// Map to aggregate project data and collect logs
	projectMap := make(map[int64]*dto.ProjectWithLogs)

	for rows.Next() {
		var projectID int64
		var projectName string
		var projectTotalPage int
		var projectPage int
		var totalPages int
		var logCount int
		var logID sql.NullInt64
		var logProjectID sql.NullInt64
		var logData sql.NullString
		var logStartPage sql.NullInt32
		var logEndPage sql.NullInt32
		var logNote sql.NullString

		err := rows.Scan(
			&projectID,
			&projectName,
			&projectTotalPage,
			&projectPage,
			&totalPages,
			&logCount,
			&logID,
			&logProjectID,
			&logData,
			&logStartPage,
			&logEndPage,
			&logNote,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project with log row: %w", err)
		}

		// Initialize project if not exists
		if _, exists := projectMap[projectID]; !exists {
			projectMap[projectID] = &dto.ProjectWithLogs{
				Project: dto.NewProjectAggregateResponse(
					projectID,
					projectName,
					totalPages,
					logCount,
					0.0, // Progress will be calculated in service layer
				),
				Logs:       make([]dto.LogEntry, 0),
				TotalPages: totalPages,
				Pages:      projectPage,
				Progress:   0.0,
			}
		}

		// Add log entry if present (LEFT JOIN may return NULL for projects without logs)
		if logID.Valid {
			var note *string
			if logNote.Valid {
				note = &logNote.String
			}

			project := projectMap[projectID]
			log := dto.LogEntry{
				ID:        logID.Int64,
				ProjectID: logProjectID.Int64,
				Data:      logData.String,
				StartPage: int(logStartPage.Int32),
				EndPage:   int(logEndPage.Int32),
				Note:      note,
				ReadPages: int(logEndPage.Int32) - int(logStartPage.Int32),
			}
			project.Logs = append(project.Logs, log)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating project with logs: %w", err)
	}

	// Convert map to slice
	projects := make([]*dto.ProjectWithLogs, 0, len(projectMap))
	for _, project := range projectMap {
		projects = append(projects, project)
	}

	return projects, nil
}
