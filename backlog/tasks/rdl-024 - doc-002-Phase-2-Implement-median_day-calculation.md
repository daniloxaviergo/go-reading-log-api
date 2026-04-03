---
id: RDL-024
title: '[doc-002 Phase 2] Implement median_day calculation'
status: To Do
assignee:
  - catarina
created_date: '2026-04-03 14:03'
updated_date: '2026-04-03 21:02'
labels:
  - phase-2
  - derived-calculation
  - arithmetic-calculation
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
Implement median_day calculation in Go matching Rails: median_day = page / days_reading.round(2). Round days_reading to 2 decimal places before division. Handle edge cases: zero days_reading returns 0.00.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 median_day = page / days_reading.round(2)
- [ ] #2 days_reading rounded to 2 decimal places before division
- [ ] #3 Zero days_reading edge case returns 0.00
- [ ] #4 Result is a float64 value
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
