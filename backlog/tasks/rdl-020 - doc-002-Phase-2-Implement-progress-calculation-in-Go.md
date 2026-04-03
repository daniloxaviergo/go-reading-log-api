---
id: RDL-020
title: '[doc-002 Phase 2] Implement progress calculation in Go'
status: To Do
assignee:
  - thomas
created_date: '2026-04-03 14:02'
updated_date: '2026-04-03 16:21'
labels:
  - phase-2
  - derived-calculation
  - go-implementation
dependencies: []
references:
  - >-
    PRD Section: Technical Decisions - Decision 1: Derived Calculations
    Implementation
  - 'PRD Section: Validation Rules - progress range'
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement the progress calculation method in Go matching Rails behavior: progress = (page / total_page) * 100 rounded to 2 decimal places. Clamp result to 0.00-100.00 range and handle edge cases (zero total_page, null values).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 progress = (page / total_page) * 100 rounded to 2 decimal places
- [x] #2 Result clamped to 0.00-100.00 range
- [x] #3 Zero total_page edge case returns 0.00
- [x] #4 Calculate method added to Project model or calculations package
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Implement the progress calculation in Go by adding a method to the Project model that calculates progress as `(page / total_page) * 100` rounded to 2 decimal places, clamped to 0.00-100.00 range.

**Approach**:
- Add a `CalculateProgress()` method to the `Project` model in `internal/domain/models/project.go`
- The method handles edge cases: zero `total_page` (returns 0.00), null/zero `page` (returns 0.00)
- Clamp result to 0.00-100.00 range as per Rails validation rules
- Use `strconv.FormatFloat` with 'f' format and 2 decimal places for consistent rounding
- Add unit tests for all edge cases

**Why this approach**:
- Keeps business logic in the domain model layer (Clean Architecture)
- Matches Rails implementation: `((page.to_f / total_page.to_f) * 100).round(2)`
- Returns pointer to float64 to match existing JSON serialization pattern
- Handles all edge cases explicitly per acceptance criteria

### 2. Files to Modify

| File | Action | Reason |
|------|--------|--------|
| `internal/domain/models/project.go` | Modify | Add `CalculateProgress()` method and update `NewProject()` to accept optional progress parameter |
| `internal/domain/models/project_test.go` | Modify | Add unit tests for progress calculation (normal case, edge cases) |
| `internal/domain/dto/project_response.go` | Modify | Add helper method to set calculated progress |
| `internal/adapter/postgres/project_repository.go` | Modify | Update repository to call `CalculateProgress()` when building responses |
| `internal/api/v1/handlers/projects_handler.go` | Modify | Ensure calculated progress is set in responses |
| `test/test_helper.go` | Modify | Update mock to support progress calculation testing |

### 3. Dependencies

**Prerequisites**:
- Task RDL-020 is independent but depends on existing domain model structure (already in place)
- No other tasks need to be completed first
- Rails implementation reference: `internal/domain/models/project.rb` (already analyzed)

**Edge cases to handle**:
- When `total_page` is 0 → return 0.00
- When `page` is 0 or null → return 0.00
- When `total_page` is negative → return 0.00 (invalid data)
- When calculated result > 100.00 → clamp to 100.00
- When calculated result < 0.00 → clamp to 0.00

### 4. Code Patterns

**Existing conventions to follow**:

1. **Project model pattern** (from `internal/domain/models/project.go`):
```go
type Project struct {
    ctx        context.Context
    ID         int64      `json:"id"`
    Name       string     `json:"name"`
    TotalPage  int        `json:"total_page"`
    StartedAt  *time.Time `json:"started_at"`
    Page       int        `json:"page"`
    Progress   *float64   `json:"progress,omitempty"`
    // ... other fields
}
```

2. **Helper functions** (from `test/test_helper.go`):
```go
func floatPtr(f float64) *float64 {
    return &f
}
```

3. **JSON serialization**:
- Use pointer to float64 for optional fields
- `omitempty` tag for optional derived fields

4. **Error handling**:
- Return 0.00 for invalid/edge cases (nil pointer not needed for default value)

### 5. Testing Strategy

**Unit tests to add** (in `internal/domain/models/project_test.go`):
1. `TestProject_CalculateProgress_Normal` - Standard calculation (50/100 = 50.00)
2. `TestProject_CalculateProgress_Rounding` - Test 2 decimal places (33/100 = 33.00, 1/3 = 33.33)
3. `TestProject_CalculateProgress_ZeroTotalPage` - Returns 0.00
4. `TestProject_CalculateProgress_ZeroPage` - Returns 0.00
5. `TestProject_CalculateProgress_ClampMax` - Values over 100% clamped to 100.00
6. `TestProject_CalculateProgress_EmptyContext` - Test with nil context

**Test approach**:
- Use `testing.T` for assertions
- Assert with tolerance for float comparison (use `math.Abs` for epsilon comparison)
- Test all acceptance criteria individually

**Integration test notes**:
- Existing tests in `test_helper.go` use `MockProjectRepository`
- No database changes needed - calculation is pure logic

### 6. Risks and Considerations

**Potential pitfalls**:
1. **Integer division**: Must convert to float before division to avoid integer truncation
2. **Division by zero**: Must check `total_page` before division
3. **Float comparison**: Unit tests should use epsilon comparison for float equality
4. **Rounding behavior**: Rails uses `round(2)` which rounds half to even (banker's rounding). Go's `strconv.FormatFloat` with 'f' uses round-half-up. May need to use `math.Round` for exact alignment.

**Decision points**:
- **Return type**: `*float64` matches existing pattern in Project model
- **Calculation timing**: Calculate when setting fields (deferred calculation) or on-demand method
  - Chosen: On-demand method for cleaner separation of concerns
- **Clamping**: Apply clamping after rounding per Rails validation rules
- **Zero total_page**: Return 0.00 instead of error/panic per acceptance criteria

**Implementation notes**:
- Rails rounds to 2 decimal places using `.round(2)` which rounds half to even
- Go's `strconv.FormatFloat` with 'f' format and 2 decimal places rounds half-up
- For exact Rails compatibility, use: `math.Round((page/total_page)*100*100) / 100`
- Clamping must happen after rounding to handle edge cases like 100.005 → 100.01 → clamp to 100.00
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
2026-04-03: Implementation complete. Added CalculateProgress() method to Project model with edge case handling. Tests verified via testing-expert subagent.
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [ ] #2 All integration tests pass use testing-expert subagent for test execution and verification
- [x] #3 go fmt and go vet pass with no errors
- [x] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Database queries optimized with proper indexes
- [ ] #8 Documentation updated in QWEN.md
- [ ] #9 New code paths include error path tests
- [ ] #10 HTTP handlers test both success and error responses
- [ ] #11 Integration tests verify actual database interactions
- [ ] #12 Tests use testing-expert subagent for test execution and verification
<!-- DOD:END -->
