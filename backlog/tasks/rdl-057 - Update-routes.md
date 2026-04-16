---
id: RDL-057
title: Update routes
status: To Do
assignee:
  - catarina
created_date: '2026-04-16 21:06'
updated_date: '2026-04-16 21:09'
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
- Update all API endpoint paths in `internal/api/v1/routes.go` to include `.json` suffix (e.g., `/api/v1/projects` → `/api/v1/projects.json`)
- Maintain existing HTTP method handling for each route
- Ensure path variables (e.g., `{id}`) are preserved before the `.json` suffix
- Update test scripts and integration tests to reference new endpoint formats

### 2. Files to Modify
- `internal/api/v1/routes.go` (primary routing configuration)
- `test/compare_responses.sh` (API comparison script)
- Any test files referencing old endpoint paths (e.g., `test/integration_test.go`)
- `AGENTS.md` or documentation files describing API endpoints

### 3. Dependencies
- Must update all references to routes before testing can pass
- Existing unit tests will need path updates to validate new structure
- Integration tests depend on updated test scripts for validation

### 4. Code Patterns
- Follow existing Gorilla Mux routing conventions (suffix-based paths)
- Maintain snake_case naming consistency in route definitions
- Preserve current handler-to-route mappings exactly, only modifying path strings
- Use consistent `.json` suffix placement across all endpoints

### 5. Testing Strategy
- **Unit tests**: Verify router configuration with `go test -v ./internal/api/v1/...`
- **Integration tests**: Update test clients to use new paths and validate responses
- **Test script**: Modify `test/compare_responses.sh` to query `.json` endpoints
- **Edge cases**: Test malformed JSON suffixes (e.g., `/projects.json/extra`) should return 404
- **Coverage**: Ensure all handler methods have path-specific test coverage

### 6. Risks and Considerations
- **Path conflicts**: Verify no route collisions between `.json` and non-`.json` versions
- **Middleware compatibility**: Check if any middleware parses full paths (unlikely but verify)
- **Documentation sync**: Update all external documentation to match new endpoint format
- **Legacy clients**: No known legacy clients in this Phase 1 migration project
- **CI/CD pipeline**: Ensure test scripts and build validation steps are updated for new paths
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
