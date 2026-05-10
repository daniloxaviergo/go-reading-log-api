---
id: RDL-157
title: >-
  [doc-14 Phase 3] Rewrite dashboard handler SpeculateActual method with flat
  JSON response
status: To Do
assignee: []
created_date: '2026-05-10 10:48'
labels:
  - handler
  - api
  - phase-3
dependencies: []
documentation:
  - doc-014
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Rewrite internal/api/v1/handlers/dashboard_handler.go SpeculateActual() method to inject SpeculateService, call service methods to generate chart config, and return flat JSON { echart: chartConfig } matching Rails response format. Remove current faults-based implementation.

Handler must handle errors gracefully, return appropriate HTTP status codes (200 OK, 500 Internal Server Error), and follow existing error handling patterns.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 SpeculateActual method injects SpeculateService via dependency injection
- [ ] #2 Handler calls service to generate ECharts configuration
- [ ] #3 Response format is flat JSON { echart: {...} } without JSON:API envelope
- [ ] #4 Error handling returns 500 status with proper error message
- [ ] #5 Handler follows existing middleware and logging patterns
- [ ] #6 Code compiles without errors
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
