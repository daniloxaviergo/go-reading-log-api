---
id: RDL-015
title: Add make command to reload database
status: To Do
assignee:
  - thomas
created_date: '2026-04-03 09:36'
updated_date: '2026-04-03 09:48'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create a command make reload when drop database of docker-compose and up the doc/database.sql. Add this file on pg container and use pg_restore inside container.
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The `make reload` command will drop and recreate the database using the existing `docs/database.sql` file. The implementation will:

1. **Add a new Makefile target** `reload` that:
   - Displays a WARNING prompt requiring user confirmation (to prevent accidental data loss)
   - Stops the Docker Compose services using `docker-compose down`
   - Removes the PostgreSQL data volume using `docker-compose down -v`
   - Starts the services again using `docker-compose up -d`
   - Waits for PostgreSQL to be ready using `pg_isready`
   - Restores the database from `docs/database.sql` using `psql` inside the PostgreSQL container
   - Verifies the restoration by checking table existence

2. **Database restoration method**: The `docs/database.sql` is a plain PostgreSQL dump (not custom format), so we'll use `psql` instead of `pg_restore`. We'll use `docker exec` to run `psql` inside the PostgreSQL container.

3. **Integration with existing Docker setup**: Use `docker-compose` commands to manage services consistently with existing commands.

### 2. Files to Modify

- **Makefile**: Add the `reload` target, confirmation prompt, and helper functions
- **docs/database.sql**: No changes needed (already exists with database schema and data)
- **docker-compose.yml**: No changes needed (already configured properly)

### 3. Dependencies

- Docker and Docker Compose must be installed and running
- Existing `docs/database.sql` file must exist
- Environment variables must be configured in `.env` (or use defaults from `.env.example`)

### 4. Code Patterns

Follow existing Makefile patterns:
- UseColors for output formatting (GREEN, RED, YELLOW, BLUE)
- Check for Docker installation with helpful error messages
- Print progress messages with appropriate colors
- Use `-v` flag for volume removal in `docker-compose down`
- Follow the naming convention: `reload` (simple and clear)

### 5. Testing Strategy

1. **Makefile syntax validation**: Verify the Makefile is valid with `make help`
2. **Integration test**: Run `make reload` in a test environment with dummy data
3. **Verify database restoration**: Check that tables and data are restored correctly by querying the database
4. **Test error handling**: Verify proper error messages when:
   - Docker is not available
   - `.env` file is missing or has invalid values
   - `docs/database.sql` doesn't exist

### 6. Risks and Considerations

**Data Loss**: The command will permanently delete all database data. A warning message with confirmation prompt will be displayed before proceeding.

**Dependencies**: 
- Requires Docker and Docker Compose to be installed
- Requires the `docs/database.sql` file to exist
- Requires proper environment variables in `.env`

**Error handling**:
- Docker not installed or not running
- Database restoration fails (SQL syntax errors)
- Container startup fails or health check times out
- Environment variable configuration issues

**User experience**:
- Clear warning about data loss before proceeding
- User confirmation required before executing
- Helpful success messages
- Troubleshooting suggestions for common failures

**Implementation details**:
- Use `docker-compose down -v` to remove volumes (clean slate)
- Wait for PostgreSQL to be ready using `pg_isready` in a loop
- Use `docker exec` to run `psql` inside the PostgreSQL container
- Pass environment variables from `.env` to the container via `docker exec -e`
- Validate `docs/database.sql` exists before attempting restoration
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
[2026-04-03] Task overview:
- Need to add `make reload` command that drops and recreates database
- Using docker-compose for service management
- Database dump is in `docs/database.sql` (plain PostgreSQL format, not custom dump format)
- Using psql (not pg_restore) since the dump is plain SQL
- Command needs to:
  1. Warn user about data loss
  2. Require confirmation
  3. Stop services
  4. Remove volumes (docker-compose down -v)
  5. Start services again
  6. Wait for PostgreSQL ready
  7. Restore database from docs/database.sql using psql in container
  8. Verify restoration

[2026-04-03] Implementation completed:
- Added `docker-reload` target to Makefile
- Added `reload` alias target for convenience
- Confirmation prompt with data loss warning
- Docker availability check
- docs/database.sql existence check
- docker-compose down to stop services
- docker-compose down -v to remove volumes
- docker-compose up -d to start services
- PostgreSQL readiness check loop (30 attempts, 2s interval)
- Database restoration using psql
- Database verification query
- Success messages and next steps guidance
- All tests passed: syntax validation, dry-run, help output
<!-- SECTION:NOTES:END -->

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
