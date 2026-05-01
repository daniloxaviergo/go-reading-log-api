package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/service/dashboard"
)

// Helper function to set a fixed test date and restore it after the test
func withFaultsServiceFixedDate(t *testing.T, fixedDate time.Time, fn func()) {
	defer dashboard.SetTestDate(time.Now())
	dashboard.SetTestDate(fixedDate)
	fn()
}

// MockUserConfig is a mock implementation of UserConfigProvider for testing
type MockUserConfig struct {
	MaxFaultsVal     int
	PredictionPctVal float64
	PagesPerDayVal   float64
}

// NewMockUserConfig creates a new MockUserConfig
func NewMockUserConfig(maxFaults int, predictionPct, pagesPerDay float64) *MockUserConfig {
	return &MockUserConfig{
		MaxFaultsVal:     maxFaults,
		PredictionPctVal: predictionPct,
		PagesPerDayVal:   pagesPerDay,
	}
}

// GetMaxFaults returns the configured max faults value
func (m *MockUserConfig) GetMaxFaults() int {
	return m.MaxFaultsVal
}

// GetPredictionPct returns the configured prediction percentage
func (m *MockUserConfig) GetPredictionPct() float64 {
	return m.PredictionPctVal
}

// GetPagesPerDay returns the configured pages per day value
func (m *MockUserConfig) GetPagesPerDay() float64 {
	return m.PagesPerDayVal
}

// =============================================================================
// Unit Tests for FaultsService
// =============================================================================

// TestFaultsService_CalculatePercentage_FullCapacity tests percentage calculation at full capacity
func TestFaultsService_CalculatePercentage_FullCapacity(t *testing.T) {
	// Act
	result := dashboard.CalculatePercentage(5, 10)

	// Assert
	assert.Equal(t, 50.00, result)
}

// TestFaultsService_CalculatePercentage_ZeroFaults tests zero faults returns 0%
func TestFaultsService_CalculatePercentage_ZeroFaults(t *testing.T) {
	// Act
	result := dashboard.CalculatePercentage(0, 10)

	// Assert
	assert.Equal(t, 0.00, result)
}

// TestFaultsService_CalculatePercentage_MaxFaults tests when faults equal max
func TestFaultsService_CalculatePercentage_MaxFaults(t *testing.T) {
	// Act
	result := dashboard.CalculatePercentage(10, 10)

	// Assert
	assert.Equal(t, 100.00, result)
}

// TestFaultsService_CalculatePercentage_ExceedsMax tests when faults exceed max
func TestFaultsService_CalculatePercentage_ExceedsMax(t *testing.T) {
	// Act
	result := dashboard.CalculatePercentage(15, 10)

	// Assert
	assert.Equal(t, 150.00, result)
}

// TestFaultsService_CalculatePercentage_ZeroMax tests zero max faults returns 0% (not NaN)
func TestFaultsService_CalculatePercentage_ZeroMax(t *testing.T) {
	// Act
	result := dashboard.CalculatePercentage(5, 0)

	// Assert
	assert.Equal(t, 0.00, result)
}

// TestFaultsService_CalculatePercentage_NegativeMax tests negative max faults returns 0% (not NaN)
func TestFaultsService_CalculatePercentage_NegativeMax(t *testing.T) {
	// Act
	result := dashboard.CalculatePercentage(5, -10)

	// Assert
	assert.Equal(t, 0.00, result)
}

// TestFaultsService_GetFaultsPercentage_Empty tests with no faults in database
func TestFaultsService_GetFaultsPercentage_Empty(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepository{
		mockGetFaultsByDateRange: func(ctx context.Context, start, end time.Time) (*dto.FaultStats, error) {
			return &dto.FaultStats{FaultCount: 0}, nil
		},
	}

	service := dashboard.NewFaultsService(mockRepo, NewMockUserConfig(10, 0.15, 25.0))
	ctx := context.Background()

	// Act
	percentage, err := service.GetFaultsPercentage(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 0.00, percentage)
}

// TestFaultsService_GetFaultsPercentage_HasData tests with faults in database
func TestFaultsService_GetFaultsPercentage_HasData(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepository{
		mockGetFaultsByDateRange: func(ctx context.Context, start, end time.Time) (*dto.FaultStats, error) {
			return &dto.FaultStats{FaultCount: 3}, nil
		},
	}

	service := dashboard.NewFaultsService(mockRepo, NewMockUserConfig(10, 0.15, 25.0))
	ctx := context.Background()

	// Act
	percentage, err := service.GetFaultsPercentage(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 30.00, percentage)
}

// TestFaultsService_GetFaultsPercentage_ConfigOverride tests with custom max faults
func TestFaultsService_GetFaultsPercentage_ConfigOverride(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepository{
		mockGetFaultsByDateRange: func(ctx context.Context, start, end time.Time) (*dto.FaultStats, error) {
			return &dto.FaultStats{FaultCount: 3}, nil
		},
	}

	service := dashboard.NewFaultsService(mockRepo, NewMockUserConfig(20, 0.15, 25.0))
	ctx := context.Background()

	// Act
	percentage, err := service.GetFaultsPercentage(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 15.00, percentage) // 3/20 = 15%
}

// TestFaultsService_GetFaultsPercentage_DBError tests database error handling
func TestFaultsService_GetFaultsPercentage_DBError(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepository{
		mockGetFaultsByDateRange: func(ctx context.Context, start, end time.Time) (*dto.FaultStats, error) {
			return nil, assert.AnError
		},
	}

	service := dashboard.NewFaultsService(mockRepo, NewMockUserConfig(10, 0.15, 25.0))
	ctx := context.Background()

	// Act
	percentage, err := service.GetFaultsPercentage(ctx)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0.00, percentage)
}

// TestFaultsService_CreateGaugeChart tests gauge chart creation
func TestFaultsService_CreateGaugeChart(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepository{}

	service := dashboard.NewFaultsService(mockRepo, NewMockUserConfig(10, 0.15, 25.0))

	// Act
	gauge := service.CreateGaugeChart(45.67)

	// Assert
	assert.NotNil(t, gauge)
	assert.Equal(t, "Fault Percentage by Weekday", gauge.Title)
	assert.Len(t, gauge.Series, 1)
	assert.Equal(t, "gauge", gauge.Series[0].Type)
	assert.Len(t, gauge.Series[0].Data, 1)
	assert.Equal(t, 45.67, gauge.Series[0].Data[0])
}

// TestFaultsService_CreateGaugeChart_ColorLow tests color determination for low percentage
func TestFaultsService_CreateGaugeChart_ColorLow(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepository{}

	service := dashboard.NewFaultsService(mockRepo, NewMockUserConfig(10, 0.15, 25.0))

	// Act
	gauge := service.CreateGaugeChart(25.0) // Below 30%

	// Assert
	assert.NotNil(t, gauge.Series[0].ItemStyle)
	color := gauge.Series[0].ItemStyle["color"]
	assert.Equal(t, "#4caf50", color) // Green
}

// TestFaultsService_CreateGaugeChart_ColorMedium tests color determination for medium percentage
func TestFaultsService_CreateGaugeChart_ColorMedium(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepository{}

	service := dashboard.NewFaultsService(mockRepo, NewMockUserConfig(10, 0.15, 25.0))

	// Act
	gauge := service.CreateGaugeChart(45.0) // Between 30-60%

	// Assert
	assert.NotNil(t, gauge.Series[0].ItemStyle)
	color := gauge.Series[0].ItemStyle["color"]
	assert.Equal(t, "#ff9800", color) // Orange
}

// TestFaultsService_CreateGaugeChart_ColorHigh tests color determination for high percentage
func TestFaultsService_CreateGaugeChart_ColorHigh(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepository{}

	service := dashboard.NewFaultsService(mockRepo, NewMockUserConfig(10, 0.15, 25.0))

	// Act
	gauge := service.CreateGaugeChart(75.0) // Above 60%

	// Assert
	assert.NotNil(t, gauge.Series[0].ItemStyle)
	color := gauge.Series[0].ItemStyle["color"]
	assert.Equal(t, "#f44336", color) // Red
}

// TestFaultsService_GetToday tests the GetToday function returns consistent date
func TestFaultsService_GetToday(t *testing.T) {
	// Arrange - Use a fixed date for deterministic testing
	fixedDate := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)

	// Act
	withFaultsServiceFixedDate(t, fixedDate, func() {
		today1 := dashboard.GetToday()
		today2 := dashboard.GetToday()

		// Assert
		assert.Equal(t, today1.Year(), today2.Year())
		assert.Equal(t, today1.Month(), today2.Month())
		assert.Equal(t, today1.Day(), today2.Day())
		assert.Equal(t, 0, today1.Hour())
		assert.Equal(t, 0, today1.Minute())
		assert.Equal(t, 0, today1.Second())
	})
}

// TestFaultsService_GetDateRangeLast30Days tests the 30-day date range calculation
func TestFaultsService_GetDateRangeLast30Days(t *testing.T) {
	// Arrange - Use a fixed date for deterministic testing
	fixedDate := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)

	// Act
	withFaultsServiceFixedDate(t, fixedDate, func() {
		start, end := dashboard.GetDateRangeLast30Days()

		// Assert
		expectedEnd := fixedDate
		expectedStart := expectedEnd.AddDate(0, 0, -30)

		assert.Equal(t, expectedEnd.Year(), end.Year())
		assert.Equal(t, expectedEnd.Month(), end.Month())
		assert.Equal(t, expectedEnd.Day(), end.Day())

		assert.Equal(t, expectedStart.Year(), start.Year())
		assert.Equal(t, expectedStart.Month(), start.Month())
		assert.Equal(t, expectedStart.Day(), start.Day())
	})
}

// TestFaultsService_CalculatePercentage_Precision tests 2 decimal precision
func TestFaultsService_CalculatePercentage_Precision(t *testing.T) {
	// Act & Assert - Test various precision scenarios
	testCases := []struct {
		name     string
		faults   int
		max      int
		expected float64
	}{
		{"Exact decimal", 1, 3, 33.33}, // 33.333... -> 33.33
		{"Rounding up", 2, 3, 66.67},   // 66.666... -> 66.67
		{"Exact hundredth", 7, 10, 70.00},
		{"Small percentage", 1, 100, 1.00},
		{"Large percentage", 99, 100, 99.00},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := dashboard.CalculatePercentage(tc.faults, tc.max)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestFaultsService_NewFaultsService tests service initialization
func TestFaultsService_NewFaultsService(t *testing.T) {
	// Arrange
	mockRepo := &MockDashboardRepository{}

	// Act
	service := dashboard.NewFaultsService(mockRepo, NewMockUserConfig(10, 0.15, 25.0))

	// Assert
	assert.NotNil(t, service)
}
