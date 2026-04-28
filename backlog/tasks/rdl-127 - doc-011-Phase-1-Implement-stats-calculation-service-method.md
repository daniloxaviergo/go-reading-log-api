---
id: RDL-127
title: '[doc-011 Phase 1] Implement stats calculation service method'
status: To Do
assignee:
  - workflow
created_date: '2026-04-28 11:16'
updated_date: '2026-04-28 12:41'
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
Add CalculateStats() method to projects_service.go that computes stats object: total_pages (sum of all project total_page), pages (sum of all project page), progress_geral (round((pages/total_pages)*100, 3)). Handle edge cases with zero projects and division by zero.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 stats.total_pages equals sum of all project total_page values
- [ ] #2 stats.pages equals sum of all project page values
- [ ] #3 stats.progress_geral calculated as round((pages/total_pages)*100, 3)
- [ ] #4 Zero projects returns stats with all values at 0
- [ ] #5 Division by zero returns 0.0 for progress_geral
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
