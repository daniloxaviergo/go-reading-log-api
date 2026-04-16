---
id: RDL-056
title: '[doc-004 Phase 3.3] Create quick reference guide'
status: Done
assignee:
  - workflow
created_date: '2026-04-15 12:15'
updated_date: '2026-04-16 20:18'
labels:
  - documentation
  - reference
  - low-priority
dependencies: []
references:
  - 'Step 3.3: Create quick reference guide'
documentation:
  - doc-004
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create a quick reference guide for developers covering all test database cleanup commands, validation rules, and common operations. Include examples of manual cleanup using make test-clean, checking for orphaned databases, and troubleshooting common issues. Ensure the guide is concise and easy to reference during development.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Commands reference with examples
- [ ] #2 Validation rules summary
- [ ] #3 Troubleshooting section
- [ ] #4 Quick lookup format
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Create a dedicated "Quick Reference Guide" section in `QWEN.md` consolidating database cleanup commands, validation rules, and troubleshooting steps into a concise, scannable format. This section will be structured under the "Common Tasks" header to maintain existing documentation flow. The approach prioritizes:
- **Table-driven summaries** for validation rules (easier scanning than paragraphs)
- **Code-block examples** for all CLI/SQL commands with clear labels
- **Warning callouts** for critical safety considerations (e.g., production DB risks)

This structure aligns with the project's existing documentation style while improving usability for developers needing quick reference during development.

### 2. Files to Modify

- `QWEN.md` (only file)

### 3. Dependencies

- None. This is purely a documentation update requiring no code changes or dependencies.

### 4. Code Patterns

- **Header structure**: Use `###` for section headers matching existing QWEN.md conventions
- **Validation rules table**: Format with columns: `Table`, `Field`, `Validation Rule`
- **Command examples**: Always wrap in ```bash or ```sql code blocks with explicit labels
- **Critical warnings**: Use bolded "⚠️ WARNING" format as seen in existing documentation

### 5. Testing Strategy

- **Manual verification**:
  - Cross-check all SQL commands against current database schema (projects/logs tables)
  - Verify `make test-clean` behavior matches actual Makefile implementation
  - Validate validation rules against implemented code logic (e.g., RDL-032, RDL-033 tasks)
- **Documentation review**: Ensure all examples are syntactically correct and contextually accurate

### 6. Risks and Considerations

- **Production database safety**: Explicitly warn that cleanup commands only apply to `reading_log_test` databases
- **Makefile dependency**: Confirm `make test-clean` exists in project's Makefile before documenting it
- **Validation rule accuracy**: Double-check rules against actual code (e.g., `page <= total_page` must be enforced in application logic)
- **Format consistency**: Maintain existing markdown style (tables, code blocks) to avoid breaking documentation flow
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added Quick Reference Guide section in QWEN.md under Common Tasks. Includes database cleanup commands (make test-clean, manual SQL cleanup), validation rules table for projects/logs tables, and troubleshooting steps for test database issues. All tests passed with cached results. No code changes introduced; documentation only update.
<!-- SECTION:FINAL_SUMMARY:END -->

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
