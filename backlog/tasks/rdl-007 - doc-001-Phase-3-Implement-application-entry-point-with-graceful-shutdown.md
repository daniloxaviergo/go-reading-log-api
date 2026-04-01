---
id: RDL-007
title: '[doc-001 Phase 3] Implement application entry point with graceful shutdown'
status: To Do
assignee:
  - catarina
created_date: '2026-04-01 00:58'
updated_date: '2026-04-01 02:38'
labels: []
dependencies: []
references:
  - 'PRD Section: Technical Decisions'
  - 'Implementation Checklist: Core Components'
  - 'PRD Section: Acceptance Criteria'
documentation:
  - doc-001
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement cmd/server.go as the application entry point with HTTP server setup using net/http router.

Implement graceful shutdown on SIGTERM signal with context-based timeout (5 seconds).

Wire up all middleware and routes including health check endpoint at /healthz.

Configure HTTP server with proper timeout settings and connection pool from config.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Application starts successfully on configured port
- [ ] #2 Graceful shutdown implemented with 5-second timeout
- [ ] #3 All routes registered correctly
- [ ] #4 Health check endpoint available at /healthz
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Implement a production-ready application entry point in `cmd/server.go` with:
- Configuration loading from `.env` file and environment variables using the existing `config.LoadConfig()` function
- PostgreSQL database connection using `pgx` driver with proper connection pooling
- HTTP server with timeout settings (read, write, idle) from config
- Graceful shutdown on SIGTERM using `context.Context` and `server.Shutdown()`
- Route registration using the existing `api.SetupRoutes()` function
- Health check endpoint at `/healthz` wired through the routes
- Structured logging configuration with `log/slog`
- Repository instantiation (PostgreSQL implementations) passed to handlers

The application will follow Go best practices for server setup with proper resource cleanup and error handling.

### 2. Files to Modify

**Primary file:**
- `cmd/server.go` - Complete rewrite to implement all required functionality

**Files to reference (read-only for implementation):**
- `internal/config/config.go` - Configuration struct and loading logic
- `internal/api/v1/routes.go` - Route registration function
- `internal/adapter/postgres/project_repository.go` - Repository implementation
- `internal/adapter/postgres/log_repository.go` - Repository implementation
- `internal/logger/logger.go` - Logger setup (may need to create if not exists)
- `internal/domain/dto/health_check_response.go` - Health check DTO
- `.env.example` - Expected environment variables

### 3. Dependencies

**Existing dependencies (already in go.mod):**
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/jackc/pgx/v5/stdlib` - database/sql bridge
- `github.com/joho/godotenv` - Environment variable loading
- `github.com/google/uuid` - Request ID generation

**No new dependencies required**

**Prerequisites:**
- PostgreSQL database must be running and accessible
- Environment variables must be configured (via `.env` or system env vars)
- Database schema must exist (tables: `projects`, `logs`)
- RDL-006 (API handlers) should be completed or in progress for handler wiring

### 4. Code Patterns

**Follow existing codebase patterns:**
- Use `context.WithTimeout` for all database operations (already done in repository implementations)
- Use `slog` for structured logging with INFO level in production
- Use `http.Server` with explicit timeout settings
- Use `server.Shutdown(ctx)` for graceful shutdown
- Wrap middleware in `middleware.Chain()` with correct order: Recovery → CORS → RequestID → Logging → Handler
- Pass context through request chain using `r.Context()`
- Use `middleware.DefaultTimeout` for request timeouts where applicable

**Naming conventions:**
- Variable names: camelCase (e.g., `server`, `dbConn`, `config`)
- Types: PascalCase (e.g., `Config`, `ProjectRepositoryImpl`)
- Constants: UPPER_SNAKE_CASE (e.g., `defaultShutdownTimeout`)

**Error handling:**
- Log errors before returning them
- Use formatted error messages with context
- Exit gracefully with non-zero status code on startup failures

### 5. Testing Strategy

**Unit tests for cmd/server.go:**
- Test server startup with valid config
- Test server startup with invalid config (wrong port, missing database)
- Test graceful shutdown duration (should complete within 5 seconds)
- Test health check endpoint responsiveness

**Integration tests (using existing test infrastructure):**
- Run against test database
- Verify all routes are registered correctly
- Test health check endpoint returns expected JSON
- Test full request lifecycle with middleware chain

**How to test:**
```bash
# Build the application
go build -o server ./cmd/server.go

# Run with test database
DATABASE_URL=postgres://testuser:testpass@localhost:5432/testdb ./server

# Test health endpoint
curl http://localhost:3000/healthz

# Test shutdown (SIGTERM)
kill -TERM $(pgrep server)
```

**What to verify:**
- Server starts and listens on configured port
- Health check endpoint returns valid JSON
- All API endpoints return expected responses
- Graceful shutdown completes within timeout
- No resource leaks (connections closed properly)

### 6. Risks and Considerations

**Blocking issues:**
- logger/logger.go doesn't exist yet - need to create logger setup function before server.go can configure it
- May need to verify health check response DTO doesn't require `ctx` field in JSON serialization
- Database repository implementations need proper `pgx.Conn` or `*sql.DB` connection

**Potential pitfalls:**
- Context timeout handling must be consistent across all layers
- Shutdown timeout must be less than Kubernetes/cluster termination grace period
- Database connection pool settings should match server max concurrency
- Health check should verify database connectivity, not just return static response

**Trade-offs:**
- Using `net/http` stdlib router instead of chi/echo (user preference)
- Context timeout of 5 seconds for database operations (matches PRD requirement)
- Graceful shutdown timeout of 5 seconds (matches PRD requirement)
- Logger initialization before config (to enable logging during startup)

**Deployment considerations:**
- Kubernetes: Set pod terminationGracePeriodSeconds > 10 seconds
- Docker: Use `docker stop` (sends SIGTERM) to test graceful shutdown
- Health checks: `/healthz` endpoint can be used for Kubernetes liveness/readiness probes
- Database: Ensure connection pool max equals or exceeds expected concurrency
<!-- SECTION:PLAN:END -->
