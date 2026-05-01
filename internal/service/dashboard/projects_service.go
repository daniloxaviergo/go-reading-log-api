package dashboard

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
)

// ProjectsServiceInterface defines the interface for ProjectsService
type ProjectsServiceInterface interface {
	GetRunningProjectsWithLogs(ctx context.Context) ([]*ProjectWithLogs, error)
	CalculateStats(ctx context.Context) (*dto.StatsData, error)
	// New method for JSON:API format response
	GetDashboardProjects(ctx context.Context) (*dto.DashboardProjectsResponse, error)
}

// PgxPoolInterface defines the interface for database pool operations (subset of pgxpool.Pool)
type PgxPoolInterface interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// pgxPoolInterface is an alias for internal use
type pgxPoolInterface = PgxPoolInterface

// ProjectsService handles project aggregation and log eager loading
type ProjectsService struct {
	repo   repository.DashboardRepository
	dbPool pgxPoolInterface
}

// NewProjectsService creates a new ProjectsService with the given repository and db pool
func NewProjectsService(repo repository.DashboardRepository, dbPool pgxPoolInterface) *ProjectsService {
	return &ProjectsService{repo: repo, dbPool: dbPool}
}

// SetDBPool sets the database pool for the service
func (s *ProjectsService) SetDBPool(pool pgxPoolInterface) {
	s.dbPool = pool
}

// ProjectWithLogs represents a project with its eager-loaded logs
type ProjectWithLogs struct {
	Project    *dto.ProjectAggregateResponse `json:"project"`
	Logs       []*dto.LogEntry               `json:"logs"`
	TotalPages int                           `json:"total_pages"`
	Pages      int                           `json:"pages"`
	Progress   float64                       `json:"progress_geral"`
}

// GetAll retrieves all projects with eager-loaded logs and aggregate calculations
// Returns projects ordered by progress descending with first 4 logs per project
func (s *ProjectsService) GetAll(ctx context.Context) ([]*ProjectWithLogs, error) {
	// Get all projects with their aggregates
	aggregates, err := s.repo.GetProjectsWithLogs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects with aggregates: %w", err)
	}

	// Create a map to store logs for each project
	logsMap := make(map[int64][]*dto.LogEntry)

	// Fetch logs for all projects sequentially
	for _, agg := range aggregates {
		// Get first 4 logs for this project, ordered by date DESC
		logs, err := s.repo.GetProjectLogs(ctx, agg.ProjectID, 4)
		if err != nil {
			return nil, fmt.Errorf("failed to get logs for project %d: %w", agg.ProjectID, err)
		}
		logsMap[agg.ProjectID] = logs
	}

	// Build result with aggregate calculations
	var results []*ProjectWithLogs

	for _, agg := range aggregates {
		logs := logsMap[agg.ProjectID]

		// Calculate total_pages from all project logs (not just first 4)
		totalPages, err := s.calculateTotalPages(ctx, agg.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate total pages for project %d: %w", agg.ProjectID, err)
		}

		// Calculate pages from all project logs (not just first 4)
		pages, err := s.calculatePages(ctx, agg.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate pages for project %d: %w", agg.ProjectID, err)
		}

		// Calculate progress_geral
		progress := s.calculateProgress(pages, totalPages)

		result := &ProjectWithLogs{
			Project:    agg,
			Logs:       logs,
			TotalPages: totalPages,
			Pages:      pages,
			Progress:   progress,
		}
		results = append(results, result)
	}

	// Sort by progress descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Progress > results[j].Progress
	})

	return results, nil
}

// calculateTotalPages calculates the total pages for a project from all logs
func (s *ProjectsService) calculateTotalPages(ctx context.Context, projectID int64) (int, error) {
	if s.dbPool == nil {
		return 0, fmt.Errorf("database pool not initialized")
	}

	query := `
		SELECT COALESCE(SUM(CASE 
			WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
			THEN end_page - start_page 
			ELSE 0 
		END), 0)
		FROM logs
		WHERE project_id = $1
	`

	var totalPages int
	err := s.dbPool.QueryRow(ctx, query, projectID).Scan(&totalPages)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total pages: %w", err)
	}

	return totalPages, nil
}

// calculatePages calculates the pages read for a project from all logs
func (s *ProjectsService) calculatePages(ctx context.Context, projectID int64) (int, error) {
	if s.dbPool == nil {
		return 0, fmt.Errorf("database pool not initialized")
	}

	query := `
		SELECT COALESCE(SUM(CASE 
			WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
			THEN end_page - start_page 
			ELSE 0 
		END), 0)
		FROM logs
		WHERE project_id = $1
	`

	var pages int
	err := s.dbPool.QueryRow(ctx, query, projectID).Scan(&pages)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate pages: %w", err)
	}

	return pages, nil
}

// calculateProgress calculates the progress percentage
func (s *ProjectsService) calculateProgress(pages, totalPages int) float64 {
	if totalPages <= 0 {
		return 0.0
	}
	progress := float64(pages) / float64(totalPages) * 100
	return math.Round(progress*1000) / 1000
}

// GetRunningProjectsWithLogs retrieves projects with status='running' and eager-loaded logs
// Filters by page != total_page to match Rails only_status(:running) scope
// Returns projects ordered by progress DESC, then logs.data DESC, then id ASC
// Includes first 4 logs per project ordered by data DESC
//
// Acceptance Criteria:
// #1 - Returns only projects with page != total_page (not finished)
// #2 - Filtering done in SQL query to match Rails only_status(:running)
// #3 - Progress calculated as (page/total_page)*100
// #4 - Projects ordered by progress DESC, logs.data DESC, then id ASC
// #5 - Division by zero handled returning 0.0
func (s *ProjectsService) GetRunningProjectsWithLogs(ctx context.Context) ([]*ProjectWithLogs, error) {
	// Get all running projects from repository (already filtered by page != total_page)
	dtoProjects, err := s.repo.GetRunningProjectsWithLogs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get running projects: %w", err)
	}

	// Process projects - no additional filtering needed as SQL handles it
	var runningProjects []*ProjectWithLogs

	for _, dtoProject := range dtoProjects {
		// Calculate progress using the service method (handles division by zero)
		// Acceptance Criteria #3: Progress calculated as (page/total_page)*100
		// Acceptance Criteria #5: Division by zero handled returning 0.0
		progress := s.calculateProgress(dtoProject.Pages, dtoProject.TotalPages)

		// Create service layer ProjectWithLogs from DTO
		project := &ProjectWithLogs{
			Project:    dtoProject.Project,
			Logs:       make([]*dto.LogEntry, len(dtoProject.Logs)),
			TotalPages: dtoProject.TotalPages,
			Pages:      dtoProject.Pages,
			Progress:   progress,
		}
		// Convert logs slice
		for i, log := range dtoProject.Logs {
			logCopy := log
			project.Logs[i] = &logCopy
		}

		runningProjects = append(runningProjects, project)
	}

	// Sort by progress DESC, logs.data DESC, then id ASC
	// Acceptance Criteria #4: Projects ordered by progress DESC, logs.data DESC, then id ASC
	// Note: Repository already orders by progress DESC, MAX(logs.data) DESC, id ASC
	// This secondary sort ensures consistency at service layer
	sort.Slice(runningProjects, func(i, j int) bool {
		if runningProjects[i].Progress != runningProjects[j].Progress {
			return runningProjects[i].Progress > runningProjects[j].Progress // DESC
		}
		// Compare latest log dates
		var latestLogI, latestLogJ string
		if len(runningProjects[i].Logs) > 0 {
			latestLogI = runningProjects[i].Logs[0].Data
			for _, log := range runningProjects[i].Logs {
				if log.Data > latestLogI {
					latestLogI = log.Data
				}
			}
		}
		if len(runningProjects[j].Logs) > 0 {
			latestLogJ = runningProjects[j].Logs[0].Data
			for _, log := range runningProjects[j].Logs {
				if log.Data > latestLogJ {
					latestLogJ = log.Data
				}
			}
		}
		if latestLogI != latestLogJ {
			return latestLogI > latestLogJ // DESC
		}
		return runningProjects[i].Project.ProjectID < runningProjects[j].Project.ProjectID // ASC
	})

	return runningProjects, nil
}

// CalculateStats computes aggregate statistics across all projects
// Returns StatsData with total_pages, pages, and progress_geral
// Edge cases: zero projects returns all zeros, division by zero returns 0.0
//
// Acceptance Criteria:
// #1 - stats.total_pages equals sum of all project total_page values
// #2 - stats.pages equals sum of all project page values
// #3 - stats.progress_geral calculated as round((pages/total_pages)*100, 3)
// #4 - Zero projects returns stats with all values at 0
// #5 - Division by zero returns 0.0 for progress_geral
func (s *ProjectsService) CalculateStats(ctx context.Context) (*dto.StatsData, error) {
	// Fetch all project aggregates from repository
	aggregates, err := s.repo.GetProjectAggregates(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get project aggregates: %w", err)
	}

	// Edge case #4: Zero projects returns stats with all values at 0
	if len(aggregates) == 0 {
		return dto.NewStatsData(), nil
	}

	// Calculate aggregates
	var totalPages int
	var pages int

	for _, agg := range aggregates {
		// #1: Sum all TotalPage values (project's total_page field)
		totalPages += agg.TotalPage

		// #2: Query the page field from projects table for each project
		// Note: GetProjectAggregates returns TotalPage but not the page field
		var projectPage int
		err := s.dbPool.QueryRow(ctx, "SELECT page FROM projects WHERE id = $1", agg.ProjectID).Scan(&projectPage)
		if err != nil {
			if err == pgx.ErrNoRows {
				// Project not found, skip (should not happen if aggregates are correct)
				continue
			}
			return nil, fmt.Errorf("failed to get page for project %d: %w", agg.ProjectID, err)
		}
		pages += projectPage
	}

	// Calculate progress_geral with division by zero protection
	// #5: Division by zero returns 0.0 for progress_geral
	progressGeral := s.calculateProgress(pages, totalPages)

	// #3: Round to 3 decimal places
	progressGeral = math.Round(progressGeral*1000) / 1000

	// Create and populate StatsData
	stats := dto.NewStatsData()
	stats.SetTotalPages(totalPages)
	stats.SetPages(pages)
	stats.SetProgressGeral(progressGeral)

	return stats, nil
}

// GetDashboardProjects returns projects in JSON:API format matching Rails response
// Response structure: { "data": [...], "stats": {...} }
// Each project: { "id": "123", "type": "projects", "attributes": {...} }
// Attributes use kebab-case field names
func (s *ProjectsService) GetDashboardProjects(ctx context.Context) (*dto.DashboardProjectsResponse, error) {
	response := dto.NewDashboardProjectsResponse()

	// Get running projects with logs
	projects, err := s.GetRunningProjectsWithLogs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get running projects: %w", err)
	}

	// Convert to JSON:API format
	for _, project := range projects {
		// Get earliest log date for started-at
		var startedAt *string
		if len(project.Logs) > 0 {
			// Sort logs by date to find earliest
			earliestLog := project.Logs[0]
			for _, log := range project.Logs {
				if log.Data < earliestLog.Data {
					earliestLog = log
				}
			}
			startedAt = &earliestLog.Data
		}

		// Calculate days-unreading from latest log
		daysUnreading := 0
		if len(project.Logs) > 0 {
			// Find latest log
			latestLog := project.Logs[0]
			for _, log := range project.Logs {
				if log.Data > latestLog.Data {
					latestLog = log
				}
			}
			// Parse latest log date and calculate days since then
			if logTime, err := time.Parse(time.RFC3339, latestLog.Data); err == nil {
				today := dto.GetToday()
				daysUnreading = int(today.Sub(logTime).Hours() / 24)
				if daysUnreading < 0 {
					daysUnreading = 0
				}
			}
		}

		// Create attributes with kebab-case fields
		attributes := dto.NewDashboardProjectAttributes()
		attributes.Name = project.Project.ProjectName
		attributes.TotalPage = project.TotalPages
		attributes.Page = project.Pages
		attributes.Progress = project.Progress
		attributes.LogsCount = project.Project.LogCount
		attributes.Status = "stopped" // Default status as per task requirements
		attributes.DaysUnreading = daysUnreading

		// Set started-at if available
		if startedAt != nil {
			attributes.SetStartedAt(*startedAt)
		}

		// Create project item
		item := dto.NewDashboardProjectItem(
			strconv.FormatInt(project.Project.ProjectID, 10), // Convert int64 to string
			attributes,
		)

		response.AddProject(*item)
	}

	// Calculate stats (simplified - only keep progress_geral, total_pages, pages)
	statsData, err := s.CalculateStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate stats: %w", err)
	}

	stats := dto.NewDashboardStats()
	stats.SetProgressGeral(statsData.ProgressGeral)
	stats.SetTotalPages(statsData.TotalPages)
	stats.SetPages(statsData.Pages)

	response.SetStats(stats)

	return response, nil
}
