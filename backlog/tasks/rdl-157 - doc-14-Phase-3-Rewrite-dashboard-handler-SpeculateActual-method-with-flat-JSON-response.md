---
id: RDL-157
title: >-
  [doc-14 Phase 3] Rewrite dashboard handler SpeculateActual method with flat
  JSON response
status: To Do
assignee:
  - thomas
created_date: '2026-05-10 10:48'
updated_date: '2026-05-10 14:25'
labels:
  - handler
  - api
  - phase-3
dependencies: []
documentation:
  - doc-014
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Rewrite internal/api/v1/handlers/dashboard_handler.go SpeculateActual() method to inject SpeculateService, call service methods to generate chart config, and return flat JSON { echart: chartConfig } matching Rails response format. Remove current faults-based implementation.

Handler must handle errors gracefully, return appropriate HTTP status codes (200 OK, 500 Internal Server Error), and follow existing error handling patterns.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 SpeculateActual method injects SpeculateService via dependency injection
- [x] #2 Handler calls service to generate ECharts configuration
- [x] #3 Response format is flat JSON { echart: {...} } without JSON:API envelope
- [x] #4 Error handling returns 500 status with proper error message
- [x] #5 Handler follows existing middleware and logging patterns
- [x] #6 Code compiles without errors
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The implementation will rewrite the `SpeculateActual` handler method in `dashboard_handler.go` to use the existing `SpeculateService` (already implemented in RDL-155) and return flat JSON response matching the Rails API format.

**Current State Analysis:**
- The current `SpeculateActual` method (lines 425-495 in `dashboard_handler.go`) uses a faults-based implementation that calculates predicted faults
- The `SpeculateService` (implemented in RDL-155) already has all the calculation logic for pages-based speculation
- The handler needs to be refactored to inject and use `SpeculateService` instead of calculating faults directly

**Implementation Strategy:**
1. **Dependency Injection**: Add `SpeculateService` to `DashboardHandler` struct and inject it via constructor
2. **Handler Rewrite**: Replace the current faults-based implementation with a call to `SpeculateService.GenerateChartConfig()`
3. **Response Format**: Return flat JSON `{ echart: {...} }` instead of JSON:API envelope to match Rails API
4. **Error Handling**: Follow existing error handling patterns (500 status with error message)

**Why This Approach:**
- Leverages existing, tested `SpeculateService` logic from RDL-155
- Follows the same pattern as `Faults` handler which also uses service layer
- Flat JSON response matches Rails API format and simplifies frontend consumption
- Maintains Clean Architecture separation (handler → service → repository)

**Key Design Decisions:**
- **Service Injection**: Add `*dashboard.SpeculateService` to `DashboardHandler` struct (not created inline like Faults)
- **Response Structure**: Use `map[string]interface{}{"echart": chartConfig}` for flat JSON
- **Content-Type**: Use `application/json` (not `application/vnd.api+json` since no JSON:API envelope)
- **Error Responses**: Return `{"error": "Internal server error"}` with 500 status

### 2. Files to Modify

**Handler Layer (Modify):**
- `internal/api/v1/handlers/dashboard_handler.go`
  - **Modify `DashboardHandler` struct**: Add `speculateService *dashboard.SpeculateService` field
  - **Modify `NewDashboardHandler` constructor**: Accept `*dashboard.SpeculateService` parameter and assign to struct
  - **Rewrite `SpeculateActual` method** (lines 425-495):
    - Remove faults-based implementation
    - Call `h.speculateService.GenerateChartConfig(ctx)`
    - Return flat JSON `{ echart: chartConfig }`
    - Handle errors with 500 status

**Routes (Modify):**
- `internal/api/v1/routes.go`
  - **Modify `SetupRoutes` function**: Add `SpeculateService` parameter
  - **Pass SpeculateService** to `NewDashboardHandler` constructor
  - **Add route registration**: `r.HandleFunc("/v1/dashboard/echart/speculate_actual.json", dashboardHandler.SpeculateActual).Methods("GET")`

**Main Server (Modify):**
- `cmd/server.go`
  - **Create SpeculateService instance**: `speculateService := dashboard.NewSpeculateService(dashboardRepo, userConfig)`
  - **Pass to SetupRoutes**: Add `speculateService` parameter

**No Changes Required:**
- `internal/service/dashboard/speculate_service.go` - Already implemented in RDL-155
- `internal/domain/dto/dashboard_response.go` - EchartConfig DTO already exists
- `internal/repository/dashboard_repository.go` - All repository methods already exist

### 3. Dependencies

**Prerequisites (Already Completed):**
- ✅ RDL-155: SpeculateService with weekday-based calculation methods implemented
- ✅ RDL-152: Repository interface methods for weekday-based mean calculations
- ✅ RDL-153: PostgreSQL repository implementations for weekday grouping queries
- ✅ RDL-151: ECharts DTOs extended with MarkPoint, MarkLine, BoundaryGap fields
- ✅ RDL-154: Mock repository implementations for unit testing

**Service Dependencies:**
- `dashboard.SpeculateService` - Must be instantiated before handler creation
- `repository.DashboardRepository` - Injected into SpeculateService
- `service.UserConfigService` - Injected into SpeculateService (for prediction percentage)

**No Blocking Issues:**
- All service layer code is complete and tested
- All repository methods are implemented
- DTOs support all required fields
- Only handler and routing layer changes needed

### 4. Code Patterns

**Follow Existing Patterns:**

1. **Handler Struct Pattern:**
```go
type DashboardHandler struct {
    repo            repository.DashboardRepository
    userConfig      *service.UserConfigService
    projectsService dashboard.ProjectsServiceInterface
    speculateService *dashboard.SpeculateService  // NEW
}
```

2. **Constructor Injection Pattern:**
```go
func NewDashboardHandler(
    repo repository.DashboardRepository,
    userConfig *service.UserConfigService,
    projectsService dashboard.ProjectsServiceInterface,
    speculateService *dashboard.SpeculateService,  // NEW
) *DashboardHandler {
    return &DashboardHandler{
        repo:            repo,
        userConfig:      userConfig,
        projectsService: projectsService,
        speculateService: speculateService,  // NEW
    }
}
```

3. **Service Usage Pattern (like Faults handler):**
```go
func (h *DashboardHandler) SpeculateActual(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    chartConfig, err := h.speculateService.GenerateChartConfig(ctx)
    if err != nil {
        slog.Error("Failed to generate speculate actual chart", "error", err)
        http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
        return
    }
    
    response := map[string]interface{}{
        "echart": chartConfig,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

4. **Route Registration Pattern:**
```go
r.HandleFunc("/v1/dashboard/echart/speculate_actual.json", dashboardHandler.SpeculateActual).Methods("GET")
```

5. **Error Logging Pattern:**
```go
slog.Error("Failed to generate speculate actual chart", "error", err)
```

**Naming Conventions:**
- Field name: `speculateService` (camelCase, matches Go conventions)
- Method name: `SpeculateActual` (PascalCase, matches existing handler methods)
- Route path: `/v1/dashboard/echart/speculate_actual.json` (snake_case, matches Rails)

**Integration with Existing Code:**
- Use `slog` for error logging (consistent with other handlers)
- Use `json.NewEncoder` for response encoding (consistent with other handlers)
- Follow Clean Architecture: handler → service → repository

### 5. Testing Strategy

**Unit Tests (RDL-158 - Separate Task):**
- Tests will be created in `test/unit/api/v1/handlers/dashboard_handler_speculate_actual_test.go`
- Test scenarios:
  - `TestSpeculateActual_Success`: Valid response with 200 status
  - `TestSpeculateActual_ServiceError`: Service error returns 500 status
  - `TestSpeculateActual_ResponseFormat`: Verify flat JSON structure with `echart` key

**Integration Tests (RDL-160 - Separate Task):**
- Tests will be created in `test/integration/api/v1/dashboard/echart_speculate_actual_test.go`
- Test scenarios:
  - `TestResponseFormat`: Verify flat JSON, no JSON:API envelope
  - `TestDateRange`: Verify 15 days of data
  - `TestSeriesNames`: Verify 'Pages' and 'Mean' series
  - `TestEmptyDatabase`: Verify zero-fill behavior
  - `TestRailsParity`: Compare with Rails API output

**Manual Testing:**
```bash
# Start server
go run ./cmd/server.go

# Test endpoint
curl http://localhost:3000/v1/dashboard/echart/speculate_actual.json

# Verify response structure
curl http://localhost:3000/v1/dashboard/echart/speculate_actual.json | jq
```

**Expected Response Format:**
```json
{
  "echart": {
    "title": "Speculated vs Actual",
    "tooltip": {
      "trigger": "axis",
      "formatter": "{a} <br/>{b}: {c}"
    },
    "legend": {
      "show": true,
      "data": ["Actual", "Speculated"]
    },
    "xAxis": {
      "type": "category",
      "data": ["10-05 (Sat)", "11-05 (Sun)", ...]
    },
    "yAxis": {
      "type": "value"
    },
    "grid": {
      "left": "3%",
      "right": "4%",
      "top": "15%",
      "bottom": "3%"
    },
    "series": [
      {
        "name": "Actual",
        "type": "line",
        "data": [0, 15, 30, ...],
        "itemStyle": { "color": "#5470C6" },
        "lineStyle": { "width": 2 }
      },
      {
        "name": "Speculated",
        "type": "line",
        "data": [0, 17, 35, ...],
        "itemStyle": { "color": "#91CC75" },
        "lineStyle": { "width": 2, "type": "dashed" }
      }
    ]
  }
}
```

**Validation Checklist:**
- ✅ Response is flat JSON with `echart` key at root
- ✅ No JSON:API envelope (`data`, `type`, `attributes`)
- ✅ Content-Type is `application/json`
- ✅ ECharts config has all required fields
- ✅ Series names are 'Actual' and 'Speculated'
- ✅ 15 data points in each series
- ✅ Date labels in 'DD-MM (Day)' format

### 6. Risks and Considerations

**Technical Risks:**

1. **Service Injection Order:**
   - **Risk**: SpeculateService must be created before DashboardHandler
   - **Mitigation**: Follow dependency order: repository → service → handler
   - **Code Order in `cmd/server.go`:**
     ```go
     // 1. Create repositories
     dashboardRepo := postgres.NewDashboardRepositoryImpl(dbPool)
     
     // 2. Create services
     speculateService := dashboard.NewSpeculateService(dashboardRepo, userConfig)
     projectsService := dashboard.NewProjectsService(dashboardRepo)
     
     // 3. Create handlers
     dashboardHandler := handlers.NewDashboardHandler(dashboardRepo, userConfig, projectsService, speculateService)
     ```

2. **Route Registration Conflict:**
   - **Risk**: Route might already exist in test files
   - **Mitigation**: Check existing routes before adding, ensure no duplicates
   - **Verification**: Run `go build` to catch compilation errors

3. **Response Format Breaking Change:**
   - **Risk**: Frontend might expect JSON:API envelope
   - **Mitigation**: Document the flat JSON format clearly, coordinate with frontend team
   - **Note**: This matches Rails API format, so frontend should already handle it

4. **Service Method Availability:**
   - **Risk**: `GenerateChartConfig` method might not exist or have wrong signature
   - **Mitigation**: Verify RDL-155 implementation is complete
   - **Expected Signature**: `func (s *SpeculateService) GenerateChartConfig(ctx context.Context) (*dto.EchartConfig, error)`

**Design Considerations:**

1. **Service Lifetime:**
   - **Decision**: SpeculateService is created once at startup and reused
   - **Rationale**: Service is stateless, thread-safe for concurrent requests
   - **Pattern**: Singleton service injected into handler (same as ProjectsService)

2. **Error Response Consistency:**
   - **Decision**: Use same error format as other handlers
   - **Format**: `{"error": "Internal server error"}`
   - **Status**: 500 Internal Server Error
   - **Logging**: Use slog.Error for error tracking

3. **Content-Type Selection:**
   - **Decision**: Use `application/json` (not `application/vnd.api+json`)
   - **Rationale**: No JSON:API envelope, so standard JSON content type
   - **Alignment**: Matches Rails API content type

**Performance Considerations:**

1. **Service Call Overhead:**
   - **Impact**: Minimal - service method is already optimized
   - **Query**: Uses single repository query for logs
   - **Memory**: 15 data points is negligible

2. **Context Timeout:**
   - **Current**: SpeculateService uses repository context (no explicit timeout)
   - **Recommendation**: Verify repository methods use 15-second timeout
   - **Pattern**: `ctx, cancel := context.WithTimeout(ctx, 15*time.Second)`

**Testing Considerations:**

1. **Mock Service for Unit Tests:**
   - **Strategy**: Create mock SpeculateService interface
   - **Test Cases**: Success case, error case, nil response case
   - **File**: `test/unit/api/v1/handlers/dashboard_handler_speculate_actual_test.go`

2. **Integration Test Data:**
   - **Strategy**: Use test fixtures with known data
   - **Verification**: Compare output against expected values
   - **File**: `test/integration/api/v1/dashboard/echart_speculate_actual_test.go`

**Rollout Plan:**

1. **Implementation Phase (This Task - RDL-157):**
   - Modify handler and routes
   - Verify compilation
   - Manual testing

2. **Testing Phase (RDL-158, RDL-160):**
   - Unit tests for handler
   - Integration tests for endpoint
   - Rails parity verification

3. **Documentation Phase (RDL-162, RDL-163):**
   - Update AGENTS.md
   - Create implementation guide
   - Update API documentation

**Acceptance Criteria Mapping:**

| AC | Implementation | Verification |
|----|----------------|--------------|
| #1 SpeculateActual injects SpeculateService | Add field to struct, inject via constructor | Code review, compilation |
| #2 Handler calls service to generate chart | Call `GenerateChartConfig(ctx)` | Unit test `TestSpeculateActual_Success` |
| #3 Response is flat JSON `{ echart: {...} }` | Use `map[string]interface{}{"echart": ...}` | Integration test `TestResponseFormat` |
| #4 Error handling returns 500 status | Check error, return `http.StatusInternalServerError` | Unit test `TestSpeculateActual_ServiceError` |
| #5 Follows middleware and logging patterns | Use `slog.Error`, standard middleware | Code review |
| #6 Code compiles without errors | `go build ./...` | Build verification |

**Definition of Done Checklist:**

- [ ] #1 All unit tests pass (RDL-158)
- [ ] #2 All integration tests pass (RDL-160)
- [ ] #3 `go fmt` and `go vet` pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md and AGENTS.md (RDL-163)
- [ ] #8 New code paths include error path tests (RDL-158)
- [ ] #9 HTTP handlers test both success and error responses (RDL-158)
- [ ] #10 Integration tests verify actual database interactions (RDL-160)

**Related Tasks:**
- RDL-155: SpeculateService implementation (DONE)
- RDL-158: Unit tests for handler (NEXT)
- RDL-159: Route registration (NEXT)
- RDL-160: Integration tests (NEXT)
- RDL-162: Implementation guide documentation (NEXT)
- RDL-163: API documentation update (NEXT)
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress

### Completed Steps (2026-05-10)

1. ✅ **Modified DashboardHandler struct** - Added `speculateService dashboard.SpeculateServiceInterface` field
2. ✅ **Modified NewDashboardHandler constructor** - Added `speculateService dashboard.SpeculateServiceInterface` parameter
3. ✅ **Rewrote SpeculateActual method** - Now uses SpeculateService.GenerateChartConfig() and returns flat JSON `{ echart: {...} }`
4. ✅ **Updated routes.go** - Added speculateService parameter and registered new route `/v1/dashboard/echart/speculate_actual.json`
5. ✅ **Updated cmd/server.go** - Created SpeculateService instance and passed to SetupRoutes
6. ✅ **Updated test/integration/test_context.go** - Updated SetupRoutes wrapper and Setup function
7. ✅ **Created SpeculateServiceInterface** - Added interface in speculate_service.go for testability
8. ✅ **Updated all test files** - Added MockSpeculateService and updated all NewDashboardHandler calls
9. ✅ **Code compiles** - `go build ./...` passes without errors
10. ✅ **go vet passes** - No issues found
11. ✅ **All tests pass** - `go test ./internal/api/v1/handlers/` passes

### Code Changes Summary

**Modified Files:**
- `internal/api/v1/handlers/dashboard_handler.go` - Handler struct, constructor, and SpeculateActual method
- `internal/api/v1/routes.go` - SetupRoutes function signature and route registration
- `cmd/server.go` - Service instantiation
- `internal/service/dashboard/speculate_service.go` - Added SpeculateServiceInterface
- `test/integration/test_context.go` - Test context setup
- `internal/api/v1/handlers/dashboard_handler_test.go` - Added MockSpeculateService
- `test/dashboard_integration_test.go` - Added MockSpeculateService and updated calls
- `test/unit/dashboard_handler_test.go` - Added MockSpeculateService and updated calls
- `test/performance/dashboard_benchmark_test.go` - Added MockSpeculateService and updated calls
- `test/integration/dashboard_mock_test.go` - Added MockSpeculateService
- `test/integration/dashboard_day_permean_integration_test.go` - Updated calls
- `test/integration/dashboard_projects_test.go` - Updated calls
- `test/integration/dashboard_stats_integration_test.go` - Updated calls
- `test/integration/error_scenarios_test.go` - Updated calls
- `test/integration/projects_integration_test.go` - Updated SetupRoutes calls
- `test/integration/projects_rails_comparison_test.go` - Updated SetupRoutes calls
- `internal/api/v1/routes_test.go` - Added MockSpeculateService and updated calls

### Acceptance Criteria Status

- ✅ #1 SpeculateActual method injects SpeculateService via dependency injection
- ✅ #2 Handler calls service to generate ECharts configuration
- ✅ #3 Response format is flat JSON { echart: {...} } without JSON:API envelope
- ✅ #4 Error handling returns 500 status with proper error message
- ✅ #5 Handler follows existing middleware and logging patterns
- ✅ #6 Code compiles without errors

### Next Steps

- Run go fmt
- Check all acceptance criteria
- Mark task as Done
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary

Rewrote the `SpeculateActual` handler method in `internal/api/v1/handlers/dashboard_handler.go` to use the existing `SpeculateService` and return flat JSON response matching the Rails API format.

## What Was Done

1. **Handler Layer Changes**:
   - Modified `DashboardHandler` struct to inject `speculateService dashboard.SpeculateServiceInterface`
   - Updated `NewDashboardHandler` constructor to accept `SpeculateServiceInterface` parameter
   - Rewrote `SpeculateActual()` method to:
     - Call `h.speculateService.GenerateChartConfig(ctx)` instead of faults-based implementation
     - Return flat JSON `{ echart: {...} }` without JSON:API envelope
     - Use `application/json` content type
     - Handle errors with 500 status and `slog.Error` logging

2. **Service Layer Changes**:
   - Added `SpeculateServiceInterface` in `internal/service/dashboard/speculate_service.go` for testability

3. **Routing Changes**:
   - Updated `SetupRoutes()` in `internal/api/v1/routes.go` to accept `SpeculateServiceInterface`
   - Registered new route: `/v1/dashboard/echart/speculate_actual.json`

4. **Server Initialization**:
   - Updated `cmd/server.go` to create `SpeculateService` instance and pass to `SetupRoutes`

5. **Test Updates**:
   - Created `MockSpeculateService` in multiple test files
   - Updated all `NewDashboardHandler` and `SetupRoutes` calls across 17 test files

## Key Changes

- **Files Modified**: 17 files
- **Response Format**: Changed from JSON:API envelope to flat JSON `{ echart: {...} }`
- **Content-Type**: Changed from `application/vnd.api+json` to `application/json`
- **Implementation**: Replaced faults-based calculation with pages-based speculation from `SpeculateService`

## Tests Run

- ✅ `go build ./...` - Compiles without errors
- ✅ `go vet ./...` - No issues found
- ✅ `go fmt ./...` - Code formatted
- ✅ `go test ./internal/api/v1/handlers/` - All tests pass
- ✅ `TestDashboardHandler_SpeculateActual` - Tests flat JSON response format

## Verification

- All 6 acceptance criteria checked:
  1. ✅ SpeculateActual method injects SpeculateService via dependency injection
  2. ✅ Handler calls service to generate ECharts configuration
  3. ✅ Response format is flat JSON { echart: {...} } without JSON:API envelope
  4. ✅ Error handling returns 500 status with proper error message
  5. ✅ Handler follows existing middleware and logging patterns
  6. ✅ Code compiles without errors

## Risks/Follow-ups

- **Breaking Change**: Response format changed from JSON:API to flat JSON (matches Rails API)
- **Frontend Impact**: Frontend should already handle Rails API format
- **Documentation**: AGENTS.md and API documentation to be updated in RDL-163
- **Testing**: Full unit and integration tests for the endpoint to be added in RDL-158 and RDL-160
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass
- [ ] #2 All integration tests pass execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md and AGENTS.md
- [ ] #8 New code paths include error path tests
- [ ] #9 HTTP handlers test both success and error responses
- [ ] #10 Integration tests verify actual database interactions
<!-- DOD:END -->
