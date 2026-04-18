---
id: RDL-069
title: >-
  [doc-005 Phase 4] Create expected values file and update PRD with final
  results
status: To Do
assignee:
  - catarina
created_date: '2026-04-18 11:48'
updated_date: '2026-04-18 16:57'
labels:
  - phase-4
  - test-automation
  - prd-update
dependencies: []
references:
  - 'PRD Section: Test Artifacts'
  - test/expected-values.go
documentation:
  - doc-005
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create test/expected-values.go with calculated expected values for all acceptance criteria tests, and update the PRD document with implementation results and verification status.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Expected values file created with all calculated test data
- [ ] #2 PRD updated with implementation results and verification status
- [ ] #3 Traceability matrix completed for all requirements
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### Implementation Plan: Create Expected Values File and Update PRD

---

## 1. Technical Approach

### Overview
This task has two main components:
1. **Create `test/expected-values.go`** - A Go test data file containing pre-calculated expected values for all acceptance criteria tests
2. **Update PRD (doc-005)** - Add Phase 4 implementation results, verification status, and traceability matrix completion

### Expected Values File Design

The `test/expected-values.go` file will:

- Define a package `testdata` with shared test value structures
- Include calculated expected values for all project fields (progress, days_unreading, median_day, finished_at, status)
- Provide reference data for comparing Go API vs Rails API responses
- Support both unit and integration test scenarios

**Key Data Structures:**
```go
type ExpectedProject struct {
    ID              int64
    Name            string
    TotalPage       int
    Page            int
    StartedAt       *time.Time
    Progress        *float64
    Status          *string
    LogsCount       *int
    DaysUnreading   *int
    MedianDay       *float64
    FinishedAt      *time.Time
    Logs            []ExpectedLog
}

type ExpectedLog struct {
    ID         int64
    ProjectID  int64
    Data       string
    StartPage  int
    EndPage    int
    Note       *string
}
```

### Test Scenarios to Cover

| Scenario | Fields to Calculate | Source |
|----------|---------------------|--------|
| Normal project with logs | All calculated fields | project-450-go.json, project-450-rails.json |
| Project with no logs | progress, status only | Derived from project data |
| Completed project | finished_at = null | page >= total_page logic |
| Zero page project | All zeros/nils | Edge case handling |

---

## 2. Files to Modify

### New Files to Create

| File | Purpose | Lines Estimate |
|------|---------|----------------|
| `test/testdata/expected-values.go` | Main expected values definitions | ~300-400 |
| `test/testdata/project-450-data.go` | Project 450 specific data | ~100-150 |
| `test/testdata/generate_expected.go` | Script to regenerate expected values | ~80-120 |

### Modified Files

| File | Changes | Reason |
|------|---------|--------|
| `backlog/docs/doc-005 - PRD-Complete-API-Response-Alignment-Project-450-Resolution.md` | Add Phase 4 section, update traceability matrix, increment version | Document implementation results and verification status |

---

## 3. Dependencies

### Prerequisites
- ✅ Go 1.25.7 available in environment
- ✅ PostgreSQL running with test database
- ✅ `test/data/project-450-go.json` exists (already captured)
- ✅ `test/data/project-450-rails.json` exists (already captured)

### Required Packages
```go
// Existing dependencies already in go.mod:
// - github.com/jackc/pgx/v5
// - github.com/joho/godotenv
// - No new dependencies required for this task
```

### Build/Run Dependencies
```bash
# For generating expected values
go run test/testdata/generate_expected.go

# For running tests
go test ./test/...
```

---

## 4. Code Patterns

### Following Existing Conventions

**Test File Structure:**
```go
package testdata

import (
    "time"
    // ... other imports
)

// Expected values for Project 450
var Project450 = ExpectedProject{
    ID: 450,
    Name: "História da Igreja VIII.1",
    TotalPage: 691,
    Page: 691,
    // ... other fields
}
```

**Field Calculation Patterns (from existing code):**

```go
// From internal/domain/models/project.go - follow these patterns:

// Progress calculation
func CalculateProgress() *float64 {
    if p.TotalPage <= 0 {
        return floatPtr(0.0)
    }
    progress := (float64(p.Page) / float64(p.TotalPage)) * 100
    if progress > 100 {
        progress = 100
    }
    return floatPtr(progress)
}

// Days unreading calculation
func CalculateDaysUnreading(logs []*dto.LogResponse) *int {
    // Logic to calculate days since last reading activity
    // Uses max(log.data) or started_at as fallback
}

// Median day calculation  
func CalculateMedianDay() *float64 {
    // median_day = page / days_reading (rounded to 2 decimals)
    // Returns 0.0 for edge cases
}
```

**JSON Serialization:**
```go
// Use json tags matching existing patterns
type ExpectedProject struct {
    ID         int64     `json:"id"`
    Name       string    `json:"name"`
    TotalPage  int       `json:"total_page"`
    Page       int       `json:"page"`
    StartedAt  *string   `json:"started_at,omitempty"`
    Progress   *float64  `json:"progress,omitempty"`
    // ... other fields
}
```

### Naming Conventions
- Package name: `testdata`
- Variable names: `Project450`, `ExpectedProjects`, `TestCases`
- Function names: `CalculateExpectedValues()`, `LoadExpectedData()`

---

## 5. Testing Strategy

### Unit Tests for Expected Values

**File:** `test/testdata/expected-values_test.go`

```go
func TestExpectedValues_Project450(t *testing.T) {
    // Verify all calculated fields match expected values
    tests := []struct {
        name string
        fn   func() error
    }{
        {
            name: "progress calculation",
            fn:   func() error { /* verify progress */ },
        },
        {
            name: "days_unreading calculation", 
            fn:   func() error { /* verify days_unreading */ },
        },
        {
            name: "median_day calculation",
            fn:   func() error { /* verify median_day */ },
        },
        {
            name: "finished_at calculation",
            fn:   func() error { /* verify finished_at */ },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if err := tt.fn(); err != nil {
                t.Errorf("%s: %v", tt.name, err)
            }
        })
    }
}
```

### Integration Tests

**File:** `test/integration/expected_values_integration_test.go`

```go
func TestExpectedValues_Integration(t *testing.T) {
    helper, err := test.SetupTestDB()
    if err != nil {
        t.Fatal(err)
    }
    defer helper.Close()
    
    // Setup test schema and data
    if err := helper.SetupTestSchema(); err != nil {
        t.Fatal(err)
    }
    
    // Insert expected test data
    if err := insertExpectedData(helper.Pool); err != nil {
        t.Fatal(err)
    }
    
    // Run comparison tests against calculated values
    compareWithExpectedValues(t, helper.Pool)
}
```

### Test Coverage Requirements

| Component | Coverage Target | Details |
|-----------|-----------------|---------|
| Expected value calculations | 100% | All edge cases covered |
| Project scenarios | 100% | Normal, empty, completed, zero page |
| Date calculations | 100% | Different formats, timezones |
| JSON serialization | 100% | All fields, omitempty behavior |

---

## 6. Risks and Considerations

### Known Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Test data becomes stale | Medium | Document generation script; add CI check for freshness |
| Timezone calculation differences | High | Use explicit timezone in test context; document assumptions |
| JSON format mismatch (Go vs Rails) | Medium | Include both formats in expected values; note differences |
| Database state inconsistent | High | Use unique test database per run; clean up properly |

### Trade-offs

1. **Hard-coded vs Generated Values**
   - Decision: Store hard-coded values for stability, provide generator script for updates
   - Rationale: Test data should be deterministic; generation script ensures reproducibility

2. **Complete vs Partial Expected Values**
   - Decision: Include all fields even if null/zero
   - Rationale: Explicit defaults prevent silent failures; easier to compare

3. **Timezone Handling**
   - Decision: Use Brazil timezone (BRT) as default for consistency with Rails behavior
   - Rationale: Matches existing implementation in `project.go`

### Blocking Issues

- None identified - this is a documentation/testing enhancement task

### Acceptance Criteria Verification

| AC | Verification Method | Expected Evidence |
|----|---------------------|-------------------|
| #1 Expected values file created | File exists at `test/expected-values.go` | Go file with test data structures |
| #2 PRD updated | Document modified with Phase 4 section | Updated markdown with results |
| #3 Traceability completed | Matrix shows all "Done" status | All requirements linked to tests |

---

## 7. Implementation Steps

### Phase 1: Create Expected Values Structure (Day 1)
- [ ] Create `test/testdata/` directory
- [ ] Define `ExpectedProject` and `ExpectedLog` structs
- [ ] Implement calculation functions for all fields
- [ ] Write `generate_expected.go` script
- [ ] Run script to generate initial values

### Phase 2: Write Unit Tests (Day 2)
- [ ] Create `expected-values_test.go`
- [ ] Test each calculated field
- [ ] Test edge cases (null, zero, empty)
- [ ] Achieve 100% coverage

### Phase 3: Write Integration Tests (Day 3)
- [ ] Create `expected_values_integration_test.go`
- [ ] Setup test database with expected data
- [ ] Run comparison tests
- [ ] Verify against Rails API values

### Phase 4: Update PRD (Day 4)
- [ ] Add "Phase 4 Implementation Results" section
- [ ] Update traceability matrix with completion status
- [ ] Increment version to 1.0.1
- [ ] Add verification evidence links

### Phase 5: Documentation (Day 5)
- [ ] Update `test/README.md` with expected values usage
- [ ] Document generation script usage
- [ ] Add troubleshooting section

---

## 8. Verification Checklist

### Before Completion
- [ ] All unit tests pass (`go test ./test/testdata/...`)
- [ ] All integration tests pass (`go test ./test/integration/...`)
- [ ] `go fmt` passes with no changes
- [ ] `go vet` reports no issues
- [ ] Coverage >= 80% for new code
- [ ] PRD markdown renders correctly

### Acceptance Criteria Checklist
- [ ] `test/expected-values.go` exists and is valid Go
- [ ] Expected values calculated correctly (verified against Rails)
- [ ] PRD updated with Phase 4 section
- [ ] Traceability matrix complete
- [ ] Version incremented to 1.0.1

---

## Summary

| Aspect | Details |
|--------|---------|
| **Files Created** | `test/testdata/expected-values.go`, `generate_expected.go` |
| **Files Modified** | `backlog/docs/doc-005 - PRD-Complete-API-Response-Alignment-Project-450-Resolution.md` |
| **Test Files Added** | `test/testdata/expected-values_test.go`, `integration/expected_values_integration_test.go` |
| **Estimated Effort** | 3-5 days |
| **Blocking Issues** | None |
| **Dependencies** | Existing test infrastructure, Go 1.25.7 |

---

*Implementation Plan created: 2026-04-18*
*Task: RDL-069 - [doc-005 Phase 4] Create expected values file and update PRD with final results*
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
- [ ] #13 Expected values validated against Rails API responses
- [ ] #14 PRD version incremented to 1.0.1
<!-- DOD:END -->
