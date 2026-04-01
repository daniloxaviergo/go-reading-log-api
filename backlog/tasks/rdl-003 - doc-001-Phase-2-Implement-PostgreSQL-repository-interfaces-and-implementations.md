---
id: RDL-003
title: >-
  [doc-001 Phase 2] Implement PostgreSQL repository interfaces and
  implementations
status: To Do
assignee:
  - catarina
created_date: '2026-04-01 00:57'
updated_date: '2026-04-01 01:41'
labels: []
dependencies: []
references:
  - 'PRD Section: Technical Decisions'
  - 'Implementation Checklist: Core Components'
  - 'PRD Section: Repository Pattern'
documentation:
  - doc-001
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Define repository interfaces in internal/repository/project_repository.go and internal/repository/log_repository.go using the repository pattern for data access abstraction.

Implement concrete PostgreSQL adapters in internal/adapter/postgres/ that use pgx/v5 for database operations with proper connection pooling configuration.

Ensure all methods accept context for timeout and cancellation propagation with 5-second timeout.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Repository interfaces defined with clean abstraction for data access
- [ ] #2 PostgreSQL implementations use pgx/v5 with connection pooling
- [ ] #3 All methods accept context with proper timeout handling
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Implement repository pattern for data access abstraction with PostgreSQL adapters:

**Repository Interfaces** (`internal/repository/`):
- Define `ProjectRepository` interface with methods: `GetByID(ctx, id)`, `GetAll(ctx)`, `GetWithLogs(ctx, id)`
- Define `LogRepository` interface with methods: `GetByID(ctx, id)`, `GetByProjectID(ctx, projectID)`, `GetAll(ctx)`
- Each method returns domain models or DTOs and an error
- Use context for timeout/cancellation propagation

**PostgreSQL Implementations** (`internal/adapter/postgres/`):
- Implement adapter structs embedding `*pgx.Conn` for database connection
- Use `pgx/v5/stdlib` bridges for database/sql compatibility
- Configure connection pooling per PRD spec: `MaxOpenConns=25`, `MaxIdleConns=25`, `ConnMaxLifetime=5m`, `ConnMaxIdleTime=1m`
- Implement each method using prepared statements with 5-second context timeout
- Handle nulls with `.sql.NullableXxx` types for nullable DB fields

**Architecture Decisions:**
- Use `pgx/v5` directly for query execution (more control than database/sql wrapper)
- Connection pooling configured in adapter initialization
- Context timeout set per method call (5 seconds per acceptance criteria)
- Return domain models from repository, convert to DTOs in handler layer

**Why this approach:**
- Clean separation between business logic and data access
- PostgreSQL-specific optimizations via pgx/v5
- Future portability (can swap implementations for testing)
- Explicit context handling for timeout propagation

### 2. Files to Modify

**New files to create:**
- `internal/repository/project_repository.go` - Project repository interface
- `internal/repository/log_repository.go` - Log repository interface
- `internal/adapter/postgres/project_repository.go` - Project repository implementation
- `internal/adapter/postgres/log_repository.go` - Log repository implementation

**Files referenced by implementations:**
- `internal/domain/models/project.go` - Input/output types
- `internal/domain/models/log.go` - Input/output types
- `internal/domain/dto/project_response.go` - Response DTOs
- `internal/domain/dto/log_response.go` - Response DTOs
- `.env.example` - Environment variable names
- `go.mod` - Existing dependencies (pgx/v5)

### 3. Dependencies

**Required (already in go.mod):**
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/jackc/pgx/v5/stdlib` - Database/sql bridge

**New dependencies needed:**
- Add `github.com/jackc/pgconn` for connection pooling configuration

**Prerequisites:**
- Domain models must be defined (RDL-002 is done - ✅)
- PostgreSQL database must exist with schema from Rails app
- Environment variables configured in `.env` (based on `.env.example`)

**Blocking issues:**
- None - this task depends on RDL-002 (models/DTOs) being complete

### 4. Code Patterns

**Interface naming:**
- `ProjectRepository` and `LogRepository` (PascalCase)
- Method names: `GetByID`, `GetAll`, `GetByProjectID`, `GetWithLogs`

**Implementation naming:**
- `ProjectRepositoryImpl` and `LogRepositoryImpl`
- Each struct contains `conn *pgx.Conn` or `db *sql.DB`

**Method signatures:**
```go
// ProjectRepository interface
GetByID(ctx context.Context, id int64) (*models.Project, error)
GetAll(ctx context.Context) ([]*models.Project, error)

// LogRepository interface  
GetByID(ctx context.Context, id int64) (*models.Log, error)
GetByProjectID(ctx context.Context, projectID int64) ([]*models.Log, error)
GetAll(ctx context.Context) ([]*models.Log, error)
```

**Error handling:**
- Return `nil, nil` on success
- Return `nil, fmt.Errorf` on error with descriptive message
- Handle `pgx.ErrNoRows` as "not found" errors

**Context usage:**
- Accept context as first parameter
- Use `context.WithTimeout(ctx, 5*time.Second)` for query execution
- Propagate context to pgx `Exec`, `Query`, `QueryRow` calls

### 5. Testing Strategy

**Unit tests:**
- Test interface contracts in `internal/repository/*_test.go`
- Test PostgreSQL implementations in `internal/adapter/postgres/*_test.go`
- Verify context timeout behavior
- Verify error handling (NoRows, connection failures)

**Integration tests:**
- Use test database (separate from development)
- Test full repository round-trips (insert and retrieve)
- Verify connection pooling configuration

**Verification steps:**
1. Run `go test ./internal/repository/...` - tests pass
2. Run `go test ./internal/adapter/postgres/...` - tests pass
3. Run `go build ./...` - no errors
4. Verify connection pool metrics via pprof or logging

### 6. Risks and Considerations

**Blocking issues:**
- None identified

**Trade-offs:**
- Using pgx/v5 directly vs database/sql wrapper (more control, slightly more verbose)
- Returning domain models from repository (repository layer responsibility)
- Context timeout per method (not at adapter level - more granular control)

**Implementation considerations:**
- All nullable DB fields require `sql.NullXxx` type handling
- `time.Time` for datetime fields, `*string` for date fields (Rails uses date type)
- Connection pooling must be configured at initialization (not per-query)
- pgx connection pool size defaults: Min=0, Max=4 (need to override)

**Deployment considerations:**
- PostgreSQL must be running before server starts
- Environment variables must be set (DB_HOST, DB_USER, DB_PASS, DB_DATABASE)
- No migration tool in Phase 1 (schema must pre-exist)
<!-- SECTION:PLAN:END -->
