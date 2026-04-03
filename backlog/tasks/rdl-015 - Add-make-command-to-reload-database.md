---
id: RDL-015
title: Add make command to reload database
status: To Do
assignee:
  - catarina
created_date: '2026-04-03 09:36'
updated_date: '2026-04-03 09:39'
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
   - Stops the PostgreSQL container
   - Drops the existing database volume (or executes DROP DATABASE)
   - Starts PostgreSQL fresh
   - Waits for database to be ready
   - Restores the database from `docs/database.sql` using `pg_restore` or `psql`
   - Verifies the restoration

2. **Database SQL format**: The `docs/database.sql` is a plain PostgreSQL dump (not custom format), so we'll use `psql` instead of `pg_restore` for restoration.

3. **Integration with existing Docker setup**: Use the same container name (`reading-log-db`) and environment variables from `.env` file.

### 2. Files to Modify

- **Makefile**: Add the `reload` target and supporting functions
- **docs/database.sql**: No changes needed (already exists with database schema and data)
- **.env**: No changes needed (configuration already exists)

### 3. Dependencies

- Docker must be installed and running
- Existing `docs/database.sql` file must exist
- Environment variables must be configured in `.env` (or use defaults from `.env.example`)

### 4. Code Patterns

Follow existing Makefile patterns:
- UseColors for output formatting (GREEN, RED, YELLOW, BLUE)
- Check for Docker installation
- Print helpful messages for user
- Use the same environment variable pattern (load from .env)
- Follow the naming convention: `reload` (not `db-reload` or `database-reload`)

### 5. Testing Strategy

1. **Unit test the Makefile syntax**: Verify the Makefile is valid
2. **Integration test**: Run `make reload` in a test environment
3. **Verify database restoration**: Check that tables and data are restored correctly
4. **Test error handling**: Verify proper error messages when Docker is not available

### 6. Risks and Considerations

**Data Loss**: The command will permanently delete all database data. A warning message will be displayed before proceeding.

**Dependencies**: 
- Requires Docker to be installed
- Requires the `docs/database.sql` file to exist
- Requires proper environment variables in `.env`

**Error handling**:
- Docker not installed
- Database restoration fails (SQL syntax errors)
- Container startup fails
- Database connection timeout

**User experience**:
- Clear warning about data loss before proceeding
- Helpful success messages
- Troubleshooting suggestions for common failures

**Implementation details**:
- Use `docker-compose down` to stop services first (safer than stopping just the database)
- Use volume removal to ensure clean slate
- Wait for PostgreSQL to be ready using `pg_isready`
- Use environment variables from `.env` for database credentials
- Provide helpful error messages with resolution suggestions
<!-- SECTION:PLAN:END -->

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
