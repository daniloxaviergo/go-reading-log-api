---
id: RDL-161
title: '[doc-14 Phase 5] Add unit tests for repository weekday grouping methods'
status: To Do
assignee:
  - catarina
created_date: '2026-05-10 10:49'
updated_date: '2026-05-10 15:50'
labels:
  - testing
  - unit-tests
  - repository
  - phase-5
dependencies: []
documentation:
  - doc-014
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create test/unit/repository/dashboard_weekday_test.go with unit tests for new repository methods: GetWeekdayPagesGrouped, GetFirstLogDate, and GetWeekdayMeanWithIntervals. Tests must validate SQL query behavior, NULL handling, and edge cases (empty results, single row).

Tests use mock database or test database with TestHelper for verification of actual SQL execution.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 TestGetWeekdayPagesGrouped validates weekday aggregation results
- [ ] #2 TestGetFirstLogDate validates nil return on empty table
- [ ] #3 TestGetWeekdayMeanWithIntervals validates 7-day interval calculation
- [ ] #4 TestEmptyResults validates empty slice returns
- [ ] #5 TestSingleRow validates single log entry handling
- [ ] #6 All tests compile and execute without errors
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
