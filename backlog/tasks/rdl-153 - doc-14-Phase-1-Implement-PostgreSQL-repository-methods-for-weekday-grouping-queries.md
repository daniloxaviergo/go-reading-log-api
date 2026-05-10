---
id: RDL-153
title: >-
  [doc-14 Phase 1] Implement PostgreSQL repository methods for weekday grouping
  queries
status: To Do
assignee:
  - workflow
created_date: '2026-05-10 10:47'
updated_date: '2026-05-10 11:43'
labels:
  - infrastructure
  - repository
  - postgres
  - phase-1
dependencies: []
documentation:
  - doc-014
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement the three new repository methods in internal/adapter/postgres/dashboard_repository.go: GetWeekdayPagesGrouped with SQL query for weekday-based page aggregation, GetFirstLogDate to retrieve earliest log timestamp, and GetWeekdayMeanWithIntervals calculating mean using 7-day interval logic.

Each method must use 15-second context timeout and handle NULL/empty data gracefully.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 GetWeekdayPagesGrouped implemented with EXTRACT(DOW FROM data) weekday grouping
- [ ] #2 GetFirstLogDate returns *time.Time (nil when no logs exist)
- [ ] #3 GetWeekdayMeanWithIntervals calculates 7-day intervals correctly
- [ ] #4 All methods use 15-second context timeout
- [ ] #5 NULL values handled gracefully in all queries
- [ ] #6 Code compiles without errors
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
