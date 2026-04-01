---
id: RDL-008
title: '[doc-001 Phase 4] Create test infrastructure and integration tests'
status: To Do
assignee:
  - workflow
created_date: '2026-04-01 00:58'
updated_date: '2026-04-01 03:09'
labels: []
dependencies: []
references:
  - 'PRD Section: Acceptance Criteria'
  - 'Implementation Checklist: Testing'
  - 'PRD Section: Traceability Matrix'
documentation:
  - doc-001
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create test infrastructure in test/ directory with test_helper.go for common utilities and database setup.

Implement integration tests in test/project_integration_test.go and test/log_integration_test.go to verify endpoints work correctly against a test database.

Write unit tests for repository implementations using mocks.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Test database setup and cleanup implemented
- [ ] #2 Integration tests for all endpoints
- [ ] #3 Repository unit tests with mocks
- [ ] #4 Health check integration test
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task establishes a comprehensive test infrastructure for the Go project following Clean Architecture principles. The approach includes:

- **Test Directory Structure**: Create `test/` directory with organized subdirectories for different test types (unit, integration, mocks)
- **Test Helper Utilities**: Create `test/test_helper.go` with common setup functions including database connection pool, schema setup/teardown, and test context management
- **Integration Tests**: Write tests for all API endpoints using a test database to verify actual HTTP behavior
- **Repository Unit Tests**: Implement tests using mock repository implementations to verify data access logic in isolation
- **Mock Framework**: Use Go's interface-based mocking pattern (no external mocking framework needed) to enable repository unit tests

**Architectural Decisions**:
- Use Go's standard `testing` package (no external testing framework needed for Phase 1)
- Employ interface-based mocking for repository tests to maintain testability
- Use PostgreSQL test database with unique schema per test run to avoid race conditions
- Follow existing code patterns (context timeouts, error handling conventions)
- Integration tests will use a dedicated test database or schema to avoid polluting dev data

**Why This Approach**:
- Clean Architecture already uses interfaces for repositories, making mocking straightforward
- Go 1.21+ has robust standard testing package
- PostgreSQL test container or in-memory approach not needed for Phase 1 (simple setup/teardown suffices)
- Matches the project's existing patterns (context propagation, error wrapping, typed errors)

### 2. Files to Modify

**New Files to Create**:
- `test/test_helper.go` - Common test utilities (database setup/teardown, test context, helper functions)
- `test/test_helper_test.go` - Tests for test helper utilities
- `test/integration/test_context.go` - Integration test context management
- `test/unit/repository_mocks.go` - Mock repository implementations for unit tests
- `test/integration/health_integration_test.go` - Health check endpoint integration tests
- `test/integration/projects_integration_test.go` - Projects endpoints integration tests
- `test/integration/logs_integration_test.go` - Logs endpoints integration tests
- `test/unit/project_repository_test.go` - Project repository unit tests with mocks
- `test/unit/log_repository_test.go` - Log repository unit tests with mocks

**Existing Files Referenced (Read-Only Research)**:
- `internal/config/config.go` - Configuration loading patterns
- `internal/adapter/postgres/project_repository.go` - Repository implementation patterns
- `internal/adapter/postgres/log_repository.go` - Repository implementation patterns
- `internal/api/v1/handlers/projects_handler.go` - Handler patterns
- `internal/api/v1/handlers/logs_handler.go` - Handler patterns
- `internal/api/v1/handlers/health_handler.go` - Handler patterns
- `internal/domain/models/project.go` - Domain model structure
- `internal/domain/models/log.go` - Domain model structure
- `internal/domain/dto/*_response.go` - Response DTO structures

### 3. Dependencies

**Prerequisites**:
- RDL-004 (Configuration management) - Required for test database configuration
- RDL-007 (Application entry point) - Required for connection pool patterns
- Existing repository implementations - Required to understand interface signatures for mocking

**No External Dependencies Required**:
- Go's standard `testing` package sufficient for Phase 1
- No mocking framework needed (interface-based mocking)
- No test database container needed (can use in-memory SQLite or direct PostgreSQL setup)

**Setup Steps Before Implementation**:
1. Ensure `test/` directory exists
2. Verify database configuration is working (run application once)
3. Confirm repository interfaces are stable (already implemented)
4. Review handler code to understand expected error responses and HTTP status codes

**Potential Dependencies for Future**:
- `github.com/DATA-DOG/go-sqlmock` - For SQL query mocking (optional, not in Phase 1 scope)
- Test containers for PostgreSQL (optional, for complex integration scenarios)

### 4. Code Patterns

**Consistent Patterns from Existing Codebase**:

1. **Error Handling**:
   - Use `fmt.Errorf` with `%w` for wrapping errors
   - Return typed errors: `nil` for success, pointer to typed error for failures
   - In HTTP handlers: return JSON with `"error": "<message>"` format

2. **Context Propagation**:
   - Always pass context through all layers
   - Use `context.WithTimeout` for database operations (5-second default)
   - Derive from request context in handlers (`r.Context()`)

3. **Repository Interface Pattern**:
   - Define interface in `internal/repository/`
   - Implement concrete struct with `Impl` suffix
   - Use `New<Name>Impl(pool *pgxpool.Pool)` constructor pattern

4. **Test File Naming**:
   - Suffix with `_test.go`
   - Group by type: `integration/` and `unit/` subdirectories
   - Mirror package structure in test files

5. **HTTP Test Patterns**:
   - Use `httptest.NewRequest` and `httptest.NewRecorder`
   - Test `ResponseRecorder.Code` for status codes
   - Unmarshal JSON response for body content verification
   - Test both success and failure scenarios

6. **Database Test Setup**:
   - Use `pgxpool.New` with test database URL
   - Run schema setup before tests, teardown after
   - Use transactions for test isolation where possible
   - Clean up temporary data after each test

### 5. Testing Strategy

**Test Types and Coverage**:

1. **Integration Tests** (in `test/integration/`):
   - **Health Check Tests**: 
     - Verify `/healthz` returns `{"status":"ok","message":"healthy"}`
     - Test with database available/unavailable
   - **Projects Endpoint Tests**:
     - `GET /api/v1/projects` returns array of projects
     - `GET /api/v1/projects/{id}` returns single project with logs
     - Test 404 for non-existent project
   - **Logs Endpoint Tests**:
     - `GET /api/v1/projects/{project_id}/logs` returns first 4 logs
     - Test 404 when project not found
     - Test empty logs array when no logs exist

2. **Repository Unit Tests** (in `test/unit/`):
   - **Project Repository Tests**:
     - `GetByID` with existing and non-existing IDs
     - `GetAll` returns all projects
     - `GetWithLogs` eager-loads logs correctly
   - **Log Repository Tests**:
     - `GetByID` with existing and non-existing IDs
     - `GetByProjectID` returns all logs for project
     - `GetAll` returns all logs

3. **Test Helper Tests** (in `test/test_helper_test.go`):
   - Verify database connection helpers work
   - Test context timeout generation
   - Verify test data cleanup functions

**Edge Cases to Cover**:
- Empty database (no projects/logs)
- Non-existent resource IDs
- Malformed request parameters
- Database connection errors
- Context cancellation/timeout

**Verification Approach**:
- Run tests with `go test ./test/... -v`
- Verify all acceptance criteria from PRD
- Ensure tests pass with test database
- No coverage targets needed in Phase 1 (RDL-009 covers coverage verification)

### 6. Risks and Considerations

**Blocking Issues**:
1. **Test Database Configuration**: The test infrastructure requires a separate test database or schema. Need to decide between:
   - Dedicated `reading_log_test` database
   - Schema-based isolation in existing database
   - Transaction-based rollback per test
   **Recommendation**: Start with `reading_log_test` database for simplicity

2. **Database State Management**: Tests must not pollute dev/prod data. 
   **Solution**: Configure separate test database via environment variable (`DB_DATABASE_TEST`)
   **Risk**: If not configured, tests could run against wrong database
   **Mitigation**: Add runtime check in test helper to prevent execution without test database

**Potential Pitfalls**:
1. **Test Performance**: Integration tests with real database will be slower than unit tests
   **Mitigation**: Keep integration tests minimal, focus on API layer; unit tests cover business logic

2. **Flaky Tests**: Database timing issues could cause flaky tests
   **Mitigation**: Use consistent timeouts, ensure proper cleanup, run tests in isolated transactions

3. **Mock Complexity**: Repository mocks may need updates if interfaces change
   **Mitigation**: Keep mock implementations in separate file; update them when interfaces change

**Trade-offs**:
1. **No External Testing Framework**: Using standard library only keeps dependencies minimal but lacks some features
   **Rationale**: Phase 1 scope is MVP; can introduce frameworks later if needed

2. **No SQL Query Mocking**: Direct database queries instead of SQL mocks
   **Rationale**: Simpler setup for Phase 1; real database ensures integration coverage

3. **No Test Database Container**: Relies on existing PostgreSQL instance
   **Rationale**: Reduces complexity; test database can be any accessible PostgreSQL

**Deployment Considerations**:
- Test database should be configured via environment variables (not in source control)
- CI/CD pipeline should set up test database before running tests
- Test suite should clean up after itself (drop test tables/schemas)
- No data migration needed for tests (use existing schema)

**Implementation Checklist** (from PRD Acceptance Criteria):
- [ ] #1 Test database setup and cleanup implemented (test_helper.go)
- [ ] #2 Integration tests for all endpoints (integration test files)
- [ ] #3 Repository unit tests with mocks (unit test files)
- [ ] #4 Health check integration test (health_integration_test.go)
<!-- SECTION:PLAN:END -->
