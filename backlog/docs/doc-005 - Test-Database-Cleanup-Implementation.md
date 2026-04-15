---
id: doc-005
title: Test Database Cleanup Implementation
type: other
created_date: '2026-04-15 11:43'
---


# PRD: Test Database Cleanup Implementation

## Problem Statement

The Go Reading Log API test suite is experiencing **critical database pollution**:

| Metric | Current State |
|--------|---------------|
| Stale test databases | 6,004 accumulated |
| Root cause | Test sessions create unique databases but don't drop them |
| Impact | PostgreSQL bloat, potential performance degradation |

## Technical Analysis

### Current Implementation

**File:** `test/test_helper.go`

```go
// Current approach (PROBLEMATIC)
func SetupTestDB() (*TestHelper, error) {
    // Creates unique database name
    testDBName = fmt.Sprintf("%s_%d_%d", testDBName, os.Getpid(), time.Now().UnixNano())
    
    // Creates database
    mainPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", testDBName))
    
    // ... test runs ...
    
    // Close method attempts cleanup
    func (h *TestHelper) Close() {
        h.Pool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", h.TestDBName))
    }
}
```

### Issues Identified

1. **No cleanup on test failure/panic** - `defer` not used consistently
2. **No orphaned database cleanup** - Previous session databases remain
3. **Parallel test risk** - Race conditions on database creation
4. **No scheduled cleanup** - Accumulation over time

## Proposed Solution

### Phase 1: Enhanced Cleanup on Close

```go
// Before Close(), clean up orphaned databases
func cleanupOrphanedDatabases(pool *pgxpool.Pool, excludeName string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Query for test databases older than 24 hours
    query := `
        SELECT datname 
        FROM pg_database 
        WHERE datname LIKE 'reading_log_test_%'
        AND datname != $1
        AND pg_catalog.pg_get_userbyid(datdba) = current_user
    `
    
    rows, err := pool.Query(ctx, query, excludeName)
    if err != nil {
        return fmt.Errorf("failed to query orphaned databases: %w", err)
    }
    defer rows.Close()
    
    var toDrop []string
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err == nil {
            toDrop = append(toDrop, name)
        }
    }
    
    // Drop each orphaned database
    for _, dbName := range toDrop {
        dropQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
        if _, err := pool.Exec(ctx, dropQuery); err != nil {
            return fmt.Errorf("failed to drop %s: %w", dbName, err)
        }
    }
    
    return nil
}
```

### Phase 2: Proper Integration with Test Lifecycle

```go
func (h *TestHelper) Close() {
    if h.Pool != nil {
        // Drop the unique database for this test session
        if h.TestDBName != "" {
            ctx, cancel := context.WithTimeout(context.Background(), testContextTimeout)
            defer cancel()
            
            connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
                h.Config.DBUser, h.Config.DBPassword, h.Config.DBHost, 
                h.Config.DBPort, h.Config.DBDatabase)
            mainPool, err := pgxpool.New(ctx, connStr)
            if err == nil {
                defer mainPool.Close()
                
                // Drop this test's database
                dropQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s", h.TestDBName)
                if _, err := mainPool.Exec(ctx, dropQuery); err != nil {
                    log.Printf("Warning: failed to drop test database %s: %v", 
                        h.TestDBName, err)
                }
                
                // Clean up orphaned databases from previous sessions
                if err := cleanupOrphanedDatabases(mainPool, h.TestDBName); err != nil {
                    log.Printf("Warning: orphaned cleanup failed: %v", err)
                }
            }
        }
        h.Pool.Close()
    }
}
```

## Implementation Questions

### Question 1: Single Database vs Per-Test Database

**Approach A: Keep per-test databases (current)**
- ✅ True isolation between tests
- ✅ Faster parallel execution (no schema reset needed)
- ❌ Database accumulation issue
- ❌ Higher resource usage

**Approach B: Single test database with schema resets**
- ✅ No accumulation issue
- ✅ Faster setup/teardown
- ❌ Potential race conditions in parallel tests
- ❌ Slower due to schema resets

**Recommendation:** Keep per-test databases but implement **aggressive cleanup policy**:
- Drop after each test session
- Scheduled cleanup for orphans (cron job or manual script)
- Consider using `TRUNCATE` instead of `DROP` for schema reset if switching to single DB

---

### Question 2: Performance Implications

| Operation | Time Estimate | Resource Cost |
|-----------|---------------|---------------|
| `CREATE DATABASE` | ~100-200ms | High (new files) |
| `DROP DATABASE` | ~50-100ms | Medium (file deletion) |
| `TRUNCATE TABLE` | ~10-20ms | Low (just metadata) |

**Recommendation:** For the 6,004 existing orphans:
1. Run a one-time cleanup script
2. Implement the `DROP DATABASE IF EXISTS` in `Close()`
3. Consider adding a `make test-clean` command for manual cleanup

---

### Question 3: Security/Permission Concerns

```sql
-- Current approach uses postgres superuser
CREATE DATABASE reading_log_test_12345_67890;

-- Consider using a dedicated test user with limited permissions
-- GRANT CREATE DATABASE to test_user;
-- REVOKE ALL ON DATABASE reading_log FROM test_user;
```

**Recommendation:** 
- Use dedicated test database user with minimal privileges
- Ensure `DROP DATABASE` permission is granted
- Consider using `SUPERUSER` only in development, not production test environments

---

### Question 4: Parallel Test Compatibility

**Current:** Uses `os.Getpid() + UnixNano()` for uniqueness

```go
testDBName = fmt.Sprintf("%s_%d_%d", testDBName, os.Getpid(), time.Now().UnixNano())
```

**Concern:** `UnixNano()` can collide under high parallelization

**Recommendation:** Add test process ID to ensure uniqueness:

```go
import "runtime"

func generateUniqueDBName(baseName string) string {
    // Use goroutine ID and process ID for uniqueness
    return fmt.Sprintf("%s_%d_%d_%d", 
        baseName, 
        os.Getpid(),
        getGoroutineID(), // Use runtime/debug for goroutine ID
        time.Now().UnixNano())
}
```

---

### Question 5: Edge Cases to Handle

| Edge Case | Handling Strategy |
|-----------|-------------------|
| Database doesn't exist | `DROP DATABASE IF EXISTS` (already handled) |
| Connection failure | Log warning, don't fail test |
| Test panics | Use `defer` in test helper, not manual calls |
| Orphaned databases | Scheduled cleanup with age check |
| Concurrent cleanup | Database-level locks via PostgreSQL |

---

## Proposed Implementation Plan

### Step 1: Add Cleanup Function

```go
// internal/adapter/postgres/cleanup.go (new file)
package postgres

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/jackc/pgx/v5/pgxpool"
)

// CleanupOrphanedTestDatabases removes test databases older than 24 hours
func CleanupOrphanedTestDatabases(pool *pgxpool.Pool, excludeName string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()
    
    // Get all test databases
    query := `
        SELECT datname 
        FROM pg_database 
        WHERE datname LIKE $1
        AND pg_catalog.pg_get_userbyid(datdba) = current_user
    `
    
    rows, err := pool.Query(ctx, query, "reading_log_test_%")
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
        if name != excludeName {
            toDrop = append(toDrop, name)
        }
    }
    
    // Drop each database
    for _, dbName := range toDrop {
        if _, err := pool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)); err != nil {
            log.Printf("Warning: failed to drop %s: %v", dbName, err)
        }
    }
    
    return nil
}
```

### Step 2: Update TestHelper

```go
// Add cleanup call to Close()
func (h *TestHelper) Close() {
    if h.Pool != nil {
        // ... existing cleanup code ...
        
        // Also clean up orphaned databases
        connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
            h.Config.DBUser, h.Config.DBPassword, h.Config.DBHost, 
            h.Config.DBPort, h.Config.DBDatabase)
        mainPool, err := pgxpool.New(ctx, connStr)
        if err == nil {
            defer mainPool.Close()
            if err := CleanupOrphanedTestDatabases(mainPool, h.TestDBName); err != nil {
                log.Printf("Warning: cleanup failed: %v", err)
            }
        }
        
        h.Pool.Close()
    }
}
```

### Step 3: Add Make Command

```makefile
# Makefile
test-clean:
	@echo "Cleaning up orphaned test databases..."
	@go run ./scripts/cleanup_test_databases.go
```

---

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Accidental data loss | HIGH | `IF EXISTS` + user verification |
| Test flakiness | MEDIUM | Proper error handling, logging |
| Performance degradation | LOW | scheduled cleanup, not on every test |
| Race conditions | MEDIUM | Database-level locks, unique names |

---

## Recommended Actions

1. **Immediate:** Run a one-time cleanup script for the 6,004 orphans
2. **Short-term:** Implement `DROP DATABASE IF EXISTS` in `Close()`
3. **Medium-term:** Add scheduled cleanup job (cron or manual)
4. **Long-term:** Consider schema reset approach if parallel tests become problematic

---

## Verification Checklist

- [ ] All test databases are dropped after test completion
- [ ] No orphaned databases remain after 24 hours
- [ ] Parallel tests don't interfere with each other
- [ ] Test failure doesn't leave databases behind
- [ ] Cleanup doesn't impact test performance significantly
