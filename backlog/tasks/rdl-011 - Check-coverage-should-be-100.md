---
id: RDL-011
title: Check coverage should be 100%
status: To Do
assignee:
  - thomas
created_date: '2026-04-01 13:29'
updated_date: '2026-04-01 13:37'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add new tests to coverage to 100%
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The goal of this task is to achieve 100% test coverage for the Go application. Based on code review, the current test coverage is incomplete in several areas. The approach will be:

1. **Analyze current coverage gaps** - Review existing test coverage reports to identify untested code paths
2. **Add unit tests for uncovered areas** - Write focused unit tests for handlers, middleware, domain models, and DTOs that lack coverage
3. **Add integration tests for edge cases** - Add tests for error scenarios, boundary conditions, and edge cases not covered by existing integration tests
4. **Verify 100% coverage** - Run `go test -coverprofile` to verify all code paths are covered

The implementation will focus on:
- Handler methods that don't have test coverage (HealthHandler, ProjectsHandler edge cases, LogsHandler error paths)
- Domain models and DTOs without unit tests
- Middleware edge cases and error handling
- Repository error scenarios
- Helper functions and edge cases in test utilities

### 2. Files to Modify

**New test files to create:**
- `internal/api/v1/handlers/health_handler_test.go` - Unit tests for HealthHandler
- `internal/api/v1/handlers/projects_handler_test.go` - Unit tests for ProjectsHandler
- `internal/api/v1/handlers/logs_handler_test.go` - Unit tests for LogsHandler
- `internal/domain/dto/health_check_response_test.go` - Unit tests for DTOs
- `internal/domain/dto/project_response_test.go` - Unit tests for ProjectResponse DTO
- `internal/domain/dto/log_response_test.go` - Unit tests for LogResponse DTO
- `internal/domain/models/project_test.go` - Unit tests for Project model
- `internal/domain/models/log_test.go` - Unit tests for Log model
- `internal/config/integration_test.go` - Integration tests for config (if needed)

**Files to enhance with additional tests:**
- `internal/api/v1/middleware/recovery_test.go` - Add tests for panic recovery with stack traces
- `test/unit/project_repository_test.go` - Add comprehensive tests for MockProjectRepository
- `test/test_helper_test.go` - Add database lifecycle integration tests (already partially covered)

**Files to verify (may need coverage enhancements):**
- `internal/api/v1/middleware/middleware_test.go` - Verify Chain middleware tests are comprehensive
- `test/integration/health_integration_test.go` - Add more edge case tests

### 3. Dependencies

**Prerequisites:**
- No external dependencies required beyond existing project dependencies
- Tests should use existing mock implementations in `test/test_helper.go`
- Existing integration test infrastructure (`test/integration/test_context.go`) should be leveraged

**Blocking issues:**
- None identified. All required dependencies (slog, pgx, gorilla/mux) are already in go.mod

**Setup steps:**
1. Run `go test ./... -coverprofile=coverage.out` to see current coverage gaps
2. Review coverage report to identify specific uncovered lines
3. Add tests incrementally, re-running coverage to verify progress

### 4. Code Patterns

**Testing conventions to follow:**
- Use table-driven tests for multiple input scenarios
- Follow existing test naming conventions: `Test[Component]_[Scenario]` (e.g., `TestHealthHandler_EdgeCases`)
- Mock external dependencies (database) for unit tests
- Use `httptest.NewRecorder()` for HTTP handler testing
- Use context with timeout for context propagation tests
- Test both success and error paths for each handler

**Handler patterns:**
- Test all error conditions (invalid IDs, not found, database errors)
- Verify response status codes match expectations
- Verify JSON response format and content
- Test context propagation through handlers

**Mock repository patterns:**
- Set errors to test error handling paths
- Use `CallCount()` methods to verify method invocation
- Test both warm (data present) and cold (no data) scenarios

**Context patterns:**
- Use `context.WithTimeout()` for tests with database operations
- Verify context cancellation and timeout behavior
- Test context propagation through middleware chains

### 5. Testing Strategy

**Unit tests to add:**
1. **HealthHandler** (high priority):
   - Test basic health check response
   - Test response JSON format
   - Test different HTTP methods (GET/POST)

2. **ProjectsHandler** (high priority):
   - Test Index with empty projects list
   - Test Index with multiple projects
   - Test Show with valid project ID
   - Test Show with non-existent project (404)
   - Test Show with invalid project ID (400)
   - Test error handling from repository

3. **LogsHandler** (high priority):
   - Test Index with no logs
   - Test Index with logs (1, 4, and >4 logs to verify limit)
   - Test Index with non-existent project (404)
   - Test Index with invalid project ID (400)
   - Test project eager-loading in response

4. **DTOs and Models** (medium priority):
   - Test struct creation and field validation
   - Test context getter/setter methods
   - Test JSON serialization

5. **Middleware** (medium priority):
   - Add panic recovery with context cleanup
   - Test middleware chain error handling
   - Test context propagation through multiple middleware

**Integration tests to add:**
1. **Configuration tests:**
   - Test LoadConfig with missing .env file
   - Test LoadConfig with invalid port values
   - Test LoadConfig with empty environment variables

2. **Repository tests:**
   - Test error scenarios (connection failure, query errors)
   - Test concurrent access patterns

**Coverage verification:**
- Use `go test -coverprofile=coverage.out ./...` to collect coverage
- Use `go tool cover -func=coverage.out` to identify remaining gaps
- Run integration tests separately with database
- Verify both unit and integration test coverage

### 6. Risks and Considerations

**Potential risks:**
1. **Database connection issues** - Integration tests require a test database. If DB_DATABASE_TEST is not set, integration tests will be skipped.
   - **Mitigation**: Document that 100% coverage may require test database setup; unit tests should achieve coverage without database

2. **Third-party package methods** - Methods from external packages (slog, pgx, gorilla/mux) cannot be tested directly
   - **Mitigation**: Focus coverage on application code only; external package coverage is not required

3. **Time-sensitive code** - Code involving time operations may need special handling for deterministic tests
   - **Mitigation**: Use mock time sources or accept time variance in tests

4. **Race conditions** - Concurrent tests may introduce flakiness
   - **Mitigation**: Use synchronization primitives (sync.Mutex, channels) and ensure proper cleanup

**Considerations:**
- **Trade-off**: 100% coverage may include coverage for trivial getters/setters. Consider whether these add value.
- **Maintenance**: New code additions should maintain coverage; consider adding coverage checks to CI/CD
- **False positives**: High coverage doesn't guarantee quality; tests should verify correct behavior, not just execution

**Post-implementation tasks:**
1. Generate final coverage report with `go tool cover -html=coverage.out`
2. Document coverage results and any exceptions
3. Consider adding coverage threshold checks to build pipeline

### Implementation Steps

1. **Phase 1: Analysis** - Run coverage analysis to identify exact gaps
2. **Phase 2: Unit tests** - Add unit tests for handlers, DTOs, models
3. **Phase 3: Integration tests** - Add integration tests for edge cases
4. **Phase 4: Verification** - Run final coverage check and document results
5. **Phase 5: Cleanup** - Remove redundant tests, consolidate if needed
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Analysis Complete (2026-04-01 13:45)

### Current Coverage State
- **Total coverage**: ~60% overall (integration tests contribute significantly)
- **No coverage**: cmd/, internal/adapter/postgres/, handlers/, DTOs, Models
- **High coverage**: middleware (97.2%), config (86.7%), logger (100%)

### Key Gaps Identified
1. **Handlers** (0% coverage):
 - HealthHandler (Healthz method)
 - ProjectsHandler (Index, Show methods)
 - LogsHandler (Index method)

2. **DTOs** (0% coverage):
 - HealthCheckResponse, ProjectResponse, LogResponse

3. **Models** (0% coverage):
 - Project, Log

4. **Adapters** (0% coverage):
 - Project repository implementations
 - Log repository implementations

5. **Integration** (partial):
 - Test context helper methods

### Implementation Strategy
1. Add unit tests for handlers using httptest
2. Add unit tests for DTOs/Models (constructors, getters, setters)
3. Add comprehensive tests for repository implementations
4. Add integration tests for edge cases
5. Run coverage verification after each phase

### Next Steps
- Start with handler tests (highest priority)
- Use testing-expert subagent for all test execution

## Phase 1 Progress - Handlers (2026-04-01 14:00)

### HealthHandler - COMPLETE
- Created: `internal/api/v1/handlers/health_handler_test.go`
- Tests: 4/4 passing
- Coverage: 100% for health_handler.go

### Tests Created:
1. `TestHealthHandler_Healthz` - Basic health check with response verification
2. `TestHealthHandler_Healthz_GetMethod` - GET method test
3. `TestHealthHandler_Healthz_PostMethod` - POST method test  
4. `TestNewHealthHandler` - Constructor test

### Next Steps:
- ProjectsHandler (high priority)
- LogsHandler (high priority)

### Coverage Status:
- HealthHandler: 100%
- ProjectsHandler: 0%
- LogsHandler: 0%
<!-- SECTION:NOTES:END -->
