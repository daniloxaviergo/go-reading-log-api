package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/service/dashboard"
)

// Helper function to set a fixed test date and restore it after the test
func withSpeculateServiceFixedDate(t *testing.T, fixedDate time.Time, fn func()) {
	defer dashboard.SetTestDate(time.Now())
	dashboard.SetTestDate(fixedDate)
	fn()
}

// =============================================================================
// TestSpeculateService_CalculateSpeculativeMean - Mean calculation tests
// =============================================================================

// TestSpeculateService_NewSpeculateService tests service initialization
func TestSpeculateService_NewSpeculateService(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepositoryExtended{}
	mockConfig := &MockUserConfigProvider{}

	// Act
	service := dashboard.NewSpeculateService(mockRepo, mockConfig)

	// Assert
	assert.NotNil(t, service)
}

// TestSpeculateService_CalculateSpeculativeMean_Normal tests normal mean calculation
func TestSpeculateService_CalculateSpeculativeMean_Normal(t *testing.T) {
	// Arrange
	predictionPct := 0.15

	// Act - Call CalculateSpeculativeMean directly (it's a static method)
	actualMean := 25.0
	specMean := dashboard.CalculateSpeculativeMean(actualMean, predictionPct)

	// Assert
	// spec_mean = actual_mean * (1 + prediction_pct)
	// spec_mean = 25.0 * (1 + 0.15) = 28.75
	expectedSpecMean := 28.75
	assert.Equal(t, expectedSpecMean, specMean)
}

// TestSpeculateService_CalculateSpeculativeMean_ZeroMean tests zero mean edge case
func TestSpeculateService_CalculateSpeculativeMean_ZeroMean(t *testing.T) {
	// Arrange
	predictionPct := 0.15

	// Act
	specMean := dashboard.CalculateSpeculativeMean(0.0, predictionPct)

	// Assert
	// Zero mean should return 0.0 (not NaN)
	assert.Equal(t, 0.0, specMean)
}

// TestSpeculateService_CalculateSpeculativeMean_NegativeMean tests negative mean edge case
func TestSpeculateService_CalculateSpeculativeMean_NegativeMean(t *testing.T) {
	// Arrange
	predictionPct := 0.15

	// Act
	specMean := dashboard.CalculateSpeculativeMean(-10.0, predictionPct)

	// Assert
	// Negative mean should return 0.0 (not NaN)
	assert.Equal(t, 0.0, specMean)
}

// TestSpeculateService_CalculateSpeculativeMean_HighPercentage tests high prediction percentage
func TestSpeculateService_CalculateSpeculativeMean_HighPercentage(t *testing.T) {
	// Arrange
	actualMean := 20.0
	highPct := 0.5 // 50% prediction

	// Act
	specMean := dashboard.CalculateSpeculativeMean(actualMean, highPct)

	// Assert
	// spec_mean = 20.0 * (1 + 0.5) = 30.0
	expectedSpecMean := 30.0
	assert.Equal(t, expectedSpecMean, specMean)
}

// TestSpeculateService_CalculateSpeculativeMean_Rounding tests decimal rounding
func TestSpeculateService_CalculateSpeculativeMean_Rounding(t *testing.T) {
	// Arrange
	actualMean := 25.123456
	predictionPct := 0.15

	// Act
	specMean := dashboard.CalculateSpeculativeMean(actualMean, predictionPct)

	// Assert
	// spec_mean = 25.123456 * 1.15 = 28.8929744
	// Should round to 3 decimal places: 28.893 (banker's rounding)
	expectedSpecMean := 28.892 // Go uses banker's rounding, so 28.8925 rounds to 28.892
	assert.Equal(t, expectedSpecMean, specMean)
}

// =============================================================================
// TestSpeculateService_GetSpeculativeMean - Integration-style mean calculation
// =============================================================================

// TestSpeculateService_GetSpeculativeMean tests the full speculative mean calculation
func TestSpeculateService_GetSpeculativeMean(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepositoryExtended{}
	mockConfig := &MockUserConfigProvider{mockPredictionPct: 0.15}

	// Create service
	service := dashboard.NewSpeculateService(mockRepo, mockConfig)

	// Setup test data - 3 projects with different weekday means
	_ = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // Sunday (weekday 0)

	mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
		return []*dto.ProjectAggregate{
			dto.NewProjectAggregate(1, "Project A", 100, 5),
			dto.NewProjectAggregate(2, "Project B", 200, 3),
			dto.NewProjectAggregate(3, "Project C", 150, 4),
		}, nil
	}

	mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
		// Return different means for each project
		switch projectID {
		case 1:
			return 25.0, nil
		case 2:
			return 30.0, nil
		case 3:
			return 20.0, nil
		default:
			return 0.0, nil
		}
	}

	// Act
	specMean, err := service.GetSpeculativeMean(context.Background())

	// Assert
	require.NoError(t, err)

	// Expected calculation:
	// actual_mean = (25.0 + 30.0 + 20.0) / 3 = 25.0
	// spec_mean = 25.0 * (1 + 0.15) = 28.75
	expectedSpecMean := 28.75
	assert.Equal(t, expectedSpecMean, specMean)
}

// TestSpeculateService_GetSpeculativeMean_EmptyProjects tests with no projects
func TestSpeculateService_GetSpeculativeMean_EmptyProjects(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepositoryExtended{}
	mockConfig := &MockUserConfigProvider{mockPredictionPct: 0.15}

	mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
		return []*dto.ProjectAggregate{}, nil
	}

	service := dashboard.NewSpeculateService(mockRepo, mockConfig)

	// Act
	specMean, err := service.GetSpeculativeMean(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 0.0, specMean)
}

// TestSpeculateService_GetSpeculativeMean_DifferentWeekday tests different weekday
func TestSpeculateService_GetSpeculativeMean_DifferentWeekday(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepositoryExtended{}
	mockConfig := &MockUserConfigProvider{mockPredictionPct: 0.15}

	mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
		return []*dto.ProjectAggregate{
			dto.NewProjectAggregate(1, "Project A", 100, 5),
		}, nil
	}

	mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
		// Return mean for Monday (weekday 1)
		return 35.0, nil
	}

	service := dashboard.NewSpeculateService(mockRepo, mockConfig)

	// Act
	specMean, err := service.GetSpeculativeMean(context.Background())

	// Assert
	require.NoError(t, err)
	// spec_mean = 35.0 * 1.15 = 40.25
	expectedSpecMean := 40.25
	assert.Equal(t, expectedSpecMean, specMean)
}

// =============================================================================
// TestSpeculateService_GenerateChartData - Chart data generation tests
// =============================================================================

// TestSpeculateService_GenerateChartData_Last15Days tests 15-day data generation
func TestSpeculateService_GenerateChartData_Last15Days(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepositoryExtended{}
	mockConfig := &MockUserConfigProvider{mockPredictionPct: 0.15}

	service := dashboard.NewSpeculateService(mockRepo, mockConfig)

	// Use a fixed date for deterministic testing
	testDate := time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC) // Tuesday

	mockRepo.mockGetLogsByDateRange = func(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
		return []*dto.LogEntry{
			// Tuesday (index 14) - 2 entries
			dto.NewLogEntry(1, testDate.Format(time.RFC3339), 0, 25, nil, nil),
			dto.NewLogEntry(2, testDate.Format(time.RFC3339), 25, 50, nil, nil),
			// Monday (index 13) - 1 entry
			dto.NewLogEntry(3, testDate.AddDate(0, 0, -1).Format(time.RFC3339), 0, 30, nil, nil),
			// Sunday (index 12) - 2 entries
			dto.NewLogEntry(4, testDate.AddDate(0, 0, -2).Format(time.RFC3339), 10, 40, nil, nil),
			dto.NewLogEntry(5, testDate.AddDate(0, 0, -2).Format(time.RFC3339), 40, 65, nil, nil),
		}, nil
	}

	// Act
	withSpeculateServiceFixedDate(t, testDate, func() {
		dataMap, err := service.GenerateChartData(context.Background())

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, dataMap)

		// Verify we have data for the indices we created logs for
		assert.Contains(t, dataMap, 14) // Tuesday (today)
		assert.Contains(t, dataMap, 13) // Monday
		assert.Contains(t, dataMap, 12) // Sunday

		// Verify page counts are correct
		// Tuesday: 25 + 25 = 50 pages
		assert.Equal(t, 50, dataMap[14])
		// Monday: 30 pages
		assert.Equal(t, 30, dataMap[13])
		// Sunday: 30 + 25 = 55 pages
		assert.Equal(t, 55, dataMap[12])
	})
}

// TestSpeculateService_GenerateChartData_MissingDays tests zero-fill for missing days
func TestSpeculateService_GenerateChartData_MissingDays(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepositoryExtended{}
	mockConfig := &MockUserConfigProvider{mockPredictionPct: 0.15}

	service := dashboard.NewSpeculateService(mockRepo, mockConfig)

	// Create logs only for specific days (not all 15 days)
	// Use a fixed date for proper testing
	testDate := time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC) // Tuesday

	mockRepo.mockGetLogsByDateRange = func(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
		return []*dto.LogEntry{
			// Only Tuesday has data (at index 14)
			dto.NewLogEntry(1, testDate.Format(time.RFC3339), 0, 25, nil, nil),
		}, nil
	}

	// Act
	withSpeculateServiceFixedDate(t, testDate, func() {
		dataMap, err := service.GenerateChartData(context.Background())

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, dataMap)

		// Should have Tuesday data at index 14 (today)
		assert.Contains(t, dataMap, 14)
		assert.Equal(t, 25, dataMap[14])

		// Other days should not be in map (will be zero-filled in GetLast15DaysData)
		// The map only contains days with actual data; zero-filling happens in consumer
	})
}

// TestSpeculateService_GenerateChartData_EmptyData tests empty log set
func TestSpeculateService_GenerateChartData_EmptyData(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepositoryExtended{}
	mockConfig := &MockUserConfigProvider{mockPredictionPct: 0.15}

	service := dashboard.NewSpeculateService(mockRepo, mockConfig)

	mockRepo.mockGetLogsByDateRange = func(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
		return []*dto.LogEntry{}, nil
	}

	// Act
	dataMap, err := service.GenerateChartData(context.Background())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, dataMap)
	assert.Len(t, dataMap, 0) // Empty map for no logs
}

// TestSpeculateService_GenerateChartData_InvalidTimestamp tests handling of invalid timestamps
func TestSpeculateService_GenerateChartData_InvalidTimestamp(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepositoryExtended{}
	mockConfig := &MockUserConfigProvider{mockPredictionPct: 0.15}

	service := dashboard.NewSpeculateService(mockRepo, mockConfig)

	mockRepo.mockGetLogsByDateRange = func(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
		return []*dto.LogEntry{
			dto.NewLogEntry(1, "invalid-timestamp", 0, 25, nil, nil),
			dto.NewLogEntry(2, "", 0, 30, nil, nil),
		}, nil
	}

	// Act
	dataMap, err := service.GenerateChartData(context.Background())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, dataMap)
	// Invalid timestamps should be skipped, resulting in empty map
	assert.Len(t, dataMap, 0)
}

// =============================================================================
// TestSpeculateService_GetLast15DaysData - Full 15-day data retrieval
// =============================================================================

// TestSpeculateService_GetLast15DaysData tests full 15-day data with zero-fill
func TestSpeculateService_GetLast15DaysData(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepositoryExtended{}
	mockConfig := &MockUserConfigProvider{mockPredictionPct: 0.15}

	service := dashboard.NewSpeculateService(mockRepo, mockConfig)

	// Create logs for specific days - Use a fixed date for proper testing
	// April 21, 2026 is a Tuesday (weekday 2)
	testDate := time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC) // Tuesday

	mockRepo.mockGetLogsByDateRange = func(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
		return []*dto.LogEntry{
			// Tuesday (index 14, today)
			dto.NewLogEntry(1, testDate.Format(time.RFC3339), 0, 25, nil, nil),
			// Monday (index 13)
			dto.NewLogEntry(2, testDate.AddDate(0, 0, -1).Format(time.RFC3339), 0, 30, nil, nil),
			// Sunday (index 12)
			dto.NewLogEntry(3, testDate.AddDate(0, 0, -2).Format(time.RFC3339), 10, 40, nil, nil),
		}, nil
	}

	// Act
	withSpeculateServiceFixedDate(t, testDate, func() {
		data, err := service.GetLast15DaysData(context.Background())

		// Assert
		require.NoError(t, err)
		assert.Len(t, data, 15) // Must have exactly 15 days

		// Verify data for days with logs
		assert.Equal(t, 25, data[14]) // Tuesday (today)
		assert.Equal(t, 30, data[13]) // Monday
		assert.Equal(t, 30, data[12]) // Sunday (30 pages)

		// Verify zero-fill for days without logs
		// Friday through Thursday (indices 5-11) should be zero
		assert.Equal(t, 0, data[11]) // Thursday
		assert.Equal(t, 0, data[5])  // Friday
	})
}

// TestSpeculateService_GetLast15DaysData_NoData tests with no logs
func TestSpeculateService_GetLast15DaysData_NoData(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepositoryExtended{}
	mockConfig := &MockUserConfigProvider{mockPredictionPct: 0.15}

	service := dashboard.NewSpeculateService(mockRepo, mockConfig)

	mockRepo.mockGetLogsByDateRange = func(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
		return []*dto.LogEntry{}, nil
	}

	// Act
	data, err := service.GetLast15DaysData(context.Background())

	// Assert
	require.NoError(t, err)
	assert.Len(t, data, 15)

	// All days should be zero-filled
	for i, pages := range data {
		assert.Equal(t, 0, pages, "Day %d should be zero-filled", i)
	}
}

// =============================================================================
// TestSpeculateService_GenerateChartConfig - Chart configuration tests
// =============================================================================

// TestSpeculateService_GenerateChartConfig tests full chart generation
func TestSpeculateService_GenerateChartConfig(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepositoryExtended{}
	mockConfig := &MockUserConfigProvider{mockPredictionPct: 0.15}

	service := dashboard.NewSpeculateService(mockRepo, mockConfig)

	// Setup test data - Use a fixed date matching Tuesday for proper testing
	// April 21, 2026 is a Tuesday (weekday 2)
	testDate := time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC) // Tuesday

	mockRepo.mockGetLogsByDateRange = func(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
		return []*dto.LogEntry{
			dto.NewLogEntry(1, testDate.Format(time.RFC3339), 0, 25, nil, nil),
		}, nil
	}

	mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
		return []*dto.ProjectAggregate{
			dto.NewProjectAggregate(1, "Project A", 100, 5),
		}, nil
	}

	mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
		return 25.0, nil
	}

	// Act
	withSpeculateServiceFixedDate(t, testDate, func() {
		chart, err := service.GenerateChartConfig(context.Background())

		// Assert
		require.NoError(t, err)
		assert.NotNil(t, chart)
		assert.Equal(t, "Speculated vs Actual", chart.Title)
		assert.NotNil(t, chart.Series)
		assert.Len(t, chart.Series, 2)

		// Verify series names
		assert.Equal(t, "Actual", chart.Series[0].Name)
		assert.Equal(t, "Speculated", chart.Series[1].Name)

		// Verify data point counts (should be 15 days)
		assert.Len(t, chart.Series[0].Data, 15)
		assert.Len(t, chart.Series[1].Data, 15)

		// Verify actual data has our log value at the correct position (Tuesday = index 14)
		// Since testDate is Tuesday (weekday 2), and today is also Tuesday,
		// the log should appear at index 14 (today)
		assert.Equal(t, 25, chart.Series[0].Data[14])

		// Verify speculated data uses the formula
		// spec_mean = 25.0 * 1.15 = 28.75
		// For Tuesday (index 14), should use actual value scaled
		var expectedSpeculated float64 = 25.0 * 1.15
		assert.Equal(t, int(expectedSpeculated+0.5), chart.Series[1].Data[14], "Speculated data should match expected value")
	})
}

// TestSpeculateService_GenerateChartConfig_ZeroData tests chart with no data
func TestSpeculateService_GenerateChartConfig_ZeroData(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepositoryExtended{}
	mockConfig := &MockUserConfigProvider{mockPredictionPct: 0.15}

	service := dashboard.NewSpeculateService(mockRepo, mockConfig)

	mockRepo.mockGetLogsByDateRange = func(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
		return []*dto.LogEntry{}, nil
	}

	mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
		return []*dto.ProjectAggregate{}, nil
	}

	// Act
	chart, err := service.GenerateChartConfig(context.Background())

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, chart)
	assert.Len(t, chart.Series, 2)

	// All data points should be zero
	for i := 0; i < 15; i++ {
		assert.Equal(t, 0, chart.Series[0].Data[i], "Actual series day %d should be zero", i)
	}
}

// =============================================================================
// TestSpeculateService_GetSpeculativeData - Speculative data retrieval
// =============================================================================

// TestSpeculateService_GetSpeculativeData tests speculative data generation
func TestSpeculateService_GetSpeculativeData(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepositoryExtended{}
	mockConfig := &MockUserConfigProvider{mockPredictionPct: 0.20} // 20% prediction

	service := dashboard.NewSpeculateService(mockRepo, mockConfig)

	// Setup test data - Use a fixed date matching Tuesday for proper testing
	// April 21, 2026 is a Tuesday (weekday 2)
	testDate := time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC) // Tuesday

	mockRepo.mockGetLogsByDateRange = func(ctx context.Context, start, end time.Time) ([]*dto.LogEntry, error) {
		return []*dto.LogEntry{
			dto.NewLogEntry(1, testDate.Format(time.RFC3339), 0, 25, nil, nil),
		}, nil
	}

	mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
		return []*dto.ProjectAggregate{
			dto.NewProjectAggregate(1, "Project A", 100, 5),
		}, nil
	}

	mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
		return 25.0, nil
	}

	// Act
	withSpeculateServiceFixedDate(t, testDate, func() {
		speculativeData, err := service.GetSpeculativeData(context.Background())

		// Assert
		require.NoError(t, err)
		assert.Len(t, speculativeData, 15)

		// Verify Sunday data (index 14)
		// Actual: 25 pages
		// Speculated: 25 * (1 + 0.20) = 30 pages
		assert.Equal(t, 30, speculativeData[14])

		// Verify zero-filled days use speculative mean
		// spec_mean = 25.0 * 1.20 = 30 (from GetSpeculativeMean)
		// For days without data, should use baseline speculative value
		// Note: GetSpeculativeMean already applies prediction percentage
		assert.Equal(t, 30, speculativeData[13]) // Saturday - uses baseline
	})
}

// =============================================================================
// TestSpeculateService_DateRangeCalculation - Date range tests
// =============================================================================

// TestSpeculateService_GetDateRangeLast15Days tests date range calculation
func TestSpeculateService_GetDateRangeLast15Days(t *testing.T) {
	// Arrange - Use a fixed date for deterministic testing
	fixedDate := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)

	// Act
	withSpeculateServiceFixedDate(t, fixedDate, func() {
		start, end := dashboard.GetDateRangeLast15Days()

		// Assert
		// End date should be today (within same day)
		assert.Equal(t, fixedDate.Year(), end.Year())
		assert.Equal(t, fixedDate.Month(), end.Month())
		assert.Equal(t, fixedDate.Day(), end.Day())

		// Start date should be 14 days ago (to include today in 15-day count)
		expectedStart := fixedDate.AddDate(0, 0, -14)

		// Note: GetToday() truncates to midnight, so we compare dates
		assert.Equal(t, expectedStart.Year(), start.Year())
		assert.Equal(t, expectedStart.Month(), start.Month())
		assert.Equal(t, expectedStart.Day(), start.Day())
	})
}

// TestSpeculateService_GetDateRangeLast15Days_VerifyCount tests that range covers exactly 15 days
func TestSpeculateService_GetDateRangeLast15Days_VerifyCount(t *testing.T) {
	// Arrange - Use a fixed reference date for verification
	fixedDate := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC) // Today's date

	// Act
	withSpeculateServiceFixedDate(t, fixedDate, func() {
		start, end := dashboard.GetDateRangeLast15Days()

		// Assert
		// Verify end date is today (truncated to midnight)
		assert.Equal(t, fixedDate.Year(), end.Year())
		assert.Equal(t, fixedDate.Month(), end.Month())
		assert.Equal(t, fixedDate.Day(), end.Day())

		// Verify start date is 14 days ago
		expectedStart := fixedDate.AddDate(0, 0, -14)
		assert.Equal(t, expectedStart.Year(), start.Year())
		assert.Equal(t, expectedStart.Month(), start.Month())
		assert.Equal(t, expectedStart.Day(), start.Day())

		// Verify the range covers exactly 15 days (including today)
		dayCount := int(end.Sub(start).Hours()/24) + 1
		assert.Equal(t, 15, dayCount)
	})
}

// =============================================================================
// TestSpeculateService_EdgeCases - Edge case tests
// =============================================================================

// TestSpeculateService_NegativePredictionPct tests negative prediction percentage handling
func TestSpeculateService_NegativePredictionPct(t *testing.T) {
	// Arrange
	actualMean := 25.0
	negativePct := -0.5 // -50% prediction

	// Act
	specMean := dashboard.CalculateSpeculativeMean(actualMean, negativePct)

	// Assert
	// spec_mean = 25.0 * (1 - 0.5) = 12.5
	expectedSpecMean := 12.5
	assert.Equal(t, expectedSpecMean, specMean)
}

// TestSpeculateService_ZeroPredictionPct tests zero prediction percentage
func TestSpeculateService_ZeroPredictionPct(t *testing.T) {
	// Arrange
	actualMean := 25.0
	zeroPct := 0.0

	// Act
	specMean := dashboard.CalculateSpeculativeMean(actualMean, zeroPct)

	// Assert
	assert.Equal(t, actualMean, specMean)
}

// TestSpeculateService_LargePredictionPct tests large prediction percentage
func TestSpeculateService_LargePredictionPct(t *testing.T) {
	// Arrange
	actualMean := 25.0
	largePct := 2.0 // 200% prediction

	// Act
	specMean := dashboard.CalculateSpeculativeMean(actualMean, largePct)

	// Assert
	// spec_mean = 25.0 * (1 + 2.0) = 75.0
	expectedSpecMean := 75.0
	assert.Equal(t, expectedSpecMean, specMean)
}

// TestSpeculateService_VerySmallMean tests very small mean values
func TestSpeculateService_VerySmallMean(t *testing.T) {
	// Arrange
	actualMean := 0.001
	predictionPct := 0.15

	// Act
	specMean := dashboard.CalculateSpeculativeMean(actualMean, predictionPct)

	// Assert
	// spec_mean = 0.001 * 1.15 = 0.00115
	// After rounding to 3 decimals: 0.001
	assert.Equal(t, 0.001, specMean)
}
