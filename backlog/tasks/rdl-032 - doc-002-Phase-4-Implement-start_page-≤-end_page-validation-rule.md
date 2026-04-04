---
id: RDL-032
title: '[doc-002 Phase 4] Implement start_page ≤ end_page validation rule'
status: To Do
assignee:
  - workflow
created_date: '2026-04-03 14:04'
updated_date: '2026-04-04 03:00'
labels:
  - phase-4
  - validation-rule
  - business-logic
dependencies: []
references:
  - 'PRD Section: Validation Rules - start_page ≤ end_page'
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement validation for log start_page ≤ end_page constraint. Create validation function in internal/validation/ package and integrate into log creation flow. Return appropriate error with error code and message when constraint violated.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Validation function checks start_page ≤ end_page
- [ ] #2 Error returned when constraint violated
- [ ] #3 Error includes error code and descriptive message
- [ ] #4 Validation logic matches Rails behavior
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The validation for `start_page ≤ end_page` was completed in RDL-030 (shared validation package). This task focuses on **integrating the validation into the log creation flow** when the POST endpoint for logs is implemented.

**Key Findings from Completed Tasks**:
- RDL-030: Validation package created with `ValidateStartEndPage()` and `ValidateLog()` functions, 35 tests pass
- RDL-031: Page validation pattern established (`ValidatePage()`, validation in handler, 400 error responses)
- RDL-031 Implementation: POST /projects endpoint with `ValidateProject()` integration, 29 tests pass total

**Architecture Decisions**:
- Reuse existing validation functions from `internal/validation/` package (no new validation code needed)
- Follow the same pattern as `ProjectsHandler.Create` for consistency
- Call validation BEFORE database operations in the handler layer
- Error responses use the same format: `{"error": "validation failed", "details": {...}}`

**Implementation Strategy**:
1. No changes to validation functions (already complete from RDL-030)
2. When log creation endpoint is implemented:
   - Parse request body to `LogRequest` DTO
   - Call `validation.ValidateStartEndPage(req.StartPage, req.EndPage)` or `ValidateLog()`
   - Return 400 Bad Request if validation fails
   - Proceed with database insert if validation passes

**Rationale**:
- Follows pattern from RDL-031 implementation (projects creation validation)
- Validation already tested and comprehensive (8 test cases for log validation)
- Separation of concerns: validation in handlers, not repositories or models
- Reuse existing error handling infrastructure

### 2. Files to Modify

#### Existing Files (No Changes - Already Complete):
| File | Status | Purpose |
|------|--------|--|
| `internal/validation/validate_log.go` | ✓ Already Complete | `ValidateStartEndPage()` and `ValidateLog()` functions |
| `internal/validation/validate_test.go` | ✓ Already Complete | 8 test cases for log validation |
| `internal/validation/errors.go` | ✓ Already Complete | `ValidationError`, `ValidationErrorList` types |

#### Files to Modify (When Log Creation is Implemented):
| File | Action | Reason |
|------|--|------|
| `internal/api/v1/handlers/logs_handler.go` | Modify | Add Create handler with validation integration |
| `internal/domain/dto/log_request.go` | Create | New DTO for log creation request body |
| `internal/api/v1/routes.go` | Modify | Add POST route for `/api/v1/projects/{project_id}/logs` |
| `test/test_helper.go` | Modify | Add Create method to `MockLogRepository` |

### 3. Dependencies

**Prerequisites from RDL-030**:
- ✅ Validation package at `internal/validation/`
- ✅ `ValidateStartEndPage(startPage, endPage)` - checks start_page <= end_page
- ✅ `ValidateLog(startPage, endPage)` - comprehensive log validation
- ✅ 8 unit tests for log validation (all passing)

**Prerequisites from RDL-031**:
- ✅ `ValidatePage(page, totalPage)` - page <= total_page pattern established
- ✅ Handler validation integration pattern (ProjectsHandler.Create)
- ✅ Error response format: `{"error": "validation failed", "details": {...}}`
- ✅ HTTP 400 for validation errors

**Dependent Tasks** (for log creation endpoint):
- Future task: Implement POST /api/v1/projects/{project_id}/logs endpoint
- Future task: Create `LogRequest` DTO
- Future task: Implement repository Create method if needed

**Blocking Issues**:
- None - validation function already exists and tested

**Setup Steps Required**:
- No additional setup needed
- Validation package is ready for use

### 4. Code Patterns

**Pattern 1: Validation in Handler (from ProjectsHandler - already implemented)**:
```go
validationErr := validation.ValidateProject(req.Page, req.TotalPage, status)
if validationErr != nil && validationErr.HasErrors() {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error":   "validation failed",
        "details": validationErr.ToMap(),
    })
    return
}
```

**Pattern 2: Log Validation (to be used in logs handler)**:
```go
validationErr := validation.ValidateStartEndPage(req.StartPage, req.EndPage)
if validationErr != nil {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error":   "validation failed",
        "details": validationErr.ToMap(),
    })
    return
}
```

**Pattern 3: Comprehensive Validation (multiple fields)**:
```go
validationErr := validation.ValidateLog(req.StartPage, req.EndPage)
if validationErr != nil && validationErr.HasErrors() {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error":   "validation failed",
        "details": validationErr.ToMap(),
    })
    return
}
```

**Naming Conventions** (established in RDL-030/RDL-031):
- Validation functions: `ValidateX` (e.g., `ValidateStartEndPage`)
- Error codes: snake_case (e.g., `start_page_exceeds_end_page`)
- JSON fields: snake_case (e.g., `start_page`, `end_page`)
- Response format: `{"error": "message", "details": {"field": "error"}}`

**Integration Points**:
- Call validation in handler BEFORE database operations
- Return `400 Bad Request` for validation errors
- Use `validationErr.ToMap()` to serialize errors
- Log errors to stderr for debugging

### 5. Testing Strategy

**Unit Tests** (when log creation handler is implemented):
- Test `ValidateStartEndPage` with valid inputs (start ≤ end, including equal values)
- Test validation failure returns 400 status
- Test error response format matches expectations
- Test error details include specific field and message

**Test Cases** (from existing validate_test.go, all 8 pass):
1. `TestValidateStartEndPage_Valid` - start=10, end=10 → passes
2. `TestValidateStartEndPage_Valid` - start=10, end=20 → passes
3. `TestValidateStartEndPage_Valid` - start=0, end=0 → passes
4. `TestValidateStartEndPage_StartNegative` - start=-1 → fails
5. `TestValidateStartEndPage_EndNegative` - end=-1 → fails
6. `TestValidateStartEndPage_StartExceedsEnd` - start=20, end=10 → fails
7. `TestValidateLog_Valid` - calls `ValidateLog(10, 20)` → passes
8. `TestValidateLog_StartExceedsEnd` - calls `ValidateLog(20, 10)` → fails

**Integration Tests** (when full POST endpoint implemented):
- Test POST with valid log data (start ≤ end, project exists) → 201 Created
- Test POST with start > end → 400 Bad Request with error details
- Test POST with negative start_page → 400 Bad Request
- Test POST with project not found → 404 Not Found
- Verify database record contains correct values

**Test Coverage Targets**:
- 100% coverage for validation logic (already achieved: 8 test cases)
- Test both success and error paths for handler integration
- Integration tests must use real database (not mocks)

**Testing Approach**:
- Unit tests in `handlers` package using `testing` package
- Integration tests in `test/integration` using test database
- Use existing test helpers from `test/test_helper.go`
- Use `testing-expert` subagent for test execution

### 6. Risks and Considerations

**Blocking Issues**:
- None - validation function already exists

**Potential Pitfalls**:
1. **LogRequest DTO**: Need to create request DTO matching database schema
2. **Project Validation**: Should verify project exists before validating log
3. **Error Code Consistency**: Error codes from RDL-030: `start_page_invalid`, `end_page_invalid`, `start_page_exceeds_end_page`
4. **HTTP Status Codes**: 400 for validation, 404 for not found, 201 for success

**Trade-offs**:
1. **Validation Location**: Handler validation (like RDL-031) vs service layer
2. **Validation Type**: Single call `ValidateLog()` vs individual field checks
3. **Error Detail**: Using `ToMap()` for field-specific errors vs generic message

**Database Constraints**:
- PostgreSQL could have `CHECK (start_page <= end_page)` constraint
- Currently no database constraint (validation is application-level)
- Application-level validation provides better user-facing error messages

**Deployment Considerations**:
- No database migrations required
- No configuration changes needed
- Backward compatible (adds new endpoint when implemented)
- Can be deployed with other Phase 4 features

### Implementation Checklist (When Log Creation is Implemented)

1. Create `LogRequest` DTO in `internal/domain/dto/log_request.go`
2. Add `Create` handler to `logs_handler.go` (similar to `ProjectsHandler.Create`)
3. Add validation call using `validation.ValidateStartEndPage()` or `ValidateLog()`
4. Add POST route in `routes.go`: `/api/v1/projects/{project_id}/logs`
5. Add `Create` method to `MockLogRepository` in `test/test_helper.go`
6. Write unit tests for handler (success and validation error cases)
7. Write integration tests (database interaction, project validation)
8. Run `go fmt`, `go vet`, and verify no issues
9. Update documentation for new endpoint

### Current Status Summary

**RDL-030 Completed** (2026-04-04):
- Validation package created with 4 files (errors.go, validate_project.go, validate_log.go, validate_test.go)
- 35 unit tests pass, go fmt and go vet pass, application builds successfully
- `ValidateStartEndPage()` implements `start_page ≤ end_page` check with error codes

**RDL-031 Completed** (2026-04-04):
- POST /projects endpoint with validation integration
- `ValidatePage(page, totalPage)` used in `ProjectsHandler.Create`
- HTTP 400 for validation errors, 201 for success
- 29 total tests pass (17 unit + 12 integration)

**RDL-032 Current State**:
- No changes needed to validation code (already complete from RDL-030)
- Validation function `ValidateStartEndPage()` ready for integration
- When POST `/api/v1/projects/{project_id}/logs` is implemented:
  - Follow existing pattern from `ProjectsHandler.Create`
  - Call `validation.ValidateStartEndPage()` before database operation
  - Use existing error response format
  - Use `error.Code` values: `start_page_invalid`, `end_page_invalid`, `start_page_exceeds_end_page`
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
