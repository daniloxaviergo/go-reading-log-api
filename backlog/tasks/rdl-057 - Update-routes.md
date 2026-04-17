---
id: RDL-057
title: Update routes
status: To Do
assignee: []
created_date: '2026-04-16 21:06'
updated_date: '2026-04-17 12:08'
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

The task requires ensuring **consistent route definitions** across the entire codebase. While `routes.go` already uses `/v1/...` (without `/api` prefix), there are inconsistencies in tests and documentation that reference `/api/v1/...`.

**Key Changes:**
1. **Verify route consistency**: Confirm all route references use `/v1/...` without `/api` prefix
2. **Update test files**: Fix any tests referencing `/api/v1` to use `/v1`
3. **Update documentation**: Align all docs with current route structure
4. **Ensure `.json` suffix**: Verify all project endpoints include `.json` suffix as per PRD

**Architecture Decision:** Follow the Rails API route structure exactly:
- `/v1/projects.json` (GET) - List all projects
- `/v1/projects/{id}.json` (GET) - Get single project  
- `/v1/projects/{project_id}/logs.json` (GET) - Get logs for project

**Why This Approach:**
- Matches Rails API route structure defined in `rails-app/config/routes.rb`
- Ensures consistent response format comparison between Go and Rails APIs
- Maintains clean URL structure without unnecessary path prefixes

---

### 2. Files to Modify

| File | Action | Description |
|------|--------|-------------|
| `internal/api/v1/routes_test.go` | **MODIFY** | Update test routes from `/api/v1/...` to `/v1/...` |
| `test/integration/projects_integration_test.go` | **MODIFY** | Fix route reference on line 30-31 (`/api/v1/projects.json` → `/v1/projects.json`) |
| `test/compare_responses.sh` | **MODIFY** | Update API URL defaults from `.../api/v1` to `.../v1` |
| `docs/endpoint-comparison-report-v1-projects.md` | **MODIFY** | Update documentation to reflect `/v1` routes |
| `docs/rdl-057-route-updates.md` | **CREATE** | Document the route consistency verification and fixes |

---

### 3. Dependencies

**Prerequisites:**
- ✅ RDL-047 completed (handlers updated for routes)
- ✅ RDL-042 completed (JSON:API response structure)
- ✅ Existing route infrastructure in place (`internal/api/v1/routes.go`)

**External Requirements:**
- Go API must be running on port 3000
- Rails API must be running on port 3001
- Both APIs must be accessible for comparison testing

---

### 4. Code Patterns

**Consistent Patterns to Follow:**

```go
// Correct route pattern (already in routes.go):
r.HandleFunc("/v1/projects.json", handler).Methods("GET")
r.HandleFunc("/v1/projects/{id}.json", handler).Methods("GET")
r.HandleFunc("/v1/projects/{project_id}/logs.json", handler).Methods("GET")

// URL configuration for comparison script:
GO_API_URL="http://localhost:3000/v1"
RAILS_API_URL="http://localhost:3001/v1"

// Test request pattern (fixed from /api/v1 to /v1):
req := httptest.NewRequest(http.MethodGet, "/v1/projects.json", nil)
```

**Naming Conventions:**
- Keep existing handler method names unchanged (`Index`, `Show`, `Create`)
- Maintain consistent path parameter naming (`{id}`, `{project_id}`)
- Use uppercase HTTP methods in route definitions
- Always include `.json` suffix for project endpoints

---

### 5. Testing Strategy

**Unit Tests:**
- Verify route registration with correct paths (`/v1/...`)
- Test handler invocation for each endpoint
- Validate error responses for invalid routes

**Integration Tests:**
- Run `test/compare_responses.sh` to verify Go and Rails API responses match
- Test all three endpoints: `/v1/projects.json`, `/v1/projects/{id}.json`, `/v1/projects/{project_id}/logs.json`
- Verify `.json` suffix is preserved in response URLs

**Test Execution Plan:**
```bash
# 1. Start services
make docker-up

# 2. Run unit tests for routes
go test -v ./internal/api/v1/routes_test.go

# 3. Run integration tests
go test -v ./test/integration/...

# 4. Run comparison script
./test/compare_responses.sh

# 5. Verify go fmt and go vet
go fmt ./...
go vet ./...
```

---

### 6. Risks and Considerations

**Blocking Issues:**
- None identified - changes are straightforward path modifications

**Trade-offs:**
- **Breaking Change:** Tests referencing `/api/v1` will need updates
  - *Mitigation:* Systematic search and replace across codebase
- **Test Script Dependencies:** Comparison script relies on correct URL configuration
  - *Mitigation:* Update `test/compare_responses.sh` with proper defaults

**Design Decisions:**
1. **Keep `.json` suffix** - Per PRD specification: "dont remove suffix `.json` only remove prefix `api`"
2. **Match Rails API exactly** - Ensures consistent response comparison
3. **Minimal changes** - Only modify route paths, keep handlers and logic unchanged

**Deployment Considerations:**
- No database migrations required
- No configuration changes needed
- Simple code change with immediate effect
- Easy rollback if issues detected

---

### 7. Acceptance Criteria Verification

| Criteria | Status | Verification Method |
|----------|--------|---------------------|
| #1 All unit tests pass | To Do | Run `go test -v ./internal/api/v1/...` |
| #2 All integration tests pass | To Do | Run `go test -v ./test/...` |
| #3 go fmt and go vet pass | To Do | Run formatting and vetting commands |
| #4 Clean Architecture layers followed | To Do | Verify route definitions in correct layer |
| #5 Error responses consistent | To Do | Verify error handling unchanged |
| #6 HTTP status codes correct | To Do | Verify handlers return correct status codes |
| #7 Database queries optimized | To Do | Verify no query changes needed |
| #8 Documentation updated | To Do | Create `docs/rdl-057-route-updates.md` |
| #9 New code paths include error tests | To Do | Verify existing error tests still pass |
| #10 Handlers test success/error | To Do | Run handler unit tests |
| #11 Integration tests verify DB | To Do | Run integration test suite |
| #12 Tests use testing-expert | To Do | Delegate to testing-expert subagent |

---

### 8. Implementation Checklist

- [ ] Audit all files for `/api/v1` references using `grep`
- [ ] Update `internal/api/v1/routes_test.go` - change `/api/v1` to `/v1`
- [ ] Update `test/integration/projects_integration_test.go` - fix route on lines 30-31
- [ ] Update `test/compare_responses.sh` - change default URLs from `.../api/v1` to `.../v1`
- [ ] Update documentation files referencing `/api/v1`
- [ ] Run unit tests and verify pass
- [ ] Run integration tests and verify pass
- [ ] Run `go fmt` and `go vet`
- [ ] Execute `test/compare_responses.sh` to validate endpoint matching
- [ ] Create documentation for route changes
- [ ] Mark task complete with all DOD criteria met
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
