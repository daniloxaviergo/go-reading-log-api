---
id: RDL-060
title: >-
  [doc-005 Phase 1] Update date parsing to support multiple formats and timezone
  configuration
status: To Do
assignee:
  - wokflow
created_date: '2026-04-18 11:46'
updated_date: '2026-04-18 11:56'
labels:
  - phase-1
  - date-calculation
  - critical
dependencies: []
references:
  - 'https://github.com/go-reading-log-api-next/internal/domain/models/project.go'
  - 'PRD Section: Decision 1'
documentation:
  - doc-005
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement multi-format date parsing in internal/domain/models/project.go to fix the 42-day discrepancy between Go and Rails API. The parseLogDate function must support YYYY-MM-DD, RFC3339, and standard datetime formats with timezone-aware comparison matching Rails' Date.today behavior.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 parseLogDate supports at least 3 date formats (YYYY-MM-DD, RFC3339, standard datetime)
- [ ] #2 CalculateDaysUnreading uses timezone-aware comparison matching Rails
- [ ] #3 Unit tests validate edge cases with different date formats
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
- [ ] #13 Code follows Go formatting standards
- [ ] #14 All new functions have unit tests with >80% coverage
<!-- DOD:END -->
