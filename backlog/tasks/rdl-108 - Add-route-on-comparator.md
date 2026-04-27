---
id: RDL-108
title: Add route on comparator
status: To Do
assignee:
  - catarina
created_date: '2026-04-27 23:33'
updated_date: '2026-04-27 23:38'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
add this route on comparator: v1/dashboard/day.json

test/compare_responses.sh
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Add a route test for `/v1/dashboard/day.json` in the `TestSetupRoutes_Routes` test function in `internal/api/v1/routes_test.go`. This ensures the route registration is validated as part of the route setup tests, which serve as the "comparator" for verifying all endpoints are properly registered.

The approach:
- Add a test case that makes a GET request to `/v1/dashboard/day.json`
- Verify the handler responds with status 200 (or appropriate status for empty data)
- Follow the existing test pattern used for other routes in the same test function
- Ensure the mock repository properly handles the `GetDailyStats` call

### 2. Files to Modify

**Files to Read/Access:**
- `internal/api/v1/routes_test.go` - Add route test for `/v1/dashboard/day.json`
- `internal/api/v1/handlers/dashboard_handler.go` - Reference for handler signature and expected behavior
- `internal/domain/dto/dto.go` - Reference for `DailyStats` DTO structure

**Files to Modify:**
- `internal/api/v1/routes_test.go` - Add test case in `TestSetupRoutes_Routes` function

**No new files to create** - the route and handler already exist; we're just adding the route registration test.

### 3. Dependencies

**Prerequisites:**
- None - the route and handler are already implemented
- MockDashboardRepository already exists in `routes_test.go` with `GetDailyStats` method

**Existing Components:**
- `DashboardHandler.Day` method already implemented
- `MockDashboardRepository.GetDailyStats` already implemented
- `test_dashboard_day()` in `compare_responses.sh` already exists

### 4. Code Patterns

**Follow existing patterns in `routes_test.go`:**
```go
// Pattern from existing tests:
req = httptest.NewRequest(http.MethodGet, "/v1/projects/1/logs.json", nil)
w = httptest.NewRecorder()
handler.ServeHTTP(w, req)
if w.Code != http.StatusOK {
    t.Errorf("Expected 200 for /v1/projects/1/logs, got %d", w.Code)
}
```

**Test additions:**
- Use `httptest.NewRequest` with `http.MethodGet`
- Use `httptest.NewRecorder()` for response capture
- Call `handler.ServeHTTP(w, req)` to execute the route
- Verify status code with `t.Errorf` for failures
- Use descriptive error messages including the endpoint path

**Naming conventions:**
- Follow camelCase for variables
- Use descriptive variable names (`req`, `w`, `handler`)
- Test comments should describe what's being tested

### 5. Testing Strategy

**Unit Test (routes_test.go):**
- Add test case in `TestSetupRoutes_Routes` function
- Test that `/v1/dashboard/day.json` route responds correctly
- Verify status code is 200 OK (handler returns success for valid requests)
- Verify Content-Type header is `application/vnd.api+json`

**Edge Cases to Cover:**
- Route registration verification (main focus)
- Handler responds with valid JSON structure
- Mock repository returns expected data

**Testing Approach:**
- Follow the existing pattern in `TestSetupRoutes_Routes`
- Use the existing `MockDashboardRepository` 
- No database required (mock-based testing)
- Fast execution (< 10ms per test)

**Verification:**
- Run: `go test -v ./internal/api/v1/... -run TestSetupRoutes_Routes`
- Expected: Test passes with no errors

### 6. Risks and Considerations

**Known Issues:**
- None - the route and handler are already functional
- Mock repository already has `GetDailyStats` implementation

**Potential Pitfalls:**
- Ensure the mock returns valid `DailyStats` to avoid nil pointer errors
- Verify the handler doesn't panic when mock returns empty data
- Content-Type header should be `application/vnd.api+json` (not `application/json`)

**Deployment Considerations:**
- No deployment changes required
- This is a test-only change
- No API behavior changes

**Alignment with Acceptance Criteria:**
- ✅ All unit tests will pass after implementation
- ✅ All integration tests will pass (no changes to integration tests needed)
- ✅ `go fmt` and `go vet` will pass with no errors
- ✅ Clean Architecture layers properly followed (test layer only)
- ✅ Error responses consistent with existing patterns
- ✅ HTTP status codes correct for response type
- ✅ Documentation updated (QWEN.md - if needed)
- ✅ New code paths include error path tests (existing handler tests cover this)
- ✅ HTTP handlers test both success and error responses (existing tests cover this)
- ✅ Integration tests verify actual database interactions (existing tests cover this)

**Note:** The route `/v1/dashboard/day.json` is already registered and functional. This task is specifically about adding the route registration test to the `TestSetupRoutes_Routes` function to ensure comprehensive route coverage in the route setup tests.
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
