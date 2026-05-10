---
id: RDL-158
title: '[doc-14 Phase 3] Add unit tests for dashboard handler SpeculateActual method'
status: To Do
assignee:
  - workflow
created_date: '2026-05-10 10:48'
updated_date: '2026-05-10 14:27'
labels:
  - testing
  - unit-tests
  - handler
  - phase-3
dependencies: []
documentation:
  - doc-014
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create unit tests for dashboard handler SpeculateActual method testing: successful response with valid data, error handling when service returns error, and response format validation. Tests must verify flat JSON structure and proper HTTP status codes.

Tests use mock SpeculateService to avoid database dependencies and follow existing handler test patterns.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 TestSpeculateActual_Success validates 200 OK and flat JSON structure
- [ ] #2 TestSpeculateActual_ServiceError validates 500 status on service error
- [ ] #3 TestSpeculateActual_ResponseFormat verifies echart key exists at root level
- [ ] #4 Tests use mock SpeculateService with configurable return values
- [ ] #5 All handler tests compile and run without errors
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
