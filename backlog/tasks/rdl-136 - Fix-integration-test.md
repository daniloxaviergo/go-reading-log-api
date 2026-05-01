---
id: RDL-136
title: Fix integration test
status: To Do
assignee:
  - Thomas
created_date: '2026-05-01 11:54'
updated_date: '2026-05-01 11:55'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Analyze the provided Go codebase to resolve integration test timeouts caused by database concurrency issues and connection pool deadlocks during cleanup. The goal is to ensure tests complete reliably within the timeout limits without hanging on database operations.

Step 1: Analyze the Stack Trace
- Examine the panic message indicating a 2s timeout during TestDashboardDayEndpoint_Integration.
- Identify the blocking goroutines (e.g., goroutine 21 waiting on IO, goroutine 39 running the test).
- Pinpoint the specific functions causing the hang (e.g., pgxpool.(*Pool).backgroundHealthCheck, test.cleanupOrphanedDatabases).

Step 2: Identify Concurrency Bottlenecks
- Review the pgx/v5 connection pool configuration used in the test helper.
- Determine if connection exhaustion or lock contention is occurring during test setup or teardown.
- Check for improper context cancellation or missing timeout configurations on database queries.

Step 3: Evaluate Test Lifecycle Management
- Inspect the TestHelper.Close method and cleanupOrphanedDatabases function.
- Identify if resources are being closed while other goroutines are still attempting to access the pool.
- Look for race conditions between the background health check and the test execution.

Step 4: Formulate Corrective Actions
- Propose specific code modifications to the test helper and pool configuration.
- Suggest adjustments to MaxConns, ConnMaxLifetime, or health check intervals to prevent blocking.
- Recommend synchronization mechanisms (e.g., mutexes, wait groups) to ensure safe cleanup.

Step 5: Verify the Solution
- Explain how the proposed changes eliminate the race condition and prevent future timeouts.
- Ensure the solution maintains test isolation and performance.
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Task Instructions
1.  **Analyze `test/test_helper.go`:**
    - Inspect the `Close` method and `cleanupOrphanedDatabases` function.
    - Verify if `context.WithTimeout` is used for all database operations during cleanup.
    - Check if `pool.Close()` is called correctly to prevent blocking on active connections.
2.  **Review Connection Pool Configuration:**
    - Examine how the `pgxpool.Config` is set up in the test helper.
    - Ensure `MaxConns` and `MinConns` are appropriate for the test environment to prevent contention.
    - Verify if `MaxConnIdleTime` or `MaxConnLifetime` are contributing to stale connection issues.
3.  **Check Test Concurrency:**
    - Identify if `t.Parallel()` is used in `dashboard_integration_test.go`.
    - Ensure tests do not share a single global database connection pool if they run in parallel (use per-test DB instances or distinct schemas if necessary).
    - Look for shared state that might cause locks during the `cleanupOrphanedDatabases` phase.
4.  **Implement Fixes:**
    - Refactor `cleanupOrphanedDatabases` to use a strict context timeout (e.g., 500ms) to prevent indefinite hangs.
    - Ensure `pool.Close()` is handled gracefully, potentially using `pool.Stat()` to verify drains or forcing connection closures if stuck.
    - Add proper error handling and logging to identify exactly where the hang occurs.
5.  **Validation:**
    - Provide the corrected code snippets.
    - Explain why the changes resolve the concurrency deadlock.
    - Suggest a command to run tests with increased verbosity (`-v`) or timeout (`-timeout=5m`) to verify the fix.
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
