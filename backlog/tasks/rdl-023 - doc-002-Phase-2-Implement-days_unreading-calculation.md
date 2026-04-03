---
id: RDL-023
title: '[doc-002 Phase 2] Implement days_unreading calculation'
status: To Do
assignee:
  - next-task
created_date: '2026-04-03 14:03'
updated_date: '2026-04-03 20:15'
labels:
  - phase-2
  - derived-calculation
  - date-calculation
dependencies: []
references:
  - >-
    PRD Section: Technical Decisions - Decision 1: Derived Calculations
    Implementation
  - 'PRD Section: Validation Rules - days_unreading rule'
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement days_unreading calculation in Go matchingRails: days_unreading = (Date.today - last_log_or_started_at).days. Handle edge cases: if no logs, use started_at; if no logs and no started_at, return 0. Method should return non-negative integer.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 days_unreading = today minus last log date or started_at
- [ ] #2 If no logs, use started_at date
- [ ] #3 If neither logs nor started_at exist, return 0
- [ ] #4 Result is non-negative integer
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
