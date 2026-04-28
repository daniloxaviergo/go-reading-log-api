---
id: RDL-111
title: '[doc-10 Phase 1] Update StatsData DTO with new fields'
status: Done
assignee:
  - next-task
created_date: '2026-04-28 00:28'
updated_date: '2026-04-28 01:32'
labels:
  - dto
  - phase-1
  - backend
dependencies: []
documentation:
  - doc-010
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add new fields to StatsData struct in dashboard_response.go: MaxDay, MeanGeral, PerMeanDay, PerSpecMeanDay. Update Validate() method to allow null values for ratio fields.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 StatsData struct includes MaxDay, MeanGeral, PerMeanDay, PerSpecMeanDay fields
- [x] #2 Validate() method accepts null values for ratio fields
- [x] #3 DTO compiles without errors
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task focuses on Phase 1 of the Dashboard Day Endpoint Alignment (doc-010). The implementation will:

**Technical Strategy:**
- Add 4 new fields to the `StatsData` DTO: `MaxDay`, `MeanGeral`, `PerMeanDay`, `PerSpecMeanDay`
- Update the `Validate()` method to allow `null` values for ratio fields (using pointer types)
- Follow existing Clean Architecture patterns and code conventions

**Key Decisions:**
1. **Use pointer types for nullable fields**: Following Go best practices for JSON null handling, ratio fields (`PerPages`, `PerMeanDay`, `PerSpecMeanDay`) will use `*float64` instead of `float64` to properly represent `null` in JSON
2. **Non-nullable fields remain as values**: `MaxDay` and `MeanGeral` will use `float64` with 0.0 as default for edge cases (no data)
3. **Validation updates**: The `Validate()` method will be updated to skip validation for nil pointer fields and only validate non-nil values

**Why this approach:**
- Matches Rails API behavior where ratio fields return `null` when denominator is 0
- Pointer types provide explicit null representation in JSON serialization
- Maintains backward compatibility with existing code while enabling new functionality
- Aligns with PRD Decision 4: "null Handling for Division by Zero"

### 2. Files to Modify

**Primary File:**
- `internal/domain/dto/dashboard_response.go`
  - Add `MaxDay *float64` field to `StatsData` struct (JSON: `max_day`)
  - Add `MeanGeral *float64` field to `StatsData` struct (JSON: `mean_geral`)
  - Add `PerMeanDay *float64` field to `StatsData` struct (JSON: `per_mean_day`)
  - Add `PerSpecMeanDay *float64` field to `StatsData` struct (JSON: `per_spec_mean_day`)
  - Change `PerPages` from `float64` to `*float64` to support null values
  - Update `NewStatsData()` to initialize pointer fields as `nil`
  - Add setter methods for new fields: `SetMaxDay()`, `SetMeanGeral()`, `SetPerMeanDay()`, `SetPerSpecMeanDay()`
  - Update `RoundToThreeDecimals()` to handle pointer fields (round only non-nil values)
  - Update `Validate()` method to:
    - Accept nil values for `PerPages`, `PerMeanDay`, `PerSpecMeanDay`
    - Validate `MaxDay` and `MeanGeral` when non-nil (must be >= 0)
    - Remove the 0-100 range constraint from `PerPages` (ratios can exceed 100%)

**Related Files (for context, not modified in this task):**
- `internal/repository/dashboard_repository.go` - Will need new interface methods in Phase 2
- `internal/adapter/postgres/dashboard_repository.go` - Will need new implementation in Phase 2
- `internal/api/v1/handlers/dashboard_handler.go` - Will need updates in Phase 2 to use new fields

### 3. Dependencies

**Prerequisites:**
- None - This is Phase 1 and can be implemented independently
- The DTO changes are foundational and do not depend on repository or handler changes

**Blocking Issues:**
- None identified

**Setup Steps:**
1. Ensure project compiles: `go build ./...`
2. Verify existing tests pass: `go test ./internal/domain/dto/...`

### 4. Code Patterns

**Naming Conventions:**
- Follow existing field naming: `MaxDay`, `MeanGeral`, `PerMeanDay`, `PerSpecMeanDay` (PascalCase for struct fields)
- JSON tags use snake_case: `max_day`, `mean_geral`, `per_mean_day`, `per_spec_mean_day`
- Setter methods follow pattern: `Set<FieldName>(value *float64) *StatsData`

**Integration Patterns:**
- **Pointer fields**: Use `*float64` for nullable fields to support JSON null serialization
- **Builder pattern**: Setter methods return `*StatsData` for method chaining (existing pattern)
- **Context embedding**: New fields do not require context (calculated values, not database entities)

**Code Style:**
- Follow existing struct field ordering: int fields first, then float fields, then pointer fields
- Use `math.Round(value*1000) / 1000` for 3-decimal rounding (existing pattern)
- Wrap errors with context: `fmt.Errorf("field validation failed: %w", err)`

### 5. Testing Strategy

**Unit Tests:** (`test/unit/dashboard_response_test.go`)

1. **Test StatsData field additions:**
   - Create StatsData with new fields populated
   - Verify JSON serialization includes new fields
   - Test method chaining with new setter methods

2. **Test null handling:**
   - Verify pointer fields serialize to `null` when nil
   - Verify pointer fields serialize to numeric value when set
   - Test JSON round-trip preserves null values

3. **Test validation:**
   - Test `Validate()` accepts nil ratio fields (no error)
   - Test `Validate()` rejects negative `MaxDay` when non-nil
   - Test `Validate()` rejects negative `MeanGeral` when non-nil
   - Test `Validate()` accepts PerPages > 100 (removed constraint)
   - Test `Validate()` with all fields nil (should pass)

4. **Test rounding:**
   - Verify `RoundToThreeDecimals()` handles nil pointer fields correctly
   - Verify non-nil pointer fields are rounded to 3 decimals

**Edge Cases to Cover:**
- All fields nil (empty StatsData)
- All fields set to zero values
- All fields set to maximum values
- Mix of nil and non-nil pointer fields
- Negative values for MaxDay and MeanGeral (should fail validation)
- PerPages values > 100 (should pass validation now)

**Test Organization:**
- Add test functions following existing naming: `TestStatsData_NewFields`, `TestStatsData_NullHandling`, `TestStatsData_ValidationUpdates`
- Use table-driven tests for validation scenarios
- Use `assert.Equal` and `require.NoError` from testify

### 6. Risks and Considerations

**Known Issues:**
- None - This is a straightforward DTO update with no complex logic

**Potential Pitfalls:**
1. **JSON serialization**: Ensure pointer fields serialize correctly to `null` vs `0.0`
   - Mitigation: Write explicit tests for JSON marshaling/unmarshaling
   
2. **Backward compatibility**: Existing code may assume `PerPages` is non-nil
   - Mitigation: Update calling code in Phase 2 to handle nil values gracefully
   
3. **Validation logic**: Removing the 0-100 constraint from PerPages may affect other parts of the system
   - Mitigation: Document this change and verify in Phase 2 that ratio > 100% is expected behavior

**Deployment Considerations:**
- No database changes required
- No API contract changes in this phase (fields are added but not yet populated)
- Safe to deploy independently - new fields will be `null` until Phase 2 implementation

**Rollback Plan:**
- Simple code revert if issues arise
- No data migration needed

**Follow-up Tasks:**
- Phase 2 (RDL-114, RDL-115): Implement repository methods to calculate and populate new fields
- Phase 3 (RDL-116): Align calculation algorithms with Rails
- Phase 4 (RDL-118, RDL-119): Update validation and null handling across the system
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary

Successfully implemented Phase 1 of the Dashboard Day Endpoint Alignment by updating the StatsData DTO with new fields and null handling for ratio fields.

## What Was Done

### 1. StatsData Struct Updates (`internal/domain/dto/dashboard_response.go`)
- Added 4 new fields to StatsData:
  - `MaxDay *float64` (JSON: `max_day`) - Maximum pages in a single day
  - `MeanGeral *float64` (JSON: `mean_geral`) - General mean across all days
  - `PerMeanDay *float64` (JSON: `per_mean_day`) - Ratio for mean day
  - `PerSpecMeanDay *float64` (JSON: `per_spec_mean_day`) - Ratio for speculative mean day
- Changed `PerPages` from `float64` to `*float64` to support null values
- Updated `NewStatsData()` to initialize pointer fields as `nil`

### 2. New Setter Methods
- `SetMaxDay(value *float64) *StatsData`
- `SetMeanGeral(value *float64) *StatsData`
- `SetPerMeanDay(value *float64) *StatsData`
- `SetPerSpecMeanDay(value *float64) *StatsData`

### 3. Updated Methods
- `RoundToThreeDecimals()` - Now handles pointer fields (rounds only non-nil values)
- `Validate()` - Updated to:
  - Accept nil values for `PerPages`, `PerMeanDay`, `PerSpecMeanDay`
  - Validate `MaxDay` and `MeanGeral` when non-nil (must be >= 0)
  - Remove the 0-100 range constraint from `PerPages` (ratios can exceed 100%)

### 4. Related Code Updates
- `internal/service/dashboard/day_service.go` - Updated to convert float64 to *float64 when calling SetPerPages
- `internal/api/v1/handlers/dashboard_handler.go` - Updated to use pointer types for PerPages

### 5. Test Updates
- `test/unit/dashboard_response_test.go` - Added comprehensive tests for:
  - New fields (TestStatsData_NewFields)
  - JSON serialization of new fields
  - Null handling for pointer fields
  - Method chaining with new setters
  - Rounding of pointer fields
  - Validation of new fields
  - PerPages null handling
- `test/unit/day_service_test.go` - Updated to handle pointer type for PerPages
- `test/fixtures/dashboard/scenarios.go` - Updated StatsExpectations to use *float64 for PerPages
- `test/dashboard_integration_test.go` - Updated to handle pointer type for PerPages
- `test/integration/error_scenarios_test.go` - Updated to handle pointer type for PerPages

## Key Changes

### Files Modified
1. `internal/domain/dto/dashboard_response.go` - Primary DTO updates
2. `internal/service/dashboard/day_service.go` - Service layer adaptation
3. `internal/api/v1/handlers/dashboard_handler.go` - Handler layer adaptation
4. `test/unit/dashboard_response_test.go` - New unit tests
5. `test/unit/day_service_test.go` - Updated tests
6. `test/fixtures/dashboard/scenarios.go` - Updated fixtures
7. `test/dashboard_integration_test.go` - Updated integration tests
8. `test/integration/error_scenarios_test.go` - Updated integration tests
9. `test/unit/project_calculations_test.go` - Removed duplicate floatPtr helper

## Testing

All tests pass:
- Unit tests: `go test ./test/unit/...` - PASS
- Integration tests: `go test ./test/integration/...` - PASS
- Build: `go build ./...` - SUCCESS
- Format: `go fmt ./...` - SUCCESS
- Vet: `go vet ./...` - SUCCESS

## Notes for Reviewers

- This is a foundational Phase 1 change; new fields will be `null` until Phase 2 implementation
- The change is backward compatible - existing code will continue to work
- The removal of the 0-100 constraint on PerPages allows ratios > 100%, which matches Rails API behavior
- Pointer types provide explicit null representation in JSON serialization

## Follow-up Tasks

- Phase 2 (RDL-114, RDL-115): Implement repository methods to calculate and populate new fields
- Phase 3 (RDL-116): Align calculation algorithms with Rails
- Phase 4 (RDL-118, RDL-119): Update validation and null handling across the system
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass
- [x] #2 All integration tests pass execution and verification
- [x] #3 go fmt and go vet pass with no errors
- [x] #4 Clean Architecture layers properly followed
- [x] #5 Error responses consistent with existing patterns
- [x] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md
- [x] #8 New code paths include error path tests
- [x] #9 HTTP handlers test both success and error responses
- [x] #10 Integration tests verify actual database interactions
<!-- DOD:END -->
