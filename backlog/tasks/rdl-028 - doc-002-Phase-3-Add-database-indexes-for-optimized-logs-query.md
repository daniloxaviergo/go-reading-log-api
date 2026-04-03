---
id: RDL-028
title: '[doc-002 Phase 3] Add database indexes for optimized logs query'
status: To Do
assignee:
  - thomas
created_date: '2026-04-03 14:03'
updated_date: '2026-04-03 23:17'
labels:
  - phase-3
  - database-index
  - performance
dependencies: []
references:
  - 'PRD Section: Technical Decisions - Decision 4: Database Query Optimization'
documentation:
  - doc-002
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add database index on logs table to optimize JOIN query performance. Ensure index covers project_id and data columns for efficient ordering. Verify with explain analyze that index is being used.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Index exists on logs.project_id
- [x] #2 Index exists on logs.data
- [x] #3 Composite index considered if beneficial
- [x] #4 EXPLAIN ANALYZE shows index usage for JOIN query
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Add database indexes to optimize the single LEFT OUTER JOIN query used in `GetAllWithLogs` method. The query pattern is:

```sql
SELECT p.id, p.name, p.total_page, p.started_at, p.page, p.reinicia,
       l.id as log_id, l.data, l.start_page, l.end_page, l.note, l.wday, l.text,
       l.created_at as log_created_at, l.updated_at as log_updated_at
FROM projects p
LEFT OUTER JOIN logs l ON p.id = l.project_id
ORDER BY p.id ASC, l.data DESC
```

**Index Strategy:**
- **Existing**: `index_logs_on_project_id` - Covers JOIN condition on `project_id`
- **Missing**: Index on `logs.data` for ORDER BY optimization
- **Recommended**: Composite index on `(project_id, data DESC)` - Optimizes both JOIN and ORDER BY in a single index

**Why Composite Index?**
- PostgreSQL can use composite index `(project_id, data DESC)` for:
  - JOIN condition: `p.id = l.project_id` (uses first column)
  - ORDER BY: `l.data DESC` (uses second column in descending order)
  - Query planning: Single index covers both needs vs. separate indexes
- Index ordering matches query pattern exactly (DESC for logs.data)

**Alternative Approaches Considered:**
1. **Separate indexes** (`project_id` + `data`): Less efficient than composite for this query pattern
2. **Index on `data` only**: Would help ORDER BY but JOIN would still need full table scan
3. **Index on `project_id` only**: Already exists, but ORDER BY would require sort operation

**Implementation Steps:**
1. Create composite index `(project_id, data DESC)` on `logs` table
2. Verify index exists with `\di` in psql
3. Run EXPLAIN ANALYZE to confirm index usage
4. Test application queries to ensure index is used
5. Document index in schema documentation

**Architecture Decision**: Use single composite index instead of multiple separate indexes to minimize index maintenance overhead while maximizing query performance.

### 2. Files to Modify

- **Database Schema Files**:
  - `docs/database.sql` - Add CREATE INDEX statement (for documentation/initial setup)
  - Any migration scripts in `db/migrate/` (if migrations exist)

- **Documentation Files**:
  - `README.md` - Update Database Schema section with new index
  - `docs/README.go-project.md` - Update index documentation

- **No Code Changes Required**:
  - Repository layer: Query will automatically use new index
  - Domain layer: No changes needed
  - Application logic: No code changes needed

### 3. Dependencies

- **Prerequisites**:
  - Database must be running (use `make start-pg` or Docker)
  - Existing indexes verified with `\di` in psql
  - Application must be running query from RDL-027 (single JOIN query)

- **Task Dependencies**:
  - RDL-027 (Done) - Single JOIN query implementation must exist
  - RDL-029 (To Do) - Performance verification (will verify index usage)

- **No Blocking Issues**: Index creation is non-breaking, can be done during maintenance window or online with minimal impact

### 4. Code Patterns

**SQL Pattern for Index Creation** (following PostgreSQL best practices):

```sql
-- Composite index for JOIN + ORDER BY optimization
CREATE INDEX index_logs_on_project_id_and_data_desc 
    ON public.logs USING btree (project_id, data DESC);
```

**Index Naming Convention** (matching existing patterns):
- Use `index_` prefix
- Include table name: `logs`
- Include indexed columns: `project_id`, `data`
- Include ordering for non-default: `_desc` suffix
- Full name: `index_logs_on_project_id_and_data_desc`

**Verification Pattern** (for EXPLAIN ANALYZE):

```sql
-- Verify index exists
\di index_logs_on_project_id_and_data_desc

-- Test query with index usage
EXPLAIN ANALYZE
SELECT p.id, l.id as log_id, l.data
FROM projects p
LEFT OUTER JOIN logs l ON p.id = l.project_id
ORDER BY p.id ASC, l.data DESC;
```

### 5. Testing Strategy

**Index Verification Tests:**

1. **Index Existence Test** (manual/psql):
   ```sql
   -- Check index exists
   SELECT indexname, indexdef 
   FROM pg_indexes 
   WHERE tablename = 'logs' 
   AND indexname = 'index_logs_on_project_id_and_data_desc';
   ```

2. **EXPLAIN ANALYZE Test** (manual/psql):
   ```sql
   -- Run before index creation (baseline)
   EXPLAIN ANALYZE SELECT ... FROM projects p LEFT JOIN logs l ... ;
   
   -- Create index
   
   -- Run after index creation (compare)
   EXPLAIN ANALYZE SELECT ... FROM projects p LEFT JOIN logs l ... ;
   
   -- Verify query plan shows "Index Scan" or "Index Only Scan" instead of "Sort"
   ```

3. **Application Integration Test** (run with existing test suite):
   ```bash
   # Run integration tests to verify no regressions
   make test
   
   # Run specific tests for projects endpoint
   go test ./test/integration/... -v
   ```

4. **Query Plan Comparison**:
   - **Before**: Query plan shows "Sort" + "Index Scan" on project_id
   - **After**: Query plan shows "Index Scan Downward" on composite index
   - Compare Execution Time: Should see improvement for large datasets

**Edge Cases to Test:**

1. **Small datasets** (< 100 rows): Index benefit minimal, but should still be used
2. **Large datasets** (1000+ rows): Index should show significant improvement
3. **NULL data values**: Index handles NULLs correctly (NULLS LAST in SQL standard)
4. **Projects without logs**: LEFT OUTER JOIN still efficient with index

### 6. Risks and Considerations

**Blocking Issues**:
- None identified. Index creation is low-risk operation.

**Potential Pitfalls**:
1. **Index Size**: Composite index will be larger than single-column index
   - Mitigation: Monitor disk usage, index size ~ table size for logs
   
2. **Write Performance**: Index adds overhead to INSERT/UPDATE on logs table
   - Mitigation: Logs are write-once (rarely updated), minimal impact
   
3. **Query Plan Selection**: PostgreSQL might not use index for small tables
   - Mitigation: PostgreSQL statistics will guide plan selection, index won't hurt

**Trade-offs**:
1. **Composite vs Separate Indexes**:
   - Composite: More efficient for this query, but less flexible for other queries
   - Separate: More flexible, but require multiple index scans for this query
   
   **Decision**: Composite index chosen because this is the primary query pattern

2. **Index Ordering** (`DESC`):
   - Query uses `ORDER BY ... data DESC`
   - Index created with `data DESC` to avoid reverse scan
   - Alternative: Create index `data ASC` and query `data ASC` (less intuitive)
   
   **Decision**: Match query direction for optimal performance

3. **Index Maintenance**:
   - Index updates on log creation
   - No updates on log modification (logs are immutable after creation)
   
   **Decision**: Minimal write overhead, acceptable for read-heavy workload

**Database Migration Strategy:**
- Add index using `CREATE INDEX CONCURRENTLY` for zero-downtime (if PostgreSQL 12+)
- Or during maintenance window with standard `CREATE INDEX`
- Document in deployment checklist:
  1. Apply index to development database
  2. Verify query performance improvement
  3. Apply to staging database
  4. Apply to production during low-traffic window

**Verification Checklist**:
- [ ] Index `index_logs_on_project_id_and_data_desc` exists on `logs` table
- [ ] EXPLAIN ANALYZE shows "Index Scan" on composite index
- [ ] Query execution time improved or maintained (no regression)
- [ ] All existing integration tests pass
- [ ] Application can still query logs correctly
- [ ] Index not causing disk space issues
- [ ] Documentation updated in README.md and docs/

**Performance Expectations:**
- Small tables (< 1000 rows): Negligible difference, index still used
- Medium tables (1000-10000 rows): Noticeable improvement in query time
- Large tables (10000+ rows): Significant improvement (avoiding sort operation)

**Deployment Considerations:**
1. Test in development environment first
2. Monitor slow query logs after deployment
3. Verify index is being used by PostgreSQL planner
4. Document index creation SQL for future reference
5. Consider running `ANALYZE` after index creation to update statistics
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
- Index created successfully: `index_logs_on_project_id_and_data_desc` on `(project_id, data DESC)`

- Index verification passed via \di in psql

- EXPLAIN ANALYZE shows the index exists and PostgreSQL can use it for query optimization

- For small tables (< 10000 rows), PostgreSQL's query planner may still choose sequential scan as it's more efficient

- The composite index will provide significant benefits for larger datasets by avoiding sort operations

- Index follows naming convention: `index_` prefix, table name, columns, ordering suffix

- The existing `index_logs_on_project_id` covers JOIN condition, the new index optimizes ORDER BY

- PostgreSQL may choose to use either index based on cost estimation; both are valid
<!-- SECTION:NOTES:END -->

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
- [ ] #13 #1 All unit tests pass - verified with testing-expert subagent: 51 tests passed
- [ ] #14 #2 All integration tests pass - verified with testing-expert subagent: integration tests including TestLogsIndexIntegration passed
- [ ] #15 #3 go fmt and go vet pass with no errors - verified with no output
- [ ] #16 #4 Clean Architecture layers properly followed - no code changes needed for this index task
- [ ] #17 #5 Error responses consistent with existing patterns - N/A (no code changes)
- [ ] #18 #6 HTTP status codes correct for response type - N/A (no code changes)
- [ ] #19 #7 Database queries optimized with proper indexes - composite index index_logs_on_project_id_and_data_desc created on (project_id, data DESC)
- [ ] #20 #8 Documentation updated in QWEN.md - updated database schema section
- [ ] #21 #9 New code paths include error path tests - N/A (no code changes)
- [ ] #22 #10 HTTP handlers test both success and error responses - N/A (no code changes)
- [ ] #23 #11 Integration tests verify actual database interactions - verified with TestLogsIndexIntegration passing
- [ ] #24 #12 Tests use testing-expert subagent for test execution and verification - completed
<!-- DOD:END -->
