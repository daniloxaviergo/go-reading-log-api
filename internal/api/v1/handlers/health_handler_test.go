package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-reading-log-api-next/internal/domain/dto"
)

func TestHealthHandler_Healthz(t *testing.T) {
	handler := NewHealthHandler()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	handler.Healthz(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	expectedContentType := "application/json"
	actualContentType := w.Header().Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, actualContentType)
	}

	var response dto.HealthCheckResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response.Status)
	}
}

func TestHealthHandler_Healthz_GetMethod(t *testing.T) {
	handler := NewHealthHandler()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	handler.Healthz(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /healthz: expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHealthHandler_Healthz_PostMethod(t *testing.T) {
	handler := NewHealthHandler()

	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	w := httptest.NewRecorder()

	handler.Healthz(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("POST /healthz: expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestNewHealthHandler(t *testing.T) {
	handler := NewHealthHandler()

	if handler == nil {
		t.Fatal("Expected non-nil handler, got nil")
	}
}
