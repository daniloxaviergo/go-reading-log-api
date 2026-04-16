---
id: RDL-050
title: '[doc-004 Phase 1.3] Add database name validation'
status: To Do
assignee: []
created_date: '2026-04-15 12:14'
updated_date: '2026-04-16 00:36'
labels:
  - validation
  - security
  - high-priority
dependencies: []
references:
  - 'R5: Database Name Validation'
documentation:
  - doc-004
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement database name validation in test/test_helper.go to prevent SQL injection and ensure format compliance. The validation should check that names match the pattern reading_log_test(_[a-zA-Z0-9_]+)?$, are limited to 63 characters, and only contain alphanumeric characters, underscores, and hyphens. Create a SafeDropDatabase wrapper function that validates before executing DROP DATABASE, and return clear error messages for invalid names.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Only alphanumeric characters, underscores, and hyphens are allowed
- [ ] #2 Names are limited to 63 characters (PostgreSQL limit)
- [ ] #3 Names must match pattern reading_log_test[_[a-zA-Z0-9_]+]
- [ ] #4 Invalid names are rejected with clear error messages
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
