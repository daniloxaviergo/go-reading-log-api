---
id: RDL-048
title: '[doc-004 Phase 1.1] Implement defer cleanup in TestHelper.Close()'
status: To Do
assignee:
  - thomas
created_date: '2026-04-15 12:14'
updated_date: '2026-04-15 12:28'
labels:
  - cleanup
  - infrastructure
  - critical
dependencies: []
references:
  - 'R1: Auto-Cleanup on Test Completion'
documentation:
  - doc-004
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update the TestHelper.Close() method in test/test_helper.go to automatically drop the test database when tests complete. The cleanup must use defer to ensure it runs even on panic, within 1 second of test completion. Implement proper error handling that doesn't throw errors if the database doesn't exist and ensures cleanup doesn't block test results. The implementation should create a separate connection pool to the main database to execute the DROP DATABASE command.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Test database is dropped within 1 second of test completion
- [ ] #2 Cleanup occurs even if test panics
- [ ] #3 No error is thrown if database doesn't exist
- [ ] #4 Cleanup doesn't block test results
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The task requires modifying the `TestHelper.Close()` method to use `defer` for automatic cleanup. The current implementation requires manual cleanup. The plan is to:

1. Add a `defer` statement that calls the cleanup logic when the function exits
2. Ensure the cleanup runs within 1 second of test completion
3. Handle cases where the database doesn't exist (no error thrown)
4. Use a separate connection pool to execute the DROP DATABASE command
5. Ensure cleanup doesn't block test results

### 2. Files to Modify

| File | Action | Description |
|------|--------|-------------|
| `test/test_helper.go` | Modify | Update `TestHelper.Close()` method to use defer for cleanup |
| `test/test_helper_test.go` | Create/Modify | Add tests to verify defer cleanup works correctly |

### 3. Dependencies

- No external dependencies required
- Uses existing `pgxpool` library
- No prerequisite tasks - this is a self-contained improvement
- No other code depends on the current `Close()` method behavior

### 4. Code Patterns

- **Defer pattern**: Use `defer` to ensure cleanup runs on function exit, including panics
- **Separate connection pool**: Create a new connection pool for DROP DATABASE to avoid issues with closing the pool being dropped
- **Error suppression**: Log errors but don't propagate them to avoid blocking test results
- **Timeout management**: Use 1-second timeout for cleanup operations
- **Safe database drop**: Use `DROP DATABASE IF EXISTS` to handle missing databases gracefully

### 5. Testing Strategy

- Write unit tests that verify `Close()` schedules cleanup
- Test that cleanup runs even when test panics
- Test that cleanup doesn't fail if database doesn't exist
- Test that cleanup completes within 1 second
- Run existing test suite to ensure no regressions

### 6. Risks and Considerations

- **Timing**: Must ensure cleanup happens within 1 second - need to verify timeout is set correctly
- **Parallel tests**: Must ensure unique database names don't cause conflicts during cleanup
- **Error handling**: Must suppress errors during cleanup to not block test results
- **Connection pool**: Must use separate pool for DROP DATABASE to avoid connection issues
- **Graceful degradation**: If cleanup fails, tests should still report results
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress - RDL-048

### Current Status
- Task: Implement defer cleanup in TestHelper.Close()
- Status: In Progress
- Date: 2026-04-15

### Analysis
The current `TestHelper.Close()` method (lines 260-279 in test_helper.go) has manual cleanup logic that:
1. Drops the test database if created
2. Closes the connection pool

However, it does NOT use `defer` to ensure cleanup runs on panic or automatic cleanup.

### Required Changes
According to acceptance criteria:
- [ ] Test database is dropped within 1 second of test completion
- [ ] Cleanup occurs even if test panics
- [ ] No error is thrown if database doesn't exist
- [ ] Cleanup doesn't block test results

### Implementation Plan
1. Wrap the cleanup logic in a `defer` statement in `Close()`
2. Ensure cleanup runs within 1 second (current timeout is 30s, needs reduction)
3. Add proper error suppression so cleanup doesn't block test results
4. Use separate connection pool for DROP DATABASE command
5. Handle cases where database doesn't exist gracefully

### Code Changes Needed
- Modify `TestHelper.Close()` method to use `defer` for cleanup
- Ensure timeout is 1 second for cleanup operations
- Suppress errors during cleanup to not block test results
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
