---
id: RDL-159
title: '[doc-14 Phase 4] Register /v1/dashboard/echart/speculate_actual.json route'
status: To Do
assignee:
  - book
created_date: '2026-05-10 10:48'
updated_date: '2026-05-10 14:46'
labels:
  - routing
  - api
  - phase-4
dependencies: []
documentation:
  - doc-014
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add route registration in internal/api/v1/routes.go: r.HandleFunc("/v1/dashboard/echart/speculate_actual.json", dashboardHandler.SpeculateActual).Methods("GET"). Verify route is properly registered and accessible through the middleware chain.

Route must follow existing routing patterns and be accessible via the configured server port.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Route registered with GET method on /v1/dashboard/echart/speculate_actual.json
- [ ] #2 Route uses dashboardHandler.SpeculateActual handler
- [ ] #3 Route follows existing middleware chain
- [ ] #4 Route compiles without errors
- [ ] #5 Route is accessible via curl test
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
