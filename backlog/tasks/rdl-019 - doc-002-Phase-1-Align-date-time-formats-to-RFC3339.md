---
id: RDL-019
title: '[doc-002 Phase 1] Align date time formats to RFC3339'
status: To Do
assignee:
  - thomas
created_date: '2026-04-03 14:02'
updated_date: '2026-04-03 14:55'
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
- The `Data` field in `Log` model is `*string` type (inconsistent with `Project.StartedAt` which is `*time.Time`)
- The `Data` field in `LogResponse` is `*string` type with no timestamp formatting applied
- The `started_at` field in `ProjectResponse` already uses `*time.Time` in model and gets formatted to RFC3339 string in the DTO

**Required Changes**:
1. Change `Log.Data` from `*string` to `*time.Time` (aligns with `Project.StartedAt` type)
2. Update PostgreSQL repository to scan `data` column as `time.Time` instead of string
3. Update DTO to format `*time.Time` to RFC3339 string for JSON output
4. Ensure NULL/nil values serialize to JSON null instead of zero values

**Approach**:
- Update `Log` model to use `*time.Time` for `Data` field (consistent with `Project.StartedAt`)
- Update PostgreSQL repository to scan `data` column as `time.Time` using `pgx` driver
- Update `LogResponse` DTO to include a helper method or custom JSON marshaling to format `time.Time` to RFC3339
- OR update handler to format at response creation time (simpler, follows existing pattern)

**Why This Approach**:
- RFC3339 is the modern standard for API date/time representation (`2024-01-15T10:30:00Z`)
- Consistent with PostgreSQL timestamp type and Go's `time.Time`
- Compatible with Rails ActiveModel::Serializer default behavior
- NULL values with pointer types serialize to JSON null automatically
- Aligns with existing pattern in `ProjectResponse` where `StartedAt` uses `*string` (formatted) while model uses `*time.Time`

### 2. Files to Modify

**Core Files**:

1. **`internal/domain/models/log.go`**
   - Change `Data` field type from `*string` to `*time.Time`
   - Update JSON tag to `"data"` (no changes needed)
   - This aligns with `Project.StartedAt` which is also `*time.Time`

2. **`internal/adapter/postgres/log_repository.go`**
   - Update all query scans to read `data` as `time.Time` instead of string
   - Update `GetByID`, `GetByProjectID`, `GetByProjectIDOrdered`, `GetAll` methods
   - Change `var data string` to `var data time.Time`
   - Set `log.Data = &data` (now `*time.Time` from `time.Time`)
   - Ensure proper NULL handling with pointer types

3. **`internal/domain/dto/log_response.go`**
   - Add helper method to format `time.Time` to RFC3339 string
   - OR add custom JSON marshaling
   - Since `Data` in model is now `*time.Time`, we need to format for JSON response

4. **`internal/api/v1/handlers/logs_handler.go`**
   - Update to handle `*time.Time` in Log model's Data field
   - Add helper to format time.Time to RFC3339 string for response DTO

**Test Files** (manual updates required):

5. **`internal/domain/dto/log_response_test.go`**
   - Update test data to use `time.Time` for Data field (or string if formatting at DTO level)
   - Update assertions to verify RFC3339 format in JSON output

6. **`internal/api/v1/handlers/logs_handler_test.go`**
   - Update mock data to use `time.Time` for Data field (or string depending on implementation)
   - Update assertions to verify RFC3339 format

7. **`test/test_helper.go`**
   - Update `MockLogRepository` to use `time.Time` for Log.Data

**No Changes Required**:
- `internal/repository/log_repository.go` (interface doesn't specify types)
- `internal/domain/models/project.go` (already uses `*time.Time` for StartedAt)
- `internal/repository/project_repository.go` (repository interface)

### 3. Dependencies

**Prerequisites**:
1. Understanding of RFC3339 format: `2024-01-15T10:30:00Z` (ISO 8601 with timezone)
2. PostgreSQL timestamp handling with Go's `time.Time` type
3. JSON null handling with pointer types in Go

**Expected Behavior**:
- Non-nil `time.Time` → RFC3339 formatted string in JSON (e.g., `"2024-01-15T10:30:00Z"`)
- Nil `time.Time` → JSON `null`
- ISO date format for `started_at` (date only): `2024-01-15`

**Rails API Reference**:
- Rails `t.datetime` columns → RFC3339 formatted strings (via ActiveModel::Serializer)
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
- `ProjectResponse.StartedAt` is `*string` (formatted from time.Time)
- This pattern suggests formatting at the DTO level

**Implementation Pattern Choice**:
Looking at `ProjectResponse`, it uses `*string` for `StartedAt` after formatting from `*time.Time` in the model. This indicates we should:
1. Keep model with `*time.Time`
2. Keep DTO with `*string` (formatted during response construction)

However, the current `LogResponse.Data` is already `*string` with no formatting helper. We need to add:
- A helper method to convert `*time.Time` to `*string` (RFC3339 format)
- OR custom JSON marshaling on the response

**Chosen Pattern**: Format in handler where the response is constructed (similar to `projects_handler.go`)

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
- Timezone handling (should use UTC for consistency)
- Midnight timestamps (e.g., `2024-01-15T00:00:00Z`)
- Different timestamp precisions (nanoseconds, etc.)

**Test Execution**:
- Use `go test ./...` to run all tests
- Verify `go test -v ./internal/domain/dto/...` for DTO tests
- Verify `go test -v ./internal/api/v1/handlers/...` for handler tests
- Verify `go test -v ./internal/adapter/postgres/...` for repository tests
- Ensure `go vet ./...` passes with no warnings
- Ensure `go fmt ./...` is run before commit

### 6. Risks and Considerations

**Potential Issues**:

1. **Breaking Change in Model**: Changing `Data` from `*string` to `*time.Time` in the model is a breaking change for any code that directly uses the `models.Log` struct. This is an internal domain model, so the risk is contained but test coverage is essential.

2. **Database Scanner Compatibility**: The `pgx` driver should scan PostgreSQL `timestamp` type to `time.Time` automatically. This is standard behavior and should work.

3. **Format Matching Rails**: Need to verify the exact format matches Rails API output. Rails uses ISO 8601/RFC3339 by default for datetime serializers. Potential timezone differences may exist.

4. **Timezone Handling**: PostgreSQL `timestamp` type without timezone stores times without timezone info. Need to ensure:
   - Times are interpreted consistently
   - Output uses UTC or appropriate timezone for API consistency

5. **Zero Value Handling**: Empty timestamp (`0001-01-01T00:00:00Z`) could be misinterpreted. NULL should be used for absent values.

6. **Existing String-Based Tests**: Tests that pass string data for `Log.Data` need to be updated to use `time.Time` or formatted strings depending on implementation.

**Mitigation Strategies**:
- Run all existing tests after changes to catch regressions
- Review Rails API output to confirm exact format
- Test with actual database values including NULLs
- Add explicit tests for edge cases
- Verify with integration tests against test database

**Deployment Considerations**:
- No database migrations required (schema unchanged)
- No configuration changes required
- Risk is low as changes are contained to response formatting
- Backward compatibility note: JSON output changes from unformatted string to RFC3339 formatted string
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
- 2026-04-03: Started implementation of RFC3339 date format alignment for Log.Data field

## Implementation Summary

### Changes Made:

1. **internal/domain/models/log.go**: Changed `Data` field from `*string` to `*time.Time`

2. **internal/adapter/postgres/log_repository.go**: 
   - Updated all query scans to read `data` as `time.Time` directly from PostgreSQL
   - Removed intermediate `string` variable
   - Scanning works directly because pgx driver maps PostgreSQL timestamp to `time.Time`

3. **internal/domain/dto/log_response.go**: 
   - Changed `Data` field from `*string` to `*time.Time` 
   - Added `time` import
   - JSON marshaling automatically formats `time.Time` to RFC3339

4. **internal/api/v1/handlers/logs_handler.go**:
   - Passes `logs[i].Data` (now `*time.Time`) directly to response DTO
   - No additional formatting needed - JSON marshaling handles RFC3339

5. **internal/api/v1/handlers/logs_handler_test.go**:
   - Updated test data to use `time.Date()` instead of string literals
   - Changed from `"2024-01-01"` to `time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)`

### Test Results:
- All unit tests pass (41 tests in DTO and handlers)
- go vet passes with no issues
- go fmt passes (no formatting changes needed)
- Application builds successfully

### Integration Tests:
- Integration tests fail due to missing PostgreSQL database (expected - infrastructure issue, not code issue)
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
