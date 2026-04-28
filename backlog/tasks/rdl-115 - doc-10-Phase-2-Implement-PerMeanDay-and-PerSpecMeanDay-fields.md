---
id: RDL-115
title: '[doc-10 Phase 2] Implement PerMeanDay and PerSpecMeanDay fields'
status: To Do
assignee: []
created_date: '2026-04-28 00:29'
labels:
  - repository
  - phase-2
  - backend
dependencies: []
documentation:
  - doc-010
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add repository methods to fetch previous period mean and speculated mean. Implement ratio calculations: per_mean_day = current_mean / previous_mean, per_spec_mean_day = current_mean / speculated_mean.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 GetPreviousPeriodMean() method implemented
- [ ] #2 Speculated mean calculation logic added
- [ ] #3 Ratio fields computed correctly
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
