---
id: RDL-049
title: '[doc-004 Phase 1.2] Implement orphaned database cleanup function'
status: To Do
assignee:
  - catarina
created_date: '2026-04-15 12:14'
updated_date: '2026-04-16 00:05'
labels:
  - cleanup
  - infrastructure
  - critical
dependencies: []
references:
  - 'R2: Orphaned Database Cleanup'
documentation:
  - doc-004
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement the cleanupOrphanedDatabases() function in test/test_helper.go to identify and drop test databases that are older than 24 hours. The function should query pg_database for databases matching the pattern reading_log_test_%, exclude the current test database, and drop each identified orphan. The cleanup must complete in under 1 minute for 6,000+ databases, exclude the current test database, log errors without failing test execution, and use context timeouts to prevent indefinite blocking.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Databases older than 24 hours are identified and dropped
- [ ] #2 Current test database is excluded from cleanup
- [ ] #3 Cleanup runs in under 1 minute for 6,000+ databases
- [ ] #4 Errors are logged but don't fail test execution
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The task requires implementing the `cleanupOrphanedDatabases()` function in `test/test_helper.go` to identify and drop test databases that are older than 24 hours. The function should:

- Query `pg_database` for databases matching the pattern `reading_log_test_%`
- Exclude the current test database from cleanup
- Drop each identified orphan database
- Complete within 1 minute for 6,000+ databases
- Log errors without failing test execution
- Use context timeouts to prevent indefinite blocking

**Architecture Decisions:**
- Use a separate connection pool to query the main database for orphaned databases
- Implement batch cleanup to process multiple databases efficiently
- Use context timeouts (60 seconds) to prevent indefinite blocking
- Log errors at warning level but continue cleanup of other databases
- Use `DROP DATABASE IF EXISTS` for safe deletion

**Why this approach:**
- 6,000+ database cleanup requires efficient batch processing
- Separate connection pool ensures we can query even if test pool is closed
- 60-second timeout balances thoroughness with speed
- Error logging without failure ensures cleanup doesn't break tests
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
