package dashboard

import (
	"context"
	"fmt"
	"time"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
)

// WeekdayFaultsService calculates and formats weekday fault data for radar chart visualization.
// It groups faults by weekday (0-6 = Sunday-Saturday) over a 6-month date range.
type WeekdayFaultsService struct {
	repo       repository.DashboardRepository
	userConfig UserConfigProvider
}

// NewWeekdayFaultsService creates a new WeekdayFaultsService with the given repository and config service.
// The service uses dependency injection for testability and follows the same pattern
// as DayService and FaultsService.
func NewWeekdayFaultsService(repo repository.DashboardRepository, userConfig UserConfigProvider) *WeekdayFaultsService {
	return &WeekdayFaultsService{
		repo:       repo,
		userConfig: userConfig,
	}
}

// GetDateRangeLast6Months returns the date range for the last 6 months.
// The end date is today (midnight) and start date is 6 months ago.
// Uses AddDate(0, -6, 0) for exact month arithmetic that handles varying month lengths.
func GetDateRangeLast6Months() (start, end time.Time) {
	end = dto.GetToday()
	start = end.AddDate(0, -6, 0)
	return start, end
}

// GetWeekdayFaults retrieves fault data grouped by weekday for the last 6 months.
// Returns a WeekdayFaults response with all 7 days (0-6) present in the map,
// defaulting to 0 for days without any faults.
func (s *WeekdayFaultsService) GetWeekdayFaults(ctx context.Context) (*dto.WeekdayFaults, error) {
	// Get date range for last 6 months
	startDate, endDate := GetDateRangeLast6Months()

	// Get weekday faults data from repository
	faults, err := s.repo.GetWeekdayFaults(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekday faults: %w", err)
	}

	return faults, nil
}

// CreateRadarChart creates an ECharts radar chart configuration for displaying weekday fault distribution.
// Implements AC-DASH-006: "Chart shows fault counts by weekday"
// Converts integer fault counts to float64 for ECharts compatibility.
func (s *WeekdayFaultsService) CreateRadarChart(faults map[int]int) *dto.EchartConfig {
	// Ensure all 7 days are present in the data array
	data := make([]interface{}, 7)
	for i := 0; i < 7; i++ {
		if val, exists := faults[i]; exists {
			data[i] = float64(val) // Convert int to float64 for ECharts
		} else {
			data[i] = 0.0
		}
	}

	// Create radar chart configuration matching ECharts v5 API
	return dto.NewEchartConfig().
		SetTitle("Faults by Weekday").
		SetTooltip(map[string]interface{}{
			"trigger": "item",
		}).
		AddSeries(*dto.NewSeries("Faults", "radar", data).
			SetItemStyle(map[string]interface{}{
				"color": "#54a8ff",
			}))
}

// ValidateOutput validates the weekday faults data against acceptance criteria.
// Checks:
//   - All 7 weekdays are present (0-6)
//   - All counts are non-negative integers
func (s *WeekdayFaultsService) ValidateOutput(faults map[int]int) error {
	// Check all 7 weekdays are present
	for i := 0; i < 7; i++ {
		if _, exists := faults[i]; !exists {
			return fmt.Errorf("weekday %d is missing from output", i)
		}
	}

	// Check all counts are non-negative
	for i, count := range faults {
		if count < 0 {
			return fmt.Errorf("weekday %d has negative count: %d", i, count)
		}
	}

	return nil
}
