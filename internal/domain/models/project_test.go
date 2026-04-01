package models

import (
	"context"
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
		MedianDay:  &now,
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
