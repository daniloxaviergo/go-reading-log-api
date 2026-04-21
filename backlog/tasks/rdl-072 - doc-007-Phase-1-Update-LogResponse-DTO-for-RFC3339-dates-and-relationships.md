---
id: RDL-072
title: '[doc-007 Phase 1] Update LogResponse DTO for RFC3339 dates and relationships'
status: To Do
assignee:
  - catarina
created_date: '2026-04-21 12:11'
updated_date: '2026-04-21 12:18'
labels:
  - refactoring
  - backend
dependencies: []
references:
  - REQ-01
  - REQ-02
documentation:
  - doc-007
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update the internal/domain/dto/log_response.go file to change the Data field from string to time.Time and add a Relationships struct. Remove the embedded Project object from the attributes to comply with JSON:API spec.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Data field is time.Time type
- [ ] #2 Relationships struct exists with project data
- [ ] #3 Project field removed from attributes
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
- [ ] #11 go fmt passes
- [ ] #12 go vet passes
<!-- DOD:END -->
