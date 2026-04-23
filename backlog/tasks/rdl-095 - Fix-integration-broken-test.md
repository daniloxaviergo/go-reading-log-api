---
id: RDL-095
title: Fix integration broken test
status: To Do
assignee:
  - thomas
created_date: '2026-04-23 12:16'
updated_date: '2026-04-23 12:20'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
When execute the test broken `go test -v -timeout=5s ./...`

--- 
The integration tests are failing with two main issues:

### Error 1: Missing Database Tables
```
ERROR: relation "projects" does not exist (SQLSTATE 42P01)
```

### Error 2: Test Timeouts
Tests timing out after 5 seconds, likely stuck waiting for database operations.

## Root Cause

The integration test files in `/home/danilo/scripts/github/go-reading-log-api-next/test/` are missing the call to `SetupTestSchema()` after creating the test database connection.

### Current (Broken) Pattern:
```go
helper, err := SetupTestDB()
require.NoError(t, err)
defer helper.Close()

// Missing: helper.SetupTestSchema() - tables not created!
```

### Correct Pattern (from unit tests):
```go
helper, err := SetupTestDB()
require.NoError(t, err)
defer helper.Close()

// Create tables before inserting data
err = helper.SetupTestSchema()
require.NoError(t, err)
```

## Solution

Add `helper.SetupTestSchema()` calls after `SetupTestDB()` in all integration test files.

### Files to Fix:

1. **test/dashboard_integration_test.go**
   - `TestDashboardDayEndpoint_Integration`
   - `TestDashboardProjectsEndpoint_Integration`
   - `TestDashboardLastDaysEndpoint_Integration`
   - `TestDashboardFaultsChart_Integration`
   - `TestDashboardSpeculateActual_Integration`
   - `TestDashboardWeekdayFaults_Integration`
   - `TestDashboardMeanProgress_Integration`
   - `TestDashboardYearlyTotal_Integration`
   - `TestDashboardEndpoints_ErrorHandling`

2. **test/integration/error_scenarios_test.go**
   - All test functions using `SetupTestDB()`

3. **Any other integration test files** that use `SetupTestDB()`

## Implementation

For each test function, add the schema setup immediately after `SetupTestDB()`:

```go
func TestYourIntegrationTest(t *testing.T) {
    helper, err := SetupTestDB()
    require.NoError(t, err)
    defer helper.Close()
    
    // ADD THIS LINE - Create database tables
    err = helper.SetupTestSchema()
    require.NoError(t, err)
    
    // ... rest of test code
}
```

## Verification

After making changes, run tests again:

```bash
go test -v -timeout=10s ./test/...
```

Expected result: All tests should pass without "relation does not exist" errors.

## Additional Notes

### Why This Happened
The `SetupTestDB()` function only creates the database connection and the database itself. It does NOT create the tables. The `SetupTestSchema()` method must be explicitly called to execute the SQL DDL statements that create:
- `projects` table
- `logs` table
- Indexes

### Timeout Issue
The 5-second timeout may also be caused by:
1. Database connection pool waiting for available connections
2. Schema creation taking time on first run
3. Network latency to PostgreSQL

Consider increasing the test timeout from 5s to 10s or 15s in `go test` command.

### Consistency with Unit Tests
The unit tests in `test/unit/dashboard_repository_test.go` already follow this pattern correctly, so use them as a reference implementation.
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The integration tests are failing because the `SetupTestDB()` function creates a database connection but does NOT create the required database tables. The tests need to explicitly call `SetupTestSchema()` after establishing the database connection.

**Root Cause Analysis:**
- `SetupTestDB()` creates the test database and connection pool
- `SetupTestSchema()` executes SQL DDL statements to create tables (projects, logs)
- Integration tests are missing the `SetupTestSchema()` call
- This causes "relation 'projects' does not exist" errors when queries execute

**Architecture Decision:**
- Follow the existing pattern used in unit tests (`dashboard_repository_test.go`)
- Add `helper.SetupTestSchema()` immediately after `SetupTestDB()` in all integration test files
- Increase test timeout from 5s to 10s to accommodate schema creation overhead

**Why This Approach:**
- Minimal code changes - just one line per test function
- Consistent with existing test infrastructure
- No architectural modifications needed
- Low risk, high impact fix

---

### 2. Files to Modify

#### Primary Fix Files:

1. **test/dashboard_integration_test.go**
   - `TestDashboardDayEndpoint_Integration`
   - `TestDashboardProjectsEndpoint_Integration`
   - `TestDashboardLastDaysEndpoint_Integration`
   - `TestDashboardFaultsChart_Integration`
   - `TestDashboardSpeculateActual_Integration`
   - `TestDashboardWeekdayFaults_Integration`
   - `TestDashboardMeanProgress_Integration`
   - `TestDashboardYearlyTotal_Integration`
   - `TestDashboardEndpoints_ErrorHandling`

2. **test/integration/error_scenarios_test.go**
   - `RunErrorScenarios` function
   - All test functions using `SetupTestDB()`

#### Verification Files (read-only for reference):

3. **test/unit/dashboard_repository_test.go** (reference pattern)
   - Shows correct usage of `SetupTestSchema()`
   - Use as template for fixing broken tests

4. **test/test_helper.go** (reference implementation)
   - Contains `SetupTestSchema()` method definition
   - Contains SQL DDL statements for table creation

---

### 3. Dependencies

**Prerequisites:**
- PostgreSQL must be running and accessible
- Test database user must have CREATE TABLE permissions
- No external dependencies required (uses existing test infrastructure)

**Blocking Issues:**
None - this is a straightforward code fix

**Setup Steps Required:**
```bash
# Ensure PostgreSQL is running
pg_isready -h localhost -p 5432

# Run tests to verify fix
go test -v -timeout=10s ./test/...
```

---

### 4. Code Patterns

**Pattern to Follow (from unit tests):**

```go
func TestYourIntegrationTest(t *testing.T) {
    if !IsTestDatabase() {
        t.Skip("Test database not configured")
    }

    helper, err := SetupTestDB()
    require.NoError(t, err)
    defer helper.Close()
    
    // CRITICAL: Create tables before inserting data
    err = helper.SetupTestSchema()
    require.NoError(t, err)
    
    // ... rest of test code (load fixtures, create handlers, etc.)
}
```

**Key Conventions:**
1. Call `SetupTestSchema()` immediately after `SetupTestDB()`
2. Use `require.NoError(t, err)` to fail fast if schema creation fails
3. The `defer helper.Close()` will handle cleanup automatically
4. Keep the same timeout pattern (30 seconds for context operations)

**Files Following This Pattern:**
- `test/unit/dashboard_repository_test.go` - Already correct, use as reference
- All new integration tests should follow this pattern

---

### 5. Testing Strategy

**Unit Testing:**
- Each test function will be verified independently
- No mock changes needed - using real database via `SetupTestDB()`

**Integration Testing:**
- Test actual database schema creation
- Verify table existence before queries execute
- Confirm fixture data loads correctly into created tables

**Edge Cases to Cover:**
1. **Empty database**: Tests should handle empty tables gracefully
2. **Concurrent tests**: Each test uses unique database name (via `getGoroutineID()`)
3. **Cleanup**: `defer helper.Close()` ensures tables are dropped after each test

**Verification Commands:**
```bash
# Run all integration tests
go test -v -timeout=10s ./test/...

# Run specific test to verify fix
go test -v -timeout=10s ./test/dashboard_integration_test.go

# Check for "relation does not exist" errors (should be zero)
go test -v ./test/... 2>&1 | grep -c "relation.*does not exist"
```

**Expected Results:**
- Zero "relation does not exist" errors
- All tests complete within timeout
- Test output shows passing tests with ✓

---

### 6. Risks and Considerations

**Low Risk Changes:**
- ✅ Minimal code modification (one line per test)
- ✅ No changes to production code
- ✅ No architectural changes
- ✅ Follows existing patterns in unit tests

**Potential Issues:**

1. **Test Timeout**: 
   - Current timeout: 5 seconds
   - Recommended: 10-15 seconds
   - Reason: Schema creation adds overhead on first run
   
2. **Database Permissions**:
   - Test user needs CREATE TABLE privileges
   - Currently using `postgres` user which has full permissions

3. **Parallel Test Isolation**:
   - Each test gets unique database name via `getGoroutineID()`
   - No risk of table name conflicts between parallel tests

**Deployment Considerations:**
- This fix only affects test infrastructure
- No production database changes required
- Can be deployed immediately without migration

**Rollback Plan:**
If issues arise, simply revert the one-line changes in each test file. The original broken state will be restored.

---

### Implementation Checklist

Before running tests, ensure:

- [ ] All `SetupTestDB()` calls are followed by `SetupTestSchema()`
- [ ] Test timeout increased from 5s to 10s
- [ ] PostgreSQL is running and accessible
- [ ] Test database user has CREATE permissions

After implementing, verify:

- [ ] No "relation does not exist" errors
- [ ] All tests pass with ✓ prefix
- [ ] Test execution completes within timeout
- [ ] Cleanup runs successfully (no orphaned databases)
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
