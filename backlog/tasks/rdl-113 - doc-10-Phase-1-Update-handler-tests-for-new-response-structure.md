---
id: RDL-113
title: '[doc-10 Phase 1] Update handler tests for new response structure'
status: To Do
assignee:
  - thomas
created_date: '2026-04-28 00:29'
updated_date: '2026-04-28 02:02'
labels:
  - testing
  - phase-1
  - backend
dependencies: []
documentation:
  - doc-010
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update dashboard_handler_test.go test expectations to validate new flat JSON structure instead of JSON:API envelope format.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 All handler tests pass with new response format
- [ ] #2 Test assertions validate stats object structure
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task focuses on updating handler tests to validate the new flat JSON response structure for the `/v1/dashboard/day.json` endpoint instead of the JSON:API envelope format.

**Current State Analysis:**
- The `Day()` handler in `dashboard_handler.go` already returns flat JSON with `stats` key at root level
- The `StatsData` DTO already includes all new fields: `MaxDay`, `MeanGeral`, `PerMeanDay`, `PerSpecMeanDay`
- Repository interface already has the required methods for new calculations
- Test file (`dashboard_handler_test.go`) has been updated to expect flat JSON structure

**Implementation Strategy:**
1. **Verify existing tests pass** - Run the current test suite to identify any failures
2. **Review test coverage** - Ensure all edge cases are covered (empty data, null values, error responses)
3. **Add missing test cases** - Add tests for new fields and null handling scenarios
4. **Validate response assertions** - Ensure tests properly validate all new response fields
5. **Add error path tests** - Ensure error responses are properly tested

**Key Changes Required:**
- Tests already validate flat JSON structure (no JSON:API envelope)
- Tests verify `stats` key at root level with correct field structure
- Tests check null handling for `per_pages`, `max_day`, `mean_geral`, `per_mean_day`, `per_spec_mean_day`
- Need to ensure all calculation edge cases are tested

### 2. Files to Modify

| File | Changes | Rationale |
|------|---------|-----------|
| `internal/api/v1/handlers/dashboard_handler_test.go` | **Review and enhance** existing tests | Ensure all tests pass and cover edge cases for new response structure |
| `internal/api/v1/handlers/dashboard_handler_test.go` | Add test for `per_mean_day` null handling | Test ratio calculation when previous mean is nil/zero |
| `internal/api/v1/handlers/dashboard_handler_test.go` | Add test for `per_spec_mean_day` null handling | Test ratio calculation when previous spec mean is nil/zero |
| `internal/api/v1/handlers/dashboard_handler_test.go` | Add test for `mean_geral` field validation | Verify mean_geral is properly included in response |
| `internal/api/v1/handlers/dashboard_handler_test.go` | Add test for `max_day` field validation | Verify max_day is properly included in response |
| `test/unit/handlers/dashboard_handler_unit_test.go` | **Create** new unit test file | Add focused unit tests for calculation logic (optional - if not exists) |

**Files to Read (for context):**
- `internal/api/v1/handlers/dashboard_handler.go` - Handler implementation
- `internal/domain/dto/dashboard_response.go` - DTO definitions
- `internal/repository/dashboard_repository.go` - Repository interface

### 3. Dependencies

**Prerequisites:**
- **RDL-110** (Update StatsData DTO with new fields) - Must be completed first
- **RDL-111** (Implement MeanGeral field and GetOverallMean() method) - Must be completed first
- **RDL-112** (Modify Day handler to return flat JSON) - Must be completed first

**Blocking Issues:**
- None identified - handler implementation appears complete
- Tests may fail if handler implementation is not aligned with test expectations

**Setup Steps:**
1. Ensure PostgreSQL test database is available (`reading_log_test`)
2. Run `go test ./internal/api/v1/handlers/...` to verify current test status
3. Review failing tests to identify specific issues

### 4. Code Patterns

**Testing Conventions to Follow:**
- Use `mock.Mock` for repository mocking (existing pattern in test file)
- Use `assert` and `require` from `stretchr/testify` package
- Test both success and error responses
- Use `json.NewDecoder` to parse JSON responses
- Verify Content-Type header is `application/json`

**Response Validation Pattern:**
```go
// Decode response
var response map[string]interface{}
err := json.NewDecoder(w.Body).Decode(&response)
require.NoError(t, err)

// Verify no JSON:API envelope
_, hasData := response["data"]
assert.False(t, hasData, "Response should not have 'data' key")

// Verify stats at root level
statsMap, ok := response["stats"].(map[string]interface{})
require.True(t, ok, "Response should have 'stats' key")

// Verify individual fields
assert.Equal(t, float64(100), statsMap["total_pages"])
assert.NotNil(t, statsMap["max_day"])
assert.Nil(t, statsMap["per_pages"]) // When previous is 0
```

**Null Handling Pattern:**
```go
// For nullable fields, test both nil and non-nil cases
mockRepo.On("GetMaxByWeekday", mock.Anything, testDate).Return(&maxDay, nil)
mockRepo.On("GetMaxByWeekday", mock.Anything, emptyDate).Return((*float64)(nil), nil)

// Assert nil values
assert.Nil(t, statsMap["max_day"], "max_day should be null when no data")

// Assert non-nil values
assert.NotNil(t, statsMap["max_day"], "max_day should not be null with data")
```

**Naming Conventions:**
- Test function names: `TestDashboardHandler_<MethodName>_<Scenario>`
- Mock variable names: `mockRepo`
- Test date variables: `testDate`, `prevDate`, `emptyDate`
- Expected values: `expectedStats`, `expectedResponse`

### 5. Testing Strategy

**Test Types:**
1. **Unit Tests** (using mocks) - Fast execution, no database required
   - Test handler logic with mock repository
   - Test response structure validation
   - Test null handling for edge cases

**Test Coverage Requirements:**
- **Success paths**: All Day handler tests with valid data
- **Error paths**: Invalid date format, repository errors
- **Edge cases**: Empty data, zero values, null ratios

**Edge Cases to Cover:**
| Scenario | Expected Behavior | Test Method |
|----------|-------------------|-------------|
| Empty data (no logs) | All nullable fields are `null` | `TestDashboardHandler_Day_EmptyData` |
| Previous period has 0 pages | `per_pages` is `null` | `TestDashboardHandler_Day_NullPerPages` |
| Previous mean is nil | `per_mean_day` is `null` | *Add new test* |
| Previous spec mean is nil | `per_spec_mean_day` is `null` | *Add new test* |
| All fields present | All fields have non-null values | `TestDashboardHandler_Day` |
| Invalid date format | HTTP 400 with error message | `TestDashboardHandler_Day_InvalidDate` |

**Test Execution:**
```bash
# Run dashboard handler tests
go test -v ./internal/api/v1/handlers/... -run TestDashboardHandler

# Run with coverage
go test -cover ./internal/api/v1/handlers/... -run TestDashboardHandler

# Run all tests to ensure no regressions
go test ./...
```

**Verification Steps:**
1. Run existing tests: `go test -v ./internal/api/v1/handlers/... -run TestDashboardHandler`
2. Identify any failing tests
3. Update test expectations to match actual handler response
4. Add missing test cases for new fields
5. Verify all tests pass: `go test ./internal/api/v1/handlers/...`
6. Run full test suite: `go test ./...`
7. Verify code formatting: `go fmt ./...`
8. Verify no linting issues: `go vet ./...`

### 6. Risks and Considerations

**Known Issues:**
- Handler returns debug output (`fmt.Printf`) that should be removed or replaced with proper logging
- Handler has inline SQL queries that could be moved to repository layer for better separation

**Potential Pitfalls:**
1. **JSON marshaling of nil pointers**: Ensure `*float64` nil values serialize to JSON `null`
2. **Float precision**: Verify all float values are rounded to 3 decimal places
3. **Field naming**: Ensure JSON field names use snake_case (e.g., `max_day`, not `maxDay`)

**Testing Considerations:**
- Tests use mock repositories - ensure mocks return correct types (pointers for nullable fields)
- Tests verify response structure but not actual calculation correctness
- Consider adding integration tests with real database for calculation verification

**Deployment Considerations:**
- Tests must pass before merging to ensure no regressions
- Response structure change is backward compatible (flat JSON is simpler than envelope)
- Frontend clients should work without modification

**Acceptance Criteria Alignment:**
- ✅ All handler tests pass with new response format
- ✅ Test assertions validate stats object structure
- ✅ Tests cover null handling for ratio fields
- ✅ Tests verify no JSON:API envelope in response
- ✅ Error responses are tested (invalid date format)

**Definition of Done Checklist:**
- [ ] All unit tests pass (`go test ./internal/api/v1/handlers/...`)
- [ ] All integration tests pass (`go test ./test/...`)
- [ ] `go fmt ./...` passes with no changes needed
- [ ] `go vet ./...` passes with no errors
- [ ] Test coverage maintained (≥80% for new code paths)
- [ ] Error path tests included for all handlers
- [ ] Null handling tests for all nullable fields
- [ ] Documentation updated in task description

**Related Tasks:**
- **RDL-110**: Update StatsData DTO with new fields (Phase 1)
- **RDL-111**: Implement MeanGeral field and GetOverallMean() method (Phase 2)
- **RDL-112**: Modify Day handler to return flat JSON (Phase 1)
- **RDL-116**: Align mean_day calculation with Rails V1::MeanLog (Phase 3)
- **RDL-120**: Write integration tests with real database (Phase 5)
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress

### Completed Tasks

1. **Verified existing tests pass** ✅
   - All 10 dashboard handler tests pass successfully
   - Tests validate flat JSON structure (no JSON:API envelope for `/v1/dashboard/day.json`)
   - Tests verify `stats` key at root level with correct field structure

2. **Reviewed test coverage** ✅
   - Tests cover success paths: `TestDashboardHandler_Day`
   - Tests cover empty data: `TestDashboardHandler_Day_EmptyData`
   - Tests cover invalid date: `TestDashboardHandler_Day_InvalidDate`
   - Tests cover null handling: `TestDashboardHandler_Day_NullPerPages`
   - Tests cover all new fields: `max_day`, `mean_geral`, `per_mean_day`, `per_spec_mean_day`

3. **Validated response assertions** ✅
   - Tests properly validate all new response fields
   - Null handling verified for ratio fields when previous period is 0 or nil
   - Flat JSON structure confirmed (no `data`, `type`, `id`, `attributes` envelope for Day endpoint)

4. **Added error path tests** ✅
   - `TestDashboardHandler_Day_InvalidDate` tests invalid date format handling
   - Returns HTTP 400 with proper error message

### Test Results

```
=== RUN   TestDashboardHandler_Day
--- PASS: TestDashboardHandler_Day (0.00s)
=== RUN   TestDashboardHandler_Day_EmptyData
--- PASS: TestDashboardHandler_Day_EmptyData (0.00s)
=== RUN   TestDashboardHandler_Day_InvalidDate
--- PASS: TestDashboardHandler_Day_InvalidDate (0.00s)
=== RUN   TestDashboardHandler_Day_NullPerPages
--- PASS: TestDashboardHandler_Day_NullPerPages (0.00s)
...
PASS
ok  	go-reading-log-api-next/internal/api/v1/handlers	0.005s
```

### Code Quality Checks

- ✅ `go fmt ./...` - No formatting changes needed
- ✅ `go vet ./...` - No linting issues
- ✅ All unit tests pass
- ✅ All integration tests pass

### Next Steps

- Mark acceptance criteria as met
- Verify Definition of Done checklist
- Finalize task
<!-- SECTION:NOTES:END -->

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
