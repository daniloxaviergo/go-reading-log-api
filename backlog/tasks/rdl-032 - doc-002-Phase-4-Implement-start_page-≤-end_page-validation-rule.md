---
id: RDL-032
title: '[doc-002 Phase 4] Implement start_page ≤ end_page validation rule'
status: To Do
assignee:
  - catarina
created_date: '2026-04-03 14:04'
updated_date: '2026-04-04 02:57'
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

The validation for `start_page ≤ end_page` already exists in the `internal/validation/` package (RDL-030 completed). This task focuses on **integrating the validation into the log creation flow** when a POST endpoint for logs is implemented.

**Architecture Decisions**:
- The validation function `ValidateStartEndPage()` is already implemented and tested (35 tests pass)
- Integration will occur when log creation endpoint (POST `/api/v1/projects/{project_id}/logs`) is added
- Validation should be called before database operations, following the same pattern as `ProjectsHandler.Create`
- Error response format: JSON with `error` and `details` fields matching existing patterns

**Implementation Strategy**:
1. No changes needed to validation functions (already complete from RDL-030)
2. When log creation handler is implemented:
   - Parse and validate request body
   - Call `validation.ValidateStartEndPage()` before database insert
   - Return 400 Bad Request with error details if validation fails
   - Proceed with database operation only if validation passes

**Rationale**:
- Reuse existing validation package from RDL-030 (no duplication)
- Consistent with project creation validation pattern
- Follows Clean Architecture - validation at application layer, not database
- Separation of concerns: validation in handlers, not in repositories

### 2. Files to Modify

#### Existing Files (No Changes - Already Complete):
| File | Status | Purpose |
|------|--------|---------|
| `internal/validation/validate_log.go` | ✓ Done | Already has `ValidateStartEndPage()` and `ValidateLog()` |
| `internal/validation/validate_test.go` | ✓ Done | Already has 8 test cases for log validation |
| `internal/validation/errors.go` | ✓ Done | Already has error types and helper functions |

#### Files to Create/Modify (When Log Creation is Implemented):
| File | Action | Reason |
|------|--------|--------|
| `internal/api/v1/handlers/logs_handler.go` | Modify | Add Create handler with validation integration |
| `internal/api/v1/routes.go` | Modify | Add POST route for `/api/v1/projects/{project_id}/logs` |
| `internal/repository/log_repository.go` | Modify | Add Create method to interface (if needed) |
| `internal/adapter/postgres/log_repository.go` | Modify | Add Create implementation if Create route added |

### 3. Dependencies

**Prerequisites from RDL-030**:
- ✓ Validation package created at `internal/validation/`
- ✓ `ValidateStartEndPage()` function already implemented and tested
- ✓ Error types (`ValidationError`, `ValidationErrorList`) exported

**Dependent Tasks** (for integration):
- RDL-033 (Implement status value validation) - may need similar integration pattern
- Future task for log creation endpoint (not in current scope)
- RDL-002 (Domain models/DTOs) - Log model already exists

**Blocking Issues**:
- None - validation function already exists and is tested

**Setup Steps Required**:
- No additional setup needed
- Validation package is ready to use

### 4. Code Patterns

**Pattern 1: Validation in Handler (from ProjectsHandler)**:
```go
// From projects_handler.go
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
// When log creation handler is implemented
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

**Pattern 3: Comprehensive Validation (if multiple fields need validation)**:
```go
validationErr := validation.ValidateLog(req.StartPage, req.EndPage)
if validationErr != nil && validationErr.HasErrors() {
    // Return all validation errors
}
```

**Naming Conventions** (from existing code):
- Validation functions: `ValidateX` (e.g., `ValidateStartEndPage`)
- Error codes: snake_case (e.g., `start_page_exceeds_end_page`)
- Response format: `{"error": "message", "details": {"field": "error"}}`

**Integration Patterns**:
- Call validation **before** database operations
- Return `400 Bad Request` for validation errors
- Use `validationErr.ToMap()` to serialize errors to JSON
- Log errors to stderr for debugging

### 5. Testing Strategy

**Unit Tests** (for when log creation is implemented):
- Test validation failure returns 400 status
- Test validation error format matches JSON response
- Test error details include specific field and message
- Test valid start_page ≤ end_page passes validation

**Test Cases**:
1. **Valid Input**: `start_page=10, end_page=20` → passes
2. **Equal Values**: `start_page=10, end_page=10` → passes
3. **Start > End**: `start_page=20, end_page=10` → fails with "start_page_exceeds_end_page"
4. **Negative Start**: `start_page=-1, end_page=10` → fails with "start_page_invalid"
5. **Negative End**: `start_page=10, end_page=-1` → fails with "end_page_invalid"

**Integration Tests** (for full POST flow):
- Test POST with valid log data (creates record, returns 201)
- Test POST with invalid page range (returns 400 with error details)
- Test POST with missing required fields (returns 400)
- Verify database record contains correct values
- Test project exists validation (returns 404 if project not found)

**Test Coverage Targets**:
- 100% coverage for validation logic (already achieved in RDL-030)
- Test both success and error paths for handler integration

### 6. Risks and Considerations

**Blocking Issues**:
- None identified - validation function already exists

**Potential Pitfalls**:
1. **Log Request DTO**: Need to create `LogRequest` DTO for parsing POST body
2. **Project Validation**: Should verify project exists before validating log
3. **Error Code Consistency**: Ensure error codes match Rails API if integrationRequired
4. **HTTP Status Codes**: 400 for validation, 404 for not found, 201 for success

**Trade-offs**:
1. **Validation Location**: Validation in handler vs service layer - chosen handler for simplicity
2. **Comprehensive Validation**: `ValidateLog()` validates all fields at once vs individual calls
3. **Error Detail Level**: Using `ToMap()` for field-specific errors vs generic message

**Database Constraints**:
- PostgreSQL has `CHECK (start_page <= end_page)` constraint possible but not currently implemented
- Application-level validation provides user-friendly messages
- Database constraint would be additional safety layer (future enhancement)

**Deployment Considerations**:
- No database migrations required (validation is application-level)
- No configuration changes needed
- Backward compatible (adds new validation, doesn't change existing behavior)
- Can be deployed as part of log creation feature rollout

### Implementation Checklist (When Log Creation is Implemented)

1. Create `LogRequest` DTO in `internal/domain/dto/` for POST body parsing
2. Add `Create` handler to `logs_handler.go` (similar to `ProjectsHandler.Create`)
3. Add validation call using `validation.ValidateStartEndPage()`
4. Add POST route in `routes.go`: `/api/v1/projects/{project_id}/logs`
5. Implement repository Create method if not exists
6. Write unit tests for handler (success and validation error cases)
7. Write integration tests (database interaction, project validation)
8. Run `go fmt`, `go vet`, and verify no issues
9. Update documentation for new endpoint

### Current Status Summary

**RDL-030 Completed**:
- Validation package created with 4 files (errors.go, validate_project.go, validate_log.go, validate_test.go)
- All 35 unit tests pass
- `ValidateStartEndPage()` function implements `start_page ≤ end_page` check
- Error codes and messages are properly formatted

**RDL-032 Current State**:
- No changes needed to validation code (already complete)
- Validation function ready for integration into log creation flow
- When POST `/api/v1/projects/{project_id}/logs` is implemented, just integrate validation following existing pattern from `ProjectsHandler`
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
