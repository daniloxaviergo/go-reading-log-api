---
id: RDL-006
title: '[doc-001 Phase 3] Implement API handlers for projects and logs endpoints'
status: To Do
assignee:
  - catarina
created_date: '2026-04-01 00:58'
updated_date: '2026-04-01 10:58'
labels: []
dependencies: []
references:
  - 'PRD Section: Acceptance Criteria'
  - 'Implementation Checklist: API Layer'
  - 'PRD Section: Key Requirements'
documentation:
  - doc-001
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement handlers in internal/api/v1/handlers/ for all required endpoints: GET /api/v1/projects, GET /api/v1/projects/:id, GET /api/v1/projects/:project_id/logs, and GET /healthz.

Each handler should use repository interfaces for data access and return proper JSON responses matching Rails API behavior.

Implement error handling to return {"error": "<resource> not found"} for missing records and {"error": "Internal server error"} for unexpected errors.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 GET /api/v1/projects returns array of projects with eager-loaded logs ordered by logs descending
- [ ] #2 GET /api/v1/projects/:id returns single project with eager-loaded logs
- [ ] #3 GET /api/v1/projects/:project_id/logs returns first 4 logs for project with project eager-loaded
- [ ] #4 Error responses match Rails API format
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The implementation will create API handlers for projects and logs endpoints following Clean Architecture principles. The handlers will:

1. Use the existing repository interfaces (`ProjectRepository`, `LogRepository`) for data access
2. Return JSON responses matching Rails API behavior
3. Implement proper error handling with standard JSON error responses

**Key decisions:**
- Use Go's standard `net/http` package for HTTP handling (consistent with existing middleware)
- Return DTOs that match Rails API response format (from serializers)
- Use error wrapping for proper error propagation
- Use `json.NewEncoder` for JSON response encoding with `Content-Type` header

**Endpoint behavior:**
- `GET /api/v1/projects`: Returns all projects with eager-loaded logs, ordered by logs descending
- `GET /api/v1/projects/:id`: Returns single project with eager-loaded logs
- `GET /api/v1/projects/:project_id/logs`: Returns first 4 logs for project, with project eager-loaded
- `GET /healthz`: Returns health check response

**Response format mapping (from Rails serializers):**
- ProjectSerializer: `id, name, started_at, progress, total_page, page, status, logs_count, days_unreading, median_day (show only), finished_at (show only)`
- LogSerializer: `id, data, start_page, end_page, note, project`
- Logs in projects: eager-loaded array of LogResponse

### 2. Files to Modify

**New files to create:**
- `internal/api/v1/handlers/projects_handler.go` - Project index/show handlers with repository injection
- `internal/api/v1/handlers/logs_handler.go` - Logs index handler (returns first 4 logs for project)
- `internal/api/v1/handlers/health_handler.go` - Health check handler
- `internal/api/v1/routes.go` - Route registration and router setup

**Dto files to update:**
- `internal/domain/dto/project_response.go` - Add `Logs []*LogResponse` field to match Rails serialization with eager loaded logs

**No other files to modify** - existing middleware, domain models, and repositories are ready

### 3. Dependencies

**Prerequisites (all already implemented):**
- ✅ Domain models (project.go, log.go)
- ✅ DTOs (project_response.go, log_response.go, health_check_response.go)
- ✅ Repository interfaces and implementations
- ✅ Middleware stack (recovery, cors, request_id, logging)
- ✅ Database connection infrastructure

**Required setup:**
1. Database connection must be established
2. Repository implementations must be instantiated with the connection
3. Handlers must be created with repository dependencies injected
4. Routes must be registered and server started

### 4. Code Patterns

**Handler structure:**
```go
type ProjectsHandler struct {
    repo ProjectRepository
}

func NewProjectsHandler(repo ProjectRepository) *ProjectsHandler {
    return &ProjectsHandler{repo: repo}
}

func (h *ProjectsHandler) Index(w http.ResponseWriter, r *http.Request) {
    // 1. Call repository
    // 2. Handle errors
    // 3. Encode JSON response
}
```

**Handler function pattern (from Rails controllers):**
```
1. Call repository method with context
2. Handle errors (404 for not found, 500 for internal)
3. Encode JSON response using json.NewEncoder
4. Set Content-Type: application/json header
```

**Error response format:**
- Missing record: `{"error": "<resource> not found"}`
- Unexpected error: `{"error": "Internal server error"}`

**JSON encoding:**
```go
w.Header().Set("Content-Type", "application/json")
encoder := json.NewEncoder(w)
encoder.Encode(response)
```

**Handling route parameters (with net/http):**
- For `/api/v1/projects/:id`, parse path manually or use `http.ServeMux` with patterns
- Extract project ID from URL path using `strings.TrimPrefix` or regex

### 5. Testing Strategy

**Unit tests** (to be written as separate task per PRD checklist):
- Test each handler method with mocked repository
- Verify JSON response format matches Rails output
- Test error cases (not found returns 404, internal errors return 500)
- Test status codes (200 for success, 404 for not found, 500 for internal error)

**Integration tests** (referenced in PRD acceptance criteria):
- Test against actual database with test data
- Verify `GET /api/v1/projects` returns projects with eager-loaded logs
- Verify logs ordering by `data DESC`
- Verify logs first 4 limit for `/api/v1/projects/:project_id/logs`
- Compare responses to Rails API output

**Test coverage targets:**
- Handler logic: 80%+ coverage
- Integration tests: cover all 4 endpoints

### 6. Risks and Considerations

**Potential issues:**
1. **Missing database fields**: Rails schema doesn't have `progress`, `status`, `logs_count`, `days_unread`, `median_day`, `finished_at` columns in projects table. These may be calculated fields in Rails or require DB schema changes.

2. **Logs ordering**: Rails uses `order('logs.data DESC')` - logs table has `data` as datetime. Need to verify sorting behavior in Go/pgx.

3. **Eager loading implementation**: The repository's `GetWithLogs` method returns a `ProjectResponse` DTO without logs array. This may need updating or a new method to properly eager-load logs.

4. **Route parameter parsing**: With `net/http` stdlib (no chi/mux), route parameters require manual path parsing. For simple patterns like `/api/v1/projects/:id`, can parse manually or use pattern matching.

5. **Response structure mismatch**: Rails' `render json: @project, include: ['logs']` returns project with `logs` array. Need to ensure `ProjectResponse` DTO includes `Logs []*LogResponse` field.

**Trade-offs:**
- Using `net/http` stdlib (no router library) for simplicity, but routing requires manual path parsing
- DTO pattern for response serialization (allows future flexibility)
- Separate handlers per resource for clean separation of concerns
- Repository dependency injection for testability

**Action items to resolve gaps:**
1. Add `Logs []*LogResponse` field to `ProjectResponse` DTO
2. Implement `GetWithLogs` repository method to fetch logs for each project
3. Add route parameter parsing helper or use pattern matching
4. Verify DB schema matches DTO fields (may need Rails migration)
<!-- SECTION:PLAN:END -->
