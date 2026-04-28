---
id: RDL-123
title: '[doc-10 Phase 6] Create Rails calculation reference documentation'
status: To Do
assignee:
  - workflow
created_date: '2026-04-28 00:31'
updated_date: '2026-04-28 06:23'
labels:
  - documentation
  - phase-6
  - backend
dependencies: []
documentation:
  - doc-010
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create docs/rails-calculation-reference.md documenting Rails V1::MeanLog, V1::MaxLog algorithms with code examples and formula explanations for developer reference.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 rails-calculation-reference.md created
- [ ] #2 V1::MeanLog algorithm documented
- [ ] #3 V1::MaxLog algorithm documented
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task involves reviewing and updating the existing `docs/rails-calculation-reference.md` documentation to ensure it accurately reflects the actual Go implementation behavior rather than just the Rails reference.

**Key Findings from Codebase Research:**

1. **Current Documentation Issues:**
   - Documentation states edge cases return `0.0` (e.g., empty logs, zero intervals)
   - Actual Go implementation returns `nil` for these cases (using pointer types)
   - This is a critical discrepancy that could confuse developers

2. **Algorithms to Document:**
   - **V1::MeanLog** (`GetMeanByWeekday`): Calculates mean pages per 7-day interval for a specific weekday
     - Formula: `total_pages / count_reads` where `count_reads = floor((log_data - begin_data) / 7 days)`
     - Returns `nil` when no logs exist or count_reads = 0
   
   - **V1::MaxLog** (`GetMaxByWeekday`): Calculates maximum pages read in a single day for a weekday
     - Formula: `MAX(end_page - start_page)` for logs matching weekday
     - Returns `nil` when no logs exist

3. **Implementation Strategy:**
   - Review existing documentation line by line
   - Update edge case behavior to match actual Go implementation (nil vs 0.0)
   - Update Go code examples to match actual implementation in `internal/adapter/postgres/dashboard_repository.go`
   - Ensure null handling documentation is accurate
   - Add missing implementation details (e.g., `GetProjectWeekdayMean` for per-project calculations)

### 2. Files to Modify

| File | Action | Reason |
|------|--------|--------|
| `docs/rails-calculation-reference.md` | **Update** | Fix discrepancies between documentation and actual implementation |

**Specific Changes Required:**

1. **MeanLog Algorithm Section:**
   - Update edge case behavior: change "return 0.0" to "return nil"
   - Update Go implementation example to match actual `GetMeanByWeekday` method
   - Fix SQL query to use `COALESCE` and proper NULL handling
   - Update example calculation to show NULL return for zero intervals

2. **MaxLog Algorithm Section:**
   - Update edge case behavior: change "return 0.0" to "return nil"
   - Update Go implementation to use pointer return type `(*float64, error)`
   - Fix SQL query to return NULL when no logs exist

3. **Null Handling Rules Section:**
   - Update table to clarify which fields return `nil` vs `0.0`
   - Add `mean_day` and `max_day` to "When to Return Null" table
   - Update JSON serialization examples to show `null` values

4. **Add New Section - Per-Project vs Global Calculations:**
   - Document `GetProjectWeekdayMean` (per-project weekday mean)
   - Document `GetMeanByWeekday` (global weekday mean - V1::MeanLog)
   - Explain when each is used in the DayService

### 3. Dependencies

- **None** - This is a documentation-only task
- Prerequisite: Task RDL-114 (Implement MaxDay field) and RDL-116 (Align mean_day calculation) should be completed to ensure implementation is stable

### 4. Code Patterns

**Documentation Format to Follow:**

```markdown
## Algorithm Name

### Purpose

Brief description of what the algorithm calculates.

### Rails Implementation Reference

```ruby
# Original Rails code for reference
```

### Algorithm Steps

1. Step 1
2. Step 2
...

### Go Implementation

```go
// Comment explaining the method
func MethodName(ctx context.Context, ...) (*float64, error) {
    // Actual implementation
}
```

### Example Calculation

**Input:**
- Parameter 1: value
- Parameter 2: value

**Calculation:**
```
step-by-step calculation
```

**Output:** `value` (or `null` for edge cases)
```

**Naming Conventions:**
- Use snake_case for field names in JSON examples
- Use PascalCase for Go method names
- Use code blocks with proper syntax highlighting

### 5. Testing Strategy

**Documentation Verification:**

1. **Code Accuracy Test:**
   - Verify every code example in documentation matches actual implementation
   - Run `grep` to find actual method implementations
   - Compare SQL queries character by character

2. **Behavior Verification:**
   - Test edge cases with actual database
   - Verify NULL vs 0.0 behavior matches documentation
   - Check JSON serialization produces correct `null` values

3. **Integration Test Reference:**
   - Review existing integration tests in `test/integration/dashboard_day_permean_integration_test.go`
   - Ensure documentation examples match test expectations
   - Verify edge case assertions (e.g., `assert.Nil(t, statsMap["per_mean_day"])`)

**Test Cases to Validate:**

| Scenario | Expected Result | Test Method |
|----------|----------------|-------------|
| Empty logs for weekday | `nil` | `GetMeanByWeekday` with no matching logs |
| Logs within same week | `nil` (count_reads = 0) | `GetMeanByWeekday` with < 7 day span |
| Valid logs spanning multiple weeks | Calculated mean | `GetMeanByWeekday` with normal data |
| No logs for max calculation | `nil` | `GetMaxByWeekday` with no matching logs |
| Valid logs for max calculation | Maximum value | `GetMaxByWeekday` with data |

### 6. Risks and Considerations

**Critical Issues:**

1. **Documentation Accuracy:**
   - **Risk:** Current documentation incorrectly states edge cases return `0.0`
   - **Impact:** Developers implementing new features may introduce bugs
   - **Mitigation:** Thoroughly verify all edge case behavior against actual implementation

2. **Rails vs Go Behavior Differences:**
   - **Risk:** Rails may return `0.0` while Go returns `nil`
   - **Impact:** API response inconsistency between Rails and Go
   - **Mitigation:** Document the Go-specific behavior clearly; note any Rails differences

3. **Null Handling Confusion:**
   - **Risk:** Developers may not understand when to use pointer types
   - **Impact:** JSON serialization issues (0 vs null)
   - **Mitigation:** Add clear examples of pointer usage and JSON output

**Trade-offs:**

- **Documentation Completeness vs Maintenance:** More detailed documentation is harder to keep up-to-date
  - Decision: Focus on algorithm accuracy; keep implementation examples minimal but correct

- **Rails Reference vs Go Implementation:** Should documentation emphasize Rails or Go?
  - Decision: Primary focus on Go implementation; Rails code as reference only

**Deployment Considerations:**

- This is a documentation-only change; no deployment risk
- Consider adding a "Last Verified Against Implementation" date to the document
- Suggest adding documentation review to the Definition of Done for future calculation-related tasks

**Acceptance Criteria Verification:**

- [x] `rails-calculation-reference.md` exists
- [ ] V1::MeanLog algorithm documented **accurately** (update required)
- [ ] V1::MaxLog algorithm documented **accurately** (update required)
- [ ] Edge case behavior matches actual implementation
- [ ] Go code examples are copy-pasteable and correct
- [ ] Null handling is clearly explained with examples

---

*Plan created: 2026-04-28*
*Based on codebase state as of: 2026-04-28*
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
