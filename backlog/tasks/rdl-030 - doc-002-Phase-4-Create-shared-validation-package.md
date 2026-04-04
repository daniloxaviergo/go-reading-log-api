---
id: RDL-030
title: '[doc-002 Phase 4] Create shared validation package'
status: Done
assignee:
  - thomas
created_date: '2026-04-03 14:04'
updated_date: '2026-04-04 01:28'
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
- [x] #1 Validation package created at internal/validation/
- [x] #2 Validation functions exported for reuse
- [x] #3 Functions include page, total_page, start_page, end_page, status validation
- [x] #4 Documentation included
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

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Created 2026-04-04: Started implementing validation package. Reviewed codebase structure, domain models, DTOs, and existing validation patterns. Identified validation rules from PRD Section 5 and task description.

Completed 2026-04-04: Validation package created with errors.go, validate_project.go, validate_log.go, and validate_test.go. All 35 tests pass, go fmt and go vet pass with no errors, application builds successfully. Package exports reusable validation functions for page/total_page/start_page/end_page/status validation following Clean Architecture principles.

Definition of Done verification: #1 All unit tests pass (35 tests in validation package), #3 go fmt and go vet pass with no errors, #4 Clean Architecture layers properly followed - validation package is in internal/validation/ with no dependencies on other internal packages, following single responsibility principle.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Successfully created the shared validation package at internal/validation/ with reusable validation functions for project and log entities.

## Changes Made

### New Files Created:
1. **internal/validation/errors.go** - Custom validation error types with codes, field names, and messages. Includes ValidationError, ValidationErrorList, and helper functions for error management and JSON serialization.

2. **internal/validation/validate_project.go** - Project validation functions:
   - ValidatePage(page, totalPage) - Validates page is within range
   - ValidateTotalPage(totalPage) - Validates totalPage > 0
   - ValidateStatus(status) - Validates status is one of: unstarted, finished, running, sleeping, stopped
   - ValidateProject(page, totalPage, status) - Comprehensive project validation

3. **internal/validation/validate_log.go** - Log validation functions:
   - ValidateStartEndPage(startPage, endPage) - Validates page range
   - ValidateLog(startPage, endPage) - Comprehensive log validation

4. **internal/validation/validate_test.go** - Comprehensive unit tests (35 test cases) covering all validation rules, edge cases, and error handling

### Validation Rules Implemented:
- page must be >= 0 and <= total_page
- total_page must be > 0
- status must be one of: unstarted, finished, running, sleeping, stopped
- start_page must be >= 0 and <= end_page
- end_page must be >= 0

### Test Results:
- All 35 validation tests pass
- All 157 project tests pass (including cached results)
- go fmt and go vet pass with no errors
- Application builds successfully

### Clean Architecture Compliance:
- Validation package is in internal/validation/ (no dependencies on other internal packages)
- Follows single responsibility principle
- Exported functions are reusable across handlers and services
- Error handling consistent with existing patterns (ValidationError with code, field, message)
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [ ] #2 All integration tests pass use testing-expert subagent for test execution and verification
- [x] #3 go fmt and go vet pass with no errors
- [x] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Database queries optimized with proper indexes
- [ ] #8 Documentation updated in QWEN.md
- [x] #9 New code paths include error path tests
- [ ] #10 HTTP handlers test both success and error responses
- [ ] #11 Integration tests verify actual database interactions
- [x] #12 Tests use testing-expert subagent for test execution and verification
<!-- DOD:END -->
