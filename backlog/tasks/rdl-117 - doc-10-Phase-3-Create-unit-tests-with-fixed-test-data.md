---
id: RDL-117
title: '[doc-10 Phase 3] Create unit tests with fixed test data'
status: To Do
assignee:
  - workflow
created_date: '2026-04-28 00:29'
updated_date: '2026-04-28 03:59'
labels:
  - testing
  - phase-3
  - backend
dependencies: []
documentation:
  - doc-010
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create unit tests in dashboard_service_test.go with fixed test data to verify mean_day calculation matches Rails V1::MeanLog exactly. Use deterministic dates and page counts.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Unit tests use fixed test data
- [ ] #2 All calculation tests pass with expected values
- [ ] #3 Tests verify Rails parity
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task requires creating comprehensive unit tests for the `DayService` with fixed, deterministic test data to verify that the `mean_day` calculation matches the Rails `V1::MeanLog` implementation exactly.

**Technical Strategy:**
- Create a new test file or enhance existing `day_service_test.go` with fixed test data scenarios
- Use deterministic dates (e.g., fixed "today" date like 2026-04-21) and fixed page counts
- Mock the `DashboardRepository` to return predictable data
- Compare Go calculations against expected values derived from the Rails algorithm

**Rails V1::MeanLog Algorithm Reference:**
```ruby
# Rails calculates mean as:
total_pages = sum of read_pages for the weekday
count_reads = number of 7-day intervals from first_read to log_date
mean = (total_pages / count_reads).round(3)
```

**Go Implementation to Test:**
```go
// Go calculates mean as:
projectMean = AVG(read_pages) for logs on current weekday
mean_day = average of all project means
```

**Key Test Scenarios:**
1. **Normal calculation** - Multiple projects with data on the target weekday
2. **Single project** - Verify calculation with one project's data
3. **Empty data** - No logs for target weekday (should return 0.0)
4. **Edge cases** - Zero pages, single log entry, large page counts
5. **Rounding verification** - Ensure 3-decimal rounding matches Rails

**Why This Approach:**
- Unit tests with mocks isolate the service logic from database
- Fixed dates ensure deterministic, reproducible test results
- Clear expected values allow verification of Rails parity
- Follows existing test patterns in the codebase (see `day_service_test.go`)

### 2. Files to Modify

**Files to Read/Analyze:**
- `test/unit/day_service_test.go` - Current test structure and patterns
- `internal/service/dashboard/day_service.go` - Implementation to test
- `rails-app/app/classes/v1/mean_log.rb` - Rails reference algorithm
- `internal/domain/dto/dashboard_response.go` - StatsData DTO structure
- `internal/repository/dashboard_repository.go` - Repository interface

**Files to Modify:**
- `test/unit/day_service_test.go` - Add new test cases with fixed test data

**New Files to Create:**
- None (enhance existing test file)

**Specific Test Cases to Add:**

1. **TestDayService_CalculateMeanDay_RailsParity** - Main test verifying Rails algorithm match
   - Fixed date: 2026-04-21 (Monday, weekday=1)
   - Fixed test data: 2-3 projects with known page counts
   - Expected mean_day calculated from Rails formula

2. **TestDayService_CalculateMeanDay_MultipleWeekdays** - Test across different weekdays
   - Test data spanning multiple weeks
   - Verify 7-day interval calculation

3. **TestDayService_CalculateMeanDay_EdgeCases** - Comprehensive edge case coverage
   - No logs for weekday
   - Single log entry
   - Zero pages read
   - Large page counts (1000+)

4. **TestDayService_CalculateWeeklyStats_FixedData** - Full integration test with fixed data
   - All fields: previous_week_pages, last_week_pages, per_pages, mean_day, spec_mean_day
   - Verify rounding to 3 decimals

### 3. Dependencies

**Prerequisites:**
- Existing test infrastructure in `test/unit/` is functional
- `MockDashboardRepository` already exists in `day_service_test.go`
- `MockUserConfigService` already exists
- `SetTestDate` / `GetTodayFunc` date injection is available

**Blocking Issues:**
- None identified - all dependencies exist

**Setup Steps:**
1. Review existing test patterns in `day_service_test.go`
2. Understand Rails V1::MeanLog algorithm from `rails-app/app/classes/v1/mean_log.rb`
3. Calculate expected values for test data manually (to verify against)

**Test Data Preparation:**
- Define fixed "today" date (e.g., 2026-04-21)
- Create deterministic log entries with known read_pages
- Calculate expected mean_day values using Rails formula
- Document expected values in test comments

### 4. Code Patterns

**Following Existing Conventions:**

1. **Test Structure:**
   ```go
   func TestDayService_SpecificMethod(t *testing.T) {
       // Setup
       mockRepo := &MockDashboardRepository{}
       mockConfig := &MockUserConfigService{mockPredictionPct: 0.15}
       dayService := dashboard.NewDayService(mockRepo, mockConfig)
       
       // Test cases
       t.Run("test case name", func(t *testing.T) {
           // Arrange
           // Act
           // Assert
       })
   }
   ```

2. **Date Injection Pattern:**
   ```go
   fixedDate := time.Date(2026, 4, 21, 12, 0, 0, 0, time.UTC)
   defer dashboard.SetTestDate(time.Now())
   dashboard.SetTestDate(fixedDate)
   ```

3. **Mock Setup Pattern:**
   ```go
   mockRepo.mockGetProjectWeekdayMean = func(ctx context.Context, projectID int64, weekday int) (float64, error) {
       if projectID == 1 {
           return 10.5, nil
       }
       return 0.0, nil
   }
   ```

4. **Assertion Pattern:**
   ```go
   assert.Equal(t, expected, actual, "description")
   assert.InDelta(t, expected, actual, tolerance, "description")
   require.NoError(t, err)
   ```

5. **Naming Conventions:**
   - Test functions: `Test<Service>_<Method>_<Scenario>`
   - Test case names: descriptive, lowercase (e.g., "normal calculation", "zero previous week")
   - Mock types: `Mock<InterfaceName>`

6. **Documentation Pattern:**
   - Add comments explaining the Rails algorithm being tested
   - Document expected values and their calculation
   - Include inline comments for complex calculations

### 5. Testing Strategy

**Test Types:**

1. **Unit Tests (Primary):**
   - Use `MockDashboardRepository` to isolate service logic
   - Test each calculation method independently
   - Focus on `CalculateMeanDay`, `CalculateWeeklyStats`

2. **Edge Case Coverage:**
   - Empty aggregates (no projects)
   - Zero pages read
   - Single log entry
   - Large page counts
   - Division by zero scenarios

3. **Rounding Verification:**
   - Test that values are rounded to 3 decimal places
   - Use `assert.InDelta(t, expected, actual, 0.001)` for float comparisons

**Edge Cases to Cover:**

| Scenario | Input | Expected Output |
|----------|-------|-----------------|
| No logs for weekday | Empty mock | 0.0 |
| Single project | One project, 10.5 mean | 10.5 |
| Multiple projects | Two projects, 10.0 and 20.0 | 15.0 |
| Zero pages | All logs have 0 pages | 0.0 |
| Large values | Pages > 1000 | Correct mean, rounded |
| Float precision | 7.123456 | 7.123 |

**Verification Approach:**

1. **Manual Calculation:**
   - Calculate expected values using Rails formula before writing tests
   - Document calculations in test comments

2. **Cross-Reference:**
   - Compare test expected values with Rails API responses (if available)
   - Use `rails-app/spec/services/v1/mean_log_spec.rb` as reference

3. **Test Execution:**
   ```bash
   go test -v ./test/unit/... -run TestDayService
   ```

4. **Coverage Goal:**
   - All code paths in `CalculateMeanDay` covered
   - All edge cases tested
   - No panics or unhandled errors

### 6. Risks and Considerations

**Known Algorithm Differences:**

⚠️ **Critical:** The Rails and Go implementations use different algorithms:

- **Rails:** `mean = total_pages / count_7day_intervals`
- **Go:** `mean = AVG(read_pages) for logs on weekday`

This means:
1. The Go implementation may NOT match Rails exactly
2. Tests should verify the current Go behavior, not necessarily Rails parity
3. If Rails parity is required, the Go implementation may need refactoring

**Potential Pitfalls:**

1. **Date Calculation Discrepancies:**
   - Rails uses `step(7).map { |d| d }.size` for interval counting
   - Go uses simple averaging
   - May produce different results for the same data

2. **Floating Point Precision:**
   - Rounding differences between Ruby and Go
   - Use `assert.InDelta` with appropriate tolerance (0.001)

3. **Weekday Calculation:**
   - Rails `wday`: 0=Sunday, 1=Monday, ..., 6=Saturday
   - Go `time.Weekday()`: 0=Sunday, 1=Monday, ..., 6=Saturday
   - Should be compatible, but verify in tests

4. **Test Data Complexity:**
   - Creating realistic test data that exercises all paths
   - Ensuring mock returns are deterministic

**Mitigation Strategies:**

1. **Document Algorithm Differences:**
   - Add comments in test file explaining the difference
   - Note which algorithm is being tested

2. **Start with Simple Tests:**
   - Test basic functionality first
   - Add complex scenarios incrementally

3. **Verify Expected Values:**
   - Manually calculate expected values
   - Double-check against Rails implementation

**Deployment Considerations:**

- None (unit tests don't affect production code)
- Tests should be run as part of CI/CD pipeline
- Add to `make test` command

**Documentation Updates:**

- Update `QWEN.md` with test coverage information
- Document any algorithm differences discovered
- Add test examples to developer guide
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
