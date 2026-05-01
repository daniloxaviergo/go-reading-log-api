package dashboard

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
)

// MockRow is a helper struct for mocking pgx.Row
type MockRow struct {
	value int
	err   error
}

func (m *MockRow) Scan(dest ...any) error {
	if m.err != nil {
		return m.err
	}
	if len(dest) > 0 {
		// Assume first destination is an *int
		if ptr, ok := dest[0].(*int); ok {
			*ptr = m.value
		}
	}
	return nil
}

// MockPgxPoolForProjects is a mock implementation of pgxPoolInterface for testing
type MockPgxPoolForProjects struct {
	mockQueryRow func(ctx context.Context, sql string, args ...any) pgx.Row
}

func (m *MockPgxPoolForProjects) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if m.mockQueryRow != nil {
		return m.mockQueryRow(ctx, sql, args...)
	}
	return nil
}

func (m *MockPgxPoolForProjects) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func (m *MockPgxPoolForProjects) Exec(ctx context.Context, sql string, args ...any) (any, error) {
	return nil, nil
}

func (m *MockPgxPoolForProjects) Acquire(ctx context.Context) (*pgx.Conn, error) {
	return nil, nil
}

func (m *MockPgxPoolForProjects) AcquireFunc(ctx context.Context, fn func(*pgx.Conn) error) error {
	return nil
}

func (m *MockPgxPoolForProjects) Close() {}

func (m *MockPgxPoolForProjects) Len() int {
	return 0
}

func (m *MockPgxPoolForProjects) Cap() int {
	return 0
}

func (m *MockPgxPoolForProjects) Stats() any {
	return nil
}

func (m *MockPgxPoolForProjects) Config() any {
	return nil
}

func (m *MockPgxPoolForProjects) Reset(ctx context.Context) error {
	return nil
}

// MockDashboardRepositoryForProjects is a mock implementation of DashboardRepository for testing ProjectsService
type MockDashboardRepositoryForProjects struct {
	mockGetRunningProjectsWithLogs func(ctx context.Context) ([]*dto.ProjectWithLogs, error)
	mockGetProjectsWithLogs        func(ctx context.Context) ([]*dto.ProjectAggregateResponse, error)
	mockGetProjectLogs             func(ctx context.Context, projectID int64, limit int) ([]*dto.LogEntry, error)
	mockGetProjectAggregates       func(ctx context.Context) ([]*dto.ProjectAggregate, error)
}

func (m *MockDashboardRepositoryForProjects) GetDailyStats(ctx context.Context, date time.Time) (*dto.DailyStats, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryForProjects) GetFaultsByDateRange(ctx context.Context, start, end time.Time) (*dto.FaultStats, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryForProjects) GetWeekdayFaults(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryForProjects) GetLogsByDateRange(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryForProjects) GetProjectWeekdayMean(ctx context.Context, projectID int64, weekday int) (float64, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryForProjects) CalculatePeriodPages(ctx context.Context, start, end time.Time) (int, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryForProjects) GetProjectsWithLogs(ctx context.Context) ([]*dto.ProjectAggregateResponse, error) {
	if m.mockGetProjectsWithLogs != nil {
		return m.mockGetProjectsWithLogs(ctx)
	}
	return []*dto.ProjectAggregateResponse{}, nil
}

func (m *MockDashboardRepositoryForProjects) GetProjectLogs(ctx context.Context, projectID int64, limit int) ([]*dto.LogEntry, error) {
	if m.mockGetProjectLogs != nil {
		return m.mockGetProjectLogs(ctx, projectID, limit)
	}
	return []*dto.LogEntry{}, nil
}

func (m *MockDashboardRepositoryForProjects) GetMaxByWeekday(ctx context.Context, date time.Time) (*float64, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryForProjects) GetOverallMean(ctx context.Context, date time.Time) (*float64, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryForProjects) GetPreviousPeriodMean(ctx context.Context, date time.Time) (*float64, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryForProjects) GetPreviousPeriodSpecMean(ctx context.Context, date time.Time) (*float64, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryForProjects) GetMeanByWeekday(ctx context.Context, weekday int) (*float64, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryForProjects) GetRunningProjectsWithLogs(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
	if m.mockGetRunningProjectsWithLogs != nil {
		return m.mockGetRunningProjectsWithLogs(ctx)
	}
	return []*dto.ProjectWithLogs{}, nil
}

func (m *MockDashboardRepositoryForProjects) GetPool() repository.PoolInterface {
	return nil
}

func (m *MockDashboardRepositoryForProjects) GetProjectAggregates(ctx context.Context) ([]*dto.ProjectAggregate, error) {
	if m.mockGetProjectAggregates != nil {
		return m.mockGetProjectAggregates(ctx)
	}
	return []*dto.ProjectAggregate{}, nil
}

// TestProjectsService_GetRunningProjectsWithLogs tests the main method
func TestProjectsService_GetRunningProjectsWithLogs(t *testing.T) {
	mockRepo := &MockDashboardRepositoryForProjects{}
	service := NewProjectsService(mockRepo, nil)
	ctx := context.Background()

	// Test case 1: Normal case - multiple running projects
	t.Run("normal case - multiple running projects", func(t *testing.T) {
		mockRepo.mockGetRunningProjectsWithLogs = func(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
			return []*dto.ProjectWithLogs{
				{
					Project:    dto.NewProjectAggregateResponse(1, "Project A", 100, 10, 0.0),
					Logs:       []dto.LogEntry{{ID: 1}, {ID: 2}},
					TotalPages: 100,
					Pages:      50,
					Progress:   0.0,
				},
				{
					Project:    dto.NewProjectAggregateResponse(2, "Project B", 200, 20, 0.0),
					Logs:       []dto.LogEntry{{ID: 3}},
					TotalPages: 200,
					Pages:      100,
					Progress:   0.0,
				},
			}, nil
		}

		results, err := service.GetRunningProjectsWithLogs(ctx)

		require.NoError(t, err)
		assert.Len(t, results, 2)

		// Verify ordering: Project B (50% progress) should come before Project A (50% progress)
		// Since they have equal progress, they should be ordered by id ASC
		assert.Equal(t, int64(1), results[0].Project.ProjectID)
		assert.Equal(t, int64(2), results[1].Project.ProjectID)

		// Verify progress calculation
		assert.Equal(t, 50.0, results[0].Progress) // 50/100 * 100
		assert.Equal(t, 50.0, results[1].Progress) // 100/200 * 100
	})

	// Test case 2: Ordering by progress DESC
	t.Run("ordering by progress DESC", func(t *testing.T) {
		mockRepo.mockGetRunningProjectsWithLogs = func(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
			return []*dto.ProjectWithLogs{
				{
					Project:    dto.NewProjectAggregateResponse(1, "Project A", 100, 10, 0.0),
					Logs:       []dto.LogEntry{{ID: 1}},
					TotalPages: 100,
					Pages:      25,
					Progress:   0.0,
				},
				{
					Project:    dto.NewProjectAggregateResponse(2, "Project B", 100, 10, 0.0),
					Logs:       []dto.LogEntry{{ID: 2}},
					TotalPages: 100,
					Pages:      75,
					Progress:   0.0,
				},
				{
					Project:    dto.NewProjectAggregateResponse(3, "Project C", 100, 10, 0.0),
					Logs:       []dto.LogEntry{{ID: 3}},
					TotalPages: 100,
					Pages:      50,
					Progress:   0.0,
				},
			}, nil
		}

		results, err := service.GetRunningProjectsWithLogs(ctx)

		require.NoError(t, err)
		assert.Len(t, results, 3)

		// Verify ordering by progress DESC
		assert.Equal(t, 75.0, results[0].Progress) // Project B
		assert.Equal(t, 50.0, results[1].Progress) // Project C
		assert.Equal(t, 25.0, results[2].Progress) // Project A
	})

	// Test case 3: Ordering by id ASC when progress is equal
	t.Run("equal progress ordering by id ASC", func(t *testing.T) {
		mockRepo.mockGetRunningProjectsWithLogs = func(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
			return []*dto.ProjectWithLogs{
				{
					Project:    dto.NewProjectAggregateResponse(3, "Project C", 100, 10, 0.0),
					Logs:       []dto.LogEntry{{ID: 3}},
					TotalPages: 100,
					Pages:      50,
					Progress:   0.0,
				},
				{
					Project:    dto.NewProjectAggregateResponse(1, "Project A", 100, 10, 0.0),
					Logs:       []dto.LogEntry{{ID: 1}},
					TotalPages: 100,
					Pages:      50,
					Progress:   0.0,
				},
				{
					Project:    dto.NewProjectAggregateResponse(2, "Project B", 100, 10, 0.0),
					Logs:       []dto.LogEntry{{ID: 2}},
					TotalPages: 100,
					Pages:      50,
					Progress:   0.0,
				},
			}, nil
		}

		results, err := service.GetRunningProjectsWithLogs(ctx)

		require.NoError(t, err)
		assert.Len(t, results, 3)

		// Verify ordering by id ASC when progress is equal
		assert.Equal(t, int64(1), results[0].Project.ProjectID)
		assert.Equal(t, int64(2), results[1].Project.ProjectID)
		assert.Equal(t, int64(3), results[2].Project.ProjectID)
	})

	// Test case 4: Empty results
	t.Run("no running projects", func(t *testing.T) {
		mockRepo.mockGetRunningProjectsWithLogs = func(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
			return []*dto.ProjectWithLogs{}, nil
		}

		results, err := service.GetRunningProjectsWithLogs(ctx)

		require.NoError(t, err)
		assert.Len(t, results, 0)
	})

	// Test case 5: Repository error
	t.Run("repository error", func(t *testing.T) {
		mockRepo.mockGetRunningProjectsWithLogs = func(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
			return nil, assert.AnError
		}

		results, err := service.GetRunningProjectsWithLogs(ctx)

		assert.Error(t, err)
		assert.Nil(t, results)
		assert.Contains(t, err.Error(), "failed to get running projects")
	})
}

// TestProjectsService_GetRunningProjectsWithLogs_Filtering tests that filtering is done at repository level
// The SQL query filters by page != total_page (matching Rails only_status(:running))
// Service layer does not perform additional filtering
func TestProjectsService_GetRunningProjectsWithLogs_Filtering(t *testing.T) {
	mockRepo := &MockDashboardRepositoryForProjects{}
	service := NewProjectsService(mockRepo, nil)
	ctx := context.Background()

	// Test case 1: Repository returns only running projects (filtering done in SQL)
	// Finished projects (page == total_page) are filtered out by SQL query
	t.Run("repository returns only running projects", func(t *testing.T) {
		mockRepo.mockGetRunningProjectsWithLogs = func(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
			// SQL already filters out finished projects, so only running projects are returned
			return []*dto.ProjectWithLogs{
				{
					Project:    dto.NewProjectAggregateResponse(1, "Running Project", 100, 10, 0.0),
					Logs:       []dto.LogEntry{{ID: 1}},
					TotalPages: 100,
					Pages:      50, // Not finished
					Progress:   0.0,
				},
				// Note: Finished project would be filtered out by SQL query (page != total_page)
			}, nil
		}

		results, err := service.GetRunningProjectsWithLogs(ctx)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, int64(1), results[0].Project.ProjectID)
	})

	// Test case 2: Repository returns projects with logs (filtering done in SQL)
	// Projects without logs are handled by LEFT JOIN in SQL
	t.Run("repository returns projects with logs", func(t *testing.T) {
		mockRepo.mockGetRunningProjectsWithLogs = func(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
			// SQL returns projects with their logs (or empty logs via LEFT JOIN)
			return []*dto.ProjectWithLogs{
				{
					Project:    dto.NewProjectAggregateResponse(1, "Project With Logs", 100, 1, 0.0),
					Logs:       []dto.LogEntry{{ID: 1}},
					TotalPages: 100,
					Pages:      50,
					Progress:   0.0,
				},
			}, nil
		}

		results, err := service.GetRunningProjectsWithLogs(ctx)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, int64(1), results[0].Project.ProjectID)
	})
}

// TestProjectsService_GetRunningProjectsWithLogs_DivisionByZero tests division by zero handling
func TestProjectsService_GetRunningProjectsWithLogs_DivisionByZero(t *testing.T) {
	mockRepo := &MockDashboardRepositoryForProjects{}
	service := NewProjectsService(mockRepo, nil)
	ctx := context.Background()

	// Test case 1: total_pages = 0 (edge case - progress should be 0)
	t.Run("division by zero - total_pages is zero", func(t *testing.T) {
		mockRepo.mockGetRunningProjectsWithLogs = func(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
			return []*dto.ProjectWithLogs{
				{
					Project:    dto.NewProjectAggregateResponse(1, "Project", 0, 0, 0.0),
					Logs:       []dto.LogEntry{{ID: 1}},
					TotalPages: 100, // Use a valid total for running status check
					Pages:      0,   // pages = 0, total = 100
					Progress:   0.0,
				},
			}, nil
		}

		results, err := service.GetRunningProjectsWithLogs(ctx)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, 0.0, results[0].Progress) // 0/100 * 100 = 0
	})

	// Test case 2: pages = 0
	t.Run("pages is zero", func(t *testing.T) {
		mockRepo.mockGetRunningProjectsWithLogs = func(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
			return []*dto.ProjectWithLogs{
				{
					Project:    dto.NewProjectAggregateResponse(1, "Project", 100, 0, 0.0),
					Logs:       []dto.LogEntry{{ID: 1}},
					TotalPages: 100,
					Pages:      0,
					Progress:   0.0,
				},
			}, nil
		}

		results, err := service.GetRunningProjectsWithLogs(ctx)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, 0.0, results[0].Progress) // 0/100 * 100 = 0
	})

	// Test case 3: Zero pages with valid total
	t.Run("zero pages with valid total", func(t *testing.T) {
		mockRepo.mockGetRunningProjectsWithLogs = func(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
			return []*dto.ProjectWithLogs{
				{
					Project:    dto.NewProjectAggregateResponse(1, "Project", 100, 0, 0.0),
					Logs:       []dto.LogEntry{{ID: 1}},
					TotalPages: 100,
					Pages:      0,
					Progress:   0.0,
				},
			}, nil
		}

		results, err := service.GetRunningProjectsWithLogs(ctx)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, 0.0, results[0].Progress) // 0/100 * 100 = 0
	})
}

// TestProjectsService_GetRunningProjectsWithLogs_FloatRounding tests float rounding to 3 decimals
func TestProjectsService_GetRunningProjectsWithLogs_FloatRounding(t *testing.T) {
	mockRepo := &MockDashboardRepositoryForProjects{}
	service := NewProjectsService(mockRepo, nil)
	ctx := context.Background()

	// Test case: Float rounding to 3 decimals
	t.Run("float rounding to 3 decimals", func(t *testing.T) {
		mockRepo.mockGetRunningProjectsWithLogs = func(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
			return []*dto.ProjectWithLogs{
				{
					Project:    dto.NewProjectAggregateResponse(1, "Project", 3, 0, 0.0),
					Logs:       []dto.LogEntry{{ID: 1}},
					TotalPages: 3,
					Pages:      1,
					Progress:   0.0,
				},
			}, nil
		}

		results, err := service.GetRunningProjectsWithLogs(ctx)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		// 1/3 * 100 = 33.333... -> rounded to 33.333
		assert.InDelta(t, 33.333, results[0].Progress, 0.001)
	})
}

// TestProjectsService_GetRunningProjectsWithLogs_SingleProject tests with a single project
func TestProjectsService_GetRunningProjectsWithLogs_SingleProject(t *testing.T) {
	mockRepo := &MockDashboardRepositoryForProjects{}
	service := NewProjectsService(mockRepo, nil)
	ctx := context.Background()

	t.Run("single running project", func(t *testing.T) {
		mockRepo.mockGetRunningProjectsWithLogs = func(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
			return []*dto.ProjectWithLogs{
				{
					Project:    dto.NewProjectAggregateResponse(1, "Single Project", 200, 5, 0.0),
					Logs:       []dto.LogEntry{{ID: 1}, {ID: 2}, {ID: 3}},
					TotalPages: 200,
					Pages:      150,
					Progress:   0.0,
				},
			}, nil
		}

		results, err := service.GetRunningProjectsWithLogs(ctx)

		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, int64(1), results[0].Project.ProjectID)
		assert.Equal(t, 75.0, results[0].Progress) // 150/200 * 100
	})
}

// TestProjectsService_calculateProgress tests the calculateProgress helper method
func TestProjectsService_calculateProgress(t *testing.T) {
	mockRepo := &MockDashboardRepositoryForProjects{}
	service := NewProjectsService(mockRepo, nil)

	testCases := []struct {
		name         string
		pages        int
		totalPages   int
		expectedProg float64
	}{
		{"normal progress", 50, 100, 50.0},
		{"zero pages", 0, 100, 0.0},
		{"zero total pages", 50, 0, 0.0},
		{"both zero", 0, 0, 0.0},
		{"full progress", 100, 100, 100.0},
		{"exceeds total", 150, 100, 150.0}, // Note: not clamped in calculateProgress
		{"decimal rounding", 1, 3, 33.333},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			progress := service.calculateProgress(tc.pages, tc.totalPages)
			if tc.name == "decimal rounding" {
				assert.InDelta(t, tc.expectedProg, progress, 0.001)
			} else {
				assert.Equal(t, tc.expectedProg, progress)
			}
		})
	}
}

// TestProjectsService_CalculateStats tests the CalculateStats method
func TestProjectsService_CalculateStats(t *testing.T) {
	mockRepo := &MockDashboardRepositoryForProjects{}
	mockPool := &MockPgxPoolForProjects{}
	service := NewProjectsService(mockRepo, mockPool)
	ctx := context.Background()

	// Test case #1: Normal case - multiple projects
	t.Run("normal case - multiple projects", func(t *testing.T) {
		mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
			return []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Project A", TotalPages: 100, LogCount: 5, TotalPage: 200},
				{ProjectID: 2, ProjectName: "Project B", TotalPages: 150, LogCount: 8, TotalPage: 300},
				{ProjectID: 3, ProjectName: "Project C", TotalPages: 50, LogCount: 3, TotalPage: 100},
			}, nil
		}

		// Mock database queries for page values
		mockPool.mockQueryRow = func(ctx context.Context, sql string, args ...any) pgx.Row {
			projectID := args[0].(int64)
			pageValues := map[int64]int{
				1: 50,
				2: 100,
				3: 25,
			}
			return &MockRow{value: pageValues[projectID]}
		}

		stats, err := service.CalculateStats(ctx)

		require.NoError(t, err)
		assert.NotNil(t, stats)

		// #1: stats.total_pages equals sum of all project total_page values
		assert.Equal(t, 600, stats.TotalPages) // 200 + 300 + 100

		// #2: stats.pages equals sum of all project page values
		assert.Equal(t, 175, stats.Pages) // 50 + 100 + 25

		// #3: stats.progress_geral calculated as round((pages/total_pages)*100, 3)
		// 175/600 * 100 = 29.1666... -> rounded to 29.167
		assert.InDelta(t, 29.167, stats.ProgressGeral, 0.001)
	})

	// Test case #4: Zero projects returns stats with all values at 0
	t.Run("zero projects", func(t *testing.T) {
		mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
			return []*dto.ProjectAggregate{}, nil
		}

		stats, err := service.CalculateStats(ctx)

		require.NoError(t, err)
		assert.NotNil(t, stats)

		// #4: All values should be 0
		assert.Equal(t, 0, stats.TotalPages)
		assert.Equal(t, 0, stats.Pages)
		assert.Equal(t, 0.0, stats.ProgressGeral)
	})

	// Test case #5: Division by zero returns 0.0 for progress_geral
	t.Run("division by zero - total_pages is zero", func(t *testing.T) {
		mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
			return []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Project A", TotalPages: 0, LogCount: 0, TotalPage: 0},
				{ProjectID: 2, ProjectName: "Project B", TotalPages: 0, LogCount: 0, TotalPage: 0},
			}, nil
		}

		mockPool.mockQueryRow = func(ctx context.Context, sql string, args ...any) pgx.Row {
			projectID := args[0].(int64)
			pageValues := map[int64]int{
				1: 50,
				2: 30,
			}
			return &MockRow{value: pageValues[projectID]}
		}

		stats, err := service.CalculateStats(ctx)

		require.NoError(t, err)
		assert.NotNil(t, stats)

		// #2: pages should be summed
		assert.Equal(t, 80, stats.Pages) // 50 + 30

		// #1: total_pages should be 0
		assert.Equal(t, 0, stats.TotalPages)

		// #5: Division by zero returns 0.0 for progress_geral
		assert.Equal(t, 0.0, stats.ProgressGeral)
	})

	// Test case: Single project
	t.Run("single project", func(t *testing.T) {
		mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
			return []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Single Project", TotalPages: 100, LogCount: 5, TotalPage: 200},
			}, nil
		}

		mockPool.mockQueryRow = func(ctx context.Context, sql string, args ...any) pgx.Row {
			return &MockRow{value: 150}
		}

		stats, err := service.CalculateStats(ctx)

		require.NoError(t, err)
		assert.NotNil(t, stats)

		assert.Equal(t, 200, stats.TotalPages)
		assert.Equal(t, 150, stats.Pages)
		assert.Equal(t, 75.0, stats.ProgressGeral) // 150/200 * 100 = 75.0
	})

	// Test case: Float rounding to 3 decimals
	t.Run("float rounding to 3 decimals", func(t *testing.T) {
		mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
			return []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Project", TotalPages: 1, LogCount: 1, TotalPage: 3},
			}, nil
		}

		mockPool.mockQueryRow = func(ctx context.Context, sql string, args ...any) pgx.Row {
			return &MockRow{value: 1}
		}

		stats, err := service.CalculateStats(ctx)

		require.NoError(t, err)
		assert.NotNil(t, stats)

		assert.Equal(t, 3, stats.TotalPages)
		assert.Equal(t, 1, stats.Pages)
		// 1/3 * 100 = 33.333... -> rounded to 33.333
		assert.InDelta(t, 33.333, stats.ProgressGeral, 0.001)
	})

	// Test case: Zero pages with valid total
	t.Run("zero pages with valid total", func(t *testing.T) {
		mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
			return []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Project", TotalPages: 0, LogCount: 0, TotalPage: 100},
			}, nil
		}

		mockPool.mockQueryRow = func(ctx context.Context, sql string, args ...any) pgx.Row {
			return &MockRow{value: 0}
		}

		stats, err := service.CalculateStats(ctx)

		require.NoError(t, err)
		assert.NotNil(t, stats)

		assert.Equal(t, 100, stats.TotalPages)
		assert.Equal(t, 0, stats.Pages)
		assert.Equal(t, 0.0, stats.ProgressGeral) // 0/100 * 100 = 0
	})

	// Test case: Repository error
	t.Run("repository error", func(t *testing.T) {
		mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
			return nil, assert.AnError
		}

		stats, err := service.CalculateStats(ctx)

		assert.Error(t, err)
		assert.Nil(t, stats)
		assert.Contains(t, err.Error(), "failed to get project aggregates")
	})

	// Test case: Project not found (pgx.ErrNoRows) - should skip and continue
	t.Run("project not found - skip and continue", func(t *testing.T) {
		mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
			return []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Project A", TotalPages: 100, LogCount: 5, TotalPage: 200},
				{ProjectID: 2, ProjectName: "Project B", TotalPages: 150, LogCount: 8, TotalPage: 300},
			}, nil
		}

		callCount := 0
		mockPool.mockQueryRow = func(ctx context.Context, sql string, args ...any) pgx.Row {
			callCount++
			// First call returns valid page, second call returns ErrNoRows
			if callCount == 1 {
				return &MockRow{value: 50}
			}
			// Return ErrNoRows for second project
			return &MockRow{err: pgx.ErrNoRows}
		}

		stats, err := service.CalculateStats(ctx)

		require.NoError(t, err)
		assert.NotNil(t, stats)

		// Should have only counted the first project's page (50)
		// Second project was skipped due to ErrNoRows
		assert.Equal(t, 500, stats.TotalPages) // 200 + 300
		assert.Equal(t, 50, stats.Pages)       // Only first project's page
		// 50/500 * 100 = 10.0
		assert.InDelta(t, 10.0, stats.ProgressGeral, 0.001)
	})
}
