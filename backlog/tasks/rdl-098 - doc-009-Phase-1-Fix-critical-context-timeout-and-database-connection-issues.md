---
id: RDL-098
title: '[doc-009 Phase 1] Fix critical context timeout and database connection issues'
status: To Do
assignee:
  - thomas
created_date: '2026-04-24 13:41'
updated_date: '2026-04-24 14:37'
labels:
  - bug
  - test-fix
  - p1-critical
dependencies: []
references:
  - REQ-01
  - REQ-02
  - REQ-05
documentation:
  - doc-009
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Fix GetTestContext() and GetTestContextWithTimeout() functions to return cancel functions instead of discarding them, preventing resource leaks. Add database availability checks with timeout to integration tests to prevent hangs during test execution.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 All 9 SpeculateService unit tests pass without panics
- [ ] #2 Context timeout tests complete within 5 seconds
- [ ] #3 Integration tests have proper database availability checks
- [ ] #4 No resource leaks detected in test execution
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task addresses critical test infrastructure issues in three areas:

**Issue A: Context Timeout Pattern (P1 Critical)**
The `GetTestContext()` and `GetTestContextWithTimeout()` functions currently discard the cancel function, causing resource leaks. The fix requires changing return types to expose the cancel function for proper cleanup.

**Current broken pattern:**
```go
func GetTestContext() context.Context {
    ctx, cancel := context.WithTimeout(context.Background(), testContextTimeout)
    _ = cancel  // ❌ Resource leak!
    return ctx
}
```

**Fix approach:** Return both context and cancel function, allowing callers to properly clean up resources.

**Issue B: Integration Test Database Availability (P1 Critical)**
Integration tests hang when database is unavailable because there's no connection health check with timeout before attempting operations. The fix adds a ping-based availability check with configurable timeout.

**Fix approach:** Implement `CheckDatabaseAvailability()` function that verifies connection before test execution, failing fast with clear error messages if database is unreachable.

**Issue C: SpeculateService Date Determinism (P2 High)**
The `GetToday()` function uses `time.Now()` directly, making tests non-deterministic. The fix introduces a date abstraction layer allowing test-specific date injection while maintaining backward compatibility.

**Fix approach:** Introduce `GetTodayFunc` variable for dependency injection, with `SetTestDate()` helper for testing.

---

### 2. Files to Modify

| File | Changes | Priority |
|------|---------|----------|
| `test/test_helper.go` | Fix `GetTestContext()`, `GetTestContextWithTimeout()` - return cancel functions; Add `CheckDatabaseAvailability()` with timeout | P1 |
| `internal/service/dashboard/day_service.go` | Add date abstraction layer (`GetTodayFunc`, `SetTestDate`) | P2 |
| `internal/service/dashboard/speculate_service.go` | Update to use abstracted date function | P2 |
| `test/test_helper_test.go` | Update tests for new context return signature; Add database availability check tests | P1 |
| `test/unit/speculate_service_test.go` | Use `SetTestDate()` for deterministic testing; Fix index assertions | P2 |

---

### 3. Dependencies

**Prerequisites:**
- None - this is a foundational fix that enables other work

**Blocking issues:**
- None identified

**Setup steps:**
1. Apply context timeout fixes first (P1)
2. Add database availability checks (P1)  
3. Implement date abstraction layer (P2)
4. Update dependent code and tests

---

### 4. Code Patterns

**Context with Cancel Pattern (Go Best Practice):**
```go
func GetTestContext() (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.Background(), testContextTimeout)
}

// Usage:
ctx, cancel := GetTestContext()
defer cancel()
// ... use ctx
```

**Database Availability Check Pattern:**
```go
func CheckDatabaseAvailability(ctx context.Context, pool *pgxpool.Pool, timeout time.Duration) error {
    pingCtx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()
    
    if err := pool.Ping(pingCtx); err != nil {
        return fmt.Errorf("database not available: %w", err)
    }
    return nil
}
```

**Date Abstraction Pattern:**
```go
// Variable for injection (defaults to real time)
var GetTodayFunc = func() time.Time {
    return truncateToMidnight(time.Now())
}

func GetToday() time.Time {
    return GetTodayFunc()
}

func SetTestDate(date time.Time) {
    GetTodayFunc = func() time.Time {
        return truncateToMidnight(date)
    }
}
```

---

### 5. Testing Strategy

**Unit Tests (P1):**
- `TestGetTestContext_ReturnsCancelFunc` - Verify cancel function is returned
- `TestGetTestContextWithTimeout_CustomTimeout` - Verify custom timeout works
- `TestCheckDatabaseAvailability_Success` - Verify healthy database passes
- `TestCheckDatabaseAvailability_Timeout` - Verify unavailable database fails fast

**Unit Tests (P2):**
- `TestSpeculateService_GetToday_Deterministic` - Verify fixed date returns consistent results
- `TestSpeculateService_SetTestDate_Injection` - Verify date injection works
- Update existing SpeculateService tests to use `SetTestDate()`

**Integration Tests:**
- Add database availability check before all integration tests
- Verify tests fail fast with clear error if database is unreachable
- Ensure no resource leaks after test completion

---

### 6. Risks and Considerations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Breaking existing tests that call `GetTestContext()` without cancel | Medium | Update all callers; add deprecation warning first |
| Database availability check adds latency to test startup | Low | Use short timeout (5s); only run when needed |
| Date abstraction introduces complexity | Low | Keep simple; well-documented; minimal surface area |
| Race conditions in parallel tests with date injection | Medium | Ensure `GetTodayFunc` is goroutine-safe; consider per-test isolation |

**Breaking Changes:**
- `GetTestContext()` and `GetTestContextWithTimeout()` change return type from `context.Context` to `(context.Context, context.CancelFunc)`
- All callers must be updated to capture and call the cancel function

**Rollback Plan:**
- Changes are isolated to test infrastructure
- No production code affected
- Easy to revert if issues discovered

---

### 7. Implementation Checklist

**Phase 1: Context Timeout Fixes (P1)**
- [ ] Modify `GetTestContext()` to return `(context.Context, context.CancelFunc)`
- [ ] Modify `GetTestContextWithTimeout()` to return `(context.Context, context.CancelFunc)`
- [ ] Update all callers in test files
- [ ] Add unit tests for new signature
- [ ] Verify no resource leaks with `go test -v -count=10`

**Phase 2: Database Availability Checks (P1)**
- [ ] Implement `CheckDatabaseAvailability()` function
- [ ] Add timeout parameter with sensible default
- [ ] Integrate into `SetupTestDB()` and `SetupTestDBWithConfig()`
- [ ] Update integration tests to use availability check
- [ ] Verify fast failure on unavailable database

**Phase 3: Date Abstraction (P2)**
- [ ] Add `GetTodayFunc` variable to `day_service.go`
- [ ] Implement `SetTestDate()` helper function
- [ ] Update `SpeculateService` to use abstracted date
- [ ] Update all affected unit tests
- [ ] Verify deterministic test results across multiple runs

**Phase 4: Validation (P3)**
- [ ] Run full test suite with `-race` flag
- [ ] Verify no resource leaks detected
- [ ] Confirm all acceptance criteria met
- [ ] Update documentation in QWEN.md
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
# Implementation Progress - RDL-098

## Status: In Progress

### Phase 1: Context Timeout Fixes (P1) - COMPLETE ✅

**Issue:** `GetTestContext()` and `GetTestContextWithTimeout()` discard cancel functions, causing resource leaks.

**Fix Applied:**
1. Modified `GetTestContext()` to return `(context.Context, context.CancelFunc)`
2. Modified `GetTestContextWithTimeout()` to return `(context.Context, context.CancelFunc)`
3. Updated all callers in test files to capture and call the cancel function

**Files Modified:**
- `test/test_helper.go` - Fixed context functions
- `test/test_helper_test.go` - Updated tests for new signature
- `test/integration/test_context.go` - Updated integration test context
- `test/performance/*.go` - Updated benchmark tests

### Phase 2: Database Availability Checks (P1) - COMPLETE ✅

**Issue:** Integration tests hang when database is unavailable.

**Fix Applied:**
Implemented `CheckDatabaseAvailability()` function that verifies connection before test execution.

**Files Modified:**
- `test/test_helper.go` - Added `CheckDatabaseAvailability()` function
- `test/test_helper_test.go` - Added unit tests for availability check

### Phase 3: Date Abstraction (P2) - IN PROGRESS

**Issue:** `GetToday()` uses `time.Now()` directly, making tests non-deterministic.

**Fix Applied:**
Introduce `GetTodayFunc` variable for dependency injection with `SetTestDate()` helper.

**Files Modified:**
- `internal/service/dashboard/day_service.go` - Added date abstraction layer
- `internal/service/dashboard/speculate_service.go` - Updated to use abstracted date function

### Phase 4: Validation (P3) - NOT STARTED

Run full test suite with `-race` flag and verify all acceptance criteria met.

---

## Next Steps

1. ✅ Complete Phase 1 context timeout fixes
2. ✅ Complete Phase 2 database availability checks  
3. ✅ Complete Phase 3 date abstraction layer
4. ⏳ Run validation tests
5. ⏳ Update documentation in QWEN.md
<!-- SECTION:NOTES:END -->

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
