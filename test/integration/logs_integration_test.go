package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-reading-log-api-next/internal/domain/dto"
)

// TestLogsIndexIntegration tests the GET /api/v1/projects/:project_id/logs endpoint
func TestLogsIndexIntegration(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	projectID := ctx.CreateTestProject(t)
	ctx.CreateTestLog(t, projectID)
	ctx.CreateTestLog(t, projectID)

	// Make HTTP request to the test server
	url := "/api/v1/projects/" + string(rune(projectID+'0')) + "/logs"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}

	var response []*dto.LogResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Should return first 4 logs (we created 2)
	if len(response) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(response))
	}
}

// TestLogsIndexEmpty tests with no logs for project
func TestLogsIndexEmpty(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	projectID := ctx.CreateTestProject(t)

	// Make HTTP request to the test server
	url := "/api/v1/projects/" + string(rune(projectID+'0')) + "/logs"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	var response []*dto.LogResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(response) != 0 {
		t.Errorf("Expected 0 logs, got %d", len(response))
	}
}

// TestLogsIndexProjectNotFound tests 404 for non-existent project
func TestLogsIndexProjectNotFound(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	// Make HTTP request to the test server
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, "/api/v1/projects/999999/logs", nil))

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}
}

// TestLogsIndexInvalidProjectID tests 400 for invalid project ID
func TestLogsIndexInvalidProjectID(t *testing.T) {
	ctx := Setup(t)
	defer ctx.Teardown(t)

	// Make HTTP request to the test server
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, "/api/v1/projects/invalid/logs", nil))

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", recorder.Code)
	}
}

// TestLogsIndexLimit tests that only first 4 logs are returned
func TestLogsIndexLimit(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	projectID := ctx.CreateTestProject(t)

	// Create more than 4 logs
	for i := 0; i < 6; i++ {
		ctx.CreateTestLog(t, projectID)
	}

	// Make HTTP request to the test server
	url := "/api/v1/projects/" + string(rune(projectID+'0')) + "/logs"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	var response []*dto.LogResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Should be limited to 4
	if len(response) != 4 {
		t.Errorf("Expected 4 logs (limited), got %d", len(response))
	}
}

// TestLogsIndexWithLogs tests eager loading of logs with data
func TestLogsIndexWithLogs(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	projectID := ctx.CreateTestProject(t)
	ctx.CreateTestLogWithNote(t, projectID)

	// Make HTTP request to the test server
	url := "/api/v1/projects/" + string(rune(projectID+'0')) + "/logs"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	var response []*dto.LogResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("Expected 1 log, got %d", len(response))
	}

	// Verify log data is present
	if response[0].Data == nil {
		t.Error("Log data should not be nil")
	}
}

// TestLogsIndexConcurrent tests concurrent read operations
func TestLogsIndexConcurrent(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	projectID := ctx.CreateTestProject(t)

	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func() {
			url := "/api/v1/projects/" + string(rune(projectID+'0')) + "/logs"
			recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

			if recorder.Code == http.StatusOK {
				done <- true
			} else {
				done <- false
			}
		}()
	}

	successCount := 0
	for i := 0; i < 5; i++ {
		if <-done {
			successCount++
		}
	}

	if successCount != 5 {
		t.Errorf("Expected 5 successful requests, got %d", successCount)
	}
}

// TestLogsIndexResponseFormat tests response JSON format
func TestLogsIndexResponseFormat(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	projectID := ctx.CreateTestProject(t)
	ctx.CreateTestLog(t, projectID)

	// Make HTTP request to the test server
	url := "/api/v1/projects/" + string(rune(projectID+'0')) + "/logs"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	body := recorder.Body.String()

	// Verify required fields are present
	requiredFields := []string{"id", "start_page", "end_page"}
	for _, field := range requiredFields {
		if !contains(body, `"`+field+`"`) {
			t.Errorf("Response missing required field: %s", field)
		}
	}

	// Verify project is eager-loaded
	if !contains(body, `"project"`) {
		t.Error("Response should include eager-loaded project")
	}
}
