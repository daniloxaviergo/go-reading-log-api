---
id: doc-008
title: 'PRD: Dashboard API Migration - Phase 2 Implementation'
type: other
created_date: '2026-04-21 15:26'
---


# Project Requirements Document

# Executive Summary

This PRD specifies the implementation of Phase 2 dashboard endpoints for the Go Reading Log API, migrating functionality from the Rails application. The dashboard provides daily statistics, project aggregates, trend data, and fault tracking through a set of RESTful JSON endpoints.

**Scope**: 8 dashboard endpoints across 3 priority tiers
**Timeline**: Phase 2 (post-Phase 1 core API completion)
**Impact**: Enables feature parity with Rails dashboard while improving performance and maintainability

---

# Key Requirements

| Requirement | Priority | Status |
|-------------|----------|--------|
| `/v1/dashboard/day.json` - Daily statistics | P1 (Blocker) | Not Started |
| `/v1/dashboard/projects.json` - Project aggregates | P1 (Blocker) | Not Started |
| `/v1/dashboard/last_days.json` - Trend data | P2 (Must-have) | Not Started |
| `/v1/dashboard/echart/faults.json` - Fault gauge | P2 (Must-have) | Not Started |
| `/v1/dashboard/echart/speculate_actual.json` - Progress chart | P2 (Must-have) | Not Started |
| `/v1/dashboard/echart/faults_week_day.json` - Weekday faults | P3 (Should-have) | Not Started |
| `/v1/dashboard/echart/mean_progress.json` - Daily progress | P3 (Should-have) | Not Started |
| `/v1/dashboard/echart/last_year_total.json` - Yearly trend | P3 (Should-have) | Not Started |
| Unit tests for all calculations | P1 (Blocker) | Not Started |
| Integration tests with test database | P2 (Must-have) | Not Started |

---

# Technical Decisions

## Decision 1: Service Layer Architecture

**Decision**: Implement dedicated service layer for dashboard calculations, mirroring Rails service classes.

**Rationale**:
- Clean separation between HTTP handlers and business logic
- Testable without HTTP stack
- Reusable across different API versions or entry points
- Matches established Rails patterns with clear migration path

**Implementation**:
```
internal/service/
├── dashboard/
│   ├── day_service.go           # Daily statistics
│   ├── projects_service.go      # Project aggregates
│   ├── last_days_service.go     # Trend data
│   ├── faults_service.go        # Fault calculations
│   └── speculate_service.go     # Predicted vs actual
└── user_config_service.go       # Configuration access
```

## Decision 2: UserConfig Implementation Strategy

**Decision**: Use configuration file with hardcoded defaults as fallback, NOT database table.

**Rationale**:
- Simplifies Phase 2 (no schema migrations required)
- Easier to manage in containerized environments
- Matches questionnaire response: "Hardcoded Defaults"
- Can migrate to database later if needed

**Implementation**:
```go
// internal/config/dashboard_config.go
type DashboardConfig struct {
    MaxFaults        int     `yaml:"max_faults"`
    PredictionPct    float64 `yaml:"prediction_pct"`
    PagesPerDay      float64 `yaml:"pages_per_day"`
}

func LoadDashboardConfig(path string) (*DashboardConfig, error)
```

## Decision 3: Time Zone Handling

**Decision**: Use server local time for all dashboard calculations, with explicit documentation.

**Rationale**:
- Matches Rails `Time.zone.today` behavior when server TZ is set
- Avoids complexity of per-user time zone storage
- Consistent with current Go implementation (no user-specific TZ)

**Implementation**:
```go
// Always use this for "today" calculations
func GetToday() time.Time {
    return time.Now().Truncate(24 * time.Hour)
}

// For date range queries, use the same reference point
func GetDateRange(days int) (start, end time.Time) {
    end = GetToday()
    start = end.AddDate(0, 0, -days)
    return start, end
}
```

## Decision 4: Fault Calculation Logic

**Decision**: Count ALL faults (not just closed) to match Rails behavior exactly.

**Rationale**:
- Questionnaire response confirmed: "Exact Parity Required"
- Rails query: `Fault.where(project_id: @project.id).count` includes all states
- Consistency with existing API contract

**Implementation**:
```go
// internal/service/dashboard/faults_service.go
func (s *FaultsService) CountForProject(ctx context.Context, projectID int64) (int, error) {
    // Match Rails: count ALL faults for project, regardless of status
    return s.repo.CountFaultsByProject(ctx, projectID)
}
```

## Decision 5: Response Format - Chart Configurations

**Decision**: Return ECharts-style JSON configurations directly from API.

**Rationale**:
- Questionnaire response confirmed: "Chart Configurations"
- Reduces client-side processing
- Consistent with Rails implementation pattern
- Enables rapid frontend development

**Implementation**:
```go
// internal/domain/dto/dashboard_response.go
type DashboardResponse struct {
    Echart *EchartConfig `json:"echart,omitempty"`
    Stats  *StatsData    `json:"stats,omitempty"`
    Logs   []LogEntry    `json:"logs,omitempty"`
}

type EchartConfig struct {
    Tooltip map[string]interface{} `json:"tooltip"`
    Series  []Series               `json:"series"`
    // ... other ECharts options
}
```

## Decision 6: Repository Pattern Extension

**Decision**: Extend repository interface with dashboard-specific methods rather than creating separate service layer.

**Rationale**:
- Reduces architectural complexity for Phase 2
- Maintains existing repository pattern
- Can refactor to dedicated services if scope grows

**Implementation**:
```go
// internal/repository/dashboard_repository.go
type DashboardRepository interface {
    // Existing methods...
    
    // Dashboard-specific aggregations
    GetDailyStats(ctx context.Context, date time.Time) (*DailyStats, error)
    GetProjectAggregates(ctx context.Context) ([]*ProjectAggregate, error)
    GetFaultsByDateRange(ctx context.Context, start, end time.Time) (int, error)
    GetWeekdayFaults(ctx context.Context, start, end time.Time) (map[int]int, error)
}
```

---

# Acceptance Criteria

## Functional Acceptance Criteria

### AC-DASH-001: Daily Dashboard Endpoint
**Given**: A user requests `/v1/dashboard/day.json`
**When**: The request is made with valid authentication (none required)
**Then**:
- Response status is 200 OK
- Response body contains `stats` object with:
  - `previous_week_pages`: integer (sum of pages from previous 7 days)
  - `last_week_pages`: integer (sum of pages from last 7 days)
  - `per_pages`: float (ratio of last_week to previous_week, 3 decimals)
  - `mean_day`: float (average pages per day for current weekday)
  - `spec_mean_day`: float (predicted average for current weekday)
- All numeric fields are non-negative
- Float values rounded to 3 decimal places

**Test Case**:
```go
func TestDailyDashboard_Calculations(t *testing.T) {
    // Setup: Create logs with known page counts
    // Execute: Call /v1/dashboard/day.json
    // Verify: Check all calculated fields match expected values
}
```

---

### AC-DASH-002: Project Dashboard Endpoint
**Given**: A user requests `/v1/dashboard/projects.json`
**When**: The request is made with valid authentication (none required)
**Then**:
- Response status is 200 OK
- Response body contains `projects` array with all projects
- Each project includes eager-loaded logs (first 4, ordered by date DESC)
- Response body contains `stats` object with:
  - `progress_geral`: float (overall completion percentage, 3 decimals)
  - `total_pages`: integer (sum of all project total_page values)
  - `pages`: integer (sum of all project page values)
- Projects ordered by progress descending
- Float values rounded to 3 decimal places

**Test Case**:
```go
func TestProjectDashboard_Aggregates(t *testing.T) {
    // Setup: Create projects with varying completion levels
    // Execute: Call /v1/dashboard/projects.json
    // Verify: Check progress_geral calculation matches manual computation
}
```

---

### AC-DASH-003: Last Days Endpoint
**Given**: A user requests `/v1/dashboard/last_days.json` with type parameter
**When**: The request includes type parameter (1=7days, 2=15days, 3=30days, 4=90days, 5=1day)
**Then**:
- Response status is 200 OK
- Response body contains `logs` array with logs from specified period
- Each log includes eager-loaded project data
- Response body contains `stats` object with:
  - `count_pages`: integer (sum of read_pages in period)
  - `speculate_pages`: integer (config-based target for period)
- Logs ordered by date descending
- Invalid type parameter returns 422 Unprocessable Entity

**Test Case**:
```go
func TestLastDays_InvalidType(t *testing.T) {
    // Setup: Request with invalid type parameter
    // Execute: Call /v1/dashboard/last_days.json?type=invalid
    // Verify: Response status is 422
}
```

---

### AC-DASH-004: Faults Chart Endpoint
**Given**: A user requests `/v1/dashboard/echart/faults.json`
**When**: The request is made with valid authentication (none required)
**Then**:
- Response status is 200 OK
- Response body contains `echart` object with gauge chart configuration
- Gauge shows fault percentage: `(faults_last_30_days / max_faults) * 100`
- Fault percentage rounded to 2 decimal places
- Max faults comes from config (default 10 if not configured)
- Zero faults returns 0% (not NaN or error)

**Test Case**:
```go
func TestFaultsChart_PercentageCalculation(t *testing.T) {
    // Setup: Create known number of faults and configure max_faults
    // Execute: Call /v1/dashboard/echart/faults.json
    // Verify: Percentage calculation is correct
}
```

---

### AC-DASH-005: Speculate vs Actual Chart
**Given**: A user requests `/v1/dashboard/echart/speculate_actual.json`
**When**: The request is made with valid authentication (none required)
**Then**:
- Response status is 200 OK
- Response body contains `echart` object with line chart configuration
- Chart shows daily pages read vs speculative mean
- Speculative mean calculated as: `actual_mean * (1 + prediction_pct)`
- Data covers last 15 days (including today)
- Missing days show 0 pages (not omitted)

**Test Case**:
```go
func TestSpeculateActual_DataCoverage(t *testing.T) {
    // Setup: Create logs with gaps in dates
    // Execute: Call /v1/dashboard/echart/speculate_actual.json
    // Verify: All 15 days present in output, missing days = 0
}
```

---

### AC-DASH-006: Faults by Weekday Chart
**Given**: A user requests `/v1/dashboard/echart/faults_week_day.json`
**When**: The request is made with valid authentication (none required)
**Then**:
- Response status is 200 OK
- Response body contains `echart` object with radar chart configuration
- Chart shows fault counts by weekday (Sunday=0 through Saturday=6)
- Data covers last 6 months
- Each weekday has integer count >= 0
- All 7 weekdays present in output

**Test Case**:
```go
func TestWeekdayFaults_Completeness(t *testing.T) {
    // Setup: Create logs spanning multiple weeks
    // Execute: Call /v1/dashboard/echart/faults_week_day.json
    // Verify: All 7 days present, counts are non-negative integers
}
```

---

### AC-DASH-007: Mean Progress Chart
**Given**: A user requests `/v1/dashboard/echart/mean_progress.json`
**When**: The request is made with valid authentication (none required)
**Then**:
- Response status is 200 OK
- Response body contains `echart` object with line chart configuration
- Shows daily progress as percentage of mean
- Progress = `(daily_pages / mean_pages) * 100 - 100`
- Data covers last 30 days
- Visual map colors based on progress ranges:
  - 0-10%: gray
  - 10-20%: cyan
  - 20-50%: blue
  - >50%: green
  - Negative: red

**Test Case**:
```go
func TestMeanProgress_ColorRanges(t *testing.T) {
    // Setup: Create logs with varying daily page counts
    // Execute: Call /v1/dashboard/echart/mean_progress.json
    // Verify: Color assignments match range definitions
}
```

---

### AC-DASH-008: Last Year Total Chart
**Given**: A user requests `/v1/dashboard/echart/last_year_total.json`
**When**: The request is made with valid authentication (none required)
**Then**:
- Response status is 200 OK
- Response body contains `echart` object with line chart configuration
- Shows total reads per week for last 52 weeks
- Each data point includes:
  - `begin_week`: ISO week start date
  - `end_week`: ISO week end date
  - `count_reads`: integer (total logs in week)
- Average line displayed
- Smooth curve rendering enabled

**Test Case**:
```go
func TestYearlyTotal_WeekBoundaries(t *testing.T) {
    // Setup: Create logs spanning multiple weeks
    // Execute: Call /v1/dashboard/echart/last_year_total.json
    // Verify: Week boundaries align with ISO week numbering
}
```

---

## Non-Functional Acceptance Criteria

### NFA-DASH-001: Performance
| Criterion | Target | Measurement |
|-----------|--------|-------------|
| First request latency | < 500ms | From process start |
| Subsequent requests | < 100ms | 95th percentile |
| Concurrent requests | > 100 QPS | With connection pooling |
| Database query time | < 100ms | For single dashboard endpoint |

### NFA-DASH-002: Reliability
| Criterion | Target | Measurement |
|-----------|--------|-------------|
| Uptime | 99.9% | Monthly average |
| Error rate | < 0.1% | Of all requests |
| Graceful degradation | On | Database unavailable |

### NFA-DASH-003: Maintainability
| Criterion | Target | Measurement |
|-----------|--------|-------------|
| Test coverage | > 85% | Line coverage |
| Code duplication | < 5% | Duplication index |
| Documentation | Complete | All public functions |

---

# Files to Modify

## New Files Created

| File Path | Purpose | Priority |
|-----------|---------|----------|
| `internal/service/dashboard/day_service.go` | Daily statistics calculation | P1 |
| `internal/service/dashboard/projects_service.go` | Project aggregate calculation | P1 |
| `internal/service/dashboard/last_days_service.go` | Trend data calculation | P2 |
| `internal/service/dashboard/faults_service.go` | Fault counting logic | P2 |
| `internal/service/dashboard/speculate_service.go` | Predicted vs actual logic | P2 |
| `internal/service/user_config_service.go` | Configuration access | P1 |
| `internal/repository/dashboard_repository.go` | Dashboard data access interface | P1 |
| `internal/adapter/postgres/dashboard_repository.go` | PostgreSQL implementation | P1 |
| `internal/api/v1/handlers/dashboard_handler.go` | HTTP handlers for all endpoints | P1 |
| `internal/domain/dto/dashboard_response.go` | Response DTOs | P1 |
| `test/dashboard_integration_test.go` | Integration tests | P2 |
| `test/dashboard_service_test.go` | Unit tests | P1 |

## Existing Files Modified

| File Path | Modification | Priority |
|-----------|--------------|----------|
| `internal/api/v1/routes.go` | Add dashboard route registrations | P1 |
| `go.mod` | Add any new dependencies | P1 |
| `.env.example` | Add dashboard configuration options | P2 |

---

# Implementation Checklist

## Phase 1: Foundation (Blocker)

- [ ] **FC-001**: Create `internal/service/user_config_service.go`
  - [ ] Implement config loading from file
  - [ ] Provide defaults for missing values
  - [ ] Add unit tests

- [ ] **FC-002**: Create `internal/repository/dashboard_repository.go`
  - [ ] Define interface with all dashboard query methods
  - [ ] Implement in `internal/adapter/postgres/dashboard_repository.go`
  - [ ] Add connection pooling configuration
  - [ ] Add unit tests

- [ ] **FC-003**: Create `internal/domain/dto/dashboard_response.go`
  - [ ] Define all response DTOs
  - [ ] Implement JSON marshaling
  - [ ] Add validation methods

- [ ] **FC-004**: Create `internal/api/v1/handlers/dashboard_handler.go`
  - [ ] Implement HTTP handlers for all 8 endpoints
  - [ ] Error handling and response formatting
  - [ ] Add unit tests

## Phase 2: Core Calculations (Blocker)

- [ ] **CC-001**: Implement `internal/service/dashboard/day_service.go`
  - [ ] Calculate previous/last week page totals
  - [ ] Calculate mean by weekday
  - [ ] Calculate speculative mean
  - [ ] Add unit tests with known values

- [ ] **CC-002**: Implement `internal/service/dashboard/projects_service.go`
  - [ ] Query all projects with eager-loaded logs
  - [ ] Calculate aggregate progress
  - [ ] Order by progress descending
  - [ ] Add unit tests

- [ ] **CC-003**: Implement `internal/service/dashboard/faults_service.go`
  - [ ] Count faults for date range
  - [ ] Calculate fault percentage
  - [ ] Handle zero faults case
  - [ ] Add unit tests

## Phase 3: Chart Configurations (Must-have)

- [ ] **CC-004**: Implement `internal/service/dashboard/speculate_service.go`
  - [ ] Compare actual vs predicted reading
  - [ ] Generate chart data points
  - [ ] Handle missing days (zero fill)
  - [ ] Add unit tests

- [ ] **CC-005**: Implement weekday fault counting
  - [ ] Group faults by weekday
  - [ ] Generate radar chart data
  - [ ] Handle 6-month date range
  - [ ] Add unit tests

- [ ] **CC-006**: Implement mean progress calculation
  - [ ] Calculate daily progress percentages
  - [ ] Apply visual map color ranges
  - [ ] Generate line chart data
  - [ ] Add unit tests

## Phase 4: Integration & Testing (Must-have)

- [ ] **IT-001**: Set up test database with sample data
  - [ ] Create test fixtures for all scenarios
  - [ ] Include edge cases (zero pages, null dates)
  - [ ] Document fixture structure

- [ ] **IT-002**: Implement integration tests
  - [ ] Test each endpoint against real database
  - [ ] Verify calculations match Rails reference
  - [ ] Test error scenarios
  - [ ] Add test coverage reporting

- [ ] **IT-003**: Performance testing
  - [ ] Benchmark all endpoints
  - [ ] Identify and fix slow queries
  - [ ] Verify connection pooling works correctly

## Phase 5: Documentation (Should-have)

- [ ] **DOC-001**: Update API documentation
  - [ ] Add dashboard endpoint definitions
  - [ ] Document response formats
  - [ ] Add example requests/responses

- [ ] **DOC-002**: Create developer guide
  - [ ] Explain calculation methodologies
  - [ ] Document configuration options
  - [ ] Provide troubleshooting guide

---

# Stakeholder Alignment

| Stakeholder | Responsibility | Verification |
|-------------|----------------|--------------|
| **Product Owner** | Approve feature set and priorities | Review Key Requirements table |
| **Engineering Lead** | Approve technical decisions | Review Technical Decisions section |
| **Backend Developer** | Implement PRD | Execute Implementation Checklist |
| **QA Team** | Test functionality | Verify Acceptance Criteria |
| **DevOps** | Deployment readiness | Verify NFA criteria |

---

# Traceability Matrix

| Requirement ID | User Story | Acceptance Criteria | Test File | Status |
|----------------|------------|---------------------|-----------|--------|
| REQ-DASH-001 | Daily dashboard statistics | AC-DASH-001 | test/dashboard_service_test.go | TODO |
| REQ-DASH-002 | Project aggregates | AC-DASH-002 | test/dashboard_service_test.go | TODO |
| REQ-DASH-003 | Trend data (last days) | AC-DASH-003 | test/dashboard_service_test.go | TODO |
| REQ-DASH-004 | Fault percentage chart | AC-DASH-004 | test/dashboard_service_test.go | TODO |
| REQ-DASH-005 | Speculate vs actual chart | AC-DASH-005 | test/dashboard_service_test.go | TODO |
| REQ-DASH-006 | Weekday fault distribution | AC-DASH-006 | test/dashboard_service_test.go | TODO |
| REQ-DASH-007 | Mean progress tracking | AC-DASH-007 | test/dashboard_service_test.go | TODO |
| REQ-DASH-008 | Yearly trend analysis | AC-DASH-008 | test/dashboard_service_test.go | TODO |

---

# Validation

## Code Quality Standards
- [ ] Go 1.25.7 compatible
- [ ] Follows existing code patterns in project
- [ ] Linting passes (`go vet ./...`)
- [ ] Formatting correct (`go fmt ./...`)

## Technical Feasibility
- [ ] All technologies proven and production-ready
- [ ] No experimental features used
- [ ] Database queries optimized
- [ ] Error handling comprehensive

## User Needs Alignment
- [ ] Endpoints match Rails API behavior exactly
- [ ] Calculation methods verified against Rails
- [ ] Response formats consistent with existing Go API
- [ ] Performance meets acceptable thresholds

---

# Out of Scope

The following items are explicitly **OUT OF SCOPE** for Phase 2:

1. **Authentication/Authorization**: Dashboard endpoints remain public (like Rails)
2. **Caching Layer**: No Redis or in-memory caching implemented
3. **Background Jobs**: All calculations synchronous
4. **User-Specific Config Storage**: Uses config file, not database
5. **Real-time Updates**: No WebSocket or SSE support
6. **Advanced Filtering**: No query parameters for date ranges (fixed periods only)
7. **Export Functionality**: No CSV/PDF export
8. **Dashboard Customization**: User cannot customize widgets

---

# Ready for Implementation

## Approval Status: ✅ **READY FOR IMPLEMENTATION**

This PRD has been:
- ✅ **Validated** against Rails implementation for accuracy
- ✅ **Clarified** through questionnaire to resolve ambiguities
- ✅ **Prioritized** into Blocker/Must-have/Should-have tiers
- ✅ **Tested** with acceptance criteria that are objective and measurable
- ✅ **Documented** with clear file locations and implementation steps

## Prerequisites for Starting Implementation

Before beginning Phase 2, ensure:

1. **Phase 1 is complete**:
   - [ ] Core API endpoints working (`/v1/projects.json`, `/v1/projects/{id}.json`)
   - [ ] Logs endpoint working (`/v1/projects/{project_id}/logs.json`)
   - [ ] Health check endpoint working (`/healthz`)

2. **Development environment ready**:
   - [ ] Go 1.25.7 installed
   - [ ] PostgreSQL running and accessible
   - [ ] Test database created (`reading_log_test`)
   - [ ] `.env` file configured

3. **Stakeholder sign-off**:
   - [ ] Product owner approved feature set
   - [ ] Engineering lead approved technical approach
   - [ ] QA team confirmed testability

## Implementation Start Command

```bash
# Verify environment
make test-setup

# Create initial PRD tracking
mkdir -p backlog/docs
cp docs/doc-template.md backlog/docs/doc-008-dashboard-migration.md

# Begin Phase 2 implementation
cd internal/service/dashboard
touch day_service.go projects_service.go last_days_service.go faults_service.go speculate_service.go
```

---

*PRD Version: 1.0*
*Created: 2026-04-21*
*Last Updated: 2026-04-21*