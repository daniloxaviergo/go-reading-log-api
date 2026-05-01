package dashboard

import (
	"context"
	"fmt"
	"math"
	"time"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
)

// SpeculateService compares actual vs predicted reading data and generates chart configurations.
// It implements AC-DASH-005: "Speculated vs Actual Comparison" with:
// - Speculative mean calculation: actual_mean * (1 + prediction_pct)
// - Last 15 days data coverage including today
// - Zero-fill for missing days
type SpeculateService struct {
	repo       repository.DashboardRepository
	userConfig UserConfigProvider
}

// NewSpeculateService creates a new SpeculateService with the given repository and config service.
// The service uses dependency injection for testability and follows the same pattern
// as DayService and FaultsService.
func NewSpeculateService(repo repository.DashboardRepository, userConfig UserConfigProvider) *SpeculateService {
	return &SpeculateService{
		repo:       repo,
		userConfig: userConfig,
	}
}

// GetDateRangeLast15Days returns the date range for the last 15 days.
// The end date is today (midnight) and start date is 14 days ago (to include today in the count).
// This ensures exactly 15 days of data coverage including today.
func GetDateRangeLast15Days() (start, end time.Time) {
	end = dto.GetToday()
	start = end.AddDate(0, 0, -14)
	return start, end
}

// CalculateSpeculativeMean calculates the speculative mean using the formula:
// spec_mean = actual_mean * (1 + prediction_pct)
//
// Implements AC-DASH-005 Requirement #2: "Speculative mean formula correct (actual * (1 + pct))"
// Returns 0.0 if actualMean is zero or negative (not NaN).
func CalculateSpeculativeMean(actualMean float64, predictionPct float64) float64 {
	if actualMean <= 0 {
		return 0.0
	}
	specMean := actualMean * (1 + predictionPct)
	return math.Round(specMean*1000) / 1000
}

// GetSpeculativeMean calculates the speculative mean for the current weekday.
// Uses all available data from project aggregates to calculate a mean,
// then applies the prediction percentage from config.
func (s *SpeculateService) GetSpeculativeMean(ctx context.Context) (float64, error) {
	// Get project aggregates for all projects
	aggregates, err := s.repo.GetProjectAggregates(ctx)
	if err != nil {
		return 0.0, fmt.Errorf("failed to get project aggregates: %w", err)
	}

	if len(aggregates) == 0 {
		return 0.0, nil
	}

	// Get today's weekday (0=Sunday, 1=Monday, ..., 6=Saturday)
	today := GetToday()
	currentWeekday := int(today.Weekday())

	// Calculate total mean and count for the current weekday
	var totalMean float64
	var count int

	for _, agg := range aggregates {
		// For each project, calculate its mean for the current weekday
		projectMean, err := s.repo.GetProjectWeekdayMean(ctx, agg.ProjectID, currentWeekday)
		if err != nil {
			return 0.0, fmt.Errorf("failed to get project weekday mean for project %d: %w", agg.ProjectID, err)
		}
		totalMean += projectMean
		count++
	}

	if count == 0 {
		return 0.0, nil
	}

	// Calculate average mean across all projects
	actualMean := totalMean / float64(count)

	// Get prediction percentage from config
	predictionPct := s.userConfig.GetPredictionPct()

	// Calculate speculative mean
	specMean := CalculateSpeculativeMean(actualMean, predictionPct)

	return specMean, nil
}

// GenerateChartData generates chart data for the last 15 days including today.
// Implements AC-DASH-005 Requirement #3: "Last 15 days data coverage including today"
// Implements AC-DASH-005 Requirement #4: "Missing days zero-filled not omitted"
//
// Returns a map of index (0-14) to page count for each of the 15 days.
// Index 14 represents today, index 0 represents 14 days ago.
func (s *SpeculateService) GenerateChartData(ctx context.Context) (map[int]int, error) {
	// Get date range for last 15 days
	startDate, endDate := GetDateRangeLast15Days()

	// Get all logs within the date range
	logs, err := s.repo.GetLogsByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs by date range: %w", err)
	}

	// Initialize data map with zero values for all 15 indices
	dataMap := make(map[int]int)

	// Calculate total pages read per day (indexed 0-14)
	// Index 14 = today, index 0 = 14 days ago
	for _, log := range logs {
		if log.Data == "" {
			continue
		}

		// Parse the timestamp
		t, err := time.Parse(time.RFC3339, log.Data)
		if err != nil {
			continue // Skip invalid timestamps
		}

		// Calculate days difference from today
		today := GetToday()
		daysDiff := int(today.Sub(t).Hours() / 24)

		// Only include logs within the last 15 days (indices 0-14)
		if daysDiff >= 0 && daysDiff < 15 {
			// Calculate pages read for this log entry
			pagesRead := log.EndPage - log.StartPage
			if pagesRead < 0 {
				pagesRead = 0 // Handle edge case where end_page < start_page
			}

			// Convert daysDiff to index (14 = today, 0 = 14 days ago)
			index := 14 - daysDiff
			dataMap[index] += pagesRead
		}
	}

	return dataMap, nil
}

// GenerateChartConfig generates an ECharts configuration for the speculative vs actual comparison.
// Creates a line chart with two series:
// - "Actual": Actual pages read per day (last 15 days)
// - "Speculated": Predicted pages per day based on config percentage
//
// Implements AC-DASH-005 Requirement #1: "Actual vs predicted comparison implemented"
func (s *SpeculateService) GenerateChartConfig(ctx context.Context) (*dto.EchartConfig, error) {
	// Get actual data for last 15 days
	actualData, err := s.GenerateChartData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate chart data: %w", err)
	}

	// Get speculative mean for the current weekday
	specMean, err := s.GetSpeculativeMean(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get speculative mean: %w", err)
	}

	// Get prediction percentage from config
	predictionPct := s.userConfig.GetPredictionPct()

	// Create 15 data points for the chart (last 15 days including today)
	// Each day gets a data point, zero-filled if no data exists
	var actualSeriesData []interface{}
	var speculatedSeriesData []interface{}

	for i := 0; i < 15; i++ {
		// Get actual data for this index, or zero if not available
		actualPages, exists := actualData[i]
		if !exists {
			actualPages = 0
		}

		// Calculate speculated pages: spec_mean * (1 + prediction_pct) for each day
		// Or use the daily mean scaled by prediction
		var speculatedPages int
		if actualPages > 0 {
			// Scale actual pages by prediction percentage
			speculatedPages = int(math.Round(float64(actualPages) * (1 + predictionPct)))
		} else {
			// Use speculative mean as baseline for zero days
			// Note: specMean already has prediction percentage applied from GetSpeculativeMean
			speculatedPages = int(math.Round(specMean))
			if speculatedPages < 0 {
				speculatedPages = 0
			}
		}

		actualSeriesData = append(actualSeriesData, actualPages)
		speculatedSeriesData = append(speculatedSeriesData, speculatedPages)
	}

	// Create the chart configuration
	chart := dto.NewEchartConfig().
		SetTitle("Speculated vs Actual").
		SetTooltip(map[string]interface{}{
			"trigger":   "axis",
			"formatter": "{a} <br/>{b}: {c}",
		})

	// Add legend
	chart.SetLegend(dto.NewLegend(true, []string{"Actual", "Speculated"}))

	// Create axis with boundary gap setting
	xAxis := dto.NewAxis("category")
	yAxis := dto.NewAxis("value")

	// Set boundary gap for x-axis
	xAxis.BoundaryGap = []bool{false, false}

	// Add grid
	grid := dto.NewGrid()
	grid.Left = "3%"
	grid.Right = "4%"
	grid.Top = "15%"
	grid.Bottom = "3%"

	// Add axes and grid to chart
	chart.SetXAxis(xAxis)
	chart.SetYAxis(yAxis)
	chart.SetGrid(grid)

	// Add Actual series (line)
	actualSeries := dto.NewSeries("Actual", "line", actualSeriesData).
		SetItemStyle(map[string]interface{}{
			"color": "#5470C6",
		}).
		SetLineStyle(map[string]interface{}{
			"width": 2,
		})
	chart.AddSeries(*actualSeries)

	// Add Speculated series (line)
	speculatedSeries := dto.NewSeries("Speculated", "line", speculatedSeriesData).
		SetItemStyle(map[string]interface{}{
			"color": "#91CC75",
		}).
		SetLineStyle(map[string]interface{}{
			"width": 2,
			"type":  "dashed",
		})
	chart.AddSeries(*speculatedSeries)

	return chart, nil
}

// GetLast15DaysData returns the last 15 days of reading data with zero-fill.
// Returns a slice of 15 integers representing pages read per day.
func (s *SpeculateService) GetLast15DaysData(ctx context.Context) ([]int, error) {
	dataMap, err := s.GenerateChartData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate chart data: %w", err)
	}

	// Create 15 data points with zero-fill for missing days
	result := make([]int, 15)

	for i := 0; i < 15; i++ {
		if pages, exists := dataMap[i]; exists {
			result[i] = pages
		} else {
			result[i] = 0 // Zero-fill missing days
		}
	}

	return result, nil
}

// GetSpeculativeData returns speculative page predictions for the last 15 days.
// Uses the speculative mean formula to predict future reading.
func (s *SpeculateService) GetSpeculativeData(ctx context.Context) ([]int, error) {
	// Get actual data for last 15 days
	actualData, err := s.GenerateChartData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate chart data: %w", err)
	}

	// Get actual mean from config (without prediction percentage applied yet)
	// We need the base mean to apply the prediction percentage once
	predictionPct := s.userConfig.GetPredictionPct()

	// Calculate speculative mean with prediction percentage
	specMeanWithPct := CalculateSpeculativeMean(s.userConfig.GetPagesPerDay(), predictionPct)

	// Create 15 data points with zero-fill for missing days
	result := make([]int, 15)

	for i := 0; i < 15; i++ {
		if pages, exists := actualData[i]; exists {
			// Scale actual pages by prediction percentage
			result[i] = int(math.Round(float64(pages) * (1 + predictionPct)))
		} else {
			// Use speculative mean as baseline for zero days
			// Apply prediction percentage to the baseline
			result[i] = int(math.Round(specMeanWithPct))
			if result[i] < 0 {
				result[i] = 0
			}
		}
	}

	return result, nil
}
