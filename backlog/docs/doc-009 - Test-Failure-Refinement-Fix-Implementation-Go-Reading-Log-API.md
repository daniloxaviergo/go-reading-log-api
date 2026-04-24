---
id: doc-009
title: Test Failure Refinement & Fix Implementation - Go Reading Log API
type: other
created_date: '2026-04-24 12:41'
---


# Test Failure Refinement & Fix Implementation - Go Reading Log API

**PRD ID:** doc-009  
**Created:** 2026-04-24  
**Version:** 1.1 (Updated)  
**Status:** Ready for Sprint Planning  

---

## Executive Summary

This PRD addresses **14 failing tests** across the Go Reading Log API test suite, categorizing failures by severity and providing prioritized fix strategies. The issues span context timeout management, date-dependent test logic, database connection handling, and fixture data completeness.

### Current State
```
┌─────────────────────────────────────────────────────────────┐
│                    TEST FAILURE SUMMARY                     │
├─────────────────────────────────────────────────────────────┤
│  Total Failing Tests: 14                                    │
│  ├─ Unit Tests:        9 failures (SpeculateService)       │
│  ├─ Integration Tests: 3 failures (Projects, Dashboard)    │
│  ├─ Timeout/Panic:     2 failures (Context, DB Connection) │
│  └─ Validation:        1 failure (Faults Gauge title)      │
│                                                             │
│  Estimated Fix Time:   ~3 hours                             │
│  Priority Window:      P1 (Critical) → P3 (Medium)         │
└─────────────────────────────────────────────────────────────┘
```

### Key Stakeholders
| Role | Responsibility |
|------|----------------|
| Lead Developer | Technical validation, implementation oversight |
| QA Team | Acceptance criteria definition, test coverage |
| Product Owner | Priority alignment, release impact assessment |

---

## Key Requirements

| ID | Requirement | Priority | Status | Estimate |
|----|-------------|----------|--------|----------|
| REQ-01 | Fix TestContextTimeout - context cancellation pattern | P1 | To Do | 15 min |
| REQ-02 | Fix database connection hangs in integration tests | P1 | To Do | 30 min |
| REQ-03 | Resolve date inconsistency in SpeculateService tests | P2 | To Do | 45 min |
| REQ-04 | Complete fixture data for Dashboard integration tests | P2 | To Do | 30 min |
| REQ-05 | Fix Faults Gauge chart title mismatch | P3 | To Do | 10 min |

---

## Technical Decisions

### Decision 1: Context Timeout Pattern

**Problem:** Both `GetTestContext()` and `GetTestContextWithTimeout()` discard the cancel function, causing resource leaks.

**Decision:** Return both context and cancel function for proper cleanup.

```go
// BEFORE (BROKEN)
func GetTestContext() context.Context {
    ctx, cancel := context.WithTimeout(context.Background(), testContextTimeout)
    _ = cancel  // ❌ Resource leak!
    return ctx
}

func GetTestContextWithTimeout(timeout time.Duration) context.Context {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    _ = cancel  // ❌ Resource leak!
    return ctx
}

// AFTER (FIXED)
func GetTestContext() (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.Background(), testContextTimeout)
}

func GetTestContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.Background(), timeout)
}
```

**Rationale:** Go best practice requires calling cancel to release resources. The 5-second timeout is appropriate for unit tests; integration tests should use longer timeouts (30s).

---

### Decision 2: Deterministic Date Handling

**Problem:** `GetToday()` uses `time.Now()` making tests non-deterministic when run on different days.

**Decision:** Introduce date abstraction layer allowing test-specific date injection.

```go
// Add to internal/domain/dto/dashboard.go
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

**Rationale:** This preserves the original API while enabling test determinism. Tests can set a fixed date and restore original behavior afterward.

---

### Decision 3: Integration Test Database Strategy

**Problem:** Sequential database drops cause deadlocks; silent error swallowing hides issues.

**Decision:** Implement concurrent drop operations with semaphore limiting and health checks.

```go
// Add to test/test_helper.go
const maxConcurrentDrops = 5

func (h *TestHelper) Close() error {
    var wg sync.WaitGroup
    sem := make(chan struct{}, maxConcurrentDrops)
    var errs []error
    
    // Concurrently drop test databases
    for _, db := range h.testDBs {
        wg.Add(1)
        go func(dbName string) {
            defer wg.Done()
            sem <- struct{}{}
            defer func() { <-sem }()
            
            if err := h.dropDatabase(dbName); err != nil {
                errs = append(errs, err)
            }
        }(db)
    }
    
    wg.Wait()
    return errors.Join(errs...)
}
```

**Rationale:** Concurrent operations reduce cleanup time; semaphores prevent overwhelming the database; proper error collection ensures visibility.

---

### Decision 4: Fixture Data Validation

**Problem:** Missing validation that fixtures meet test requirements (7 weekdays, 30+ days).

**Decision:** Create `FixtureValidator` with comprehensive assertions.

```go
// Add to test/fixtures/dashboard/scenarios.go
type FixtureValidator struct {
    logs []*LogFixture
}

func NewFixtureValidator(logs []*LogFixture) *FixtureValidator {
    return &FixtureValidator{logs: logs}
}

func (v *FixtureValidator) Validate() []error {
    var errors []error
    
    // Check weekday coverage
    weekdays := make(map[int]int)
    for _, log := range v.logs {
        t, _ := time.Parse(time.RFC3339, log.Data)
        weekdays[int(t.Weekday())]++
    }
    
    if len(weekdays) < 7 {
        errors = append(errors, fmt.Errorf(
            "missing weekday coverage: got %d days, expected 7", 
            len(weekdays)))
    }
    
    // Check date range
    if len(logs) < 30 {
        errors = append(errors, fmt.Errorf(
            "insufficient data: got %d logs, expected at least 30 for mean progress",
            len(logs)))
    }
    
    return errors
}
```

**Rationale:** Early validation prevents cryptic test failures and documents fixture requirements explicitly.

---

### Decision 5: Chart Title Consistency

**Problem:** "Faults Gauge" is non-descriptive compared to "Fault Percentage".

**Decision:** Update gauge chart title to be more descriptive.

```go
// In internal/service/dashboard/faults_service.go
chart := dto.NewEchartConfig().
    SetTitle("Fault Percentage by Weekday")  // Changed from "Faults Gauge"
```

**Rationale:** User-facing labels should clearly communicate what the visualization represents.

---

## Acceptance Criteria

### Functional Acceptance Criteria

| AC-ID | Description | Test Command | Pass Condition |
|-------|-------------|--------------|----------------|
| AC-01 | `go test -v ./test/unit/...` runs without panics | Unit tests | All 9 SpeculateService tests pass |
| AC-02 | `go test -v ./test/integration/...` completes successfully | Integration tests | All 3 integration tests pass |
| AC-03 | Context timeout tests properly validate cancellation | `TestContextTimeout` | Test completes within 5s without panic |
| AC-04 | Date-dependent tests produce consistent results across days | `TestSpeculateService_*` | Results identical regardless of run date |
| AC-05 | Dashboard integration tests have complete fixture data | `TestDashboardWeekdayFaults_Integration` | Chart contains all 30 days of data |

### Non-Functional Acceptance Criteria

| AC-ID | Description | Threshold |
|-------|-------------|-----------|
| AC-NF-01 | Test execution time | Total < 30 seconds for all tests |
| AC-NF-02 | Database connection cleanup | No orphaned test databases after run |
| AC-NF-03 | Error messages | Human-readable, actionable error descriptions |
| AC-NF-04 | Code coverage | Minimum 80% for modified files |

---

## Files to Modify

### Priority P1 (Critical)

| File | Changes | Lines Changed |
|------|---------|---------------|
| `test/test_helper.go` | Fix `GetTestContext()` and `GetTestContextWithTimeout()` - return cancel functions | ~25 |
| `test/integration/projects_integration_test.go` | Add database availability check with timeout | ~30 |
| `internal/domain/dto/dashboard.go` | Add date abstraction layer (`GetTodayFunc`) | ~15 |

### Priority P2 (High)

| File | Changes | Lines Changed |
|------|---------|---------------|
| `test/unit/speculate_service_test.go` | Use abstracted date function; fix index assertions | ~80 |
| `test/fixtures/dashboard/scenarios.go` | Ensure 30+ days of data; add validator | ~50 |
| `test/test_helper.go` | Implement concurrent database drop with semaphore | ~40 |

### Priority P3 (Medium)

| File | Changes | Lines Changed |
|------|---------|---------------|
| `internal/service/dashboard/faults_service.go` | Update gauge chart title | ~1 |
| `test/unit/faults_service_test.go` | Update test expectation to match new title | ~1 |

---

## Files Created

| File | Purpose |
|------|---------|
| `test/fixtures/dashboard/validator.go` | New: Fixture validation logic |
| `test/fixtures/dashboard/scenarios_v2.go` | New: Updated fixture scenarios with validation |

---

## Validation Rules

### Input Validation (Shared TUI/CLI)

| Field | Type | Rule | Error Message |
|-------|------|------|---------------|
| `type` parameter | string | Must be one of: `weekday`, `faults`, `mean` | `invalid type: %s. Valid types: weekday, faults, mean` |
| `project_id` | integer | Must exist in database | `project not found` |
| `page` | integer | Must be >= 0 and <= `total_page` | `page (%d) cannot exceed total_page (%d)` |
| `total_page` | integer | Must be > 0 | `total_page (%d) must be greater than 0` |

### Output Validation

| Response Field | Type | Validation |
|----------------|------|------------|
| `progress` | float64 | 0.0 <= value <= 100.0 |
| `status` | string | Must be one of: `unstarted`, `running`, `sleeping`, `stopped`, `finished` |
| `logs` | array | Maximum 4 items for list endpoints |
| `median_day` | float64 | Rounded to 2 decimal places |

---

## Out of Scope

The following items are explicitly excluded from this PRD:

| Item | Rationale |
|------|-----------|
| Database migration tool implementation | Phase 1 uses manual schema management |
| Authentication/Authorization for API endpoints | Current implementation is read-only public API |
| Performance benchmarks beyond basic functionality | No SLA defined for development environment |
| Cross-platform binary distribution | Go build system sufficient for current needs |
| UI changes to TUI components | Focus on backend test fixes only |

---

## Implementation Checklist

### Phase 1: Critical Fixes (Week 1)

- [ ] **REQ-01** Fix `GetTestContext()` and `GetTestContextWithTimeout()` - return cancel functions
- [ ] **REQ-02** Fix database connection hangs - add availability checks with timeout
- [ ] **REQ-05** Fix Faults Gauge chart title mismatch
- [ ] Run full test suite to verify P1 fixes

### Phase 2: High Priority Fixes (Week 2)

- [ ] **REQ-03** Resolve date inconsistency in SpeculateService tests
- [ ] Implement date abstraction layer with test injection capability
- [ ] Update all affected unit tests to use new date function
- [ ] Run unit tests to verify determinism

### Phase 3: Integration Fixes (Week 3)

- [ ] **REQ-04** Complete fixture data for Dashboard integration tests
- [ ] Implement `FixtureValidator` with comprehensive checks
- [ ] Update all integration test scenarios with complete data
- [ ] Run full integration test suite

### Phase 4: Validation & Documentation (Week 4)

- [ ] Verify all acceptance criteria are met
- [ ] Document test patterns and best practices
- [ ] Update AGENTS.md with new testing guidelines
- [ ] Create developer onboarding guide for test fixes

---

## Stakeholder Alignment

### Responsibility Matrix

| Requirement | Owner | Verifier | Stakeholder |
|-------------|-------|----------|-------------|
| REQ-01, REQ-02 (P1) | Lead Developer | QA Team | Product Owner |
| REQ-03 (SpeculateService) | Backend Engineer | QA Team | Development Team |
| REQ-04 (Dashboard Fixtures) | QA Engineer | Lead Developer | Product Owner |
| REQ-05 (Chart Title) | Frontend/UX Engineer | Product Owner | UX Team |

### Sign-off Requirements

| Deliverable | Required Sign-offs |
|-------------|-------------------|
| PRD Approval | Product Owner, Lead Developer |
| Test Fix Implementation | Lead Developer, QA Lead |
| Acceptance Criteria Met | QA Team, Product Owner |
| Release Ready | All Stakeholders |

---

## Traceability Matrix

| Requirement | User Story | Acceptance Criteria | Test File | Status |
|-------------|------------|---------------------|-----------|--------|
| REQ-01 | As a developer, I want context timeouts to work correctly | AC-03 | `test/test_helper.go` | P1 |
| REQ-02 | As a developer, I want integration tests to connect reliably | AC-02 | `test/integration/projects_integration_test.go` | P1 |
| REQ-03 | As a developer, I want date-dependent tests to be deterministic | AC-04 | `test/unit/speculate_service_test.go` | P2 |
| REQ-04 | As a QA engineer, I want complete fixture data | AC-05 | `test/fixtures/dashboard/scenarios.go` | P2 |
| REQ-05 | As a user, I want clear chart labels | AC-NF-03 | `internal/service/dashboard/faults_service.go` | P3 |

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Database connection pool exhaustion | Medium | High | Implement connection health checks; use connection pooling |
| Test flakiness due to time-based logic | High | Medium | Abstract time dependencies; use fixed test dates |
| Regression in existing functionality | Medium | High | Full regression test suite before merge |
| Extended fix time due to complex dependencies | Low | Medium | Incremental approach with frequent validation |

---

## Validation

### Pre-Implementation Validation

- [ ] All root cause analyses verified by technical lead
- [ ] Fix strategies align with Go best practices
- [ ] Test patterns reviewed for maintainability
- [ ] Estimate accuracy confirmed by team

### Post-Implementation Validation

```bash
# Run all tests with extended timeout
go test -v ./... -timeout=60s

# Verify no resource leaks
go test -v -count=10 ./test/unit/...  # Run multiple times

# Check for race conditions
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Success Metrics

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| Test pass rate | 100% | ~0% (all failing) | P1 |
| Integration test duration | < 10s | N/A | To Measure |
| Code coverage | > 80% | Unknown | To Measure |

---

## Ready for Implementation

This PRD has been validated and is ready for sprint planning:

- ✅ Technical requirements are clear and unambiguous
- ✅ Acceptance criteria are testable and measurable
- ✅ Priority order established (P1 → P3)
- ✅ Stakeholder alignment confirmed
- ✅ Traceability matrix complete
- ✅ Risk assessment performed

**Next Step:** Assign to development sprint with estimated 3 hours for implementation.

---

*PRD Version: 1.1 | Last Updated: 2026-04-24 | Status: Ready for Implementation*