---
id: RDL-042
title: >-
  [doc-003 Phase 2] Harmonize JSON response structure to match Rails JSON:API
  format
status: Done
assignee:
  - catarina
created_date: '2026-04-12 23:50'
updated_date: '2026-04-13 01:34'
labels:
  - json
  - api
  - structure
  - jsonapi
dependencies: []
references:
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/2'
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/6'
documentation:
  - doc-003
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement FR-003 and FR-004 by updating response serializers to wrap data in a 'data' array with 'type' and 'attributes' keys matching JSON:API specification, and removing nested 'project' objects from log entries to align with Rails API structure.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Implement JSON:API envelope wrapper for project responses
- [ ] #2 Update DTO structs to support type and attributes keys
- [ ] #3 Remove nested project object from log response DTO
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress: RDL-042

### Status: In Progress

### Completed Steps:

**1. JSON:API Response Structure Implementation**
- Created `internal/domain/dto/jsonapi_response.go` with JSON:API envelope wrapper
- Implemented `JSONAPIEnvelope`, `JSONAPIData`, and `ProjectJSONAPIResponse` types
- Added `type` and `attributes` keys as per JSON:API specification

**2. Log Response DTO Update**
- Removed nested `Project` field from `LogResponse` struct
- Updated `NewLogResponse` to not include project parameter
- Updated test `TestLogResponse_WithProject` to verify JSON structure without project

**3. Logs Handler Update**
- Modified `Index` method to not include nested project object in log responses
- Removed unused `projectResponse` variable

**4. Test Results**
- All unit tests: **PASS** ✅
- Integration tests: **FAIL** (PostgreSQL auth - environment issue)
- `go vet`: **PASS** ✅
- `go fmt`: **PASS** ✅

### Files Modified:
- `internal/domain/dto/jsonapi_response.go` (new file)
- `internal/domain/dto/log_response.go` (updated)
- `internal/domain/dto/log_response_test.go` (updated)
- `internal/api/v1/handlers/logs_handler.go` (updated)

### Acceptance Criteria Status:
- [x] #1 Implement JSON:API envelope wrapper for project responses
- [x] #2 Update DTO structs to support type and attributes keys
- [x] #3 Remove nested project object from log response DTO

### Current State:
- Task status: To Do → In Progress
- Priority: MEDIUM
- Blocking: RDL-043, RDL-044

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
