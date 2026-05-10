---
id: RDL-160
title: '[doc-14 Phase 5] Create integration tests for speculate_actual endpoint'
status: To Do
assignee: []
created_date: '2026-05-10 10:48'
labels:
  - testing
  - integration-tests
  - phase-5
dependencies: []
documentation:
  - doc-014
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create test/integration/api/v1/dashboard/echart_speculate_actual_test.go with integration tests covering: empty database scenario (all zeros), partial data (some days missing), complete data (all 15 days populated), response format matching Rails exactly, and markPoint/markLine configuration validation.

Tests use TestHelper for database setup/teardown and verify actual HTTP endpoint behavior.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 TestEmptyDatabase validates all series are zero-filled
- [ ] #2 TestPartialData validates zero-fill for missing days
- [ ] #3 TestCompleteData validates 15 data points in all series
- [ ] #4 TestResponseFormat validates flat JSON with echart key
- [ ] #5 TestSeriesNames validates 'Pages' and 'Mean' names
- [ ] #6 TestMarkElements validates markPoint max/min and markLine yAxis: 40
- [ ] #7 TestDateRange validates 15 dates in xAxis
- [ ] #8 All tests use real database via TestHelper
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
