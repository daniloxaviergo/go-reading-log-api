---
id: RDL-042
title: >-
  [doc-003 Phase 2] Harmonize JSON response structure to match Rails JSON:API
  format
status: To Do
assignee: []
created_date: '2026-04-12 23:50'
labels:
  - json
  - api
  - structure
  - jsonapi
dependencies: []
references:
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/2'
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/6'
documentation:
  - doc-003
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement FR-003 and FR-004 by updating response serializers to wrap data in a 'data' array with 'type' and 'attributes' keys matching JSON:API specification, and removing nested 'project' objects from log entries to align with Rails API structure.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Implement JSON:API envelope wrapper for project responses
- [ ] #2 Update DTO structs to support type and attributes keys
- [ ] #3 Remove nested project object from log response DTO
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
