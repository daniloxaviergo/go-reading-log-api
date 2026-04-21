---
id: RDL-086
title: '[doc-008 Phase 3] Implement weekday fault counting with radar chart data'
status: To Do
assignee: []
created_date: '2026-04-21 15:51'
labels:
  - phase-3
  - service
  - weekdays
dependencies: []
references:
  - REQ-DASH-006
  - AC-DASH-006
  - Implementation Checklist Phase 3
documentation:
  - doc-008
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement weekday fault grouping logic generating radar chart data. Group faults by weekday (Sunday=0 through Saturday=6), cover last 6 months, ensure all 7 weekdays present in output with integer counts >= 0.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Faults grouped by weekday (0-6)
- [ ] #2 6-month date range covered
- [ ] #3 All 7 weekdays present in output
- [ ] #4 Integer counts non-negative
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
