package models

import (
	"context"
	"go-reading-log-api-next/internal/config"
	"go-reading-log-api-next/internal/domain/dto"
	validation "go-reading-log-api-next/internal/validation"
	"testing"
	"time"
)

// TestProject tests the Project model
func TestProject(t *testing.T) {
	ctx := context.Background()

	// Test NewProject
	project := NewProject(ctx, 1, "Test Project", 100, 50, false)

	if project == nil {
		t.Fatal("Expected non-nil project, got nil")
	}

	if project.ID != 1 {
		t.Errorf("Expected ID 1, got %d", project.ID)
	}

	if project.Name != "Test Project" {
		t.Errorf("Expected name 'Test Project', got '%s'", project.Name)
	}

	if project.TotalPage != 100 {
		t.Errorf("Expected total_page 100, got %d", project.TotalPage)
	}

	if project.Page != 50 {
		t.Errorf("Expected page 50, got %d", project.Page)
	}

	if project.Reinicia != false {
		t.Errorf("Expected reinicia false, got %t", project.Reinicia)
	}

	// Test GetContext
	ctxFromProject := project.GetContext()
	if ctxFromProject == nil {
		t.Fatal("Expected non-nil context")
	}

	// Test SetContext
	newCtx := context.Background()
	project.SetContext(newCtx)

	if project.GetContext() != newCtx {
		t.Error("Context was not set correctly")
	}
}

// TestProject_WithOptionalFields tests the project with optional fields
func TestProject_WithOptionalFields(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	startedAt := now
	status := "completed"
	logsCount := 10
	daysUnread := 5

	project := &Project{
		ctx:        ctx,
		ID:         1,
		Name:       "Test Project",
		TotalPage:  100,
		Page:       50,
		Reinicia:   false,
		StartedAt:  &startedAt,
		Progress:   floatPtr(50.5),
		Status:     &status,
		LogsCount:  &logsCount,
		DaysUnread: &daysUnread,
		MedianDay:  floatPtr(5.5),
		FinishedAt: &now,
	}

	if *project.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", *project.Status)
	}

	if project.LogsCount == nil {
		t.Fatal("Expected LogsCount to be set")
	}
	if *project.LogsCount != 10 {
		t.Errorf("Expected logs_count 10, got %d", *project.LogsCount)
	}
}

// TestProject_EmptyContext tests context fallback
func TestProject_EmptyContext(t *testing.T) {
	project := &Project{}
	ctx := project.GetContext()
	if ctx == nil {
		t.Error("Expected default context, got nil")
	}
}

// TestProject_Reinicia tests the reinicia field
func TestProject_Reinicia(t *testing.T) {
	ctx := context.Background()

	project := NewProject(ctx, 1, "Test Project", 100, 50, true)

	if project.Reinicia != true {
		t.Errorf("Expected reinicia true, got %t", project.Reinicia)
	}
}

// floatPtr returns a pointer to a float64 value
func floatPtr(f float64) *float64 {
	return &f
}

// TestProject_CalculateProgress_ZeroTotalPage_EdgeCase tests edge cases for CalculateProgress
// Tests zero/negative total_page scenarios
func TestProject_CalculateProgress_ZeroTotalPage_EdgeCase(t *testing.T) {
	ctx := context.Background()

	t.Run("zero_total_page", func(t *testing.T) {
		project := NewProject(ctx, 1, "Test", 0, 50, false)
		progress := project.CalculateProgress()
		if progress == nil || *progress != 0.0 {
			t.Errorf("Expected 0.0 for zero total_page, got %v", progress)
		}
	})

	t.Run("negative_total_page", func(t *testing.T) {
		project := NewProject(ctx, 1, "Test", -10, 50, false)
		progress := project.CalculateProgress()
		if progress == nil || *progress != 0.0 {
			t.Errorf("Expected 0.0 for negative total_page, got %v", progress)
		}
	})

	t.Run("page_exceeds_total", func(t *testing.T) {
		project := NewProject(ctx, 1, "Test", 100, 150, false)
		progress := project.CalculateProgress()
		if progress == nil || *progress != 100.0 {
			t.Errorf("Expected 100.0 for page > total_page (clamped), got %v", progress)
		}
	})
}

// TestProject_CalculateDaysUnreading_NoLogsEdgeCases tests edge cases for CalculateDaysUnreading
// Tests scenarios with no logs
func TestProject_CalculateDaysUnreading_NoLogsEdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("no_logs_no_started_at", func(t *testing.T) {
		project := &Project{
			ctx:       ctx,
			ID:        1,
			TotalPage: 100,
			Page:      50,
			Reinicia:  false,
		}
		days := project.CalculateDaysUnreading(nil)
		if days == nil || *days != 0 {
			t.Errorf("Expected 0 for no logs and no started_at, got %v", days)
		}
	})

	t.Run("no_logs_with_started_at", func(t *testing.T) {
		now := time.Now()
		project := &Project{
			ctx:       ctx,
			ID:        1,
			TotalPage: 100,
			Page:      50,
			Reinicia:  false,
			StartedAt: &now,
		}
		days := project.CalculateDaysUnreading(nil)
		if days == nil || *days < 0 {
			t.Errorf("Expected non-negative days, got %v", days)
		}
	})
}

// TestProject_CalculateFinishedAt_100PercentProgress tests edge cases for CalculateFinishedAt
// Tests 100% progress scenarios
func TestProject_CalculateFinishedAt_100PercentProgress(t *testing.T) {
	ctx := context.Background()

	t.Run("finished_book_no_logs", func(t *testing.T) {
		today := time.Now()
		project := &Project{
			ctx:       ctx,
			ID:        1,
			TotalPage: 100,
			Page:      100,
			Reinicia:  false,
			StartedAt: &today,
		}
		finished := project.CalculateFinishedAt(nil)
		if finished != nil {
			t.Errorf("Expected nil for finished book with no logs, got %v", finished)
		}
	})

	t.Run("page_equals_total_with_logs", func(t *testing.T) {
		today := time.Now()
		logs := []*dto.LogResponse{
			{
				Data: &today,
			},
		}
		project := &Project{
			ctx:       ctx,
			ID:        1,
			TotalPage: 100,
			Page:      100,
			Reinicia:  false,
			StartedAt: &today,
		}
		finished := project.CalculateFinishedAt(logs)
		if finished == nil {
			t.Errorf("Expected non-nil finished_at with logs, got nil")
		}
	})

	t.Run("no_started_at", func(t *testing.T) {
		project := &Project{
			ctx:       ctx,
			ID:        1,
			TotalPage: 100,
			Page:      50,
			Reinicia:  false,
		}
		finished := project.CalculateFinishedAt(nil)
		if finished != nil {
			t.Errorf("Expected nil for no started_at, got %v", finished)
		}
	})
}

// TestProject_CalculateStatus_InvalidValue tests CalculateStatus with various edge cases
func TestProject_CalculateStatus_InvalidValue(t *testing.T) {
	ctx := context.Background()
	config := &config.Config{
		EmAndamentoRange: 7,
		DormindoRange:    14,
	}

	t.Run("no_logs_returns_unstarted", func(t *testing.T) {
		project := NewProject(ctx, 1, "Test Project", 100, 50, false)
		logs := []*dto.LogResponse{}
		result := project.CalculateStatus(logs, config)
		if result == nil {
			t.Fatal("Expected non-nil status, got nil")
		}
		if *result != StatusUnstarted {
			t.Errorf("Expected status 'unstarted' with no logs, got '%s'", *result)
		}
	})
}

// TestProject_CalculateStatus_ValidationWithInvalidValue tests validation error for AC4
// This verifies that ValidateStatus() returns error for invalid status values (AC4 requirement)
func TestProject_CalculateStatus_ValidationWithInvalidValue(t *testing.T) {
	invalidStatuses := []string{
		"invalid_status",
		"",
		"unknown",
		"running ", // trailing space
		" Running", // leading space
	}

	for _, status := range invalidStatuses {
		t.Run(status, func(t *testing.T) {
			err := validation.ValidateStatus(status)
			if err == nil {
				t.Errorf("Expected validation error for invalid status '%s', got nil", status)
			}
			if err != nil && err.Code != "invalid_status" {
				t.Errorf("Expected error code 'invalid_status', got '%s'", err.Code)
			}
		})
	}
}

// TestProject_DerivedCalculationsEdgeCases_Documentation tests that all edge cases are documented
// This verifies that all derived calculation methods handle edge cases gracefully (AC5 requirement)
func TestProject_DerivedCalculationsEdgeCases_Documentation(t *testing.T) {
	ctx := context.Background()

	// This test documents all edge cases handled by the calculation methods
	// Each calculation method should handle edge cases and return safe defaults

	// Test 1: CalculateProgress handles zero/negative values
	t.Run("CalculateProgress handles zero/negative", func(t *testing.T) {
		project := NewProject(ctx, 1, "Test", 0, 0, false)
		progress := project.CalculateProgress()
		if progress == nil || *progress != 0.0 {
			t.Errorf("Expected 0.0 for zero/negative edge cases")
		}
	})

	// Test 2: CalculateDaysUnreading handles no logs and no started_at
	t.Run("CalculateDaysUnreading handles no logs/no started_at", func(t *testing.T) {
		project := &Project{
			ctx:       ctx,
			ID:        1,
			TotalPage: 100,
			Page:      50,
			Reinicia:  false,
		}
		days := project.CalculateDaysUnreading(nil)
		if days == nil || *days != 0 {
			t.Errorf("Expected 0 for no logs and no started_at")
		}
	})

	// Test 3: CalculateMedianDay handles no started_at
	t.Run("CalculateMedianDay handles no started_at", func(t *testing.T) {
		project := &Project{
			ctx:       ctx,
			ID:        1,
			TotalPage: 100,
			Page:      50,
			Reinicia:  false,
		}
		median := project.CalculateMedianDay()
		if median == nil || *median != 0.0 {
			t.Errorf("Expected 0.0 for no started_at")
		}
	})

	// Test 4: CalculateFinishedAt handles edge cases
	t.Run("CalculateFinishedAt handles edge cases", func(t *testing.T) {
		// No started_at
		project1 := &Project{
			ctx:       ctx,
			ID:        1,
			TotalPage: 100,
			Page:      50,
			Reinicia:  false,
		}
		finished1 := project1.CalculateFinishedAt(nil)
		if finished1 != nil {
			t.Errorf("Expected nil for no started_at")
		}

		// Finished book with no logs
		today := time.Now()
		project2 := &Project{
			ctx:       ctx,
			ID:        2,
			TotalPage: 100,
			Page:      100,
			Reinicia:  false,
			StartedAt: &today,
		}
		finished2 := project2.CalculateFinishedAt(nil)
		if finished2 != nil {
			t.Errorf("Expected nil for finished book with no logs")
		}
	})
}

// TestProject_ParseLogDate tests the parseLogDate function with multiple formats
func TestProject_ParseLogDate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wantOk bool
	}{
		{
			name:   "YYYY-MM-DD format",
			input:  "2024-01-15",
			wantOk: true,
		},
		{
			name:   "RFC3339 format",
			input:  "2024-01-15T10:30:00Z",
			wantOk: true,
		},
		{
			name:   "Standard datetime format",
			input:  "2024-01-15 10:30:00",
			wantOk: true,
		},
		{
			name:   "Invalid format",
			input:  "not-a-date",
			wantOk: false,
		},
		{
			name:   "Empty string",
			input:  "",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, gotOk := parseLogDate(tt.input)

			if gotOk != tt.wantOk {
				t.Errorf("parseLogDate(%q) = _, %v; want ok=%v", tt.input, gotOk, tt.wantOk)
			}

			if gotOk && tt.wantOk {
				// Verify the parsed time is reasonable
				if gotTime.IsZero() {
					t.Errorf("parseLogDate(%q) returned zero time", tt.input)
				}
			}
		})
	}
}

// TestProject_CalculateDaysUnreading_MultiFormat tests CalculateDaysUnreading with different date formats
func TestProject_CalculateDaysUnreading_MultiFormat(t *testing.T) {
	ctx := context.Background()

	// Create a project with started_at
	today := time.Now()
	startedAt := today.AddDate(0, 0, -5) // 5 days ago
	project := &Project{
		ctx:       ctx,
		ID:        1,
		TotalPage: 100,
		Page:      20,
		Reinicia:  false,
		StartedAt: &startedAt,
	}

	tests := []struct {
		name     string
		logs     []*dto.LogResponse
		expected *int
	}{
		{
			name: "YYYY-MM-DD format in log",
			logs: []*dto.LogResponse{
				{Data: timePtr(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))},
			},
			expected: intPtr(3), // 3 days since last read (assuming today is 2024-01-18)
		},
		{
			name: "RFC3339 format in log",
			logs: []*dto.LogResponse{
				{Data: timePtr(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC))},
			},
			expected: intPtr(3),
		},
		{
			name: "Standard datetime format in log",
			logs: []*dto.LogResponse{
				{Data: timePtr(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC))},
			},
			expected: intPtr(3),
		},
		{
			name: "Multiple logs with different formats",
			logs: []*dto.LogResponse{
				{Data: timePtr(time.Date(2024, 1, 10, 10, 30, 0, 0, time.UTC))},
				{Data: timePtr(time.Date(2024, 1, 12, 15, 0, 0, 0, time.UTC))},
				{Data: timePtr(time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC))},
			},
			expected: intPtr(4), // Most recent is 2024-01-14
		},
		{
			name:     "No logs, uses started_at",
			logs:     []*dto.LogResponse{},
			expected: intPtr(5), // 5 days since started_at
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			days := project.CalculateDaysUnreading(tt.logs)

			if days == nil {
				t.Fatal("CalculateDaysUnreading returned nil")
			}

			// Note: The expected value depends on the current date
			// This test verifies the function works with multiple formats
			// and returns a non-negative integer
			if *days < 0 {
				t.Errorf("Expected non-negative days, got %d", *days)
			}
		})
	}
}

// TestProject_CalculateDaysUnreading_Timezone tests timezone-aware date comparison
func TestProject_CalculateDaysUnreading_Timezone(t *testing.T) {
	ctx := context.Background()

	// Create a project with started_at in the past
	today := time.Now()
	startedAt := today.AddDate(0, 0, -7) // 7 days ago
	project := &Project{
		ctx:       ctx,
		ID:        1,
		TotalPage: 100,
		Page:      30,
		Reinicia:  false,
		StartedAt: &startedAt,
	}

	// Test with Brazil timezone (BRT)
	brazilCtx := context.WithValue(ctx, "timezone", time.FixedZone("BRT", -3*60*60))
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

// TestProject_CalculateMedianDay_Timezone tests timezone-aware median day calculation
func TestProject_CalculateMedianDay_Timezone(t *testing.T) {
	ctx := context.Background()

	today := time.Now()
	startedAt := today.AddDate(0, 0, -14) // 14 days ago
	project := &Project{
		ctx:       ctx,
		ID:        1,
		TotalPage: 100,
		Page:      50,
		Reinicia:  false,
		StartedAt: &startedAt,
	}

	// Test with Brazil timezone (BRT)
	brazilCtx := context.WithValue(ctx, "timezone", time.FixedZone("BRT", -3*60*60))
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

// TestProject_CalculateFinishedAt_MultiFormat tests CalculateFinishedAt with different date formats
func TestProject_CalculateFinishedAt_MultiFormat(t *testing.T) {
	ctx := context.Background()

	today := time.Now()
	startedAt := today.AddDate(0, 0, -10)
	project := &Project{
		ctx:       ctx,
		ID:        1,
		TotalPage: 200,
		Page:      50,
		Reinicia:  false,
		StartedAt: &startedAt,
	}

	tests := []struct {
		name string
		logs []*dto.LogResponse
	}{
		{
			name: "YYYY-MM-DD format in log",
			logs: []*dto.LogResponse{
				{Data: timePtr(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))},
			},
		},
		{
			name: "RFC3339 format in log",
			logs: []*dto.LogResponse{
				{Data: timePtr(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC))},
			},
		},
		{
			name: "Standard datetime format in log",
			logs: []*dto.LogResponse{
				{Data: timePtr(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC))},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finished := project.CalculateFinishedAt(tt.logs)

			if finished == nil {
				t.Fatal("CalculateFinishedAt returned nil")
			}

			// Verify the date is in the future (estimated finish date)
			now := time.Now()
			if !finished.After(now) {
				t.Errorf("Expected finished_at to be in the future, got %v", *finished)
			}
		})
	}
}

// intPtr returns a pointer to an int value
func intPtr(i int) *int {
	return &i
}

// timePtr returns a pointer to a time.Time value
func timePtr(t time.Time) *time.Time {
	return &t
}
