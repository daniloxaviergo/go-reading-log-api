---
id: RDL-085
title: >-
  [doc-008 Phase 3] Implement SpeculateService with predicted vs actual
  comparison
status: To Do
assignee: []
created_date: '2026-04-21 15:50'
labels:
  - phase-3
  - service
  - prediction
dependencies: []
references:
  - REQ-DASH-008
  - AC-DASH-005
  - Implementation Checklist Phase 3
documentation:
  - doc-008
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement internal/service/dashboard/speculate_service.go comparing actual vs predicted reading. Generate chart data for last 15 days including today, calculate speculative mean as actual_mean * (1 + prediction_pct), and zero-fill missing days.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Actual vs predicted comparison implemented
- [ ] #2 Speculative mean formula correct (actual * (1 + pct))
- [ ] #3 Last 15 days data coverage including today
- [ ] #4 Missing days zero-filled not omitted
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
