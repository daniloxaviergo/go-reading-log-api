---
id: RDL-095
title: Fix integration broken test
status: To Do
assignee: []
created_date: '2026-04-23 12:16'
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
