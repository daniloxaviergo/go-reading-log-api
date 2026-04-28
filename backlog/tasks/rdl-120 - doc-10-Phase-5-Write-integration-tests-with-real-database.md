---
id: RDL-120
title: '[doc-10 Phase 5] Write integration tests with real database'
status: To Do
assignee:
  - book
created_date: '2026-04-28 00:30'
updated_date: '2026-04-28 04:48'
labels:
  - integration-testing
  - phase-5
  - backend
dependencies: []
documentation:
  - doc-010
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create integration tests in dashboard_integration_test.go using TestHelper for database setup. Test all new fields with real PostgreSQL queries and verify calculation accuracy.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Integration tests pass with real database
- [ ] #2 All new fields tested with fixtures
- [ ] #3 Test coverage >= 80% for new code
<!-- AC:END -->

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
