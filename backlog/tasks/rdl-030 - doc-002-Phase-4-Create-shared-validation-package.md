---
id: RDL-030
title: '[doc-002 Phase 4] Create shared validation package'
status: To Do
assignee:
  - catarina
created_date: '2026-04-03 14:04'
updated_date: '2026-04-04 00:43'
labels:
  - phase-4
  - validation-package
  - code-structure
dependencies: []
references:
  - 'PRD Section: Technical Decisions - Decision 5: Shared Validation Logic'
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create a validation package at `internal/validation/` with functions for project and log validation rules. Include page ≤ total_page, start_page ≤ end_page, and status value validation. Package should export reusable validation functions.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Validation package created at internal/validation/
- [ ] #2 Validation functions exported for reuse
- [ ] #3 Functions include page, total_page, start_page, end_page, status validation
- [ ] #4 Documentation included
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The validation package will be implemented following Clean Architecture principles with a focus on testability and reusability.

**Architecture Decisions**:
- Create `internal/validation/` package with Go modules for each validation rule
- Export functions following Go conventions (`ValidateX` pattern)
- Use error wrapping with descriptive messages for context
- Separate validation logic from business logic (models) to maintain single responsibility
- Consider database constraints where appropriate (e.g., page <= total_page could be enforced at DB level)

**Rationale**:
- Clean separation of concerns - validation rules are reusable across handlers and services
- Testable in isolation without database dependencies
- Consistent error handling across the application
- Follows existing codebase patterns (similar to `internal/config/` and `internal/logger/`)

### 2. Files to Modify

#### New Files to Create:
| File | Purpose |
|------|---------|
| `internal/validation/errors.go` | Custom validation error types with codes and messages |
| `internal/validation/validate_project.go` | Project-related validations (page, status, etc.) |
| `internal/validation/validate_log.go` | Log-related validations (start_page, end_page) |
| `internal/validation/validate_test.go` | Unit tests for all validation functions |
| `internal/validation/errors_test.go` | Tests for custom error types |

#### Existing Files to Reference:
| File | Purpose |
|------|---------|
| `internal/domain/models/project.go` | Existing validation logic (CalculateStatus, etc.) |
| `internal/config/config.go` | Configuration patterns (Getters, defaults) |
| `internal/domain/dto/*.go` | DTO structures for validation targets |

### 3. Dependencies

**Prerequisites**:
- RDL-004 must be complete (configuration management) - validation errors may need config for range validations
- RDL-006 must be complete (handlers) - validation will be called from handlers
- RDL-002 must be complete (domain models) - validation works with existing DTOs/models

**Blocking Issues**:
- None identified

**Setup Steps Required**:
- Create `internal/validation/` directory
- Run `go mod tidy` after file creation

### 4. Code Patterns

**Pattern 1: Error Structure**
```go
type ValidationError struct {
    Code    string
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

**Pattern 2: Validation Functions**
```go
func ValidatePage(page int, totalPage int) *ValidationError {
    if page > totalPage {
        return &ValidationError{
            Code:    "page_exceeds_total",
            Field:   "page",
            Message: fmt.Sprintf("page (%d) cannot exceed total_page (%d)", page, totalPage),
        }
    }
    return nil
}
```

**Pattern 3: Status Value Validation**
```go
func IsValidStatus(status string) bool {
    switch status {
    case StatusUnstarted, StatusFinished, StatusRunning, StatusSleeping, StatusStopped:
        return true
    }
    return false
}
```

**Naming Conventions**:
- Validation functions: `ValidateX` or `XIsValid` (e.g., `ValidatePage`, `IsValidPage`)
- Custom errors: `ValidationError` with specific types as needed
- Error codes: snake_case (e.g., `page_exceeds_total`, `start_page_after_end_page`)

**Integration with Existing Code**:
- Call validation in handlers before database operations
- Return 400 Bad Request for validation errors
- Continue using existing error response format: `{"error": "message"}`

### 5. Testing Strategy

**Unit Tests (Using `go testing`)**:
- **Validation Function Tests**: Each `ValidateX` function tested with edge cases
  - Valid values (page = total_page, page < total_page)
  - Invalid values (page > total_page)
  - Zero/negative values (boundary conditions)
  - Null pointer handling

- **Error Tests**: Custom error types tested for:
  - Error message format
  - Error code values
  - Error wrapping behavior

**Test Coverage Targets**:
- 100% line coverage for validation functions
- Test all validation rules from PRD Section 5
- Edge case coverage: nil pointers, zero values, empty strings

**Testing Approach**:
- Test files in same package: `validation_test.go`
- Use table-driven tests for validation functions
- Test both valid and invalid inputs
- Verify error returned for invalid cases

### 6. Risks and Considerations

**Blocking Issues**:
- None identified

**Potential Pitfalls**:
1. **Validation vs Database Constraints**: page <= total_page could be enforced at DB level (CHECK constraint), but validation package provides application-level validation with user-friendly messages
2. **Circular Dependencies**: Validation package should not depend on other internal packages to avoid circular imports
3. **Error Handling Consistency**: Must match existing error response format used in handlers

**Trade-offs**:
1. **Validation Location**: Choose validation package over model methods for reusability across CRUD operations
2. **Error Types**: Use simple error interface rather than complex error hierarchy for quick implementation
3. **Status Validation**: Simple string check (not enum) to match existing status constants in `models` package

**Deployment Considerations**:
- No database migrations required (validation is application-level)
- No configuration changes needed
- Backward compatible (adds new validation, doesn't change existing behavior)
- Can be deployed as part of Phase 4 rollout

### Implementation Checklist

1. Create `internal/validation/` directory structure
2. Implement `errors.go` with `ValidationError` type
3. Implement `validate_project.go` with project validation functions:
   - `ValidatePage(page, totalPage)` - page <= total_page
   - `ValidateStatus(status)` - valid status values
4. Implement `validate_log.go` with log validation functions:
   - `ValidateStartEndPage(startPage, endPage)` - start_page <= end_page
5. Write unit tests for all validation functions
6. Run `go fmt` and `go vet` to ensure code quality
7. Update task progress (mark as In Progress → Done after coding)
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
