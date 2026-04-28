---
id: RDL-128
title: >-
  [doc-011 Phase 1] Add GetRunningProjectsWithLogs repository query with SQL
  JOIN
status: To Do
assignee:
  - thomas
created_date: '2026-04-28 11:16'
updated_date: '2026-04-28 13:05'
labels:
  - feature
  - backend
  - phase-1
dependencies: []
documentation:
  - doc-011
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update internal/adapter/postgres/dashboard_repository.go with GetRunningProjectsWithLogs() method implementing SQL query with JOIN to eager-load first 4 logs per project (ordered by date DESC). Add progress ordering in SQL using CASE statement and handle NULL values with COALESCE.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 SQL query joins projects with logs table
- [ ] #2 Logs limited to first 4 per project ordered by data DESC
- [ ] #3 Progress ordering implemented via SQL CASE statement
- [ ] #4 NULL values handled with COALESCE
- [ ] #5 Query returns all required project and log fields
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This implementation will create an efficient single-query repository method that uses SQL JOINs and window functions to eager-load the first 4 logs per project. The approach follows Clean Architecture principles and matches the existing codebase patterns.

**Technical Strategy:**
- Implement `GetRunningProjectsWithLogs()` in `internal/adapter/postgres/dashboard_repository.go`
- Use a CTE (Common Table Expression) with `ROW_NUMBER()` window function to limit logs to first 4 per project
- Order projects by progress using SQL `CASE` statement for NULL handling
- Use `COALESCE` for all NULL value handling to ensure consistent results
- Return `[]*dto.ProjectWithLogs` matching the existing DTO structure

**Why This Approach:**
- Single query eliminates N+1 query problem (current implementation calls `GetProjectLogs` for each project)
- Window functions are PostgreSQL-native and efficient
- Matches existing patterns in `GetProjectsWithLogs` method
- Service layer already expects this signature and handles filtering/sorting

**Architecture Decisions:**
- Keep repository focused on data access (SQL queries)
- Service layer handles business logic (status filtering, progress calculation)
- Handler layer handles HTTP concerns (response formatting, error handling)
- Follow existing 15-second timeout pattern for dashboard queries

---

### 2. Files to Modify

#### Files to Modify

| File | Changes | Priority |
|------|---------|----------|
| `internal/adapter/postgres/dashboard_repository.go` | Replace current `GetRunningProjectsWithLogs` implementation with single SQL JOIN query using CTE and window functions | P1 |
| `internal/service/dashboard/projects_service.go` | No changes needed - service already calls repository method correctly | - |
| `internal/api/v1/handlers/dashboard_handler.go` | No changes needed - handler already uses service layer correctly | - |

#### Files to Review (No Changes Required)

| File | Purpose |
|------|---------|
| `internal/repository/dashboard_repository.go` | Interface definition - already has `GetRunningProjectsWithLogs` method signature |
| `internal/domain/dto/dashboard_response.go` | DTOs already defined (`ProjectWithLogs`, `ProjectAggregateResponse`, `LogEntry`) |
| `internal/service/dashboard/projects_service.go` | Service layer already implements filtering and sorting logic |

---

### 3. Dependencies

**Prerequisites:**
- DashboardRepository interface already defined in `internal/repository/dashboard_repository.go`
- DTOs (`ProjectWithLogs`, `ProjectAggregateResponse`, `LogEntry`) already exist in `internal/domain/dto/dashboard_response.go`
- Service layer (`ProjectsService`) already implemented and calling repository method
- PostgreSQL with pgx/v5 driver (already in use)

**Existing Infrastructure:**
- `GetProjectsWithLogs` method exists as reference implementation
- `GetProjectLogs` method exists (will be replaced by single query)
- Context timeout pattern established (15 seconds for dashboard queries)
- Error handling patterns established (fmt.Errorf with wrapping)

**No External Dependencies Required:**
- All required packages already imported
- No new database tables or schema changes needed
- No new configuration required

---

### 4. Code Patterns

**SQL Query Pattern:**
```go
query := `
    WITH log_ranked AS (
        SELECT 
            l.id,
            l.project_id,
            l.data,
            l.start_page,
            l.end_page,
            l.note,
            ROW_NUMBER() OVER (PARTITION BY l.project_id ORDER BY l.data DESC) as rn
        FROM logs l
    )
    SELECT 
        p.id as project_id,
        p.name as project_name,
        p.total_page as project_total_page,
        p.page as project_page,
        COALESCE(SUM(CASE WHEN l.start_page IS NOT NULL AND l.end_page IS NOT NULL 
            THEN l.end_page - l.start_page ELSE 0 END), 0) as total_pages,
        COUNT(l.id) as log_count,
        lr.id as log_id,
        lr.project_id as log_project_id,
        lr.data as log_data,
        lr.start_page as log_start_page,
        lr.end_page as log_end_page,
        lr.note as log_note
    FROM projects p
    LEFT JOIN log_ranked lr ON p.id = lr.project_id AND lr.rn <= 4
    LEFT JOIN logs l ON p.id = l.project_id
    GROUP BY p.id, p.name, p.total_page, p.page, 
             lr.id, lr.project_id, lr.data, lr.start_page, lr.end_page, lr.note
    ORDER BY 
        CASE 
            WHEN p.total_page = 0 THEN 0 
            ELSE p.page::float / p.total_page::float 
        END DESC,
        p.id ASC
`
```

**Error Handling Pattern:**
```go
if err != nil {
    return nil, fmt.Errorf("failed to get running projects with logs: %w", err)
}
```

**Context Timeout Pattern:**
```go
ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
defer cancel()
```

**NULL Handling Pattern:**
- Use `COALESCE` for all aggregate functions
- Use `LEFT JOIN` to handle projects without logs
- Check for `pgx.ErrNoRows` for empty result sets

**Scanning Pattern:**
- Use `rows.Scan` with explicit column mapping
- Handle nullable fields with pointer types
- Close rows with `defer rows.Close()`
- Check `rows.Err()` after iteration

---

### 5. Testing Strategy

**Unit Tests (`internal/service/dashboard/projects_service_test.go`):**
- Test `GetRunningProjectsWithLogs` with mock repository
- Verify status filtering logic (running vs finished vs stopped)
- Test progress calculation edge cases (division by zero, zero projects)
- Test sorting order (progress DESC, id ASC)
- Test empty result handling

**Integration Tests (`test/integration/dashboard_projects_test.go`):**
- Test with real database and sample data fixtures
- Verify SQL JOIN returns correct data structure
- Test first 4 logs per project (ordered by data DESC)
- Test progress ordering with various data scenarios
- Test NULL handling (projects without logs, NULL notes)
- Test performance (single query vs N+1)
- Rails parity validation (compare with Rails endpoint)

**Edge Cases to Cover:**
1. Projects with no logs (empty logs array)
2. Projects with fewer than 4 logs
3. Projects with exactly 4 logs
4. Projects with more than 4 logs (only first 4 returned)
5. Projects with zero total_page (division by zero)
6. Projects with NULL started_at
7. Mixed progress values (0%, 50%, 100%)
8. Projects with equal progress (ordered by id ASC)

**Test Data Setup:**
- Use `test/test_helper.go` for database setup/teardown
- Create fixtures with varying log counts
- Include edge cases in test data
- Clean up after each test using `defer helper.Close()`

---

### 6. Risks and Considerations

**Technical Risks:**
1. **SQL Complexity**: CTE with window functions and multiple JOINs may be complex to debug
   - *Mitigation*: Test query incrementally in PostgreSQL client first
   - *Mitigation*: Add detailed logging for query execution

2. **Performance**: Large datasets may impact query performance
   - *Mitigation*: Verify existing indexes on `logs(project_id, data DESC)`
   - *Mitigation*: Use EXPLAIN ANALYZE to verify query plan
   - *Mitigation*: Set appropriate timeout (15 seconds)

3. **Data Consistency**: Multiple JOINs may produce duplicate rows if not careful
   - *Mitigation*: Use proper GROUP BY clause
   - *Mitigation*: Test with known data to verify row counts

4. **NULL Handling**: Inconsistent NULL handling may cause calculation errors
   - *Mitigation*: Use COALESCE consistently for all aggregates
   - *Mitigation*: Test with NULL values explicitly

**Blocking Issues:**
- None identified - all dependencies are in place

**Trade-offs:**
- **Single Query vs Multiple Queries**: Single query is more efficient but more complex
  - Chose single query for performance (eliminates N+1 problem)
- **Window Functions vs Application-side Filtering**: Window functions push filtering to database
  - Chose window functions for efficiency and cleaner code

**Deployment Considerations:**
- No database migrations required
- No configuration changes required
- Backward compatible - method signature unchanged
- Can be deployed without downtime

**Rollback Plan:**
- Revert to previous implementation if issues arise
- Previous implementation (N+1 queries) still functional as fallback

---

### Implementation Steps Summary

1. **Update Repository Method**: Replace current implementation with single SQL JOIN query
2. **Verify Query**: Test SQL query in PostgreSQL client with sample data
3. **Run Unit Tests**: Verify service layer logic with existing tests
4. **Run Integration Tests**: Verify end-to-end behavior with real database
5. **Performance Test**: Verify query execution time meets targets (< 50ms)
6. **Code Review**: Ensure Clean Architecture compliance and code quality
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
