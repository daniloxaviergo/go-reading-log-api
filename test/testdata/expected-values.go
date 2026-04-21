package testdata

import (
	"context"
	"math"
	"time"

	"go-reading-log-api-next/internal/domain/dto"
)

// ExpectedProject represents the expected structure of a project response
// with all calculated fields populated for comparison testing.
type ExpectedProject struct {
	ID         int64         `json:"id"`
	Name       string        `json:"name"`
	TotalPage  int           `json:"total_page"`
	Page       int           `json:"page"`
	StartedAt  *string       `json:"started_at,omitempty"`
	Progress   *float64      `json:"progress,omitempty"`
	Status     *string       `json:"status,omitempty"`
	LogsCount  *int          `json:"logs_count,omitempty"`
	DaysUnread *int          `json:"days_unreading,omitempty"`
	MedianDay  *float64      `json:"median_day,omitempty"`
	FinishedAt *string       `json:"finished_at,omitempty"`
	Logs       []ExpectedLog `json:"logs,omitempty"`
}

// ExpectedLog represents the expected structure of a log entry
type ExpectedLog struct {
	ID        int64   `json:"id"`
	ProjectID int64   `json:"project_id"`
	Data      string  `json:"data"`
	StartPage int     `json:"start_page"`
	EndPage   int     `json:"end_page"`
	Note      *string `json:"note,omitempty"`
}

// floatPtr returns a pointer to a float64 value
func floatPtr(f float64) *float64 {
	return &f
}

// intPtr returns a pointer to an int value
func intPtr(i int) *int {
	return &i
}

// stringPtr returns a pointer to a string value
func stringPtr(s string) *string {
	return &s
}

// timePtr returns a pointer to a time.Time value
func timePtr(t time.Time) *time.Time {
	return &t
}

// timeToString converts time.Time to RFC3339 string pointer
func timeToString(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

// Project450 contains expected values for Project ID 450
// Based on: test/data/project-450-go.json and test/data/project-450-rails.json
var Project450 = ExpectedProject{
	ID:        450,
	Name:      "História da Igreja VIII.1",
	TotalPage: 691,
	Page:      691,
	// StartedAt calculated from project data: 2026-02-19T00:00:00Z
	StartedAt: timeToString(time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC)),
	// Progress: (691 / 691) * 100 = 100.00%
	Progress: floatPtr(100.0),
	// Status: finished (page >= total_page)
	Status: stringPtr("finished"),
	// Logs count from Rails API
	LogsCount: intPtr(38),
	// Days unreading: Calculated from last log date to today
	// Based on Rails calculation using Date.today in BRT timezone
	DaysUnread: intPtr(58),
	// Median day: page / days_reading (rounded to 2 decimals)
	// Note: This value depends on the exact days calculation
	MedianDay: floatPtr(11.91),
	// FinishedAt: null for completed projects with no logs
	// For Project 450, finished_at is calculated from last log date
	FinishedAt: nil,
}

// ExpectedProjects contains all expected project values for testing
var ExpectedProjects = map[int64]*ExpectedProject{
	450: &Project450,
}

// CalculateExpectedValues computes all derived fields for a project
// This function follows the same logic as internal/domain/models/project.go
func CalculateExpectedValues(ctx context.Context, project *dto.ProjectResponse) *ExpectedProject {
	// Create expected project with base values
	var startedAtTime *time.Time
	if project.StartedAt != nil {
		t, _ := time.Parse(time.RFC3339, *project.StartedAt)
		startedAtTime = &t
	}

	expected := &ExpectedProject{
		ID:        project.ID,
		Name:      project.Name,
		TotalPage: project.TotalPage,
		Page:      project.Page,
		StartedAt: func() *string {
			if startedAtTime != nil {
				return timeToString(*startedAtTime)
			}
			return nil
		}(),
	}

	// Calculate progress
	expected.Progress = CalculateProgress(project.TotalPage, project.Page)

	// Calculate status
	expected.Status = CalculateStatus(project.Page, project.TotalPage, project.Logs)

	// Calculate logs count
	expected.LogsCount = intPtr(len(project.Logs))

	// Calculate days unreading
	expected.DaysUnread = CalculateDaysUnreading(project.Logs, startedAtTime, ctx)

	// Calculate median day
	expected.MedianDay = CalculateMedianDay(project.Page, startedAtTime, ctx)

	// Calculate finished at
	expected.FinishedAt = CalculateFinishedAt(project.Page, project.TotalPage, project.Logs, startedAtTime, ctx)

	// Convert logs to expected format
	expected.Logs = make([]ExpectedLog, len(project.Logs))
	for i, log := range project.Logs {
		// ProjectID is not available in LogResponse DTO anymore
		// It should be passed separately or derived from context
		// Defaulting to 0 for now
		projectID := int64(0)
		expected.Logs[i] = ExpectedLog{
			ID:        log.ID,
			ProjectID: projectID,
			Data:      log.Data.Format(time.RFC3339),
			StartPage: log.StartPage,
			EndPage:   log.EndPage,
			Note:      log.Note,
		}
	}

	return expected
}

// CalculateProgress computes progress percentage
func CalculateProgress(totalPage, page int) *float64 {
	if totalPage <= 0 {
		return floatPtr(0.0)
	}
	if page <= 0 {
		return floatPtr(0.0)
	}

	progress := (float64(page) / float64(totalPage)) * 100.0
	// Round to 2 decimal places
	rounded := math.Round(progress*100) / 100

	// Clamp to 0.00-100.00 range
	if rounded < 0.0 {
		rounded = 0.0
	}
	if rounded > 100.0 {
		rounded = 100.0
	}

	return floatPtr(rounded)
}

// CalculateStatus determines project status
func CalculateStatus(page, totalPage int, logs []*dto.LogResponse) *string {
	// Check for finished first
	if page >= totalPage {
		return stringPtr("finished")
	}

	// Check for unstarted (no logs)
	if len(logs) == 0 {
		return stringPtr("unstarted")
	}

	// For running/sleeping/stopped, we need days_unreading
	// This is a simplified version - full implementation in CalculateDaysUnreading
	days := CalculateDaysUnreading(logs, nil, context.Background())
	if days == nil {
		return stringPtr("stopped")
	}

	// Default thresholds (matching config defaults)
	const emAndamentoRange = 7
	const dormindoRange = 14

	if *days <= emAndamentoRange {
		return stringPtr("running")
	}
	if *days <= dormindoRange {
		return stringPtr("sleeping")
	}
	return stringPtr("stopped")
}

// CalculateDaysUnreading calculates days since last reading activity
func CalculateDaysUnreading(logs []*dto.LogResponse, startedAt *time.Time, ctx context.Context) *int {
	if len(logs) == 0 && startedAt == nil {
		zero := 0
		return &zero
	}

	// Find most recent log date
	var lastReadDate time.Time
	found := false

	for _, log := range logs {
		if log.Data != nil {
			// Use time.Time directly since Data is now *time.Time
			t := *log.Data
			if !found || t.After(lastReadDate) {
				lastReadDate = t
				found = true
			}
		}
	}

	// Fallback to started_at if no log found
	if !found && startedAt != nil {
		lastReadDate = *startedAt
		found = true
	}

	if !found {
		zero := 0
		return &zero
	}

	// Calculate days with timezone support
	now := time.Now()
	tzLocation := getTimezoneFromContext(ctx)

	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tzLocation)
	lastReadDate = time.Date(lastReadDate.Year(), lastReadDate.Month(), lastReadDate.Day(), 0, 0, 0, 0, tzLocation)

	diff := nowDate.Sub(lastReadDate)
	days := int(diff.Hours() / 24)

	if days < 0 {
		zero := 0
		return &zero
	}

	return &days
}

// CalculateMedianDay calculates pages per day reading rate
func CalculateMedianDay(page int, startedAt *time.Time, ctx context.Context) *float64 {
	if startedAt == nil {
		return nil
	}

	tzLocation := getTimezoneFromContext(ctx)

	now := time.Now()
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tzLocation)
	startedAtTime := time.Date(startedAt.Year(), startedAt.Month(), startedAt.Day(), 0, 0, 0, 0, tzLocation)

	diff := nowDate.Sub(startedAtTime)
	daysReading := int(diff.Hours() / 24)

	if daysReading <= 0 {
		return nil
	}

	medianDay := float64(page) / float64(daysReading)
	rounded := math.Round(medianDay*100) / 100

	return &rounded
}

// CalculateFinishedAt calculates estimated completion date
func CalculateFinishedAt(page, totalPage int, logs []*dto.LogResponse, startedAt *time.Time, ctx context.Context) *string {
	if startedAt == nil {
		return nil
	}

	if page <= 0 {
		return nil
	}

	if page >= totalPage {
		// For finished books, return most recent log date or nil
		var latestDate time.Time
		found := false
		for _, log := range logs {
			if log.Data != nil {
				// Use time.Time directly since Data is now *time.Time
				t := *log.Data
				if !found || t.After(latestDate) {
					latestDate = t
					found = true
				}
			}
		}
		if !found {
			return nil
		}
		return timeToString(latestDate)
	}

	if len(logs) == 0 {
		return nil
	}

	medianDay := CalculateMedianDay(page, startedAt, ctx)
	if medianDay == nil || *medianDay <= 0 {
		return nil
	}

	daysToFinish := (float64(totalPage) - float64(page)) / *medianDay
	daysToFinishRounded := int(math.Round(daysToFinish))

	now := time.Now()
	tzLocation := getTimezoneFromContext(ctx)
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tzLocation)

	finishedAt := nowDate.AddDate(0, 0, daysToFinishRounded)
	return timeToString(finishedAt)
}

// ParseLogDate attempts to parse a date string using multiple formats
func ParseLogDate(dateStr string) (time.Time, bool) {
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, true
		}
	}

	return time.Time{}, false
}

// getTimezoneFromContext retrieves timezone from context or returns BRT fallback
func getTimezoneFromContext(ctx context.Context) *time.Location {
	if tz, ok := ctx.Value("timezone").(*time.Location); ok && tz != nil {
		return tz
	}
	return time.FixedZone("BRT", -3*60*60)
}

// GenerateExpectedValues creates expected values for a project from raw JSON data
// This function is used by the generator script to create test data
func GenerateExpectedValues(projectID int64, name string, totalPage, page int, startedAtStr string) *ExpectedProject {
	startedAt, _ := time.Parse(time.RFC3339, startedAtStr)

	return &ExpectedProject{
		ID:        projectID,
		Name:      name,
		TotalPage: totalPage,
		Page:      page,
		StartedAt: timeToString(startedAt),
		Progress:  CalculateProgress(totalPage, page),
		Status:    CalculateStatus(page, totalPage, nil),
		LogsCount: intPtr(0),
		DaysUnread: func() *int {
			if startedAt.IsZero() {
				zero := 0
				return &zero
			}
			now := time.Now()
			tzLocation := time.FixedZone("BRT", -3*60*60)
			nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tzLocation)
			startedAtDate := time.Date(startedAt.Year(), startedAt.Month(), startedAt.Day(), 0, 0, 0, 0, tzLocation)
			days := int(nowDate.Sub(startedAtDate).Hours() / 24)
			if days < 0 {
				return intPtr(0)
			}
			return intPtr(days)
		}(),
		MedianDay: func() *float64 {
			if startedAt.IsZero() {
				result := 0.0
				return &result
			}
			now := time.Now()
			tzLocation := time.FixedZone("BRT", -3*60*60)
			nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tzLocation)
			startedAtDate := time.Date(startedAt.Year(), startedAt.Month(), startedAt.Day(), 0, 0, 0, 0, tzLocation)
			daysReading := int(nowDate.Sub(startedAtDate).Hours() / 24)
			if daysReading <= 0 {
				result := 0.0
				return &result
			}
			medianDay := float64(page) / float64(daysReading)
			rounded := math.Round(medianDay*100) / 100
			return &rounded
		}(),
		FinishedAt: func() *string {
			if page >= totalPage {
				return nil // Finished books have null finished_at when no logs
			}
			if page <= 0 || startedAt.IsZero() {
				return nil
			}
			now := time.Now()
			tzLocation := time.FixedZone("BRT", -3*60*60)
			nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tzLocation)
			daysToFinish := (float64(totalPage) - float64(page)) / *CalculateMedianDay(page, &startedAt, context.WithValue(context.Background(), "timezone", tzLocation))
			finishedAt := nowDate.AddDate(0, 0, int(math.Round(daysToFinish)))
			return timeToString(finishedAt)
		}(),
		Logs: nil,
	}
}
