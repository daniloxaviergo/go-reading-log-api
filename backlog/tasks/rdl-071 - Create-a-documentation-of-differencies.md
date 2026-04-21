---
id: RDL-071
title: Create a documentation of differencies
status: To Do
assignee: []
created_date: '2026-04-21 10:35'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create a documentation of differencies of jsons and write in docs/diff_show_logs.md
http://0.0.0.0:3001/v1/projects/450/logs.json -> Rails-Api
http://0.0.0.0:3000/v1/projects/450/logs.json -> Go-Api

Dont change rails-app
The fix should be in golang
<!-- SECTION:DESCRIPTION:END -->

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
