---
id: RDL-051
title: '[doc-004 Phase 1.4] Add make test-clean command'
status: To Do
assignee:
  - catarina
created_date: '2026-04-15 12:14'
updated_date: '2026-04-16 00:57'
labels:
  - build
  - automation
  - medium-priority
dependencies: []
references:
  - 'R4: Make Command for Manual Cleanup'
documentation:
  - doc-004
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add a make test-clean command to the Makefile that provides a manual cleanup mechanism for orphaned test databases. The command should execute a standalone cleanup script, provide progress feedback during execution, and handle errors gracefully without crashing. Include a test-cleanup alias for convenience. Ensure the Makefile targets use colorized output consistent with existing commands.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Command is available in Makefile
- [ ] #2 It drops all orphaned test databases
- [ ] #3 It provides progress feedback
- [ ] #4 It handles errors gracefully
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
