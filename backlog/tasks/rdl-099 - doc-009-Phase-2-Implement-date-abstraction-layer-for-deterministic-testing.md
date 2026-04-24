---
id: RDL-099
title: '[doc-009 Phase 2] Implement date abstraction layer for deterministic testing'
status: To Do
assignee:
  - book
created_date: '2026-04-24 13:41'
updated_date: '2026-04-24 14:37'
labels:
  - feature
  - test-fix
  - p2-high
dependencies: []
references:
  - REQ-03
  - Decision 2
documentation:
  - doc-009
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create date abstraction layer in internal/domain/dto/dashboard.go with GetTodayFunc variable allowing test-specific date injection. Update all SpeculateService unit tests to use the abstracted date function and fix index assertions to ensure deterministic test results regardless of run date.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Date-dependent tests produce consistent results across different days
- [ ] #2 SpeculateService tests use abstracted date function
- [ ] #3 All 9 SpeculateService unit tests pass deterministically
- [ ] #4 Test execution time remains under 30 seconds
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
