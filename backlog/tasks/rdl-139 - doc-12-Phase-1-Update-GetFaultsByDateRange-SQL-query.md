---
id: RDL-139
title: '[doc-12 Phase 1] Update GetFaultsByDateRange SQL query'
status: To Do
assignee:
  - thomas
created_date: '2026-05-01 15:07'
updated_date: '2026-05-01 16:19'
labels:
  - bugfix
  - repository
  - phase-1
dependencies: []
documentation:
  - doc-012
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement CTE-based SQL query in GetFaultsByDateRange to count days with zero pages read instead of counting log entries. Use generate_series to create all dates in range, LEFT JOIN with daily_read aggregation, and CASE statement to handle NULL values.
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The current `GetFaultsByDateRange` method incorrectly counts log entries instead of counting days with zero pages read. This implementation will fix the SQL query to match the Rails logic exactly by:

**Core Logic Change:**
- A "fault" = a day with zero pages read (sum of `read_pages` = 0 for that day)
- Days with no logs at all = 1 fault
- Days with logs but `sum(end_page - start_page) = 0` = 1 fault
- Days with logs where `sum(end_page - start_page) > 0` = 0 faults

**SQL Query Strategy:**
1. Use a CTE (Common Table Expression) approach with two parts:
   - `daily_read`: Aggregate all logs by date, calculating sum of pages read per day
   - `all_dates`: Use PostgreSQL's `generate_series()` to create all dates in the range
2. LEFT JOIN the two CTEs to include days without any logs
3. COUNT days where `daily_pages IS NULL OR daily_pages = 0`

**Why CTE Approach:**
- Clear separation of concerns (aggregation vs. date generation)
- Matches Rails logic of "for each day in range, check if any reading occurred"
- Handles edge cases: empty database, single day ranges, NULL values
- PostgreSQL's `generate_series()` is efficient for date ranges

**Edge Cases to Handle:**
- Empty database (all days in range = faults)
- Single day range (start_date == end_date)
- Logs with NULL start_page or end_page (treat as 0 pages)
- Multiple logs on same day (aggregate by sum)
- Logs with zero pages (end_page == start_page)

### 2. Files to Modify

#### Primary Implementation Files
| File | Change Type | Description |
|------|-------------|-------------|
| `internal/adapter/postgres/dashboard_repository.go` | Modify | Update `GetFaultsByDateRange` method with CTE-based SQL query |

**Specific Changes to `dashboard_repository.go`:**
```go
// Replace current implementation:
// SELECT COUNT(*) FROM logs WHERE data::date BETWEEN $1 AND $2

// With new CTE-based implementation:
WITH daily_read AS (
    SELECT 
        data::date as log_date,
        SUM(CASE 
            WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
            THEN end_page - start_page 
            ELSE 0 
        END) as daily_pages
    FROM logs
    WHERE data::date BETWEEN $1 AND $2
    GROUP BY data::date
),
all_dates AS (
    SELECT generate_series($1::date, $2::date, '1 day'::interval)::date as log_date
)
SELECT COUNT(*) as fault_count
FROM all_dates ad
LEFT JOIN daily_read dr ON ad.log_date = dr.log_date
WHERE dr.daily_pages IS NULL OR dr.daily_pages = 0
```

#### Test Files to Update
| File | Change Type | Description |
|------|-------------|-------------|
| `test/unit/dashboard_repository_test.go` | Rewrite | Update `TestDashboardRepository_GetFaultsByDateRange` with correct expectations |
| `test/dashboard_integration_test.go` | Update | Update integration test expectations in `TestDashboardFaults_EmptyDatabase_Integration` and `TestDashboardFaults_SingleLogEntry_Integration` |
| `test/api/v1/handlers/dashboard_handler_test.go` | Update | Update mock repository expectations for faults calculation |

#### Files to Read (No Changes Required)
| File | Purpose |
|------|---------|
| `internal/repository/dashboard_repository.go` | Review interface signature |
| `internal/domain/dto/dashboard_response.go` | Review FaultStats DTO structure |
| `test/test_helper.go` | Review test database setup utilities |
| `test/fixtures/dashboard_fixtures.go` | Review test data fixture patterns |

### 3. Dependencies

#### Prerequisites
1. **PostgreSQL 13+**: Required for `generate_series()` function support
2. **Test Database**: `reading_log_test` must be running and accessible
3. **Existing Task Completion**: This is Phase 1 of doc-012; Phase 2 (RDL-140) updates `GetWeekdayFaults`

#### Related Tasks
| Task ID | Relationship | Status |
|---------|--------------|--------|
| RDL-140 | Sequential dependency (Phase 2) | To Do |
| RDL-137 | Sequential dependency (Phase 1 - GetFaultsByDateRange with CTE) | To Do |
| RDL-138 | Sequential dependency (Phase 2 - Mock updates) | To Do |

#### Environment Setup
```bash
# Ensure test database is ready
make test-setup

# Verify PostgreSQL is running
pg_isready -h localhost -p 5432

# Run existing tests to establish baseline
go test -v ./test/unit/... -run TestDashboardRepository_GetFaultsByDateRange
```

### 4. Code Patterns

#### SQL Query Pattern
Follow the existing CTE pattern used in `GetProjectsWithLogs`:
```go
query := `
    WITH cte_name AS (
        -- CTE definition
    )
    SELECT ...
    FROM cte_name
    ...
`
```

#### Error Handling Pattern
Maintain consistency with existing repository methods:
```go
func (r *DashboardRepositoryImpl) GetFaultsByDateRange(ctx context.Context, start, end time.Time) (*dto.FaultStats, error) {
    ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
    defer cancel()

    query := `...`

    var stats dto.FaultStats
    err := r.pool.QueryRow(ctx, query, start, end).Scan(&stats.FaultCount)
    if err != nil {
        return nil, fmt.Errorf("failed to get faults by date range: %w", err)
    }

    return &stats, nil
}
```

#### Comment Pattern
Add inline comments explaining the SQL logic:
```go
// GetFaultsByDateRange returns the count of faults (days with zero pages read) within a date range
// Uses CTE to:
// 1. Aggregate daily reading activity
// 2. Generate all dates in range using generate_series
// 3. LEFT JOIN to count days with no reading activity
// A "fault" = a day where sum(read_pages) = 0
```

#### Naming Conventions
- Keep method signature unchanged: `GetFaultsByDateRange(ctx, start, end)`
- Maintain variable names: `start`, `end`, `stats`, `fault_count`
- SQL column aliases: `log_date`, `daily_pages`, `fault_count`

### 5. Testing Strategy

#### Unit Tests (test/unit/dashboard_repository_test.go)

**Test Case 1: No logs in range (all days are faults)**
```go
// Given: 7-day range with no logs
// When: GetFaultsByDateRange called
// Then: Returns 7 faults (all 7 days have zero pages)
```

**Test Case 2: All days have reading (0 faults)**
```go
// Given: 5-day range with logs on every day
// Each day has at least one log with pages > 0
// When: GetFaultsByDateRange called
// Then: Returns 0 faults
```

**Test Case 3: Mixed scenario (some faults)**
```go
// Given: 7-day range
// Day 1: 10 pages read → 0 faults
// Day 2: 0 pages read (no logs) → 1 fault
// Day 3: 20 pages read → 0 faults
// Day 4: 0 pages read (log with start=0, end=0) → 1 fault
// Day 5: 15 pages read → 0 faults
// Days 6-7: no logs → 2 faults
// When: GetFaultsByDateRange called
// Then: Returns 4 faults
```

**Test Case 4: Single day range**
```go
// Given: Single day (start == end)
// No logs on that day
// When: GetFaultsByDateRange called
// Then: Returns 1 fault
```

**Test Case 5: NULL values handling**
```go
// Given: Log with start_page = NULL, end_page = 10
// When: GetFaultsByDateRange called
// Then: NULL treated as 0 pages, day counts as fault
```

#### Integration Tests (test/dashboard_integration_test.go)

**Test Case 1: Empty database**
```go
// Given: Empty database, 30-day range
// When: HTTP endpoint called
// Then: Returns 30 faults (all days are faults)
```

**Test Case 2: Single log entry**
```go
// Given: 7-day range, 1 log on day 3 with 10 pages
// When: HTTP endpoint called
// Then: Returns 6 faults (7 days - 1 day with reading)
```

**Test Case 3: Multiple logs same day**
```go
// Given: 3-day range
// Day 2: 3 logs with pages 10, 15, 5 (total 30 pages)
// When: HTTP endpoint called
// Then: Returns 2 faults (Day 2 has reading, Days 1 and 3 are faults)
```

#### Test Execution Order
1. Run unit tests first to validate SQL query logic
2. Run integration tests to validate HTTP endpoint behavior
3. Run full test suite to ensure no regressions

```bash
# Run unit tests
go test -v ./test/unit/... -run TestDashboardRepository_GetFaultsByDateRange

# Run integration tests
go test -v ./test/... -run TestDashboardFaults

# Run full suite
go test ./...
```

### 6. Risks and Considerations

#### Known Risks

**Risk 1: Performance Impact**
- **Probability**: Low
- **Impact**: Medium
- **Mitigation**: 
  - Use EXPLAIN ANALYZE to compare query performance before/after
  - `generate_series()` is efficient for typical date ranges (< 90 days)
  - Add database index on `data::date` if not present

**Risk 2: Breaking Existing Tests**
- **Probability**: High (by design)
- **Impact**: Low
- **Mitigation**:
  - All failing tests are expected (they validate incorrect behavior)
  - Tests will be rewritten with correct expectations
  - Document test changes in commit messages

**Risk 3: Edge Case Handling**
- **Probability**: Medium
- **Impact**: Medium
- **Mitigation**:
  - Comprehensive test coverage for edge cases
  - NULL value handling verified in SQL CASE statement
  - Single day range tested explicitly

#### Technical Considerations

**Date Casting Behavior:**
- `data` column is VARCHAR in schema
- Casting with `data::date` works but may have performance impact
- Consider adding generated column or index for production optimization

**Timezone Handling:**
- Server local time used for date casting (matches current implementation)
- Ensure `TZ` environment variable is set correctly
- Rails comparison requires matching timezone configuration

**Connection Pooling:**
- 15-second timeout applied via `dashboardContextTimeout`
- CTE query should complete well within timeout for typical ranges
- Monitor for timeout issues with very large date ranges (> 1 year)

#### Deployment Considerations

**Rollback Plan:**
1. Keep original SQL query in version control
2. If issues arise, revert to previous implementation
3. Investigate failures in isolated branch

**Monitoring:**
- Track query execution time in production logs
- Monitor for increased database load
- Alert if fault percentage deviates significantly from expected values

**Documentation Updates Required:**
- Update AGENTS.md with new fault calculation explanation
- Add SQL query breakdown to docs/faults-calculation-explanation.md (Phase 5)
- Update API documentation to clarify fault definition

#### Acceptance Criteria Verification

| Criterion | Verification Method |
|-----------|---------------------|
| CTE-based SQL query implemented | Code review of `dashboard_repository.go` |
| Counts days with zero pages | Unit tests with controlled scenarios |
| Uses generate_series | SQL query inspection |
| LEFT JOIN for days without logs | Unit tests with empty database |
| Handles NULL values | Unit test with NULL start_page/end_page |
| Inclusive date range | Test with single day range |
| Context timeout applied | Code review of timeout implementation |
| Error handling consistent | Code review of error wrapping |

---

**Implementation Ready**: ✅ This plan provides sufficient detail for another agent to implement the fix without reading additional documentation. All edge cases, test scenarios, and risk mitigations are documented.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress - RDL-139: Update GetFaultsByDateRange SQL query

### Completed Tasks

#### 1. Updated SQL Query in `internal/adapter/postgres/dashboard_repository.go`
- Replaced simple log count query with CTE-based implementation
- Added `daily_read` CTE to aggregate daily reading activity
- Added `all_dates` CTE using `generate_series()` to create all dates in range
- LEFT JOIN to include days without logs
- COUNT days where `daily_pages IS NULL OR daily_pages = 0`

#### 2. Updated Unit Tests in `test/unit/dashboard_repository_test.go`
- Rewrote `TestDashboardRepository_GetFaultsByDateRange` with correct expectations
- Added `TestDashboardRepository_GetFaultsByDateRange_NoLogs` - all days are faults
- Added `TestDashboardRepository_GetFaultsByDateRange_AllDaysHaveReading` - 0 faults
- Added `TestDashboardRepository_GetFaultsByDateRange_SingleDay` - single day range
- Added `TestDashboardRepository_GetFaultsByDateRange_ZeroPagesRead` - logs with 0 pages
- Added `ptr()` helper function for string pointers

#### 3. Updated Integration Tests in `test/dashboard_integration_test.go`
- Updated comments to reflect RDL-139 fix completion
- Changed log statements to assertions for better test validation
- Fixed `TestDashboardFaultsChart_Integration` to use `dto.SetTestDate()` for deterministic testing
- Updated expected fault percentage from 50% to 160% (16 faults / 10 maxFaults)

#### 4. Added Helper Function in `internal/domain/dto/dashboard.go`
- Added `ResetTestDate()` function to reset the global date function after tests

### Test Results

All tests pass:
- ✅ Unit tests: `TestDashboardRepository_GetFaultsByDateRange` (5 test cases)
- ✅ Integration tests: All `TestDashboardFaults*` tests pass
- ✅ `go fmt` - passes
- ✅ `go vet` - passes
- ✅ Build - succeeds

### Key Changes

1. **SQL Query Logic**: Changed from counting log entries to counting days with zero pages read
2. **Edge Cases Handled**:
   - Empty database (all days are faults)
   - Single day range
   - Logs with zero pages (end_page == start_page)
   - Multiple logs on same day (count as 1 day with reading)
   - NULL values in start_page/end_page

### Files Modified
- `internal/adapter/postgres/dashboard_repository.go` - Main implementation
- `test/unit/dashboard_repository_test.go` - Unit tests
- `test/dashboard_integration_test.go` - Integration tests
- `internal/domain/dto/dashboard.go` - Added ResetTestDate() helper
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass
- [x] #2 All integration tests pass execution and verification
- [x] #3 go fmt and go vet pass with no errors
- [x] #4 Clean Architecture layers properly followed
- [x] #5 Error responses consistent with existing patterns
- [x] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md
- [x] #8 New code paths include error path tests
- [x] #9 HTTP handlers test both success and error responses
- [ ] #10 Integration tests verify actual database interactions
<!-- DOD:END -->
