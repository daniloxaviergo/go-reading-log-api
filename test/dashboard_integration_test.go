package test

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pg "go-reading-log-api-next/internal/adapter/postgres"
	"go-reading-log-api-next/internal/api/v1/handlers"
	"go-reading-log-api-next/internal/domain/dto"
	"go-reading-log-api-next/internal/repository"
	"go-reading-log-api-next/internal/service"
	"go-reading-log-api-next/internal/service/dashboard"
	dashboardFixtures "go-reading-log-api-next/test/fixtures/dashboard"
)

// MockProjectsService is a mock implementation of ProjectsServiceInterface
// It returns the scenario data that was loaded into the database
type MockProjectsService struct {
	scenario *dashboardFixtures.Scenario
	pool     *pgxpool.Pool
}

func (m *MockProjectsService) GetRunningProjectsWithLogs(ctx context.Context) ([]*dashboard.ProjectWithLogs, error) {
	if m.pool == nil || m.scenario == nil {
		return []*dashboard.ProjectWithLogs{}, nil
	}

	// Query actual data from database to match what the real service would return
	var projects []*dashboard.ProjectWithLogs

	for _, proj := range m.scenario.Projects {
		if proj.Status != "running" {
			continue
		}

		// Get project data from database
		var page, totalPage int
		err := m.pool.QueryRow(ctx, "SELECT page, total_page FROM projects WHERE id = $1", proj.ID).Scan(&page, &totalPage)
		if err != nil {
			continue
		}

		// Get logs for this project (limit 4, ordered by date DESC)
		var logs []*dto.LogEntry
		rows, err := m.pool.Query(ctx, `
			SELECT id, project_id, TO_CHAR(data, 'YYYY-MM-DD"T"HH24:MI:SS'), start_page, end_page, COALESCE(note, '')
			FROM logs WHERE project_id = $1 ORDER BY data DESC LIMIT 4
		`, proj.ID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var log dto.LogEntry
				var note string
				err := rows.Scan(&log.ID, &log.ProjectID, &log.Data, &log.StartPage, &log.EndPage, &note)
				if err == nil {
					log.Note = &note
					logs = append(logs, &log)
				}
			}
		}

		// Calculate progress
		progress := 0.0
		if totalPage > 0 {
			progress = math.Round(float64(page)/float64(totalPage)*100*1000) / 1000
		}

		project := &dashboard.ProjectWithLogs{
			Project: &dto.ProjectAggregateResponse{
				ProjectID:   proj.ID,
				ProjectName: proj.Name,
				TotalPages:  totalPage,
				LogCount:    len(logs),
				Progress:    progress,
			},
			Logs:       logs,
			TotalPages: totalPage,
			Pages:      page,
			Progress:   progress,
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (m *MockProjectsService) CalculateStats(ctx context.Context) (*dto.StatsData, error) {
	if m.pool == nil || m.scenario == nil {
		return dto.NewStatsData(), nil
	}

	// Calculate actual stats from database
	var stats = dto.NewStatsData()

	// Get total pages and count from all projects
	var totalCapacity, totalPage int
	err := m.pool.QueryRow(ctx, "SELECT COALESCE(SUM(total_page), 0), COALESCE(SUM(page), 0) FROM projects").Scan(&totalCapacity, &totalPage)
	if err == nil && totalCapacity > 0 {
		stats.ProgressGeral = math.Round(float64(totalPage)/float64(totalCapacity)*100*1000) / 1000
		stats.TotalPages = totalPage
	}

	return stats, nil
}

func (m *MockProjectsService) GetDashboardProjects(ctx context.Context) (*dto.DashboardProjectsResponse, error) {
	if m.pool == nil {
		return dto.NewDashboardProjectsResponse(), nil
	}

	response := dto.NewDashboardProjectsResponse()

	// Get running projects with logs from database
// Rails logic: running = page != total_page (not finished)
	rows, err := m.pool.Query(ctx, `
		SELECT DISTINCT p.id, p.name, p.total_page, p.page
		FROM projects p
		INNER JOIN logs l ON p.id = l.project_id
		WHERE p.page != p.total_page
		ORDER BY p.id
	`)
	if err != nil {
		return response, nil
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string
		var totalPage, page int
		if err := rows.Scan(&id, &name, &totalPage, &page); err != nil {
			continue
		}

		// Calculate progress
		progress := 0.0
		if totalPage > 0 {
			progress = math.Round(float64(page)/float64(totalPage)*100*1000) / 1000
		}

		// Get logs count
		var logsCount int
		m.pool.QueryRow(ctx, "SELECT COUNT(*) FROM logs WHERE project_id = $1", id).Scan(&logsCount)

		// Get started-at from earliest log
		var startedAt *string
		var earliestDate string
		if err := m.pool.QueryRow(ctx, "SELECT TO_CHAR(MIN(data), 'YYYY-MM-DD\"T\"HH24:MI:SS') FROM logs WHERE project_id = $1", id).Scan(&earliestDate); err == nil && earliestDate != "" {
			startedAt = &earliestDate
		}

		// Get days-unreading from latest log
		var daysUnreading int
		var latestDate string
		if err := m.pool.QueryRow(ctx, "SELECT TO_CHAR(MAX(data), 'YYYY-MM-DD\"T\"HH24:MI:SS') FROM logs WHERE project_id = $1", id).Scan(&latestDate); err == nil && latestDate != "" {
			// Calculate days since latest log
			if logTime, err := time.Parse(time.RFC3339, latestDate); err == nil {
				today := dto.GetToday()
				daysUnreading = int(today.Sub(logTime).Hours() / 24)
				if daysUnreading < 0 {
					daysUnreading = 0
				}
			}
		}

		// Create attributes
		attributes := dto.NewDashboardProjectAttributes()
		attributes.Name = name
		attributes.TotalPage = totalPage
		attributes.Page = page
		attributes.Progress = progress
		attributes.LogsCount = logsCount
		attributes.Status = "stopped"
		attributes.DaysUnreading = daysUnreading
		if startedAt != nil {
			attributes.SetStartedAt(*startedAt)
		}

		// Add project to response
		response.AddProject(*dto.NewDashboardProjectItem(fmt.Sprintf("%d", id), attributes))
	}

	// Calculate stats
	var totalCapacity, totalPage int
	err = m.pool.QueryRow(ctx, "SELECT COALESCE(SUM(total_page), 0), COALESCE(SUM(page), 0) FROM projects").Scan(&totalCapacity, &totalPage)
	stats := dto.NewDashboardStats()
	stats.SetTotalPages(totalPage)
	stats.SetPages(totalPage) // Simplified for mock
	if totalCapacity > 0 {
		stats.SetProgressGeral(math.Round(float64(totalPage)/float64(totalCapacity)*100*1000) / 1000)
	}
	response.SetStats(stats)

	return response, nil
}

// TestDashboardDayEndpoint_Integration tests /v1/dashboard/day.json
func TestDashboardDayEndpoint_Integration(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	helper, err := SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	// Create database tables before inserting data
	err = helper.SetupTestSchema()
	require.NoError(t, err)

	// Load test data
	fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
	scenario := dashboardFixtures.ScenarioMultipleProjects()

	err = fixtureManager.LoadScenario(scenario)
	require.NoError(t, err)

	// Create handler with real dependencies
	repo, err := createTestRepository(helper.Pool)
	require.NoError(t, err)

	userConfig, err := service.LoadDashboardConfig("")
	if err != nil {
		userConfig = service.NewUserConfigService(service.GetDefaultConfig())
	}

	dashboardHandler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{
		scenario: scenario,
		pool:     helper.Pool,
	})

	// Test GET /v1/dashboard/day.json
	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json", nil)
	recorder := httptest.NewRecorder()

	dashboardHandler.Day(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	// Parse response using helper function
	response, err := parseDashboardResponse(recorder.Body.Bytes())
	require.NoError(t, err)

	// Validate structure
	assert.NotNil(t, response.Stats)
	// PerPages should be nil when no previous period data is available (new behavior)
	assert.Nil(t, response.Stats.PerPages, "PerPages should be null when no previous period data")

	// Validate expected values
	err = validateExpectedValues(*response, scenario.Expected)
	require.NoError(t, err)
}

// TestDashboardLastDaysEndpoint_Integration tests /v1/dashboard/last_days.json
func TestDashboardLastDaysEndpoint_Integration(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	helper, err := SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	// Create database tables before inserting data
	err = helper.SetupTestSchema()
	require.NoError(t, err)

	fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)

	// Create scenario with logs in date range
	scenario := &dashboardFixtures.Scenario{
		Name:        "Last Days Test",
		Description: "Logs within last 30 days",
		Projects: []*dashboardFixtures.ProjectFixture{
			{ID: 1, Name: "Test Project", TotalPage: 200, Page: 50},
		},
		Logs: []*dashboardFixtures.LogFixture{
			// Create logs within last 7 days (for LastDays endpoint)
			{
				ID:        1,
				ProjectID: 1,
				Data:      time.Now().AddDate(0, 0, -3),
				StartPage: 0,
				EndPage:   25,
				WDay:      int(time.Now().AddDate(0, 0, -3).Weekday()),
			},
			{
				ID:        2,
				ProjectID: 1,
				Data:      time.Now().AddDate(0, 0, -2),
				StartPage: 25,
				EndPage:   50,
				WDay:      int(time.Now().AddDate(0, 0, -2).Weekday()),
			},
		},
	}

	err = fixtureManager.LoadScenario(scenario)
	require.NoError(t, err)

	repo, err := createTestRepository(helper.Pool)
	require.NoError(t, err)

	userConfig, err := service.LoadDashboardConfig("")
	if err != nil {
		userConfig = service.NewUserConfigService(service.GetDefaultConfig())
	}

	dashboardHandler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{
		scenario: scenario,
		pool:     helper.Pool,
	})

	// Test different type parameters
	typeParams := []string{"1", "2", "3", "4", "5"}

	for _, tp := range typeParams {
		t.Run("type_"+tp, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/last_days.json?type="+tp, nil)
			recorder := httptest.NewRecorder()

			dashboardHandler.LastDays(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			// Parse response using helper function
			response, err := parseDashboardResponse(recorder.Body.Bytes())
			require.NoError(t, err)

			// Verify logs are present and ordered
			assert.NotEmpty(t, response.Logs)
		})
	}

	// Test invalid type parameter
	t.Run("invalid_type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/last_days.json?type=99", nil)
		recorder := httptest.NewRecorder()

		dashboardHandler.LastDays(recorder, req)

		// Should return 422 for invalid type
		assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	})
}

// TestDashboardFaultsChart_Integration tests /v1/dashboard/echart/faults.json
func TestDashboardFaultsChart_Integration(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	helper, err := SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	// Create database tables before inserting data
	err = helper.SetupTestSchema()
	require.NoError(t, err)

	fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
	scenario := dashboardFixtures.ScenarioFaultsByWeekday()

	err = fixtureManager.LoadScenario(scenario)
	require.NoError(t, err)

	repo, err := createTestRepository(helper.Pool)
	require.NoError(t, err)

	userConfig, err := service.LoadDashboardConfig("")
	if err != nil {
		userConfig = service.NewUserConfigService(service.GetDefaultConfig())
	}

	dashboardHandler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{
		scenario: scenario,
		pool:     helper.Pool,
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/faults.json", nil)
	recorder := httptest.NewRecorder()

	dashboardHandler.Faults(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	// Parse response using helper function
	response, err := parseDashboardResponse(recorder.Body.Bytes())
	require.NoError(t, err)

	// Verify echart configuration
	assert.NotNil(t, response.Echart)
	assert.NotEmpty(t, response.Echart.Series)

	// Verify gauge chart structure
	series := response.Echart.Series[0]
	assert.Equal(t, "gauge", series.Type)
	assert.Len(t, series.Data, 1)

	// Verify percentage value
	percentage := series.Data[0].(float64)
	assert.GreaterOrEqual(t, percentage, 0.0)
	assert.LessOrEqual(t, percentage, 100.0)
}

// TestDashboardSpeculateActual_Integration tests /v1/dashboard/echart/speculate_actual.json
func TestDashboardSpeculateActual_Integration(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	helper, err := SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	// Create database tables before inserting data
	err = helper.SetupTestSchema()
	require.NoError(t, err)

	fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)

	// Scenario with logs across 15 days
	scenario := &dashboardFixtures.Scenario{
		Name:        "Speculate Actual Test",
		Description: "Logs for speculate vs actual chart",
		Projects: []*dashboardFixtures.ProjectFixture{
			{ID: 1, Name: "Test Project", TotalPage: 200, Page: 50},
		},
		Logs: func() []*dashboardFixtures.LogFixture {
			var logs []*dashboardFixtures.LogFixture
			baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

			// Create logs for last 15 days
			for i := 0; i < 15; i++ {
				logs = append(logs, &dashboardFixtures.LogFixture{
					ID:        int64(i + 1),
					ProjectID: 1,
					Data:      baseDate.AddDate(0, 0, -i),
					StartPage: i * 10,
					EndPage:   (i + 1) * 10,
					WDay:      int(baseDate.AddDate(0, 0, -i).Weekday()),
				})
			}
			return logs
		}(),
	}

	err = fixtureManager.LoadScenario(scenario)
	require.NoError(t, err)

	repo, err := createTestRepository(helper.Pool)
	require.NoError(t, err)

	userConfig, err := service.LoadDashboardConfig("")
	if err != nil {
		userConfig = service.NewUserConfigService(service.GetDefaultConfig())
	}

	dashboardHandler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{
		scenario: scenario,
		pool:     helper.Pool,
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/speculate_actual.json", nil)
	recorder := httptest.NewRecorder()

	dashboardHandler.SpeculateActual(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	// Parse response using helper function
	response, err := parseDashboardResponse(recorder.Body.Bytes())
	require.NoError(t, err)

	// Verify line chart with 15 data points
	assert.NotNil(t, response.Echart)
	assert.NotEmpty(t, response.Echart.Series)

	// Should have actual and speculate series
	assert.Len(t, response.Echart.Series, 2)

	for _, s := range response.Echart.Series {
		assert.Equal(t, "line", s.Type)
		assert.Len(t, s.Data, 15) // 15 days
	}
}

// TestDashboardWeekdayFaults_Integration tests /v1/dashboard/echart/faults_week_day.json
func TestDashboardWeekdayFaults_Integration(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	helper, err := SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	// Create database tables before inserting data
	err = helper.SetupTestSchema()
	require.NoError(t, err)

	fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
	scenario := dashboardFixtures.ScenarioFaultsByWeekday()

	err = fixtureManager.LoadScenario(scenario)
	require.NoError(t, err)

	repo, err := createTestRepository(helper.Pool)
	require.NoError(t, err)

	userConfig, err := service.LoadDashboardConfig("")
	if err != nil {
		userConfig = service.NewUserConfigService(service.GetDefaultConfig())
	}

	dashboardHandler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{
		scenario: scenario,
		pool:     helper.Pool,
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/faults_week_day.json", nil)
	recorder := httptest.NewRecorder()

	dashboardHandler.WeekdayFaults(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	// Parse response using helper function
	response, err := parseDashboardResponse(recorder.Body.Bytes())
	require.NoError(t, err)

	// Verify radar chart structure
	assert.NotNil(t, response.Echart)
	assert.NotEmpty(t, response.Echart.Series)

	series := response.Echart.Series[0]
	assert.Equal(t, "radar", series.Type)

	// Should have 7 data points (one for each weekday)
	assert.Len(t, series.Data, 7)

	// Verify fault counts for each weekday (should match the fixture distribution)
	// The fixture creates: Sunday=3, Monday=2, Tuesday=1, Wednesday=2, Thursday=1, Friday=3, Saturday=4
	expectedFaults := []int{3, 2, 1, 2, 1, 3, 4}
	actualFaults := make([]int, len(series.Data))
	for i, v := range series.Data {
		actualFaults[i] = int(v.(float64)) // Value represents fault count for weekday i
	}

	assert.Equal(t, expectedFaults, actualFaults)
}

// TestDashboardMeanProgress_Integration tests /v1/dashboard/echart/mean_progress.json
func TestDashboardMeanProgress_Integration(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	helper, err := SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	// Create database tables before inserting data
	err = helper.SetupTestSchema()
	require.NoError(t, err)

	fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)

	// Scenario with varying daily progress
	scenario := &dashboardFixtures.Scenario{
		Name:        "Mean Progress Test",
		Description: "Varying daily progress for visual map testing",
		Projects: []*dashboardFixtures.ProjectFixture{
			{ID: 1, Name: "Test Project", TotalPage: 200, Page: 50},
		},
		Logs: func() []*dashboardFixtures.LogFixture {
			var logs []*dashboardFixtures.LogFixture
			baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

			// Create logs with varying page counts

			for i := 0; i < 30; i++ {
				logs = append(logs, &dashboardFixtures.LogFixture{
					ID:        int64(i + 1),
					ProjectID: 1,
					Data:      baseDate.AddDate(0, 0, -i),
					StartPage: (i * 5) % 100,
					EndPage:   ((i + 1) * 5) % 100,
					WDay:      int(baseDate.AddDate(0, 0, -i).Weekday()),
				})
			}
			return logs
		}(),
	}

	err = fixtureManager.LoadScenario(scenario)
	require.NoError(t, err)

	repo, err := createTestRepository(helper.Pool)
	require.NoError(t, err)

	userConfig, err := service.LoadDashboardConfig("")
	if err != nil {
		userConfig = service.NewUserConfigService(service.GetDefaultConfig())
	}

	dashboardHandler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{
		scenario: scenario,
		pool:     helper.Pool,
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/mean_progress.json", nil)
	recorder := httptest.NewRecorder()

	dashboardHandler.MeanProgress(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	// Parse response using helper function
	response, err := parseDashboardResponse(recorder.Body.Bytes())
	require.NoError(t, err)

	// Verify line chart with 30 data points
	assert.NotNil(t, response.Echart)
	assert.NotEmpty(t, response.Echart.Series)

	series := response.Echart.Series[0]
	assert.Equal(t, "line", series.Type)
	assert.Len(t, series.Data, 30)

	// Verify color assignments (visual map)
	// Check that colors are assigned based on progress ranges
	for _, point := range series.Data {
		if dataPoint, ok := point.(map[string]interface{}); ok {
			if value, exists := dataPoint["value"]; exists {
				progress := value.(float64)

				// Verify color is assigned
				if color, exists := dataPoint["itemStyle"].(map[string]interface{})["color"]; exists {
					assert.NotEmpty(t, color)

					// Verify color matches expected range
					verifyColorForProgress(t, progress, color.(string))
				}
			}
		}
	}
}

// TestDashboardYearlyTotal_Integration tests /v1/dashboard/echart/last_year_total.json
func TestDashboardYearlyTotal_Integration(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	helper, err := SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	// Create database tables before inserting data
	err = helper.SetupTestSchema()
	require.NoError(t, err)

	fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)

	// Scenario with logs spanning 52 weeks
	scenario := &dashboardFixtures.Scenario{
		Name:        "Yearly Total Test",
		Description: "Logs spanning 52 weeks for yearly trend",
		Projects: []*dashboardFixtures.ProjectFixture{
			{ID: 1, Name: "Test Project", TotalPage: 200, Page: 50},
		},
		Logs: func() []*dashboardFixtures.LogFixture {
			var logs []*dashboardFixtures.LogFixture
			baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

			// Create logs for each week over 52 weeks
			for w := 0; w < 52; w++ {
				weekStart := baseDate.AddDate(0, 0, w*7)

				// Create multiple logs per week
				for d := 0; d < 7; d++ {
					logs = append(logs, &dashboardFixtures.LogFixture{
						ID:        int64(w*7 + d + 1),
						ProjectID: 1,
						Data:      weekStart.AddDate(0, 0, d),
						StartPage: (w * 10) % 100,
						EndPage:   ((w + 1) * 10) % 100,
						WDay:      int(weekStart.AddDate(0, 0, d).Weekday()),
					})
				}
			}
			return logs
		}(),
	}

	err = fixtureManager.LoadScenario(scenario)
	require.NoError(t, err)

	repo, err := createTestRepository(helper.Pool)
	require.NoError(t, err)

	userConfig, err := service.LoadDashboardConfig("")
	if err != nil {
		userConfig = service.NewUserConfigService(service.GetDefaultConfig())
	}

	dashboardHandler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{
		scenario: scenario,
		pool:     helper.Pool,
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/last_year_total.json", nil)
	recorder := httptest.NewRecorder()

	dashboardHandler.YearlyTotal(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	// Parse response using helper function
	response, err := parseDashboardResponse(recorder.Body.Bytes())
	require.NoError(t, err)

	// Verify line chart with 52 data points (weekly aggregates)
	assert.NotNil(t, response.Echart)
	assert.NotEmpty(t, response.Echart.Series)

	series := response.Echart.Series[0]
	assert.Equal(t, "line", series.Type)
	assert.Len(t, series.Data, 52)

	// Verify each data point has week boundaries
	for _, point := range series.Data {
		if dataPoint, ok := point.(map[string]interface{}); ok {
			assert.Contains(t, dataPoint, "begin_week")
			assert.Contains(t, dataPoint, "end_week")
			assert.Contains(t, dataPoint, "count_reads")
		}
	}
}

// TestDashboardEndpoints_ErrorHandling tests error scenarios
func TestDashboardEndpoints_ErrorHandling(t *testing.T) {
	if !IsTestDatabase() {
		t.Skip("Test database not configured")
	}

	helper, err := SetupTestDB()
	require.NoError(t, err)
	defer helper.Close()

	// Create database tables before inserting data
	err = helper.SetupTestSchema()
	require.NoError(t, err)

	// Clear data to test empty state
	err = helper.ClearTestData()
	require.NoError(t, err)

	repo, err := createTestRepository(helper.Pool)
	require.NoError(t, err)

	userConfig, err := service.LoadDashboardConfig("")
	if err != nil {
		userConfig = service.NewUserConfigService(service.GetDefaultConfig())
	}

	// Use empty scenario for error handling tests
	emptyScenario := &dashboardFixtures.Scenario{
		Projects: []*dashboardFixtures.ProjectFixture{},
		Logs:     []*dashboardFixtures.LogFixture{},
	}

	dashboardHandler := handlers.NewDashboardHandler(repo, userConfig, &MockProjectsService{
		scenario: emptyScenario,
		pool:     helper.Pool,
	})

	testCases := []struct {
		name     string
		endpoint string
		method   string
	}{
		{"Day Endpoint Empty", "/v1/dashboard/day.json", "GET"},
		{"Projects Endpoint Empty", "/v1/dashboard/projects.json", "GET"},
		{"Last Days Invalid Type", "/v1/dashboard/last_days.json?type=invalid", "GET"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.endpoint, nil)
			recorder := httptest.NewRecorder()

			// Route to appropriate handler
			switch {
			case tc.endpoint == "/v1/dashboard/day.json":
				dashboardHandler.Day(recorder, req)
			case tc.endpoint == "/v1/dashboard/projects.json":
				dashboardHandler.Projects(recorder, req)
			case tc.endpoint == "/v1/dashboard/last_days.json?type=invalid":
				dashboardHandler.LastDays(recorder, req)
			}

			// For "Last Days Invalid Type", expect 422 Unprocessable Entity
			if tc.name == "Last Days Invalid Type" {
				assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
				return // Don't try to parse error response
			}

			// For empty data, should still return 200 with zero values
			assert.Equal(t, http.StatusOK, recorder.Code)

			// Parse response using helper function
			response, err := parseDashboardResponse(recorder.Body.Bytes())
			require.NoError(t, err)

			// Verify zero/empty responses are valid
			if tc.endpoint == "/v1/dashboard/day.json" {
				assert.NotNil(t, response.Stats)
				assert.Equal(t, 0, response.Stats.PreviousWeekPages)
				assert.Equal(t, 0, response.Stats.LastWeekPages)
			}
		})
	}
}

// validateExpectedValues compares actual response with expected values
func validateExpectedValues(response dto.DashboardResponse, expected *dashboardFixtures.ExpectedResults) error {
	if expected == nil {
		return nil
	}

	if expected.Stats != nil {
		if response.Stats != nil {
			// Allow small floating point differences for PerPages (pointer type)
			if response.Stats.PerPages != nil && expected.Stats.PerPages != nil {
				if !floatEqual(*response.Stats.PerPages, *expected.Stats.PerPages, 0.001) {
					return fmt.Errorf("per_pages mismatch: got %f, expected %f",
						*response.Stats.PerPages, *expected.Stats.PerPages)
				}
			}
			if !floatEqual(response.Stats.MeanDay, expected.Stats.MeanDay, 0.001) {
				return fmt.Errorf("mean_day mismatch: got %f, expected %f",
					response.Stats.MeanDay, expected.Stats.MeanDay)
			}
			if !floatEqual(response.Stats.SpecMeanDay, expected.Stats.SpecMeanDay, 0.001) {
				return fmt.Errorf("spec_mean_day mismatch: got %f, expected %f",
					response.Stats.SpecMeanDay, expected.Stats.SpecMeanDay)
			}
			if !floatEqual(response.Stats.ProgressGeral, expected.Stats.ProgressGeral, 0.001) {
				return fmt.Errorf("progress_geral mismatch: got %f, expected %f",
					response.Stats.ProgressGeral, expected.Stats.ProgressGeral)
			}
		}
	}

	return nil
}

func floatEqual(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}

// verifyColorForProgress determines expected color based on progress range
func verifyColorForProgress(t *testing.T, progress float64, actualColor string) {
	var expectedColor string

	switch {
	case progress >= 0 && progress < 10:
		expectedColor = "#95a5a6" // gray
	case progress >= 10 && progress < 20:
		expectedColor = "#1abc9c" // cyan
	case progress >= 20 && progress < 50:
		expectedColor = "#3498db" // blue
	case progress >= 50:
		expectedColor = "#2ecc71" // green
	default: // negative
		expectedColor = "#e74c3c" // red
	}

	assert.Equal(t, expectedColor, actualColor,
		"Color mismatch for progress %.1f%%: got %s, expected %s",
		progress, actualColor, expectedColor)
}

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

	// Parse Echart - fields are directly in attributes
	response.Echart = &dto.EchartConfig{}
	if title, exists := attributes["title"]; exists {
		response.Echart.Title = title.(string)
	}
	if tooltip, exists := attributes["tooltip"]; exists {
		response.Echart.Tooltip = tooltip.(map[string]interface{})
	}
	if series, exists := attributes["series"].([]interface{}); exists {
		response.Echart.Series = make([]dto.Series, len(series))
		for i, s := range series {
			seriesMap := s.(map[string]interface{})
			itemStyle := make(map[string]interface{})
			if style, ok := seriesMap["itemStyle"].(map[string]interface{}); ok {
				itemStyle = style
			}
			response.Echart.Series[i] = dto.Series{
				Name:      seriesMap["name"].(string),
				Type:      seriesMap["type"].(string),
				Data:      seriesMap["data"].([]interface{}),
				ItemStyle: itemStyle,
			}
		}
	}

	return response, nil
}

// createTestRepository creates a dashboard repository for testing
func createTestRepository(pool *pgxpool.Pool) (repository.DashboardRepository, error) {
	return pg.NewDashboardRepositoryImpl(pool), nil
}
