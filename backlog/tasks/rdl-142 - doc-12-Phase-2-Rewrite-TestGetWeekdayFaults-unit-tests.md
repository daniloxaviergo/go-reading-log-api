---
id: RDL-142
title: '[doc-12 Phase 2] Rewrite TestGetWeekdayFaults unit tests'
status: To Do
assignee:
  - thomas
created_date: '2026-05-01 15:07'
updated_date: '2026-05-01 17:14'
labels:
  - bugfix
  - testing
  - phase-2
dependencies: []
documentation:
  - doc-012
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Rewrite all TestGetWeekdayFaults test cases to validate correct weekday grouping. Create test cases for: all weekdays represented, some weekdays missing, zero faults for all days, and all days are faults. Update expected weekday distribution maps to match Rails grouping logic.
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task involves rewriting the `TestGetWeekdayFaults` unit tests to validate correct weekday grouping logic that matches the Rails implementation. The fault calculation logic was fixed in RDL-139/RDL-140 to count days with zero pages read (not log entries).

**Technical Strategy:**
- The repository method `GetWeekdayFaults` uses a CTE-based SQL query that:
  1. Aggregates daily reading activity (`daily_read` CTE)
  2. Generates all dates in range using `generate_series` (`all_dates` CTE)
  3. LEFT JOINs to identify days with no reading activity
  4. Groups by weekday (0-6 = Sunday-Saturday) and counts faults
- Tests must validate that the weekday distribution map correctly reflects days WITHOUT reading activity
- Follow the Rails logic: a "fault" = a day where `sum(read_pages) = 0`

**Why this approach:**
- Aligns with the existing CTE-based implementation in `dashboard_repository.go`
- Matches the Rails `V1::Dashboard::WeekdayFaults` calculation logic
- Ensures all 7 weekdays (0-6) are present in the result with default value of 0

### 2. Files to Modify

| File | Change Type | Description |
|------|-------------|-------------|
| `test/unit/dashboard_repository_test.go` | Modify | Rewrite `TestDashboardRepository_GetWeekdayFaults` and `TestDashboardRepository_GetWeekdayFaults_EmptyRange` with comprehensive test cases |
| `test/unit/dashboard_repository_test.go` | Add | New test cases: `TestDashboardRepository_GetWeekdayFaults_AllWeekdaysRepresented`, `TestDashboardRepository_GetWeekdayFaults_SomeWeekdaysMissing`, `TestDashboardRepository_GetWeekdayFaults_ZeroFaults`, `TestDashboardRepository_GetWeekdayFaults_AllDaysAreFaults` |

### 3. Dependencies

| Dependency | Status | Notes |
|------------|--------|-------|
| RDL-139 (GetFaultsByDateRange SQL fix) | ✅ Done | Prerequisite - SQL query fix completed |
| RDL-140 (GetWeekdayFaults SQL fix) | ✅ Done | Prerequisite - SQL query fix completed |
| RDL-141 (Context timeout verification) | ✅ Done | Prerequisite - Context handling verified |
| Test database setup | ✅ Available | Uses `test.SetupTestDB()` helper |

**No blocking dependencies** - All prerequisite tasks are complete.

### 4. Code Patterns

Follow existing patterns in the codebase:

**Test Structure:**
```go
func TestDashboardRepository_GetWeekdayFaults_<Scenario>(t *testing.T) {
    // Setup
    helper, err := test.SetupTestDB()
    require.NoError(t, err)
    defer helper.Close()

    err = helper.SetupTestSchema()
    require.NoError(t, err)

    repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

    // Create test project
    ctx := context.Background()
    _, err = helper.Pool.Exec(ctx, 
        "INSERT INTO projects (id, name, total_page, page) VALUES ($1, $2, $3, $4)",
        1, "Test Project", 100, 0)
    require.NoError(t, err)

    // Create test logs (use helper function)
    err = createTestLogs(helper.Pool, []testLog{...})
    require.NoError(t, err)

    // Execute
    stats, err := repo.GetWeekdayFaults(context.Background(), startDate, endDate)

    // Verify
    assert.NoError(t, err)
    assert.NotNil(t, stats)
    assert.NotNil(t, stats.Faults)
    // Validate specific weekday counts
    assert.Equal(t, expectedCount, stats.Faults[weekday])
}
```

**Helper Functions (already exist):**
- `createTestLogs(pool, []testLog)` - Creates log entries
- `ptr(string)` - Returns pointer to string for nullable fields

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

#### Test Case 1: All Weekdays Represented
- **Scenario**: 4-week period (28 days = 4 of each weekday)
- **Setup**: Skip reading on all Mondays (weekday 1)
- **Expected**: `map[1:4, 0:0, 2:0, 3:0, 4:0, 5:0, 6:0]`
- **Purpose**: Validates correct weekday grouping when one weekday has all faults

#### Test Case 2: Some Weekdays Missing
- **Scenario**: Irregular date range that doesn't cover all weekdays equally
- **Setup**: 10-day range starting on Wednesday (covers Wed, Thu, Fri, Sat, Sun, Mon, Tue, Wed, Thu, Fri)
- **Expected**: Specific distribution based on actual dates
- **Purpose**: Validates handling of partial weekday coverage

#### Test Case 3: Zero Faults for All Days
- **Scenario**: All days in range have reading activity
- **Setup**: Create logs for every day in a 7-day range
- **Expected**: `map[0:0, 1:0, 2:0, 3:0, 4:0, 5:0, 6:0]`
- **Purpose**: Validates no false positives when all days have reading

#### Test Case 4: All Days Are Faults
- **Scenario**: No logs in the date range
- **Setup**: Create test project but no logs
- **Expected**: Distribution based on number of days per weekday in range
- **Example**: 8-day range (Sat, Sun, Mon, Tue, Wed, Thu, Fri, Sat) → `map[6:2, 0:1, 1:1, 2:1, 3:1, 4:1, 5:1]`
- **Purpose**: Validates all days counted as faults when no reading activity

**Edge Cases to Cover:**
- Single-day range (start == end)
- Month boundary dates
- Leap year February
- NULL values in start_page/end_page
- Logs with zero pages (end_page == start_page)

**Validation Approach:**
1. Use `require.NoError(t, err)` for error handling
2. Use `assert.NotNil(t, stats)` for result validation
3. Use `assert.Equal(t, expected, actual)` for exact value matching
4. Verify all 7 weekdays present in result map

### 6. Risks and Considerations

**Known Risks:**
1. **Date Calculation Complexity**: Weekday distribution depends on exact date range. Must carefully calculate expected values based on actual calendar dates.
   - **Mitigation**: Use Go's `time.Weekday()` to verify test date weekdays
   - **Mitigation**: Add inline comments explaining date calculations

2. **Test Data Isolation**: Each test creates its own database state
   - **Mitigation**: Use `defer helper.Close()` and unique project IDs per test
   - **Mitigation**: Tests use independent date ranges to avoid overlap

3. **Rails Parity Validation**: Hard to verify exact Rails behavior without running Rails app
   - **Mitigation**: Follow the SQL logic documented in doc-012
   - **Mitigation**: Use CTE-based query that matches Rails logic exactly

**Considerations:**
- **PostgreSQL DOW**: Ensure correct mapping (0=Sunday, not 0=Monday)
- **Inclusive Date Range**: BETWEEN $1 AND $2 includes both endpoints
- **Weekday Distribution**: Must account for partial weeks at range boundaries
- **Null Handling**: Repository returns 0 for weekdays with no faults (not missing keys)

**Rollback Plan:**
- If tests fail unexpectedly, revert to previous test version
- Investigate by running individual test cases with verbose output
- Compare expected vs actual weekday distributions

**Success Criteria:**
- All new test cases pass
- Existing tests continue to pass (no regressions)
- Test coverage remains at 100% for `GetWeekdayFaults` method
- Code passes `go fmt` and `go vet`
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
