---
id: RDL-035
title: '[doc-002 Phase 5] Test edge cases for derived calculations'
status: To Do
assignee:
  - thomas
created_date: '2026-04-03 14:05'
updated_date: '2026-04-04 05:16'
labels:
  - phase-5
  - edge-cases
  - testing
dependencies: []
references:
  - 'PRD Section: Acceptance Criteria - AC4'
  - AC7
  - 'PRD Section: Validation Rules - edge cases'
documentation:
  - doc-002
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create comprehensive tests for derived calculation edge cases: zero total_page (progress), no logs (days_unreading), 100% progress (finished_at), and invalid status values. Verify all calculations handle errors gracefully and return expected defaults.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Zero total_page returns 0.00 progress
- [ ] #2 No logs uses started_at for days_unreading or returns 0
- [ ] #3 100% progress returns null finished_at
- [ ] #4 Invalid status values handled with error
- [ ] #5 All edge cases documented
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Create comprehensive test coverage for derived calculation edge cases as specified in the acceptance criteria. The tests will cover:

**Test Scenarios:**
1. **Zero total_page**: Verify progress returns 0.00 when total_page is 0 or negative
2. **No logs scenario**: Verify days_unreading calculation handles missing logs (uses started_at or returns 0)
3. **100% progress scenario**: Verify finished_at returns null for finished books (page >= total_page)
4. **Invalid status values**: Verify validation rejects invalid status values with proper error response
5. **Edge case documentation**: Ensure all edge cases are documented in test comments

**Approach:**
- Unit tests in `internal/domain/models/project_test.go` - test calculation methods in isolation
- Integration tests in `test/integration/` - verify HTTP endpoints handle edge cases correctly
- Validation tests in `internal/validation/validate_test.go` - already exists, add more status edge cases
- Test helper functions for creating test data with edge case values

**Why this approach:**
- Clean Architecture requires separation of concerns (models vs adapters vs handlers)
- Unit tests verify calculation logic without database dependencies
- Integration tests verify HTTP handlers return proper error responses
- Follow existing test patterns in the codebase

### 2. Files to Modify

#### Unit Tests (internal/domain/models/)
| File | Action | Reason |
|------|--------|-- ------|
| `project_test.go` | Modify | Add comprehensive edge case tests for derived calculations |
| `project.go` | Read only | Reference existing calculation methods for test scenarios |

#### Validation Tests (internal/validation/)
| File | Action | Reason |
|------|--------|-- ------|
| `validate_test.go` | Modify | Add edge case tests for invalid status values |

#### Integration Tests (test/integration/)
| File | Action | Reason |
|------|--------|-- ------|
| `projects_integration_test.go` | Create | Test HTTP endpoints with edge case data |
| `projects_integration_test.go` | Modify | Add tests for status validation and derived fields |

#### Test Helpers (test/)
| File | Action | Reason |
|------|-- ------|-- ------|
| `testdata/edge_cases.json` | Create | JSON fixtures for edge case test data |

### 3. Dependencies

**Prerequisites (Already Completed):**
- Task RDL-034 (JSON response comparison) - ✅ Done
- Task RDL-020 (progress calculation) - ✅ Done
- Task RDL-021 (status ranges config) - ✅ Done
- Task RDL-022 (status determination) - ✅ Done
- Task RDL-023 (days_unreading calculation) - ✅ Done
- Task RDL-024 (median_day calculation) - ✅ Done
- Task RDL-025 (finished_at calculation) - ✅ Done

**Required Infrastructure:**
- Existing `CalculateProgress()`, `CalculateStatus()`, `CalculateDaysUnreading()` methods already handle edge cases
- Validation package with `ValidateStatus()` to test invalid status values
- Test database with test data for integration tests

**No additional setup required** - all prerequisites are in place.

**Implementation Notes:**
- Edge cases already handled in calculation methods (return 0.00 or nil as appropriate)
- Invalid status validation already exists in `validate_test.go`
- HTTP handlers already return 400 for validation errors

### 4. Code Patterns

**Unit Test Pattern (from existing tests):**
```go
func TestProject_CalculateXyz_EdgeCase(t *testing.T) {
    ctx := context.Background()
    
    // Setup edge case: zero total_page
    project := &Project{
        ID:        1,
        Name:      "Test Project",
        TotalPage: 0,  // Edge case: zero total_page
        Page:      50,
    }
    
    result := project.CalculateProgress()
    
    if result == nil {
        t.Fatal("Expected non-nil result, got nil")
    }
    
    if *result != 0.0 {
        t.Errorf("Expected 0.00 for zero total_page, got %.2f", *result)
    }
}
```

**Integration Test Pattern (from existing tests):**
```go
func TestCreate_Project_ValidationError_Status(t *testing.T) {
    req := dto.ProjectRequest{
        Name:      "Test Project",
        TotalPage: 100,
        Page:      50,
        Status:    "invalid_status",  // Edge case: invalid status
    }
    
    body, _ := json.Marshal(req)
    recorder := httptest.NewRecorder()
    
    handler.Create(recorder, httptest.NewRequest(http.MethodPost, "/api/v1/projects", bytes.NewBuffer(body)))
    
    if recorder.Code != http.StatusBadRequest {
        t.Errorf("Expected status 400 for invalid status, got %d", recorder.Code)
    }
}
```

**Validation Test Pattern (existing):**
```go
func TestValidateStatus_Invalid(t *testing.T) {
    err := ValidateStatus("invalid_status")
    
    if err == nil {
        t.Fatal("Expected validation error, got nil")
    }
    
    if err.Code != "invalid_status" {
        t.Errorf("Expected error code 'invalid_status', got '%s'", err.Code)
    }
}
```

### 5. Testing Strategy

**Unit Tests (internal/domain/models/project_test.go):**

1. **TestZeroTotalPage**: Verify progress returns 0.00 when total_page is 0
   - Test with total_page = 0, page > 0
   - Test with total_page = 0, page = 0
   - Test with total_page < 0 (negative)

2. **TestNoLogsDaysUnreading**: Verify days_unreading with no logs
   - Test with no logs and no started_at → returns 0
   - Test with no logs but has started_at → returns days since started_at
   - Test with nil logs slice

3. **Test100PercentProgressFinishedAt**: Verify finished_at returns null when 100% complete
   - Test with page >= total_page and logs → returns last log date
   - Test with page >= total_page and no logs → returns nil
   - Test with page < total_page → calculates estimated finish date

4. **TestInvalidStatusValidation**: Verify invalid status values are rejected
   - Test with "invalid_status" string
   - Test with empty string ""
   - Test with random case variations like "Running" (should fail - case sensitive)

5. **TestEdgeCaseDocumentation**: Verify all edge cases are documented
   - All test functions have comments describing the edge case
   - Edge cases cover: zero values, negative values, nil values, boundary values

**Integration Tests (test/integration/):**

1. **TestProjectsHandler_Create_ValidationError**: Test HTTP handler returns 400 for validation errors
   - Test invalid status value
   - Test page > total_page
   - Test total_page <= 0

2. **TestProjectsHandler_Index_EdgeCases**: Test index endpoint with edge case data
   - Project with zero total_page
   - Project with no logs
   - Project with 100% progress

**Test Execution:**
- Run all tests: `go test ./... -v`
- Run only model tests: `go test ./internal/domain/models/... -v`
- Run only validation tests: `go test ./internal/validation/... -v`
- Integration tests: `go test ./test/integration/... -v`

### 6. Risks and Considerations

**Known Edge Cases Already Handled:**
- Zero/negative total_page in `CalculateProgress()` → returns 0.00
- Zero/negative page in `CalculateProgress()` → returns 0.00
- No logs in `CalculateStatus()` → returns "unstarted"
- No logs in `CalculateDaysUnreading()` → returns 0 or uses started_at
- No started_at in `CalculateMedianDay()` → returns 0.00
- Finished book in `CalculateFinishedAt()` → returns last log date or nil

**Potential Pitfalls:**
1. **Floating point precision**: Use `math.Round()` for 2 decimal places, test with tolerance if needed
2. **Date calculations**: Edge cases near midnight may have off-by-one errors; use the same date-only comparison pattern as existing code
3. **Null vs empty slice**: Distinguish between `nil` slice and empty slice `[]*LogResponse{}`

**Trade-offs:**
- Test coverage vs. development time: Focus on acceptance criteria scenarios (AC1-AC4) rather than exhaustive combinations
- Integration tests require database setup; unit tests are faster but less comprehensive

**Deployment Considerations:**
- No database schema changes required
- No migration scripts needed
- All changes are test-only (no production code modifications)

**Acceptance Criteria Verification:**
- AC1: Zero total_page returns 0.00 progress → test `TestCalculateProgress_ZeroTotalPage`
- AC2: No logs uses started_at or returns 0 → test `TestCalculateDaysUnreading_NoLogsNoStartedAt`
- AC3: 100% progress returns null finished_at (or last log date) → test `TestCalculateFinishedAt_FinishedBook`
- AC4: Invalid status values handled with error → test `TestValidateStatus_Invalid`
- AC5: All edge cases documented → document in test function comments
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
