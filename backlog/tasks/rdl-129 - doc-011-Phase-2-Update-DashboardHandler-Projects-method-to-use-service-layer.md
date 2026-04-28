---
id: RDL-129
title: '[doc-011 Phase 2] Update DashboardHandler Projects method to use service layer'
status: To Do
assignee: []
created_date: '2026-04-28 11:16'
labels:
  - feature
  - backend
  - phase-2
dependencies: []
documentation:
  - doc-011
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Modify internal/api/v1/handlers/dashboard_handler.go Projects() method to use new projects service instead of direct repository calls. Implement response formatting matching Rails structure with projects array and stats object at root level. Add error handling and structured logging.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Handler calls service GetRunningProjectsWithLogs method
- [ ] #2 Handler calls service CalculateStats method
- [ ] #3 Response structure matches Rails (projects array + stats object)
- [ ] #4 Database errors return 500 Internal Server Error with logging
- [ ] #5 Empty data returns 200 OK with empty arrays
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
