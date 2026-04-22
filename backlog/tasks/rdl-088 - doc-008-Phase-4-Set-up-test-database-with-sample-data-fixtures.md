---
id: RDL-088
title: '[doc-008 Phase 4] Set up test database with sample data fixtures'
status: To Do
assignee:
  - thomas
created_date: '2026-04-21 15:51'
updated_date: '2026-04-22 11:11'
labels:
  - phase-4
  - testing
  - fixtures
dependencies: []
references:
  - NFA-DASH-001
  - IT-001
  - Acceptance Criteria All
documentation:
  - doc-008
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create comprehensive test fixtures for dashboard testing. Include scenarios covering edge cases: zero pages, null dates, multiple projects, varying completion levels, faults across different weekdays, and logs spanning required date ranges.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Test fixtures for all dashboard scenarios created
- [ ] #2 Edge cases covered (zero pages, null dates)
- [ ] #3 Multiple projects with varying completion levels
- [ ] #4 Faults distributed across different weekdays
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task focuses on creating comprehensive test fixtures for dashboard testing, covering all edge cases and scenarios defined in doc-008 Phase 4.

**Architecture Strategy:**
- Create a fixture management system that can populate test databases with controlled data
- Implement scenario-based fixtures for each dashboard endpoint type
- Ensure fixtures cover all acceptance criteria from AC-DASH-001 through AC-DASH-008
- Use the existing `TestHelper` infrastructure for database setup/teardown

**Fixture Categories:**
1. **Basic Project Fixtures** - Single project with varied completion levels
2. **Multiple Project Fixtures** - 3-5 projects with different states (unstarted, running, sleeping, stopped, finished)
3. **Edge Case Fixtures** - Zero pages, null dates, single log entries
4. **Fault Distribution Fixtures** - Logs spanning different weekdays across 6-month period
5. **Trend Data Fixtures** - Logs covering 7/15/30/90 day periods for last_days endpoint

**Implementation Pattern:**
```
test/data/fixtures/
├── dashboard/
│   ├── fixtures.go          # Fixture manager and helper functions
│   ├── projects.go          # Project fixture generators
│   ├── logs.go              # Log fixture generators
│   ├── faults.go            # Fault distribution fixtures
│   └── scenarios.go         # Pre-built scenario combinations
└── testdata/
    ├── project-450-go.json  # Existing reference data
    └── expected-values.yaml # New: Expected calculation results
```

**Key Design Decisions:**
- Use `time.Time` with specific dates (not "today") for reproducible tests
- Create fixtures that allow incremental building (start simple, add complexity)
- Support both unit test mocks and integration test database population
- Document expected values alongside fixtures for validation

---

### 2. Files to Modify

#### New Files to Create

| File Path | Purpose | Lines | Priority |
|-----------|---------|-------|----------|
| `test/data/fixtures/dashboard/fixtures.go` | Core fixture manager with setup/teardown helpers | ~150 | P1 |
| `test/data/fixtures/dashboard/projects.go` | Project fixture generators (single/multiple) | ~200 | P1 |
| `test/data/fixtures/dashboard/logs.go` | Log fixture generators with date control | ~250 | P1 |
| `test/data/fixtures/dashboard/faults.go` | Fault distribution fixtures by weekday | ~180 | P2 |
| `test/data/fixtures/dashboard/scenarios.go` | Pre-built test scenarios | ~300 | P2 |
| `test/data/fixtures/testdata/expected-values.yaml` | Expected calculation results for validation | ~100 | P1 |
| `test/dashboard_integration_test.go` | Integration tests using fixtures (replaces empty dir) | ~400 | P2 |

#### Files to Extend

| File Path | Extension | Priority |
|-----------|-----------|----------|
| `test/test_helper.go` | Add `CreateTestProjectWithLogs()` helper | P1 |
| `internal/domain/dto/dashboard_response.go` | Add `ValidateExpectedValues()` method | P2 |

---

### 3. Dependencies

**Prerequisites:**
- [x] Phase 1 core API complete (projects, logs endpoints)
- [x] Dashboard services implemented (`DayService`, `ProjectsService`, etc.)
- [x] Repository interface defined (`DashboardRepository`)
- [x] Test infrastructure exists (`TestHelper`, `SetupTestDB`)

**External Dependencies:**
- `github.com/stretchr/testify/assert` - Already in use
- `github.com/stretchr/testify/require` - Already in use
- `gopkg.in/yaml.v3` - Already in use (for expected values)

**No New Dependencies Required**

---

### 4. Code Patterns

#### Pattern 1: Fixture Manager

```go
// test/data/fixtures/dashboard/fixtures.go
package fixtures

import (
    "time"
    "go-reading-log-api-next/test"
)

// DashboardFixtures manages all dashboard-related test data
type DashboardFixtures struct {
    helper *test.TestHelper
    
    // Pre-created data for reuse across tests
    projects map[int64]*ProjectFixture
    logs     map[int64][]*LogFixture
}

// ProjectFixture represents a single project with its associated data
type ProjectFixture struct {
    ID        int64
    Name      string
    TotalPage int
    Page      int
    StartedAt *time.Time
    Status    string
    
    // Associated logs
    Logs []*LogFixture
}

// LogFixture represents a log entry with full control over fields
type LogFixture struct {
    ID        int64
    ProjectID int64
    Data      time.Time  // Explicit date for reproducibility
    StartPage int
    EndPage   int
    Note      *string
    WDay      int  // Weekday (0-6)
}

// NewDashboardFixtures creates a new fixture manager
func NewDashboardFixtures(helper *test.TestHelper) *DashboardFixtures {
    return &DashboardFixtures{
        helper:   helper,
        projects: make(map[int64]*ProjectFixture),
        logs:     make(map[int64][]*LogFixture),
    }
}
```

#### Pattern 2: Scenario-Based Fixtures

```go
// test/data/fixtures/dashboard/scenarios.go
package fixtures

import (
    "time"
)

// Scenario represents a complete test scenario with all required data
type Scenario struct {
    Name        string
    Description string
    Projects    []*ProjectFixture
    Logs        []*LogFixture
    Expected    *ExpectedResults  // Calculated values for validation
}

// ExpectedResults holds pre-calculated values for validation
type ExpectedResults struct {
    Stats       *StatsExpectations
    EchartData  map[string]interface{}
    LogsCount   int
}

// StatsExpectations defines expected statistical values
type StatsExpectations struct {
    PreviousWeekPages int
    LastWeekPages     int
    PerPages          float64
    MeanDay           float64
    SpecMeanDay       float64
    ProgressGeral     float64
}

// Pre-built Scenarios

// ScenarioZeroPages: Project with zero pages (edge case)
func ScenarioZeroPages() *Scenario {
    return &Scenario{
        Name:        "Zero Pages Edge Case",
        Description: "Project with total_page=0 for edge case testing",
        Projects: []*ProjectFixture{
            {
                ID:        1,
                Name:      "Empty Project",
                TotalPage: 0,
                Page:      0,
                Status:    "unstarted",
            },
        },
        Logs: []*LogFixture{},
        Expected: &ExpectedResults{
            Stats: &StatsExpectations{
                PreviousWeekPages: 0,
                LastWeekPages:     0,
                PerPages:          0.0,
                MeanDay:           0.0,
                SpecMeanDay:       0.0,
                ProgressGeral:     0.0,
            },
        },
    }
}

// ScenarioCompleteBook: Fully completed project
func ScenarioCompleteBook() *Scenario {
    return &Scenario{
        Name:        "Complete Book",
        Description: "Project with all pages read",
        Projects: []*ProjectFixture{
            {
                ID:        2,
                Name:      "Completed Book",
                TotalPage: 300,
                Page:      300,
                Status:    "finished",
            },
        },
        Logs: []*LogFixture{
            {
                ID:        1,
                ProjectID: 2,
                Data:      time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
                StartPage: 0,
                EndPage:   300,
                WDay:      1, // Monday
            },
        },
        Expected: &ExpectedResults{
            Stats: &StatsExpectations{
                ProgressGeral: 100.0,
            },
        },
    }
}

// ScenarioMultipleProjects: Multiple projects with varying completion
func ScenarioMultipleProjects() *Scenario {
    return &Scenario{
        Name:        "Multiple Projects",
        Description: "3 projects: unstarted, running, finished",
        Projects: []*ProjectFixture{
            {
                ID:        10,
                Name:      "Unstarted Project",
                TotalPage: 200,
                Page:      0,
                Status:    "unstarted",
            },
            {
                ID:        11,
                Name:      "Running Project",
                TotalPage: 200,
                Page:      50,
                Status:    "running",
            },
            {
                ID:        12,
                Name:      "Finished Project",
                TotalPage: 200,
                Page:      200,
                Status:    "finished",
            },
        },
        Logs: []*LogFixture{
            // Logs for running project
            {ID: 100, ProjectID: 11, Data: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), StartPage: 0, EndPage: 25, WDay: 1},
            {ID: 101, ProjectID: 11, Data: time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC), StartPage: 25, EndPage: 50, WDay: 2},
        },
        Expected: &ExpectedResults{
            Stats: &StatsExpectations{
                ProgressGeral: 12.5, // (0 + 25 + 200) / (200 + 200 + 200) * 100
            },
        },
    }
}
```

#### Pattern 3: Fault Distribution Fixtures

```go
// test/data/fixtures/dashboard/faults.go
package fixtures

import (
    "time"
)

// FaultFixture represents a fault entry with weekday control
type FaultFixture struct {
    ID        int64
    ProjectID int64
    LogID     int64
    Date      time.Time
    WDay      int  // 0=Sunday, 6=Saturday
    Note      string
}

// FaultDistributionConfig defines expected fault distribution
type FaultDistributionConfig struct {
    TotalFaults   int
    ByWeekday     map[int]int  // Map of weekday -> count
    DateRange     struct {
        Start time.Time
        End   time.Time
    }
}

// ScenarioFaultsByWeekday: Faults distributed across all weekdays
func ScenarioFaultsByWeekday() *Scenario {
    // Define fault distribution (example: more faults on weekends)
    distribution := map[int]int{
        0: 3, // Sunday
        1: 2, // Monday
        2: 1, // Tuesday
        3: 2, // Wednesday
        4: 1, // Thursday
        5: 3, // Friday
        6: 4, // Saturday
    }
    
    totalFaults := 0
    for _, count := range distribution {
        totalFaults += count
    }
    
    // Generate fault fixtures
    var faults []*FaultFixture
    logID := int64(200)
    for weekday, count := range distribution {
        // Create multiple faults for this weekday
        for i := 0; i < count; i++ {
            // Calculate date within last 30 days for consistent testing
            baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
            targetDate := baseDate.AddDate(0, 0, -weekday)
            
            faults = append(faults, &FaultFixture{
                ID:        int64(logID),
                ProjectID: 1,
                LogID:     int64(logID),
                Date:      targetDate,
                WDay:      weekday,
                Note:      fmt.Sprintf("Fault on weekday %d", weekday),
            })
            logID++
        }
    }
    
    return &Scenario{
        Name:        "Faults by Weekday",
        Description: "Faults distributed across all 7 weekdays",
        Projects: []*ProjectFixture{
            {ID: 1, Name: "Faulty Project", TotalPage: 200, Page: 50},
        },
        Logs: func() []*LogFixture {
            var logs []*LogFixture
            for _, f := range faults {
                logs = append(logs, &LogFixture{
                    ID:        f.LogID,
                    ProjectID: f.ProjectID,
                    Data:      f.Date,
                    StartPage: 0,
                    EndPage:   10,
                    WDay:      f.WDay,
                })
            }
            return logs
        }(),
        Expected: &ExpectedResults{
            Stats: &StatsExpectations{
                // Fault percentage calculation
            },
        },
    }
}
```

#### Pattern 4: Integration Test Structure

```go
// test/dashboard_integration_test.go
package test

import (
    "net/http"
    "net/http/httptest"
    "testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
    "go-reading-log-api-next/internal/api/v1/handlers"
    "go-reading-log-api-next/internal/config"
    "go-reading-log-api-next/internal/adapter/postgres"
    "go-reading-log-api-next/internal/repository"
    "go-reading-log-api-next/test/data/fixtures/dashboard"
)

// TestDashboardEndpoints_Integration tests all 8 dashboard endpoints
func TestDashboardEndpoints_Integration(t *testing.T) {
    if !IsTestDatabase() {
        t.Skip("Test database not configured")
    }

    // Setup test database
    helper, err := SetupTestDB()
    require.NoError(t, err)
    defer helper.Close()

    // Load fixtures
    fixtureManager := dashboard.NewDashboardFixtures(helper)
    
    // Run each scenario
    scenarios := []struct {
        name     string
        scenario *dashboard.Scenario
        endpoint string
    }{
        {"Day Endpoint", dashboard.ScenarioMultipleProjects(), "/v1/dashboard/day.json"},
        {"Projects Endpoint", dashboard.ScenarioMultipleProjects(), "/v1/dashboard/projects.json"},
        {"Last Days (7)", dashboard.ScenarioMultipleProjects(), "/v1/dashboard/last_days.json?type=1"},
        {"Faults Chart", dashboard.ScenarioFaultsByWeekday(), "/v1/dashboard/echart/faults.json"},
        // Add remaining scenarios...
    }

    for _, sc := range scenarios {
        t.Run(sc.name, func(t *testing.T) {
            // Setup fixture data
            err := fixtureManager.LoadScenario(sc.scenario)
            require.NoError(t, err)

            // Make request
            req := httptest.NewRequest(http.MethodGet, sc.endpoint, nil)
            recorder := httptest.NewRecorder()
            
            // Get handler (need to expose via package or create here)
            cfg := config.LoadConfig()
            pool := helper.Pool
            repo := postgres.NewDashboardRepositoryImpl(pool)
            
            userConfig, err := service.LoadDashboardConfig("config/dashboard.yaml")
            if err != nil {
                userConfig = service.NewUserConfigService(service.GetDefaultConfig())
            }
            
            dashboardHandler := handlers.NewDashboardHandler(repo, userConfig)
            
            // Route and handle
            router := setupTestRouter(dashboardHandler)
            router.ServeHTTP(recorder, req)

            // Verify response
            assert.Equal(t, http.StatusOK, recorder.Code)
            
            // Parse and validate response
            var response dto.DashboardResponse
            err = json.Unmarshal(recorder.Body.Bytes(), &response)
            require.NoError(t, err)
            
            // Validate against expected values
            err = response.ValidateExpectedValues(sc.scenario.Expected)
            assert.NoError(t, err)
        })
    }
}

// setupTestRouter creates a test router with dashboard routes
func setupTestRouter(handler *handlers.DashboardHandler) *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/v1/dashboard/day.json", handler.GetDay).Methods("GET")
    router.HandleFunc("/v1/dashboard/projects.json", handler.GetProjects).Methods("GET")
    router.HandleFunc("/v1/dashboard/last_days.json", handler.GetLastDays).Methods("GET")
    router.HandleFunc("/v1/dashboard/echart/faults.json", handler.GetFaultsChart).Methods("GET")
    router.HandleFunc("/v1/dashboard/echart/speculate_actual.json", handler.GetSpeculateActualChart).Methods("GET")
    router.HandleFunc("/v1/dashboard/echart/faults_week_day.json", handler.GetWeekdayFaultsChart).Methods("GET")
    router.HandleFunc("/v1/dashboard/echart/mean_progress.json", handler.GetMeanProgressChart).Methods("GET")
    router.HandleFunc("/v1/dashboard/echart/last_year_total.json", handler.GetYearlyTotalChart).Methods("GET")
    return router
}
```

---

### 5. Testing Strategy

#### Unit Tests (test/unit/dashboard_service_test.go - extend existing)

**Coverage Goals:**
- All calculation methods with known input values
- Edge cases: zero values, null pointers, empty slices
- Float precision: verify 3-decimal rounding
- Error handling paths

**Test Categories:**

```go
// test/unit/dashboard_service_test.go (extensions)

// TestDayService_Calculations_KnownValues tests with predetermined inputs
func TestDayService_Calculations_KnownValues(t *testing.T) {
    // Setup: Create service with mock that returns known values
    mockRepo := &MockDashboardRepository{}
    mockConfig := &MockUserConfigService{mockPredictionPct: 0.15}
    
    // Define expected inputs
    knownAggregates := []*dto.ProjectAggregate{
        {ProjectID: 1, TotalPages: 100, LogCount: 50},
        {ProjectID: 2, TotalPages: 200, LogCount: 100},
    }
    
    mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
        return knownAggregates, nil
    }
    
    // Set up period page returns (deterministic)
    mockRepo.prevWeekPages = 75
    mockRepo.lastWeekPages = 100
    
    mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
        if projectID == 1 { return 8.5, nil }
        if projectID == 2 { return 9.5, nil }
        return 0, nil
    }
    
    // Execute
    service := dashboard.NewDayService(mockRepo, mockConfig)
    stats, err := service.CalculateWeeklyStats(context.Background())
    
    // Verify: Expected values calculated manually
    // previous_week_pages = 75
    // last_week_pages = 100
    // per_pages = (100/75)*100 = 133.333
    // mean_day = (8.5+9.5)/2 = 9.0
    // spec_mean_day = 9.0 * 1.15 = 10.35
    
    require.NoError(t, err)
    assert.Equal(t, 75, stats.PreviousWeekPages)
    assert.Equal(t, 100, stats.LastWeekPages)
    assert.Equal(t, 133.333, stats.PerPages)
    assert.Equal(t, 9.0, stats.MeanDay)
    assert.Equal(t, 10.35, stats.SpecMeanDay)
}

// TestDayService_EdgeCases tests edge case handling
func TestDayService_EdgeCases(t *testing.T) {
    testCases := []struct {
        name     string
        setup    func() (*MockDashboardRepository, *MockUserConfigService)
        expected *StatsExpectations
    }{
        {
            name: "Empty aggregates",
            setup: func() (*MockDashboardRepository, *MockUserConfigService) {
                mockRepo := &MockDashboardRepository{}
                mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
                    return []*dto.ProjectAggregate{}, nil
                }
                return mockRepo, &MockUserConfigService{mockPredictionPct: 0.15}
            },
            expected: &StatsExpectations{
                PreviousWeekPages: 0,
                LastWeekPages:     0,
                PerPages:          0.0,
                MeanDay:           0.0,
                SpecMeanDay:       0.0,
            },
        },
        {
            name: "Zero previous week",
            setup: func() (*MockDashboardRepository, *MockUserConfigService) {
                mockRepo := &MockDashboardRepository{}
                mockRepo.mockGetProjectAggregates = func(ctx context.Context) ([]*dto.ProjectAggregate, error) {
                    return []*dto.ProjectAggregate{{ProjectID: 1, TotalPages: 100, LogCount: 50}}, nil
                }
                mockRepo.prevWeekPages = 0
                mockRepo.lastWeekPages = 50
                return mockRepo, &MockUserConfigService{mockPredictionPct: 0.15}
            },
            expected: &StatsExpectations{
                PerPages: 0.0, // Should not divide by zero
            },
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            mockRepo, mockConfig := tc.setup()
            service := dashboard.NewDayService(mockRepo, mockConfig)
            
            stats, err := service.CalculateWeeklyStats(context.Background())
            
            require.NoError(t, err)
            assert.Equal(t, tc.expected.PreviousWeekPages, stats.PreviousWeekPages)
            assert.Equal(t, tc.expected.LastWeekPages, stats.LastWeekPages)
            assert.Equal(t, tc.expected.PerPages, stats.PerPages)
        })
    }
}
```

#### Integration Tests (test/dashboard_integration_test.go - new file)

**Coverage Goals:**
- All 8 dashboard endpoints with real database
- Full data pipeline: handler → service → repository → database
- Response format validation against JSON:API spec
- Error scenarios: missing data, invalid queries

**Test Structure:**

```go
// test/dashboard_integration_test.go

package test

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
    "go-reading-log-api-next/internal/api/v1/handlers"
    "go-reading-log-api-next/internal/config"
    "go-reading-log-api-next/internal/domain/dto"
    "go-reading-log-api-next/internal/service"
    "go-reading-log-api-next/test/data/fixtures/dashboard"
)

// TestDashboardDayEndpoint_Integration tests /v1/dashboard/day.json
func TestDashboardDayEndpoint_Integration(t *testing.T) {
    if !IsTestDatabase() {
        t.Skip("Test database not configured")
    }

    helper, err := SetupTestDB()
    require.NoError(t, err)
    defer helper.Close()

    // Load test data
    fixtureManager := dashboard.NewDashboardFixtures(helper)
    scenario := dashboard.ScenarioMultipleProjects()
    
    err = fixtureManager.LoadScenario(scenario)
    require.NoError(t, err)

    // Create handler with real dependencies
    cfg := config.LoadConfig()
    repo, err := createTestRepository(helper.Pool)
    require.NoError(t, err)
    
    userConfig, err := service.LoadDashboardConfig("")
    if err != nil {
        userConfig = service.NewUserConfigService(service.GetDefaultConfig())
    }
    
    dashboardHandler := handlers.NewDashboardHandler(repo, userConfig)

    // Test GET /v1/dashboard/day.json
    req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json", nil)
    recorder := httptest.NewRecorder()
    
    dashboardHandler.GetDay(recorder, req.Request)

    assert.Equal(t, http.StatusOK, recorder.Code)
    
    var response dto.DashboardResponse
    err = json.Unmarshal(recorder.Body.Bytes(), &response)
    require.NoError(t, err)
    
    // Validate structure
    assert.NotNil(t, response.Stats)
    assert.Equal(t, 133.333, response.Stats.PerPages) // From scenario expectations
    
    // Validate expected values
    err = validateExpectedValues(response, scenario.Expected)
    require.NoError(t, err)
}

// TestDashboardProjectsEndpoint_Integration tests /v1/dashboard/projects.json
func TestDashboardProjectsEndpoint_Integration(t *testing.T) {
    if !IsTestDatabase() {
        t.Skip("Test database not configured")
    }

    helper, err := SetupTestDB()
    require.NoError(t, err)
    defer helper.Close()

    fixtureManager := dashboard.NewDashboardFixtures(helper)
    scenario := dashboard.ScenarioMultipleProjects()
    
    err = fixtureManager.LoadScenario(scenario)
    require.NoError(t, err)

    repo, err := createTestRepository(helper.Pool)
    require.NoError(t, err)
    
    userConfig, err := service.LoadDashboardConfig("")
    if err != nil {
        userConfig = service.NewUserConfigService(service.GetDefaultConfig())
    }
    
    dashboardHandler := handlers.NewDashboardHandler(repo, userConfig)

    req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/projects.json", nil)
    recorder := httptest.NewRecorder()
    
    dashboardHandler.GetProjects(recorder, req.Request)

    assert.Equal(t, http.StatusOK, recorder.Code)
    
    var response dto.DashboardResponse
    err = json.Unmarshal(recorder.Body.Bytes(), &response)
    require.NoError(t, err)
    
    // Verify projects array exists and is sorted by progress
    assert.NotNil(t, response.Logs) // Projects returned as logs in current implementation
    
    // Validate expected values
    err = validateExpectedValues(response, scenario.Expected)
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

    fixtureManager := dashboard.NewDashboardFixtures(helper)
    
    // Create scenario with logs in date range
    scenario := &dashboard.Scenario{
        Name:        "Last Days Test",
        Description: "Logs within last 30 days",
        Projects: []*dashboard.ProjectFixture{
            {ID: 1, Name: "Test Project", TotalPage: 200, Page: 50},
        },
        Logs: []*dashboard.LogFixture{
            // Create logs within last 30 days
            {
                ID:        1,
                ProjectID: 1,
                Data:      time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
                StartPage: 0,
                EndPage:   25,
                WDay:      1,
            },
            {
                ID:        2,
                ProjectID: 1,
                Data:      time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC),
                StartPage: 25,
                EndPage:   50,
                WDay:      2,
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
    
    dashboardHandler := handlers.NewDashboardHandler(repo, userConfig)

    // Test different type parameters
    typeParams := []string{"1", "2", "3", "4", "5"}
    
    for _, tp := range typeParams {
        t.Run("type_"+tp, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/last_days.json?type="+tp, nil)
            recorder := httptest.NewRecorder()
            
            dashboardHandler.GetLastDays(recorder, req.Request)

            assert.Equal(t, http.StatusOK, recorder.Code)
            
            var response dto.DashboardResponse
            err = json.Unmarshal(recorder.Body.Bytes(), &response)
            require.NoError(t, err)
            
            // Verify logs are present and ordered
            assert.NotEmpty(t, response.Logs)
        })
    }
    
    // Test invalid type parameter
    t.Run("invalid_type", func(t *testing.T) {
        req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/last_days.json?type=99", nil)
        recorder := httptest.NewRecorder()
        
        dashboardHandler.GetLastDays(recorder, req.Request)

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

    fixtureManager := dashboard.NewDashboardFixtures(helper)
    scenario := dashboard.ScenarioFaultsByWeekday()
    
    err = fixtureManager.LoadScenario(scenario)
    require.NoError(t, err)

    repo, err := createTestRepository(helper.Pool)
    require.NoError(t, err)
    
    userConfig, err := service.LoadDashboardConfig("")
    if err != nil {
        userConfig = service.NewUserConfigService(service.GetDefaultConfig())
    }
    
    dashboardHandler := handlers.NewDashboardHandler(repo, userConfig)

    req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/faults.json", nil)
    recorder := httptest.NewRecorder()
    
    dashboardHandler.GetFaultsChart(recorder, req.Request)

    assert.Equal(t, http.StatusOK, recorder.Code)
    
    var response dto.DashboardResponse
    err = json.Unmarshal(recorder.Body.Bytes(), &response)
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

    fixtureManager := dashboard.NewDashboardFixtures(helper)
    
    // Scenario with logs across 15 days
    scenario := &dashboard.Scenario{
        Name:        "Speculate Actual Test",
        Description: "Logs for speculate vs actual chart",
        Projects: []*dashboard.ProjectFixture{
            {ID: 1, Name: "Test Project", TotalPage: 200, Page: 50},
        },
        Logs: func() []*dashboard.LogFixture {
            var logs []*dashboard.LogFixture
            baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
            
            // Create logs for last 15 days
            for i := 0; i < 15; i++ {
                logs = append(logs, &dashboard.LogFixture{
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
    
    dashboardHandler := handlers.NewDashboardHandler(repo, userConfig)

    req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/speculate_actual.json", nil)
    recorder := httptest.NewRecorder()
    
    dashboardHandler.GetSpeculateActualChart(recorder, req.Request)

    assert.Equal(t, http.StatusOK, recorder.Code)
    
    var response dto.DashboardResponse
    err = json.Unmarshal(recorder.Body.Bytes(), &response)
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

    fixtureManager := dashboard.NewDashboardFixtures(helper)
    scenario := dashboard.ScenarioFaultsByWeekday()
    
    err = fixtureManager.LoadScenario(scenario)
    require.NoError(t, err)

    repo, err := createTestRepository(helper.Pool)
    require.NoError(t, err)
    
    userConfig, err := service.LoadDashboardConfig("")
    if err != nil {
        userConfig = service.NewUserConfigService(service.GetDefaultConfig())
    }
    
    dashboardHandler := handlers.NewDashboardHandler(repo, userConfig)

    req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/faults_week_day.json", nil)
    recorder := httptest.NewRecorder()
    
    dashboardHandler.GetWeekdayFaultsChart(recorder, req.Request)

    assert.Equal(t, http.StatusOK, recorder.Code)
    
    var response dto.DashboardResponse
    err = json.Unmarshal(recorder.Body.Bytes(), &response)
    require.NoError(t, err)
    
    // Verify radar chart structure
    assert.NotNil(t, response.Echart)
    assert.NotEmpty(t, response.Echart.Series)
    
    series := response.Echart.Series[0]
    assert.Equal(t, "radar", series.Type)
    
    // Should have 7 data points (one for each weekday)
    assert.Len(t, series.Data, 7)
    
    // Verify all weekdays present
    expectedWeekdays := []int{0, 1, 2, 3, 4, 5, 6}
    actualWeekdays := make([]int, len(series.Data))
    for i, v := range series.Data {
        actualWeekdays[i] = int(v.(float64)) // Value represents weekday
    }
    
    assert.ElementsMatch(t, expectedWeekdays, actualWeekdays)
}

// TestDashboardMeanProgress_Integration tests /v1/dashboard/echart/mean_progress.json
func TestDashboardMeanProgress_Integration(t *testing.T) {
    if !IsTestDatabase() {
        t.Skip("Test database not configured")
    }

    helper, err := SetupTestDB()
    require.NoError(t, err)
    defer helper.Close()

    fixtureManager := dashboard.NewDashboardFixtures(helper)
    
    // Scenario with varying daily progress
    scenario := &dashboard.Scenario{
        Name:        "Mean Progress Test",
        Description: "Varying daily progress for visual map testing",
        Projects: []*dashboard.ProjectFixture{
            {ID: 1, Name: "Test Project", TotalPage: 200, Page: 50},
        },
        Logs: func() []*dashboard.LogFixture {
            var logs []*dashboard.LogFixture
            baseDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
            
            // Create logs with varying page counts
            dailyPages := []int{5, 10, 15, 20, 25, 30, 35, 40, 45, 50,
                              55, 60, 65, 70, 75, 80, 85, 90, 95, 100,
                              105, 110, 115, 120, 125, 130, 135, 140, 145, 150}
            
            for i := 0; i < 30; i++ {
                logs = append(logs, &dashboard.LogFixture{
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
    
    dashboardHandler := handlers.NewDashboardHandler(repo, userConfig)

    req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/mean_progress.json", nil)
    recorder := httptest.NewRecorder()
    
    dashboardHandler.GetMeanProgressChart(recorder, req.Request)

    assert.Equal(t, http.StatusOK, recorder.Code)
    
    var response dto.DashboardResponse
    err = json.Unmarshal(recorder.Body.Bytes(), &response)
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

    fixtureManager := dashboard.NewDashboardFixtures(helper)
    
    // Scenario with logs spanning 52 weeks
    scenario := &dashboard.Scenario{
        Name:        "Yearly Total Test",
        Description: "Logs spanning 52 weeks for yearly trend",
        Projects: []*dashboard.ProjectFixture{
            {ID: 1, Name: "Test Project", TotalPage: 200, Page: 50},
        },
        Logs: func() []*dashboard.LogFixture {
            var logs []*dashboard.LogFixture
            baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
            
            // Create logs for each week over 52 weeks
            for w := 0; w < 52; w++ {
                weekStart := baseDate.AddDate(0, 0, w*7)
                
                // Create multiple logs per week
                for d := 0; d < 7; d++ {
                    logs = append(logs, &dashboard.LogFixture{
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
    
    dashboardHandler := handlers.NewDashboardHandler(repo, userConfig)

    req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/last_year_total.json", nil)
    recorder := httptest.NewRecorder()
    
    dashboardHandler.GetYearlyTotalChart(recorder, req.Request)

    assert.Equal(t, http.StatusOK, recorder.Code)
    
    var response dto.DashboardResponse
    err = json.Unmarshal(recorder.Body.Bytes(), &response)
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

    // Clear data to test empty state
    err = helper.ClearTestData()
    require.NoError(t, err)

    repo, err := createTestRepository(helper.Pool)
    require.NoError(t, err)
    
    userConfig, err := service.LoadDashboardConfig("")
    if err != nil {
        userConfig = service.NewUserConfigService(service.GetDefaultConfig())
    }
    
    dashboardHandler := handlers.NewDashboardHandler(repo, userConfig)

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
                dashboardHandler.GetDay(recorder, req.Request)
            case tc.endpoint == "/v1/dashboard/projects.json":
                dashboardHandler.GetProjects(recorder, req.Request)
            case tc.endpoint == "/v1/dashboard/last_days.json":
                dashboardHandler.GetLastDays(recorder, req.Request)
            }

            // For empty data, should still return 200 with zero values
            assert.Equal(t, http.StatusOK, recorder.Code)
            
            var response dto.DashboardResponse
            err = json.Unmarshal(recorder.Body.Bytes(), &response)
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
func validateExpectedValues(response dto.DashboardResponse, expected *dashboard.ExpectedResults) error {
    if expected == nil {
        return nil
    }

    if expected.Stats != nil {
        if response.Stats != nil {
            // Allow small floating point differences
            if !floatEqual(response.Stats.PerPages, expected.Stats.PerPages, 0.001) {
                return fmt.Errorf("per_pages mismatch: got %f, expected %f", 
                    response.Stats.PerPages, expected.Stats.PerPages)
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

// createTestRepository creates a dashboard repository for testing
func createTestRepository(pool *pgxpool.Pool) (repository.DashboardRepository, error) {
    return postgres.NewDashboardRepositoryImpl(pool), nil
}
```

---

### 6. Risks and Considerations

#### Known Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| **Fixture data complexity** | High - difficult to maintain | Start simple; use scenario composition; document fixture relationships |
| **Time-dependent calculations** | Medium - tests may fail over time | Use fixed dates (not "today"); document date assumptions |
| **Database cleanup race conditions** | Medium - parallel test interference | Use unique DB names per test; verify cleanup completes |
| **Expected values drift** | Low - calculations change | Version expected values; review changes with PRs |

#### Trade-offs

1. **Fixture Reusability vs. Specificity**: 
   - *Decision*: Create reusable base fixtures + scenario-specific overrides
   - *Rationale*: Reduces duplication while allowing test-specific customization

2. **Database Population vs. Mocking**:
   - *Decision*: Use real database for integration tests; mocks only for unit tests
   - *Rationale*: Ensures full pipeline testing; catches SQL query issues

3. **Date Handling**:
   - *Decision*: Use absolute dates (e.g., "2024-01-15") not relative to "today"
   - *Rationale*: Reproducible tests; fails fast if date logic broken

4. **Expected Values Storage**:
   - *Decision*: Store in YAML alongside fixtures for human readability
   - *Rationale*: Easy to review changes; version-controlled; self-documenting

#### Acceptance Criteria Mapping

| AC Requirement | Test Coverage | Status |
|----------------|---------------|--------|
| AC-DASH-001: Daily statistics | `TestDashboardDayEndpoint_Integration` | ✅ Planned |
| AC-DASH-002: Project aggregates | `TestDashboardProjectsEndpoint_Integration` | ✅ Planned |
| AC-DASH-003: Last days trend | `TestDashboardLastDaysEndpoint_Integration` | ✅ Planned |
| AC-DASH-004: Faults gauge chart | `TestDashboardFaultsChart_Integration` | ✅ Planned |
| AC-DASH-005: Speculate vs actual | `TestDashboardSpeculateActual_Integration` | ✅ Planned |
| AC-DASH-006: Weekday faults radar | `TestDashboardWeekdayFaults_Integration` | ✅ Planned |
| AC-DASH-007: Mean progress chart | `TestDashboardMeanProgress_Integration` | ✅ Planned |
| AC-DASH-008: Yearly trend chart | `TestDashboardYearlyTotal_Integration` | ✅ Planned |

#### Edge Cases Covered

1. **Zero pages**: `ScenarioZeroPages()` - tests division by zero handling
2. **Null dates**: Logs with explicit nil/zero timestamps
3. **Multiple projects**: `ScenarioMultipleProjects()` - 3+ projects
4. **Varying completion**: Unstarted (0%), Running (25%), Finished (100%)
5. **Fault distribution**: All 7 weekdays covered
6. **Date range gaps**: Missing days filled with zeros

---

### Implementation Checklist

**Phase 1: Core Fixture System (Blocker)**
- [ ] Create `test/data/fixtures/dashboard/fixtures.go` - Base fixture manager
- [ ] Create `test/data/fixtures/dashboard/projects.go` - Project generators
- [ ] Create `test/data/fixtures/dashboard/logs.go` - Log generators
- [ ] Add `LoadScenario()` method to populate test database
- [ ] Write unit tests for fixture managers

**Phase 2: Scenario Implementations (Blocker)**
- [ ] Implement `ScenarioZeroPages()`
- [ ] Implement `ScenarioCompleteBook()`
- [ ] Implement `ScenarioMultipleProjects()`
- [ ] Implement `ScenarioFaultsByWeekday()`
- [ ] Create expected values YAML for each scenario

**Phase 3: Integration Tests (Blocker)**
- [ ] Create `test/dashboard_integration_test.go`
- [ ] Implement tests for all 8 endpoints
- [ ] Add error handling tests
- [ ] Verify response validation

**Phase 4: Validation & Documentation (Must-have)**
- [ ] Run all tests and verify pass
- [ ] Update QWEN.md with fixture usage guide
- [ ] Document expected values calculation methodology
- [ ] Create example test demonstrating fixture usage

---

### Quick Start for Implementation

```bash
# 1. Create directory structure
mkdir -p test/data/fixtures/dashboard
mkdir -p test/data/fixtures/testdata

# 2. Implement base fixtures (fixtures.go)
# 3. Implement scenarios (scenarios.go)
# 4. Implement integration tests (dashboard_integration_test.go)
# 5. Run tests
go test -v ./test/... -run "TestDashboard"

# 6. Verify coverage
go test -coverprofile=coverage.out ./test/dashboard_integration_test.go
go tool cover -html=coverage.out
```

---

*Implementation Plan Last Updated: 2026-04-22*
*Plan Author: Architect Agent*
*Task ID: RDL-088*
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass
- [ ] #2 All integration tests pass execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md
- [ ] #8 New code paths include error path tests
- [ ] #9 HTTP handlers test both success and error responses
- [ ] #10 Integration tests verify actual database interactions
<!-- DOD:END -->
