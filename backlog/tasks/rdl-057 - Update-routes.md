---
id: RDL-057
title: Update routes
status: To Do
assignee: []
created_date: '2026-04-16 21:06'
updated_date: '2026-04-17 10:35'
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
Update all API route paths to include `.json` suffix as specified in PRD requirements. This involves modifying route definitions in Gorilla Mux to standardize endpoint naming conventions across all endpoints. The changes will ensure consistency between the code implementation and documentation, aligning with the requirement for `.json` suffixes on all project-related endpoints.

### 2. Files to Modify
- `internal/api/v1/routes.go`: Update route paths to include `.json` suffix (e.g., `/projects/{project_id}/logs` → `/projects/{project_id}/logs.json`)
- `test/compare_responses.sh`: Adjust endpoint URLs in comparison script to match new structure
- `docs/README.go-project.md`: Update API endpoint examples to reflect correct `.json` suffixes
- Any integration test files referencing old endpoint paths (e.g., `test/integration/projects_test.go`)

### 3. Dependencies
- No external dependencies; all changes are internal to the codebase
- Ensure all related tests pass after route updates
- Verify that no existing client applications depend on old URL structure (unlikely in Phase 1 migration)

### 4. Code Patterns
- Maintain consistent naming for all endpoints using `.json` suffixes (e.g., `/projects.json`, `/projects/{id}.json`)
- Follow existing Gorilla Mux routing conventions with explicit path definitions
- Use snake_case for path parameters as per current project standards

### 5. Testing Strategy
- Update unit tests to verify new route paths are correctly registered in router setup
- Modify integration tests to use updated endpoint URLs and validate responses
- Run `testing-expert` subagent to execute all tests with coverage verification
- Validate HTTP status codes and response formats for both success/error cases
- Confirm that `go fmt` and `go vet` pass without errors after changes

### 6. Risks and Considerations
- Potential breaking changes if any internal tools/scripts reference old URLs; verify all test scripts are updated
- Documentation must be synchronized with code changes to prevent confusion during future development
- No public-facing clients exist yet, so impact is limited to internal testing and documentation
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
