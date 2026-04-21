package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"go-reading-log-api-next/internal/domain/dto"
)

// TestLogsIndexIntegration tests the GET /v1/projects/:project_id/logs.json endpoint
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
	url := "/v1/projects/" + strconv.Itoa(int(projectID)) + "/logs.json"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}

	logs := ctx.ParseLogResponseArray(t, recorder.Body.String())

	// Should return first 4 logs (we created 2)
	if len(logs) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(logs))
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
	url := "/v1/projects/" + strconv.Itoa(int(projectID)) + "/logs.json"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	logs := ctx.ParseLogResponseArray(t, recorder.Body.String())

	if len(logs) != 0 {
		t.Errorf("Expected 0 logs, got %d", len(logs))
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
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, "/v1/projects/999999/logs.json", nil))

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
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, "/v1/projects/invalid/logs.json", nil))

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
	url := "/v1/projects/" + strconv.Itoa(int(projectID)) + "/logs.json"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	logs := ctx.ParseLogResponseArray(t, recorder.Body.String())

	// Should be limited to 4
	if len(logs) != 4 {
		t.Errorf("Expected 4 logs (limited), got %d", len(logs))
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
	url := "/v1/projects/" + strconv.Itoa(int(projectID)) + "/logs.json"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	logs := ctx.ParseLogResponseArray(t, recorder.Body.String())

	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}

	// Verify log data is present
	if logs[0].Data == nil {
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
			url := "/v1/projects/" + strconv.Itoa(int(projectID)) + "/logs.json"
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
	url := "/v1/projects/" + strconv.Itoa(int(projectID)) + "/logs.json"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	body := recorder.Body.String()

	// Verify required fields are present
	requiredFields := []string{"id", "start_page", "end_page"}
	for _, field := range requiredFields {
		if !contains(body, `"`+field+`"`) {
			t.Errorf("Response missing required field: %s", field)
		}
	}

	// Verify project is in relationships (not embedded in attributes)
	if !contains(body, `"relationships"`) {
		t.Error("Response should include relationships")
	}

	if !contains(body, `"project"`) {
		t.Error("Response should include project relationship")
	}

	// Verify JSON:API content type
	if recorder.Header().Get("Content-Type") != "application/vnd.api+json" {
		t.Errorf("Expected Content-Type 'application/vnd.api+json', got '%s'", recorder.Header().Get("Content-Type"))
	}
}

// TestLogsIndexJSONAPIStructure tests the JSON:API envelope structure
func TestLogsIndexJSONAPIStructure(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	projectID := ctx.CreateTestProject(t)
	ctx.CreateTestLog(t, projectID)

	// Make HTTP request to the test server
	url := "/v1/projects/" + strconv.Itoa(int(projectID)) + "/logs.json"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}

	body := recorder.Body.String()

	// Validate JSON:API schema structure
	ctx.ValidateJSONAPIStructure(t, body)

	// Verify envelope contains data array
	var envelope dto.JSONAPIEnvelope
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		t.Fatalf("Failed to parse JSON:API envelope: %v", err)
	}

	// Check that data is an array
	if dataArray, ok := envelope.Data.([]interface{}); !ok {
		t.Error("Response 'data' field should be an array")
	} else if len(dataArray) == 0 {
		t.Error("Response 'data' array should not be empty")
	} else {
		// Check first item has required JSON:API fields
		if dataObj, ok := dataArray[0].(map[string]interface{}); ok {
			requiredFields := []string{"type", "id", "attributes"}
			for _, field := range requiredFields {
				if _, exists := dataObj[field]; !exists {
					t.Errorf("Data item missing required field: %s", field)
				}
			}

			// Verify attributes contains expected log fields
			if attrs, ok := dataObj["attributes"].(map[string]interface{}); ok {
				logFields := []string{"data", "start_page", "end_page", "note"}
				for _, field := range logFields {
					if _, exists := attrs[field]; !exists {
						t.Errorf("Attributes missing required field: %s", field)
					}
				}
			}

			// Verify relationships exist
			if rel, exists := dataObj["relationships"]; !exists || rel == nil {
				t.Error("Data item missing relationships")
			}
		}
	}

	// Verify included array exists and contains project
	if len(envelope.Included) == 0 {
		t.Error("Response 'included' array should not be empty (project should be included)")
	} else {
		// Check that included resource is a project
		if includedObj, ok := envelope.Included[0].(map[string]interface{}); ok {
			if includedObj["type"] != "projects" {
				t.Errorf("Expected included type 'projects', got '%v'", includedObj["type"])
			}
		}
	}
}

// TestLogsIndexRFC3339DateFormat tests that dates are in RFC3339 format
func TestLogsIndexRFC3339DateFormat(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	projectID := ctx.CreateTestProject(t)
	ctx.CreateTestLog(t, projectID)

	// Make HTTP request to the test server
	url := "/v1/projects/" + strconv.Itoa(int(projectID)) + "/logs.json"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	body := recorder.Body.String()

	// Parse the envelope
	var envelope dto.JSONAPIEnvelope
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		t.Fatalf("Failed to parse JSON:API envelope: %v", err)
	}

	// Check dates in data array
	if dataArray, ok := envelope.Data.([]interface{}); ok {
		for i, item := range dataArray {
			if dataObj, ok := item.(map[string]interface{}); ok {
				if attrs, ok := dataObj["attributes"].(map[string]interface{}); ok {
					if dateVal, exists := attrs["data"]; exists {
						dateStr, ok := dateVal.(string)
						if !ok {
							t.Errorf("Item %d: 'data' field is not a string", i)
							continue
						}

						// Verify RFC3339 format
						ctx.VerifyRFC3339Date(t, dateStr)

						// Also verify it can be parsed as time.Time
						parsedTime, err := time.Parse(time.RFC3339, dateStr)
						if err != nil {
							t.Errorf("Item %d: Failed to parse date '%s' as RFC3339: %v", i, dateStr, err)
						} else if parsedTime.IsZero() {
							t.Errorf("Item %d: Parsed time is zero", i)
						}
					} else {
						t.Errorf("Item %d: Missing 'data' field in attributes", i)
					}
				}
			}
		}
	}
}

// TestLogsIndexPayloadSize tests that response payload size is within acceptable limits
func TestLogsIndexPayloadSize(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	projectID := ctx.CreateTestProject(t)

	// Create multiple logs to test size with more data
	for i := 0; i < 3; i++ {
		ctx.CreateTestLog(t, projectID)
	}

	// Make HTTP request to the test server
	url := "/v1/projects/" + strconv.Itoa(int(projectID)) + "/logs.json"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	body := recorder.Body.String()
	size := ctx.CalculatePayloadSize(t, body)

	// Define acceptable size limit (approximate, based on expected response structure)
	// A single log with project included should be under 5KB
	const maxPayloadSize = 5000 // 5KB

	if size > maxPayloadSize {
		t.Errorf("Payload size (%d bytes) exceeds acceptable limit (%d bytes)", size, maxPayloadSize)
	}

	// Log the size for reference
	t.Logf("Response payload size: %d bytes", size)
}

// TestLogsIndexRelationshipReference tests that project is referenced, not embedded
func TestLogsIndexRelationshipReference(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	projectID := ctx.CreateTestProject(t)
	ctx.CreateTestLog(t, projectID)

	// Make HTTP request to the test server
	url := "/v1/projects/" + strconv.Itoa(int(projectID)) + "/logs.json"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	body := recorder.Body.String()

	// Parse the envelope
	var envelope dto.JSONAPIEnvelope
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		t.Fatalf("Failed to parse JSON:API envelope: %v", err)
	}

	// Extract relationships from first data item
	relationships := ctx.ExtractRelationships(t, body)

	// Verify project relationship exists
	if projRel, exists := relationships["project"]; !exists {
		t.Error("Response missing 'project' relationship")
	} else if projRel == nil {
		t.Error("Project relationship is nil")
	}

	// Verify project data is in included array (not embedded in attributes)
	included := ctx.ExtractIncludedResources(t, body)

	if len(included) == 0 {
		t.Error("No resources in 'included' array - project should be included")
	} else {
		// Check that at least one included resource is a project
		foundProject := false
		for _, inc := range included {
			if incObj, ok := inc.(map[string]interface{}); ok {
				if incObj["type"] == "projects" {
					foundProject = true
					break
				}
			}
		}
		if !foundProject {
			t.Error("No project found in 'included' array")
		}
	}

	// Verify attributes do NOT contain embedded project object
	// (project should only be in relationships/included)
	if dataArray, ok := envelope.Data.([]interface{}); ok && len(dataArray) > 0 {
		if dataObj, ok := dataArray[0].(map[string]interface{}); ok {
			if attrs, ok := dataObj["attributes"].(map[string]interface{}); ok {
				// Attributes should not have a 'project' key
				// (project is in relationships, not attributes)
				if _, hasProjectInAttrs := attrs["project"]; hasProjectInAttrs {
					t.Error("Project should not be embedded in attributes (should be in relationships)")
				}
			}
		}
	}
}

// TestLogsIndexEmptyJSONAPIStructure tests JSON:API structure for empty results
func TestLogsIndexEmptyJSONAPIStructure(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	projectID := ctx.CreateTestProject(t)
	// Don't create any logs

	// Make HTTP request to the test server
	url := "/v1/projects/" + strconv.Itoa(int(projectID)) + "/logs.json"
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, url, nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	body := recorder.Body.String()

	// Validate JSON:API schema structure even for empty results
	ctx.ValidateJSONAPIStructure(t, body)

	// Parse the envelope
	var envelope dto.JSONAPIEnvelope
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		t.Fatalf("Failed to parse JSON:API envelope: %v", err)
	}

	// Data should be an empty array, not null
	if dataArray, ok := envelope.Data.([]interface{}); !ok {
		t.Error("Response 'data' field should be an array (even if empty)")
	} else if len(dataArray) != 0 {
		t.Errorf("Expected 0 logs, got %d", len(dataArray))
	}

	// Included is optional in JSON:API spec, but if present should not be nil
	// When no logs exist, included might be omitted or empty
	if envelope.Included != nil && len(envelope.Included) > 0 {
		// If included has items, verify they are projects
		for _, inc := range envelope.Included {
			if incObj, ok := inc.(map[string]interface{}); ok {
				if incObj["type"] != "projects" {
					t.Errorf("Expected included type 'projects', got '%v'", incObj["type"])
				}
			}
		}
	}
	// Note: We don't fail if Included is nil or empty because when no logs exist,
	// the project data might not be included (optimization)
}
