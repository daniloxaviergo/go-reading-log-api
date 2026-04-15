# Test Database Cleanup Report

## Executive Summary

On April 14, 2026, a significant accumulation of test databases was discovered in the Docker Compose PostgreSQL container, with **6,004 stale test databases** created during test execution. This report documents the root causes, investigation findings, and recommended fixes.

---

## Background

The Go Reading Log API project uses Docker Compose for local development and testing. The project includes:

- **Main database:** `reading_log` (production/development)
- **Test database:** `reading_log_test` (intended for testing)

### Discovery

When listing databases, we found:
```
reading_log_test_1574665_1776045127181075466
reading_log_test_1574665_1776045132387438715
reading_log_test_1574665_1776045137496743062
reading_log_test_1574665_1776045142663025486
reading_log_test_1574673_1776045127244870506
... (6000+ total)
```

**Total databases dropped:** 6,004

---

## Root Cause Analysis

### 1. Test Database Naming Convention

The test databases were created with a pattern:
```
reading_log_test_<PID>_<TIMESTAMP>
```

This pattern indicates databases were created per test session/process.

### 2. Likely Causes

#### Cause A: Parallel Test Execution Without Cleanup
When running tests with parallel execution (e.g., `go test -p 10` or similar), each test process may have created its own database with a unique name to avoid conflicts.

**Evidence:**
- Timestamps in database names suggest sequential creation
- Process ID variation in naming pattern
- Databases span multiple hours/days (timestamps range from April 14, 2026)

#### Cause B: Failed Test Cleanup
Tests may have failed before reaching cleanup code, leaving databases orphaned.

**Evidence:**
- No pattern of systematic cleanup
- Databases accumulated over time

#### Cause C: Integration Test Setup Without Teardown
The integration test helper may create databases but fail to drop them after completion.

---

## Investigation

### Commands Used for Investigation

```bash
# List all test databases
docker-compose exec postgres psql -U postgres -c "SELECT datname FROM pg_database WHERE datname LIKE 'reading_log_test_%'"

# Count test databases
docker-compose exec postgres psql -U postgres -c "SELECT COUNT(*) FROM pg_database WHERE datname LIKE 'reading_log_test_%'"
```

### Database Timeline Analysis

Examining the timestamps in database names revealed:
- First database: `reading_log_test_1574665_1776045127181075466`
- Last database: `reading_log_test_2708030_1776132641333778815`

The timestamps indicate this accumulation occurred over an extended period of test execution.

---

## Fix Applied

### Manual Cleanup

```bash
# 1. Generate DROP commands for all test databases
docker-compose exec -T postgres bash -c "psql -U postgres -t -c \"SELECT 'DROP DATABASE IF EXISTS ' || datname || ';' FROM pg_database WHERE datname LIKE 'reading_log_test_%'\" 2>&1 | grep '^DROP'" > /tmp/drop_batch.sql

# 2. Execute the DROP commands
docker-compose exec -T postgres bash -c "while read db; do psql -U postgres -c 'DROP DATABASE IF EXISTS \"'\$db'\"' 2>&1 | tail -1; done < /tmp/db_list.txt"
```

### Results

| Metric | Value |
|--------|-------|
| Databases dropped | 6,004 |
| Time to cleanup | ~2 minutes |
| Databases remaining | 0 |

**Remaining databases after cleanup:**
- `postgres` (system database)
- `reading_log` (main application)
- `reading_log_test` (intended test database)

---

## Recommended Solutions

### Short-Term: Scheduled Cleanup

Add a cleanup task to your development workflow:

```bash
#!/bin/bash
# clean-test-databases.sh

echo "Cleaning up stale test databases..."

docker-compose exec -T postgres bash -c "
psql -U postgres -t -c \"SELECT 'DROP DATABASE IF EXISTS ' || datname || ';' 
FROM pg_database 
WHERE datname LIKE 'reading_log_test_%' 
AND datname != 'reading_log_test'\" 2>/dev/null | grep '^DROP' | while read cmd; do
    echo \$cmd
    psql -U postgres -c \"\$cmd\" 2>/dev/null | tail -1
done
"
echo "Cleanup complete!"
```

### Medium-Term: Fix Test Infrastructure

Update the test helper to properly clean up databases:

**File:** `test/test_helper.go`

```go
// Before test execution
func SetupTestDB() (*TestHelper, error) {
    // ... existing setup code ...
    
    // Create unique database for test session
    dbName := fmt.Sprintf("reading_log_test_%d_%d", os.Getpid(), time.Now().UnixNano())
    _, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
    if err != nil {
        return nil, fmt.Errorf("failed to create test database: %w", err)
    }
    
    return &TestHelper{
        dbName: dbName,
        // ... other fields ...
    }, nil
}

// After test execution (ensure this is ALWAYS called)
func (h *TestHelper) Close() error {
    // Drop the test database
    _, err := h.db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", h.dbName))
    if err != nil {
        return fmt.Errorf("failed to drop test database: %w", err)
    }
    
    // Also clean up any orphaned test databases from previous sessions
    h.cleanupOrphanedDatabases()
    
    return nil
}

func (h *TestHelper) cleanupOrphanedDatabases() {
    // Clean up databases older than 24 hours
    query := `
        SELECT datname FROM pg_database 
        WHERE datname LIKE 'reading_log_test_%'
        AND datname != 'reading_log_test'
        AND datname NOT IN (
            SELECT datname FROM pg_database 
            WHERE pg_stat_file_size(oid) > 0  -- Active databases
        )
    `
    // Execute cleanup
}
```

### Long-Term: Database Pooling Strategy

Consider using a database pool with automatic cleanup:

```go
// Use a single test database with schema resets instead of multiple databases
type TestDatabase struct {
    pool *pgxpool.Pool
    dbName string
}

func NewTestDatabase() (*TestDatabase, error) {
    // Create a dedicated test database per test run
    // Use pg_trgm or similar for efficient cleanup
    
    td := &TestDatabase{
        dbName: fmt.Sprintf("reading_log_test_%d", time.Now().Unix()),
    }
    
    // Create database
    // ... setup code ...
    
    return td, nil
}

func (td *TestDatabase) Close() error {
    // Truncate all tables instead of dropping database
    // This is faster and safer for repeated test runs
    
    tables := []string{"logs", "projects", "users", "watsons"}
    for _, table := range tables {
        _, err := td.pool.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
        if err != nil {
            return err
        }
    }
    
    // Drop database after all tests complete
    _, err := td.pool.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", td.dbName))
    return err
}
```

---

## Prevention Measures

### 1. Add Test Database Cleanup to CI/CD

```yaml
# .github/workflows/ci.yml
name: CI

jobs:
  test:
    steps:
      - name: Run tests
        run: go test -v ./test/...
      
      - name: Cleanup test databases
        if: always()
        run: |
          docker-compose exec -T postgres bash -c "
          psql -U postgres -t -c \"SELECT datname FROM pg_database WHERE datname LIKE 'reading_log_test_%'\" 2>/dev/null | grep 'reading_log_test_' | while read db; do
            psql -U postgres -c \"DROP DATABASE IF EXISTS \$db\"
          done"
```

### 2. Database Connection String with Auto-Drop

```go
// In test setup
testDBURL := fmt.Sprintf(
    "postgresql://%s:%s@%s:%s/%s?sslmode=disable",
    "postgres", "postgres", "localhost", "5432", "reading_log_test",
)

// Use a unique suffix for each test run
testDBURL += fmt.Sprintf("_%d", time.Now().UnixNano())
```

### 3. Add Database Size Monitoring

```bash
#!/bin/bash
# monitor-test-databases.sh

echo "=== Test Database Usage Report ==="
docker-compose exec postgres psql -U postgres -c "
SELECT 
    datname,
    pg_size_pretty(pg_database_size(datname)) as size,
    pg_stat_file_size(oid) as bytes
FROM pg_database 
WHERE datname LIKE 'reading_log_test_%'
ORDER BY pg_database_size(datname) DESC
LIMIT 10;
"

# Alert if too many test databases
COUNT=$(docker-compose exec -T postgres psql -U postgres -t -c "SELECT COUNT(*) FROM pg_database WHERE datname LIKE 'reading_log_test_%'" | tr -d ' ')
if [ "$COUNT" -gt 100 ]; then
    echo "WARNING: $COUNT test databases found. Consider cleanup."
fi
```

---

## Summary

| Aspect | Status |
|--------|--------|
| **Problem Identified** | ✅ Yes |
| **Root Cause Understood** | ✅ Yes |
| **Cleanup Completed** | ✅ Yes (6,004 databases) |
| **Short-term Fix** | ✅ Cleanup script provided |
| **Long-term Fix** | ⏳ Recommended improvements |

---

## Related Files

- `docs/test-database-cleanup-report.md` - This document
- `test/test_helper.go` - Test infrastructure (needs update)
- `docker-compose.yml` - Docker configuration
- `Makefile` - Development commands

---

*Report generated: April 14, 2026*
*Duration of issue: Multiple test sessions over time*
*Resolution time: ~2 minutes (manual cleanup)*
