package dto

import (
	"context"
	"encoding/json"
	"testing"
)

// TestHealthCheckResponse tests the HealthCheckResponse DTO
func TestHealthCheckResponse(t *testing.T) {
	// Test NewHealthCheckResponse
	response := NewHealthCheckResponse("healthy")

	if response == nil {
		t.Fatal("Expected non-nil response, got nil")
	}

	if response.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response.Status)
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

// TestHealthCheckResponse_JSON tests JSON serialization
func TestHealthCheckResponse_JSON(t *testing.T) {
	response := NewHealthCheckResponse("healthy")
	response.SetContext(context.Background())

	// Test JSON marshaling
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Test JSON unmarshaling
	var decoded HealthCheckResponse
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if decoded.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", decoded.Status)
	}
}

// TestHealthCheckResponse_WithMessage tests the response with optional message
func TestHealthCheckResponse_WithMessage(t *testing.T) {
	response := HealthCheckResponse{
		ctx:     context.Background(),
		Status:  "healthy",
		Message: "All systems operational",
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verify message is included
	if !containsString(string(jsonData), "All systems operational") {
		t.Errorf("Expected message in JSON, got: %s", string(jsonData))
	}
}
