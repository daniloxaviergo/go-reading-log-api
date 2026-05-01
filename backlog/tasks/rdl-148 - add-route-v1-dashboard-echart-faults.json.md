---
id: RDL-148
title: add route v1/dashboard/echart/faults.json
status: To Do
assignee:
  - catarina
created_date: '2026-05-01 19:29'
updated_date: '2026-05-01 19:35'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
check if the echart faults is implemented and add route v1/dashboard/echart/faults.json
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The task requires adding the route `/v1/dashboard/echart/faults.json` to the API. After thorough codebase research, the implementation is already complete at the handler and service layers - only the route registration is missing.

**Implementation Strategy:**
- Add route registration in `internal/api/v1/routes.go` following the existing pattern for other ECharts endpoints
- Add route verification test in `internal/api/v1/routes_test.go`
- Verify all existing tests pass (unit, integration, performance)

**Why this approach:**
- The `DashboardHandler.Faults()` method is already fully implemented
- The `FaultsService` is complete with `GetFaultsPercentage()` and `CreateGaugeChart()` methods
- Unit and integration tests already exist and validate the functionality
- Only the route wiring is missing, making this a minimal, low-risk change

**Architecture Alignment:**
- Follows Clean Architecture: Handler layer calls service layer, which calls repository layer
- Uses existing dependency injection pattern with `UserConfigProvider`
- Response format matches JSON:API envelope pattern used by other ECharts endpoints

### 2. Files to Modify

| File | Action | Description |
|------|--------|-------------|
| `internal/api/v1/routes.go` | **Modify** | Add route registration: `r.HandleFunc("/v1/dashboard/echart/faults.json", dashboardHandler.Faults).Methods("GET")` |
| `internal/api/v1/routes_test.go` | **Modify** | Add test case to verify the route responds correctly (similar to existing route tests) |

**No new files required** - all handler, service, repository, and DTO implementations already exist.

### 3. Dependencies

**Prerequisites (All Complete):**
- ✅ `DashboardHandler.Faults()` method implemented in `internal/api/v1/handlers/dashboard_handler.go`
- ✅ `FaultsService` implemented in `internal/service/dashboard/faults_service.go`
- ✅ `DashboardRepository.GetFaultsByDateRange()` method exists
- ✅ `UserConfigService` with `GetMaxFaults()` and `GetPredictionPct()` methods
- ✅ `dto.EchartConfig` and `dto.JSONAPIEnvelope` DTOs for response formatting
- ✅ Unit test `TestDashboardHandler_Faults` exists
- ✅ Integration test `TestDashboardFaultsChart_Integration` exists
- ✅ Performance benchmark `BenchmarkDashboardFaults` exists

**No blocking dependencies** - all infrastructure is in place.

### 4. Code Patterns

**Route Registration Pattern:**
```go
// Follow existing ECharts endpoint pattern in routes.go
r.HandleFunc("/v1/dashboard/echart/faults.json", dashboardHandler.Faults).Methods("GET")
```

**Positioning:**
- Add route in the "Dashboard endpoints" section
- Place after `/v1/dashboard/projects.json`
- Group with other ECharts endpoints for maintainability

**Testing Pattern:**
```go
// Follow existing route test pattern in routes_test.go
req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/faults.json", nil)
w = httptest.NewRecorder()
handler.ServeHTTP(w, req)
if w.Code != http.StatusOK {
    t.Errorf("Expected 200 for /v1/dashboard/echart/faults.json, got %d", w.Code)
}
```

**Naming Conventions:**
- Route path: snake_case with `.json` extension (consistent with all v1 endpoints)
- Handler method: PascalCase (`Faults`)
- Service methods: PascalCase (`GetFaultsPercentage`, `CreateGaugeChart`)

### 5. Testing Strategy

**Unit Tests (Already Exist):**
- `TestDashboardHandler_Faults` in `internal/api/v1/handlers/dashboard_handler_test.go`
  - Validates handler returns correct JSON:API envelope
  - Tests gauge chart configuration
  - Validates percentage calculation (8 faults → 80%)

**Integration Tests (Already Exist):**
- `TestDashboardFaultsChart_Integration` in `test/dashboard_integration_test.go`
  - Tests with real database using fixture data
  - Validates 16 faults / 10 maxFaults = 160%
  - Uses `ScenarioFaultsChartCorrect()` fixture

**Route Verification Test (To Add):**
- Add to `internal/api/v1/routes_test.go` in `TestSetupRoutes_Routes`
- Verify route responds with 200 OK
- Ensures route is properly registered in router

**Verification Steps:**
1. Run `go test ./internal/api/v1/...` - all unit tests pass
2. Run `go test ./test/...` - integration tests pass
3. Run `go test -bench=BenchmarkDashboardFaults ./test/performance/...` - performance benchmark passes
4. Run `go fmt ./...` and `go vet ./...` - no errors

### 6. Risks and Considerations

**Low Risk Factors:**
- ✅ Minimal change (single line addition to routes.go)
- ✅ All underlying implementation complete and tested
- ✅ Follows established patterns in codebase
- ✅ No database schema changes required
- ✅ No breaking changes to existing APIs

**Considerations:**
1. **Route Ordering**: Add route in logical position with other dashboard endpoints to maintain code readability
2. **Test Coverage**: Existing tests cover handler and integration; only route verification test needs addition
3. **Documentation**: Update AGENTS.md or API documentation to reflect new endpoint (if required by acceptance criteria)

**Edge Cases Handled by Existing Code:**
- Zero faults returns 0% (not NaN) per AC-DASH-004
- Missing config values use defaults (maxFaults = 10)
- Invalid date ranges handled by repository layer
- Context timeout propagation (15-second dashboard context timeout)

**Rollback Plan:**
- If issues arise, simply remove the route registration line from `routes.go`
- No database migrations or configuration changes to revert

**Acceptance Criteria Checklist:**
- [ ] Route `/v1/dashboard/echart/faults.json` registered and accessible
- [ ] Returns JSON:API envelope with gauge chart configuration
- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] `go fmt` and `go vet` pass
- [ ] Clean Architecture layers followed
- [ ] Error responses consistent with existing patterns
- [ ] HTTP status codes correct (200 OK, 500 Internal Server Error)
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
