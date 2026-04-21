---
id: RDL-084
title: '[doc-008 Phase 2] Implement FaultsService with percentage calculation'
status: To Do
assignee: []
created_date: '2026-04-21 15:50'
labels:
  - phase-2
  - service
  - faults
dependencies: []
references:
  - REQ-DASH-007
  - AC-DASH-004
  - 'Decision 4: Fault Calculation Logic'
documentation:
  - doc-008
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement internal/service/dashboard/faults_service.go counting ALL faults (regardless of status) for date range, calculating fault percentage as (faults_last_30_days / max_faults) * 100. Handle zero faults case returning 0% not NaN.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Counts all faults regardless of status (matches Rails)
- [ ] #2 Percentage calculation correct with 2 decimal precision
- [ ] #3 Zero faults returns 0% not NaN/error
- [ ] #4 Max faults from config with default fallback
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
