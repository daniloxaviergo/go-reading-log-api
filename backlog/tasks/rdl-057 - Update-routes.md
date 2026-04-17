---
id: RDL-057
title: Update routes
status: To Do
assignee: []
created_date: '2026-04-16 21:06'
updated_date: '2026-04-17 00:32'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The endpoints should be equals responses
- /v1/projects/{project_id}.json
- /v1/projects/{project_id}/logs.json
- /v1/projects.json

update the prd to correct url, should be `.json` at the end
- /v1/projects/{project_id}.json
- /v1/projects/{project_id}/logs.json
- /v1/projects.json

update for new routes: test/compare_responses.sh
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Update API endpoint paths to consistently include `.json` suffixes as specified in PRD documentation. The current logs endpoint (`/api/v1/projects/{project_id}/logs`) lacks the required suffix while other endpoints (projects and project-by-id) already have it. This change will align all public-facing API routes with JSON response expectations.

Key steps:
- Modify route definitions to append `.json` to affected paths
- Ensure middleware and handler patterns remain consistent with Clean Architecture principles
- Update test suites to validate new route behavior
- Synchronize documentation with endpoint changes

### 2. Files to Modify

- `internal/api/v1/routes.go`: 
  - Change `/api/v1/projects/{project_id}/logs` → `/api/v1/projects/{project_id}/logs.json`
  - Verify other routes (e.g., `/api/v1/projects.json`, `/api/v1/projects/:id.json`) already have correct suffixes
- `test/compare_responses.sh`:
  - Update all curl commands to use `.json` suffix on relevant endpoints
- All test files referencing old route patterns (e.g., `internal/api/v1/handlers/*_test.go`, integration tests):
  - Replace hardcoded paths with new `.json`-suffixed versions
- `QWEN.md` (API documentation file):
  - Update endpoint examples to reflect correct `.json` suffixes

### 3. Dependencies

- All route changes must be completed before running test suites
- Database schema remains unchanged (no SQL modifications required)
- Existing project initialization code (`cmd/server.go`) doesn't need updates since routes are defined in separate routing layer

### 4. Code Patterns

- Follow existing naming conventions for API endpoints:
  - Versioned paths under `/api/v1/`
  - Resource endpoints always use `.json` suffix for JSON responses
  - Path parameters remain snake_case (e.g., `{project_id}`)
- Maintain current middleware chain structure (Recovery, CORS, RequestID)
- Use Gorilla Mux's strict route matching pattern consistency

### 5. Testing Strategy

- **Unit tests**: Update all handler test cases to use new route paths
  - Verify success/error responses for `/projects/{id}.json` and `/projects/{id}/logs.json`
  - Confirm correct HTTP status codes (200 for success, 404/400 for errors)
- **Integration tests**:
  - Run `testing-expert` subagent to execute full test suite
  - Validate actual HTTP responses against expected JSON structures
  - Verify database query behavior remains unchanged through routing layer
- **Edge cases**:
  - Test invalid project IDs with new route paths
  - Check handling of malformed `.json` requests (e.g., `/projects/1.jsonx`)
- **Coverage**: Ensure all new code paths include error condition tests

### 6. Risks and Considerations

- **Client compatibility**: Existing API consumers will need to update their requests to include `.json` suffixes. This is intentional per PRD requirements.
- **Test maintenance**: All test files must be updated simultaneously to prevent CI failures
- **Documentation sync**: Ensure `QWEN.md` reflects exact endpoint formats to avoid confusion
- **Deployment impact**: No downtime required since this is a pure routing change with backward-compatible path adjustments
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
