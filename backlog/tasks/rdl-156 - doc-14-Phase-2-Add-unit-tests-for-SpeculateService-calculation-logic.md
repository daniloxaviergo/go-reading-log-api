---
id: RDL-156
title: '[doc-14 Phase 2] Add unit tests for SpeculateService calculation logic'
status: To Do
assignee:
  - catarina
created_date: '2026-05-10 10:47'
updated_date: '2026-05-10 13:20'
labels:
  - testing
  - unit-tests
  - service
  - phase-2
dependencies: []
documentation:
  - doc-014
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create test/unit/service/dashboard/speculate_service_test.go with comprehensive unit tests for all SpeculateService methods. Test scenarios include: CalculateHistoricalMean with various data distributions, CalculateSpeculativeMean edge cases (zero mean, nil data), GenerateXAxisLabels format validation, and zero-fill logic for missing days.

Tests must use mock repository and cover all edge cases defined in acceptance criteria.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 TestCalculateHistoricalMean covers normal data, empty data, and single-read scenarios
- [ ] #2 TestCalculateSpeculativeMean validates 10% buffer and zero-mean edge case
- [ ] #3 TestGenerateXAxisLabels verifies 'DD-MMM (Day)' format for all 15 dates
- [ ] #4 TestZeroFillLogic validates missing days are filled with zero values
- [ ] #5 TestCalculateHistoricalMean_WeekdayGrouping validates weekday-specific calculations
- [ ] #6 All tests achieve >80% code coverage
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
