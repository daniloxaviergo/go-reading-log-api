---
id: RDL-014
title: Add docker-compose
status: To Do
assignee:
  - thomas
created_date: '2026-04-01 19:36'
updated_date: '2026-04-01 19:51'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
add docker-compose to up @rails-app
changes envs to two applications connect in the same pg
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task will create a `docker-compose.yml` file to containerize both the Go API and the Rails application, allowing them to share the same PostgreSQL database instance. The approach will be:

- Create a docker-compose.yml file that defines three services: `postgres`, `go-api`, and `rails-api`
- The Go API will use the existing Dockerfile pattern from cmd/server.go
- The Rails API will use its existing Dockerfile
- Configure both applications to connect to the shared PostgreSQL container using environment variables
- Use Docker networks for service discovery and isolation
- Create a shared `.env` configuration file for consistent database connection settings

### 2. Files to Modify

| Action | File | Reason |
|--------|------|--------|
| **Create** | `docker-compose.yml` | Main orchestration file for all services |
| **Modify** | `.env` | Update DB_HOST to use container service name (`postgres`) |
| **Modify** | `.env.example` | Add example configuration for docker-compose environment |
| **Create** | `.dockerignore` (optional) | Exclude unnecessary files from Docker builds |

### 3. Dependencies

- **Prerequisites:**
  - Docker and Docker Compose must be installed on the development machine
  - Existing Dockerfiles for both applications (Go and Rails)
  - PostgreSQL 15 container image (already used in Makefile)

- **Configuration Requirements:**
  - Both applications must use environment variables for database configuration
  - PostgreSQL must accept connections from container networks (sslmode=disable or cert-based)
  - Database already exists (`reading_log`) - no migrations needed for Phase 1

- **Blocking Issues:**
  - None identified

### 4. Code Patterns

**Docker Compose Service Structure:**
```yaml
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASS}
      POSTGRES_DB: ${DB_DATABASE}
    ports:
      - "5432:5432"
  
  go-api:
    build:
      context: .
      dockerfile: Dockerfile  # Will need to create
    environment:
      DB_HOST: postgres  # Service name for DNS resolution
      DB_PORT: 5432
    depends_on:
      - postgres
  
  rails-api:
    build:
      context: ./rails-app
      dockerfile: Dockerfile
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
    depends_on:
      - postgres
```

**Environment Variable Consistency:**
- Use same environment variable names in both applications (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_DATABASE`)
- Keep database configuration in `.env` file at project root
- Application-specific settings (SERVER_PORT, PORT) should use appropriate defaults

### 5. Testing Strategy

**Unit Testing:**
- Verify docker-compose.yml syntax with `docker-compose config`
- Test service startup order with `docker-compose up`
- Verify both applications can connect to PostgreSQL
- Confirm no port conflicts between services

**Integration Testing:**
- Start all services: `docker-compose up -d`
- Test Go API health endpoint: `curl http://localhost:3000/healthz`
- Test Rails API endpoints: `curl http://localhost:3000/api/v1/projects` (or Rails port)
- Verify database connectivity: `docker exec reading-log-db psql -U postgres -d reading_log`
- Test cross-container communication via Docker network

**Database Verification:**
```bash
# Connect to PostgreSQL container
docker exec -it reading-log-db psql -U postgres -d reading_log

# Verify tables exist
\d projects
\d logs
```

### 6. Risks and Considerations

**Potential Issues:**
1. **Port Conflicts:** Both Go and Rails APIs default to port 3000 - must configure different ports
   - *Mitigation:* Use `SERVER_PORT=3000` for Go, `PORT=3001` for Rails

2. **Database Schema Mismatch:** If Rails app expects different schema than Go app
   - *Mitigation:* Verify both applications work with existing schema before deployment

3. **SSL/TLS Configuration:** Go app uses `sslmode=disable`, Rails may need same
   - *Mitigation:* Add `sslmode=disable` to connection strings if needed

4. **Container Network Isolation:** Services may not discover each other if not on same network
   - *Mitigation:* Explicitly define shared network in docker-compose.yml

**Deployment Considerations:**
- This setup is for development/local testing only
- Production deployment would require:
  - SSL/TLS certificates for PostgreSQL
  - secrets management for database credentials
  - Health checks for production readiness
  - Volume mounts for data persistence

**Next Steps (Post-RDL-014):**
- Consider adding docker-compose for development-only services (Redis for Sidekiq)
- Document usage in QWEN.md and README.md
- Add make targets: `make docker-up`, `make docker-down`, `make docker-logs`
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
- Created docker-compose.yml with postgres, go-api, and rails-api services

- Created .dockerignore files for both applications

- Updated .env.example with Docker-compatible configuration

- Updated .env for Docker environment (DB_HOST=postgres)

- Created Dockerfile for Go application (multi-stage build)

- Updated rails-app Dockerfile to use environment variables for DB and PORT

- Added docker-compose targets to Makefile: docker-up, docker-down, docker-logs, docker-ps, docker-stop-pg

- All 89 tests passed

- go fmt and go vet passed with no errors

- Port conflict mitigation: Go API uses 3000, Rails API uses 3001
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Implementation Summary

### What Changed
This task adds Docker Compose support to the Go Reading Log API project, enabling containerized deployment of both the Go API and Rails API sharing the same PostgreSQL database.

### Files Created/Modified

**New Files:**
- `docker-compose.yml` - Orchestrates three services: postgres, go-api, and rails-api
- `Dockerfile` - Multi-stage build for Go API application
- `.dockerignore` - Excludes unnecessary files from Go app Docker build
- `rails-app/.dockerignore` - Excludes unnecessary files from Rails app Docker build

**Modified Files:**
- `.env` - Updated with Docker-compatible DB_HOST=postgres
- `.env.example` - Added Docker configuration examples
- `rails-app/Dockerfile` - Added environment variables for DB and PORT
- `Makefile` - Added docker-compose targets (docker-up, docker-down, docker-logs, docker-ps, docker-stop-pg)
- `QWEN.md` - Added Docker Compose documentation section

### Key Features
- PostgreSQL 15 container with persistent volume
- Go API on port 3000 with health checks and graceful shutdown
- Rails API on port 3001 (port conflict resolved)
- Service discovery via Docker network (postgres service name)
- Health checks for PostgreSQL container
- Multi-stage Docker builds for optimized images

### Testing
- All 89 tests passed (unit and integration)
- `go fmt` and `go vet` pass with no errors
- Build successful: `make build`

### Benefits
- Simplified local development with single command: `make docker-up`
- Consistent environment across development and CI
- Isolated database with persistent volume
- Easy service management with Makefile commands
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [ ] #2 All integration tests pass use testing-expert subagent for test execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Database queries optimized with proper indexes
- [ ] #8 Documentation updated in QWEN.md
- [ ] #9 New code paths include error path tests
- [ ] #10 HTTP handlers test both success and error responses
- [ ] #11 Integration tests verify actual database interactions
- [ ] #12 Tests use testing-expert subagent for test execution and verification
<!-- DOD:END -->
