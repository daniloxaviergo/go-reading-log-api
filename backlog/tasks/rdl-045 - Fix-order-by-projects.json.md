---
id: RDL-045
title: Fix order by projects.json
status: To Do
assignee:
  - catarina
created_date: '2026-04-14 09:53'
updated_date: '2026-04-14 09:53'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add correct order query database, look rails-app
execute test/compare_responses.sh to endpoint v1/projects.json

The Go API and Rails API return completely different projects from the same database query. The Go API returns project ID 1 ("Filocalia"), while the Rails API returns project ID 450 ("História da Igreja VIII.1").

**Go Response (First Project):**
```json
{
  "days_unreading": 3354,
  "finished_at": null,
  "id": 1,
  "logs_count": 50,
  "name": "Filocalia",
  "page": 655,
  "progress": null,
  "started_at": "2017-02-04T00:00:00Z",
  "status": "stopped",
  "total_page": 1267
}
```

**Rails Response (First Project):**
```json
{
  "days_unreading": 10,
  "finished_at": null,
  "id": 450,
  "logs_count": 38,
  "name": "História da Igreja VIII.1",
  "page": 691,
  "progress": 100.0,
  "started_at": "2026-02-19",
  "status": "finished",
  "total_page": 691
}
```
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
