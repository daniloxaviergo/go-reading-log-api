---
id: RDL-074
title: '[doc-007 Phase 2] Implement JSON marshaling for relationships and string IDs'
status: To Do
assignee:
  - catarina
created_date: '2026-04-21 12:11'
updated_date: '2026-04-21 13:18'
labels:
  - serialization
  - backend
dependencies: []
references:
  - REQ-03
  - Decision 4
documentation:
  - doc-007
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement custom JSON marshaling logic to ensure IDs are serialized as strings (JSON:API compliant) and relationships/project data are correctly structured in the response payload.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Top-level ID is string
- [ ] #2 Relationships.project.data.id is string
- [ ] #3 included array populated correctly
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task implements JSON:API compliant serialization for the logs endpoint, specifically focusing on relationship data and string ID formatting.

**Architecture Decision**: The implementation uses Go's standard `encoding/json` package with struct tags for field naming and custom logic for type conversion. No custom `MarshalJSON()` methods are required because:
- Standard library handles string/integer conversion via `strconv.FormatInt()`
- Struct tags provide clean JSON field mapping
- The approach aligns with the existing codebase pattern used in `ProjectResponse`

**Key Components**:
- **RelationshipData struct**: Simple container for relationship metadata (ID, Type)
- **Relationships struct**: Container for project relationship reference
- **LogResponse update**: Add Relationships field and ensure ID is serialized as string
- **Handler logic**: Populate relationships with string IDs and build included array

**Why This Approach**:
- Minimal code changes required
- Leverages existing DTO structures
- No breaking changes to internal domain models
- Compatible with existing test infrastructure
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass
- [ ] #2 All integration tests pass execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md
- [ ] #8 New code paths include error path tests
- [ ] #9 HTTP handlers test both success and error responses
- [ ] #10 Integration tests verify actual database interactions
- [ ] #11 Benchmark tests pass
<!-- DOD:END -->
