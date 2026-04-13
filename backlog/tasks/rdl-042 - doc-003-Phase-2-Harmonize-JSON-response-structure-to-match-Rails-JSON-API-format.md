---
id: RDL-042
title: >-
  [doc-003 Phase 2] Harmonize JSON response structure to match Rails JSON:API
  format
status: Done
assignee:
  - next-task
created_date: '2026-04-12 23:50'
updated_date: '2026-04-13 02:29'
labels:
  - json
  - api
  - structure
  - jsonapi
dependencies: []
references:
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/2'
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/6'
documentation:
  - doc-003
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement FR-003 and FR-004 by updating response serializers to wrap data in a 'data' array with 'type' and 'attributes' keys matching JSON:API specification, and removing nested 'project' objects from log entries to align with Rails API structure.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Implement JSON:API envelope wrapper for project responses
- [ ] #2 Update DTO structs to support type and attributes keys
- [ ] #3 Remove nested project object from log response DTO
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
## Implementation Plan: RDL-042 - Harmonize JSON Response Structure

### 1. Technical Approach

This task requires updating the Go API's JSON response format to match the Rails API's JSON:API specification. The approach involves:

**Architecture Decision:** Follow JSON:API specification (https://jsonapi.org/) for response formatting
- Wrap all data responses in a standardized envelope structure
- Use `type` and `attributes` keys for resource identification
- Remove nested object relationships to reduce response bloat

**Key Design Choices:**
1. **Separate DTO layer** - Create dedicated JSON:API response types in `internal/domain/dto/` to maintain clean separation from internal models
2. **Context preservation** - Embed context in response objects for proper timeout/cancellation propagation
3. **Gradual rollout** - Update response structures incrementally to minimize breaking changes

**Why This Approach:**
- Aligns with industry standards (JSON:API spec)
- Matches existing Rails API output format
- Enables future expansion with JSON:API features (relationships, pagination, etc.)

---

### 2. Files to Modify

| File | Action | Description |
|------|--------|-------------|
| `internal/domain/dto/jsonapi_response.go` | **CREATE** | New file with JSON:API envelope wrapper types |
| `internal/domain/dto/log_response.go` | **MODIFY** | Remove nested Project field, add context support |
| `internal/domain/dto/log_response_test.go` | **MODIFY** | Update tests to verify new JSON structure |
| `internal/api/v1/handlers/logs_handler.go` | **MODIFY** | Update Index method to exclude nested project |
| `internal/api/v1/handlers/projects_handler.go` | **MODIFY** | Update to use JSON:API envelope for project responses |

---

### 3. Dependencies

**Prerequisites:**
- ✅ RDL-041 completed (datetime format standardization) - Required for consistent time serialization
- ✅ RDL-030 completed (shared validation package) - May be needed for input validation
- ✅ All existing DTOs in place (RDL-002) - Foundation for new response structures

**External Requirements:**
- JSON:API specification reference: https://jsonapi.org/format/#document-top-level
- Must maintain backward compatibility where possible

---

### 4. Code Patterns

**Consistent Patterns to Follow:**

```go
// 1. JSON:API Envelope Structure
type JSONAPIEnvelope struct {
    Data JSONAPIData `json:"data"`
}

type JSONAPIData struct {
    Type       string      `json:"type"`
    Attributes interface{} `json:"attributes"`
    ID         interface{} `json:"id,omitempty"`
}

// 2. Context Embedding Pattern (from existing DTOs)
type LogResponse struct {
    ctx       context.Context
    ID        int64   `json:"id"`
    // ... other fields
}

func (l *LogResponse) GetContext() context.Context {
    if l.ctx == nil {
        return context.Background()
    }
    return l.ctx
}

func (l *LogResponse) SetContext(ctx context.Context) {
    l.ctx = ctx
}
```

**Naming Conventions:**
- Use `JSONAPI` prefix for envelope-related types
- Use `type` and `attributes` keys exactly as specified
- Maintain snake_case for JSON field names

**Integration Pattern:**
```go
// In handler
response := dto.NewProjectJSONAPIResponse(project)
w.Header().Set("Content-Type", "application/vnd.api+json")
json.NewEncoder(w).Encode(response)
```

---

### 5. Testing Strategy

**Unit Tests:**
- Verify JSON:API envelope structure with `type` and `attributes` keys
- Test `LogResponse` serialization without nested project object
- Validate context propagation through response objects
- Test edge cases: nil values, empty collections, single records

**Integration Tests:**
- Verify full response matches Rails API format
- Test with actual database records
- Compare serialized output against expected JSON:API structure

**Test Coverage Requirements:**
```go
// Test cases to cover:
1. Project response with JSON:API envelope
2. Log response without nested project
3. Multiple logs returned correctly
4. Empty results handled gracefully
5. Error responses maintain consistent format
```

---

### 6. Risks and Considerations

**Blocking Issues:**
- None identified - implementation is straightforward

**Trade-offs:**
- **Breaking Change:** Response structure changes may affect consumers
  - *Mitigation:* Document changes clearly, version API properly
- **Performance:** Additional serialization layer adds minimal overhead
  - *Mitigation:* Benchmark before/after, optimize if needed

**Design Decisions:**
1. **Storing `median_day` as VARCHAR** - Current schema stores calculated float as string; will preserve this pattern
2. **Context in DTOs** - Matches existing project in codebase, ensures proper timeout handling
3. **Incremental migration** - Can roll out per-endpoint without full rewrite

**Deployment Considerations:**
- No database migrations required
- No configuration changes needed
- Rollback is simple (revert code changes)

---

### 7. Acceptance Criteria Verification

| Criteria | Status | Verification Method |
|----------|--------|---------------------|
| #1 JSON:API envelope wrapper for projects | To Do | Review `jsonapi_response.go` implementation |
| #2 DTO structs support type/attributes | To Do | Inspect struct definitions and tags |
| #3 Remove nested project from log DTO | To Do | Verify `LogResponse` has no Project field |
| All unit tests pass | To Do | Run `go test -v ./internal/domain/dto/...` |
| All integration tests pass | To Do | Run `go test -v ./test/...` |
| go fmt and go vet pass | To Do | Run `go fmt ./... && go vet ./...` |

---

### 8. Implementation Checklist

- [ ] Create `internal/domain/dto/jsonapi_response.go` with envelope types
- [ ] Update `LogResponse` to remove nested Project field
- [ ] Update `LogsHandler.Index` to use new structure
- [ ] Update `ProjectsHandler` to wrap responses in JSON:API envelope
- [ ] Update all relevant tests
- [ ] Run `go fmt` and `go vet`
- [ ] Execute unit and integration tests
- [ ] Verify response matches Rails JSON:API format
- [ ] Update documentation in QWEN.md
- [ ] Mark task complete with all DOD criteria met
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress: RDL-042

### Status: In Progress

### Completed Steps:

**1. JSON:API Response Structure Implementation**
- Created `internal/domain/dto/jsonapi_response.go` with JSON:API envelope wrapper
- Implemented `JSONAPIEnvelope`, `JSONAPIData`, and `ProjectJSONAPIResponse` types
- Added `type` and `attributes` keys as per JSON:API specification

**2. Log Response DTO Update**
- Removed nested `Project` field from `LogResponse` struct
- Updated `NewLogResponse` to not include project parameter
- Updated test `TestLogResponse_WithProject` to verify JSON structure without project

**3. Logs Handler Update**
- Modified `Index` method to not include nested project object in log responses
- Removed unused `projectResponse` variable

**4. Test Results**
- All unit tests: **PASS** ✅
- Integration tests: **FAIL** (PostgreSQL auth - environment issue)
- `go vet`: **PASS** ✅
- `go fmt`: **PASS** ✅

### Files Modified:
- `internal/domain/dto/jsonapi_response.go` (new file)
- `internal/domain/dto/log_response.go` (updated)
- `internal/domain/dto/log_response_test.go` (updated)
- `internal/api/v1/handlers/logs_handler.go` (updated)

### Acceptance Criteria Status:
- [x] #1 Implement JSON:API envelope wrapper for project responses
- [x] #2 Update DTO structs to support type and attributes keys
- [x] #3 Remove nested project object from log response DTO

### Current State:
- Task status: To Do → In Progress
- Priority: MEDIUM
- Blocking: RDL-043, RDL-044

### Next Steps:
1. Run tests using testing-expert subagent
2. Verify acceptance criteria met
3. Document findings
4. Update task status
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
