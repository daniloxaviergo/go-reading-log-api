---
id: RDL-045
title: Fix order by projects.json
status: To Do
assignee:
  - thomas
created_date: '2026-04-14 09:53'
updated_date: '2026-04-14 10:02'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add correct order query database, look rails-app
execute test/compare_responses.sh to endpoint v1/projects.json

The Go API and Rails API return completely different projects from the same database query. The Go API returns project ID 1 ("Filocalia"), while the Rails API returns project ID 450 ("História da Igreja VIII.1").

**Go Response (First Project):**
```json
{
  "days_unreading": 3354,
  "finished_at": null,
  "id": 1,
  "logs_count": 50,
  "name": "Filocalia",
  "page": 655,
  "progress": null,
  "started_at": "2017-02-04T00:00:00Z",
  "status": "stopped",
  "total_page": 1267
}
```

**Rails Response (First Project):**
```json
{
  "days_unreading": 10,
  "finished_at": null,
  "id": 450,
  "logs_count": 38,
  "name": "História da Igreja VIII.1",
  "page": 691,
  "progress": 100.0,
  "started_at": "2026-02-19",
  "status": "finished",
  "total_page": 691
}
```
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The issue is that the Go API and Rails API return projects in different orders when querying `/api/v1/projects`. 

**Root Cause Analysis:**
- **Rails API**: Orders projects by `logs.data DESC` (most recent log date first) using `Project.eager_load(:logs).where(id: projects.map(&:id)).order('logs.data DESC')`
- **Go API**: Currently orders by `p.id ASC` in the JOIN query (`ORDER BY p.id ASC, l.data DESC`)

This causes the first project in the list to be different:
- Rails returns: Project ID 450 ("História da Igreja VIII.1") - has recent logs
- Go returns: Project ID 1 ("Filocalia") - oldest project, no recent logs

**Solution:**
Modify the `GetAllWithLogs` method in `ProjectRepositoryImpl` to order results by the most recent log date (`logs.data DESC`) instead of project ID. This matches the Rails API behavior.

**Implementation Strategy:**
1. Update the SQL query in `GetAllWithLogs` to order by `l.data DESC` (with NULLS LAST for projects without logs)
2. Update the `GetWithLogs` method to also use consistent ordering
3. Ensure the ordering is applied to the joined result set
4. Update tests to verify the ordering behavior

### 2. Files to Modify

| File | Changes | Reason |
|------|---------|--------|
| `internal/adapter/postgres/project_repository.go` | Modify `GetAllWithLogs` query to order by `l.data DESC` | Fix ordering to match Rails API |
| `internal/adapter/postgres/project_repository.go` | Modify `GetWithLogs` to use consistent ordering | Ensure single project lookup is also ordered correctly |
| `internal/api/v1/handlers/projects_handler.go` | No changes needed | Handler already uses repository correctly |
| `test/compare_responses.sh` | No changes needed | Test script will validate the fix |

### 3. Dependencies

| Dependency | Status | Notes |
|------------|--------|-------|
| Database indexes | Already exists | `index_logs_on_project_id_and_data_desc` already created in schema |
| Go version | 1.25.7 | No changes needed |
| pgx/v5 | Already used | No changes needed |

### 4. Code Patterns

**Current Query (Incorrect Ordering):**
```go
query := `
    SELECT 
        p.id, p.name, p.total_page, p.started_at, p.page, p.reinicia,
        l.id as log_id, l.data::text as data_text, l.start_page, l.end_page, l.note, l.wday, l.text,
        l.created_at as log_created_at, l.updated_at as log_updated_at
    FROM projects p
    LEFT OUTER JOIN logs l ON p.id = l.project_id
    ORDER BY p.id ASC, l.data DESC  -- ❌ Orders by project ID first
`
```

**Fixed Query (Correct Ordering):**
```go
query := `
    SELECT 
        p.id, p.name, p.total_page, p.started_at, p.page, p.reinicia,
        l.id as log_id, l.data::text as data_text, l.start_page, l.end_page, l.note, l.wday, l.text,
        l.created_at as log_created_at, l.updated_at as log_updated_at
    FROM projects p
    LEFT OUTER JOIN logs l ON p.id = l.project_id
    ORDER BY l.data DESC NULLS LAST  -- ✅ Orders by log date first (matches Rails)
`
```

**Key Considerations:**
1. `NULLS LAST` ensures projects without logs still appear in results
2. The ordering must be applied to the joined result set before grouping
3. Single project queries (`GetWithLogs`) should also use log-based ordering for consistency

### 5. Testing Strategy

**Unit Tests:**
- Verify `GetAllWithLogs` returns projects in correct order
- Verify `GetWithLogs` returns single project with correct log ordering
- Test edge cases: projects without logs, single log, multiple logs

**Integration Tests:**
- Compare full API response with Rails API using `test/compare_responses.sh`
- Verify first project in response matches between Go and Rails APIs
- Verify all projects are returned in same order

**Test Commands:**
```bash
# Run unit tests
go test -v ./internal/adapter/postgres/...

# Run integration tests
go test -v ./internal/api/v1/...

# Compare API responses
./test/compare_responses.sh
```

### 6. Risks and Considerations

| Risk | Mitigation |
|------|------------|
| **Performance**: Sorting by `logs.data` on large datasets | Index `index_logs_on_project_id_and_data_desc` already exists, ensuring efficient sorting |
| **Projects without logs**: `NULLS LAST` handles this but verify behavior | Test with projects that have no logs to ensure they appear at end |
| **Breaking change**: Client code may rely on current ordering | Document the fix in QWEN.md; this aligns with Rails API which is the source of truth |
| **Test failures**: Existing tests may fail due to ordering change | Update test assertions to expect new ordering; use `compare_responses.sh` as validation |

**Blocking Issues:**
- None identified. This is a bug fix to align with existing Rails API behavior.

**Trade-offs:**
- The current implementation in `GetAllWithLogs` groups logs by project after querying, which works correctly with the new ordering
- No schema changes required
- No migration needed

---

## Final Acceptance Criteria Checklist

After implementation, verify:

- [ ] **#1** All unit tests pass (use testing-expert subagent for test execution and verification)
- [ ] **#2** All integration tests pass (use testing-expert subagent for test execution and verification)
- [ ] **#3** `go fmt` and `go vet` pass with no errors
- [ ] **#4** Clean Architecture layers properly followed
- [ ] **#5** Error responses consistent with existing patterns
- [ ] **#6** HTTP status codes correct for response type
- [ ] **#7** Database queries optimized with proper indexes
- [ ] **#8** Documentation updated in QWEN.md
- [ ] **#9** New code paths include error path tests
- [ ] **#10** HTTP handlers test both success and error responses
- [ ] **#11** Integration tests verify actual database interactions
- [ ] **#12** Tests use testing-expert subagent for test execution and verification

**Manual Verification Step:**
```bash
# Run comparison script to verify fix
./test/compare_responses.sh

# Expected: All tests should pass, first project should match between Go and Rails APIs
```
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
