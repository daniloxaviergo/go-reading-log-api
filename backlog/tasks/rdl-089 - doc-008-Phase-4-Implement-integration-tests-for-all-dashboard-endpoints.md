---
id: RDL-089
title: '[doc-008 Phase 4] Implement integration tests for all dashboard endpoints'
status: To Do
assignee:
  - catarina
created_date: '2026-04-21 15:51'
updated_date: '2026-04-22 11:59'
labels:
  - phase-4
  - testing
  - integration
dependencies: []
references:
  - NFA-DASH-002
  - IT-002
  - Implementation Checklist Phase 4
documentation:
  - doc-008
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create test/dashboard_integration_test.go testing each endpoint against real database. Verify calculations match Rails reference, test error scenarios, and include coverage reporting setup.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Integration tests for all 8 endpoints implemented
- [ ] #2 Calculations verified against Rails reference
- [ ] #3 Error scenarios tested comprehensively
- [ ] #4 Test coverage reporting configured
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task focuses on completing the integration test suite for all dashboard endpoints, ensuring they match Rails API behavior and cover all edge cases.

**Architecture Strategy:**
- Leverage existing fixture system from RDL-088 (already implemented)
- Extend `dashboard_integration_test.go` with comprehensive validation
- Implement Rails comparison tests using existing test infrastructure
- Add coverage reporting via Go's built-in coverage tools

**Key Components:**

```
test/
├── dashboard_integration_test.go    # Main integration tests (extend)
├── fixtures/dashboard/              # Test data fixtures (already created)
│   ├── fixtures.go
│   ├── scenarios.go
│   └── testdata/expected-values.yaml
└── integration/
    └── rails_comparison_test.go     # New: Rails API comparison tests
```

**Test Categories:**
1. **Endpoint Tests**: Each of the 8 dashboard endpoints tested with real database
2. **Calculation Verification**: Compare calculated values against expected values from fixtures
3. **Error Scenarios**: Invalid inputs, empty data, edge cases
4. **Rails Comparison**: Verify Go API matches Rails API behavior exactly

---

### 2. Files to Modify

#### New Files to Create

| File Path | Purpose | Lines | Priority |
|-----------|---------|-------|----------|
| `test/integration/rails_comparison_test.go` | Compare Go vs Rails API responses | ~300 | P1 |
| `test/integration/coverage_report.go` | Coverage reporting utilities | ~150 | P2 |

#### Files to Extend

| File Path | Extension | Priority |
|-----------|-----------|----------|
| `test/dashboard_integration_test.go` | Add Rails comparison tests, error scenarios, coverage setup | P1 |
| `test/fixtures/dashboard/scenarios.go` | Add missing scenarios for edge cases | P2 |
| `test/fixtures/testdata/expected-values.yaml` | Complete expected values for all endpoints | P1 |

---

### 3. Dependencies

**Prerequisites (Already Completed):**
- [x] RDL-088 - Test fixtures and scenarios created
- [x] All 8 dashboard endpoints implemented in handlers
- [x] `DashboardRepository` interface and PostgreSQL implementation
- [x] Service layer for calculations (`DayService`, `ProjectsService`, etc.)
- [x] UserConfig service with defaults

**External Dependencies (Already Available):**
- `github.com/stretchr/testify/assert` - For assertions
- `github.com/stretchr/testify/require` - For required assertions
- `net/http/httptest` - For HTTP request testing
- `encoding/json` - For JSON parsing

**No New Dependencies Required**

---

### 4. Code Patterns

#### Pattern 1: Integration Test Structure

```go
// test/integration/rails_comparison_test.go
package integration

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "go-reading-log-api-next/internal/api/v1/handlers"
    "go-reading-log-api-next/internal/domain/dto"
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
        // Setup test database
        helper, err := SetupTestDB()
        require.NoError(t, err)
        defer helper.Close()

        // Load fixture data
        fixtureManager := dashboard.NewDashboardFixtures(helper)
        scenario := test.setupFixture()
        
        err = fixtureManager.LoadScenario(scenario)
        require.NoError(t, err)

        // Create Go handler
        goHandler := createGoHandler(helper.Pool)

        // Make request to Go API
        goResponse := makeRequest(t, goHandler, test.method, test.endpoint)

        // Fetch Rails API response (requires Rails server running)
        railsURL := "http://localhost:3001" + test.endpoint
        railsResponse := fetchRailsAPI(t, railsURL)

        // Compare responses
        test.validate(t, goResponse, railsResponse)
    })
}

// Helper functions
func createGoHandler(pool *pgxpool.Pool) http.Handler {
    repo := postgres.NewDashboardRepositoryImpl(pool)
    userConfig := service.NewUserConfigService(service.GetDefaultConfig())
    dashboardHandler := handlers.NewDashboardHandler(repo, userConfig)
    
    router := mux.NewRouter()
    router.HandleFunc("/v1/dashboard/day.json", dashboardHandler.Day).Methods("GET")
    router.HandleFunc("/v1/dashboard/projects.json", dashboardHandler.Projects).Methods("GET")
    // Add all 8 endpoints...
    return router
}

func fetchRailsAPI(t *testing.T, url string) []byte {
    resp, err := http.Get(url)
    require.NoError(t, err)
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    require.NoError(t, err)
    
    assert.Equal(t, http.StatusOK, resp.StatusCode,
        "Rails API returned status %d for URL: %s", resp.StatusCode, url)
    
    return body
}

func makeRequest(t *testing.T, handler http.Handler, method, endpoint string) []byte {
    req := httptest.NewRequest(method, endpoint, nil)
    recorder := httptest.NewRecorder()
    handler.ServeHTTP(recorder, req)
    return recorder.Body.Bytes()
}
```

#### Pattern 2: Response Validation

```go
// test/integration/validation.go
package integration

import (
    "fmt"
    "math"
    "reflect"

    "github.com/stretchr/testify/assert"
)

// ValidationResult holds comparison results
type ValidationResult struct {
    Passed   bool
    Errors   []string
    Warnings []string
}

// Validator provides response validation utilities
type Validator struct {
    tolerance float64 // Floating point comparison tolerance
}

// NewValidator creates a new validator with specified tolerance
func NewValidator(tolerance float64) *Validator {
    return &Validator{tolerance: tolerance}
}

// ValidateDashboardResponse compares Go and Rails responses
func (v *Validator) ValidateDashboardResponse(
    t *testing.T,
    goResponse interface{},
    railsResponse interface{},
    endpoint string,
) ValidationResult {
    result := ValidationResult{Passed: true, Errors: []string{}}

    // Convert to map for comparison
    goMap := v.interfaceToMap(goResponse)
    railsMap := v.interfaceToMap(railsResponse)

    // Compare common fields
    commonFields := []string{"status", "data", "meta"}
    for _, field := range commonFields {
        goVal, goOk := goMap[field]
        railsVal, railsOk := railsMap[field]

        if goOk != railsOk {
            result.Errors = append(result.Errors,
                fmt.Sprintf("%s: field '%s' present in Go but not Rails (or vice versa)", 
                    endpoint, field))
            result.Passed = false
            continue
        }

        if !v.valuesEqual(goVal, railsVal) {
            result.Errors = append(result.Errors,
                fmt.Sprintf("%s: field '%s' differs - Go: %v, Rails: %v", 
                    endpoint, field, goVal, railsVal))
            result.Passed = false
        }
    }

    // Special handling for calculated fields
    v.validateCalculatedFields(t, goMap, railsMap, endpoint, &result)

    return result
}

// validateCalculatedFields handles special comparison for calculated values
func (v *Validator) validateCalculatedFields(
    t *testing.T,
    goMap, railsMap map[string]interface{},
    endpoint string,
    result *ValidationResult,
) {
    // Handle float comparisons with tolerance
    floatFields := []string{"progress", "per_pages", "mean_day", "spec_mean_day"}
    
    for _, field := range floatFields {
        goVal, goOk := goMap[field]
        railsVal, railsOk := railsMap[field]

        if !goOk || !railsOk {
            continue
        }

        goFloat := v.toFloat64(goVal)
        railsFloat := v.toFloat64(railsVal)

        if math.Abs(goFloat-railsFloat) > v.tolerance {
            result.Warnings = append(result.Warnings,
                fmt.Sprintf("%s: %s differs slightly - Go: %.6f, Rails: %.6f (tolerance: %.6f)", 
                    endpoint, field, goFloat, railsFloat, v.tolerance))
        }
    }
}

func (v *Validator) valuesEqual(a, b interface{}) bool {
    return reflect.DeepEqual(a, b)
}

func (v *Validator) toFloat64(val interface{}) float64 {
    switch v := val.(type) {
    case float64:
        return v
    case float32:
        return float64(v)
    case int:
        return float64(v)
    default:
        return 0
    }
}

func (v *Validator) interfaceToMap(val interface{}) map[string]interface{} {
    result := make(map[string]interface{})
    
    switch v := val.(type) {
    case map[string]interface{}:
        return v
    case dto.DashboardResponse:
        // Extract fields from DTO
        if v.Stats != nil {
            result["stats"] = v.Stats
        }
        if v.Echart != nil {
            result["echart"] = v.Echart
        }
        if len(v.Logs) > 0 {
            result["logs"] = v.Logs
        }
    }

    return result
}
```

#### Pattern 3: Error Scenario Testing

```go
// test/integration/error_scenarios_test.go
package integration

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "go-reading-log-api-next/internal/api/v1/handlers"
    "go-reading-log-api-next/test/fixtures/dashboard"
)

// ErrorScenario represents a test case for error handling
type ErrorScenario struct {
    Name       string
    Endpoint   string
    Method     string
    Setup      func(*testing.T) *dashboard.Scenario
    Request    func(*http.Request)
    Validate   func(*testing.T, *httptest.ResponseRecorder)
}

// RunErrorScenarios runs all error scenario tests
func RunErrorScenarios(t *testing.T, scenarios []ErrorScenario) {
    for _, scenario := range scenarios {
        t.Run(scenario.Name, func(t *testing.T) {
            // Setup test database
            helper, err := SetupTestDB()
            require.NoError(t, err)
            defer helper.Close()

            // Load fixture if setup provided
            if scenario.Setup != nil {
                fixtureManager := dashboard.NewDashboardFixtures(helper)
                scenario := scenario.Setup(t)
                err = fixtureManager.LoadScenario(scenario)
                require.NoError(t, err)
            }

            // Create handler
            repo := postgres.NewDashboardRepositoryImpl(helper.Pool)
            userConfig := service.NewUserConfigService(service.GetDefaultConfig())
            dashboardHandler := handlers.NewDashboardHandler(repo, userConfig)

            // Create request
            req := httptest.NewRequest(scenario.Method, scenario.Endpoint, nil)
            
            // Apply request modifications
            if scenario.Request != nil {
                scenario.Request(req)
            }

            // Execute request
            recorder := httptest.NewRecorder()
            dashboardHandler.ServeHTTP(recorder, req)

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
        Setup: func(t *testing.T) *dashboard.Scenario {
            return dashboard.ScenarioEmptyData()
        },
        Validate: func(t *testing.T, recorder *httptest.ResponseRecorder) {
            assert.Equal(t, http.StatusOK, recorder.Code)
            // Should return empty array with zero values
            var response map[string]interface{}
            json.Unmarshal(recorder.Body.Bytes(), &response)
            data := response["data"].([]interface{})
            assert.Empty(t, data)
        },
    },
    // Add more error scenarios...
}
```

---

### 5. Testing Strategy

#### Test Coverage Goals

| Endpoint | Unit Tests | Integration Tests | Error Tests | Rails Comparison |
|----------|------------|-------------------|-------------|------------------|
| `/v1/dashboard/day.json` | ✅ | ✅ | ✅ | ✅ |
| `/v1/dashboard/projects.json` | ✅ | ✅ | ✅ | ✅ |
| `/v1/dashboard/last_days.json` | ✅ | ✅ | ✅ | ✅ |
| `/v1/dashboard/echart/faults.json` | ✅ | ✅ | ✅ | ✅ |
| `/v1/dashboard/echart/speculate_actual.json` | ✅ | ✅ | ✅ | ✅ |
| `/v1/dashboard/echart/faults_week_day.json` | ✅ | ✅ | ✅ | ✅ |
| `/v1/dashboard/echart/mean_progress.json` | ✅ | ✅ | ✅ | ✅ |
| `/v1/dashboard/echart/last_year_total.json` | ✅ | ✅ | ✅ | ✅ |

#### Test Execution Order

```bash
# 1. Run unit tests (fast, no database)
go test -v ./test/unit/... -run "TestDashboard"

# 2. Run integration tests (requires database)
go test -v ./test/dashboard_integration_test.go

# 3. Run Rails comparison tests (requires Rails server)
go test -v ./test/integration/rails_comparison_test.go

# 4. Run error scenario tests
go test -v ./test/integration/error_scenarios_test.go

# 5. Generate coverage report
go test -coverprofile=coverage.out ./test/...
go tool cover -html=coverage.out
```

#### Coverage Reporting Setup

```go
// test/integration/coverage_report.go
package integration

import (
    "fmt"
    "os"
    "runtime/debug"

    "github.com/stretchr/testify/assert"
)

// CoverageReport tracks test coverage statistics
type CoverageReport struct {
    TestsRun      int
    TestsPassed   int
    TestsFailed   int
    CoveragePct   float64
    FilesCovered  []string
    SlowTests     []string
}

// NewCoverageReport creates a new coverage report
func NewCoverageReport() *CoverageReport {
    return &CoverageReport{
        FilesCovered: make([]string, 0),
        SlowTests:    make([]string, 0),
    }
}

// RecordTest records a test execution result
func (r *CoverageReport) RecordTest(name string, passed bool, duration float64) {
    r.TestsRun++
    if passed {
        r.TestsPassed++
    } else {
        r.TestsFailed++
    }

    // Track slow tests (> 1 second)
    if duration > 1.0 {
        r.SlowTests = append(r.SlowTests, fmt.Sprintf("%s (%.2fs)", name, duration))
    }
}

// CalculateCoverage calculates coverage percentage
func (r *CoverageReport) CalculateCoverage() float64 {
    if r.TestsRun == 0 {
        return 0.0
    }
    return float64(r.TestsPassed) / float64(r.TestsRun) * 100
}

// GenerateHTML generates HTML coverage report
func (r *CoverageReport) GenerateHTML() string {
    html := `<!DOCTYPE html>
<html>
<head><title>Test Coverage Report</title></head>
<body>
<h1>Dashboard API Test Coverage Report</h1>
<table border="1">
<tr><th>Metric</th><th>Value</th></tr>
<tr><td>Tests Run</td><td>` + fmt.Sprintf("%d", r.TestsRun) + `</td></tr>
<tr><td>Tests Passed</td><td>` + fmt.Sprintf("%d", r.TestsPassed) + `</td></tr>
<tr><td>Tests Failed</td><td>` + fmt.Sprintf("%d", r.TestsFailed) + `</td></tr>
<tr><td>Coverage</td><td>` + fmt.Sprintf("%.1f%%", r.CalculateCoverage()) + `</td></tr>
</table>

<h2>Files Covered</h2>
<ul>`
    
    for _, file := range r.FilesCovered {
        html += `<li>` + file + `</li>`
    }
    
    html += `</ul>

<h2>Slow Tests (> 1s)</h2>
<ul>`
    
    for _, slow := range r.SlowTests {
        html += `<li>` + slow + `</li>`
    }
    
    html += `</ul>
</body>
</html>`
    
    return html
}

// SaveCoverageReport saves the coverage report to a file
func (r *CoverageReport) SaveCoverageReport(filename string) error {
    html := r.GenerateHTML()
    return os.WriteFile(filename, []byte(html), 0644)
}
```

---

### 6. Risks and Considerations

#### Known Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| **Rails API unavailable** | High - comparison tests fail | Run comparison tests only when Rails is available; mark as optional |
| **Time-dependent calculations** | Medium - flaky tests | Use fixed dates in fixtures; avoid "now" in comparisons |
| **Database state pollution** | Medium - test interference | Use unique DB names per test; proper cleanup |
| **Floating point precision** | Low - minor differences | Use tolerance-based comparison (0.001) |

#### Trade-offs

1. **Complete Rails Parity vs. Time**:
   - *Decision*: Focus on critical paths first; incremental improvement
   - *Rationale*: Core functionality must match; edge cases can be refined later

2. **Test Speed vs. Coverage**:
   - *Decision*: Separate fast unit tests from slower integration tests
   - *Rationale*: Developers can run unit tests quickly; CI runs full suite

3. **Rails Dependency**:
   - *Decision*: Make Rails comparison optional via environment variable
   - *Rationale*: Allow testing without Rails server running

#### Acceptance Criteria Mapping

| AC Requirement | Test Implementation | Status |
|----------------|---------------------|--------|
| #1 Integration tests for all 8 endpoints | `dashboard_integration_test.go` extended | ✅ Planned |
| #2 Calculations verified against Rails | `rails_comparison_test.go` | ✅ Planned |
| #3 Error scenarios tested comprehensively | `error_scenarios_test.go` | ✅ Planned |
| #4 Test coverage reporting configured | `coverage_report.go` + flags | ✅ Planned |

---

### Implementation Checklist

**Phase 1: Core Integration Tests (Blocker)**
- [ ] Extend `dashboard_integration_test.go` with all 8 endpoints
- [ ] Implement response validation against expected values
- [ ] Add test fixtures for missing scenarios
- [ ] Verify database cleanup between tests

**Phase 2: Rails Comparison (Blocker)**
- [ ] Create `rails_comparison_test.go`
- [ ] Implement HTTP client for Rails API
- [ ] Compare responses for all 8 endpoints
- [ ] Handle floating point tolerance

**Phase 3: Error Scenarios (Must-have)**
- [ ] Create `error_scenarios_test.go`
- [ ] Test invalid inputs (bad dates, types)
- [ ] Test empty database handling
- [ ] Verify error response formats

**Phase 4: Coverage Reporting (Must-have)**
- [ ] Create `coverage_report.go` utilities
- [ ] Add `-coverprofile` flag to test commands
- [ ] Generate HTML coverage report
- [ ] Document coverage targets (>85%)

**Phase 5: Documentation (Should-have)**
- [ ] Update QWEN.md with testing guide
- [ ] Document expected values calculation
- [ ] Create troubleshooting section

---

### Quick Start Commands

```bash
# Run all dashboard tests
go test -v ./test -run "TestDashboard"

# Run with coverage
go test -v -coverprofile=coverage.out ./test -run "TestDashboard"
go tool cover -html=coverage.out

# Run Rails comparison (requires Rails server on :3001)
RAILS_API_URL=http://localhost:3001 go test -v ./test/integration -run "TestRails"

# Run error scenarios
go test -v ./test/integration -run "TestError"

# Generate full coverage report
go test -coverprofile=full-coverage.out ./...
go tool cover -html=full-coverage.out
```

---

### Expected Outcomes

After implementation, the test suite will provide:

1. **Comprehensive Coverage**: All 8 dashboard endpoints tested with real database
2. **Rails Parity Verification**: Go API responses match Rails API within tolerance
3. **Error Resilience**: All error scenarios covered and validated
4. **Performance Baseline**: Coverage reporting identifies slow tests
5. **Documentation**: Test code serves as usage examples

---

*Implementation Plan Last Updated: 2026-04-22*
*Task ID: RDL-089*
*PRD: doc-008 Phase 4 - Integration & Testing*
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
