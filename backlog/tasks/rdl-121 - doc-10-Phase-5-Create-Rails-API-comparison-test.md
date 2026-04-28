---
id: RDL-121
title: '[doc-10 Phase 5] Create Rails API comparison test'
status: To Do
assignee:
  - thomas
created_date: '2026-04-28 00:30'
updated_date: '2026-04-28 05:25'
labels:
  - comparison-testing
  - phase-5
  - backend
dependencies: []
documentation:
  - doc-010
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Build comparison test that queries both Go and Rails APIs with same parameters and verifies responses match exactly for all stats fields.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Comparison test queries both APIs
- [ ] #2 All fields match between Go and Rails responses
- [ ] #3 Test documents any discrepancies
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task requires creating Go-based integration tests that compare responses between the Go API and the legacy Rails API for the projects/logs endpoints. The approach follows the existing pattern in `test/integration/rails_comparison_test.go` (which tests dashboard endpoints) but adapts it for the projects/logs endpoints.

**Key Technical Decisions:**
- **HTTP-based comparison**: Query both APIs via HTTP rather than direct database access to test the complete request/response cycle
- **Graceful skipping**: Tests skip when `RAILS_API_URL` environment variable is not set (Rails API not running)
- **Tolerance for floating-point**: Use delta comparisons (0.01 tolerance) for calculated fields like `progress`, `median_day`
- **1-day tolerance for `days_unreading`**: Account for timezone/time calculation differences between Go and Rails
- **JSON normalization**: Normalize JSON structure before comparison to handle field ordering differences
- **Fixture-based setup**: Use dashboard fixtures or create project-specific fixtures to ensure consistent test data

**Architecture:**
- Create new test file: `test/integration/projects_rails_comparison_test.go`
- Reuse existing `IntegrationTestContext` from `test_context.go`
- Reuse existing `SetupTestDB()` and fixture patterns
- Follow existing test naming conventions (`TestXxxComparison`)

### 2. Files to Modify

**New Files to Create:**
- `test/integration/projects_rails_comparison_test.go` - Main comparison test file for projects/logs endpoints

**Files to Read/Reference:**
- `test/integration/rails_comparison_test.go` - Reference for comparison test patterns
- `test/integration/test_context.go` - Test context and HTTP helpers
- `test/integration/projects_integration_test.go` - Projects endpoint test patterns
- `test/integration/logs_integration_test.go` - Logs endpoint test patterns
- `test/compare_responses.sh` - Reference for comparison logic and tolerance handling
- `test/fixtures/dashboard/fixtures.go` - Fixture patterns for test data setup
- `internal/domain/dto/project.go` - ProjectResponse DTO structure
- `internal/domain/dto/log.go` - LogResponse DTO structure

**Files to Modify:**
- `Makefile` - Add `compare-responses` target to run comparison tests
- `.env.test` - Add `RAILS_API_URL` environment variable placeholder

### 3. Dependencies

**Prerequisites:**
- PostgreSQL test database must be configured (`reading_log_test`)
- Rails API must be running on port 3001 (for actual comparison)
- Both Go and Rails databases must have identical fixture data

**Environment Variables:**
- `RAILS_API_URL` - Base URL for Rails API (e.g., `http://localhost:3001`)
- Standard test database variables from `.env.test`

**Existing Infrastructure to Leverage:**
- `test.TestHelper` - Database setup/cleanup
- `integration.IntegrationTestContext` - HTTP test helpers
- Dashboard fixture manager (or create project-specific fixtures)
- Existing JSON parsing helpers (`ParseProjectResponse`, `ParseLogResponse`)

### 4. Code Patterns

**Test Structure Pattern:**
```go
func TestProjectsIndexRailsComparison(t *testing.T) {
    if !IsTestDatabase() {
        t.Skip("Test database not configured")
    }
    
    railsURL := os.Getenv("RAILS_API_URL")
    if railsURL == "" {
        t.Skip("RAILS_API_URL not set - skipping Rails comparison test")
    }
    
    ctx := Setup(t)
    defer ctx.Teardown(t)
    
    // Setup fixture data
    // Query Go API
    // Query Rails API
    // Compare responses
}
```

**Response Comparison Pattern:**
```go
// Compare structures
assert.Equal(t, len(goProjects), len(railsProjects), "Project count mismatch")

// Compare calculated fields with tolerance
assert.InDelta(t, railsProgress, goProgress, 0.01, "Progress mismatch")

// Compare integer fields exactly
assert.Equal(t, railsLogsCount, goLogsCount, "Logs count mismatch")

// Compare dates (allow 1-day tolerance for days_unreading)
assert.InDelta(t, railsDaysUnreading, goDaysUnreading, 1, "Days unreading exceeds tolerance")
```

**Naming Conventions:**
- Test functions: `Test<Endpoint><ComparisonType>` (e.g., `TestProjectsIndexRailsComparison`)
- Test cases: Descriptive names in `t.Run()` subtests
- Variables: Follow Go conventions (camelCase for locals, PascalCase for exported)

**Integration with Existing Code:**
- Use `SetupRoutes()` from `test_context.go` to create test server
- Use `ctx.MakeRequest()` for Go API requests
- Use `http.Get()` for Rails API requests
- Use `ctx.ParseProjectResponse()` and `ctx.ParseLogResponse()` for parsing

### 5. Testing Strategy

**Test Coverage:**

1. **Projects Index Endpoint (`GET /v1/projects.json`)**
   - Test with multiple projects
   - Test with empty database
   - Verify all calculated fields: `progress`, `status`, `logs_count`, `days_unreading`, `median_day`, `finished_at`
   - Verify JSON:API envelope structure
   - Verify field names match Rails (snake_case)

2. **Projects Show Endpoint (`GET /v1/projects/{id}.json`)**
   - Test with existing project
   - Test with non-existent project (404)
   - Test with invalid ID (400)
   - Verify calculated fields match Rails
   - Verify logs are eager-loaded (first 4)

3. **Logs Index Endpoint (`GET /v1/projects/{id}/logs.json`)**
   - Test with multiple logs
   - Test with empty logs
   - Test with non-existent project (404)
   - Verify logs limited to 4
   - Verify ordering (by date DESC)
   - Verify project relationship structure

**Edge Cases to Cover:**
- Empty database (no projects)
- Project with no logs
- Project with more than 4 logs
- Null/missing `started_at`
- Completed projects (`page >= total_page`)
- Projects with zero progress
- Date/time format consistency (RFC3339)

**Tolerance Rules:**
- Floating-point fields (`progress`, `median_day`): ±0.01
- `days_unreading`: ±1 day (timezone calculation differences)
- All other fields: Exact match required

**Test Execution:**
```bash
# Run all comparison tests
RAILS_API_URL=http://localhost:3001 go test -v ./test/integration/... -run ".*Comparison.*"

# Run specific comparison test
RAILS_API_URL=http://localhost:3001 go test -v ./test/integration/... -run "TestProjectsIndexRailsComparison"

# Run without Rails API (should skip gracefully)
go test -v ./test/integration/... -run ".*Comparison.*"
```

**Verification Checklist:**
- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] `go fmt` and `go vet` pass
- [ ] Clean Architecture layers followed
- [ ] Error responses consistent with existing patterns
- [ ] HTTP status codes correct
- [ ] Tests skip gracefully when Rails API unavailable

### 6. Risks and Considerations

**Known Issues/Blocking:**
- **Rails API availability**: Tests require Rails API running on port 3001. Tests must skip gracefully when unavailable.
- **Database synchronization**: Test data must exist in both Go and Rails databases (they share the same PostgreSQL instance, so this is handled).
- **Timezone differences**: `days_unreading` calculation may differ by 1 day due to timezone handling. Document this tolerance.

**Potential Pitfalls:**
- **Floating-point precision**: Go and Rails may calculate `progress` or `median_day` with slightly different precision. Use tolerance comparisons.
- **Date format inconsistencies**: Ensure both APIs return RFC3339 format. Verify `started_at`, `data` (log date), and `finished_at` fields.
- **JSON field ordering**: JSON comparison should normalize structure, not rely on string equality.
- **Fixture data consistency**: Ensure fixture data produces predictable results for comparison.

**Deployment/Rollout:**
- Add `compare-responses` target to Makefile for easy execution
- Document in AGENTS.md how to run comparison tests
- Add `RAILS_API_URL` to `.env.test` as optional variable
- Consider adding to CI/CD pipeline (with Rails API container in docker-compose)

**Documentation Updates:**
- Update AGENTS.md with comparison test execution instructions
- Add section to README about Rails comparison testing
- Document tolerance rules and edge cases in test file comments

**Performance Considerations:**
- Comparison tests make HTTP requests to both APIs; ensure timeout is reasonable (5-10 seconds)
- Test database setup/teardown adds overhead; consider running comparison tests separately from unit tests
- Fixture loading should be efficient; reuse existing fixture patterns

**Future Enhancements:**
- Consider adding automated comparison test execution in CI/CD
- Add visual diff report generation for failed comparisons
- Extend to cover POST/PUT endpoints in Phase 2 (when log creation is implemented)
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
