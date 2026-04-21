---
id: RDL-082
title: '[doc-008 Phase 2] Implement DayService with weekly statistics calculations'
status: To Do
assignee:
  - workflow
created_date: '2026-04-21 15:50'
updated_date: '2026-04-21 19:30'
labels:
  - phase-2
  - service
  - calculation
dependencies: []
references:
  - REQ-DASH-005
  - AC-DASH-001
  - Implementation Checklist Phase 2
documentation:
  - doc-008
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement internal/service/dashboard/day_service.go calculating previous_week_pages, last_week_pages, per_pages ratio, mean_day, and spec_mean_day. Use GetToday() for consistent date references and ensure all float values rounded to 3 decimal places.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Previous week page total calculated correctly
- [ ] #2 Last week page total calculated correctly
- [ ] #3 Per pages ratio computed with 3 decimal precision
- [ ] #4 Mean day by weekday calculated accurately
- [ ] #5 Speculative mean derived from actual mean
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
