package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/domain/models"
	"go-reading-log-api-next/test"
)

// TestLogsHandler_Index tests the Index handler with no logs
func TestLogsHandler_Index(t *testing.T) {
	mockLogRepo := test.NewMockLogRepository()
	mockProjectRepo := test.NewMockProjectRepository()

	// Add project to project repo
	project := &models.Project{
		ID:        1,
		Name:      "Test Project",
		TotalPage: 100,
		Page:      50,
	}
	mockProjectRepo.AddProject(project)

	handler := NewLogsHandler(mockLogRepo, mockProjectRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects/1/logs", nil)
	req = mux.SetURLVars(req, map[string]string{"project_id": "1"})

	w := httptest.NewRecorder()

	handler.Index(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []*dto.LogResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response == nil {
		t.Error("Expected empty array, got nil")
	}

	if len(response) != 0 {
		t.Errorf("Expected 0 logs, got %d", len(response))
	}
}

// TestLogsHandler_IndexWithLogs tests the Index handler with logs
func TestLogsHandler_IndexWithLogs(t *testing.T) {
	mockLogRepo := test.NewMockLogRepository()
	mockProjectRepo := test.NewMockProjectRepository()

	// Add project to project repo
	project := &models.Project{
		ID:        1,
		Name:      "Test Project",
		TotalPage: 100,
		Page:      50,
	}
	mockProjectRepo.AddProject(project)

	// Add logs to log repo - Data is *string in Log model
	data1 := "2024-01-01"
	data2 := "2024-01-02"
	data3 := "2024-01-03"
	data4 := "2024-01-04"
	data5 := "2024-01-05"

	logs := []*models.Log{
		{ID: 1, Data: &data1, StartPage: 1, EndPage: 10},
		{ID: 2, Data: &data2, StartPage: 11, EndPage: 20},
		{ID: 3, Data: &data3, StartPage: 21, EndPage: 30},
		{ID: 4, Data: &data4, StartPage: 31, EndPage: 40},
		{ID: 5, Data: &data5, StartPage: 41, EndPage: 50},
	}
	mockLogRepo.AddLogsForProject(1, logs)

	handler := NewLogsHandler(mockLogRepo, mockProjectRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects/1/logs", nil)
	req = mux.SetURLVars(req, map[string]string{"project_id": "1"})

	w := httptest.NewRecorder()

	handler.Index(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []*dto.LogResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Should limit to 4 logs
	if len(response) != 4 {
		t.Errorf("Expected 4 logs (limited), got %d", len(response))
	}
}

// TestLogsHandler_IndexWithLessThanLimit tests the Index handler with less than 4 logs
func TestLogsHandler_IndexWithLessThanLimit(t *testing.T) {
	mockLogRepo := test.NewMockLogRepository()
	mockProjectRepo := test.NewMockProjectRepository()

	project := &models.Project{
		ID:        1,
		Name:      "Test Project",
		TotalPage: 100,
		Page:      50,
	}
	mockProjectRepo.AddProject(project)

	data1 := "2024-01-01"
	data2 := "2024-01-02"

	logs := []*models.Log{
		{ID: 1, Data: &data1, StartPage: 1, EndPage: 10},
		{ID: 2, Data: &data2, StartPage: 11, EndPage: 20},
	}
	mockLogRepo.AddLogsForProject(1, logs)

	handler := NewLogsHandler(mockLogRepo, mockProjectRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects/1/logs", nil)
	req = mux.SetURLVars(req, map[string]string{"project_id": "1"})

	w := httptest.NewRecorder()

	handler.Index(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []*dto.LogResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(response))
	}
}

// TestLogsHandler_IndexWithOneLog tests the Index handler with exactly 1 log
func TestLogsHandler_IndexWithOneLog(t *testing.T) {
	mockLogRepo := test.NewMockLogRepository()
	mockProjectRepo := test.NewMockProjectRepository()

	project := &models.Project{
		ID:        1,
		Name:      "Test Project",
		TotalPage: 100,
		Page:      50,
	}
	mockProjectRepo.AddProject(project)

	data1 := "2024-01-01"
	logs := []*models.Log{
		{ID: 1, Data: &data1, StartPage: 1, EndPage: 10},
	}
	mockLogRepo.AddLogsForProject(1, logs)

	handler := NewLogsHandler(mockLogRepo, mockProjectRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects/1/logs", nil)
	req = mux.SetURLVars(req, map[string]string{"project_id": "1"})

	w := httptest.NewRecorder()

	handler.Index(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []*dto.LogResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("Expected 1 log, got %d", len(response))
	}
}

// TestLogsHandler_IndexNotFound tests the Index handler with non-existent project
func TestLogsHandler_IndexNotFound(t *testing.T) {
	mockLogRepo := test.NewMockLogRepository()
	mockProjectRepo := test.NewMockProjectRepository()

	handler := NewLogsHandler(mockLogRepo, mockProjectRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects/999/logs", nil)
	req = mux.SetURLVars(req, map[string]string{"project_id": "999"})

	w := httptest.NewRecorder()

	handler.Index(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	body := strings.TrimSpace(w.Body.String())
	expectedBody := `{"error": "project not found"}`
	if body != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, body)
	}
}

// TestLogsHandler_IndexInvalidProjectID tests the Index handler with invalid project ID
func TestLogsHandler_IndexInvalidProjectID(t *testing.T) {
	mockLogRepo := test.NewMockLogRepository()
	mockProjectRepo := test.NewMockProjectRepository()

	handler := NewLogsHandler(mockLogRepo, mockProjectRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects/invalid/logs", nil)
	req = mux.SetURLVars(req, map[string]string{"project_id": "invalid"})

	w := httptest.NewRecorder()

	handler.Index(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	body := strings.TrimSpace(w.Body.String())
	expectedBody := `{"error": "Invalid project ID"}`
	if body != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, body)
	}
}

// TestLogsHandler_IndexProjectRepoError tests the Index handler with project repository error
func TestLogsHandler_IndexProjectRepoError(t *testing.T) {
	mockLogRepo := test.NewMockLogRepository()
	mockProjectRepo := test.NewMockProjectRepository()
	mockProjectRepo.SetError(fmt.Errorf("database error"))

	handler := NewLogsHandler(mockLogRepo, mockProjectRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects/1/logs", nil)
	req = mux.SetURLVars(req, map[string]string{"project_id": "1"})

	w := httptest.NewRecorder()

	handler.Index(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if !strings.Contains(w.Body.String(), "Internal server error") {
		t.Errorf("Expected 'Internal server error' in response, got %s", w.Body.String())
	}
}

// TestLogsHandler_IndexLogRepoError tests the Index handler with log repository error
func TestLogsHandler_IndexLogRepoError(t *testing.T) {
	mockLogRepo := test.NewMockLogRepository()
	mockProjectRepo := test.NewMockProjectRepository()

	project := &models.Project{
		ID:        1,
		Name:      "Test Project",
		TotalPage: 100,
		Page:      50,
	}
	mockProjectRepo.AddProject(project)
	mockLogRepo.SetError(fmt.Errorf("database error"))

	handler := NewLogsHandler(mockLogRepo, mockProjectRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects/1/logs", nil)
	req = mux.SetURLVars(req, map[string]string{"project_id": "1"})

	w := httptest.NewRecorder()

	handler.Index(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if !strings.Contains(w.Body.String(), "Internal server error") {
		t.Errorf("Expected 'Internal server error' in response, got %s", w.Body.String())
	}
}

// TestLogsHandler_IndexEmptyProject tests the Index handler with empty project response
func TestLogsHandler_IndexEmptyProject(t *testing.T) {
	mockLogRepo := test.NewMockLogRepository()
	mockProjectRepo := test.NewMockProjectRepository()

	handler := NewLogsHandler(mockLogRepo, mockProjectRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects/1/logs", nil)
	req = mux.SetURLVars(req, map[string]string{"project_id": "1"})

	w := httptest.NewRecorder()

	handler.Index(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	body := strings.TrimSpace(w.Body.String())
	expectedBody := `{"error": "project not found"}`
	if body != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, body)
	}
}

// TestFormatTimePtr tests the formatTimePtr helper function
func TestFormatTimePtr(t *testing.T) {
	t.Run("nil time", func(t *testing.T) {
		result := formatTimePtr(nil)
		if result != nil {
			t.Errorf("Expected nil, got %v", result)
		}
	})

	t.Run("valid time", func(t *testing.T) {
		tm := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		result := formatTimePtr(&tm)
		if result == nil {
			t.Fatal("Expected non-nil result")
		}

		expected := "2024-01-01T12:00:00Z"
		if *result != expected {
			t.Errorf("Expected %s, got %s", expected, *result)
		}
	})
}

// TestNewLogsHandler tests the constructor
func TestNewLogsHandler(t *testing.T) {
	mockLogRepo := test.NewMockLogRepository()
	mockProjectRepo := test.NewMockProjectRepository()

	handler := NewLogsHandler(mockLogRepo, mockProjectRepo)

	if handler == nil {
		t.Fatal("Expected non-nil handler, got nil")
	}
}
