---
id: RDL-038
title: Create a report of all routes of rails app
status: To Do
assignee: []
created_date: '2026-04-08 12:26'
updated_date: '2026-04-08 12:29'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Execute rails route on container rails-api and make a report of all routes
Save in docs/rails_routes.md
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The task requires generating a report of all Rails API routes by executing `rails routes` command and saving the output to `docs/rails_routes.md`.

- **Method**: Execute `rails routes` command inside the running Rails container (`reading-log-rails-api`)
- **Output format**: Markdown format with structured presentation of routes
- **Storage**: Save to `docs/rails_routes.md` file

The approach will:
1. Execute `rails routes` in the rails-api container
2. Parse and format the output into a readable Markdown document
3. Save to docs/rails_routes.md

### 2. Files to Modify

- **Create**: `docs/rails_routes.md` - New file containing the route report in Markdown format

### 3. Dependencies

- Docker must be installed and running
- Rails container (`reading-log-rails-api`) must be started
- PostgreSQL database should be available for Rails to operate properly
- Rails environment must be properly configured (database connection)

**Prerequisites**:
- Run `docker-compose up rails-api -d` to start the Rails container
- Ensure Rails can access the database

### 4. Code Patterns

- **Markdown format**: Use proper Markdown headers, tables, and code blocks
- **Route structure**: Include HTTP verb, path, controller#action, and any relevant notes
- **Organization**: Group routes by namespace/module (v1, dashboard, echart)
- **Format**: Use tables for clean presentation of route information
- **Code blocks**: Use Ruby syntax highlighting for controller references

### 5. Testing Strategy

**Verification steps**:
1. Verify Docker containers are running: `docker-compose ps`
2. Verify Rails container can execute commands: `docker exec reading-log-rails-api rails routes --help`
3. Execute `rails routes` and capture output
4. Verify output file exists and contains route information
5. Validate Markdown format is correct

**Test approach**:
- Manual verification via Docker commands
- Check file exists and has proper content
- Verify all expected routes from `config/routes.rb` appear in the report

### 6. Risks and Considerations

- **Container availability**: If Rails container isn't running, the command will fail
- **Database connection**: Rails routes command requires database connectivity
- **Output formatting**: Raw `rails routes` output needs to be converted to Markdown format
- **Route coverage**: Ensure all routes from `config/routes.rb` are captured
- **File location**: Confirm `docs/` directory exists before saving

**Potential issues**:
- If postgres container is not healthy, Rails may fail to connect
- Rails container might need to be rebuilt if dependencies changed
- Output might include environment-specific information that should be filtered
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
