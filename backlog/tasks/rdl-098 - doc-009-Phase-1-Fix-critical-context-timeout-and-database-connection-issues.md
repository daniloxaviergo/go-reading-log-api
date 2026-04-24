---
id: RDL-098
title: '[doc-009 Phase 1] Fix critical context timeout and database connection issues'
status: To Do
assignee:
  - catarina
created_date: '2026-04-24 13:41'
updated_date: '2026-04-24 14:15'
labels:
  - bug
  - test-fix
  - p1-critical
dependencies: []
references:
  - REQ-01
  - REQ-02
  - REQ-05
documentation:
  - doc-009
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Fix GetTestContext() and GetTestContextWithTimeout() functions to return cancel functions instead of discarding them, preventing resource leaks. Add database availability checks with timeout to integration tests to prevent hangs during test execution.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 All 9 SpeculateService unit tests pass without panics
- [ ] #2 Context timeout tests complete within 5 seconds
- [ ] #3 Integration tests have proper database availability checks
- [ ] #4 No resource leaks detected in test execution
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
