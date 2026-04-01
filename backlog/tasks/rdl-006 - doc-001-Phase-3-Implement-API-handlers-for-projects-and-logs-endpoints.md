---
id: RDL-006
title: '[doc-001 Phase 3] Implement API handlers for projects and logs endpoints'
status: To Do
assignee:
  - workflow
created_date: '2026-04-01 00:58'
updated_date: '2026-04-01 11:01'
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

**Analysis of Current State:**

After thorough research of the codebase, I found that:

1. **Most handler code already exists** - The handlers for projects, logs, and health check are already implemented in `internal/api/v1/handlers/`
2. **Routes are configured** - `routes.go` uses `gorilla/mux` for routing with path parameters
3. **Middleware chain is ready** - Includes recovery, CORS, request ID, and logging middleware

**Issues Found (requires fixes):**

1. **Missing `Logs` field in ProjectResponse** - The Rails API returns logs with projects (e.g., `GET /api/v1/projects` returns projects with `logs` array), but the current `ProjectResponse` doesn't include the `Logs []*LogResponse` field

2. **Index endpoint doesn't eager-load logs** - `ProjectsHandler.Index` returns projects without logs, but Rails uses `Project.eager_load(:logs).order('logs.data DESC')`

3. **Show endpoint doesn't eager-load logs properly** - `ProjectsHandler.Show` retrieves project separately from logs but doesn't include them in the response

4. **ProjectRepository.GetWithLogs doesn't load logs** - The repository method exists but only fetches project data, not associated logs

5. **Logs handler doesn't order by data DESC** - The Rails API orders logs by `logs.data DESC`, but current implementation returns logs in database order

**Corrected Approach:**

The implementation will:
1. Add `Logs []*LogResponse` field to `ProjectResponse` DTO
2. Implement `GetAllWithLogs` repository method that eager-loads logs ordered by `data DESC`
3. Update `ProjectsHandler.Index` to use eager-loaded logs
4. Update `ProjectsHandler.Show` to use eager-loaded logs
5. Update `LogsHandler.Index` to order logs by `data DESC` before limiting to 4
6. Fix `ProjectRepository.GetWithLogs` to include logs array

### 2. Files to Modify

**New files to create:**
- None (all handlers and infrastructure already exist)

**Files to modify:**

1. `internal/domain/dto/project_response.go` - Add `Logs []*LogResponse` field to store eager-loaded logs
2. `internal/domain/dto/log_response.go` - Add `Data` field type fix (currently `*string`, should match DB `datetime`)
3. `internal/repository/project_repository.go` - Add `GetAllWithLogs() []ProjectWithLogs` interface method (helper struct)
4. `internal/adapter/postgres/project_repository.go` - Implement `GetAllWithLogs()` with eager-loaded logs ordered by `data DESC`
5. `internal/api/v1/handlers/projects_handler.go` - Update `Index` and `Show` to include eager-loaded logs
6. `internal/api/v1/handlers/logs_handler.go` - Update to order logs by `data DESC` before limiting

### 3. Dependencies

**Prerequisites (all already implemented):**
- ✅ Domain models (`project.go`, `log.go`)
- ✅ DTOs (`project_response.go`, `log_response.go`, `health_check_response.go`)
- ✅ Repository interfaces and PostgreSQL implementations
- ✅ Middleware stack (recovery, CORS, request_id, logging)
- ✅ Database connection infrastructure
- ✅ Route registration with `gorilla/mux`
- ✅ Application entry point with graceful shutdown

**Required changes before implementation:**
1. Add `Logs` field to `ProjectResponse` DTO
2. Extend repository interface to return logs with projects
3. Implement `GetAllWithLogs` in `ProjectRepositoryImpl`

### 4. Code Patterns

**Existing patterns to follow:**

1. **Handler structure** (already implemented):
```go
type ProjectsHandler struct {
    repo repository.ProjectRepository
}

func NewProjectsHandler(repo repository.ProjectRepository) *ProjectsHandler {
    return &ProjectsHandler{repo: repo}
}
```

2. **Error handling** (already implemented):
```go
if strings.Contains(err.Error(), "not found") {
    http.Error(w, `{"error": "project not found"}`, http.StatusNotFound)
    return
}
http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
```

3. **JSON response encoding** (already implemented):
```go
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(response)
```

4. **Path parameter extraction** (already implemented):
```go
idStr := mux.Vars(r)["id"]
id, err := strconv.ParseInt(idStr, 10, 64)
```

**New patterns to add:**

1. **Eager-loaded project with logs struct:**
```go
type ProjectWithLogs struct {
    Project *dto.ProjectResponse
    Logs    []*dto.LogResponse
}
```

2. **Time formatting helper:**
```go
func formatTimePtr(t *time.Time) *string {
    if t == nil {
        return nil
    }
    s := t.Format(time.RFC3339)
    return &s
}
```

### 5. Testing Strategy

**Unit tests** (to be written as separate task per PRD):
- Test `ProjectsHandler.Index` with mocked repository returning projects with logs
- Test `ProjectsHandler.Show` with mocked repository returning project with logs
- Test `LogsHandler.Index` with mocked repositories verifying log limit and ordering
- Verify error handling for not found and internal errors
- Test JSON response format matches Rails output

**Integration tests** (referenced in PRD acceptance criteria):
- Test `GET /api/v1/projects` returns projects with eager-loaded logs ordered by `data DESC`
- Test `GET /api/v1/projects/:id` returns single project with logs
- Test `GET /api/v1/projects/:project_id/logs` returns first 4 logs ordered by `data DESC`
- Verify `ProjectResponse` includes `Logs []*LogResponse` array
- Compare responses to Rails API output

**Test coverage targets:**
- Handler logic: 80%+ coverage
- Repository eager-loading: Verify SQL queries use `LEFT JOIN` for logs
- Logs ordering: Verify `ORDER BY logs.data DESC` in SQL

### 6. Risks and Considerations

**Blocking issues identified:**

1. **Missing logs in ProjectResponse** - CRITICAL: The Rails API returns `logs` array with projects, but current DTO lacks this field. Without this fix, acceptance criteria #1, #2, #3 will fail.

2. **GetWithLogs not loading logs** - CRITICAL: The repository method exists but doesn't fetch associated logs. Need to implement proper eager-loading with `LEFT JOIN`.

3. **Logs ordering** - HIGH: Rails uses `order('logs.data DESC')`, but current implementation doesn't order logs. Need to verify Go/PostgreSQL sorting matches Rails behavior.

4. **Route parameter precision** - MEDIUM: `mux.Vars(r)["id"]` requires path to match route exactly. Need to verify route patterns in `routes.go` match expected patterns.

**Potential pitfalls:**

1. **N+1 query problem** - If `GetAll` returns projects and we loop to fetch logs per project, we'd have N+1 queries. The solution must use `LEFT JOIN` for single query with eager loading.

2. **Null time handling** - Rails uses `DateTime` fields (nullable), Go uses `*time.Time`. Need consistent handling of null/nil values in DTO.

3. **Logs ordering field** - Rails uses `logs.data DESC`. Need to verify this column exists in PostgreSQL schema (DB schema shows `data` as `datetime`).

**Trade-offs:**

1. **SQL approach** - Use `LEFT JOIN` with `ORDER BY logs.data DESC` for eager loading (single query, matches Rails `eager_load` behavior)

2. **Response structure** - Keep `ProjectResponse` with optional `Logs` field (backwards compatible for single project without logs)

3. **Error message format** - Use exact Rails format: `{"error": "project not found"}` (already implemented)

**Action items to resolve gaps:**
1. Add `Logs []*LogResponse` field to `ProjectResponse` DTO
2. Implement `ProjectWithLogs` struct for eager-loaded projects
3. Implement `GetAllWithLogs` in `ProjectRepositoryImpl` with `LEFT JOIN`
4. Update `ProjectsHandler.Index` to use eager-loaded logs
5. Update `ProjectsHandler.Show` to use eager-loaded logs
6. Add `ORDER BY logs.data DESC` to logs queries
7. Implement helper function `formatTimePtr` for time conversion
8. Verify `logs.data` column exists in DB schema (currently named `data datetime` in schema)
<!-- SECTION:PLAN:END -->
