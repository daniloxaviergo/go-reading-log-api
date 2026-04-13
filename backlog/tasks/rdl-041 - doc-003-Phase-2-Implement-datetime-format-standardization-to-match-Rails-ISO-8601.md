---
id: RDL-041
title: >-
  [doc-003 Phase 2] Implement datetime format standardization to match Rails ISO
  8601
status: To Do
assignee: []
created_date: '2026-04-12 23:50'
updated_date: '2026-04-13 00:17'
labels:
  - datetime
  - format
  - standardization
dependencies: []
references:
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/3'
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/7'
documentation:
  - doc-003
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement FR-002 by updating all datetime field marshaling in the Go API to use the exact ISO 8601 format with timezone offset (+00:00) matching the Rails API, including handling both 'Z' and '+00:00' suffixes during unmarshaling to ensure client compatibility.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Update all time.Time fields in DTOs to use +00:00 format
- [ ] #2 Implement custom MarshalJSON methods for datetime fields
- [ ] #3 Update unmarshaling to accept both Z and +00:00 formats
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress: RDL-041

### Status: In Progress

### Completed Steps:

**1. Task Analysis**
- Reviewed RDL-041 task details
- Understanding datetime format requirements for Rails ISO 8601 compatibility
- Identifying files that need modification for datetime format standardization

**2. Files to Modify:**
- `internal/domain/dto/*.go` - Update time.Time field marshaling
- Implement custom MarshalJSON methods for datetime fields
- Update unmarshaling to accept both Z and +00:00 formats

**3. Current State:**
- Task status: To Do → In Progress
- Priority: HIGH
- Blocking: RDL-042, RDL-043, RDL-044 (Phase 2 tasks depend on this)

### Next Steps:
1. Update DTOs with custom datetime marshaling
2. Implement ISO 8601 format with timezone offset (+00:00)
3. Run tests using testing-expert subagent
4. Verify acceptance criteria
5. Document findings
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
