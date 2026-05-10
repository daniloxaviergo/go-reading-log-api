---
id: RDL-152
title: >-
  [doc-14 Phase 1] Add repository interface methods for weekday-based mean
  calculations
status: To Do
assignee:
  - book
created_date: '2026-05-10 10:46'
updated_date: '2026-05-10 11:18'
labels:
  - infrastructure
  - repository
  - phase-1
dependencies: []
documentation:
  - doc-014
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Define new methods in DashboardRepository interface: GetWeekdayPagesGrouped(ctx, startDate, endDate) for fetching pages grouped by weekday, GetFirstLogDate(ctx) for getting the earliest log date, and GetWeekdayMeanWithIntervals(ctx, weekday, currentDate) for calculating mean with 7-day interval logic.

These methods support the weekday-based historical mean calculation algorithm required by the speculate_actual endpoint.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 GetWeekdayPagesGrouped method defined in interface with correct signature
- [ ] #2 GetFirstLogDate method defined in interface
- [ ] #3 GetWeekdayMeanWithIntervals method defined in interface
- [ ] #4 All interface methods documented with comments
- [ ] #5 Interface compiles without errors
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass
- [ ] #2 All integration tests pass execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md and AGENTS.md
- [ ] #8 New code paths include error path tests
- [ ] #9 HTTP handlers test both success and error responses
- [ ] #10 Integration tests verify actual database interactions
<!-- DOD:END -->
