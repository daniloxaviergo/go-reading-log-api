package dto

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

// TestLogResponse tests the LogResponse DTO
func TestLogResponse(t *testing.T) {
	data := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	response := NewLogResponse(1, &data, 10, 20)

	if response == nil {
		t.Fatal("Expected non-nil response, got nil")
	}

	if response.ID != 1 {
		t.Errorf("Expected ID 1, got %d", response.ID)
	}

	if *response.Data != data {
		t.Errorf("Expected data to match input time")
	}

	if response.StartPage != 10 {
		t.Errorf("Expected start_page 10, got %d", response.StartPage)
	}

	if response.EndPage != 20 {
		t.Errorf("Expected end_page 20, got %d", response.EndPage)
	}

	// Test GetContext
	ctx := response.GetContext()
	if ctx == nil {
		t.Fatal("Expected non-nil context")
	}

	// Test SetContext
	newCtx := context.Background()
	response.SetContext(newCtx)

	if response.GetContext() != newCtx {
		t.Error("Context was not set correctly")
	}
}

// TestLogResponse_JSON tests JSON serialization
func TestLogResponse_JSON(t *testing.T) {
	data := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	response := NewLogResponse(1, &data, 10, 20)
	response.SetContext(context.Background())

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verify RFC3339 format in JSON
	expectedDataStr := data.Format(time.RFC3339)
	if !contains(string(jsonData), expectedDataStr) {
		t.Errorf("Expected RFC3339 data '%s' in JSON: %s", expectedDataStr, string(jsonData))
	}

	// Test JSON unmarshaling
	var decoded LogResponse
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if decoded.ID != 1 {
		t.Errorf("Expected ID 1, got %d", decoded.ID)
	}

	if decoded.Data == nil {
		t.Error("Expected non-nil Data after unmarshaling")
	} else if decoded.Data.Format(time.RFC3339) != expectedDataStr {
		t.Errorf("Expected data '%s', got '%s'", expectedDataStr, decoded.Data.Format(time.RFC3339))
	}
}

// TestLogResponse_WithNote tests the response with optional note
func TestLogResponse_WithNote(t *testing.T) {
	data := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	note := "This is a note"
	response := &LogResponse{
		ID:        1,
		Data:      &data,
		StartPage: 10,
		EndPage:   20,
		Note:      &note,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verify note is included
	if !contains(string(jsonData), "This is a note") {
		t.Errorf("Expected note in JSON: %s", string(jsonData))
	}
}

// TestLogResponse_EmptyContext tests context fallback
func TestLogResponse_EmptyContext(t *testing.T) {
	response := &LogResponse{}
	ctx := response.GetContext()
	if ctx == nil {
		t.Error("Expected default context, got nil")
	}
}

// TestLogResponse_AtributesStructure tests that attributes contain expected fields only
func TestLogResponse_AtributesStructure(t *testing.T) {
	data := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	response := &LogResponse{
		ID:        1,
		Data:      &data,
		StartPage: 10,
		EndPage:   20,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verify attributes contain expected fields
	expectedFields := []string{"id", "data", "start_page", "end_page"}
	for _, field := range expectedFields {
		if !contains(string(jsonData), `"`+field+`"`) {
			t.Errorf("Expected field '%s' in JSON attributes: %s", field, string(jsonData))
		}
	}

	// Verify relationships are NOT in attributes (they're at envelope level now)
	if contains(string(jsonData), `"relationships"`) {
		t.Errorf("Relationships should not be in attributes (should be at envelope level): %s", string(jsonData))
	}
}

// indexof returns the index of substr in s, or -1 if not found
func indexof(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// contains checks if s contains substr
func contains(s, substr string) bool {
	return indexof(s, substr) >= 0
}
