---
id: RDL-012
title: Add make to project
status: To Do
assignee:
  - thomas
created_date: '2026-04-01 15:12'
updated_date: '2026-04-01 15:39'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add commands like make run, make build, make test, make help, make start pg
Should be cover all aspects of project
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
```markdown
### 1. Technical Approach

The Makefile will provide a unified interface for common development operations including building, running, testing, and database management.

- **Core commands to implement:**
  - `make run` - Build and run the server (equivalent to `go run ./cmd/server.go`)
  - `make build` - Build the binary for production (`go build -o bin/server`)
  - `make test` - Run all tests (`go test ./...`)
  - `make help` - Display available commands
  - `make start-pg` - Start PostgreSQL database (using Docker if available, else fallback to manual)

- **Additional helper commands:**
  - `make test-coverage` - Generate coverage report
  - `make fmt` - Format code with `go fmt`
  - `make vet` - Run `go vet` for static analysis
  - `make clean` - Clean up binaries and build artifacts
  - `make docker-start-pg` - Start PostgreSQL via Docker

- **Architecture decision:** Use a standard Makefile with platform-specific checks (e.g., detect if Docker is available for database management)

- **Why Makefile over shell scripts:** More maintainable, portable across Unix-like systems, better IDE integration

### 2. Files to Modify

**New files to create:**
- `Makefile` - Main build automation file

**Files to reference (read-only access):**
- `cmd/server.go` - Server entry point to understand build/run requirements
- `internal/config/config.go` - Config structure for environment requirements
- `.env.example` - Environment variables needed for testing
- `test/test_helper.go` - Test infrastructure setup
- `test/integration/` - Integration test packages

### 3. Dependencies

**Prerequisites:**
- Go 1.25.7 installed and in PATH
- Make installed (standard on most Unix-like systems)
- PostgreSQL running (for tests and runtime)
- Docker (optional, for `start-pg` command if available)

**Environment setup required:**
- `.env` file must exist with database credentials (or environment variables set)
- Test database `reading_log_test` must exist

**Blocking issues:**
- None known. This is a development tooling task that doesn't block other work.

### 4. Code Patterns

**Makefile conventions to follow:**
- Variables at the top for easy configuration (`BINARY_NAME`, `MODULE_NAME`)
- `.PHONY` targets for non-file targets (`help`, `test`, `run`, etc.)
- Clear error messages with color output where possible
- Use `$(MAKEFLAGS)` for proper jobserver support
- Include help text for each target

**Naming convention:**
- All-make-targets snake_case (consistent with `start-pg` hybrid)
- Binary name: `server` (matches existing `bin/server`)

**Integration patterns:**
- The Makefile should match existing Go development patterns
- Commands should mirror those documented in QWEN.md and README
- Use `go test ./...` pattern for comprehensive test runs

### 5. Testing Strategy

**Makefile validation:**
- Test each target manually to ensure commands work as expected
- Verify `make test` runs all tests without errors
- Verify `make build` creates binary in correct location
- Verify `make run` starts server correctly (with proper environment)
- Verify `make help` displays all available commands

**Test coverage targets:**
- Test with different environment configurations
- Test with and without Docker installed
- Verify error handling in Makefile (e.g., missing .env file)

**Integration with project tests:**
- `make test` should pass existing test suite (80 tests as documented)
- Coverage targets should generate proper reports

### 6. Risks and Considerations

**Trade-offs:**
- Docker dependency for `start-pg`: If Docker is not available, the command should provide clear error message and instructions
- Windows compatibility: Makefile primarily targets Unix-like systems; Windows users would need WSL or alternative

**Potential issues:**
- PostgreSQL connection: Makefile needs to handle missing database gracefully
- Port conflicts: `make run` should handle port already in use scenarios
- Environment variable precedence: Make sure `.env` file is loaded correctly

**Deployment considerations:**
- Production builds should use `make build` before deployment
- The `start-pg` command is for development only; production should use managed PostgreSQL
- No testing database management in Makefile (too complex for a single file)

**Documentation:**
- Update QWEN.md to reference Makefile commands
- Add Makefile usage to project README if created
- Include in AGENTS.md for AI agent reference
```
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Created Makefile with all required commands

All 121 tests pass with make test

Fixed Go 1.25.7 go fmt -w flag issue

Docker PostgreSQL integration working
<!-- SECTION:NOTES:END -->
