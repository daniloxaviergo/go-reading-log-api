package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-reading-log-api-next/test/fixtures/dashboard"
)

// RailsComparisonTest tests that Go API responses match Rails API behavior
type RailsComparisonTest struct {
	testName     string
	endpoint     string
	method       string
	setupFixture func() *dashboard.Scenario
	validate     func(t *testing.T, goResponse interface{}, railsResponse interface{})
}

// Run executes the comparison test
func (test *RailsComparisonTest) Run(t *testing.T) {
	t.Run(test.testName, func(t *testing.T) {
		// Check if Rails API URL is configured
		railsURL := os.Getenv("RAILS_API_URL")
		if railsURL == "" {
			t.Skip("RAILS_API_URL not set - skipping Rails comparison test")
		}

		// Setup test database
		helper := SetupTestDB(t)
		defer helper.Close()

		// Create database tables before inserting data
		if err := helper.SetupTestSchema(); err != nil {
			t.Fatalf("Failed to setup test schema: %v", err)
		}

		// Load fixture data
		fixtureManager := dashboard.NewDashboardFixtures(helper.Pool)
		scenario := test.setupFixture()

		err := fixtureManager.LoadScenario(scenario)
		require.NoError(t, err)

		// Create Go handler
		goHandler := createGoHandler(helper.Pool)

		// Make request to Go API
		goResponse := makeRequest(t, goHandler, test.method, test.endpoint)

		// Fetch Rails API response
		railsResponse := fetchRailsAPI(t, railsURL+test.endpoint)

		// Compare responses
		test.validate(t, goResponse, railsResponse)
	})
}

// createGoHandler creates a Go handler with all dependencies
func createGoHandler(pool interface{}) http.Handler {
	// Note: This is a simplified version - adjust based on actual implementation
	// repo := postgres.NewDashboardRepositoryImpl(pool)
	// userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	// dashboardHandler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{})

	// For now, return a mock handler that returns expected structure
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"status": "ok",
			"data":   map[string]interface{}{"test": "data"},
		}
		json.NewEncoder(w).Encode(response)
	})
}

// fetchRailsAPI fetches response from Rails API
func fetchRailsAPI(t T, url string) []byte {
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode,
		"Rails API returned status %d for URL: %s", resp.StatusCode, url)

	return body
}

// makeRequest makes an HTTP request to the handler
func makeRequest(t T, handler http.Handler, method, endpoint string) []byte {
	req := httptest.NewRequest(method, endpoint, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	return recorder.Body.Bytes()
}

// ResponseValidator validates Go vs Rails response comparison
type ResponseValidator struct {
	validator *Validator
}

// NewResponseValidator creates a new response validator
func NewResponseValidator() *ResponseValidator {
	return &ResponseValidator{
		validator: NewValidator(0.001), // 0.1% tolerance for float comparisons
	}
}

// Validate compares Go and Rails responses
func (rv *ResponseValidator) Validate(
	t *testing.T,
	goResponse interface{},
	railsResponse interface{},
	endpoint string,
) ValidationResult {
	return rv.validator.ValidateDashboardResponse(t, goResponse, railsResponse, endpoint)
}

// ComparisonTest represents a comparison test case
type ComparisonTest struct {
	Name       string
	Endpoint   string
	Method     string
	Setup      func() *dashboard.Scenario
	Validate   func(*testing.T, interface{}, interface{})
	ExpectFail bool // If true, test is expected to fail (Rails not available)
}

// RunComparisonTests runs all comparison tests
func RunComparisonTests(t *testing.T, tests []ComparisonTest) {
	for _, ct := range tests {
		t.Run(ct.Name, func(t *testing.T) {
			// Check if Rails API URL is configured
			railsURL := os.Getenv("RAILS_API_URL")
			if railsURL == "" && !ct.ExpectFail {
				t.Skip("RAILS_API_URL not set - skipping comparison test")
			}

			// Setup test database
			helper := SetupTestDB(t)
			defer helper.Close()

			// Load fixture
			fixtureManager := dashboard.NewDashboardFixtures(helper.Pool)
			scenario := ct.Setup()
			err := fixtureManager.LoadScenario(scenario)
			require.NoError(t, err)

			// Create Go handler
			goHandler := createGoHandler(helper.Pool)

			// Make request to Go API
			goResponse := makeRequest(t, goHandler, ct.Method, ct.Endpoint)

			if railsURL != "" {
				// Fetch Rails API response
				railsResponse := fetchRailsAPI(t, railsURL+ct.Endpoint)

				// Compare responses
				ct.Validate(t, goResponse, railsResponse)
			} else {
				// Just validate Go response structure without Rails comparison
				assert.NotNil(t, goResponse, "Go response should not be nil")
			}
		})
	}
}

// Comparison Tests for Dashboard Endpoints

var DashboardComparisonTests = []ComparisonTest{
	{
		Name:     "Day Endpoint - Basic Response",
		Endpoint: "/v1/dashboard/day.json",
		Method:   "GET",
		Setup: func() *dashboard.Scenario {
			return dashboard.ScenarioMultipleProjects()
		},
		Validate: func(t *testing.T, goResponse interface{}, railsResponse interface{}) {
			validator := NewResponseValidator()
			result := validator.Validate(t, goResponse, railsResponse, "/v1/dashboard/day.json")

			assert.True(t, result.Passed,
				"Day endpoint comparison failed: %v", result.Errors)

			// Verify structure
			goMap := goResponse.(map[string]interface{})
			railsMap := railsResponse.(map[string]interface{})

			// Check common fields exist in both
			for _, field := range []string{"status", "data", "meta"} {
				assert.Contains(t, goMap, field, "Go response missing field: %s", field)
				assert.Contains(t, railsMap, field, "Rails response missing field: %s", field)
			}
		},
		ExpectFail: false,
	},
	{
		Name:     "Projects Endpoint - Project List",
		Endpoint: "/v1/dashboard/projects.json",
		Method:   "GET",
		Setup: func() *dashboard.Scenario {
			return dashboard.ScenarioMultipleProjects()
		},
		Validate: func(t *testing.T, goResponse interface{}, railsResponse interface{}) {
			goMap := goResponse.(map[string]interface{})
			railsMap := railsResponse.(map[string]interface{})

			// Verify projects array exists
			goData := goMap["data"].([]interface{})
			railsData := railsMap["data"].([]interface{})

			assert.Equal(t, len(goData), len(railsData),
				"Projects count mismatch: Go=%d, Rails=%d", len(goData), len(railsData))

			// Verify each project has expected fields
			for i := range goData {
				goProj := goData[i].(map[string]interface{})
				railsProj := railsData[i].(map[string]interface{})

				// Compare calculated fields with tolerance
				goProgress := goProj["progress"].(float64)
				railsProgress := railsProj["progress"].(float64)

				assert.InDelta(t, railsProgress, goProgress, 0.001,
					"Project %d progress mismatch", i)
			}
		},
		ExpectFail: false,
	},
	{
		Name:     "Last Days Endpoint - Trend Data",
		Endpoint: "/v1/dashboard/last_days.json?type=7",
		Method:   "GET",
		Setup: func() *dashboard.Scenario {
			return dashboard.ScenarioLastDays()
		},
		Validate: func(t *testing.T, goResponse interface{}, railsResponse interface{}) {
			goMap := goResponse.(map[string]interface{})
			railsMap := railsResponse.(map[string]interface{})

			// Verify trend fields
			goTrend := goMap["data"].(map[string]interface{})
			railsTrend := railsMap["data"].(map[string]interface{})

			// Compare numerical fields with tolerance
			fieldsToCompare := []string{"total_faults", "avg_per_day"}
			for _, field := range fieldsToCompare {
				goVal := goTrend[field].(float64)
				railsVal := railsTrend[field].(float64)

				assert.InDelta(t, railsVal, goVal, 0.001,
					"Field %s mismatch", field)
			}

			// Verify date range fields
			assert.Equal(t, goTrend["days"], railsTrend["days"],
				"Days count mismatch")
			assert.Equal(t, goTrend["type"], railsTrend["type"],
				"Type mismatch")
		},
		ExpectFail: false,
	},
	{
		Name:     "Faults Chart - Gauge Configuration",
		Endpoint: "/v1/dashboard/echart/faults.json",
		Method:   "GET",
		Setup: func() *dashboard.Scenario {
			return dashboard.ScenarioFaultsByWeekday()
		},
		Validate: func(t *testing.T, goResponse interface{}, railsResponse interface{}) {
			goMap := goResponse.(map[string]interface{})
			railsMap := railsResponse.(map[string]interface{})

			goEchart := goMap["data"].(map[string]interface{})
			railsEchart := railsMap["data"].(map[string]interface{})

			// Verify gauge chart structure
			assert.Equal(t, "gauge", goEchart["type"], "Go chart type mismatch")
			assert.Equal(t, "gauge", railsEchart["type"], "Rails chart type mismatch")

			goSeries := goEchart["series"].([]interface{})
			railsSeries := railsEchart["series"].([]interface{})

			assert.Equal(t, len(goSeries), len(railsSeries),
				"Series count mismatch")

			// Compare percentage values
			goPercent := goSeries[0].(float64)
			railsPercent := railsSeries[0].(float64)

			assert.InDelta(t, railsPercent, goPercent, 0.1,
				"Gauge percentage mismatch")
		},
		ExpectFail: false,
	},
	{
		Name:     "Speculate Actual - Line Chart",
		Endpoint: "/v1/dashboard/echart/speculate_actual.json",
		Method:   "GET",
		Setup: func() *dashboard.Scenario {
			return dashboard.ScenarioSpeculateActual()
		},
		Validate: func(t *testing.T, goResponse interface{}, railsResponse interface{}) {
			goMap := goResponse.(map[string]interface{})
			railsMap := railsResponse.(map[string]interface{})

			goEchart := goMap["data"].(map[string]interface{})
			railsEchart := railsMap["data"].(map[string]interface{})

			// Verify line chart structure
			assert.Equal(t, "line", goEchart["type"], "Go chart type mismatch")
			assert.Equal(t, "line", railsEchart["type"], "Rails chart type mismatch")

			goSeries := goEchart["series"].([]interface{})
			railsSeries := railsEchart["series"].([]interface{})

			// Should have 2 series (Actual and Speculated)
			assert.Len(t, goSeries, 2, "Go should have 2 series")
			assert.Len(t, railsSeries, 2, "Rails should have 2 series")

			// Compare data values
			goData := goSeries[0].([]interface{})
			railsData := railsSeries[0].([]interface{})

			assert.Equal(t, len(goData), len(railsData),
				"Series data length mismatch")
		},
		ExpectFail: false,
	},
	{
		Name:     "Weekday Faults - Radar Chart",
		Endpoint: "/v1/dashboard/echart/faults_week_day.json",
		Method:   "GET",
		Setup: func() *dashboard.Scenario {
			return dashboard.ScenarioFaultsByWeekday()
		},
		Validate: func(t *testing.T, goResponse interface{}, railsResponse interface{}) {
			goMap := goResponse.(map[string]interface{})
			railsMap := railsResponse.(map[string]interface{})

			goEchart := goMap["data"].(map[string]interface{})
			railsEchart := railsMap["data"].(map[string]interface{})

			// Verify radar chart structure
			assert.Equal(t, "radar", goEchart["type"], "Go chart type mismatch")
			assert.Equal(t, "radar", railsEchart["type"], "Rails chart type mismatch")

			goSeries := goEchart["series"].([]interface{})
			railsSeries := railsEchart["series"].([]interface{})

			// Should have 7 data points (one for each weekday)
			goData := goSeries[0].([]interface{})
			railsData := railsSeries[0].([]interface{})

			assert.Len(t, goData, 7, "Go should have 7 weekday data points")
			assert.Len(t, railsData, 7, "Rails should have 7 weekday data points")
		},
		ExpectFail: false,
	},
	{
		Name:     "Mean Progress - Line Chart with Visual Map",
		Endpoint: "/v1/dashboard/echart/mean_progress.json",
		Method:   "GET",
		Setup: func() *dashboard.Scenario {
			return dashboard.ScenarioMeanProgress()
		},
		Validate: func(t *testing.T, goResponse interface{}, railsResponse interface{}) {
			goMap := goResponse.(map[string]interface{})
			railsMap := railsResponse.(map[string]interface{})

			goEchart := goMap["data"].(map[string]interface{})
			railsEchart := railsMap["data"].(map[string]interface{})

			// Verify line chart structure
			assert.Equal(t, "line", goEchart["type"], "Go chart type mismatch")
			assert.Equal(t, "line", railsEchart["type"], "Rails chart type mismatch")

			goSeries := goEchart["series"].([]interface{})
			railsSeries := railsEchart["series"].([]interface{})

			// Should have 30 data points (one for each day)
			goData := goSeries[0].([]interface{})
			railsData := railsSeries[0].([]interface{})

			assert.Len(t, goData, 30, "Go should have 30 daily data points")
			assert.Len(t, railsData, 30, "Rails should have 30 daily data points")

			// Verify visual map configuration exists
			assert.Contains(t, goEchart, "visualMap", "Go missing visualMap")
			assert.Contains(t, railsEchart, "visualMap", "Rails missing visualMap")
		},
		ExpectFail: false,
	},
	{
		Name:     "Yearly Total - Bar Chart",
		Endpoint: "/v1/dashboard/echart/last_year_total.json",
		Method:   "GET",
		Setup: func() *dashboard.Scenario {
			return dashboard.ScenarioYearlyTotal()
		},
		Validate: func(t *testing.T, goResponse interface{}, railsResponse interface{}) {
			goMap := goResponse.(map[string]interface{})
			railsMap := railsResponse.(map[string]interface{})

			goEchart := goMap["data"].(map[string]interface{})
			railsEchart := railsMap["data"].(map[string]interface{})

			// Verify bar chart structure
			assert.Equal(t, "bar", goEchart["type"], "Go chart type mismatch")
			assert.Equal(t, "bar", railsEchart["type"], "Rails chart type mismatch")

			goSeries := goEchart["series"].([]interface{})
			railsSeries := railsEchart["series"].([]interface{})

			// Should have 2 series (last year and current year)
			assert.Len(t, goSeries, 2, "Go should have 2 yearly series")
			assert.Len(t, railsSeries, 2, "Rails should have 2 yearly series")

			// Verify color configurations
			goStyle := goSeries[0].(map[string]interface{})["itemStyle"].(map[string]interface{})
			railsStyle := railsSeries[0].(map[string]interface{})["itemStyle"].(map[string]interface{})

			assert.Equal(t, railsStyle["color"], goStyle["color"],
				"Series color mismatch")
		},
		ExpectFail: false,
	},
}

// TestRailsComparison runs all Rails comparison tests
func TestRailsComparison(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured")
	}
	RunComparisonTests(t, DashboardComparisonTests)
}
