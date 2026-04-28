---
id: RDL-111
title: '[doc-10 Phase 1] Update StatsData DTO with new fields'
status: To Do
assignee: []
created_date: '2026-04-28 00:28'
labels:
  - dto
  - phase-1
  - backend
dependencies: []
documentation:
  - doc-010
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add new fields to StatsData struct in dashboard_response.go: MaxDay, MeanGeral, PerMeanDay, PerSpecMeanDay. Update Validate() method to allow null values for ratio fields.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 StatsData struct includes MaxDay, MeanGeral, PerMeanDay, PerSpecMeanDay fields
- [ ] #2 Validate() method accepts null values for ratio fields
- [ ] #3 DTO compiles without errors
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
