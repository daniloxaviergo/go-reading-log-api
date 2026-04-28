---
id: RDL-125
title: >-
  [doc-011 Phase 1] Update dashboard_response.go DTOs with ProjectWithLogs and
  StatsData
status: To Do
assignee: []
created_date: '2026-04-28 11:15'
labels:
  - feature
  - backend
  - phase-1
dependencies: []
documentation:
  - doc-011
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update internal/domain/dto/dashboard_response.go to add/verify ProjectWithLogs struct with all required fields (id, name, total_page, page, started_at, progress, status, logs_count, days_unreading, median_day, finished_at, logs) and StatsData struct for progress_geral, total_pages, pages. Implement JSON marshaling for the response structure matching Rails format.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 ProjectWithLogs struct contains all 12 required project fields plus logs array
- [ ] #2 StatsData struct contains progress_geral, total_pages, pages fields
- [ ] #3 JSON marshaling produces flat structure with projects array and stats object at root
- [ ] #4 Float values support 3 decimal precision
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
