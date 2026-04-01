package dto

import (
	"context"
	"encoding/json"
	"testing"
)

// TestLogResponse tests the LogResponse DTO
func TestLogResponse(t *testing.T) {
	data := "2024-01-01"
	response := NewLogResponse(1, &data, 10, 20)

	if response == nil {
		t.Fatal("Expected non-nil response, got nil")
	}

	if response.ID != 1 {
		t.Errorf("Expected ID 1, got %d", response.ID)
	}

	if *response.Data != "2024-01-01" {
		t.Errorf("Expected data '2024-01-01', got '%s'", *response.Data)
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
	data := "2024-01-01"
	response := NewLogResponse(1, &data, 10, 20)
	response.SetContext(context.Background())

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
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

	if *decoded.Data != "2024-01-01" {
		t.Errorf("Expected data '2024-01-01', got '%s'", *decoded.Data)
	}
}

// TestLogResponse_WithNote tests the response with optional note
func TestLogResponse_WithNote(t *testing.T) {
	data := "2024-01-01"
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

// TestLogResponse_WithProject tests the response with project
func TestLogResponse_WithProject(t *testing.T) {
	data := "2024-01-01"
	project := NewProjectResponse(1, "Test Project", nil, 100, 50)

	response := &LogResponse{
		ID:        1,
		Data:      &data,
		StartPage: 10,
		EndPage:   20,
		Project:   project,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verify project is included
	if !contains(string(jsonData), `"project"`) {
		t.Errorf("Expected project in JSON: %s", string(jsonData))
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

// contains checks if s contains substr
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
