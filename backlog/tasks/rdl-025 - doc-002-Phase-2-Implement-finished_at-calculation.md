---
id: RDL-025
title: '[doc-002 Phase 2] Implement finished_at calculation'
status: To Do
assignee: []
created_date: '2026-04-03 14:03'
labels:
  - phase-2
  - derived-calculation
  - date-calculation
dependencies: []
references:
  - >-
    PRD Section: Technical Decisions - Decision 1: Derived Calculations
    Implementation
documentation:
  - doc-002
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement finished_at calculation in Go matching Rails: calculate future date when book will be finished based on reading rate (median_day). If progress is 100% or pages remaining is 0, return nil/null. Otherwise, calculate: days_to_finish = (total_page - page) / median_day, then finished_at = today + days_to_finish days.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 finished_at = today + (total_page - page) / median_day days
- [ ] #2 100% progress edge case returns null
- [ ] #3 Pages remaining = 0 edge case returns null
- [ ] #4 Date calculated as future date in days
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
