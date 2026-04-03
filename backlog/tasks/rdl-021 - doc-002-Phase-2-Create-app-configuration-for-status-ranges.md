---
id: RDL-021
title: '[doc-002 Phase 2] Create app configuration for status ranges'
status: To Do
assignee: []
created_date: '2026-04-03 14:02'
labels:
  - phase-2
  - configuration
  - setup
dependencies: []
references:
  - 'PRD Section: Technical Decisions - Decision 2: Configuration Values'
  - 'PRD Section: Files to Modify - config.go'
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create a Go configuration structure in `internal/config/config.go` with `em_andamento_range` (7 days default) and `dormindo_range` (14 days default) values matching Rails configuration. Add methods to access these values from status calculation logic.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 em_andamento_range = 7 days
- [ ] #2 dormindo_range = 14 days
- [ ] #3 Access methods for configuration values
- [ ] #4 Configuration loads from environment variables or defaults
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
<!-- DOD:END -->
