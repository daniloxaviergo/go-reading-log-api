package models

import (
	"context"
	"testing"
)

// TestLog tests the Log model
func TestLog(t *testing.T) {
	ctx := context.Background()

	// Test NewLog
	log := NewLog(ctx, 1, 100, 10, 20)

	if log == nil {
		t.Fatal("Expected non-nil log, got nil")
	}

	if log.ID != 1 {
		t.Errorf("Expected ID 1, got %d", log.ID)
	}

	if log.ProjectID != 100 {
		t.Errorf("Expected project_id 100, got %d", log.ProjectID)
	}

	if log.StartPage != 10 {
		t.Errorf("Expected start_page 10, got %d", log.StartPage)
	}

	if log.EndPage != 20 {
		t.Errorf("Expected end_page 20, got %d", log.EndPage)
	}

	// Test GetContext
	ctxFromLog := log.GetContext()
	if ctxFromLog == nil {
		t.Fatal("Expected non-nil context")
	}

	// Test SetContext
	newCtx := context.Background()
	log.SetContext(newCtx)

	if log.GetContext() != newCtx {
		t.Error("Context was not set correctly")
	}
}

// TestLog_WithOptionalFields tests the log with optional fields
func TestLog_WithOptionalFields(t *testing.T) {
	ctx := context.Background()
	data := "2024-01-01"
	note := "This is a note"

	log := &Log{
		ctx:       ctx,
		ID:        1,
		ProjectID: 100,
		Data:      &data,
		StartPage: 10,
		EndPage:   20,
		Wday:      1,
		Note:      &note,
	}

	if *log.Data != "2024-01-01" {
		t.Errorf("Expected data '2024-01-01', got '%s'", *log.Data)
	}

	if *log.Note != "This is a note" {
		t.Errorf("Expected note 'This is a note', got '%s'", *log.Note)
	}

	if log.Wday != 1 {
		t.Errorf("Expected wday 1, got %d", log.Wday)
	}
}

// TestLog_EmptyContext tests context fallback
func TestLog_EmptyContext(t *testing.T) {
	log := &Log{}
	ctx := log.GetContext()
	if ctx == nil {
		t.Error("Expected default context, got nil")
	}
}

// TestLog_Wday tests the wday field
func TestLog_Wday(t *testing.T) {
	ctx := context.Background()

	log := NewLog(ctx, 1, 100, 10, 20)
	log.Wday = 5 // Friday

	if log.Wday != 5 {
		t.Errorf("Expected wday 5, got %d", log.Wday)
	}
}
