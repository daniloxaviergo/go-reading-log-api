---
id: RDL-066
title: '[doc-005 Phase 4] Complete regression test suite for API response alignment'
status: To Do
assignee:
  - workflow
created_date: '2026-04-18 11:47'
updated_date: '2026-04-18 14:37'
labels:
  - phase-4
  - regression-testing
  - comprehensive
dependencies: []
references:
  - 'PRD Section: Acceptance Criteria'
  - test/compare_responses.sh
documentation:
  - doc-005
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create comprehensive regression tests in test/compare_responses.sh and internal/api/v1/handlers/projects_handler_test.go that verify all acceptance criteria are met, including days_unreading tolerance, finished_at calculation, and JSON structure compliance.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Automated comparison tests for days_unreading match Rails within 1 day tolerance
- [ ] #2 finished_at calculation tests cover edge cases
- [ ] #3 JSON:API compliance verified programmatically
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
- [ ] #13 Test coverage >80% for modified code
- [ ] #14 Tests run in CI/CD pipeline
<!-- DOD:END -->
