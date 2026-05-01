package dashboard

import (
	"context"
	"fmt"
	"math"
	"time"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
)

// MeanProgressService calculates mean progress statistics with visual map colors.
// It implements the requirements from AC-DASH-007:
// - Calculate daily progress as (daily_pages / mean_pages) * 100 - 100
// - Apply visual map color ranges (gray 0-10%, cyan 10-20%, blue 20-50%, green >50%, red negative)
// - Cover last 30 days of data
type MeanProgressService struct {
	repo       repository.DashboardRepository
	userConfig UserConfigProvider
}

// NewMeanProgressService creates a new MeanProgressService with the given repository and config service.
func NewMeanProgressService(repo repository.DashboardRepository, userConfig UserConfigProvider) *MeanProgressService {
	return &MeanProgressService{
		repo:       repo,
		userConfig: userConfig,
	}
}

// GetDateRangeLast30Days returns the date range for the last 30 days.
// The end date is today (midnight) and start date is 29 days ago (inclusive = 30 days).
// Note: This function conflicts with faults_service.go's version which uses -30 days.
// We use a different name to avoid collision, or could refactor to share.
func GetDateRangeLast30DaysMeanProgress() (start, end time.Time) {
	end = dto.GetToday()
	start = end.AddDate(0, 0, -29)
	return start, end
}

// CalculateDailyProgress calculates the daily progress percentage.
// Formula: (daily_pages / mean_pages) * 100 - 100
// Returns 0.0 if meanPages is zero (avoids division by zero).
func CalculateDailyProgress(dailyPages, meanPages float64) float64 {
	if meanPages == 0 {
		return 0.0
	}
	progress := (dailyPages/meanPages)*100 - 100
	return math.Round(progress*1000) / 1000
}

// GetColorForProgress returns the color string based on progress percentage.
// Color ranges (half-open intervals):
//   - Red (#ff4d4f): negative (< 0%)
//   - Gray (#959595): 0% to < 10%
//   - Cyan (#1890ff): 10% to < 20%
//   - Blue (#108ee9): 20% to < 50%
//   - Green (#67c23a): >= 50%
func GetColorForProgress(progress float64) string {
	switch {
	case progress < 0:
		return "#ff4d4f" // Red
	case progress < 10:
		return "#959595" // Gray
	case progress < 20:
		return "#1890ff" // Cyan
	case progress < 50:
		return "#108ee9" // Blue
	default: // >= 50
		return "#67c23a" // Green
	}
}

// GetMeanProgressData calculates mean progress data for the last 30 days.
// Returns a slice of ProgressDay entries with calculated progress and colors.
// Always returns exactly 30 data points, one for each day in the range.
func (s *MeanProgressService) GetMeanProgressData(ctx context.Context) ([]*dto.ProgressDay, error) {
	// Get date range for last 30 days
	startDate, endDate := GetDateRangeLast30DaysMeanProgress()

	// Get logs from repository
	logs, err := s.repo.GetLogsByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs by date range: %w", err)
	}

	// Calculate mean pages per day across all projects in the date range
	totalPages := 0
	for _, log := range logs {
		pageDiff := log.EndPage - log.StartPage
		if pageDiff > 0 {
			totalPages += pageDiff
		}
	}

	// Mean pages = total pages / 30 days
	meanPages := float64(totalPages) / 30.0

	// Build a map of logs by date for quick lookup
	logsByDate := make(map[string]*dto.LogEntry)
	for _, log := range logs {
		// Parse the log date string and extract just the date part
		if logTime, err := time.Parse(time.RFC3339, log.Data); err == nil {
			dateKey := time.Date(logTime.Year(), logTime.Month(), logTime.Day(), 0, 0, 0, 0, time.UTC).Format("2006-01-02")
			logsByDate[dateKey] = log
		}
	}

	// Build progress data for all 30 days
	progressData := make([]*dto.ProgressDay, 30)
	for i := 0; i < 30; i++ {
		dayDate := endDate.AddDate(0, 0, -29+i)
		dayDate = time.Date(dayDate.Year(), dayDate.Month(), dayDate.Day(), 0, 0, 0, 0, time.UTC)
		dateKey := dayDate.Format("2006-01-02")

		var dailyPages float64
		var progress float64
		var color string

		// Check if we have logs for this day
		if log, exists := logsByDate[dateKey]; exists {
			dailyPages = float64(log.EndPage - log.StartPage)
			if dailyPages < 0 {
				dailyPages = 0
			}

			// Calculate progress if we have meaningful data
			if meanPages > 0 {
				progress = CalculateDailyProgress(dailyPages, meanPages)
			} else {
				progress = 0.0
			}
		} else {
			// No logs for this day
			dailyPages = 0
			progress = 0.0
		}

		// Get color for progress
		color = GetColorForProgress(progress)

		// Create progress day entry with RFC3339 formatted date string
		progressData[i] = &dto.ProgressDay{
			Date:       dayDate.Format(time.RFC3339),
			DailyPages: dailyPages,
			Progress:   progress,
			Color:      color,
		}
	}

	return progressData, nil
}

// GenerateChartConfig generates an ECharts configuration for mean progress visualization.
// Creates a line chart with color-coded data points using visual map.
func (s *MeanProgressService) GenerateChartConfig(ctx context.Context) (*dto.EchartConfig, error) {
	// Get mean progress data
	progressData, err := s.GetMeanProgressData(ctx)
	if err != nil {
		return nil, err
	}

	// Create chart configuration
	chart := dto.NewEchartConfig().
		SetTitle("Mean Progress").
		SetTooltip(map[string]interface{}{
			"trigger": "axis",
		})

	// Add legend
	legend := dto.NewLegend(true, []string{"Progress"})
	chart.SetLegend(legend)

	// Create axis configurations
	xAxis := dto.NewAxis("category")
	xAxis.Name = "Date"
	chart.SetXAxis(xAxis)

	yAxis := dto.NewAxis("value")
	yAxis.Name = "Progress (%)"
	chart.SetYAxis(yAxis)

	// Prepare data for series
	data := make([]interface{}, len(progressData))
	colors := make([]string, len(progressData))

	for i, pd := range progressData {
		data[i] = pd.Progress
		colors[i] = pd.Color
	}

	// Create series with visual map configuration
	series := dto.NewSeries("Progress", "line", data)

	// Configure item style with color array for visual map
	itemStyle := map[string]interface{}{
		"color": colors,
	}
	series.SetItemStyle(itemStyle)

	chart.AddSeries(*series)

	return chart, nil
}
