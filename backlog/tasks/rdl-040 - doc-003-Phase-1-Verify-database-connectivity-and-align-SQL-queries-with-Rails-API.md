---
id: RDL-040
title: >-
  [doc-003 Phase 1] Verify database connectivity and align SQL queries with
  Rails API
status: To Do
assignee: []
created_date: '2026-04-12 23:50'
labels:
  - database
  - query
  - alignment
dependencies: []
references:
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/1'
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/7'
documentation:
  - doc-003
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement FR-001 and FR-007 by verifying the database connection string in .env points to the correct 'reading_log' database and ensuring all SQL queries in the adapter layer strictly replicate Rails ActiveRecord logic to guarantee identical result sets between the Go and Rails APIs.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Verify .env DB_DATABASE is set to 'reading_log'
- [ ] #2 Run integration tests confirming go_count equals rails_count
- [ ] #3 Audit all SQL queries in internal/adapter/postgres/queries.go against Rails ActiveRecord counterparts
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
