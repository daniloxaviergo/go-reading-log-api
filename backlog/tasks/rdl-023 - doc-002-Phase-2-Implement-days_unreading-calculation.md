---
id: RDL-023
title: '[doc-002 Phase 2] Implement days_unreading calculation'
status: To Do
assignee:
  - thomas
created_date: '2026-04-03 14:03'
updated_date: '2026-04-03 20:52'
labels:
  - phase-2
  - derived-calculation
  - date-calculation
dependencies: []
references:
  - >-
    PRD Section: Technical Decisions - Decision 1: Derived Calculations
    Implementation
  - 'PRD Section: Validation Rules - days_unreading rule'
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement days_unreading calculation in Go matchingRails: days_unreading = (Date.today - last_log_or_started_at).days. Handle edge cases: if no logs, use started_at; if no logs and no started_at, return 0. Method should return non-negative integer.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 days_unreading = today minus last log date or started_at
- [x] #2 If no logs, use started_at date
- [x] #3 If neither logs nor started_at exist, return 0
- [x] #4 Result is non-negative integer
- [ ] #5 #1 All unit tests pass
- [ ] #6 #2 All integration tests pass
- [ ] #7 #3 go fmt and go vet pass
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Implement the `days_unreading` calculation method in Go that matches Rails behavior. The calculation is: `days_unreading = (Date.today - last_log_or_started_at).days`.

**Key Findings from Codebase Review:**
1. The Rails database schema (`rails-app/db/schema.rb`) shows the `projects` table does NOT have computed fields like `days_unread`, `status`, `progress`, etc.
2. These fields are calculated in Rails via ActiveModelSerializer virtual methods
3. The current Go PostgreSQL repository incorrectly reads these fields from the database in its SELECT queries
4. The `CalculateDaysUnreading` method exists in `internal/domain/models/project.go` but has incorrect edge case handling

**Key Changes Required:**
1. Fix `CalculateDaysUnreading` in `internal/domain/models/project.go` to:
   - Return `0` instead of `nil` when no logs and no `started_at` exist
   - Clamp negative results to 0 for future dates
2. Update PostgreSQL repository queries to NOT read computed fields from database
3. Ensure `CalculateDaysUnreading` is called when building response DTOs in the repository

**Technical Implementation:**
- The calculation uses date-only comparison to match Rails `Date.today` behavior
- If logs exist, use the most recent log's `data` field (YYYY-MM-DD format)
- If no logs but `started_at` exists, use `started_at`
- If neither exists, return 0 (not nil)
- Calculate difference in days using date arithmetic
- Clamp to 0 if result would be negative (future dates)

**Why this approach:**
- Matches Rails implementation exactly: `days_unreading = (Date.today - base_data).to_i`
- Edge case handling ensures no nil values in response
- Returns non-negative integer as required
- Follows existing patterns in the codebase (similar to progress/status calculations)
- Ensures calculated fields are computed on-read, not read from database

### 2. Files to Modify

| File | Action | Reason |
|------|--|--|
| `internal/domain/models/project.go` | Modify | Fix `CalculateDaysUnreading()` to return 0 instead of nil when no logs and no started_at; add non-negative validation |
| `internal/domain/models/project_test.go` | Modify | Update existing test `TestProject_CalculateDaysUnreading_NoLogsNoStartedAt` to expect 0 instead of nil; add test for future dates |
| `internal/adapter/postgres/project_repository.go` | Modify | Remove computed fields (progress, status, logs_count, days_unread, median_day, finished_at) from SELECT queries; add calls to CalculateDaysUnreading() when building responses; ensure fields are calculated on-read |
| `internal/api/v1/handlers/projects_handler.go` | Verify | Ensure days_unreading is included in responses from repository |
| `internal/api/v1/handlers/logs_handler.go` | Verify | Ensure days_unreading is included in project eager-loaded in log responses |

### 3. Dependencies

**Prerequisites:**
1. Task RDL-022 completed (status determination logic) - already done status
2. Existing `CalculateDaysUnreading` method exists but needs edge case fix
3. Access to logs through ProjectWithLogs structure
4. Configuration package available (already completed in RDL-022)

**Required Existing Infrastructure:**
- `internal/domain/models` package with Project model
- `internal/domain/dto` package with LogResponse structure
- `internal/config` package with status range configuration (already implemented)
- PostgreSQL repository with eager-loaded logs support (already implemented, needs query updates)

**No additional setup required** - all prerequisites are in place.

### 4. Code Patterns

**Existing conventions to follow:**

1. **Project model pattern** (from `internal/domain/models/project.go`):
```go
type Project struct {
    ctx        context.Context
    ID         int64      `json:"id"`
    Name       string     `json:"name"`
    TotalPage  int        `json:"total_page"`
    StartedAt  *time.Time `json:"started_at"`
    Page       int        `json:"page"`
    // ... other fields
}

// CalculateDaysUnreading calculates the number of days since the last reading activity
func (p *Project) CalculateDaysUnreading(logs []*dto.LogResponse) *int {
    // Implementation here
}
```

2. **Helper function pattern** (already exists in project.go):
```go
func stringPtr(s string) *string {
    return &s
}
```

3. **Date handling pattern** (from progress calculation):
```go
nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
lastReadDate := time.Date(lastRead.Year(), lastRead.Month(), lastRead.Day(), 0, 0, 0, 0, time.UTC)
diff := nowDate.Sub(lastReadDate)
days := int(diff.Hours() / 24)
```

**Key changes to make in CalculateDaysUnreading:**
```go
// Fix edge case: return 0 instead of nil
if len(logs) == 0 && p.StartedAt == nil {
    zero := 0
    return &zero
}

// Clamp to 0 if negative (future dates)
if days < 0 {
    return &zero
}
```

**Repository changes to make:**
- Remove computed fields from SELECT queries:
  ```go
  // Current (incorrect):
  SELECT id, name, total_page, started_at, page, reinicia, progress, status, logs_count, days_unread, median_day, finished_at
  // Correct:
  SELECT id, name, total_page, started_at, page, reinicia
  ```
- Add calculation calls when building responses:
  ```go
  // Calculate days_unreading
  daysUnread := domainProject.CalculateDaysUnreading(logs)
  project.DaysUnread = daysUnread
  ```

### 5. Testing Strategy

**Unit tests to update/add** (in `internal/domain/models/project_test.go`):

1. `TestProject_CalculateDaysUnreading_NoLogsNoStartedAt` - **Update existing test** to expect 0 instead of nil
2. `TestProject_CalculateDaysUnreading_FutureDate` - Add test for when last read date is in the future (should return 0)
3. `TestProject_CalculateDaysUnreading_NegativeResult` - Add test to ensure negative values don't occur
4. `TestProject_CalculateDaysUnreading_WithLogs` - Verify existing log-based calculation still works
5. `TestProject_CalculateDaysUnreading_StartedAtOnly` - Verify started_at fallback still works

**Test data setup:**
- Project with no logs and no started_at → expect 0
- Project with logs today → expect 0 days
- Project with logs in future → expect 0 (not negative)
- Project with no logs but started_at → expect_days_since_started_at
- Project with logs from various dates → expect_days_since_most_recent_log

**Approach:**
- Use `testing.T` for assertions
- Compare returned days value with expected values
- Test all edge cases: no data, future dates, boundary conditions
- Verify returned pointer is not nil for valid cases

**Integration tests to verify:**
- Repository correctly calculates days_unreading instead of reading from database
- Handler returns correct JSON response with calculated days_unreading
- Full response includes all derived fields consistently

### 6. Risks and Considerations

**Potential pitfalls:**
1. **Return nil vs 0**: Current implementation returns nil when no data available, but task requires 0
   - Solution: Modify the edge case return value, ensure consistent non-nil behavior

2. **Future dates**: If last_log is in the future, calculation could return negative
   - Solution: Clamp to 0 if negative (use `max(0, days)` pattern)

3. **Timezone handling**: Rails Date.today vs time.Now() timezone differences
   - Solution: Use date-only comparison with UTC timezone matching existing pattern

4. **Database field assumption**: Repository currently reads computed fields from database (incorrect)
   - Solution: Remove computed fields from SELECT queries, calculate on-read

5. **Consistency with status calculation**: `CalculateStatus()` uses `CalculateDaysUnreading()`, ensure both are working together
   - Solution: Test status determination with the updated days_unreading calculation

**Design decisions:**
- **Return type**: `*int` (pointer to int) to allow 0 as a valid value while still being able to handle edge cases
- **Non-negative guarantee**: Always return >= 0, even for future dates
- **Consistency**: Follow the same pattern as CalculateStatus which uses the same method
- **On-read calculation**: Calculate fields when building responses, not read from database (matches Rails pattern)

**Edge cases to handle:**
1. No logs, no started_at → return 0
2. Log with no data field → skip and check next log
3. All logs have no data field → use started_at if available, else return 0
4. Future dates → return 0 (not negative)
5. Very old dates → return large positive integer (no special handling needed)
6. All logs with valid data → use most recent log

**Database field clarification:**
- The database schema shows no computed columns exist
- These fields are calculated in Rails via ActiveModelSerializer
- The Go implementation should follow the same pattern:
  - Calculate fields on-read in the repository/handler layer
  - Do NOT read from database columns
  - Return calculated values in DTOs

**Recommendation:**
1. First fix `CalculateDaysUnreading` to handle edge cases correctly
2. Update PostgreSQL repository to remove computed fields from SELECT queries
3. Add calculation calls when building responses
4. Update tests to verify on-read calculation
5. Verify all acceptance criteria pass

**Note:** This task should also consider the related issue in the PostgreSQL repository where computed fields are incorrectly being read from the database. The SELECT queries in `project_repository.go` should be updated to remove `progress, status, logs_count, days_unread, median_day, finished_at` since these don't exist in the database schema and should be calculated in Go.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Updated task status to In Progress after reviewing implementation. Current CalculateDaysUnreading implementation has two issues: 1) Returns nil when no logs and no started_at, but task requires 0. 2) Does not clamp negative results to 0 for future dates. Need to fix both issues in project.go and update corresponding tests.

Implementation progress:

1. Fixed CalculateDaysUnreading in project.go: edge case returns 0 instead of nil, clamps negative values to 0

2. Updated tests in project_test.go: added future date tests, updated edge case expectation

3. Updated PostgreSQL repository: removed computed fields from SELECT queries, added calculation calls, added nil check for startedAt

4. All 157 tests pass, go fmt and go vet pass, build succeeds

5. Bug fix: fixed nil pointer dereference in GetWithLogs when startedAt is NULL
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [ ] #2 All integration tests pass use testing-expert subagent for test execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Database queries optimized with proper indexes
- [ ] #8 Documentation updated in QWEN.md
- [ ] #9 New code paths include error path tests
- [ ] #10 HTTP handlers test both success and error responses
- [ ] #11 Integration tests verify actual database interactions
- [ ] #12 Tests use testing-expert subagent for test execution and verification
- [ ] #13 All unit tests pass - 27 tests in models package, 157 total tests pass
- [ ] #14 All integration tests pass - 21 tests in integration package pass
- [ ] #15 go fmt and go vet pass with no errors
- [ ] #16 Clean Architecture layers properly followed - domain models, DTOs, repository adapters
- [ ] #17 Error responses consistent with existing patterns
- [ ] #18 HTTP status codes correct for response type
- [ ] #19 Database queries optimized with proper indexes - only base fields selected
- [ ] #20 New code paths include error path tests
- [ ] #21 HTTP handlers test both success and error responses
- [ ] #22 Integration tests verify actual database interactions
- [ ] #23 Tests use testing-expert subagent for test execution and verification
<!-- DOD:END -->
