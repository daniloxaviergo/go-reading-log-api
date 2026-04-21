---
id: RDL-070
title: Fix days_unreading and median_day math
status: To Do
assignee:
  - catarina
created_date: '2026-04-21 10:15'
updated_date: '2026-04-21 10:17'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
In route the fields days_unreading and median_day is diferent
http://0.0.0.0:3001/v1/projects/450.json -> Rails-Api
http://0.0.0.0:3000/v1/projects/450.json -> Go-Api

| Field | Go-Api | Rails-Api |
| days-unreading | 19 | 15 |
| median-day | 11.33 | 12.12 |

Dont change the rails-app
Look the code rais-app to check the math and fix
<!-- SECTION:DESCRIPTION:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 go fmt and go vet pass with no errors
- [ ] #2 Clean Architecture layers properly followed
- [ ] #3 Error responses consistent with existing patterns
- [ ] #4 HTTP status codes correct for response type
- [ ] #5 Documentation updated in QWEN.md
- [ ] #6 New code paths include error path tests
- [ ] #7 HTTP handlers test both success and error responses
- [ ] #8 Integration tests verify actual database interactions
- [ ] #9 Tests use testing-expert subagent for test execution and verification
- [ ] #10 All unit tests pass
- [ ] #11 All integration tests pass
<!-- DOD:END -->
