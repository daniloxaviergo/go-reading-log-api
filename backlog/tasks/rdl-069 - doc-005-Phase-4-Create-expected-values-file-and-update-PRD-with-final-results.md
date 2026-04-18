---
id: RDL-069
title: >-
  [doc-005 Phase 4] Create expected values file and update PRD with final
  results
status: To Do
assignee:
  - catarina
created_date: '2026-04-18 11:48'
updated_date: '2026-04-18 16:12'
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
### 1. Technical Approach

This task involves creating test infrastructure artifacts and updating documentation to support API response validation testing.

**Key Components:**

1. **Expected Values File (`test/expected-values.go`):**
   - Create a Go test utility file that defines expected values for all acceptance criteria
   - Include pre-calculated values derived from Rails API responses
   - Support both unit and integration test scenarios
   - Provide helper functions for comparing actual vs expected values

2. **PRD Update:**
   - Add implementation results section documenting completed work
   - Include verification status for each acceptance criterion
   - Document any deviations or known issues
   - Update traceability matrix with completed items

3. **Test Data Management:**
   - Leverage existing recorded API responses in `test/data/`
   - Create expected values based on Rails API as source of truth
   - Ensure test data is versioned and reproducible

**Architecture Decision:**
- Use Go's table-driven test pattern for maintainability
- Keep expected values immutable (generated from Rails API)
- Provide clear error messages when actual vs expected differ

---

### 2. Files to Modify

#### New Files to Create:

| File | Purpose |
|------|---------|
| `test/expected-values.go` | Main expected values file with all calculated test data |
| `test/data/project-450-go.json` (update) | Record current Go API response for regression testing |
| `test/data/project-450-rails.json` (update) | Record current Rails API response as source of truth |

#### Files to Reference (Read-Only for this task):

| File | Purpose |
|------|---------|
| `internal/domain/dto/project_response.go` | ProjectResponse DTO structure |
| `internal/domain/dto/log_response.go` | LogResponse DTO structure |
| `internal/domain/models/project.go` | Project model with calculation methods |
| `test/compare_responses.sh` | Existing comparison script |

---

### 3. Dependencies

**Prerequisites:**
- ✅ Go 1.25.7 installed and configured
- ✅ PostgreSQL test database accessible
- ✅ Both Go API and Rails API running for data capture
- ✅ Existing test infrastructure (`test/test_helper.go`)

**Blocking Issues:**
- None - this is a documentation/test artifact task

**Setup Steps:**
1. Ensure both APIs are running and accessible
2. Verify `test/data/` directory exists
3. Run `compare_responses.sh` to capture current state
4. Generate expected values from captured data

---

### 4. Code Patterns

**Expected Values File Structure:**

```go
package test

import (
    "go-reading-log-api-next/internal/domain/dto"
)

// ExpectedProjectValues contains pre-calculated expected values for project 450
type ExpectedProjectValues struct {
    ID         int64   `json:"id"`
    Name       string  `json:"name"`
    TotalPage  int     `json:"total_page"`
    Page       int     `json:"page"`
    Progress   float64 `json:"progress"`
    Status     string  `json:"status"`
    LogsCount  int     `json:"logs_count"`
    DaysUnread int     `json:"days_unreading"`
    MedianDay  float64 `json:"median_day"`
}

// ExpectedValues holds all expected values for acceptance criteria
type ExpectedValues struct {
    Project450 *ExpectedProjectValues
    // Add more projects as needed
}

// GetExpectedValues returns the complete set of expected values
func GetExpectedValues() *ExpectedValues {
    return &ExpectedValues{
        Project450: &ExpectedProjectValues{
            ID:         450,
            Name:       "História da Igreja VIII.1",
            TotalPage:  691,
            Page:       691,
            Progress:   100.0,
            Status:     "finished",
            LogsCount:  38,
            DaysUnread: 16,  // From Rails API
            MedianDay:  11.91,
        },
    }
}
```

**Comparison Helper Pattern:**

```go
// CompareProjectResponse compares actual vs expected project values
func CompareProjectResponse(actual *dto.ProjectResponse, expected *ExpectedProjectValues) *ComparisonResult {
    result := &ComparisonResult{Passed: true, Differences: []string{}}
    
    if actual.ID != expected.ID {
        result.Differences = append(result.Differences, 
            fmt.Sprintf("ID: expected %d, got %d", expected.ID, actual.ID))
        result.Passed = false
    }
    
    // Compare calculated fields with tolerance
    if actual.Progress != nil {
        diff := *actual.Progress - expected.Progress
        if diff < 0 { diff = -diff }
        if diff > 0.01 {
            result.Differences = append(result.Differences, 
                fmt.Sprintf("Progress: expected %.2f, got %.2f", expected.Progress, *actual.Progress))
            result.Passed = false
        }
    }
    
    return result
}
```

---

### 5. Testing Strategy

**Unit Tests (to be created in `test/unit/`):**

| Test File | Coverage |
|-----------|----------|
| `expected_values_test.go` | Verify expected values match Rails API |
| `comparison_helpers_test.go` | Test comparison logic |

**Integration Tests (to be created in `test/integration/`):**

| Test File | Coverage |
|-----------|----------|
| `expected_values_integration_test.go` | Full pipeline: DB → API → Expected Values |

**Test Execution:**
```bash
# Run all tests with coverage
go test -v ./test/...

# Run specific test file
go test -v ./test/unit/expected_values_test.go

# Generate coverage report
go test -coverprofile=coverage.out ./test/...
```

---

### 6. Risks and Considerations

**Known Issues:**
1. **Date Calculation Tolerance:** The `days_unreading` field has a 1-day tolerance due to timezone differences between Go (UTC) and Rails (BRT)
2. **JSON:API vs Flat JSON:** Go API returns flat JSON while Rails uses JSON:API envelope - comparison logic must handle both formats
3. **Floating Point Precision:** `median_day` and `progress` may have minor floating point differences (< 0.01)

**Trade-offs:**
- Expected values are generated from Rails API as "source of truth" - this assumes Rails is correct
- Some fields (like `started_at` format) may differ between APIs due to implementation choices
- Test data should be periodically regenerated to capture API changes

**Deployment Considerations:**
- No database migrations required for this task
- No breaking changes to existing API contracts
- Test artifacts can be safely added without affecting production
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
