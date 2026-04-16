---
id: RDL-054
title: '[doc-004 Phase 3.1] Update AGENTS.md with cleanup procedures'
status: To Do
assignee:
  - thomas
created_date: '2026-04-15 12:15'
updated_date: '2026-04-16 19:39'
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

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach
Update AGENTS.md by adding a new "Cleanup Procedures" section under "Common Tasks". Break down into four subsections matching acceptance criteria:
- Auto-cleanup mechanism (TestHelper.Close() usage)
- Orphan cleanup process (database adapter function)
- Manual cleanup commands (`make test-clean`, DB recreation steps)
- Parallel test safety (goroutine ID database naming)
Include code snippets, warnings about production DB risks, and cross-references to related tasks.

### 2. Files to Modify
- `/home/danilo/scripts/github/go-reading-log-api-next/AGENTS.md`

### 3. Dependencies
- RDL-048: `TestHelper.Close()` auto-cleanup implementation (already complete)
- RDL-049: Orphaned database cleanup function (implemented in adapter)
- RDL-051: `make test-clean` command added to Makefile
- RDL-052: Goroutine ID-based database naming for parallel tests

### 4. Code Patterns
- Use markdown code blocks with triple backticks for commands and Go snippets
- Format steps as numbered lists for clarity
- Highlight critical warnings in **bold** (e.g., "DO NOT run on production databases")
- Reference task IDs and PRD documents where applicable

### 5. Testing Strategy
- Manual review by team members to verify accuracy of commands and procedures
- Cross-check all examples against current codebase behavior (e.g., verify `make test-clean` exists in Makefile)
- No automated tests required since this is documentation-only change

### 6. Risks and Considerations
- **Critical risk**: Explicitly state that database cleanup commands must ONLY be used on test databases (`reading_log_test`). Include warnings like "WARNING: Running DROP DATABASE on production will cause permanent data loss"
- Differentiate between safe commands (`make test-clean`) vs manual SQL operations
- Ensure all references to task IDs (RDL-048 etc.) and file paths are current
- Verify parallel test safety section aligns with RDL-052 implementation details
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [x] #2 All integration tests pass use testing-expert subagent for test execution and verification
- [x] #3 go fmt and go vet pass with no errors
- [x] #4 Clean Architecture layers properly followed
- [x] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Database queries optimized with proper indexes
- [ ] #8 Documentation updated in QWEN.md
- [ ] #9 New code paths include error path tests
- [ ] #10 HTTP handlers test both success and error responses
- [ ] #11 Integration tests verify actual database interactions
- [ ] #12 Tests use testing-expert subagent for test execution and verification
<!-- DOD:END -->
