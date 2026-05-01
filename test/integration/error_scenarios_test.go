package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pg "go-reading-log-api-next/internal/adapter/postgres"
	"go-reading-log-api-next/internal/api/v1/handlers"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/service"
	dashboardFixtures "go-reading-log-api-next/test/fixtures/dashboard"
)

// extractPath extracts the path portion from an endpoint URL, removing query parameters
// This ensures endpoints like "/v1/dashboard/day.json?date=invalid" match the case "/v1/dashboard/day.json"
func extractPath(endpoint string) string {
	if idx := strings.Index(endpoint, "?"); idx != -1 {
		return endpoint[:idx]
	}
	return endpoint
}

// ErrorScenario represents a test case for error handling
type ErrorScenario struct {
	Name     string
	Endpoint string
	Method   string
	Setup    func(*testing.T) *dashboardFixtures.Scenario
	Request  func(*http.Request)
	Validate func(*testing.T, *httptest.ResponseRecorder)
}

// RunErrorScenarios runs all error scenario tests
func RunErrorScenarios(t *testing.T, scenarios []ErrorScenario) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured")
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			helper := SetupTestDB(t)
			defer helper.Close()

			// Create database tables before inserting data
			err := helper.SetupTestSchema()
			require.NoError(t, err)

			// Load fixture if setup provided
			if scenario.Setup != nil {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				scenario := scenario.Setup(t)
				err := fixtureManager.LoadScenario(scenario)
				require.NoError(t, err)
			}

			// Create handler
			repo := pg.NewDashboardRepositoryImpl(helper.Pool)
			userConfig := service.NewUserConfigService(service.GetDefaultConfig())
			dashboardHandler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{})

			// Create request
			req := httptest.NewRequest(scenario.Method, scenario.Endpoint, nil)

			// Apply request modifications
			if scenario.Request != nil {
				scenario.Request(req)
			}

			// Execute request based on endpoint
			recorder := httptest.NewRecorder()

			switch {
			case extractPath(scenario.Endpoint) == "/v1/dashboard/day.json":
				dashboardHandler.Day(recorder, req)
			case extractPath(scenario.Endpoint) == "/v1/dashboard/projects.json":
				dashboardHandler.Projects(recorder, req)
			case extractPath(scenario.Endpoint) == "/v1/dashboard/last_days.json":
				dashboardHandler.LastDays(recorder, req)
			case extractPath(scenario.Endpoint) == "/v1/dashboard/echart/faults.json":
				dashboardHandler.Faults(recorder, req)
			case extractPath(scenario.Endpoint) == "/v1/dashboard/echart/speculate_actual.json":
				dashboardHandler.SpeculateActual(recorder, req)
			case extractPath(scenario.Endpoint) == "/v1/dashboard/echart/faults_week_day.json":
				dashboardHandler.WeekdayFaults(recorder, req)
			case extractPath(scenario.Endpoint) == "/v1/dashboard/echart/mean_progress.json":
				dashboardHandler.MeanProgress(recorder, req)
			case extractPath(scenario.Endpoint) == "/v1/dashboard/echart/last_year_total.json":
				dashboardHandler.YearlyTotal(recorder, req)
			default:
				t.Fatalf("Unknown endpoint: %s", scenario.Endpoint)
			}

			// Validate response
			scenario.Validate(t, recorder)
		})
	}
}

// Error Scenarios for Dashboard Endpoints

var DashboardErrorScenarios = []ErrorScenario{
	{
		Name:     "Day Endpoint - Invalid Date",
		Endpoint: "/v1/dashboard/day.json?date=invalid",
		Method:   "GET",
		Validate: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			// Verify error response format
			var response map[string]interface{}
			json.Unmarshal(recorder.Body.Bytes(), &response)
			assert.Contains(t, response, "error")
			assert.Contains(t, response, "details")
		},
	},
	{
		Name:     "Last Days - Invalid Type",
		Endpoint: "/v1/dashboard/last_days.json?type=99",
		Method:   "GET",
		Validate: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
		},
	},
	{
		Name:     "Projects Endpoint - Empty Database",
		Endpoint: "/v1/dashboard/projects.json",
		Method:   "GET",
		Setup: func(t *testing.T) *dashboardFixtures.Scenario {
			return dashboardFixtures.ScenarioEmptyData()
		},
		Validate: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			// Parse response using helper function
			response, err := parseDashboardResponse(recorder.Body.Bytes())
			require.NoError(t, err)
			// Should return empty array with zero values
			assert.Empty(t, response.Logs)
		},
	},
	{
		Name:     "Day Endpoint - Empty Database",
		Endpoint: "/v1/dashboard/day.json",
		Method:   "GET",
		Setup: func(t *testing.T) *dashboardFixtures.Scenario {
			return dashboardFixtures.ScenarioEmptyData()
		},
		Validate: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			// Parse response using helper function
			response, err := parseDashboardResponse(recorder.Body.Bytes())
			require.NoError(t, err)
			// Should return zero values for empty data
			assert.NotNil(t, response.Stats)
			assert.Equal(t, 0, response.Stats.PreviousWeekPages)
			assert.Equal(t, 0, response.Stats.LastWeekPages)
		},
	},
	{
		Name:     "Mean Progress - Empty Database",
		Endpoint: "/v1/dashboard/echart/mean_progress.json",
		Method:   "GET",
		Setup: func(t *testing.T) *dashboardFixtures.Scenario {
			return dashboardFixtures.ScenarioEmptyData()
		},
		Validate: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, recorder.Code)
			// Parse response using helper function
			response, err := parseDashboardResponse(recorder.Body.Bytes())
			require.NoError(t, err)
			// Should return valid chart config even with empty data
			assert.NotNil(t, response.Echart)
		},
	},
}

// parseDashboardResponse parses a JSON:API envelope and extracts DashboardResponse
// Handles both direct response and JSON:API envelope formats for backward compatibility
// parseDashboardResponse parses a JSON:API envelope and extracts DashboardResponse
// Handles both direct response and JSON:API envelope formats for backward compatibility
func parseDashboardResponse(body []byte) (*dto.DashboardResponse, error) {
	var envelope map[string]interface{}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	// Check if this is a JSON:API envelope or direct response
	data, hasData := envelope["data"].(map[string]interface{})
	if !hasData {
		// Not a JSON:API envelope, try to parse as direct response
		var response dto.DashboardResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}
		return &response, nil
	}

	attributes, hasAttributes := data["attributes"].(map[string]interface{})
	if !hasAttributes {
		return nil, fmt.Errorf("expected 'attributes' field in data")
	}

	response := &dto.DashboardResponse{}

	// Parse DailyStats (what handlers actually return) into StatsData
	if stats, ok := attributes["stats"].(map[string]interface{}); ok {
		response.Stats = &dto.StatsData{}
		if v, exists := stats["per_pages"]; exists {
			if v != nil {
				val := v.(float64)
				response.Stats.PerPages = &val
			}
		}
		if v, exists := stats["mean_day"]; exists {
			response.Stats.MeanDay = v.(float64)
		}
		if v, exists := stats["spec_mean_day"]; exists {
			response.Stats.SpecMeanDay = v.(float64)
		}
		if v, exists := stats["progress_geral"]; exists {
			response.Stats.ProgressGeral = v.(float64)
		}
	}

	// Also handle DailyStats format (total_pages, log_count)
	if _, ok := attributes["total_pages"]; ok {
		if response.Stats == nil {
			response.Stats = &dto.StatsData{}
		}
		if v, exists := attributes["total_pages"]; exists {
			response.Stats.TotalPages = int(v.(float64))
		}
	}

	// Parse Logs
	if logs, ok := attributes["logs"].([]interface{}); ok {
		response.Logs = make([]dto.LogEntry, len(logs))
		for i, log := range logs {
			logMap := log.(map[string]interface{})
			response.Logs[i] = dto.LogEntry{
				ID:        int64(logMap["id"].(float64)),
				Data:      logMap["data"].(string),
				StartPage: int(logMap["start_page"].(float64)),
				EndPage:   int(logMap["end_page"].(float64)),
			}
		}
	}

	// Parse Echart - either under "echart" key or directly in attributes
	if echart, ok := attributes["echart"].(map[string]interface{}); ok {
		response.Echart = &dto.EchartConfig{}
		if title, exists := echart["title"]; exists {
			response.Echart.Title = title.(string)
		}
		if tooltip, exists := echart["tooltip"]; exists {
			response.Echart.Tooltip = tooltip.(map[string]interface{})
		}
		if series, exists := echart["series"]; exists {
			if seriesArr, ok := series.([]interface{}); ok {
				response.Echart.Series = make([]dto.Series, len(seriesArr))
				for i, s := range seriesArr {
					if sMap, ok := s.(map[string]interface{}); ok {
						seriesData := make([]interface{}, 0)
						if data, exists := sMap["data"]; exists {
							if dataArr, ok := data.([]interface{}); ok {
								seriesData = dataArr
							}
						}
						response.Echart.Series[i] = dto.Series{
							Name: sMap["name"].(string),
							Type: sMap["type"].(string),
							Data: seriesData,
						}
					}
				}
			}
		}
	} else if _, ok := attributes["title"]; ok {
		// Echart config is directly in attributes (e.g., MeanProgress endpoint)
		response.Echart = &dto.EchartConfig{}
		if title, exists := attributes["title"]; exists {
			response.Echart.Title = title.(string)
		}
		if tooltip, exists := attributes["tooltip"]; exists {
			response.Echart.Tooltip = tooltip.(map[string]interface{})
		}
		if series, exists := attributes["series"]; exists {
			if seriesArr, ok := series.([]interface{}); ok {
				response.Echart.Series = make([]dto.Series, len(seriesArr))
				for i, s := range seriesArr {
					if sMap, ok := s.(map[string]interface{}); ok {
						seriesData := make([]interface{}, 0)
						if data, exists := sMap["data"]; exists {
							if dataArr, ok := data.([]interface{}); ok {
								seriesData = dataArr
							}
						}
						response.Echart.Series[i] = dto.Series{
							Name: sMap["name"].(string),
							Type: sMap["type"].(string),
							Data: seriesData,
						}
					}
				}
			}
		}
	}

	return response, nil
}

// TestErrorScenarios runs all error scenario tests
func TestErrorScenarios(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured")
	}
	RunErrorScenarios(t, DashboardErrorScenarios)
}
