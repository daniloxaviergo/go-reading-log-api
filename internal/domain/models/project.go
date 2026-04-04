package models

import (
	"context"
	"math"
	"time"

	"go-reading-log-api-next/internal/config"
	"go-reading-log-api-next/internal/domain/dto"
)

// Project represents a reading project domain model
type Project struct {
	ctx        context.Context
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	TotalPage  int        `json:"total_page"`
	StartedAt  *time.Time `json:"started_at"`
	Page       int        `json:"page"`
	Reinicia   bool       `json:"reinicia"`
	Progress   *float64   `json:"progress,omitempty"`
	Status     *string    `json:"status,omitempty"`
	LogsCount  *int       `json:"logs_count,omitempty"`
	DaysUnread *int       `json:"days_unreading,omitempty"`
	MedianDay  *float64   `json:"median_day,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

// Status constants for project status determination
const (
	StatusUnstarted = "unstarted"
	StatusFinished  = "finished"
	StatusRunning   = "running"
	StatusSleeping  = "sleeping"
	StatusStopped   = "stopped"
)

// NewProject creates a new Project with context
func NewProject(ctx context.Context, id int64, name string, totalPage int, page int, reinicia bool) *Project {
	return &Project{
		ctx:       ctx,
		ID:        id,
		Name:      name,
		TotalPage: totalPage,
		Page:      page,
		Reinicia:  reinicia,
	}
}

// stringPtr returns a pointer to a string value
func stringPtr(s string) *string {
	return &s
}

// CalculateDaysUnreading calculates the number of days since the last reading activity
// Uses the last log's data field if available, otherwise uses StartedAt
// Returns the number of days as a non-negative integer
func (p *Project) CalculateDaysUnreading(logs []*dto.LogResponse) *int {
	// If no logs and no started_at, return 0 (cannot calculate, but return 0)
	if len(logs) == 0 && p.StartedAt == nil {
		zero := 0
		return &zero
	}

	// Find the most recent log date
	var lastReadDate time.Time
	found := false

	for _, log := range logs {
		if log.Data != nil {
			// Parse the log data field (expected format: YYYY-MM-DD)
			if t, err := time.Parse("2006-01-02", *log.Data); err == nil {
				lastReadDate = t
				found = true
				break
			}
		}
	}

	// If no log with data found, use started_at
	if !found && p.StartedAt != nil {
		lastReadDate = *p.StartedAt
		found = true
	}

	// If still no date, return 0
	if !found {
		zero := 0
		return &zero
	}

	// Calculate days unreading
	now := time.Now()
	// Use date-only comparison to match Rails behavior (Date.today)
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	lastReadDate = time.Date(lastReadDate.Year(), lastReadDate.Month(), lastReadDate.Day(), 0, 0, 0, 0, time.UTC)

	// Calculate difference in days
	diff := nowDate.Sub(lastReadDate)
	days := int(diff.Hours() / 24)

	// Clamp to 0 if negative (future dates)
	if days < 0 {
		zero := 0
		return &zero
	}

	return &days
}

// CalculateStatus determines the project status based on logs count and days_unreading
// Status determination logic:
//   - unstarted: No logs exist for the project
//   - finished: page >= total_page (reading complete)
//   - running: days_unreading <= em_andamento_range (default 7 days)
//   - sleeping: days_unreading <= dormindo_range (default 14 days)
//   - stopped: All other cases
func (p *Project) CalculateStatus(logs []*dto.LogResponse, config *config.Config) *string {
	// 1. Check for unstarted (no logs)
	if len(logs) == 0 {
		return stringPtr(StatusUnstarted)
	}

	// 2. Check for finished (page >= total_page)
	if p.Page >= p.TotalPage {
		return stringPtr(StatusFinished)
	}

	// 3. Calculate days_unreading
	daysUnreading := p.CalculateDaysUnreading(logs)

	// If we can't calculate days, return stopped
	if daysUnreading == nil {
		return stringPtr(StatusStopped)
	}

	// 4. Check running (days <= em_andamento_range)
	if *daysUnreading <= config.GetEmAndamentoRange() {
		return stringPtr(StatusRunning)
	}

	// 5. Check sleeping (days <= dormindo_range)
	if *daysUnreading <= config.GetDormindoRange() {
		return stringPtr(StatusSleeping)
	}

	// 6. All other cases → stopped
	return stringPtr(StatusStopped)
}

// GetContext returns the embedded context
func (p *Project) GetContext() context.Context {
	if p.ctx == nil {
		return context.Background()
	}
	return p.ctx
}

// SetContext sets the context for the Project
func (p *Project) SetContext(ctx context.Context) {
	p.ctx = ctx
}

// CalculateProgress calculates the progress percentage as (page / total_page) * 100
// rounded to 2 decimal places, clamped to 0.00-100.00 range.
// Returns 0.00 for edge cases (zero/negative total_page, null/zero page).
func (p *Project) CalculateProgress() *float64 {
	// Handle edge cases
	if p.TotalPage <= 0 {
		result := 0.0
		return &result
	}

	if p.Page <= 0 {
		result := 0.0
		return &result
	}

	// Calculate progress: (page / total_page) * 100
	progress := (float64(p.Page) / float64(p.TotalPage)) * 100.0

	// Round to 2 decimal places using math.Round (round half-up)
	// Multiply by 100, round, then divide by 100
	rounded := math.Round(progress*100) / 100

	// Clamp to 0.00-100.00 range
	if rounded < 0.0 {
		rounded = 0.0
	}
	if rounded > 100.0 {
		rounded = 100.0
	}

	return &rounded
}

// CalculateMedianDay calculates median_day as (page / days_reading).round(2)
// where days_reading is the number of days since started_at
// Returns 0.00 for edge cases (zero/negative days_reading, no started_at)
func (p *Project) CalculateMedianDay() *float64 {
	// Handle edge case: no started_at date
	if p.StartedAt == nil {
		result := 0.0
		return &result
	}

	// Calculate days_reading (days since started_at)
	now := time.Now()
	// Use date-only comparison to match Rails behavior (Date.today)
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startedAt := time.Date(p.StartedAt.Year(), p.StartedAt.Month(), p.StartedAt.Day(), 0, 0, 0, 0, time.UTC)

	// Calculate difference in days
	diff := nowDate.Sub(startedAt)
	daysReading := int(diff.Hours() / 24)

	// Handle edge case: zero or negative days_reading
	if daysReading <= 0 {
		result := 0.0
		return &result
	}

	// Calculate: (page / days_reading).round(2)
	// Note: Rails rounds the result, not days_reading
	// page and days_reading are converted to float, divided, then rounded to 2 decimal places
	pageFloat := float64(p.Page)
	daysReadingFloat := float64(daysReading)

	// Divide first, then round the result to 2 decimal places
	medianDay := pageFloat / daysReadingFloat
	rounded := math.Round(medianDay*100) / 100

	return &rounded
}

// CalculateLogsCount calculates the number of log entries in the logs array
// Matches Rails behavior: def logs_count; logs.size; end
// Returns the count as a non-negative integer
func (p *Project) CalculateLogsCount(logs []*dto.LogResponse) *int {
	count := len(logs)
	return &count
}

// CalculateFinishedAt calculates the estimated date when the book will be finished
// based on the reading rate (median_day).
// Formula: days_to_finish = (total_page - page) / median_day, finished_at = today + days_to_finish
// Returns nil if:
//   - progress is 100% (page >= total_page) and no logs exist
//   - page <= 0 (would cause division by zero)
//   - no started_at date (can't calculate median_day)
func (p *Project) CalculateFinishedAt(logs []*dto.LogResponse) *time.Time {
	// Handle edge case: no started_at date (can't calculate median_day)
	if p.StartedAt == nil {
		return nil
	}

	// Handle edge case: page <= 0 (would cause division by zero)
	if p.Page <= 0 {
		return nil
	}

	// Handle edge case: page >= total_page (finished book)
	// In this case, return the most recent log's data field as a date
	if p.Page >= p.TotalPage {
		// Find the most recent log with a data field
		var latestDate time.Time
		found := false
		for _, log := range logs {
			if log.Data != nil {
				// Parse the log data field (expected format: YYYY-MM-DD)
				if t, err := time.Parse("2006-01-02", *log.Data); err == nil {
					if !found || t.After(latestDate) {
						latestDate = t
						found = true
					}
				}
			}
		}
		// If no logs with data found, return nil
		if !found {
			return nil
		}
		return &latestDate
	}

	// Calculate median_day first (this handles all edge cases)
	medianDay := p.CalculateMedianDay()

	// Handle edge case: median_day is 0 or nil (would cause division by zero)
	if medianDay == nil || *medianDay <= 0 {
		return nil
	}

	// Calculate: days_to_finish = (total_page - page) / median_day
	totalPageFloat := float64(p.TotalPage)
	pageFloat := float64(p.Page)
	daysToFinish := (totalPageFloat - pageFloat) / *medianDay

	// Round to nearest integer to match Rails behavior
	daysToFinishRounded := int(math.Round(daysToFinish))

	// Calculate: finished_at = today + days_to_finish days
	now := time.Now()
	// Use date-only comparison to match Rails behavior (Date.today)
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	finishedAt := nowDate.AddDate(0, 0, daysToFinishRounded)
	return &finishedAt
}
