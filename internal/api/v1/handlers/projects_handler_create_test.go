package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/domain/models"
	"go-reading-log-api-next/internal/validation"
	"go-reading-log-api-next/test"
)

// TestProjectsHandler_Create validates success cases
func TestProjectsHandler_Create(t *testing.T) {
	// Test case 1: Valid project with page < total_page
	t.Run("valid_project_page_less_than_total", func(t *testing.T) {
		mockRepo := test.NewMockProjectRepository()

		// Pre-add a project with ID 1
		project := &models.Project{
			ID:        1,
			Name:      "Test Project",
			TotalPage: 100,
			Page:      50,
		}
		mockRepo.AddProject(project)

		handler := NewProjectsHandler(mockRepo)

		reqBody := dto.ProjectRequest{
			Name:      "New Project",
			TotalPage: 200,
			Page:      100,
			Reinicia:  false,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
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

		// Verify data structure - after json.Unmarshal, single objects become map[string]interface{}
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

		// Verify ID is string
		idVal, ok := dataMap["id"]
		if !ok {
			t.Error("Missing 'id' field in data object")
		} else if _, ok := idVal.(string); !ok {
			t.Errorf("Expected ID to be string type, got %T", idVal)
		}

		// Verify attributes structure
		attrsMap, ok := dataMap["attributes"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected Attributes to be an object")
		}

		if nameStr, ok := attrsMap["name"].(string); !ok || nameStr != "New Project" {
			t.Errorf("Expected name 'New Project', got '%v'", attrsMap["name"])
		}

		if totalPage, ok := attrsMap["total_page"].(float64); !ok || int(totalPage) != 200 {
			t.Errorf("Expected total_page 200, got %v", attrsMap["total_page"])
		}

		if pageVal, ok := attrsMap["page"].(float64); !ok || int(pageVal) != 100 {
			t.Errorf("Expected page 100, got %v", attrsMap["page"])
		}
	})

	// Test case 2: Valid project with page = total_page
	t.Run("valid_project_page_equals_total", func(t *testing.T) {
		mockRepo := test.NewMockProjectRepository()

		handler := NewProjectsHandler(mockRepo)

		reqBody := dto.ProjectRequest{
			Name:      "Equal Page Project",
			TotalPage: 100,
			Page:      100,
			Reinicia:  false,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}
	})

	// Test case 3: Valid project with page = 0
	t.Run("valid_project_page_zero", func(t *testing.T) {
		mockRepo := test.NewMockProjectRepository()

		handler := NewProjectsHandler(mockRepo)

		reqBody := dto.ProjectRequest{
			Name:      "Zero Page Project",
			TotalPage: 100,
			Page:      0,
			Reinicia:  false,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}
	})
}

// TestProjectsHandler_CreateValidationErrors tests validation error cases
func TestProjectsHandler_CreateValidationErrors(t *testing.T) {
	// Test case 1: page > total_page (should fail validation)
	t.Run("validation_error_page_exceeds_total", func(t *testing.T) {
		mockRepo := test.NewMockProjectRepository()

		handler := NewProjectsHandler(mockRepo)

		reqBody := dto.ProjectRequest{
			Name:      "Invalid Project",
			TotalPage: 50,
			Page:      100, // page > total_page - should fail
			Reinicia:  false,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d for validation error, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response["error"] != "validation failed" {
			t.Errorf("Expected error 'validation failed', got '%v'", response["error"])
		}

		details, ok := response["details"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected details to be a map, got %T", response["details"])
		}

		// Check that page validation error is present
		if _, exists := details["page"]; !exists {
			t.Error("Expected 'page' in validation details")
		}
	})

	// Test case 2: negative page value
	t.Run("validation_error_negative_page", func(t *testing.T) {
		mockRepo := test.NewMockProjectRepository()

		handler := NewProjectsHandler(mockRepo)

		reqBody := dto.ProjectRequest{
			Name:      "Invalid Project",
			TotalPage: 100,
			Page:      -10, // negative page - should fail
			Reinicia:  false,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d for validation error, got %d", http.StatusBadRequest, w.Code)
		}
	})

	// Test case 3: zero total_page
	t.Run("validation_error_zero_total_page", func(t *testing.T) {
		mockRepo := test.NewMockProjectRepository()

		handler := NewProjectsHandler(mockRepo)

		reqBody := dto.ProjectRequest{
			Name:      "Invalid Project",
			TotalPage: 0, // zero total_page - should fail
			Page:      0,
			Reinicia:  false,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		handler.Create(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d for validation error, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// TestProjectsHandler_CreateInvalidJSON tests invalid JSON body
func TestProjectsHandler_CreateInvalidJSON(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()

	handler := NewProjectsHandler(mockRepo)

	invalidBody := []byte(`{invalid json`)

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for invalid JSON, got %d", http.StatusBadRequest, w.Code)
	}

	// Check the response contains the expected error message (allowing for JSON formatting differences)
	responseBody := w.Body.String()
	if !containsString(responseBody, `"error":"invalid request body"`) && !containsString(responseBody, `"error": "invalid request body"`) {
		t.Errorf("Expected body to contain error message, got %s", responseBody)
	}
}

// TestProjectsHandler_CreateRepositoryError tests repository errors
func TestProjectsHandler_CreateRepositoryError(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()
	mockRepo.SetError(fmt.Errorf("database connection failed"))

	handler := NewProjectsHandler(mockRepo)

	// Use a page > total_page to fail validation, so we can test the repo error path
	// But wait - we need to test the repo error path after validation passes
	// Since validation passes for valid data, we use valid data and let repo error occur
	reqBody := dto.ProjectRequest{
		Name:      "Test Project",
		TotalPage: 100,
		Page:      50, // Valid: page <= total_page
		Reinicia:  false,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Create(w, req)

	// Validation passes, then repo error occurs
	// The handler should return 500 with "Internal server error"
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d for repository error, got %d", http.StatusInternalServerError, w.Code)
	}

	// Verify the response body
	responseBody := w.Body.String()
	if !containsString(responseBody, `"error":"Internal server error"`) && !containsString(responseBody, `"error": "Internal server error"`) {
		t.Errorf("Expected body to contain 'Internal server error', got %s", responseBody)
	}
}

// TestProjectsHandler_CreateEmptyBody tests empty request body
func TestProjectsHandler_CreateEmptyBody(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()

	handler := NewProjectsHandler(mockRepo)

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for empty body, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestProjectsHandler_CreateValidationIntegration tests that validation is actually called
func TestProjectsHandler_CreateValidationIntegration(t *testing.T) {
	// Test that the validation error format is correct
	mockRepo := test.NewMockProjectRepository()

	handler := NewProjectsHandler(mockRepo)

	// This should fail validation: page > total_page
	reqBody := dto.ProjectRequest{
		Name:      "Validation Test",
		TotalPage: 50,
		Page:      100,
		Reinicia:  false,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify error structure
	errorMsg, ok := response["error"].(string)
	if !ok || errorMsg != "validation failed" {
		t.Errorf("Expected error message 'validation failed', got '%v'", response["error"])
	}

	details, ok := response["details"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected details to be a map, got %T", response["details"])
	}

	// Verify page validation error message contains expected text
	pageErr, ok := details["page"].(string)
	if !ok {
		t.Fatalf("Expected 'page' error to be a string, got %T", details["page"])
	}

	// The validation message should contain "exceeds" or similar
	if !validationContains(pageErr, "exceeds") && !validationContains(pageErr, "cannot") {
		t.Errorf("Expected validation error message to contain 'exceeds' or 'cannot', got '%s'", pageErr)
	}
}

// validationContains is a helper to check if a string contains a substring
func validationContains(s, substr string) bool {
	return len(s) >= len(substr) && contains(s, substr)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestProjectsHandler_CreateWithStartedAt tests project creation with started_at
func TestProjectsHandler_CreateWithStartedAt(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()

	handler := NewProjectsHandler(mockRepo)

	// Valid started_at in RFC3339 format
	now := "2025-01-15T10:30:00Z"
	reqBody := dto.ProjectRequest{
		Name:      "Project with Date",
		TotalPage: 100,
		Page:      50,
		StartedAt: &now,
		Reinicia:  false,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
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

	// Verify data structure - after json.Unmarshal, single objects become map[string]interface{}
	dataMap, ok := envelope.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected Data to be a JSON object, got %T", envelope.Data)
	}

	// Verify attributes structure
	attrsMap, ok := dataMap["attributes"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected Attributes to be an object")
	}

	if nameStr, ok := attrsMap["name"].(string); !ok || nameStr != "Project with Date" {
		t.Errorf("Expected name 'Project with Date', got '%v'", attrsMap["name"])
	}
}

// TestProjectsHandler_CreateWithInvalidDate tests project creation with invalid date format
func TestProjectsHandler_CreateWithInvalidDate(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()

	handler := NewProjectsHandler(mockRepo)

	invalidDate := "not-a-date"
	reqBody := dto.ProjectRequest{
		Name:      "Project with Invalid Date",
		TotalPage: 100,
		Page:      50,
		StartedAt: &invalidDate,
		Reinicia:  false,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d for invalid date, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["error"] != "invalid date format" {
		t.Errorf("Expected error 'invalid date format', got '%v'", response["error"])
	}
}

// TestProjectsHandler_CreateReiniciaField tests project creation with reinicia flag
func TestProjectsHandler_CreateReiniciaField(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()

	handler := NewProjectsHandler(mockRepo)

	reqBody := dto.ProjectRequest{
		Name:      "Project with Reinicia",
		TotalPage: 100,
		Page:      50,
		Reinicia:  true,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
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

	// Verify data structure - after json.Unmarshal, single objects become map[string]interface{}
	dataMap, ok := envelope.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected Data to be a JSON object, got %T", envelope.Data)
	}

	// Verify attributes structure
	attrsMap, ok := dataMap["attributes"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected Attributes to be an object")
	}

	// Reinicia is not included in ProjectResponse - just verify other fields
	if nameStr, ok := attrsMap["name"].(string); !ok || nameStr != "Project with Reinicia" {
		t.Errorf("Expected name 'Project with Reinicia', got '%v'", attrsMap["name"])
	}
}

// TestProjectsHandler_CreateWithoutRequiredFields tests missing required fields
func TestProjectsHandler_CreateWithoutRequiredFields(t *testing.T) {
	// Test with empty name
	t.Run("empty_name", func(t *testing.T) {
		mockRepo := test.NewMockProjectRepository()

		handler := NewProjectsHandler(mockRepo)

		// Empty name is technically allowed by the struct but might be invalid for DB
		reqBody := dto.ProjectRequest{
			Name:      "",
			TotalPage: 100,
			Page:      50,
			Reinicia:  false,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		handler.Create(w, req)

		// The handler doesn't validate empty name, so it depends on DB constraint
		// For now, we just verify the handler processes it
		if w.Code != http.StatusCreated && w.Code != http.StatusInternalServerError {
			t.Errorf("Unexpected status code %d", w.Code)
		}
	})

	// Test with missing required fields (total_page, page)
	t.Run("missing_total_page", func(t *testing.T) {
		mockRepo := test.NewMockProjectRepository()

		handler := NewProjectsHandler(mockRepo)

		reqBody := map[string]interface{}{
			"name":     "Test",
			"page":     50,
			"reinicia": false,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/v1/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		handler.Create(w, req)

		// JSON unmarshals missing int fields as 0, which will fail validation
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d for invalid data, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// TestNewProjectsHandler_Create tests that the handler is properly initialized
func TestNewProjectsHandler_Create(t *testing.T) {
	mockRepo := test.NewMockProjectRepository()
	handler := NewProjectsHandler(mockRepo)

	if handler == nil {
		t.Fatal("Expected non-nil handler, got nil")
	}

	// Verify handler has repository
	if handler.repo == nil {
		t.Error("Expected handler to have a repository")
	}
}

// TestProjectsHandler_CreateMethodSignature tests that Create method exists
func TestProjectsHandler_CreateMethodSignature(t *testing.T) {
	handler := NewProjectsHandler(test.NewMockProjectRepository())

	// Verify handler has Create method
	// This is a compile-time check - if Create method doesn't exist, this won't compile
	_ = handler.Create
}

// TestProjectsHandler_CreateValidationHelper tests the validation package integration
func TestProjectsHandler_CreateValidationHelper(t *testing.T) {
	// Test validation directly
	tests := []struct {
		name       string
		page       int
		totalPage  int
		status     string
		shouldPass bool
		errorField string
	}{
		{"valid page < total", 50, 100, "unstarted", true, ""},
		{"valid page = total", 100, 100, "unstarted", true, ""},
		{"valid page = 0", 0, 100, "unstarted", true, ""},
		{"invalid page > total", 150, 100, "unstarted", false, "page"},
		{"invalid page < 0", -10, 100, "unstarted", false, "page"},
		{"invalid total_page = 0", 50, 0, "unstarted", false, "total_page"},
		{"valid with status", 50, 100, "running", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateProject(tt.page, tt.totalPage, tt.status)

			if tt.shouldPass {
				// Should pass validation - no errors expected
				if err != nil && err.HasErrors() {
					t.Errorf("Expected no validation errors for valid data, got: %v", err)
				}
			} else {
				// Should fail validation - errors expected
				if err == nil || !err.HasErrors() {
					t.Errorf("Expected validation errors for invalid data")
				}
				if tt.errorField != "" && err != nil {
					details := err.ToMap()
					if details != nil {
						if _, exists := details[tt.errorField]; !exists {
							t.Errorf("Expected '%s' in validation details", tt.errorField)
						}
					}
				}
			}
		})
	}
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
