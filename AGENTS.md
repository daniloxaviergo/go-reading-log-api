# Go Reading Log API - Project Context

## Project Overview

This is a RESTful backend API service built in Go following **Clean Architecture** principles. It serves as a migration from an existing Rails application, providing endpoints for managing reading projects and their associated logs.

Use MCP backlog.

### Key Technologies

- **Language:** Go 1.25.7
- **Database:** PostgreSQL with connection pooling (pgx/v5)
- **Web Framework:** Standard library `net/http` with Gorilla Mux for routing
- **Architecture:** Clean Architecture with layered separation of concerns

### Project Status

This is a Phase 1 migration project. Notable characteristics:
- No database migration tool (schema managed manually)
- Direct PostgreSQL queries without an ORM
- Clean Architecture structure for maintainability and testability

## Application Architecture

The project follows Clean Architecture with these layers:

```
cmd/              → Entry point (server.go)
internal/api/     → HTTP layer (handlers, middleware, routes)
internal/domain/  → Business logic (models, DTOs)
internal/repository/ → Repository interfaces
internal/adapter/ → Infrastructure (PostgreSQL implementations)
test/             → Test infrastructure
```

### Layer Responsibilities

| Layer | Responsibility | Key Components |
|-------|----------------|----------------|
| **cmd/** | Application entry point | `server.go` - main(), server setup |
| **api/** | HTTP layer | Handlers, middleware, routing |
| **domain/** | Business logic | Models, DTOs |
| **repository/** | Repository interfaces | ProjectRepository, LogRepository interfaces |
| **adapter/** | Data access | PostgreSQL repository implementations |
| **config/** | Configuration | Environment variable loading |
| **logger/** | Logging | Structured logging with slog |

## Building and Running

### Prerequisites

- Go 1.25.7 or later
- PostgreSQL 13 or later
- PostgreSQL must be running and accessible

### Environment Setup

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your database credentials:
   ```bash
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_database_user
   DB_PASS=your_database_password
   DB_DATABASE=reading_log
   ```

3. Create the databases:
   ```sql
   CREATE DATABASE reading_log;
   CREATE DATABASE reading_log_test;
   ```

### Running the Server

**Direct Run (Local Development):**
```bash
# Build the application
go build -o server ./cmd

# Run the server
./server

# Or run directly
go run ./cmd/server.go
```

**Docker Compose (Containerized):**
```bash
# Start all services (PostgreSQL, Go API, Rails API)
make docker-up

# Stop all services
make docker-down

# View logs
make docker-logs

# List containers
make docker-ps
```

The Go API starts on `http://0.0.0.0:3000` and the Rails API on `http://0.0.0.0:3001`.

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./internal/config/...
go test ./test/...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out
```

### Development Commands

```bash
# Format code
go fmt ./...

# Run linter
go vet ./...

# Build for production
go build -o bin/server ./cmd/server.go
```

### Docker Compose Commands

```bash
# Start all services via docker-compose
make docker-up

# Stop all services via docker-compose
make docker-down

# View logs from all services
make docker-logs

# List running containers
make docker-ps

# Stop only PostgreSQL container
make docker-stop-pg
```

## API Endpoints

The API uses version `/v1/` prefix (not `/api/v1/`). All endpoints return JSON responses.

### Base URL
```
http://localhost:3000
```

### Route Prefix
All API routes use `/v1/` prefix (note: not `/api/v1/`).

---

## Endpoints

### Health Check

| Property | Value |
|----------|-------|
| **Method** | GET |
| **Path** | `/healthz` |
| **Description** | Returns health status of the API service |
| **Authentication** | None |
| **Response Code** | 200 OK |

**Request:**
```bash
curl http://localhost:3000/healthz
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "message": "API is running"
}
```

---

### Projects Endpoints

#### List All Projects

| Property | Value |
|----------|-------|
| **Method** | GET |
| **Path** | `/v1/projects.json` |
| **Description** | Returns all projects with eager-loaded logs (first 4) and calculated fields |
| **Authentication** | None |
| **Response Code** | 200 OK |

**Request:**
```bash
curl http://localhost:3000/v1/projects.json
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Project Name",
    "total_page": 200,
    "page": 50,
    "started_at": "2024-01-15T10:30:00Z",
    "progress": 25.0,
    "status": "running",
    "logs_count": 4,
    "days_unreading": 5,
    "median_day": 10.0,
    "finished_at": "2024-02-15T00:00:00Z",
    "logs": [
      {
        "id": 1,
        "data": "2024-01-15T10:30:00",
        "start_page": 0,
        "end_page": 25,
        "note": "Morning reading",
        "project": {
          "id": 1,
          "name": "Project Name",
          "total_page": 200,
          "page": 50
        }
      }
    ]
  }
]
```

#### Get Project by ID

| Property | Value |
|----------|-------|
| **Method** | GET |
| **Path** | `/v1/projects/{id}.json` |
| **Description** | Returns a single project by ID with eager-loaded logs and calculated fields |
| **Authentication** | None |
| **Response Code** | 200 OK, 404 Not Found |

**Request:**
```bash
curl http://localhost:3000/v1/projects/1.json
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Project Name",
  "total_page": 200,
  "page": 50,
  "started_at": "2024-01-15T10:30:00Z",
  "progress": 25.0,
  "status": "running",
  "logs_count": 4,
  "days_unreading": 5,
  "median_day": 10.0,
  "finished_at": "2024-02-15T00:00:00Z",
  "logs": [
    {
      "id": 1,
      "data": "2024-01-15T10:30:00",
      "start_page": 0,
      "end_page": 25,
      "note": "Morning reading"
    }
  ]
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "project not found"
}
```

#### Create Project

| Property | Value |
|----------|-------|
| **Method** | POST |
| **Path** | `/v1/projects.json` |
| **Description** | Creates a new reading project |
| **Authentication** | None |
| **Response Code** | 201 Created, 400 Bad Request |

**Request:**
```bash
curl -X POST http://localhost:3000/v1/projects.json \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Reading Project",
    "total_page": 200,
    "page": 0,
    "started_at": "2024-01-15T10:30:00Z",
    "reinicia": false
  }'
```

**Request Body Schema:**
```json
{
  "name": "string (required, max 255)",
  "total_page": "integer (required, must be > 0)",
  "page": "integer (required, must be <= total_page)",
  "started_at": "string (optional, RFC3339 format)",
  "reinicia": "boolean (optional, default: false)"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "name": "My Reading Project",
  "total_page": 200,
  "page": 0,
  "started_at": "2024-01-15T10:30:00Z"
}
```

**Error Response (400 Bad Request) - Validation Failed:**
```json
{
  "error": "validation failed",
  "details": {
    "page": "page (100) cannot exceed total_page (50)",
    "total_page": "total_page (0) must be greater than 0"
  }
}
```

**Error Response (400 Bad Request) - Invalid Date:**
```json
{
  "error": "invalid date format",
  "details": {
    "started_at": "must be in RFC3339 format"
  }
}
```

---

### Logs Endpoints

#### List Logs for Project

| Property | Value |
|----------|-------|
| **Method** | GET |
| **Path** | `/v1/projects/{project_id}/logs.json` |
| **Description** | Returns first 4 logs for a project, ordered by date DESC |
| **Authentication** | None |
| **Response Code** | 200 OK, 404 Not Found |

**Request:**
```bash
curl http://localhost:3000/v1/projects/1/logs.json
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "data": "2024-01-15T10:30:00",
    "start_page": 0,
    "end_page": 25,
    "note": "Morning reading",
    "project": {
      "id": 1,
      "name": "Project Name",
      "total_page": 200,
      "page": 50,
      "started_at": "2024-01-15T10:30:00Z",
      "status": "running",
      "progress": 25.0
    }
  },
  {
    "id": 2,
    "data": "2024-01-14T09:00:00",
    "start_page": 25,
    "end_page": 50,
    "note": "Evening reading",
    "project": {
      "id": 1,
      "name": "Project Name",
      "total_page": 200,
      "page": 50,
      "started_at": "2024-01-15T10:30:00Z",
      "status": "running",
      "progress": 25.0
    }
  }
]
```

**Error Response (404 Not Found):**
```json
{
  "error": "project not found"
}
```

---

## Calculated Fields

The API computes several derived fields for each project:

| Field | Type | Description | Formula |
|-------|------|-------------|---------|
| `progress` | float | Percentage of book completed | `(page / total_page) * 100` |
| `status` | string | Project status | Determined by page/total_page and started_at |
| `logs_count` | int | Number of log entries | `len(logs)` |
| `days_unreading` | int | Days since last reading activity | Calculated from logs data |
| `median_day` | float | Pages per day (rounded to 2 decimals) | `round(page / days_reading, 2)` |
| `finished_at` | datetime | Estimated completion date | Based on median_day calculation |

### Status Values

The `status` field can have one of these values:
- `unstarted` - Project not yet started
- `running` - Currently reading
- `sleeping` - Paused reading
- `stopped` - Stopped reading
- `finished` - Completed the book

---

## Error Handling

All error responses follow this format:

```json
{
  "error": "error_type_or_message",
  "details": {
    "field_name": "human-readable error description"
  }
}
```

### HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created (POST) |
| 400 | Bad Request - validation errors or invalid data |
| 404 | Not Found - resource doesn't exist |
| 500 | Internal Server Error |

### Validation Error Types

| Code | Field | Description |
|------|-------|-------------|
| `page_invalid` | page | Page number is negative |
| `page_exceeds_total` | page | Page exceeds total_page |
| `total_page_invalid` | total_page | Total page is zero or negative |
| `invalid_status` | status | Status is not one of the valid values |
| `start_page_exceeds_end` | start_page | Start page exceeds end page |

---

## Phase 1 Limitations

The current implementation (Phase 1) provides **read-only** access to logs:

- ✅ GET `/v1/projects.json` - List projects
- ✅ GET `/v1/projects/{id}.json` - Get project details
- ✅ POST `/v1/projects.json` - Create projects
- ✅ GET `/v1/projects/{project_id}/logs.json` - List logs

**Not implemented (Phase 2):**
- ❌ POST `/v1/projects/{project_id}/logs.json` - Create log
- ❌ PUT `/v1/logs/{id}.json` - Update log
- ❌ DELETE `/v1/logs/{id}.json` - Delete log

---

## Quick Reference

### Complete curl Examples

```bash
# Health check
curl http://localhost:3000/healthz

# List all projects
curl http://localhost:3000/v1/projects.json

# Get project by ID
curl http://localhost:3000/v1/projects/1.json

# Create a new project
curl -X POST http://localhost:3000/v1/projects.json \
  -H "Content-Type: application/json" \
  -d '{"name":"My Book","total_page":300,"page":0}'

# Get logs for a project
curl http://localhost:3000/v1/projects/1/logs.json
```

## Environment Variables

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DB_USER` | PostgreSQL username | `postgres` |
| `DB_PASS` | PostgreSQL password | `secret123` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_DATABASE` | Database name | `reading_log` |

### Optional Variables (with defaults)

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `3000` | Server listening port |
| `SERVER_HOST` | `0.0.0.0` | Server listening host |
| `LOG_LEVEL` | `info` | Logging level: debug, info, warn, error |
| `LOG_FORMAT` | `text` | Log format: text or json |

### Docker Compose Configuration

When using Docker Compose, both applications connect to a shared PostgreSQL container:

| Variable | Description | Docker Value |
|----------|-------------|--------------|
| `DB_HOST` | PostgreSQL hostname | `postgres` (service name) |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL username | `postgres` |
| `DB_PASS` | PostgreSQL password | `postgres` |
| `DB_DATABASE` | Database name | `reading_log` |
| `SERVER_PORT` | Go API port | `3000` |
| `PORT` | Rails API port | `3001` |

**Port Conflict Resolution:** The Go API uses port 3000, while the Rails API uses port 3001 to avoid conflicts.

## Database Schema

### Projects Table

```sql
CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    total_page INT NOT NULL DEFAULT 0,
    started_at TIMESTAMP WITH TIME ZONE,
    page INT NOT NULL DEFAULT 0,
    reinicia BOOLEAN NOT NULL DEFAULT false,
    progress VARCHAR(255),
    status VARCHAR(255),
    logs_count INT DEFAULT 0,
    days_unread INT DEFAULT 0,
    median_day VARCHAR(255),  -- Stores calculated float as string (e.g., "5.5")
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

**Note:** The `median_day` field is stored as a VARCHAR but contains the calculated float value as a string. The API returns this as a `float64` in the response.

### Logs Table

```sql
CREATE TABLE logs (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    data TIMESTAMP WITHOUT TIME ZONE,
    start_page INT NOT NULL DEFAULT 0,
    end_page INT NOT NULL DEFAULT 0,
    wday INT NOT NULL DEFAULT 0,
    note TEXT,
    text TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for optimized JOIN and ORDER BY queries
CREATE INDEX index_logs_on_project_id ON logs(project_id);
CREATE INDEX index_logs_on_project_id_and_data_desc ON logs(project_id, data DESC);
```

### Users Table

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Watsons Table

```sql
CREATE TABLE watsons (
    id BIGSERIAL PRIMARY KEY,
    start_at TIMESTAMP WITH TIME ZONE,
    end_at TIMESTAMP WITH TIME ZONE,
    minutes INT,
    external_id VARCHAR(255),
    log_id BIGINT,
    project_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX index_watsons_on_log_id ON watsons(log_id);
CREATE INDEX index_watsons_on_project_id ON watsons(project_id);
```

## Development Conventions

### Context Usage

All database operations use a context with a 5-second timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

project, err := repo.GetByID(ctx, id)
```

### Error Handling

Consistent error wrapping pattern:

```go
return nil, fmt.Errorf("failed to get project: %w", err)
```

### Repository Pattern

Repository interfaces define the contract in `internal/repository/`:

```go
type ProjectRepository interface {
    GetByID(ctx context.Context, id int64) (*models.Project, error)
    GetAll(ctx context.Context) ([]*models.Project, error)
    GetWithLogs(ctx context.Context, id int64) (*dto.ProjectResponse, error)
}
```

Adapters implement the interface in `internal/adapter/postgres/`.

### Dependency Injection

Repositories are injected into handlers:

```go
projectRepo := postgres.NewProjectRepositoryImpl(dbPool)
handler := handlers.NewProjectsHandler(projectRepo)
```

### Logging

Uses Go's `log/slog` package with structured logging:

```go
log.Info("Starting server...", "host", cfg.ServerHost, "port", cfg.ServerPort)
log.Error("Database connection failed", "error", err)
```

### Testing Strategy

- **Unit tests:** Use mock repositories to test business logic without database
- **Integration tests:** Use `TestHelper` from `test/test_helper.go` for database setup/teardown

```go
// Unit test example
mockRepo := test.NewMockProjectRepository()
handler := handlers.NewProjectsHandler(mockRepo)

// Integration test example
helper, err := test.SetupTestDB()
if err != nil {
    t.Fatal(err)
}
defer helper.Close()
```

## Code Patterns

### Derived Calculation Methods

The `Project` model includes several derived calculation methods that compute values dynamically:

```go
// CalculateLogsCount calculates logs_count as len(logs)
// Matches Rails behavior: def logs_count; logs.size; end
func (p *Project) CalculateLogsCount(logs []*dto.LogResponse) *int

// CalculateMedianDay calculates median_day as (page / days_reading).round(2)
// where days_reading is the number of days since started_at
// Returns 0.00 for edge cases (zero/negative days_reading, no started_at)
func (p *Project) CalculateMedianDay() *float64
```

**Formula:** `median_day = (page / days_reading).round(2)` where `days_reading = (today - started_at).days`

### Handler Pattern

Handlers receive HTTP requests and return responses:

```go
func (h *ProjectsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id := getPathParam(r, "id")
    project, err := h.repo.GetByID(r.Context(), id)
    if err != nil {
        h.handleError(w, err)
        return
    }
    h.respond(w, http.StatusOK, project)
}
```

### Middleware Chain

Middleware is chained in `cmd/server.go`:

```go
middlewareChain := middleware.Chain(router,
    middleware.RecoveryMiddleware,
    middleware.CORSMiddleware,
    middleware.RequestIDMiddleware,
    middleware.LoggingMiddleware,
)
```

## Important Notes

### Go Version

The project uses Go 1.25.7. Verify this is intentional or adjust as needed.

### No Database Migrations

Phase 1 has no migration tool. Schema management is done manually. For production, consider adding:
- [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- [pressly/goose](https://github.com/pressly/goose)

### SSL Mode

The connection string uses `sslmode=disable`. For production, use:
```
?sslmode=verify-full&sslrootcert=/path/to/ca.pem
```

### Error Handling

All database errors use PostgreSQL error codes via `pgx`:
- `pgx.ErrNoRows` - No rows found
- Other `pgx` errors - Database operation failures

### Context Propagation

Context is embedded in domain models for timeout and cancellation propagation.

## Common Tasks

### Adding a New Endpoint

1. Create handler in `internal/api/v1/handlers/`
2. Add route in `internal/api/v1/routes.go`
3. Add middleware if needed in `internal/api/v1/middleware/`

### Adding a New Model

1. Create domain model in `internal/domain/models/`
2. Create DTO in `internal/domain/dto/` if needed for API response
3. Add repository interface in `internal/repository/`
4. Add adapter implementation in `internal/adapter/postgres/`

### Adding Middleware

1. Create middleware function in `internal/api/v1/middleware/`
2. Add to middleware chain in `cmd/server.go`

### Cleanup Procedures

#### Auto-cleanup mechanism
The integration test helper automatically cleans up test databases after each test run using `TestHelper.Close()`. This ensures no leftover data affects subsequent tests. The cleanup process includes:
- Dropping the current test database (e.g., `reading_log_test_12345_67890_123456789`)
- Cleaning up orphaned test databases from previous test runs by dropping all databases matching the pattern `reading_log_test_%`

Example usage in integration tests:
```go
helper, err := test.SetupTestDB()
if err != nil {
    t.Fatal(err)
}
defer helper.Close()

// Test cases here
```

#### Server Log Management

When using `make run`, server output is written to `server.log`. To prevent excessive log growth:
- Run `rm server.log` to manually clear the log
- Add a `clean-log` target to your Makefile for convenience:
  ```makefile
clean-log:
	rm -f server.log
  ```

#### Database naming conventions
All test databases follow a unique naming convention to ensure parallel test safety:

- Format: `reading_log_test_<pid>_<goroutine_id>_<timestamp>`
- Example: `reading_log_test_12345_67890_1620000000`
- The name includes:
  - Process ID (`os.Getpid()`)
  - Goroutine ID (extracted from stack trace)
  - Current timestamp in UnixNano format

This ensures each test runs in its own isolated database, preventing conflicts during parallel execution.

#### SQL cleanup patterns
The following SQL commands are used for database cleanup:

1. **Drop orphaned test databases**:
   ```sql
   SELECT datname FROM pg_database WHERE datname LIKE 'reading_log_test_%' AND datname != 'current_test_db';
   DROP DATABASE IF EXISTS <database_name>;
   ```

2. **Reset entire test database** (used in `make test-clean`):
   ```sql
   DROP DATABASE IF EXISTS reading_log_test;
   CREATE DATABASE reading_log_test;
   ```

3. **Clean specific tables within a database**:
   ```sql
   TRUNCATE TABLE logs CASCADE;
   TRUNCATE TABLE projects CASCADE;
   ```

⚠️ **WARNING:** Never run manual cleanup commands on production databases (`reading_log`). Always verify you're using the `reading_log_test` database.

#### Rationale for per-test database strategy
- **Isolation**: Each test runs in its own database, preventing side effects between tests.
- **Parallel execution**: Unique database names allow multiple tests to run simultaneously without interference.
- **Simplicity**: Dropping entire databases is safer than schema resets which could accidentally affect other environments.
- **Consistency**: Ensures a clean state for every test run, regardless of previous failures.

#### Troubleshooting steps
- **Test database conflicts**: If tests fail due to existing databases, manually drop all `reading_log_test_%` databases:
  ```bash
  psql -c "SELECT datname FROM pg_database WHERE datname LIKE 'reading_log_test_%';" | grep reading_log_test | xargs -I {} psql -c "DROP DATABASE IF EXISTS {};"
  ```
- **Slow cleanup**: If cleanup takes too long, ensure no other processes are using the test databases. Use `pg_terminate_backend` to kill lingering connections.
- **Permission issues**: Verify the database user has sufficient privileges to drop databases.

## Troubleshooting

### Database Connection Failed

```bash
# Check database is running
pg_isready -h localhost -p 5432

# Check connection string
echo "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=disable"
```

### Port Already in Use

```bash
# Find process using port 3000
lsof -i :3000

# Or use a different port
SERVER_PORT=8080 go run ./cmd/server.go
```

### Tests Failing

```bash
# Ensure test database exists
psql -U postgres -c "CREATE DATABASE reading_log_test;"

# Run tests with verbose output
go test -v ./...
```

### Docker Compose Troubleshooting

```bash
# Check docker-compose configuration
docker-compose config

# Check container status
docker-compose ps

# View specific service logs
docker-compose logs -f go-api
docker-compose logs -f rails-api

# Rebuild containers after code changes
docker-compose up -d --build

# Connect to PostgreSQL container
docker exec -it reading-log-db psql -U postgres -d reading_log
```

## Related Files

- `docs/README.go-project.md` - Detailed project structure documentation
- `rails-app/` - Original Rails application (reference)
- `AGENTS.md` - Project agent guidelines (MCP workflow)

## Qwen Settings

See `.qwen/settings.json` for editor/IDE configuration settings.

---

*Last updated: 2026-04-03*
