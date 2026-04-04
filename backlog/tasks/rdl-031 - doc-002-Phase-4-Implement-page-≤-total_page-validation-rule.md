---
id: RDL-031
title: '[doc-002 Phase 4] Implement page ≤ total_page validation rule'
status: To Do
assignee:
  - catarina
created_date: '2026-04-03 14:04'
updated_date: '2026-04-04 01:32'
labels:
  - phase-4
  - validation-rule
  - business-logic
dependencies: []
references:
  - 'PRD Section: Validation Rules - page ≤ total_page'
  - 'PRD Section: Files to Modify - project_repository.go'
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement validation for project page ≤ total_page constraint. Create validation function in internal/validation/ package and integrate into project creation/update flow. Return appropriate error with error code and message.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Validation function checks page ≤ total_page
- [ ] #2 Error returned when constraint violated
- [ ] #3 Error includes error code and descriptive message
- [ ] #4 Validation logic matches Rails behavior
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The implementation will create a POST endpoint for projects that includes page ≤ total_page validation using the existing validation package from RDL-030.

**Strategy**:
- Create a new `Create` handler in `ProjectsHandler` that accepts JSON payload for project creation
- Integrate validation package to check `page ≤ total_page` constraint before database insert
- Return proper error responses with error code and message when validation fails
- Ensure validation error handling follows existing patterns in the codebase

**Architecture Decisions**:
- Follow Clean Architecture: validation in `internal/validation/`, handler in `internal/api/v1/handlers/`
- Use existing `ValidationError` type from validation package for consistency
- Validate in handler layer before repository call (prevents invalid data from reaching database)
- Return HTTP 400 Bad Request for validation errors (matching Rails behavior)

**Rationale**:
- The validation package already has `ValidatePage(page, totalPage)` implemented and tested
- No existing POST endpoints exist, so this follows new endpoint pattern
- Validation must happen before database insert to prevent constraint violations
- Error response format must match existing error patterns for consistency

### 2. Files to Modify

#### New Files to Create:
| File | Purpose |
|------|---------|
| `internal/api/v1/handlers/projects_handler_create_test.go` | Unit tests for Create handler |
| `test/integration/projects_create_integration_test.go` | Integration tests for POST /projects |

#### Existing Files to Modify:
| File | Change Type | Reason |
|------|-------------|--------|
| `internal/api/v1/handlers/projects_handler.go` | Modify | Add `Create` handler method with validation integration |
| `internal/domain/dto/project_request.go` | Create | New struct for incoming project creation payload |
| `internal/repository/project_repository.go` | Modify | Add `Create` method interface |

#### Existing Files to Reference (No Changes):
| File | Purpose |
|------|---------|
| `internal/validation/validate_project.go` | Validation function `ValidatePage(page, totalPage)` already exists |
| `internal/validation/errors.go` | `ValidationError` struct and helper functions |
| `internal/adapter/postgres/project_repository.go` | PostgreSQL implementation reference |

### 3. Dependencies

**Prerequisites**:
- ✅ RDL-030 complete - validation package exists with `ValidatePage` function
- ✅ RDL-002 complete - domain models and DTOs exist
- ✅ RDL-006 complete - existing handlers follow established patterns

**Blocking Issues**:
- None - all prerequisites are met

**Setup Steps Required**:
1. Create `project_request.go` DTO for request payload
2. Add `Create` method to `ProjectRepository` interface
3. Implement `Create` in `ProjectRepositoryImpl`
4. Add `Create` handler to `ProjectsHandler`
5. Update routes in `internal/api/v1/routes.go`

### 4. Code Patterns

**Pattern 1: Request DTO Structure** (based on `project_response.go`)
```go
type ProjectRequest struct {
    Name       string `json:"name"`
    TotalPage  int    `json:"total_page"`
    Page       int    `json:"page"`
    StartedAt  *string `json:"started_at,omitempty"`
    Reinicia   bool   `json:"reinicia,omitempty"`
}
```

**Pattern 2: Validation Integration in Handler**
```go
// Validate before processing
if err := validation.ValidateProject(req.Page, req.TotalPage, ""); err != nil {
    if err.HasErrors() {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "error": "validation failed",
            "details": err.ToMap(),
        })
        return
    }
}
```

**Pattern 3: Error Response Format** (matching existing handlers)
```json
{
  "error": "validation failed",
  "details": {
    "page": "page (100) cannot exceed total_page (50)"
  }
}
```

**Naming Conventions**:
- Struct fields: snake_case JSON keys matching Rails API
- Handler methods: Title case (e.g., `Create`)
- Error codes: snake_case (e.g., `page_exceeds_total`)

**Integration Points**:
- Call validation in handler BEFORE repository insert
- Use existing error response format from handlers
- Follow repository interface pattern for database operations

### 5. Testing Strategy

**Unit Tests (handlers package)**:
- **Validation Tests**: Test `ValidatePage` integration in handler
  - Valid: page ≤ total_page (should proceed)
  - Invalid: page > total_page (should return 400)
  - Edge: page = total_page (should proceed)
  - Edge: page = 0 (should proceed if total_page > 0)

- **Handler Response Tests**:
  - Success: returns 201 Created with project JSON
  - Validation error: returns 400 Bad Request with error details
  - Database error: returns 500 Internal Server Error

**Integration Tests (test/integration package)**:
- **End-to-End Tests**:
  - Create project with valid data (page ≤ total_page)
  - Create project with invalid data (page > total_page) - expects 400
  - Verify validation error response format matches expectations
  - Test database persistence: retrieve created project via GET

**Test Coverage Targets**:
- All validation rule combinations (page values relative to total_page)
- Error response format verification
- HTTP status code correctness (201 for success, 400 for validation)
- Edge cases: zero values, negative values, equal values

**Testing Approach**:
- Unit tests in `handlers` package using `testing` package
- Integration tests in `test/integration` using test database
- Use existing test helpers from `test/test_helper.go`
- Mock repository for unit tests, real PostgreSQL for integration

### 6. Risks and Considerations

**Blocking Issues**:
- None identified

**Potential Pitfalls**:
1. **Database Constraint vs Application Validation**: The database may have a CHECK constraint for page ≤ total_page, but we should validate in application layer for user-friendly error messages
2. **Rails API Differences**: Rails may use different validation logic; must verify behavior matches Rails (see Rails tests in spec/models/project_spec.rb)
3. **Timing**: Validation should happen BEFORE database insert to prevent constraint violations

**Trade-offs**:
1. **Validation Location**: Application-level validation provides better error messages than database constraint errors
2. **Error Response Format**: Using `{"error": "...", "details": {...}}` format matches existing error patterns
3. **HTTP Status Codes**: 400 for validation errors, 201 for successful creation (standard REST conventions)

**Deployment Considerations**:
- No database migrations required (validation is application-level)
- Backward compatible (adds new endpoint, doesn't change existing GET endpoints)
- Can be deployed as part of Phase 4 rollout
- No configuration changes needed

**Validation Logic Details**:
- The existing `ValidatePage(page, totalPage)` function from RDL-030 already validates: page >= 0 AND page <= totalPage
- Error code: `page_exceeds_total` when page > totalPage
- Error code: `page_invalid` when page < 0
- Rails behavior: Must match Rails validation (check existing specs)
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
