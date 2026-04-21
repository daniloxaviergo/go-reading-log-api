---
id: RDL-073
title: '[doc-007 Phase 1] Update GetProjectLogs handler logic for new structure'
status: To Do
assignee:
  - catarina
created_date: '2026-04-21 12:11'
updated_date: '2026-04-21 12:52'
labels:
  - refactoring
  - backend
dependencies: []
references:
  - REQ-02
  - Decision 3
documentation:
  - doc-007
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Modify the internal/api/v1/handlers/logs_handler.go file to update the GetProjectLogs function. Ensure it correctly populates the relationships and included arrays instead of embedding full project objects in each log entry.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Handler returns valid JSON:API structure
- [ ] #2 Relationships populated correctly
- [ ] #3 Included array contains project data
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
- [ ] #11 No breaking changes to route signature
<!-- DOD:END -->
