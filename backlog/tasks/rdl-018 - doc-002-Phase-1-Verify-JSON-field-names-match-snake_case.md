---
id: RDL-018
title: '[doc-002 Phase 1] Verify JSON field names match (snake_case)'
status: To Do
assignee:
  - workflow
created_date: '2026-04-03 14:02'
updated_date: '2026-04-03 14:08'
labels:
  - phase-1
  - field-alignment
  - code-quality
dependencies: []
references:
  - 'PRD Section: Technical Decisions - Decision 3: Date/Time Format Alignment'
  - 'PRD Section: Files to Modify - project_response.go'
  - log_response.go
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Review Go DTO response structures in `internal/domain/dto/` package and confirm all field names use snake_case matching Rails API JSON output. Update struct tags if needed to ensure JSON keys match exactly.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 All field names in project_response.go match Rails API snake_case format
- [ ] #2 All field names in log_response.go match Rails API snake_case format
- [ ] #3 Null handling verified for optional date fields (started_at, finished_at)
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
