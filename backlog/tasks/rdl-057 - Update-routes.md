---
id: RDL-057
title: Update routes
status: To Do
assignee:
  - catarina
created_date: '2026-04-16 21:06'
updated_date: '2026-04-16 21:06'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The endpoints should be equals responses
- /v1/projects/{project_id}.json
- /v1/projects/{project_id}/logs.json
- /v1/projects.json

update the prd to correct url, should be `.json` at the end
- /v1/projects/{project_id}.json
- /v1/projects/{project_id}/logs.json
- /v1/projects.json

update for new routes: test/compare_responses.sh
<!-- SECTION:DESCRIPTION:END -->

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
