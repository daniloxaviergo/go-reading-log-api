---
id: RDL-075
title: '[doc-007 Phase 3] Update integration tests for JSON:API schema compliance'
status: To Do
assignee:
  - catarina
created_date: '2026-04-21 12:11'
updated_date: '2026-04-21 13:31'
labels:
  - testing
  - backend
dependencies: []
references:
  - AC-FUNC-01
  - AC-NFUNC-02
documentation:
  - doc-007
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update the test/integration/logs_endpoint_test.go file to validate the new JSON:API response structure, including checks for RFC3339 date format, relationship references, and payload size reduction.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Tests validate JSON:API schema
- [ ] #2 Date format checked for RFC3339
- [ ] #3 Payload size verified
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
- [ ] #11 100% coverage for modified files
<!-- DOD:END -->
