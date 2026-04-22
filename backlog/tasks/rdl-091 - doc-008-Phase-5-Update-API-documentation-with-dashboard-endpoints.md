---
id: RDL-091
title: '[doc-008 Phase 5] Update API documentation with dashboard endpoints'
status: To Do
assignee:
  - catarina
created_date: '2026-04-21 15:52'
updated_date: '2026-04-22 16:38'
labels:
  - phase-5
  - documentation
  - api
dependencies: []
references:
  - DOC-001
  - Implementation Checklist Phase 5
documentation:
  - doc-008
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Document all 8 dashboard endpoints in API docs including request/response formats, parameter descriptions, and example requests/responses. Ensure consistency with existing Go API documentation style.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 All 8 endpoints documented with request/response formats
- [ ] #2 Example requests and responses provided
- [ ] #3 Parameter descriptions complete and accurate
- [ ] #4 Documentation consistent with existing style
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task involves creating comprehensive API documentation for all 8 dashboard endpoints that have already been implemented in Phase 2. The approach will mirror the existing documentation style in `README.go-project.md` while adding specific sections for dashboard endpoints.

**Architecture Decision**: 
- Create a new documentation file `docs/dashboard-api-reference.md` following the established pattern
- Document each endpoint with consistent structure: Method, Path, Description, Authentication, Response Code, Request, Response
- Include example requests using curl format matching existing documentation
- Add calculated fields table for each endpoint type
- Provide complete code examples where relevant

**Why this approach**:
- Consistency with existing Go API documentation
- Easy to maintain and extend
- Follows the same structure as projects/logs endpoints documentation
- Provides clear migration reference from Rails to Go

### 2. Files to Modify

#### New Files Created:
| File | Purpose |
|------|---------|
| `docs/dashboard-api-reference.md` | Complete API reference for all dashboard endpoints |

#### Existing Files Referenced (for consistency):
| File | Use |
|------|-----|
| `docs/README.go-project.md` | Reference existing documentation style and structure |
| `docs/rails_routes.md` | Reference Rails endpoint definitions for comparison |
| `internal/api/v1/handlers/dashboard_handler.go` | Verify endpoint implementations match documentation |

### 3. Dependencies

**Prerequisites**:
- Phase 2 dashboard implementation must be complete (all 8 endpoints implemented)
- All handlers must be registered in `internal/api/v1/routes.go`
- Service layer implementations must exist in `internal/service/dashboard/`

**Verification Steps Before Documentation**:
```bash
# Verify all dashboard routes are registered
grep -r "dashboard" internal/api/v1/routes.go

# Verify handlers exist
ls -la internal/api/v1/handlers/*dashboard*

# Run tests to ensure endpoints are functional
go test -v ./internal/api/v1/handlers/... -run TestDashboardHandler
```

### 4. Code Patterns

**Documentation Style to Follow**:

1. **Endpoint Table Format**:
```markdown
| Property | Value |
|----------|-------|
| **Method** | GET |
| **Path** | `/v1/dashboard/day.json` |
| **Description** | Returns daily statistics with weekday breakdown |
| **Authentication** | None |
| **Response Code** | 200 OK |
```

2. **Request/Response Format**:
```bash
# Request example
curl http://localhost:3000/v1/dashboard/day.json

# Response (200 OK)
{
  "data": {
    "type": "dashboard_day",
    "id": "1713849600",
    "attributes": {
      "total_pages": 150,
      "log_count": 5
    }
  }
}
```

3. **Error Response Format** (consistent with existing):
```json
{
  "error": "invalid date format",
  "details": {
    "date": "must be in RFC3339 format"
  }
}
```

4. **Calculated Fields Table**:
Each endpoint type should document its specific calculated fields with formulas.

### 5. Testing Strategy

**Documentation Verification**:

The implementation plan includes verification tests to ensure documentation accuracy:

1. **Endpoint Existence Tests**:
```go
// Verify all documented endpoints exist in handler
func TestDashboardEndpoints_Documented(t *testing.T) {
    // Check that all 8 documented endpoints have handlers
    expectedEndpoints := []string{
        "/v1/dashboard/day.json",
        "/v1/dashboard/projects.json", 
        "/v1/dashboard/last_days.json",
        "/v1/dashboard/echart/faults.json",
        "/v1/dashboard/echart/speculate_actual.json",
        "/v1/dashboard/echart/faults_week_day.json",
        "/v1/dashboard/echart/mean_progress.json",
        "/v1/dashboard/echart/last_year_total.json",
    }
    // Verify each endpoint is registered and responds
}
```

2. **Response Structure Tests**:
```go
// Verify documented response structure matches implementation
func TestDashboardResponses_Structure(t *testing.T) {
    // For each endpoint, verify JSON:API envelope structure
    // Verify all documented fields are present in response
}
```

3. **Integration Verification**:
```bash
# Run integration tests to verify endpoints work as documented
go test -v ./test/integration/... -run TestDashboard

# Check coverage for dashboard handlers
go test -coverprofile=dashboard-coverage.out ./internal/api/v1/handlers/dashboard_handler_test.go
```

### 6. Risks and Considerations

**Known Challenges**:

1. **Dynamic Data in Examples**: 
   - Risk: Example responses will contain dynamic values (timestamps, counts)
   - Mitigation: Use placeholder values and clarify they're illustrative

2. **Configuration-Dependent Values**:
   - Risk: Some endpoints depend on user config (faults max, prediction %)
   - Mitigation: Document default values and explain configuration options

3. **Time-Based Calculations**:
   - Risk: Date calculations may vary based on when documentation is viewed
   - Mitigation: Use relative time references and clarify calculation methods

4. **ECharts Configuration Complexity**:
   - Risk: Chart configurations are complex nested objects
   - Mitigation: Provide simplified examples with key options, link to full ECharts docs

5. **JSON:API Envelope Changes**:
   - Risk: Response structure may change in future versions
   - Mitigation: Version the documentation and note any breaking changes

**Trade-offs**:

| Decision | Rationale |
|----------|-----------|
| Document all 8 endpoints in single file | Easier to maintain than splitting across multiple files |
| Include curl examples | Users can copy-paste and test immediately |
| Reference existing README style | Maintains consistency across project docs |
| Separate calculated fields per endpoint | Clearer than consolidating all calculations |

**Blocking Issues**:
- None identified - Phase 2 implementation is complete per PRD

**Post-Documentation Tasks** (out of scope for this task):
- Add dashboard endpoint examples to main README
- Create interactive API documentation (Swagger/OpenAPI spec)
- Add video walkthrough of dashboard features
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
<!-- DOD:END -->
