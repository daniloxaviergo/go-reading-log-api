---
id: RDL-152
title: >-
  [doc-14 Phase 1] Add repository interface methods for weekday-based mean
  calculations
status: To Do
assignee:
  - thomas
created_date: '2026-05-10 10:46'
updated_date: '2026-05-10 11:39'
labels:
  - infrastructure
  - repository
  - phase-1
dependencies: []
documentation:
  - doc-014
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Define new methods in DashboardRepository interface: GetWeekdayPagesGrouped(ctx, startDate, endDate) for fetching pages grouped by weekday, GetFirstLogDate(ctx) for getting the earliest log date, and GetWeekdayMeanWithIntervals(ctx, weekday, currentDate) for calculating mean with 7-day interval logic.

These methods support the weekday-based historical mean calculation algorithm required by the speculate_actual endpoint.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 GetWeekdayPagesGrouped method defined in interface with correct signature
- [ ] #2 GetFirstLogDate method defined in interface
- [ ] #3 GetWeekdayMeanWithIntervals method defined in interface
- [ ] #4 All interface methods documented with comments
- [ ] #5 Interface compiles without errors
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task defines three new methods in the `DashboardRepository` interface to support weekday-based historical mean calculations required by the `speculate_actual` endpoint. The implementation follows Clean Architecture principles and aligns with existing repository patterns.

**Technical Strategy:**

1. **GetWeekdayPagesGrouped** - Fetches pages grouped by weekday within a date range
   - SQL approach: Use `EXTRACT(DOW FROM data::timestamp)` to group by weekday (0-6)
   - Returns aggregated data: weekday, total_pages, log_count for each weekday
   - Supports date range filtering for historical analysis
   - Returns empty map when no data exists (not nil/error)

2. **GetFirstLogDate** - Retrieves the earliest log timestamp in the database
   - Simple `MIN(data::timestamp)` query across all logs
   - Returns `*time.Time` (nil pointer when no logs exist)
   - Used to calculate the starting point for 7-day interval calculations

3. **GetWeekdayMeanWithIntervals** - Calculates mean using 7-day interval logic
   - Algorithm (V1::MeanLog from Rails):
     - Filter logs by target weekday (DOW 0-6)
     - Calculate `total_pages = SUM(end_page - start_page)`
     - Find `begin_data` (first log) and `log_data` (most recent log)
     - Calculate `count_reads = FLOOR((log_data - begin_data) / 7 days)`
     - Return `mean_day = total_pages / count_reads` (rounded to 3 decimals)
   - Returns `*float64` (nil when no data or zero intervals)
   - Matches existing `GetMeanByWeekday` implementation pattern

**Architecture Decisions:**

- **Interface-first design**: Define methods in interface first (this task), implement in PostgreSQL adapter (RDL-153)
- **Nullable return types**: Use pointer types (`*time.Time`, `*float64`) to handle missing data gracefully
- **Consistent error handling**: Return `(nil, nil)` for empty data, `(value, nil)` for success, `(nil, error)` for failures
- **Context timeout**: All methods use 15-second timeout (consistent with existing `dashboardContextTimeout`)

**Why this approach:**

- Follows existing Clean Architecture patterns (interface in `internal/repository/`, implementation in `internal/adapter/postgres/`)
- Aligns with Rails `V1::MeanLog` calculation algorithm
- Enables the service layer (RDL-155) to perform weekday-based historical mean calculations
- Maintains backward compatibility (no changes to existing methods)

**Trade-offs:**

- ✅ Type safety through interface definitions
- ✅ Consistent with existing repository patterns (GetMaxByWeekday, GetMeanByWeekday)
- ✅ Nullable pointers handle edge cases gracefully
- ⚠️ Three separate queries may have performance impact (acceptable for dashboard endpoint)

---

### 2. Files to Modify

| File | Changes | Rationale |
|------|---------|-----------|
| `internal/repository/dashboard_repository.go` | **MODIFY** - Add three new method signatures to `DashboardRepository` interface | Interface definition layer; follows existing pattern for repository methods |

**New Interface Methods to Add:**

```go
// GetWeekdayPagesGrouped returns pages aggregated by weekday within a date range
// Groups logs by weekday (0-6 = Sunday-Saturday) and calculates total pages and log count
// Returns empty map when no data exists (not nil/error)
GetWeekdayPagesGrouped(ctx context.Context, startDate, endDate time.Time) (map[int]dto.WeekdayPages, error)

// GetFirstLogDate returns the earliest log timestamp in the database
// Returns nil pointer when no logs exist
GetFirstLogDate(ctx context.Context) (*time.Time, error)

// GetWeekdayMeanWithIntervals calculates mean pages per 7-day interval for a specific weekday
// Algorithm: total_pages / count_reads where count_reads = floor((log_data - begin_data) / 7 days)
// Returns nil for no data or zero intervals (consistent with GetMeanByWeekday)
GetWeekdayMeanWithIntervals(ctx context.Context, weekday int, currentDate time.Time) (*float64, error)
```

**Files to Create:** None (this is Phase 1 interface definition only)

**Files to Delete:** None

**DTO Dependencies:**

The `GetWeekdayPagesGrouped` method returns `map[int]dto.WeekdayPages`, which requires a new DTO struct. This DTO must be defined in `internal/domain/dto/dashboard_response.go`:

```go
// WeekdayPages represents aggregated page data for a specific weekday
type WeekdayPages struct {
    Weekday   int     `json:"weekday"`   // 0-6 (Sunday-Saturday)
    TotalPages int    `json:"total_pages"`
    LogCount  int     `json:"log_count"`
    Mean      float64 `json:"mean"` // total_pages / log_count
}
```

---

### 3. Dependencies

**Prerequisites:**

- Task RDL-151 (Extend ECharts DTOs) must be completed first
  - Provides `BoundaryGap` field and mark element support needed for chart configuration
  - Ensures DTO infrastructure is ready for new response structures

**Required Before Implementation:**

- None - this task only defines interface methods
- No database schema changes required
- No PostgreSQL implementation needed (RDL-153)

**Sequential Dependencies:**

- This task must be completed BEFORE:
  - RDL-153 (Implement PostgreSQL repository methods) - needs interface definition first
  - RDL-154 (Add mock repository implementations) - needs interface methods to mock
  - RDL-155 (Implement SpeculateService methods) - depends on repository interface

**Parallel Dependencies:**

- Can be implemented in parallel with RDL-151 (DTO extensions) since they operate on different layers
- Interface definition is independent of implementation details

**External Dependencies:**

- None - all changes are internal to the Go codebase

---

### 4. Code Patterns

**Naming Conventions:**

- Interface method names: PascalCase (e.g., `GetWeekdayPagesGrouped`, `GetFirstLogDate`)
- Parameter names: camelCase (e.g., `startDate`, `endDate`, `currentDate`)
- DTO struct names: PascalCase (e.g., `WeekdayPages`)
- JSON tags: snake_case (e.g., `total_pages`, `log_count`, `weekday`)

**Interface Definition Pattern:**

Follow existing repository interface pattern from `dashboard_repository.go`:

```go
// Method signature pattern
MethodName(ctx context.Context, param1 Type1, param2 Type2) (ReturnType, error)

// Documentation comment pattern
// MethodName describes what the method does
// Additional details about algorithm or behavior
// Returns nil/special values for edge cases
MethodName(params) (return, error)
```

**Context Usage:**

All methods must accept `context.Context` as first parameter:

```go
GetWeekdayPagesGrouped(ctx context.Context, startDate, endDate time.Time) (map[int]dto.WeekdayPages, error)
```

**Error Handling Pattern:**

- Return `(nil, nil)` for empty data (no error)
- Return `(value, nil)` for successful queries
- Return `(nil, error)` for database errors

Example:
```go
// For GetFirstLogDate
if err == pgx.ErrNoRows {
    return nil, nil // No logs exist, not an error
}
if err != nil {
    return nil, fmt.Errorf("failed to get first log date: %w", err)
}
return &firstLogDate, nil
```

**DTO Definition Pattern:**

Follow existing DTO pattern from `dashboard_response.go`:

```go
// WeekdayPages represents aggregated page data for a specific weekday
type WeekdayPages struct {
    ctx        context.Context
    Weekday    int     `json:"weekday"`
    TotalPages int     `json:"total_pages"`
    LogCount   int     `json:"log_count"`
    Mean       float64 `json:"mean"`
}

// NewWeekdayPages creates a new WeekdayPages instance
func NewWeekdayPages(weekday int, totalPages, logCount int, mean float64) *WeekdayPages {
    return &WeekdayPages{
        Weekday:    weekday,
        TotalPages: totalPages,
        LogCount:   logCount,
        Mean:       mean,
    }
}

// Validate validates the WeekdayPages struct
func (w *WeekdayPages) Validate() error {
    if w == nil {
        return fmt.Errorf("weekday pages is nil")
    }
    if w.Weekday < 0 || w.Weekday > 6 {
        return fmt.Errorf("weekday must be between 0 and 6")
    }
    if w.LogCount < 0 {
        return fmt.Errorf("log count cannot be negative")
    }
    return nil
}
```

**Integration Patterns:**

- **Interface Layer**: `internal/repository/dashboard_repository.go` - method signatures only
- **DTO Layer**: `internal/domain/dto/dashboard_response.go` - WeekdayPages struct
- **Implementation Layer**: `internal/adapter/postgres/dashboard_repository.go` - RDL-153 task
- **Mock Layer**: `test/testutil/mock_dashboard_repository.go` - RDL-154 task

---

### 5. Testing Strategy

**Unit Tests:**

Since this task only defines interface methods (no implementation), unit testing focuses on:

1. **Interface Compilation Test** - Verify interface methods compile without errors
   - **File**: `test/unit/repository/dashboard_repository_interface_test.go` (create new file)
   - **Test**: `TestDashboardRepository_InterfaceCompilation` - Ensure interface satisfies Go compilation

2. **DTO Validation Tests** - Test WeekdayPages DTO validation
   - **File**: `test/unit/domain/dto/dashboard_response_test.go` (add to existing file)
   - **Tests**:
     - `TestWeekdayPages_StructDefinition` - Verify struct fields and JSON tags
     - `TestWeekdayPages_NewWeekdayPages` - Test constructor
     - `TestWeekdayPages_Validate_Valid` - Test validation with valid data
     - `TestWeekdayPages_Validate_Invalid` - Test validation edge cases
       - Weekday < 0 or > 6
       - Negative log count
       - Nil struct
     - `TestWeekdayPages_JSONMarshaling` - Verify JSON serialization

**Edge Cases to Cover:**

- `GetWeekdayPagesGrouped`: Empty date range, no logs in range, all weekdays present, some weekdays missing
- `GetFirstLogDate`: Empty database (no logs), single log entry, multiple logs
- `GetWeekdayMeanWithIntervals`: No logs for weekday, zero intervals (logs within same 7-day period), single log entry

**Integration Tests:**

Not required for this Phase 1 task (interface definition only). Integration testing will occur in:
- RDL-153 (PostgreSQL implementation)
- RDL-154 (Mock implementation)

**Testing Approach:**

- Follow existing unit test patterns in `test/unit/repository/` and `test/unit/domain/dto/`
- Use `t.Parallel()` for independent tests
- Test DTO validation logic thoroughly
- Verify interface compiles without errors

**Example Test Case:**

```go
func TestWeekdayPages_Validate_Valid(t *testing.T) {
    t.Parallel()
    
    wp := dto.NewWeekdayPages(2, 100, 5, 20.0)
    err := wp.Validate()
    
    require.NoError(t, err)
}

func TestWeekdayPages_Validate_Invalid_Weekday(t *testing.T) {
    t.Parallel()
    
    wp := dto.NewWeekdayPages(7, 100, 5, 20.0) // Invalid weekday
    err := wp.Validate()
    
    require.Error(t, err)
    assert.Contains(t, err.Error(), "weekday must be between 0 and 6")
}
```

**Test Coverage Goals:**

- 100% line coverage for new DTO code
- 100% branch coverage for validation logic
- All edge cases covered

---

### 6. Risks and Considerations

**Known Issues:**

- None identified - this is interface definition work only

**Potential Pitfalls:**

1. **DTO Dependency**: `GetWeekdayPagesGrouped` returns `map[int]dto.WeekdayPages`, requiring the DTO to be defined first
   - **Mitigation**: Define `WeekdayPages` struct in `internal/domain/dto/dashboard_response.go` before adding interface method
   - **Order**: 1) Add DTO, 2) Add interface method

2. **Method Signature Consistency**: Must match existing repository method patterns
   - **Mitigation**: Review existing methods (GetMeanByWeekday, GetMaxByWeekday) for consistency
   - **Check**: Context as first parameter, error as last return value, pointer types for nullable returns

3. **Documentation Clarity**: Interface methods must be well-documented for implementers
   - **Mitigation**: Follow existing documentation pattern with algorithm details
   - **Include**: Parameter descriptions, return value behavior, edge case handling

**Design Decisions:**

- **Return Type for GetWeekdayPagesGrouped**: `map[int]dto.WeekdayPages` vs `[]dto.WeekdayPages`
  - **Decision**: `map[int]dto.WeekdayPages` for O(1) weekday lookup
  - **Rationale**: Service layer needs to access specific weekdays by index (0-6)
  
- **Nullable Return Types**: Use pointer types (`*time.Time`, `*float64`) for optional values
  - **Rationale**: Consistent with existing methods (GetMaxByWeekday, GetOverallMean)
  - **Behavior**: Return `nil` when no data exists, not error

**Deployment Considerations:**

- No deployment impact - interface changes are backward compatible
- No database migrations required
- No API contract changes (this is infrastructure layer only)

**Rollout Strategy:**

- Phase 1 (this task): Interface definition only
- Phase 2 (RDL-153): PostgreSQL implementation
- Phase 3 (RDL-154): Mock implementation for testing
- Phase 4 (RDL-155): Service layer usage

**Code Review Checklist:**

- [ ] Interface method signatures follow existing patterns
- [ ] Documentation comments are clear and complete
- [ ] DTO struct defined with correct JSON tags
- [ ] Context parameter is first in all method signatures
- [ ] Error is last return value in all method signatures
- [ ] Nullable return types use pointers (`*time.Time`, `*float64`)
- [ ] Unit tests cover DTO validation edge cases
- [ ] Code compiles without errors
- [ ] `go fmt` passes with no errors
- [ ] `go vet` passes with no errors

**Performance Impact:**

- Negligible - interface definition only, no runtime impact
- Actual performance determined by PostgreSQL implementation (RDL-153)

**Related Tasks:**

- RDL-151: DTO extensions (prerequisite - must complete first)
- RDL-153: PostgreSQL implementation (depends on this task)
- RDL-154: Mock implementation (depends on this task)
- RDL-155: Service layer methods (depends on this task)

---

**Implementation Estimate:** 1-2 hours
- DTO definition (WeekdayPages struct): 30 minutes
- Interface method definitions: 30 minutes
- Unit tests for DTO: 30 minutes
- Code review and refinement: 15 minutes

**Ready for Implementation:** ✅ Yes

**Notes:**
- This task is Phase 1 of doc-014 (Dashboard Enhancements)
- Interface-first approach enables parallel development of implementation and mocks
- DTO definition is critical dependency for interface compilation
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress

### Completed ✅

#### 1. DTO Definition (WeekdayPages struct)
- **File**: `internal/domain/dto/dashboard_response.go`
- Added `WeekdayPages` struct with fields:
  - `Weekday int` (0-6, Sunday-Saturday)
  - `TotalPages int`
  - `LogCount int`
  - `Mean float64` (rounded to 3 decimals)
- Added constructor: `NewWeekdayPages(weekday, totalPages, logCount, mean)`
- Added builder methods: `SetWeekday()`, `SetTotalPages()`, `SetLogCount()`, `SetMean()`
- Added `Validate()` method with comprehensive validation
- Added context embedding with `GetContext()` and `SetContext()` methods

#### 2. Interface Method Definitions
- **File**: `internal/repository/dashboard_repository.go`
- Added three new interface methods:
  1. `GetWeekdayPagesGrouped(ctx, startDate, endDate) (map[int]dto.WeekdayPages, error)`
  2. `GetFirstLogDate(ctx) (*time.Time, error)`
  3. `GetWeekdayMeanWithIntervals(ctx, weekday, currentDate) (*float64, error)`
- All methods properly documented with comments explaining algorithm and edge cases

#### 3. PostgreSQL Implementation (Stubs)
- **File**: `internal/adapter/postgres/dashboard_repository.go`
- Implemented all three methods with full SQL queries:
  - `GetWeekdayPagesGrouped`: Groups logs by weekday, calculates totals and means
  - `GetFirstLogDate`: Returns MIN(data::timestamp) from logs table
  - `GetWeekdayMeanWithIntervals`: Implements V1::MeanLog algorithm with 7-day intervals
- All methods use 15-second context timeout
- Proper error handling with `(nil, nil)` for empty data

#### 4. Mock Implementations
- Updated all mock repositories to implement new interface methods:
  - `test/testutil/mock_dashboard_repository.go` - Full mock with testdouble
  - `internal/service/dashboard/projects_service_test.go` - MockDashboardRepositoryForProjects
  - `internal/api/v1/handlers/dashboard_handler_test.go` - MockDashboardRepository
  - `test/unit/day_service_test.go` - MockDashboardRepository
  - `internal/api/v1/routes_test.go` - MockDashboardRepository
  - `test/unit/weekday_faults_service_test.go` - MockDashboardRepositoryWeekdayFaults

#### 5. Unit Tests
- **File**: `test/unit/domain/dto/weekday_pages_test.go`
- Created 33 comprehensive test cases covering:
  - Struct definition and field verification
  - Constructor with valid/zero values
  - Context embedding
  - Builder methods with chaining
  - Validation (valid data, all edge cases)
  - JSON marshaling/unmarshaling
  - Map usage integration

### Verification
- ✅ Code compiles without errors
- ✅ `go fmt` passes with no errors
- ✅ `go vet` passes for modified packages
- ✅ All WeekdayPages unit tests pass (33 tests)
- ✅ Interface compiles correctly
- ✅ All mock implementations updated

### Next Steps
- Task is ready for acceptance criteria verification
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass
- [ ] #2 All integration tests pass execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md and AGENTS.md
- [ ] #8 New code paths include error path tests
- [ ] #9 HTTP handlers test both success and error responses
- [ ] #10 Integration tests verify actual database interactions
<!-- DOD:END -->
