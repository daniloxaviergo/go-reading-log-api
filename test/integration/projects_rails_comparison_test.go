package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-reading-log-api-next/internal/adapter/postgres"
	"go-reading-log-api-next/internal/service"
	"go-reading-log-api-next/internal/service/dashboard"
)

// ============================================================================
// Rails Comparison Tests for Projects and Logs Endpoints
// ============================================================================
// These tests compare Go API responses with the legacy Rails API to ensure
// feature parity during the migration. Tests query both APIs via HTTP and
// verify that responses match for all calculated and computed fields.
//
// Prerequisites:
// - RAILS_API_URL environment variable must be set (e.g., http://localhost:3001)
// - Both Go and Rails databases must have identical fixture data
// - PostgreSQL test database must be configured
//
// Execution:
//   RAILS_API_URL=http://localhost:3001 go test -v ./test/integration/... -run ".*Comparison.*"
//
// Tolerance Rules:
// - Floating-point fields (progress, median_day): ±0.01
// - days_unreading: ±1 day (timezone calculation differences)
// - All other fields: Exact match required
// ============================================================================

// ProjectsRailsComparisonTest represents a comparison test case for projects endpoints
type ProjectsRailsComparisonTest struct {
	testName     string
	endpoint     string
	method       string
	setupProject func(t *testing.T, ctx *IntegrationTestContext) // Setup function for test data
	validate     func(t *testing.T, goResponse interface{}, railsResponse interface{})
	skipIfEmpty  bool // Skip test if no projects exist in database
}

// Run executes a single comparison test
func (test *ProjectsRailsComparisonTest) Run(t *testing.T) {
	t.Run(test.testName, func(t *testing.T) {
		// Check if Rails API URL is configured
		railsURL := os.Getenv("RAILS_API_URL")
		if railsURL == "" {
			t.Skip("RAILS_API_URL not set - skipping Rails comparison test")
		}

		// Check if test database is configured
		if !IsTestDatabase() {
			t.Skip("Test database not configured - skipping Rails comparison test")
		}

		// Setup test database
		helper := SetupTestDB(t)
		defer helper.Close()

		if err := helper.SetupTestSchema(); err != nil {
			t.Fatalf("Failed to setup test schema: %v", err)
		}

		// Create Go handler with PostgreSQL repositories
		projectRepo := postgres.NewProjectRepositoryImpl(helper.Pool)
		logRepo := postgres.NewLogRepositoryImpl(helper.Pool)
		dashboardRepo := postgres.NewDashboardRepositoryImpl(helper.Pool)
		userConfig, err := service.LoadDashboardConfig("dashboard_config.yaml")
		if err != nil {
			userConfig = service.NewUserConfigService(service.GetDefaultConfig())
		}
		goHandler := SetupRoutes(projectRepo, logRepo, dashboardRepo, userConfig, dashboard.ProjectsServiceInterface(&MockProjectsService{}))

		// Setup test data via fixture function
		if test.setupProject != nil {
			ctx := &IntegrationTestContext{
				TestHelper: helper,
				Server:     httptest.NewServer(goHandler),
				Client:     &http.Client{},
			}
			defer ctx.Server.Close()
			test.setupProject(t, ctx)
		}

		// Make request to Go API
		goResponse := makeRequest(t, goHandler, test.method, test.endpoint)

		// Fetch Rails API response
		railsResponse := fetchRailsAPI(t, railsURL+test.endpoint)

		// Compare responses
		test.validate(t, goResponse, railsResponse)
	})
}

// ============================================================================
// Test Data Setup Functions
// ============================================================================

// setupEmptyDatabase creates an empty database (no projects)
func setupEmptyDatabase(t *testing.T, ctx *IntegrationTestContext) {
	// Database is already empty after SetupTestDB
}

// setupSingleProject creates a single project with basic data
func setupSingleProject(t *testing.T, ctx *IntegrationTestContext) {
	// Create a server for this context if not already created
	if ctx.Server == nil {
		projectRepo := postgres.NewProjectRepositoryImpl(ctx.TestHelper.Pool)
		logRepo := postgres.NewLogRepositoryImpl(ctx.TestHelper.Pool)
		dashboardRepo := postgres.NewDashboardRepositoryImpl(ctx.TestHelper.Pool)
		userConfig, err := service.LoadDashboardConfig("dashboard_config.yaml")
		if err != nil {
			userConfig = service.NewUserConfigService(service.GetDefaultConfig())
		}
		ctx.Server = httptest.NewServer(SetupRoutes(projectRepo, logRepo, dashboardRepo, userConfig, dashboard.ProjectsServiceInterface(&MockProjectsService{})))
	}
	projectID := ctx.CreateTestProject(t)
	ctx.CreateTestLog(t, projectID)
}

// setupMultipleProjects creates multiple projects with varying data
func setupMultipleProjects(t *testing.T, ctx *IntegrationTestContext) {
	// Create a server for this context if not already created
	if ctx.Server == nil {
		projectRepo := postgres.NewProjectRepositoryImpl(ctx.TestHelper.Pool)
		logRepo := postgres.NewLogRepositoryImpl(ctx.TestHelper.Pool)
		dashboardRepo := postgres.NewDashboardRepositoryImpl(ctx.TestHelper.Pool)
		userConfig, err := service.LoadDashboardConfig("dashboard_config.yaml")
		if err != nil {
			userConfig = service.NewUserConfigService(service.GetDefaultConfig())
		}
		ctx.Server = httptest.NewServer(SetupRoutes(projectRepo, logRepo, dashboardRepo, userConfig, dashboard.ProjectsServiceInterface(&MockProjectsService{})))
	}
	// Project 1: Basic project with logs
	project1ID := ctx.CreateTestProject(t)
	ctx.CreateTestLog(t, project1ID)
	ctx.CreateTestLog(t, project1ID)

	// Project 2: Project with more logs (tests limit of 4)
	project2ID := ctx.CreateTestProject(t)
	for i := 0; i < 6; i++ {
		ctx.CreateTestLog(t, project2ID)
	}

	// Project 3: Project with no logs
	ctx.CreateTestProject(t)
}

// setupProjectWithCustomDates creates a project with specific dates for testing
func setupProjectWithCustomDates(t *testing.T, ctx *IntegrationTestContext) {
	// Create a server for this context if not already created
	if ctx.Server == nil {
		projectRepo := postgres.NewProjectRepositoryImpl(ctx.TestHelper.Pool)
		logRepo := postgres.NewLogRepositoryImpl(ctx.TestHelper.Pool)
		dashboardRepo := postgres.NewDashboardRepositoryImpl(ctx.TestHelper.Pool)
		userConfig, err := service.LoadDashboardConfig("dashboard_config.yaml")
		if err != nil {
			userConfig = service.NewUserConfigService(service.GetDefaultConfig())
		}
		ctx.Server = httptest.NewServer(SetupRoutes(projectRepo, logRepo, dashboardRepo, userConfig, dashboard.ProjectsServiceInterface(&MockProjectsService{})))
	}
	// Create project with specific started_at
	ctxID := ctx.TestHelper.GetContext()
	query := `
		INSERT INTO projects (id, name, total_page, page, started_at, reinicia)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var projectID int64
	startedAt := "2024-01-15T10:30:00Z"
	err := ctx.TestHelper.Pool.QueryRow(ctxID, query, 1, "Test Project", 200, 50, startedAt, false).Scan(&projectID)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	// Create logs with specific dates
	logQuery := `
		INSERT INTO logs (id, project_id, data, start_page, end_page, wday)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	logDates := []string{
		"2024-01-15T10:30:00",
		"2024-01-16T09:00:00",
		"2024-01-18T14:00:00",
	}
	for i, logDate := range logDates {
		_, err := ctx.TestHelper.Pool.Exec(ctxID, logQuery, i+1, projectID, logDate, i*25, (i+1)*25, 1)
		if err != nil {
			t.Fatalf("Failed to create test log: %v", err)
		}
	}
}

// ============================================================================
// Response Fetching Functions
// ============================================================================

// ============================================================================
// Validation Functions
// ============================================================================

// validateProjectsIndex validates the projects index endpoint comparison
func validateProjectsIndex(t *testing.T, goResponse interface{}, railsResponse interface{}) {
	goBody, ok := goResponse.([]byte)
	require.True(t, ok, "Go response should be []byte")

	railsBody, ok := railsResponse.([]byte)
	require.True(t, ok, "Rails response should be []byte")

	// Parse both responses
	var goProjects []map[string]interface{}
	var railsProjects []map[string]interface{}

	// Handle JSON:API envelope format
	goEnvelope := parseProjectsJSONAPIEnvelope(t, goBody)
	railsEnvelope := parseProjectsJSONAPIEnvelope(t, railsBody)

	// Extract data arrays
	goProjects = extractProjectsDataArray(t, goEnvelope)
	railsProjects = extractProjectsDataArray(t, railsEnvelope)

	// Compare counts
	assert.Equal(t, len(railsProjects), len(goProjects),
		"Project count mismatch: Go=%d, Rails=%d", len(goProjects), len(railsProjects))

	if len(goProjects) == 0 {
		// Empty database - test passed
		return
	}

	// Compare each project
	for i := range goProjects {
		goProj := goProjects[i]
		railsProj := railsProjects[i]

		// Compare basic fields (exact match)
		assert.Equal(t, railsProj["name"], goProj["name"], "Project %d: name mismatch", i)
		assert.Equal(t, railsProj["total_page"], goProj["total_page"], "Project %d: total_page mismatch", i)
		assert.Equal(t, railsProj["page"], goProj["page"], "Project %d: page mismatch", i)

		// Compare calculated fields with tolerance
		if goProg, ok := goProj["progress"].(float64); ok {
			if railsProg, ok := railsProj["progress"].(float64); ok {
				assert.InDelta(t, railsProg, goProg, 0.01,
					"Project %d: progress mismatch", i)
			}
		}

		// Compare logs_count (exact match)
		if goLogsCount, ok := goProj["logs_count"].(float64); ok {
			if railsLogsCount, ok := railsProj["logs_count"].(float64); ok {
				assert.Equal(t, int(railsLogsCount), int(goLogsCount),
					"Project %d: logs_count mismatch", i)
			}
		}

		// Compare days_unreading with 1-day tolerance
		if goDays, ok := goProj["days_unreading"].(float64); ok {
			if railsDays, ok := railsProj["days_unreading"].(float64); ok {
				assert.InDelta(t, railsDays, goDays, 1,
					"Project %d: days_unreading exceeds 1-day tolerance", i)
			}
		}

		// Compare median_day with tolerance
		if goMedian, ok := goProj["median_day"].(float64); ok {
			if railsMedian, ok := railsProj["median_day"].(float64); ok {
				assert.InDelta(t, railsMedian, goMedian, 0.01,
					"Project %d: median_day mismatch", i)
			}
		}

		// Compare status (exact match)
		assert.Equal(t, railsProj["status"], goProj["status"],
			"Project %d: status mismatch", i)
	}
}

// validateProjectsShow validates the projects show endpoint comparison
func validateProjectsShow(t *testing.T, goResponse interface{}, railsResponse interface{}) {
	goBody, ok := goResponse.([]byte)
	require.True(t, ok, "Go response should be []byte")

	railsBody, ok := railsResponse.([]byte)
	require.True(t, ok, "Rails response should be []byte")

	// Parse both responses
	goEnvelope := parseProjectsJSONAPIEnvelope(t, goBody)
	railsEnvelope := parseProjectsJSONAPIEnvelope(t, railsBody)

	// Extract single objects
	goProj := extractProjectsSingleObject(t, goEnvelope)
	railsProj := extractProjectsSingleObject(t, railsEnvelope)

	// Compare basic fields
	assert.Equal(t, railsProj["name"], goProj["name"], "name mismatch")
	assert.Equal(t, railsProj["total_page"], goProj["total_page"], "total_page mismatch")
	assert.Equal(t, railsProj["page"], goProj["page"], "page mismatch")

	// Compare calculated fields with tolerance
	if goProg, ok := goProj["progress"].(float64); ok {
		if railsProg, ok := railsProj["progress"].(float64); ok {
			assert.InDelta(t, railsProg, goProg, 0.01, "progress mismatch")
		}
	}

	if goDays, ok := goProj["days_unreading"].(float64); ok {
		if railsDays, ok := railsProj["days_unreading"].(float64); ok {
			assert.InDelta(t, railsDays, goDays, 1, "days_unreading exceeds tolerance")
		}
	}

	// Compare status
	assert.Equal(t, railsProj["status"], goProj["status"], "status mismatch")
}

// validateLogsIndex validates the logs index endpoint comparison
func validateLogsIndex(t *testing.T, goResponse interface{}, railsResponse interface{}) {
	goBody, ok := goResponse.([]byte)
	require.True(t, ok, "Go response should be []byte")

	railsBody, ok := railsResponse.([]byte)
	require.True(t, ok, "Rails response should be []byte")

	// Parse both responses
	goEnvelope := parseProjectsJSONAPIEnvelope(t, goBody)
	railsEnvelope := parseProjectsJSONAPIEnvelope(t, railsBody)

	// Extract data arrays
	goLogs := extractProjectsDataArray(t, goEnvelope)
	railsLogs := extractProjectsDataArray(t, railsEnvelope)

	// Compare counts (limited to 4)
	assert.Equal(t, len(railsLogs), len(goLogs),
		"Log count mismatch: Go=%d, Rails=%d", len(goLogs), len(railsLogs))

	if len(goLogs) == 0 {
		return
	}

	// Compare each log
	for i := range goLogs {
		goLog := goLogs[i]
		railsLog := railsLogs[i]

		// Compare basic fields
		assert.Equal(t, railsLog["start_page"], goLog["start_page"],
			"Log %d: start_page mismatch", i)
		assert.Equal(t, railsLog["end_page"], goLog["end_page"],
			"Log %d: end_page mismatch", i)

		// Compare note (may be null)
		assert.Equal(t, railsLog["note"], goLog["note"],
			"Log %d: note mismatch", i)
	}
}

// ============================================================================
// Comparison Test Cases
// ============================================================================

// ProjectsIndexComparisonTests defines all comparison tests for the projects index endpoint
var ProjectsIndexComparisonTests = []ProjectsRailsComparisonTest{
	{
		testName:     "Empty Database - No Projects",
		endpoint:     "/v1/projects.json",
		method:       "GET",
		setupProject: setupEmptyDatabase,
		validate:     validateProjectsIndex,
		skipIfEmpty:  false,
	},
	{
		testName:     "Single Project with Logs",
		endpoint:     "/v1/projects.json",
		method:       "GET",
		setupProject: setupSingleProject,
		validate:     validateProjectsIndex,
		skipIfEmpty:  false,
	},
	{
		testName:     "Multiple Projects with Varying Logs",
		endpoint:     "/v1/projects.json",
		method:       "GET",
		setupProject: setupMultipleProjects,
		validate:     validateProjectsIndex,
		skipIfEmpty:  false,
	},
}

// ProjectsShowComparisonTests defines all comparison tests for the projects show endpoint
var ProjectsShowComparisonTests = []ProjectsRailsComparisonTest{
	{
		testName:     "Project Show - Basic Response",
		endpoint:     "/v1/projects/1.json",
		method:       "GET",
		setupProject: setupSingleProject,
		validate:     validateProjectsShow,
		skipIfEmpty:  false,
	},
	{
		testName:     "Project Show - With Custom Dates",
		endpoint:     "/v1/projects/1.json",
		method:       "GET",
		setupProject: setupProjectWithCustomDates,
		validate:     validateProjectsShow,
		skipIfEmpty:  false,
	},
}

// LogsIndexComparisonTests defines all comparison tests for the logs index endpoint
var LogsIndexComparisonTests = []ProjectsRailsComparisonTest{
	{
		testName:     "Logs Index - Empty Project",
		endpoint:     "/v1/projects/1/logs.json",
		method:       "GET",
		setupProject: setupEmptyDatabase,
		validate:     validateLogsIndex,
		skipIfEmpty:  false,
	},
	{
		testName:     "Logs Index - Multiple Logs",
		endpoint:     "/v1/projects/1/logs.json",
		method:       "GET",
		setupProject: setupSingleProject,
		validate:     validateLogsIndex,
		skipIfEmpty:  false,
	},
	{
		testName:     "Logs Index - More Than 4 Logs (Limit Test)",
		endpoint:     "/v1/projects/2/logs.json",
		method:       "GET",
		setupProject: setupMultipleProjects,
		validate:     validateLogsIndex,
		skipIfEmpty:  false,
	},
}

// ============================================================================
// Test Functions
// ============================================================================

// TestProjectsIndexRailsComparison tests the GET /v1/projects.json endpoint
func TestProjectsIndexRailsComparison(t *testing.T) {
	for _, test := range ProjectsIndexComparisonTests {
		test.Run(t)
	}
}

// TestProjectsShowRailsComparison tests the GET /v1/projects/:id.json endpoint
func TestProjectsShowRailsComparison(t *testing.T) {
	for _, test := range ProjectsShowComparisonTests {
		test.Run(t)
	}
}

// TestLogsIndexRailsComparison tests the GET /v1/projects/:id/logs.json endpoint
func TestLogsIndexRailsComparison(t *testing.T) {
	for _, test := range LogsIndexComparisonTests {
		test.Run(t)
	}
}

// ============================================================================
// Additional Validation Tests
// ============================================================================

// TestRailsComparisonJSONAPICompliance verifies JSON:API format consistency
func TestRailsComparisonJSONAPICompliance(t *testing.T) {
	railsURL := os.Getenv("RAILS_API_URL")
	if railsURL == "" {
		t.Skip("RAILS_API_URL not set - skipping Rails comparison test")
	}

	if !IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	t.Run("Projects Index - JSON:API Envelope", func(t *testing.T) {
		helper := SetupTestDB(t)
		defer helper.Close()

		if err := helper.SetupTestSchema(); err != nil {
			t.Fatalf("Failed to setup test schema: %v", err)
		}

		// Setup test data
		ctx := &IntegrationTestContext{
			TestHelper: helper,
		}
		ctx.CreateTestProject(t)

		// Create Go handler with proper repositories
		projectRepo := postgres.NewProjectRepositoryImpl(helper.Pool)
		logRepo := postgres.NewLogRepositoryImpl(helper.Pool)
		dashboardRepo := postgres.NewDashboardRepositoryImpl(helper.Pool)
		userConfig, err := service.LoadDashboardConfig("dashboard_config.yaml")
		if err != nil {
			userConfig = service.NewUserConfigService(service.GetDefaultConfig())
		}
		goHandler := SetupRoutes(projectRepo, logRepo, dashboardRepo, userConfig, dashboard.ProjectsServiceInterface(&MockProjectsService{}))

		// Fetch both responses
		goBody := makeRequest(t, goHandler, "GET", "/v1/projects.json")
		railsBody := fetchRailsAPI(t, railsURL+"/v1/projects.json")

		// Verify both have JSON:API envelope structure
		goEnvelope := parseProjectsJSONAPIEnvelope(t, goBody)
		railsEnvelope := parseProjectsJSONAPIEnvelope(t, railsBody)

		// Check required JSON:API fields
		for _, envelope := range []map[string]interface{}{goEnvelope, railsEnvelope} {
			assert.Contains(t, envelope, "data", "Missing 'data' field")
		}
	})
}

// TestRailsComparisonErrorResponses verifies error response consistency
func TestRailsComparisonErrorResponses(t *testing.T) {
	railsURL := os.Getenv("RAILS_API_URL")
	if railsURL == "" {
		t.Skip("RAILS_API_URL not set - skipping Rails comparison test")
	}

	if !IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	t.Run("Project Not Found - 404", func(t *testing.T) {
		helper := SetupTestDB(t)
		defer helper.Close()

		if err := helper.SetupTestSchema(); err != nil {
			t.Fatalf("Failed to setup test schema: %v", err)
		}

		projectRepo := postgres.NewProjectRepositoryImpl(helper.Pool)
		logRepo := postgres.NewLogRepositoryImpl(helper.Pool)
		dashboardRepo := postgres.NewDashboardRepositoryImpl(helper.Pool)
		userConfig, err := service.LoadDashboardConfig("dashboard_config.yaml")
		if err != nil {
			userConfig = service.NewUserConfigService(service.GetDefaultConfig())
		}
		goHandler := SetupRoutes(projectRepo, logRepo, dashboardRepo, userConfig, dashboard.ProjectsServiceInterface(&MockProjectsService{}))

		// Fetch both 404 responses
		goBody := makeRequest(t, goHandler, "GET", "/v1/projects/999999.json")
		railsBody := fetchRailsAPI(t, railsURL+"/v1/projects/999999.json")

		// Both should return error responses
		var goError map[string]interface{}
		var railsError map[string]interface{}

		json.Unmarshal(goBody, &goError)
		json.Unmarshal(railsBody, &railsError)

		// Verify error structure
		assert.Contains(t, goError, "error", "Go response should have 'error' field")
		assert.Contains(t, railsError, "error", "Rails response should have 'error' field")
	})
}

// ============================================================================
// JSON Parsing Helpers (Projects/Logs specific)
// ============================================================================

// parseProjectsJSONAPIEnvelope parses a JSON:API envelope from raw bytes
// Uses unique name to avoid conflict with existing functions in rails_comparison_test.go
func parseProjectsJSONAPIEnvelope(t *testing.T, body []byte) map[string]interface{} {
	var envelope map[string]interface{}
	err := json.Unmarshal(body, &envelope)
	require.NoError(t, err, "Failed to parse JSON:API envelope")
	return envelope
}

// extractProjectsDataArray extracts the data array from a JSON:API envelope
func extractProjectsDataArray(t *testing.T, envelope map[string]interface{}) []map[string]interface{} {
	data, ok := envelope["data"]
	require.True(t, ok, "Response missing 'data' field")

	dataArray, ok := data.([]interface{})
	require.True(t, ok, "Response 'data' field should be an array")

	result := make([]map[string]interface{}, len(dataArray))
	for i, item := range dataArray {
		itemMap, ok := item.(map[string]interface{})
		require.True(t, ok, "Data item should be an object")

		// Extract attributes if present (JSON:API format)
		if attrs, ok := itemMap["attributes"].(map[string]interface{}); ok {
			result[i] = attrs
		} else {
			// Use the item directly (flat format)
			result[i] = itemMap
		}
	}

	return result
}

// extractProjectsSingleObject extracts a single object from a JSON:API envelope
func extractProjectsSingleObject(t *testing.T, envelope map[string]interface{}) map[string]interface{} {
	data, ok := envelope["data"]
	require.True(t, ok, "Response missing 'data' field")

	// Handle both object and array formats
	switch v := data.(type) {
	case map[string]interface{}:
		if attrs, ok := v["attributes"].(map[string]interface{}); ok {
			return attrs
		}
		return v
	case []interface{}:
		if len(v) > 0 {
			if itemMap, ok := v[0].(map[string]interface{}); ok {
				if attrs, ok := itemMap["attributes"].(map[string]interface{}); ok {
					return attrs
				}
				return itemMap
			}
		}
	}

	return make(map[string]interface{})
}
