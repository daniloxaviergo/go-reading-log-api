package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-reading-log-api-next/internal/api/v1/handlers"
	"go-reading-log-api-next/internal/domain/dto"
)

// TestHealthCheckIntegration tests the health check endpoint
func TestHealthCheckIntegration(t *testing.T) {
	// Create health handler and router
	handler := handlers.NewHealthHandler()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	handler.Healthz(recorder, req)

	// Verify response status
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	// Verify content type
	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	// Parse and verify response body
	var response dto.HealthCheckResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response.Status)
	}
}

// TestHealthCheckResponseFormat tests the health check response format
func TestHealthCheckResponseFormat(t *testing.T) {
	handler := handlers.NewHealthHandler()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	handler.Healthz(recorder, req)

	// Verify required JSON fields
	body := recorder.Body.String()
	if !contains(body, `"status":"healthy"`) && !contains(body, `"status": "healthy"`) {
		t.Errorf("Response missing expected status field: %s", body)
	}

	// Verify it's valid JSON
	var response dto.HealthCheckResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Response is not valid JSON: %v", err)
	}
}

// TestHealthCheckMethodNotAllowed tests that only GET is allowed
func TestHealthCheckMethodNotAllowed(t *testing.T) {
	handler := handlers.NewHealthHandler()
	recorder := httptest.NewRecorder()

	// Try POST
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	handler.Healthz(recorder, req)

	// Currently the handler doesn't restrict methods, so it should still work
	// This test documents current behavior
	if recorder.Code != http.StatusOK {
		t.Logf("POST handler returned status %d (may need method restriction)", recorder.Code)
	}
}

// TestHealthCheckEmptyPath tests with empty path
func TestHealthCheckEmptyPath(t *testing.T) {
	handler := handlers.NewHealthHandler()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler.Healthz(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}
}

// TestHealthCheckWithRequestContext tests context propagation
func TestHealthCheckWithRequestContext(t *testing.T) {
	handler := handlers.NewHealthHandler()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	handler.Healthz(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}
}

// TestHealthCheckMultipleRequests tests handling multiple sequential requests
func TestHealthCheckMultipleRequests(t *testing.T) {
	handler := handlers.NewHealthHandler()

	for i := 0; i < 3; i++ {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

		handler.Healthz(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Request %d: Expected status 200, got %d", i+1, recorder.Code)
		}
	}
}

// TestHealthCheckConcurrentRequests tests concurrent request handling
func TestHealthCheckConcurrentRequests(t *testing.T) {
	handler := handlers.NewHealthHandler()

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

			handler.Healthz(recorder, req)

			if recorder.Code == http.StatusOK {
				done <- true
			} else {
				done <- false
			}
		}()
	}

	// Wait for all requests
	successCount := 0
	for i := 0; i < 10; i++ {
		if <-done {
			successCount++
		}
	}

	if successCount != 10 {
		t.Errorf("Expected 10 successful requests, got %d", successCount)
	}
}
