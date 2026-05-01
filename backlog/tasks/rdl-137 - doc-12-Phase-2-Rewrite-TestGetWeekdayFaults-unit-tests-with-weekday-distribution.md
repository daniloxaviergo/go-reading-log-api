---
id: RDL-137
title: >-
  [doc-12 Phase 2] Rewrite TestGetWeekdayFaults unit tests with weekday
  distribution
status: To Do
assignee: []
created_date: '2026-05-01 14:59'
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
Rewrite all TestGetWeekdayFaults test cases to validate correct weekday grouping. Create test cases for: all weekdays represented, some weekdays missing, zero faults for all days, and all days are faults. Update expected weekday distribution maps to match Rails grouping logic.
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
