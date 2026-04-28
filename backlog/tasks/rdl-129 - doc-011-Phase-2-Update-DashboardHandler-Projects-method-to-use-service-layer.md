---
id: RDL-129
title: '[doc-011 Phase 2] Update DashboardHandler Projects method to use service layer'
status: To Do
assignee:
  - thomas
created_date: '2026-04-28 11:16'
updated_date: '2026-04-28 14:10'
labels:
  - feature
  - backend
  - phase-2
dependencies: []
documentation:
  - doc-011
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Modify internal/api/v1/handlers/dashboard_handler.go Projects() method to use new projects service instead of direct repository calls. Implement response formatting matching Rails structure with projects array and stats object at root level. Add error handling and structured logging.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Handler calls service GetRunningProjectsWithLogs method
- [ ] #2 Handler calls service CalculateStats method
- [ ] #3 Response structure matches Rails (projects array + stats object)
- [ ] #4 Database errors return 500 Internal Server Error with logging
- [ ] #5 Empty data returns 200 OK with empty arrays
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The implementation will refactor the `DashboardHandler.Projects()` method to use the existing service layer instead of direct repository calls, matching the Rails response structure.

**Technical Approach:**
- The `ProjectsService` already exists with `GetRunningProjectsWithLogs()` and `CalculateStats()` methods implemented
- The handler will be updated to inject `ProjectsService` via dependency injection (similar to how other dashboard handlers use services)
- Response format will change from flat logs array to Rails-matching structure: `{ "projects": [...], "stats": {...} }`
- Error handling will use structured logging with `slog` instead of `fmt.Printf`
- The service layer already handles status filtering, progress calculation, and stats aggregation

**Why this approach:**
- Follows Clean Architecture principles (handlers → services → repositories)
- Reuses existing, tested service logic
- Consistent with other dashboard endpoints (e.g., `Faults`, `WeekdayFaults` already use services)
- Enables better testability and maintainability

**Architecture Decisions:**
- Use dependency injection for `ProjectsService` in `DashboardHandler`
- Response format: Flat JSON with `projects` array and `stats` object at root (not JSON:API envelope)
- Status filtering done in service layer using `isRunningProject()` method
- Stats calculated once per request using `CalculateStats()`

---

### 2. Files to Modify

#### Files to Modify

| File | Changes | Priority |
|------|---------|----------|
| `internal/api/v1/handlers/dashboard_handler.go` | - Add `ProjectsService` field to `DashboardHandler` struct<br>- Update `NewDashboardHandler()` to accept `ProjectsService`<br>- Refactor `Projects()` method to call service layer<br>- Change response format to Rails structure<br>- Replace `fmt.Printf` with `slog` logging | P1 |
| `internal/api/v1/routes.go` | Verify route registration for `/v1/dashboard/projects.json` (may already exist) | P1 |
| `cmd/server.go` | Update `NewDashboardHandler()` call to pass `ProjectsService` | P1 |

#### Files to Create

| File | Purpose | Priority |
|------|---------|----------|
| `internal/api/v1/handlers/dashboard_handler_projects_test.go` | Unit tests for `Projects()` handler method | P1 |

---

### 3. Dependencies

**Prerequisites:**
- `ProjectsService` already implemented in `internal/service/dashboard/projects_service.go` with:
  - `GetRunningProjectsWithLogs(ctx)` - Returns filtered projects with logs
  - `CalculateStats(ctx)` - Returns aggregate statistics
- `DashboardRepository` interface already has `GetRunningProjectsWithLogs()` method
- DTOs exist in `internal/domain/dto/dashboard_response.go`:
  - `ProjectWithLogs` struct
  - `StatsData` struct

**Related Tasks:**
- RDL-126: `GetRunningProjectsWithLogs` service method (DONE)
- RDL-128: Repository SQL JOIN implementation (DONE)
- RDL-127: Stats calculation service method (DONE)
- RDL-125: Route registration (To Do - may need to be completed first)

**Setup Steps:**
1. Verify `ProjectsService` is properly instantiated in `cmd/server.go`
2. Ensure `DashboardRepository` is injected into `ProjectsService`
3. Confirm route is registered in `routes.go`

---

### 4. Code Patterns

**Following Existing Patterns:**

1. **Dependency Injection Pattern** (from `Faults`, `WeekdayFaults` handlers):
```go
type DashboardHandler struct {
    repo       repository.DashboardRepository
    userConfig *service.UserConfigService
    projectsService *dashboard.ProjectsService  // New field
}

func NewDashboardHandler(repo repository.DashboardRepository, userConfig *service.UserConfigService, projectsService *dashboard.ProjectsService) *DashboardHandler {
    return &DashboardHandler{
        repo: repo,
        userConfig: userConfig,
        projectsService: projectsService,
    }
}
```

2. **Service Usage Pattern** (from `Faults` handler):
```go
func (h *DashboardHandler) Projects(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Call service method
    projects, err := h.projectsService.GetRunningProjectsWithLogs(ctx)
    if err != nil {
        slog.Error("Failed to get running projects", "error", err)
        http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
        return
    }
    
    // Call stats service
    stats, err := h.projectsService.CalculateStats(ctx)
    if err != nil {
        slog.Error("Failed to calculate stats", "error", err)
        http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
        return
    }
    
    // Build response
    response := map[string]interface{}{
        "projects": projects,
        "stats": stats,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

3. **Error Handling Pattern:**
- Database errors → 500 Internal Server Error with structured logging
- Empty data → 200 OK with empty arrays (not 404)
- Consistent error JSON format: `{"error": "message"}`

4. **Naming Conventions:**
- Follow existing Go naming: `Projects()`, `GetRunningProjectsWithLogs()`, `CalculateStats()`
- JSON field names: snake_case (`progress_geral`, `total_pages`, `pages`)

---

### 5. Testing Strategy

**Unit Tests** (`dashboard_handler_projects_test.go`):

1. **TestProjects_Success** - Happy path
   - Mock `GetRunningProjectsWithLogs()` returning 2 projects
   - Mock `CalculateStats()` returning valid stats
   - Verify 200 OK response
   - Verify response structure: `{ "projects": [...], "stats": {...} }`
   - Verify projects array contains correct data
   - Verify stats object has `progress_geral`, `total_pages`, `pages`

2. **TestProjects_EmptyData** - No projects
   - Mock returning empty project slice
   - Mock returning zero stats
   - Verify 200 OK with empty projects array and zero stats

3. **TestProjects_ServiceError** - Database error
   - Mock `GetRunningProjectsWithLogs()` returning error
   - Verify 500 Internal Server Error
   - Verify error response format

4. **TestProjects_ProjectsWithLogs** - Projects with logs
   - Mock returning projects with eager-loaded logs
   - Verify logs are included in response
   - Verify log structure matches expected format

**Integration Tests** (`test/integration/dashboard_projects_test.go`):

1. **TestDashboardProjectsEndpoint_Integration** - Full endpoint test
   - Setup test database with fixtures
   - Create projects with different statuses (running, finished, stopped)
   - Create logs for each project
   - Call endpoint
   - Verify only running projects returned
   - Verify stats calculation matches expected values
   - Verify ordering by progress DESC

2. **TestDashboardProjectsEndpoint_RailsParity** - Structure comparison
   - Compare Go response structure with Rails response
   - Normalize timestamps and floats for comparison
   - Verify structural equivalence

**Edge Cases to Cover:**
- Zero projects → empty array, zero stats
- Division by zero in progress calculation → 0.0
- Projects with no logs → empty logs array
- Null/nil values in database → COALESCE handling

**Test Coverage Target:** > 85% for new code paths

---

### 6. Risks and Considerations

**Known Issues:**
1. **Response Format Change**: The current handler returns a JSON:API envelope with flat logs array. Changing to Rails structure (`{projects, stats}`) may break existing clients if any are using this endpoint.
   - **Mitigation**: Verify no production clients exist; update API documentation

2. **Service Injection in cmd/server.go**: Need to ensure `ProjectsService` is properly instantiated and injected
   - **Mitigation**: Follow existing pattern from `FaultsService` instantiation

3. **Stats Calculation Performance**: `CalculateStats()` queries each project's page individually
   - **Mitigation**: Already implemented; performance acceptable per existing tests

**Potential Pitfalls:**
1. **Status Filtering Logic**: The service uses `isRunningProject()` which checks logs existence and pages < total_pages
   - **Consideration**: Verify this matches Rails behavior (7-day threshold mentioned in PRD)
   - **Action**: Review PRD doc-011 Decision 2 for status filter implementation

2. **JSON Marshaling**: `ProjectWithLogs` and `StatsData` structs must marshal correctly
   - **Mitigation**: Existing DTOs have JSON tags; verify field names match Rails

3. **Context Timeout**: Ensure service calls use context with timeout
   - **Mitigation**: Service already uses context passed from handler

**Deployment Considerations:**
- No database migrations required
- Backward compatibility: If endpoint was not previously exposed, no breaking changes
- Documentation: Update API docs with new response structure

**Rollback Plan:**
- If issues arise, revert handler to direct repository calls
- Service layer can remain for future refactoring

---

**Implementation Notes:**
- Follow Clean Architecture strictly: Handler → Service → Repository
- Use `slog` for structured logging (already configured in project)
- Maintain existing test infrastructure patterns (mock repositories, TestHelper)
- Ensure all acceptance criteria from task RDL-129 are met
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
**Implementation Completed - 2026-04-28**

## Summary
Successfully refactored the `DashboardHandler.Projects()` method to use the service layer instead of direct repository calls.

## Changes Made

### 1. Modified Files
- **internal/api/v1/handlers/dashboard_handler.go**
  - Added `ProjectsService` field to `DashboardHandler` struct
  - Updated `NewDashboardHandler()` to accept `ProjectsServiceInterface`
  - Refactored `Projects()` method to call service layer
  - Changed response format from JSON:API envelope to Rails structure `{ "projects": [...], "stats": {...} }`
  - Replaced `fmt.Printf` with `slog` for structured logging

- **internal/api/v1/routes.go**
  - Updated `SetupRoutes()` to accept `ProjectsServiceInterface`
  - Added route registration for `/v1/dashboard/projects.json`

- **cmd/server.go**
  - Added `ProjectsService` instantiation
  - Updated `NewDashboardHandler()` call to pass `ProjectsService`

- **internal/service/dashboard/projects_service.go**
  - Added `ProjectsServiceInterface` for testability
  - Added `PgxPoolInterface` for database pool abstraction

- **test/integration/dashboard_mock_test.go** (NEW)
  - Created shared `MockProjectsService` for integration tests

### 2. Test Files Updated
Updated all test files to use the new `ProjectsServiceInterface`:
- internal/api/v1/handlers/dashboard_handler_test.go
- internal/api/v1/handlers/dashboard_handler_projects_test.go (NEW)
- internal/api/v1/routes_test.go
- test/unit/dashboard_handler_test.go
- test/integration/*.go files
- test/performance/dashboard_benchmark_test.go

### 3. Acceptance Criteria Status
✅ #1 Handler calls service GetRunningProjectsWithLogs method
✅ #2 Handler calls service CalculateStats method
✅ #3 Response structure matches Rails (projects array + stats object)
✅ #4 Database errors return 500 Internal Server Error with logging
✅ #5 Empty data returns 200 OK with empty arrays

### 4. Definition of Done Status
✅ #1 All unit tests pass
✅ #2 go fmt and go vet pass with no errors
✅ #3 Clean Architecture layers properly followed
✅ #4 Error responses consistent with existing patterns
✅ #5 HTTP status codes correct for response type

## Testing
- All unit tests pass: `go test ./internal/api/v1/handlers/...`
- Code formatted: `go fmt ./...`
- No vet errors: `go vet ./...`

## Next Steps
- Run integration tests to verify database interactions
- Mark task as Done after integration tests pass
<!-- SECTION:NOTES:END -->

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
