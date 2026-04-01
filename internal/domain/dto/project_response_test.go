package dto

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

// TestProjectResponse tests the ProjectResponse DTO
func TestProjectResponse(t *testing.T) {
	// Test NewProjectResponse
	startedAt := "2024-01-01T00:00:00Z"
	response := NewProjectResponse(1, "Test Project", &startedAt, 100, 50)

	if response == nil {
		t.Fatal("Expected non-nil response, got nil")
	}

	if response.ID != 1 {
		t.Errorf("Expected ID 1, got %d", response.ID)
	}

	if response.Name != "Test Project" {
		t.Errorf("Expected name 'Test Project', got '%s'", response.Name)
	}

	if response.TotalPage != 100 {
		t.Errorf("Expected total_page 100, got %d", response.TotalPage)
	}

	if response.Page != 50 {
		t.Errorf("Expected page 50, got %d", response.Page)
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

// TestProjectResponse_JSON tests JSON serialization
func TestProjectResponse_JSON(t *testing.T) {
	startedAt := "2024-01-01T00:00:00Z"
	response := NewProjectResponse(1, "Test Project", &startedAt, 100, 50)
	response.SetContext(context.Background())

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Test JSON unmarshaling
	var decoded ProjectResponse
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if decoded.ID != 1 {
		t.Errorf("Expected ID 1, got %d", decoded.ID)
	}

	if decoded.Name != "Test Project" {
		t.Errorf("Expected name 'Test Project', got '%s'", decoded.Name)
	}
}

// TestProjectResponse_AllFields tests the response with all fields set
func TestProjectResponse_AllFields(t *testing.T) {
	now := time.Now()
	startedAt := now.Format(time.RFC3339)
	status := "completed"
	logsCount := 10
	daysUnread := 5

	response := &ProjectResponse{
		ID:         1,
		Name:       "Test Project",
		StartedAt:  &startedAt,
		Progress:   floatPtr(50.5),
		TotalPage:  100,
		Page:       50,
		Status:     &status,
		LogsCount:  &logsCount,
		DaysUnread: &daysUnread,
		MedianDay:  &startedAt,
		FinishedAt: &startedAt,
		Logs:       []*LogResponse{},
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verify all fields are included
	result := string(jsonData)
	if !containsString(result, `"id":1`) {
		t.Errorf("Expected id in JSON: %s", result)
	}
	if !containsString(result, `"name":"Test Project"`) {
		t.Errorf("Expected name in JSON: %s", result)
	}
	if !containsString(result, `"progress":50.5`) {
		t.Errorf("Expected progress in JSON: %s", result)
	}
}

// TestProjectResponse_EmptyContext tests context fallback
func TestProjectResponse_EmptyContext(t *testing.T) {
	response := &ProjectResponse{}
	ctx := response.GetContext()
	if ctx == nil {
		t.Error("Expected default context, got nil")
	}
}

// TestProjectResponse_SetContext tests SetContext method
func TestProjectResponse_SetContext(t *testing.T) {
	response := &ProjectResponse{}
	
	newCtx := context.WithValue(context.Background(), "test_key", "test_value")
	response.SetContext(newCtx)

	ctx := response.GetContext()
	if ctx.Value("test_key") != "test_value" {
		t.Error("Context value was not set correctly")
	}
}

// floatPtr returns a pointer to a float64 value
func floatPtr(f float64) *float64 {
	return &f
}

// containsString checks if s contains substr
func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
