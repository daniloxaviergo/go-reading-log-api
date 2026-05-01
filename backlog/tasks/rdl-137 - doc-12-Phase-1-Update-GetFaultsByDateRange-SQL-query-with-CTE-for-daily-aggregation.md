---
id: RDL-137
title: >-
  [doc-12 Phase 1] Update GetFaultsByDateRange SQL query with CTE for daily
  aggregation
status: To Do
assignee: []
created_date: '2026-05-01 14:59'
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
Implement CTE-based SQL query in GetFaultsByDateRange to count days with zero pages read instead of counting log entries. Use generate_series to create all dates in range, LEFT JOIN with daily_read aggregation, and CASE statement to handle NULL values. Include comprehensive inline comments explaining the fault calculation logic.
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
