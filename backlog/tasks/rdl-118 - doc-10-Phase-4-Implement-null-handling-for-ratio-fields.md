---
id: RDL-118
title: '[doc-10 Phase 4] Implement null handling for ratio fields'
status: To Do
assignee:
  - thomas
created_date: '2026-04-28 00:30'
updated_date: '2026-04-28 04:30'
labels:
  - null-handling
  - phase-4
  - backend
dependencies: []
documentation:
  - doc-010
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update per_pages logic to return null when previous_week_pages = 0. Apply same null handling to per_mean_day and per_spec_mean_day when denominator is 0 or nil.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 per_pages returns null when previous_week_pages = 0
- [x] #2 Ratio fields return null when denominator = 0
- [ ] #3 JSON serialization handles null correctly
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task implements null handling for ratio fields (`per_pages`, `per_mean_day`, `per_spec_mean_day`) to match Rails API behavior. The current implementation returns `0.0` when denominators are zero, but the Rails API returns `null` in these cases.

**Key Changes:**
- **per_pages**: Currently returns `0.0` when `previous_week_pages = 0`. Will be changed to return `null` (nil pointer).
- **per_mean_day**: Already has null handling in the handler layer, but needs verification that it returns `null` when `previous_mean = 0` or `nil`.
- **per_spec_mean_day**: Already has null handling in the handler layer, but needs verification that it returns `null` when `speculated_mean = 0` or `nil`.

**Architecture:**
- The DTO `StatsData` already uses pointer types (`*float64`) for nullable fields: `PerPages`, `PerMeanDay`, `PerSpecMeanDay`
- JSON serialization automatically converts nil pointers to `null` in Go
- The change is primarily in the business logic layer (services/handlers) to set these fields to `nil` instead of `0.0`

**Why this approach:**
- Uses Go's native pointer nil semantics for JSON null serialization
- Minimal code changes required (just set to nil instead of 0.0)
- Follows existing pattern already used for `per_mean_day` and `per_spec_mean_day` in the handler
- Aligns with Rails API behavior for API parity

### 2. Files to Modify

**Core Implementation:**
1. `internal/service/dashboard/day_service.go`
   - Modify `CalculatePerPagesRatio` method to return `nil` (not `0.0`) when `previousWeekPages = 0`
   - Update `CalculateWeeklyStats` to handle the nil return value correctly

2. `internal/api/v1/handlers/dashboard_handler.go`
   - Verify `per_mean_day` null handling: already sets `statsData.PerMeanDay = nil` when `prevMean <= 0`
   - Verify `per_spec_mean_day` null handling: already sets `statsData.PerSpecMeanDay = nil` when `prevSpecMean <= 0`
   - No changes likely needed, but verification required

**Test Files:**
3. `test/unit/day_service_test.go`
   - Update test "zero previous week pages" to expect `stats.PerPages = nil` instead of `*stats.PerPages = 0.0`
   - Update test "CalculatePerPagesRatio" with zero previous week to expect `0.0` from the method (but handler/service wraps it as nil)
   - Add new test cases specifically for null return behavior

4. `test/integration/dashboard_day_permean_integration_test.go`
   - Verify existing tests already assert `assert.Nil(t, statsMap["per_mean_day"])` for zero cases
   - Add test case for `per_pages` null handling with integration test

**Documentation:**
5. `docs/rails-calculation-reference.md`
   - Already documents null handling rules (verified in research)
   - No changes needed

6. `QWEN.md`
   - Update with implementation details and decisions

### 3. Dependencies

**Prerequisites:**
- None - this is a bug fix that aligns existing behavior with Rails API
- No database schema changes required
- No new dependencies required

**Related Tasks:**
- RDL-112: Modified Day handler to return flat JSON (completed)
- RDL-117: Created unit tests with fixed test data (completed)
- RDL-119: Update DTO validation for null values (To Do - may need coordination)

**Blocking Issues:**
- None known

### 4. Code Patterns

**Null Handling Pattern (Existing):**
```go
// In dashboard_handler.go - per_mean_day
if prevMeanErr == nil && prevMean != nil && *prevMean > 0 {
    ratio := math.Round(float64(statsData.MeanDay)/float64(*prevMean)*1000) / 1000
    statsData.PerMeanDay = &ratio
} else {
    statsData.PerMeanDay = nil  // Returns null in JSON
}
```

**New Pattern for per_pages:**
```go
// In day_service.go - CalculatePerPagesRatio
func (s *DayService) CalculatePerPagesRatio(lastWeekPages, previousWeekPages int) *float64 {
    if previousWeekPages == 0 {
        return nil  // Returns null in JSON
    }
    ratio := float64(lastWeekPages) / float64(previousWeekPages) * 100
    rounded := math.Round(ratio*1000) / 1000
    return &rounded
}
```

**JSON Serialization:**
- Go's `encoding/json` automatically converts `nil` pointers to `null` in JSON
- No custom marshaling required
- DTO fields already use `*float64` with `omitempty` tag

**Naming Conventions:**
- Follow existing naming: `CalculatePerPagesRatio`, `CalculateWeeklyStats`
- Return type changes from `float64` to `*float64` for nullable return

**Integration Patterns:**
- Service layer returns `*float64` (nullable)
- Handler layer sets DTO field directly from service return
- DTO field type is `*float64` with `omitempty` JSON tag

### 5. Testing Strategy

**Unit Tests (`test/unit/day_service_test.go`):**
1. **Test CalculatePerPagesRatio with zero previous week:**
   ```go
   t.Run("zero previous week returns nil", func(t *testing.T) {
       ratio := dayService.CalculatePerPagesRatio(100, 0)
       assert.Nil(t, ratio, "should return nil when previous week is 0")
   })
   ```

2. **Test CalculatePerPagesRatio with normal values:**
   ```go
   t.Run("normal ratio returns pointer", func(t *testing.T) {
       ratio := dayService.CalculatePerPagesRatio(150, 100)
       assert.NotNil(t, ratio)
       assert.InDelta(t, 150.0, *ratio, 0.001)
   })
   ```

3. **Test CalculateWeeklyStats null handling:**
   ```go
   t.Run("zero previous week pages results in null per_pages", func(t *testing.T) {
       // Setup mocks with prevWeekPages = 0
       stats, err := dayService.CalculateWeeklyStats(ctx)
       require.NoError(t, err)
       assert.Nil(t, stats.PerPages, "per_pages should be null when previous week is 0")
   })
   ```

**Integration Tests (`test/integration/dashboard_day_permean_integration_test.go`):**
1. **Test per_pages null handling with real database:**
   ```go
   t.Run("PerPagesNullHandling", func(t *testing.T) {
       // Create logs with no previous week data
       // Call handler
       // Assert statsMap["per_pages"] is nil
       assert.Nil(t, statsMap["per_pages"])
   })
   ```

2. **Verify existing per_mean_day and per_spec_mean_day null handling:**
   - Run existing tests to confirm they pass
   - Tests already assert `assert.Nil(t, statsMap["per_mean_day"])` for edge cases

**Edge Cases to Cover:**
- `previous_week_pages = 0` → `per_pages = null`
- `previous_week_pages > 0` and `last_week_pages = 0` → `per_pages = 0.0` (valid ratio)
- Empty database → all ratios = `null`
- Both weeks zero → `per_pages = null`

**Test Execution:**
```bash
# Run unit tests
go test -v ./test/unit/day_service_test.go

# Run integration tests
go test -v ./test/integration/dashboard_day_permean_integration_test.go

# Run all dashboard tests
go test -v ./test/unit/... -run Dashboard

# Check coverage
go test -cover ./test/unit/day_service_test.go
```

### 6. Risks and Considerations

**Known Issues:**
- **Current Bug**: `per_pages` returns `0.0` instead of `null` when `previous_week_pages = 0`
  - This breaks API parity with Rails
  - Frontend may expect `null` to handle "no data" vs "zero progress" differently

**Potential Pitfalls:**
1. **Breaking Change**: Frontend consumers may have handled `0.0` and not `null`
   - Mitigation: This is the intended behavior to match Rails API
   - Frontend should already handle `null` for other ratio fields

2. **Test Coverage**: Existing tests expect `0.0`, will fail after change
   - Mitigation: Update all affected tests before merging

3. **DTO Validation**: `StatsData.Validate()` may need updates for null values
   - Current validation allows nil `PerPages`: `if s.PerPages != nil { ... }`
   - No changes needed - validation already handles nil correctly

**Trade-offs:**
- **None**: This is a bug fix with no alternative approach
- Returning `null` is the correct behavior per Rails API specification

**Deployment Considerations:**
- No database migrations required
- No configuration changes required
- API response format changes (0.0 → null) - documented in release notes
- Backward compatible for clients that handle both `0.0` and `null`

**Rollout Plan:**
1. Update implementation
2. Update unit tests
3. Update integration tests
4. Run full test suite
5. Verify against Rails API responses
6. Document changes in QWEN.md
7. Deploy with release notes

**Verification Steps:**
1. Run `go test ./...` - all tests pass
2. Run `go fmt ./...` - code formatted
3. Run `go vet ./...` - no issues
4. Manual API test:
   ```bash
   # With no previous week data
   curl "http://localhost:3000/v1/dashboard/day.json?date=2024-01-15T10:00:00Z"
   # Expected: "per_pages": null
   ```
5. Compare with Rails API response for same data
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
