---
id: RDL-158
title: '[doc-14 Phase 3] Add unit tests for dashboard handler SpeculateActual method'
status: To Do
assignee:
  - thomas
created_date: '2026-05-10 10:48'
updated_date: '2026-05-10 14:42'
labels:
  - testing
  - unit-tests
  - handler
  - phase-3
dependencies: []
documentation:
  - doc-014
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create unit tests for dashboard handler SpeculateActual method testing: successful response with valid data, error handling when service returns error, and response format validation. Tests must verify flat JSON structure and proper HTTP status codes.

Tests use mock SpeculateService to avoid database dependencies and follow existing handler test patterns.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 TestSpeculateActual_Success validates 200 OK and flat JSON structure
- [x] #2 TestSpeculateActual_ServiceError validates 500 status on service error
- [ ] #3 TestSpeculateActual_ResponseFormat verifies echart key exists at root level
- [ ] #4 Tests use mock SpeculateService with configurable return values
- [ ] #5 All handler tests compile and run without errors
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The implementation will create comprehensive unit tests for the `SpeculateActual` handler method in `dashboard_handler.go`. The tests will verify the handler's behavior using a mock `SpeculateService` to avoid database dependencies, following the existing handler test patterns established in `dashboard_handler_test.go`.

**Current State Analysis:**
- The `SpeculateActual` handler method (lines 467-483 in `dashboard_handler.go`) is already implemented in RDL-157
- It uses dependency injection with `SpeculateServiceInterface` 
- It returns flat JSON `{ echart: {...} }` without JSON:API envelope
- A basic test `TestDashboardHandler_SpeculateActual` already exists in `dashboard_handler_test.go` (lines 385-426)

**Implementation Strategy:**
1. **Expand Test Coverage**: The existing test only covers the success case. We need to add:
   - Service error handling test (500 status code)
   - Response format validation test (verify flat JSON structure)
   - Edge cases (nil chart config, empty series, etc.)

2. **Test Organization**: Create dedicated test functions that follow the naming convention:
   - `TestDashboardHandler_SpeculateActual_Success` - Valid response with 200 OK
   - `TestDashboardHandler_SpeculateActual_ServiceError` - Service error returns 500
   - `TestDashboardHandler_SpeculateActual_ResponseFormat` - Verify echart key at root level

3. **Mock Setup**: Use the existing `MockSpeculateService` that's already defined in the test file

**Why This Approach:**
- Follows existing test patterns in the codebase (see `TestDashboardHandler_Day`, `TestDashboardHandler_Faults`)
- Uses mock service to isolate handler logic from service/repository layers
- Ensures Clean Architecture separation (handler tests don't need database)
- Provides comprehensive coverage for success and error paths

### 2. Files to Modify

**Files to Read/Verify:**
- `internal/api/v1/handlers/dashboard_handler.go` - Review SpeculateActual method implementation
- `internal/api/v1/handlers/dashboard_handler_test.go` - Review existing test patterns and MockSpeculateService
- `internal/service/dashboard/speculate_service.go` - Review SpeculateServiceInterface methods
- `internal/domain/dto/dashboard_response.go` - Review EchartConfig structure

**Files to Modify:**
- `internal/api/v1/handlers/dashboard_handler_test.go`
  - **Add `TestDashboardHandler_SpeculateActual_Success`**: Split existing test to explicitly test success scenario with detailed assertions
  - **Add `TestDashboardHandler_SpeculateActual_ServiceError`**: Test service error returns 500 status with proper error message
  - **Add `TestDashboardHandler_SpeculateActual_ResponseFormat`**: Verify flat JSON structure with echart key at root level
  - **Add `TestDashboardHandler_SpeculateActual_EmptySeries`**: Edge case - service returns chart with empty series
  - **Add `TestDashboardHandler_SpeculateActual_NilChartConfig`**: Edge case - service returns nil chart config

**No Changes Required:**
- `internal/api/v1/handlers/dashboard_handler.go` - Handler implementation is complete
- `internal/service/dashboard/speculate_service.go` - Service interface is complete
- `internal/domain/dto/dashboard_response.go` - DTOs are complete
- `test/unit/api/v1/handlers/` - Tests will be added to existing file

### 3. Dependencies

**Prerequisites (Already Completed):**
- ✅ RDL-157: SpeculateActual handler implementation with flat JSON response
- ✅ RDL-155: SpeculateService with GenerateChartConfig method
- ✅ RDL-154: Mock repository implementations for unit testing
- ✅ MockSpeculateService already exists in `dashboard_handler_test.go`

**Test Dependencies:**
- `github.com/stretchr/testify/assert` - Assertion library
- `github.com/stretchr/testify/mock` - Mock framework
- `github.com/stretchr/testify/require` - Required assertions
- `net/http/httptest` - HTTP test utilities
- `context` - Context for service calls

**No Blocking Issues:**
- All handler code is complete
- All service interfaces are defined
- Mock service is already implemented
- Only test code needs to be added

### 4. Code Patterns

**Follow Existing Test Patterns:**

1. **Test Function Naming:**
```go
func TestDashboardHandler_<MethodName>_<Scenario>(t *testing.T) {
    // Test implementation
}
```

2. **Mock Setup Pattern:**
```go
mockRepo := &MockDashboardRepository{}
userConfig := service.NewUserConfigService(service.GetDefaultConfig())
mockProjectsService := &MockProjectsService{}
mockSpeculateService := &MockSpeculateService{}
handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService, mockSpeculateService)
```

3. **Service Mock Pattern:**
```go
// Success case
expectedChart := dto.NewEchartConfig().SetTitle("Test")
mockSpeculateService.On("GenerateChartConfig", mock.Anything).Return(expectedChart, nil)

// Error case
mockSpeculateService.On("GenerateChartConfig", mock.Anything).Return(nil, errors.New("service error"))
```

4. **HTTP Test Pattern:**
```go
req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/speculate_actual.json", nil)
w := httptest.NewRecorder()

handler.SpeculateActual(w, req)

assert.Equal(t, http.StatusOK, w.Code)
assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
```

5. **Response Validation Pattern:**
```go
var response map[string]interface{}
err := json.NewDecoder(w.Body).Decode(&response)
require.NoError(t, err)

// Verify no JSON:API envelope
_, hasData := response["data"]
assert.False(t, hasData, "Response should not have 'data' key")

// Verify echart key at root
echartMap, ok := response["echart"].(map[string]interface{})
require.True(t, ok, "Response should have 'echart' key")
```

6. **Mock Verification Pattern:**
```go
mockSpeculateService.AssertExpectations(t)
```

**Error Response Pattern:**
```go
// Expected error response format
assert.Equal(t, http.StatusInternalServerError, w.Code)
assert.Contains(t, w.Body.String(), "Internal server error")
```

**Naming Conventions:**
- Test function names: `TestDashboardHandler_<Method>_<Scenario>` (PascalCase)
- Mock variables: `mockSpeculateService` (camelCase)
- Expected values: `expectedChart` (camelCase)
- Response variables: `response`, `echartMap` (camelCase)

### 5. Testing Strategy

**Test Scenarios to Cover:**

1. **TestDashboardHandler_SpeculateActual_Success** (Acceptance Criteria #1)
   - **Setup**: Mock service returns valid EchartConfig with 2 series
   - **Expected**: 200 OK status, application/json content-type
   - **Assertions**:
     - Response has `echart` key at root level
     - No JSON:API envelope keys (`data`, `type`, `attributes`)
     - Chart has title "Speculated vs Actual"
     - Chart has 2 series: "Actual" and "Speculated"
     - Each series has 15 data points
   - **Mock**: `mockSpeculateService.On("GenerateChartConfig", mock.Anything).Return(expectedChart, nil)`

2. **TestDashboardHandler_SpeculateActual_ServiceError** (Acceptance Criteria #2)
   - **Setup**: Mock service returns error
   - **Expected**: 500 Internal Server Error status
   - **Assertions**:
     - Status code is 500
     - Response body contains error message
     - No chart data in response
   - **Mock**: `mockSpeculateService.On("GenerateChartConfig", mock.Anything).Return(nil, errors.New("service error"))`

3. **TestDashboardHandler_SpeculateActual_ResponseFormat** (Acceptance Criteria #3)
   - **Setup**: Mock service returns valid chart config
   - **Expected**: Flat JSON structure validation
   - **Assertions**:
     - Root level has `echart` key
     - No `data`, `type`, `id`, `attributes` keys (JSON:API envelope)
     - Content-Type is `application/json` (not `application/vnd.api+json`)
     - Echart object has required fields: title, legend, series, xAxis, yAxis
   - **Mock**: Same as success case

4. **TestDashboardHandler_SpeculateActual_EmptySeries** (Edge Case)
   - **Setup**: Mock service returns chart with empty series array
   - **Expected**: 200 OK (handler doesn't validate service output)
   - **Assertions**:
     - Status code is 200
     - Response structure is valid
     - Series array is empty
   - **Mock**: Return chart with `Series: make([]Series, 0)`

5. **TestDashboardHandler_SpeculateActual_NilChartConfig** (Edge Case)
   - **Setup**: Mock service returns nil chart config
   - **Expected**: 200 OK with null echart (handler doesn't validate)
   - **Assertions**:
     - Status code is 200
     - echart field is null or empty object
   - **Mock**: `mockSpeculateService.On("GenerateChartConfig", mock.Anything).Return((*dto.EchartConfig)(nil), nil)`

**Edge Cases to Cover:**
- Service returns nil error but nil chart config
- Service returns chart with missing optional fields (tooltip, grid, etc.)
- Service returns chart with only one series
- Context cancellation (optional, if time permits)

**Test Execution:**
```bash
# Run all dashboard handler tests
go test -v ./internal/api/v1/handlers/ -run TestDashboardHandler_SpeculateActual

# Run with coverage
go test -cover ./internal/api/v1/handlers/ -run TestDashboardHandler_SpeculateActual

# Run all tests to ensure no regressions
go test ./...
```

**Expected Test Output:**
```
=== RUN   TestDashboardHandler_SpeculateActual_Success
--- PASS: TestDashboardHandler_SpeculateActual_Success (0.00s)
=== RUN   TestDashboardHandler_SpeculateActual_ServiceError
--- PASS: TestDashboardHandler_SpeculateActual_ServiceError (0.00s)
=== RUN   TestDashboardHandler_SpeculateActual_ResponseFormat
--- PASS: TestDashboardHandler_SpeculateActual_ResponseFormat (0.00s)
=== RUN   TestDashboardHandler_SpeculateActual_EmptySeries
--- PASS: TestDashboardHandler_SpeculateActual_EmptySeries (0.00s)
=== RUN   TestDashboardHandler_SpeculateActual_NilChartConfig
--- PASS: TestDashboardHandler_SpeculateActual_NilChartConfig (0.00s)
PASS
```

**Validation Checklist:**
- ✅ All 5 test scenarios implemented
- ✅ Tests use mock SpeculateService (no database dependency)
- ✅ Tests verify flat JSON structure (no JSON:API envelope)
- ✅ Tests verify HTTP status codes (200, 500)
- ✅ Tests verify Content-Type header
- ✅ Tests verify response structure (echart key at root)
- ✅ Tests verify error handling
- ✅ All mocks properly verified with `AssertExpectations`
- ✅ Tests follow existing naming conventions
- ✅ Tests compile and run without errors

### 6. Risks and Considerations

**Technical Risks:**

1. **Test File Location**:
   - **Current**: Tests will be added to `internal/api/v1/handlers/dashboard_handler_test.go`
   - **Alternative**: Create separate file `speculate_actual_test.go`
   - **Decision**: Add to existing file to maintain consistency with other handler tests (Day, Projects, Faults all in same file)

2. **Mock Service Implementation**:
   - **Risk**: MockSpeculateService might not be fully implemented
   - **Mitigation**: Verify existing implementation in `dashboard_handler_test.go` (lines 17-37)
   - **Status**: ✅ Already implemented with all interface methods

3. **Response Structure Changes**:
   - **Risk**: EchartConfig structure might change
   - **Mitigation**: Tests validate structure dynamically using `map[string]interface{}`
   - **Impact**: Low - tests are flexible enough to handle optional fields

4. **Test Isolation**:
   - **Risk**: Tests might interfere with each other
   - **Mitigation**: Each test creates fresh mocks and handler instance
   - **Pattern**: Follow existing test pattern (see `TestDashboardHandler_Day`)

**Design Considerations:**

1. **Test Granularity**:
   - **Decision**: Separate tests for success, error, and format validation
   - **Rationale**: Clearer failure messages, easier maintenance
   - **Alternative**: Single comprehensive test (less maintainable)

2. **Edge Case Coverage**:
   - **Decision**: Include empty series and nil chart config tests
   - **Rationale**: Handler should handle gracefully even if service returns unexpected data
   - **Note**: Handler doesn't validate service output, so these test handler behavior, not service validation

3. **Assertion Library**:
   - **Decision**: Use `assert` for most assertions, `require` for critical checks
   - **Rationale**: `require` stops test on critical failure (e.g., JSON decode error)
   - **Pattern**: `require.NoError(t, err)` for decode, `assert.Equal` for values

**Acceptance Criteria Mapping:**

| AC | Implementation | Verification |
|----|----------------|--------------|
| #1 TestSpeculateActual_Success validates 200 OK and flat JSON | TestDashboardHandler_SpeculateActual_Success | Test passes with 200 status and echart key |
| #2 TestSpeculateActual_ServiceError validates 500 status | TestDashboardHandler_SpeculateActual_ServiceError | Test passes with 500 status |
| #3 TestSpeculateActual_ResponseFormat verifies echart key | TestDashboardHandler_SpeculateActual_ResponseFormat | Test verifies echart at root level |
| #4 Tests use mock SpeculateService | All tests use MockSpeculateService | Code review |
| #5 All handler tests compile and run | Run `go test ./internal/api/v1/handlers/` | Build and test verification |

**Definition of Done Checklist:**

- [ ] #1 All unit tests pass (`go test ./internal/api/v1/handlers/`)
- [ ] #2 All integration tests pass (no impact on integration tests)
- [ ] #3 `go fmt` and `go vet` pass with no errors
- [ ] #4 Clean Architecture layers properly followed (handler uses mock service)
- [ ] #5 Error responses consistent with existing patterns (500 status, error message)
- [ ] #6 HTTP status codes correct (200 for success, 500 for error)
- [ ] #7 Documentation updated (QWEN.md and AGENTS.md - separate task)
- [ ] #8 New code paths include error path tests (ServiceError test)
- [ ] #9 HTTP handlers test both success and error responses (Success + ServiceError tests)
- [ ] #10 Integration tests verify actual database interactions (separate task RDL-160)

**Related Tasks:**
- RDL-157: SpeculateActual handler implementation (DONE)
- RDL-159: Register route (NEXT)
- RDL-160: Integration tests for endpoint (NEXT)
- RDL-162: Implementation guide documentation (NEXT)
- RDL-163: API documentation update (NEXT)

**Implementation Notes:**

1. **Test File Structure**: Add tests after existing `TestDashboardHandler_SpeculateActual` test (lines 385-426)
2. **Code Organization**: Group new tests together for easy maintenance
3. **Comments**: Add brief comments explaining each test scenario
4. **Mock Setup**: Reuse MockSpeculateService pattern from existing tests
5. **Assertions**: Use descriptive error messages in assertions for easier debugging

**Example Test Implementation:**

```go
// TestDashboardHandler_SpeculateActual_Success validates 200 OK and flat JSON structure
func TestDashboardHandler_SpeculateActual_Success(t *testing.T) {
    mockRepo := &MockDashboardRepository{}
    userConfig := service.NewUserConfigService(service.GetDefaultConfig())
    mockProjectsService := &MockProjectsService{}
    mockSpeculateService := &MockSpeculateService{}
    handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService, mockSpeculateService)

    // Mock the service to return a valid chart config
    expectedChart := dto.NewEchartConfig().
        SetTitle("Speculated vs Actual").
        SetTooltip(map[string]interface{}{"trigger": "axis"})
    expectedChart.SetLegend(dto.NewLegend(true, []string{"Actual", "Speculated"}))

    // Add series with 15 data points each
    actualData := make([]interface{}, 15)
    specData := make([]interface{}, 15)
    for i := 0; i < 15; i++ {
        actualData[i] = float64(i * 10)
        specData[i] = float64(i * 11)
    }
    expectedChart.AddSeries(*dto.NewSeries("Actual", "line", actualData))
    expectedChart.AddSeries(*dto.NewSeries("Speculated", "line", specData))

    mockSpeculateService.On("GenerateChartConfig", mock.Anything).Return(expectedChart, nil)

    req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/speculate_actual.json", nil)
    w := httptest.NewRecorder()

    handler.SpeculateActual(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

    // Verify flat JSON structure (no JSON:API envelope)
    var response map[string]interface{}
    err := json.NewDecoder(w.Body).Decode(&response)
    require.NoError(t, err)

    // Verify no JSON:API envelope keys present
    _, hasData := response["data"]
    _, hasType := response["type"]
    _, hasAttributes := response["attributes"]
    assert.False(t, hasData, "Response should not have 'data' key")
    assert.False(t, hasType, "Response should not have 'type' key")
    assert.False(t, hasAttributes, "Response should not have 'attributes' key")

    // Verify echart key at root level
    echartMap, ok := response["echart"].(map[string]interface{})
    require.True(t, ok, "Response should have 'echart' key at root level")

    title, ok := echartMap["title"].(string)
    require.True(t, ok)
    assert.Equal(t, "Speculated vs Actual", title)

    seriesArr, ok := echartMap["series"].([]interface{})
    require.True(t, ok)
    assert.Len(t, seriesArr, 2)

    mockSpeculateService.AssertExpectations(t)
}

// TestDashboardHandler_SpeculateActual_ServiceError validates 500 status on service error
func TestDashboardHandler_SpeculateActual_ServiceError(t *testing.T) {
    mockRepo := &MockDashboardRepository{}
    userConfig := service.NewUserConfigService(service.GetDefaultConfig())
    mockProjectsService := &MockProjectsService{}
    mockSpeculateService := &MockSpeculateService{}
    handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService, mockSpeculateService)

    // Mock service to return error
    mockSpeculateService.On("GenerateChartConfig", mock.Anything).
        Return(nil, errors.New("service error"))

    req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/speculate_actual.json", nil)
    w := httptest.NewRecorder()

    handler.SpeculateActual(w, req)

    assert.Equal(t, http.StatusInternalServerError, w.Code)
    assert.Contains(t, w.Body.String(), "Internal server error")

    mockSpeculateService.AssertExpectations(t)
}

// TestDashboardHandler_SpeculateActual_ResponseFormat verifies echart key exists at root level
func TestDashboardHandler_SpeculateActual_ResponseFormat(t *testing.T) {
    mockRepo := &MockDashboardRepository{}
    userConfig := service.NewUserConfigService(service.GetDefaultConfig())
    mockProjectsService := &MockProjectsService{}
    mockSpeculateService := &MockSpeculateService{}
    handler := NewDashboardHandler(mockRepo, userConfig, mockProjectsService, mockSpeculateService)

    // Mock the service to return a valid chart config
    expectedChart := dto.NewEchartConfig().SetTitle("Speculated vs Actual")
    expectedChart.SetLegend(dto.NewLegend(true, []string{"Actual", "Speculated"}))
    expectedChart.AddSeries(*dto.NewSeries("Actual", "line", []interface{}{10, 20, 30}))
    expectedChart.AddSeries(*dto.NewSeries("Speculated", "line", []interface{}{11, 22, 33}))

    mockSpeculateService.On("GenerateChartConfig", mock.Anything).Return(expectedChart, nil)

    req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/echart/speculate_actual.json", nil)
    w := httptest.NewRecorder()

    handler.SpeculateActual(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    // Verify flat JSON structure
    var response map[string]interface{}
    err := json.NewDecoder(w.Body).Decode(&response)
    require.NoError(t, err)

    // Verify echart key exists at root level
    assert.Contains(t, response, "echart", "Response should have 'echart' key at root level")

    // Verify no JSON:API envelope
    assert.NotContains(t, response, "data", "Response should not have JSON:API 'data' key")
    assert.NotContains(t, response, "type", "Response should not have JSON:API 'type' key")
    assert.NotContains(t, response, "attributes", "Response should not have JSON:API 'attributes' key")

    // Verify Content-Type is application/json (not application/vnd.api+json)
    assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

    mockSpeculateService.AssertExpectations(t)
}
```

**Final Summary:**

This implementation plan outlines the creation of comprehensive unit tests for the `SpeculateActual` handler method. The tests will verify:
- Successful response with valid data (200 OK, flat JSON structure)
- Error handling when service returns error (500 status)
- Response format validation (echart key at root level, no JSON:API envelope)
- Edge cases (empty series, nil chart config)

All tests will use mock `SpeculateService` to avoid database dependencies and follow existing handler test patterns in the codebase.
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass
- [ ] #2 All integration tests pass execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md and AGENTS.md
- [ ] #8 New code paths include error path tests
- [ ] #9 HTTP handlers test both success and error responses
- [ ] #10 Integration tests verify actual database interactions
<!-- DOD:END -->
