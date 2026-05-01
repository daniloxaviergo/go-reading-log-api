---
id: doc-012
title: 'PRD: Fix Faults Calculation Logic - Rails Parity for Dashboard API'
type: other
created_date: '2026-05-01 14:53'
---


---
id: doc-009
title: 'PRD: Fix Faults Calculation Logic - Rails Parity for Dashboard API'
type: other
created_date: '2026-05-01'
---

# Project Requirements Document

# Executive Summary

**Critical Bug Fix**: The current Go implementation of the `/v1/dashboard/echart/faults.json` endpoint has a fundamental logic error. It counts log entries instead of counting days with zero pages read, which is what the Rails application does.

**Why Necessary**: Dashboard statistics showing fault percentage are completely inaccurate. Users see incorrect fault metrics that don't match the Rails application behavior.

**Scope**: Fix `GetFaultsByDateRange` and `GetWeekdayFaults` repository methods to match Rails logic exactly.

**Impact**: HIGH - This is a blocker for Phase 2 dashboard completion. All existing tests validate incorrect behavior and must be rewritten.

**Timeline**: Must be fixed before Phase 2 completion.

---

# Key Requirements

| Requirement | Priority | Status |
|-------------|----------|--------|
| Fix `GetFaultsByDateRange` to count days with zero pages | P1 (Blocker) | TODO |
| Fix `GetWeekdayFaults` to count zero-page days by weekday | P1 (Blocker) | TODO |
| Use inclusive date range (BETWEEN $1 AND $2) | P1 (Blocker) | TODO |
| Use server local time for date boundaries | P1 (Blocker) | TODO |
| Rewrite all unit tests with correct expectations | P1 (Blocker) | TODO |
| Rewrite all integration tests with correct expectations | P1 (Blocker) | TODO |
| Add Rails comparison tests for validation | P2 (Must-have) | TODO |

---

# Technical Decisions

## Decision 1: Fault Definition

**Decision**: A "fault" = a day with zero pages read (sum of `read_pages` = 0 for that day).

**Rationale**:
- Rails implementation clearly shows: `next 1 if read_pages.sum.zero?`
- For each day in date range, check if any reading occurred
- If no pages were read that day, count as 1 fault
- Matches user expectation: "faults" = missed reading days

**Implementation**:
```sql
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

## Decision 2: Date Range Inclusivity

**Decision**: Use inclusive date range on both ends (BETWEEN $1 AND $2).

**Rationale**:
- Rails uses Ruby range: `(date_start..date_end)` which is inclusive
- Query must match: `data::date BETWEEN $1 AND $2`
- Ensures edge case dates are handled identically

**Implementation**:
```go
func GetDateRangeLast30Days() (start, end time.Time) {
    end = dto.GetToday()  // Today at midnight (inclusive)
    start = end.AddDate(0, 0, -30)  // 30 days ago (inclusive)
    return start, end
}
```

## Decision 3: Timezone Handling

**Decision**: Use server local time for date casting, matching current Go implementation.

**Rationale**:
- Questionnaire response: "Server local time"
- Current implementation already uses `dto.GetToday()` which returns `time.Now().Truncate(24 * time.Hour)`
- Rails uses `Time.zone.today` but server TZ should be configured to match

**Implementation**:
```go
// Ensure server TZ matches Rails TZ
// In .env: TZ=America/Sao_Paulo (or your timezone)

func GetToday() time.Time {
    return time.Now().Truncate(24 * time.Hour)
}
```

## Decision 4: Test Strategy

**Decision**: Rewrite ALL existing tests with correct expectations.

**Rationale**:
- Questionnaire response: "Rewrite all tests"
- Current tests validate incorrect behavior
- Better to have 100% correct tests than mix of old/new
- Risk: More work upfront, but prevents confusion later

**Implementation**:
```go
// Example test case
func TestGetFaultsByDateRange_ZeroPagesDay(t *testing.T) {
    // Setup: Create logs for 5 days, skip 2 days
    // Expected: 2 faults (the 2 days with no logs)
    // NOT: 5 faults (counting log entries)
}
```

---

# Acceptance Criteria

## Functional Acceptance Criteria

### AC-FAULT-001: Faults Count Calculation

**Given**: A date range of 30 days  
**When**: The `GetFaultsByDateRange` method is called  
**Then**:
- Returns count of days where `sum(read_pages) = 0`
- Days with one or more logs where `sum(end_page - start_page) > 0` are NOT counted
- Days with logs but `sum(end_page - start_page) = 0` ARE counted as faults
- Days with no logs at all ARE counted as faults
- Result matches Rails `V1::Dashboard::Faults` exactly

**Test Case**:
```go
func TestGetFaultsByDateRange_RailsParity(t *testing.T) {
    // Create 30-day period
    // Day 1: 10 pages read → 0 faults
    // Day 2: 0 pages read (no logs) → 1 fault
    // Day 3: 0 pages read (log with start=0, end=0) → 1 fault
    // Day 4: 20 pages read → 0 faults
    // Expected: 2 faults total
}
```

---

### AC-FAULT-002: Weekday Faults Calculation

**Given**: A date range spanning multiple weeks  
**When**: The `GetWeekdayFaults` method is called  
**Then**:
- Returns map of weekday (0-6) → count of zero-page days
- Each weekday shows how many times that day had no reading
- All 7 weekdays present in result (default 0 if none)
- Result matches Rails weekday grouping logic

**Test Case**:
```go
func TestGetWeekdayFaults_RailsParity(t *testing.T) {
    // 4-week period (28 days = 4 of each weekday)
    // Skip reading on all Mondays (4 faults for weekday 1)
    // Read every other day
    // Expected: map[1:4, 0:0, 2:0, 3:0, 4:0, 5:0, 6:0]
}
```

---

### AC-FAULT-003: Date Boundary Handling

**Given**: Edge case dates (start_date == end_date)  
**When**: Methods are called with same start and end date  
**Then**:
- Returns 1 fault if no logs on that day
- Returns 0 faults if logs exist with pages > 0
- Handles single-day ranges correctly

**Test Case**:
```go
func TestGetFaultsByDateRange_SingleDay(t *testing.T) {
    today := dto.GetToday()
    
    // No logs today
    faults, _ := repo.GetFaultsByDateRange(ctx, today, today)
    assert.Equal(t, 1, faults.FaultCount)
    
    // Add log with 10 pages
    // faults, _ = repo.GetFaultsByDateRange(ctx, today, today)
    // assert.Equal(t, 0, faults.FaultCount)
}
```

---

### AC-FAULT-004: NULL Value Handling

**Given**: Logs with NULL start_page or end_page  
**When**: Faults are calculated  
**Then**:
- NULL values treated as 0 pages read
- Day with NULL pages counts as fault
- SQL CASE statement handles NULL gracefully

**Test Case**:
```go
func TestGetFaultsByDateRange_NULLValues(t *testing.T) {
    // Create log with start_page = NULL, end_page = 10
    // Expected: Counts as 0 pages (NULL treated as 0)
    // Day counts as 1 fault
}
```

---

## Non-Functional Acceptance Criteria

### NFA-FAULT-001: Query Performance

| Criterion | Target | Measurement |
|-----------|--------|-------------|
| Query execution time | < 100ms | For 30-day range |
| Query complexity | O(days) | Uses generate_series |
| Index usage | YES | Uses data::date index |

### NFA-FAULT-002: Code Quality

| Criterion | Target | Measurement |
|-----------|--------|-------------|
| Test coverage | 100% | For fault calculation logic |
| Code duplication | 0% | No duplicated SQL logic |
| Documentation | Complete | All public functions documented |

---

# Files to Modify

## Files Requiring Changes

| File Path | Change Type | Priority | Rationale |
|-----------|-------------|----------|-----------|
| `internal/adapter/postgres/dashboard_repository.go` | Modify | P1 | Fix `GetFaultsByDateRange` SQL query |
| `internal/adapter/postgres/dashboard_repository.go` | Modify | P1 | Fix `GetWeekdayFaults` SQL query |
| `test/unit/dashboard_repository_test.go` | Rewrite | P1 | Update all fault-related tests |
| `test/dashboard_integration_test.go` | Rewrite | P1 | Update integration test expectations |
| `internal/api/v1/handlers/dashboard_handler_test.go` | Rewrite | P1 | Update mock expectations |
| `api/v1/routes_test.go` | Rewrite | P1 | Update mock implementations |
| `test/unit/dashboard_response_test.go` | Modify | P2 | Verify JSON serialization unchanged |

## Files Created

| File Path | Purpose | Priority |
|-----------|---------|----------|
| `test/rails_comparison/faults_comparison_test.go` | Side-by-side Rails/Go validation | P2 |
| `docs/faults-calculation-explanation.md` | Technical documentation | P2 |

---

# Validation Rules

## SQL Query Validation

```sql
-- Correct pattern for fault counting
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

**Validation Checks**:
- ✅ Uses `generate_series` to create all dates in range
- ✅ LEFT JOIN ensures days without logs are counted
- ✅ CASE statement handles NULL values
- ✅ BETWEEN ensures inclusive date range
- ✅ Aggregation uses `SUM(read_pages)` not COUNT(*)

---

# Out of Scope

The following items are explicitly **OUT OF SCOPE**:

1. **Other dashboard endpoints**: Only faults calculation is being fixed
2. **Database schema changes**: No table modifications required
3. **API response format changes**: JSON structure remains identical
4. **UserConfig changes**: max_faltas configuration unchanged
5. **Performance optimization**: Current query complexity acceptable
6. **Caching layer**: No caching changes required

---

# Implementation Checklist

## Phase 1: Repository Implementation (Blocker)

- [ ] **RI-001**: Update `GetFaultsByDateRange` SQL query
  - [ ] Implement CTE with daily_read aggregation
  - [ ] Add generate_series for all dates
  - [ ] LEFT JOIN to count zero-page days
  - [ ] Add comprehensive inline comments

- [ ] **RI-002**: Update `GetWeekdayFaults` SQL query
  - [ ] Implement same daily aggregation logic
  - [ ] Group by weekday (EXTRACT DOW)
  - [ ] Ensure all 7 weekdays in result
  - [ ] Add comprehensive inline comments

- [ ] **RI-003**: Verify context timeout handling
  - [ ] Ensure 15-second timeout applied
  - [ ] Check error wrapping consistent

## Phase 2: Unit Tests (Blocker)

- [ ] **UT-001**: Rewrite `TestGetFaultsByDateRange`
  - [ ] Test case: No logs in range (all days faults)
  - [ ] Test case: All days have reading (0 faults)
  - [ ] Test case: Mixed scenario (some faults)
  - [ ] Test case: Single day range
  - [ ] Test case: NULL values handling

- [ ] **UT-002**: Rewrite `TestGetWeekdayFaults`
  - [ ] Test case: All weekdays represented
  - [ ] Test case: Some weekdays missing
  - [ ] Test case: Zero faults for all days
  - [ ] Test case: All days are faults

- [ ] **UT-003**: Update mock repositories
  - [ ] Update `MockDashboardRepository` expectations
  - [ ] Update `MockDashboardRepositoryForProjects` expectations

## Phase 3: Integration Tests (Blocker)

- [ ] **IT-001**: Update `TestDashboardFaultsChart_Integration`
  - [ ] Fix test data setup
  - [ ] Update expected fault counts
  - [ ] Verify percentage calculation

- [ ] **IT-002**: Update `TestDashboardWeekdayFaults_Integration`
  - [ ] Fix test data setup
  - [ ] Update expected weekday distribution
  - [ ] Verify radar chart data

- [ ] **IT-003**: Add edge case integration tests
  - [ ] Test with empty database
  - [ ] Test with single log entry
  - [ ] Test with month boundary dates

## Phase 4: Rails Comparison Tests (Should-have)

- [ ] **RCT-001**: Create Rails comparison test harness
  - [ ] Setup Rails test data fixtures
  - [ ] Execute Rails faults calculation
  - [ ] Execute Go faults calculation
  - [ ] Compare results

- [ ] **RCT-002**: Add comparison test cases
  - [ ] 30-day range with random gaps
  - [ ] 6-month range for weekday faults
  - [ ] Edge case: Leap year February

## Phase 5: Documentation (Should-have)

- [ ] **DOC-001**: Create faults calculation documentation
  - [ ] Explain Rails logic in detail
  - [ ] Show SQL query breakdown
  - [ ] Provide examples with diagrams

- [ ] **DOC-002**: Update API documentation
  - [ ] Clarify fault definition
  - [ ] Add calculation examples

---

# Stakeholder Alignment

| Stakeholder | Responsibility | Verification |
|-------------|----------------|--------------|
| **Engineering Lead** | Approve SQL query changes | Review technical decisions |
| **Backend Developer** | Implement repository fixes | Execute implementation checklist |
| **QA Team** | Validate Rails parity | Run comparison tests |
| **Product Owner** | Accept fix as complete | Verify dashboard accuracy |
| **DevOps** | Monitor query performance | Check execution time metrics |

---

# Traceability Matrix

| Requirement ID | User Story | Acceptance Criteria | Test File | Status |
|----------------|------------|---------------------|-----------|--------|
| REQ-FAULT-001 | Count days with zero pages read | AC-FAULT-001 | test/unit/dashboard_repository_test.go | TODO |
| REQ-FAULT-002 | Group faults by weekday | AC-FAULT-002 | test/unit/dashboard_repository_test.go | TODO |
| REQ-FAULT-003 | Handle date boundaries correctly | AC-FAULT-003 | test/dashboard_integration_test.go | TODO |
| REQ-FAULT-004 | Handle NULL values gracefully | AC-FAULT-004 | test/unit/dashboard_repository_test.go | TODO |
| REQ-FAULT-005 | Match Rails implementation exactly | AC-FAULT-001 | test/rails_comparison/faults_comparison_test.go | TODO |

---

# Validation

## Code Quality Standards
- [ ] Go 1.25.7 compatible
- [ ] SQL queries optimized with EXPLAIN ANALYZE
- [ ] Linting passes (`go vet ./...`)
- [ ] Formatting correct (`go fmt ./...`)
- [ ] No SQL injection vulnerabilities

## Technical Feasibility
- [ ] PostgreSQL `generate_series` available (PG 13+)
- [ ] CTE syntax supported
- [ ] Date casting (`::date`) works with VARCHAR column
- [ ] Connection pooling handles concurrent queries

## User Needs Alignment
- [ ] Dashboard statistics now accurate
- [ ] Fault percentage reflects actual missed days
- [ ] Rails and Go outputs match exactly
- [ ] No breaking changes to API contract

## Test Coverage Validation
- [ ] All fault calculation paths tested
- [ ] Edge cases covered (NULL, empty range, single day)
- [ ] Rails comparison tests pass
- [ ] Integration tests validate real database behavior

---

# Ready for Implementation

## Approval Status: ✅ **READY FOR IMPLEMENTATION**

This PRD has been:
- ✅ **Validated** against Rails implementation for accuracy
- ✅ **Clarified** through questionnaire to resolve ambiguities
- ✅ **Prioritized** as Blocker for Phase 2 completion
- ✅ **Tested** with acceptance criteria that are objective and measurable
- ✅ **Documented** with exact SQL queries and implementation steps
- ✅ **Reviewed** with palha subagent for technical accuracy

## Prerequisites for Starting Implementation

Before beginning this fix, ensure:

1. **Development environment ready**:
   - [ ] Go 1.25.7 installed
   - [ ] PostgreSQL running (13+)
   - [ ] Test database created (`reading_log_test`)
   - [ ] `.env` file configured with TZ setting

2. **Rails app available for comparison** (optional but recommended):
   - [ ] Rails app running or accessible
   - [ ] Test fixtures can be shared between Rails and Go

3. **Stakeholder sign-off**:
   - [ ] Engineering lead approved SQL changes
   - [ ] QA team prepared to run comparison tests

## Implementation Start Command

```bash
# Verify environment
make test-setup

# Run existing tests to establish baseline
go test -v ./internal/adapter/postgres/... -run TestGetFaultsByDateRange

# Begin implementation
cd internal/adapter/postgres
# Edit dashboard_repository.go

# Run tests after each change
go test -v ./internal/adapter/postgres/... -run TestGetFaultsByDateRange
```

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| SQL query performance degradation | Low | Medium | EXPLAIN ANALYZE before/after |
| Tests fail unexpectedly | Medium | Low | Comprehensive test coverage |
| Rails parity not achieved | Low | High | Comparison tests validate |
| Breaking existing functionality | Low | High | Integration tests catch issues |

## Rollback Plan

If issues arise:
1. Revert `dashboard_repository.go` changes
2. Restore backed-up test files
3. Investigate failures in isolated branch
4. Re-implement with fixes

---

*PRD Version: 1.0*
*Created: 2026-05-01*
*Last Updated: 2026-05-01*
*Author: PRD Refinement Specialist*
*Status: Ready for Implementation*