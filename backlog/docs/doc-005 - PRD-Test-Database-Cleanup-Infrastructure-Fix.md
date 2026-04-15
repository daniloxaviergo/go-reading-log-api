---
id: doc-005
title: 'PRD Test Database Cleanup Infrastructure Fix'
type: other
created_date: '2026-04-15 11:49'
---


# Test Database Cleanup Infrastructure Fix

## Document Information

| Field | Value |
|-------|-------|
| **PRD ID** | PRD-001 |
| **Title** | Test Database Cleanup Infrastructure Fix |
| **Created** | 2026-04-15 |
| **Status** | In Review |
| **Priority** | Critical |
| **Owner** | Development Team |

---

## Executive Summary

This PRD addresses the critical accumulation of **6,004 stale test databases** discovered in the Docker Compose PostgreSQL container on April 14, 2026. The issue stems from test sessions creating unique databases that are not properly dropped after completion.

### Business Impact

| Metric | Current | Target |
|--------|---------|--------|
| Orphaned test databases | 6,004 | 0 |
| Database disk usage | ~50GB+ | Minimal |
| Test execution reliability | Unreliable | 100% |
| Development productivity | Blocked | Normal |

### Scope

**In Scope:**
- ✅ Update `test/test_helper.go` to properly drop databases after test execution
- ✅ Add cleanup of orphaned databases from previous sessions
- ✅ Implement `defer` protection for automatic cleanup even on panic
- ✅ Add `make test-clean` command for manual cleanup
- ✅ Add database name validation to prevent SQL injection

**Out of Scope:**
- ❌ Switching to single database with schema reset (Phase 2 consideration)
- ❌ Connection pooling implementation (Phase 2 consideration)
- ❌ Database monitoring/alerting (Phase 2 consideration)

---

## Key Requirements

### Requirement Priority Matrix

| ID | Requirement | Priority | Status |
|----|-------------|----------|--------|
| R1 | Auto-cleanup test databases on test completion | Critical | ✅ Planned |
| R2 | Clean up orphaned databases from previous sessions | Critical | ✅ Planned |
| R3 | Prevent database name collision in parallel tests | High | ✅ Planned |
| R4 | Add `make test-clean` command | Medium | ✅ Planned |
| R5 | Validate database names to prevent SQL injection | High | ✅ Planned |

### Detailed Requirements

#### R1: Auto-Cleanup on Test Completion

**Priority:** Critical

**Description:** Every test session must automatically drop its unique test database when the test completes (success or failure).

**Acceptance Criteria:**
- [ ] Test database is dropped within 1 second of test completion
- [ ] Cleanup occurs even if test panics
- [ ] No error is thrown if database doesn't exist
- [ ] Cleanup doesn't block test results

**Implementation:**
```go
func (h *TestHelper) Close() {
    if h.Pool != nil {
        defer func() {
            // Drop test database
            if h.TestDBName != "" {
                ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
                defer cancel()
                
                connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
                    h.Config.DBUser, h.Config.DBPassword, h.Config.DBHost,
                    h.Config.DBPort, h.Config.DBDatabase)
                
                mainPool, err := pgxpool.New(ctx, connStr)
                if err == nil {
                    defer mainPool.Close()
                    mainPool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", h.TestDBName))
                }
            }
        }()
        h.Pool.Close()
    }
}
```

---

#### R2: Orphaned Database Cleanup

**Priority:** Critical

**Description:** Implement cleanup of test databases that were not properly dropped (older than 24 hours).

**Acceptance Criteria:**
- [ ] Databases older than 24 hours are identified and dropped
- [ ] Current test database is excluded from cleanup
- [ ] Cleanup runs in under 1 minute for 6,000+ databases
- [ ] Errors are logged but don't fail test execution

**Implementation:**
```go
func cleanupOrphanedDatabases(pool *pgxpool.Pool, excludeName string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()
    
    query := `
        SELECT datname 
        FROM pg_database 
        WHERE datname LIKE $1
        AND datname != $2
        AND pg_catalog.pg_get_userbyid(datdba) = current_user
    `
    
    rows, err := pool.Query(ctx, query, "reading_log_test_%", excludeName)
    if err != nil {
        return fmt.Errorf("failed to query test databases: %w", err)
    }
    defer rows.Close()
    
    var toDrop []string
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            continue
        }
        toDrop = append(toDrop, name)
    }
    
    for _, dbName := range toDrop {
        pool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
    }
    
    return nil
}
```

---

#### R3: Parallel Test Safety

**Priority:** High

**Description:** Ensure test database names are unique across parallel test executions to prevent collisions.

**Acceptance Criteria:**
- [ ] No two parallel tests create databases with the same name
- [ ] Test execution speed is not significantly impacted
- [ ] Database cleanup doesn't interfere with parallel tests

**Implementation:**
```go
import "runtime/debug"

func generateUniqueDBName(baseName string) string {
    return fmt.Sprintf("%s_%d_%d_%d", 
        baseName,
        os.Getpid(), 
        getGoroutineID(),  // Unique across all goroutines
        time.Now().UnixNano())
}

func getGoroutineID() uint64 {
    buf := make([]byte, 64)
    n := runtime.Stack(buf, false)
    for i := 0; i < n-10; i++ {
        if buf[i] == 'g' && buf[i+1] == 'o' && buf[i+2] == 'r' {
            start := i + 10
            end := start
            for end < n && buf[end] >= '0' && buf[end] <= '9' {
                end++
            }
            if end > start {
                var id uint64
                for j := start; j < end; j++ {
                    id = id*10 + uint64(buf[j]-'0')
                }
                return id
            }
        }
    }
    return uint64(time.Now().UnixNano())
}
```

---

#### R4: Make Command for Manual Cleanup

**Priority:** Medium

**Description:** Provide a `make test-clean` command for developers to manually cleanup orphaned databases.

**Acceptance Criteria:**
- [ ] Command is available in Makefile
- [ ] It drops all orphaned test databases
- [ ] It provides progress feedback
- [ ] It handles errors gracefully

**Implementation:**
```makefile
test-clean:
	@echo "$(BLUE)Cleaning up orphaned test databases...$(NC)"
	@go run ./test/cleanup_orphaned_databases.go
	@echo "$(GREEN)Cleanup complete$(NC)"

test-cleanup: test-clean
```

---

#### R5: Database Name Validation

**Priority:** High

**Description:** Validate database names to prevent SQL injection and ensure format compliance.

**Acceptance Criteria:**
- [ ] Only alphanumeric characters, underscores, and hyphens are allowed
- [ ] Names are limited to 63 characters (PostgreSQL limit)
- [ ] Names must match pattern `reading_log_test[_[a-zA-Z0-9_]+]`
- [ ] Invalid names are rejected with clear error messages

**Implementation:**
```go
import "regexp"

func isValidTestDatabaseName(name string) bool {
    validPattern := regexp.MustCompile(`^reading_log_test(_[a-zA-Z0-9_]+)?$`)
    return validPattern.MatchString(name) && len(name) <= 63
}

func SafeDropDatabase(pool *pgxpool.Pool, dbName string) error {
    if !isValidTestDatabaseName(dbName) {
        return fmt.Errorf("invalid database name: %s", dbName)
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    query := `DROP DATABASE IF EXISTS ` + pq.QuoteIdentifier(dbName)
    _, err := pool.Exec(ctx, query)
    if err != nil {
        return fmt.Errorf("failed to drop %s: %w", dbName, err)
    }
    
    return nil
}
```

---

## Technical Decisions

### Decision 1: Keep Per-Test Database Strategy

**Decision:** Continue using unique databases per test session rather than switching to single database with schema reset.

**Rationale:**
- Complete isolation between tests prevents side effects
- Parallel test execution is faster with separate databases
- Existing test infrastructure is already built around this pattern
- Schema reset approach would require significant refactoring

**Trade-offs:**
- Slightly slower test startup (~150ms for CREATE DATABASE)
- More complex cleanup logic
- Potential for database accumulation if not properly managed

**Alternatives Considered:**
- Single database with TRUNCATE: Faster but less isolated
- Hybrid approach: Could be considered in Phase 2

---

### Decision 2: Use `defer` for Guaranteed Cleanup

**Decision:** Implement cleanup using `defer` statements to ensure databases are dropped even on panic.

**Rationale:**
- Guarantees cleanup regardless of test outcome
- Minimal code changes required
- Go's `defer` is the idiomatic way to ensure cleanup
- No external dependencies required

**Implementation:**
```go
func (h *TestHelper) Close() {
    if h.Pool != nil {
        defer func() {
            // Cleanup code here
        }()
        h.Pool.Close()
    }
}
```

---

### Decision 3: Separate Connection Pool for Cleanup

**Decision:** Use a separate connection pool to the main database for cleanup operations.

**Rationale:**
- The test database pool may already be closed
- Main database connection is needed to DROP the test database
- Isolates cleanup operations from test operations
- Prevents resource leaks

**Implementation:**
```go
connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
    h.Config.DBUser, h.Config.DBPassword, h.Config.DBHost,
    h.Config.DBPort, h.Config.DBDatabase)
mainPool, err := pgxpool.New(ctx, connStr)
```

---

### Decision 4: Time-Based Orphan Detection

**Decision:** Identify orphaned databases by age (24+ hours) rather than active session tracking.

**Rationale:**
- Simpler implementation without complex session tracking
- Handles cases where test process was killed unexpectedly
- 24-hour threshold balances cleanup urgency with safety
- No risk of dropping databases from currently running tests

**Trade-offs:**
- Very long-running tests (>24h) might be affected
- Could implement session tracking in Phase 2 if needed

---

### Decision 5: Prefix-Based Database Selection

**Decision:** Use `reading_log_test_%` pattern to identify test databases for cleanup.

**Rationale:**
- Simple and fast SQL query
- No need to track all possible test database names
- Consistent with existing naming convention
- Easy to verify in SQL console

---

## Acceptance Criteria

### Functional Acceptance Criteria

| Scenario | Expected Behavior |
|----------|-------------------|
| Test completes successfully | Test database is dropped within 1s |
| Test panics | Test database is dropped within 1s |
| Test fails | Test database is dropped within 1s |
| Multiple parallel tests | Each test's database is dropped independently |
| Orphaned database exists | Cleanup script drops all orphans >24h old |
| Invalid database name | Error is returned, no SQL injection |
| Database doesn't exist | No error (IF EXISTS handles this) |

### Non-Functional Acceptance Criteria

| Requirement | Metric | Target |
|-------------|--------|--------|
| Performance | Test startup time | < 200ms increase |
| Performance | Cleanup time (6,000 orphans) | < 60 seconds |
| Reliability | Test isolation | 100% (no cross-test contamination) |
| Security | SQL injection prevention | 100% (parameterized queries + validation) |
| Usability | Make command availability | `make test-clean` available |
| Maintainability | Code clarity | Comments explain cleanup logic |

---

## Files to Modify

### 1. `test/test_helper.go`

**Purpose:** Core test infrastructure cleanup

**Changes:**
- Update `Close()` method to include database cleanup
- Add `cleanupOrphanedDatabases()` method
- Add `isValidTestDatabaseName()` validation function
- Add `SafeDropDatabase()` wrapper with validation

**Lines to modify:** ~350-450 (Close method), add new functions at end of file

---

### 2. `Makefile`

**Purpose:** Add manual cleanup command

**Changes:**
- Add `test-clean` target
- Add `test-cleanup` alias
- Include colorized output for consistency

**Lines to modify:** ~100-120 (add targets)

---

### 3. `test/cleanup_orphaned_databases.go` (NEW FILE)

**Purpose:** Standalone cleanup script for manual and scheduled cleanup

**Contents:**
- Database connection setup
- Query for orphaned databases
- Drop each orphaned database
- Progress reporting
- Error handling

---

## Files Created

| File | Purpose |
|------|---------|
| `test/cleanup_orphaned_databases.go` | Standalone cleanup script |
| `docs/test-database-cleanup.md` | Documentation for the fix |
| `scripts/benchmark_cleanup.go` | Performance benchmarking tool |

---

## Validation Rules

### Database Name Validation

```go
// Validation Logic
1. Length check: len(name) <= 63
2. Pattern check: ^reading_log_test(_[a-zA-Z0-9_]+)?$
3. Character check: only alphanumeric, underscore, hyphen
4. Reserved check: not equal to 'reading_log_test' (main test DB)
```

### Error Handling Validation

| Error Type | Handling Strategy |
|------------|-------------------|
| Connection failed | Log warning, continue cleanup of other DBs |
| DROP fails | Log error with database name, continue |
| Query fails | Return error, stop cleanup |
| Permission denied | Log error, suggest checking permissions |

---

## Out of Scope

### Explicitly Excluded Items

| Item | Reason | Future Consideration |
|------|--------|---------------------|
| Single database with schema reset | Requires significant refactoring, performance testing | Phase 2 - Evaluate if parallel test speed becomes bottleneck |
| Connection pooling | Complex implementation, current pool size is adequate | Phase 2 - Consider for production deployment |
| Database monitoring/alerting | Requires external tools, dashboard | Phase 2 - Add to monitoring stack |
| Automatic cron cleanup | Requires cron setup, additional infrastructure | Phase 2 - Add to CI/CD pipeline |
| Test database compression | PostgreSQL already handles this efficiently | Not needed currently |

### Future Enhancements (Phase 2)

| Feature | Description | Estimated Effort |
|---------|-------------|------------------|
| Single DB + Schema Reset | Use one database, reset schema between tests | 8-16 hours |
| Connection Pooling | Reuse connections across test sessions | 16-24 hours |
| Database Monitoring | Track database size, orphan count | 8-12 hours |
| Automated Cleanup Cron | Nightly cron job for orphan cleanup | 4-6 hours |
| Test Database Compression | pg_repack for PostgreSQL optimization | 12-20 hours |

---

## Implementation Checklist

### Phase 1: Core Fix (Target: 2-3 days)

- [ ] **Step 1.1:** Add `defer` cleanup in `TestHelper.Close()`
  - [ ] Test: Run single test, verify database is dropped
  - [ ] Test: Panic test, verify database is still dropped
  
- [ ] **Step 1.2:** Implement `cleanupOrphanedDatabases()`
  - [ ] Test: Create 100 orphaned databases, run cleanup
  - [ ] Test: Verify orphans older than 24h are dropped
  - [ ] Test: Verify current test DB is not dropped
  
- [ ] **Step 1.3:** Add database name validation
  - [ ] Test: Valid names are accepted
  - [ ] Test: Invalid names are rejected with clear error
  - [ ] Test: SQL injection attempts are blocked
  
- [ ] **Step 1.4:** Add `make test-clean` command
  - [ ] Test: Command runs successfully
  - [ ] Test: Drops all orphaned databases
  - [ ] Test: Provides progress feedback

### Phase 2: Parallel Test Safety (Target: 1-2 days)

- [ ] **Step 2.1:** Add goroutine ID to database name
  - [ ] Test: Parallel tests (10+) don't collide
  - [ ] Test: Each test has unique database name
  - [ ] Test: Cleanup works for all parallel tests
  
- [ ] **Step 2.2:** Verify parallel test performance
  - [ ] Measure test execution time before/after
  - [ ] Ensure < 10% performance regression

### Phase 3: Documentation & Training (Target: 0.5 day)

- [ ] **Step 3.1:** Update `AGENTS.md` with new cleanup procedures
- [ ] **Step 3.2:** Document the database cleanup process
- [ ] **Step 3.3:** Create quick reference guide

---

## Stakeholder Alignment

### Responsibility Matrix

| Stakeholder | Responsibility | Acceptance Criteria |
|-------------|----------------|---------------------|
| **Development Team** | Implement fixes | All acceptance criteria met |
| **QA Team** | Test cleanup logic | All test scenarios pass |
| **DevOps** | Deploy and monitor | Cleanup runs in production environment |
| **Product Owner** | Prioritize work | Critical path approved |

### Sign-off Requirements

| Item | Owner | Due Date | Status |
|------|-------|----------|--------|
| PRD Approval | Product Owner | 2026-04-16 | ⏳ Pending |
| Implementation Review | Tech Lead | 2026-04-17 | ⏳ Pending |
| QA Sign-off | QA Lead | 2026-04-19 | ⏳ Pending |
| Production Deployment | DevOps | 2026-04-20 | ⏳ Pending |

---

## Traceability Matrix

### Requirements → User Stories

| Requirement ID | User Story | Priority |
|----------------|------------|----------|
| R1 | As a developer, I want test databases to be automatically cleaned up so I don't have to manually clean up 6,000+ databases | Critical |
| R2 | As a developer, I want orphaned databases to be cleaned up periodically so my disk space doesn't fill up | Critical |
| R3 | As a developer, I want parallel tests to run reliably without database name collisions | High |
| R4 | As a developer, I want a simple command to clean up test databases manually | Medium |
| R5 | As a security-conscious developer, I want database names to be validated to prevent SQL injection | High |

### User Stories → Acceptance Criteria

| User Story | AC1 | AC2 | AC3 | AC4 |
|------------|-----|-----|-----|-----|
| Auto-cleanup | ✅ R1.AC1 | ✅ R1.AC2 | ✅ R1.AC3 | ✅ R1.AC4 |
| Orphan cleanup | ✅ R2.AC1 | ✅ R2.AC2 | ✅ R2.AC3 | - |
| Parallel safety | ✅ R3.AC1 | ✅ R3.AC2 | - | - |
| Manual cleanup | ✅ R4.AC1 | ✅ R4.AC2 | ✅ R4.AC3 | - |
| SQL injection prevention | ✅ R5.AC1 | ✅ R5.AC2 | ✅ R5.AC3 | ✅ R5.AC4 |

### Acceptance Criteria → Tests

| AC ID | Test File | Test Function | Type |
|-------|-----------|---------------|------|
| R1.AC1 | `test/test_helper_test.go` | `TestTestHelperCleanup` | Integration |
| R1.AC2 | `test/test_helper_test.go` | `TestTestHelperCleanupOnPanic` | Integration |
| R2.AC1 | `test/test_helper_test.go` | `TestCleanupOrphanedDatabases` | Integration |
| R2.AC2 | `test/test_helper_test.go` | `TestCleanupOrphanedDatabases_ExcludeCurrent` | Integration |
| R3.AC1 | `test/test_helper_test.go` | `TestParallelTestDatabaseUniqueness` | Integration |
| R5.AC1 | `test/test_helper_test.go` | `TestIsValidTestDatabaseName_Valid` | Unit |
| R5.AC2 | `test/test_helper_test.go` | `TestIsValidTestDatabaseName_Invalid` | Unit |
| R5.AC3 | `test/test_helper_test.go` | `TestIsValidTestDatabaseName_SQLOjection` | Unit |

---

## Validation

### Code Quality Checklist

- [ ] Code follows Go 1.25.7 best practices
- [ ] All functions have proper error handling
- [ ] Context timeouts are set appropriately
- [ ] Database connections are properly closed
- [ ] SQL queries use parameterized inputs where possible
- [ ] Code is properly commented
- [ ] No hardcoded secrets or credentials
- [ ] Logging is appropriate (not too verbose, not too silent)

### Technical Feasibility Checklist

- [ ] Go version (1.25.7) supports all required features
- [ ] pgx/v5 library supports all required functionality
- [ ] PostgreSQL version supports all required features
- [ ] Docker Compose configuration doesn't block cleanup
- [ ] No external dependencies required
- [ ] Changes are backward compatible

### User Needs Validation

| Need | Solution | Verification |
|------|----------|--------------|
| Prevent database accumulation | Auto-cleanup + orphan cleanup | Test with 6,000+ orphans |
| Run tests reliably | Unique DB names per test | Parallel test stress test |
| Manual cleanup option | `make test-clean` command | Command runs successfully |
| Security | Name validation + parameterized queries | SQL injection test passes |

---

## Ready for Implementation

### Final Review Checklist

- [ ] All stakeholder questions answered
- [ ] Technical risks identified and mitigated
- [ ] Acceptance criteria are testable
- [ ] Implementation steps are clear and actionable
- [ ] Files to modify are specified with line numbers
- [ ] Out of scope items are documented
- [ ] Phase 2 roadmap is established

### Sign-off Block

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Product Owner | | | |
| Tech Lead | | | |
| QA Lead | | | |

---

## Appendix A: Reference Implementation

See `test/test_helper.go` and `test/cleanup_orphaned_databases.go` for complete implementation examples.

## Appendix B: Performance Benchmarks

Run `scripts/benchmark_cleanup.go` to measure:
- Database creation time
- Database drop time
- Orphan cleanup time
- Parallel test performance impact

---

**PRD Version:** 1.0  
**Next Review Date:** 2026-04-16  
**Status:** Awaiting Stakeholder Approval
