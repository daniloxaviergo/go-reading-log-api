---
id: RDL-100
title: >-
  [doc-009 Phase 3] Complete fixture data validation for Dashboard integration
  tests
status: To Do
assignee:
  - catarina
created_date: '2026-04-24 13:42'
updated_date: '2026-04-24 15:23'
labels:
  - feature
  - test-fix
  - p2-high
dependencies: []
references:
  - REQ-04
  - Decision 4
documentation:
  - doc-009
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement FixtureValidator with comprehensive checks ensuring 7 weekday coverage and minimum 30 days of data. Update all Dashboard integration test scenarios with complete fixture data and add validator to prevent cryptic test failures.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Dashboard integration tests have complete fixture data covering all 7 weekdays
- [ ] #2 FixtureValidator catches missing or insufficient data before test execution
- [ ] #3 All 3 Dashboard integration tests pass with valid fixtures
- [ ] #4 Chart contains all 30 days of data for mean progress calculation
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The task requires implementing a **FixtureValidator** with comprehensive validation checks for Dashboard integration test fixtures. The validator will ensure:

1. **7 weekday coverage**: All 7 days of the week must be represented in the fixture data
2. **Minimum 30 days of data**: At least 30 log entries to support mean progress calculations

The implementation follows a "fail-fast" pattern where validation occurs before test execution, providing clear error messages about missing or insufficient data.

**Architecture Decision:**
- Create a standalone `FixtureValidator` type with method-based validation
- Return detailed error lists rather than failing on first error (allows fixing multiple issues at once)
- Integrate with existing `Scenario` struct in the fixtures package
- Use Go's error wrapping for context preservation

**Why this approach:**
- Early validation prevents cryptic test failures deep in test execution
- Multiple errors reported simultaneously improve developer experience
- Minimal coupling with existing code - validator is a pure validation component

---

### 2. Files to Modify

| File | Action | Description |
|------|--------|-------------|
| `test/fixtures/dashboard/validator.go` | **CREATE** | New file containing `FixtureValidator` struct and validation methods |
| `test/fixtures/dashboard/scenarios.go` | **MODIFY** | Add validator integration to existing scenarios; ensure all scenarios pass validation |
| `test/integration/error_scenarios_test.go` | **MODIFY** | Add validation call before test execution in error scenarios |
| `test/integration/rails_comparison_test.go` | **MODIFY** | Add validation for comparison tests |
| `test/fixtures/dashboard/fixtures.go` | **MODIFY** | Add `Validate()` method to `DashboardFixtures` |

---

### 3. Dependencies

**Prerequisites:**
- ✅ Existing fixture infrastructure (`test/fixtures/dashboard/scenarios.go`)
- ✅ Dashboard repository and handlers (already implemented in Phase 2)
- ✅ TestHelper with database setup/teardown (RDL-102 - concurrent drop)

**No blocking dependencies** - this task can proceed independently but should be completed before full test suite runs.

---

### 4. Code Patterns

**Validation Pattern:**
```go
type FixtureValidator struct {
    logs []*LogFixture
}

func (v *FixtureValidator) Validate() []error {
    var errors []error
    
    // Collect all validation errors
    if err := v.validateWeekdayCoverage(); err != nil {
        errors = append(errors, err)
    }
    if err := v.validateDataRange(); err != nil {
        errors = append(errors, err)
    }
    
    return errors
}
```

**Integration Pattern:**
- Validator is instantiated per-scenario
- Validation runs in test setup phase
- Test fails immediately with detailed error output if validation fails

**Error Message Format:**
```go
return fmt.Errorf("fixture validation failed: %s (got: %d, expected: %d)", 
    "insufficient data", len(logs), 30)
```

---

### 5. Testing Strategy

**Unit Tests for Validator:**
- Test each validation method independently
- Test with edge cases (exactly 7 days, exactly 30 logs)
- Test with missing data scenarios

**Integration Test Updates:**
- All existing Dashboard integration tests must use validated fixtures
- Add validation check at start of `RunErrorScenarios` and `RunComparisonTests`
- Verify that all 8 dashboard endpoints have valid fixture data

**Test Coverage Requirements:**
- `TestFixtureValidator_WeekdayCoverage`: Validates 7-day requirement
- `TestFixtureValidator_DataRange`: Validates 30-day minimum
- `TestFixtureValidator_Combined`: Validates multiple failures reported together

---

### 6. Risks and Considerations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Existing tests break due to validation failures | Medium | Run validator on all existing scenarios first; fix issues incrementally |
| Validation adds significant overhead | Low | Validator is lightweight (<1ms execution) |
| Some tests intentionally use minimal data | Medium | Create "minimal valid" scenario for edge case tests |

**Key Design Decisions:**
1. **Validation runs before test execution** - Prevents partial test runs with invalid data
2. **All errors collected, not just first** - Developer can fix multiple issues in one iteration
3. **Validator is opt-in via helper methods** - Doesn't force validation on all tests immediately

---

### Implementation Steps

1. **Create `validator.go`** with `FixtureValidator` struct and methods:
   - `ValidateWeekdayCoverage()` - Check for 7 unique weekdays
   - `ValidateDataRange()` - Check for minimum 30 log entries
   - `Validate()` - Run all validations and return combined errors

2. **Update existing scenarios** in `scenarios.go`:
   - Review each scenario for compliance
   - Add missing data to meet requirements
   - Document validation requirements in scenario comments

3. **Integrate validator into test infrastructure**:
   - Add validation call in `RunErrorScenarios`
   - Add validation call in `RunComparisonTests`
   - Ensure clear error messages on failure

4. **Run full test suite** to identify any remaining issues
5. **Update documentation** in AGENTS.md with new validation requirements
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
