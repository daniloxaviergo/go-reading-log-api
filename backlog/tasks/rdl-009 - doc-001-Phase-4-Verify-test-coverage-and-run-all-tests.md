---
id: RDL-009
title: '[doc-001 Phase 4] Verify test coverage and run all tests'
status: To Do
assignee:
  - catarina
created_date: '2026-04-01 00:58'
updated_date: '2026-04-01 04:32'
labels: []
dependencies: []
references:
  - 'PRD Section: Acceptance Criteria'
  - 'Implementation Checklist: Testing'
documentation:
  - doc-001
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Run `go test ./...` to execute all tests against the test database.

Verify test coverage exceeds 80% on core packages using `go test -coverpkg=./... ./...`.

Fix any failing tests to ensure all acceptance criteria are met.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 All tests pass with go test ./...
- [ ] #2 Test coverage exceeds 80% on core packages
- [ ] #3 Tests run against test database successfully
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
## Implementation Plan: Verify Test Coverage and Run All Tests

### 1. Technical Approach

This task focuses on **verifying** existing test coverage rather than adding new tests. The approach will be:

1. **Test Execution**: Run `go test ./...` to verify all existing tests pass
2. **Coverage Analysis**: Use `go test -coverpkg=./... ./...` to measure coverage
3. **Database Integration**: Run tests against a test database to verify integration
4. **Gap Analysis**: Identify any failing tests or coverage gaps
5. **Documentation**: Report coverage metrics and any issues found

The implementation will focus on:
- Setting up proper test database environment
- Running coverage analysis with appropriate packages
- Verifying all tests pass in CI-ready conditions
- Reporting results in a standardized format

### 2. Files to Modify

**No new files creation required** - this is a verification task.

**Files to Verify/Analyze**:
- `/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go` - Test utilities
- `/home/danilo/scripts/github/go-reading-log-api-next/test/integration/` - Integration tests
- `/home/danilo/scripts/github/go-reading-log-api-next/test/unit/` - Unit tests
- `/home/danilo/scripts/github/go-reading-log-api-next/internal/config/config_test.go` - Config tests
- `/home/danilo/scripts/github/go-reading-log-api-next/internal/logger/logger_test.go` - Logger tests
- `/home/danilo/scripts/github/go-reading-log-api-next/internal/adapter/postgres/` - Repository implementations

**Test Files**:
- `test/integration/health_integration_test.go`
- `test/integration/projects_integration_test.go`
- `test/integration/logs_integration_test.go`
- `test/integration/test_context.go`
- `test/unit/project_repository_test.go`
- `test/unit/log_repository_test.go`
- `test/test_helper_test.go`

### 3. Dependencies

**Prerequisites**:
- Go 1.25.7+ installed and in PATH
- PostgreSQL server accessible (for integration tests)
- Test database configured via environment variables:
  - `DB_HOST` (default: localhost)
  - `DB_PORT` (default: 5432)
  - `DB_USER` (default: postgres)
  - `DB_PASS` (default: empty)
  - `DB_DATABASE` (default: reading_log)
  - `DB_DATABASE_TEST` (optional, falls back to DB_DATABASE + _test)

**Environment Setup**:
```bash
# Required for integration tests
export DB_HOST=${DB_HOST:-localhost}
export DB_PORT=${DB_PORT:-5432}
export DB_USER=${DB_USER:-postgres}
export DB_PASS=${DB_PASS:-}
export DB_DATABASE=${DB_DATABASE:-reading_log}
export DB_DATABASE_TEST=${DB_DATABASE_TEST:-${DB_DATABASE}_test}
```

### 4. Code Patterns

**Testing Patterns Used**:

1. **Test Helpers**: Shared database setup/teardown in `test/test_helper.go`
2. **Integration Context**: Test context with HTTP server in `test/integration/test_context.go`
3. **Mock Repositories**: In-memory mock implementations for unit tests
4. **Context with Timeout**: All DB operations use context with 5-second timeout
5. **Assertion Style**: Standard Go testing with if/err patterns

**Consistency Requirements**:
- All tests should use `t.Helper()` for helper functions
- All integration tests skip if no test database configured
- All tests use context with proper timeout
- Error messages follow Go conventions

### 5. Testing Strategy

**Test Execution Plan**:

1. **Unit Tests** (no database required):
   ```bash
   go test -v ./internal/config/...
   go test -v ./internal/logger/...
   go test -v ./internal/domain/...
   ```

2. **Integration Tests** (requires test database):
   ```bash
   go test -v ./test/...
   ```

3. **Coverage Report**:
   ```bash
   go test -coverpkg=./... -coverprofile=coverage.out ./...
   go tool cover -func=coverage.out
   ```

4. **Per-Package Coverage**:
   - Target: >80% on core packages
   - Measure: `internal/config`, `internal/logger`, `internal/domain`
   - Measure: `internal/adapter/postgres`, `internal/api/v1/handlers`

**Coverage Analysis**:
- Calculate coverage for each package
- Identify packages below 80% threshold
- Focus on core domain logic over test infrastructure

### 6. Risks and Considerations

**Potential Issues**:

1. **Test Database Dependencies**:
   - Integration tests will skip if DB not configured (documented behavior)
   - Need to verify `IsTestDatabase()` check functions correctly
   - Test schema setup must be idempotent

2. **Coverage Measurement**:
   - `-coverpkg=./...` may not correctly identify all packages
   - Need to verify coverage includes actual application code
   - Test infrastructure code should be excluded from coverage

3. **Test Isolation**:
   - Integration tests share test database
   - Clear test data between test suites
   - Verify rollback/cleanup works properly

4. **False Positives**:
   - Some tests may pass but not verify actual functionality
   - Mock tests don't verify database integration
   - Need to verify integration tests run against real DB

**Success Criteria**:
- All `go test ./...` commands complete with exit code 0
- Coverage report shows >80% on core packages
- Integration tests run successfully against test database
- No test failures or panics

**Reporting Requirements**:
1. Test execution output
2. Coverage metrics per package
3. List of any failing tests
4. List of packages below 80% coverage threshold
<!-- SECTION:PLAN:END -->
