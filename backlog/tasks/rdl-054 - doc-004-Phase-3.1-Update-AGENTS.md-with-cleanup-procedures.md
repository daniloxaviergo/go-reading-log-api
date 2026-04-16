---
id: RDL-054
title: '[doc-004 Phase 3.1] Update AGENTS.md with cleanup procedures'
status: To Do
assignee: []
created_date: '2026-04-15 12:15'
updated_date: '2026-04-16 11:03'
labels:
  - documentation
  - guide
  - low-priority
dependencies: []
references:
  - 'Step 3.1: Update AGENTS.md'
documentation:
  - doc-004
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update the AGENTS.md documentation file to include the new test database cleanup procedures. Document the auto-cleanup mechanism, orphan cleanup process, manual cleanup commands, and parallel test safety measures. Include examples of how developers should use the new cleanup functionality and reference the relevant PRD for implementation details.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Document auto-cleanup mechanism
- [x] #2 Document orphan cleanup process
- [x] #3 Document manual cleanup commands
- [x] #4 Include parallel test safety measures
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress - RDL-054: Update AGENTS.md with cleanup procedures

### Status: In Progress
### Date: 2026-04-16

### What Was Done:
1. **Read the current AGENTS.md file** to understand the existing structure
2. **Identified sections to add:**
   - Test database cleanup procedures
   - Auto-cleanup mechanism details
   - Orphan cleanup process
   - Manual cleanup commands
   - Parallel test safety measures

3. **Drafted documentation updates:**
   - Added new section "Test Database Cleanup"
   - Documented auto-cleanup mechanism using defer
   - Documented orphan cleanup function
   - Documented manual cleanup commands (make test-clean)
   - Documented parallel test safety with goroutine ID

### Next Steps:
1. Verify acceptance criteria are met
2. Check Definition of Done items
3. Finalize task documentation

### Blockers:
- None identified
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [x] #2 All integration tests pass use testing-expert subagent for test execution and verification
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
