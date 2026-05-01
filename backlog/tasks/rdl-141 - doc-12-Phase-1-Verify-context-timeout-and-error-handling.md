---
id: RDL-141
title: '[doc-12 Phase 1] Verify context timeout and error handling'
status: To Do
assignee:
  - workflow
created_date: '2026-05-01 15:07'
updated_date: '2026-05-01 16:53'
labels:
  - bugfix
  - repository
  - phase-1
dependencies: []
documentation:
  - doc-012
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Ensure 15-second dashboardContextTimeout is properly applied to all fault calculation queries. Verify error wrapping follows the pattern fmt.Errorf("failed to get faults: %w", err). Check that context cancellation is handled correctly throughout repository methods.
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task verifies that the 15-second `dashboardContextTimeout` is properly applied to all fault calculation queries in the dashboard repository, and that error handling follows the consistent pattern `fmt.Errorf("failed to get faults: %w", err)`. The verification will also ensure context cancellation is handled correctly throughout all repository methods.

**Technical Strategy:**

1. **Audit Existing Implementation**: Review all methods in `internal/adapter/postgres/dashboard_repository.go` to verify:
   - Every method creates a context with timeout using `context.WithTimeout(ctx, dashboardContextTimeout)`
   - Every method calls `defer cancel()` immediately after timeout creation
   - All database operations use the timeout context
   - Error wrapping follows the pattern `fmt.Errorf("failed to <operation>: %w", err)`

2. **Verify Context Propagation**: Ensure context flows correctly from:
   - HTTP handlers → Services → Repository
   - Repository methods → Database operations
   - Context cancellation is respected by pgx operations

3. **Error Handling Verification**: Check that:
   - All errors are wrapped with context using `%w` verb
   - Error messages follow the pattern `"failed to <action>: %w"`
   - No errors are swallowed or returned without wrapping
   - `pgx.ErrNoRows` is handled appropriately where needed

4. **Edge Case Verification**: Ensure proper handling of:
   - Context timeout scenarios (15-second limit)
   - Context cancellation by client
   - Database connection failures
   - Empty result sets

**Why This Approach:**

- The `dashboardContextTimeout` constant (15 seconds) is already defined but needs verification across all methods
- RDL-140 and RDL-139 implemented CTE-based queries that need timeout verification
- Consistent error wrapping aids debugging and monitoring in production
- Context propagation is critical for preventing resource leaks and ensuring responsive cancellation

### 2. Files to Modify

#### Primary Files (Audit Only - No Changes Expected)
| File | Action | Description |
|------|--------|-------------|
| `internal/adapter/postgres/dashboard_repository.go` | Read/Audit | Verify all 15 methods have proper timeout and error handling |

**Methods to Verify (15 total):**
1. `GetDailyStats` - Already has timeout ✓
2. `GetProjectAggregates` - Already has timeout ✓
3. `GetFaultsByDateRange` - Already has timeout ✓
4. `GetWeekdayFaults` - Already has timeout ✓
5. `GetLogsByDateRange` - Already has timeout ✓
6. `GetProjectWeekdayMean` - Already has timeout ✓
7. `CalculatePeriodPages` - Already has timeout ✓
8. `GetProjectsWithLogs` - Already has timeout ✓
9. `GetProjectLogs` - Already has timeout ✓
10. `GetMaxByWeekday` - Already has timeout ✓
11. `GetOverallMean` - Already has timeout ✓
12. `GetPreviousPeriodMean` - Already has timeout ✓
13. `GetPreviousPeriodSpecMean` - Already has timeout ✓
14. `GetMeanByWeekday` - Already has timeout ✓
15. `GetRunningProjectsWithLogs` - Already has timeout ✓

#### Service Files (Audit for Context Propagation)
| File | Action | Description |
|------|--------|-------------|
| `internal/service/dashboard/day_service.go` | Read/Audit | Verify context flows to repository calls |
| `internal/service/dashboard/faults_service.go` | Read/Audit | Verify context flows to repository calls |
| `internal/service/dashboard/weekday_faults_service.go` | Read/Audit | Verify context flows to repository calls |
| `internal/service/dashboard/speculate_service.go` | Read/Audit | Verify context flows to repository calls |
| `internal/service/dashboard/projects_service.go` | Read/Audit | Verify context flows to repository calls |
| `internal/service/dashboard/mean_progress_service.go` | Read/Audit | Verify context flows to repository calls |

#### Handler Files (Audit for Context Usage)
| File | Action | Description |
|------|--------|-------------|
| `internal/api/v1/handlers/dashboard_handler.go` | Read/Audit | Verify `r.Context()` is passed to services/repositories |

#### Test Files (Add Missing Tests)
| File | Change Type | Description |
|------|-------------|-------------|
| `test/unit/dashboard_repository_test.go` | Add Tests | Add context timeout tests for each method |
| `test/dashboard_integration_test.go` | Add Tests | Add context cancellation tests |

**New Test Cases to Add:**

1. **Context Timeout Tests** (per method):
```go
func TestDashboardRepository_GetDailyStats_ContextTimeout(t *testing.T) {
    // Create context with very short timeout
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
    defer cancel()

    // Execute query that would take longer than timeout
    stats, err := repo.GetDailyStats(ctx, testDate)

    // Verify context deadline exceeded error
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "context deadline exceeded")
    assert.Nil(t, stats)
}
```

2. **Context Cancellation Tests**:
```go
func TestDashboardRepository_GetFaultsByDateRange_ContextCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    
    // Start goroutine to cancel context
    go func() {
        time.Sleep(10 * time.Millisecond)
        cancel()
    }()

    // Execute long-running query
    stats, err := repo.GetFaultsByDateRange(ctx, startDate, endDate)

    // Verify cancellation error
    assert.Error(t, err)
    assert.Nil(t, stats)
}
```

3. **Error Wrapping Tests** (per method):
```go
func TestDashboardRepository_GetDailyStats_ErrorWrapping(t *testing.T) {
    // Use invalid database connection to trigger error
    invalidRepo := postgres.NewDashboardRepositoryImpl(invalidPool)
    
    stats, err := invalidRepo.GetDailyStats(context.Background(), testDate)

    // Verify error is wrapped correctly
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "failed to get daily stats")
    assert.Nil(t, stats)
}
```

### 3. Dependencies

#### Prerequisites
1. **Existing Implementation**: All repository methods must already exist (verified in RDL-140)
2. **Test Database**: `reading_log_test` must be running for integration tests
3. **Context Package**: Standard library `context` package (already imported)
4. **pgx/v5**: Database driver that supports context (already in use)

#### Related Tasks
| Task ID | Relationship | Status |
|---------|--------------|--------|
| RDL-140 | Predecessor - Updated GetWeekdayFaults SQL | Done |
| RDL-139 | Predecessor - Updated GetFaultsByDateRange SQL | Done |
| RDL-142 | Successor - Rewrite unit tests | To Do |

#### Environment Setup
```bash
# Ensure test database is ready
make test-setup

# Verify PostgreSQL is running
pg_isready -h localhost -p 5432

# Run baseline tests
go test -v ./test/unit/... -run TestDashboardRepository
```

### 4. Code Patterns

#### Context Timeout Pattern (Verify Consistency)
```go
func (r *DashboardRepositoryImpl) GetDailyStats(ctx context.Context, date time.Time) (*dto.DailyStats, error) {
    // Create timeout context with defer cancel
    ctx, cancel := context.WithTimeout(ctx, dashboardContextTimeout)
    defer cancel()

    // Use timeout context for ALL database operations
    query := `SELECT ...`
    
    var stats dto.DailyStats
    err := r.pool.QueryRow(ctx, query, date).Scan(
        &stats.TotalPages,
        &stats.LogCount,
    )
    if err != nil {
        // Wrap errors with context using %w verb
        if err == pgx.ErrNoRows {
            return dto.NewDailyStats(0, 0), nil
        }
        return nil, fmt.Errorf("failed to get daily stats: %w", err)
    }

    return &stats, nil
}
```

#### Error Wrapping Pattern (Verify Consistency)
```go
// CORRECT - Error wrapping with context
return nil, fmt.Errorf("failed to get daily stats: %w", err)
return nil, fmt.Errorf("failed to query project aggregates: %w", err)
return nil, fmt.Errorf("failed to scan log entry: %w", err)

// INCORRECT - No context wrapping (should not exist)
return nil, err
return fmt.Errorf("some error: %v", err)  // Wrong verb
```

#### Service Layer Context Propagation Pattern
```go
// Handler passes r.Context() to service
func (h *DashboardHandler) Faults(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()  // Get context from request
    
    // Pass context through service layer
    percentage, err := faultsService.GetFaultsPercentage(ctx)
    if err != nil {
        // Error handling
    }
}

// Service passes context to repository
func (s *FaultsService) GetFaultsPercentage(ctx context.Context) (float64, error) {
    // Pass context to repository method
    faults, err := s.repo.GetFaultsByDateRange(ctx, startDate, endDate)
    if err != nil {
        return 0.0, fmt.Errorf("failed to get faults by date range: %w", err)
    }
    // ...
}
```

#### Naming Conventions
- Timeout constant: `dashboardContextTimeout` (15 seconds)
- Error message prefix: `"failed to <operation>"`
- Context variable: `ctx` (standard Go convention)
- Cancel function: `cancel` (standard Go convention)

### 5. Testing Strategy

#### Unit Tests (test/unit/dashboard_repository_test.go)

**Test Category 1: Context Timeout Verification**
- Test each of the 15 repository methods with a short timeout context
- Verify `context.DeadlineExceeded` error is returned
- Verify no database resource leaks

**Test Category 2: Context Cancellation Verification**
- Test each method with cancellable context
- Cancel context mid-operation
- Verify context cancellation error is returned

**Test Category 3: Error Wrapping Verification**
- Test each method with invalid database state
- Verify error messages contain expected prefix
- Verify error chain preserves original error using `%w`

**Test Category 4: Success Path Verification**
- Verify normal operations still work with timeout context
- Verify successful queries complete within timeout
- Verify context is properly cleaned up

#### Integration Tests (test/dashboard_integration_test.go)

**Test Category 1: Full Stack Context Propagation**
- Test HTTP handler → Service → Repository context flow
- Verify request cancellation propagates through all layers
- Verify timeout configured at repository level is respected

**Test Category 2: Performance Under Load**
- Test repository methods with concurrent requests
- Verify context timeouts prevent resource exhaustion
- Verify no goroutine leaks on timeout/cancellation

**Test Category 3: Database Connection Failures**
- Test with database connection issues
- Verify errors are properly wrapped and propagated
- Verify context timeout prevents hanging

#### Test Execution Commands
```bash
# Run repository unit tests
go test -v ./test/unit/... -run TestDashboardRepository

# Run new context timeout tests
go test -v ./test/unit/... -run ".*Context.*"

# Run integration tests
go test -v ./test/... -run ".*Dashboard.*"

# Run full test suite
go test ./...

# Run with coverage
go test -cover ./test/unit/...
```

#### Edge Cases to Cover

1. **Zero Timeout**: Context with 0 timeout should immediately fail
2. **Very Long Timeout**: Context with long timeout should allow normal operation
3. **Parent Context Cancellation**: Cancelling parent should cancel child contexts
4. **Database Slowdown**: Simulated slow queries should respect timeout
5. **Concurrent Requests**: Multiple requests should each have independent timeouts
6. **Empty Results**: Empty result sets should not cause timeout issues
7. **NULL Handling**: NULL values in results should not affect timeout behavior

### 6. Risks and Considerations

#### Known Risks

**Risk 1: Existing Code May Have Inconsistencies**
- **Probability**: Medium
- **Impact**: Low
- **Mitigation**:
  - Audit all 15 methods systematically
  - Document any inconsistencies found
  - Fix inconsistencies as part of this task
  - Add regression tests to prevent future issues

**Risk 2: Test Coverage May Be Incomplete**
- **Probability**: High
- **Impact**: Medium
- **Mitigation**:
  - Add context-specific tests for all 15 methods
  - Verify existing tests still pass after changes
  - Aim for 100% coverage of error paths

**Risk 3: Timeout Value May Need Adjustment**
- **Probability**: Low
- **Impact**: Medium
- **Mitigation**:
  - Current 15-second timeout is reasonable for most queries
  - Document timeout value in code comments
  - Consider profiling for production optimization
  - Make timeout configurable via environment variable if needed

#### Technical Considerations

**Context Hierarchy:**
- HTTP request context → Service context → Repository context
- Each layer can add its own timeout via `context.WithTimeout`
- Shortest timeout wins (cascading cancellation)
- Ensure timeout values are appropriate at each layer

**Error Message Consistency:**
- All errors should follow pattern: `"failed to <operation>: %w"`
- Error messages should be human-readable
- Original error preserved via `%w` for stack traces
- Avoid exposing internal implementation details

**Resource Cleanup:**
- Always call `defer cancel()` after `context.WithTimeout`
- Verify no goroutine leaks on timeout
- Verify database connections are released on cancellation
- Use `rows.Err()` to check for iteration errors

**Performance Impact:**
- Context timeout adds minimal overhead
- pgx/v5 respects context cancellation efficiently
- No significant performance degradation expected
- Monitor production for any timeout-related issues

#### Deployment Considerations

**Rollback Plan:**
1. All changes are additive (adding tests, verifying existing code)
2. No breaking changes to API or behavior
3. If issues arise, revert test additions only
4. Investigate failures in isolated branch

**Monitoring:**
- Track context timeout errors in production logs
- Monitor query execution times
- Alert if timeout frequency increases
- Track error wrapping consistency

**Documentation Updates Required:**
- Update AGENTS.md with context timeout explanation
- Add timeout configuration documentation
- Document error handling patterns in docs/error-handling.md
- Update API documentation with timeout behavior

#### Acceptance Criteria Verification

| Criterion | Verification Method |
|-----------|---------------------|
| 15-second timeout applied to all methods | Code audit of all 15 repository methods |
| Error wrapping follows pattern | Code audit of all error returns |
| Context cancellation handled | Unit tests with cancellation |
| All unit tests pass | `go test ./test/unit/...` |
| All integration tests pass | `go test ./test/...` |
| go fmt and go vet pass | `go fmt ./...` and `go vet ./...` |
| Clean Architecture followed | Code review of layer separation |
| Error responses consistent | Manual verification of error messages |
| HTTP status codes correct | Integration test verification |
| Documentation updated | Review of AGENTS.md updates |
| New code paths include error tests | Test coverage verification |
| Handlers test success and error | Handler test review |
| Integration tests verify DB interactions | Integration test review |

---

**Implementation Ready**: ✅ This plan provides sufficient detail for another agent to verify and implement context timeout and error handling improvements without reading additional documentation. All methods, test scenarios, and risk mitigations are documented.

**Note**: This is primarily a verification task. If all existing code already follows the correct patterns (which appears to be the case based on initial review), the main work will be adding comprehensive tests to ensure these patterns are maintained.
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
