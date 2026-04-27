---
id: RDL-106
title: Fix all tests
status: To Do
assignee:
  - Thomas
created_date: '2026-04-27 18:44'
updated_date: '2026-04-27 18:46'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Fix all tests
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Fix all tests
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
# Implementation Plan: Fix All Tests

## Test Failures Identified

### 1. TestDashboardWeekdayFaults_Integration (test/dashboard_integration_test.go:411)
- **Issue**: Expected weekday fault data [0,1,2,3,4,5,6] but got [0,0,0,0,0,0,0]
- **Root Cause**: Dashboard service not returning correct fault data for weekdays

### 2. TestDashboardMeanProgress_Integration (test/dashboard_integration_test.go:487)
- **Issue**: Expected 30 items in series but got 1
- **Root Cause**: Mean progress endpoint returning single value instead of 30-day series

### 3. TestDashboardYearlyTotal_Integration (test/dashboard_integration_test.go:584)
- **Issue**: Expected "line" chart type but got "bar", expected 52 items but got 1
- **Root Cause**: Yearly total endpoint returning wrong chart type and data points

### 4. TestDashboardEndpoints_ErrorHandling/Last_Days_Invalid_Type (test/dashboard_integration_test.go:655)
- **Issue**: Received unexpected end of JSON error
- **Root Cause**: Error handling not returning proper JSON error response

### 5. TestErrorScenarios/Mean_Progress_-_Empty_Database (test/integration/error_scenarios_test.go:176)
- **Issue**: Expected value not nil assertion failed
- **Root Cause**: Mean progress endpoint not handling empty database scenario correctly

## Execution Steps

1. **Examine failing test code** to understand expected behavior
2. **Review dashboard service implementation** to identify bugs
3. **Fix weekday faults calculation** - ensure correct data is returned
4. **Fix mean progress series** - return 30-day data points
5. **Fix yearly total chart type** - return line chart with 52 weeks
6. **Fix error handling** - return proper JSON error responses
7. **Run tests** to verify fixes
8. **Run go fmt and go vet** to ensure code quality
9. **Check all acceptance criteria**
<!-- SECTION:NOTES:END -->

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
