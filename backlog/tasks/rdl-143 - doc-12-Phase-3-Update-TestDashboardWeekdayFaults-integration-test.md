---
id: RDL-143
title: '[doc-12 Phase 3] Update TestDashboardWeekdayFaults integration test'
status: To Do
assignee:
  - thomas
created_date: '2026-05-01 15:08'
updated_date: '2026-05-01 17:38'
labels:
  - bugfix
  - testing
  - phase-3
dependencies: []
documentation:
  - doc-012
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Fix TestDashboardWeekdayFaults integration test by updating test data setup to create reading patterns with specific weekday gaps. Update expected weekday distribution to match Rails logic. Verify radar chart data reflects correct fault counts per weekday.
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task involves updating the `TestDashboardWeekdayFaults_Integration` integration test to correctly validate the fault calculation after the SQL query fixes implemented in RDL-139 and RDL-140. The integration test currently uses `ScenarioFaultsByWeekday()` fixture data which was designed for the old (incorrect) logic that counted log entries as faults.

**Technical Strategy:**
- Update the integration test to use a fixture scenario that creates specific weekday gaps (days WITHOUT reading)
- The test must validate that the radar chart data reflects correct fault counts per weekday
- Follow the same pattern as the unit tests in RDL-142, but at the integration level with HTTP handler verification
- The fixture must create a date range where specific weekdays have no reading activity, and those days should be counted as faults

**Why this approach:**
- Aligns with the CTE-based SQL query implementation in `dashboard_repository.go`
- Matches the Rails `V1::Dashboard::WeekdayFaults` calculation logic
- Ensures the integration test validates the complete stack: HTTP handler → service → repository → database
- Follows the pattern established in RDL-142 for unit tests

**Key Understanding:**
- A "fault" = a day where `sum(read_pages) = 0` (either no logs or logs with end_page - start_page = 0)
- The SQL query uses `generate_series` to create all dates in the range, then LEFT JOINs to count days without reading
- The last 6 months from "today" is the default date range for the weekday faults endpoint

### 2. Files to Modify

| File | Change Type | Description |
|------|-------------|-------------|
| `test/dashboard_integration_test.go` | Modify | Update `TestDashboardWeekdayFaults_Integration` test - update fixture usage and expected weekday distribution |
| `test/fixtures/dashboard/scenarios.go` | Modify | Update `ScenarioFaultsByWeekday()` to create proper fault scenarios (days WITHOUT reading) instead of log entries |

**Files to Review (no changes expected):**
- `internal/adapter/postgres/dashboard_repository.go` - Repository implementation (already fixed in RDL-140)
- `internal/api/v1/handlers/dashboard_handler.go` - Handler implementation (no changes needed)
- `internal/domain/dto/dashboard_response.go` - DTO definitions (no changes needed)

### 3. Dependencies

| Dependency | Status | Notes |
|------------|--------|-------|
| RDL-139 (GetFaultsByDateRange SQL fix) | ✅ Done | Prerequisite - SQL query fix completed |
| RDL-140 (GetWeekdayFaults SQL fix) | ✅ Done | Prerequisite - SQL query fix completed |
| RDL-141 (Context timeout verification) | ✅ Done | Prerequisite - Context handling verified |
| RDL-142 (Unit tests rewrite) | ✅ Done | Prerequisite - Unit tests completed |
| Test database setup | ✅ Available | Uses `test.SetupTestDB()` helper |

**No blocking dependencies** - All prerequisite tasks are complete.

### 4. Code Patterns

Follow existing patterns in the codebase:

**Integration Test Structure:**
```go
func TestDashboardWeekdayFaults_Integration(t *testing.T) {
    if !IsTestDatabase() {
        t.Skip("Test database not configured")
    }

    helper, err := SetupTestDB()
    require.NoError(t, err)
    defer helper.Close()

    // Create database tables
    err = helper.SetupTestSchema()
    require.NoError(t, err)

    // Load test data using fixture
    fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
    scenario := dashboardFixtures.ScenarioFaultsByWeekday()
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

    // Test GET /v1/dashboard/echart/faults_week_day.json
    req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/faults_week_day.json", nil)
    recorder := httptest.NewRecorder()

    dashboardHandler.WeekdayFaults(recorder, req)

    assert.Equal(t, http.StatusOK, recorder.Code)

    // Parse and validate response
    response, err := parseDashboardResponse(recorder.Body.Bytes())
    require.NoError(t, err)

    // Verify radar chart structure and data
    assert.NotNil(t, response.Echart)
    series := response.Echart.Series[0]
    assert.Equal(t, "radar", series.Type)
    assert.Len(t, series.Data, 7) // 7 weekdays

    // Validate specific weekday fault counts
    actualFaults := make([]int, len(series.Data))
    for i, v := range series.Data {
        actualFaults[i] = int(v.(float64))
    }
    assert.Equal(t, expectedFaults, actualFaults)
}
```

**Fixture Pattern (ScenarioFaultsByWeekday):**
```go
func ScenarioFaultsByWeekday() *Scenario {
    // Define which weekdays should have faults (no reading)
    // Example: Create a 6-month period where specific weekdays have no reading
    
    // For each weekday that should have faults:
    // - Generate dates that fall on that weekday within the 6-month range
    // - DO NOT create logs for those dates (they become fault days)
    
    // For weekdays that should NOT have faults:
    // - Create logs on those dates to ensure reading activity
    
    // The key insight: faults = days WITHOUT logs, not log entries
}
```

**PostgreSQL Weekday Mapping (DOW):**
- 0 = Sunday
- 1 = Monday
- 2 = Tuesday
- 3 = Wednesday
- 4 = Thursday
- 5 = Friday
- 6 = Saturday

### 5. Testing Strategy

**Test Coverage Requirements:**

#### Update ScenarioFaultsByWeekday Fixture
The current `ScenarioFaultsByWeekday()` creates log entries that were incorrectly interpreted as faults. It needs to be rewritten to:

1. **Create a 6-month date range** (matching the endpoint's default range)
2. **Define specific weekdays that should have faults** (no reading activity)
3. **Create logs ONLY on non-fault days** to ensure those days are NOT counted as faults
4. **Calculate expected fault distribution** based on the number of days per weekday in the range minus the days with reading

**Example Scenario:**
- Date range: 6 months from today (e.g., Oct 2025 to Apr 2026)
- Fault weekdays: Monday (1) and Wednesday (3) - no reading on these days
- Non-fault weekdays: Sunday (0), Tuesday (2), Thursday (4), Friday (5), Saturday (6) - create logs on these days
- Expected: Count of Mondays and Wednesdays in the 6-month range = faults for those weekdays

#### Update TestDashboardWeekdayFaults_Integration
The integration test needs to:

1. **Use the updated fixture** - Call `ScenarioFaultsByWeekday()` which now creates proper fault scenarios
2. **Calculate expected weekday distribution** - Based on the fixture's design
3. **Verify HTTP response** - Check status code 200 OK
4. **Parse radar chart data** - Extract fault counts for each weekday
5. **Validate specific values** - Compare expected vs actual fault counts per weekday

**Edge Cases to Cover:**
- Empty database (all days are faults)
- All days have reading (0 faults for all weekdays)
- Single weekday has all faults
- Irregular distribution across weekdays

**Validation Approach:**
1. Use `require.NoError(t, err)` for error handling
2. Use `assert.NotNil(t, response)` for result validation
3. Use `assert.Equal(t, expected, actual)` for exact value matching
4. Verify all 7 weekdays present in radar chart data (0-6)
5. Verify fault counts match the expected distribution

**Integration Test Checklist:**
- [ ] HTTP handler returns 200 OK
- [ ] Response contains valid JSON:API envelope
- [ ] Echart configuration present with radar chart type
- [ ] 7 data points in series (one per weekday)
- [ ] Fault counts match expected distribution
- [ ] All values are non-negative integers

### 6. Risks and Considerations

**Known Risks:**
1. **Date Range Calculation Complexity**: The endpoint uses "last 6 months" from "today", which is dynamic. Tests must either:
   - Use a fixed "today" date via `dto.SetTestDate()`
   - Calculate expected values based on the actual current date
   - **Mitigation**: Use `dto.SetTestDate()` to fix the reference date for predictable test results

2. **Fixture Data Complexity**: Creating a scenario with specific weekday faults requires careful date calculation
   - **Mitigation**: Use Go's `time.Weekday()` to verify test date weekdays
   - **Mitigation**: Add inline comments explaining date calculations and expected distributions
   - **Mitigation**: Create a helper function to calculate expected weekday distribution for a given date range

3. **Test Data Isolation**: Integration tests use the test database
   - **Mitigation**: Use `defer helper.Close()` for cleanup
   - **Mitigation**: Each test creates its own isolated database state

4. **Rails Parity Validation**: Hard to verify exact Rails behavior without running Rails app
   - **Mitigation**: Follow the SQL logic documented in doc-012
   - **Mitigation**: Use the same CTE-based query pattern as the repository

**Considerations:**
- **Dynamic Date Range**: The endpoint calculates "last 6 months" dynamically - must fix the reference date for tests
- **PostgreSQL DOW**: Ensure correct mapping (0=Sunday, not 0=Monday)
- **Inclusive Date Range**: BETWEEN $1 AND $2 includes both endpoints
- **Weekday Distribution**: Must account for partial months at range boundaries
- **Null Handling**: All 7 weekdays should be present in the response (default 0 if no faults)

**Rollback Plan:**
- If tests fail unexpectedly, revert to previous test version
- Investigate by running individual test cases with verbose output (`go test -v`)
- Compare expected vs actual weekday distributions
- Check if the fixture data is being loaded correctly into the database

**Success Criteria:**
- `TestDashboardWeekdayFaults_Integration` passes with correct expected values
- Fixture `ScenarioFaultsByWeekday()` creates predictable fault distribution
- All existing integration tests continue to pass (no regressions)
- `go fmt` and `go vet` pass with no errors
- Test documentation updated with explanation of the fault calculation logic

**Definition of Done:**
- [ ] All unit tests pass
- [ ] All integration tests pass execution and verification
- [ ] `go fmt` and `go vet` pass with no errors
- [ ] Clean Architecture layers properly followed
- [ ] Error responses consistent with existing patterns
- [ ] HTTP status codes correct for response type
- [ ] Documentation updated in QWEN.md
- [ ] New code paths include error path tests
- [ ] HTTP handlers test both success and error responses
- [ ] Integration tests verify actual database interactions
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress - RDL-143

### Understanding the Task
- **Goal**: Fix `TestDashboardWeekdayFaults_Integration` by updating test data setup
- **Problem**: Current `ScenarioFaultsByWeekday()` creates log entries that are incorrectly interpreted as faults
- **Solution**: Rewrite fixture to create proper fault scenarios where specific weekdays have NO reading activity

### Key Understanding
- A "fault" = a day where `sum(read_pages) = 0` (no logs or logs with zero pages)
- The SQL query uses CTE with `generate_series` to create all dates, then LEFT JOINs to find days without reading
- Weekday mapping: 0=Sunday, 1=Monday, ..., 6=Saturday

### Implementation Steps
1. ✅ Review task details and acceptance criteria
2. ✅ Examine current test implementation and fixture
3. ✅ Understand SQL query logic in `GetWeekdayFaults`
4. ✅ Update `ScenarioFaultsByWeekday()` fixture to create proper fault scenarios
   - Created 6-month date range (Oct 1, 2025 to Apr 1, 2026)
   - Logs on Sun/Tue/Thu/Sat (no faults)
   - NO logs on Mon/Wed/Fri (these become faults)
5. ✅ Update integration test with expected values
   - Set fixed "today" date for predictable results
   - Calculate expected faults based on fixture design
   - Validate each weekday's fault count
6. 🔄 Run tests and verify
7. ⏳ Check acceptance criteria
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
