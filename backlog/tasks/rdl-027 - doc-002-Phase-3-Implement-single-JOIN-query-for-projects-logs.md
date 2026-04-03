---
id: RDL-027
title: '[doc-002 Phase 3] Implement single JOIN query for projects + logs'
status: To Do
assignee: []
created_date: '2026-04-03 14:03'
labels:
  - phase-3
  - query-optimization
  - database
dependencies: []
references:
  - 'PRD Section: Technical Decisions - Decision 4: Database Query Optimization'
  - 'PRD Section: Files to Modify - project_repository.go'
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Replace current N+1 queries with a single LEFT OUTER JOIN query in project_repository.go matching Rails eager loading. Query: SELECT p.id, p.name, p.total_page, p.started_at, p.page, p.reinicia, l.id as log_id, l.data, l.start_page, l.end_page, l.note FROM projects p LEFT JOIN logs l ON p.id = l.project_id ORDER BY p.id, l.data DESC
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Single JOIN query replaces N+1 pattern
- [ ] #2 Ordering matches Rails (projects.id, logs.data DESC)
- [ ] #3 LEFT OUTER JOIN used to include projects without logs
- [ ] #4 Query executes in expected time
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
