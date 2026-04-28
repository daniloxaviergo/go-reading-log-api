---
id: RDL-112
title: '[doc-10 Phase 1] Modify Day handler to return flat JSON'
status: To Do
assignee:
  - thomas
created_date: '2026-04-28 00:28'
updated_date: '2026-04-28 01:37'
labels:
  - handler
  - phase-1
  - backend
dependencies: []
documentation:
  - doc-010
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update Day() method in dashboard_handler.go to return flat JSON structure with stats object at root level instead of JSON:API envelope. Remove data, type, id, attributes wrapper.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Response format is flat JSON with stats key
- [ ] #2 No JSON:API envelope present in response
- [ ] #3 Content-type remains application/json
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The implementation will modify the `Day()` handler in `dashboard_handler.go` to return a flat JSON structure instead of the current JSON:API envelope format. This involves:

**Response Format Change:**
- **Current:** `{ "data": { "type": "dashboard_day", "id": "...", "attributes": { "stats": {...} } } }`
- **Target:** `{ "stats": { ...all stats fields... } }`

**Key Changes:**
1. Remove JSON:API envelope wrapper (`data`, `type`, `id`, `attributes`)
2. Return `StatsData` directly as the root JSON object
3. Change Content-Type from `application/vnd.api+json` to `application/json`
4. Fix `per_pages` calculation to return `null` when `previous_week_pages` = 0 (currently returns hardcoded 133.333)
5. Add implementation for missing fields: `max_day`, `mean_geral`, `per_mean_day`, `per_spec_mean_day`

**Architecture Decisions:**
- Keep the `StatsData` DTO structure as-is (already has all required fields)
- Add new repository methods to fetch data for missing calculations
- Follow existing Clean Architecture patterns (Handler → Service → Repository)
- Maintain backward compatibility for existing fields

**Why this approach:**
- The PRD (doc-010) explicitly requires Rails parity for the response format
- Simpler JSON structure is easier for clients to consume
- Existing clients expect flat `stats` object, not JSON:API envelope
- The `StatsData` DTO already includes all new fields, minimizing DTO changes

### 2. Files to Modify

**Handler Layer:**
- `internal/api/v1/handlers/dashboard_handler.go`
  - Modify `Day()` method to return flat JSON structure
  - Remove JSON:API envelope wrapper
  - Change Content-Type header to `application/json`
  - Fix `per_pages` null handling logic
  - Add calls to new repository methods for `max_day`, `mean_geral`, `per_mean_day`, `per_spec_mean_day`

**Repository Interface:**
- `internal/repository/dashboard_repository.go`
  - Add `GetMaxByWeekday(ctx context.Context, date time.Time) (*float64, error)`
  - Add `GetOverallMean(ctx context.Context, date time.Time) (*float64, error)`
  - Add `GetPreviousPeriodMean(ctx context.Context, date time.Time) (*float64, error)`
  - Add `GetPreviousPeriodSpecMean(ctx context.Context, date time.Time) (*float64, error)`

**Repository Implementation:**
- `internal/adapter/postgres/dashboard_repository.go`
  - Implement `GetMaxByWeekday()` - Query max pages for target weekday
  - Implement `GetOverallMean()` - Calculate average of all weekday means
  - Implement `GetPreviousPeriodMean()` - Get mean from 7 days prior
  - Implement `GetPreviousPeriodSpecMean()` - Get speculative mean from 7 days prior

**DTO Layer:**
- `internal/domain/dto/dashboard_response.go`
  - No changes needed - `StatsData` already has all required fields (`MaxDay`, `MeanGeral`, `PerMeanDay`, `PerSpecMeanDay` as nullable pointers)
  - `Validate()` method already handles null values correctly

**Test Files:**
- `internal/api/v1/handlers/dashboard_handler_test.go`
  - Update `TestDashboardHandler_Day` to expect flat JSON structure
  - Update `TestDashboardHandler_Day_EmptyData` to expect flat JSON structure
  - Add new test `TestDashboardHandler_Day_NullPerPages` for null handling
  - Update mock expectations for new repository methods

### 3. Dependencies

**Prerequisites:**
- RDL-111 (Update StatsData DTO with new fields) - **Status: Done** (verified fields exist in `StatsData`)
- Database must support the new queries (no schema changes required)

**Blocked By:**
- None - This task can proceed independently

**Sequential Dependencies:**
- This task (RDL-112) must complete before:
  - RDL-113 (Update handler tests for new response structure) - Tests need to match new format
  - RDL-114 (Implement MaxDay field) - This task includes MaxDay implementation
  - RDL-115 (Implement PerMeanDay and PerSpecMeanDay fields) - This task includes these implementations

**Parallel Work:**
- Repository method implementations can be done in parallel with handler changes
- Test updates can be done after handler changes are complete

### 4. Code Patterns

**Following Existing Conventions:**

1. **Error Handling:**
   ```go
   if err != nil {
       fmt.Printf("Handler error: %v\n", err)
       http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
       return
   }
   ```

2. **Context Timeout:**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   ```

3. **JSON Response:**
   ```go
   w.Header().Set("Content-Type", "application/json")
   json.NewEncoder(w).Encode(response)
   ```

4. **Null Handling for Ratios:**
   ```go
   if previous_value == 0 {
       ratio = nil
   } else {
       calculated := current / previous
       ratio = &calculated
   }
   ```

5. **Float Rounding (3 decimals):**
   ```go
   math.Round(value*1000) / 1000
   ```

6. **Repository Pattern:**
   - Handler calls repository methods
   - Repository uses `pgxpool` with context timeout
   - Errors wrapped with context: `fmt.Errorf("failed to X: %w", err)`

**Naming Conventions:**
- Use `snake_case` for JSON field names (already correct in `StatsData`)
- Use `PascalCase` for Go struct fields
- Repository method names: `GetXxxYyy` pattern (e.g., `GetMaxByWeekday`)

**Integration Patterns:**
- Dependency injection via constructor (`NewDashboardHandler`)
- Mock repositories for unit tests (`MockDashboardRepository`)
- Context propagated through all layers

### 5. Testing Strategy

**Unit Tests (Mock-based):**

1. **TestResponseStructure:**
   - Verify response is flat JSON with `stats` key at root
   - Verify no `data`, `type`, `id`, `attributes` keys present
   - Verify Content-Type is `application/json`

2. **TestNullPerPages:**
   - When `previous_week_pages` = 0, verify `per_pages` is `null`
   - Verify JSON serialization produces `"per_pages": null`

3. **TestMaxDayCalculation:**
   - Mock `GetMaxByWeekday` to return specific value
   - Verify `max_day` field in response matches expected

4. **TestMeanGeralCalculation:**
   - Mock `GetOverallMean` to return specific value
   - Verify `mean_geral` field in response matches expected

5. **TestPerMeanDayCalculation:**
   - Mock `GetPreviousPeriodMean` to return specific value
   - Verify `per_mean_day` ratio calculation is correct
   - Verify `per_mean_day` is `null` when previous mean is 0

**Integration Tests (Database-based):**

1. **TestDayEndpointWithRealData:**
   - Use `TestHelper` from `test/test_helper.go`
   - Create fixtures with known log data
   - Verify response structure and calculated values match expectations

2. **TestDayEndpointNullHandling:**
   - Create fixtures with no previous period data
   - Verify `per_pages`, `per_mean_day`, `per_spec_mean_day` are all `null`

3. **TestDayEndpointWithDateParameter:**
   - Test with `?date=` query parameter
   - Verify calculations use the specified date's weekday

**Edge Cases to Cover:**

1. Empty database (no logs)
2. No logs for target weekday
3. No previous period data
4. Division by zero scenarios
5. Invalid date format (already tested)
6. Nil pointer handling for nullable fields

**Test Coverage Goals:**
- ≥80% line coverage for modified code (NFC-DASH-004)
- All new repository methods have unit tests
- All new calculation logic has unit tests with fixed test data

### 6. Risks and Considerations

**Known Issues:**

1. **Breaking Change for Clients:**
   - **Risk:** Existing clients expecting JSON:API envelope may break
   - **Mitigation:** PRD (doc-010) confirms this is intentional for Rails parity
   - **Decision:** Frontend team must verify compatibility (stakeholder alignment in PRD)

2. **Null Handling in JSON:**
   - **Risk:** Go's `*float64` with `nil` value serializes to `null` correctly
   - **Mitigation:** Already tested in existing code - `StatsData.PerPages` uses same pattern
   - **Verification:** Add explicit test for JSON null serialization

3. **Performance Impact:**
   - **Risk:** Additional repository queries may increase response time
   - **Mitigation:** 
     - Use efficient SQL queries with proper indexes
     - Consider combining queries where possible
     - Verify response time < 500ms (NFC-DASH-001)
   - **Database Indexes:** Verify `index_logs_on_project_id_and_data_desc` exists

4. **Calculation Accuracy:**
   - **Risk:** New calculations may not match Rails exactly
   - **Mitigation:** 
     - Study Rails V1::MeanLog implementation (see PRD doc-010)
     - Create comparison test against Rails API (RDL-121)
     - Use fixed test data with known expected output

**Trade-offs:**

1. **JSON:API Compliance vs Simplicity:**
   - **Trade-off:** Losing JSON:API spec compliance for simpler structure
   - **Decision:** Justified in PRD Decision 1 - backward compatibility is more important

2. **Extra Fields Preservation:**
   - **Trade-off:** Keeping Go-specific fields (`progress_geral`, `total_pages`, etc.)
   - **Decision:** PRD Decision 2 - maintain forward compatibility, can deprecate later

**Deployment Considerations:**

1. **Rollback Plan:**
   - Keep previous handler code in version control
   - Feature flag not needed - change is backward incompatible by design
   - Deploy during low-traffic window

2. **Monitoring:**
   - Monitor error rates after deployment
   - Track response time metrics
   - Alert on client-side parsing errors

3. **Documentation Updates:**
   - Update API documentation to reflect new response format
   - Document Go-specific extensions (extra fields)
   - Create Rails calculation reference document (RDL-123)

**Blocking Issues:**

- **None identified** - All prerequisites are met
- RDL-111 is Done (StatsData has all required fields)
- Database schema requires no changes
- No external dependencies required

**Acceptance Criteria Mapping:**

| AC | Implementation Step |
|----|---------------------|
| AC-DASH-001 (Response Structure) | Modify `Day()` to return flat JSON |
| AC-DASH-005 (Per Pages Null) | Fix `per_pages` logic to return `null` |
| AC-DASH-002 (Max Day) | Implement `GetMaxByWeekday()` repository method |
| AC-DASH-003 (Mean Geral) | Implement `GetOverallMean()` repository method |
| AC-DASH-004 (Per Mean Day) | Implement `GetPreviousPeriodMean()` repository method |
| NFC-DASH-002 (3 Decimal Precision) | Use existing `RoundToThreeDecimals()` method |
| NFC-DASH-004 (Test Coverage) | Write unit + integration tests with ≥80% coverage |
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
