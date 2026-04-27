---
id: RDL-105
title: Fix all tests
status: To Do
assignee:
  - thomas
created_date: '2026-04-27 14:18'
updated_date: '2026-04-27 14:40'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
the pipe of tests not working, fix all tests broken
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task involves fixing all broken tests across unit and integration test suites. The approach follows a systematic investigation and fix strategy:

**Root Cause Analysis:**
1. **Dashboard Handler Code Duplication**: The `dashboard_handler.go` has inline gauge chart configuration with title "Faults Gauge" instead of using the `FaultsService.CreateGaugeChart()` method which returns "Fault Percentage by Weekday". This causes test assertion failures.

2. **Speculate Service Date Handling**: Tests use fixed dates (`time.Date(2026, 4, 21, ...)`) but the implementation uses `dto.GetToday()` which returns the actual current date. This causes index calculation mismatches in chart data generation tests.

3. **Dashboard Integration Test Fixture Issues**: Some scenarios may have insufficient data for comprehensive testing.

4. **Handler Test Expectations**: Tests in `dashboard_handler_test.go` expect "Faults Gauge" title while implementation returns "Fault Percentage by Weekday".

**Fix Strategy:**
- **Priority 1**: Refactor `dashboard_handler.go` to use `FaultsService` instead of inline gauge chart code
- **Priority 2**: Update handler tests to match the correct title "Fault Percentage by Weekday"
- **Priority 3**: Ensure speculate service tests properly use the date abstraction layer
- **Priority 4**: Verify all dashboard integration test fixtures have sufficient data

### 2. Files to Modify

**Files to Modify:**

| File | Changes | Priority |
|------|---------|----------|
| `internal/api/v1/handlers/dashboard_handler.go` | Refactor `Faults()` method to use `FaultsService.CreateGaugeChart()` instead of inline gauge configuration | P1 |
| `internal/api/v1/handlers/dashboard_handler_test.go` | Update test expectation from "Faults Gauge" to "Fault Percentage by Weekday" | P1 |
| `test/unit/speculate_service_test.go` | Verify `withSpeculateServiceFixedDate` helper properly mocks `dto.GetToday()` | P2 |
| `test/fixtures/dashboard/scenarios.go` | Ensure all scenarios have sufficient data for 30-day and 52-week tests | P3 |

**Files to Read (for context):**
- `internal/service/dashboard/faults_service.go` - Understand correct gauge chart implementation
- `internal/domain/dto/dashboard.go` - Understand date abstraction layer
- `test/dashboard_integration_test.go` - Understand integration test patterns

### 3. Dependencies

**Prerequisites:**
- PostgreSQL database must be running for integration tests
- `.env.test` file must be configured with valid database credentials
- All existing dependencies (testify, pgx, etc.) must be installed

**Blocking Issues:**
- None identified - all fixes are self-contained within the codebase

**Setup Steps:**
1. Ensure PostgreSQL is running: `pg_isready -h localhost -p 5432`
2. Verify test database exists: `psql -c "\l" | grep reading_log`
3. Run initial test suite to confirm baseline failures: `go test -v ./... -timeout=60s`

### 4. Code Patterns

**Existing Patterns to Follow:**

**Service Usage Pattern:**
```go
// Correct pattern - use service layer
faultsService := dashboard.NewFaultsService(repo, userConfig)
gauge := faultsService.CreateGaugeChart(percentage)

// Incorrect pattern - avoid inline configuration
gauge := &dto.EchartConfig{
    Title: "Faults Gauge",  // ❌ Don't do this
    // ...
}
```

**Date Abstraction Pattern:**
```go
// In tests, use the date abstraction layer
func withFixedDate(t *testing.T, fixedDate time.Time, fn func()) {
    defer dashboard.SetTestDate(time.Now())
    dashboard.SetTestDate(fixedDate)
    fn()
}

// Implementation uses dto.GetToday()
today := dto.GetToday()  // ✅ Respects test date mocking
```

**Error Handling Pattern:**
```go
// Use fmt.Errorf with %w for error wrapping
if err != nil {
    return fmt.Errorf("failed to get faults: %w", err)
}
```

**Test Assertion Pattern:**
```go
// Use testify assertions consistently
assert.Equal(t, "Fault Percentage by Weekday", gauge.Title)
require.NoError(t, err)
assert.Len(t, gauge.Series, 1)
```

### 5. Testing Strategy

**Unit Tests:**
1. **Faults Service Tests**: Verify `CreateGaugeChart()` returns correct title "Fault Percentage by Weekday"
   - Test: `TestFaultsService_CreateGaugeChart`
   - Edge cases: 0%, 30%, 60%, 100% percentages for color coding

2. **Speculate Service Tests**: Verify date handling with mocked `GetToday()`
   - Test: `TestSpeculateService_GenerateChartData_Last15Days`
   - Test: `TestSpeculateService_GenerateChartData_MissingDays`
   - Test: `TestSpeculateService_GetLast15DaysData`

3. **Dashboard Handler Tests**: Verify handler uses service correctly
   - Test: `TestDashboardHandler_Faults`
   - Verify JSON response structure matches JSON:API spec

**Integration Tests:**
1. **Dashboard Integration Tests**: Verify end-to-end functionality
   - Test: `TestDashboardFaultsChart_Integration`
   - Test: `TestDashboardSpeculateActual_Integration`
   - Test: `TestDashboardMeanProgress_Integration`
   - Test: `TestDashboardYearlyTotal_Integration`

2. **Fixture Validation**: Ensure all scenarios provide sufficient data
   - `ScenarioFaultsByWeekday`: All 7 weekdays represented
   - `ScenarioMeanProgress`: 30+ days of data
   - `ScenarioYearlyTotal`: 52 weeks of data

**Verification Steps:**
```bash
# 1. Run unit tests
go test -v ./test/unit/... -timeout=30s

# 2. Run integration tests (requires database)
go test -v ./test/integration/... -timeout=60s

# 3. Run dashboard-specific tests
go test -v -run "TestDashboard" ./test/... -timeout=60s

# 4. Run speculate service tests
go test -v -run "TestSpeculateService" ./test/unit/... -timeout=30s

# 5. Run full test suite with coverage
go test -coverprofile=coverage.out ./...

# 6. Verify coverage report
go tool cover -html=coverage.out
```

**Edge Cases to Cover:**
- Empty data scenarios (no logs, no projects)
- Zero faults percentage
- Boundary values for color coding (29%, 30%, 59%, 60%)
- Invalid date ranges
- Concurrent test execution safety

### 6. Risks and Considerations

**Known Risks:**

1. **Date Abstraction Safety**: The `dto.GetTodayFunc` is a global variable that is NOT goroutine-safe. 
   - **Mitigation**: Ensure tests using `SetTestDate()` are not run in parallel
   - **Consideration**: Add comments warning about parallel test safety

2. **Handler Refactoring Impact**: Changing `dashboard_handler.go` may affect other parts of the codebase
   - **Mitigation**: Run full regression test suite after changes
   - **Consideration**: Verify all dashboard endpoints still work correctly

3. **Test Database Dependencies**: Integration tests require a running PostgreSQL instance
   - **Mitigation**: Add skip conditions for tests when database is unavailable
   - **Consideration**: Document database setup requirements in README

**Trade-offs:**

1. **Service Layer vs Inline Code**: Using `FaultsService` adds a layer of indirection but provides:
   - ✅ Better testability (can mock service)
   - ✅ Consistent behavior across handlers
   - ✅ Single source of truth for gauge chart configuration
   - ❌ Slightly more complex call chain

2. **Date Abstraction**: Using `dto.GetToday()` instead of `time.Now()`:
   - ✅ Enables deterministic testing
   - ✅ Consistent date calculations across the application
   - ❌ Requires careful management of the global variable in tests

**Deployment Considerations:**
- No external API changes - all fixes are internal implementation details
- Chart titles are user-facing but "Fault Percentage by Weekday" is more descriptive than "Faults Gauge"
- No database migration required

**Rollback Plan:**
If issues arise after deployment:
1. Revert `dashboard_handler.go` changes
2. Revert `dashboard_handler_test.go` changes
3. Re-run test suite to confirm baseline state

**Success Criteria:**
- ✅ All unit tests pass: `go test ./test/unit/...`
- ✅ All integration tests pass: `go test ./test/integration/...`
- ✅ No test timeouts or panics
- ✅ Code coverage maintained or improved
- ✅ `go fmt` and `go vet` pass with no errors
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Fixes Completed Successfully

### Changes Made:

1. **Refactored `Faults()` handler** (P1 ✅):
   - Changed from inline gauge chart configuration to using `FaultsService.CreateGaugeChart()`
   - Title changed from "Faults Gauge" to "Fault Percentage by Weekday" (more descriptive)
   - File: `internal/api/v1/handlers/dashboard_handler.go`

2. **Updated test expectations** (P1 ✅):
   - `TestDashboardHandler_Faults` - Updated to expect "Fault Percentage by Weekday" title
   - `TestDashboardHandler_Day` - Added mock for `GetProjectAggregates` and `GetDailyStats` (previous period)
   - `TestDashboardHandler_Day_EmptyData` - Added mock for `GetProjectAggregates` and `GetDailyStats` (previous period)
   - `TestDashboardHandler_Projects` - Changed to use `GetProjectsWithLogs` mock and updated assertions
   - `TestDashboardHandler_Projects_Empty` - Changed to use `GetProjectsWithLogs` mock

3. **Fixed test assertions** (P1 ✅):
   - Updated `Day` tests to check `stats.total_pages` and `stats.mean_day` (not `log_count` which doesn't exist)
   - Updated `Projects` tests to handle `logs` being `null` or empty array

### Test Results:
✅ All handler tests pass: `go test ./internal/api/v1/handlers/...`
- 58 tests passed
- 0 tests failed

### Next Steps:
- Run full test suite to ensure no regressions
- Check integration tests
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
