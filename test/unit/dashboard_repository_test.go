package unit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-reading-log-api-next/internal/domain/dto"

	"go-reading-log-api-next/internal/adapter/postgres"
	"go-reading-log-api-next/test"
)

// TestDashboardRepository_GetDailyStats tests the GetDailyStats method
func TestDashboardRepository_GetDailyStats(t *testing.T) {
	// Setup test database
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	// Setup test schema (create tables)
	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Create test data - use project 1 which is created by default in test setup
	testDateStr := "2024-01-15"

	// Ensure project 1 exists before creating logs
	ctx := context.Background()
	query := `
		INSERT INTO projects (id, name, total_page, page)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`
	_, err = helper.Pool.Exec(ctx, query, 1, "Test Project", 100, 0)
	require.NoError(t, err)

	err = createTestLogs(helper.Pool, []testLog{
		{ProjectID: 1, Data: &testDateStr, StartPage: 10, EndPage: 20},
		{ProjectID: 1, Data: &testDateStr, StartPage: 20, EndPage: 30},
	})
	require.NoError(t, err)

	// Execute
	testDate, _ := time.Parse("2006-01-02", testDateStr)
	stats, err := repo.GetDailyStats(context.Background(), testDate)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 20, stats.TotalPages) // (20-10) + (30-20) = 20
	assert.Equal(t, 2, stats.LogCount)
}

// TestDashboardRepository_GetDailyStats_EmptyDate tests with no data for the date
func TestDashboardRepository_GetDailyStats_EmptyDate(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	// Setup test schema (create tables)
	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Use a date with no test data
	testDateStr := "2024-01-20"
	testDate, _ := time.Parse("2006-01-02", testDateStr)

	// Execute
	stats, err := repo.GetDailyStats(context.Background(), testDate)

	// Verify - should return zero values, not error
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 0, stats.TotalPages)
	assert.Equal(t, 0, stats.LogCount)
}

// TestDashboardRepository_GetProjectAggregates tests the GetProjectAggregates method
func TestDashboardRepository_GetProjectAggregates(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	// Setup test schema (create tables)
	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Create test projects and logs
	testDateStr := "2024-01-15"

	// First, create the projects that will be referenced by logs
	ctx := context.Background()
	query := `
		INSERT INTO projects (id, name, total_page, page)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`

	// Create Project 1
	_, err = helper.Pool.Exec(ctx, query, 1, "Project 1", 100, 0)
	require.NoError(t, err)

	// Create Project 2
	_, err = helper.Pool.Exec(ctx, query, 2, "Project 2", 50, 0)
	require.NoError(t, err)

	// Project 1 with 2 logs
	err = createTestLogs(helper.Pool, []testLog{
		{ProjectID: 1, Data: &testDateStr, StartPage: 0, EndPage: 10},
		{ProjectID: 1, Data: &testDateStr, StartPage: 10, EndPage: 20},
	})
	require.NoError(t, err)

	// Project 2 with 1 log
	err = createTestLogs(helper.Pool, []testLog{
		{ProjectID: 2, Data: &testDateStr, StartPage: 5, EndPage: 15},
	})
	require.NoError(t, err)

	// Execute
	aggregates, err := repo.GetProjectAggregates(context.Background())

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, aggregates)
	assert.Len(t, aggregates, 2) // Should have 2 projects

	// Check Project 1 aggregate
	project1 := findAggregateByProjectID(aggregates, 1)
	assert.NotNil(t, project1)
	assert.Equal(t, int64(1), project1.ProjectID)
	assert.Equal(t, "Project 1", project1.ProjectName)
	assert.Equal(t, 30, project1.TotalPages) // 10 + 20 = 30 (sum of end_page)
	assert.Equal(t, 2, project1.LogCount)

	// Check Project 2 aggregate
	project2 := findAggregateByProjectID(aggregates, 2)
	assert.NotNil(t, project2)
	assert.Equal(t, int64(2), project2.ProjectID)
	assert.Equal(t, "Project 2", project2.ProjectName)
	assert.Equal(t, 15, project2.TotalPages) // 15 = sum of end_page
	assert.Equal(t, 1, project2.LogCount)
}

// TestDashboardRepository_GetFaultsByDateRange tests the GetFaultsByDateRange method
func TestDashboardRepository_GetFaultsByDateRange(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	// Setup test schema (create tables)
	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Create test data spanning a date range
	startDateStr := "2024-01-15"
	endDateStr := "2024-01-22"
	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	// Create test project first
	ctx := context.Background()
	query := `
		INSERT INTO projects (id, name, total_page, page)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`
	_, err = helper.Pool.Exec(ctx, query, 1, "Test Project", 100, 0)
	require.NoError(t, err)

	err = createTestLogs(helper.Pool, []testLog{
		{ProjectID: 1, Data: &startDateStr, StartPage: 0, EndPage: 10},
		{ProjectID: 1, Data: &startDateStr, StartPage: 10, EndPage: 20},
		{ProjectID: 1, Data: &startDateStr, StartPage: 20, EndPage: 30},
	})
	require.NoError(t, err)

	// Execute
	stats, err := repo.GetFaultsByDateRange(context.Background(), startDate, endDate)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 3, stats.FaultCount) // 3 logs in date range
}

// TestDashboardRepository_GetWeekdayFaults tests the GetWeekdayFaults method
func TestDashboardRepository_GetWeekdayFaults(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	// Setup test schema (create tables)
	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Create test data spanning multiple weekdays
	// Monday is day 1 in PostgreSQL EXTRACT(DOW)
	startDateStr := "2024-01-15" // Monday
	endDateStr := "2024-01-22"   // Next Monday
	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	// Create test project first
	ctx := context.Background()
	query := `
		INSERT INTO projects (id, name, total_page, page)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`
	_, err = helper.Pool.Exec(ctx, query, 1, "Test Project", 100, 0)
	require.NoError(t, err)

	err = createTestLogs(helper.Pool, []testLog{
		{ProjectID: 1, Data: &startDateStr, StartPage: 0, EndPage: 10}, // Monday (DOW=1)
		{ProjectID: 1, Data: &startDateStr, StartPage: 0, EndPage: 10}, // Monday again (DOW=1) - same day
	})
	require.NoError(t, err)

	// Execute
	stats, err := repo.GetWeekdayFaults(context.Background(), startDate, endDate)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.NotNil(t, stats.Faults)

	// Check specific days (PostgreSQL DOW: 0=Sunday, 1=Monday, etc.)
	assert.Equal(t, 2, stats.Faults[1]) // Monday (2 logs on same day)
	assert.Equal(t, 0, stats.Faults[0]) // Sunday (no data)
	assert.Equal(t, 0, stats.Faults[3]) // Wednesday
	assert.Equal(t, 0, stats.Faults[4]) // Thursday
	assert.Equal(t, 0, stats.Faults[5]) // Friday
	assert.Equal(t, 0, stats.Faults[6]) // Saturday
}

// TestDashboardRepository_GetWeekdayFaults_EmptyRange tests with no data in range
func TestDashboardRepository_GetWeekdayFaults_EmptyRange(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	// Setup test schema (create tables)
	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Use a date range with no test data
	startDateStr := "2024-01-20"
	endDateStr := "2024-01-27"
	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	// Create test project first
	ctx := context.Background()
	query := `
		INSERT INTO projects (id, name, total_page, page)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`
	_, err = helper.Pool.Exec(ctx, query, 1, "Test Project", 100, 0)
	require.NoError(t, err)

	// Execute
	stats, err := repo.GetWeekdayFaults(context.Background(), startDate, endDate)

	// Verify - should return all zeros
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.NotNil(t, stats.Faults)
	for i := 0; i < 7; i++ {
		assert.Equal(t, 0, stats.Faults[i], "Day %d should have 0 faults", i)
	}
}

// Helper function to find aggregate by project ID
func findAggregateByProjectID(aggregates []*dto.ProjectAggregate, projectID int64) *dto.ProjectAggregate {
	for _, agg := range aggregates {
		if agg.ProjectID == projectID {
			return agg
		}
	}
	return nil
}

// TestDashboardRepository_GetMaxByWeekday tests the GetMaxByWeekday method
func TestDashboardRepository_GetMaxByWeekday(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Create test project
	ctx := context.Background()
	query := `
		INSERT INTO projects (id, name, total_page, page)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`
	_, err = helper.Pool.Exec(ctx, query, 1, "Test Project", 100, 0)
	require.NoError(t, err)

	// Create logs for Monday (weekday=1) with different page counts
	// 2024-01-15 is a Monday
	mondayDateStr1 := "2024-01-15 10:00:00"
	mondayDateStr2 := "2024-01-15 14:00:00"
	mondayDateStr3 := "2024-01-22 09:00:00" // Another Monday

	err = createTestLogs(helper.Pool, []testLog{
		{ProjectID: 1, Data: &mondayDateStr1, StartPage: 0, EndPage: 20},  // 20 pages
		{ProjectID: 1, Data: &mondayDateStr2, StartPage: 20, EndPage: 50}, // 30 pages
		{ProjectID: 1, Data: &mondayDateStr3, StartPage: 0, EndPage: 15},  // 15 pages
	})
	require.NoError(t, err)

	// Execute - query for Monday (weekday=1)
	monday, _ := time.Parse("2006-01-02", "2024-01-15")
	maxPages, err := repo.GetMaxByWeekday(context.Background(), monday)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, maxPages)
	assert.Equal(t, 30.0, *maxPages) // Max of 20, 30, 15 = 30
}

// TestDashboardRepository_GetMaxByWeekday_EmptyWeekday tests with no data for the weekday
func TestDashboardRepository_GetMaxByWeekday_EmptyWeekday(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Create test project
	ctx := context.Background()
	query := `
		INSERT INTO projects (id, name, total_page, page)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`
	_, err = helper.Pool.Exec(ctx, query, 1, "Test Project", 100, 0)
	require.NoError(t, err)

	// Create logs only for Monday (weekday=1)
	mondayDateStr := "2024-01-15 10:00:00"
	err = createTestLogs(helper.Pool, []testLog{
		{ProjectID: 1, Data: &mondayDateStr, StartPage: 0, EndPage: 20},
	})
	require.NoError(t, err)

	// Execute - query for Sunday (weekday=0) - no data
	sunday, _ := time.Parse("2006-01-02", "2024-01-14")
	maxPages, err := repo.GetMaxByWeekday(context.Background(), sunday)

	// Verify - should return nil, not error
	assert.NoError(t, err)
	assert.Nil(t, maxPages)
}

// TestDashboardRepository_GetMaxByWeekday_MultipleProjects tests across multiple projects
func TestDashboardRepository_GetMaxByWeekday_MultipleProjects(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Create test projects
	ctx := context.Background()
	query := `
		INSERT INTO projects (id, name, total_page, page)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`
	_, err = helper.Pool.Exec(ctx, query, 1, "Project 1", 100, 0)
	require.NoError(t, err)
	_, err = helper.Pool.Exec(ctx, query, 2, "Project 2", 50, 0)
	require.NoError(t, err)

	// Create logs for Monday (weekday=1) across both projects
	mondayDateStr1 := "2024-01-15 10:00:00"
	mondayDateStr2 := "2024-01-15 14:00:00"

	err = createTestLogs(helper.Pool, []testLog{
		{ProjectID: 1, Data: &mondayDateStr1, StartPage: 0, EndPage: 30}, // Project 1: 30 pages
		{ProjectID: 2, Data: &mondayDateStr2, StartPage: 0, EndPage: 50}, // Project 2: 50 pages
	})
	require.NoError(t, err)

	// Execute - query for Monday (weekday=1)
	monday, _ := time.Parse("2006-01-02", "2024-01-15")
	maxPages, err := repo.GetMaxByWeekday(context.Background(), monday)

	// Verify - should return max across all projects
	assert.NoError(t, err)
	assert.NotNil(t, maxPages)
	assert.Equal(t, 50.0, *maxPages) // Max of 30, 50 = 50
}

// TestDashboardRepository_GetMaxByWeekday_ZeroPages tests with zero pages read
func TestDashboardRepository_GetMaxByWeekday_ZeroPages(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Create test project
	ctx := context.Background()
	query := `
		INSERT INTO projects (id, name, total_page, page)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`
	_, err = helper.Pool.Exec(ctx, query, 1, "Test Project", 100, 0)
	require.NoError(t, err)

	// Create logs where end_page = start_page (0 pages read)
	mondayDateStr := "2024-01-15 10:00:00"
	err = createTestLogs(helper.Pool, []testLog{
		{ProjectID: 1, Data: &mondayDateStr, StartPage: 10, EndPage: 10}, // 0 pages
	})
	require.NoError(t, err)

	// Execute - query for Monday (weekday=1)
	monday, _ := time.Parse("2006-01-02", "2024-01-15")
	maxPages, err := repo.GetMaxByWeekday(context.Background(), monday)

	// Verify - should return 0.0
	assert.NoError(t, err)
	assert.NotNil(t, maxPages)
	assert.Equal(t, 0.0, *maxPages)
}

// TestDashboardRepository_GetMaxByWeekday_InvalidData tests edge case with start_page > end_page
func TestDashboardRepository_GetMaxByWeekday_InvalidData(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Create test project
	ctx := context.Background()
	query := `
		INSERT INTO projects (id, name, total_page, page)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`
	_, err = helper.Pool.Exec(ctx, query, 1, "Test Project", 100, 0)
	require.NoError(t, err)

	// Create logs with some invalid data (start_page > end_page)
	mondayDateStr1 := "2024-01-15 10:00:00"
	mondayDateStr2 := "2024-01-15 14:00:00"

	err = createTestLogs(helper.Pool, []testLog{
		{ProjectID: 1, Data: &mondayDateStr1, StartPage: 50, EndPage: 30}, // Invalid: -20 pages (should be 0)
		{ProjectID: 1, Data: &mondayDateStr2, StartPage: 0, EndPage: 25},  // Valid: 25 pages
	})
	require.NoError(t, err)

	// Execute - query for Monday (weekday=1)
	monday, _ := time.Parse("2006-01-02", "2024-01-15")
	maxPages, err := repo.GetMaxByWeekday(context.Background(), monday)

	// Verify - should return 25 (the valid entry), CASE handles invalid as 0
	assert.NoError(t, err)
	assert.NotNil(t, maxPages)
	assert.Equal(t, 25.0, *maxPages)
}

// TestDashboardRepository_GetMaxByWeekday_EmptyDatabase tests with no logs at all
func TestDashboardRepository_GetMaxByWeekday_EmptyDatabase(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Don't create any logs

	// Execute - query for any weekday
	monday, _ := time.Parse("2006-01-02", "2024-01-15")
	maxPages, err := repo.GetMaxByWeekday(context.Background(), monday)

	// Verify - should return nil
	assert.NoError(t, err)
	assert.Nil(t, maxPages)
}

// TestDashboardRepository_GetMaxByWeekday_SingleEntry tests with single log entry
func TestDashboardRepository_GetMaxByWeekday_SingleEntry(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Create test project
	ctx := context.Background()
	query := `
		INSERT INTO projects (id, name, total_page, page)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO NOTHING
	`
	_, err = helper.Pool.Exec(ctx, query, 1, "Test Project", 100, 0)
	require.NoError(t, err)

	// Create single log entry
	mondayDateStr := "2024-01-15 10:00:00"
	err = createTestLogs(helper.Pool, []testLog{
		{ProjectID: 1, Data: &mondayDateStr, StartPage: 0, EndPage: 42},
	})
	require.NoError(t, err)

	// Execute - query for Monday (weekday=1)
	monday, _ := time.Parse("2006-01-02", "2024-01-15")
	maxPages, err := repo.GetMaxByWeekday(context.Background(), monday)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, maxPages)
	assert.Equal(t, 42.0, *maxPages)
}

// testLog represents a log entry for testing
type testLog struct {
	ID        int64
	ProjectID int64
	Data      *string
	StartPage int
	EndPage   int
}

// createTestLogs creates test log entries in the database
func createTestLogs(pool *pgxpool.Pool, logs []testLog) error {
	ctx := context.Background()

	for _, log := range logs {
		query := `
			INSERT INTO logs (project_id, data, start_page, end_page)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`

		var id int64
		err := pool.QueryRow(ctx, query,
			log.ProjectID,
			log.Data,
			log.StartPage,
			log.EndPage,
		).Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to create test log: %w", err)
		}
	}

	return nil
}
