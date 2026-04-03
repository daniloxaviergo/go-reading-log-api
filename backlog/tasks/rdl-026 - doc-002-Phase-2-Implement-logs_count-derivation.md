---
id: RDL-026
title: '[doc-002 Phase 2] Implement logs_count derivation'
status: To Do
assignee: []
created_date: '2026-04-03 14:03'
labels:
  - phase-2
  - derived-calculation
  - array-count
dependencies: []
references:
  - >-
    PRD Section: Technical Decisions - Decision 1: Derived Calculations
    Implementation
  - 'PRD Section: Validation Rules - logs_count rule'
documentation:
  - doc-002
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement logs_count derivation to match Rails behavior. Count the number of log entries in the logs array (logs_count = logs.size). Ensure this field is always present in the response JSON even when logs array is empty.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 logs_count = len(logs)
- [ ] #2 logs_count included in JSON even if empty array
- [ ] #3 logs_count is an integer type
- [ ] #4 Matches Rails logs.size behavior
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
