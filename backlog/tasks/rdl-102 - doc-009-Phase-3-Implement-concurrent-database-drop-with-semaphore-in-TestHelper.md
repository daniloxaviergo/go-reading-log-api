---
id: RDL-102
title: >-
  [doc-009 Phase 3] Implement concurrent database drop with semaphore in
  TestHelper
status: To Do
assignee:
  - thomas
created_date: '2026-04-24 13:42'
updated_date: '2026-04-24 18:11'
labels:
  - bug
  - test-fix
  - p2-high
dependencies: []
references:
  - Decision 3
  - test/test_helper.go
documentation:
  - doc-009
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement concurrent database cleanup in test/test_helper.go using goroutines and semaphores to prevent deadlocks during sequential database drops. Add health checks and proper error collection to ensure visibility of cleanup failures.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Concurrent database drops complete without deadlocks
- [ ] #2 Maximum 5 concurrent drop operations enforced via semaphore
- [ ] #3 All orphaned test databases are properly cleaned up
- [ ] #4 Error collection provides visibility into cleanup failures
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task requires converting the sequential database cleanup in `TestHelper.Close()` to a concurrent implementation using goroutines and a semaphore pattern to prevent overwhelming the PostgreSQL server with too many simultaneous DROP DATABASE operations.

**Architecture Decision:**
- Use Go's `sync.WaitGroup` to wait for all concurrent drop operations to complete
- Use a buffered channel as a semaphore (capacity 5) to limit concurrent database drops
- Collect errors from all goroutines and return them aggregated using `errors.Join()`
- Each DROP DATABASE operation runs in its own goroutine with its own context timeout

**Why This Approach:**
1. **Performance**: Sequential drops can take significant time when many orphaned databases exist; concurrency reduces total cleanup time
2. **Resource Protection**: Semaphore prevents connection pool exhaustion and database server overload
3. **Reliability**: Error collection ensures visibility into which specific drops failed without failing the entire cleanup
4. **Backward Compatibility**: The Close() method signature remains unchanged (void return), maintaining existing test code compatibility

**Implementation Strategy:**
```
┌─────────────────────────────────────────────────────────────┐
│              Concurrent Database Drop Architecture           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                    Main Thread                        │   │
│  │  - Collect databases to drop                         │   │
│  │  - Launch goroutines (limited by semaphore)          │   │
│  │  - Wait for all completions                          │   │
│  └──────────────────────┬───────────────────────────────┘   │
│                         │                                    │
│         ┌───────────────┼───────────────┐                   │
│         ▼               ▼               ▼                   │
│    ┌──────────┐   ┌──────────┐   ┌──────────┐              │
│    │  Goroutine 1 │   │  Goroutine 2 │   │  Goroutine N │              │
│    │  (DROP DB) │   │  (DROP DB) │   │  (DROP DB) │              │
│    └─────┬──────┘   └─────┬──────┘   └─────┬──────┘              │
│          │                │                │                  │
│          ▼                ▼                ▼                  │
│    ┌──────────────────────────────────────────────┐           │
│    │              Semaphore (max 5)               │           │
│    │  ┌─────┬─────┬─────┬─────┬─────┐            │           │
│    │  │ D1  │ D2  │ D3  │ D4  │ D5  │ ← Active  │           │
│    │  └─────┴─────┴─────┴─────┴─────┘            │           │
│    │  Available slots: ████░░░░░░░░░░░░          │           │
│    └──────────────────────────────────────────────┘           │
│                         │                                    │
│         ┌───────────────┼───────────────┐                   │
│         ▼               ▼               ▼                   │
│    ┌──────────┐   ┌──────────┐   ┌──────────┐              │
│    │  Error 1  │   │  Error 2  │   │  Error N  │              │
│    └─────┬──────┘   └─────┬──────┘   └─────┬──────┘              │
│          │                │                │                  │
│          └────────────────┼────────────────┘                   │
│                           ▼                                   │
│                    ┌────────────┐                            │
│                    │ errors.Join()│                           │
│                    └────────────┘                            │
└─────────────────────────────────────────────────────────────┘
```

---

### 2. Files to Modify

| File | Action | Description |
|------|--------|-------------|
| `test/test_helper.go` | **MODIFY** | Rewrite `TestHelper.Close()` method with concurrent database drops using semaphore pattern |
| `test/test_helper_test.go` | **NO CHANGE NEEDED** | Existing tests should pass with new implementation (same external behavior) |

**Files NOT to Modify:**
- `internal/domain/` - No domain changes required
- `internal/api/` - No API changes required  
- `cmd/server.go` - No server changes required

---

### 3. Dependencies

**Prerequisites:**
- ✅ Go standard library (`sync`, `errors`, `context`) - Already available
- ✅ `pgxpool` - Already imported and used in test_helper.go
- ✅ `test/test_helper.go` - Current implementation exists and is well-documented

**No External Dependencies Required:**
This implementation uses only Go standard library packages. No additional go.mod dependencies needed.

**Blocking Issues:** None identified. This is a self-contained refactoring within the test infrastructure.

---

### 4. Code Patterns

#### Pattern 1: Semaphore with Buffered Channel
```go
const maxConcurrentDrops = 5

// Semaphore that allows max 5 concurrent operations
sem := make(chan struct{}, maxConcurrentDrops)

// Acquire semaphore (blocks if at capacity)
sem <- struct{}{}

// Release semaphore when done
<-sem
```

#### Pattern 2: Error Collection with WaitGroup
```go
var wg sync.WaitGroup
var errs []error
var mu sync.Mutex // Mutex for thread-safe error append

for _, db := range databasesToDrop {
    wg.Add(1)
    go func(dbName string) {
        defer wg.Done()
        
        // Acquire semaphore
        sem <- struct{}{}
        defer func() { <-sem }()
        
        // Perform operation
        if err := dropDatabase(dbName); err != nil {
            mu.Lock()
            errs = append(errs, err)
            mu.Unlock()
        }
    }(db)
}

wg.Wait()

// Aggregate errors
if len(errs) > 0 {
    return errors.Join(errs...)
}
return nil
```

#### Pattern 3: Context with Timeout per Operation
```go
// Each DROP DATABASE gets its own context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

_, err := pool.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
```

#### Pattern 4: Safe Error Logging Without Blocking
```go
// Log errors but don't fail the test - cleanup should not block test results
if err != nil {
    log.Printf("[Cleanup] Database %s drop failed: %v", dbName, err)
    // Don't return error to avoid blocking test results
}
```

---

### 5. Testing Strategy

#### Unit Tests (Existing - No Changes Required)
The existing `test_helper_test.go` tests should continue to work without modification because:
- The `Close()` method signature remains identical (no return value)
- The external behavior is preserved (cleanup still happens, just concurrently)
- Error handling is improved (errors are logged, not silently ignored)

**Test Coverage to Verify:**
```bash
# Run TestHelper tests to verify no regression
go test -v ./test/... -run "TestHelper"

# Expected: All existing tests pass without modification
```

#### Integration Tests (Existing - No Changes Required)
All integration tests that use `defer helper.Close()` should continue working:
- `dashboard_integration_test.go`
- `projects_integration_test.go`
- `rails_comparison_test.go`
- `expected_values_integration_test.go`
- `error_scenarios_test.go`

**Verification Command:**
```bash
# Run all integration tests
go test -v ./test/integration/... -timeout=60s

# Expected: All integration tests pass with improved cleanup reliability
```

#### Performance Tests (New - Recommended)
Consider adding a benchmark to measure concurrent vs sequential performance:

```go
func BenchmarkTestHelperCleanup(b *testing.B) {
    // Setup multiple test databases
    helpers := make([]*TestHelper, b.N)
    for i := 0; i < b.N; i++ {
        h, err := SetupTestDB()
        if err != nil {
            b.Fatal(err)
        }
        helpers[i] = h
    }
    
    b.ResetTimer()
    
    // Measure cleanup time
    for _, h := range helpers {
        h.Close()
    }
}
```

---

### 6. Risks and Considerations

#### Risk 1: Connection Pool Exhaustion
**Description:** Creating multiple connection pools concurrently could exhaust available connections.

**Mitigation:**
- Limit concurrent drops to 5 via semaphore
- Reuse existing pool where possible (already implemented)
- Each goroutine creates its own short-lived pool for DROP operations

#### Risk 2: Test Flakiness from Concurrent Operations
**Description:** Race conditions in cleanup could cause intermittent test failures.

**Mitigation:**
- WaitGroup ensures all goroutines complete before Close() returns
- Mutex protects error collection from race conditions
- Each operation uses its own context (no shared state)

#### Risk 3: Error Visibility Loss
**Description:** Current implementation silently ignores errors; concurrent version could lose error information.

**Mitigation:**
- Collect ALL errors using `errors.Join()`
- Log each individual failure with database name
- Return aggregated error for visibility (while not blocking test results)

#### Risk 4: Backward Compatibility
**Description:** Some existing code might depend on Close() behavior.

**Mitigation:**
- Keep same method signature (void return)
- Same timing characteristics (still uses defer)
- Only change internal implementation details

---

### 7. Implementation Steps

**Step 1: Add imports (if not already present)**
```go
import (
    "sync"
    "errors"
    // ... existing imports
)
```

**Step 2: Define semaphore constant**
```go
const maxConcurrentDrops = 5
```

**Step 3: Rewrite `TestHelper.Close()` method**
- Replace sequential cleanup with concurrent implementation
- Use WaitGroup to track goroutines
- Use buffered channel as semaphore
- Collect and aggregate errors
- Maintain same logging behavior

**Step 4: Verify no test modifications needed**
- Run existing tests without changing them
- Confirm all tests pass

**Step 5: Run full test suite**
```bash
go test -v ./test/...
```

---

### 8. Acceptance Criteria Alignment

| AC | Requirement | Implementation Addressed |
|----|-------------|-------------------------|
| #1 | Concurrent database drops complete without deadlocks | ✅ WaitGroup + semaphore prevents deadlock |
| #2 | Maximum 5 concurrent drop operations enforced via semaphore | ✅ `sem := make(chan struct{}, 5)` |
| #3 | All orphaned test databases are properly cleaned up | ✅ Same logic, just concurrent execution |
| #4 | Error collection provides visibility into cleanup failures | ✅ Errors collected via mutex, logged individually |

---

### Summary

**Total Files Modified:** 1 (`test/test_helper.go`)

**Lines of Code Changed:** ~60 lines (main Close() method and helper function)

**Complexity:** Medium - involves concurrent programming patterns but uses standard library only

**Risk Level:** Low - well-isolated change, backward compatible, existing tests validate behavior

**Estimated Implementation Time:** 30-45 minutes
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass
- [ ] #2 All integration tests pass execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md
- [ ] #8 New code paths include error path tests
- [ ] #9 HTTP handlers test both success and error responses
- [ ] #10 Integration tests verify actual database interactions
<!-- DOD:END -->
