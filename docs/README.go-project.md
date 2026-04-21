# Go Reading Log API - Project Structure Documentation

This document provides a comprehensive overview of the Go Reading Log API project structure, architecture, and development setup.

## Table of Contents

- [Overview](#overview)
- [Application Architecture](#application-architecture)
- [Directory Structure](#directory-structure)
- [Environment Variables](#environment-variables)
- [Database Schema](#database-schema)
- [Run Commands](#run-commands)
- [Developer Onboarding](#developer-onboarding)

---

## Overview

The Go Reading Log API is a RESTful backend service built following **Clean Architecture** principles. It provides endpoints for managing reading projects and their associated logs, serving as a migration from an existing Rails application.

### Key Features

- RESTful API with versioned endpoints (`/api/v1/`)
- PostgreSQL database with connection pooling
- Structured logging using Go's `log/slog` package
- Comprehensive error handling
- Middleware chain for cross-cutting concerns (CORS, recovery, request ID, logging)

---

## Application Architecture

The project follows **Clean Architecture** with the following layered structure:

```
┌─────────────────────────────────────────────────────────────────┐
│                        cmd/ (Entry Point)                       │
│                     main.go / server.go                         │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    internal/api/                                │
│              HTTP Handlers & Routing (Controller)               │
│         - v1/handlers/  - Request handling logic                │
│         - v1/middleware/ - HTTP middleware                      │
│         - v1/routes.go  - Route registration                    │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                   internal/domain/                              │
│              Business Logic & Domain Models                     │
│         - models/       - Domain entities                       │
│         - dto/          - Data Transfer Objects                 │
│         - repository/   - Repository interfaces                 │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                   internal/adapter/                             │
│              Infrastructure / Data Access Layer                 │
│         - postgres/     - PostgreSQL implementations            │
└─────────────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

| Layer | Responsibility | Key Components |
|-------|----------------|----------------|
| **cmd/** | Application entry point | `server.go` - main() function, server setup |
| **api/** | HTTP layer | Handlers, middleware, routing, request/response objects |
| **domain/** | Business logic | Models, DTOs, repository interfaces |
| **adapter/** | Data access | PostgreSQL implementations, repository adapters |
| **config/** | Configuration | Environment variable loading, config struct |
| **logger/** | Logging | Structured logging setup using slog |

---

## Directory Structure

```
go-reading-log-api-next/
├── cmd/                          # Application entry points
│   └── server.go                 # Main application entry point
│
├── internal/                     # Private application code
│   ├── adapter/                  # Database adapters
│   │   └── postgres/             # PostgreSQL implementation
│   │       ├── project_repository.go
│   │       └── log_repository.go
│   │
│   ├── api/                      # HTTP layer
│   │   └── v1/                   # API version 1
│   │       ├── handlers/         # Request handlers
│   │       │   ├── health_handler.go
│   │       │   ├── projects_handler.go
│   │       │   └── logs_handler.go
│   │       ├── middleware/       # HTTP middleware
│   │       │   ├── cors.go
│   │       │   ├── logging.go
│   │       │   ├── middleware.go
│   │       │   ├── recovery.go
│   │       │   └── request_id.go
│   │       └── routes.go         # Router setup
│   │
│   ├── config/                   # Configuration
│   │   ├── config.go
│   │   └── config_test.go
│   │
│   ├── domain/                   # Business logic
│   │   ├── dto/                  # Data Transfer Objects
│   │   │   ├── health_check_response.go
│   │   │   ├── log_response.go
│   │   │   └── project_response.go
│   │   └── models/               # Domain models
│   │       ├── project.go
│   │       └── log.go
│   │
│   ├── repository/               # Repository interfaces
│   │   ├── project_repository.go
│   │   └── log_repository.go
│   │
│   └── logger/                   # Logging setup
│       ├── logger.go
│       └── logger_test.go
│
├── pkg/                          # Public reusable packages (future)
│
├── test/                         # Test infrastructure
│   ├── integration/              # Integration tests
│   ├── unit/                     # Unit tests
│   ├── test_helper.go            # Test utilities and mocks
│   └── test_helper_test.go       # Test helper tests
│
├── docs/                         # Documentation
│   └── README.go-project.md      # This file
│
├── rails-app/                    # Original Rails application (reference)
│
├── .env.example                  # Environment variable template
├── go.mod                        # Go module definition
├── go.sum                        # Go dependencies checksum
├── coverage.out                  # Test coverage report
├── coverage.html                 # HTML coverage report
└── AGENTS.md                     # Project agent guidelines
```

### File Descriptions

#### cmd/server.go
Main application entry point. Responsible for:
- Loading configuration
- Initializing logger
- Setting up database connection with pooling
- Creating repository instances
- Registering routes with middleware chain
- Starting HTTP server with graceful shutdown

#### internal/api/v1/handlers/
Contains request handler implementations:
- `health_handler.go` - Health check endpoint (`/healthz`)
- `projects_handler.go` - Project CRUD operations
- `logs_handler.go` - Log CRUD operations

#### internal/api/v1/middleware/
HTTP middleware components:
- `recovery.go` - Panic recovery middleware
- `cors.go` - CORS configuration middleware
- `request_id.go` - Request ID generation/propagation
- `logging.go` - Request logging middleware
- `middleware.go` - Middleware chaining utilities

**Note:** The middleware directory also contains test files (`*_test.go`) for the middleware components.

#### internal/domain/models/
Core domain entities with context support:
- `project.go` - Project domain model
- `log.go` - Log entry domain model

#### internal/domain/dto/
Data Transfer Objects for API responses:
- `health_check_response.go` - Health check response
- `project_response.go` - Project response with computed fields
- `log_response.go` - Log response structure

#### internal/adapter/postgres/
PostgreSQL implementations of repository interfaces:
- `project_repository.go` - Project repository implementation
- `log_repository.go` - Log repository implementation

#### test/test_helper.go
Test utilities and mock implementations:
- `TestHelper` - Database setup/teardown for integration tests
- `MockProjectRepository` - Mock for unit testing
- `MockLogRepository` - Mock for unit testing

---

## Environment Variables

Configuration is loaded from environment variables with sensible defaults.

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

### Test Database

For testing, use:
- `DB_DATABASE_TEST` - Test database name (defaults to `<DB_DATABASE>_test`)

### .env.example

```bash
# Database Configuration (PostgreSQL)
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_database_user
DB_PASS=your_database_password
DB_DATABASE=your_database_name

# Server Configuration
SERVER_PORT=3000
SERVER_HOST=0.0.0.0

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=text
```

---

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

---

## Database Schema

The application uses PostgreSQL with the following tables. Note that the Go implementation uses a more complete schema than the original Rails schema, with computed columns populated by PostgreSQL queries.

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
    median_day VARCHAR(255),
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

**Note on computed columns:** The `progress`, `status`, `logs_count`, `days_unread`, `median_day`, and `finished_at` columns are populated by PostgreSQL queries in the application. They represent computed values from the Rails application.

### Logs Table

```sql
CREATE TABLE logs (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    data VARCHAR(255),
    start_page INT NOT NULL DEFAULT 0,
    end_page INT NOT NULL DEFAULT 0,
    wday INT NOT NULL DEFAULT 0,
    note TEXT,
    text TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_logs_project_id ON logs(project_id);
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

**Note:** The Rails `schema.rb` file shows a simpler schema without computed columns. The Go application uses an extended schema that includes additional columns populated by database queries.

---

## Run Commands

### Starting the Server

```bash
# Build the application
go build -o server ./cmd

# Run the server
./server
```

Or run directly:

```bash
go run ./cmd/server.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./internal/config/...

# Run integration tests (requires database connection)
go test ./test/...
```

### Test Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# View HTML coverage report
go tool cover -html=coverage.out
```

### Development Commands

```bash
# Format code
go fmt ./...

# Run linter (if configured)
go vet ./...

# Build for production
go build -o bin/server ./cmd/server.go
```

---

## Developer Onboarding

### Prerequisites

- Go 1.25.7 or later
- PostgreSQL 13 or later
- Make sure PostgreSQL is running and accessible

### Getting Started

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-reading-log-api-next
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Create the database**
   ```bash
   # Using psql or any PostgreSQL client
   CREATE DATABASE reading_log;
   CREATE DATABASE reading_log_test;
   ```

4. **Run the server**
   ```bash
   go run ./cmd/server.go
   ```

5. **Verify the API**
   ```bash
   # Health check
   curl http://localhost:3000/healthz
   # Should return: {"status":"ok"}

   # Get all projects
   curl http://localhost:3000/api/v1/projects.json

   # Get a specific project by ID (e.g., 1)
   curl http://localhost:3000/api/v1/projects/1.json

   # Get logs for a project (e.g., project ID 1)
   curl http://localhost:3000/api/v1/projects/1/logs.json
   ```

### Code Patterns

#### Context Usage

All database operations accept a context with a 5-second timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

project, err := repo.GetByID(ctx, id)
```

#### Error Handling

Consistent error format throughout the application:

```go
return nil, fmt.Errorf("failed to get project: %w", err)
```

#### Repository Pattern

Repository interfaces define the contract:

```go
type ProjectRepository interface {
    GetByID(ctx context.Context, id int64) (*models.Project, error)
    GetAll(ctx context.Context) ([]*models.Project, error)
    GetWithLogs(ctx context.Context, id int64) (*dto.ProjectResponse, error)
}
```

Adapters implement the interface:

```go
type ProjectRepositoryImpl struct {
    pool *pgxpool.Pool
}

func (r *ProjectRepositoryImpl) GetByID(ctx context.Context, id int64) (*models.Project, error) {
    // Implementation
}
```

#### Dependency Injection

Repositories are injected into handlers:

```go
projectRepo := postgres.NewProjectRepositoryImpl(dbPool)
handler := handlers.NewProjectsHandler(projectRepo)
```

### Testing Strategy

#### Unit Tests

Use mock repositories to test business logic without database:

```go
mockRepo := test.NewMockProjectRepository()
handler := handlers.NewProjectsHandler(mockRepo)
```

#### Integration Tests

Use `TestHelper` for database setup/teardown:

```go
helper, err := test.SetupTestDB()
if err != nil {
    t.Fatal(err)
}
defer helper.Close()

err = helper.SetupTestSchema()
// ... test code ...
err = helper.CleanupTestSchema()
```

### Common Tasks

#### Adding a New Endpoint

1. Create handler in `internal/api/v1/handlers/`
2. Add route in `internal/api/v1/routes.go`
3. Add middleware if needed in `internal/api/v1/middleware/`

#### Adding a New Model

1. Create domain model in `internal/domain/models/`
2. Create DTO in `internal/domain/dto/` if needed for API response
3. Add repository interface in `internal/repository/`
4. Add adapter implementation in `internal/adapter/postgres/`

#### Adding Middleware

1. Create middleware function in `internal/api/v1/middleware/`
2. Add to middleware chain in `cmd/server.go`

---

## Important Notes

### Go Version

The project uses Go 1.25.7. This is a future version - verify this is intentional or adjust as needed.

### No Database Migrations

The PRD notes that there is no migration tool (Phase 1). Schema management is done manually. For production deployments, consider adding a migration tool like:
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

Context is embedded in domain models for timeout and cancellation propagation:

```go
func (p *Project) GetContext() context.Context {
    if p.ctx == nil {
        return context.Background()
    }
    return p.ctx
}
```

---

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

---

## Related Documentation

- [PRD Section: Files to Modify](backlog/docs/) - Original implementation plan
- [Rails Schema](rails-app/db/schema.rb) - Source of truth for database structure
- [Implementation Checklist: Documentation](backlog/docs/) - Documentation requirements
- [Key Requirements](backlog/docs/) - Project requirements

---

*Last updated: 2026-04-01*
