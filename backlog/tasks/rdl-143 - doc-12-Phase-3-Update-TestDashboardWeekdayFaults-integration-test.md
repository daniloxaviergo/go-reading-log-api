---
id: RDL-143
title: '[doc-12 Phase 3] Update TestDashboardWeekdayFaults integration test'
status: To Do
assignee:
  - catarina
created_date: '2026-05-01 15:08'
updated_date: '2026-05-01 17:24'
labels:
  - bugfix
  - testing
  - phase-3
dependencies: []
documentation:
  - doc-012
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Fix TestDashboardWeekdayFaults integration test by updating test data setup to create reading patterns with specific weekday gaps. Update expected weekday distribution to match Rails logic. Verify radar chart data reflects correct fault counts per weekday.
<!-- SECTION:DESCRIPTION:END -->

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
