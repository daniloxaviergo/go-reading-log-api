package dashboard

import (
	"context"
	"fmt"
	"math"
	"time"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
)

// UserConfigProvider is an interface for accessing user configuration
// This allows for easier testing and dependency injection
type UserConfigProvider interface {
	GetPredictionPct() float64
	GetMaxFaults() int
	GetPagesPerDay() float64
}

// DayService calculates weekly statistics for dashboard views
type DayService struct {
	repo       repository.DashboardRepository
	userConfig UserConfigProvider
}

// NewDayService creates a new DayService with the given repository and config service
func NewDayService(repo repository.DashboardRepository, userConfig UserConfigProvider) *DayService {
	return &DayService{
		repo:       repo,
		userConfig: userConfig,
	}
}

// Re-export date abstraction from dto package for backward compatibility.
// All services should import from the shared location (dto/dashboard.go)
// These are kept here to maintain existing API, but new code should use dto.GetToday()
var (
	GetTodayFunc = dto.GetTodayFunc
	GetToday     = dto.GetToday
	SetTestDate  = dto.SetTestDate
)

// CalculateWeeklyStats calculates all weekly statistics for the dashboard.
// Uses GetToday() for consistent date references and ensures all float values
// are rounded to 3 decimal places as per AC-DASH-001.
//
// Calculations:
//   - previous_week_pages: Sum of pages from 14-7 days ago (previous week)
//   - last_week_pages: Sum of pages from 7 days ago to today (current week so far)
//   - per_pages: Ratio of last_week_pages / previous_week_pages * 100
//   - mean_day: Average pages per day for the current weekday across all weeks
//   - spec_mean_day: mean_day * (1 + prediction_pct from config)
func (s *DayService) CalculateWeeklyStats(ctx context.Context) (*dto.StatsData, error) {
	// Use GetToday() for consistent date references
	today := GetToday()

	// Calculate date ranges
	// Previous week: 14 days ago to 7 days ago
	prevWeekStart := today.AddDate(0, 0, -14)
	prevWeekEnd := today.AddDate(0, 0, -7)

	// Last week (current week so far): 7 days ago to today
	lastWeekStart := today.AddDate(0, 0, -7)
	lastWeekEnd := today

	// Get project aggregates for the entire period (all data)
	aggregates, err := s.repo.GetProjectAggregates(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get project aggregates: %w", err)
	}

	// Calculate previous week pages (14-7 days ago)
	previousWeekPages, err := s.CalculatePeriodPages(ctx, prevWeekStart, prevWeekEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate previous week pages: %w", err)
	}

	// Calculate last week pages (7 days ago to today)
	lastWeekPages, err := s.CalculatePeriodPages(ctx, lastWeekStart, lastWeekEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate last week pages: %w", err)
	}

	// Calculate per_pages ratio with 3 decimal precision (returns nil when previousWeekPages = 0)
	perPages := s.CalculatePerPagesRatio(lastWeekPages, previousWeekPages)

	// Calculate mean_day from all available data
	meanDay, err := s.CalculateMeanDay(ctx, aggregates)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate mean day: %w", err)
	}

	// Calculate spec_mean_day from config
	specMeanDay := s.CalculateSpecMeanDay(meanDay)

	// Get total pages and current pages from aggregates for progress calculation
	totalPages, pages := s.CalculateProgress(aggregates)

	// Build response with all calculated values
	// perPages is already *float64 (nil when previousWeekPages = 0)
	response := dto.NewStatsData().
		SetPreviousWeekPages(previousWeekPages).
		SetLastWeekPages(lastWeekPages).
		SetPerPages(perPages).
		SetMeanDay(meanDay).
		SetSpecMeanDay(specMeanDay).
		SetTotalPages(totalPages).
		SetPages(pages)

	// Round all float values to 3 decimal places as per AC-DASH-001
	response.RoundToThreeDecimals()

	return response, nil
}

// CalculatePeriodPages calculates the total pages read within a date range
func (s *DayService) CalculatePeriodPages(ctx context.Context, start, end time.Time) (int, error) {
	return s.repo.CalculatePeriodPages(ctx, start, end)
}

// CalculatePerPagesRatio calculates the ratio of last week pages to previous week pages.
// Returns nil if previous week pages is zero (to match Rails API behavior).
func (s *DayService) CalculatePerPagesRatio(lastWeekPages, previousWeekPages int) *float64 {
	if previousWeekPages == 0 {
		return nil
	}
	ratio := float64(lastWeekPages) / float64(previousWeekPages) * 100
	rounded := math.Round(ratio*1000) / 1000
	return &rounded
}

// CalculateMeanDay calculates the average pages per day for the current weekday.
// Uses all available data from project aggregates to calculate a mean.
func (s *DayService) CalculateMeanDay(ctx context.Context, aggregates []*dto.ProjectAggregate) (float64, error) {
	if len(aggregates) == 0 {
		return 0.0, nil
	}

	// Get today's weekday (0=Sunday, 1=Monday, ..., 6=Saturday)
	today := GetToday()
	currentWeekday := int(today.Weekday())

	// Calculate total pages and count for the current weekday
	var totalPages float64
	var count int

	for _, agg := range aggregates {
		// For each project, we need to calculate its mean for the current weekday
		// This requires querying logs for this specific project and weekday
		projectMean, err := s.repo.GetProjectWeekdayMean(ctx, agg.ProjectID, currentWeekday)
		if err != nil {
			return 0.0, fmt.Errorf("failed to get project weekday mean for project %d: %w", agg.ProjectID, err)
		}
		totalPages += projectMean
		count++
	}

	if count == 0 {
		return 0.0, nil
	}

	mean := totalPages / float64(count)
	return math.Round(mean*1000) / 1000, nil
}

// CalculateSpecMeanDay calculates the speculative mean day from config
func (s *DayService) CalculateSpecMeanDay(meanDay float64) float64 {
	predictionPct := s.userConfig.GetPredictionPct()
	specMean := meanDay * (1 + predictionPct)
	return math.Round(specMean*1000) / 1000
}

// CalculateProgress calculates overall progress from aggregates
func (s *DayService) CalculateProgress(aggregates []*dto.ProjectAggregate) (int, int) {
	var totalPages int
	var pages int

	for _, agg := range aggregates {
		totalPages += agg.TotalPages
		pages += agg.LogCount // Using log count as proxy for current progress
	}

	return totalPages, pages
}
