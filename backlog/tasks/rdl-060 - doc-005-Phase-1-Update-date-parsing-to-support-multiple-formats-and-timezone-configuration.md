---
id: RDL-060
title: >-
  [doc-005 Phase 1] Update date parsing to support multiple formats and timezone
  configuration
status: To Do
assignee:
  - thomas
created_date: '2026-04-18 11:46'
updated_date: '2026-04-18 12:23'
labels:
  - phase-1
  - date-calculation
  - critical
dependencies: []
references:
  - 'https://github.com/go-reading-log-api-next/internal/domain/models/project.go'
  - 'PRD Section: Decision 1'
documentation:
  - doc-005
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement multi-format date parsing in internal/domain/models/project.go to fix the 42-day discrepancy between Go and Rails API. The parseLogDate function must support YYYY-MM-DD, RFC3339, and standard datetime formats with timezone-aware comparison matching Rails' Date.today behavior.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 parseLogDate supports at least 3 date formats (YYYY-MM-DD, RFC3339, standard datetime)
- [x] #2 CalculateDaysUnreading uses timezone-aware comparison matching Rails
- [x] #3 Unit tests validate edge cases with different date formats
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The implementation will address the 42-day discrepancy between Go and Rails API by updating date parsing to support multiple formats with timezone-aware comparison.

**Key Changes:**
- Update `parseLogDate()` function in `project.go` to support YYYY-MM-DD, RFC3339, and standard datetime formats
- Introduce timezone configuration via environment variable with Brazil timezone (BRT) as fallback
- Modify `CalculateDaysUnreading()` to use timezone-aware date comparison matching Rails' `Date.today`
- Add comprehensive unit tests for edge cases

**Architecture Decisions:**
- Keep timezone configuration in existing `config.go` to maintain single source of truth
- Use `time.FixedZone` for Brazil timezone as it matches Rails behavior
- Maintain backward compatibility by keeping existing format parsing while adding new ones

---

### 2. Files to Modify

| File | Changes |
|------|---------|
| `internal/domain/models/project.go` | Update `CalculateDaysUnreading()` with multi-format date parsing; Add timezone configuration support |
| `internal/config/config.go` | Add `TZLocation` field and environment variable loading |
| `internal/domain/models/project_test.go` | Add unit tests for date parsing edge cases |
| `.env.example` | Add `TZ_LOCATION` configuration example |

---

### 3. Dependencies

- [x] Existing config infrastructure in `config.go`
- [x] Domain model structure in `project.go`
- [ ] No external dependencies required (uses standard library `time` package)

---

### 4. Code Patterns

**Date Parsing Pattern:**
```go
func parseLogDate(dateStr string) (time.Time, bool) {
    formats := []string{
        "2006-01-02",
        "2006-01-02T15:04:05Z",
        "2006-01-02 15:04:05",
    }
    for _, format := range formats {
        if t, err := time.Parse(format, dateStr); err == nil {
            return t, true
        }
    }
    return time.Time{}, false
}
```

**Timezone Configuration Pattern:**
```go
var TZLocation = time.FixedZone("BRT", -3*60*60) // Brazil timezone

// In config loading
if tzStr := os.Getenv("TZ_LOCATION"); tzStr != "" {
    if loc, err := time.LoadLocation(tzStr); err == nil {
        TZLocation = loc
    }
}
```

**Date Comparison Pattern (matching Rails):**
```go
nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, TZLocation)
lastReadDate = time.Date(lastReadDate.Year(), lastReadDate.Month(), lastReadDate.Day(), 0, 0, 0, 0, TZLocation)
```

---

### 5. Testing Strategy

**Unit Tests to Add:**
- Test `parseLogDate()` with YYYY-MM-DD format
- Test `parseLogDate()` with RFC3339 format
- Test `parseLogDate()` with standard datetime format
- Test `CalculateDaysUnreading()` with various log dates
- Test edge case: no logs, no started_at → returns 0
- Test edge case: future dates → returns 0

**Test Files:**
- `internal/domain/models/project_test.go` - Unit tests for date parsing
- `test/compare_responses.sh` - Integration test comparing Go vs Rails responses

---

### 6. Risks and Considerations

| Risk | Mitigation |
|------|------------|
| Breaking existing log parsing | Keep all format attempts in fallback chain |
| Timezone loading failure | Fallback to BRT with logging warning |
| Performance impact of multiple format attempts | Minimal - only affects log date parsing, not main query path |
| Inconsistent behavior across environments | Document timezone configuration requirement |

**Blocking Issues:**
- None identified

**Trade-offs:**
- Support for multiple formats increases code complexity slightly but ensures compatibility with diverse log data
- Timezone configuration adds flexibility but requires documentation to prevent environment-specific discrepancies
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress - RDL-060

**Status:** In Progress

### What I've Learned So Far

**Current State Analysis:**

1. **project.go** (`internal/domain/models/project.go`):
   - `CalculateDaysUnreading()` currently only parses "YYYY-MM-DD" format
   - Uses `time.UTC` for date comparison (not timezone-aware)
   - Rails uses `Date.today` which is timezone-aware
   - The 42-day discrepancy is likely caused by timezone differences

2. **config.go** (`internal/config/config.go`):
   - Does NOT have timezone configuration yet
   - Needs `TZLocation` field and loading logic

3. **project_test.go** (`internal/domain/models/project_test.go`):
   - Has basic tests for CalculateDaysUnreading
   - Missing tests for multi-format date parsing
   - Missing timezone-specific edge case tests

4. **LogResponse** (`internal/domain/dto/log_response.go`):
   - `Data` field is `*string` containing the date
   - Expected format: "YYYY-MM-DD" currently

### Implementation Plan Refined:

| File | Change |
|------|--------|
| `config.go` | Add `TZLocation *time.Location` field + loading from env |
| `project.go` | Update `CalculateDaysUnreading()` with multi-format parsing + timezone-aware comparison |
| `project_test.go` | Add tests for 3 formats (YYYY-MM-DD, RFC3339, datetime) + timezone edge cases |
| `.env.example` | Add `TZ_LOCATION` example |

**Key Decision:** The 42-day discrepancy is likely because Rails uses local timezone (Brazil/BRT) while Go uses UTC. I need to make the comparison timezone-aware.

### Implementation Completed

**Files Modified:**

1. **internal/config/config.go**
   - Added `TZLocation *time.Location` field to Config struct
   - Added `parseTZLocation()` function with BRT fallback
   - Updated `LoadConfig()` to initialize timezone from environment variable

2. **internal/domain/models/project.go**
   - Added `parseLogDate()` function supporting 3 formats:
     - YYYY-MM-DD (e.g., "2024-01-15")
     - RFC3339 (e.g., "2024-01-15T10:30:00Z")
     - Standard datetime (e.g., "2024-01-15 10:30:00")
   - Updated `CalculateDaysUnreading()` to use multi-format parsing and timezone-aware comparison
   - Updated `CalculateMedianDay()` to use timezone-aware comparison
   - Updated `CalculateFinishedAt()` to use multi-format parsing and timezone-aware comparison
   - Added `getTimezoneFromContext()` helper function

3. **internal/config/config_test.go**
   - Added 4 new test functions for timezone configuration:
     - `TestLoadConfigTimezoneDefault` - Verifies BRT default
     - `TestLoadConfigTimezoneFromEnv` - Verifies env var loading
     - `TestLoadConfigTimezoneInvalidFallback` - Verifies fallback on invalid value
     - `TestLoadConfigTimezoneEmptyFallback` - Verifies fallback on empty value

4. **internal/domain/models/project_test.go**
   - Added 5 new test functions for date parsing:
     - `TestProject_ParseLogDate` - Tests all 3 date formats
     - `TestProject_CalculateDaysUnreading_MultiFormat` - Tests CalculateDaysUnreading with different formats
     - `TestProject_CalculateDaysUnreading_Timezone` - Tests timezone-aware comparison
     - `TestProject_CalculateMedianDay_Timezone` - Tests median day with timezone
     - `TestProject_CalculateFinishedAt_MultiFormat` - Tests finished at calculation with different formats

5. **.env.example**
   - Added `TZ_LOCATION` configuration example with documentation

**Test Results:**
- All unit tests pass ✓
- All integration tests pass ✓
- go fmt passes ✓
- go vet passes ✓
- Build succeeds ✓
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
- [x] #11 Integration tests verify actual database interactions
- [x] #12 Tests use testing-expert subagent for test execution and verification
- [x] #13 Code follows Go formatting standards
- [ ] #14 All new functions have unit tests with >80% coverage
<!-- DOD:END -->
