---
id: RDL-080
title: '[doc-008 Phase 1] Create DashboardResponse DTOs with JSON marshaling'
status: To Do
assignee:
  - workflow
created_date: '2026-04-21 15:49'
updated_date: '2026-04-21 16:56'
labels:
  - phase-1
  - dto
  - api
dependencies: []
references:
  - REQ-DASH-003
  - AC-DASH-003
  - 'Decision 5: Response Format - Chart Configurations'
documentation:
  - doc-008
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement internal/domain/dto/dashboard_response.go defining all response structures including DashboardResponse, EchartConfig, StatsData, and LogEntry. Ensure proper JSON field tags and implement validation methods for each DTO.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 All response DTOs defined with correct JSON tags
- [ ] #2 EchartConfig supports ECharts-style configurations
- [ ] #3 StatsData includes all required aggregate fields
- [ ] #4 Validation methods implemented for each DTO
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
