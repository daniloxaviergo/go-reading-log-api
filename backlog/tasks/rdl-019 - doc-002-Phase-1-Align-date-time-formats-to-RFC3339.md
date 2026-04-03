---
id: RDL-019
title: '[doc-002 Phase 1] Align date time formats to RFC3339'
status: To Do
assignee:
  - catarina
created_date: '2026-04-03 14:02'
updated_date: '2026-04-03 14:15'
labels:
  - phase-1
  - date-format
  - code-quality
dependencies: []
references:
  - 'PRD Section: Technical Decisions - Decision 3: Date/Time Format Alignment'
  - 'PRD Section: Files to Modify - log_response.go'
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update all date/time fields in `log_response.go` to use RFC3339 format for timestamps and ISO date format for started_at field. Ensure NULL date fields serialize to JSON null instead of zero values.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Timestamp fields formatted as RFC3339 (e.g. 2024-01-15T10:30:00Z)
- [ ] #2 Date fields formatted as ISO date (e.g. 2024-01-15)
- [ ] #3 NULL database values serialize to JSON null
- [ ] #4 Format matches Rails API output exactly
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
