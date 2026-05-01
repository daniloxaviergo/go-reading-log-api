---
id: RDL-138
title: '[doc-12 Phase 3] Update TestDashboardFaultsChart integration test'
status: To Do
assignee:
  - workflow
created_date: '2026-05-01 15:02'
updated_date: '2026-05-01 15:39'
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
Fix TestDashboardFaultsChart integration test by updating test data setup to create realistic reading patterns with gaps. Update expected fault counts to match Rails calculation (counting zero-page days, not log entries). Verify percentage calculation uses correct fault count.
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task updates the `TestDashboardFaultsChart_Integration` integration test to correctly validate the faults calculation after the SQL query fix (RDL-139). The current test uses `ScenarioFaultsByWeekday()` which creates log entries that are incorrectly interpreted as "faults" when they should represent days WITH reading activity.

**Key Changes:**
1. **Fix test data setup**: Create a new scenario fixture that explicitly defines which days have reading activity vs. which days are faults (zero-page days)
2. **Update expected fault counts**: Calculate the correct fault count based on days with zero pages read, not log entry count
3. **Verify percentage calculation**: Ensure the gauge chart percentage is calculated from the correct fault count

**Why this approach:**
- The Rails implementation counts days with zero pages read as faults, not log entries
- The current test fixture creates log entries on specific weekdays, but these logs represent reading activity, NOT faults
- After RDL-139 fixes the SQL query to use `generate_series` and count zero-page days, the test must validate the correct behavior

**Architecture:**
- Test data fixture: `test/fixtures/dashboard/scenarios.go` - Add new scenario with explicit fault days
- Integration test: `test/dashboard_integration_test.go` - Update test expectations and validation
- Repository fix (dependency): `internal/adapter/postgres/dashboard_repository.go` - RDL-139 must be completed first

---

### 2. Files to Modify

| File | Change Type | Description |
|------|-------------|-------------|
| `test/fixtures/dashboard/scenarios.go` | **Modify** | Add new scenario `ScenarioFaultsChartCorrect()` with explicit fault day distribution |
| `test/dashboard_integration_test.go` | **Modify** | Update `TestDashboardFaultsChart_Integration` to use new scenario and validate correct fault count |
| `test/dashboard_integration_test.go` | **Modify** | Add assertion to verify fault percentage calculation matches expected value |

**New Scenario Fixture:**
```go
// ScenarioFaultsChartCorrect: Creates a 30-day period with specific fault distribution
// Days with logs = reading activity (not faults)
// Days without logs = faults (zero-page days)
func ScenarioFaultsChartCorrect() *Scenario {
    // 30-day period: Jan 1-30, 2024
    // Reading on: Jan 1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29 (15 days)
    // Faults on: Jan 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30 (15 days)
    // Expected faults: 15 out of 30 days = 50%
}
```

---

### 3. Dependencies

**Critical Prerequisites:**
1. **RDL-139** - "Update GetFaultsByDateRange SQL query with CTE for daily aggregation" MUST be completed first
   - This task fixes the SQL query to count days with zero pages, not log entries
   - Without this fix, the integration test will fail because the repository still returns incorrect data

2. **RDL-137** - "Add edge case integration tests for fault calculation" (recommended)
   - Provides additional test coverage for edge cases
   - Can be completed in parallel or before this task

**Setup Steps:**
1. Ensure test database is running: `make docker-up` or local PostgreSQL
2. Verify RDL-139 SQL query fix is merged
3. Run existing tests to establish baseline: `go test -v ./test/... -run TestDashboardFaultsChart_Integration`

---

### 4. Code Patterns

**Follow Existing Test Patterns:**
- Use `SetupTestDB()` for test database setup/teardown
- Use `dashboardFixtures.NewDashboardFixtures()` to load scenario data
- Use `handlers.NewDashboardHandler()` with repository and config dependencies
- Use `parseDashboardResponse()` to parse JSON response
- Use `assert.Equal()`, `assert.GreaterOrEqual()`, `assert.LessOrEqual()` from testify

**Scenario Fixture Pattern:**
```go
func ScenarioFaultsChartCorrect() *Scenario {
    baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    var logs []*LogFixture
    
    // Create logs for reading days (every other day)
    for i := 0; i < 15; i++ {
        logDate := baseDate.AddDate(0, 0, i*2) // Days 0, 2, 4, 6, ...
        logs = append(logs, &LogFixture{
            ID:        int64(i + 1),
            ProjectID: 1,
            Data:      logDate,
            StartPage: 0,
            EndPage:   10,
            WDay:      int(logDate.Weekday()),
        })
    }
    
    return &Scenario{
        Name:        "Faults Chart Correct",
        Description: "30-day period with 15 reading days and 15 fault days",
        Projects: []*ProjectFixture{
            {ID: 1, Name: "Faults Test Project", TotalPage: 200, Page: 50},
        },
        Logs: logs,
        Expected: &ExpectedResults{
            FaultCount: 15, // 30 days - 15 reading days = 15 faults
        },
    }
}
```

**Test Assertion Pattern:**
```go
// Verify fault percentage calculation
faultPercentage := series.Data[0].(float64)
expectedPercentage := 50.0 // 15 faults / 30 days * 100
assert.InDelta(t, expectedPercentage, faultPercentage, 0.01, 
    "Fault percentage should be 50%% (15 faults out of 30 days)")
```

---

### 5. Testing Strategy

**Test Coverage:**
1. **Happy Path**: 30-day period with 15 reading days, 15 fault days → 50% fault rate
2. **Edge Case**: Single day with no logs → 100% fault rate
3. **Edge Case**: All days have reading → 0% fault rate

**Test Cases:**
```go
// Main test: Verify fault count and percentage
t.Run("30-day period with 50% faults", func(t *testing.T) {
    // Setup: Load ScenarioFaultsChartCorrect
    // Execute: Call /v1/dashboard/echart/faults.json
    // Verify: 
    //   - Fault count = 15 (days with zero pages)
    //   - Percentage = 50% (15/30 * 100)
})

// Edge case: All days have reading
t.Run("All days have reading", func(t *testing.T) {
    // Setup: Create logs for all 30 days
    // Execute: Call endpoint
    // Verify: Fault count = 0, Percentage = 0%
})

// Edge case: No reading at all
t.Run("No reading activity", func(t *testing.T) {
    // Setup: No logs in database
    // Execute: Call endpoint
    // Verify: Fault count = 30, Percentage = 100%
})
```

**Validation Steps:**
1. Run unit tests: `go test -v ./test/unit/...`
2. Run integration tests: `go test -v ./test/... -run TestDashboardFaultsChart_Integration`
3. Verify test coverage: `go test -cover ./test/...`
4. Run full test suite: `go test ./...`

**Edge Cases to Cover:**
- Empty database (no logs)
- Single day date range
- Month boundary dates
- Leap year February
- Logs with zero pages (start_page = end_page)

---

### 6. Risks and Considerations

**Blocking Issues:**
1. **RDL-139 must be completed first**: The SQL query fix is a prerequisite. Without it, the repository returns incorrect fault counts (log count instead of zero-page day count).
   - **Mitigation**: Coordinate with team to ensure RDL-139 is merged before starting RDL-138

2. **Test data fixture complexity**: Creating realistic reading patterns with gaps requires careful date calculation.
   - **Mitigation**: Use explicit date calculations with `time.Date()` and `AddDate()` for reproducibility

3. **Percentage calculation verification**: The gauge chart percentage depends on both fault count and total days in range.
   - **Mitigation**: Add explicit calculation in test to verify: `percentage = (fault_count / total_days) * 100`

**Potential Pitfalls:**
- **Timezone issues**: Ensure test dates use UTC to avoid timezone-related discrepancies
- **Date range inclusivity**: Verify both start and end dates are included in the count (BETWEEN is inclusive)
- **Weekday calculation**: PostgreSQL DOW (0=Sunday, 6=Saturday) must match Go's `time.Weekday()`

**Deployment Considerations:**
- This is a test-only change; no production deployment required
- Ensure all existing integration tests still pass after the change
- Update API documentation if fault calculation behavior changes

**Rollback Plan:**
- If issues arise, revert the test file changes
- Keep the new scenario fixture for future use
- Investigate failures in a separate branch

---

**Implementation Checklist:**
- [ ] Verify RDL-139 SQL query fix is completed and merged
- [ ] Add `ScenarioFaultsChartCorrect()` to `test/fixtures/dashboard/scenarios.go`
- [ ] Update `TestDashboardFaultsChart_Integration` to use new scenario
- [ ] Add assertion for fault count validation
- [ ] Add assertion for percentage calculation validation
- [ ] Run integration tests to verify changes
- [ ] Run full test suite to ensure no regressions
- [ ] Update task description with test results
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
