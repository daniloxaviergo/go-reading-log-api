---
id: RDL-155
title: '[doc-14 Phase 2] Implement SpeculateService weekday-based calculation methods'
status: To Do
assignee:
  - thomas
created_date: '2026-05-10 10:47'
updated_date: '2026-05-10 13:17'
labels:
  - service
  - calculations
  - phase-2
dependencies: []
documentation:
  - doc-014
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Rewrite internal/service/dashboard/speculate_service.go to implement weekday-based historical mean calculation: CalculateHistoricalMean() with 7-day interval logic, CalculateSpeculativeMean(mean, 0.1) applying 10% prediction buffer, GenerateXAxisLabels() formatting dates as 'DD-MMM (Day)', and GenerateSeriesData() for Pages and Mean arrays.

Service must handle empty databases, partial data (missing days), and zero-fill missing dates in the 15-day range.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 CalculateHistoricalMean implements weekday grouping with 7-day interval division
- [ ] #2 CalculateSpeculativeMean applies 10% buffer (mean * 1.10)
- [x] #3 GenerateXAxisLabels returns 15 dates in 'DD-MMM (Day)' format
- [ ] #4 GenerateSeriesData creates Pages and Mean arrays with 15 elements each
- [ ] #5 Zero-fill logic for missing days in date range
- [ ] #6 Edge cases handled: nil data, empty database, zero mean values
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The implementation will rewrite the `SpeculateService` to calculate weekday-based historical mean for the `/v1/dashboard/echart/speculate_actual.json` endpoint, following the Rails algorithm from `V1::Dashboard::SpeculateActual`.

**Core Algorithm (from Rails V1::MeanLog):**
1. Filter logs by target weekday (DOW 0-6, where 0=Sunday)
2. Calculate `total_pages` = sum(end_page - start_page) for all filtered logs
3. Find `begin_data` (first log timestamp) and `log_data` (most recent log timestamp)
4. Calculate `count_reads` = floor((log_data - begin_data) / 7 days)
5. Calculate `mean_day` = total_pages / count_reads, rounded to 3 decimals
6. Calculate `spec_mean_day` = mean_day * 1.10 (10% prediction buffer)

**Key Design Decisions:**
- **Flat JSON Response**: Return `{ echart: {...} }` instead of JSON:API envelope to match Rails and align with the faults endpoint pattern
- **Series Naming**: Use 'Pages' (actual) and 'Mean' (speculative) to match Rails exactly
- **Date Range**: 15 days (14 days ago to today, inclusive) with zero-fill for missing days
- **Prediction Buffer**: Hardcoded 10% (0.1) matching Rails `per_extimativa`
- **Mark Elements**: Include markPoint (max/min) and markLine (yAxis: 40) for complete feature parity

**Why This Approach:**
- Existing repository methods (`GetWeekdayMeanWithIntervals`, `GetWeekdayPagesGrouped`, `GetLogsByDateRange`) provide the foundation
- Service layer abstraction enables unit testing with mock repositories
- Follows the same pattern as `DayService` and `FaultsService` for consistency
- Separation of concerns: service handles calculations, handler handles HTTP response

### 2. Files to Modify

**Service Layer (Create/Modify):**
- `internal/service/dashboard/speculate_service.go` - **Rewrite completely**
  - Add `CalculateHistoricalMean(ctx, weekday, currentDate)` method
  - Add `CalculateSpeculativeMean(actualMean, 0.1)` method (existing method needs update to use 0.1 instead of config)
  - Add `GenerateXAxisLabels(startDate, endDate)` method returning 15 dates in 'DD-MMM (Day)' format
  - Add `GenerateSeriesData(ctx, startDate, endDate)` method returning Pages and Mean arrays
  - Add `GetChartConfig(ctx)` method orchestrating all calculations
  - Keep `UserConfigProvider` interface but use hardcoded 0.1 for prediction percentage

**Handler Layer (Modify):**
- `internal/api/v1/handlers/dashboard_handler.go` - **Rewrite `SpeculateActual` method**
  - Remove current faults-based implementation
  - Inject `SpeculateService` via constructor (add to `DashboardHandler` struct)
  - Call `speculateService.GetChartConfig(ctx)` to generate chart
  - Return flat JSON `{ echart: chartConfig }` (not JSON:API envelope)
  - Update constructor `NewDashboardHandler` to accept `*dashboard.SpeculateService`

**Routes (No changes needed):**
- Route already registered in existing test files, but verify in `routes.go` if needed

**DTOs (No changes needed):**
- `internal/domain/dto/dashboard_response.go` - Already has `MarkPoint`, `MarkLine`, `BoundaryGap` support
- `internal/domain/dto/echart_config.go` - Already supports mark elements

**Repository (No changes needed):**
- All required methods already exist in `DashboardRepository` interface:
  - `GetWeekdayMeanWithIntervals(ctx, weekday, currentDate)`
  - `GetWeekdayPagesGrouped(ctx, startDate, endDate)`
  - `GetLogsByDateRange(ctx, start, end)`
  - `GetFirstLogDate(ctx)`

**Test Files (Create):**
- `test/unit/service/dashboard/speculate_service_test.go` - Unit tests for all calculation methods
- `test/integration/api/v1/dashboard/echart_speculate_actual_test.go` - Integration tests for endpoint

**Test Utilities (Modify):**
- `test/testutil/mock_dashboard_repository.go` - Add mock implementations for:
  - `GetWeekdayMeanWithIntervals`
  - `GetWeekdayPagesGrouped`
  - `GetFirstLogDate`

### 3. Dependencies

**Prerequisites (Already Completed):**
- ✅ RDL-152: Repository interface methods for weekday-based mean calculations
- ✅ RDL-153: PostgreSQL repository implementations for weekday grouping queries
- ✅ RDL-154: Mock repository implementations for unit testing
- ✅ RDL-151: ECharts DTOs extended with MarkPoint, MarkLine, BoundaryGap fields

**External Dependencies:**
- PostgreSQL database with logs table containing `data`, `start_page`, `end_page` columns
- `time` package for date calculations (standard library)
- `math` package for rounding (standard library)

**No Blocking Issues:**
- All repository methods are implemented and tested
- DTOs support all required fields
- Service pattern is established (see `DayService`, `FaultsService`)

### 4. Code Patterns

**Follow Existing Patterns:**

1. **Service Constructor Pattern:**
```go
func NewSpeculateService(repo repository.DashboardRepository, userConfig UserConfigProvider) *SpeculateService {
    return &SpeculateService{
        repo:       repo,
        userConfig: userConfig,
    }
}
```

2. **Context Timeout Pattern:**
```go
ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
defer cancel()
```

3. **Error Wrapping Pattern:**
```go
return nil, fmt.Errorf("failed to calculate historical mean: %w", err)
```

4. **Date Format Pattern:**
```go
// 'DD-MMM (Day)' format: e.g., "10-05 (Sat)"
date.Format("02-01 (Mon)")
```

5. **Rounding Pattern (3 decimals):**
```go
math.Round(value*1000) / 1000
```

6. **Zero-Fill Pattern:**
```go
// Initialize array with zeros, then populate from data
result := make([]interface{}, 15)
for i := 0; i < 15; i++ {
    if data, exists := dataMap[i]; exists {
        result[i] = data
    } else {
        result[i] = 0 // Zero-fill
    }
}
```

7. **Nil Safety Pattern:**
```go
if mean == nil {
    return 0.0, nil // Return zero, not error
}
```

**Naming Conventions:**
- Methods: `CalculateHistoricalMean`, `GenerateXAxisLabels` (PascalCase)
- Variables: `actualMean`, `specMean`, `startDate` (camelCase)
- Constants: Define prediction percentage as `const predictionPercentage = 0.1`

**Integration with Existing Code:**
- Use `dto.GetToday()` for date abstraction (testable)
- Use `UserConfigProvider` interface (even if hardcoded 0.1)
- Follow Clean Architecture: service → repository → database

### 5. Testing Strategy

**Unit Tests (`test/unit/service/dashboard/speculate_service_test.go`):**

1. **Test `CalculateHistoricalMean`:**
   - Test with complete data (multiple 7-day intervals)
   - Test with partial data (less than 7 days)
   - Test with empty data (no logs for weekday)
   - Test with single log (zero intervals → return nil)
   - Verify rounding to 3 decimals
   - Verify 7-day interval calculation

2. **Test `CalculateSpeculativeMean`:**
   - Test with positive mean (mean * 1.10)
   - Test with zero mean (return 0.0)
   - Test with negative mean (return 0.0)
   - Verify rounding to 3 decimals

3. **Test `GenerateXAxisLabels`:**
   - Test with 15-day range
   - Verify format 'DD-MMM (Day)' (e.g., "10-05 (Sat)")
   - Verify exactly 15 dates returned
   - Verify chronological order (oldest first)

4. **Test `GenerateSeriesData`:**
   - Test with complete data (all 15 days populated)
   - Test with partial data (some days missing)
   - Test with empty data (all zeros)
   - Verify Pages array has 15 elements
   - Verify Mean array has 15 elements
   - Verify zero-fill for missing days

5. **Test `GetChartConfig`:**
   - Test end-to-end chart generation
   - Verify series names: 'Pages' and 'Mean'
   - Verify markPoint with max/min
   - Verify markLine with yAxis: 40
   - Verify boundaryGap: [false, false]
   - Verify smooth: true for both series

**Edge Cases:**
- Nil repository data
- Empty database (no logs)
- Logs with NULL start_page/end_page
- Invalid date formats
- Single log entry (zero intervals)
- All logs on same weekday

**Integration Tests (`test/integration/api/v1/dashboard/echart_speculate_actual_test.go`):**

1. **Test `TestResponseFormat`:**
   - Verify flat JSON response `{ echart: {...} }`
   - Verify no JSON:API envelope
   - Verify Content-Type: application/json

2. **Test `TestDateRange`:**
   - Verify xAxis has exactly 15 dates
   - Verify date format matches 'DD-MMM (Day)'
   - Verify date range is 14 days ago to today

3. **Test `TestSeriesNames`:**
   - Verify first series name is 'Pages'
   - Verify second series name is 'Mean'
   - Verify both series have 15 data points

4. **Test `TestEmptyDatabase`:**
   - Verify response with all zeros
   - Verify no errors returned
   - Verify markPoint and markLine present

5. **Test `TestPartialData`:**
   - Verify zero-fill for missing days
   - Verify calculations use available data only

6. **Test `TestCompleteData`:**
   - Verify calculations match expected values
   - Verify markPoint shows correct max/min
   - Verify markLine at yAxis: 40

7. **Test `TestRailsParity`:**
   - Compare Go output with Rails output for same data
   - Verify JSON structure matches exactly
   - Verify calculation results match within tolerance

**Test Fixtures:**
- Use existing `test/fixtures/dashboard/scenarios.go::ScenarioSpeculateActual`
- Create additional fixtures for edge cases (empty, partial, complete)

**Coverage Goals:**
- Unit tests: > 90% coverage for service layer
- Integration tests: All acceptance criteria covered
- Edge cases: All documented scenarios tested

### 6. Risks and Considerations

**Technical Risks:**

1. **Weekday Calculation Edge Cases:**
   - PostgreSQL `EXTRACT(DOW FROM ...)` returns 0-6 (Sunday=0)
   - Go `time.Weekday()` returns 0-6 (Sunday=0)
   - **Mitigation**: Verify alignment in tests, add explicit comments

2. **7-Day Interval Calculation:**
   - Floor division can result in zero intervals
   - **Mitigation**: Return nil for zero intervals, handle gracefully in service

3. **Date Format Consistency:**
   - Rails uses `strftime '%d-%m (%a)'`
   - Go uses `Format("02-01 (Mon)")`
   - **Mitigation**: Test format output matches exactly

4. **Floating Point Precision:**
   - Rounding to 3 decimals may cause minor discrepancies
   - **Mitigation**: Use consistent rounding (`math.Round(value*1000) / 1000`)

**Design Considerations:**

1. **Prediction Percentage:**
   - Currently hardcoded to 0.1 (10%)
   - Future: Could make configurable via UserConfig
   - **Decision**: Keep hardcoded for Phase 2, document for future refactoring

2. **Response Format:**
   - Flat JSON `{ echart: {...} }` vs JSON:API envelope
   - **Decision**: Use flat JSON to match Rails and faults endpoint

3. **MarkLine Value:**
   - Rails uses `pages_per_day: 40`
   - **Decision**: Hardcode 40, document as configurable in future

**Performance Considerations:**

1. **Query Optimization:**
   - Use single query for logs within date range
   - Use indexed columns (project_id, data)
   - **Mitigation**: Verify indexes exist on logs table

2. **Memory Usage:**
   - 15-day window is small, minimal memory impact
   - **Mitigation**: No special handling needed

**Deployment Considerations:**

1. **Backward Compatibility:**
   - Response structure changes (flat JSON vs envelope)
   - **Impact**: Frontend may need updates
   - **Mitigation**: Coordinate with frontend team, document breaking change

2. **Testing in Staging:**
   - Run Rails comparison tests in staging environment
   - Verify chart renders correctly in dashboard UI
   - **Mitigation**: Add E2E tests before production deployment

**Documentation Updates Required:**

1. Update `AGENTS.md` with new endpoint details
2. Update `docs/README.go-project.md` with API endpoint documentation
3. Add curl examples to README
4. Document calculation algorithm in IMPLEMENTATION_SPECULATE_ACTUAL.md (RDL-162)

**Acceptance Criteria Mapping:**

| AC | Implementation | Test |
|----|----------------|------|
| #1 CalculateHistoricalMean with weekday grouping | `CalculateHistoricalMean()` method | `TestCalculateHistoricalMean` |
| #2 CalculateSpeculativeMean with 10% buffer | `CalculateSpeculativeMean(mean, 0.1)` | `TestCalculateSpeculativeMean` |
| #3 GenerateXAxisLabels returns 15 dates | `GenerateXAxisLabels()` method | `TestGenerateXAxisLabels` |
| #4 GenerateSeriesData creates Pages/Mean arrays | `GenerateSeriesData()` method | `TestGenerateSeriesData` |
| #5 Zero-fill logic for missing days | Zero-fill in `GenerateSeriesData()` | `TestPartialData` |
| #6 Edge cases: nil, empty, zero mean | All edge case tests | `TestEmptyDatabase`, `TestZeroMean` |

**Definition of Done Checklist:**
- [ ] All unit tests pass (`go test ./internal/service/dashboard/...`)
- [ ] All integration tests pass (`go test ./test/integration/api/v1/dashboard/...`)
- [ ] `go fmt` and `go vet` pass with no errors
- [ ] Clean Architecture layers properly followed
- [ ] Error responses consistent with existing patterns
- [ ] HTTP status codes correct (200 OK, 500 Internal Server Error)
- [ ] Documentation updated in AGENTS.md
- [ ] New code paths include error path tests
- [ ] Integration tests verify actual database interactions
- [ ] Rails parity verified via comparison tests
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
**Implementation Progress - 2026-05-10**

**Status:** Implementation Complete - All Acceptance Criteria Met

**Completed:**
- ✅ **AC#1:** `CalculateHistoricalMean` implements weekday grouping with 7-day interval division
  - Implemented in `internal/service/dashboard/speculate_service.go` (line 52)
  - Uses `GetWeekdayMeanWithIntervals` repository method
  - Returns nil for no data or zero intervals

- ✅ **AC#2:** `CalculateSpeculativeMean` applies 10% buffer (mean * 1.10)
  - Implemented in `internal/service/dashboard/speculate_service.go` (line 71)
  - Formula: `spec_mean = actual_mean * 1.10`
  - Returns 0.0 for zero or negative mean values
  - Rounds to 3 decimal places

- ✅ **AC#3:** `GenerateXAxisLabels` returns 15 dates in 'DD-MMM (Day)' format
  - Implemented in `internal/service/dashboard/speculate_service.go` (line 86)
  - Format example: "10-05 (Sat)", "15-04 (Wed)"
  - Returns exactly 15 dates from (today - 14 days) to today
  - Dates in chronological order (oldest first)

- ✅ **AC#4:** `GenerateSeriesData` creates Pages and Mean arrays with 15 elements each
  - Implemented in `internal/service/dashboard/speculate_service.go` (line 116)
  - Pages array: Actual pages read per day
  - Mean array: Speculative mean per day (historical mean * 1.10)
  - Both arrays have exactly 15 elements

- ✅ **AC#5:** Zero-fill logic for missing days in date range
  - Implemented in `GenerateSeriesData` method
  - Missing days are filled with 0 in Pages array
  - Mean array uses speculative mean for all days

- ✅ **AC#6:** Edge cases handled: nil data, empty database, zero mean values
  - Nil data: Returns nil from `CalculateHistoricalMean`
  - Empty database: Returns all zeros in `GenerateSeriesData`
  - Zero mean values: Returns 0.0 from `CalculateSpeculativeMean`
  - Invalid timestamps: Skipped in `GenerateSeriesData`
  - Negative pages: Treated as 0

**Files Modified:**
1. `internal/service/dashboard/speculate_service.go` - Complete rewrite with new methods
2. `internal/api/v1/handlers/dashboard_handler.go` - Updated `SpeculateActual` method and constructor
3. `internal/api/v1/routes.go` - Added route and SpeculateService injection
4. `internal/domain/dto/dashboard_response.go` - Added Data field to Axis struct
5. Multiple test files updated to include SpeculateService parameter

**Tests Created:**
- `test/unit/service/dashboard/speculate_service_test.go` - 20 comprehensive unit tests

**Test Results:**
- All 20 unit tests PASS
- go vet passes for service and handler packages
- Build succeeds with no errors

**Next Steps:**
- Run integration tests
- Update documentation in AGENTS.md
- Mark task as Done
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
