package integration

import (
	"context"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/service/dashboard"
)

// MockProjectsService is a mock implementation of ProjectsServiceInterface
// It returns actual data from the database to match what the real service would return
type MockProjectsService struct {
	pool *pgxpool.Pool
}

func (m *MockProjectsService) GetRunningProjectsWithLogs(ctx context.Context) ([]*dashboard.ProjectWithLogs, error) {
	if m.pool == nil {
		return []*dashboard.ProjectWithLogs{}, nil
	}

	// Query actual running projects from database
// Rails logic: running = page != total_page (not finished)
	rows, err := m.pool.Query(ctx, `
		SELECT id, name, total_page, page 
		FROM projects 
		WHERE page != total_page
	`)
	if err != nil {
		return []*dashboard.ProjectWithLogs{}, nil
	}
	defer rows.Close()

	var projects []*dashboard.ProjectWithLogs
	for rows.Next() {
		var projectID int64
		var projectName string
		var totalPage, page int
		err := rows.Scan(&projectID, &projectName, &totalPage, &page)
		if err != nil {
			continue
		}

		// Calculate progress
		progress := 0.0
		if totalPage > 0 {
			progress = math.Round(float64(page)/float64(totalPage)*100*1000) / 1000
		}

		// Get first 4 logs for this project, ordered by date DESC
		var logs []*dto.LogEntry
		logRows, err := m.pool.Query(ctx, `
			SELECT id, project_id, TO_CHAR(data, 'YYYY-MM-DD"T"HH24:MI:SS'), start_page, end_page, COALESCE(note, '')
			FROM logs 
			WHERE project_id = $1 
			ORDER BY data DESC 
			LIMIT 4
		`, projectID)
		if err == nil {
			defer logRows.Close()
			for logRows.Next() {
				var log dto.LogEntry
				var note string
				err := logRows.Scan(&log.ID, &log.ProjectID, &log.Data, &log.StartPage, &log.EndPage, &note)
				if err == nil {
					log.Note = &note
					logs = append(logs, &log)
				}
			}
		}

		project := &dashboard.ProjectWithLogs{
			Project: &dto.ProjectAggregateResponse{
				ProjectID:   projectID,
				ProjectName: projectName,
				TotalPages:  totalPage,
				LogCount:    len(logs),
				Progress:    progress,
			},
			Logs:       logs,
			TotalPages: totalPage,
			Pages:      page,
			Progress:   progress,
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (m *MockProjectsService) CalculateStats(ctx context.Context) (*dto.StatsData, error) {
	if m.pool == nil {
		return dto.NewStatsData(), nil
	}

	// Calculate actual stats from database
	stats := dto.NewStatsData()

	// Get total pages and count from all projects
	var totalCapacity, totalPage int
	err := m.pool.QueryRow(ctx, "SELECT COALESCE(SUM(total_page), 0), COALESCE(SUM(page), 0) FROM projects").Scan(&totalCapacity, &totalPage)
	if err == nil && totalCapacity > 0 {
		stats.ProgressGeral = math.Round(float64(totalPage)/float64(totalCapacity)*100*1000) / 1000
		stats.TotalPages = totalPage
	}

	return stats, nil
}

func (m *MockProjectsService) GetDashboardProjects(ctx context.Context) (*dto.DashboardProjectsResponse, error) {
	if m.pool == nil {
		return dto.NewDashboardProjectsResponse(), nil
	}

	response := dto.NewDashboardProjectsResponse()

	// Get all projects from database
	rows, err := m.pool.Query(ctx, "SELECT id, name, total_page, page FROM projects")
	if err != nil {
		return response, nil
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string
		var totalPage, page int
		if err := rows.Scan(&id, &name, &totalPage, &page); err != nil {
			continue
		}

		// Calculate progress
		progress := 0.0
		if totalPage > 0 {
			progress = math.Round(float64(page)/float64(totalPage)*100*1000) / 1000
		}

		// Get logs count
		var logsCount int
		m.pool.QueryRow(ctx, "SELECT COUNT(*) FROM logs WHERE project_id = $1", id).Scan(&logsCount)

		// Create attributes
		attributes := dto.NewDashboardProjectAttributes()
		attributes.Name = name
		attributes.TotalPage = totalPage
		attributes.Page = page
		attributes.Progress = progress
		attributes.LogsCount = logsCount
		attributes.Status = "stopped"
		attributes.DaysUnreading = 0

		// Add project to response
		response.AddProject(*dto.NewDashboardProjectItem(fmt.Sprintf("%d", id), attributes))
	}

	// Calculate stats
	var totalCapacity, totalPage int
	err = m.pool.QueryRow(ctx, "SELECT COALESCE(SUM(total_page), 0), COALESCE(SUM(page), 0) FROM projects").Scan(&totalCapacity, &totalPage)
	stats := dto.NewDashboardStats()
	stats.SetTotalPages(totalPage)
	stats.SetPages(totalPage)
	if totalCapacity > 0 {
		stats.SetProgressGeral(math.Round(float64(totalPage)/float64(totalCapacity)*100*1000) / 1000)
	}
	response.SetStats(stats)

	return response, nil
}

// NewMockProjectsService creates a new MockProjectsService with the given database pool
func NewMockProjectsService(pool *pgxpool.Pool) *MockProjectsService {
	return &MockProjectsService{pool: pool}
}
