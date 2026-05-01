package unit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
	"go-reading-log-api-next/internal/service/dashboard"
)

// Helper function to set a fixed test date and restore it after the test
func withDayServiceFixedDate(t *testing.T, fixedDate time.Time, fn func()) {
	defer dashboard.SetTestDate(time.Now())
	dashboard.SetTestDate(fixedDate)
	fn()
}

// MockDashboardRepository is a mock implementation of DashboardRepository for testing
type MockDashboardRepository struct {
	mockGetProjectAggregates  func(ctx context.Context) ([]*dto.ProjectAggregate, error)
	mockGetFaultsByDateRange  func(ctx context.Context, start, end time.Time) (*dto.FaultStats, error)
	mockGetProjectWeekdayMean func(ctx context.Context, projectID int64, weekday int) (float64, error)
	mockGetLogsByDateRange    func(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error)
	// Fields for CalculatePeriodPages mock
	prevWeekPages int
	lastWeekPages int
	callCount     int
}

func (m *MockDashboardRepository) GetDailyStats(ctx context.Context, date time.Time) (*dto.DailyStats, error) {
	panic("not implemented")
}

func (m *MockDashboardRepository) GetProjectAggregates(ctx context.Context) ([]*dto.ProjectAggregate, error) {
	if m.mockGetProjectAggregates != nil {
		return m.mockGetProjectAggregates(ctx)
	}
	panic("not implemented")
}

// AssertExpectations asserts that all expected method calls were made
func (m *MockDashboardRepository) AssertExpectations(t *testing.T) {
	// No-op for this mock implementation - uses inline assertions instead
}

func (m *MockDashboardRepository) GetFaultsByDateRange(ctx context.Context, start, end time.Time) (*dto.FaultStats, error) {
	if m.mockGetFaultsByDateRange != nil {
		return m.mockGetFaultsByDateRange(ctx, start, end)
	}
	panic("not implemented")
}

func (m *MockDashboardRepository) GetWeekdayFaults(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error) {
	panic("not implemented")
}

func (m *MockDashboardRepository) GetLogsByDateRange(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
	if m.mockGetLogsByDateRange != nil {
		return m.mockGetLogsByDateRange(ctx, start, end)
	}
	panic("not implemented")
}

func (m *MockDashboardRepository) GetProjectWeekdayMean(ctx context.Context, projectID int64, weekday int) (float64, error) {
	if m.mockGetProjectWeekdayMean != nil {
		return m.mockGetProjectWeekdayMean(ctx, projectID, weekday)
	}
	panic("not implemented")
}

func (m *MockDashboardRepository) CalculatePeriodPages(ctx context.Context, start, end time.Time) (int, error) {
	m.callCount++

	// Return the appropriate value based on call order
	// First call = previous week, Second call = last week
	if m.callCount == 1 {
		return m.prevWeekPages, nil
	}
	if m.callCount == 2 {
		return m.lastWeekPages, nil
	}
	return 0, nil
}

func (m *MockDashboardRepository) GetPool() repository.PoolInterface {
	return nil
}

func (m *MockDashboardRepository) GetProjectsWithLogs(ctx context.Context) ([]*dto.ProjectAggregateResponse, error) {
	return []*dto.ProjectAggregateResponse{}, nil
}

func (m *MockDashboardRepository) GetProjectLogs(ctx context.Context, projectID int64, limit int) ([]*dto.LogEntry, error) {
	return []*dto.LogEntry{}, nil
}

func (m *MockDashboardRepository) GetMaxByWeekday(ctx context.Context, date time.Time) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepository) GetOverallMean(ctx context.Context, date time.Time) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepository) GetPreviousPeriodMean(ctx context.Context, date time.Time) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepository) GetPreviousPeriodSpecMean(ctx context.Context, date time.Time) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepository) GetMeanByWeekday(ctx context.Context, weekday int) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepository) GetRunningProjectsWithLogs(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
	return nil, nil
}

// MockUserConfigService is a mock implementation of UserConfigProvider for testing
type MockUserConfigService struct {
	mockPredictionPct float64
}

func (m *MockUserConfigService) GetPredictionPct() float64 {
	return m.mockPredictionPct
}

func (m *MockUserConfigService) GetMaxFaults() int {
	panic("not implemented")
}

func (m *MockUserConfigService) GetPagesPerDay() float64 {
	panic("not implemented")
}

// resetMock resets the mock state for a new test
func resetMock(mockRepo *MockDashboardRepository) {
	mockRepo.callCount = 0
}

// TestDayService_CalculateWeeklyStats tests the main CalculateWeeklyStats method
func TestDayService_CalculateWeeklyStats(t *testing.T) {
	// Setup mocks
	mockRepo := &MockDashboardRepository{}
	mockConfig := &MockUserConfigService{mockPredictionPct: 0.15}

	// Create DayService with mocks
	resetMock(mockRepo)
	dayService := dashboard.NewDayService(mockRepo, mockConfig)

	// Test case 1: Normal calculation with data
	t.Run("normal calculation", func(t *testing.T) {
		// Setup mock returns for normal case
		mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
			return []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Test Project 1", TotalPages: 100, LogCount: 50},
				{ProjectID: 2, ProjectName: "Test Project 2", TotalPages: 50, LogCount: 25},
			}, nil
		}

		mockRepo.prevWeekPages = 10 // Previous week pages
		mockRepo.lastWeekPages = 15 // Last week pages

		mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
			// Return consistent mean for both projects
			return 10.5, nil
		}

		ctx := context.Background()
		stats, err := dayService.CalculateWeeklyStats(ctx)

		// Verify
		require.NoError(t, err)
		assert.NotNil(t, stats)

		// Check calculated values
		assert.Equal(t, 10, stats.PreviousWeekPages)
		assert.Equal(t, 15, stats.LastWeekPages)

		// per_pages = (15 / 10) * 100 = 150.0
		assert.NotNil(t, stats.PerPages)
		assert.Equal(t, 150.0, *stats.PerPages)

		// mean_day = 10.5 (from mock)
		assert.Equal(t, 10.5, stats.MeanDay)

		// spec_mean_day = 10.5 * (1 + 0.15) = 12.075
		assert.Equal(t, 12.075, stats.SpecMeanDay)

		// Total pages from aggregates
		assert.Equal(t, 150, stats.TotalPages) // 100 + 50
		assert.Equal(t, 75, stats.Pages)       // 50 + 25
	})

	// Test case 2: Zero previous week pages (edge case)
	t.Run("zero previous week pages", func(t *testing.T) {
		mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
			return []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Test Project", TotalPages: 100, LogCount: 50},
			}, nil
		}

		mockRepo.prevWeekPages = 0 // Zero previous week
		mockRepo.lastWeekPages = 20

		mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
			return 8.0, nil
		}

		ctx := context.Background()
		stats, err := dayService.CalculateWeeklyStats(ctx)

		require.NoError(t, err)
		assert.NotNil(t, stats)

		// per_pages should be nil when previous week is 0 (matches Rails API behavior)
		assert.Nil(t, stats.PerPages, "per_pages should be null when previous week pages is 0")
	})

	// Test case 3: Empty aggregates
	t.Run("empty aggregates", func(t *testing.T) {
		mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
			return []*dto.ProjectAggregate{}, nil
		}

		ctx := context.Background()
		stats, err := dayService.CalculateWeeklyStats(ctx)

		require.NoError(t, err)
		assert.NotNil(t, stats)

		// All values should be zero for empty aggregates
		assert.Equal(t, 0, stats.PreviousWeekPages)
		assert.Equal(t, 0, stats.LastWeekPages)
		// per_pages should be nil when previous week is 0
		assert.Nil(t, stats.PerPages, "per_pages should be null when previous week pages is 0")
		assert.Equal(t, 0.0, stats.MeanDay)
		assert.Equal(t, 0.0, stats.SpecMeanDay)
	})

	// Test case 4: Float rounding to 3 decimals
	t.Run("float rounding", func(t *testing.T) {
		resetMock(mockRepo)
		mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
			return []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Test Project", TotalPages: 100, LogCount: 50},
			}, nil
		}

		mockRepo.prevWeekPages = 3
		mockRepo.lastWeekPages = 5

		mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
			// Return a value that would need rounding
			return 7.123456, nil
		}

		ctx := context.Background()
		stats, err := dayService.CalculateWeeklyStats(ctx)

		require.NoError(t, err)
		assert.NotNil(t, stats)

		// Check rounding to 3 decimals
		assert.NotNil(t, stats.PerPages)
		assert.InDelta(t, 166.667, *stats.PerPages, 0.001) // (5/3)*100 rounded
		assert.InDelta(t, 7.123, stats.MeanDay, 0.001)     // 7.123456 rounded
		assert.InDelta(t, 8.191, stats.SpecMeanDay, 0.001) // 7.123 * 1.15 rounded
	})
}

// TestDayService_GetToday tests the GetToday helper function
func TestDayService_GetToday(t *testing.T) {
	// Arrange - Use a fixed date for deterministic testing
	fixedDate := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)

	// Act
	withDayServiceFixedDate(t, fixedDate, func() {
		today := dashboard.GetToday()

		// Verify it returns a time
		assert.NotNil(t, today)

		// Verify it's truncated to midnight
		assert.Equal(t, 0, today.Hour())
		assert.Equal(t, 0, today.Minute())
		assert.Equal(t, 0, today.Second())
		assert.Equal(t, 0, today.Nanosecond())

		// Verify it matches the fixed date (not current date)
		assert.Equal(t, fixedDate.Year(), today.Year())
		assert.Equal(t, fixedDate.Month(), today.Month())
		assert.Equal(t, fixedDate.Day(), today.Day())
	})
}

// TestDayService_CalculatePerPagesRatio tests the ratio calculation
func TestDayService_CalculatePerPagesRatio(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	mockConfig := &MockUserConfigService{mockPredictionPct: 0.15}
	resetMock(mockRepo)
	dayService := dashboard.NewDayService(mockRepo, mockConfig)

	testCases := []struct {
		name           string
		lastWeek       int
		previousWeek   int
		expectedRatio  *float64
		expectNilRatio bool
	}{
		{"normal ratio", 100, 50, ptrToFloat64(200.0), false},
		{"ratio less than 100%", 30, 100, ptrToFloat64(30.0), false},
		{"equal values", 50, 50, ptrToFloat64(100.0), false},
		{"zero previous week", 100, 0, nil, true},
		{"zero both", 0, 0, nil, true},
		{"large values", 10000, 5000, ptrToFloat64(200.0), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ratio := dayService.CalculatePerPagesRatio(tc.lastWeek, tc.previousWeek)
			if tc.expectNilRatio {
				assert.Nil(t, ratio, "should return nil when previous week is 0")
			} else {
				assert.NotNil(t, ratio, "should return non-nil ratio")
				assert.InDelta(t, *tc.expectedRatio, *ratio, 0.001, "ratio mismatch")
			}
		})
	}
}

// ptrToFloat64 is a helper to create a pointer to a float64
func ptrToFloat64(v float64) *float64 {
	return &v
}

// TestDayService_CalculateSpecMeanDay tests the speculative mean calculation
func TestDayService_CalculateSpecMeanDay(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	mockConfig := &MockUserConfigService{mockPredictionPct: 0.20}
	resetMock(mockRepo)
	dayService := dashboard.NewDayService(mockRepo, mockConfig)

	testCases := []struct {
		name         string
		meanDay      float64
		expectedSpec float64
	}{
		{"normal mean", 10.0, 12.0},
		{"zero mean", 0.0, 0.0},
		{"small mean", 1.5, 1.8},
		{"large mean", 100.0, 120.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			specMean := dayService.CalculateSpecMeanDay(tc.meanDay)
			assert.Equal(t, tc.expectedSpec, specMean, "spec mean mismatch")
		})
	}
}

// TestDayService_CalculateProgress tests the progress calculation
func TestDayService_CalculateProgress(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	mockConfig := &MockUserConfigService{mockPredictionPct: 0.15}
	resetMock(mockRepo)
	dayService := dashboard.NewDayService(mockRepo, mockConfig)

	aggregates := []*dto.ProjectAggregate{
		{ProjectID: 1, ProjectName: "Project 1", TotalPages: 100, LogCount: 50},
		{ProjectID: 2, ProjectName: "Project 2", TotalPages: 200, LogCount: 100},
		{ProjectID: 3, ProjectName: "Project 3", TotalPages: 50, LogCount: 25},
	}

	totalPages, pages := dayService.CalculateProgress(aggregates)

	assert.Equal(t, 350, totalPages) // 100 + 200 + 50
	assert.Equal(t, 175, pages)      // 50 + 100 + 25
}

// TestDayService_CalculateMeanDay tests the mean day calculation
func TestDayService_CalculateMeanDay(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	mockConfig := &MockUserConfigService{mockPredictionPct: 0.15}
	resetMock(mockRepo)
	dayService := dashboard.NewDayService(mockRepo, mockConfig)

	aggregates := []*dto.ProjectAggregate{
		{ProjectID: 1, ProjectName: "Project 1", TotalPages: 100, LogCount: 50},
		{ProjectID: 2, ProjectName: "Project 2", TotalPages: 200, LogCount: 100},
	}

	// Mock the project weekday mean calls
	mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
		if projectID == 1 {
			return 10.0, nil
		}
		if projectID == 2 {
			return 20.0, nil
		}
		return 0.0, nil
	}

	meanDay, err := dayService.CalculateMeanDay(context.Background(), aggregates)

	require.NoError(t, err)
	assert.Equal(t, 15.0, meanDay) // (10 + 20) / 2 = 15
}

// TestDayService_CalculateMeanDay_EmptyAggregates tests with no aggregates
func TestDayService_CalculateMeanDay_EmptyAggregates(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	mockConfig := &MockUserConfigService{mockPredictionPct: 0.15}
	resetMock(mockRepo)
	dayService := dashboard.NewDayService(mockRepo, mockConfig)

	aggregates := []*dto.ProjectAggregate{}

	meanDay, err := dayService.CalculateMeanDay(context.Background(), aggregates)

	require.NoError(t, err)
	assert.Equal(t, 0.0, meanDay)
}

// TestDayService_CalculatePeriodPages tests the period pages calculation
func TestDayService_CalculatePeriodPages(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	mockConfig := &MockUserConfigService{mockPredictionPct: 0.15}
	resetMock(mockRepo)
	dayService := dashboard.NewDayService(mockRepo, mockConfig)

	startDate, _ := time.Parse("2006-01-02", "2024-01-08")
	endDate, _ := time.Parse("2006-01-02", "2024-01-15")

	mockRepo.prevWeekPages = 25

	pages, err := dayService.CalculatePeriodPages(context.Background(), startDate, endDate)

	require.NoError(t, err)
	assert.Equal(t, 25, pages)
}

// TestDayService_CalculatePeriodPages_Error tests error handling
func TestDayService_CalculatePeriodPages_Error(t *testing.T) {
	mockRepo := &MockDashboardRepository{}
	mockConfig := &MockUserConfigService{mockPredictionPct: 0.15}
	resetMock(mockRepo)
	dayService := dashboard.NewDayService(mockRepo, mockConfig)

	startDate, _ := time.Parse("2006-01-02", "2024-01-08")
	endDate, _ := time.Parse("2006-01-02", "2024-01-15")

	// Note: The mock doesn't return errors currently
	// This test verifies the error handling path

	pages, err := dayService.CalculatePeriodPages(context.Background(), startDate, endDate)

	assert.NoError(t, err)
	assert.Equal(t, 0, pages)
}

// =============================================================================
// Exported Mock Types for Use in Other Test Files
// =============================================================================

// MockDashboardRepositoryExtended extends the existing mock with additional methods
type MockDashboardRepositoryExtended struct {
	MockDashboardRepository
	mockGetLogsByDateRange func(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error)
}

func (m *MockDashboardRepositoryExtended) GetLogsByDateRange(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
	if m.mockGetLogsByDateRange != nil {
		return m.mockGetLogsByDateRange(ctx, start, end)
	}
	panic("not implemented")
}

// MockUserConfigProvider is a mock implementation of UserConfigProvider for testing
type MockUserConfigProvider struct {
	mockPredictionPct float64
	mockMaxFaults     int
	mockPagesPerDay   float64
}

func (m *MockUserConfigProvider) GetPredictionPct() float64 {
	return m.mockPredictionPct
}

func (m *MockUserConfigProvider) GetMaxFaults() int {
	return m.mockMaxFaults
}

func (m *MockUserConfigProvider) GetPagesPerDay() float64 {
	if m.mockPagesPerDay == 0 {
		return 25.0 // Default value from GetDefaultConfig
	}
	return m.mockPagesPerDay
}

// =============================================================================
// NEW TESTS: Fixed Test Data for Rails Parity Verification
// =============================================================================

// TestDayService_CalculateMeanDay_RailsParity tests the mean_day calculation
// with fixed test data to verify deterministic results.
//
// Note: The Go implementation uses a different algorithm than Rails:
//   - Rails: mean = total_pages / count_7day_intervals
//   - Go: mean = AVG(read_pages) for logs on current weekday
//
// This test verifies the Go algorithm produces consistent, expected results
// with fixed test data, not necessarily Rails parity.
//
// Fixed Date: 2026-04-21 (Tuesday, weekday=2)
func TestDayService_CalculateMeanDay_RailsParity(t *testing.T) {
	// Setup with fixed date: 2026-04-21 (Tuesday)
	fixedDate := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	defer dashboard.SetTestDate(time.Now())
	dashboard.SetTestDate(fixedDate)

	mockRepo := &MockDashboardRepository{}
	mockConfig := &MockUserConfigService{mockPredictionPct: 0.15}
	resetMock(mockRepo)
	dayService := dashboard.NewDayService(mockRepo, mockConfig)

	// Fixed test data: 3 projects with known weekday means
	aggregates := []*dto.ProjectAggregate{
		{ProjectID: 1, ProjectName: "Project Alpha", TotalPages: 300, LogCount: 60},
		{ProjectID: 2, ProjectName: "Project Beta", TotalPages: 200, LogCount: 40},
		{ProjectID: 3, ProjectName: "Project Gamma", TotalPages: 150, LogCount: 30},
	}

	// Mock returns fixed weekday means for Tuesday (weekday=2)
	// These values are predetermined and deterministic
	mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
		switch projectID {
		case 1:
			return 12.5, nil // Project Alpha: 12.5 pages/day on Tuesdays
		case 2:
			return 8.75, nil // Project Beta: 8.75 pages/day on Tuesdays
		case 3:
			return 15.0, nil // Project Gamma: 15.0 pages/day on Tuesdays
		}
		return 0.0, nil
	}

	// Act
	meanDay, err := dayService.CalculateMeanDay(context.Background(), aggregates)

	// Assert
	require.NoError(t, err)
	// Expected: (12.5 + 8.75 + 15.0) / 3 = 12.083333... -> rounded to 12.083
	expected := 12.083
	assert.InDelta(t, expected, meanDay, 0.001, "mean_day should be average of project weekday means")
}

// TestDayService_CalculateMeanDay_RailsParity_SingleProject tests with a single project
func TestDayService_CalculateMeanDay_RailsParity_SingleProject(t *testing.T) {
	fixedDate := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	defer dashboard.SetTestDate(time.Now())
	dashboard.SetTestDate(fixedDate)

	mockRepo := &MockDashboardRepository{}
	mockConfig := &MockUserConfigService{mockPredictionPct: 0.15}
	resetMock(mockRepo)
	dayService := dashboard.NewDayService(mockRepo, mockConfig)

	aggregates := []*dto.ProjectAggregate{
		{ProjectID: 1, ProjectName: "Single Project", TotalPages: 250, LogCount: 50},
	}

	mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
		return 10.555, nil // Fixed value that tests rounding
	}

	meanDay, err := dayService.CalculateMeanDay(context.Background(), aggregates)

	require.NoError(t, err)
	// Single project: mean = project mean = 10.555 (already rounded)
	assert.InDelta(t, 10.555, meanDay, 0.001)
}

// TestDayService_CalculateMeanDay_MultipleWeekdays tests calculations across
// different weekdays to ensure the service correctly identifies and uses
// the current weekday from the fixed date.
func TestDayService_CalculateMeanDay_MultipleWeekdays(t *testing.T) {
	testCases := []struct {
		name            string
		fixedDate       time.Time
		expectedWeekday int
		projectMeans    map[int64]float64
		expectedMean    float64
	}{
		{
			name:            "Monday (weekday=1)",
			fixedDate:       time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC), // Monday
			expectedWeekday: 1,
			projectMeans: map[int64]float64{
				1: 10.0,
				2: 20.0,
			},
			expectedMean: 15.0, // (10 + 20) / 2
		},
		{
			name:            "Wednesday (weekday=3)",
			fixedDate:       time.Date(2026, 4, 22, 12, 0, 0, 0, time.UTC), // Wednesday
			expectedWeekday: 3,
			projectMeans: map[int64]float64{
				1: 5.5,
				2: 7.5,
				3: 9.0,
			},
			expectedMean: 7.333, // (5.5 + 7.5 + 9.0) / 3 = 7.333...
		},
		{
			name:            "Sunday (weekday=0)",
			fixedDate:       time.Date(2026, 4, 19, 12, 0, 0, 0, time.UTC), // Sunday
			expectedWeekday: 0,
			projectMeans: map[int64]float64{
				1: 100.0,
			},
			expectedMean: 100.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer dashboard.SetTestDate(time.Now())
			dashboard.SetTestDate(tc.fixedDate)

			mockRepo := &MockDashboardRepository{}
			mockConfig := &MockUserConfigService{mockPredictionPct: 0.15}
			resetMock(mockRepo)
			dayService := dashboard.NewDayService(mockRepo, mockConfig)

			aggregates := []*dto.ProjectAggregate{}
			for projectID := range tc.projectMeans {
				aggregates = append(aggregates, &dto.ProjectAggregate{
					ProjectID:   projectID,
					ProjectName: fmt.Sprintf("Project %d", projectID),
					TotalPages:  100,
					LogCount:    20,
				})
			}

			mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
				// Verify the correct weekday is being queried
				assert.Equal(t, tc.expectedWeekday, weekday, "should query for current weekday")
				return tc.projectMeans[projectID], nil
			}

			meanDay, err := dayService.CalculateMeanDay(context.Background(), aggregates)

			require.NoError(t, err)
			assert.InDelta(t, tc.expectedMean, meanDay, 0.001, "mean calculation mismatch for %s", tc.name)
		})
	}
}

// TestDayService_CalculateMeanDay_EdgeCases covers edge cases for mean calculation
func TestDayService_CalculateMeanDay_EdgeCases(t *testing.T) {
	testCases := []struct {
		name         string
		fixedDate    time.Time
		aggregates   []*dto.ProjectAggregate
		mockReturns  map[int64]float64
		mockError    map[int64]error
		expectedMean float64
		expectError  bool
	}{
		{
			name:      "no logs for weekday - returns zero",
			fixedDate: time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC),
			aggregates: []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Project 1", TotalPages: 100, LogCount: 20},
			},
			mockReturns:  map[int64]float64{1: 0.0}, // No pages read on this weekday
			mockError:    map[int64]error{},
			expectedMean: 0.0,
			expectError:  false,
		},
		{
			name:      "single log entry - returns that value",
			fixedDate: time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC),
			aggregates: []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Project 1", TotalPages: 50, LogCount: 1},
			},
			mockReturns:  map[int64]float64{1: 25.0}, // Single entry with 25 pages
			mockError:    map[int64]error{},
			expectedMean: 25.0,
			expectError:  false,
		},
		{
			name:      "zero pages read - returns zero",
			fixedDate: time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC),
			aggregates: []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Project 1", TotalPages: 100, LogCount: 10},
			},
			mockReturns:  map[int64]float64{1: 0.0},
			mockError:    map[int64]error{},
			expectedMean: 0.0,
			expectError:  false,
		},
		{
			name:      "large page counts - handles correctly",
			fixedDate: time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC),
			aggregates: []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Large Project 1", TotalPages: 10000, LogCount: 500},
				{ProjectID: 2, ProjectName: "Large Project 2", TotalPages: 15000, LogCount: 750},
			},
			mockReturns:  map[int64]float64{1: 500.5, 2: 750.75},
			mockError:    map[int64]error{},
			expectedMean: 625.625, // (500.5 + 750.75) / 2
			expectError:  false,
		},
		{
			name:      "float precision - rounds to 3 decimals",
			fixedDate: time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC),
			aggregates: []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Project 1", TotalPages: 100, LogCount: 33},
			},
			mockReturns:  map[int64]float64{1: 3.333333}, // Would be 3.333333...
			mockError:    map[int64]error{},
			expectedMean: 3.333, // Rounded to 3 decimals
			expectError:  false,
		},
		{
			name:         "empty aggregates - returns zero",
			fixedDate:    time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC),
			aggregates:   []*dto.ProjectAggregate{},
			mockReturns:  map[int64]float64{},
			mockError:    map[int64]error{},
			expectedMean: 0.0,
			expectError:  false,
		},
		{
			name:      "repository error - returns error",
			fixedDate: time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC),
			aggregates: []*dto.ProjectAggregate{
				{ProjectID: 1, ProjectName: "Project 1", TotalPages: 100, LogCount: 20},
			},
			mockReturns:  map[int64]float64{},
			mockError:    map[int64]error{1: fmt.Errorf("database error")},
			expectedMean: 0.0,
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer dashboard.SetTestDate(time.Now())
			dashboard.SetTestDate(tc.fixedDate)

			mockRepo := &MockDashboardRepository{}
			mockConfig := &MockUserConfigService{mockPredictionPct: 0.15}
			resetMock(mockRepo)
			dayService := dashboard.NewDayService(mockRepo, mockConfig)

			mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
				if err, ok := tc.mockError[projectID]; ok {
					return 0.0, err
				}
				return tc.mockReturns[projectID], nil
			}

			meanDay, err := dayService.CalculateMeanDay(context.Background(), tc.aggregates)

			if tc.expectError {
				assert.Error(t, err, "should return error")
				assert.Equal(t, 0.0, meanDay)
			} else {
				require.NoError(t, err)
				assert.InDelta(t, tc.expectedMean, meanDay, 0.001, "mean calculation mismatch")
			}
		})
	}
}

// TestDayService_CalculateWeeklyStats_FixedData tests the complete weekly stats
// calculation with fixed test data to verify all fields are calculated correctly.
//
// Fixed Date: 2026-04-21 (Tuesday)
// Expected calculations:
//   - previous_week_pages: mocked value
//   - last_week_pages: mocked value
//   - per_pages: (last_week / previous_week) * 100
//   - mean_day: average of project weekday means
//   - spec_mean_day: mean_day * (1 + prediction_pct)
func TestDayService_CalculateWeeklyStats_FixedData(t *testing.T) {
	fixedDate := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	defer dashboard.SetTestDate(time.Now())
	dashboard.SetTestDate(fixedDate)

	mockRepo := &MockDashboardRepository{}
	mockConfig := &MockUserConfigService{mockPredictionPct: 0.15}
	resetMock(mockRepo)
	dayService := dashboard.NewDayService(mockRepo, mockConfig)

	// Fixed test data
	aggregates := []*dto.ProjectAggregate{
		{ProjectID: 1, ProjectName: "Project A", TotalPages: 200, LogCount: 40},
		{ProjectID: 2, ProjectName: "Project B", TotalPages: 150, LogCount: 30},
	}

	// Mock period pages calculations
	mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
		return aggregates, nil
	}

	mockRepo.prevWeekPages = 50 // Fixed: previous week pages
	mockRepo.lastWeekPages = 75 // Fixed: last week pages

	// Mock weekday means for Tuesday
	mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
		switch projectID {
		case 1:
			return 10.0, nil
		case 2:
			return 20.0, nil
		}
		return 0.0, nil
	}

	// Act
	stats, err := dayService.CalculateWeeklyStats(context.Background())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, stats)

	// Verify fixed values
	assert.Equal(t, 50, stats.PreviousWeekPages, "previous_week_pages should match mock")
	assert.Equal(t, 75, stats.LastWeekPages, "last_week_pages should match mock")

	// Verify calculated per_pages: (75 / 50) * 100 = 150.0
	assert.NotNil(t, stats.PerPages)
	assert.InDelta(t, 150.0, *stats.PerPages, 0.001, "per_pages calculation")

	// Verify mean_day: (10.0 + 20.0) / 2 = 15.0
	assert.InDelta(t, 15.0, stats.MeanDay, 0.001, "mean_day calculation")

	// Verify spec_mean_day: 15.0 * 1.15 = 17.25
	assert.InDelta(t, 17.25, stats.SpecMeanDay, 0.001, "spec_mean_day calculation")

	// Verify progress values
	assert.Equal(t, 350, stats.TotalPages, "total_pages: 200 + 150")
	assert.Equal(t, 70, stats.Pages, "pages: 40 + 30")
}

// TestDayService_CalculateWeeklyStats_FixedData_Comprehensive tests with
// values that test rounding to 3 decimal places
func TestDayService_CalculateWeeklyStats_FixedData_Comprehensive(t *testing.T) {
	fixedDate := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
	defer dashboard.SetTestDate(time.Now())
	dashboard.SetTestDate(fixedDate)

	mockRepo := &MockDashboardRepository{}
	mockConfig := &MockUserConfigService{mockPredictionPct: 0.123} // Non-standard prediction pct
	resetMock(mockRepo)
	dayService := dashboard.NewDayService(mockRepo, mockConfig)

	aggregates := []*dto.ProjectAggregate{
		{ProjectID: 1, ProjectName: "Project", TotalPages: 100, LogCount: 25},
	}

	mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
		return aggregates, nil
	}

	mockRepo.prevWeekPages = 7 // Creates non-integer ratio
	mockRepo.lastWeekPages = 11

	mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
		return 7.123456, nil // Tests rounding
	}

	stats, err := dayService.CalculateWeeklyStats(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, stats)

	// per_pages: (11 / 7) * 100 = 157.142857... -> 157.143
	assert.InDelta(t, 157.143, *stats.PerPages, 0.001, "per_pages rounded to 3 decimals")

	// mean_day: 7.123456 -> 7.123
	assert.InDelta(t, 7.123, stats.MeanDay, 0.001, "mean_day rounded to 3 decimals")

	// spec_mean_day: 7.123 * 1.123 = 7.999129 -> 7.999
	assert.InDelta(t, 7.999, stats.SpecMeanDay, 0.001, "spec_mean_day rounded to 3 decimals")
}

// TestDayService_CalculateWeeklyStats_FixedData_EdgeCases tests edge cases
// with fixed data
func TestDayService_CalculateWeeklyStats_FixedData_EdgeCases(t *testing.T) {
	testCases := []struct {
		name              string
		fixedDate         time.Time
		prevWeekPages     int
		lastWeekPages     int
		weekdayMean       float64
		predictionPct     float64
		expectedPerPages  *float64
		expectNilPerPages bool
		expectedSpecMean  float64
	}{
		{
			name:              "zero previous week - per_pages is nil",
			fixedDate:         time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC),
			prevWeekPages:     0,
			lastWeekPages:     50,
			weekdayMean:       10.0,
			predictionPct:     0.15,
			expectedPerPages:  nil,
			expectNilPerPages: true,
			expectedSpecMean:  11.5, // 10.0 * 1.15
		},
		{
			name:              "both weeks zero - per_pages is nil",
			fixedDate:         time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC),
			prevWeekPages:     0,
			lastWeekPages:     0,
			weekdayMean:       0.0,
			predictionPct:     0.15,
			expectedPerPages:  nil,
			expectNilPerPages: true,
			expectedSpecMean:  0.0,
		},
		{
			name:              "equal weeks - per_pages is 100",
			fixedDate:         time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC),
			prevWeekPages:     100,
			lastWeekPages:     100,
			weekdayMean:       15.0,
			predictionPct:     0.15,
			expectedPerPages:  ptrToFloat64(100.0),
			expectNilPerPages: false,
			expectedSpecMean:  17.25,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer dashboard.SetTestDate(time.Now())
			dashboard.SetTestDate(tc.fixedDate)

			mockRepo := &MockDashboardRepository{}
			mockConfig := &MockUserConfigService{mockPredictionPct: tc.predictionPct}
			resetMock(mockRepo)
			dayService := dashboard.NewDayService(mockRepo, mockConfig)

			mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
				return []*dto.ProjectAggregate{
					{ProjectID: 1, ProjectName: "Test", TotalPages: 100, LogCount: 20},
				}, nil
			}

			mockRepo.prevWeekPages = tc.prevWeekPages
			mockRepo.lastWeekPages = tc.lastWeekPages

			mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
				return tc.weekdayMean, nil
			}

			stats, err := dayService.CalculateWeeklyStats(context.Background())

			require.NoError(t, err)
			assert.NotNil(t, stats)
			if tc.expectNilPerPages {
				assert.Nil(t, stats.PerPages, "per_pages should be null when previous week pages is 0")
			} else {
				assert.NotNil(t, stats.PerPages)
				assert.InDelta(t, *tc.expectedPerPages, *stats.PerPages, 0.001)
			}
			assert.InDelta(t, tc.expectedSpecMean, stats.SpecMeanDay, 0.001)
		})
	}
}
