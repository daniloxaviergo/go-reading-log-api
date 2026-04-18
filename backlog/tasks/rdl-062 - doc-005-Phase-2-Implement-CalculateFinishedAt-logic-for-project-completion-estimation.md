---
id: RDL-062
title: >-
  [doc-005 Phase 2] Implement CalculateFinishedAt logic for project completion
  estimation
status: To Do
assignee:
  - workflow
created_date: '2026-04-18 11:46'
updated_date: '2026-04-18 12:50'
labels:
  - phase-2
  - finished-at
  - median-day
dependencies: []
references:
  - 'PRD Section: Key Requirements REQ-002'
  - internal/domain/models/project.go
documentation:
  - doc-005
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement or verify the CalculateFinishedAt() method in internal/domain/models/project.go that projects completion date based on median_day calculation. The method should return a calculated date when page < total_page and no logs exist, returning null appropriately for edge cases.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 CalculateFinishedAt returns calculated date when page < total_page
- [ ] #2 CalculateFinishedAt returns null when page >= total_page and no logs exist
- [ ] #3 AC-REQ-002.1 and AC-REQ-002.2 acceptance criteria verified
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
- [ ] #13 Formula: finished_at = started_at + (total_page - page) / median_day days
- [ ] #14 Edge cases handled: zero median_day, negative days
<!-- DOD:END -->
