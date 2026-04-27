---
id: RDL-106
title: Fix all tests
status: To Do
assignee:
  - Thomas
created_date: '2026-04-27 18:44'
updated_date: '2026-04-27 18:48'
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
# Implementation Progress

## Issues Fixed

### 1. WeekdayFaults - Fixed
- Service correctly returns all 7 weekdays (0-6)
- Issue: Repository not returning data from test database

### 2. MeanProgress - In Progress
- Returns empty slice when no logs found
- Test expects 30 data points with proper chart config
- Fix: Generate 30-day data points even with empty data

### 3. YearlyTotal - Needs Fix
- Currently returns bar chart with 1 data point per year
- Test expects line chart with 52 weekly data points
- Fix: Implement weekly aggregation over 52 weeks

### 4. LastDays Error Handling - Needs Fix
- Returns empty JSON instead of proper error response
- Fix: Return proper JSON error response

### 5. MeanProgress Empty Database - Needs Fix
- Returns nil Echart when no data
- Test expects valid chart config even with empty data
- Fix: Return empty but valid chart config

## Next Steps
1. Fix MeanProgress to return 30 data points
2. Fix YearlyTotal to return 52-week line chart
3. Fix error handling in LastDays
4. Fix MeanProgress empty database scenario
5. Run tests to verify fixes
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
