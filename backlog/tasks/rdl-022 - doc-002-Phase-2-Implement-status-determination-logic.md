---
id: RDL-022
title: '[doc-002 Phase 2] Implement status determination logic'
status: To Do
assignee:
  - workflow
created_date: '2026-04-03 14:02'
updated_date: '2026-04-03 17:18'
labels:
  - phase-2
  - status-logic
  - business-rules
dependencies: []
references:
  - >-
    PRD Section: Technical Decisions - Decision 1: Derived Calculations
    Implementation
  - 'PRD Section: Validation Rules - status values'
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement status determination logic in Go matching Rails ActiveModelSerializer status method. Status depends on days_unreading ranges (configured values) and logs count: unstarted (no logs started), finished (logs count = total_page), running (days_unreading ≤ em_andamento_range), sleeping (days_unreading ≤ dormindo_range), stopped (all other cases).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Unstarted: No logs or no log with data
- [ ] #2 Finished: Logs count equals total_page
- [ ] #3 Running: days_unreading ≤ em_andamento_range
- [ ] #4 Sleeping: days_unreading ≤ dormindo_range
- [ ] #5 Stopped: All other cases
- [ ] #6 Method implemented in Project model or calculations package
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
