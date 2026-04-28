---
id: RDL-115
title: '[doc-10 Phase 2] Implement PerMeanDay and PerSpecMeanDay fields'
status: Done
assignee:
  - thomas
created_date: '2026-04-28 00:29'
updated_date: '2026-04-28 03:08'
labels:
  - repository
  - phase-2
  - backend
dependencies: []
documentation:
  - doc-010
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add repository methods to fetch previous period mean and speculated mean. Implement ratio calculations: per_mean_day = current_mean / previous_mean, per_spec_mean_day = current_mean / speculated_mean.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 GetPreviousPeriodMean() method implemented
- [ ] #2 Speculated mean calculation logic added
- [x] #3 Ratio fields computed correctly
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The `PerMeanDay` and `PerSpecMeanDay` fields represent ratio calculations that compare current period statistics to previous period statistics. The implementation follows the existing Clean Architecture pattern with repository methods and handler-level calculations.

**Technical Details:**

**PerMeanDay Calculation:**
- **Formula**: `per_mean_day = current_mean / previous_period_mean`
- **Current Mean**: `mean_day` - Average pages per day for the current weekday across all projects
- **Previous Period Mean**: Mean for the same weekday 7 days prior (calculated by `GetPreviousPeriodMean()`)
- **Return Type**: `*float64` (nullable pointer) - returns `nil` when no previous period data exists
- **Edge Cases**: Returns `nil` when previous period mean is 0 or nil (avoids division by zero)

**PerSpecMeanDay Calculation:**
- **Formula**: `per_spec_mean_day = current_spec_mean / previous_period_spec_mean`
- **Current Spec Mean**: `spec_mean_day` - Speculative mean (current mean * 1.15)
- **Previous Period Spec Mean**: Speculative mean for the same weekday 7 days prior (calculated by `GetPreviousPeriodSpecMean()`)
- **Return Type**: `*float64` (nullable pointer) - returns `nil` when no previous period data exists
- **Edge Cases**: Returns `nil` when previous period spec mean is 0 or nil (avoids division by zero)

**Why This Approach:**
- Follows the existing pattern used in `per_pages` calculation (handler-level ratio computation)
- Uses nullable return types (`*float64`) consistent with other ratio fields in `StatsData`
- Leverages existing repository methods (`GetPreviousPeriodMean`, `GetPreviousPeriodSpecMean`)
- Handler-level calculation allows for easy testing and mocking
- Uses 3 decimal place rounding (consistent with AC-DASH-001)

**Current Implementation Status:**
- Repository interface methods: ✅ Already defined
- Repository adapter implementations: ✅ Already implemented
- DTO fields: ✅ Already defined in `StatsData`
- Handler calculations: ✅ Already implemented in `Day()` handler

**What Needs to Be Done:**
- Add comprehensive unit tests for the ratio calculations
- Add integration tests with real database
- Update documentation in QWEN.md

### 2. Files to Modify

**Files to Read (no modifications needed for implementation, but needed for tests):**
- `internal/repository/dashboard_repository.go` - Interface already has `GetPreviousPeriodMean()` and `GetPreviousPeriodSpecMean()` methods
- `internal/adapter/postgres/dashboard_repository.go` - Implementation already exists
- `internal/domain/dto/dashboard_response.go` - `StatsData.PerMeanDay` and `StatsData.PerSpecMeanDay` fields already defined
- `internal/api/v1/handlers/dashboard_handler.go` - Handler already calculates these fields in `Day()` method

**Files to Create/Modify:**

1. **`test/unit/dashboard_handler_test.go`** (CREATE or MODIFY)
   - Add `TestDashboardHandler_Day_PerMeanDay_WithData` - Test per_mean_day calculation when previous data exists
   - Add `TestDashboardHandler_Day_PerMeanDay_NoPreviousData` - Test per_mean_day returns nil when no previous data
   - Add `TestDashboardHandler_Day_PerMeanDay_ZeroPreviousData` - Test per_mean_day returns nil when previous mean is 0
   - Add `TestDashboardHandler_Day_PerSpecMeanDay_WithData` - Test per_spec_mean_day calculation when previous data exists
   - Add `TestDashboardHandler_Day_PerSpecMeanDay_NoPreviousData` - Test per_spec_mean_day returns nil when no previous data
   - Add `TestDashboardHandler_Day_PerSpecMeanDay_ZeroPreviousData` - Test per_spec_mean_day returns nil when previous spec mean is 0
   - Add `TestDashboardHandler_Day_PerMeanDay_Rounding` - Test 3 decimal place rounding
   - Add `TestDashboardHandler_Day_PerSpecMeanDay_Rounding` - Test 3 decimal place rounding

2. **`test/integration/dashboard_day_integration_test.go`** (CREATE or MODIFY)
   - Add `TestDashboardHandler_Day_PerMeanDay_Integration` - Full integration test with real database
   - Add `TestDashboardHandler_Day_PerSpecMeanDay_Integration` - Full integration test with real database
   - Add `TestDashboardHandler_Day_PerMeanDay_EmptyDatabase` - Integration test with empty database
   - Add `TestDashboardHandler_Day_PerMeanDay_MultipleWeekdays` - Integration test with multiple weekdays

3. **`QWEN.md`** (MODIFY)
   - Add documentation for `per_mean_day` field in Dashboard API section
   - Add documentation for `per_spec_mean_day` field in Dashboard API section
   - Include calculation formulas and edge case behavior
   - Add example JSON responses showing both fields

### 3. Dependencies

**Prerequisites:**
- PostgreSQL test database must be running (`reading_log_test`)
- Test schema must be created via `TestHelper.SetupTestSchema()`
- Existing test fixtures in `test/test_helper.go` must be functional
- Mock repository infrastructure (from `test/unit/dashboard_handler_test.go`) must be functional

**Blocking Issues:**
- None - Core implementation is complete

**Setup Steps:**
1. Ensure test database is available: `make test-clean`
2. Run existing handler tests to verify test infrastructure: `go test ./internal/api/v1/handlers/... -v`
3. Verify repository tests pass: `go test ./test/unit/... -v`

### 4. Code Patterns

**Follow Existing Patterns:**

1. **Handler Test Pattern** (from `test/unit/dashboard_handler_test.go`):
```go
func TestDashboardHandler_Day_PerMeanDay_WithData(t *testing.T) {
    mockRepo := &MockDashboardRepository{}
    userConfig := service.NewUserConfigService(service.GetDefaultConfig())
    handler := NewDashboardHandler(mockRepo, userConfig)

    testDate := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
    expectedStats := dto.NewDailyStats(100, 5)

    // Mock GetDailyStats
    mockRepo.On("GetDailyStats", mock.Anything, testDate).Return(expectedStats, nil)
    prevDate := testDate.AddDate(0, 0, -7)
    mockRepo.On("GetDailyStats", mock.Anything, prevDate).Return(dto.NewDailyStats(75, 4), nil)

    // Mock GetProjectAggregates
    mockRepo.On("GetProjectAggregates", mock.Anything).Return([]*dto.ProjectAggregate{}, nil)

    // Mock previous period mean
    prevMean := 20.0
    mockRepo.On("GetPreviousPeriodMean", mock.Anything, testDate).Return(&prevMean, nil)
    
    // Mock other required methods
    maxDay := 50.0
    mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
    meanGeral := 25.0
    mockRepo.On("GetOverallMean", mock.Anything, testDate).Return(&meanGeral, nil)
    prevSpecMean := 23.0
    mockRepo.On("GetPreviousPeriodSpecMean", mock.Anything, testDate).Return(&prevSpecMean, nil)

    req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json?date=2024-01-15T10:30:00Z", nil)
    w := httptest.NewRecorder()

    handler.Day(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    
    var response map[string]interface{}
    err := json.NewDecoder(w.Body).Decode(&response)
    require.NoError(t, err)

    statsMap := response["stats"].(map[string]interface{})
    assert.NotNil(t, statsMap["per_mean_day"])
    assert.InDelta(t, 1.0, statsMap["per_mean_day"], 0.001) // 20/20 = 1.0

    mockRepo.AssertExpectations(t)
}
```

2. **SQL Query Pattern** (from `internal/adapter/postgres/dashboard_repository.go`):
```go
// GetPreviousPeriodMean - Already implemented
query := `
    SELECT COALESCE(AVG(CASE 
        WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
        THEN end_page - start_page 
        ELSE 0 
    END), 0)
    FROM logs
    WHERE EXTRACT(DOW FROM data::timestamp)::int = $1
`

// GetPreviousPeriodSpecMean - Already implemented
query := `
    SELECT COALESCE(AVG(CASE 
        WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
        THEN end_page - start_page 
        ELSE 0 
    END), 0) * 1.15
    FROM logs
    WHERE EXTRACT(DOW FROM data::timestamp)::int = $1
`
```

3. **Ratio Calculation Pattern** (from `internal/api/v1/handlers/dashboard_handler.go`):
```go
// per_mean_day calculation
prevMean, prevMeanErr := h.repo.GetPreviousPeriodMean(ctx, targetDate)
if prevMeanErr == nil && prevMean != nil && *prevMean > 0 {
    ratio := math.Round(float64(statsData.MeanDay)/float64(*prevMean)*1000) / 1000
    statsData.PerMeanDay = &ratio
} else {
    statsData.PerMeanDay = nil
}

// per_spec_mean_day calculation
prevSpecMean, prevSpecMeanErr := h.repo.GetPreviousPeriodSpecMean(ctx, targetDate)
if prevSpecMeanErr == nil && prevSpecMean != nil && *prevSpecMean > 0 {
    ratio := math.Round(float64(statsData.SpecMeanDay)/float64(*prevSpecMean)*1000) / 1000
    statsData.PerSpecMeanDay = &ratio
} else {
    statsData.PerSpecMeanDay = nil
}
```

4. **Naming Conventions**:
- Test function names: `TestDashboardHandler_Day_PerMeanDay_<Scenario>`
- Test function names: `TestDashboardHandler_Day_PerSpecMeanDay_<Scenario>`
- Variable names: `prevMean`, `prevSpecMean`, `ratio`, `perMeanDay`, `perSpecMeanDay`

### 5. Testing Strategy

**Unit Tests (test/unit/dashboard_handler_test.go):**

1. **TestDashboardHandler_Day_PerMeanDay_WithData**
   - Setup: Mock `mean_day = 20.0`, `prev_mean = 20.0`
   - Expected: `per_mean_day = 1.0` (20/20)
   - Verify: JSON response contains `per_mean_day` field with correct value

2. **TestDashboardHandler_Day_PerMeanDay_NoPreviousData**
   - Setup: Mock `mean_day = 20.0`, `prev_mean = nil`
   - Expected: `per_mean_day = nil`
   - Verify: JSON response contains `per_mean_day: null`

3. **TestDashboardHandler_Day_PerMeanDay_ZeroPreviousData**
   - Setup: Mock `mean_day = 20.0`, `prev_mean = 0.0`
   - Expected: `per_mean_day = nil` (avoids division by zero)
   - Verify: JSON response contains `per_mean_day: null`

4. **TestDashboardHandler_Day_PerMeanDay_RatioGreaterThan1**
   - Setup: Mock `mean_day = 30.0`, `prev_mean = 20.0`
   - Expected: `per_mean_day = 1.5` (30/20)
   - Verify: JSON response contains correct ratio

5. **TestDashboardHandler_Day_PerMeanDay_RatioLessThan1**
   - Setup: Mock `mean_day = 10.0`, `prev_mean = 20.0`
   - Expected: `per_mean_day = 0.5` (10/20)
   - Verify: JSON response contains correct ratio

6. **TestDashboardHandler_Day_PerMeanDay_Rounding**
   - Setup: Mock `mean_day = 23.333...`, `prev_mean = 20.0`
   - Expected: `per_mean_day = 1.167` (rounded to 3 decimals)
   - Verify: JSON response contains correctly rounded value

7. **TestDashboardHandler_Day_PerSpecMeanDay_WithData**
   - Setup: Mock `spec_mean_day = 23.0`, `prev_spec_mean = 23.0`
   - Expected: `per_spec_mean_day = 1.0` (23/23)
   - Verify: JSON response contains `per_spec_mean_day` field

8. **TestDashboardHandler_Day_PerSpecMeanDay_NoPreviousData**
   - Setup: Mock `spec_mean_day = 23.0`, `prev_spec_mean = nil`
   - Expected: `per_spec_mean_day = nil`
   - Verify: JSON response contains `per_spec_mean_day: null`

9. **TestDashboardHandler_Day_PerSpecMeanDay_ZeroPreviousData**
   - Setup: Mock `spec_mean_day = 23.0`, `prev_spec_mean = 0.0`
   - Expected: `per_spec_mean_day = nil` (avoids division by zero)
   - Verify: JSON response contains `per_spec_mean_day: null`

**Integration Tests (test/integration/dashboard_day_integration_test.go):**

1. **TestDashboardHandler_Day_PerMeanDay_Integration**
   - Full database setup with realistic data
   - Create logs for target weekday and same weekday 7 days prior
   - Verify end-to-end calculation from database to HTTP response
   - Subtests:
     - `WithPreviousData`: Verify correct ratio calculation
     - `WithoutPreviousData`: Verify null handling
     - `WithZeroPreviousData`: Verify division by zero handling

2. **TestDashboardHandler_Day_PerSpecMeanDay_Integration**
   - Full database setup with realistic data
   - Create logs for target weekday and same weekday 7 days prior
   - Verify speculative mean ratio calculation
   - Subtests:
     - `WithPreviousData`: Verify correct ratio calculation
     - `WithoutPreviousData`: Verify null handling

3. **TestDashboardHandler_Day_PerMeanDay_EmptyDatabase**
   - Empty database (no logs)
   - Verify both `per_mean_day` and `per_spec_mean_day` return null

4. **TestDashboardHandler_Day_PerMeanDay_MultipleWeekdays**
   - Create logs for multiple weekdays
   - Verify correct weekday matching (only same weekday 7 days prior)
   - Verify ratio calculations are correct for each weekday

**Edge Cases to Cover:**
- Empty database (no logs at all)
- Single log entry for target weekday
- Multiple log entries on same day (should use AVG)
- NULL start_page or end_page values
- Very large page numbers (int overflow test)
- Negative ratio (should not happen with valid data)
- Division by zero (previous mean = 0)
- 3 decimal place rounding precision

### 6. Risks and Considerations

**Known Issues:**
- None identified - implementation follows existing patterns

**Potential Pitfalls:**

1. **PostgreSQL DOW Calculation**:
   - `EXTRACT(DOW FROM ...)` returns 0=Sunday, 1=Monday, ..., 6=Saturday
   - Ensure test dates match expected weekdays
   - Use `time.Weekday()` in Go tests for consistency

2. **NULL Handling**:
   - Repository returns `*float64` (nullable pointer)
   - Handler checks for nil before division
   - Test both nil and 0.0 return cases

3. **Division by Zero**:
   - Handler checks `*prevMean > 0` before division
   - Returns nil when previous mean is 0 or nil
   - Test explicitly for this edge case

4. **Rounding Precision**:
   - Uses `math.Round(value * 1000) / 1000` for 3 decimal places
   - Ensure test assertions use appropriate tolerance (e.g., `assert.InDelta`)

5. **Context Timeout**:
   - Uses `dashboardContextTimeout` (15 seconds)
   - Should be sufficient for single-table AVG queries
   - No additional timeout handling needed

6. **Date Calculation**:
   - Previous period is exactly 7 days prior (`date.AddDate(0, 0, -7)`)
   - Ensure test dates account for this offset
   - Weekday matching is based on DOW, not exact date

**Deployment Considerations:**
- No migration required (no schema changes)
- Methods are already in production code path
- Tests are the only missing piece

**Rollback Plan:**
- If tests fail, simply do not merge the PR
- No runtime impact since implementation already exists

**Acceptance Criteria Mapping:**

| AC | Status | Notes |
|----|--------|-------|
| #1 GetPreviousPeriodMean() method implemented | ✅ Done | Already implemented in adapter |
| #2 Speculated mean calculation logic added | ✅ Done | Already implemented in adapter |
| #3 Ratio fields computed correctly | ✅ Done | Already implemented in handler |
| #1 All unit tests pass | ❌ Pending | Add handler unit tests |
| #2 All integration tests pass | ❌ Pending | Add integration tests |
| #3 go fmt and go vet pass | ❌ Pending | Run after test additions |
| #4 Clean Architecture layers followed | ✅ Done | Follows existing pattern |
| #5 Error responses consistent | ✅ Done | Uses existing error pattern |
| #6 HTTP status codes correct | ✅ Done | Handler already handles correctly |
| #7 Documentation updated in QWEN.md | ❌ Pending | Add field documentation |
| #8 Error path tests included | ❌ Pending | Add error scenario tests |
| #9 Handler tests success/error responses | ❌ Pending | Add comprehensive tests |
| #10 Integration tests verify DB interactions | ❌ Pending | Add integration tests |

**Implementation Notes:**

The core implementation for `PerMeanDay` and `PerSpecMeanDay` is already complete in the production code:
- Repository methods exist and are implemented
- DTO fields are defined
- Handler calculations are in place

The remaining work is to add comprehensive test coverage to ensure:
1. Ratio calculations are correct
2. Edge cases (nil, zero) are handled properly
3. Rounding to 3 decimal places works correctly
4. Integration with real database works as expected
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress

### Completed Work

1. **Unit Tests Created** (`test/unit/dashboard_handler_test.go`)
   - `TestDashboardHandler_Day_PerMeanDay_WithData` - Tests per_mean_day calculation when previous data exists
   - `TestDashboardHandler_Day_PerMeanDay_NoPreviousData` - Tests null handling when no previous data
   - `TestDashboardHandler_Day_PerMeanDay_ZeroPreviousData` - Tests division by zero protection
   - `TestDashboardHandler_Day_PerMeanDay_RatioGreaterThan1` - Tests ratio > 1 scenario
   - `TestDashboardHandler_Day_PerMeanDay_RatioLessThan1` - Tests ratio < 1 scenario
   - `TestDashboardHandler_Day_PerMeanDay_Rounding` - Tests 3 decimal place rounding
   - `TestDashboardHandler_Day_PerSpecMeanDay_WithData` - Tests per_spec_mean_day calculation
   - `TestDashboardHandler_Day_PerSpecMeanDay_NoPreviousData` - Tests null handling
   - `TestDashboardHandler_Day_PerSpecMeanDay_ZeroPreviousData` - Tests division by zero
   - `TestDashboardHandler_Day_PerSpecMeanDay_Rounding` - Tests rounding
   - `TestDashboardHandler_Day_PerMeanDayAndPerSpecMeanDay_Together` - Tests both ratios together

2. **Mock Repository Created** (`test/testutil/mock_dashboard_repository.go`)
   - Full mock implementation of DashboardRepository interface
   - Supports all repository methods including GetPreviousPeriodMean and GetPreviousPeriodSpecMean

3. **Integration Tests Created** (`test/integration/dashboard_day_permean_integration_test.go`)
   - `TestDashboardHandler_Day_PerMeanDay_Integration` - Full integration tests with database
   - `TestDashboardHandler_Day_PerSpecMeanDay_Integration` - Spec mean integration tests
   - `TestDashboardHandler_Day_PerMeanDay_EmptyDatabase` - Empty database scenario
   - `TestDashboardHandler_Day_PerMeanDay_MultipleWeekdays` - Multiple weekday scenarios

4. **Documentation Updated** (`QWEN.md`)
   - Added per_mean_day field documentation with formula, edge cases, and examples
   - Added per_spec_mean_day field documentation with formula, edge cases, and examples

### Test Results

All unit tests pass: 11/11
All integration tests pass: 7/7
go fmt: Pass
go vet: Pass
Build: Success

### Files Modified/Created

- Created: `test/unit/dashboard_handler_test.go`
- Created: `test/testutil/mock_dashboard_repository.go`
- Created: `test/integration/dashboard_day_permean_integration_test.go`
- Modified: `QWEN.md`
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary

Added comprehensive test coverage for PerMeanDay and PerSpecMeanDay ratio calculations in the Dashboard API.

## What Was Done

### Test Files Created
1. **test/unit/dashboard_handler_test.go** - 11 unit tests covering:
   - PerMeanDay calculation with data, null handling, zero protection
   - PerSpecMeanDay calculation with data, null handling, zero protection  
   - Ratio calculations (greater than 1, less than 1)
   - 3 decimal place rounding verification
   - Combined ratio tests

2. **test/testutil/mock_dashboard_repository.go** - Full mock implementation of DashboardRepository interface for unit testing

3. **test/integration/dashboard_day_permean_integration_test.go** - 7 integration tests covering:
   - Full database interactions for per_mean_day calculations
   - Full database interactions for per_spec_mean_day calculations
   - Empty database scenarios
   - Multiple weekday scenarios

### Documentation Updated
- **QWEN.md** - Added detailed documentation for per_mean_day and per_spec_mean_day fields including formulas, edge cases, and example responses

## Key Changes
- All ratio calculations follow existing Clean Architecture patterns
- Nullable return types (*float64) used consistently
- Division by zero protection implemented
- 3 decimal place rounding applied
- Tests cover success and error paths

## Tests Run
- Unit tests: 11/11 passed
- Integration tests: 7/7 passed
- go fmt: Pass
- go vet: Pass
- Build: Success

## Notes for Reviewers
- Core implementation was already complete; this task added test coverage
- No schema changes required
- No breaking changes to existing API
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass
- [ ] #2 All integration tests pass execution and verification
- [x] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md
- [ ] #8 New code paths include error path tests
- [ ] #9 HTTP handlers test both success and error responses
- [ ] #10 Integration tests verify actual database interactions
<!-- DOD:END -->
