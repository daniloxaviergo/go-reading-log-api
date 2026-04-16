---
id: RDL-055
title: '[doc-004 Phase 3.2] Document database cleanup process'
status: To Do
assignee:
  - thomas
created_date: '2026-04-15 12:15'
updated_date: '2026-04-16 20:00'
labels:
  - documentation
  - architecture
  - low-priority
dependencies: []
references:
  - 'Decision 4: Time-Based Orphan Detection'
  - 'Decision 5: Prefix-Based Database Selection'
documentation:
  - doc-004
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create comprehensive documentation for the database cleanup process covering the 24-hour orphan detection strategy, database name validation rules, and cleanup SQL patterns. Document the rationale for keeping the per-test database strategy rather than switching to schema reset, and include troubleshooting steps for common cleanup issues.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Document 24-hour orphan detection strategy
- [x] #2 Document database name validation rules
- [x] #3 Document SQL cleanup patterns
- [x] #4 Include troubleshooting steps
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Document database cleanup procedures by expanding the "Cleanup Procedures" section in `AGENTS.md` with detailed technical specifications. The documentation will:
- Explain the time-based orphan detection logic implemented in `cleanup.go`
- Clarify database naming conventions and validation rules
- Provide SQL examples for critical cleanup operations
- Detail why per-test databases are preferred over schema resets (parallel test safety)
- Include troubleshooting guidance for common issues

### 2. Files to Modify

- `AGENTS.md` (update "Cleanup Procedures" section)
- Add new reference file: `docs/database-cleanup.md` (for detailed SQL patterns and decision rationale)

### 3. Dependencies

- Existing code in:
  - `internal/adapter/postgres/cleanup.go` (orphan detection logic)
  - `test/test_helper.go` (database name generation)
  - `Makefile` (test-clean command implementation)
- Decision documents: Decision 4 (Time-Based Orphan Detection) and Decision 5 (Prefix-Based Database Selection)

### 4. Code Patterns

- Follow existing AGENTS.md structure with:
  - Clear section headers for each cleanup component
  - Code blocks for SQL commands and Go code snippets
  - Warning callouts for production database safety
- Use consistent terminology matching the codebase (e.g., "orphaned records" vs "unlinked entries")

### 5. Testing Strategy

- **Accuracy Verification**: Cross-reference documentation against:
  - Actual implementation in `cleanup.go` (time thresholds, foreign key checks)
  - Database name generation logic in `test_helper.go`
  - Makefile commands for test cleanup
- **Peer Review**: Have another developer validate documentation accuracy
- **Real-world Validation**: Test manual cleanup procedures using:
  ```bash
  make test-clean && go run ./cmd/cleanup.go orphan
  ```

### 6. Risks and Considerations

- **Documentation Drift**: New cleanup features may not be reflected in docs if code changes without updates
- **Production Safety**: Must emphasize strict warnings against running cleanup on `reading_log` (production DB)
- **Query Performance**: SQL patterns should include index usage notes to prevent full table scans
- **Parallel Test Conflicts**: Clarify that unique database naming is critical for concurrent test execution
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [x] #2 All integration tests pass use testing-expert subagent for test execution and verification
- [x] #3 go fmt and go vet pass with no errors
- [x] #4 Clean Architecture layers properly followed
- [x] #5 Error responses consistent with existing patterns
- [x] #6 HTTP status codes correct for response type
- [x] #7 Database queries optimized with proper indexes
- [x] #8 Documentation updated in QWEN.md
- [x] #9 New code paths include error path tests
- [ ] #10 HTTP handlers test both success and error responses
- [ ] #11 Integration tests verify actual database interactions
- [ ] #12 Tests use testing-expert subagent for test execution and verification
<!-- DOD:END -->
