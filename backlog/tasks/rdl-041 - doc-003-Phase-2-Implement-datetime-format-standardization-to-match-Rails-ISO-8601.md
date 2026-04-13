---
id: RDL-041
title: >-
  [doc-003 Phase 2] Implement datetime format standardization to match Rails ISO
  8601
status: Done
assignee:
  - thomas
created_date: '2026-04-12 23:50'
updated_date: '2026-04-13 01:13'
labels:
  - datetime
  - format
  - standardization
dependencies: []
references:
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/3'
  - 'https://github.com/go-reading-log-api-next/go-reading-log-api-next/issues/7'
documentation:
  - doc-003
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement FR-002 by updating all datetime field marshaling in the Go API to use the exact ISO 8601 format with timezone offset (+00:00) matching the Rails API, including handling both 'Z' and '+00:00' suffixes during unmarshaling to ensure client compatibility.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Update all time.Time fields in DTOs to use +00:00 format
- [ ] #2 Implement custom MarshalJSON methods for datetime fields
- [ ] #3 Update unmarshaling to accept both Z and +00:00 formats
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

**Current State:**
- Go API uses `time.Time` fields in models, serialized to JSON strings
- Current format: `time.RFC3339` produces `2026-04-13T00:00:00Z` format
- Rails API uses ISO 8601 with explicit timezone offset `+00:00`
- Need to match Rails format exactly: `2026-04-13T00:00:00+00:00`

**Technical Challenge:**
Go's `time.RFC3339` constant produces `Z` suffix for UTC, while Rails uses `+00:00` suffix. These are semantically equivalent but string-incompatible for client parsing.

**Solution:**
1. Create custom datetime formatting that produces `+00:00` suffix
2. Update all `time.Time` → string conversion to use custom format
3. Update unmarshaling to accept both `Z` and `+00:00` formats
4. Apply to all datetime fields: `started_at`, `created_at`, `updated_at`, `finished_at`, `data`

**Format Strategy:**
- Output: `2026-04-13T00:00:00+00:00` (ISO 8601 with explicit offset)
- Input: Accept both `Z` and `+00:00` (for backward/forward compatibility)
- Implementation: Custom `MarshalJSON` on DTOs or helper functions

**Why Custom MarshalJSON:**
- Centralized format control
- Consistent behavior across all handlers
- Easy to update if format changes again
- Follows existing pattern in codebase

### 2. Files to Modify

**Core Implementation Files:**

| File | Changes | Reason |
|------|---------|--------|
| `internal/domain/dto/project_response.go` | Add custom `MarshalJSON` with `+00:00` format | Project datetime fields need standardization |
| `internal/domain/dto/log_response.go` | Add custom `MarshalJSON` with `+00:00` format | Log datetime fields need standardization |
| `internal/domain/dto/jsonapi_response.go` | Update JSON:API envelope datetime handling | Nested objects also need format consistency |
| `internal/api/v1/handlers/projects_handler.go` | Update `started_at` parsing to accept `+00:00` | Input validation needs to handle new format |
| `internal/api/v1/handlers/logs_handler.go` | Update `data` parsing to accept `+00:00` | Input validation needs to handle new format |

**Test Files:**

| File | Changes | Reason |
|------|---------|--------|
| `internal/domain/dto/project_response_test.go` | Update tests to verify `+00:00` format | Ensure format compliance |
| `internal/domain/dto/log_response_test.go` | Update tests to verify `+00:00` format | Ensure format compliance |
| `internal/api/v1/handlers/projects_handler_test.go` | Update test data to use `+00:00` format | Test data consistency |
| `internal/api/v1/handlers/logs_handler_test.go` | Update test data to use `+00:00` format | Test data consistency |

**No Changes Required:**
- `internal/domain/models/*.go` - Models keep `time.Time` types
- `internal/adapter/postgres/*.go` - Database layer unchanged
- `internal/repository/*.go` - Repository interfaces unchanged

### 3. Dependencies

**Prerequisites:**
1. **RDL-040 Complete**: Database connectivity verified (Blocker)
2. **RDL-019 Complete**: RFC3339 format work provides foundation
3. **Rails API Access**: Need access to Rails API for format verification
4. **Go 1.25.7**: Current project version with full time package support

**Blocking Issues:**
- None identified - task is self-contained

**Setup Requirements:**
- Ensure `TZ=UTC` environment variable set for consistent timezone handling
- Verify Rails API is running to capture exact datetime format
- Run `go mod tidy` after any dependency changes

### 4. Code Patterns

**Pattern 1: Custom MarshalJSON for Datetime Fields**

```go
type ProjectResponse struct {
    // ... existing fields ...
    StartedAt  *string `json:"started_at"`
    FinishedAt *string `json:"finished_at"`
}

func (p *ProjectResponse) MarshalJSON() ([]byte, error) {
    // Custom implementation that ensures +00:00 format
    // Uses time.Time.Format with custom layout
}
```

**Pattern 2: ISO 8601 Format with Offset**

```go
// Format that produces: 2026-04-13T00:00:00+00:00
layout := "2006-01-02T15:04:05-07:00"
formatted := t.Format(layout) // For UTC: 2026-04-13T00:00:00+00:00
```

**Pattern 3: Dual Format Unmarshaling**

```go
func parseDateTime(s string) (*time.Time, error) {
    // Try Z suffix format first
    if t, err := time.Parse(time.RFC3339, s); err == nil {
        return &t, nil
    }
    // Try +00:00 format
    if t, err := time.Parse("2006-01-02T15:04:05-07:00", s); err == nil {
        return &t, nil
    }
    return nil, fmt.Errorf("invalid datetime format: %s", s)
}
```

**Pattern 4: Consistent Helper Functions**

```go
// datetime_formatter.go
package dto

// FormatTimeRFC3339 formats time.Time to ISO 8601 with +00:00 offset
func FormatTimeRFC3339(t time.Time) string {
    return t.Format("2006-01-02T15:04:05+00:00")
}

// FormatTimePtrRFC3339 formats *time.Time to ISO 8601 with +00:00 offset
func FormatTimePtrRFC3339(t *time.Time) *string {
    if t == nil {
        return nil
    }
    s := t.Format("2006-01-02T15:04:05+00:00")
    return &s
}
```

**Conventions to Follow:**
1. All datetime strings in JSON must use `+00:00` suffix (never `Z`)
2. Unmarshaling must accept both `Z` and `+00:00` for compatibility
3. Keep `time.Time` in models, format at DTO boundary
4. Use pointer types for optional datetimes to ensure JSON null
5. Follow existing naming conventions (`FormatTimePtrRFC3339`)

### 5. Testing Strategy

**Unit Tests:**

| Test | Description | Pass Criteria |
|------|-------------|---------------|
| `TestProjectResponse_JSON_RailsFormat` | Verify JSON output uses `+00:00` | String contains `+00:00`, not `Z` |
| `TestLogResponse_JSON_RailsFormat` | Verify log JSON output uses `+00:00` | String contains `+00:00`, not `Z` |
| `TestParseDateTime_ZSuffix` | Verify `Z` suffix parsing works | Returns valid `time.Time` |
| `TestParseDateTime_OffsetSuffix` | Verify `+00:00` suffix parsing works | Returns valid `time.Time` |
| `TestFormatTimePtrRFC3339_Nil` | Verify nil handling | Returns `nil` |
| `TestFormatTimePtrRFC3339_Valid` | Verify valid time formatting | Returns `+00:00` formatted string |

**Integration Tests:**

| Test | Description | Pass Criteria |
|------|-------------|---------------|
| `TestProjectsHandler_Index_DateFormat` | Full endpoint test | Response matches Rails format exactly |
| `TestProjectsHandler_Show_DateFormat` | Single project test | All datetime fields use `+00:00` |
| `TestLogsHandler_Index_DateFormat` | Logs endpoint test | `data` field uses `+00:00` format |
| `TestDateTime_RoundTrip` | Marshal + Unmarshal | Original and restored times match |

**Test Execution:**
```bash
# Run all datetime-related tests
go test -v ./internal/domain/dto/... -run "DateTime|Format"
go test -v ./internal/api/v1/handlers/... -run "DateTime|Format"

# Run with coverage
go test -coverprofile=datetime-coverage.out ./internal/domain/dto/...
go tool cover -html=datetime-coverage.out

# Run all tests
go test ./...
```

**Comparison with Rails API:**
```bash
# Capture Rails API response
curl -s http://localhost:3001/api/v1/projects.json | jq '.[0].started_at' > rails_started_at.txt

# Capture Go API response  
curl -s http://localhost:3000/api/v1/projects.json | jq '.[0].started_at' > go_started_at.txt

# Compare formats
echo "Rails: $(cat rails_started_at.txt)"
echo "Go:    $(cat go_started_at.txt)"
# Both should show: "2026-04-13T00:00:00+00:00"
```

**Edge Cases to Test:**
1. `0001-01-01T00:00:00Z` (zero time) → Should handle gracefully
2. `nil` timestamps → Should serialize to JSON null
3. Timezone-aware times → Should convert to UTC before formatting
4. Millisecond precision → Should preserve or truncate consistently

### 6. Risks and Considerations

**Risk 1: Breaking Client Compatibility**
- **Description**: Clients expecting `Z` suffix may break
- **Mitigation**: Accept both `Z` and `+00:00` during unmarshaling (backward compatible)
- **Impact**: Low - Both formats are valid ISO 8601, widely supported

**Risk 2: Rails API Format Mismatch**
- **Description**: Exact Rails format may differ from assumed `+00:00`
- **Mitigation**: Verify against actual Rails API response before implementation
- **Impact**: Medium - May require adjustment if Rails uses different format

**Risk 3: Timezone Conversion Errors**
- **Description**: Times may shift if timezone handling is incorrect
- **Mitigation**: Always convert to UTC before formatting; use `time.UTC` explicitly
- **Impact**: High - Data integrity issue; must test with known timestamps

**Risk 4: Inconsistent Field Handling**
- **Description**: Some fields may use `Z`, others `+00:00`
- **Mitigation**: Centralize formatting logic in helper functions
- **Impact**: Medium - Inconsistent user experience

**Risk 5: Test Coverage Gaps**
- **Description**: Existing tests may not validate format strictly
- **Mitigation**: Add explicit format validation tests; compare with Rails
- **Impact**: Medium - May miss format issues in production

**Design Considerations:**

1. **Centralized Formatter**: Create `internal/domain/dto/datetime_formatter.go` to avoid duplication
2. **Configuration**: Consider making format configurable (production vs development)
3. **Documentation**: Update API docs to specify exact datetime format
4. **Migration Path**: Consider deprecation warning if `Z` format was previous standard

**Deployment Checklist:**
- [ ] Verify Rails API datetime format matches assumption
- [ ] Run full test suite before deployment
- [ ] Deploy to staging environment first
- [ ] Compare staging vs production API responses
- [ ] Monitor client compatibility after deployment
- [ ] Have rollback plan if widespread issues occur

**Post-Implementation:**
- Update API documentation with exact datetime format
- Add datetime format validation to CI/CD pipeline
- Consider adding datetime format tests to performance benchmarks
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress: RDL-041

### Status: In Progress

### Completed Steps:

**1. Task Analysis**
- Reviewed RDL-041 task details
- Understanding datetime format requirements for Rails ISO 8601 compatibility
- Identifying files that need modification for datetime format standardization

**2. Files to Modify:**
- `internal/domain/dto/*.go` - Update time.Time field marshaling
- Implement custom MarshalJSON methods for datetime fields
- Update unmarshaling to accept both Z and +00:00 formats

**3. Current State:**
- Task status: To Do → In Progress
- Priority: HIGH
- Blocking: RDL-042, RDL-043, RDL-044 (Phase 2 tasks depend on this)

### Next Steps:
1. Update DTOs with custom datetime marshaling
2. Implement ISO 8601 format with timezone offset (+00:00)
3. Run tests using testing-expert subagent
4. Verify acceptance criteria
5. Document findings
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
