---
id: RDL-145
title: '[doc-12 Phase 4] Add Rails comparison test cases'
status: To Do
assignee:
  - workflow
created_date: '2026-05-01 15:08'
updated_date: '2026-05-01 18:16'
labels:
  - testing
  - rails-parity
  - phase-4
dependencies: []
documentation:
  - doc-012
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add comparison test cases: 30-day range with random gaps, 6-month range for weekday faults validation, and edge case with leap year February. Verify Rails and Go outputs match exactly for all test scenarios.
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task requires creating Rails comparison test cases to validate that the Go implementation of fault calculation logic matches the Rails API exactly. The approach involves:

**Test Strategy:**
- Create three specific test scenarios as defined in the task description
- Each scenario will test the `GetFaultsByDateRange` and `GetWeekdayFaults` repository methods
- Compare Go outputs against expected Rails outputs (calculated based on Rails logic documented in doc-012)
- Tests will be integration tests using the real test database

**Three Test Scenarios:**
1. **30-day range with random gaps**: Tests fault calculation with irregular reading patterns
2. **6-month range for weekday faults validation**: Tests `GetWeekdayFaults` with extended date range
3. **Edge case with leap year February**: Tests date boundary handling including Feb 29

**Implementation Pattern:**
- Follow existing patterns from `test/dashboard_integration_test.go`
- Use existing fixture infrastructure (`test/fixtures/dashboard/`)
- Leverage existing test helper (`test/test_helper.go`)
- Create new scenario fixtures in `test/fixtures/dashboard/scenarios.go`
- Add integration tests in a new file `test/faults_rails_comparison_test.go`

**Why This Approach:**
- Reuses existing test infrastructure (fixtures, helpers, database setup)
- Maintains consistency with existing test patterns
- Provides clear validation against Rails logic without requiring Rails API to be running
- Tests are deterministic with fixed dates and predictable outcomes

### 2. Files to Modify

**Files to Create:**
1. `test/faults_rails_comparison_test.go` - Main comparison test file with three test scenarios
   - `TestFaultsComparison_30DayRandomGaps` - 30-day range test
   - `TestFaultsComparison_6MonthWeekdayValidation` - 6-month weekday faults test
   - `TestFaultsComparison_LeapYearFebruary` - Leap year edge case test

**Files to Modify:**
1. `test/fixtures/dashboard/scenarios.go` - Add three new scenario functions:
   - `ScenarioFaults30DayRandomGaps()` - Creates 30-day period with random gaps
   - `ScenarioFaults6MonthWeekday()` - Creates 6-month period for weekday validation
   - `ScenarioFaultsLeapYearFebruary()` - Creates February leap year scenario

**Files to Read/Reference:**
1. `test/dashboard_integration_test.go` - Reference for integration test patterns
2. `test/fixtures/dashboard/scenarios.go` - Add new scenarios following existing patterns
3. `test/fixtures/dashboard/fixtures.go` - Understand fixture loading mechanism
4. `internal/adapter/postgres/dashboard_repository.go` - Reference for repository methods being tested
5. `backlog/docs/doc-012 - PRD-Fix-Faults-Calculation-Logic-Rails-Parity-for-Dashboard-API.md` - Rails logic reference
6. `test/test_helper.go` - Database setup/teardown patterns

### 3. Dependencies

**Prerequisites:**
- RDL-139: `GetFaultsByDateRange` SQL query fix (DONE)
- RDL-140: `GetWeekdayFaults` SQL query fix (DONE)
- Test database configured (`reading_log_test`)
- PostgreSQL running with test schema

**Existing Infrastructure to Leverage:**
- `test/test_helper.go` - `SetupTestDB()`, `SetupTestSchema()`, `ClearTestData()`
- `test/fixtures/dashboard/` - Fixture management and scenario definitions
- `test/dashboard_integration_test.go` - Integration test patterns and helper functions
- `dto.GetToday()` and `dto.SetTestDate()` for date control in tests

**No Blocking Dependencies:**
- All required repository methods are implemented (RDL-139, RDL-140 completed)
- Test infrastructure exists and is functional
- Rails logic is documented in doc-012 for expected value calculations

### 4. Code Patterns

**Test Structure Pattern (following `dashboard_integration_test.go`):**
```go
func TestFaultsComparison_<ScenarioName>(t *testing.T) {
    if !IsTestDatabase() {
        t.Skip("Test database not configured")
    }

    // Set fixed date for reproducibility
    dto.SetTestDate(<fixed_date>)
    defer dto.ResetTestDate()

    helper, err := SetupTestDB()
    require.NoError(t, err)
    defer helper.Close()

    // Create database tables
    err = helper.SetupTestSchema()
    require.NoError(t, err)

    // Load fixture data
    fixtureManager := dashboardFixtures.NewDashboardFixtures(helper.Pool)
    scenario := dashboardFixtures.Scenario<ScenarioName>()
    err = fixtureManager.LoadScenario(scenario)
    require.NoError(t, err)

    // Create repository
    repo, err := createTestRepository(helper.Pool)
    require.NoError(t, err)

    // Test repository method
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    result, err := repo.GetFaultsByDateRange(ctx, startDate, endDate)
    require.NoError(t, err)

    // Validate against expected Rails output
    assert.Equal(t, expectedFaultCount, result.FaultCount)
}
```

**Scenario Fixture Pattern (following existing scenarios):**
```go
func ScenarioFaults30DayRandomGaps() *Scenario {
    baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    var logs []*LogFixture

    // Create logs with random gaps (specific pattern for reproducibility)
    // Example: reading on days 0, 3, 5, 8, 10, 13, 15, 18, 20, 23, 25, 28 (12 reading days)
    // Expected faults: 30 - 12 = 18 faults

    logID := int64(500)
    readingDays := []int{0, 3, 5, 8, 10, 13, 15, 18, 20, 23, 25, 28}
    for i, day := range readingDays {
        logDate := baseDate.AddDate(0, 0, day)
        logs = append(logs, &LogFixture{
            ID:        logID,
            ProjectID: 1,
            Data:      logDate,
            StartPage: i * 10,
            EndPage:   (i + 1) * 10,
            WDay:      int(logDate.Weekday()),
        })
        logID++
    }

    return &Scenario{
        Name:        "30-Day Random Gaps",
        Description: "30-day period with 12 reading days and 18 fault days",
        Projects: []*ProjectFixture{
            {ID: 1, Name: "30-Day Test Project", TotalPage: 200, Page: 120},
        },
        Logs: logs,
        Expected: &ExpectedResults{
            EchartData: map[string]interface{}{
                "fault_count": 18, // 30 - 12 = 18
            },
        },
    }
}
```

**Naming Conventions:**
- Test functions: `TestFaultsComparison_<ScenarioName>`
- Scenario functions: `ScenarioFaults<ScenarioName>()`
- Use snake_case for JSON fields, camelCase for Go structs
- Fixed dates for reproducibility (e.g., `time.Date(2024, 1, 1, ...)`)

**Integration Patterns:**
- Use `context.WithTimeout` with 15-second timeout (matching `dashboardContextTimeout`)
- Always call `defer dto.ResetTestDate()` after `dto.SetTestDate()`
- Use `require.NoError(t, err)` for setup errors, `assert.Equal()` for validations
- Clean up test data between test cases if needed

### 5. Testing Strategy

**Test Types:**
1. **Integration Tests** - Primary test type using real database
   - Test actual SQL query execution
   - Validate database interactions
   - Verify context timeout handling
   - Check error handling

**Test Coverage:**
1. **Scenario 1: 30-Day Random Gaps**
   - Date range: Jan 1-30, 2024 (30 days)
   - Reading days: 12 specific days with logs
   - Expected faults: 18 days without reading
   - Validates: Basic fault counting with irregular patterns

2. **Scenario 2: 6-Month Weekday Validation**
   - Date range: Oct 1, 2025 - Apr 1, 2026 (~183 days)
   - Reading days: Specific weekdays (e.g., Sun, Tue, Thu, Sat)
   - Expected faults: All Mon, Wed, Fri days (~78-84 faults)
   - Validates: Weekday fault distribution with `GetWeekdayFaults`

3. **Scenario 3: Leap Year February**
   - Date range: Feb 1-29, 2024 (29 days, 2024 is leap year)
   - Reading days: Specific pattern (e.g., every other day)
   - Expected faults: Calculated based on reading pattern
   - Validates: Leap year date handling, Feb 29 inclusion

**Edge Cases to Cover:**
- Empty result sets (no logs in range)
- Single day ranges
- Month boundaries
- Year boundaries
- Leap year February 29
- NULL value handling in logs
- Zero-page logs (start_page = end_page)

**Validation Approach:**
- Compare actual fault counts against pre-calculated expected values
- Validate all 7 weekdays present in `GetWeekdayFaults` result
- Verify context timeout is respected
- Check error handling for invalid date ranges

**Test Execution:**
```bash
# Run all faults comparison tests
go test -v ./test/... -run "TestFaultsComparison"

# Run specific scenario
go test -v ./test/... -run "TestFaultsComparison_30DayRandomGaps"

# Run with coverage
go test -cover ./test/... -run "TestFaultsComparison"
```

### 6. Risks and Considerations

**Known Risks:**
1. **Date Calculation Complexity** - Risk of off-by-one errors in date range calculations
   - Mitigation: Use pre-calculated expected values based on documented Rails logic
   - Mitigation: Include detailed test comments explaining date range boundaries

2. **Weekday Distribution Variance** - 6-month range may have slightly different weekday counts
   - Mitigation: Use exact date counting in test setup (countWeekdaysInRange helper)
   - Mitigation: Document exact expected values in scenario fixtures

3. **Timezone Handling** - Server timezone affects date boundaries
   - Mitigation: Use `dto.SetTestDate()` with UTC timezone for reproducibility
   - Mitigation: Document TZ requirement in test comments

4. **Database State Pollution** - Tests may interfere with each other
   - Mitigation: Each test creates fresh database via `SetupTestDB()`
   - Mitigation: Use `defer helper.Close()` for cleanup

**Trade-offs:**
- **Expected Values vs Rails API**: Tests use pre-calculated expected values instead of live Rails API comparison
  - Pros: Tests run without Rails API, faster execution, more deterministic
  - Cons: Expected values must be manually calculated and verified
  - Decision: This approach is chosen because Rails API may not always be available

**Deployment Considerations:**
- Tests should be run as part of CI/CD pipeline
- Tests require PostgreSQL with test database
- Tests should be marked as integration tests (not unit tests)

**Documentation Requirements:**
- Add test scenario descriptions to `docs/TESTING.md` or similar
- Document expected values calculation methodology
- Include examples of date range calculations

**Acceptance Criteria Alignment:**
- ✅ All unit tests pass - Tests use standard `testing` package
- ✅ All integration tests pass - Tests use real database via `SetupTestDB()`
- ✅ go fmt and go vet pass - Standard Go tooling
- ✅ Clean Architecture layers followed - Tests repository layer via adapter
- ✅ Error responses consistent - Tests verify error handling
- ✅ HTTP status codes correct - N/A (repository-level tests)
- ✅ Documentation updated - Add to QWEN.md with test examples
- ✅ New code paths include error path tests - Include error handling validations
- ✅ HTTP handlers test both success and error - Repository-level, not handler-level
- ✅ Integration tests verify actual database interactions - Yes, uses real PostgreSQL

**Related Tasks:**
- RDL-139: GetFaultsByDateRange SQL query fix (dependency - DONE)
- RDL-140: GetWeekdayFaults SQL query fix (dependency - DONE)
- RDL-146: Create faults calculation documentation (related documentation task)
- RDL-147: Update API documentation with fault definition (related documentation task)
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
