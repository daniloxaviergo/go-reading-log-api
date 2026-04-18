---
id: RDL-064
title: '[doc-005 Phase 3] Implement JSON:API response wrapper for v1 endpoints'
status: To Do
assignee:
  - thomas
created_date: '2026-04-18 11:47'
updated_date: '2026-04-18 13:36'
labels:
  - phase-3
  - json-api
  - breaking-change
dependencies: []
references:
  - 'PRD Section: Decision 2'
  - internal/api/v1/handlers/projects_handler.go
documentation:
  - doc-005
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement JSON:API response wrapper structure in internal/api/v1/handlers/projects_handler.go. The response must use the root wrapper format {data: {...}} with ID as string type according to JSON:API 1.0 specification.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 JSON:API wrapper format implemented for v1 endpoints
- [ ] #2 ID field serialized as string type
- [ ] #3 AC-REQ-004.1 verified: Response has data/attributes structure
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task requires updating the Go API's v1 endpoints to return JSON:API formatted responses that match the Rails API specification. The approach involves:

**Architecture Decision:** Leverage existing JSON:API wrapper types in `internal/domain/dto/jsonapi_response.go` and integrate them into HTTP handlers.

**Key Design Choices:**
1. **Reuse existing DTOs** - The `ProjectJSONAPIResponse` type already exists and wraps `ProjectResponse`; we'll extend it to support collections
2. **Handler-level wrapping** - Wrap response data at the handler level using `JSONAPIEnvelope` before serialization
3. **ID as string** - Convert integer IDs to strings per JSON:API 1.0 specification using `strconv.FormatInt`
4. **Content-Type header** - Use `application/vnd.api+json` for JSON:API compliant responses
5. **Collection handling** - For list endpoints, wrap array items in individual `JSONAPIData` objects

**Why This Approach:**
- Minimal code duplication by reusing existing DTO infrastructure
- Clean separation between domain models and API response format
- Easy to maintain consistency with Rails JSON:API output
- Non-breaking change possible (can support both formats if needed)

---

### 2. Files to Modify

| File | Action | Description |
|------|--------|-------------|
| `internal/api/v1/handlers/projects_handler.go` | **MODIFY** | Update Index, Show, Create methods to wrap responses in JSON:API envelope; convert ID to string |
| `internal/api/v1/handlers/logs_handler.go` | **MODIFY** | Update Index method to wrap log responses in JSON:API envelope |
| `internal/domain/dto/jsonapi_response.go` | **MODIFY** | Add support for collections (arrays); ensure ID is serialized as string |
| `internal/api/v1/routes.go` | **OPTIONAL** | May need to add versioning or content-negotiation support |
| `test/compare_responses.sh` | **MODIFY** | Update comparison logic to expect JSON:API envelope structure |
| `docs/api-changes.md` | **CREATE** | Document JSON:API response format changes for client migration |

---

### 3. Dependencies

**Prerequisites:**
- ✅ RDL-042 completed - JSON:API wrapper types already exist in `internal/domain/dto/jsonapi_response.go`
- ✅ RDL-063 completed - `median_day` field included in ProjectResponse DTO
- ✅ RDL-062 completed - `CalculateFinishedAt` logic implemented
- ✅ RDL-061 completed - Timezone configuration support added

**External Requirements:**
- JSON:API 1.0 specification reference: https://jsonapi.org/format/
- Must match Rails API response structure for compatibility

---

### 4. Code Patterns

**Consistent Patterns to Follow:**

```go
// 1. JSON:API Envelope Structure (already exists)
type JSONAPIEnvelope struct {
    Data JSONAPIData `json:"data"`
}

type JSONAPIData struct {
    Type       string      `json:"type"`
    Attributes interface{} `json:"attributes"`
    ID         interface{} `json:"id,omitempty"`
}

// 2. Pattern for single resource response (Show/Get)
func (h *ProjectsHandler) Show(w http.ResponseWriter, r *http.Request) {
    // ... get project from repo ...
    
    // Wrap in JSON:API envelope
    envelope := dto.NewJSONAPIEnvelope(dto.JSONAPIData{
        Type:       "projects",
        ID:         strconv.FormatInt(project.ID, 10), // ID as string per JSON:API spec
        Attributes: project,
    })
    
    w.Header().Set("Content-Type", "application/vnd.api+json")
    json.NewEncoder(w).Encode(envelope)
}

// 3. Pattern for collection response (Index/List)
func (h *ProjectsHandler) Index(w http.ResponseWriter, r *http.Request) {
    // ... get projects from repo ...
    
    // Convert each project to JSON:API data object
    dataObjects := make([]dto.JSONAPIData, len(projects))
    for i, p := range projects {
        dataObjects[i] = dto.JSONAPIData{
            Type:       "projects",
            ID:         strconv.FormatInt(p.ID, 10), // ID as string
            Attributes: p,
        }
    }
    
    // Wrap collection in envelope
    envelope := dto.NewJSONAPIEnvelope(dataObjects)
    
    w.Header().Set("Content-Type", "application/vnd.api+json")
    json.NewEncoder(w).Encode(envelope)
}

// 4. Pattern for Create (single resource, 201 status)
func (h *ProjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
    // ... create project ...
    
    envelope := dto.NewJSONAPIEnvelope(dto.JSONAPIData{
        Type:       "projects",
        ID:         strconv.FormatInt(createdProject.ID, 10),
        Attributes: createdProject,
    })
    
    w.Header().Set("Content-Type", "application/vnd.api+json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(envelope)
}
```

**Naming Conventions:**
- Use `type` and `attributes` keys exactly as specified
- Keep snake_case for JSON field names via struct tags
- Use `application/vnd.api+json` Content-Type header

**Integration Pattern:**
```go
// In handler methods, after getting data from repo:
response := dto.NewProjectJSONAPIResponse(project)
// Or for collections, build array of JSONAPIData
```

---

### 5. Testing Strategy

**Unit Tests to Add/Update:**

```go
// Test JSON:API envelope structure for single resource
func TestProjectsHandler_Show_JSONAPI(t *testing.T) {
    // ... setup with mock repo ...
    
    handler.Show(w, req)
    
    var envelope dto.JSONAPIEnvelope
    json.NewDecoder(w.Body).Decode(&envelope)
    
    // Verify envelope structure
    if envelope.Data.Type != "projects" {
        t.Errorf("Expected type 'projects', got '%s'", envelope.Data.Type)
    }
    
    // Verify ID is string
    if _, ok := envelope.Data.ID.(string); !ok {
        t.Error("Expected ID to be string type")
    }
    
    // Verify attributes contain expected fields
    attrs, ok := envelope.Data.Attributes.(*dto.ProjectResponse)
    if !ok {
        t.Fatal("Expected Attributes to be ProjectResponse")
    }
    
    if attrs.Name != "Expected Name" {
        t.Errorf("Expected name 'Expected Name', got '%s'", attrs.Name)
    }
}

// Test JSON:API envelope structure for collection
func TestProjectsHandler_Index_JSONAPI(t *testing.T) {
    // ... setup with multiple projects in mock repo ...
    
    handler.Index(w, req)
    
    var envelope dto.JSONAPIEnvelope
    json.NewDecoder(w.Body).Decode(&envelope)
    
    // Verify data is array of JSONAPIData
    dataArray, ok := envelope.Data.Attributes.([]dto.JSONAPIData)
    if !ok {
        t.Fatal("Expected Attributes to be array of JSONAPIData")
    }
    
    if len(dataArray) != expectedCount {
        t.Errorf("Expected %d projects, got %d", expectedCount, len(dataArray))
    }
    
    // Verify all IDs are strings
    for _, item := range dataArray {
        if _, ok := item.ID.(string); !ok {
            t.Error("All IDs must be strings")
        }
    }
}
```

**Integration Tests:**
- Compare full response against Rails API format using `test/compare_responses.sh`
- Verify all calculated fields (progress, finished_at, median_day, days_unreading)
- Test with actual database records via `test.SetupTestDB()`

**Test Coverage Requirements:**
```go
// Test cases to cover:
1. Single project response with JSON:API envelope
2. Collection response with JSON:API envelope  
3. ID serialized as string type
4. Empty collection handled gracefully
5. Error responses maintain consistent format
6. All calculated fields present in attributes
7. Content-Type header is application/vnd.api+json
8. HTTP status codes correct (200 for GET, 201 for POST)
9. Error responses don't use JSON:API envelope (keep existing format)
```

---

### 6. Risks and Considerations

**Blocking Issues:**
- None identified - implementation is straightforward using existing DTOs

**Trade-offs:**
1. **Breaking Change:** Response structure changes from flat JSON to JSON:API envelope
   - *Mitigation:* Document clearly in migration guide; consider versioning strategy
2. **Performance:** Minimal overhead from additional wrapper struct
   - *Mitigation:* Benchmark to ensure < 100ms impact
3. **Client Compatibility:** Existing clients may need updates
   - *Mitigation:* Clear migration documentation, deprecation timeline

**Design Decisions:**
1. **ID as String** - Required by JSON:API 1.0 spec; use `strconv.FormatInt` for conversion
2. **Content-Type** - Use `application/vnd.api+json` to indicate JSON:API compliance
3. **Envelope Structure** - Match Rails Active Model Serializers format exactly
4. **Error Responses** - Keep existing error format (don't wrap in JSON:API) for simplicity

**Deployment Considerations:**
- No database migrations required
- No configuration changes needed
- Rollback is simple (revert code changes)
- Consider gradual rollout or feature flag for backward compatibility

---

### 7. Acceptance Criteria Verification

| Criteria | Status | Verification Method |
|----------|--------|---------------------|
| #1 JSON:API wrapper format implemented | To Do | Review handler code uses `JSONAPIEnvelope` |
| #2 ID field serialized as string | To Do | Verify `strconv.FormatInt` used for IDs |
| #3 AC-REQ-004.1 verified | To Do | Compare response structure against PRD spec |
| All unit tests pass | To Do | Run `go test -v ./internal/api/v1/handlers/...` |
| All integration tests pass | To Do | Run `go test -v ./test/...` |
| go fmt and go vet pass | To Do | Run `go fmt ./... && go vet ./...` |

---

### 8. Implementation Checklist

- [ ] **Step 1:** Update `internal/api/v1/handlers/projects_handler.go`
  - [ ] Modify `Index` to wrap response in JSON:API envelope
  - [ ] Modify `Show` to wrap response in JSON:API envelope  
  - [ ] Modify `Create` to wrap response in JSON:API envelope
  - [ ] Convert all integer IDs to strings using `strconv.FormatInt`
  - [ ] Set `Content-Type: application/vnd.api+json`

- [ ] **Step 2:** Update `internal/api/v1/handlers/logs_handler.go`
  - [ ] Modify `Index` to wrap response in JSON:API envelope
  - [ ] Convert log IDs to strings

- [ ] **Step 3:** Verify/extend `internal/domain/dto/jsonapi_response.go`
  - [ ] Ensure `ProjectJSONAPIResponse` works correctly
  - [ ] Add collection support if needed
  - [ ] Verify string ID serialization

- [ ] **Step 4:** Update tests
  - [ ] Update `projects_handler_test.go` with JSON:API verification
  - [ ] Update `logs_handler_test.go` with JSON:API verification
  - [ ] Update `test/compare_responses.sh` to expect envelope structure

- [ ] **Step 5:** Quality checks
  - [ ] Run `go fmt ./...`
  - [ ] Run `go vet ./...`
  - [ ] Run unit tests: `go test -v ./internal/api/v1/handlers/...`
  - [ ] Run integration tests: `go test -v ./test/...`

- [ ] **Step 6:** Documentation
  - [ ] Update QWEN.md with JSON:API response format
  - [ ] Create migration guide for API clients
  - [ ] Document breaking changes

- [ ] **Step 7:** Verification
  - [ ] Compare Go API response against Rails API format
  - [ ] Verify all acceptance criteria met
  - [ ] Confirm Definition of Done items completed
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress - RDL-064

### Status: In Progress

I've reviewed the task requirements and existing codebase. The task requires implementing JSON:API response wrapper for v1 endpoints following Decision 2 in doc-005.

**Key Requirements:**
1. Implement JSON:API wrapper format `{data: {...}}` for v1 endpoints
2. Serialize ID field as string type per JSON:API 1.0 specification
3. Update projects_handler.go and logs_handler.go

**Existing Infrastructure:**
- `internal/domain/dto/jsonapi_response.go` already exists with `JSONAPIEnvelope`, `JSONAPIData`, and `ProjectJSONAPIResponse` types
- Need to integrate these into HTTP handlers

**Implementation Plan:**
1. Update `projects_handler.go` - Index, Show, Create methods to wrap responses
2. Update `logs_handler.go` - Index method to wrap responses
3. Ensure ID conversion to string using `strconv.FormatInt`
4. Set Content-Type header to `application/vnd.api+json`

**Next Steps:**
- Implement JSON:API wrapping in projects handler
- Implement JSON:API wrapping in logs handler  
- Update tests to verify JSON:API structure
- Run go fmt and go vet
- Execute all tests

### Blockers/Notes:
- None identified - implementation is straightforward using existing DTOs
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
- [ ] #13 Breaking change documented in migration guide
- [ ] #14 Versioning strategy defined
<!-- DOD:END -->
