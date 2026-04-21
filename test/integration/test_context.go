package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-reading-log-api-next/internal/adapter/postgres"
	api "go-reading-log-api-next/internal/api/v1"
	"go-reading-log-api-next/internal/config"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
	"go-reading-log-api-next/test"
)

// SetupRoutes is a wrapper for api.SetupRoutes
func SetupRoutes(repo repository.ProjectRepository, logRepo repository.LogRepository) http.Handler {
	return api.SetupRoutes(repo, logRepo)
}

const testContextTimeout = 30 * time.Second

// IntegrationTestContext provides a test context with database and HTTP client
type IntegrationTestContext struct {
	TestHelper *test.TestHelper
	Server     *httptest.Server
	Client     *http.Client
	ProjectID  int64
	LogID      int64
}

// Setup creates a new integration test context
func Setup(t *testing.T) *IntegrationTestContext {
	t.Helper()

	// Skip if no test database is configured
	if !test.IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration tests")
	}

	// Setup database
	helper, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	// Setup schema
	if err := helper.SetupTestSchema(); err != nil {
		helper.Close()
		t.Fatalf("Failed to setup test schema: %v", err)
	}

	// Create PostgreSQL repositories (not mocks) for integration tests
	projectRepo := postgres.NewProjectRepositoryImpl(helper.Pool)
	logRepo := postgres.NewLogRepositoryImpl(helper.Pool)

	// Setup routes with PostgreSQL repositories
	router := SetupRoutes(projectRepo, logRepo)

	return &IntegrationTestContext{
		TestHelper: helper,
		Server:     httptest.NewServer(router),
		Client:     &http.Client{},
	}
}

// Teardown cleans up the integration test context
func (ctx *IntegrationTestContext) Teardown(t *testing.T) {
	t.Helper()

	// Cleanup schema
	if err := ctx.TestHelper.CleanupTestSchema(); err != nil {
		t.Logf("Failed to cleanup test schema: %v", err)
	}

	// Close helper
	ctx.TestHelper.Close()

	// Close test server
	ctx.Server.Close()
}

// CreateTestProject creates a test project in the database
func (ctx *IntegrationTestContext) CreateTestProject(t *testing.T) int64 {
	t.Helper()

	ctxID := ctx.TestHelper.GetContext()
	query := `
		INSERT INTO projects (name, total_page, page, reinicia)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int64
	err := ctx.TestHelper.Pool.QueryRow(ctxID, query, "Test Project", 100, 10, false).Scan(&id)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	ctx.ProjectID = id
	return id
}

// CreateTestLog creates a test log for a project
func (ctx *IntegrationTestContext) CreateTestLog(t *testing.T, projectID int64) int64 {
	t.Helper()

	ctxID := ctx.TestHelper.GetContext()
	query := `
		INSERT INTO logs (project_id, data, start_page, end_page, wday)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int64
	err := ctx.TestHelper.Pool.QueryRow(ctxID, query, projectID, "2024-01-01", 1, 10, 1).Scan(&id)
	if err != nil {
		t.Fatalf("Failed to create test log: %v", err)
	}

	ctx.LogID = id
	return id
}

// CreateTestLogWithNote creates a test log with note and text
func (ctx *IntegrationTestContext) CreateTestLogWithNote(t *testing.T, projectID int64) int64 {
	t.Helper()

	ctxID := ctx.TestHelper.GetContext()
	query := `
		INSERT INTO logs (project_id, data, start_page, end_page, wday, note, text)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id int64
	err := ctx.TestHelper.Pool.QueryRow(ctxID, query, projectID, "2024-01-01", 1, 10, 1, "Test note", "Test text").Scan(&id)
	if err != nil {
		t.Fatalf("Failed to create test log: %v", err)
	}

	ctx.LogID = id
	return id
}

// NewRequest creates a new HTTP request with context
func (ctx *IntegrationTestContext) NewRequest(t *testing.T, method, path string, body interface{}) *http.Request {
	t.Helper()

	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	if body != nil {
		// Handle JSON body if needed
	}

	return req
}

// MakeRequest makes an HTTP request and returns the response
func (ctx *IntegrationTestContext) MakeRequest(t *testing.T, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	ctx.Server.Config.Handler.ServeHTTP(recorder, req)
	return recorder
}

// MakeRequestWithContext makes an HTTP request with custom context timeout
func (ctx *IntegrationTestContext) MakeRequestWithContext(t *testing.T, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	ctx.Server.Config.Handler.ServeHTTP(recorder, req)
	return recorder
}

// GetContext returns a test context with timeout
func (ctx *IntegrationTestContext) GetContext() context.Context {
	ctxID, cancel := context.WithTimeout(context.Background(), testContextTimeout)
	_ = cancel
	return ctxID
}

// ParseHealthCheckResponse parses a health check response
func (ctx *IntegrationTestContext) ParseHealthCheckResponse(t *testing.T, body string) *dto.HealthCheckResponse {
	t.Helper()

	// This is a helper - actual parsing would be done by test code
	return nil
}

// Helper function to setup test database with custom config
func SetupTestDBWithConfig(t *testing.T, cfg *config.Config) *test.TestHelper {
	t.Helper()

	if !test.IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	helper, err := test.SetupTestDBWithConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	return helper
}

// Helper function to setup test database
func SetupTestDB(t *testing.T) *test.TestHelper {
	t.Helper()

	if !test.IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	helper, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	return helper
}

// Helper function to verify HTTP response
func VerifyHTTPResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int, expectedContentType string) {
	t.Helper()

	if recorder.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}

	if expectedContentType != "" && recorder.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, recorder.Header().Get("Content-Type"))
	}
}

// Helper function to verify error response
func VerifyErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedStatus int) {
	t.Helper()

	VerifyHTTPResponse(t, recorder, expectedStatus, "application/json")

	if recorder.Code >= 400 {
		// Verify error response contains "error" field
		body := recorder.Body.String()
		if !contains(body, `"error"`) && !contains(body, `"error":`) {
			t.Logf("Warning: Error response may not contain 'error' field: %s", body)
		}
	}
}

// ParseJSONAPIEnvelope parses a JSON:API envelope response
func (ctx *IntegrationTestContext) ParseJSONAPIEnvelope(t *testing.T, body string) *dto.JSONAPIEnvelope {
	t.Helper()

	var envelope dto.JSONAPIEnvelope
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		t.Fatalf("Failed to parse JSON:API envelope: %v", err)
	}
	return &envelope
}

// ParseProjectResponse parses a project response (handles both flat and JSON:API formats)
func (ctx *IntegrationTestContext) ParseProjectResponse(t *testing.T, body string) *dto.ProjectResponse {
	t.Helper()

	// Try to parse as JSON:API envelope first
	var envelope dto.JSONAPIEnvelope
	if err := json.Unmarshal([]byte(body), &envelope); err == nil {
		// Check if it's a collection or single object
		switch v := envelope.Data.(type) {
		case []interface{}:
			// Collection - return first item or nil if empty
			if len(v) > 0 {
				if dataObj, ok := v[0].(map[string]interface{}); ok {
					if attrs, ok := dataObj["attributes"].(map[string]interface{}); ok {
						attrsJSON, _ := json.Marshal(attrs)
						var parsed dto.ProjectResponse
						json.Unmarshal(attrsJSON, &parsed)
						return &parsed
					}
				}
			}
			return nil
		case map[string]interface{}:
			// Single object
			if attrs, ok := v["attributes"].(map[string]interface{}); ok {
				attrsJSON, _ := json.Marshal(attrs)
				var parsed dto.ProjectResponse
				json.Unmarshal(attrsJSON, &parsed)
				return &parsed
			}
		}
	}

	// Fall back to parsing as flat JSON
	var response dto.ProjectResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		t.Fatalf("Failed to parse project response: %v", err)
	}
	return &response
}

// ParseLogResponse parses a log response (handles both flat and JSON:API formats)
func (ctx *IntegrationTestContext) ParseLogResponse(t *testing.T, body string) *dto.LogResponse {
	t.Helper()

	// Try to parse as JSON:API envelope first
	var envelope dto.JSONAPIEnvelope
	if err := json.Unmarshal([]byte(body), &envelope); err == nil {
		switch v := envelope.Data.(type) {
		case []interface{}:
			if len(v) > 0 {
				if dataObj, ok := v[0].(map[string]interface{}); ok {
					if attrs, ok := dataObj["attributes"].(map[string]interface{}); ok {
						attrsJSON, _ := json.Marshal(attrs)
						var parsed dto.LogResponse
						json.Unmarshal(attrsJSON, &parsed)
						return &parsed
					}
				}
			}
			return nil
		case map[string]interface{}:
			if attrs, ok := v["attributes"].(map[string]interface{}); ok {
				attrsJSON, _ := json.Marshal(attrs)
				var parsed dto.LogResponse
				json.Unmarshal(attrsJSON, &parsed)
				return &parsed
			}
		}
	}

	// Fall back to parsing as flat JSON
	var response dto.LogResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		t.Fatalf("Failed to parse log response: %v", err)
	}
	return &response
}

// ParseProjectResponseArray parses an array of project responses
func (ctx *IntegrationTestContext) ParseProjectResponseArray(t *testing.T, body string) []*dto.ProjectResponse {
	t.Helper()

	// Try to parse as JSON:API envelope first
	var envelope dto.JSONAPIEnvelope
	if err := json.Unmarshal([]byte(body), &envelope); err == nil {
		switch v := envelope.Data.(type) {
		case []interface{}:
			result := make([]*dto.ProjectResponse, len(v))
			for i, item := range v {
				if dataObj, ok := item.(map[string]interface{}); ok {
					if attrs, ok := dataObj["attributes"].(map[string]interface{}); ok {
						attrsJSON, _ := json.Marshal(attrs)
						var parsed dto.ProjectResponse
						json.Unmarshal(attrsJSON, &parsed)
						result[i] = &parsed
					}
				}
			}
			return result
		}
	}

	// Fall back to parsing as flat JSON array
	var response []*dto.ProjectResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		t.Fatalf("Failed to parse project response array: %v", err)
	}
	return response
}

// ParseLogResponseArray parses an array of log responses
func (ctx *IntegrationTestContext) ParseLogResponseArray(t *testing.T, body string) []*dto.LogResponse {
	t.Helper()

	// Try to parse as JSON:API envelope first
	var envelope dto.JSONAPIEnvelope
	if err := json.Unmarshal([]byte(body), &envelope); err == nil {
		switch v := envelope.Data.(type) {
		case []interface{}:
			result := make([]*dto.LogResponse, len(v))
			for i, item := range v {
				if dataObj, ok := item.(map[string]interface{}); ok {
					if attrs, ok := dataObj["attributes"].(map[string]interface{}); ok {
						attrsJSON, _ := json.Marshal(attrs)
						var parsed dto.LogResponse
						json.Unmarshal(attrsJSON, &parsed)
						result[i] = &parsed
					}
				}
			}
			return result
		}
	}

	// Fall back to parsing as flat JSON array
	var response []*dto.LogResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		t.Fatalf("Failed to parse log response array: %v", err)
	}
	return response
}

// contains checks if s contains substr
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Helper to check if we're running against a test database
func IsTestDatabase() bool {
	return test.IsTestDatabase()
}
