---
id: RDL-112
title: '[doc-10 Phase 1] Modify Day handler to return flat JSON'
status: To Do
assignee:
  - catarina
created_date: '2026-04-28 00:28'
updated_date: '2026-04-28 01:32'
labels:
  - handler
  - phase-1
  - backend
dependencies: []
documentation:
  - doc-010
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update Day() method in dashboard_handler.go to return flat JSON structure with stats object at root level instead of JSON:API envelope. Remove data, type, id, attributes wrapper.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Response format is flat JSON with stats key
- [ ] #2 No JSON:API envelope present in response
- [ ] #3 Content-type remains application/json
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
<!-- DOD:END -->
