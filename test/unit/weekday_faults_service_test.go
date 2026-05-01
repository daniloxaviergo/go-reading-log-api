package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
	"go-reading-log-api-next/internal/service/dashboard"
)

// Helper function to set a fixed test date and restore it after the test
func withWeekdayFaultsServiceFixedDate(t *testing.T, fixedDate time.Time, fn func()) {
	defer dashboard.SetTestDate(time.Now())
	dashboard.SetTestDate(fixedDate)
	fn()
}

// MockDashboardRepositoryWeekdayFaults is a mock implementation of DashboardRepository for testing weekday faults
type MockDashboardRepositoryWeekdayFaults struct {
	mockGetWeekdayFaults func(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error)
}

func (m *MockDashboardRepositoryWeekdayFaults) GetDailyStats(ctx context.Context, date time.Time) (*dto.DailyStats, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryWeekdayFaults) GetProjectAggregates(ctx context.Context) ([]*dto.ProjectAggregate, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryWeekdayFaults) GetFaultsByDateRange(ctx context.Context, start, end time.Time) (*dto.FaultStats, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryWeekdayFaults) GetWeekdayFaults(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error) {
	if m.mockGetWeekdayFaults != nil {
		return m.mockGetWeekdayFaults(ctx, start, end)
	}
	panic("not implemented")
}

func (m *MockDashboardRepositoryWeekdayFaults) GetLogsByDateRange(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryWeekdayFaults) GetProjectWeekdayMean(ctx context.Context, projectID int64, weekday int) (float64, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryWeekdayFaults) CalculatePeriodPages(ctx context.Context, start, end time.Time) (int, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryWeekdayFaults) GetPool() repository.PoolInterface {
	return nil
}

func (m *MockDashboardRepositoryWeekdayFaults) GetProjectsWithLogs(ctx context.Context) ([]*dto.ProjectAggregateResponse, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryWeekdayFaults) GetProjectLogs(ctx context.Context, projectID int64, limit int) ([]*dto.LogEntry, error) {
	panic("not implemented")
}

func (m *MockDashboardRepositoryWeekdayFaults) GetMaxByWeekday(ctx context.Context, date time.Time) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepositoryWeekdayFaults) GetOverallMean(ctx context.Context, date time.Time) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepositoryWeekdayFaults) GetPreviousPeriodMean(ctx context.Context, date time.Time) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepositoryWeekdayFaults) GetPreviousPeriodSpecMean(ctx context.Context, date time.Time) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepositoryWeekdayFaults) GetMeanByWeekday(ctx context.Context, weekday int) (*float64, error) {
	return nil, nil
}

func (m *MockDashboardRepositoryWeekdayFaults) GetRunningProjectsWithLogs(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
	return nil, nil
}

// MockUserConfigProviderWeekday is a mock implementation of UserConfigProvider for testing
type MockUserConfigProviderWeekday struct {
	mockPredictionPct float64
	mockMaxFaults     int
	mockPagesPerDay   float64
}

func (m *MockUserConfigProviderWeekday) GetPredictionPct() float64 {
	return m.mockPredictionPct
}

func (m *MockUserConfigProviderWeekday) GetMaxFaults() int {
	return m.mockMaxFaults
}

func (m *MockUserConfigProviderWeekday) GetPagesPerDay() float64 {
	return m.mockPagesPerDay
}

// TestWeekdayFaultsService_GetWeekdayFaults tests the main GetWeekdayFaults method
func TestWeekdayFaultsService_GetWeekdayFaults(t *testing.T) {
	// Setup mocks
	mockRepo := &MockDashboardRepositoryWeekdayFaults{}
	mockConfig := &MockUserConfigProviderWeekday{
		mockPredictionPct: 0.15,
		mockMaxFaults:     10,
		mockPagesPerDay:   25.0,
	}

	// Create WeekdayFaultsService with mocks
	weekdayService := dashboard.NewWeekdayFaultsService(mockRepo, mockConfig)

	// Test case 1: Normal data with all 7 days present
	t.Run("normal data with all days", func(t *testing.T) {
		mockRepo.mockGetWeekdayFaults = func(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error) {
			// Verify date range is approximately 6 months (180-186 days)
			daysDiff := int(end.Sub(start).Hours() / 24)
			assert.GreaterOrEqual(t, daysDiff, 175)
			assert.LessOrEqual(t, daysDiff, 190)

			// Verify end date is approximately today
			now := time.Now()
			assert.Equal(t, now.Year(), end.Year())
			assert.Equal(t, now.Month(), end.Month())
			assert.Equal(t, now.Day(), end.Day())

			return dto.NewWeekdayFaults(map[int]int{
				0: 5, // Sunday
				1: 8, // Monday
				2: 3, // Tuesday
				3: 7, // Wednesday
				4: 2, // Thursday
				5: 9, // Friday
				6: 4, // Saturday
			}), nil
		}

		ctx := context.Background()
		result, err := weekdayService.GetWeekdayFaults(ctx)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 5, result.Faults[0]) // Sunday
		assert.Equal(t, 4, result.Faults[6]) // Saturday
	})

	// Test case 2: Missing some days (repository fills in zeros for missing days)
	t.Run("missing days in data", func(t *testing.T) {
		mockRepo.mockGetWeekdayFaults = func(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error) {
			// Return data with some days missing - repository should fill in 0s
			// Note: The mock returns exactly what we specify; the real repository
			// would fill in missing days with 0. For this test, we simulate that.
			faults := map[int]int{
				0: 5,
				1: 0, // Filled in by repository
				2: 3,
				3: 0, // Filled in by repository
				4: 2,
				5: 0, // Filled in by repository
				6: 4,
			}
			return dto.NewWeekdayFaults(faults), nil
		}

		ctx := context.Background()
		result, err := weekdayService.GetWeekdayFaults(ctx)

		require.NoError(t, err)
		assert.NotNil(t, result)
		// Repository ensures all 7 days are present with default value of 0
		expectedDays := []int{0, 1, 2, 3, 4, 5, 6}
		for _, i := range expectedDays {
			assert.Contains(t, result.Faults, i, "Day %d should be present", i)
			assert.GreaterOrEqual(t, result.Faults[i], 0, "Day %d count should be non-negative", i)
		}
	})

	// Test case 3: Empty data (no faults in range)
	t.Run("empty data", func(t *testing.T) {
		mockRepo.mockGetWeekdayFaults = func(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error) {
			return dto.NewWeekdayFaults(map[int]int{
				0: 0, 1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0,
			}), nil
		}

		ctx := context.Background()
		result, err := weekdayService.GetWeekdayFaults(ctx)

		require.NoError(t, err)
		assert.NotNil(t, result)
		for i := 0; i < 7; i++ {
			assert.Equal(t, 0, result.Faults[i])
		}
	})

	// Test case 4: Repository error handling
	t.Run("repository error", func(t *testing.T) {
		mockRepo.mockGetWeekdayFaults = func(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error) {
			return nil, assert.AnError
		}

		ctx := context.Background()
		result, err := weekdayService.GetWeekdayFaults(ctx)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get weekday faults")
	})
}

// TestWeekdayFaultsService_CreateRadarChart tests the CreateRadarChart method
func TestWeekdayFaultsService_CreateRadarChart(t *testing.T) {
	mockRepo := &MockDashboardRepositoryWeekdayFaults{}
	mockConfig := &MockUserConfigProviderWeekday{
		mockPredictionPct: 0.15,
		mockMaxFaults:     10,
		mockPagesPerDay:   25.0,
	}

	weekdayService := dashboard.NewWeekdayFaultsService(mockRepo, mockConfig)

	testCases := []struct {
		name      string
		faults    map[int]int
		expectLen int
	}{
		{
			name: "all 7 days present",
			faults: map[int]int{
				0: 5, 1: 8, 2: 3, 3: 7, 4: 2, 5: 9, 6: 4,
			},
			expectLen: 7,
		},
		{
			name: "missing some days",
			faults: map[int]int{
				0: 5, 2: 3, 4: 2, 6: 4,
			},
			expectLen: 7, // Should still have 7 slots with 0 for missing
		},
		{
			name:      "empty map",
			faults:    map[int]int{},
			expectLen: 7,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			radarChart := weekdayService.CreateRadarChart(tc.faults)

			assert.NotNil(t, radarChart)
			assert.Equal(t, "Faults by Weekday", radarChart.Title)
			assert.NotNil(t, radarChart.Series)
			assert.Len(t, radarChart.Series, 1)

			series := radarChart.Series[0]
			assert.Equal(t, "Faults", series.Name)
			assert.Equal(t, "radar", series.Type)
			assert.Len(t, series.Data, tc.expectLen)

			// Verify data contains only float64 values (from int conversion)
			for _, val := range series.Data {
				_, ok := val.(float64)
				assert.True(t, ok, "Expected float64 value in radar chart data")
			}
		})
	}
}

// TestWeekdayFaultsService_ValidateOutput tests the ValidateOutput method
func TestWeekdayFaultsService_ValidateOutput(t *testing.T) {
	mockRepo := &MockDashboardRepositoryWeekdayFaults{}
	mockConfig := &MockUserConfigProviderWeekday{
		mockPredictionPct: 0.15,
		mockMaxFaults:     10,
		mockPagesPerDay:   25.0,
	}

	weekdayService := dashboard.NewWeekdayFaultsService(mockRepo, mockConfig)

	testCases := []struct {
		name   string
		faults map[int]int
		valid  bool
	}{
		{
			name: "valid data all days present",
			faults: map[int]int{
				0: 5, 1: 8, 2: 3, 3: 7, 4: 2, 5: 9, 6: 4,
			},
			valid: true,
		},
		{
			name: "missing one day",
			faults: map[int]int{
				0: 5, 1: 8, 2: 3, 3: 7, 4: 2, 5: 9,
				// Missing day 6 (Saturday)
			},
			valid: false,
		},
		{
			name: "negative count",
			faults: map[int]int{
				0: 5, 1: -2, 2: 3, 3: 7, 4: 2, 5: 9, 6: 4,
			},
			valid: false,
		},
		{
			name: "zero counts (valid)",
			faults: map[int]int{
				0: 0, 1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0,
			},
			valid: true,
		},
		{
			name:   "empty map",
			faults: map[int]int{},
			valid:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := weekdayService.ValidateOutput(tc.faults)

			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestWeekdayFaultsService_GetDateRangeLast6Months tests the date range calculation
func TestWeekdayFaultsService_GetDateRangeLast6Months(t *testing.T) {
	start, end := dashboard.GetDateRangeLast6Months()

	// Verify end date is today (truncated to midnight)
	now := time.Now()
	assert.Equal(t, now.Year(), end.Year())
	assert.Equal(t, now.Month(), end.Month())
	assert.Equal(t, now.Day(), end.Day())
	assert.Equal(t, 0, end.Hour())
	assert.Equal(t, 0, end.Minute())

	// Verify start date is 6 months ago
	expectedStart := end.AddDate(0, -6, 0)
	assert.Equal(t, expectedStart.Year(), start.Year())
	assert.Equal(t, expectedStart.Month(), start.Month())
	assert.Equal(t, expectedStart.Day(), start.Day())

	// Verify the range is approximately 6 months (180-186 days)
	daysDiff := int(end.Sub(start).Hours() / 24)
	assert.GreaterOrEqual(t, daysDiff, 180)
	assert.LessOrEqual(t, daysDiff, 186)
}

// TestWeekdayFaultsService_GetToday tests the GetToday helper function
func TestWeekdayFaultsService_GetToday(t *testing.T) {
	// Arrange - Use a fixed date for deterministic testing
	fixedDate := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)

	// Act
	withWeekdayFaultsServiceFixedDate(t, fixedDate, func() {
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

// TestWeekdayFaultsService_Integration tests the complete flow
func TestWeekdayFaultsService_Integration(t *testing.T) {
	// This test verifies the integration between all service methods
	mockRepo := &MockDashboardRepositoryWeekdayFaults{}
	mockConfig := &MockUserConfigProviderWeekday{
		mockPredictionPct: 0.15,
		mockMaxFaults:     10,
		mockPagesPerDay:   25.0,
	}

	weekdayService := dashboard.NewWeekdayFaultsService(mockRepo, mockConfig)

	// Setup mock to return test data
	mockRepo.mockGetWeekdayFaults = func(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error) {
		return dto.NewWeekdayFaults(map[int]int{
			0: 5, 1: 8, 2: 3, 3: 7, 4: 2, 5: 9, 6: 4,
		}), nil
	}

	ctx := context.Background()

	// Get weekday faults
	withWeekdayFaultsServiceFixedDate(t, time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC), func() {
		faults, err := weekdayService.GetWeekdayFaults(ctx)
		require.NoError(t, err)

		// Validate output
		err = weekdayService.ValidateOutput(faults.Faults)
		require.NoError(t, err)

		// Create radar chart
		radarChart := weekdayService.CreateRadarChart(faults.Faults)

		// Verify chart configuration
		assert.NotNil(t, radarChart)
		assert.Equal(t, "Faults by Weekday", radarChart.Title)
		assert.Len(t, radarChart.Series, 1)
		assert.Equal(t, "radar", radarChart.Series[0].Type)
		assert.Len(t, radarChart.Series[0].Data, 7)

		// Verify all values are non-negative floats
		for _, val := range radarChart.Series[0].Data {
			floatVal, ok := val.(float64)
			assert.True(t, ok)
			assert.GreaterOrEqual(t, floatVal, 0.0)
		}
	})
}
