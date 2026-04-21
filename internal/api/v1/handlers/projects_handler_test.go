package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/domain/models"
	"go-reading-log-api-next/test"
)

// TestProjectsHandler_Index tests the Index handler with empty projects list
func TestProjectsHandler_Index(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()
	handler := NewProjectsHandler(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
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

	// Verify data is an array of JSONAPIData
	// Note: Empty slices may be encoded as []interface{} due to Go's json package behavior
	dataArray, ok := envelope.Data.([]dto.JSONAPIData)
	if !ok {
		// Try interface{} slice for empty collections
		if ifaceArr, ok := envelope.Data.([]interface{}); ok && len(ifaceArr) == 0 {
			// Empty array is valid
			return
		}
		t.Fatalf("Expected Data to be array of JSONAPIData, got %T", envelope.Data)
	}

	if len(dataArray) != 0 {
		t.Errorf("Expected 0 projects, got %d", len(dataArray))
	}

	// Verify all IDs are strings
	for _, item := range dataArray {
		if _, ok := item.ID.(string); !ok {
			t.Error("All project IDs must be strings")
		}
	}
}

// TestProjectsHandler_IndexWithProjects tests the Index handler with multiple projects
func TestProjectsHandler_IndexWithProjects(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()

	project1 := &models.Project{
		ID:        1,
		Name:      "Project 1",
		TotalPage: 100,
		Page:      50,
	}
	project2 := &models.Project{
		ID:        2,
		Name:      "Project 2",
		TotalPage: 200,
		Page:      100,
	}

	mockRepo.AddProject(project1)
	mockRepo.AddProject(project2)

	handler := NewProjectsHandler(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
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
				t.Errorf("Expected 2 projects, got %d", len(ifaceArr))
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
		t.Errorf("Expected 2 projects, got %d", len(dataArray))
	}

	// Verify all IDs are strings
	for _, item := range dataArray {
		if _, ok := item.ID.(string); !ok {
			t.Error("All project IDs must be strings")
		}
	}
}

// TestProjectsHandler_IndexRepositoryError tests the Index handler with repository error
func TestProjectsHandler_IndexRepositoryError(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()
	mockRepo.SetError(fmt.Errorf("database connection failed"))

	handler := NewProjectsHandler(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	w := httptest.NewRecorder()

	handler.Index(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	expectedBody := `{"error": "Internal server error"}`
	if strings.TrimSpace(w.Body.String()) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

// TestProjectsHandler_Show tests the Show handler with valid project ID
func TestProjectsHandler_Show(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()

	responseDTO := &dto.ProjectResponse{
		ID:        1,
		Name:      "Test Project",
		TotalPage: 100,
	}

	mockRepo.AddProjectWithLogs(responseDTO)

	handler := NewProjectsHandler(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	w := httptest.NewRecorder()

	handler.Show(w, req)

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

	// Verify data structure - after json.Unmarshal, the interface{} becomes map[string]interface{}
	// We need to verify the structure manually
	dataMap, ok := envelope.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected Data to be a JSON object, got %T", envelope.Data)
	}

	// Verify required fields exist
	if _, ok := dataMap["type"]; !ok {
		t.Error("Missing 'type' field in data object")
	}
	if _, ok := dataMap["attributes"]; !ok {
		t.Error("Missing 'attributes' field in data object")
	}

	// Verify ID is string (could be int or string depending on JSON)
	idVal, ok := dataMap["id"]
	if !ok {
		t.Error("Missing 'id' field in data object")
	} else {
		// ID should be string per JSON:API spec
		if _, ok := idVal.(string); !ok {
			t.Errorf("Expected ID to be string type, got %T", idVal)
		}
	}

	// Verify attributes structure
	attrsMap, ok := dataMap["attributes"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected Attributes to be an object")
	}

	if nameStr, ok := attrsMap["name"].(string); !ok || nameStr != "Test Project" {
		t.Errorf("Expected name 'Test Project', got '%v'", attrsMap["name"])
	}

	if totalPage, ok := attrsMap["total_page"].(float64); !ok || int(totalPage) != 100 {
		t.Errorf("Expected total_page 100, got %v", attrsMap["total_page"])
	}
}

// TestProjectsHandler_ShowNotFound tests the Show handler with non-existent project
func TestProjectsHandler_ShowNotFound(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()

	handler := NewProjectsHandler(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})

	w := httptest.NewRecorder()

	handler.Show(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	expectedBody := `{"error": "project not found"}`
	if strings.TrimSpace(w.Body.String()) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

// TestProjectsHandler_ShowInvalidID tests the Show handler with invalid project ID
func TestProjectsHandler_ShowInvalidID(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()

	handler := NewProjectsHandler(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})

	w := httptest.NewRecorder()

	handler.Show(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedBody := `{"error": "Invalid project ID"}`
	if strings.TrimSpace(w.Body.String()) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

// TestProjectsHandler_ShowRepositoryError tests the Show handler with repository error
func TestProjectsHandler_ShowRepositoryError(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()
	mockRepo.SetError(fmt.Errorf("database error"))

	handler := NewProjectsHandler(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	w := httptest.NewRecorder()

	handler.Show(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	expectedBody := `{"error": "Internal server error"}`
	if strings.TrimSpace(w.Body.String()) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

// TestProjectsHandler_ShowRepositoryNotFoundError tests the Show handler with not found error
func TestProjectsHandler_ShowRepositoryNotFoundError(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()
	mockRepo.SetError(fmt.Errorf("project with ID 1 not found"))

	handler := NewProjectsHandler(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/projects/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	w := httptest.NewRecorder()

	handler.Show(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	expectedBody := `{"error": "project not found"}`
	if strings.TrimSpace(w.Body.String()) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

// TestNewProjectsHandler tests the constructor
func TestNewProjectsHandler(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()
	handler := NewProjectsHandler(mockRepo)

	if handler == nil {
		t.Fatal("Expected non-nil handler, got nil")
	}
}
