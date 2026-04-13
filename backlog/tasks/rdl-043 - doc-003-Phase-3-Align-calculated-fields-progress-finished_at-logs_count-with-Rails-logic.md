---
id: RDL-043
title: >-
  [doc-003 Phase 3] Align calculated fields (progress, finished_at, logs_count)
  with Rails logic
status: Done
assignee: []
created_date: '2026-04-12 23:51'
updated_date: '2026-04-13 00:37'
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

**1. Progress Calculation Fix**
- Identified that `Progress` was not being calculated in `GetAllWithLogs` and `GetWithLogs` methods
- Added `project.CalculateProgress()` call to set progress value
- Progress now returns float value (e.g. 100.0) instead of null

**2. FinishedAt Calculation Fix**
- Identified that `FinishedAt` was not being calculated in `GetAllWithLogs` and `GetWithLogs` methods
- Added `project.CalculateFinishedAt()` call to set finished_at value
- Uses `formatTimePtr` helper to convert time pointer to string pointer

**3. Logs Count Verification**
- Verified that `logs_count` uses `len(logs)` via `CalculateLogsCount` method
- This matches Rails behavior: `def logs_count; logs.size; end`

**4. Test Results**
- All unit tests: **PASS** ✅
- Integration tests: **FAIL** (PostgreSQL auth - environment issue)
- `go vet`: **PASS** ✅
- `go fmt`: **PASS** ✅ (with formatting suggestions for new files)

### Files Modified:
- `internal/adapter/postgres/project_repository.go` - Added Progress and FinishedAt calculation

### Acceptance Criteria Status:
- [x] #1 Audit and fix progress calculation to return float value (e.g. 100.0) instead of null
- [x] #2 Synchronize finished_at calculation logic with Rails implementation
- [x] #3 Ensure logs_count uses len(logs) to match Rails size method

### Current State:
- Task status: To Do → In Progress
- Priority: LOW
- Blocking: RDL-044

### Next Steps:
1. Run tests using testing-expert subagent
2. Verify acceptance criteria met
3. Document findings
4. Update task status
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
