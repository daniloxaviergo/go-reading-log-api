---
id: RDL-163
title: >-
  [doc-14 Phase 6] Update API documentation with speculate_actual endpoint
  details
status: To Do
assignee:
  - workflow
created_date: '2026-05-10 10:49'
updated_date: '2026-05-10 16:55'
labels:
  - documentation
  - api-docs
  - phase-6
dependencies: []
documentation:
  - doc-014
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update API documentation (README.md or docs/README.md) with complete endpoint details for GET /v1/dashboard/echart/speculate_actual.json including: request format, response schema, example JSON response, and all acceptance criteria (AC-001 to AC-015) documented.

Documentation must match the format of existing endpoint documentation and include curl examples.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Endpoint documented in API reference section
- [ ] #2 Request format shows GET method with no parameters
- [ ] #3 Response schema includes echart object structure
- [ ] #4 Example JSON response with all fields documented
- [ ] #5 Curl example included for testing
- [ ] #6 Acceptance criteria AC-001 to AC-015 listed
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
