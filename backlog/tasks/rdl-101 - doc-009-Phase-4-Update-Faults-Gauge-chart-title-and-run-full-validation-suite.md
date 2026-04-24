---
id: RDL-101
title: >-
  [doc-009 Phase 4] Update Faults Gauge chart title and run full validation
  suite
status: To Do
assignee:
  - catarina
created_date: '2026-04-24 13:42'
updated_date: '2026-04-24 17:20'
labels:
  - documentation
  - test-fix
  - p3-medium
dependencies: []
references:
  - REQ-05
  - Decision 5
documentation:
  - doc-009
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update gauge chart title from 'Faults Gauge' to 'Fault Percentage by Weekday' for better user clarity. Run complete test suite to verify all acceptance criteria are met, document test patterns, and update AGENTS.md with new testing guidelines.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Gauge chart title is more descriptive and user-friendly
- [ ] #2 All 14 failing tests are now passing (100% pass rate)
- [ ] #3 Test execution time is under 30 seconds total
- [ ] #4 Code coverage meets minimum 80% threshold for modified files
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task addresses **two distinct requirements** from doc-009 Phase 4:

**Requirement 1: Gauge Chart Title Update**
- Change the gauge chart title from "Faults Gauge" to "Fault Percentage by Weekday"
- This is a user-facing label improvement for better clarity
- Requires updating the service and corresponding test expectations

**Requirement 2: Full Validation Suite Execution**
- Run complete test suite to verify all acceptance criteria are met
- Document test patterns for future reference
- Update AGENTS.md with new testing guidelines

**Architecture Decision:**
- Minimal code changes - single string literal update in `faults_service.go`
- Test updates - single assertion expectation change
- Documentation updates - comprehensive testing guidelines in AGENTS.md

**Why this approach:**
- The title change is purely cosmetic/user-facing
- No business logic changes required
- Low risk, high value for user experience
- Test suite validation ensures no regressions

---

### 2. Files to Modify

| File | Action | Description |
|------|--------|-------------|
| `internal/service/dashboard/faults_service.go` | **MODIFY** | Update gauge chart title from "Faults Gauge" to "Fault Percentage by Weekday" (line ~107) |
| `test/unit/faults_service_test.go` | **MODIFY** | Update test expectation for gauge chart title (line ~248) |
| `AGENTS.md` | **MODIFY** | Add comprehensive testing guidelines section documenting test patterns, best practices, and troubleshooting |

---

### 3. Dependencies

**Prerequisites:**
- ✅ RDL-100 complete - Fixture validation infrastructure in place
- ✅ Date abstraction layer implemented (RDL-099) - `dto.GetTodayFunc` available
- ✅ TestHelper with proper context timeout handling (RDL-098) - Fixed in PRD doc-009

**No blocking dependencies** - this task can proceed independently once RDL-100 is complete.

---

### 4. Code Patterns

**Title Update Pattern:**
```go
// In faults_service.go - Line ~107
return dto.NewEchartConfig().
    SetTitle("Fault Percentage by Weekday").  // Changed from "Faults Gauge"
    SetTooltip(map[string]interface{}{
        "formatter": "{a} <br/>{b} : {c}%",
    }).
    AddSeries(...)
```

**Test Update Pattern:**
```go
// In faults_service_test.go - Line ~248
assert.Equal(t, "Fault Percentage by Weekday", gauge.Title)  // Updated assertion
```

**Testing Guidelines Pattern (AGENTS.md):**
- Document test organization (unit vs integration)
- Document fixture patterns and validation
- Document date/time testing strategies
- Document database cleanup procedures

---

### 5. Testing Strategy

**Unit Tests to Verify:**

1. **Gauge Chart Title Test**
   ```bash
   go test -v ./test/unit/faults_service_test.go -run TestFaultsService_CreateGaugeChart
   ```
   - Verifies title is "Fault Percentage by Weekday"
   - Verifies chart configuration is valid
   - Verifies series data is correct

2. **Full Unit Test Suite**
   ```bash
   go test -v ./test/unit/...
   ```
   - All 14 faults service tests should pass
   - All other unit tests should pass
   - Expected: 100% pass rate

3. **Integration Tests**
   ```bash
   go test -v ./test/integration/...
   ```
   - Dashboard integration tests with validated fixtures
   - Projects integration tests
   - Expected: 100% pass rate (after RDL-100 fixes)

4. **Full Test Suite**
   ```bash
   go test -v ./...
   ```
   - All packages tested
   - Total execution time < 30 seconds
   - Code coverage > 80%

**Test Coverage Requirements:**
- `faults_service_test.go`: 14 tests, all passing
- Dashboard integration tests: 3 tests, all passing
- Projects integration tests: 3 tests, all passing
- Other unit tests: all passing

---

### 6. Risks and Considerations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Test expectation mismatch | Low | Update test assertion to match new title; verify with `go test -v` |
| Missing documentation | Medium | Comprehensive AGENTS.md update covering all testing aspects |
| Pre-existing test failures | High | RDL-100 should have addressed these; verify with full test run |
| Test execution time exceeds 30s | Low | Current infrastructure is optimized; monitor during run |

**Key Design Decisions:**

1. **Title Change Justification:**
   - "Faults Gauge" was non-descriptive
   - "Fault Percentage by Weekday" clearly communicates:
     - What is measured (fault percentage)
     - How it's grouped (by weekday)
     - User can immediately understand the visualization

2. **Minimal Changes Approach:**
   - Single string literal update = low risk
   - No logic changes = no regression risk
   - Easy to rollback if issues arise

3. **Documentation Comprehensive Coverage:**
   - Test organization and naming conventions
   - Fixture patterns and validation
   - Date/time testing strategies
   - Database cleanup procedures
   - Common troubleshooting steps

---

### Implementation Steps

1. **Update Faults Service Title** (`internal/service/dashboard/faults_service.go`)
   - Locate line ~107 in `CreateGaugeChart` method
   - Change `SetTitle("Faults Gauge")` to `SetTitle("Fault Percentage by Weekday")`
   - Verify no other references to old title exist

2. **Update Test Assertion** (`test/unit/faults_service_test.go`)
   - Locate test `TestFaultsService_CreateGaugeChart`
   - Update expected title assertion from "Faults Gauge" to "Fault Percentage by Weekday"
   - Run test to verify pass: `go test -v ./test/unit/... -run TestFaultsService_CreateGaugeChart`

3. **Run Full Unit Test Suite**
   ```bash
   go test -v ./test/unit/...
   ```
   - Verify all 14 faults service tests pass
   - Verify all other unit tests pass
   - Record any failures for investigation

4. **Run Integration Tests**
   ```bash
   go test -v ./test/integration/...
   ```
   - Verify dashboard integration tests pass (with validated fixtures from RDL-100)
   - Verify projects integration tests pass
   - Document any pre-existing failures

5. **Run Full Test Suite with Coverage**
   ```bash
   go test -v -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
   ```
   - Verify total execution time < 30 seconds
   - Verify code coverage > 80% for modified files
   - Generate coverage report for documentation

6. **Update AGENTS.md**
   - Add new section: "Testing Guidelines"
   - Document test patterns (unit, integration, fixture)
   - Document date/time testing strategies
   - Document database cleanup procedures
   - Add troubleshooting section

7. **Final Validation**
   ```bash
   go fmt ./...
   go vet ./...
   ```
   - Verify no formatting issues
   - Verify no static analysis errors

---

### Acceptance Criteria Verification

| AC | Description | Verification Method |
|----|-------------|---------------------|
| AC-01 | Gauge chart title is more descriptive and user-friendly | Code review + manual verification of title text |
| AC-02 | All 14 failing tests are now passing (100% pass rate) | `go test -v ./test/unit/...` - verify all pass |
| AC-03 | Test execution time is under 30 seconds total | `time go test ./...` - verify < 30s |
| AC-04 | Code coverage meets minimum 80% threshold for modified files | `go test -coverprofile=coverage.out` - verify > 80% |

### Definition of Done Verification

| DoD Item | Description | Verification |
|----------|-------------|--------------|
| #1 All unit tests pass | `go test ./test/unit/...` succeeds | Run command, verify no failures |
| #2 All integration tests pass | `go test ./test/integration/...` succeeds | Run command, verify no failures |
| #3 go fmt and go vet pass | No output/errors from commands | Run both commands, verify clean |
| #4 Clean Architecture layers properly followed | Code follows existing patterns | Code review |
| #5 Error responses consistent | Follow existing error patterns | Code review |
| #6 HTTP status codes correct | Match endpoint specifications | Code review |
| #7 Documentation updated in AGENTS.md | Testing guidelines section added | Review AGENTS.md |
| #8 New code paths include error path tests | Existing tests cover all paths | Review test coverage |
| #9 HTTP handlers test success/error | Existing tests comprehensive | Review test coverage |
| #10 Integration tests verify DB interactions | Existing tests validated | Review test implementation |

---

### Expected Results

**Before Changes:**
- Title: "Faults Gauge" (non-descriptive)
- Test expectation: "Faults Gauge"
- AGENTS.md: No testing guidelines section

**After Changes:**
- Title: "Fault Percentage by Weekday" (clear, descriptive)
- Test expectation: "Fault Percentage by Weekday"
- AGENTS.md: Comprehensive testing guidelines added
- All tests: 100% pass rate
- Execution time: < 30 seconds
- Code coverage: > 80%

---

### Rollback Plan

If issues arise:
1. Revert title change in `faults_service.go`
2. Revert test assertion in `faults_service_test.go`
3. No database or complex state changes to revert
4. Minimal code footprint = easy rollback

---

### Notes for Implementation

- **Low Risk Task**: This is primarily a documentation and validation task with minimal code changes
- **RDL-100 Prerequisite**: RDL-100 must be complete for integration tests to pass
- **Test Helper Fixes**: Ensure RDL-098 context timeout fixes are in place
- **Date Abstraction**: Ensure RDL-099 date layer is working correctly
- **Documentation Quality**: AGENTS.md update should be comprehensive and actionable
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
