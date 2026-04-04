---
id: RDL-031
title: '[doc-002 Phase 4] Implement page ≤ total_page validation rule'
status: To Do
assignee:
  - catarina
created_date: '2026-04-03 14:04'
updated_date: '2026-04-04 01:29'
labels:
  - phase-4
  - validation-rule
  - business-logic
dependencies: []
references:
  - 'PRD Section: Validation Rules - page ≤ total_page'
  - 'PRD Section: Files to Modify - project_repository.go'
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement validation for project page ≤ total_page constraint. Create validation function in internal/validation/ package and integrate into project creation/update flow. Return appropriate error with error code and message.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Validation function checks page ≤ total_page
- [ ] #2 Error returned when constraint violated
- [ ] #3 Error includes error code and descriptive message
- [ ] #4 Validation logic matches Rails behavior
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
