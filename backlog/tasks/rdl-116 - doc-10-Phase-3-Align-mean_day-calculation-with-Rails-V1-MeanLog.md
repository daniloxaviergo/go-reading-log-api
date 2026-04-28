---
id: RDL-116
title: '[doc-10 Phase 3] Align mean_day calculation with Rails V1::MeanLog'
status: To Do
assignee:
  - catarina
created_date: '2026-04-28 00:29'
updated_date: '2026-04-28 03:09'
labels:
  - calculation
  - phase-3
  - backend
dependencies: []
documentation:
  - doc-010
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Study Rails V1::MeanLog implementation and replicate exact algorithm in Go. Formula: total_pages / count_reads where count_reads = number of 7-day intervals since begin_data. Round to 3 decimals.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 mean_day calculation matches Rails output exactly
- [ ] #2 Algorithm uses 7-day intervals from begin_data
- [ ] #3 Values rounded to 3 decimal places
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
