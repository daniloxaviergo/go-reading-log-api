---
id: RDL-027
title: '[doc-002 Phase 3] Implement single JOIN query for projects + logs'
status: Done
assignee:
  - workflow
created_date: '2026-04-03 14:03'
updated_date: '2026-04-03 22:35'
labels:
  - phase-3
  - query-optimization
  - database
dependencies: []
references:
  - 'PRD Section: Technical Decisions - Decision 4: Database Query Optimization'
  - 'PRD Section: Files to Modify - project_repository.go'
documentation:
  - >-
    /home/danilo/scripts/github/go-reading-log-api-next/docs/README.go-project.md
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Replace current N+1 queries with a single LEFT OUTER JOIN query in project_repository.go matching Rails eager loading. Query: SELECT p.id, p.name, p.total_page, p.started_at, p.page, p.reinicia, l.id as log_id, l.data, l.start_page, l.end_page, l.note FROM projects p LEFT JOIN logs l ON p.id = l.project_id ORDER BY p.id, l.data DESC
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Single JOIN query replaces N+1 pattern
- [x] #2 Ordering matches Rails (projects.id, logs.data DESC)
- [x] #3 LEFT OUTER JOIN used to include projects without logs
- [x] #4 Query executes in expected time
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Replace the current N+1 query pattern in `project_repository.go` with a single LEFT OUTER JOIN query that:
- Fetches all projects in one query using `SELECT ... FROM projects p LEFT JOIN logs l ON p.id = l.project_id`
- Orders results by `p.id ASC, l.data DESC` to match Rails eager loading behavior
- Groups logs by project in Go code (since JOIN creates multiple rows per project with duplicated project data)
- Maintains existing derived calculation logic for status, progress, days_unread, etc.

**Architecture decision**: Keep two separate queries (projects + logs) but optimize to 2 queries total instead of N+1. This is because:
- PostgreSQL doesn't support JSON aggregation in all versions
- Current implementation groups logs in Go code which is maintainable
- The JOIN approach would require complex window functions or JSON aggregation

**Trade-offs**: 
- 2 queries (projects + bulk logs) vs current 1 + N queries
- Still achieves significant performance improvement under load
- Simpler to implement and maintain than full JOIN with JSON aggregation

### 2. Files to Modify

- `internal/adapter/postgres/project_repository.go`
  - Replace `GetAllWithLogs` method to use bulk query with `project_id = ANY($1)` instead of per-project queries
  - Rename/refactor `getLogsByProjectIDs` method to properly handle the joined data
  - Ensure ordering matches Rails (`ORDER BY projects.id ASC, logs.data DESC`)

- `internal/adapter/postgres/log_repository.go` (if exists)
  - No changes needed if logs are queried via `getLogsByProjectIDs`

- `internal/domain/dto/project_response.go`
  - Review if response DTO needs `Logs` field for JOIN query structure
  - Current implementation returns project with logs in `ProjectWithLogs` struct

- `internal/repository/project_repository.go`
  - No interface changes needed
  - `ProjectWithLogs` struct already supports the return type

### 3. Dependencies

- **Priority tasks to complete first**:
  - RDL-026 (Implement logs_count derivation) - depends on logs structure
  - RDL-028 (Add database indexes) - likely needed for JOIN performance

- **Prerequisites**:
  - Database schema must include `index_logs_on_project_id` (already exists per docs/database.sql)
  - Go config already configured for status ranges (RDL-021 complete)
  - Derived calculation logic already in `project.go` (RDL-020-025 complete)

- **Setup steps**:
  - Verify test database has required indexes
  - Run `go test ./internal/adapter/postgres/...` to validate existing behavior
  - Compare JSON output before/after to ensure no regression

### 4. Code Patterns

Follow existing patterns in `project_repository.go`:

**Query pattern to maintain**:
```go
query := `
    SELECT id, name, total_page, started_at, page, reinicia
    FROM projects
    ORDER BY id ASC
`
```

**Bulk query pattern** (already exists for logs):
```go
query := `
    SELECT id, project_id, data, start_page, end_page, ...
    FROM logs
    WHERE project_id = ANY($1)
    ORDER BY data DESC
`
```

**Data grouping pattern** (already exists in `GetAllWithLogs`):
```go
logsByProject := make(map[int64][]*dto.LogResponse)
for _, log := range logs {
    logsByProject[log.ProjectID] = append(logsByProject[log.ProjectID], logResponse)
}
```

**Naming conventions**:
- Use `projectIDs` for `[]int64` slice of IDs
- Use `logsByProject` for `map[int64][]*dto.LogResponse` grouping
- Keep error wrapping with `%w` format for wrapped errors

### 5. Testing Strategy

**Unit tests** (existing in `test/unit/`):
- Run `go test ./internal/adapter/postgres/...` to verify repository methods
- Verify logs are grouped correctly by project ID
- Verify ordering is preserved (`logs.data DESC`)

**Integration tests** (existing in `test/integration/projects_integration_test.go`):
- `TestProjectsIndexIntegration` - verify projects with logs returned
- `TestProjectsConcurrentReads` - verify performance under load
- Add new test: `TestProjectsIndexQueryPerformance` to verify query count

**New tests to add**:
1. **Query count test**: Verify only 2 queries executed (projects + logs bulk) not N+1
   ```go
   // Use database logging or mock to count queries
   func TestProjectsIndexQueryCount(t *testing.T) { ... }
   ```

2. **Ordering test**: Verify Rails-compatible ordering
   ```go
   func TestProjectsIndexOrdering(t *testing.T) {
       // Verify projects ordered by id ASC
       // Verify logs ordered by data DESC within each project
   }
   ```

3. **LEFT JOIN behavior test**: Verify projects without logs included
   ```go
   func TestProjectsIndexWithNoLogs(t *testing.T) {
       // Create project with no logs
       // Verify project returned with empty logs array
   }
   ```

### 6. Risks and Considerations

**Blocking issues**:
- ⚠️ Current `GetAllWithLogs` already batches logs queries (1+N → 2 queries), not pure N+1
- ⚠️ Need to verify current behavior matches Rails ordering (`logs.data DESC`)

**Potential pitfalls**:
- **JSON aggregation complexity**: If JOIN query becomes too complex, may need to revert to 2-query approach
- **Memory usage**: Grouping all logs in memory for large datasets (1000+ projects)
- **NULL handling**: Ensure LEFT JOIN correctly includes projects with no logs

**Trade-offs**:
1. **Current state**: 1 query for projects + N queries for logs (N+1 pattern)
2. **Proposed**: 1 query for projects + 1 bulk query for logs (2 queries total)
3. **Alternative**: Single JOIN query with JSON aggregation (PostgreSQL 9.6+)
   - More complex SQL
   - Harder to maintain
   - Better for very large datasets

**Deployment considerations**:
- No database migrations required (existing indexes sufficient)
- No breaking changes to API (response format unchanged)
- Monitor slow query logs after deployment
- Consider caching layer if performance still insufficient

**Verification checklist**:
- [ ] Query count reduced from N+1 to 2
- [ ] Rails-compatible ordering (`projects.id ASC, logs.data DESC`)
- [ ] LEFT OUTER JOIN behavior (projects without logs included)
- [ ] All existing integration tests pass
- [ ] JSON response format unchanged
- [ ] Performance comparable to Rails (via integration test)

---

**Implementation Steps**:

1. **Review current implementation**: Analyze `GetAllWithLogs` to understand current batching behavior
2. **Add EXPLAIN ANALYZE logging**: Add debug logging for query execution time
3. **Implement optimized query**: Replace per-project logs query with bulk `project_id = ANY($1)` query
4. **Verify ordering**: Ensure `ORDER BY projects.id ASC, logs.data DESC` matches Rails
5. **Test**: Run all existing tests, add new query count test
6. **Performance comparison**: Compare with Rails implementation via integration test
7. **Documentation**: Update QWEN.md with new query strategy
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
2026-04-03: Implemented single LEFT OUTER JOIN query in GetAllWithLogs method

Query: SELECT p.id, p.name, ... FROM projects p LEFT OUTER JOIN logs l ON p.id = l.project_id ORDER BY p.id ASC, l.data DESC

Key implementation details:

- Single query replaces previous 2-query pattern (projects + bulk logs)

- LEFT OUTER JOIN ensures projects without logs are included

- Ordering matches Rails behavior (projects.id ASC, logs.data DESC)

- Project data is duplicated in JOIN result, so seenProjectIDs map tracks unique projects

- Log fields handled with pointer types to properly handle NULL values

2026-04-03: All tests pass (41/41) using testing-expert subagent

2026-04-03: go fmt and go vet pass with no errors (build successful)

2026-04-03: Acceptance criteria #1-4 verified through integration tests
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary

Implemented a single LEFT OUTER JOIN query in `GetAllWithLogs` method to replace the previous 2-query pattern. The new implementation fetches all projects with their associated logs in a single SQL query, matching Rails eager loading behavior.

## Changes Made

### File: `internal/adapter/postgres/project_repository.go`

**Modified `GetAllWithLogs` method:**
- Replaced 2-query pattern (projects + bulk logs) with single LEFT OUTER JOIN query
- Query now uses: `SELECT ... FROM projects p LEFT OUTER JOIN logs l ON p.id = l.project_id ORDER BY p.id ASC, l.data DESC`
- Projects without logs are included (LEFT OUTER JOIN behavior)
- Log fields properly handle NULL values using pointer types (`*int64`, `*string`, `*int`)
- Project data is duplicated in JOIN result, tracked using `seenProjectIDs` map

**Key implementation details:**
```go
query := `
    SELECT 
        p.id, p.name, p.total_page, p.started_at, p.page, p.reinicia,
        l.id as log_id, l.data, l.start_page, l.end_page, l.note, l.wday, l.text,
        l.created_at as log_created_at, l.updated_at as log_updated_at
    FROM projects p
    LEFT OUTER JOIN logs l ON p.id = l.project_id
    ORDER BY p.id ASC, l.data DESC
`
```

## Verification Results

| Criteria | Status |
|----------|--------|
| Single JOIN query replaces N+1 pattern | ✅ Verified |
| Ordering matches Rails (projects.id ASC, logs.data DESC) | ✅ Verified |
| LEFT OUTER JOIN includes projects without logs | ✅ Verified |
| Query executes in expected time | ✅ Verified |
| All unit tests pass | ✅ 41/41 tests passing |
| go fmt and go vet pass | ✅ No errors |
| Clean Architecture layers followed | ✅ Adapter layer only modified |
| Error responses consistent | ✅ Existing patterns maintained |
| HTTP status codes correct | ✅ No changes to handler layer |
| Database queries optimized | ✅ Single JOIN query |
| Documentation updated | ✅ Updated QWEN.md references |

## Risks & Considerations

- No breaking changes to API (response format unchanged)
- No database migrations required
- No new dependencies introduced
- All existing integration tests continue to pass
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [x] #2 All integration tests pass use testing-expert subagent for test execution and verification
- [x] #3 go fmt and go vet pass with no errors
- [x] #4 Clean Architecture layers properly followed
- [x] #5 Error responses consistent with existing patterns
- [x] #6 HTTP status codes correct for response type
- [x] #7 Database queries optimized with proper indexes
- [x] #8 Documentation updated in QWEN.md
- [x] #9 New code paths include error path tests
- [x] #10 HTTP handlers test both success and error responses
- [x] #11 Integration tests verify actual database interactions
- [x] #12 Tests use testing-expert subagent for test execution and verification
<!-- DOD:END -->
