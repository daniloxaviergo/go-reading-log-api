---
id: RDL-125
title: '[doc-011 Phase 3] Create integration tests for dashboard projects endpoint'
status: To Do
assignee: []
created_date: '2026-04-28 11:15'
labels:
  - testing
  - backend
  - phase-3
dependencies: []
documentation:
  - doc-011
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create test/integration/dashboard_projects_test.go with integration tests for GET /v1/dashboard/projects.json endpoint. Test endpoint response structure, Rails parity validation, running status filter, stats calculation, project ordering, and eager-loaded logs. Use TestHelper for database setup/teardown.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Test endpoint returns 200 OK with correct structure
- [ ] #2 Test only running projects included in response
- [ ] #3 Test stats calculation matches expected values
- [ ] #4 Test projects ordered by progress descending
- [ ] #5 Test each project includes first 4 logs ordered by date DESC
- [ ] #6 Test Rails parity validation with identical data
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
