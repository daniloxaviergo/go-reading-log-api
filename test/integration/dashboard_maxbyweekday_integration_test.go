package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-reading-log-api-next/internal/adapter/postgres"
	"go-reading-log-api-next/test"
)

// TestDashboardRepository_GetMaxByWeekday_Integration tests the GetMaxByWeekday method
// with real database interactions
func TestDashboardRepository_GetMaxByWeekday_Integration(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Setup: Create multiple projects with logs across different weekdays
	ctx := context.Background()

	// Create projects
	projectsQuery := `
		INSERT INTO projects (id, name, total_page, page)
		VALUES ($1, $2, $3, $4),
		       ($5, $6, $7, $8),
		       ($9, $10, $11, $12)
		ON CONFLICT (id) DO NOTHING
	`
	_, err = helper.Pool.Exec(ctx, projectsQuery,
		1, "Project Alpha", 300, 100,
		2, "Project Beta", 200, 50,
		3, "Project Gamma", 150, 25,
	)
	require.NoError(t, err)

	// Create logs for different weekdays
	// Monday (DOW=1): 2024-01-15, 2024-01-22
	// Tuesday (DOW=2): 2024-01-16, 2024-01-23
	// Wednesday (DOW=3): 2024-01-17
	logsData := []struct {
		projectID int64
		data      string
		startPage int
		endPage   int
	}{
		// Monday logs
		{1, "2024-01-15 10:00:00", 0, 30},  // 30 pages
		{1, "2024-01-22 14:00:00", 30, 70}, // 40 pages
		{2, "2024-01-15 09:00:00", 0, 25},  // 25 pages
		{2, "2024-01-22 11:00:00", 25, 55}, // 30 pages
		{3, "2024-01-15 16:00:00", 0, 20},  // 20 pages

		// Tuesday logs
		{1, "2024-01-16 10:00:00", 70, 100}, // 30 pages
		{2, "2024-01-16 14:00:00", 55, 80},  // 25 pages

		// Wednesday logs
		{1, "2024-01-17 09:00:00", 100, 130}, // 30 pages
	}

	for _, ld := range logsData {
		query := `
			INSERT INTO logs (project_id, data, start_page, end_page)
			VALUES ($1, $2, $3, $4)
		`
		_, err = helper.Pool.Exec(ctx, query, ld.projectID, ld.data, ld.startPage, ld.endPage)
		require.NoError(t, err)
	}

	// Test 1: Query for Monday (DOW=1) - expect max of 40 pages
	t.Run("Monday_MaxPages", func(t *testing.T) {
		monday, _ := time.Parse("2006-01-02", "2024-01-15")
		maxPages, err := repo.GetMaxByWeekday(ctx, monday)

		assert.NoError(t, err)
		assert.NotNil(t, maxPages)
		assert.Equal(t, 40.0, *maxPages) // Max of 30, 40, 25, 30, 20 = 40
	})

	// Test 2: Query for Tuesday (DOW=2) - expect max of 30 pages
	t.Run("Tuesday_MaxPages", func(t *testing.T) {
		tuesday, _ := time.Parse("2006-01-02", "2024-01-16")
		maxPages, err := repo.GetMaxByWeekday(ctx, tuesday)

		assert.NoError(t, err)
		assert.NotNil(t, maxPages)
		assert.Equal(t, 30.0, *maxPages) // Max of 30, 25 = 30
	})

	// Test 3: Query for Wednesday (DOW=3) - expect max of 30 pages
	t.Run("Wednesday_MaxPages", func(t *testing.T) {
		wednesday, _ := time.Parse("2006-01-02", "2024-01-17")
		maxPages, err := repo.GetMaxByWeekday(ctx, wednesday)

		assert.NoError(t, err)
		assert.NotNil(t, maxPages)
		assert.Equal(t, 30.0, *maxPages) // Only one entry: 30 pages
	})

	// Test 4: Query for Sunday (DOW=0) - expect nil (no data)
	t.Run("Sunday_NoData", func(t *testing.T) {
		sunday, _ := time.Parse("2006-01-02", "2024-01-14")
		maxPages, err := repo.GetMaxByWeekday(ctx, sunday)

		assert.NoError(t, err)
		assert.Nil(t, maxPages)
	})
}

// TestDashboardRepository_GetMaxByWeekday_Performance tests query performance with large dataset
func TestDashboardRepository_GetMaxByWeekday_Performance(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	// Setup: Create large dataset (1000+ logs)
	ctx := context.Background()

	// Create 10 projects
	projectsQuery := `
		INSERT INTO projects (id, name, total_page, page)
		VALUES 
		(1, 'Project 1', 1000, 0),
		(2, 'Project 2', 1000, 0),
		(3, 'Project 3', 1000, 0),
		(4, 'Project 4', 1000, 0),
		(5, 'Project 5', 1000, 0),
		(6, 'Project 6', 1000, 0),
		(7, 'Project 7', 1000, 0),
		(8, 'Project 8', 1000, 0),
		(9, 'Project 9', 1000, 0),
		(10, 'Project 10', 1000, 0)
		ON CONFLICT (id) DO NOTHING
	`
	_, err = helper.Pool.Exec(ctx, projectsQuery)
	require.NoError(t, err)

	// Create 1000 logs spread across weekdays
	// Monday dates: 2024-01-15, 2024-01-22, 2024-01-29, etc.
	mondayDates := []string{
		"2024-01-15", "2024-01-22", "2024-01-29",
		"2024-02-05", "2024-02-12", "2024-02-19",
		"2024-02-26", "2024-03-04", "2024-03-11",
		"2024-03-18", "2024-03-25", "2024-04-01",
		"2024-04-08", "2024-04-15", "2024-04-22",
	}

	for i := 0; i < 1000; i++ {
		projectID := int64((i % 10) + 1)
		dateIndex := i % len(mondayDates)
		date := mondayDates[dateIndex]
		startPage := i * 2
		endPage := startPage + 10

		query := `
			INSERT INTO logs (project_id, data, start_page, end_page)
			VALUES ($1, $2, $3, $4)
		`
		_, err = helper.Pool.Exec(ctx, query, projectID, date+" 10:00:00", startPage, endPage)
		require.NoError(t, err)
	}

	// Execute and measure performance
	monday, _ := time.Parse("2006-01-02", "2024-01-15")
	startTime := time.Now()
	maxPages, err := repo.GetMaxByWeekday(ctx, monday)
	elapsed := time.Since(startTime)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, maxPages)
	assert.Equal(t, 10.0, *maxPages) // Each entry reads 10 pages (endPage - startPage = 10)

	// Performance assertion: should complete in under 100ms
	assert.Less(t, elapsed.Milliseconds(), int64(100), "Query should complete in under 100ms")
}

// TestDashboardRepository_GetMaxByWeekday_LargePageNumbers tests with very large page numbers
func TestDashboardRepository_GetMaxByWeekday_LargePageNumbers(t *testing.T) {
	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	err = helper.SetupTestSchema()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

	ctx := context.Background()

	// Create test project
	_, err = helper.Pool.Exec(ctx, `
		INSERT INTO projects (id, name, total_page, page)
		VALUES (1, 'Test Project', 1000000, 0)
	`)
	require.NoError(t, err)

	// Create logs with large page numbers
	query := `
		INSERT INTO logs (project_id, data, start_page, end_page)
		VALUES 
		(1, '2024-01-15 10:00:00', 0, 50000),
		(1, '2024-01-22 10:00:00', 50000, 150000)
	`
	_, err = helper.Pool.Exec(ctx, query)
	require.NoError(t, err)

	// Execute
	monday, _ := time.Parse("2006-01-02", "2024-01-15")
	maxPages, err := repo.GetMaxByWeekday(ctx, monday)

	// Verify - should handle large numbers correctly
	assert.NoError(t, err)
	assert.NotNil(t, maxPages)
	assert.Equal(t, 100000.0, *maxPages) // 150000 - 50000 = 100000
}
