# Go Reading Log API

[![Go Version](https://img.shields.io/badge/Go-1.25.7-blue)](#)
[![License](https://img.shields.io/badge/License-MIT-green)](#)
[![Build Status](https://img.shields.io/badge/Status-Phase_1-orange)](#)

> **Phase 1 - Read-Only API**: This is a RESTful backend service built in Go following Clean Architecture principles. It serves as a migration from an existing Rails application, providing endpoints for managing reading projects and their associated logs.

---

## 📖 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Prerequisites](#-prerequisites)
- [Getting Started](#-getting-started)
- [Environment Setup](#-environment-setup)
- [Running the Application](#-running-the-application)
- [API Documentation](#-api-documentation)
- [Testing](#-testing)
- [Project Structure](#-project-structure)
- [Configuration](#-configuration)
- [Architecture](#-architecture)
- [Database Schema](#-database-schema)
- [Troubleshooting](#-troubleshooting)
- [Related Documentation](#-related-documentation)

---

## 🎯 Overview

The Go Reading Log API is a production-ready RESTful backend service that provides:

- **Read-only endpoints** for projects and logs (Phase 1)
- **PostgreSQL** database with connection pooling
- **Structured logging** using Go's `log/slog` package
- **Comprehensive error handling** with consistent patterns
- **Middleware chain** for CORS, request ID, logging, and panic recovery

> **Note**: This is a **Phase 1** project with read-only API endpoints. Create/Update/Delete operations will be added in Phase 2.

---

## ✨ Features

- ✅ RESTful API with versioned endpoints (`/api/v1/`)
- ✅ PostgreSQL database with connection pooling (pgx/v5)
- ✅ Structured logging using Go's `log/slog`
- ✅ Clean Architecture implementation
- ✅ Comprehensive test coverage (100+ tests)
- ✅ Middleware chain for cross-cutting concerns
- ✅ Graceful shutdown handling

---

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

| Tool | Version | Notes |
|------|---------|-------|
| Go | 1.25.7 or later | [Install Go](https://golang.org/doc/install) |
| PostgreSQL | 13 or later | [Install PostgreSQL](https://www.postgresql.org/download/) |
| Make | Any recent version | For convenience commands |

### Verify installations

```bash
go version
psql --version
```

---

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd go-reading-log-api-next
```

### 2. Set up environment variables

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` with your database credentials:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_database_user
DB_PASS=your_database_password
DB_DATABASE=reading_log

# Optional
SERVER_PORT=3000
SERVER_HOST=0.0.0.0
LOG_LEVEL=info
LOG_FORMAT=text
```

### 3. Create the databases

```bash
# Using psql or any PostgreSQL client
CREATE DATABASE reading_log;
CREATE DATABASE reading_log_test;
```

Or using Docker:

```bash
docker run --name reading-log-db -p 5432:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=reading_log \
  -d postgres:15
```

### 4. Run the server

```bash
# Using make (recommended)
make run

# Or using go run directly
go run ./cmd/server.go
```

The server will start on `http://0.0.0.0:3000`.

### 5. Verify the API

```bash
# Health check
curl http://localhost:3000/healthz
# Expected: {"status":"ok"}
```

---

## 🔧 Environment Setup

### Required Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DB_USER` | PostgreSQL username | `postgres` |
| `DB_PASS` | PostgreSQL password | `secret123` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_DATABASE` | Database name | `reading_log` |

### Optional Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `3000` | Server listening port |
| `SERVER_HOST` | `0.0.0.0` | Server listening host |
| `LOG_LEVEL` | `info` | Logging level: debug, info, warn, error |
| `LOG_FORMAT` | `text` | Log format: text or json |

### Test Database

For testing, the test database name defaults to `<DB_DATABASE>_test`. You can override with:

| Variable | Description | Example |
|----------|-------------|---------|
| `DB_DATABASE_TEST` | Test database name | `reading_log_test` |

---

## ▶️ Running the Application

### Development Mode

```bash
# Start the server
make run

# Or run directly
go run ./cmd/server.go
```

### Build for Production

```bash
# Build binary
make build

# Or manually
go build -o server ./cmd

# Run the binary
./server
```

### Make Commands

| Command | Description |
|---------|-------------|
| `make run` | Build and run the server |
| `make build` | Build the binary to `bin/server` |
| `make test` | Run all tests |
| `make test-verbose` | Run tests with verbose output |
| `make test-coverage` | Run tests and generate coverage report |
| `make fmt` | Format code with `go fmt` |
| `make vet` | Run `go vet` for static analysis |
| `make clean` | Clean up build artifacts |
| `make start-pg` | Start PostgreSQL via Docker |
| `make docker-start-pg` | Explicitly start PostgreSQL via Docker |

---

## 📡 API Documentation

The API is versioned under `/api/v1/`. All responses are in JSON format.

### Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/healthz` | GET | Health check endpoint |
| `/api/v1/projects` | GET | List all projects |
| `/api/v1/projects/:id` | GET | Get project by ID |
| `/api/v1/projects/:project_id/logs` | GET | Get logs for a project |

### Example Requests

#### Health Check

```bash
curl http://localhost:3000/healthz
```

**Response:**

```json
{
  "status": "ok"
}
```

#### List Projects

```bash
curl http://localhost:3000/api/v1/projects
```

**Response:**

```json
{
  "projects": [
    {
      "id": 1,
      "name": "The Great Gatsby",
      "total_page": 180,
      "page": 45,
      "started_at": "2024-01-01T00:00:00Z",
      "progress": "25%",
      "status": "reading",
      "logs_count": 5,
      "days_unread": 3,
      "median_day": "2024-01-05T00:00:00Z",
      "finished_at": null,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-05T00:00:00Z"
    }
  ]
}
```

#### Get Project by ID

```bash
curl http://localhost:3000/api/v1/projects/1
```

**Response:**

```json
{
  "project": {
    "id": 1,
    "name": "The Great Gatsby",
    "total_page": 180,
    "page": 45,
    "started_at": "2024-01-01T00:00:00Z",
    "progress": "25%",
    "status": "reading",
    "logs_count": 5,
    "days_unread": 3,
    "median_day": "2024-01-05T00:00:00Z",
    "finished_at": null,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-05T00:00:00Z"
  }
}
```

#### Get Logs for a Project

```bash
curl http://localhost:3000/api/v1/projects/1/logs
```

**Response:**

```json
{
  "logs": [
    {
      "id": 1,
      "project_id": 1,
      "data": "2024-01-01",
      "start_page": 1,
      "end_page": 20,
      "wday": 0,
      "note": "Started reading",
      "text": "Chapter 1",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### Error Responses

| Status Code | Error Type | Description |
|-------------|------------|-------------|
| `404` | Not Found | Resource not found |
| `500` | Internal Server Error | Database or server error |

**Error Response Format:**

```json
{
  "error": "Resource not found"
}
```

---

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Generate coverage report
make test-coverage
```

### Using go test directly

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

### Test Coverage

| Package | Coverage |
|---------|----------|
| `internal/api/v1` | 100.0% |
| `internal/api/v1/handlers` | 96.9% |
| `internal/api/v1/middleware` | 100.0% |
| `internal/config` | 86.7% |
| `internal/domain/dto` | 93.3% |
| `internal/domain/models` | 100.0% |
| `internal/logger` | 100.0% |
| `test` | 35.0% |

**Total: 105 tests passing**

---

## 📁 Project Structure

The project follows **Clean Architecture** with clear separation of concerns:

```
go-reading-log-api-next/
├── cmd/                          # Application entry points
│   └── server.go                 # Main application entry point
│
├── internal/                     # Private application code
│   ├── adapter/                  # Database adapters
│   │   └── postgres/             # PostgreSQL implementation
│   │
│   ├── api/                      # HTTP layer
│   │   └── v1/                   # API version 1
│   │       ├── handlers/         # Request handlers
│   │       ├── middleware/       # HTTP middleware
│   │       └── routes.go         # Router setup
│   │
│   ├── config/                   # Configuration
│   │   ├── config.go
│   │   └── config_test.go
│   │
│   ├── domain/                   # Business logic
│   │   ├── dto/                  # Data Transfer Objects
│   │   ├── models/               # Domain models
│   │   └── repository/           # Repository interfaces
│   │
│   └── logger/                   # Logging setup
│       ├── logger.go
│       └── logger_test.go
│
├── test/                         # Test infrastructure
│   ├── integration/              # Integration tests
│   ├── unit/                     # Unit tests
│   └── test_helper.go            # Test utilities and mocks
│
├── docs/                         # Detailed documentation
│   └── README.go-project.md      # Project structure documentation
│
├── rails-app/                    # Original Rails application (reference)
│
├── backlog/                      # Backlog.md task management
│
├── Makefile                      # Development commands
├── .env.example                  # Environment variable template
├── go.mod                        # Go module definition
└── README.md                     # This file
```

### Layer Responsibilities

| Layer | Responsibility | Key Components |
|-------|----------------|----------------|
| **cmd/** | Application entry point | `server.go` - main(), server setup |
| **api/** | HTTP layer | Handlers, middleware, routing |
| **domain/** | Business logic | Models, DTOs, repository interfaces |
| **adapter/** | Data access | PostgreSQL implementations |
| **config/** | Configuration | Environment variable loading |
| **logger/** | Logging | Structured logging with slog |

---

## ⚙️ Configuration

### Configuration Sources

Configuration is loaded from environment variables with sensible defaults.

### Configurable Options

| Option | Environment Variable | Default | Description |
|--------|---------------------|---------|-------------|
| Server Port | `SERVER_PORT` | `3000` | HTTP server port |
| Server Host | `SERVER_HOST` | `0.0.0.0` | HTTP server host |
| Log Level | `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| Log Format | `LOG_FORMAT` | `text` | Log format (text, json) |
| DB Host | `DB_HOST` | `localhost` | PostgreSQL host |
| DB Port | `DB_PORT` | `5432` | PostgreSQL port |
| DB User | `DB_USER` | (required) | PostgreSQL username |
| DB Password | `DB_PASS` | (required) | PostgreSQL password |
| DB Database | `DB_DATABASE` | `reading_log` | Database name |

### Database Connection String

The connection string is built automatically from environment variables:

```go
postgresql://user:password@host:port/database?sslmode=disable
```

> **Production Note**: Use SSL for production deployments. Replace `sslmode=disable` with `sslmode=verify-full&sslrootcert=/path/to/ca.pem`.

---

## 🏗️ Architecture

### Clean Architecture Layers

```
┌──────────────────────────────────────────────────────────────┐
│                    cmd/ (Entry Point)                        │
│                  main.go / server.go                         │
└──────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────┐
│                   internal/api/                              │
│         HTTP Handlers & Routing (Controller)                 │
└──────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────┐
│                  internal/domain/                            │
│            Business Logic & Domain Models                    │
└──────────────────────────────────────────────────────────────┘
                            ↓
┌──────────────────────────────────────────────────────────────┐
│                  internal/adapter/                           │
│              Infrastructure / Data Access                    │
└──────────────────────────────────────────────────────────────┘
```

### Dependency Flow

Dependencies flow **inward** through the layers:

1. **cmd/** depends on all layers
2. **api/** depends on domain and adapter
3. **domain/** defines interfaces (repository contracts)
4. **adapter/** implements domain interfaces

### Key Patterns

- **Dependency Injection**: Repositories injected into handlers
- **Context Propagation**: Timeout/cancellation via `context.Context`
- **Error Wrapping**: Consistent `fmt.Errorf("...: %w", err)` pattern
- **Middleware Chain**: Composition of middleware functions

---

## 🗄️ Database Schema

The application uses PostgreSQL with the following tables:

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

> **Note**: The computed columns (`progress`, `status`, `logs_count`, etc.) are populated by PostgreSQL queries in the application.

---

## 🛠️ Troubleshooting

### Database Connection Failed

```bash
# Check database is running
pg_isready -h localhost -p 5432

# Check connection string
echo "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=disable"

# Verify database exists
psql -U postgres -c "SELECT datname FROM pg_database WHERE datname = 'reading_log';"
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

### Common Issues

| Issue | Solution |
|-------|----------|
| `connection refused` | Start PostgreSQL server |
| `password authentication failed` | Check `DB_USER` and `DB_PASS` in `.env` |
| `database does not exist` | Create the database with `CREATE DATABASE` |
| `address already in use` | Change `SERVER_PORT` or stop conflicting process |
| `module not found` | Run `go mod tidy` to download dependencies |

---

## 📚 Related Documentation

| Document | Description |
|----------|-------------|
| [docs/README.go-project.md](docs/README.go-project.md) | Detailed project structure and architecture documentation |
| [QWEN.md](QWEN.md) | Project context for AI assistant (Qwen) |
| [AGENTS.md](AGENTS.md) | Project agent guidelines (MCP workflow) |
| [Makefile](Makefile) | Development commands and scripts |

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Workflow

```bash
# 1. Run tests to ensure everything passes
make test

# 2. Format code
make fmt

# 3. Run static analysis
make vet

# 4. Make your changes

# 5. Run tests again
make test

# 6. Commit and push
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Original Rails application: `rails-app/`
- Go community for excellent tools and documentation
- Clean Architecture principles by Robert C. Martin

---

## 📞 Support

For support, please open an issue in the GitHub repository or contact the development team.

---

*Last updated: 2026-04-01* |
*Phase 1 - Read-Only API*
