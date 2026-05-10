---
id: RDL-151
title: >-
  [doc-14 Phase 1] Extend ECharts DTOs with MarkPoint, MarkLine, and BoundaryGap
  fields
status: To Do
assignee: []
created_date: '2026-05-10 10:46'
updated_date: '2026-05-10 10:52'
labels:
  - infrastructure
  - dto
  - phase-1
dependencies: []
documentation:
  - doc-014
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add MarkPoint and MarkLine structs to internal/domain/dto/echart_config.go to support Rails-matching chart configuration. Add BoundaryGap field to Axis struct. Update NewEchartConfig() function to support mark elements configuration.

This enables the speculate_actual endpoint to return complete ECharts configuration with max/min markers and pages_per_day reference line.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 MarkPoint struct with data array containing type and name fields
- [ ] #2 MarkLine struct with data array containing name and yAxis fields
- [ ] #3 BoundaryGap field added to Axis struct as []bool
- [ ] #4 NewEchartConfig() updated to accept markPoint and markLine parameters
- [ ] #5 All DTOs compile without errors
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
