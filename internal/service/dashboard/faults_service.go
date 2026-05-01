package dashboard

import (
	"context"
	"fmt"
	"math"
	"time"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
)

// FaultsService calculates fault statistics and percentage for dashboard visualization.
// It counts ALL faults (regardless of status) within a date range and calculates
// the fault percentage as (faults_last_30_days / max_faults) * 100.
type FaultsService struct {
	repo       repository.DashboardRepository
	userConfig UserConfigProvider
}

// NewFaultsService creates a new FaultsService with the given repository and config service.
// The service uses dependency injection for testability and follows the same pattern
// as DayService and ProjectsService.
func NewFaultsService(repo repository.DashboardRepository, userConfig UserConfigProvider) *FaultsService {
	return &FaultsService{
		repo:       repo,
		userConfig: userConfig,
	}
}

// GetDateRangeLast30Days returns the date range for the last 30 days.
// The end date is today (midnight) and start date is 30 days ago.
func GetDateRangeLast30Days() (start, end time.Time) {
	end = dto.GetToday()
	start = end.AddDate(0, 0, -30)
	return start, end
}

// CalculatePercentage calculates the fault percentage as (faults / maxFaults) * 100.
// Implements AC-DASH-004 Requirement #3: "Zero faults returns 0% not NaN/error"
// Returns 0.0 if maxFaults is zero or negative (not NaN).
func CalculatePercentage(faults int, maxFaults int) float64 {
	if maxFaults <= 0 {
		return 0.0
	}
	percentage := (float64(faults) / float64(maxFaults)) * 100
	// Round to 2 decimal places as per AC-DASH-004 Requirement #2
	return math.Round(percentage*100) / 100
}

// GetFaultsPercentage calculates the fault percentage for the last 30 days.
// Counts ALL logs in the last 30 days and calculates percentage against max_faults.
// Returns the calculated percentage with 2 decimal precision.
func (s *FaultsService) GetFaultsPercentage(ctx context.Context) (float64, error) {
	// Get date range for last 30 days
	startDate, endDate := GetDateRangeLast30Days()

	// Get faults count from repository
	faults, err := s.repo.GetFaultsByDateRange(ctx, startDate, endDate)
	if err != nil {
		return 0.0, fmt.Errorf("failed to get faults by date range: %w", err)
	}

	// Get max faults from config
	maxFaults := s.userConfig.GetMaxFaults()
	if maxFaults == 0 {
		maxFaults = 10 // Default value as per AC-DASH-004 Requirement #4
	}

	// Calculate percentage with zero-handling
	percentage := CalculatePercentage(faults.FaultCount, maxFaults)

	return percentage, nil
}

// CreateGaugeChart creates an ECharts gauge chart configuration for displaying fault percentage.
// Implements AC-DASH-004: "Gauge shows fault percentage"
// Uses color coding based on percentage:
//   - Green (#4caf50) for < 30%
//   - Orange (#ff9800) for 30-60%
//   - Red (#f44336) for > 60%
func (s *FaultsService) CreateGaugeChart(percentage float64) *dto.EchartConfig {
	// Determine color based on percentage
	color := determineGaugeColor(percentage)

	// Create gauge chart configuration matching ECharts v5 API
	return dto.NewEchartConfig().
		SetTitle("Fault Percentage by Weekday").
		SetTooltip(map[string]interface{}{
			"formatter": "{a} <br/>{b} : {c}%",
		}).
		AddSeries(*dto.NewSeries("Faults", "gauge", []interface{}{percentage}).
			SetItemStyle(map[string]interface{}{
				"color": color,
			}))
}

// determineGaugeColor returns the color string based on percentage for visual feedback.
func determineGaugeColor(percentage float64) string {
	switch {
	case percentage < 30:
		return "#4caf50" // Green - low faults
	case percentage < 60:
		return "#ff9800" // Orange - moderate faults
	default:
		return "#f44336" // Red - high faults
	}
}

// GetFaultsCount returns the count of faults in the last 30 days.
// Useful for debugging and direct access to the raw count.
func (s *FaultsService) GetFaultsCount(ctx context.Context) (int, error) {
	startDate, endDate := GetDateRangeLast30Days()
	faults, err := s.repo.GetFaultsByDateRange(ctx, startDate, endDate)
	if err != nil {
		return 0, fmt.Errorf("failed to get faults by date range: %w", err)
	}
	return faults.FaultCount, nil
}

// GetMaxFaults returns the maximum faults allowed from config.
func (s *FaultsService) GetMaxFaults() int {
	maxFaults := s.userConfig.GetMaxFaults()
	if maxFaults == 0 {
		return 10 // Default value
	}
	return maxFaults
}
