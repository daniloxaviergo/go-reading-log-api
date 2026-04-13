---
id: RDL-043
title: >-
  [doc-003 Phase 3] Align calculated fields (progress, finished_at, logs_count)
  with Rails logic
status: To Do
assignee: []
created_date: '2026-04-12 23:51'
updated_date: '2026-04-13 00:24'
labels:
  - calculation
  - logic
  - synchronization
dependencies: []
references:
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/4'
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/5'
documentation:
  - doc-003
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement FR-005 and FR-006 by auditing and synchronizing the calculation logic for progress (percentage), finished_at (completion date), and logs_count (array length) to ensure they match the Rails API's business logic exactly.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Audit and fix progress calculation to return float value (e.g. 100.0) instead of null
- [ ] #2 Synchronize finished_at calculation logic with Rails implementation
- [ ] #3 Ensure logs_count uses len(logs) to match Rails size method
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress: RDL-043

### Status: In Progress

### Completed Steps:

**1. Task Analysis**
- Reviewed RDL-043 task details
- Understanding calculated field requirements (progress, finished_at, logs_count)
- Identifying files that need modification for calculation logic alignment

**2. Files to Review:**
- `internal/domain/models/project.go` - Progress and finished_at calculation
- `internal/domain/dto/project_response.go` - Logs count calculation
- `internal/adapter/postgres/project_repository.go` - Database query logic

**3. Current State:**
- Task status: To Do → In Progress
- Priority: LOW
- Depends on: RDL-041, RDL-042 completion

### Next Steps:
1. Review and audit progress calculation logic
2. Review and synchronize finished_at calculation logic
3. Verify logs_count uses len(logs)
4. Run tests using testing-expert subagent
5. Verify acceptance criteria
6. Document findings
<!-- SECTION:NOTES:END -->

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
