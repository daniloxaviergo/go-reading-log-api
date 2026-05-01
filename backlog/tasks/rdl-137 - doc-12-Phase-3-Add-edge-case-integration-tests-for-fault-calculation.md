---
id: RDL-137
title: '[doc-12 Phase 3] Add edge case integration tests for fault calculation'
status: To Do
assignee:
  - workflow
created_date: '2026-05-01 14:59'
updated_date: '2026-05-01 15:18'
labels:
  - bugfix
  - testing
  - phase-3
dependencies: []
documentation:
  - doc-012
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create integration tests for edge cases: empty database (all days faults), single log entry (one day with reading, rest faults), and month boundary dates (跨 month ranges). Ensure all edge cases are covered with real database tests.
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task focuses on adding comprehensive edge case integration tests for fault calculation. The implementation will:

1. **Create three new integration test functions** in `test/dashboard_integration_test.go`:
   - `TestDashboardFaults_EmptyDatabase`: Tests when database has no logs (all days in range are faults)
   - `TestDashboardFaults_SingleLogEntry`: Tests when only one day has reading activity (rest are faults)
   - `TestDashboardFaults_MonthBoundary`: Tests date ranges that span across month boundaries

2. **Leverage existing test infrastructure**:
   - Use `test.SetupTestDB()` for test database setup
   - Use `helper.SetupTestSchema()` for table creation
   - Use existing helper functions like `createTestLogs()` pattern
   - Follow the existing test structure from `TestDashboardFaultsChart_Integration`

3. **Test the correct fault definition** (per doc-012 PRD):
   - A "fault" = a day with zero pages read (sum of read_pages = 0 for that day)
   - Days with no logs at all ARE counted as faults
   - Days with logs but sum(end_page - start_page) = 0 ARE counted as faults
   - Days with one or more logs where sum > 0 are NOT counted as faults

4. **Verify SQL query behavior** (when updated per RDL-139):
   - Uses CTE with daily_read aggregation
   - Uses generate_series for all dates in range
   - LEFT JOIN to count zero-page days
   - Inclusive date range (BETWEEN $1 AND $2)

### 2. Files to Modify

| File | Action | Purpose |
|------|--------|---------|
| `test/dashboard_integration_test.go` | Modify | Add 3 new integration test functions for edge cases |

**New Test Functions to Add:**

1. `TestDashboardFaults_EmptyDatabase_Integration`
   - Tests empty database state (no logs at all)
   - Verifies all days in date range are counted as faults
   - Example: 7-day range with no logs = 7 faults

2. `TestDashboardFaults_SingleLogEntry_Integration`
   - Tests single log entry scenario
   - One day has reading, rest of days are faults
   - Example: 7-day range with 1 log on day 3 = 6 faults

3. `TestDashboardFaults_MonthBoundary_Integration`
   - Tests date ranges spanning month boundaries
   - Handles edge cases like Jan 28 - Feb 5
   - Verifies correct fault counting across month transitions

### 3. Dependencies

**Prerequisites:**
- RDL-139: Update GetFaultsByDateRange SQL query with CTE for daily aggregation (must be completed first)
- RDL-137 depends on the corrected SQL query implementation
- Test database must be configured (`reading_log_test`)

**Related Tasks:**
- RDL-139: Update GetFaultsByDateRange SQL query (Phase 1)
- RDL-138: Update mock repository expectations (Phase 2)
- RDL-137: Add edge case integration tests (Phase 3) - **THIS TASK**
- RDL-140: Update GetWeekdayFaults SQL query (Phase 1)
- RDL-141: Verify context timeout and error handling (Phase 1)

**Test Infrastructure Requirements:**
- PostgreSQL test database (`reading_log_test`)
- Test schema created via `helper.SetupTestSchema()`
- Fixture data for scenarios (can use inline log creation)

### 4. Code Patterns

**Follow Existing Test Patterns:**

```go
// Pattern from TestDashboardFaultsChart_Integration
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

    // Load test data
    fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
    scenario := dashboardFixtures.ScenarioFaultsByWeekday()
    err = fixtureManager.LoadScenario(scenario)
    require.NoError(t, err)

    // Create handler with real dependencies
    repo, err := createTestRepository(helper.Pool)
    require.NoError(t, err)

    // Execute test
    req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/faults.json", nil)
    recorder := httptest.NewRecorder()

    dashboardHandler.Faults(recorder, req)

    // Verify
    assert.Equal(t, http.StatusOK, recorder.Code)
    // Additional assertions...
}
```

**Test Data Creation Pattern:**
```go
// Use inline log creation for edge cases
testDateStr := "2024-01-15"
ctx := context.Background()
query := `INSERT INTO projects (id, name, total_page, page) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING`
_, err = helper.Pool.Exec(ctx, query, 1, "Test Project", 100, 0)
require.NoError(t, err)

err = createTestLogs(helper.Pool, []testLog{
    {ProjectID: 1, Data: &testDateStr, StartPage: 0, EndPage: 10},
})
require.NoError(t, err)
```

**Naming Conventions:**
- Test function names: `TestDashboardFaults_<EdgeCase>_Integration`
- Use descriptive scenario names in comments
- Follow Go testing conventions (t.Run for subtests if needed)

### 5. Testing Strategy

**Test Coverage:**

1. **Empty Database Test** (`TestDashboardFaults_EmptyDatabase_Integration`)
   - Setup: Clear all logs, create empty database state
   - Date Range: 7 days (e.g., 2024-01-15 to 2024-01-21)
   - Expected: 7 faults (all days are faults)
   - Edge Cases:
     - Single day range (1 fault)
     - 30-day range (30 faults)
     - Verify handler returns 200 OK with correct fault count

2. **Single Log Entry Test** (`TestDashboardFaults_SingleLogEntry_Integration`)
   - Setup: Create 1 log entry on day 3 of a 7-day range
   - Date Range: 7 days (e.g., 2024-01-15 to 2024-01-21)
   - Expected: 6 faults (only day 3 has reading)
   - Edge Cases:
     - Log with 0 pages (end_page = start_page) = fault
     - Log with > 0 pages = not a fault
     - Multiple logs on same day = 1 day with reading

3. **Month Boundary Test** (`TestDashboardFaults_MonthBoundary_Integration`)
   - Setup: Create logs spanning month boundary
   - Date Range: Jan 28, 2024 to Feb 5, 2024 (9 days)
   - Expected: Correct fault count across month transition
   - Edge Cases:
     - End of month (Jan 31)
     - Start of new month (Feb 1)
     - Leap year February (Feb 29)
     - Verify date casting works correctly

**Test Execution:**
```bash
# Run specific edge case tests
go test -v ./test/... -run TestDashboardFaults_EmptyDatabase_Integration
go test -v ./test/... -run TestDashboardFaults_SingleLogEntry_Integration
go test -v ./test/... -run TestDashboardFaults_MonthBoundary_Integration

# Run all fault-related tests
go test -v ./test/... -run TestDashboardFaults
```

**Edge Cases to Cover:**

| Edge Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| Empty database | No logs exist | All days in range = faults |
| Single log entry | One day has reading | Rest of days = faults |
| Zero pages log | end_page = start_page | Counts as fault |
| Month boundary | Jan 31 → Feb 1 | Correct date handling |
| Leap year | Feb 29 in leap year | Correct date handling |
| Single day range | start_date = end_date | 1 day evaluated |
| NULL values | start_page or end_page NULL | Treated as 0 pages |

**Validation Approach:**
- Verify HTTP status code (200 OK)
- Parse JSON response
- Assert fault count matches expected
- Verify response structure matches DTO schema
- Check error handling for edge cases

### 6. Risks and Considerations

**Known Issues:**
1. **Current SQL query is incorrect**: The existing `GetFaultsByDateRange` counts all logs as faults instead of counting days with zero pages. This task assumes RDL-139 will be completed first to fix the SQL query.
   - **Mitigation**: Tests will be written for the correct behavior; if current implementation is tested, they will fail (which is expected)

2. **Test database cleanup**: Orphaned test databases can cause test failures
   - **Mitigation**: Use `defer helper.Close()` and ensure proper cleanup

3. **Date handling timezone issues**: Server timezone must match expected test dates
   - **Mitigation**: Use explicit timezone in test data (UTC) and verify server TZ setting

**Potential Pitfalls:**
1. **generate_series performance**: For very large date ranges, generate_series can be slow
   - **Consideration**: Test with reasonable ranges (7-30 days) for edge cases

2. **Date casting from VARCHAR**: The `data` column is VARCHAR, casting to date requires proper handling
   - **Consideration**: Verify SQL query uses `data::date` correctly

3. **Inclusive vs exclusive ranges**: BETWEEN is inclusive on both ends
   - **Consideration**: Ensure test expectations account for inclusive range

**Deployment Considerations:**
- Tests should pass in CI/CD pipeline
- No breaking changes to existing API contracts
- Test data cleanup required between test runs

**Acceptance Criteria Verification:**
- ✅ All unit tests pass (`go test ./...`)
- ✅ All integration tests pass (`go test -v ./test/...`)
- ✅ `go fmt` and `go vet` pass with no errors
- ✅ Clean Architecture layers properly followed
- ✅ Error responses consistent with existing patterns
- ✅ HTTP status codes correct for response type
- ✅ Documentation updated in QWEN.md (if needed)
- ✅ New code paths include error path tests
- ✅ HTTP handlers test both success and error responses
- ✅ Integration tests verify actual database interactions

**Test Data Examples:**

```go
// Empty Database Test
// Date Range: 2024-01-15 to 2024-01-21 (7 days)
// Expected Faults: 7

// Single Log Entry Test
// Date Range: 2024-01-15 to 2024-01-21 (7 days)
// Log on: 2024-01-17 (day 3) with 10 pages
// Expected Faults: 6

// Month Boundary Test
// Date Range: 2024-01-28 to 2024-02-05 (9 days)
// Logs on: Jan 30 (5 pages), Feb 2 (10 pages)
// Expected Faults: 7 (days without reading)
```
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
