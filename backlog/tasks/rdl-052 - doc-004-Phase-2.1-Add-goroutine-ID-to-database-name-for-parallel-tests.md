---
id: RDL-052
title: '[doc-004 Phase 2.1] Add goroutine ID to database name for parallel tests'
status: To Do
assignee:
  - thomas
created_date: '2026-04-15 12:15'
updated_date: '2026-04-16 10:42'
labels:
  - parallel
  - concurrency
  - high-priority
dependencies: []
references:
  - 'R3: Parallel Test Safety'
documentation:
  - doc-004
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Modify the database name generation logic in test/test_helper.go to include a unique goroutine identifier alongside the process ID and timestamp. This ensures that parallel test executions don't create databases with duplicate names. The implementation should extract the goroutine ID from the runtime stack trace and append it to the database name prefix. Update the unique name generation function to use this enhanced approach.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 No two parallel tests create databases with the same name
- [ ] #2 Test execution speed is not significantly impacted
- [ ] #3 Database cleanup doesn't interfere with parallel tests
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
