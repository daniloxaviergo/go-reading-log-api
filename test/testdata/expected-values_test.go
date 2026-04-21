package testdata

import (
	"context"
	"math"
	"testing"
	"time"

	"go-reading-log-api-next/internal/domain/dto"
)

// TestExpectedValues_Project450 tests all calculated fields for Project 450
func TestExpectedValues_Project450(t *testing.T) {
	tests := []struct {
		name    string
		fn      func() error
		decimal int // Number of decimal places for float comparison
	}{
		{
			name: "progress calculation",
			fn: func() error {
				expectedProgress := 100.0
				actualProgress := *Project450ExpectedValues.Progress
				if !floatEqual(actualProgress, expectedProgress, 2) {
					t.Errorf("expected progress %v, got %v", expectedProgress, actualProgress)
				}
				return nil
			},
			decimal: 2,
		},
		{
			name: "status is finished",
			fn: func() error {
				if *Project450ExpectedValues.Status != "finished" {
					t.Errorf("expected status 'finished', got '%s'", *Project450ExpectedValues.Status)
				}
				return nil
			},
		},
		{
			name: "logs_count matches",
			fn: func() error {
				expectedCount := Project450LogCount
				actualCount := *Project450ExpectedValues.LogsCount
				if actualCount != expectedCount {
					t.Errorf("expected logs_count %d, got %d", expectedCount, actualCount)
				}
				return nil
			},
		},
		{
			name: "days_unreading is calculated",
			fn: func() error {
				if Project450ExpectedValues.DaysUnread == nil {
					t.Errorf("days_unreading should not be nil")
				}
				if *Project450ExpectedValues.DaysUnread < 0 {
					t.Errorf("days_unreading should be non-negative, got %d", *Project450ExpectedValues.DaysUnread)
				}
				return nil
			},
		},
		{
			name: "median_day is calculated",
			fn: func() error {
				if Project450ExpectedValues.MedianDay == nil {
					t.Errorf("median_day should not be nil")
				}
				if *Project450ExpectedValues.MedianDay <= 0 {
					t.Errorf("median_day should be positive, got %v", *Project450ExpectedValues.MedianDay)
				}
				return nil
			},
		},
		{
			name: "finished_at is null for completed project",
			fn: func() error {
				if Project450ExpectedValues.FinishedAt != nil {
					t.Errorf("finished_at should be nil for completed project, got %s", *Project450ExpectedValues.FinishedAt)
				}
				return nil
			},
		},
		{
			name: "logs are populated",
			fn: func() error {
				if len(Project450ExpectedValues.Logs) == 0 {
					t.Errorf("logs should not be empty")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); err != nil {
				t.Errorf("%s: %v", tt.name, err)
			}
		})
	}
}

// TestCalculateProgress tests the CalculateProgress function
func TestCalculateProgress(t *testing.T) {
	tests := []struct {
		name      string
		totalPage int
		page      int
		expected  float64
	}{
		{
			name:      "50% progress",
			totalPage: 200,
			page:      100,
			expected:  50.0,
		},
		{
			name:      "100% progress",
			totalPage: 100,
			page:      100,
			expected:  100.0,
		},
		{
			name:      "0% progress (zero page)",
			totalPage: 100,
			page:      0,
			expected:  0.0,
		},
		{
			name:      "0% progress (zero total_page)",
			totalPage: 0,
			page:      50,
			expected:  0.0,
		},
		{
			name:      "clamped to 100%",
			totalPage: 100,
			page:      200,
			expected:  100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateProgress(tt.totalPage, tt.page)
			if result == nil {
				t.Fatal("result should not be nil")
			}
			if !floatEqual(*result, tt.expected, 2) {
				t.Errorf("expected %v, got %v", tt.expected, *result)
			}
		})
	}
}

// TestCalculateDaysUnreading tests the CalculateDaysUnreading function
func TestCalculateDaysUnreading(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		logs      []*dto.LogResponse
		startedAt *time.Time
		expectNil bool
		minDays   int
		maxDays   int
	}{
		{
			name:      "no logs, no started_at",
			logs:      nil,
			startedAt: func() *time.Time { return nil }(),
			expectNil: false,
			minDays:   0,
			maxDays:   0,
		},
		{
			name: "with log date",
			logs: []*dto.LogResponse{
				{Data: stringPtr("2026-04-18")},
			},
			startedAt: func() *time.Time {
				t := time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC)
				return &t
			}(),
			expectNil: false,
			minDays:   0,
			maxDays:   100,
		},
		{
			name: "with started_at fallback",
			logs: []*dto.LogResponse{},
			startedAt: func() *time.Time {
				t := time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC)
				return &t
			}(),
			expectNil: false,
			minDays:   0,
			maxDays:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateDaysUnreading(tt.logs, tt.startedAt, ctx)

			if tt.expectNil && result != nil {
				t.Errorf("expected nil, got %v", *result)
			}
			if !tt.expectNil && result == nil {
				t.Error("expected non-nil result")
			}
			if result != nil {
				if *result < tt.minDays || *result > tt.maxDays {
					t.Errorf("days_unreading %d out of range [%d, %d]", *result, tt.minDays, tt.maxDays)
				}
			}
		})
	}
}

// TestCalculateMedianDay tests the CalculateMedianDay function
func TestCalculateMedianDay(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		page      int
		startedAt *time.Time
		expectNil bool
		minValue  float64
		maxValue  float64
	}{
		{
			name:      "zero page",
			page:      0,
			startedAt: func() *time.Time { return nil }(),
			expectNil: true,
		},
		{
			name:      "no started_at",
			page:      100,
			startedAt: func() *time.Time { return nil }(),
			expectNil: true,
		},
		{
			name: "valid calculation",
			page: 100,
			startedAt: func() *time.Time {
				t := time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC)
				return &t
			}(),
			expectNil: false,
			minValue:  0.0,
			maxValue:  100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateMedianDay(tt.page, tt.startedAt, ctx)

			if tt.expectNil && result != nil {
				t.Errorf("expected nil, got %v", *result)
			}
			if !tt.expectNil && result == nil {
				t.Error("expected non-nil result")
			}
			if result != nil {
				if *result < tt.minValue || *result > tt.maxValue {
					t.Errorf("median_day %v out of range [%v, %v]", *result, tt.minValue, tt.maxValue)
				}
			}
		})
	}
}

// TestCalculateStatus tests the CalculateStatus function
func TestCalculateStatus(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		totalPage int
		logs      []*dto.LogResponse
		expected  string
	}{
		{
			name:      "finished (page >= total)",
			page:      100,
			totalPage: 100,
			logs:      nil,
			expected:  "finished",
		},
		{
			name:      "unstarted (no logs)",
			page:      50,
			totalPage: 100,
			logs:      []*dto.LogResponse{},
			expected:  "unstarted",
		},
		{
			name:      "running (days <= 7)",
			page:      50,
			totalPage: 100,
			logs:      []*dto.LogResponse{{Data: stringPtr("2026-04-18")}},
			expected:  "running",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateStatus(tt.page, tt.totalPage, tt.logs)
			if result == nil {
				t.Fatal("result should not be nil")
			}
			if *result != tt.expected {
				t.Errorf("expected status '%s', got '%s'", tt.expected, *result)
			}
		})
	}
}

// TestGenerateExpectedValues tests the GenerateExpectedValues function
func TestGenerateExpectedValues(t *testing.T) {
	expected := GenerateExpectedValues(
		450,
		"História da Igreja VIII.1",
		691,
		691,
		"2026-02-19T00:00:00Z",
	)

	if expected.ID != 450 {
		t.Errorf("expected ID 450, got %d", expected.ID)
	}
	if expected.Name != "História da Igreja VIII.1" {
		t.Errorf("expected name 'História da Igreja VIII.1', got '%s'", expected.Name)
	}
	if expected.TotalPage != 691 {
		t.Errorf("expected total_page 691, got %d", expected.TotalPage)
	}
	if expected.Page != 691 {
		t.Errorf("expected page 691, got %d", expected.Page)
	}

	// Verify calculated fields
	if expected.Progress == nil {
		t.Error("progress should not be nil")
	} else {
		if !floatEqual(*expected.Progress, 100.0, 2) {
			t.Errorf("expected progress 100.0, got %v", *expected.Progress)
		}
	}

	if expected.Status == nil {
		t.Error("status should not be nil")
	} else {
		if *expected.Status != "finished" {
			t.Errorf("expected status 'finished', got '%s'", *expected.Status)
		}
	}
}

// Helper function for float comparison
func floatEqual(a, b float64, precision int) bool {
	multiplier := math.Pow(10, float64(precision))
	return math.Round(a*multiplier) == math.Round(b*multiplier)
}
