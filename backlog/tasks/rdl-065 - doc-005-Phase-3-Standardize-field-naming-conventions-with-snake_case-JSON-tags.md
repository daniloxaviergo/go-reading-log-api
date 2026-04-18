---
id: RDL-065
title: >-
  [doc-005 Phase 3] Standardize field naming conventions with snake_case JSON
  tags
status: To Do
assignee:
  - book
created_date: '2026-04-18 11:47'
updated_date: '2026-04-18 14:21'
labels:
  - phase-3
  - naming-convention
  - json-tags
dependencies: []
references:
  - 'PRD Section: Decision 3'
  - internal/domain/dto/project_response.go
documentation:
  - doc-005
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update internal/domain/dto/project_response.go struct definitions to use snake_case JSON field names via struct tags while maintaining Go convention in code. Ensure all fields have explicit json:"field_name" tags for consistency.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 All DTO structs have explicit JSON field tags
- [ ] #2 Field names follow snake_case convention
- [ ] #3 No kebab-case fields in Go API responses
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
- [ ] #13 go vet reports no struct tag issues
- [ ] #14 Consistent with existing codebase patterns
<!-- DOD:END -->
