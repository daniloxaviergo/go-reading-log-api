---
id: RDL-006
title: '[doc-001 Phase 3] Implement API handlers for projects and logs endpoints'
status: To Do
assignee:
  - catarina
created_date: '2026-04-01 00:58'
updated_date: '2026-04-01 02:31'
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
2. Return JSON responses matching Rails API behavior (via serializers)
3. Implement proper error handling with standard JSON error responses

**Key decisions:**
- Use Go's standard `net/http` package for HTTP handling (consistent with existing middleware)
- Follow the same pattern as middleware: handler functions that wrap with context, error handling, and JSON encoding
- Use `github.com/google/uuid` for request tracing (already in go.mod)

**Endpoint behavior:**
- `GET /api/v1/projects`: Returns all projects with eager-loaded logs, ordered by logs descending
- `GET /api/v1/projects/:id`: Returns single project with eager-loaded logs
- `GET /api/v1/projects/:project_id/logs`: Returns first 4 logs for project, with project eager-loaded
- `GET /healthz`: Returns health check response

### 2. Files to Modify

**New files to create:**
- `internal/api/v1/handlers/projects_handler.go` - Project index/show handlers
- `internal/api/v1/handlers/logs_handler.go` - Logs index handler
- `internal/api/v1/handlers/health_handler.go` - Health check handler
- `internal/api/v1/routes.go` - Route registration and router setup

**No files to modify** - existing infrastructure (middleware, domain models, repositories) is ready

### 3. Dependencies

**Prerequisites (all already implemented):**
- ✅ `internal/domain/models/project.go` - Project model
- ✅ `internal/domain/models/log.go` - Log model
- ✅ `internal/domain/dto/project_response.go` - Project response DTO
- ✅ `internal/domain/dto/log_response.go` - Log response DTO
- ✅ `internal/domain/dto/health_check_response.go` - Health check response DTO
- ✅ `internal/repository/project_repository.go` - Project repository interface
- ✅ `internal/repository/log_repository.go` - Log repository interface
- ✅ `internal/adapter/postgres/project_repository.go` - Project repository implementation
- ✅ `internal/adapter/postgres/log_repository.go` - Log repository implementation
- ✅ `internal/api/v1/middleware/*.go` - Middleware stack (recovery, cors, request_id, logging)

**Required setup:**
- Database connection must be established before handlers are created
- Repository implementations must be injected into handlers (dependency injection pattern)

### 4. Code Patterns

**Handler function pattern (from Rails controllers):**
```
1. Extract parameters from context/route
2. Call repository method
3. Map domain model to response DTO if needed
4. Write JSON response with appropriate status code
5. Handle errors with consistent JSON format
```

**Error response format:**
- Missing record: `{"error": "<resource> not found"}`
- Unexpected error: `{"error": "Internal server error"}`

**JSON encoding:**
- Use `encoding/json` package for response serialization
- Set `Content-Type: application/json` header
- Use `json.NewEncoder(w).Encode()` for responses

**Handler structure:**
```go
type ProjectsHandler struct {
    repo ProjectRepository
}

func NewProjectsHandler(repo ProjectRepository) *ProjectsHandler {
    return &ProjectsHandler{repo: repo}
}

func (h *ProjectsHandler) Index(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    // ... implementation
}
```

### 5. Testing Strategy

**Unit tests** (to be written as separate task):
- Test each handler method with mocked repository
- Verify JSON response format
- Test error cases (not found, internal errors)
- Test status codes (200, 404, 500)

**Integration tests** (referenced in PRD):
- Test against actual database with test data
- Verify endpoint behavior matches Rails app
- Test eager loading (logs are returned with projects)

**Test coverage targets:**
- Handler logic: 80%+ coverage
- Integration tests: cover all endpoints

### 6. Risks and Considerations

**Potential issues:**
1. **Missing database fields**: Rails schema doesn't have `progress`, `status`, `logs_count`, `days_unread`, `median_day`, `finished_at` columns in projects table. These may be calculated fields or require DB changes.

2. **Logging order behavior**: Rails uses `order('logs.data DESC')` but logs table has `data` as datetime. Need to verify sorting behavior in Go/pgx.

3. **Eager loading**: PostgreSQL repository uses separate queries for logs. May need eager loading optimization if performance issues arise.

4. **No route parameters yet**: The existing `cmd/server.go` doesn't have route parameter extraction. May need to add path routing (e.g., using `net/http` with manual path parsing or switch to chi/mux for cleaner routing).

**Trade-offs:**
- Keeping `net/http` stdlib (no chi/mux router) for simplicity, but routing requires manual path parsing
- Using DTOs instead of direct model serialization for future flexibility
- Separate handlers per resource (projects, logs, health) for clean separation
<!-- SECTION:PLAN:END -->
