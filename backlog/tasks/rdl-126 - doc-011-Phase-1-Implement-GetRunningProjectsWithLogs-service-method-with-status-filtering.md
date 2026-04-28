---
id: RDL-126
title: >-
  [doc-011 Phase 1] Implement GetRunningProjectsWithLogs service method with
  status filtering
status: To Do
assignee:
  - workflow
created_date: '2026-04-28 11:16'
updated_date: '2026-04-28 12:01'
labels:
  - feature
  - backend
  - phase-1
dependencies: []
documentation:
  - doc-011
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create internal/service/dashboard/projects_service.go with GetRunningProjectsWithLogs() method that filters projects by calculated 'running' status using 7-day threshold. Implement progress calculation (page/total_page*100) and ordering logic (progress DESC, id ASC).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 GetRunningProjectsWithLogs returns only projects with status='running'
- [ ] #2 Status calculation uses 7-day threshold for running status
- [ ] #3 Progress calculated as (page/total_page)*100
- [ ] #4 Projects ordered by progress DESC, then id ASC
- [ ] #5 Division by zero handled returning 0.0
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
