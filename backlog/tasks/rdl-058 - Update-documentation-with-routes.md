---
id: RDL-058
title: Update documentation with routes
status: To Do
assignee:
  - Thomas
created_date: '2026-04-17 20:43'
updated_date: '2026-04-18 00:30'
labels: []
dependencies: []
ordinal: 2000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
read the file handle/routes.go and update all documentation
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The task requires updating documentation to reflect the current API routes and their implementations. This involves:

- **Analyzing the current routing structure** - Reading `routes.go` to understand all registered endpoints
- **Documenting each route's purpose** - Describing what each endpoint does, HTTP method, path, and response format
- **Mapping handlers to routes** - Connecting each route to its handler implementation
- **Describing request/response formats** - Documenting the JSON structure for each endpoint
- **Adding usage examples** - Providing curl examples or similar for each route

### 2. Files to Modify

**Documentation files to update:**
- `/home/danilo/scripts/github/go-reading-log-api-next/docs/README.go-project.md` (primary documentation)
- `/home/danilo/scripts/github/go-reading-log-api-next/QWEN.md` (Qwen-specific docs)

**Reference files (read-only for understanding):**
- `/home/danilo/scripts/github/go-reading-log-api-next/internal/api/v1/routes.go` - Route definitions
- `/home/danilo/scripts/github/go-reading-log-api-next/internal/api/v1/handlers/` - Handler implementations
- `/home/danilo/scripts/github/go-reading-log-api-next/internal/domain/dto/` - Request/Response DTOs

### 3. Dependencies

**Prerequisites:**
- Understanding of Clean Architecture patterns used in this project
- Knowledge of Go HTTP handlers and Gorilla Mux routing
- Familiarity with JSON:API specification (for response formats)
- Reference to Rails API behavior for compatibility expectations

**No code dependencies blocking this task** - documentation update can proceed independently.

### 4. Code Patterns

**Consistent patterns to document:**
1. **Route naming convention:** `/api/v1/{resource}.json` or `/api/v1/{resource}/{id}.json`
2. **Handler structure:** Each handler file contains a struct with repository dependencies
3. **Response format:** Standard JSON responses matching Rails API format
4. **Error handling:** Consistent error response structure `{"error": "message"}`
5. **Status codes:** 200 for success, 201 for creation, 400 for validation errors, 404 for not found, 500 for server errors

**Documentation conventions:**
- Use Markdown tables for route summaries
- Include curl examples for each endpoint
- Document request body structure with JSON snippets
- Document response structure with JSON snippets
- Note any calculated fields in responses

### 5. Testing Strategy

**For documentation accuracy verification:**
1. **Route coverage check** - Verify all routes in `routes.go` are documented
2. **Handler alignment** - Confirm handler methods match route definitions
3. **Response format validation** - Compare documented response formats with actual DTO implementations
4. **Example testing** - Validate curl examples work against running server

### 6. Risks and Considerations

**Potential issues:**
- Documentation may become stale if not kept in sync with code changes
- Some calculated fields (progress, status, median_day) are computed at runtime - need to document the formulas
- Phase 1 is read-only for logs; Phase 2 will add POST/DELETE - documentation should clarify current limitations
- The `logs_count` and other derived fields are computed from related data - explain the calculation logic

**Design decisions to document:**
- Why `/v1/` prefix instead of `/api/v1/` in some routes
- Why logs endpoint uses `{project_id}` path parameter
- Why only first 4 logs are returned (performance consideration)
- Reason for JSON:API envelope format choice

---
**Implementation Steps:**

1. Read and analyze all route definitions in `routes.go`
2. Review each handler implementation to understand request/response details
3. Check DTOs for exact field names and types
4. Create comprehensive documentation section covering:
   - Overview of API structure
   - Each endpoint with method, path, description
   - Request body examples
   - Response examples
   - Error responses
5. Add curl examples for easy testing
6. Document calculated fields and their formulas
7. Include notes about Phase 1 limitations (read-only logs)
8. Review against existing documentation to avoid duplication

**No code changes required** - this is purely a documentation task.
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
