---
id: RDL-049
title: '[doc-004 Phase 1.2] Implement orphaned database cleanup function'
status: To Do
assignee: []
created_date: '2026-04-15 12:14'
updated_date: '2026-04-16 00:02'
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
