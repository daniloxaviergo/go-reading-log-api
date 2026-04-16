---
id: RDL-051
title: '[doc-004 Phase 1.4] Add make test-clean command'
status: To Do
assignee:
  - workflow
created_date: '2026-04-15 12:14'
updated_date: '2026-04-16 09:44'
labels:
  - build
  - automation
  - medium-priority
dependencies: []
references:
  - 'R4: Make Command for Manual Cleanup'
documentation:
  - doc-004
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add a make test-clean command to the Makefile that provides a manual cleanup mechanism for orphaned test databases. The command should execute a standalone cleanup script, provide progress feedback during execution, and handle errors gracefully without crashing. Include a test-cleanup alias for convenience. Ensure the Makefile targets use colorized output consistent with existing commands.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Command is available in Makefile
- [ ] #2 It drops all orphaned test databases
- [ ] #3 It provides progress feedback
- [ ] #4 It handles errors gracefully
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The implementation will create a standalone cleanup script that can be invoked both by the Makefile and directly by developers. This approach provides:

- **Flexibility:** The script can run standalone or via Makefile
- **Reusability:** Same code used for both manual and automated cleanup
- **Maintainability:** Clear separation between cleanup logic and build system
- **Error Handling:** Graceful handling of database connection issues and SQL errors

The cleanup script will:
1. Load configuration from `.env.test` (same as test helper)
2. Connect to the main database
3. Query for orphaned test databases (pattern: `reading_log_test_%`)
4. Exclude the current test database if specified
5. Drop each orphaned database with `DROP DATABASE IF EXISTS`
6. Print progress feedback for each database
7. Handle errors without crashing (continue with remaining databases)

### 2. Files to Modify

| File | Action | Description |
|------|--------|-------------|
| `test/cleanup_orphaned_databases.go` | **Create** | Standalone cleanup script with progress feedback |
| `Makefile` | **Modify** | Fix `test-clean` target, add `test-cleanup` alias |
| `test/test_helper_test.go` | **Modify** | Add unit tests for the new cleanup script |

**Files to Create:**
- `test/cleanup_orphaned_databases.go` - Main cleanup script (150-200 lines)

**Files to Modify:**
- `Makefile` - Lines 104-111 (replace existing test-clean implementation)
- `Makefile` - Add `test-cleanup` alias after test-clean target

### 3. Dependencies

**Prerequisites:**
- Go 1.25.7 (already in use by the project)
- pgx/v5 library (already installed)
- PostgreSQL running and accessible
- `.env.test` file with database credentials

**No new dependencies required.**

### 4. Code Patterns

**Follow existing patterns from the codebase:**

1. **Configuration Loading:** Use `godotenv.Load(".env.test")` and `config.LoadConfig()`
2. **Database Connection:** Use `pgxpool.New()` with context timeout
3. **Error Handling:** Wrap errors with `fmt.Errorf("...: %w", err)` and log without failing
4. **Progress Feedback:** Use `fmt.Printf()` with color codes from Makefile
5. **SQL Queries:** Use parameterized queries where possible, `IF EXISTS` for safety
6. **Context Usage:** 60-second timeout for cleanup operations

**Naming Conventions:**
- Function names: `cleanupOrphanedDatabases`, `isValidTestDatabaseName`
- Variable names: `testDBName`, `mainPool`, `orphanedDBs`
- File name: `cleanup_orphaned_databases.go`

**Integration Pattern:**
```go
// The cleanup script will:
// 1. Load .env.test
// 2. Connect to main database
// 3. Query for orphaned databases
// 4. Drop each database with DROP DATABASE IF EXISTS
// 5. Print progress: "Dropping reading_log_test_123..."
// 6. Log errors but continue
// 7. Exit with code 0 (graceful)
```

### 5. Testing Strategy

**Unit Tests (add to `test_helper_test.go`):**

| Test Function | Purpose | Coverage |
|---------------|---------|----------|
| `TestCleanupScript_BasicCleanup` | Verify script runs without error | 80% |
| `TestCleanupScript_ConcurrentDBs` | Test with multiple parallel databases | 70% |
| `TestCleanupScript_Performance` | Verify cleanup completes in < 60s | 60% |

**Integration Tests:**
- Run `make test-clean` and verify exit code is 0
- Create orphaned databases manually, run cleanup, verify they're gone
- Test with invalid database names (should handle gracefully)

**Test Execution:**
```
make test-clean
# Expected: Exit code 0, progress messages printed

# Verify cleanup worked
psql -U postgres -d reading_log -c "SELECT datname FROM pg_database WHERE datname LIKE 'reading_log_test_%';"
# Expected: Empty or only current test DB listed
```

### 6. Risks and Considerations

**Blocking Issues:**
- None identified. Implementation uses existing patterns from `test_helper.go`.

**Potential Pitfalls:**
1. **Database Connection Pool Exhaustion:** Cleanup uses a separate pool to avoid conflicts
2. **Long-Running Queries:** 60-second timeout prevents hanging
3. **Permission Issues:** Current user must have DROP DATABASE privileges
4. **Race Conditions:** Cleanup excludes current test DB by name

**Trade-offs:**
- Using `DROP DATABASE IF EXISTS` is safer than checking existence first (avoids race condition)
- Progress feedback is printed to stdout (not logged) for immediate visibility
- Cleanup continues on error (doesn't stop at first failure) to maximize cleanup

**Deployment Considerations:**
- No migration required (no schema changes)
- No downtime required
- Can be run at any time, but best during low-usage periods
- Consider running `VACUUM ANALYZE` after cleanup for PostgreSQL optimization

**Verification Checklist:**
- [ ] `make test-clean` runs without error
- [ ] Progress messages are printed to console
- [ ] Orphaned databases are dropped
- [ ] Current test database is NOT dropped
- [ ] Script exits with code 0 even if some drops fail
- [ ] `make test-cleanup` alias works identically to `make test-clean`

### 7. Implementation Steps

**Step 1: Create `test/cleanup_orphaned_databases.go`**
- Import required packages (context, fmt, os, pgxpool, godotenv)
- Load configuration from `.env.test`
- Connect to main database
- Query orphaned databases
- Drop each with progress feedback
- Handle errors gracefully

**Step 2: Update `Makefile`**
- Replace existing `test-clean` target with cleaner implementation
- Add `test-cleanup` alias
- Ensure colorized output consistency

**Step 3: Add Tests**
- Add unit tests for cleanup script
- Add integration test for `make test-clean` command

**Step 4: Verification**
- Run `make test-clean` manually
- Create test databases, verify cleanup works
- Check exit codes and error handling
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
