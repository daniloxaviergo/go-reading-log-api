---
id: RDL-009
title: '[doc-001 Phase 4] Verify test coverage and run all tests'
status: To Do
assignee:
  - catarina
created_date: '2026-04-01 00:58'
updated_date: '2026-04-01 11:55'
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

Verify test coverage exceeds 100% on core packages using `go test -coverpkg=./... ./...`.

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
### 1. Technical Approach

This task focuses on **verifying** existing test coverage rather than adding new tests. The approach will be:

1. **Test Execution**: Run `go test ./...` to verify all existing tests pass
2. **Coverage Analysis**: Use `go test -coverpkg=./... -coverprofile=coverage.out ./...` to measure coverage
3. **Database Integration**: Run tests against a test database to verify integration tests work
4. **Gap Analysis**: Identify any failing tests or coverage gaps
5. **Reporting**: Generate coverage report and document results

**Technical Decisions**:
- Use `go test ./...` for comprehensive test execution
- Use `-coverpkg=./...` to ensure all packages are included in coverage
- Use `-coverprofile` to generate detailed coverage reports
- Set timeout for all database context operations (already implemented in test_helper.go)
- Skip integration tests if test database not configured (already in place)

**Why This Approach**:
- Go's standard testing utilities are sufficient for Phase 1
- No external coverage tools needed (go tool cover built-in)
- Test infrastructure already in place from RDL-008
- Clean separation between unit and integration tests

### 2. Files to Modify

**No new files creation required** - this is a verification task.

**Files to Verify/Analyze**:

**Test Infrastructure**:
- `/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper.go` - Test utilities and mock implementations
- `/home/danilo/scripts/github/go-reading-log-api-next/test/test_helper_test.go` - Test helper unit tests

**Unit Tests**:
- `/home/danilo/scripts/github/go-reading-log-api-next/test/unit/project_repository_test.go` - Project repository mock tests
- `/home/danilo/scripts/github/go-reading-log-api-next/test/unit/log_repository_test.go` - Log repository mock tests

**Integration Tests**:
- `/home/danilo/scripts/github/go-reading-log-api-next/test/integration/test_context.go` - Integration test context management
- `/home/danilo/scripts/github/go-reading-log-api-next/test/integration/health_integration_test.go` - Health check integration tests
- `/home/danilo/scripts/github/go-reading-log-api-next/test/integration/projects_integration_test.go` - Projects endpoint integration tests
- `/home/danilo/scripts/github/go-reading-log-api-next/test/integration/logs_integration_test.go` - Logs endpoint integration tests

**Application Code to Verify**:
- `/home/danilo/scripts/github/go-reading-log-api-next/internal/config/config.go` - Configuration (unit tested)
- `/home/danilo/scripts/github/go-reading-log-api-next/internal/logger/logger.go` - Logger (unit tested)
- `/home/danilo/scripts/github/go-reading-log-api-next/internal/domain/models/*.go` - Domain models
- `/home/danilo/scripts/github/go-reading-log-api-next/internal/adapter/postgres/*.go` - Repository implementations
- `/home/danilo/scripts/github/go-reading-log-api-next/internal/api/v1/handlers/*.go` - API handlers

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
  - `DB_DATABASE_TEST` (optional, falls back to DB_DATABASE with `_test` suffix)

**Environment Setup Before Testing**:

If PostgreSQL is available:
```bash
# Create test database
createdb reading_log_test

# Or set environment variable
export DB_DATABASE_TEST=reading_log_test
```

If PostgreSQL is not available:
- Unit tests will still run successfully
- Integration tests will be skipped with informative message
- This is expected behavior (not a failure)

### 4. Code Patterns

**Testing Patterns Used**:

1. **Test Helpers**: Shared database setup/teardown in `test/test_helper.go`
   - `SetupTestDB()` - Creates connection with test database
   - `SetupTestSchema()` - Creates tables for testing
   - `CleanupTestSchema()` - Drops tables after tests
   - `ClearTestData()` - Truncates data between tests
   - Context with 5-second timeout for all DB operations

2. **Integration Test Context**: Test context with HTTP server in `test/integration/test_context.go`
   - `Setup()` - Creates new test context with database and HTTP server
   - `Teardown()` - Cleans up database and closes server
   - `CreateTestProject()` - Creates test project in database
   - `MakeRequest()` - Makes HTTP request to test server

3. **Mock Repositories**: In-memory mock implementations for unit tests
   - `MockProjectRepository` - in-memory project store
   - `MockLogRepository` - in-memory log store
   - Call tracking for verification
   - Error injection for testing

4. **Assertion Style**: Standard Go testing with if/err patterns
   - Use `t.Helper()` for helper functions
   - Use `t.Fatalf()` for fatal errors
   - Use `t.Errorf()` for non-fatal errors
   - Verify both success and failure paths

**Consistency Requirements**:
- All tests use `t.Helper()` for helper functions (already implemented)
- Integration tests skip if no test database configured (already implemented)
- All tests use context with proper timeout (already implemented in test_helper.go)
- Error messages follow Go conventions (using `%w` for wrapping)

### 5. Testing Strategy

**Test Execution Plan**:

1. **Unit Tests** (no database required):
   ```bash
   go test -v ./test/unit/...
   go test -v ./test/test_helper_test.go
   go test -v ./internal/config/...
   go test -v ./internal/logger/...
   ```
   Expected: All unit tests pass (no DB connection needed)

2. **Integration Tests** (requires test database):
   ```bash
   go test -v ./test/integration/...
   ```
   Expected: Tests skip if DB not configured, otherwise pass

3. **Coverage Report**:
   ```bash
   go test -coverpkg=./... -coverprofile=coverage.out ./...
   go tool cover -func=coverage.out
   go tool cover -html=coverage.out  # Generate HTML report
   ```

4. **Per-Package Coverage Analysis**:
   - Target: >80% on core packages
   - Measure: `internal/config`, `internal/logger`, `internal/domain`
   - Measure: `internal/adapter/postgres`, `internal/api/v1/handlers`

**CoverageTargets**:
| Package | Target | Current (estimated) |
|---------|--------|---------------------|
| `internal/config` | >80% | ~100% (unit tested) |
| `internal/logger` | >80% | ~100% (unit tested) |
| `internal/domain` | >80% | ~80% (models tested via mocks) |
| `internal/adapter/postgres` | >80% | ~70% (integration tested) |
| `internal/api/v1/handlers` | >80% | ~75% (integration tested) |

**Test Categories**:
- **Unit Tests**: Mock-based, no database, fast execution
- **Integration Tests**: Real database, verify actual behavior
- **Helper Tests**: Test the test infrastructure itself

### 6. Risks and Considerations

**Potential Issues**:

1. **Test Database Dependencies**:
   - Integration tests will skip if DB not configured (documented behavior in `IsTestDatabase()`)
   - Test schema setup must be idempotent (uses `CREATE TABLE IF NOT EXISTS`)
   - Need to verify `IsTestDatabase()` check functions correctly (already implemented)

**Mitigation**: Tests automatically detect missing database and skip with informative message

2. **Coverage Measurement**:
   - `-coverpkg=./...` may not correctly identify all packages
   - Need to verify coverage includes actual application code
   - Test infrastructure code should be excluded from coverage

**Mitigation**: Use `-coverpkg` to explicitly specify packages, review coverage report

3. **Test Isolation**:
   - Integration tests share test database
   - Clear test data between test suites
   - Verify rollback/cleanup works properly

**Mitigation**: Each test uses its own data, cleanup happens in `Teardown()`

4. **False Positives**:
   - Some tests may pass but not verify actual functionality
   - Mock tests don't verify database integration
   - Need to verify integration tests run against real DB

**Mitigation**: Integration tests verify HTTP endpoints with real database connections

**Success Criteria**:
- [ ] All `go test ./...` commands complete with exit code 0
- [ ] Coverage report shows >80% on core packages
- [ ] Integration tests run successfully against test database (or skip if no DB)
- [ ] No test failures or panics
- [ ] Coverage metrics documented

**Reporting Requirements**:
1. Test execution output (pass/fail per package)
2. Coverage metrics per package
3. List of any failing tests
4. List of packages below 80% coverage threshold
5. Summary of environment setup (DB configured/not configured)

**Implementation Steps**:

```bash
# 1. Run all tests
go test ./... -v > test_output.txt 2>&1

# 2. Run coverage analysis
go test -coverpkg=./... -coverprofile=coverage.out ./...

# 3. Generate coverage report
go tool cover -func=coverage.out > coverage_report.txt

# 4. If HTML report needed
go tool cover -html=coverage.out -o coverage.html

# 5. Run unit tests specifically (no DB)
go test ./test/unit/... -v

# 6. Run integration tests (with DB if available)
go test ./test/integration/... -v

# 7. Document results
cat test_output.txt
cat coverage_report.txt
```
<!-- SECTION:PLAN:END -->
