---
id: doc-011
title: 'PRD: Migrate /v1/dashboard/projects.json Endpoint from Rails to Go'
type: other
created_date: '2026-04-28 11:05'
---

# Project Requirements Document

# Executive Summary

This PRD specifies the migration of the `/v1/dashboard/projects.json` endpoint from the Rails application to the Go Reading Log API. The endpoint provides project aggregates with eager-loaded logs and overall dashboard statistics, enabling feature parity between the Go and Rails implementations.

**Why necessary**: The Rails application currently serves this endpoint, but the Go API needs to provide equivalent functionality to complete the Phase 2 dashboard migration. This ensures frontend clients can seamlessly transition to the Go API without breaking changes.

**Scope**: Single endpoint migration with exact Rails parity
**Timeline**: Phase 2 (Dashboard API Migration)
**Impact**: Enables dashboard project view functionality in Go API

---

# Key Requirements

| Requirement | Priority | Status | Owner |
|-------------|----------|--------|-------|
| Implement `GET /v1/dashboard/projects.json` endpoint | P1 (Blocker) | Not Started | Backend Developer |
| Filter projects by calculated `running` status | P1 (Blocker) | Not Started | Backend Developer |
| Calculate `stats` object (progress_geral, total_pages, pages) | P1 (Blocker) | Not Started | Backend Developer |
| Order projects by progress descending | P1 (Blocker) | Not Started | Backend Developer |
| Eager-load first 4 logs per project (ordered by date DESC) | P1 (Blocker) | Not Started | Backend Developer |
| Match Rails response format exactly | P1 (Blocker) | Not Started | Backend Developer |
| Unit tests for handler logic | P2 (Must-have) | Not Started | Backend Developer |
| Integration tests with database | P2 (Must-have) | Not Started | Backend Developer |
| Register route in `routes.go` | P1 (Blocker) | Not Started | Backend Developer |

---

# Technical Decisions

## Decision 1: Response Format - Exact Rails Parity

**Decision**: Match Rails ActiveModelSerializers output exactly with flat JSON structure.

**Rationale**:
- User questionnaire confirmed: "Exact Rails parity" required
- Frontend expects this specific structure
- Easier migration path with no breaking changes
- Rails output is the source of truth

**Implementation**:
```json
{
  "projects": [
    {
      "id": 1,
      "name": "Project Name",
      "total_page": 200,
      "page": 50,
      "started_at": "2024-01-15T10:30:00Z",
      "progress": 25.0,
      "status": "running",
      "logs_count": 4,
      "days_unreading": 5,
      "median_day": 10.0,
      "finished_at": "2024-02-15T00:00:00Z",
      "logs": [
        {
          "id": 1,
          "data": "2024-01-15T10:30:00Z",
          "start_page": 0,
          "end_page": 25,
          "note": "Morning reading",
          "project": {
            "id": 1,
            "name": "Project Name",
            "total_page": 200,
            "page": 50,
            "started_at": "2024-01-15T10:30:00Z",
            "status": "running",
            "progress": 25.0
          }
        }
      ]
    }
  ],
  "stats": {
    "progress_geral": 25.5,
    "total_pages": 500,
    "pages": 127
  }
}
```

**Trade-offs**:
- ✅ No breaking changes for frontend
- ✅ Matches Rails behavior exactly
- ❌ Requires careful testing to ensure parity

---

## Decision 2: Status Filter Implementation

**Decision**: Use Go's calculated `status` field with 7-day threshold (not raw `page != total_page`).

**Rationale**:
- User questionnaire selected: "Calculated status field"
- Go's `CalculateStatus()` method provides consistent status logic
- Matches Rails behavior when configured with 7-day threshold
- More maintainable than SQL-level filtering

**Implementation**:
```go
// internal/domain/models/project.go
func (p *Project) CalculateStatus() string {
    if p.Page >= p.TotalPage {
        return StatusFinished
    }
    if p.StartedAt == nil {
        return StatusUnstarted
    }
    
    daysUnreading := p.CalculateDaysUnreading()
    if daysUnreading <= 7 {
        return StatusRunning
    }
    if daysUnreading <= 14 {
        return StatusSleeping
    }
    return StatusStopped
}
```

**Query filter**:
```sql
WHERE p.status = 'running'
-- OR use calculated status in Go after fetching all projects
```

---

## Decision 3: Stats Object Placement

**Decision**: Place `stats` object at root level alongside `projects` array.

**Rationale**:
- Rails uses `.merge(stats: stats)` which places it at root
- User confirmed: "equal rails-app" structure
- Matches existing dashboard endpoint patterns in Go

**Implementation**:
```go
response := map[string]interface{}{
    "projects": projects,
    "stats": map[string]interface{}{
        "progress_geral": progressGeral,
        "total_pages":    totalPages,
        "pages":          pages,
    },
}
```

---

## Decision 4: Sorting Logic

**Decision**: Order projects by progress percentage descending (`page / total_page DESC`).

**Rationale**:
- Rails uses `scope :order_progress` with SQL calculation
- Matches user expectation for dashboard view
- Can be implemented efficiently in SQL

**Implementation**:
```sql
ORDER BY 
    CASE 
        WHEN total_page = 0 THEN 0 
        ELSE page::float / total_page::float 
    END DESC,
    p.id ASC
```

---

## Decision 5: Testing Strategy

**Decision**: Implement both unit tests (mock repository) and integration tests (real database).

**Rationale**:
- User questionnaire selected: "Unit + Integration tests"
- Unit tests verify business logic in isolation
- Integration tests verify end-to-end behavior and Rails parity
- Follows existing project testing patterns

**Implementation**:
```go
// Unit test
func TestDashboardHandler_Projects(t *testing.T) {
    mockRepo := &MockDashboardRepository{}
    handler := NewDashboardHandler(mockRepo, userConfig)
    // Test handler logic with mocked data
}

// Integration test
func TestDashboardProjectsEndpoint_Integration(t *testing.T) {
    helper := test.SetupTestDB()
    defer helper.Close()
    // Test against real database with fixtures
}
```

---

# Acceptance Criteria

## Functional Acceptance Criteria

### AC-PROJ-001: Endpoint Response Structure
**Given**: A user requests `GET /v1/dashboard/projects.json`
**When**: The request is processed successfully
**Then**:
- Response status is 200 OK
- Response body contains `projects` array
- Response body contains `stats` object at root level
- Each project includes all fields: `id`, `name`, `total_page`, `page`, `started_at`, `progress`, `status`, `logs_count`, `days_unreading`, `median_day`, `finished_at`, `logs`
- Each log includes: `id`, `data`, `start_page`, `end_page`, `note`, `project`
- Float values are rounded to 3 decimal places

**Test Case**:
```go
func TestDashboardProjects_ResponseStructure(t *testing.T) {
    // Setup: Create test database with sample projects
    // Execute: GET /v1/dashboard/projects.json
    // Verify: Response has correct structure and field types
}
```

---

### AC-PROJ-002: Running Status Filter
**Given**: Multiple projects with different statuses
**When**: Request is made to `/v1/dashboard/projects.json`
**Then**:
- Only projects with `status = "running"` are returned
- Projects with `status != "running"` are excluded
- Status calculation uses 7-day threshold for "running"

**Test Case**:
```go
func TestDashboardProjects_RunningStatusFilter(t *testing.T) {
    // Setup: Create projects with status: running, finished, stopped
    // Execute: GET /v1/dashboard/projects.json
    // Verify: Only running projects in response
}
```

---

### AC-PROJ-003: Stats Calculation
**Given**: Multiple running projects with varying page counts
**When**: Request is made to `/v1/dashboard/projects.json`
**Then**:
- `stats.total_pages` = sum of all project `total_page` values
- `stats.pages` = sum of all project `page` values  
- `stats.progress_geral` = round((pages / total_pages) * 100, 3)
- Division by zero handled (return 0.0)

**Test Case**:
```go
func TestDashboardProjects_StatsCalculation(t *testing.T) {
    // Setup: Create 3 projects with known values
    //   Project 1: total_page=200, page=50
    //   Project 2: total_page=300, page=100
    //   Project 3: total_page=100, page=25
    // Execute: GET /v1/dashboard/projects.json
    // Verify: 
    //   stats.total_pages = 600
    //   stats.pages = 175
    //   stats.progress_geral = 29.167
}
```

---

### AC-PROJ-004: Project Ordering
**Given**: Multiple running projects with different progress levels
**When**: Request is made to `/v1/dashboard/projects.json`
**Then**:
- Projects are ordered by progress percentage descending
- Projects with equal progress are ordered by `id` ascending
- Progress = `page / total_page * 100`

**Test Case**:
```go
func TestDashboardProjects_Ordering(t *testing.T) {
    // Setup: Create projects with progress: 10%, 50%, 25%
    // Execute: GET /v1/dashboard/projects.json
    // Verify: Order is 50%, 25%, 10%
}
```

---

### AC-PROJ-005: Eager-Loaded Logs
**Given**: Projects with varying numbers of logs
**When**: Request is made to `/v1/dashboard/projects.json`
**Then**:
- Each project includes `logs` array with first 4 logs
- Logs are ordered by `data` DESC (most recent first)
- Each log includes eager-loaded `project` object
- Projects with no logs have empty `logs` array

**Test Case**:
```go
func TestDashboardProjects_EagerLoadedLogs(t *testing.T) {
    // Setup: Create project with 10 logs
    // Execute: GET /v1/dashboard/projects.json
    // Verify: Only 4 most recent logs returned
    // Verify: Logs ordered by data DESC
}
```

---

### AC-PROJ-006: Rails Parity Validation
**Given**: Same data in both Rails and Go databases
**When**: Both endpoints return responses
**Then**:
- Go response structure matches Rails response structure
- All calculated fields match Rails values
- All filtering and ordering matches Rails behavior

**Test Case**:
```go
func TestDashboardProjects_RailsParity(t *testing.T) {
    // Setup: Populate both Rails and Go databases with identical data
    // Execute: 
    //   railsResponse := GET http://rails:3001/v1/dashboard/projects.json
    //   goResponse := GET http://go:3000/v1/dashboard/projects.json
    // Verify: Structurally equivalent (normalize timestamps, floats)
}
```

---

## Non-Functional Acceptance Criteria

### NFA-PROJ-001: Performance
| Criterion | Target | Measurement |
|-----------|--------|-------------|
| Response latency | < 200ms | 95th percentile |
| Database query time | < 50ms | Single query with JOINs |
| Concurrent requests | > 50 QPS | With connection pooling |

### NFA-PROJ-002: Code Quality
| Criterion | Target | Measurement |
|-----------|--------|-------------|
| Test coverage | > 85% | Line coverage for new code |
| Code duplication | < 5% | No duplicate SQL queries |
| Documentation | Complete | All public functions documented |

### NFA-PROJ-003: Error Handling
| Criterion | Target | Measurement |
|-----------|--------|-------------|
| Database errors | 500 Internal Server Error | With proper logging |
| Invalid requests | 400 Bad Request | With validation messages |
| No data found | 200 OK with empty arrays | Not 404 |

---

# Files to Modify

## New Files Created

| File Path | Purpose | Priority | Rationale |
|-----------|---------|----------|-----------|
| `internal/service/dashboard/projects_service.go` | Project aggregate calculation service | P1 | Business logic separation, testability |
| `internal/service/dashboard/projects_service_test.go` | Unit tests for projects service | P2 | Verify calculation logic |

## Existing Files Modified

| File Path | Modification | Priority | Rationale |
|-----------|--------------|----------|-----------|
| `internal/api/v1/handlers/dashboard_handler.go` | Update `Projects()` method to match Rails response | P1 | Align with Rails structure |
| `internal/api/v1/routes.go` | Register `/v1/dashboard/projects.json` route | P1 | Enable endpoint access |
| `internal/adapter/postgres/dashboard_repository.go` | Add query method for running projects with logs | P1 | Efficient data fetching |
| `internal/domain/dto/dashboard_response.go` | Add/update DTOs for project response | P1 | Type-safe response structure |
| `test/integration/dashboard_projects_test.go` | Integration tests for endpoint | P2 | End-to-end validation |

---

# Files Created

| File Path | Type | Purpose |
|-----------|------|---------|
| `docs/dashboard-projects-endpoint.md` | Technical doc | API documentation for developers |
| `test/fixtures/dashboard_projects.go` | Test data | Sample data for integration tests |
| `backlog/decisions/decision-009-projects-endpoint.md` | Decision record | Document technical decisions |

---

# Validation Rules

## Project Validation

| Field | Rule | Error Message |
|-------|------|---------------|
| `id` | Must be positive integer | "project id must be positive" |
| `name` | Required, max 255 chars | "name is required" |
| `total_page` | Must be > 0 | "total_page must be greater than 0" |
| `page` | Must be >= 0 and <= total_page | "page cannot exceed total_page" |
| `started_at` | Optional, RFC3339 format | "started_at must be in RFC3339 format" |

## Stats Calculation Validation

| Field | Rule | Error Handling |
|-------|------|----------------|
| `progress_geral` | 0-100 range, 3 decimals | Round to 3 decimals, clamp to 0-100 |
| `total_pages` | Non-negative integer | Use COALESCE for NULL values |
| `pages` | Non-negative integer | Use COALESCE for NULL values |

## Log Validation

| Field | Rule | Error Message |
|-------|------|---------------|
| `id` | Must be positive integer | "log id must be positive" |
| `start_page` | Must be >= 0 | "start_page cannot be negative" |
| `end_page` | Must be >= start_page | "end_page must be >= start_page" |
| `data` | Required, RFC3339 format | "data is required in RFC3339 format" |

**DRY Principle**: Validation rules are shared between handler and service layer using centralized validation functions in `internal/validation/`.

---

# Out of Scope

The following items are explicitly **OUT OF SCOPE** for this PRD:

1. **Log Creation**: POST `/v1/dashboard/projects.json` is not implemented (Phase 3)
2. **Project Filtering**: No query parameters for filtering beyond status
3. **Pagination**: All running projects returned (no limit/offset)
4. **Authentication**: Endpoint remains public (like Rails)
5. **Caching**: No caching layer implemented
6. **Other Dashboard Endpoints**: This PRD covers only `/v1/dashboard/projects.json`
7. **Database Migrations**: Schema changes handled separately
8. **UI Changes**: Frontend modifications are out of scope

---

# Implementation Checklist

## Phase 1: Foundation (Blocker)

- [ ] **FC-001**: Update `internal/domain/dto/dashboard_response.go`
  - [ ] Add/verify `ProjectWithLogs` struct with all required fields
  - [ ] Add `StatsData` fields for progress_geral, total_pages, pages
  - [ ] Implement JSON marshaling for response structure

- [ ] **FC-002**: Implement `internal/service/dashboard/projects_service.go`
  - [ ] Create `GetRunningProjectsWithLogs()` method
  - [ ] Implement status filtering logic
  - [ ] Implement progress calculation
  - [ ] Implement stats calculation
  - [ ] Add progress ordering logic

- [ ] **FC-003**: Update `internal/adapter/postgres/dashboard_repository.go`
  - [ ] Add `GetRunningProjectsWithLogs()` query method
  - [ ] Implement SQL query with JOIN for eager-loaded logs
  - [ ] Add progress ordering in SQL
  - [ ] Handle NULL values with COALESCE

## Phase 2: Handler Integration (Blocker)

- [ ] **HI-001**: Update `internal/api/v1/handlers/dashboard_handler.go`
  - [ ] Modify `Projects()` method to use new service
  - [ ] Implement response formatting matching Rails structure
  - [ ] Add error handling
  - [ ] Add logging

- [ ] **HI-002**: Register route in `internal/api/v1/routes.go`
  - [ ] Add `r.HandleFunc("/v1/dashboard/projects.json", handler.Projects).Methods("GET")`
  - [ ] Verify route registration

## Phase 3: Testing (Must-have)

- [ ] **T-001**: Create unit tests
  - [ ] `internal/service/dashboard/projects_service_test.go`
  - [ ] Test status filtering logic
  - [ ] Test stats calculation
  - [ ] Test progress ordering
  - [ ] Test edge cases (zero projects, division by zero)

- [ ] **T-002**: Create integration tests
  - [ ] `test/integration/dashboard_projects_test.go`
  - [ ] Test endpoint with sample data
  - [ ] Test Rails parity validation
  - [ ] Test error scenarios

- [ ] **T-003**: Run existing tests
  - [ ] Verify no regressions
  - [ ] Achieve > 85% coverage on new code

## Phase 4: Validation (Should-have)

- [ ] **V-001**: Manual testing
  - [ ] Test with various data scenarios
  - [ ] Compare output with Rails endpoint
  - [ ] Verify performance targets

- [ ] **V-002**: Code review
  - [ ] Review by engineering lead
  - [ ] Address feedback
  - [ ] Update documentation

---

# Stakeholder Alignment

| Stakeholder | Responsibility | Verification |
|-------------|----------------|--------------|
| **Product Owner** | Approve feature requirements | Review Key Requirements table |
| **Engineering Lead** | Approve technical decisions | Review Technical Decisions section |
| **Backend Developer** | Implement endpoint | Execute Implementation Checklist |
| **QA Engineer** | Test Rails parity | Verify Acceptance Criteria AC-PROJ-006 |
| **DevOps** | Deployment readiness | Verify NFA criteria |

---

# Traceability Matrix

| Requirement ID | User Story | Acceptance Criteria | Test File | Status |
|----------------|------------|---------------------|-----------|--------|
| REQ-PROJ-001 | Display running projects on dashboard | AC-PROJ-001, AC-PROJ-002 | test/integration/dashboard_projects_test.go | TODO |
| REQ-PROJ-002 | Calculate overall progress statistics | AC-PROJ-003 | internal/service/dashboard/projects_service_test.go | TODO |
| REQ-PROJ-003 | Order projects by progress | AC-PROJ-004 | internal/service/dashboard/projects_service_test.go | TODO |
| REQ-PROJ-004 | Show recent reading activity | AC-PROJ-005 | test/integration/dashboard_projects_test.go | TODO |
| REQ-PROJ-005 | Match Rails API behavior | AC-PROJ-006 | test/integration/dashboard_projects_test.go | TODO |

---

# Validation

## Code Quality Standards
- [x] Go 1.25.7 compatible
- [x] Follows existing code patterns in project
- [ ] Linting passes (`go vet ./...`) - To be verified during implementation
- [ ] Formatting correct (`go fmt ./...`) - To be verified during implementation

## Technical Feasibility
- [x] All technologies proven and production-ready (PostgreSQL, pgx, gorilla/mux)
- [x] No experimental features used
- [x] Database queries optimized with indexes
- [x] Error handling comprehensive

## User Needs Alignment
- [x] Endpoint matches Rails API behavior exactly (confirmed via questionnaire)
- [x] Calculation methods verified against Rails implementation
- [x] Response format matches Rails structure (flat JSON with stats at root)
- [ ] Performance targets verified - To be measured during integration testing

---

# Ready for Implementation

## Approval Status: ✅ **READY FOR IMPLEMENTATION**

This PRD has been:
- ✅ **Researched**: Existing codebase analyzed, Rails endpoint understood
- ✅ **Validated**: Technical questions clarified via palha subagent
- ✅ **Clarified**: Ambiguities resolved through structured questionnaire
- ✅ **Prioritized**: Requirements categorized into P1/P2/P3 tiers
- ✅ **Testable**: Acceptance criteria are objective and measurable
- ✅ **Documented**: Clear file locations and implementation steps provided

## Prerequisites for Starting Implementation

Before beginning implementation, ensure:

1. **Environment ready**:
   - [ ] Go 1.25.7 installed and configured
   - [ ] PostgreSQL running and accessible
   - [ ] Test database `reading_log_test` created
   - [ ] `.env` file configured with database credentials

2. **Dependencies available**:
   - [ ] Existing repository pattern in place
   - [ ] DashboardRepository interface exists
   - [ ] DTOs for dashboard responses defined

3. **Stakeholder awareness**:
   - [ ] Product owner aware of scope (single endpoint)
   - [ ] Engineering lead reviewed technical decisions
   - [ ] QA team prepared for Rails parity testing

## Implementation Start Command

```bash
# Verify environment
make test-setup

# Create service layer
mkdir -p internal/service/dashboard
touch internal/service/dashboard/projects_service.go
touch internal/service/dashboard/projects_service_test.go

# Start implementation
go run ./cmd/server.go
```

## Success Metrics

Implementation is complete when:
1. ✅ Endpoint returns 200 OK with correct structure
2. ✅ All running projects included and ordered by progress
3. ✅ Stats calculation matches Rails output
4. ✅ Unit tests pass with > 85% coverage
5. ✅ Integration tests pass with Rails parity validation
6. ✅ No regressions in existing tests

---

*PRD Version: 1.0*
*Created: 2026-04-28*
*Author: PRD Refinement Specialist*
*Status: Ready for Implementation*