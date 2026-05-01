---
id: RDL-138
title: '[doc-12 Phase 2] Update mock repository expectations'
status: To Do
assignee: []
created_date: '2026-05-01 15:02'
labels:
  - bugfix
  - testing
  - phase-2
dependencies: []
documentation:
  - doc-012
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update MockDashboardRepository and MockDashboardRepositoryForProjects to return correct fault counts matching the new SQL logic. Ensure mock implementations return days-with-zero-pages counts instead of log-entry counts. Update all test files that use these mocks.
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
