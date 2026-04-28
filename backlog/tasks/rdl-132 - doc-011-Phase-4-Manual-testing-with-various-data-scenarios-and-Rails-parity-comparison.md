---
id: RDL-132
title: >-
  [doc-011 Phase 4] Manual testing with various data scenarios and Rails parity
  comparison
status: To Do
assignee:
  - catarina
created_date: '2026-04-28 11:17'
updated_date: '2026-04-28 15:00'
labels:
  - validation
  - backend
  - phase-4
dependencies: []
documentation:
  - doc-011
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Perform manual testing with various data scenarios (empty projects, single project, multiple projects with different statuses). Compare Go endpoint output with Rails endpoint output to verify structural equivalence and calculated field parity. Verify performance targets (< 200ms latency, < 50ms query time).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Test with empty projects returns 200 OK with empty array
- [ ] #2 Test with single running project returns correct data
- [ ] #3 Test with mixed statuses returns only running projects
- [ ] #4 Go response structure matches Rails response structure
- [ ] #5 All calculated fields match Rails values within tolerance
- [ ] #6 Response latency < 200ms at 95th percentile
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
