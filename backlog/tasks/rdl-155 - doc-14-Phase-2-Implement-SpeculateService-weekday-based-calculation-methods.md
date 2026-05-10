---
id: RDL-155
title: '[doc-14 Phase 2] Implement SpeculateService weekday-based calculation methods'
status: To Do
assignee: []
created_date: '2026-05-10 10:47'
labels:
  - service
  - calculations
  - phase-2
dependencies: []
documentation:
  - doc-014
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Rewrite internal/service/dashboard/speculate_service.go to implement weekday-based historical mean calculation: CalculateHistoricalMean() with 7-day interval logic, CalculateSpeculativeMean(mean, 0.1) applying 10% prediction buffer, GenerateXAxisLabels() formatting dates as 'DD-MMM (Day)', and GenerateSeriesData() for Pages and Mean arrays.

Service must handle empty databases, partial data (missing days), and zero-fill missing dates in the 15-day range.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 CalculateHistoricalMean implements weekday grouping with 7-day interval division
- [ ] #2 CalculateSpeculativeMean applies 10% buffer (mean * 1.10)
- [ ] #3 GenerateXAxisLabels returns 15 dates in 'DD-MMM (Day)' format
- [ ] #4 GenerateSeriesData creates Pages and Mean arrays with 15 elements each
- [ ] #5 Zero-fill logic for missing days in date range
- [ ] #6 Edge cases handled: nil data, empty database, zero mean values
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
