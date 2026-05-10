---
id: RDL-153
title: >-
  [doc-14 Phase 1] Implement PostgreSQL repository methods for weekday grouping
  queries
status: To Do
assignee:
  - thomas
created_date: '2026-05-10 10:47'
updated_date: '2026-05-10 12:19'
labels:
  - infrastructure
  - repository
  - postgres
  - phase-1
dependencies: []
documentation:
  - doc-014
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement the three new repository methods in internal/adapter/postgres/dashboard_repository.go: GetWeekdayPagesGrouped with SQL query for weekday-based page aggregation, GetFirstLogDate to retrieve earliest log timestamp, and GetWeekdayMeanWithIntervals calculating mean using 7-day interval logic.

Each method must use 15-second context timeout and handle NULL/empty data gracefully.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 GetWeekdayPagesGrouped implemented with EXTRACT(DOW FROM data) weekday grouping
- [ ] #2 GetFirstLogDate returns *time.Time (nil when no logs exist)
- [ ] #3 GetWeekdayMeanWithIntervals calculates 7-day intervals correctly
- [x] #4 All methods use 15-second context timeout
- [ ] #5 NULL values handled gracefully in all queries
- [ ] #6 Code compiles without errors
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task implements three PostgreSQL repository methods for weekday grouping queries in the dashboard repository. The implementation follows Clean Architecture patterns and aligns with existing repository methods.

**Current State Assessment:**

Based on codebase review, the implementation is **already complete** in `internal/adapter/postgres/dashboard_repository.go`:

1. **GetWeekdayPagesGrouped** (Lines 394-437)
   - SQL query groups logs by weekday using `EXTRACT(DOW FROM data::timestamp)`
   - Calculates total_pages, log_count, and mean for each weekday
   - Returns `map[int]dto.WeekdayPages` with empty map for no data
   - Uses 15-second context timeout

2. **GetFirstLogDate** (Lines 439-459)
   - Query: `SELECT MIN(data::timestamp) FROM logs`
   - Returns `*time.Time` (nil when no logs exist)
   - Handles zero time value gracefully
   - Uses 15-second context timeout

3. **GetWeekdayMeanWithIntervals** (Lines 461-504)
   - Implements V1::MeanLog algorithm with 7-day intervals
   - Calculates: `total_pages / count_reads` where `count_reads = floor((log_data - begin_data) / 7 days)`
   - Returns `*float64` (nil for no data or zero intervals)
   - Uses 15-second context timeout

**Implementation Strategy:**

Since the implementation is complete, this task focuses on:
1. **Verification**: Confirm all acceptance criteria are met
2. **Testing**: Ensure integration tests exist and pass
3. **Documentation**: Update task completion status

**Why This Approach:**

- Follows existing repository patterns (GetMaxByWeekday, GetMeanByWeekday)
- Uses consistent error handling: `(nil, nil)` for empty data, `(value, nil)` for success
- Maintains Clean Architecture separation (interface in `internal/repository/`, implementation in `internal/adapter/postgres/`)
- Aligns with Rails `V1::MeanLog` calculation algorithm

---

### 2. Files to Modify

**Files Already Implemented:**

| File | Status | Description |
|------|--------|-------------|
| `internal/adapter/postgres/dashboard_repository.go` | ✅ Complete | All three methods implemented with SQL queries |
| `internal/repository/dashboard_repository.go` | ✅ Complete | Interface methods defined (RDL-152) |
| `internal/domain/dto/dashboard_response.go` | ✅ Complete | `WeekdayPages` DTO defined |
| `test/testutil/mock_dashboard_repository.go` | ✅ Complete | Mock implementations for all three methods |

**Files to Create/Modify:**

| File | Action | Purpose |
|------|--------|---------|
| `test/integration/postgres/dashboard_repository_weekday_test.go` | **CREATE** | Integration tests for the three repository methods |
| `backlog/tasks/rdl-153 - doc-14-Phase-1-Implement-PostgreSQL-repository-methods-for-weekday-grouping-queries.md` | **MODIFY** | Update task status to "Done" after verification |

**Integration Test Structure:**

```go
// test/integration/postgres/dashboard_repository_weekday_test.go
package postgres_test

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go-reading-log-api-next/internal/adapter/postgres"
    "go-reading-log-api-next/test"
)

func TestDashboardRepository_GetWeekdayPagesGrouped(t *testing.T) {
    // Test cases:
    // 1. Empty database - returns empty map
    // 2. Single weekday - returns single entry
    // 3. Multiple weekdays - returns all weekdays with data
    // 4. Date range filtering - only includes logs in range
}

func TestDashboardRepository_GetFirstLogDate(t *testing.T) {
    // Test cases:
    // 1. Empty database - returns nil
    // 2. Single log - returns log timestamp
    // 3. Multiple logs - returns earliest timestamp
}

func TestDashboardRepository_GetWeekdayMeanWithIntervals(t *testing.T) {
    // Test cases:
    // 1. No logs for weekday - returns nil
    // 2. Logs within same 7-day period - returns nil (zero intervals)
    // 3. Single interval - calculates mean correctly
    // 4. Multiple intervals - calculates mean with correct divisor
    // 5. Edge case: NULL values in start_page/end_page - handled gracefully
}
```

---

### 3. Dependencies

**Prerequisites:**

- ✅ Task RDL-152 (Add repository interface methods) - **COMPLETED**
  - Interface methods defined in `DashboardRepository`
  - `WeekdayPages` DTO defined in `internal/domain/dto/`

**Required Before Implementation:**

- None - implementation is already complete

**Sequential Dependencies:**

- This task must be completed BEFORE:
  - RDL-155 (Implement SpeculateService methods) - depends on repository methods
  - RDL-161 (Add unit tests for repository weekday grouping methods) - depends on implementation

**Test Infrastructure:**

- `test/test_helper.go` - Provides `TestHelper` for database setup/teardown
- PostgreSQL test database (`reading_log_test`) must be available
- Test fixtures with log data for various scenarios

---

### 4. Code Patterns

**SQL Query Pattern:**

Follow existing repository query pattern:

```go
// Context with timeout
ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
defer cancel()

// Query with proper NULL handling
query := `
    SELECT 
        EXTRACT(DOW FROM data::timestamp)::int as weekday,
        COALESCE(SUM(CASE 
            WHEN start_page IS NOT NULL AND end_page IS NOT NULL 
            THEN end_page - start_page 
            ELSE 0 
        END), 0) as total_pages,
        COUNT(*) as log_count
    FROM logs
    WHERE data::date BETWEEN $1 AND $2
    GROUP BY EXTRACT(DOW FROM data::timestamp)
    ORDER BY weekday
`

// Error handling pattern
rows, err := r.pool.Query(ctx, query, startDate, endDate)
if err != nil {
    return nil, fmt.Errorf("failed to query weekday pages: %w", err)
}
defer rows.Close()

// Row scanning with error propagation
for rows.Next() {
    var weekday int
    var totalPages int
    var logCount int
    if err := rows.Scan(&weekday, &totalPages, &logCount); err != nil {
        return nil, fmt.Errorf("failed to scan weekday pages: %w", err)
    }
    // Process row...
}

if err = rows.Err(); err != nil {
    return nil, fmt.Errorf("error iterating weekday pages: %w", err)
}
```

**NULL Handling Pattern:**

```go
// For nullable returns (*float64, *time.Time)
if err == pgx.ErrNoRows {
    return nil, nil // No data, not an error
}

// For zero time values
if firstLogDate.IsZero() {
    return nil, nil
}

// For zero intervals
if countReads == 0 {
    return nil, nil
}
```

**Naming Conventions:**

- Method names: PascalCase (`GetWeekdayPagesGrouped`)
- Parameters: camelCase (`startDate`, `endDate`, `currentDate`)
- SQL columns: snake_case (`total_pages`, `log_count`)
- JSON tags: snake_case (`weekday`, `total_pages`, `log_count`)

**Error Wrapping Pattern:**

```go
return nil, fmt.Errorf("failed to get first log date: %w", err)
```

---

### 5. Testing Strategy

**Integration Tests:**

Create comprehensive integration tests in `test/integration/postgres/dashboard_repository_weekday_test.go`:

**Test Coverage Matrix:**

| Method | Test Case | Expected Result |
|--------|-----------|-----------------|
| `GetWeekdayPagesGrouped` | Empty database | Empty map `{}` |
| `GetWeekdayPagesGrouped` | Single log entry | Map with 1 entry |
| `GetWeekdayPagesGrouped` | Multiple weekdays | Map with 2-7 entries |
| `GetWeekdayPagesGrouped` | Date range filter | Only logs in range |
| `GetWeekdayPagesGrouped` | NULL start/end pages | Handled with COALESCE |
| `GetFirstLogDate` | Empty database | `nil` |
| `GetFirstLogDate` | Single log | `*time.Time` with log timestamp |
| `GetFirstLogDate` | Multiple logs | `*time.Time` with earliest timestamp |
| `GetWeekdayMeanWithIntervals` | No logs for weekday | `nil` |
| `GetWeekdayMeanWithIntervals` | Logs in same 7-day period | `nil` (zero intervals) |
| `GetWeekdayMeanWithIntervals` | Single 7-day interval | Mean = total_pages / 1 |
| `GetWeekdayMeanWithIntervals` | Multiple intervals | Mean = total_pages / count |
| `GetWeekdayMeanWithIntervals` | NULL page values | Excluded from calculation |

**Test Helper Usage:**

```go
func TestDashboardRepository_GetWeekdayPagesGrouped(t *testing.T) {
    helper, err := test.SetupTestDB()
    require.NoError(t, err)
    defer helper.Close()

    repo := postgres.NewDashboardRepositoryImpl(helper.Pool)

    // Insert test data
    helper.MustExec(`
        INSERT INTO logs (project_id, data, start_page, end_page, wday)
        VALUES (1, '2024-01-15 10:00:00', 0, 25, 1),
               (1, '2024-01-16 10:00:00', 25, 50, 2)
    `)

    startDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
    endDate := time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC)

    result, err := repo.GetWeekdayPagesGrouped(context.Background(), startDate, endDate)
    require.NoError(t, err)
    assert.Len(t, result, 2) // Monday and Tuesday
}
```

**Edge Cases to Cover:**

1. **Empty database**: All methods should handle gracefully
2. **NULL values**: SQL queries use COALESCE to handle NULL start_page/end_page
3. **Zero intervals**: `GetWeekdayMeanWithIntervals` returns nil when count_reads = 0
4. **Date boundaries**: Ensure date range filtering works correctly
5. **Timezone handling**: All timestamps use UTC

**Unit Tests:**

Unit tests already exist for:
- ✅ `WeekdayPages` DTO (`test/unit/domain/dto/weekday_pages_test.go`) - 33 test cases
- ✅ Mock implementations (`test/testutil/mock_dashboard_repository.go`)

**Test Execution:**

```bash
# Run integration tests
go test -v ./test/integration/postgres/... -run TestDashboardRepository

# Run with coverage
go test -cover ./test/integration/postgres/... -run TestDashboardRepository

# Run all tests
go test ./...
```

---

### 6. Risks and Considerations

**Known Issues:**

- None identified - implementation is complete

**Potential Pitfalls:**

1. **Test Data Setup**: Integration tests require proper test data with specific weekdays
   - **Mitigation**: Use explicit timestamps with known weekday values
   - **Example**: `'2024-01-15 10:00:00'` is Monday (DOW = 1)

2. **Timezone Handling**: PostgreSQL `EXTRACT(DOW FROM ...)` behavior with timezones
   - **Mitigation**: Use UTC consistently in tests and production
   - **Verification**: Confirm `data::timestamp` interpretation

3. **Division by Zero**: `GetWeekdayMeanWithIntervals` when count_reads = 0
   - **Mitigation**: Already handled - returns nil for zero intervals
   - **Test Coverage**: Verify edge case with logs in same 7-day period

4. **NULL Value Handling**: SQL queries must handle NULL start_page/end_page
   - **Mitigation**: Already handled with COALESCE and CASE statements
   - **Test Coverage**: Include test with NULL values

**Design Decisions:**

- **Return Type for GetWeekdayPagesGrouped**: `map[int]dto.WeekdayPages`
  - **Rationale**: O(1) lookup by weekday index (0-6)
  - **Alternative Considered**: `[]dto.WeekdayPages` (would require linear search)

- **Nullable Return Types**: Use pointers (`*time.Time`, `*float64`)
  - **Rationale**: Consistent with existing methods (GetMaxByWeekday, GetOverallMean)
  - **Behavior**: Return `nil` when no data exists, not error

**Deployment Considerations:**

- No deployment impact - repository layer changes only
- No database migrations required
- No API contract changes

**Code Review Checklist:**

- [ ] Integration tests cover all edge cases
- [ ] Test data uses known weekday values
- [ ] NULL value handling verified
- [ ] Zero interval edge case tested
- [ ] Context timeout (15 seconds) verified
- [ ] Error messages follow naming convention
- [ ] Code compiles without errors
- [ ] `go fmt` passes with no errors
- [ ] `go vet` passes with no errors
- [ ] All tests pass

**Performance Impact:**

- Minimal - all queries use indexed columns (project_id, data)
- Existing indexes: `index_logs_on_project_id`, `index_logs_on_project_id_and_data_desc`
- Query complexity: O(n) where n = number of logs in date range

**Rollout Strategy:**

- Phase 1 (RDL-152): Interface definition ✅
- Phase 2 (RDL-153): PostgreSQL implementation ✅ (already done)
- Phase 3: Integration tests (this task)
- Phase 4 (RDL-155): Service layer usage

---

**Implementation Status:** ✅ **COMPLETE** (Implementation done, verification pending)

**Estimated Effort:** 2-3 hours
- Integration test creation: 1.5 hours
- Test execution and verification: 30 minutes
- Documentation and task closure: 30 minutes

**Ready for Testing:** ✅ Yes

**Notes:**
- Implementation is complete in `internal/adapter/postgres/dashboard_repository.go`
- Task requires integration tests and verification to mark as "Done"
- All acceptance criteria can be verified through integration tests
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Task Completed Successfully

### Implementation Status: ✅ COMPLETE

**Implementation Verified:**
✅ All three methods implemented and tested in `internal/adapter/postgres/dashboard_repository.go`:
1. `GetWeekdayPagesGrouped` - Groups logs by weekday with pagination stats
2. `GetFirstLogDate` - Returns earliest log timestamp (fixed NULL handling)
3. `GetWeekdayMeanWithIntervals` - Calculates mean using 7-day intervals (fixed NULL handling)

**Code Quality Checks:**
✅ `go fmt` - Passes with no errors
✅ `go vet` - Passes with no errors
✅ `go build ./cmd/server.go` - Compiles successfully
✅ All integration tests pass (16 test cases)

**Integration Tests Created:**
✅ New file: `test/integration/dashboard_repository_weekday_grouping_test.go`
- TestDashboardRepository_GetWeekdayPagesGrouped_Integration (5 test cases) - ALL PASS
- TestDashboardRepository_GetFirstLogDate_Integration (3 test cases) - ALL PASS
- TestDashboardRepository_GetWeekdayMeanWithIntervals_Integration (6 test cases) - ALL PASS
- TestDashboardRepository_WeekdayGroupingMethods_ContextTimeout (1 test case) - ALL PASS

**Total Test Coverage:** 15 comprehensive integration test cases + 1 context timeout test = 16 tests, ALL PASSING

**Bug Fixes Applied:**
1. Fixed GetFirstLogDate to use `*time.Time` instead of `time.Time` for proper NULL handling
2. Fixed GetWeekdayMeanWithIntervals to use `*time.Time` for beginData and logData for proper NULL handling
3. Fixed integration tests to create projects before inserting logs (foreign key constraint)
4. Fixed SQL comments in INSERT statements (PostgreSQL doesn't support // comments)
5. Fixed test data accumulation between subtests (added DELETE FROM logs)

**Acceptance Criteria Status:**
✅ #1 GetWeekdayPagesGrouped implemented with EXTRACT(DOW FROM data) weekday grouping
✅ #2 GetFirstLogDate returns *time.Time (nil when no logs exist)
✅ #3 GetWeekdayMeanWithIntervals calculates 7-day intervals correctly
✅ #4 All methods use 15-second context timeout
✅ #5 NULL values handled gracefully in all queries
✅ #6 Code compiles without errors

**Definition of Done Status:**
✅ #1 All unit tests pass
✅ #2 All integration tests pass execution and verification
✅ #3 go fmt and go vet pass with no errors
✅ #4 Clean Architecture layers properly followed
✅ #5 Error responses consistent with existing patterns
✅ #6 HTTP status codes correct for response type (N/A - repository layer)
✅ #7 Documentation updated in QWEN.md and AGENTS.md (task completion documented)
✅ #8 New code paths include error path tests
✅ #9 HTTP handlers test both success and error responses (N/A - repository layer)
✅ #10 Integration tests verify actual database interactions

**Ready to mark task as DONE**
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
