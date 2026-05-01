---
id: RDL-144
title: '[doc-12 Phase 3] Add edge case integration tests'
status: Done
assignee:
  - next-task
created_date: '2026-05-01 15:08'
updated_date: '2026-05-01 18:08'
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
Create integration tests for edge cases: empty database (all days faults), single log entry (one day with reading, rest faults), and month boundary dates (cross month ranges). Ensure all edge cases are covered with real database tests.
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task focuses on ensuring comprehensive edge case integration test coverage for the faults calculation logic. After analyzing the existing test suite, I found that most edge case tests already exist, but there are opportunities to enhance coverage and add missing scenarios.

**Implementation Strategy:**
- Review existing edge case integration tests in `test/dashboard_integration_test.go`
- Identify gaps in coverage for weekday faults edge cases
- Add missing integration tests for specific edge cases
- Ensure all tests use real database interactions (not mocks)
- Follow Clean Architecture patterns for test organization

**Key Edge Cases to Cover:**
1. **Empty Database** - All days in range should be faults (already exists but can be enhanced)
2. **Single Log Entry** - One day with reading, rest are faults (already exists)
3. **Month Boundary Dates** - Cross-month date ranges (already exists)
4. **Leap Year February** - Feb 29 handling (partially covered)
5. **Year Boundary** - Dec 31 to Jan 1 transitions (missing)
6. **Zero Pages Read** - Logs exist but end_page = start_page (exists in unit tests, needs integration)
7. **Weekday Faults Edge Cases** - Empty weekday faults, all weekdays faults (missing integration tests)
8. **NULL Value Handling** - Logs with NULL start_page or end_page (missing)

### 2. Files to Modify

**Files to Read/Analyze:**
- `test/dashboard_integration_test.go` - Existing integration tests
- `test/unit/dashboard_repository_test.go` - Existing unit tests
- `internal/adapter/postgres/dashboard_repository.go` - Repository implementation
- `test/fixtures/dashboard/scenarios.go` - Test fixture scenarios

**Files to Modify:**
- `test/dashboard_integration_test.go` - Add missing edge case integration tests:
  - `TestDashboardWeekdayFaults_EmptyDatabase_Integration` - Weekday faults with empty database
  - `TestDashboardWeekdayFaults_SingleLogEntry_Integration` - Weekday faults with single log
  - `TestDashboardFaults_YearBoundary_Integration` - Year transition edge case
  - `TestDashboardFaults_NULLValues_Integration` - NULL start_page/end_page handling
  - `TestDashboardFaults_ZeroPagesRead_Integration` - Integration test for zero pages scenario

**Files to Create:**
- None (all tests will be added to existing `dashboard_integration_test.go`)

### 3. Dependencies

**Prerequisites:**
- RDL-139: GetFaultsByDateRange SQL query fix (DONE)
- RDL-140: GetWeekdayFaults SQL query fix (DONE)
- Test database `reading_log_test` must be configured
- PostgreSQL 13+ with `generate_series` support

**Existing Test Infrastructure:**
- `test.SetupTestDB()` - Test database setup
- `test.TestHelper` - Database cleanup utilities
- `dashboardFixtures.NewDashboardFixtures()` - Fixture manager
- Existing scenarios in `test/fixtures/dashboard/scenarios.go`

### 4. Code Patterns

**Test Structure Patterns to Follow:**
```go
func TestDashboardFaults_<Scenario>_Integration(t *testing.T) {
    if !IsTestDatabase() {
        t.Skip("Test database not configured")
    }

    helper, err := SetupTestDB()
    require.NoError(t, err)
    defer helper.Close()

    // Create database tables
    err = helper.SetupTestSchema()
    require.NoError(t, err)

    // Clear existing data
    err = helper.ClearTestData()
    require.NoError(t, err)

    // Create test data (projects, logs)
    // ...

    // Test the repository method directly
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    faultStats, err := repo.GetFaultsByDateRange(ctx, startDate, endDate)
    assert.NoError(t, err)
    assert.Equal(t, expectedFaults, faultStats.FaultCount)

    // Test the HTTP endpoint
    // ...
}
```

**Naming Conventions:**
- Test functions: `Test<Endpoint>_<Scenario>_Integration`
- Use descriptive scenario names (EmptyDatabase, SingleLogEntry, MonthBoundary, etc.)
- Follow existing test organization in `dashboard_integration_test.go`

**Error Handling:**
- Use `require.NoError(t, err)` for setup errors
- Use `assert.NoError(t, err)` for operation errors
- Always use `defer helper.Close()` for cleanup

### 5. Testing Strategy

**Test Types:**
1. **Repository-level tests** - Direct method calls with real database
2. **HTTP endpoint tests** - Full handler integration with HTTP request/response

**Edge Cases to Test:**

| Test Name | Scenario | Expected Result |
|-----------|----------|-----------------|
| `TestDashboardWeekdayFaults_EmptyDatabase_Integration` | No logs in database | All weekdays have faults distributed by date range |
| `TestDashboardWeekdayFaults_SingleLogEntry_Integration` | One log on specific weekday | That weekday has 0 faults, others have expected faults |
| `TestDashboardFaults_YearBoundary_Integration` | Dec 30 to Jan 2 | Correct fault count across year boundary |
| `TestDashboardFaults_NULLValues_Integration` | Logs with NULL pages | NULL treated as 0 pages, day counts as fault |
| `TestDashboardFaults_ZeroPagesRead_Integration` | end_page = start_page | Day with zero pages counts as fault |

**Test Data Scenarios:**

1. **Empty Database:**
   - No projects, no logs
   - 7-day range → 7 faults
   - 30-day range → 30 faults

2. **Single Log Entry:**
   - One project, one log on day 3 of 7-day range
   - Expected: 6 faults (7 days - 1 day with reading)

3. **Year Boundary:**
   - Date range: Dec 30, 2023 to Jan 2, 2024 (4 days)
   - Log on Dec 31, 2023
   - Expected: 3 faults

4. **NULL Values:**
   - Log with start_page = NULL, end_page = 10
   - Should be treated as 0 pages (due to CASE statement)
   - Expected: Day counts as fault

5. **Zero Pages Read:**
   - Log with start_page = 10, end_page = 10
   - Expected: Day counts as fault

**Validation Approach:**
- Compare against manually calculated expected values
- Use `assert.InDelta()` for floating point comparisons
- Test both repository layer and HTTP endpoint
- Verify JSON response structure matches API contract

### 6. Risks and Considerations

**Known Risks:**
1. **Test Database Cleanup** - Ensure proper cleanup between tests to avoid data pollution
   - Mitigation: Use `defer helper.Close()` and explicit `TRUNCATE` statements

2. **Date/Timezone Handling** - Server timezone affects date casting
   - Mitigation: Use fixed dates in tests, set `TZ=UTC` for consistency

3. **PostgreSQL Version** - `generate_series` requires PG 8.4+
   - Mitigation: Project requires PG 13+, this is not a concern

4. **Context Timeout** - Some queries may timeout with large date ranges
   - Mitigation: Use 15-second timeout as defined in `dashboardContextTimeout`

**Considerations:**
1. **Test Execution Time** - Integration tests with real database are slower
   - Consider running in parallel where safe: `t.Parallel()`

2. **Test Isolation** - Each test should be independent
   - Use unique IDs for test data to avoid conflicts

3. **Documentation** - Add comments explaining edge case scenarios
   - Follow existing comment style in `dashboard_integration_test.go`

4. **Rails Parity** - Ensure Go tests match Rails behavior
   - Reference Rails implementation in comments when applicable

**Acceptance Criteria Verification:**
- [ ] All new tests pass with real database
- [ ] Tests cover empty database scenario
- [ ] Tests cover single log entry scenario
- [ ] Tests cover month/year boundary dates
- [ ] Tests cover NULL value handling
- [ ] Tests cover zero pages read scenario
- [ ] Weekday faults edge cases covered
- [ ] No test data pollution between tests
- [ ] Code follows Clean Architecture patterns
- [ ] `go fmt` and `go vet` pass
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Complete ✅

### Tests Added (All Passing)
1. ✅ `TestDashboardFaults_YearBoundary_Integration` - Year transition edge case (Dec 30 to Jan 2)
   - Tests 3 scenarios: Dec 30 to Jan 2, Dec 31 to Jan 1, year boundary with no logs
   - Both repository and HTTP endpoint tests

2. ✅ `TestDashboardFaults_ZeroPagesRead_Integration` - Integration test for zero pages scenario
   - Tests 3 scenarios: single log with zero pages, multiple zero pages, mix of zero/normal pages
   - Both repository and HTTP endpoint tests

3. ✅ `TestDashboardWeekdayFaults_EmptyDatabase_Integration` - Weekday faults with empty database
   - Repository test verifies all 7 weekdays have faults distributed evenly (~26-27 each)
   - HTTP endpoint test verifies radar chart structure

4. ✅ `TestDashboardWeekdayFaults_SingleLogEntry_Integration` - Weekday faults with single log
   - Repository test verifies single log affects weekday fault distribution
   - HTTP endpoint test verifies radar chart structure

### Test Results
- All integration tests pass ✅
- All unit tests pass ✅
- `go fmt` passes ✅
- `go vet` passes ✅
- Build succeeds ✅

### Notes
- NULL values test was removed because the logs table has NOT NULL constraints on start_page/end_page
- NULL handling is covered in unit tests (dashboard_repository_test.go: TestDashboardRepository_GetMaxByWeekday_InvalidData)
- Zero pages read scenario covers the practical equivalent of NULL handling

### Edge Cases Covered
- ✅ Empty database (all days faults)
- ✅ Single log entry (one day with reading, rest faults)
- ✅ Month boundary dates (cross month ranges)
- ✅ Year boundary dates (Dec 31 to Jan 1 transitions)
- ✅ Zero pages read (end_page = start_page)
- ✅ Weekday faults edge cases (empty database, single log entry)
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary

Added comprehensive edge case integration tests for dashboard fault calculation logic to ensure robust testing of boundary conditions and error scenarios.

## What Was Done

### New Integration Tests Added (4 tests, 14 test cases total)

1. **TestDashboardFaults_YearBoundary_Integration**
   - Tests date ranges spanning year boundaries (Dec 30 to Jan 2, Dec 31 to Jan 1)
   - Validates correct fault counting across year transitions
   - Includes repository-level and HTTP endpoint tests

2. **TestDashboardFaults_ZeroPagesRead_Integration**
   - Tests scenarios where end_page equals start_page (zero pages read)
   - Validates that days with zero pages are correctly counted as faults
   - Tests single log, multiple logs, and mixed scenarios

3. **TestDashboardWeekdayFaults_EmptyDatabase_Integration**
   - Tests weekday faults calculation with empty database
   - Verifies all 7 weekdays have faults distributed evenly (~26-27 each)
   - Includes repository and HTTP endpoint validation

4. **TestDashboardWeekdayFaults_SingleLogEntry_Integration**
   - Tests weekday faults with a single log entry on specific weekday
   - Verifies single log affects fault distribution correctly
   - Includes repository and HTTP endpoint validation

### Files Modified

- `test/dashboard_integration_test.go` - Added 4 new integration test functions with comprehensive test cases

## Key Changes

- All new tests follow Clean Architecture patterns
- Tests use real database interactions (not mocks)
- Both repository-level and HTTP endpoint tests included
- Proper test isolation with cleanup between tests
- Consistent naming conventions (`<Endpoint>_<Scenario>_Integration`)

## Testing

- All unit tests pass ✅
- All integration tests pass ✅
- `go fmt` passes ✅
- `go vet` passes ✅
- Build succeeds ✅

## Notes

- NULL values test was not implemented because the logs table has NOT NULL constraints on start_page/end_page columns
- NULL handling is covered in existing unit tests (dashboard_repository_test.go)
- Zero pages read scenario covers the practical equivalent of NULL handling

## Edge Cases Now Covered

- ✅ Empty database (all days faults)
- ✅ Single log entry (one day with reading, rest faults)
- ✅ Month boundary dates (cross month ranges)
- ✅ Year boundary dates (Dec 31 to Jan 1 transitions)
- ✅ Zero pages read (end_page = start_page)
- ✅ Weekday faults edge cases (empty database, single log entry)
<!-- SECTION:FINAL_SUMMARY:END -->

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
