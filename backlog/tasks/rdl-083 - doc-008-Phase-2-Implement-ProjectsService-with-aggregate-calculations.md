---
id: RDL-083
title: '[doc-008 Phase 2] Implement ProjectsService with aggregate calculations'
status: To Do
assignee:
  - book
created_date: '2026-04-21 15:50'
updated_date: '2026-04-21 20:24'
labels:
  - phase-2
  - service
  - aggregate
dependencies: []
references:
  - REQ-DASH-006
  - AC-DASH-002
  - Implementation Checklist Phase 2
documentation:
  - doc-008
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement internal/service/dashboard/projects_service.go querying all projects with eager-loaded logs (first 4, ordered by date DESC), calculating progress_geral, total_pages, and pages aggregates. Order results by progress descending.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 All projects retrieved with eager-loaded logs
- [ ] #2 Log ordering correct (first 4, date DESC)
- [ ] #3 Progress aggregate calculated correctly
- [ ] #4 Total pages and pages aggregates accurate
- [ ] #5 Results ordered by progress descending
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
