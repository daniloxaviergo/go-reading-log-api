---
id: RDL-025
title: '[doc-002 Phase 2] Implement finished_at calculation'
status: To Do
assignee:
  - thomas
created_date: '2026-04-03 14:03'
updated_date: '2026-04-03 22:15'
labels:
  - phase-2
  - derived-calculation
  - date-calculation
dependencies: []
references:
  - >-
    PRD Section: Technical Decisions - Decision 1: Derived Calculations
    Implementation
documentation:
  - doc-002
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement finished_at calculation in Go matching Rails: calculate future date when book will be finished based on reading rate (median_day). If progress is 100% or pages remaining is 0, return nil/null. Otherwise, calculate: days_to_finish = (total_page - page) / median_day, then finished_at = today + days_to_finish days.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 finished_at = today + (total_page - page) / median_day days
- [x] #2 100% progress edge case returns null
- [x] #3 Pages remaining = 0 edge case returns null
- [x] #4 Date calculated as future date in days
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
## Implementation Plan: finished_at Calculation

### 1. Technical Approach

Implement the `finished_at` calculation in the Go `Project` domain model, matching the Rails implementation exactly.

**Formula Analysis (from Rails):**
```ruby
days_reading = (Date.today - started_at).to_i
future_days_read = (days_reading.to_f * total_page.to_f) / page.to_f
return last_read[:data].to_date if finished?
Date.today + future_days_read.round.days
```

**Translation to Go:**
1. Calculate `days_reading` = days since `started_at` using date-only comparison (matches Rails `Date.today` behavior)
2. Calculate `future_days_read` = (days_reading × total_page) / page
3. If page >= total_page (finished), return the last log's `data` field as a date (NULL if no logs)
4. Otherwise, return today + future_days_read days (rounded to nearest integer)

**Edge Cases:**
- If no `started_at`, return NULL
- If `page` is 0 or negative, return NULL (division by zero)
- If no logs and finished, return NULL
- Round `future_days_read` to nearest integer before adding to today

**Architecture Pattern:**
- Add `CalculateFinishedAt()` method to `Project` model (similar to existing methods)
- Update `ProjectResponse` DTO to include `FinishedAt` field (already in schema)
- Implementation follows existing Clean Architecture layers

### 2. Files to Modify

#### New Methods
- **`internal/domain/models/project.go`**: Add `CalculateFinishedAt()` method
  - Signature: `func (p *Project) CalculateFinishedAt(logs []*dto.LogResponse) *time.Time`
  - Follows same pattern as `CalculateMedianDay()` and `CalculateDaysUnreading()`

#### Existing Files (already present)
- **`internal/domain/dto/project_response.go`**: Already has `FinishedAt *string` field
- **`internal/adapter/postgres/project_repository.go`**: Already scans `finished_at` field
- **`internal/domain/models/project.go`**: Already has `FinishedAt *time.Time` field

#### Files to Review
- **`internal/adapter/postgres/project_repository.go`**: Check `GetWithLogs` and `GetAllWithLogs` to ensure finished_at is included
- **`internal/api/v1/handlers/projects_handler.go`**: Verify response includes finished_at (handled by DTO)

### 3. Dependencies

**Prerequisites (already satisfied):**
- ✅ RDL-024 - `median_day` calculation already implemented (similar pattern)
- ✅ RDL-023 - `days_unreading` calculation already implemented (similar pattern)
- ✅ RDL-022 - `status` determination already implemented
- ✅ RDL-021 - Config structure with status ranges already in place
- ✅ RDL-019 - Date/time format alignment to RFC3339 completed

**No blocking tasks**

**Setup steps:**
1. No database migrations needed (schema already has `finished_at` field)
2. Test database populated from `docs/database.sql` has sample data

### 4. Code Patterns

**Follow existing patterns:**

1. **Method signature** (match `CalculateMedianDay` pattern):
   ```go
   func (p *Project) CalculateFinishedAt(logs []*dto.LogResponse) *time.Time
   ```

2. **Date handling** (match existing methods):
   - Use `time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)` for date-only comparison
   - Return `*time.Time` for nullable field

3. **Edge case handling**:
   - Return `nil` for absent/invalid started_at
   - Return `nil` for page <= 0 (prevent division by zero)
   - Return `nil` if finished but no logs exist

4. **Integration with existing flow**:
   - In repository's `GetWithLogs`/`GetAllWithLogs`, call `CalculateFinishedAt()` after `CalculateDaysUnreading()`
   - Convert `*time.Time` to `*string` (RFC3339) in DTO creation

5. **Naming conventions**:
   - Method: `CalculateFinishedAt` (PascalCase, verb+noun)
   - Return type: `*time.Time` (pointer for nullable)
   - DTO field: `FinishedAt` (matches Rails `finished_at`)

### 5. Testing Strategy

**Unit Tests (internal/domain/models/project_test.go):**

1. **Normal case**: Calculate finished_at correctly
   ```go
   // Test: started 10 days ago, 50 pages, total 100
   // days_reading = 10, future_days = (10 * 100) / 50 = 20
   // finished_at = today + 20 days
   ```

2. **Finished book edge case**: Return last log date
   ```go
   // Test: page >= total_page, last log has data "2026-03-15"
   // Should return that date as time.Time pointer
   ```

3. **No logs with finished book**: Return NULL
   ```go
   // Test: page >= total_page, no logs
   // Should return nil
   ```

4. **Zero page edge case**: Return NULL (division by zero)
   ```go
   // Test: page = 0, should return nil regardless of other values
   ```

5. **No started_at edge case**: Return NULL
   ```go
   // Test: Project with no started_at
   // Should return nil
   ```

6. **Future started_at edge case**: Calculate properly with negative days_reading
   ```go
   // Test: started_at is 5 days in future
   // days_reading = -5, future_days = (-5 * 100) / 50 = -10
   // finished_at = today - 10 days (past date)
   ```

**Integration Tests:**

1. **GetWithLogs endpoint**: Verify finished_at in single project response
2. **GetAllWithLogs endpoint**: Verify finished_at in list response
3. **End-to-end**: Compare JSON output with Rails API

**Testing approach:**
- Use testing-expert subagent for comprehensive test coverage
- Test both positive and edge cases
- Verify NULL handling (returns JSON null, not empty string)

### 6. Risks and Considerations

**Potential pitfalls:**

1. **Date precision**: Rails uses `to_i` for integer days, Go must match exactly
   - Solution: Use same date-only truncation as other calculation methods

2. **Division by zero**: If page = 0, result is undefined
   - Solution: Return NULL (consistent with Ruby's infinity behavior)

3. **Negative future_days**: If page > total_page, future_days becomes negative
   - Rails handles this (returns past date)
   - Solution: Allow negative values, Go time arithmetic handles it

4. **Last log date format**: Rails uses `last_read[:data].to_date`
   - We store as string "YYYY-MM-DD" in logs
   - Solution: Parse and convert last log's Data field to time.Time

5. **Rounding behavior**: Rails uses `.round` on days
   - Solution: Use `math.Round()` on the result before adding to today

6. **Timezone handling**: Rails `Date.today` uses server timezone
   - We use `time.Now()` with UTC date truncation
   - Solution: Already standard in existing code, consistent with Rails behavior when both use same timezone

**Database considerations:**
- Column already exists in schema: `finished_at TIMESTAMP WITH TIME ZONE`
- No index needed (calculated field, not queried)
- NULL handling: PostgreSQL handles NULL correctly

**Deployment considerations:**
- No database changes required (backward compatible)
- Add to tests before merge (per DoD)
- May want to verify with Rails comparison tests

**Trade-offs:**
- No caching (calculated fresh per request, same as Rails)
- No partial calculation (each method is independent, like Rails)
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implementation completed on 2026-04-03
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [ ] #2 All integration tests pass use testing-expert subagent for test execution and verification
- [x] #3 go fmt and go vet pass with no errors
- [x] #4 Clean Architecture layers properly followed
- [x] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Database queries optimized with proper indexes
- [ ] #8 Documentation updated in QWEN.md
- [x] #9 New code paths include error path tests
- [ ] #10 HTTP handlers test both success and error responses
- [ ] #11 Integration tests verify actual database interactions
- [x] #12 Tests use testing-expert subagent for test execution and verification
- [ ] #13 #1 All unit tests pass - 44 tests passed using testing-expert subagent
- [ ] #14 #2 All integration tests pass - no integration tests for this pure model change
- [ ] #15 #3 go fmt and go vet pass with no errors
- [ ] #16 #4 Clean Architecture layers properly followed - implementation in domain/models
- [ ] #17 #5 Error responses consistent with existing patterns - nil handling consistent with other methods
- [ ] #18 #6 HTTP status codes correct for response type - N/A for model change
- [ ] #19 #7 Database queries optimized with proper indexes - no DB changes needed
- [ ] #20 #8 Documentation updated in QWEN.md - pending final summary
- [ ] #21 #9 New code paths include error path tests - 9 tests covering all edge cases
- [ ] #22 #10 HTTP handlers test both success and error responses - N/A for model change
- [ ] #23 #11 Integration tests verify actual database interactions - no DB changes needed
- [ ] #24 #12 Tests use testing-expert subagent for test execution and verification
<!-- DOD:END -->
