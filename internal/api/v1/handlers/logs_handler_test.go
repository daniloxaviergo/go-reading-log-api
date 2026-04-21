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

	// Verify Content-Type header
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/vnd.api+json" {
		t.Errorf("Expected Content-Type 'application/vnd.api+json', got '%s'", contentType)
	}

	// Decode JSON:API envelope
	var envelope dto.JSONAPIEnvelope
	if err := json.NewDecoder(w.Body).Decode(&envelope); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify data is an array (for empty collection)
	// Note: Empty slices may be encoded as []interface{} due to Go's json package behavior
	dataArray, ok := envelope.Data.([]dto.JSONAPIData)
	if !ok {
		// Try interface{} slice for empty collections
		if ifaceArr, ok := envelope.Data.([]interface{}); ok && len(ifaceArr) == 0 {
			// Empty array is valid - just verify it's an array
			return
		}
		t.Fatalf("Expected Data to be array of JSONAPIData, got %T", envelope.Data)
	}

	if len(dataArray) != 0 {
		t.Errorf("Expected 0 logs, got %d", len(dataArray))
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

	// Verify Content-Type header
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/vnd.api+json" {
		t.Errorf("Expected Content-Type 'application/vnd.api+json', got '%s'", contentType)
	}

	// Decode JSON:API envelope
	var envelope dto.JSONAPIEnvelope
	if err := json.NewDecoder(w.Body).Decode(&envelope); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify data structure - after json.Unmarshal, arrays become []interface{} containing map[string]interface{}
	dataArray, ok := envelope.Data.([]dto.JSONAPIData)
	if !ok {
		// Try interface{} slice for collections
		if ifaceArr, ok := envelope.Data.([]interface{}); ok {
			if len(ifaceArr) != 4 {
				t.Errorf("Expected 4 logs (limited), got %d", len(ifaceArr))
			}
			// Verify all IDs are strings
			for _, item := range ifaceArr {
				if idMap, ok := item.(map[string]interface{}); ok {
					if _, ok := idMap["id"].(string); !ok {
						t.Errorf("Expected ID to be string, got %T", idMap["id"])
					}
				}
			}
			return
		}
		t.Fatalf("Expected Data to be array of JSONAPIData, got %T", envelope.Data)
	}

	// Should limit to 4 logs
	if len(dataArray) != 4 {
		t.Errorf("Expected 4 logs (limited), got %d", len(dataArray))
	}

	// Verify all IDs are strings
	for _, item := range dataArray {
		if _, ok := item.ID.(string); !ok {
			t.Error("All log IDs must be strings")
		}
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

	// Verify Content-Type header
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/vnd.api+json" {
		t.Errorf("Expected Content-Type 'application/vnd.api+json', got '%s'", contentType)
	}

	// Decode JSON:API envelope
	var envelope dto.JSONAPIEnvelope
	if err := json.NewDecoder(w.Body).Decode(&envelope); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify data structure - after json.Unmarshal, arrays become []interface{} containing map[string]interface{}
	dataArray, ok := envelope.Data.([]dto.JSONAPIData)
	if !ok {
		// Try interface{} slice for collections
		if ifaceArr, ok := envelope.Data.([]interface{}); ok {
			if len(ifaceArr) != 2 {
				t.Errorf("Expected 2 logs, got %d", len(ifaceArr))
			}
			// Verify all IDs are strings
			for _, item := range ifaceArr {
				if idMap, ok := item.(map[string]interface{}); ok {
					if _, ok := idMap["id"].(string); !ok {
						t.Errorf("Expected ID to be string, got %T", idMap["id"])
					}
				}
			}
			return
		}
		t.Fatalf("Expected Data to be array of JSONAPIData, got %T", envelope.Data)
	}

	if len(dataArray) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(dataArray))
	}

	// Verify all IDs are strings
	for _, item := range dataArray {
		if _, ok := item.ID.(string); !ok {
			t.Error("All log IDs must be strings")
		}
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

	// Verify Content-Type header
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/vnd.api+json" {
		t.Errorf("Expected Content-Type 'application/vnd.api+json', got '%s'", contentType)
	}

	// Decode JSON:API envelope
	var envelope dto.JSONAPIEnvelope
	if err := json.NewDecoder(w.Body).Decode(&envelope); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify data structure - after json.Unmarshal, arrays become []interface{} containing map[string]interface{}
	dataArray, ok := envelope.Data.([]dto.JSONAPIData)
	if !ok {
		// Try interface{} slice for collections
		if ifaceArr, ok := envelope.Data.([]interface{}); ok {
			if len(ifaceArr) != 1 {
				t.Errorf("Expected 1 log, got %d", len(ifaceArr))
			}
			// Verify all IDs are strings
			for _, item := range ifaceArr {
				if idMap, ok := item.(map[string]interface{}); ok {
					if _, ok := idMap["id"].(string); !ok {
						t.Errorf("Expected ID to be string, got %T", idMap["id"])
					}
				}
			}
			return
		}
		t.Fatalf("Expected Data to be array of JSONAPIData, got %T", envelope.Data)
	}

	if len(dataArray) != 1 {
		t.Errorf("Expected 1 log, got %d", len(dataArray))
	}

	// Verify all IDs are strings
	for _, item := range dataArray {
		if _, ok := item.ID.(string); !ok {
			t.Error("All log IDs must be strings")
		}
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
