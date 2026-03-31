# Rails to Go Migration - Phase 1 Design

**Date**: 2026-03-31  
**Scope**: Phase 1 - Core Reading Log API  
**Status**: Approved

---

## 1. Overview

Migrate the Rails-based Reading Log API to Go, starting with core data endpoints for projects and logs. The Go application will be a drop-in replacement for the Rails API, maintaining exact endpoint compatibility while using a simplified stack (PostgreSQL only, no Redis/Sidekiq).

### 1.1 Scope

**Phase 1 includes:**
- `/api/v1/projects` - List and show projects
- `/api/v1/projects/:project_id/logs` - List logs for a project

**Phase 1 excludes:**
- Users endpoint (can be added later if needed)
- Dashboard/echart endpoints (deferred to Phase 2)
- Background jobs (all operations synchronous)
- Watson tracking (deferred to Phase 2)

---

## 2. Architecture

### 2.1 Project Structure

```
go-reading-log-api/
├── cmd/
│   └── server.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── v1/
│   │   │   ├── projects.go    # GET /api/v1/projects, GET /api/v1/projects/:id
│   │   │   └── logs.go        # GET /api/v1/projects/:id/logs
│   │   └── middleware/
│   │       └── cors.go        # CORS middleware
│   ├── models/
│   │   ├── project.go         # Project struct
│   │   └── log.go             # Log struct
│   ├── repository/
│   │   ├── project_repo.go    # Project DB operations
│   │   └── log_repo.go        # Log DB operations
│   └── config/
│       └── config.go          # App configuration
├── pkg/
│   └── errors.go              # Custom error handling
├── .env                       # Environment variables
├── go.mod
└── main.go
```

### 2.2 Component Responsibilities

| Component | Responsibility | Dependencies |
|-----------|----------------|--------------|
| `models` | Go structs mirroring DB tables | None |
| `repository` | SQL queries via `database/sql` | PostgreSQL |
| `api` | HTTP handlers, JSON serialization | `net/http`, `models`, `repository` |
| `middleware` | CORS, request logging | `net/http` |

---

## 3. Data Models

### 3.1 Project

| Field | Type | Description |
|-------|------|-------------|
| `id` | int64 | Primary key |
| `name` | string | Book/project name |
| `total_page` | int | Total pages in book |
| `started_at` | time.Time | Date reading started |
| `page` | int | Current page |
| `reinicia` | bool | Reset flag |

### 3.2 Log

| Field | Type | Description |
|-------|------|-------------|
| `id` | int64 | Primary key |
| `project_id` | int64 | Foreign key to projects |
| `data` | time.Time | Date of reading session |
| `start_page` | int | Starting page |
| `end_page` | int | Ending page |
| `wday` | int | Day of week (0-6) |
| `note` | string | Optional notes |
| `text` | string | Optional text content |

---

## 4. API Endpoints

### 4.1 List Projects

```
GET /api/v1/projects
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Book Title",
    "total_page": 200,
    "started_at": "2025-01-15",
    "page": 50,
    "reinicia": false
  }
]
```

### 4.2 Show Project

```
GET /api/v1/projects/:id
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Book Title",
  "total_page": 200,
  "started_at": "2025-01-15",
  "page": 50,
  "reinicia": false
}
```

**Response (404 Not Found):**
```json
{"error": "Project not found"}
```

### 4.3 List Project Logs

```
GET /api/v1/projects/:project_id/logs
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "project_id": 1,
    "data": "2025-01-15T10:30:00Z",
    "start_page": 1,
    "end_page": 25,
    "wday": 6,
    "note": "Chapter 1",
    "text": null
  }
]
```

**Response (404 Not Found):**
```json
{"error": "Project not found"}
```

---

## 5. Error Handling

| Status | Condition | Response |
|--------|-----------|----------|
| 200 | Success | Request body |
| 404 | Record not found | `{"error": "<resource> not found"}` |
| 422 | Validation error | `{"error": "Validation failed"}` |
| 500 | Internal error | `{"error": "Internal server error"}` (with logging) |

---

## 6. Configuration

Environment variables (loaded from `.env`):

| Variable | Description | Required |
|----------|-------------|----------|
| `DB_HOST` | PostgreSQL host | Yes |
| `DB_PORT` | PostgreSQL port (default: 5432) | No |
| `DB_USER` | PostgreSQL user | Yes |
| `DB_PASS` | PostgreSQL password | Yes |
| `DB_NAME` | Database name | Yes |
| `SERVER_PORT` | HTTP server port (default: 3000) | No |

---

## 7. Database Connection

### 7.1 Connection String

```go
postgres://user:password@host:port/dbname?sslmode=disable
```

### 7.2 Driver

- Use `github.com/lib/pq` PostgreSQL driver

### 7.3 Schema

No migration tool (Phase 1). Database schema must exist before running:
- Table `projects` with columns from Section 3.1
- Table `logs` with columns from Section 3.2

---

## 8. Testing

### 8.1 Unit Tests

- Repository layer: Test SQL queries with mock DB
- Models: Test struct methods (if any)

### 8.2 Integration Tests

- API endpoints: Test HTTP handlers with test server
- Use test database (SQLite for speed or separate PG db)

### 8.3 Test Commands

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
```

---

## 9. Build & Run

### 9.1 Dependencies

```bash
go mod init go-reading-log-api
go get github.com/lib/pq
```

### 9.2 Build

```bash
go build -o bin/server cmd/server.go
```

### 9.3 Run

```bash
# From rails-app directory
./bin/server

# Or run directly
go run cmd/server.go
```

### 9.4 Environment Setup

```bash
# Copy .env.example if it exists
cp .env.example .env

# Start server
go run cmd/server.go
```

---

## 10. Success Criteria

Phase 1 is complete when:
1. All 3 endpoints return same JSON as Rails app
2. Tests pass with >90% coverage
3. Performance comparable to Rails (or better)
4. Documentation in `README.md` updated

---

## 11. Future Phases

### Phase 2: Dashboard API
- Migrate dashboard endpoints
- Add complex analytics logic

### Phase 3: Users & Authentication
- Add users endpoint
- Implement auth if needed

### Phase 4: Background Jobs
- Integrate with message queue (NATS/NSQ)
- Add background processing

### Phase 5: Watson Tracking
- Add watsons CRUD endpoints

---

## 12. Notes

- **No migration tool**: Manual schema management for Phase 1
- **No auth**: API is open (matches Rails app if unauthenticated)
- **Single DB connection**: No connection pooling initially
- **Synchronous only**: No background jobs
