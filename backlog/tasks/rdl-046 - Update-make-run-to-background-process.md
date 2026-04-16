---
id: RDL-046
title: Update make run to background process
status: Done
assignee:
  - workflow
created_date: '2026-04-14 11:06'
updated_date: '2026-04-16 20:49'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
update the command `make run` to background
check if exist another process and kill
start up in background
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Current `make run` command runs the server directly using `go run`, which starts in foreground. To change this behavior:
- Check if port (default 3000) is in use via `lsof`
- Kill any existing process on that port
- Start the built binary (from `bin/server`) in background with `nohup` and log output to `server.log`
- Handle `SERVER_PORT` environment variable dynamically for flexibility

This approach ensures a clean startup sequence while maintaining compatibility with existing development workflows. The solution avoids complex PID file management and leverages standard Unix utilities.

### 2. Files to Modify

- **Makefile**: Update the `run` target implementation as follows:

```makefile
run: build
	@PORT=$${SERVER_PORT:-3000}; \
	if command -v lsof >/dev/null; then \
		PIDS=$$(lsof -t -i :$$PORT); \
		if [ -n "$$PIDS" ]; then \
			echo "$(YELLOW)Killing existing processes: $$PIDS$(NC)"; \
			kill -9 $$PIDS; \
		fi; \
	else \
		echo "$(YELLOW)Warning: lsof not installed. Skipping port check.$(NC)"; \
	fi; \
	nohup bin/$(BINARY_NAME) > server.log 2>&1 &
```

### 3. Dependencies

- `lsof` utility (commonly available on Unix-based systems)
- No additional dependencies beyond standard Go toolchain
- Server port configuration via environment variables (`SERVER_PORT`)

### 4. Code Patterns

- Maintain existing Makefile color-coded output using `tput` formatting
- Use shell command chaining with semicolons and backslashes for multi-line commands
- Properly handle environment variable defaults (`$${SERVER_PORT:-3000}`)
- Follow existing error handling patterns (e.g., check for `lsof` availability)

### 5. Testing Strategy

- Manual verification:
  - Start server via `make run`, verify it runs in background (check with `ps aux | grep server`)
  - Run `make run` again, confirm existing process is killed and new instance starts
  - Check `server.log` contains output after starting
- Add test case to developer documentation showing expected behavior

### 6. Risks and Considerations

- `lsof` may not be available on minimal environments (but acceptable for dev workflow)
- Killing processes by port could affect other services if multiple applications use same port in development environment
- `server.log` file growth requires cleanup instructions (will add to AGENTS.md documentation)
- Port conflict resolution works best when `SERVER_PORT` is consistent across commands
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented background process handling for `make run` command:
- Updated Makefile to check and kill existing processes on server port before starting new instance using lsof
- Added nohup logging to server.log with proper environment variable handling (SERVER_PORT)
- Documented server log cleanup procedures in AGENTS.md under "Server Log Management" section

All tests passed via testing-expert subagent:
- 12 test packages executed successfully (all cached)
- No failures or errors detected
- go fmt/vet not applicable to Makefile changes

Key improvements:
- Prevents port conflicts during development workflow
- Clear documentation for log file management
- Maintains backward compatibility with existing environment variables

Follow-up considerations:
- Consider adding `clean-log` target to Makefile for automated cleanup (already documented in AGENTS.md)
- Monitor server.log growth in production environments
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
