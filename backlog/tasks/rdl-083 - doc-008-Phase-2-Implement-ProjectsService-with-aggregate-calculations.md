---
id: RDL-083
title: '[doc-008 Phase 2] Implement ProjectsService with aggregate calculations'
status: To Do
assignee:
  - catarina
created_date: '2026-04-21 15:50'
updated_date: '2026-04-21 20:26'
labels:
  - phase-2
  - service
  - aggregate
dependencies: []
references:
  - REQ-DASH-006
  - AC-DASH-002
  - Implementation Checklist Phase 2
documentation:
  - doc-008
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement internal/service/dashboard/projects_service.go querying all projects with eager-loaded logs (first 4, ordered by date DESC), calculating progress_geral, total_pages, and pages aggregates. Order results by progress descending.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 All projects retrieved with eager-loaded logs
- [ ] #2 Log ordering correct (first 4, date DESC)
- [ ] #3 Progress aggregate calculated correctly
- [ ] #4 Total pages and pages aggregates accurate
- [ ] #5 Results ordered by progress descending
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task requires implementing a `ProjectsService` in the dashboard service layer that aggregates project data with calculated fields. The implementation will follow Clean Architecture principles:

**Architecture Decision:** Create a dedicated `ProjectsService` in `internal/service/dashboard/` that:
- Queries all projects using the repository layer
- Eager-loads first 4 logs per project (ordered by date DESC)
- Calculates aggregate metrics: `progress_geral`, `total_pages`, `pages`
- Orders results by progress descending
- Returns structured DTOs for API consumption

**Key Design Choices:**
1. **Service Layer Pattern**: Place business logic in `internal/service/dashboard/` to maintain separation from HTTP handlers and repository concerns
2. **Eager Loading Strategy**: Use a single JOIN query with LIMIT subquery pattern to efficiently fetch logs while maintaining performance
3. **Aggregate Calculation**: Implement calculations in the service layer using Go math operations rather than SQL aggregates for precision and testability
4. **Ordering**: Sort by calculated progress percentage in descending order

**Alternative Considered:** Using database-level aggregation (GROUP BY) - rejected because:
- Requires complex subqueries for eager-loaded logs
- Less flexible for future calculation changes
- Harder to unit test without database

---

### 2. Files to Modify

| Action | File Path | Description |
|--------|-----------|-------------|
| **Create** | `internal/service/dashboard/projects_service.go` | Main service implementation with aggregate calculations |
| **Create** | `internal/service/dashboard/projects_service_test.go` | Unit tests for service logic |
| **Modify** | `internal/api/v1/handlers/dashboards_handler.go` | Wire service into handler (if not already done) |
| **Review** | `internal/repository/postgres/dashboard_repository.go` | Verify existing query patterns match requirements |

---

### 3. Dependencies

**Prerequisites:**
- [x] `internal/service/dashboard/` directory structure created
- [x] `DashboardRepository` interface and PostgreSQL implementation complete (RDL-079)
- [x] Existing log aggregation queries verified working
- [x] Progress calculation logic standardized (RDL-043)

**Blocking Issues:**
None - this is a new service implementation building on existing repository infrastructure.

**Setup Steps:**
1. Ensure `internal/service/dashboard/` directory exists
2. Verify `dashboard_repository.go` has `GetAllWithLogs` method
3. Confirm test database is available for integration tests

---

### 4. Code Patterns

**Follow Existing Patterns:**

```go
// Pattern from day_service.go - similar structure expected
type ProjectsService struct {
    repo DashboardRepository
}

func NewProjectsService(repo DashboardRepository) *ProjectsService {
    return &ProjectsService{repo: repo}
}

// Method signature pattern
func (s *ProjectsService) GetAll(ctx context.Context) ([]*dto.ProjectResponse, error) {
    // Implementation
}
```

**Naming Conventions:**
- Service file: `projects_service.go` (lowercase, underscore separator)
- Test file: `projects_service_test.go`
- Struct name: `ProjectsService`
- Constructor: `NewProjectsService(...)`
- Method names: PascalCase (e.g., `GetAll`, `CalculateAggregates`)

**Integration Pattern:**
```go
// In handler
projectService := service.NewProjectsService(dashboardRepo)
projects, err := projectService.GetAll(r.Context())
```

---

### 5. Testing Strategy

**Unit Tests (`projects_service_test.go`):**

| Test Case | Description |
|-----------|-------------|
| `TestProjectsService_GetAll_Success` | Verify all projects retrieved with logs and aggregates |
| `TestProjectsService_GetAll_EmptyDatabase` | Handle empty projects table gracefully |
| `TestProjectsService_CalculateAggregates` | Verify aggregate calculations (progress_geral, total_pages, pages) |
| `TestProjectsService_SortingByProgress` | Confirm results ordered by progress DESC |
| `TestProjectsService_LogLimit` | Verify only first 4 logs per project loaded |

**Integration Tests:**
- Use existing `test.TestHelper` for database setup/teardown
- Create fixture data with known progress values
- Verify actual database queries match expectations
- Test edge cases: zero pages, null started_at, missing logs

**Coverage Targets:**
- Unit tests: 100% coverage of service methods
- Integration tests: Cover all acceptance criteria scenarios

---

### 6. Risks and Considerations

| Risk | Mitigation |
|------|------------|
| **Performance**: N+1 queries for logs | Use eager loading with single JOIN query pattern from existing code |
| **Progress calculation precision** | Follow existing `CalculateProgress` logic from project model |
| **Empty log arrays** | Handle gracefully with empty slice instead of null |
| **Sorting null progress** | Define explicit behavior (sort to end or handle as 0) |
| **Memory usage for large datasets** | Consider pagination if projects > 1000 |

**Trade-offs:**
1. **Single query vs multiple queries**: Chose single JOIN with LIMIT for efficiency, though subquery may be needed for precise "first 4 logs" logic
2. **Calculation location**: In-service Go calculations rather than SQL for testability and flexibility

**Acceptance Criteria Verification:**
- [ ] #1: Verify `logs` array populated with exactly 4 entries per project (or fewer if less exist)
- [ ] #2: Confirm logs ordered by `data DESC` within each project
- [ ] #3: Validate `progress_geral` matches `(pages / total_pages) * 100`
- [ ] #4: Ensure `total_pages` and `pages` match database values exactly
- [ ] #5: Confirm results sorted by progress descending (highest first)
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass
- [ ] #2 All integration tests pass execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md
- [ ] #8 New code paths include error path tests
- [ ] #9 HTTP handlers test both success and error responses
- [ ] #10 Integration tests verify actual database interactions
<!-- DOD:END -->
