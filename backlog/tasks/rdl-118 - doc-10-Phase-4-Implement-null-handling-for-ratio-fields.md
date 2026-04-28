---
id: RDL-118
title: '[doc-10 Phase 4] Implement null handling for ratio fields'
status: To Do
assignee:
  - workflow
created_date: '2026-04-28 00:30'
updated_date: '2026-04-28 04:09'
labels:
  - null-handling
  - phase-4
  - backend
dependencies: []
documentation:
  - doc-010
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update per_pages logic to return null when previous_week_pages = 0. Apply same null handling to per_mean_day and per_spec_mean_day when denominator is 0 or nil.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 per_pages returns null when previous_week_pages = 0
- [ ] #2 Ratio fields return null when denominator = 0
- [ ] #3 JSON serialization handles null correctly
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
