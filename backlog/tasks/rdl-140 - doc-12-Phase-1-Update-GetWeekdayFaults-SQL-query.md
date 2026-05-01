---
id: RDL-140
title: '[doc-12 Phase 1] Update GetWeekdayFaults SQL query'
status: To Do
assignee:
  - workflow
created_date: '2026-05-01 15:07'
updated_date: '2026-05-01 16:33'
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
Implement SQL query in GetWeekdayFaults to count zero-page days grouped by weekday (0-6). Use same daily aggregation logic as GetFaultsByDateRange, add EXTRACT(DOW) for weekday grouping, and ensure all 7 weekdays appear in result with default 0 values.
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The current `GetWeekdayFaults` method incorrectly counts log entries grouped by weekday instead of counting days with zero pages read (faults) grouped by weekday. This implementation will fix the SQL query to match the Rails logic by applying the same CTE-based daily aggregation approach used in the fixed `GetFaultsByDateRange` method (RDL-139).

**Core Logic Change:**
- A "fault" = a day with zero pages read (sum of `read_pages` = 0 for that day)
- Group faults by weekday (0-6 = Sunday-Saturday)
- Days with no logs at all = fault for that weekday
- Days with logs but `sum(end_page - start_page) = 0` = fault for that weekday
- Days with logs where `sum(end_page - start_page) > 0` = no fault

**SQL Query Strategy:**
1. Use a CTE (Common Table Expression) approach with two parts:
   - `daily_read`: Aggregate all logs by date, calculating sum of pages read per day
   - `all_dates`: Use PostgreSQL's `generate_series()` to create all dates in the range
2. LEFT JOIN the two CTEs to include days without any logs
3. Filter to count days where `daily_pages IS NULL OR daily_pages = 0` (faults)
4. Group faults by weekday using `EXTRACT(DOW FROM log_date)`
5. Ensure all 7 weekdays (0-6) are present in the result map with default value of 0

**Why CTE Approach (matching GetFaultsByDateRange):**
- Consistent with the fix applied in RDL-139 for GetFaultsByDateRange
- Clear separation of concerns (aggregation vs. date generation)
- Matches Rails logic of "for each day in range, check if any reading occurred"
- Handles edge cases: empty database, single day ranges, NULL values
- PostgreSQL's `generate_series()` is efficient for date ranges

**Edge Cases to Handle:**
- Empty database (all days in range = faults for their respective weekdays)
- Single day range (start_date == end_date)
- Logs with NULL start_page or end_page (treated as 0 pages)
- Multiple logs on same day (aggregate by sum, count as 1 day)
- Logs with zero pages (end_page == start_page)
- Date ranges spanning multiple weeks (multiple occurrences of same weekday)

### 2. Files to Modify

#### Primary Implementation Files
| File | Change Type | Description |
|------|-------------|-------------|
| `internal/adapter/postgres/dashboard_repository.go` | Modify | Update `GetWeekdayFaults` method with CTE-based SQL query |

**Specific Changes to `dashboard_repository.go`:**

Replace current implementation:
```go
query := `
    SELECT 
        EXTRACT(DOW FROM data::timestamp)::int as weekday,
        COUNT(*) as fault_count
    FROM logs
    WHERE data::date BETWEEN $1 AND $2
    GROUP BY EXTRACT(DOW FROM data::timestamp)
    ORDER BY weekday
`
```

With new CTE-based implementation:
```go
query := `
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
    SELECT 
        EXTRACT(DOW FROM ad.log_date)::int as weekday,
        COUNT(*) as fault_count
    FROM all_dates ad
    LEFT JOIN daily_read dr ON ad.log_date = dr.log_date
    WHERE dr.daily_pages IS NULL OR dr.daily_pages = 0
    GROUP BY EXTRACT(DOW FROM ad.log_date)
    ORDER BY weekday
`
```

#### Test Files to Update
| File | Change Type | Description |
|------|-------------|-------------|
| `test/unit/dashboard_repository_test.go` | Update | Fix `TestDashboardRepository_GetWeekdayFaults` to expect correct fault counts by weekday |
| `test/dashboard_integration_test.go` | Review | Check if any integration tests for weekday faults need updates |
| `test/api/v1/handlers/dashboard_handler_test.go` | Review | Update mock repository expectations if weekday faults endpoint is tested |

#### Files to Read (No Changes Required)
| File | Purpose |
|------|---------|
| `internal/repository/dashboard_repository.go` | Review interface signature |
| `internal/domain/dto/dashboard_response.go` | Review WeekdayFaults DTO structure |
| `test/test_helper.go` | Review test database setup utilities |
| `test/fixtures/dashboard_fixtures.go` | Review test data fixture patterns |

### 3. Dependencies

#### Prerequisites
1. **PostgreSQL 13+**: Required for `generate_series()` function support
2. **Test Database**: `reading_log_test` must be running and accessible
3. **RDL-139 Completion**: The GetFaultsByDateRange CTE pattern must be implemented first (already done)

#### Related Tasks
| Task ID | Relationship | Status |
|---------|--------------|--------|
| RDL-139 | Prerequisite - CTE pattern established | Done |
| RDL-141 | Sequential dependency (Phase 2 - context timeout verification) | To Do |
| RDL-142 | Sequential dependency (Phase 3 - unit test rewrite) | To Do |

#### Environment Setup
```bash
# Ensure test database is ready
make test-setup

# Verify PostgreSQL is running
pg_isready -h localhost -p 5432

# Run existing tests to establish baseline
go test -v ./test/unit/... -run TestDashboardRepository_GetWeekdayFaults
```

### 4. Code Patterns

#### SQL Query Pattern
Follow the exact CTE pattern used in `GetFaultsByDateRange` (RDL-139) for consistency:
```go
query := `
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
    SELECT 
        EXTRACT(DOW FROM ad.log_date)::int as weekday,
        COUNT(*) as fault_count
    FROM all_dates ad
    LEFT JOIN daily_read dr ON ad.log_date = dr.log_date
    WHERE dr.daily_pages IS NULL OR dr.daily_pages = 0
    GROUP BY EXTRACT(DOW FROM ad.log_date)
    ORDER BY weekday
`
```

#### Error Handling Pattern
Maintain consistency with existing repository methods:
```go
func (r *DashboardRepositoryImpl) GetWeekdayFaults(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error) {
    ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
    defer cancel()

    query := `...`

    rows, err := r.pool.Query(ctx, query, start, end)
    if err != nil {
        return nil, fmt.Errorf("failed to query weekday faults: %w", err)
    }
    defer rows.Close()

    result := make(map[int]int)
    for rows.Next() {
        var weekday int
        var count int
        if err := rows.Scan(&weekday, &count); err != nil {
            return nil, fmt.Errorf("failed to scan weekday fault: %w", err)
        }
        result[weekday] = count
    }

    // Ensure all 7 days are present (0-6) with default value of 0
    for i := 0; i < 7; i++ {
        if _, exists := result[i]; !exists {
            result[i] = 0
        }
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating weekday faults: %w", err)
    }

    return dto.NewWeekdayFaults(result), nil
}
```

#### Comment Pattern
Add inline comments explaining the SQL logic (matching RDL-139 style):
```go
// GetWeekdayFaults returns fault distribution by weekday (0-6 = Sunday-Saturday)
// Uses CTE to:
// 1. Aggregate daily reading activity (daily_read CTE)
// 2. Generate all dates in range using generate_series (all_dates CTE)
// 3. LEFT JOIN to identify days with no reading activity
// 4. Group by weekday and count faults
// A "fault" = a day where sum(read_pages) = 0 (either no logs or logs with zero pages)
// Ensures all 7 weekdays (0-6) are present in result with default value of 0
```

#### Naming Conventions
- Keep method signature unchanged: `GetWeekdayFaults(ctx, start, end)`
- Maintain variable names: `start`, `end`, `result`, `weekday`, `count`
- SQL column aliases: `log_date`, `daily_pages`, `weekday`, `fault_count`
- Weekday values: 0=Sunday, 1=Monday, ..., 6=Saturday (PostgreSQL DOW convention)

### 5. Testing Strategy

#### Unit Tests (test/unit/dashboard_repository_test.go)

**Test Case 1: Single weekday with reading (current test)**
```go
// Given: 8-day range from Monday 2024-01-15 to Monday 2024-01-22
// Logs only on Monday 2024-01-15 (2 logs, but same day = 1 day with reading)
// Weekdays in range: Mon(1), Tue(2), Wed(3), Thu(4), Fri(5), Sat(6), Sun(0), Mon(1)
// Days with reading: 1 (Monday 2024-01-15)
// Faults: 
//   - Monday: 1 fault (Monday 2024-01-22 has no reading)
//   - All other weekdays: 1 fault each (no logs on those days)
// When: GetWeekdayFaults called
// Then: Returns faults map with 1 fault for each weekday (0-6)
```

**Test Case 2: Empty database (all weekdays have faults)**
```go
// Given: 7-day range with no logs
// Each weekday appears once in the range
// When: GetWeekdayFaults called
// Then: Returns 1 fault for each weekday (0-6)
```

**Test Case 3: All weekdays have reading (0 faults)**
```go
// Given: 7-day range (Sunday to Saturday)
// Each day has at least one log with pages > 0
// When: GetWeekdayFaults called
// Then: Returns 0 faults for all weekdays (0-6)
```

**Test Case 4: Multiple weeks in range**
```go
// Given: 14-day range (2 full weeks)
// Logs only on Mondays (2 Mondays with reading)
// Weekday distribution:
//   - Monday: 2 days in range, 2 with reading = 0 faults
//   - All other weekdays: 2 days each, 0 with reading = 2 faults each
// When: GetWeekdayFaults called
// Then: Returns Monday=0, all others=2
```

**Test Case 5: Logs with zero pages**
```go
// Given: 3-day range
// Day 1: Log with start_page=0, end_page=0 (zero pages)
// Day 2: Log with start_page=10, end_page=20 (10 pages)
// Day 3: No logs
// When: GetWeekdayFaults called
// Then: Day 1 = 1 fault (zero pages), Day 2 = 0 faults, Day 3 = 1 fault
```

#### Integration Tests (test/dashboard_integration_test.go)

**Test Case 1: Empty database weekday distribution**
```go
// Given: Empty database, 30-day range
// When: HTTP endpoint /v1/dashboard/weekday-faults called
// Then: Returns faults distributed across weekdays based on date range
```

**Test Case 2: Single log entry weekday**
```go
// Given: 7-day range, 1 log on Wednesday with 10 pages
// When: HTTP endpoint called
// Then: Wednesday has 0 faults, all other weekdays have 1 fault
```

**Test Case 3: Multiple logs same weekday across weeks**
```go
// Given: 14-day range (2 weeks)
// Logs on Monday week 1 (10 pages), Monday week 2 (15 pages)
// When: HTTP endpoint called
// Then: Monday has 0 faults, all other weekdays have 2 faults
```

#### Test Execution Order
1. Run unit tests first to validate SQL query logic
2. Run integration tests to validate HTTP endpoint behavior
3. Run full test suite to ensure no regressions

```bash
# Run unit tests
go test -v ./test/unit/... -run TestDashboardRepository_GetWeekdayFaults

# Run integration tests
go test -v ./test/... -run TestDashboardWeekdayFaults

# Run full suite
go test ./...
```

### 6. Risks and Considerations

#### Known Risks

**Risk 1: Test Expectations Need Update**
- **Probability**: High (by design)
- **Impact**: Low
- **Mitigation**:
  - Current tests expect incorrect behavior (counting logs, not faults)
  - Tests will be updated with correct expectations
  - Document test changes in commit messages

**Risk 2: Date Range Edge Cases**
- **Probability**: Medium
- **Impact**: Medium
- **Mitigation**:
  - Test single day ranges explicitly
  - Test ranges spanning partial weeks
  - Test ranges spanning multiple full weeks

**Risk 3: Performance with Large Date Ranges**
- **Probability**: Low
- **Impact**: Medium
- **Mitigation**:
  - Use EXPLAIN ANALYZE to compare query performance
  - `generate_series()` is efficient for typical ranges (< 90 days)
  - Add context timeout (already implemented: 15 seconds)

#### Technical Considerations

**Date Casting Behavior:**
- `data` column is VARCHAR in schema
- Casting with `data::date` and `data::timestamp` works but may have performance impact
- Consider adding generated column or index for production optimization

**Weekday Calculation (PostgreSQL DOW):**
- 0 = Sunday, 1 = Monday, ..., 6 = Saturday
- Must match Rails `EXTRACT(DOW FROM ...)` behavior exactly
- Test with known dates to verify DOW values

**Timezone Handling:**
- Server local time used for date casting (matches current implementation)
- Ensure `TZ` environment variable is set correctly
- Rails comparison requires matching timezone configuration

**NULL Value Handling:**
- SQL CASE statement handles NULL start_page/end_page
- LEFT JOIN handles days with no logs (daily_pages = NULL)
- Both NULL and 0 daily_pages count as faults

#### Deployment Considerations

**Rollback Plan:**
1. Keep original SQL query in version control
2. If issues arise, revert to previous implementation
3. Investigate failures in isolated branch

**Monitoring:**
- Track query execution time in production logs
- Monitor for increased database load
- Alert if weekday fault distribution deviates significantly from expected values

**Documentation Updates Required:**
- Update AGENTS.md with new weekday fault calculation explanation
- Add SQL query breakdown to docs/faults-calculation-explanation.md (Phase 5)
- Update API documentation to clarify weekday fault definition

#### Acceptance Criteria Verification

| Criterion | Verification Method |
|-----------|---------------------|
| CTE-based SQL query implemented | Code review of `dashboard_repository.go` |
| Counts days with zero pages by weekday | Unit tests with controlled scenarios |
| Uses generate_series | SQL query inspection |
| LEFT JOIN for days without logs | Unit tests with empty database |
| Groups by weekday (0-6) | Unit tests verify weekday distribution |
| All 7 weekdays present in result | Unit tests check all keys 0-6 exist |
| Handles NULL values | Unit test with NULL start_page/end_page |
| Context timeout applied | Code review of timeout implementation |
| Error handling consistent | Code review of error wrapping |

---

**Implementation Ready**: ✅ This plan provides sufficient detail for another agent to implement the fix without reading additional documentation. All edge cases, test scenarios, and risk mitigations are documented.

**Note**: The test `TestDashboardRepository_GetWeekdayFaults` currently expects `stats.Faults[1] = 2` (counting 2 logs on Monday). After the fix, it should expect `stats.Faults[1] = 1` (1 Monday without reading in the 8-day range). This is expected and tests will be updated in RDL-142.
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
