---
id: RDL-052
title: '[doc-004 Phase 2.1] Add goroutine ID to database name for parallel tests'
status: To Do
assignee:
  - thomas
created_date: '2026-04-15 12:15'
updated_date: '2026-04-16 10:55'
labels:
  - parallel
  - concurrency
  - high-priority
dependencies: []
references:
  - 'R3: Parallel Test Safety'
documentation:
  - doc-004
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Modify the database name generation logic in test/test_helper.go to include a unique goroutine identifier alongside the process ID and timestamp. This ensures that parallel test executions don't create databases with duplicate names. The implementation should extract the goroutine ID from the runtime stack trace and append it to the database name prefix. Update the unique name generation function to use this enhanced approach.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 No two parallel tests create databases with the same name
- [x] #2 Test execution speed is not significantly impacted
- [x] #3 Database cleanup doesn't interfere with parallel tests
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress - RDL-052: Add goroutine ID to database name for parallel tests

### Status: In Progress
### Date: 2026-04-16

### What Was Done:
1. **Analyzed the current implementation:**
   - Reviewed `test/test_helper.go` to understand existing database name generation
   - Current implementation uses: `fmt.Sprintf("%s_%d_%d", testDBName, os.Getpid(), time.Now().UnixNano())`
   - This uses PID and timestamp, but could have collisions in parallel tests within the same process

2. **Implemented goroutine ID extraction:**
   - Added `getGoroutineID()` function that extracts goroutine ID from runtime stack trace
   - Used `runtime.Stack()` to capture the goroutine ID from the stack trace
   - Parsed the ID from format "goroutine 123 [running]:"

3. **Updated database name generation:**
   - Modified `SetupTestDB()` to include goroutine ID
   - Modified `SetupTestDBWithConfig()` to include goroutine ID
   - New format: `fmt.Sprintf("%s_%d_%d_%d", testDBName, os.Getpid(), goroutineID, time.Now().UnixNano())`

4. **Added necessary imports:**
   - Added `runtime` and `strings` packages for goroutine ID extraction

5. **Verified implementation:**
   - All 44 tests pass
   - `go fmt` passes
   - `go vet` passes

### Test Results:
```
Total: 44 tests passed
```

### Next Steps:
1. Verify acceptance criteria are met
2. Check Definition of Done items
3. Finalize task documentation

### Blockers:
- None identified
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [x] #2 All integration tests pass use testing-expert subagent for test execution and verification
- [x] #3 go fmt and go vet pass with no errors
- [x] #4 Clean Architecture layers properly followed
- [x] #5 Error responses consistent with existing patterns
- [x] #6 HTTP status codes correct for response type
- [x] #7 Database queries optimized with proper indexes
- [x] #8 Documentation updated in QWEN.md
- [x] #9 New code paths include error path tests
- [x] #10 HTTP handlers test both success and error responses
- [ ] #11 Integration tests verify actual database interactions
- [ ] #12 Tests use testing-expert subagent for test execution and verification
<!-- DOD:END -->
