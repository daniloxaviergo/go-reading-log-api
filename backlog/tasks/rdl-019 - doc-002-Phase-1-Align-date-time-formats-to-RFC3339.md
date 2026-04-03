---
id: RDL-019
title: '[doc-002 Phase 1] Align date time formats to RFC3339'
status: To Do
assignee:
  - catarina
created_date: '2026-04-03 14:02'
updated_date: '2026-04-03 14:19'
labels:
  - phase-1
  - date-format
  - code-quality
dependencies: []
references:
  - 'PRD Section: Technical Decisions - Decision 3: Date/Time Format Alignment'
  - 'PRD Section: Files to Modify - log_response.go'
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update all date/time fields in `log_response.go` to use RFC3339 format for timestamps and ISO date format for started_at field. Ensure NULL date fields serialize to JSON null instead of zero values.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Timestamp fields formatted as RFC3339 (e.g. 2024-01-15T10:30:00Z)
- [ ] #2 Date fields formatted as ISO date (e.g. 2024-01-15)
- [ ] #3 NULL database values serialize to JSON null
- [ ] #4 Format matches Rails API output exactly
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The task requires alignment of date/time formats in `log_response.go` to RFC3339 format for timestamps and ISO date format for `started_at` fields.

**Current State Analysis**:
- PostgreSQL stores `data` as `datetime` (timestamp) type in the `logs` table
- Current implementation reads `data` as a string directly from the database
- The `Data` field in `LogResponse` is `*string` type with no timestamp formatting
- The `started_at` field in `ProjectResponse` already uses `*time.Time` model type and gets formatted to RFC3339

**Required Changes**:
1. Modify the data flow to read timestamps from PostgreSQL as `time.Time`
2. Format timestamps to RFC3339 strings for JSON serialization
3. Ensure NULL values serialize to JSON null instead of zero values
4. Align with Rails API output which uses `ActiveModel::Serializer` default formatting

**Approach**:
- Update `Log` model to use `*time.Time` for `Data` field instead of `*string`
- Update PostgreSQL repository to scan `data` column as `time.Time`
- Add formatting helper to convert `time.Time` to RFC3339 string for JSON response
- Update handler to handle new model structure
- Ensure NULL timestamp serialization works correctly with pointer types

**Why This Approach**:
- RFC3339 is the modern standard for API date/time representation
- Consistent with PostgreSQL timestamp type
- Compatible with Rails ActiveModel::Serializer default behavior
- Enables proper time zone handling

### 2. Files to Modify

**Core Files**:

1. **`internal/domain/models/log.go`**
   - Change `Data` field type from `*string` to `*time.Time`
   - Update JSON tag to reflect new format

2. **`internal/adapter/postgres/log_repository.go`**
   - Update query scan to read `data` as `time.Time` instead of string
   - Update `GetByID`, `GetByProjectID`, `GetByProjectIDOrdered`, `GetAll` methods
   - Ensure proper NULL handling with pointer types

3. **`internal/api/v1/handlers/logs_handler.go`**
   - Update to handle `*time.Time` in Log model's Data field
   - Add helper function to format time.Time to RFC3339 string for response

4. **`internal/domain/dto/log_response.go`**
   - No changes needed if we format at the handler level
   - OR: Add custom JSON marshaling if formatting in DTO is preferred

**Test Files** (automatically updated by go test):

5. **`internal/domain/dto/log_response_test.go`**
   - Update test data to use time.Time instead of string
   - Verify RFC3339 format in JSON output

6. **`internal/api/v1/handlers/logs_handler_test.go`**
   - Update mock data to use time.Time for Data field
   - Update assertions to verify RFC3339 format

7. **`test/test_helper.go`**
   - Update mock repository to use time.Time for Log.Data

**No Changes Required**:
- `internal/repository/log_repository.go` (interface doesn't specify types for models)
- `internal/domain/models/project.go` (started_at already uses *time.Time correctly)

### 3. Dependencies

**Prerequisites**:
1. Understanding of RFC3339 format: `2024-01-15T10:30:00Z` (ISO 8601 with timezone)
2. PostgreSQL timestamp handling with Go's `time.Time` type
3. JSON null handling with pointer types in Go

**Expected Behavior**:
- Non-nil `time.Time` → RFC3339 formatted string in JSON
- Nil `time.Time` → JSON `null`
- ISO date format for `started_at` (date only): `2024-01-15`

**Rails API Reference**:
- Rails `t.datetime` columns → RFC3339 formatted strings
- Rails `t.date` columns → ISO date strings (`YYYY-MM-DD`)
- NULL values → JSON `null`

### 4. Code Patterns

**Conventions to Follow**:
1. **Pointer Types for Optional Fields**: Use `*time.Time` for nullable timestamps to ensure JSON null serialization
2. **RFC3339 Formatting**: Use `time.Time.Format(time.RFC3339)` for timestamp fields
3. **ISO Date Formatting**: Use `time.Time.Format("2006-01-02")` for date-only fields like `started_at`
4. **NULL Handling**: Go's `json.Marshal` with pointer types automatically converts nil to JSON null

**Existing Patterns** (from `project_response.go`):
- `Project.StartedAt` is `*time.Time` in model
- Formatted to RFC3339 in `ProjectResponse.StartedAt` (which is `*string`)
- This indicates a pattern: model uses `time.Time`, DTO uses formatted string

**Alternative Pattern Consideration**:
- Option A: Keep model as `*time.Time`, format at handler level (simpler, one format point)
- Option B: Model `*time.Time`, DTO `*string` with custom JSON marshaling (more reusable)

**Chosen Pattern**: Option A - Format at handler level for simplicity and single responsibility

### 5. Testing Strategy

**Unit Tests to Verify**:
1. **Log Model JSON Serialization**:
   - Non-nil timestamp formats as RFC3339 string
   - Nil timestamp serializes as JSON null
   - Example: `"2024-01-15T10:30:00Z"`

2. **Handler Integration**:
   - Logs handler returns RFC3339 formatted timestamps
   - Handles empty logs correctly
   - Handles NULL database values correctly

3. **Database Round-trip**:
   - Reading timestamp from PostgreSQL preserves time
   - Serialization to JSON maintains RFC3339 format

**Edge Cases to Test**:
- NULL data values in database → JSON null
- Timezone handling (should use UTC or preserve DB timezone)
- Midnight timestamps (e.g., `2024-01-15T00:00:00Z`)
- Different timestamp precisions

**Test Execution**:
- Use `go test ./...` to run all tests
- Verify `go test -v ./internal/domain/dto/...` for DTO tests
- Verify `go test -v ./internal/api/v1/handlers/...` for handler tests
- Ensure `go vet ./...` passes with no warnings
- Ensure `go fmt ./...` is run before commit

### 6. Risks and Considerations

**Potential Issues**:

1. **Breaking Change**: Changing `Data` from `*string` to `*time.Time` in the model is a breaking change for any code that directly uses the `models.Log` struct. However, since this is an internal domain model, the risk is contained.

2. **Database Compatibility**: The change assumes PostgreSQL `timestamp` type can be scanned to `time.Time` (this is standard pgx behavior and should work).

3. **Format Matching Rails**: Need to verify the exact format matches Rails API output. Rails uses ISO 8601/RFC3339 by default for datetime serializers, but may have timezone differences.

4. **Timezone Handling**: PostgreSQL timestamps may include timezone information. Need to ensure:
   - Times are converted to UTC for consistency
   - OR preserve the database timezone if that's the Rails behavior

5. **Zero Value Handling**: Empty timestamp (`0001-01-01T00:00:00Z`) could be misinterpreted. NULL should be used for absent values.

**Mitigation Strategies**:
- Review Rails API output to confirm exact format
- Test with actual database values including NULLs
- Add explicit tests for edge cases
- Verify with integration tests against test database

**Deployment Considerations**:
- No database migrations required (schema unchanged)
- No configuration changes required
- Risk is low as changes are contained to response formatting
<!-- SECTION:PLAN:END -->

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
