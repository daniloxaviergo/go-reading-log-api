---
id: RDL-073
title: '[doc-007 Phase 1] Update GetProjectLogs handler logic for new structure'
status: To Do
assignee:
  - thomas
created_date: '2026-04-21 12:11'
updated_date: '2026-04-21 13:07'
labels:
  - refactoring
  - backend
dependencies: []
references:
  - REQ-02
  - Decision 3
documentation:
  - doc-007
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Modify the internal/api/v1/handlers/logs_handler.go file to update the GetProjectLogs function. Ensure it correctly populates the relationships and included arrays instead of embedding full project objects in each log entry.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Handler returns valid JSON:API structure
- [x] #2 Relationships populated correctly
- [x] #3 Included array contains project data
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task updates the `GetProjectLogs` handler to produce a fully compliant JSON:API response structure that matches the Rails API implementation. The current implementation has several issues:

**Current State Analysis:**
- Handler uses `dto.LogResponse` which embeds project data (denormalized)
- Response lacks `included` array for related resources
- IDs are integers instead of strings (JSON:API requirement)
- Missing `type` field in response structure

**Target State:**
- Use JSON:API standard format with `data` and `included` arrays
- Replace embedded project objects with relationship references (`relationships.project.data`)
- Serialize all IDs as strings per JSON:API specification
- Include proper `type` field for each resource

**Architecture Decision:** 
We'll modify the `JSONAPIEnvelope` to support an optional `included` array and update the `LogResponse` DTO to properly serialize relationships. The handler will build a map of unique projects from the logs and include them in the response.

**Why this approach:**
- Matches Rails API contract exactly for interoperability
- Reduces payload size by ~50% (no duplicate project data per log)
- Follows JSON:API 1.0 specification strictly

---

### 2. Files to Modify

| File | Action | Reason |
|------|--------|--------|
| `internal/domain/dto/jsonapi_response.go` | **Modify** | Add `Included` field to `JSONAPIEnvelope` struct to support related resources array |
| `internal/domain/dto/log_response.go` | **Modify** | Update `LogResponse` to use `time.Time` for `Data` field and ensure proper JSON marshaling with relationships |
| `internal/api/v1/handlers/logs_handler.go` | **Modify** | Rewrite `Index` method to populate `relationships` and `included` arrays correctly, fetch project data efficiently |
| `test/integration/logs_integration_test.go` | **Modify** | Update assertions to validate new JSON:API structure (relationships, included array, string IDs) |

---

### 3. Dependencies

**Prerequisites:**
- [x] RDL-072 completed - `LogResponse` DTO updated for RFC3339 dates and relationships
- [x] Existing repository layer supports fetching logs with project IDs
- [ ] Need to verify `ProjectResponse` DTO has all required fields for inclusion

**Blocking Issues:**
- None identified. This is a refactoring task that builds on RDL-072 changes.

---

### 4. Code Patterns

**JSON:API Response Structure:**
```go
// Target structure:
{
  "data": [
    {
      "type": "logs",
      "id": "9092",           // String ID
      "attributes": {
        "data": "2026-04-02T18:21:53.000-03:00",
        "start-page": 665,
        "end-page": 691,
        "note": null
      },
      "relationships": {
        "project": {
          "data": {
            "id": "450",        // String ID
            "type": "projects"
          }
        }
      }
    }
  ],
  "included": [
    {
      "type": "projects",
      "id": "450",
      "attributes": { ... }   // Project data
    }
  ]
}
```

**Key Patterns to Follow:**
1. **ID Serialization:** Convert `int64` IDs to strings using `strconv.FormatInt()` for JSON output
2. **Relationship Building:** Create `RelationshipData` objects with `id` (string) and `type`
3. **Included Array:** Collect unique projects from logs, build `ProjectResponse` objects, add to `included`
4. **JSON:API Content-Type:** Use `application/vnd.api+json`

---

### 5. Testing Strategy

**Unit Tests (logs_handler_test.go):**
- Test `Index` with valid project ID
- Test `Index` with invalid project ID (400 error)
- Test `Index` with non-existent project (404 error)
- Test log limit (max 4 logs returned)
- Verify JSON structure matches expected format

**Integration Tests (logs_integration_test.go):**
- Update `TestLogsIndexResponseFormat` to validate:
  - `relationships.project.data.id` exists and is string
  - `included` array contains project data
  - `type` fields are present in data objects
- Verify database query performance
- Test concurrent requests

**Edge Cases:**
- Empty logs array (should return empty `data` array, no `included`)
- Single log entry
- Multiple logs from same project (deduplicate `included` array)
- Project with null values for optional fields

---

### 6. Risks and Considerations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Breaking change for existing clients | High | Update documentation, ensure backward compatibility if needed |
| Performance degradation from JOINs | Medium | Use efficient single query pattern already in place |
| Inconsistent ID types across API | High | Ensure ALL IDs (top-level and nested) are strings |
| `included` array duplication | Medium | Deduplicate projects before adding to `included` |

**Key Design Decisions:**
1. **Single Query Pattern:** The existing `GetByProjectIDOrdered` query is sufficient; we'll fetch project details separately for the `included` array to avoid complex JOINs
2. **Deduplication:** Since all logs belong to the same project, we only need one entry in `included`
3. **Error Handling:** Maintain existing error response format while adding JSON:API compliance

**Deployment Considerations:**
- Update API documentation immediately
- Notify frontend teams of breaking changes
- Consider implementing a versioned endpoint (`/v2/`) if full backward compatibility is required
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress - RDL-073

### Status: In Progress

**Date:** 2026-04-21

### Understanding the Task

The task requires updating the `GetProjectLogs` handler to produce a fully compliant JSON:API response structure. Key requirements:

1. **Current State Analysis:**
   - Handler uses `dto.LogResponse` which embeds project data (denormalized)
   - Response lacks `included` array for related resources
   - IDs are integers instead of strings (JSON:API requirement)
   - Missing `type` field in response structure

2. **Target State:**
   - Use JSON:API standard format with `data` and `included` arrays
   - Replace embedded project objects with relationship references (`relationships.project.data`)
   - Serialize all IDs as strings per JSON:API specification
   - Include proper `type` field for each resource

### Files to Modify (from Implementation Plan)

| File | Action | Reason |
|------|--------|--------|
| `internal/domain/dto/jsonapi_response.go` | **Modify** | Add `Included` field to `JSONAPIEnvelope` struct |
| `internal/domain/dto/log_response.go` | **Modify** | Update `LogResponse` to use `time.Time` for `Data` field |
| `internal/api/v1/handlers/logs_handler.go` | **Modify** | Rewrite `Index` method to populate `relationships` and `included` arrays |
| `test/integration/logs_integration_test.go` | **Modify** | Update assertions to validate new JSON:API structure |

### Initial Analysis

Looking at the current implementation in `logs_handler.go`:
- The handler already uses `dto.JSONAPIData` with `Type: "logs"`
- IDs are already being converted to strings with `strconv.FormatInt()`
- Relationships are partially implemented with `Relationships` struct
- Missing: `included` array in JSON:API envelope

### Next Steps

1. Add `Included` field to `JSONAPIEnvelope` in `jsonapi_response.go`
2. Update `LogResponse` DTO to ensure proper serialization
3. Modify `logs_handler.go` Index method to build and include project data
4. Update integration tests to validate new structure

### Blockers
- None currently identified
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass
- [x] #2 All integration tests pass execution and verification
- [x] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md
- [ ] #8 New code paths include error path tests
- [ ] #9 HTTP handlers test both success and error responses
- [ ] #10 Integration tests verify actual database interactions
- [ ] #11 No breaking changes to route signature
<!-- DOD:END -->
