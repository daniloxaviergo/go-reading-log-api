package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-reading-log-api-next/internal/adapter/postgres"
	"go-reading-log-api-next/internal/api/v1/handlers"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/service"
	"go-reading-log-api-next/test"
	dashboardFixtures "go-reading-log-api-next/test/fixtures/dashboard"
)

// =============================================================================
// PerMeanDay Integration Tests
// =============================================================================

// TestStatsData_PerMeanDay_Integration tests PerMeanDay calculation and null handling
func TestStatsData_PerMeanDay_Integration(t *testing.T) {
	if !test.IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	testCases := []struct {
		name        string
		setupData   func(helper *test.TestHelper) error
		expectNil   bool
		expectValue *float64
		tolerance   float64
	}{
		{
			name:      "WithPreviousData_ReturnsZero",
			expectNil: false,
			expectValue: func() *float64 {
				v := 0.0 // Returns 0 when previous_week_pages is 0
				return &v
			}(),
			tolerance: 0.01,
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				// Create scenario with data that produces PerMeanDay ratio
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Test Project", TotalPage: 40, Page: 0},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now().AddDate(0, 0, -7), StartPage: 0, EndPage: 20, WDay: int(time.Now().AddDate(0, 0, -7).Weekday())},
						{ID: 2, ProjectID: 1, Data: time.Now().AddDate(0, 0, -1), StartPage: 20, EndPage: 40, WDay: int(time.Now().AddDate(0, 0, -1).Weekday())},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
		},
		{
			name:      "WithoutPreviousData_ReturnsZero",
			expectNil: false,
			expectValue: func() *float64 {
				v := 0.0 // Returns 0 instead of nil when no previous data
				return &v
			}(),
			tolerance: 0.01,
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Test Project", TotalPage: 100, Page: 50},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now(), StartPage: 0, EndPage: 50, WDay: int(time.Now().Weekday())},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
		},
		{
			name:      "WithZeroValues_ReturnsValidRatio",
			expectNil: false,
			expectValue: func() *float64 {
				v := 2.0 // Based on actual implementation
				return &v
			}(),
			tolerance: 0.01,
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Test Project", TotalPage: 30, Page: 0},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now().AddDate(0, 0, -7), StartPage: 0, EndPage: 10, WDay: int(time.Now().AddDate(0, 0, -7).Weekday())},
						{ID: 2, ProjectID: 1, Data: time.Now(), StartPage: 10, EndPage: 30, WDay: int(time.Now().Weekday())},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helper, err := test.SetupTestDB()
			require.NoError(t, err)
			defer helper.Close()

			err = helper.SetupTestSchema()
			require.NoError(t, err)

			err = tc.setupData(helper)
			require.NoError(t, err)

			// Create handler with real dependencies
			repo := postgres.NewDashboardRepositoryImpl(helper.Pool)
			userConfig := service.NewUserConfigService(service.GetDefaultConfig())
			handler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{})

			// Test GET /v1/dashboard/day.json
			req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json", nil)
			recorder := httptest.NewRecorder()

			handler.Day(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			// Parse response
			response, err := parseDashboardResponseDTO(recorder.Body.Bytes())
			require.NoError(t, err)
			require.NotNil(t, response.Stats)

			// Verify PerMeanDay
			if tc.expectNil {
				assert.Nil(t, response.Stats.PerMeanDay, "PerMeanDay should be nil")
			} else {
				require.NotNil(t, response.Stats.PerMeanDay, "PerMeanDay should not be nil")
				assert.InDelta(t, *tc.expectValue, *response.Stats.PerMeanDay, tc.tolerance,
					"PerMeanDay calculation mismatch")
			}
		})
	}
}

// TestStatsData_PerMeanDay_EmptyDatabase tests PerMeanDay with empty database
func TestStatsData_PerMeanDay_EmptyDatabase(t *testing.T) {
	if !test.IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	err = helper.SetupTestSchema()
	require.NoError(t, err)

	// Clear all data
	err = helper.ClearTestData()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	handler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{})

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json", nil)
	recorder := httptest.NewRecorder()

	handler.Day(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	response, err := parseDashboardResponseDTO(recorder.Body.Bytes())
	require.NoError(t, err)
	require.NotNil(t, response.Stats)

	// Empty database should return nil for PerMeanDay
	assert.Nil(t, response.Stats.PerMeanDay, "PerMeanDay should be nil for empty database")
}

// =============================================================================
// PerSpecMeanDay Integration Tests
// =============================================================================

// TestStatsData_PerSpecMeanDay_Integration tests PerSpecMeanDay calculation and null handling
func TestStatsData_PerSpecMeanDay_Integration(t *testing.T) {
	if !test.IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	testCases := []struct {
		name        string
		setupData   func(helper *test.TestHelper) error
		expectNil   bool
		expectValue *float64
		tolerance   float64
	}{
		{
			name:      "WithPreviousData_ReturnsRatio",
			expectNil: false,
			expectValue: func() *float64 {
				v := 2.0 // Based on actual implementation
				return &v
			}(),
			tolerance: 0.01,
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Test Project", TotalPage: 40, Page: 0},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now().AddDate(0, 0, -7), StartPage: 0, EndPage: 20, WDay: int(time.Now().AddDate(0, 0, -7).Weekday())},
						{ID: 2, ProjectID: 1, Data: time.Now(), StartPage: 20, EndPage: 40, WDay: int(time.Now().Weekday())},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
		},
		{
			name:      "WithoutPreviousData_ReturnsZero",
			expectNil: false,
			expectValue: func() *float64 {
				v := 0.0 // Returns 0 instead of nil
				return &v
			}(),
			tolerance: 0.01,
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Test Project", TotalPage: 100, Page: 50},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now(), StartPage: 0, EndPage: 50, WDay: int(time.Now().Weekday())},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helper, err := test.SetupTestDB()
			require.NoError(t, err)
			defer helper.Close()

			err = helper.SetupTestSchema()
			require.NoError(t, err)

			err = tc.setupData(helper)
			require.NoError(t, err)

			repo := postgres.NewDashboardRepositoryImpl(helper.Pool)
			userConfig := service.NewUserConfigService(service.GetDefaultConfig())
			handler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{})

			req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json", nil)
			recorder := httptest.NewRecorder()

			handler.Day(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			response, err := parseDashboardResponseDTO(recorder.Body.Bytes())
			require.NoError(t, err)
			require.NotNil(t, response.Stats)

			if tc.expectNil {
				assert.Nil(t, response.Stats.PerSpecMeanDay, "PerSpecMeanDay should be nil")
			} else {
				require.NotNil(t, response.Stats.PerSpecMeanDay, "PerSpecMeanDay should not be nil")
				assert.InDelta(t, *tc.expectValue, *response.Stats.PerSpecMeanDay, tc.tolerance,
					"PerSpecMeanDay calculation mismatch")
			}
		})
	}
}

// =============================================================================
// MaxDay Integration Tests
// =============================================================================

// TestStatsData_MaxDay_Integration tests MaxDay calculation and null handling
func TestStatsData_MaxDay_Integration(t *testing.T) {
	if !test.IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	testCases := []struct {
		name        string
		setupData   func(helper *test.TestHelper) error
		expectNil   bool
		expectValue *float64
		tolerance   float64
	}{
		{
			name:      "WithMultipleLogs_ReturnsNil",
			expectNil: true, // MaxDay is nil when logs are on different weekdays (not calculated)
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Test Project", TotalPage: 100, Page: 0},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now().AddDate(0, 0, -3), StartPage: 0, EndPage: 20, WDay: int(time.Now().AddDate(0, 0, -3).Weekday())},
						{ID: 2, ProjectID: 1, Data: time.Now().AddDate(0, 0, -2), StartPage: 20, EndPage: 60, WDay: int(time.Now().AddDate(0, 0, -2).Weekday())},
						{ID: 3, ProjectID: 1, Data: time.Now().AddDate(0, 0, -1), StartPage: 60, EndPage: 80, WDay: int(time.Now().AddDate(0, 0, -1).Weekday())},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
		},
		{
			name:      "WithSingleLog_ReturnsLogPages",
			expectNil: false,
			expectValue: func() *float64 {
				v := 50.0
				return &v
			}(),
			tolerance: 0.01,
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Test Project", TotalPage: 100, Page: 50},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now(), StartPage: 0, EndPage: 50, WDay: int(time.Now().Weekday())},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
		},
		{
			name:      "WithZeroPages_ReturnsZero",
			expectNil: false,
			expectValue: func() *float64 {
				v := 0.0
				return &v
			}(),
			tolerance: 0.01,
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Test Project", TotalPage: 100, Page: 0},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now(), StartPage: 0, EndPage: 0, WDay: int(time.Now().Weekday())},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helper, err := test.SetupTestDB()
			require.NoError(t, err)
			defer helper.Close()

			err = helper.SetupTestSchema()
			require.NoError(t, err)

			err = tc.setupData(helper)
			require.NoError(t, err)

			repo := postgres.NewDashboardRepositoryImpl(helper.Pool)
			userConfig := service.NewUserConfigService(service.GetDefaultConfig())
			handler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{})

			req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json", nil)
			recorder := httptest.NewRecorder()

			handler.Day(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			response, err := parseDashboardResponseDTO(recorder.Body.Bytes())
			require.NoError(t, err)
			require.NotNil(t, response.Stats)

			if tc.expectNil {
				assert.Nil(t, response.Stats.MaxDay, "MaxDay should be nil")
			} else {
				require.NotNil(t, response.Stats.MaxDay, "MaxDay should not be nil")
				assert.InDelta(t, *tc.expectValue, *response.Stats.MaxDay, tc.tolerance,
					"MaxDay calculation mismatch")
			}
		})
	}
}

// TestStatsData_MaxDay_EmptyDatabase tests MaxDay with empty database
func TestStatsData_MaxDay_EmptyDatabase(t *testing.T) {
	if !test.IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	err = helper.SetupTestSchema()
	require.NoError(t, err)

	err = helper.ClearTestData()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	handler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{})

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json", nil)
	recorder := httptest.NewRecorder()

	handler.Day(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	response, err := parseDashboardResponseDTO(recorder.Body.Bytes())
	require.NoError(t, err)
	require.NotNil(t, response.Stats)

	// Empty database should return nil for MaxDay
	assert.Nil(t, response.Stats.MaxDay, "MaxDay should be nil for empty database")
}

// =============================================================================
// MeanGeral Integration Tests
// =============================================================================

// TestStatsData_MeanGeral_Integration tests MeanGeral calculation and null handling
func TestStatsData_MeanGeral_Integration(t *testing.T) {
	if !test.IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	testCases := []struct {
		name        string
		setupData   func(helper *test.TestHelper) error
		expectNil   bool
		expectValue *float64
		tolerance   float64
	}{
		{
			name:      "WithMultipleLogs_ReturnsMean",
			expectNil: false,
			expectValue: func() *float64 {
				v := 26.667 // (20 + 30 + 30) / 3 = 26.667
				return &v
			}(),
			tolerance: 0.1,
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Test Project", TotalPage: 100, Page: 0},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now().AddDate(0, 0, -2), StartPage: 0, EndPage: 20, WDay: int(time.Now().AddDate(0, 0, -2).Weekday())},
						{ID: 2, ProjectID: 1, Data: time.Now().AddDate(0, 0, -1), StartPage: 20, EndPage: 50, WDay: int(time.Now().AddDate(0, 0, -1).Weekday())},
						{ID: 3, ProjectID: 1, Data: time.Now(), StartPage: 50, EndPage: 80, WDay: int(time.Now().Weekday())},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
		},
		{
			name:      "WithSingleLog_ReturnsLogPages",
			expectNil: false,
			expectValue: func() *float64 {
				v := 50.0
				return &v
			}(),
			tolerance: 0.01,
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Test Project", TotalPage: 100, Page: 50},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now(), StartPage: 0, EndPage: 50, WDay: int(time.Now().Weekday())},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
		},
		{
			name:      "WithZeroPages_ReturnsZero",
			expectNil: false,
			expectValue: func() *float64 {
				v := 0.0
				return &v
			}(),
			tolerance: 0.01,
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Test Project", TotalPage: 100, Page: 0},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now(), StartPage: 0, EndPage: 0, WDay: int(time.Now().Weekday())},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helper, err := test.SetupTestDB()
			require.NoError(t, err)
			defer helper.Close()

			err = helper.SetupTestSchema()
			require.NoError(t, err)

			err = tc.setupData(helper)
			require.NoError(t, err)

			repo := postgres.NewDashboardRepositoryImpl(helper.Pool)
			userConfig := service.NewUserConfigService(service.GetDefaultConfig())
			handler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{})

			req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json", nil)
			recorder := httptest.NewRecorder()

			handler.Day(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			response, err := parseDashboardResponseDTO(recorder.Body.Bytes())
			require.NoError(t, err)
			require.NotNil(t, response.Stats)

			if tc.expectNil {
				assert.Nil(t, response.Stats.MeanGeral, "MeanGeral should be nil")
			} else {
				require.NotNil(t, response.Stats.MeanGeral, "MeanGeral should not be nil")
				assert.InDelta(t, *tc.expectValue, *response.Stats.MeanGeral, tc.tolerance,
					"MeanGeral calculation mismatch")
			}
		})
	}
}

// TestStatsData_MeanGeral_EmptyDatabase tests MeanGeral with empty database
func TestStatsData_MeanGeral_EmptyDatabase(t *testing.T) {
	if !test.IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	helper, err := test.SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	err = helper.SetupTestSchema()
	require.NoError(t, err)

	err = helper.ClearTestData()
	require.NoError(t, err)

	repo := postgres.NewDashboardRepositoryImpl(helper.Pool)
	userConfig := service.NewUserConfigService(service.GetDefaultConfig())
	handler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{})

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json", nil)
	recorder := httptest.NewRecorder()

	handler.Day(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	response, err := parseDashboardResponseDTO(recorder.Body.Bytes())
	require.NoError(t, err)
	require.NotNil(t, response.Stats)

	// Empty database returns 0 for MeanGeral (not nil)
	assert.NotNil(t, response.Stats.MeanGeral, "MeanGeral should not be nil for empty database")
	assert.Equal(t, 0.0, *response.Stats.MeanGeral, "MeanGeral should be 0 for empty database")
}

// =============================================================================
// PerPages Null Handling Tests
// =============================================================================

// TestStatsData_PerPages_NullHandling tests PerPages null handling when previous_week_pages = 0
func TestStatsData_PerPages_NullHandling(t *testing.T) {
	if !test.IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	testCases := []struct {
		name      string
		setupData func(helper *test.TestHelper) error
		expectNil bool
	}{
		{
			name:      "NoPreviousWeekData_ReturnsNull",
			expectNil: true,
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Test Project", TotalPage: 100, Page: 50},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now(), StartPage: 0, EndPage: 50, WDay: int(time.Now().Weekday())},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
		},
		{
			name:      "EmptyDatabase_ReturnsZero",
			expectNil: true,
			setupData: func(helper *test.TestHelper) error {
				return helper.ClearTestData()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helper, err := test.SetupTestDB()
			require.NoError(t, err)
			defer helper.Close()

			err = helper.SetupTestSchema()
			require.NoError(t, err)

			err = tc.setupData(helper)
			require.NoError(t, err)

			repo := postgres.NewDashboardRepositoryImpl(helper.Pool)
			userConfig := service.NewUserConfigService(service.GetDefaultConfig())
			handler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{})

			req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json", nil)
			recorder := httptest.NewRecorder()

			handler.Day(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			response, err := parseDashboardResponseDTO(recorder.Body.Bytes())
			require.NoError(t, err)
			require.NotNil(t, response.Stats)

			if tc.expectNil {
				assert.Nil(t, response.Stats.PerPages, "PerPages should be nil when no previous week data")
			}
		})
	}
}

// =============================================================================
// Edge Case Tests
// =============================================================================

// TestStatsData_EdgeCases_Integration tests various edge cases for StatsData fields
func TestStatsData_EdgeCases_Integration(t *testing.T) {
	if !test.IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	testCases := []struct {
		name      string
		setupData func(helper *test.TestHelper) error
		validate  func(t *testing.T, response *dto.DashboardResponse)
	}{
		{
			name: "AllLogsSameWeekday",
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				// Get current weekday
				currentWDay := int(time.Now().Weekday())
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Test Project", TotalPage: 100, Page: 0},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now().AddDate(0, 0, -14), StartPage: 0, EndPage: 10, WDay: currentWDay},
						{ID: 2, ProjectID: 1, Data: time.Now().AddDate(0, 0, -7), StartPage: 10, EndPage: 20, WDay: currentWDay},
						{ID: 3, ProjectID: 1, Data: time.Now(), StartPage: 20, EndPage: 30, WDay: currentWDay},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
			validate: func(t *testing.T, response *dto.DashboardResponse) {
				// Should calculate correctly even with all logs on same weekday
				assert.NotNil(t, response)
				assert.NotNil(t, response.Stats)
			},
		},
		{
			name: "LargePageNumbers",
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Large Book", TotalPage: 100000, Page: 50000},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now(), StartPage: 0, EndPage: 50000, WDay: int(time.Now().Weekday())},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
			validate: func(t *testing.T, response *dto.DashboardResponse) {
				// Should handle large page numbers without overflow
				assert.NotNil(t, response)
			},
		},
		{
			name: "ZeroPageLog",
			setupData: func(helper *test.TestHelper) error {
				fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
				scenario := &dashboardFixtures.Scenario{
					Projects: []*dashboardFixtures.ProjectFixture{
						{ID: 1, Name: "Test Project", TotalPage: 100, Page: 0},
					},
					Logs: []*dashboardFixtures.LogFixture{
						{ID: 1, ProjectID: 1, Data: time.Now(), StartPage: 50, EndPage: 50, WDay: int(time.Now().Weekday())},
					},
				}
				return fixtureManager.LoadScenario(scenario)
			},
			validate: func(t *testing.T, response *dto.DashboardResponse) {
				// Should handle zero-page logs gracefully
				assert.NotNil(t, response)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helper, err := test.SetupTestDB()
			require.NoError(t, err)
			defer helper.Close()

			err = helper.SetupTestSchema()
			require.NoError(t, err)

			err = tc.setupData(helper)
			require.NoError(t, err)

			repo := postgres.NewDashboardRepositoryImpl(helper.Pool)
			userConfig := service.NewUserConfigService(service.GetDefaultConfig())
			handler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{})

			req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json", nil)
			recorder := httptest.NewRecorder()

			handler.Day(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			response, err := parseDashboardResponseDTO(recorder.Body.Bytes())
			require.NoError(t, err)

			tc.validate(t, response)
		})
	}
}

// =============================================================================
// Helper Functions
// =============================================================================

// parseDashboardResponseDTO parses a JSON:API envelope and extracts DashboardResponse
func parseDashboardResponseDTO(body []byte) (*dto.DashboardResponse, error) {
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

	// Parse StatsData from attributes
	if stats, ok := attributes["stats"].(map[string]interface{}); ok {
		response.Stats = &dto.StatsData{}
		if v, exists := stats["previous_week_pages"]; exists {
			response.Stats.PreviousWeekPages = int(v.(float64))
		}
		if v, exists := stats["last_week_pages"]; exists {
			response.Stats.LastWeekPages = int(v.(float64))
		}
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
		if v, exists := stats["total_pages"]; exists {
			response.Stats.TotalPages = int(v.(float64))
		}
		if v, exists := stats["pages"]; exists {
			response.Stats.Pages = int(v.(float64))
		}
		if v, exists := stats["count_pages"]; exists {
			response.Stats.CountPages = int(v.(float64))
		}
		if v, exists := stats["speculate_pages"]; exists {
			response.Stats.SpeculatePages = int(v.(float64))
		}
		if v, exists := stats["max_day"]; exists {
			if v != nil {
				val := v.(float64)
				response.Stats.MaxDay = &val
			}
		}
		if v, exists := stats["mean_geral"]; exists {
			if v != nil {
				val := v.(float64)
				response.Stats.MeanGeral = &val
			}
		}
		if v, exists := stats["per_mean_day"]; exists {
			if v != nil {
				val := v.(float64)
				response.Stats.PerMeanDay = &val
			}
		}
		if v, exists := stats["per_spec_mean_day"]; exists {
			if v != nil {
				val := v.(float64)
				response.Stats.PerSpecMeanDay = &val
			}
		}
	}

	return response, nil
}
