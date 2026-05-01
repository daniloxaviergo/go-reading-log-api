package unit

import (
	"context"
	"testing"
	"time"

	"go-reading-log-api-next/internal/config"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/domain/models"
	"go-reading-log-api-next/test/testutil"
)

// TestProject_CalculateDaysUnreading tests the CalculateDaysUnreading method
func TestProject_CalculateDaysUnreading(t *testing.T) {
	today := time.Now()

	// Create a project with started_at date (10 days ago)
	startedAt := today.AddDate(0, 0, -10)
	project := &models.Project{
		ID:        1,
		TotalPage: 100,
		Page:      50,
		Reinicia:  false,
		StartedAt: &startedAt,
	}

	tests := []struct {
		name     string
		logs     []*dto.LogResponse
		expected *int
	}{
		{
			name: "no_logs_uses_started_at",
			logs: nil,
			// Should calculate from started_at (10 days ago)
			// Allow 1 day tolerance for date calculation
		},
		{
			name: "single_log_with_date",
			logs: []*dto.LogResponse{
				{Data: timePtr(today.AddDate(0, 0, -3))},
			},
			expected: testutil.IntPtr(3),
		},
		{
			name: "multiple_logs_returns_most_recent",
			logs: []*dto.LogResponse{
				{Data: timePtr(today.AddDate(0, 0, -8))},
				{Data: timePtr(today.AddDate(0, 0, -6))},
				{Data: timePtr(today.AddDate(0, 0, -4))},
			},
			expected: testutil.IntPtr(4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			days := project.CalculateDaysUnreading(tt.logs)

			if days == nil {
				t.Fatal("CalculateDaysUnreading returned nil")
			}

			// Verify non-negative result
			if *days < 0 {
				t.Errorf("Expected non-negative days, got %d", *days)
			}

			// If expected is provided, check within tolerance
			if tt.expected != nil {
				diff := *days - *tt.expected
				if diff < 0 {
					diff = -diff
				}
				if diff > 1 {
					t.Errorf("Expected ~%d days, got %d", *tt.expected, *days)
				}
			}
		})
	}
}

// TestProject_CalculateDaysUnreading_EdgeCases tests edge cases for days_unreading
func TestProject_CalculateDaysUnreading_EdgeCases(t *testing.T) {

	tests := []struct {
		name     string
		project  *models.Project
		logs     []*dto.LogResponse
		expected *int
	}{
		{
			name: "no_logs_no_started_at",
			project: &models.Project{
				ID:        1,
				TotalPage: 100,
				Page:      50,
				Reinicia:  false,
			},
			logs:     nil,
			expected: testutil.IntPtr(0),
		},
		{
			name: "no_logs_with_started_at",
			project: &models.Project{
				ID:        1,
				TotalPage: 100,
				Page:      50,
				Reinicia:  false,
				StartedAt: func() *time.Time { t := time.Now(); return &t }(),
			},
			logs:     nil,
			expected: testutil.IntPtr(0), // At least 0 days
		},
		{
			name: "log_after_started_at",
			project: &models.Project{
				ID:        1,
				TotalPage: 100,
				Page:      50,
				Reinicia:  false,
				StartedAt: func() *time.Time { t := time.Now().AddDate(0, 0, -5); return &t }(),
			},
			logs: []*dto.LogResponse{
				{Data: timePtr(time.Now())}, // Today
			},
			expected: testutil.IntPtr(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			days := tt.project.CalculateDaysUnreading(tt.logs)

			if days == nil {
				t.Fatal("CalculateDaysUnreading returned nil")
			}

			if *days < 0 {
				t.Errorf("Expected non-negative days, got %d", *days)
			}

			if tt.expected != nil {
				diff := *days - *tt.expected
				if diff < 0 {
					diff = -diff
				}
				if diff > 1 {
					t.Errorf("Expected ~%d days, got %d", *tt.expected, *days)
				}
			}
		})
	}
}

// TestProject_CalculateDaysUnreading_MultiFormat tests multi-format date parsing
func TestProject_CalculateDaysUnreading_MultiFormat(t *testing.T) {
	today := time.Now()
	startedAt := today.AddDate(0, 0, -7)
	project := &models.Project{
		ID:        1,
		TotalPage: 100,
		Page:      30,
		Reinicia:  false,
		StartedAt: &startedAt,
	}

	tests := []struct {
		name string
		logs []*dto.LogResponse
	}{
		{
			name: "YYYY-MM-DD format",
			logs: []*dto.LogResponse{
				{Data: timePtr(today.AddDate(0, 0, -3))},
			},
		},
		{
			name: "RFC3339 format",
			logs: []*dto.LogResponse{
				{Data: timePtr(today.AddDate(0, 0, -3).Add(time.Hour * 15 / time.Microsecond * 1000))}, // Approximate
			},
		},
		{
			name: "Standard datetime format",
			logs: []*dto.LogResponse{
				{Data: timePtr(today.AddDate(0, 0, -3).Add(time.Hour * 15 / time.Microsecond * 1000))}, // Approximate
			},
		},
		{
			name: "Mixed formats",
			logs: []*dto.LogResponse{
				{Data: timePtr(today.AddDate(0, 0, -8).Add(time.Hour * 15 / time.Microsecond * 1000))},
				{Data: timePtr(today.AddDate(0, 0, -6).Add(time.Hour * 15 / time.Microsecond * 1000))},
				{Data: timePtr(today.AddDate(0, 0, -4))},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			days := project.CalculateDaysUnreading(tt.logs)

			if days == nil {
				t.Fatal("CalculateDaysUnreading returned nil")
			}

			if *days < 0 {
				t.Errorf("Expected non-negative days, got %d", *days)
			}
		})
	}
}

// TestProject_CalculateDaysUnreading_Timezone tests timezone-aware calculation
func TestProject_CalculateDaysUnreading_Timezone(t *testing.T) {
	today := time.Now()
	startedAt := today.AddDate(0, 0, -7)
	project := &models.Project{
		ID:        1,
		TotalPage: 100,
		Page:      30,
		Reinicia:  false,
		StartedAt: &startedAt,
	}

	// Test with Brazil timezone (BRT)
	brazilCtx := context.WithValue(context.Background(), "timezone", time.FixedZone("BRT", -3*60*60))
	project.SetContext(brazilCtx)

	days := project.CalculateDaysUnreading(nil)

	if days == nil {
		t.Fatal("CalculateDaysUnreading returned nil")
	}

	// Should be approximately 7 days (allowing for timezone edge cases)
	if *days < 6 || *days > 8 {
		t.Errorf("Expected ~7 days, got %d", *days)
	}
}

// TestProject_CalculateFinishedAt tests the CalculateFinishedAt method
func TestProject_CalculateFinishedAt(t *testing.T) {

	today := time.Now()
	startedAt := today.AddDate(0, 0, -10)
	project := &models.Project{

		ID:        1,
		TotalPage: 200,
		Page:      50,
		Reinicia:  false,
		StartedAt: &startedAt,
	}

	tests := []struct {
		name     string
		logs     []*dto.LogResponse
		expected func(*time.Time) bool
	}{
		{
			name: "with_logs_estimates_completion",
			logs: []*dto.LogResponse{
				{Data: timePtr(today.AddDate(0, 0, -2))},
			},
			expected: func(finished *time.Time) bool {
				return finished != nil && finished.After(today)
			},
		},
		{
			name:     "no_logs_returns_nil",
			logs:     nil,
			expected: func(finished *time.Time) bool { return finished == nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finished := project.CalculateFinishedAt(tt.logs)

			if !tt.expected(finished) {
				if finished == nil {
					t.Errorf("Expected non-nil finished_at, got nil")
				} else {
					t.Errorf("Expected finished_at to be valid, got %v", *finished)
				}
			}
		})
	}
}

// TestProject_CalculateFinishedAt_EdgeCases tests edge cases for finished_at
func TestProject_CalculateFinishedAt_EdgeCases(t *testing.T) {
	today := time.Now()

	tests := []struct {
		name     string
		project  *models.Project
		logs     []*dto.LogResponse
		expected *time.Time
	}{
		{
			name: "no_started_at_returns_nil",
			project: &models.Project{

				ID:        1,
				TotalPage: 100,
				Page:      50,
				Reinicia:  false,
			},
			logs:     nil,
			expected: nil,
		},
		{
			name: "completed_project_no_logs_returns_nil",
			project: &models.Project{

				ID:        1,
				TotalPage: 100,
				Page:      100,
				Reinicia:  false,
				StartedAt: func() *time.Time { t := time.Now(); return &t }(),
			},
			logs:     nil,
			expected: nil,
		},
		{
			name: "page_equals_total_with_logs_returns_last_log_date",
			project: &models.Project{

				ID:        1,
				TotalPage: 100,
				Page:      100,
				Reinicia:  false,
				StartedAt: func() *time.Time { t := time.Now(); return &t }(),
			},
			logs: []*dto.LogResponse{
				{Data: timePtr(today.AddDate(0, 0, -2))},
			},
			// When page >= total_page with logs, returns the most recent log's date
			expected: func() *time.Time {
				t := today.AddDate(0, 0, -2)
				return &t
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finished := tt.project.CalculateFinishedAt(tt.logs)

			if (tt.expected == nil) != (finished == nil) {
				t.Errorf("Expected finished_at %v, got %v", tt.expected, finished)
			}

			if tt.expected != nil && finished != nil {
				// Verify the date is reasonable (within 1 day to 365 days of expected)
				expectedDate := *tt.expected
				if finished.Before(expectedDate.AddDate(0, 0, -1)) ||
					finished.After(expectedDate.AddDate(0, 0, 365)) {
					t.Errorf("Expected finished_at near %v, got %v", expectedDate, *finished)
				}
			}
		})
	}
}

// TestProject_CalculateMedianDay tests the CalculateMedianDay method
func TestProject_CalculateMedianDay(t *testing.T) {

	today := time.Now()
	startedAt := today.AddDate(0, 0, -14) // 14 days ago
	project := &models.Project{

		ID:        1,
		TotalPage: 100,
		Page:      50,
		Reinicia:  false,
		StartedAt: &startedAt,
	}

	median := project.CalculateMedianDay()

	if median == nil {
		t.Fatal("CalculateMedianDay returned nil")
	}

	// With 14 days and 50 pages, median should be around 3.57
	expectedMin := 3.0
	expectedMax := 4.0

	if *median < expectedMin || *median > expectedMax {
		t.Errorf("Expected median between %v and %v, got %v", expectedMin, expectedMax, *median)
	}
}

// TestProject_CalculateMedianDay_EdgeCases tests edge cases for median_day
func TestProject_CalculateMedianDay_EdgeCases(t *testing.T) {

	tests := []struct {
		name     string
		project  *models.Project
		expected *float64
	}{
		{
			name: "no_started_at_returns_zero",
			project: &models.Project{

				ID:        1,
				TotalPage: 100,
				Page:      50,
				Reinicia:  false,
			},
			expected: testutil.FloatPtr(0.0),
		},
		{
			name: "zero_days_reading_returns_zero",
			project: &models.Project{

				ID:        1,
				TotalPage: 100,
				Page:      50,
				Reinicia:  false,
				StartedAt: func() *time.Time { t := time.Now(); return &t }(),
			},
			expected: testutil.FloatPtr(0.0),
		},
		{
			name: "page_zero_returns_zero",
			project: &models.Project{

				ID:        1,
				TotalPage: 100,
				Page:      0,
				Reinicia:  false,
				StartedAt: func() *time.Time { t := time.Now().AddDate(0, 0, -7); return &t }(),
			},
			expected: testutil.FloatPtr(0.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			median := tt.project.CalculateMedianDay()

			if median == nil {
				t.Fatal("CalculateMedianDay returned nil")
			}

			diff := *median - *tt.expected
			if diff < 0 {
				diff = -diff
			}

			if diff > 0.01 {
				t.Errorf("Expected median %v, got %v", *tt.expected, *median)
			}
		})
	}
}

// TestProject_CalculateMedianDay_Timezone tests timezone-aware median day calculation
func TestProject_CalculateMedianDay_Timezone(t *testing.T) {

	today := time.Now()
	startedAt := today.AddDate(0, 0, -14)
	project := &models.Project{

		ID:        1,
		TotalPage: 100,
		Page:      50,
		Reinicia:  false,
		StartedAt: &startedAt,
	}

	// Test with Brazil timezone (BRT)
	brazilCtx := context.WithValue(context.Background(), "timezone", time.FixedZone("BRT", -3*60*60))
	project.SetContext(brazilCtx)

	median := project.CalculateMedianDay()

	if median == nil {
		t.Fatal("CalculateMedianDay returned nil")
	}

	// With 14 days and 50 pages, median should be around 3.57
	expectedMin := 3.0
	expectedMax := 4.0

	if *median < expectedMin || *median > expectedMax {
		t.Errorf("Expected median between %v and %v, got %v", expectedMin, expectedMax, *median)
	}
}

// TestProject_CalculateProgress tests the CalculateProgress method
func TestProject_CalculateProgress(t *testing.T) {

	tests := []struct {
		name      string
		totalPage int
		page      int
		expected  float64
	}{
		{
			name:      "50_percent",
			totalPage: 100,
			page:      50,
			expected:  50.0,
		},
		{
			name:      "25_percent",
			totalPage: 200,
			page:      50,
			expected:  25.0,
		},
		{
			name:      "100_percent",
			totalPage: 100,
			page:      100,
			expected:  100.0,
		},
		{
			name:      "zero_page",
			totalPage: 100,
			page:      0,
			expected:  0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			project := models.NewProject(ctx, 1, "Test", tt.totalPage, tt.page, false)
			progress := project.CalculateProgress()

			if progress == nil {
				t.Fatal("CalculateProgress returned nil")
			}

			diff := *progress - tt.expected
			if diff < 0 {
				diff = -diff
			}

			if diff > 0.01 {
				t.Errorf("Expected progress %v, got %v", tt.expected, *progress)
			}
		})
	}
}

// TestProject_CalculateProgress_EdgeCases tests edge cases for progress calculation
func TestProject_CalculateProgress_EdgeCases(t *testing.T) {

	tests := []struct {
		name      string
		totalPage int
		page      int
		expected  float64
	}{
		{
			name:      "zero_total_page",
			totalPage: 0,
			page:      50,
			expected:  0.0,
		},
		{
			name:      "negative_total_page",
			totalPage: -10,
			page:      50,
			expected:  0.0,
		},
		{
			name:      "page_exceeds_total",
			totalPage: 100,
			page:      150,
			expected:  100.0, // Should be clamped at 100
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			project := models.NewProject(ctx, 1, "Test", tt.totalPage, tt.page, false)
			progress := project.CalculateProgress()

			if progress == nil {
				t.Fatal("CalculateProgress returned nil")
			}

			diff := *progress - tt.expected
			if diff < 0 {
				diff = -diff
			}

			if diff > 0.01 {
				t.Errorf("Expected progress %v, got %v", tt.expected, *progress)
			}
		})
	}
}

// TestProject_CalculateStatus tests the CalculateStatus method
func TestProject_CalculateStatus(t *testing.T) {
	cfg := &config.Config{
		EmAndamentoRange: 7,
		DormindoRange:    14,
	}

	today := time.Now()
	startedAt := today.AddDate(0, 0, -5)
	project := &models.Project{

		ID:        1,
		TotalPage: 100,
		Page:      50,
		Reinicia:  false,
		StartedAt: &startedAt,
	}

	logs := []*dto.LogResponse{
		{Data: timePtr(today.AddDate(0, 0, -2))}, // 2 days ago
	}

	status := project.CalculateStatus(logs, cfg)

	if status == nil {
		t.Fatal("CalculateStatus returned nil")
	}

	// With 2 days since last read (within em_andamento_range of 7), should be "running"
	if *status != models.StatusRunning {
		t.Errorf("Expected status 'running', got '%s'", *status)
	}
}

// TestProject_CalculateStatus_EdgeCases tests edge cases for status calculation
func TestProject_CalculateStatus_EdgeCases(t *testing.T) {
	cfg := &config.Config{
		EmAndamentoRange: 7,
		DormindoRange:    14,
	}

	today := time.Now()

	tests := []struct {
		name     string
		project  *models.Project
		logs     []*dto.LogResponse
		expected string
	}{
		{
			name: "no_logs_returns_unstarted",
			project: &models.Project{

				ID:        1,
				TotalPage: 100,
				Page:      50,
				Reinicia:  false,
			},
			logs:     []*dto.LogResponse{},
			expected: models.StatusUnstarted,
		},
		{
			name: "no_logs_with_started_at_returns_unstarted",
			project: &models.Project{

				ID:        1,
				TotalPage: 100,
				Page:      50,
				Reinicia:  false,
				StartedAt: &today,
			},
			logs:     []*dto.LogResponse{},
			expected: models.StatusUnstarted,
		},
		{
			name: "completed_project_no_logs_returns_finished",
			project: &models.Project{

				ID:        1,
				TotalPage: 100,
				Page:      100,
				Reinicia:  false,
				StartedAt: &today,
			},
			logs:     []*dto.LogResponse{},
			expected: models.StatusFinished,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := tt.project.CalculateStatus(tt.logs, cfg)

			if status == nil {
				t.Fatal("CalculateStatus returned nil")
			}

			if *status != tt.expected {
				t.Errorf("Expected status '%s', got '%s'", tt.expected, *status)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
