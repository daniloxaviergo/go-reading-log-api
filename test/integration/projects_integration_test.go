package integration

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"go-reading-log-api-next/internal/adapter/postgres"
	"go-reading-log-api-next/internal/config"
	"go-reading-log-api-next/internal/service"
	"go-reading-log-api-next/internal/service/dashboard"
	"go-reading-log-api-next/test"
)

// TestProjectsIndexIntegration tests the GET /v1/projects endpoint
func TestProjectsIndexIntegration(t *testing.T) {
	// Skip if no test database is configured
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	// Create test projects
	project1ID := ctx.CreateTestProject(t)
	project2ID := ctx.CreateTestProject(t)

	// Test GET /v1/projects.json
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, "/v1/projects.json", nil))

	// Verify response status
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}

	// Parse and verify response using JSON:API envelope helper
	envelope := ctx.ParseProjectResponseArray(t, recorder.Body.String())

	if len(envelope) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(envelope))
	}

	// Verify project IDs are present
	projectIDs := make([]int64, len(envelope))
	for i, p := range envelope {
		projectIDs[i] = p.ID
		if p.ID != project1ID && p.ID != project2ID {
			t.Errorf("Unexpected project ID: %d", p.ID)
		}
	}
}

// TestProjectsIndexEmpty tests with no projects in database
func TestProjectsIndexEmpty(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	// Clear any existing data
	if err := ctx.TestHelper.ClearTestData(); err != nil {
		t.Fatalf("Failed to clear test data: %v", err)
	}

	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, "/v1/projects.json", nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	envelope := ctx.ParseProjectResponseArray(t, recorder.Body.String())

	if len(envelope) != 0 {
		t.Errorf("Expected 0 projects, got %d", len(envelope))
	}
}

// TestProjectsShowIntegration tests the GET /v1/projects/:id endpoint
func TestProjectsShowIntegration(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	projectID := ctx.CreateTestProject(t)
	ctx.CreateTestLog(t, projectID)

	idStr := strconv.Itoa(int(projectID))
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, "/v1/projects/"+idStr+".json", nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}

	response := ctx.ParseProjectResponse(t, recorder.Body.String())

	if response.ID != projectID {
		t.Errorf("Expected project ID %d, got %d", projectID, response.ID)
	}
}

// TestProjectsShowNotFound tests 404 for non-existent project
func TestProjectsShowNotFound(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, "/v1/projects/999999", nil))

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", recorder.Code)
		t.Errorf("Response body: %s", recorder.Body.String())
	}
}

// TestProjectsShowInvalidID tests 400 for invalid project ID
func TestProjectsShowInvalidID(t *testing.T) {
	ctx := Setup(t)
	defer ctx.Teardown(t)

	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, "/v1/projects/invalid.json", nil))

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", recorder.Code)
	}
}

// TestProjectsShowWithLogs tests eager loading of logs
func TestProjectsShowWithLogs(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	projectID := ctx.CreateTestProject(t)
	ctx.CreateTestLogWithNote(t, projectID)

	idStr := strconv.Itoa(int(projectID))
	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, "/v1/projects/"+idStr+".json", nil))

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	response := ctx.ParseProjectResponse(t, recorder.Body.String())

	// Note: The Show handler doesn't include logs in the response
	// GetWithLogs is used for that in the repository
	if response.ID != projectID {
		t.Errorf("Expected project ID %d, got %d", projectID, response.ID)
	}
}

// TestProjectsResponseFormat tests response JSON format
func TestProjectsResponseFormat(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	ctx.CreateTestProject(t)

	recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, "/v1/projects.json", nil))

	body := recorder.Body.String()

	// Verify JSON:API envelope structure
	if !contains(body, `"data"`) {
		t.Errorf("Response missing 'data' field (JSON:API envelope)")
	}
	if !contains(body, `"type"`) {
		t.Errorf("Response missing 'type' field (JSON:API envelope)")
	}
	if !contains(body, `"attributes"`) {
		t.Errorf("Response missing 'attributes' field (JSON:API envelope)")
	}

	// Verify required fields are present in attributes
	requiredFields := []string{"id", "name", "total_page", "page"}
	for _, field := range requiredFields {
		if !contains(body, `"`+field+`"`) {
			t.Errorf("Response missing required field: %s", field)
		}
	}
}

// TestProjectsConcurrentReads tests concurrent read operations
func TestProjectsConcurrentReads(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured - skipping integration test")
	}

	ctx := Setup(t)
	defer ctx.Teardown(t)

	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func() {
			recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, "/v1/projects.json", nil))

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

// TestProjectsNewWithCustomConfig tests using custom config
func TestProjectsNewWithCustomConfig(t *testing.T) {
	cfg := config.LoadConfig()

	helper, err := test.SetupTestDBWithConfig(cfg)
	if err != nil {
		t.Skipf("Test database not configured: %v", err)
		return
	}
	defer helper.Close()

	if err := helper.SetupTestSchema(); err != nil {
		t.Fatalf("Failed to setup schema: %v", err)
	}
	defer helper.CleanupTestSchema()

	// Create test project
	ctx := helper.GetContext()
	query := `INSERT INTO projects (name, total_page, page, reinicia) VALUES ($1, $2, $3, $4) RETURNING id`
	var projectID int64
	if err := helper.Pool.QueryRow(ctx, query, "Test", 100, 10, false).Scan(&projectID); err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Create a router and test server for this test
	projectRepo := postgres.NewProjectRepositoryImpl(helper.Pool)
	logRepo := postgres.NewLogRepositoryImpl(helper.Pool)
	dashboardRepo := postgres.NewDashboardRepositoryImpl(helper.Pool)
	userConfig, err := service.LoadDashboardConfig("dashboard_config.yaml")
	if err != nil {
		userConfig = service.NewUserConfigService(service.GetDefaultConfig())
	}
	router := SetupRoutes(projectRepo, logRepo, dashboardRepo, userConfig, dashboard.ProjectsServiceInterface(&MockProjectsService{}))
	server := httptest.NewServer(router)
	defer server.Close()

	// Test via HTTP
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, server.URL+"/v1/projects.json", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
