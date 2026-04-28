---
id: RDL-116
title: '[doc-10 Phase 3] Align mean_day calculation with Rails V1::MeanLog'
status: To Do
assignee:
  - thomas
created_date: '2026-04-28 00:29'
updated_date: '2026-04-28 03:14'
labels:
  - calculation
  - phase-3
  - backend
dependencies: []
documentation:
  - doc-010
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Study Rails V1::MeanLog implementation and replicate exact algorithm in Go. Formula: total_pages / count_reads where count_reads = number of 7-day intervals since begin_data. Round to 3 decimals.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 mean_day calculation matches Rails output exactly
- [ ] #2 Algorithm uses 7-day intervals from begin_data
- [ ] #3 Values rounded to 3 decimal places
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The current `mean_day` calculation uses a simple average (`total_pages / log_count`), but the Rails V1::MeanLog algorithm uses a weighted calculation based on 7-day intervals from the first log entry.

**Algorithm to implement:**
1. Filter all logs by the target weekday (DOW 0-6)
2. Calculate `total_pages = sum(end_page - start_page)` for all filtered logs
3. Find `begin_data` (timestamp of first log) and `log_data` (timestamp of most recent log)
4. Calculate `count_reads = floor((log_data - begin_data) / 7 days)` - number of complete 7-day intervals
5. Calculate `mean_day = total_pages / count_reads`, rounded to 3 decimals
6. Edge cases: return 0.0 if no logs exist or if count_reads is zero

**Why this approach:**
- Matches Rails V1::MeanLog implementation exactly (documented in `docs/rails-calculation-reference.md`)
- Provides more accurate "pages per week" metric rather than simple average
- Follows existing Clean Architecture patterns (repository interface + adapter implementation)
- Maintains consistency with other dashboard calculations

**Architecture:**
- Add new repository method `GetMeanByWeekday(ctx, weekday)` to DashboardRepository interface
- Implement in `DashboardRepositoryImpl` with a single SQL query that calculates all values
- Update handler to call new method instead of using `LogCount` for mean calculation

### 2. Files to Modify

**Interface Definition:**
- `internal/repository/dashboard_repository.go`
  - Add method: `GetMeanByWeekday(ctx context.Context, weekday int) (*float64, error)`

**Implementation:**
- `internal/adapter/postgres/dashboard_repository.go`
  - Add `GetMeanByWeekday` implementation with SQL query:
    ```sql
    SELECT 
        COALESCE(SUM(CASE WHEN start_page IS NOT NULL AND end_page IS NOT NULL THEN end_page - start_page ELSE 0 END), 0) as total_pages,
        MIN(data::timestamp) as begin_data,
        MAX(data::timestamp) as log_data,
        COUNT(*) as log_count
    FROM logs
    WHERE EXTRACT(DOW FROM data::timestamp)::int = $1
    ```
  - Calculate 7-day intervals in Go: `daysDiff := logData.Sub(beginData).Hours() / 24`, `countReads := int(daysDiff / 7)`
  - Return `total_pages / countReads` rounded to 3 decimals
  - Handle edge cases: empty result → nil, zero intervals → nil

**Handler Update:**
- `internal/api/v1/handlers/dashboard_handler.go`
  - In `Day()` handler, replace current mean_day calculation:
    ```go
    // Current (incorrect):
    if stats.LogCount > 0 {
        statsData.MeanDay = math.Round(float64(stats.TotalPages)/float64(stats.LogCount)*1000) / 1000
    }
    
    // New (correct):
    meanDay, err := h.repo.GetMeanByWeekday(ctx, int(targetDate.Weekday()))
    if err == nil && meanDay != nil {
        statsData.MeanDay = *meanDay
    } else {
        statsData.MeanDay = 0.0
    }
    ```

**Tests to Update:**
- `internal/api/v1/handlers/dashboard_handler_test.go`
  - Update existing `Day` handler tests to use correct expected values based on new algorithm
  - Add test cases for edge cases: empty logs, single week of data (zero intervals)

### 3. Dependencies

**Prerequisites:**
- RDL-115 (Phase 2) - Already completed (PerMeanDay and PerSpecMeanDay fields implemented)
- RDL-112 (Phase 1) - Already completed (Day handler structure in place)
- Existing database schema with `logs` table containing `data`, `start_page`, `end_page` columns

**No blocking issues** - The implementation is independent and can be done in isolation.

### 4. Code Patterns

**Follow existing patterns:**

1. **Repository Interface Pattern:**
   ```go
   type DashboardRepository interface {
       GetMeanByWeekday(ctx context.Context, weekday int) (*float64, error)
   }
   ```

2. **Adapter Implementation Pattern:**
   - Use `dashboardContextTimeout = 15 * time.Second` for context timeout
   - Use `context.WithTimeout` with defer cancel()
   - Wrap errors with context: `fmt.Errorf("failed to get mean by weekday: %w", err)`
   - Return `nil` for no data (not error)

3. **Rounding Pattern:**
   - Round to 3 decimals: `math.Round(value*1000) / 1000`
   - Matches Rails `.round(3)` behavior

4. **Weekday Mapping:**
   - Use Go's `time.Weekday` constants: `int(date.Weekday())`
   - DOW 0 = Sunday, DOW 1 = Monday, ..., DOW 6 = Saturday
   - Matches PostgreSQL's `EXTRACT(DOW FROM timestamp)`

5. **Error Handling:**
   - Return `nil, nil` for no data (consistent with `GetMaxByWeekday`, `GetOverallMean`)
   - Return error only for actual database errors

### 5. Testing Strategy

**Unit Tests (dashboard_handler_test.go):**
1. **TestDay_Success** - Verify mean_day calculation with fixed test data:
   - Setup: Create logs spanning multiple weeks for a specific weekday
   - Example: 3 Monday logs over 14 days with total 80 pages
   - Expected: `count_reads = 2`, `mean_day = 80 / 2 = 40.0`
   - Verify JSON response contains correct `mean_day`

2. **TestDay_EmptyLogs** - Verify zero mean when no logs exist:
   - Setup: No logs for target weekday
   - Expected: `mean_day = 0.0`

3. **TestDay_ZeroIntervals** - Verify zero mean when logs within same week:
   - Setup: Multiple logs for same weekday within 7-day period
   - Expected: `count_reads = 0`, `mean_day = 0.0`

**Integration Tests:**
1. Use `TestHelper` from `test/test_helper.go` for database setup
2. Insert deterministic test data with known dates and page counts
3. Call actual `/v1/dashboard/day.json` endpoint
4. Compare response `mean_day` against manually calculated expected value
5. Verify rounding to 3 decimal places

**Edge Cases to Cover:**
- Empty logs table
- Single log entry (zero intervals)
- Logs spanning exactly 7 days (1 interval)
- Logs spanning multiple weeks (multiple intervals)
- Logs with NULL start_page or end_page (should be excluded from sum)

### 6. Risks and Considerations

**Known Issues:**
- **Timezone handling**: The `data` column is stored as `VARCHAR`, not timestamp. The query casts to `timestamp` which may have timezone implications. Ensure consistency with Rails behavior (documented in `rails-calculation-reference.md`).
- **Floating point precision**: Rounding to 3 decimals may still have minor floating point differences. Use `math.Round(value*1000) / 1000` pattern consistently.

**Potential Pitfalls:**
1. **Division by zero**: Ensure `count_reads == 0` returns `nil` (not panic)
2. **NULL handling**: Logs with NULL `start_page` or `end_page` should be excluded from `total_pages` calculation
3. **Weekday alignment**: Ensure Go's `time.Weekday()` (0=Sunday) matches PostgreSQL's `EXTRACT(DOW FROM ...)` (0=Sunday)

**Testing Considerations:**
- Compare Go implementation output against Rails API output using same test data
- Use the comparison test setup from RDL-121 (Rails API comparison test) for verification
- Document any discrepancies found for Phase 6 (RDL-123)

**Deployment:**
- No migration required - only code changes
- Backward compatible - response format unchanged, only calculation improved
- Monitor for any client-side issues if mean_day values change significantly

**Rollback Plan:**
- If issues arise, revert to previous calculation using `LogCount`
- Keep old code in comments for reference during transition
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
