---
id: RDL-074
title: '[doc-007 Phase 2] Implement JSON marshaling for relationships and string IDs'
status: To Do
assignee:
  - workflow
created_date: '2026-04-21 12:11'
updated_date: '2026-04-21 13:21'
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

### 2. Files to Modify

| File | Action | Description |
|:-----|:-------|:------------|
| `internal/domain/dto/log_response.go` | **No changes needed** | Already has `RelationshipData`, `Relationships`, and proper struct tags |
| `internal/api/v1/handlers/logs_handler.go` | **Modify** | Update handler to ensure string IDs and populate relationships correctly |
| `test/integration/test_context.go` | **Verify** | Ensure parsing handles JSON:API envelope correctly |
| `test/integration/logs_integration_test.go` | **Verify** | Add specific tests for string ID validation |

**Files Already Correct**:
- `internal/domain/dto/jsonapi_response.go` - Contains `JSONAPIData`, `NewIncludedProject` helpers
- `internal/api/v1/handlers/projects_handler.go` - Reference implementation for JSON:API envelope wrapping

### 3. Dependencies

**Prerequisites (Already Met)**:
- [x] RDL-072 - LogResponse DTO updated for RFC3339 dates and relationships
- [x] RDL-073 - GetProjectLogs handler logic updated for new structure
- [x] `internal/domain/dto/log_response.go` - Contains `RelationshipData` and `Relationships` structs
- [x] `internal/domain/dto/jsonapi_response.go` - Contains JSON:API envelope helpers

**No Additional Setup Required**: The infrastructure is already in place from Phase 1 completion.

### 4. Code Patterns

**Pattern 1: String ID Serialization**
```go
// Use strconv.FormatInt for all IDs in JSON:API responses
ID: strconv.FormatInt(logs[i].ID, 10) // "123" instead of 123
```

**Pattern 2: Relationship Structure**
```go
Relationships: &dto.Relationships{
    Project: &dto.RelationshipData{
        ID:   strconv.FormatInt(project.ID, 10), // String ID
        Type: "projects",                         // Resource type
    },
},
```

**Pattern 3: Included Array Population**
```go
included = append(included, dto.NewIncludedProject(projectResponse))
// Returns: {"type": "projects", "id": "123", "attributes": {...}}
```

**Pattern 4: JSON:API Envelope Wrapping**
```go
envelope := dto.NewJSONAPIEnvelopeWithIncluded(dataObjects, included)
w.Header().Set("Content-Type", "application/vnd.api+json")
json.NewEncoder(w).Encode(envelope)
```

### 5. Testing Strategy

**Unit Tests (Existing - Verify Pass)**:
- `TestLogsHandler_Index` - Empty logs collection
- `TestLogsHandler_IndexWithLogs` - Multiple logs with limit validation
- `TestLogsHandler_IndexWithLessThanLimit` - Fewer than 4 logs
- `TestLogsHandler_IndexWithOneLog` - Single log entry
- `TestFormatTimePtr` - Time formatting helper

**Integration Tests (Existing - Verify Pass)**:
- `TestLogsIndexIntegration` - Full endpoint with database
- `TestLogsIndexEmpty` - No logs scenario
- `TestLogsIndexProjectNotFound` - 404 handling
- `TestLogsIndexInvalidProjectID` - 400 handling
- `TestLogsIndexLimit` - 4-log limit enforcement
- `TestLogsIndexWithLogs` - Logs with notes
- `TestLogsIndexConcurrent` - Concurrent access
- `TestLogsIndexResponseFormat` - JSON:API format validation

**New Tests to Add**:
1. **String ID Validation**: Explicit check that all IDs are strings in serialized JSON
2. **Relationship Structure**: Verify `relationships.project.data.id` exists and is string
3. **Included Array Verification**: Confirm included array contains project data with correct structure

### 6. Risks and Considerations

**Blocking Issues**: None identified. Implementation follows established patterns.

**Trade-offs**:
- **Consistency**: Must ensure all endpoints (projects, logs) use identical ID serialization strategy
- **Backward Compatibility**: JSON:API format is a breaking change from flat JSON, but this is intentional per PRD requirements
- **Performance**: String conversion has negligible overhead; primary concern is response size reduction via relationship references

**Design Decisions**:
1. **No Custom MarshalJSON**: Standard struct tags sufficient for current needs
2. **String IDs Everywhere**: All public-facing IDs use strings to comply with JSON:API spec and avoid JavaScript integer precision issues
3. **Included Array**: Project data included in logs response to enable relationship resolution without additional requests
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
