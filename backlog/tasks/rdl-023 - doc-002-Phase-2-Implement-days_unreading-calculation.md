---
id: RDL-023
title: '[doc-002 Phase 2] Implement days_unreading calculation'
status: To Do
assignee:
  - catarina
created_date: '2026-04-03 14:03'
updated_date: '2026-04-03 20:18'
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
- [ ] #1 days_unreading = today minus last log date or started_at
- [ ] #2 If no logs, use started_at date
- [ ] #3 If neither logs nor started_at exist, return 0
- [ ] #4 Result is non-negative integer
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Implement the `days_unreading` calculation method in Go that matches Rails behavior. The calculation is: `days_unreading = (Date.today - last_log_or_started_at).days`.

**Key Changes Required:**
- Modify the `CalculateDaysUnreading` method in `internal/domain/models/project.go` to return `0` instead of `nil` when no logs and no `started_at` exist
- Add validation to ensure the result is non-negative (handle future dates gracefully)
- Ensure the method is called from the appropriate repository/handler layer to populate the `days_unreading` field

**Technical Implementation:**
- The calculation uses date-only comparison to match Rails `Date.today` behavior
- If logs exist, use the most recent log's `data` field (YYYY-MM-DD format)
- If no logs but `started_at` exists, use `started_at`
- If neither exists, return 0 (not nil)
- Calculate difference in days using date arithmetic

**Why this approach:**
- Matches Rails implementation exactly: `days_unreading = (Date.today - base_data).to_i`
- Edge case handling ensures no nil values in response
- Returns non-negative integer as required
- Follows existing patterns in the codebase (similar to progress/status calculations)

### 2. Files to Modify

| File | Action | Reason |
|------|--------|--------|
| `internal/domain/models/project.go` | Modify | Fix `CalculateDaysUnreading()` to return 0 instead of nil when no logs and no started_at; add non-negative validation |
| `internal/domain/models/project_test.go` | Modify | Update existing test `TestProject_CalculateDaysUnreading_NoLogsNoStartedAt` to expect 0 instead of nil; add test for future dates |
| `internal/adapter/postgres/project_repository.go` | Verify | Confirm CalculateDaysUnreading is called where needed; if database field is used directly, ensure calculation is done |
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
- PostgreSQL repository with eager-loaded logs support (already implemented)

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
    // ... other fields including DaysUnread *int
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

**Key changes to make:**
```go
// Current: Returns nil when no logs and no started_at
// Required: Return 0 instead of nil

if len(logs) == 0 && p.StartedAt == nil {
    return &[]int{0}[0] // or create a zero value pointer directly
}
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
- Verify returned pointer is never nil

### 6. Risks and Considerations

**Potential pitfalls:**
1. **Return nil vs 0**: Current implementation returns nil when no data available, but task requires 0
   - Solution: Modify the edge case return value, ensure consistent non-nil behavior

2. **Future dates**: If last_log is in the future, calculation could return negative
   - Solution: Clamp to 0 if negative (use `max(0, days)` pattern)

3. **Timezone handling**: Rails Date.today vs time.Now() timezone differences
   - Solution: Use date-only comparison with UTC timezone matching existing pattern

4. **Database field vs calculation**: The repository might be using the database field directly instead of calculating
   - Solution: Verify the repository calls CalculateDaysUnreading; if not, add the call

**Design decisions:**
- **Return type**: `*int` (pointer to int) to allow 0 as a valid value while still being able to return nil if needed for other cases
- **Non-negative guarantee**: Always return >= 0, even for future dates
- **Consistency**: Follow the same pattern as CalculateStatus which uses the same method

**Edge cases to handle:**
1. No logs, no started_at → return 0
2. Log with no data field → skip and check next log
3. All logs have no data field → use started_at if available, else return 0
4. Future dates → return 0 (not negative)
5. Very old dates → return large positive integer (no special handling needed)

**Database field consideration:**
- The database has a `days_unread` column in the projects table
- The repository currently reads this field from the database
- QUESTION: Should we be calculating this on-the-fly or using database field?
- Based on RDL-022 pattern and PRD Decision 1, calculations should be done in Go, not stored in DB
- If database field exists, it may need to be set by the application when saving, OR the application should calculate it on-read
- Based on status calculation pattern, it seems calculations are done on-read, so we should follow the same pattern

**Recommendation:**
- Verify if the database `days_unread` field should be populated or if calculation on-read is preferred
- If calculation on-read, modify repository to call CalculateDaysUnreading when building responses
- If database field, ensure it's set correctly when writing to database
<!-- SECTION:PLAN:END -->

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
<!-- DOD:END -->
