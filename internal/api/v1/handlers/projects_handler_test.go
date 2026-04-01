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

	var response []*dto.ProjectResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response == nil {
		t.Error("Expected empty array, got nil")
	}

	if len(response) != 0 {
		t.Errorf("Expected 0 projects, got %d", len(response))
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

	var response []*dto.ProjectResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(response))
	}

	if response[0].ID != 1 {
		t.Errorf("Expected first project ID to be 1, got %d", response[0].ID)
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

	var response dto.ProjectResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.ID != 1 {
		t.Errorf("Expected project ID to be 1, got %d", response.ID)
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
