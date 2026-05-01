---
id: RDL-140
title: '[doc-12 Phase 1] Update GetWeekdayFaults SQL query'
status: To Do
assignee: []
created_date: '2026-05-01 15:07'
labels:
  - bugfix
  - repository
  - phase-1
dependencies: []
documentation:
  - doc-012
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement SQL query in GetWeekdayFaults to count zero-page days grouped by weekday (0-6). Use same daily aggregation logic as GetFaultsByDateRange, add EXTRACT(DOW) for weekday grouping, and ensure all 7 weekdays appear in result with default 0 values.
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
