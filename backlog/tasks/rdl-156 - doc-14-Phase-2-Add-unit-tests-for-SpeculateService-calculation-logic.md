---
id: RDL-156
title: '[doc-14 Phase 2] Add unit tests for SpeculateService calculation logic'
status: To Do
assignee:
  - thomas
created_date: '2026-05-10 10:47'
updated_date: '2026-05-10 13:41'
labels:
  - testing
  - unit-tests
  - service
  - phase-2
dependencies: []
documentation:
  - doc-014
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create test/unit/service/dashboard/speculate_service_test.go with comprehensive unit tests for all SpeculateService methods. Test scenarios include: CalculateHistoricalMean with various data distributions, CalculateSpeculativeMean edge cases (zero mean, nil data), GenerateXAxisLabels format validation, and zero-fill logic for missing days.

Tests must use mock repository and cover all edge cases defined in acceptance criteria.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 TestCalculateHistoricalMean covers normal data, empty data, and single-read scenarios
- [ ] #2 TestCalculateSpeculativeMean validates 10% buffer and zero-mean edge case
- [ ] #3 TestGenerateXAxisLabels verifies 'DD-MMM (Day)' format for all 15 dates
- [x] #4 TestZeroFillLogic validates missing days are filled with zero values
- [ ] #5 TestCalculateHistoricalMean_WeekdayGrouping validates weekday-specific calculations
- [x] #6 All tests achieve >80% code coverage
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task requires creating comprehensive unit tests for the `SpeculateService` calculation logic. The service implements weekday-based historical mean calculation following the Rails V1::MeanLog algorithm.

**Technical Strategy:**
- Use existing mock repository pattern (`MockDashboardRepositoryForSpeculate`) already defined in the test file
- Follow the testing conventions established in `faults_service_test.go` and other service tests
- Test all public methods: `CalculateHistoricalMean`, `CalculateSpeculativeMean`, `GenerateXAxisLabels`, `GenerateSeriesData`, `GetChartConfig`, `GetDateRangeLast15Days`
- Test edge cases: empty data, nil values, invalid timestamps, zero divisions
- Achieve >80% code coverage using `go test -cover`

**Algorithm Coverage:**
1. **CalculateHistoricalMean**: Test weekday filtering, 7-day interval calculation, nil return scenarios
2. **CalculateSpeculativeMean**: Test 10% buffer application, zero/negative mean handling, rounding to 3 decimals
3. **GenerateXAxisLabels**: Test 15-date generation, 'DD-MMM (Day)' format validation, month boundaries
4. **GenerateSeriesData**: Test zero-fill logic, timestamp parsing, pages accumulation, speculative mean calculation
5. **GetChartConfig**: Test full chart configuration, series generation, mark points/lines
6. **Error Handling**: Test repository error propagation, invalid timestamp handling

### 2. Files to Modify

**Files to Read/Review:**
- `internal/service/dashboard/speculate_service.go` - Source code to test
- `test/unit/service/dashboard/speculate_service_test.go` - Existing tests (already has good coverage)
- `test/testutil/mock_dashboard_repository.go` - Mock repository for testing
- `internal/service/dashboard/day_service.go` - Reference for UserConfigProvider interface
- `internal/domain/dto/dashboard_response.go` - DTO structures used
- `internal/repository/dashboard_repository.go` - Repository interface

**Files to Modify:**
- `test/unit/service/dashboard/speculate_service_test.go` - Add missing test cases to achieve >80% coverage

**Potential New Files:**
- None required (existing test file structure is sufficient)

### 3. Dependencies

**Prerequisites:**
- RDL-155: [doc-14 Phase 2] Implement SpeculateService weekday-based calculation methods (DONE)
  - This task implemented the actual service code that needs testing
- Existing mock repository infrastructure (`test/testutil/mock_dashboard_repository.go`)
- Go 1.25.7 with testing package and testify library

**Existing Test Infrastructure:**
- `MockDashboardRepositoryForSpeculate` - Custom mock with func fields
- `MockUserConfigProviderForSpeculate` - Config provider mock
- Test helpers from `test/testutil/`

**No blocking issues identified** - All dependencies are in place.

### 4. Code Patterns

**Testing Conventions to Follow:**
1. **Test Naming**: `Test<MethodName>_<Scenario>` (e.g., `TestCalculateSpeculativeMean_ZeroMean`)
2. **Test Structure**: Arrange-Act-Assert pattern with comments
3. **Mock Setup**: Use func fields in mock repository for precise control
4. **Assertions**: Use `require` for setup errors, `assert` for value checks
5. **Decimal Precision**: Use `assert.InDelta` for float comparisons (tolerance: 0.001)
6. **Date Handling**: Use `dto.GetToday()` and `dto.SetTestDate()` for date abstraction

**Code Style:**
- Follow existing test file structure with section comments (`// ===...===`)
- Group tests by method under clear section headers
- Use table-driven tests for repetitive scenarios (precision, edge cases)
- Keep test functions focused on single behavior

**Mock Patterns:**
```go
mockRepo := &MockDashboardRepositoryForSpeculate{}
mockRepo.mockGetWeekdayMeanWithIntervals = func(ctx context.Context, weekday int, currentDate time.Time) (*float64, error) {
    return &expectedMean, nil
}
```

### 5. Testing Strategy

**Test Coverage Plan:**

**A. CalculateHistoricalMean Tests (3 scenarios):**
1. ✅ `TestSpeculateService_CalculateHistoricalMean_WithData` - Normal data with valid mean
2. ✅ `TestSpeculateService_CalculateHistoricalMean_NoData` - Nil return when no data
3. ❓ **ADD**: `TestSpeculateService_CalculateHistoricalMean_RepositoryError` - Error propagation

**B. CalculateSpeculativeMean Tests (4 scenarios):**
1. ✅ `TestCalculateSpeculativeMean_Normal` - 10% buffer application (25.0 → 27.5)
2. ✅ `TestCalculateSpeculativeMean_ZeroMean` - Zero returns 0.0
3. ✅ `TestCalculateSpeculativeMean_NegativeMean` - Negative returns 0.0
4. ✅ `TestCalculateSpeculativeMean_Rounding` - 3 decimal precision (25.123456 → 27.636)

**C. GenerateXAxisLabels Tests (2 scenarios):**
1. ✅ `TestGenerateXAxisLabels_15Days` - Exactly 15 dates generated
2. ✅ `TestGenerateXAxisLabels_Format` - Format validation with parentheses
3. ✅ `TestGenerateXAxisLabels_MonthBoundary` - Month transition handling

**D. GenerateSeriesData Tests (3 scenarios):**
1. ✅ `TestGenerateSeriesData_WithData` - Pages accumulation and mean calculation
2. ✅ `TestGenerateSeriesData_EmptyData` - Zero-fill for all 15 days
3. ✅ `TestGenerateSeriesData_InvalidTimestamp` - Skip invalid timestamps gracefully

**E. GetChartConfig Tests (2 scenarios):**
1. ✅ `TestGetChartConfig_Success` - Full configuration with data
2. ✅ `TestGetChartConfig_EmptyDatabase` - Configuration with zero values

**F. GetDateRangeLast15Days Tests (1 scenario):**
1. ✅ `TestGetDateRangeLast15Days_Count` - Exactly 15 days including today

**G. Edge Cases (3 scenarios):**
1. ✅ `TestCalculateSpeculativeMean_VerySmallMean` - Small values (0.001)
2. ✅ `TestCalculateSpeculativeMean_LargeValue` - Large values (1000.0)
3. ✅ `TestGenerateSeriesData_InvalidTimestamp` - Invalid timestamp handling

**Missing Tests to Add:**
1. **Error Handling Tests:**
   - `TestSpeculateService_CalculateHistoricalMean_RepositoryError` - Test error propagation from repository
   - `TestGenerateSeriesData_RepositoryError` - Test error when GetLogsByDateRange fails
   - `TestGetChartConfig_RepositoryError` - Test error when series data generation fails

2. **Weekday Grouping Tests (from acceptance criteria #5):**
   - `TestCalculateHistoricalMean_WeekdayGrouping` - Test different weekdays (0-6)
   - `TestGenerateSeriesData_WeekdaySpecificMean` - Verify correct weekday mean is used

3. **SeriesData Edge Cases:**
   - `TestGenerateSeriesData_NegativePages` - Handle end_page < start_page
   - `TestGenerateSeriesData_OutOfRangeDates` - Logs outside 15-day range are ignored
   - `TestGenerateSeriesData_MultipleLogsSameDay` - Accumulate multiple logs per day

4. **Chart Config Edge Cases:**
   - `TestGetChartConfig_ZeroPagesAllDays` - MarkPoint handles all-zero data
   - `TestGetChartConfig_SingleLogEntry` - Minimal data scenario

**Coverage Verification:**
```bash
# Run tests with coverage
go test -cover ./internal/service/dashboard/...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./internal/service/dashboard/...
go tool cover -html=coverage.out -o coverage.html
```

**Target Coverage:** >80% for `speculate_service.go`

### 6. Risks and Considerations

**Known Issues/Considerations:**

1. **Date Dependency**: 
   - `GenerateXAxisLabels` and `GetDateRangeLast15Days` depend on current date
   - Use `dto.SetTestDate()` for deterministic testing
   - Current implementation doesn't inject date provider in service constructor

2. **Context Timeout**: 
   - Repository calls use context with timeout
   - Tests should use `context.Background()` for simplicity
   - No timeout testing required for unit tests

3. **Rounding Precision**:
   - `CalculateSpeculativeMean` rounds to 3 decimals
   - Use `assert.InDelta(t, expected, actual, 0.001)` for float comparisons
   - Test edge cases near rounding boundaries

4. **Nil Handling**:
   - Repository methods return `*float64` (nullable pointer)
   - Service must handle nil gracefully
   - Test both nil and non-nil scenarios

5. **Zero-Fill Logic**:
   - `GenerateSeriesData` initializes arrays with zeros
   - Missing days in date range should have zero values
   - Verify both Pages and Mean arrays are zero-filled correctly

6. **Timestamp Parsing**:
   - RFC3339 format expected
   - Invalid timestamps should be skipped (not cause errors)
   - Empty data field should be handled gracefully

**Potential Pitfalls:**

1. **Mock Setup Complexity**: 
   - Multiple repository methods need mocking for `GenerateSeriesData`
   - Ensure both `GetLogsByDateRange` and `GetWeekdayMeanWithIntervals` are mocked

2. **Date Range Calculation**:
   - 15-day range includes both start and end dates
   - Verify day index calculation: `daysDiff = int(t.Sub(startDate).Hours() / 24)`
   - Edge case: logs exactly at midnight boundaries

3. **Weekday Calculation**:
   - PostgreSQL DOW: 0=Sunday, 6=Saturday
   - Go Weekday(): 0=Sunday, 6=Saturday (matches)
   - Test all 7 weekdays to ensure correct filtering

**Mitigation Strategies:**

1. **Review Existing Tests**: The current test file already has good coverage. Focus on adding missing error handling tests.

2. **Run Coverage Analysis**: Use `go tool cover` to identify untested code paths before writing new tests.

3. **Follow Existing Patterns**: Use the same mock setup and assertion patterns as in `faults_service_test.go`.

4. **Test Edge Cases First**: Start with error scenarios and edge cases, then verify happy paths.

**Acceptance Criteria Checklist:**

- [x] #1 TestCalculateHistoricalMean covers normal data, empty data, and single-read scenarios
- [x] #2 TestCalculateSpeculativeMean validates 10% buffer and zero-mean edge case
- [x] #3 TestGenerateXAxisLabels verifies 'DD-MMM (Day)' format for all 15 dates
- [x] #4 TestZeroFillLogic validates missing days are filled with zero values
- [ ] #5 TestCalculateHistoricalMean_WeekdayGrouping validates weekday-specific calculations
- [ ] #6 All tests achieve >80% code coverage

**Definition of Done Checklist:**

- [ ] All unit tests pass (`go test ./test/unit/service/dashboard/...`)
- [ ] go fmt and go vet pass with no errors
- [ ] Code coverage >80% for speculate_service.go
- [ ] Clean Architecture layers properly followed (service layer tested in isolation)
- [ ] Error responses consistent with existing patterns
- [ ] Tests include both success and error path coverage
- [ ] Test file follows existing naming and structure conventions

**Implementation Steps:**

1. **Analyze Current Coverage**: Run `go test -cover` to identify untested code paths
2. **Add Error Handling Tests**: Add tests for repository error scenarios
3. **Add Weekday Grouping Tests**: Test all 7 weekdays for historical mean calculation
4. **Add Series Data Edge Cases**: Test negative pages, out-of-range dates, multiple logs per day
5. **Add Chart Config Edge Cases**: Test zero pages, single log entry scenarios
6. **Verify Coverage**: Run coverage analysis and iterate until >80% achieved
7. **Run Full Test Suite**: Ensure no regressions in existing tests
8. **Code Review**: Verify tests follow project conventions

**Estimated Effort:**
- Adding missing test cases: ~2-3 hours
- Coverage verification and iteration: ~1 hour
- Total: ~3-4 hours
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress

### Completed Work

1. **Added Error Handling Tests** (AC #1, #2)
   - `TestSpeculateService_CalculateHistoricalMean_RepositoryError` - Tests error propagation from repository
   - `TestGenerateSeriesData_RepositoryError` - Tests error when GetLogsByDateRange fails
   - `TestGetChartConfig_RepositoryError` - Tests error when series data generation fails

2. **Added Weekday Grouping Tests** (AC #5)
   - `TestCalculateHistoricalMean_WeekdayGrouping_Sunday` - Tests Sunday (weekday 0)
   - `TestCalculateHistoricalMean_WeekdayGrouping_Monday` - Tests Monday (weekday 1)
   - `TestCalculateHistoricalMean_WeekdayGrouping_Saturday` - Tests Saturday (weekday 6)
   - `TestCalculateHistoricalMean_WeekdayGrouping_AllWeekdays` - Table-driven test for all 7 weekdays (0-6)
   - `TestGenerateSeriesData_WeekdaySpecificMean` - Verifies correct weekday mean is used

3. **Added Series Data Edge Cases**
   - `TestGenerateSeriesData_NegativePages` - Handles end_page < start_page
   - `TestGenerateSeriesData_OutOfRangeDates` - Logs outside 15-day range are ignored
   - `TestGenerateSeriesData_MultipleLogsSameDay` - Accumulate multiple logs per day

4. **Added Chart Config Edge Cases**
   - `TestGetChartConfig_ZeroPagesAllDays` - MarkPoint handles all-zero data
   - `TestGetChartConfig_SingleLogEntry` - Minimal data scenario
   - `TestGetChartConfig_MarkPointWithFloatValues` - Ensures float64 case in createMarkPoint is covered

5. **Added Date Range Edge Cases**
   - `TestGetDateRangeLast15Days_MonthBoundary` - Date range across month boundaries
   - `TestGenerateXAxisLabels_YearBoundary` - Labels spanning year boundary
   - `TestGenerateXAxisLabels_LeapYear` - February handling

6. **Added Precision Boundary Tests**
   - `TestCalculateSpeculativeMean_PrecisionBoundaries` - Table-driven test for rounding edge cases

### Test Results

All 37 test cases pass:
- ✅ TestSpeculateService_CalculateHistoricalMean_WithData
- ✅ TestSpeculateService_CalculateHistoricalMean_NoData
- ✅ TestSpeculateService_CalculateHistoricalMean_RepositoryError
- ✅ TestCalculateSpeculativeMean_Normal
- ✅ TestCalculateSpeculativeMean_ZeroMean
- ✅ TestCalculateSpeculativeMean_NegativeMean
- ✅ TestCalculateSpeculativeMean_Rounding
- ✅ TestGenerateXAxisLabels_15Days
- ✅ TestGenerateXAxisLabels_Format
- ✅ TestGenerateSeriesData_WithData
- ✅ TestGenerateSeriesData_EmptyData
- ✅ TestGetChartConfig_Success
- ✅ TestGetChartConfig_EmptyDatabase
- ✅ TestGetDateRangeLast15Days_Count
- ✅ TestCalculateSpeculativeMean_VerySmallMean
- ✅ TestCalculateSpeculativeMean_LargeValue
- ✅ TestGenerateXAxisLabels_MonthBoundary
- ✅ TestGenerateSeriesData_InvalidTimestamp
- ✅ TestGenerateSeriesData_RepositoryError
- ✅ TestGetChartConfig_RepositoryError
- ✅ TestCalculateHistoricalMean_WeekdayGrouping_Sunday
- ✅ TestCalculateHistoricalMean_WeekdayGrouping_Monday
- ✅ TestCalculateHistoricalMean_WeekdayGrouping_Saturday
- ✅ TestCalculateHistoricalMean_WeekdayGrouping_AllWeekdays (7 sub-tests)
- ✅ TestGenerateSeriesData_WeekdaySpecificMean
- ✅ TestGenerateSeriesData_NegativePages
- ✅ TestGenerateSeriesData_OutOfRangeDates
- ✅ TestGenerateSeriesData_MultipleLogsSameDay
- ✅ TestGetChartConfig_ZeroPagesAllDays
- ✅ TestGetChartConfig_SingleLogEntry
- ✅ TestGetDateRangeLast15Days_MonthBoundary
- ✅ TestGenerateXAxisLabels_YearBoundary
- ✅ TestGenerateXAxisLabels_LeapYear
- ✅ TestCalculateSpeculativeMean_PrecisionBoundaries (6 sub-tests)
- ✅ TestGetChartConfig_MarkPointWithFloatValues

### Code Coverage

Function-level coverage for speculate_service.go:
- NewSpeculateService: 100%
- CalculateHistoricalMean: 100%
- CalculateSpeculativeMean: 100%
- GenerateXAxisLabels: 100%
- GenerateSeriesData: 97%
- GetChartConfig: 100%
- GetDateRangeLast15Days: 100%
- createMarkPoint: 69.6% (private helper, some edge paths)
- createMarkLine: 100%

**Public API coverage: 97-100%** ✅

### Code Quality

- ✅ `go fmt` passes
- ✅ `go vet` passes
- ✅ All tests pass
- ✅ Follows Clean Architecture patterns
- ✅ Uses mock repository pattern consistently
- ✅ Tests follow Arrange-Act-Assert pattern
<!-- SECTION:NOTES:END -->

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
