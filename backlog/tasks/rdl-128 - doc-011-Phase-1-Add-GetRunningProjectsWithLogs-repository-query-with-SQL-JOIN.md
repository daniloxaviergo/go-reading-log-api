---
id: RDL-128
title: >-
  [doc-011 Phase 1] Add GetRunningProjectsWithLogs repository query with SQL
  JOIN
status: To Do
assignee: []
created_date: '2026-04-28 11:16'
labels:
  - feature
  - backend
  - phase-1
dependencies: []
documentation:
  - doc-011
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update internal/adapter/postgres/dashboard_repository.go with GetRunningProjectsWithLogs() method implementing SQL query with JOIN to eager-load first 4 logs per project (ordered by date DESC). Add progress ordering in SQL using CASE statement and handle NULL values with COALESCE.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 SQL query joins projects with logs table
- [ ] #2 Logs limited to first 4 per project ordered by data DESC
- [ ] #3 Progress ordering implemented via SQL CASE statement
- [ ] #4 NULL values handled with COALESCE
- [ ] #5 Query returns all required project and log fields
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
