---
id: RDL-114
title: '[doc-10 Phase 2] Implement MaxDay field and repository method'
status: To Do
assignee:
  - thomas
created_date: '2026-04-28 00:29'
updated_date: '2026-04-28 02:27'
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
Add GetMaxByWeekday() repository method in dashboard_repository.go to query maximum pages read for a specific weekday. Implementation: max(pages_read_on_each_occurrence_of_weekday).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 GetMaxByWeekday() method implemented in adapter
- [x] #2 Interface method added to repository contract
- [ ] #3 Returns maximum pages for target weekday
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The `GetMaxByWeekday` method calculates the maximum pages read in a single day for a specific weekday across all projects. The implementation follows the existing Clean Architecture pattern:

**Technical Details:**
- **SQL Query**: Uses `MAX(CASE WHEN start_page IS NOT NULL AND end_page IS NOT NULL THEN end_page - start_page ELSE 0 END)` to find the maximum pages read in a single log entry for the target weekday
- **Weekday Calculation**: Uses `EXTRACT(DOW FROM data::timestamp)::int` where 0=Sunday, 1=Monday, ..., 6=Saturday
- **Return Type**: Returns `*float64` (nullable pointer) to handle cases where no data exists for the weekday
- **Edge Cases**: Returns `nil` when no logs exist for the target weekday (consistent with other mean calculation methods)

**Why This Approach:**
- Follows the existing pattern used in `GetOverallMean`, `GetPreviousPeriodMean`, etc.
- Uses nullable return type (`*float64`) consistent with `StatsData.MaxDay` field
- Handles NULL values gracefully using SQL MAX function
- Uses 15-second context timeout (consistent with `dashboardContextTimeout`)

### 2. Files to Modify

**Files to Read (no modifications needed for implementation, but needed for tests):**
- `internal/repository/dashboard_repository.go` - Interface already defined
- `internal/adapter/postgres/dashboard_repository.go` - Implementation already exists
- `internal/domain/dto/dashboard_response.go` - `StatsData.MaxDay` field already defined
- `internal/api/v1/handlers/dashboard_handler.go` - Handler already uses the method

**Files to Create/Modify:**

1. **`test/unit/dashboard_repository_test.go`** (MODIFY)
   - Add `TestDashboardRepository_GetMaxByWeekday` - Test with data for target weekday
   - Add `TestDashboardRepository_GetMaxByWeekday_EmptyWeekday` - Test with no data for weekday
   - Add `TestDashboardRepository_GetMaxByWeekday_MultipleProjects` - Test across multiple projects
   - Add `TestDashboardRepository_GetMaxByWeekday_NegativeValues` - Test edge case with start_page > end_page

2. **`test/integration/dashboard_repository_integration_test.go`** (CREATE)
   - Integration test with real database interactions
   - Test data setup with multiple projects and logs
   - Verify actual SQL query execution and result mapping

3. **`QWEN.md`** (MODIFY)
   - Add documentation for the `GetMaxByWeekday` repository method
   - Add `max_day` field to the Dashboard API documentation section
   - Include SQL query pattern reference

### 3. Dependencies

**Prerequisites:**
- PostgreSQL test database must be running (`reading_log_test`)
- Test schema must be created via `TestHelper.SetupTestSchema()`
- Existing test fixtures in `test/test_helper.go` must be functional

**Blocking Issues:**
- None - Core implementation is complete

**Setup Steps:**
1. Ensure test database is available: `make test-clean`
2. Run existing tests to verify test infrastructure: `go test ./test/unit/... -v`
3. Verify handler tests pass: `go test ./internal/api/v1/handlers/... -v`

### 4. Code Patterns

**Follow Existing Patterns:**

1. **Repository Test Pattern** (from `test/unit/dashboard_repository_test.go`):
```go
func TestDashboardRepository_GetMaxByWeekday(t *testing.T) {
    helper, err := test.SetupTestDB()
    require.NoError(t, err)
    defer helper.Close()

    err = helper.SetupTestSchema()
    require.NoError(t, err)

    repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

    // Create test data
    // Execute
    // Verify
}
```

2. **SQL Query Pattern**:
```go
query := `
    SELECT MAX(CASE 
        WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
        THEN end_page - start_page 
        ELSE 0 
    END)
    FROM logs
    WHERE EXTRACT(DOW FROM data::timestamp)::int = $1
`
```

3. **Error Handling Pattern**:
```go
if err != nil {
    if err == pgx.ErrNoRows {
        return nil, nil // Return nil for no data
    }
    return nil, fmt.Errorf("failed to get max by weekday: %w", err)
}
```

4. **Naming Conventions**:
- Test function names: `TestDashboardRepository_GetMaxByWeekday_<Scenario>`
- Test log helper: Use existing `createTestLogs` helper function
- Variable names: `maxPages`, `testDate`, `weekday`

### 5. Testing Strategy

**Unit Tests (test/unit/dashboard_repository_test.go):**

1. **TestDashboardRepository_GetMaxByWeekday**
   - Create logs with different page counts for Monday (weekday=1)
   - Verify MAX returns the highest value (e.g., 50 pages)
   - Verify no error is returned

2. **TestDashboardRepository_GetMaxByWeekday_EmptyWeekday**
   - Query for a weekday with no logs (e.g., Sunday)
   - Verify returns `nil` (not error)
   - Verify no error is returned

3. **TestDashboardRepository_GetMaxByWeekday_MultipleProjects**
   - Create logs across multiple projects for the same weekday
   - Verify MAX returns highest across all projects
   - Example: Project 1 max=30, Project 2 max=50 → Return 50

4. **TestDashboardRepository_GetMaxByWeekday_ZeroPages**
   - Create logs where end_page = start_page (0 pages read)
   - Verify returns 0.0

5. **TestDashboardRepository_GetMaxByWeekday_InvalidData**
   - Create logs where start_page > end_page (negative pages)
   - Verify the CASE statement handles this correctly (returns 0)

**Integration Tests (test/integration/dashboard_repository_integration_test.go):**

1. **TestDashboardRepository_GetMaxByWeekday_Integration**
   - Full database setup with realistic data
   - Multiple projects, multiple weekdays
   - Verify end-to-end flow from handler to database

2. **TestDashboardRepository_GetMaxByWeekday_Performance**
   - Test with large dataset (1000+ logs)
   - Verify query performance is acceptable (< 100ms)

3. **TestDashboardRepository_GetMaxByWeekday_Concurrent**
   - Test concurrent access to the method
   - Verify no race conditions or connection pool issues

**Edge Cases to Cover:**
- Empty database (no logs at all)
- Single log entry for weekday
- Multiple log entries on same day (should return max of all)
- NULL start_page or end_page values
- Very large page numbers (int overflow test)

### 6. Risks and Considerations

**Known Issues:**
- None identified - implementation follows existing patterns

**Potential Pitfalls:**

1. **PostgreSQL DOW Calculation**:
   - `EXTRACT(DOW FROM ...)` returns 0=Sunday, 1=Monday, ..., 6=Saturday
   - Ensure test dates match expected weekdays
   - Use `time.Weekday()` in Go tests for consistency

2. **NULL Handling**:
   - SQL MAX returns NULL when all values are NULL
   - Go scans NULL as nil pointer (correct behavior)
   - Test both NULL and 0.0 return cases

3. **Context Timeout**:
   - Uses `dashboardContextTimeout` (15 seconds)
   - Should be sufficient for single-table MAX query
   - No additional timeout handling needed

4. **Data Type Consistency**:
   - `StatsData.MaxDay` is `*float64`
   - Repository returns `*float64`
   - Ensure type consistency throughout the chain

**Deployment Considerations:**
- No migration required (no schema changes)
- Method is already in production code path
- Tests are the only missing piece

**Rollback Plan:**
- If tests fail, simply do not merge the PR
- No runtime impact since implementation already exists

**Acceptance Criteria Mapping:**

| AC | Status | Notes |
|----|--------|-------|
| #1 GetMaxByWeekday() method implemented in adapter | ✅ Done | Already implemented |
| #2 Interface method added to repository contract | ✅ Done | Already defined |
| #3 Returns maximum pages for target weekday | ✅ Done | Verified in code |
| #1 All unit tests pass | ❌ Pending | Add repository unit tests |
| #2 All integration tests pass | ❌ Pending | Add integration tests |
| #3 go fmt and go vet pass | ❌ Pending | Run after test additions |
| #4 Clean Architecture layers followed | ✅ Done | Follows existing pattern |
| #5 Error responses consistent | ✅ Done | Uses existing error pattern |
| #6 HTTP status codes correct | ✅ Done | Handler already handles correctly |
| #7 Documentation updated in QWEN.md | ❌ Pending | Add method documentation |
| #8 Error path tests included | ❌ Pending | Add error scenario tests |
| #9 Handler tests success/error responses | ✅ Done | Already in handler tests |
| #10 Integration tests verify DB interactions | ❌ Pending | Add integration tests |
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress - COMPLETE

### Completed Tasks

1. **Unit Tests Added** (`test/unit/dashboard_repository_test.go`)
   - `TestDashboardRepository_GetMaxByWeekday` - Tests with data for target weekday
   - `TestDashboardRepository_GetMaxByWeekday_EmptyWeekday` - Tests with no data for weekday
   - `TestDashboardRepository_GetMaxByWeekday_MultipleProjects` - Tests across multiple projects
   - `TestDashboardRepository_GetMaxByWeekday_ZeroPages` - Tests with zero pages read
   - `TestDashboardRepository_GetMaxByWeekday_InvalidData` - Tests edge case handling
   - `TestDashboardRepository_GetMaxByWeekday_EmptyDatabase` - Tests empty database
   - `TestDashboardRepository_GetMaxByWeekday_SingleEntry` - Tests single log entry

2. **Integration Tests Added** (`test/integration/dashboard_maxbyweekday_integration_test.go`)
   - `TestDashboardRepository_GetMaxByWeekday_Integration` - Full database interactions with subtests
   - `TestDashboardRepository_GetMaxByWeekday_Performance` - Performance test with 1000+ logs
   - `TestDashboardRepository_GetMaxByWeekday_LargePageNumbers` - Large page number handling

3. **Documentation Updated** (`QWEN.md`)
   - Added `max_day` field documentation to Calculated Fields section
   - Added Dashboard Repository Methods section with method table
   - Added detailed GetMaxByWeekday implementation documentation

4. **Test Utilities Added** (`test/testutil/helpers.go`)
   - Created shared helper functions (FloatPtr, IntPtr, Int64Ptr, BoolPtr, StringPtr)
   - Fixed duplicate function errors in existing test files

### Test Results
- All 7 unit tests pass
- All 3 integration tests pass (4 subtests included)
- go fmt completed successfully
- go vet passes for all modified files

### Files Modified
- `test/unit/dashboard_repository_test.go` - Added 7 unit tests
- `test/integration/dashboard_maxbyweekday_integration_test.go` - Created new file with 3 integration tests
- `QWEN.md` - Added documentation for GetMaxByWeekday method and max_day field
- `test/testutil/helpers.go` - Created new file with shared test helpers
- `test/unit/project_calculations_test.go` - Updated to use shared helpers
- `test/unit/dashboard_response_test.go` - Updated to use shared helpers
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
