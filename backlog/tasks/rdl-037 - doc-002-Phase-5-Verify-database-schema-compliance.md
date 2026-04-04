---
id: RDL-037
title: '[doc-002 Phase 5] Verify database schema compliance'
status: To Do
assignee:
  - catarina
created_date: '2026-04-03 14:05'
updated_date: '2026-04-04 07:25'
labels:
  - phase-5
  - database-verification
  - schema
dependencies: []
references:
  - 'PRD Section: Traceability Matrix'
  - 'PRD Section: Acceptance Criteria - NF1'
documentation:
  - doc-002
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Run database schema verification to ensure all constraints are properly defined and indexes exist. Verify database-level constraints for page ≤ total_page and start_page ≤ end_page match validation logic. Ensure schema matches implementation expectations.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Database constraints match validation rules
- [ ] #2 Index exists for logs JOIN optimization
- [ ] #3 Schema documented and verified
- [ ] #4 No schema drift from PRD requirements
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

**Status**: This task is a verification/review task - it does NOT require code changes to the database or validation logic. Instead, it requires creating automated tests that verify the existing code matches the PRD requirements.

**Objective**: Verify that database constraints (page ≤ total_page, start_page ≤ end_page) and existing validation logic are correctly implemented and tested.

**Approach**:

1. **Review Database Schema**: Examine the `docs/database.sql` file to identify existing database-level constraints on:
   - `projects.page` vs `projects.total_page`
   - `logs.start_page` vs `logs.end_page`

2. **Review Validation Package**: Verify the `internal/validation/` package contains:
   - `ValidatePage(page int, totalPage int)` - ensures page ≤ total_page
   - `ValidateStartEndPage(startPage int, endPage int)` - ensures start_page ≤ end_page
   - Comprehensive test coverage

3. **Review Repository Layer**: Verify that the PostgreSQL repository:
   - Uses database-level constraints where available
   - Applies validation before database operations
   - Returns appropriate errors for constraint violations

4. **Write Verification Tests**: Create tests that:
   - Verify database constraints match validation logic
   - Test validation behavior matches Rails API behavior
   - Document any schema/implementation gaps

**Why this approach**: The PRD requirements have already been implemented (RDL-030, RDL-031, RDL-032). This task verifies the implementation and fills any test gaps.

### 2. Files to Modify

**No existing files need modification** - all validation logic is already implemented.

**New files to create** (for verification/testing):

| File | Purpose |
|------|---------|
| `test/database_schema_test.go` | Verify database constraints match validation rules |
| `test/validation_test.go` | Integration tests for validation constraints |
| `docs/database_constraints.md` | Document database-level constraints |

**Files to review** (read-only verification):

| File | Review Purpose |
|------|----------------|
| `docs/database.sql` | Identify existing constraints |
| `internal/validation/validate_project.go` | Verify `ValidatePage` implementation |
| `internal/validation/validate_log.go` | Verify `ValidateStartEndPage` implementation |
| `internal/validation/validate_test.go` | Verify test coverage |
| `internal/adapter/postgres/project_repository.go` | Verify validation integration |
| `internal/domain/models/project.go` | Verify calculation methods match validation |

### 3. Dependencies

**Prerequisites**:

1. **RDL-030 Complete**: Shared validation package must exist (already done)
2. **RDL-031 Complete**: `page ≤ total_page` validation must be implemented (already done)
3. **RDL-032 Complete**: `start_page ≤ end_page` validation must be implemented (already done)

**No external dependencies** required - all validation logic is self-contained.

**Blocking Issues**: None - this is a verification task for existing functionality.

### 4. Code Patterns

**Validation Patterns to Verify**:

1. **Error Structure**: Must use `ValidationError` with:
   - `Code`: Machine-readable error (e.g., `"page_exceeds_total"`)
   - `Field`: Field name (e.g., `"page"`)
   - `Message`: Descriptive error (includes actual values)

2. **Validation Functions**: Single responsibility, return `*ValidationError` or `nil`

3. **Error Messages**: Include actual values for debugging:
   - ✅ Correct: `"page (%d) cannot exceed total_page (%d)"`
   - ❌ Avoid: Generic messages without values

**Example Verification Tests**:

```go
// Test that validates page ≤ total_page constraint
func TestValidation_PageExceedsTotal_DatabaseMatch(t *testing.T) {
    // Verify validation error matches database constraint error
    err := validation.ValidatePage(150, 100)
    
    if err == nil {
        t.Fatal("Expected validation error for page > total_page")
    }
    
    if err.Code != "page_exceeds_total" {
        t.Errorf("Expected error code 'page_exceeds_total', got '%s'", err.Code)
    }
    
    // Verify error message contains the values
    if !strings.Contains(err.Message, "150") || !strings.Contains(err.Message, "100") {
        t.Errorf("Error message should include actual values: %s", err.Message)
    }
}
```

### 5. Testing Strategy

**Verification Test Categories**:

**Category 1: Database Constraint Verification**
- Test: Verify `docs/database.sql` contains no foreign key constraints for page validation
- Test: Verify no CHECK constraints exist (validation is application-level)
- Test: Verify `projectspage` column type allows negative values (validation prevents invalid values)

**Category 2: Validation Logic Verification**
- Test: `ValidatePage(100, 100)` returns `nil` (equal values allowed)
- Test: `ValidatePage(101, 100)` returns error (exceeds total)
- Test: `ValidatePage(-1, 100)` returns error (negative values)
- Test: `ValidateStartEndPage(10, 20)` returns `nil` (start < end)
- Test: `ValidateStartEndPage(20, 10)` returns error (start > end)
- Test: `ValidateStartEndPage(10, 10)` returns `nil` (equal values)

**Category 3: Integration Verification**
- Test: Repository uses validation before database operations
- Test: Repository returns validation errors to callers
- Test: HTTP handlers properly serialize validation errors

**Category 4: Rails Compatibility Verification**
- Test: Validation logic matches Rails behavior (zero/negative handling)
- Test: Error messages match Rails format (where applicable)

**Test Coverage Requirements**:
- All validation functions: 100% coverage
- Integration tests: All happy and error paths
- Verification tests: Critical paths documented

### 6. Risks and Considerations

**Risks**:

1. **Schema/Validation Gap**: If `docs/database.sql` shows no constraints but PRD expects them, this indicates a gap.

2. **Test Coverage Gap**: Verification tests may reveal missing test cases.

3. **Rails Behavior Differences**: Rails validation may handle edge cases differently.

**Considerations**:

1. **Database Constraints**: PostgreSQL allows negative values at database level; application validation prevents them. This is intentional - database is permissive, application enforces rules.

2. **Validation Location**: All validation is application-level (not database-level). This matches PRD requirements.

3. **Documentation**: Task must produce `docs/database_constraints.md` for schema documentation.

4. **Test Results**: All verification tests must pass with 100% coverage to meet DoD.

**Gap Analysis Template** (to document findings):

| Requirement | PRD Rule | Implementation | Status |
|-------------|----------|----------------|--------|
| Page validation | page ≤ total_page | `ValidatePage()` exists | ✅ Verified |
| Log page validation | start_page ≤ end_page | `ValidateStartEndPage()` exists | ✅ Verified |
| Database constraints | None documented | No constraints found | ⚠️ Documented |

### 7. Verification Checklist

Use this checklist during verification:

**Database Schema**:
- [ ] Review `docs/database.sql` for projects table structure
- [ ] Review `docs/database.sql` for logs table structure
- [ ] Document any existing CHECK constraints
- [ ] Document any FOREIGN KEY constraints related to page validation

**Validation Logic**:
- [ ] Verify `ValidatePage()` matches PRD requirement
- [ ] Verify `ValidateStartEndPage()` matches PRD requirement
- [ ] Verify error codes match expected format
- [ ] Verify error messages include values

**Test Coverage**:
- [ ] Verify all validation edge cases are tested
- [ ] Verify negative values are handled correctly
- [ ] Verify boundary values (equal values) work correctly

**Integration**:
- [ ] Verify repository uses validation
- [ ] Verify handlers return validation errors correctly
- [ ] Verify error responses match existing patterns

**Documentation**:
- [ ] `docs/database_constraints.md` created
- [ ] Schema limitations documented
- [ ] Validation methodology explained

### 8. Acceptance Criteria Mapping

The task has 4 acceptance criteria - each maps to verification activities:

| AC | Verification Activity | Evidence |
|----|----------------------|----------|
| #1 Database constraints match validation rules | Review `docs/database.sql`, compare to `ValidatePage()` and `ValidateStartEndPage()` | `docs/database_constraints.md` |
| #2 Index exists for logs JOIN optimization | Verify indexes in `docs/database.sql`, check repository uses efficient queries | PR #28/29 documentation |
| #3 Schema documented and verified | Create `docs/database_constraints.md` with findings | Documentation file |
| #4 No schema drift from PRD requirements | Compare database structure to PRD | Gap analysis in documentation |

### 9. Deliverables

**Products to produce**:

1. **Verification Tests** (`test/database_schema_test.go`, `test/validation_test.go`)
   - Database constraint verification tests
   - Validation logic integration tests
   - Rails compatibility tests

2. **Documentation** (`docs/database_constraints.md`)
   - Database structure summary
   - Constraint analysis
   - Gap documentation
   - Validation methodology

3. **Test Results Report**
   - All verification tests pass
   - 100% coverage for validation package
   - No blocking issues

4. **Gap Analysis**
   - Document any mismatches between PRD and implementation
   -Recommendations for fixes (if any gaps found)

### 10. Recommended Workflow

**Phase 1: Review & Analysis** (2-4 hours)
- Review `docs/database.sql`
- Review `internal/validation/` package
- Compare to PRD requirements
- Identify gaps

**Phase 2: Verification Tests** (4-6 hours)
- Write database constraint verification tests
- Write validation logic integration tests
- Write edge case tests
- Execute tests

**Phase 3: Documentation** (2-3 hours)
- Write `docs/database_constraints.md`
- Document findings
- Document any gaps

**Phase 4: Final Verification** (1-2 hours)
- Run all tests
- Verify coverage
- Update backlog task with findings
- Request review

**Total Estimate**: 9-15 hours (1-2 days)

### 11. Success Criteria

Task is complete when:

1. ✅ All verification tests pass with 100% coverage
2. ✅ `docs/database_constraints.md` exists and documents:
   - Database structure
   - Constraint analysis
   - Validation methodology
   - Gap documentation
3. ✅ No critical gaps found between PRD and implementation
4. ✅ All DoD criteria met (test counts, coverage, formatting)

### 12. Risks and Considerations Summary

**Low Risk Task**: This is a verification task - no code changes required unless gaps are found.

**Gap Handling**: If gaps are found:
- Document gaps in `docs/database_constraints.md`
- Create follow-up tasks for fixes
- Update this task with gap analysis

**Testing Focus**: Verification tests are more important than code changes - ensure tests comprehensively cover validation scenarios.
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
