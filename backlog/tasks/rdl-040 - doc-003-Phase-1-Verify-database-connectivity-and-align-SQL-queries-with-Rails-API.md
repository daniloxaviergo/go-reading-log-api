---
id: RDL-040
title: >-
  [doc-003 Phase 1] Verify database connectivity and align SQL queries with
  Rails API
status: To Do
assignee:
  - thomas
created_date: '2026-04-12 23:50'
updated_date: '2026-04-13 00:13'
labels:
  - database
  - query
  - alignment
dependencies: []
references:
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/1'
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/7'
documentation:
  - doc-003
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement FR-001 and FR-007 by verifying the database connection string in .env points to the correct 'reading_log' database and ensuring all SQL queries in the adapter layer strictly replicate Rails ActiveRecord logic to guarantee identical result sets between the Go and Rails APIs.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Verify .env DB_DATABASE is set to 'reading_log'
- [ ] #2 Run integration tests confirming go_count equals rails_count
- [ ] #3 Audit all SQL queries in internal/adapter/postgres/queries.go against Rails ActiveRecord counterparts
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task focuses on **Phase 1: Database Alignment** from the PRD (RDL-039). The goal is to verify that the Go API connects to the correct database and that SQL queries produce identical results to the Rails API.

#### Key Technical Decisions:

**1.1 Database Connection Verification**
- The .env configuration must point to `reading_log` database
- The `LoadConfig()` function in `internal/config/config.go` already has `DB_DATABASE` defaulting to `reading_log`
- Need to verify this is actually being used by the PostgreSQL adapter

**1.2 SQL Query Alignment**
The Rails API uses ActiveRecord which generates SQL automatically. We need to ensure the raw SQL in `project_repository.go` replicates this behavior exactly.

**Rails ActiveQuery Pattern:**
```ruby
# Rails uses eager_load which generates LEFT OUTER JOIN
Project.eager_load(:logs).order('logs.data DESC')
```

**Generated SQL (expected):**
```sql
SELECT 
  projects.*, 
  logs.* 
FROM projects 
LEFT OUTER JOIN logs ON projects.id = logs.project_id 
ORDER BY logs.data DESC
```

**Current Go Implementation:**
The `GetAllWithLogs` method in `project_repository.go` already implements this pattern:
```go
query := `
    SELECT 
        p.id, p.name, p.total_page, p.started_at, p.page, p.reinicia,
        l.id as log_id, l.data::text as data_text, l.start_page, l.end_page, l.note, l.wday, l.text,
        l.created_at as log_created_at, l.updated_at as log_updated_at
    FROM projects p
    LEFT OUTER JOIN logs l ON p.id = l.project_id
    ORDER BY p.id ASC, l.data DESC
`
```

**1.3 Alignment Strategy:**
- Match Rails `eager_load` behavior with explicit `LEFT OUTER JOIN`
- Ensure column selection matches Rails serializer attributes
- Order by `logs.data DESC` to match Rails `-> { order(data: :desc) }` scope
- Handle null logs (projects without logs) via LEFT OUTER JOIN

### 2. Files to Modify

| File | Action | Reason |
|------|--------|--------|
| `internal/adapter/postgres/project_repository.go` | Review/Verify | Ensure SQL queries match Rails ActiveRecord logic exactly |
| `internal/adapter/postgres/log_repository.go` | Review/Verify | Ensure log queries match Rails ActiveRecord logic |
| `internal/config/config.go` | Review/Verify | Confirm DB_DATABASE defaults to `reading_log` |
| `.env` / `.env.example` | Update if needed | Ensure DB_DATABASE=reading_log is explicitly set |
| `test/compare_responses.sh` | Reference | Use for verification after implementation |
| `docs/database-alignment-report.md` | Create | Document findings and verification results |

**No new files required** - this is a verification and alignment task, not a feature implementation.

**Note:** The task description mentions `internal/adapter/postgres/queries.go` but this file doesn't exist. SQL queries are embedded in `project_repository.go` and `log_repository.go`.

### 3. Dependencies

**Prerequisites:**
- [ ] Task RDL-003 completed (PostgreSQL repository implementations) - **ALREADY DONE**
- [ ] Task RDL-007 completed (Application entry point) - **ALREADY DONE**
- [ ] PostgreSQL database `reading_log` must exist and be populated
- [ ] Rails API must be running on port 3001 for comparison testing
- [ ] Go API must be running on port 3000 for comparison testing

**Blocking Issues:**
- None identified - this task unblocks Phase 2 (Datetime Standardization) and Phase 3 (JSON Structure Harmonization)

### 4. Code Patterns

**4.1 SQL Query Pattern (from Rails to Go):**

Rails uses `eager_load` which creates LEFT OUTER JOIN:
```ruby
# Rails
Project.eager_load(:logs).order('logs.data DESC')
```

Go equivalent using pgx:
```go
query := `
    SELECT p.id, p.name, p.total_page, p.started_at, p.page, p.reinicia,
           l.id as log_id, l.data::text as data_text, l.start_page, l.end_page, l.note, l.wday, l.text
    FROM projects p
    LEFT OUTER JOIN logs l ON p.id = l.project_id
    ORDER BY p.id ASC, l.data DESC
`
```

**Key Alignment Points:**
1. `LEFT OUTER JOIN` ensures projects without logs are included
2. `ORDER BY p.id ASC, l.data DESC` matches Rails eager loading behavior
3. `data::text` cast avoids binary format scanning issues
4. Explicit column selection matches Rails serializer attributes

**4.2 Project Serializer Attributes (Rails):**
```ruby
class ProjectSerializer < ActiveModel::Serializer
  attributes :id, :name, :total_page, :started_at, :page, :reinicia,
             :progress, :status, :logs_count, :days_unreading, :median_day, :finished_at
end
```

**4.3 Log Serializer Attributes (Rails):**
```ruby
class LogSerializer < ActiveModel::Serializer
  attributes :id, :data, :start_page, :end_page, :note
  belongs_to :project
end
```

### 5. Testing Strategy

**5.1 Unit Tests (already exist):**
- `internal/adapter/postgres/*_test.go` - Repository unit tests
- `internal/config/*_test.go` - Configuration tests
- Run with: `go test -v ./internal/...`

**5.2 Integration Tests:**
- Use `test/test_helper.go` for database setup/teardown
- Create test projects and logs
- Verify queries return expected results

**5.3 Comparison Tests:**
- Execute `test/compare_responses.sh` script
- Compare Go API vs Rails API responses
- Verify identical project IDs and counts
- Verify identical calculated fields (progress, status, logs_count, days_unreading, median_day, finished_at)

**5.4 Test Scenarios:**
1. **Empty database:** Verify empty array returned
2. **Single project with logs:** Verify complete data returned
3. **Multiple projects with varying log counts:** Verify all projects returned
4. **Project without logs:** Verify LEFT OUTER JOIN includes project with nil log fields
5. **Calculation accuracy:** Verify progress, status, days_unreading match Rails

### 6. Risks and Considerations

**6.1 Known Risks:**

| Risk | Impact | Mitigation |
|------|--------|------------|
| Database connection points to wrong database | CRITICAL - Wrong data | Verify .env DB_DATABASE=reading_log; add connection logging |
| SQL query differs from Rails ActiveRecord | HIGH - Inconsistent results | Explicitly match Rails eager_load behavior |
| Datetime format mismatch | MEDIUM - Parsing issues | Ensure RFC3339 format with timezone offset |
| Calculated fields differ | MEDIUM - Business logic errors | Verify against Rails implementation |

**6.2 Trade-offs:**

1. **Direct SQL vs ORM:** Using raw SQL (pgx) gives control but requires manual replication of ActiveRecord logic
   - *Decision:* Accept this for Phase 1; consider ORM migration in Phase 2

2. **LEFT OUTER JOIN vs INNER JOIN:** LEFT OUTER JOIN includes projects without logs
   - *Decision:* Keep LEFT OUTER JOIN to match Rails `eager_load` behavior

3. **Ordering:** Rails orders by `logs.data DESC`, Go orders by `p.id ASC, l.data DESC`
   - *Decision:* This is correct - projects ordered by ID, logs within each project ordered by data DESC

**6.3 Verification Checklist:**

Before marking task complete:
- [ ] `make docker-up` successfully starts both Go and Rails APIs
- [ ] `test/compare_responses.sh` passes all tests
- [ ] Go API returns same project count as Rails API
- [ ] Project IDs match between APIs
- [ ] All calculated fields match within tolerance (0.01)
- [ ] No breaking changes to existing functionality
- [ ] `go vet` passes without errors
- [ ] `go fmt` applied to all modified files
- [ ] Documentation updated in QWEN.md

**6.4 Rollback Plan:**
If verification fails:
1. Revert to previous commit
2. Review comparison report for specific failures
3. Fix individual issues iteratively
4. Re-run comparison script after each fix

### 7. Implementation Steps

**Phase 1: Verification (This Task - RDL-040)**

1. Verify .env configuration
   ```bash
   # Check .env file
   grep DB_DATABASE .env
   # Should show: DB_DATABASE=reading_log
   ```

2. Review and verify SQL queries in `project_repository.go` and `log_repository.go`
   - Confirm `GetAllWithLogs` uses LEFT OUTER JOIN
   - Confirm ORDER BY matches Rails behavior
   - Confirm column selection matches serializer

3. Run comparison tests
   ```bash
   make docker-up
   ./test/compare_responses.sh
   ```

4. Document findings in `docs/database-alignment-report.md`

**Phase 2: Fixes (If Needed)**

If issues found during verification:
1. Fix SQL queries in `project_repository.go` or `log_repository.go`
2. Update DTOs in `internal/domain/dto/` if needed
3. Re-run comparison tests

**Phase 3: Sign-off**

1. All tests pass
2. Comparison script reports 0 failures
3. Tech Lead approval
4. Task closed with verification evidence

---

## Summary

This implementation plan addresses **Task RDL-040: Verify database connectivity and align SQL queries with Rails API**.

**Primary Goal:** Ensure the Go API connects to the `reading_log` database and executes SQL queries identical to Rails ActiveRecord, guaranteeing consistent data between both APIs.

**Key Activities:**
1. Verify .env configuration points to correct database
2. Audit SQL queries in `internal/adapter/postgres/project_repository.go` and `log_repository.go` against Rails `eager_load` behavior
3. Run comprehensive comparison tests using `test/compare_responses.sh`
4. Document any discrepancies and fixes
5. Obtain stakeholder sign-off

**Expected Outcome:** A verified alignment between Go and Rails API database connectivity and query logic, unblocking Phase 2 (Datetime Standardization) and Phase 3 (JSON Structure Harmonization).
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress: RDL-040

### Status: In Progress

### Completed Steps:

**1. Environment Verification**
- Checked .env configuration for correct database name
- Verified DB_DATABASE is set to 'reading_log' ✓

**2. SQL Query Audit**
- Reviewed `internal/adapter/postgres/project_repository.go`
- Reviewed `internal/adapter/postgres/log_repository.go`
- Comparing against Rails ActiveRecord logic

**3. Test Execution**
- Unit tests: **PASS** (all pass)
- Integration tests: **FAIL** (PostgreSQL authentication issues - environment setup required)
- `go vet` and `go fmt`: Need to run

### Test Results Summary:

```
PASS: go-reading-log-api-next/internal/api/v1
PASS: go-reading-log-api-next/internal/config
PASS: go-reading-log-api-next/internal/domain/dto
PASS: go-reading-log-api-next/internal/domain/models
PASS: go-reading-log-api-next/test/unit
FAIL: go-reading-log-api-next/test/integration (PostgreSQL connection issues)
```

### Current State:
- Task status: To Do → In Progress
- Priority: HIGH
- Blocking: RDL-041, RDL-042, RDL-043, RDL-044

### Next Steps:
1. Run `go vet` and `go fmt` checks
2. Verify acceptance criteria met where possible
3. Document test results
4. Update task status appropriately
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
<!-- DOD:END -->
