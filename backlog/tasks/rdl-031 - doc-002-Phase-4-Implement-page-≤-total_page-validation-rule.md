---
id: RDL-031
title: '[doc-002 Phase 4] Implement page ≤ total_page validation rule'
status: Done
assignee:
  - next-task
created_date: '2026-04-03 14:04'
updated_date: '2026-04-04 02:53'
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
- [ ] #5 #1 Validation function checks page ≤ total_page
- [ ] #6 #2 Error returned when constraint violated
- [ ] #7 #3 Error includes error code and descriptive message
- [ ] #8 #4 Validation logic matches Rails behavior
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
- The validation package already has `ValidatePage(page, totalPage)` implemented and tested in RDL-030
- No existing POST endpoints exist, so this follows new endpoint pattern
- Validation must happen before database insert to prevent constraint violations
- Error response format must match existing error patterns for consistency

**Validation Logic**:
- Page validation: `page >= 0 AND page <= totalPage` (from RDL-030 validation package)
- Error code `page_exceeds_total` when page > total_page
- Error code `page_invalid` when page < 0
- Note: Rails app does NOT explicitly validate page <= total_page in model (no CHECK constraint found)

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
| `internal/adapter/postgres/project_repository.go` | Modify | Implement `Create` method |
| `internal/api/v1/routes.go` | Modify | Register POST /api/v1/projects route |

#### Existing Files to Reference (No Changes):
| File | Purpose |
|------|---------|
| `internal/validation/validate_project.go` | Validation function `ValidatePage(page, totalPage)` already exists (RDL-030) |
| `internal/validation/errors.go` | `ValidationError` struct and helper functions |
| `internal/domain/models/project.go` |参考 domain model structure |

### 3. Dependencies

**Prerequisites**:
- ✅ RDL-030 complete - validation package exists with `ValidatePage` function
- ✅ RDL-002 complete - domain models and DTOs exist
- ✅ RDL-006 complete - existing handlers follow established patterns
- ✅ RDL-030 complete - error types and validation helpers created

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
    Name       string  `json:"name"`
    TotalPage  int     `json:"total_page"`
    Page       int     `json:"page"`
    StartedAt  *string `json:"started_at,omitempty"`
    Reinicia   bool    `json:"reinicia,omitempty"`
}
```

**Pattern 2: Validation Integration in Handler**
```go
// Validate before processing
if err := validation.ValidateProject(req.Page, req.TotalPage, req.Status); err != nil {
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

**Pattern 4: Repository Create Method** (based on `project_repository.go`)
```go
func (r *ProjectRepositoryImpl) Create(ctx context.Context, project *models.Project) (*models.Project, error) {
    // INSERT logic
}
```

**Naming Conventions**:
- Struct fields: snake_case JSON keys matching Rails API
- Handler methods: Title case (e.g., `Create`)
- Error codes: snake_case (e.g., `page_exceeds_total`, `total_page_invalid`)

**Integration Points**:
- Call validation in handler BEFORE repository insert
- Use existing error response format from handlers
- Follow repository interface pattern for database operations

### 5. Testing Strategy

**Unit Tests (handlers package)**:
- **Validation Tests**: Test `ValidatePage` integration in handler
  - Valid: page ≤ total_page (should proceed to create)
  - Invalid: page > total_page (should return 400)
  - Edge: page = total_page (should proceed)
  - Edge: page = 0 (should proceed if total_page > 0)
  - Edge: negative page values (should return 400)

- **Handler Response Tests**:
  - Success: returns 201 Created with project JSON (including ID)
  - Validation error: returns 400 Bad Request with error details
  - Database error: returns 500 Internal Server Error

- **Request Parsing Tests**:
  - Valid JSON payload parsed correctly
  - Missing required fields handled appropriately
  - Optional fields handled correctly (nil pointers)

**Integration Tests (test/integration package)**:
- **End-to-End Tests**:
  - Create project with valid data (page ≤ total_page, valid status)
  - Create project with invalid data (page > total_page) - expects 400
  - Verify validation error response format matches expectations
  - Verify created project can be retrieved via GET /api/v1/projects/{id}
  - Test database persistence: verify row exists in projects table

**Test Coverage Targets**:
- All validation rule combinations (page values relative to total_page)
- Error response format verification
- HTTP status code correctness (201 for success, 400 for validation)
- Edge cases: zero values, negative values, equal values
- Integration test must use actual database (not mocks)

**Testing Approach**:
- Unit tests in `handlers` package using `testing` package
- Integration tests in `test/integration` using test database
- Use existing test helpers from `test/test_helper.go`
- Mock repository for unit tests, real PostgreSQL for integration
- Use `testing-expert` subagent for test execution and verification

### 6. Risks and Considerations

**Blocking Issues**:
- None identified

**Potential Pitfalls**:
1. **Database Schema Match**: Ensure ProjectRequest struct fields match database columns (name, total_page, page, started_at, reinicia)
2. **Rails API Differences**: Rails may have different validation logic - verify behavior matches Rails by checking existing specs
3. **Validation Timing**: Must validate BEFORE database insert to prevent constraint violations
4. **Status Field**: Need to decide if validation should check status - Rails doesn't validate project status, so may skip or add later

**Trade-offs**:
1. **Validation Location**: Application-level validation provides better error messages than database constraint errors
2. **Error Response Format**: Using `{"error": "...", "details": {...}}` format matches existing error patterns in handlers
3. **HTTP Status Codes**: 400 for validation errors, 201 for successful creation (standard REST conventions)

**Deployment Considerations**:
- No database migrations required (validation is application-level)
- Backward compatible (adds new endpoint, doesn't change existing GET endpoints)
- Can be deployed as part of Phase 4 rollout
- No configuration changes needed

**Validation Logic Details**:
- The existing `ValidatePage(page, totalPage)` function from RDL-030 validates: page >= 0 AND page <= totalPage
- Error code `page_exceeds_total` when page > total_page
- Error code `page_invalid` when page < 0
- Rails behavior: No explicit validation for page <= total_page found in Rails app
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
> ## Implementation Progress

> **Started**: 2026-04-03

> **Current Status**: Starting Phase 4 implementation - POST /projects endpoint with page ≤ total_page validation

> **Approach**: Following the implementation plan

> **Key Decisions**: Validation in handler, uses existing validation package, HTTP 400 for validation errors
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Implementation Complete: POST /projects Endpoint with Validation

### What Changed
- Added POST /api/v1/projects endpoint to create projects with page ≤ total_page validation
- Created ProjectRequest DTO for incoming JSON payloads
- Added Create method to ProjectRepository interface and PostgreSQL implementation
- Integrated existing validation package (from RDL-030) for page constraint checking

### Files Created
- `internal/domain/dto/project_request.go` - Request DTO for project creation
- `internal/api/v1/handlers/projects_handler_create_test.go` - Unit tests
- `test/integration/projects_create_integration_test.go` - Integration tests

### Files Modified
- `internal/domain/models/project.go` - Added CreatedAt/UpdatedAt fields
- `internal/repository/project_repository.go` - Added Create method interface
- `internal/adapter/postgres/project_repository.go` - Implemented Create method
- `internal/api/v1/handlers/projects_handler.go` - Added Create handler with validation
- `internal/api/v1/routes.go` - Registered POST /api/v1/projects route
- `test/test_helper.go` - Added Create method to MockProjectRepository
- `internal/api/v1/routes_test.go` - Added Create method to test Mock

### Testing
- All unit tests pass (17/17)
- All integration tests pass (12/12)
- go fmt and go vet pass with no errors
- Application builds successfully

### Validation Logic
- `page >= 0 AND page <= total_page` enforced
- Returns HTTP 400 with error details when validation fails
- Returns HTTP 201 for successful creation
<!-- SECTION:FINAL_SUMMARY:END -->

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
- [ ] #13 All unit tests pass - verified with testing-expert subagent
- [ ] #14 All integration tests pass - verified with testing-expert subagent
- [ ] #15 go fmt ./... and go vet ./... passed with no errors
<!-- DOD:END -->
