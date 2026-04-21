package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestProjectsCreateIntegration tests the POST /v1/projects endpoint
func TestProjectsCreateIntegration(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	// Create a valid project request
	reqBody := map[string]interface{}{
		"name":       "Integration Test Project",
		"total_page": 200,
		"page":       50,
		"reinicia":   false,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/projects.json", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := ctx.MakeRequest(t, req)

	// Verify response status
	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}

	// Parse and verify response
	response := ctx.ParseProjectResponse(t, recorder.Body.String())

	// Verify response fields
	if response.Name != "Integration Test Project" {
		t.Errorf("Expected name 'Integration Test Project', got '%s'", response.Name)
	}

	if response.TotalPage != 200 {
		t.Errorf("Expected total_page 200, got %d", response.TotalPage)
	}

	if response.Page != 50 {
		t.Errorf("Expected page 50, got %d", response.Page)
	}

	if response.ID == 0 {
		t.Error("Expected non-zero project ID")
	}

	// started_at is optional and may be nil if not provided in request
}

// TestProjectsCreateValidationErrors tests validation error cases
func TestProjectsCreateValidationErrors(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	tests := []struct {
		name           string
		reqBody        map[string]interface{}
		expectedStatus int
		errorField     string
		errorContains  string
	}{
		{
			name: "page_exceeds_total",
			reqBody: map[string]interface{}{
				"name":       "Invalid Project",
				"total_page": 50,
				"page":       100, // page > total_page
				"reinicia":   false,
			},
			expectedStatus: http.StatusBadRequest,
			errorField:     "page",
			errorContains:  "exceed",
		},
		{
			name: "negative_page",
			reqBody: map[string]interface{}{
				"name":       "Invalid Project",
				"total_page": 100,
				"page":       -10, // negative page
				"reinicia":   false,
			},
			expectedStatus: http.StatusBadRequest,
			errorField:     "page",
			errorContains:  "negative",
		},
		{
			name: "zero_total_page",
			reqBody: map[string]interface{}{
				"name":       "Invalid Project",
				"total_page": 0, // zero total_page
				"page":       0,
				"reinicia":   false,
			},
			expectedStatus: http.StatusBadRequest,
			errorField:     "total_page",
			errorContains:  "greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.reqBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/projects.json", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			recorder := ctx.MakeRequest(t, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, recorder.Code)
				t.Errorf("Response body: %s", recorder.Body.String())
			}

			// Verify error response format
			var response map[string]interface{}
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if response["error"] != "validation failed" {
				t.Errorf("Expected error 'validation failed', got '%v'", response["error"])
			}

			details, ok := response["details"].(map[string]interface{})
			if !ok {
				t.Fatalf("Expected details to be a map, got %T", response["details"])
			}

			if fieldErr, exists := details[tt.errorField]; exists {
				if errStr, ok := fieldErr.(string); ok {
					if !contains(errStr, tt.errorContains) {
						t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorContains, errStr)
					}
				}
			} else {
				t.Errorf("Expected '%s' in validation details", tt.errorField)
			}
		})
	}
}

// TestProjectsCreateWithStartedAt tests project creation with started_at date
func TestProjectsCreateWithStartedAt(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	// Use a valid RFC3339 date
	startedAt := "2024-01-15T10:30:00Z"

	reqBody := map[string]interface{}{
		"name":       "Project with Date",
		"total_page": 100,
		"page":       50,
		"started_at": &startedAt,
		"reinicia":   false,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/projects.json", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := ctx.MakeRequest(t, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}

	response := ctx.ParseProjectResponse(t, recorder.Body.String())

	if response.Name != "Project with Date" {
		t.Errorf("Expected name 'Project with Date', got '%s'", response.Name)
	}

	if response.StartedAt == nil {
		t.Error("Expected started_at to be present")
	}
}

// TestProjectsCreateInvalidDate tests project creation with invalid date format
func TestProjectsCreateInvalidDate(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	invalidDate := "not-a-date"

	reqBody := map[string]interface{}{
		"name":       "Project with Invalid Date",
		"total_page": 100,
		"page":       50,
		"started_at": &invalidDate,
		"reinicia":   false,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/projects.json", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := ctx.MakeRequest(t, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid date, got %d", http.StatusBadRequest, recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["error"] != "invalid date format" {
		t.Errorf("Expected error 'invalid date format', got '%v'", response["error"])
	}
}

// TestProjectsCreateWithReinicia tests project creation with reinicia flag
func TestProjectsCreateWithReinicia(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	reqBody := map[string]interface{}{
		"name":       "Project with Reinicia",
		"total_page": 100,
		"page":       50,
		"reinicia":   true,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/projects.json", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := ctx.MakeRequest(t, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}

	response := ctx.ParseProjectResponse(t, recorder.Body.String())

	// Reinicia is not included in ProjectResponse - just verify other fields
	if response.Name != "Project with Reinicia" {
		t.Errorf("Expected name 'Project with Reinicia', got '%s'", response.Name)
	}
}

// TestProjectsCreateInvalidJSON tests invalid JSON body
func TestProjectsCreateInvalidJSON(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	invalidBody := []byte(`{invalid json`)

	req := httptest.NewRequest(http.MethodPost, "/v1/projects.json", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	recorder := ctx.MakeRequest(t, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusBadRequest, recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["error"] != "invalid request body" {
		t.Errorf("Expected error 'invalid request body', got '%v'", response["error"])
	}
}

// TestProjectsCreateEmptyBody tests empty request body
func TestProjectsCreateEmptyBody(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	req := httptest.NewRequest(http.MethodPost, "/v1/projects.json", nil)
	req.Header.Set("Content-Type", "application/json")
	recorder := ctx.MakeRequest(t, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for empty body, got %d", http.StatusBadRequest, recorder.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["error"] != "invalid request body" {
		t.Errorf("Expected error 'invalid request body', got '%v'", response["error"])
	}
}

// TestProjectsCreateRetrieve tests that created project can be retrieved
func TestProjectsCreateRetrieve(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	// Create a project
	reqBody := map[string]interface{}{
		"name":       "Retrieve Test Project",
		"total_page": 100,
		"page":       25,
		"reinicia":   false,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/projects.json", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := ctx.MakeRequest(t, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d for creation, got %d", http.StatusCreated, recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}

	created := ctx.ParseProjectResponse(t, recorder.Body.String())

	// Retrieve the created project
	getReq := httptest.NewRequest(http.MethodGet, "/v1/projects/"+string(rune('0'+created.ID))+".json", nil)
	getRecorder := ctx.MakeRequest(t, getReq)

	if getRecorder.Code != http.StatusOK {
		t.Errorf("Expected status %d for GET, got %d", http.StatusOK, getRecorder.Code)
		t.Errorf("Response body: %s", getRecorder.Body.String())
	}

	retrieved := ctx.ParseProjectResponse(t, getRecorder.Body.String())

	if retrieved.ID != created.ID {
		t.Errorf("Expected retrieved project ID %d, got %d", created.ID, retrieved.ID)
	}

	if retrieved.Name != created.Name {
		t.Errorf("Expected retrieved name '%s', got '%s'", created.Name, retrieved.Name)
	}
}

// TestProjectsCreateMultiple tests creating multiple projects
func TestProjectsCreateMultiple(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	// Clear any existing data first
	if err := ctx.TestHelper.ClearTestData(); err != nil {
		t.Fatalf("Failed to clear test data: %v", err)
	}

	createdIDs := make([]int64, 0, 3)

	// Create 3 projects
	for i := 1; i <= 3; i++ {
		reqBody := map[string]interface{}{
			"name":       "Project " + string(rune('0'+i)),
			"total_page": 100 * i,
			"page":       10 * i,
			"reinicia":   false,
		}
		body, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/v1/projects.json", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		recorder := ctx.MakeRequest(t, req)

		if recorder.Code != http.StatusCreated {
			t.Errorf("Expected status %d for creation %d, got %d", http.StatusCreated, i, recorder.Code)
			t.Errorf("Response body: %s", recorder.Body.String())
		}

		response := ctx.ParseProjectResponse(t, recorder.Body.String())

		if response.ID == 0 {
			t.Errorf("Expected non-zero project ID for project %d", i)
		}

		createdIDs = append(createdIDs, response.ID)
	}

	// Verify all projects can be listed
	req := httptest.NewRequest(http.MethodGet, "/v1/projects.json", nil)
	recorder := ctx.MakeRequest(t, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status %d for GET all, got %d", http.StatusOK, recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}

	response := ctx.ParseProjectResponseArray(t, recorder.Body.String())

	if len(response) != 3 {
		t.Errorf("Expected 3 projects, got %d", len(response))
	}

	// Verify all created IDs are present
	for _, id := range createdIDs {
		found := false
		for _, p := range response {
			if p.ID == id {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Created project ID %d not found in list", id)
		}
	}
}

// TestProjectsCreateConcurrent tests concurrent project creation
func TestProjectsCreateConcurrent(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	done := make(chan bool, 3)
	for i := 0; i < 3; i++ {
		go func(id int) {
			reqBody := map[string]interface{}{
				"name":       "Concurrent Project " + string(rune('0'+id)),
				"total_page": 100,
				"page":       10,
				"reinicia":   false,
			}
			body, err := json.Marshal(reqBody)
			if err != nil {
				done <- false
				return
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/projects.json", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			recorder := ctx.MakeRequest(t, req)

			if recorder.Code == http.StatusCreated {
				done <- true
			} else {
				done <- false
			}
		}(i)
	}

	successCount := 0
	for i := 0; i < 3; i++ {
		if <-done {
			successCount++
		}
	}

	if successCount != 3 {
		t.Errorf("Expected 3 successful creations, got %d", successCount)
	}
}

// TestProjectsCreateValidationErrorFormat tests error response format
func TestProjectsCreateValidationErrorFormat(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	reqBody := map[string]interface{}{
		"name":       "Invalid Project",
		"total_page": 50,
		"page":       100, // page > total_page
		"reinicia":   false,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/projects.json", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := ctx.MakeRequest(t, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	// Verify error response format
	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify top-level error field
	errorField, ok := response["error"].(string)
	if !ok {
		t.Fatalf("Expected 'error' to be a string, got %T", response["error"])
	}
	if errorField != "validation failed" {
		t.Errorf("Expected error 'validation failed', got '%s'", errorField)
	}

	// Verify details field
	details, ok := response["details"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'details' to be a map, got %T", response["details"])
	}

	// Verify page field in details
	if pageErr, exists := details["page"]; exists {
		if errStr, ok := pageErr.(string); ok {
			// Error message should contain validation details
			if len(errStr) == 0 {
				t.Error("Expected non-empty error message for page field")
			}
		} else {
			t.Fatalf("Expected page error to be a string, got %T", pageErr)
		}
	} else {
		t.Error("Expected 'page' field in validation details")
	}
}

// TestProjectsCreateWithNullStartedAt tests project creation with null started_at
func TestProjectsCreateWithNullStartedAt(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	reqBody := map[string]interface{}{
		"name":       "Project with Null StartedAt",
		"total_page": 100,
		"page":       50,
		"started_at": nil, // null started_at
		"reinicia":   false,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/projects.json", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := ctx.MakeRequest(t, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}

	response := ctx.ParseProjectResponse(t, recorder.Body.String())

	if response.Name != "Project with Null StartedAt" {
		t.Errorf("Expected name 'Project with Null StartedAt', got '%s'", response.Name)
	}

	// started_at should be nil when not provided
	if response.StartedAt != nil {
		t.Errorf("Expected started_at to be nil, got %v", *response.StartedAt)
	}
}

// TestProjectsCreateStatusCodeHeaders tests response headers
func TestProjectsCreateStatusCodeHeaders(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	reqBody := map[string]interface{}{
		"name":       "Header Test Project",
		"total_page": 100,
		"page":       50,
		"reinicia":   false,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/projects.json", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := ctx.MakeRequest(t, req)

	// Verify status code
	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	// Verify Content-Type header (should be JSON:API format now)
	contentType := recorder.Header().Get("Content-Type")
	if !contains(contentType, "application/vnd.api+json") {
		t.Errorf("Expected Content-Type to contain 'application/vnd.api+json', got '%s'", contentType)
	}
}

// TestProjectsCreateBadRequestHeaders tests error response headers
func TestProjectsCreateBadRequestHeaders(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	reqBody := map[string]interface{}{
		"name":       "Invalid Project",
		"total_page": 50,
		"page":       100, // page > total_page
		"reinicia":   false,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/projects.json", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := ctx.MakeRequest(t, req)

	// Verify status code
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	// Verify Content-Type header for error response
	contentType := recorder.Header().Get("Content-Type")
	if !contains(contentType, "application/json") {
		t.Errorf("Expected Content-Type to contain 'application/json', got '%s'", contentType)
	}
}
