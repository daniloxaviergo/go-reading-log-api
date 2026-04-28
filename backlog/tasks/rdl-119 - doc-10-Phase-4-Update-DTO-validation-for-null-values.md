---
id: RDL-119
title: '[doc-10 Phase 4] Update DTO validation for null values'
status: To Do
assignee:
  - thomas
created_date: '2026-04-28 00:30'
updated_date: '2026-04-28 04:45'
labels:
  - validation
  - phase-4
  - backend
dependencies: []
documentation:
  - doc-010
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update StatsData.Validate() to allow null values for ratio fields. Create tests validating null handling scenarios for per_pages, per_mean_day, per_spec_mean_day.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Validate() accepts null for ratio fields
- [x] #2 Tests cover all null scenarios
- [x] #3 No validation errors for valid null values
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task focuses on verifying and documenting that `StatsData.Validate()` correctly accepts null values for ratio fields, and ensuring comprehensive test coverage for all null handling scenarios. Based on research of the codebase and the completed RDL-118 task, the validation logic already supports null values for pointer fields (`PerPages`, `PerMeanDay`, `PerSpecMeanDay`, `MaxDay`, `MeanGeral`).

**Current State Analysis:**
- The `StatsData.Validate()` method already handles null pointer fields correctly by only validating when the pointer is non-nil
- RDL-118 implemented the null handling logic in the service layer to return `nil` instead of `0.0`
- Existing tests cover basic null scenarios, but may need expansion for comprehensive coverage

**Technical Approach:**
1. **Verification Phase**: Confirm that `Validate()` correctly accepts null values for all ratio fields
2. **Test Enhancement**: Add comprehensive test cases covering all null scenarios for `PerPages`, `PerMeanDay`, `PerSpecMeanDay`
3. **Documentation**: Document the validation behavior and null handling patterns

**Why This Approach:**
- The validation logic already exists and works correctly (verified in code review)
- The focus is on ensuring complete test coverage to prevent regressions
- Following the pattern established in RDL-118 for null handling

### 2. Files to Modify

**Test Files (Primary Focus):**

1. **`test/unit/dashboard_response_test.go`**
   - **Review**: Existing test `TestStatsData_NewFields_Validation` already covers nil new fields
   - **Review**: Existing test `TestStatsData_PerPages_NullHandling` covers PerPages null scenario
   - **Add**: New comprehensive test `TestStatsData_RatioFields_NullValidation` that:
     - Tests all three ratio fields (`PerPages`, `PerMeanDay`, `PerSpecMeanDay`) with null values
     - Tests combinations of null and non-null values
     - Tests validation passes when all ratio fields are null
     - Tests validation passes when some ratio fields are null and others have values
   - **Add**: Test `TestStatsData_AllNullRatioFields` specifically for the acceptance criteria

2. **`test/unit/day_service_test.go`** (if exists and related)
   - **Review**: Verify tests for `CalculatePerPagesRatio` with zero previous week return nil
   - **Add**: Integration with validation tests to ensure service layer null returns pass DTO validation

**Documentation Files:**

3. **`QWEN.md`** (or project documentation)
   - **Add**: Section documenting null validation behavior for StatsData
   - **Add**: Examples of valid null scenarios
   - **Add**: Reference to RDL-118 and RDL-119 relationship

**No changes needed to:**
- `internal/domain/dto/dashboard_response.go` - Validation logic already correct
- `internal/service/dashboard/day_service.go` - Null handling implemented in RDL-118
- `internal/api/v1/handlers/dashboard_handler.go` - Already handles null correctly

### 3. Dependencies

**Prerequisites:**
- RDL-118 must be completed (Status: Done ✅)
  - RDL-118 implemented the null handling logic in the service layer
  - RDL-119 focuses on validation and test coverage

**Related Tasks:**
- RDL-111: Update StatsData DTO with new fields (Done)
  - Added the nullable pointer fields: `MaxDay`, `MeanGeral`, `PerMeanDay`, `PerSpecMeanDay`
- RDL-118: Implement null handling for ratio fields (Done)
  - Implemented service layer logic to return nil for zero denominators

**No blocking issues:**
- Validation logic is already correct
- No code changes needed in production code
- Only test enhancement and documentation required

### 4. Code Patterns

**Validation Pattern (Existing):**
```go
// In StatsData.Validate()
// Validate pointer fields when non-nil
if s.PerPages != nil {
    if *s.PerPages < 0 {
        return fmt.Errorf("per_pages cannot be negative")
    }
}

// PerMeanDay: Must be non-negative when set
if s.PerMeanDay != nil {
    if *s.PerMeanDay < 0 {
        return fmt.Errorf("per_mean_day cannot be negative")
    }
}

// PerSpecMeanDay: Must be non-negative when set
if s.PerSpecMeanDay != nil {
    if *s.PerSpecMeanDay < 0 {
        return fmt.Errorf("per_spec_mean_day cannot be negative")
    }
}
```

**Test Pattern (Following Existing):**
```go
// TestStatsData_RatioFields_NullValidation
func TestStatsData_RatioFields_NullValidation(t *testing.T) {
    testCases := []struct {
        name        string
        stats       *dto.StatsData
        expectError bool
    }{
        {
            name:        "all ratio fields nil - should pass",
            stats:       dto.NewStatsData(),
            expectError: false,
        },
        {
            name: "PerPages nil, others with values",
            stats: func() *dto.StatsData {
                perMeanDay := 1.5
                perSpecMeanDay := 2.0
                return dto.NewStatsData().
                    SetPerMeanDay(&perMeanDay).
                    SetPerSpecMeanDay(&perSpecMeanDay)
            }(),
            expectError: false,
        },
        // ... more test cases
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            err := tc.stats.Validate()
            if tc.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

**Naming Conventions:**
- Follow existing test naming: `TestStatsData_<Feature>_<Scenario>`
- Use descriptive test case names in table-driven tests
- Follow the pattern from `TestStatsData_NewFields_Validation`

**Integration Patterns:**
- Service layer returns `*float64` (nullable)
- DTO validation accepts nil pointers without error
- JSON serialization converts nil to `null` automatically

### 5. Testing Strategy

**Unit Tests (`test/unit/dashboard_response_test.go`):**

1. **TestStatsData_RatioFields_NullValidation** (NEW)
   - Table-driven test covering all null scenarios for ratio fields
   - Test cases:
     - All ratio fields nil → validation passes
     - PerPages nil, PerMeanDay and PerSpecMeanDay with values → validation passes
     - PerMeanDay nil, PerPages and PerSpecMeanDay with values → validation passes
     - PerSpecMeanDay nil, PerPages and PerMeanDay with values → validation passes
     - All ratio fields with valid non-negative values → validation passes
     - Any ratio field with negative value → validation fails

2. **TestStatsData_AllNullRatioFields_Validation** (NEW)
   - Specific test for acceptance criteria #1 and #3
   - Creates StatsData with all ratio fields as nil
   - Verifies validation passes with no errors
   - Verifies JSON serialization produces null values

3. **TestStatsData_MixedNullAndValue_RatioFields** (NEW)
   - Tests mixed scenarios where some ratio fields are null and others have values
   - Ensures validation correctly handles partial null scenarios
   - Tests edge case: null PerPages with valid PerMeanDay

**Edge Cases to Cover:**
- All ratio fields nil (empty StatsData)
- PerPages nil (previous_week_pages = 0 case from RDL-118)
- PerMeanDay nil (previous_mean = 0 case)
- PerSpecMeanDay nil (speculated_mean = 0 case)
- Combination of null and non-null ratio fields
- Negative values when fields are set (should fail validation)

**Test Execution:**
```bash
# Run specific test file
go test -v ./test/unit/dashboard_response_test.go

# Run StatsData tests only
go test -v ./test/unit/... -run TestStatsData

# Run with coverage
go test -cover ./test/unit/dashboard_response_test.go

# Run all unit tests
go test ./test/unit/...
```

**Acceptance Criteria Verification:**
- AC1: Validate() accepts null for ratio fields
  - Test: `TestStatsData_AllNullRatioFields_Validation` verifies no validation error
- AC2: Tests cover all null scenarios
  - Test: `TestStatsData_RatioFields_NullValidation` covers all combinations
- AC3: No validation errors for valid null values
  - Test: All null scenario tests verify `assert.NoError(t, err)`

### 6. Risks and Considerations

**Known Issues:**
- **None**: The validation logic is already correct; this task is about test coverage and documentation

**Potential Pitfalls:**
1. **Test Redundancy**: Existing tests may already cover some scenarios
   - Mitigation: Review existing tests before adding new ones to avoid duplication
   - Existing tests to review:
     - `TestStatsData_NewFields_Validation` - covers nil new fields
     - `TestStatsData_PerPages_NullHandling` - covers PerPages null
     - `TestValidationScenarios_Comprehensive` - covers basic validation paths

2. **Test Coverage Gaps**: May miss edge cases with mixed null/value scenarios
   - Mitigation: Use table-driven tests to systematically cover all combinations

3. **Documentation Alignment**: Ensure QWEN.md reflects the actual validation behavior
   - Mitigation: Document based on verified code behavior, not assumptions

**Trade-offs:**
- **None**: This is a test enhancement task with no production code changes

**Deployment Considerations:**
- No deployment changes required (no production code modifications)
- Tests can be merged independently
- Documentation updates are informational

**Verification Steps:**
1. Review existing tests to identify coverage gaps
2. Add new test cases for missing scenarios
3. Run full test suite: `go test ./test/unit/...`
4. Verify coverage: `go test -cover ./test/unit/dashboard_response_test.go`
5. Run `go fmt ./...` and `go vet ./...`
6. Update QWEN.md with validation behavior documentation
7. Verify all acceptance criteria are met

**Definition of Done Checklist:**
- [ ] All unit tests pass (`go test ./test/unit/...`)
- [ ] All integration tests pass (`go test ./test/integration/...`)
- [ ] `go fmt ./...` passes with no errors
- [ ] `go vet ./...` passes with no issues
- [ ] Clean Architecture layers followed (DTO validation is in domain layer)
- [ ] Error responses consistent with existing patterns (no changes to error format)
- [ ] Documentation updated in QWEN.md
- [ ] New code paths include error path tests (validation error paths)
- [ ] Tests verify both success and error responses
- [ ] Acceptance criteria #1-3 are met

**Summary of Changes:**
- **Production Code**: No changes (validation already correct)
- **Test Code**: Add 2-3 new test functions to `test/unit/dashboard_response_test.go`
- **Documentation**: Update QWEN.md with null validation behavior
- **Total Impact**: Low risk, high value for test coverage and documentation
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress

### Phase 1: Research & Analysis ✅
- Reviewed `internal/domain/dto/dashboard_response.go` - Validation logic already correctly handles null pointer fields
- Reviewed `test/unit/dashboard_response_test.go` - Existing tests cover basic null scenarios
- Identified test coverage gaps for ratio fields (`PerPages`, `PerMeanDay`, `PerSpecMeanDay`)

### Phase 2: Test Enhancement ✅
- Added `TestStatsData_RatioFields_NullValidation` - 10 test cases covering all null scenarios
- Added `TestStatsData_AllNullRatioFields_Validation` - Specific test for acceptance criteria
- Added `TestStatsData_MixedNullAndValue_RatioFields` - 7 test cases for mixed scenarios
- Added `TestStatsData_RatioFields_JSONSerialization` - JSON serialization verification
- All new tests pass ✅
- `go fmt` and `go vet` pass ✅
- All unit tests pass ✅

### Phase 3: Documentation (In Progress)
- Need to update QWEN.md with null validation behavior documentation

### Test Coverage Summary
- **Total new test cases**: 18 test cases across 4 new test functions
- **Coverage areas**:
  - All ratio fields nil
  - Individual ratio field null scenarios
  - Mixed null/value combinations
  - Negative value validation (error paths)
  - JSON serialization of null values
<!-- SECTION:NOTES:END -->

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
