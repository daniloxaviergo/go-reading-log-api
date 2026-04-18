---
id: RDL-061
title: >-
  [doc-005 Phase 1] Implement timezone configuration support for date
  calculations
status: To Do
assignee:
  - workflow
created_date: '2026-04-18 11:46'
updated_date: '2026-04-18 12:26'
labels:
  - phase-1
  - timezone
  - critical
dependencies: []
references:
  - 'PRD Section: Decision 4'
  - internal/config/config.go
documentation:
  - doc-005
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add timezone configuration to internal/config/config.go and update project calculation methods to use configured timezone instead of UTC. The TZLocation variable must be configurable via environment variable with fallback to Brazil timezone, ensuring Date.today behavior matches Rails.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 TZLocation configurable via environment variable with BRT fallback
- [ ] #2 Date calculations use configured timezone, not UTC
- [ ] #3 AC-REQ-006.1 verified: Test with different timezone settings passes
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [ ] #2 All integration tests pass use testing-expert subagent for test execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Database queries optimized with proper indexes
- [ ] #8 Documentation updated in QWEN.md
- [ ] #9 New code paths include error path tests
- [ ] #10 HTTP handlers test both success and error responses
- [ ] #11 Integration tests verify actual database interactions
- [ ] #12 Tests use testing-expert subagent for test execution and verification
- [ ] #13 Configuration loaded at startup with validation
- [ ] #14 Environment variable documented in .env.example
<!-- DOD:END -->
