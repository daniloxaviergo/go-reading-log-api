---
id: RDL-119
title: '[doc-10 Phase 4] Update DTO validation for null values'
status: To Do
assignee:
  - book
created_date: '2026-04-28 00:30'
updated_date: '2026-04-28 04:32'
labels:
  - validation
  - phase-4
  - backend
dependencies: []
documentation:
  - doc-010
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update StatsData.Validate() to allow null values for ratio fields. Create tests validating null handling scenarios for per_pages, per_mean_day, per_spec_mean_day.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Validate() accepts null for ratio fields
- [ ] #2 Tests cover all null scenarios
- [ ] #3 No validation errors for valid null values
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
