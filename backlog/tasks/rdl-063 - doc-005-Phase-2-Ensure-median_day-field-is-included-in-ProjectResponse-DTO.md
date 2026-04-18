---
id: RDL-063
title: '[doc-005 Phase 2] Ensure median_day field is included in ProjectResponse DTO'
status: To Do
assignee:
  - workflow
created_date: '2026-04-18 11:47'
updated_date: '2026-04-18 13:06'
labels:
  - phase-2
  - median-day
  - dto
dependencies: []
references:
  - 'PRD Section: Key Requirements REQ-003'
  - internal/domain/dto/project_response.go
documentation:
  - doc-005
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update internal/domain/dto/project_response.go to ensure the median_day field is properly exposed in all project API responses. The field should be a float64 pointer that calculates pages per day rounded to 2 decimal places.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 median_day field present in ProjectResponse struct
- [ ] #2 Field serialized correctly to JSON with proper rounding
- [ ] #3 AC-REQ-003.1 verified: Inspect JSON response structure shows median_day
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
- [ ] #13 Field uses float64 pointer type for optional serialization
- [ ] #14 Rounding to 2 decimal places implemented
<!-- DOD:END -->
