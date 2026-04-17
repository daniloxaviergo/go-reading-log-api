---
id: RDL-057
title: Update routes
status: To Do
assignee: []
created_date: '2026-04-16 21:06'
updated_date: '2026-04-17 11:00'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The endpoints should be equals responses
- /v1/projects/{project_id}.json
- /v1/projects/{project_id}/logs.json
- /v1/projects.json

remove prefix api
/api/v1/projects.json -> /v1/projects.json

dont remove suffix `.json` only remove prefix `api`

update for new routes: test/compare_responses.sh
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach
- Update router base path from `/api/v1` to `/v1` in `internal/api/v1/routes.go`
- Replace all hardcoded `/api/v1` references with `/v1` in tests and scripts
- Verify no other components depend on the old path structure (e.g., middleware, client integrations)
- Adopt consistent versioning pattern: `/v1/{resource}` without redundant "api" prefix

### 2. Files to Modify
- `internal/api/v1/routes.go`: Change `router.PathPrefix("/api/v1")` → `router.PathPrefix("/v1")`
- `test/compare_responses.sh`: Update all curl commands from `/api/v1/...` to `/v1/...`
- `test/integration/projects_test.go`: Adjust HTTP request paths
- `docs/API.md` (or `QWEN.md` if exists): Update endpoint examples and documentation
- `AGENTS.md`: Refresh API reference section with new paths

### 3. Dependencies
- Must update all tests before merging changes to avoid test failures
- No external dependencies, but requires coordination with client teams if they use the API (though this is a migration project)
- Ensure existing database indexes and query structures remain unchanged

### 4. Code Patterns
- Follow Clean Architecture routing pattern: versioned paths defined in API layer only
- Maintain snake_case naming for resource paths (e.g., `/v1/projects`)
- Use consistent middleware stack without path-specific logic (CORS/recovery/loggers should work transparently)

### 5. Testing Strategy
- Run `go test -v ./...` to verify all unit tests pass with new routes
- Execute `test/compare_responses.sh` against both Go and Rails APIs to validate response consistency
- Add specific integration tests for:
  - `/v1/projects.json`
  - `/v1/projects/{id}.json`
  - `/v1/projects/{id}/logs`
- Verify error responses use consistent status codes (404 for missing resources, 500 for DB errors)
- Use `testing-expert` subagent to run test suite and validate coverage

### 6. Risks and Considerations
- Breaking change for any clients using `/api/v1` paths - must update client code simultaneously
- Verify no hardcoded path references in configuration files (e.g., `.env` or Docker config)
- Check if SSL termination/proxy settings need adjustment in production deployments
- Confirm `make docker-up` and `make run` commands reference correct port/path combinations
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
