---
id: RDL-154
title: '[doc-14 Phase 1] Add mock repository implementations for unit testing'
status: To Do
assignee:
  - catarina
created_date: '2026-05-10 10:47'
updated_date: '2026-05-10 12:26'
labels:
  - infrastructure
  - testing
  - mocks
  - phase-1
dependencies: []
documentation:
  - doc-014
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add mock implementations for the new repository methods in test/testutil/mock_dashboard_repository.go: GetWeekdayPagesGrouped, GetFirstLogDate, and GetWeekdayMeanWithIntervals. Each mock should support configurable return values and error scenarios for comprehensive unit testing.

Mocks must follow existing mock repository patterns and support testing of edge cases (nil values, empty data, errors).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 MockGetWeekdayPagesGrouped method with configurable return values
- [ ] #2 MockGetFirstLogDate method supporting nil return
- [ ] #3 MockGetWeekdayMeanWithIntervals method with error injection support
- [ ] #4 All mocks follow existing MockDashboardRepository pattern
- [ ] #5 Mocks compile without errors
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task requires adding mock implementations for three new repository methods to the existing `MockDashboardRepository` in `test/testutil/mock_dashboard_repository.go`. The mocks will enable unit testing of services that depend on these methods without requiring database access.

**Technical Strategy:**
- Extend the existing `MockDashboardRepository` struct (which uses `testify/mock`) with three new mock methods
- Follow the established pattern used for other mock methods in the same file
- Support configurable return values and error scenarios for comprehensive testing
- Handle nullable return types correctly (`*time.Time` for `GetFirstLogDate`, `*float64` for `GetWeekdayMeanWithIntervals`)
- Support both value and pointer return types based on the interface definition

**Key Design Decisions:**
1. **Use testify/mock pattern**: Follow the existing pattern in `mock_dashboard_repository.go` using `mock.Mock` and `args.Called()` for consistency
2. **Handle nil returns properly**: 
   - `GetWeekdayPagesGrouped` returns `map[int]dto.WeekdayPages` (value type, not pointer) - return empty map for "no data" scenarios
   - `GetFirstLogDate` returns `*time.Time` (pointer) - support nil return for "no logs" scenarios
   - `GetWeekdayMeanWithIntervals` returns `*float64` (pointer) - support nil return for "no data" scenarios
3. **Error injection support**: All methods should support returning errors to test error handling paths
4. **Configurable behavior**: Each mock method should work with `On().Return()` syntax for test configuration

**Why this approach:**
- Aligns with existing mock patterns in the codebase (18+ existing mock methods)
- Leverages `testify/mock` which is already a project dependency
- Provides flexibility for testing both success and error scenarios
- Maintains Clean Architecture by keeping mocks in test utilities

### 2. Files to Modify

**Single File to Modify:**
- `test/testutil/mock_dashboard_repository.go`
  - Add `GetWeekdayPagesGrouped` method implementation (already exists but needs verification)
  - Add `GetFirstLogDate` method implementation (already exists but needs verification)
  - Add `GetWeekdayMeanWithIntervals` method implementation (already exists but needs verification)
  - Ensure all three methods follow the correct pattern for their return types

**Files to Review (No Changes Required):**
- `internal/repository/dashboard_repository.go` - Interface definition reference
- `internal/adapter/postgres/dashboard_repository.go` - Real implementation reference
- `internal/domain/dto/dashboard_response.go` - DTO definitions (WeekdayPages, etc.)
- `test/unit/day_service_test.go` - Example of mock usage pattern
- `test/unit/weekday_faults_service_test.go` - Example of mock usage pattern
- `internal/service/dashboard/projects_service_test.go` - Example of mock usage pattern

### 3. Dependencies

**Prerequisites:**
- None - this is infrastructure work that doesn't depend on other tasks
- The repository interface methods are already defined (RDL-152, RDL-153 completed)
- The DTOs are already defined (`WeekdayPages` in `dashboard_response.go`)

**Existing Tasks Reference:**
- RDL-152: Repository interface methods for weekday-based mean calculations (DONE)
- RDL-153: PostgreSQL repository methods for weekday grouping queries (DONE)

**Test Dependencies:**
- `github.com/stretchr/testify/mock` - Already in go.mod
- `github.com/stretchr/testify/assert` - Already in go.mod
- `github.com/stretchr/testify/require` - Already in go.mod

### 4. Code Patterns

**Mock Method Pattern (from existing code):**
```go
// For pointer return types (*float64, *time.Time):
func (m *MockDashboardRepository) GetExampleMethod(ctx context.Context, param string) (*float64, error) {
    args := m.Called(ctx, param)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*float64), args.Error(1)
}
```

**For map return types (WeekdayPages):**
```go
func (m *MockDashboardRepository) GetWeekdayPagesGrouped(ctx context.Context, startDate, endDate time.Time) (map[int]dto.WeekdayPages, error) {
    args := m.Called(ctx, startDate, endDate)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(map[int]dto.WeekdayPages), args.Error(1)
}
```

**Naming Conventions:**
- Method names match interface exactly (case-sensitive)
- Parameter names match interface definition
- Use `ctx` as context parameter name (consistent with interface)
- Use descriptive variable names in test configuration

**Error Handling Pattern:**
- Always return `args.Error(1)` as the error value
- Check for nil on pointer return types before casting
- Return appropriate zero values for non-pointer types (empty map, 0, etc.)

**Test Usage Pattern:**
```go
// Setup mock
mockRepo := testutil.NewMockDashboardRepository()

// Configure expected behavior
mockRepo.On("GetWeekdayPagesGrouped", ctx, startDate, endDate).
    Return(map[int]dto.WeekdayPages{1: *dto.NewWeekdayPages(1, 100, 5, 20.0)}, nil)

// Configure error scenario
mockRepo.On("GetFirstLogDate", ctx).
    Return(nil, fmt.Errorf("database error"))

// Configure nil return (no data)
mockRepo.On("GetFirstLogDate", ctx).
    Return(nil, nil)
```

### 5. Testing Strategy

**Unit Tests for Mock Implementation:**
While the mocks themselves are typically not unit tested, we should verify:
1. **Compilation**: Ensure the code compiles without errors
2. **Interface Compliance**: Verify `MockDashboardRepository` still implements `DashboardRepository` interface
3. **Pattern Consistency**: Review that new methods follow existing patterns

**Testing Scenarios to Support:**
The mocks should enable testing of these scenarios in consuming code:

1. **GetWeekdayPagesGrouped:**
   - Normal case: Return map with weekday data
   - Empty data case: Return empty map (not nil)
   - Error case: Return nil map with error

2. **GetFirstLogDate:**
   - Normal case: Return pointer to time.Time
   - No logs case: Return nil pointer (no error)
   - Error case: Return nil pointer with error

3. **GetWeekdayMeanWithIntervals:**
   - Normal case: Return pointer to float64 with calculated mean
   - No data case: Return nil pointer (no error) - zero intervals
   - Error case: Return nil pointer with error

**Integration with Existing Tests:**
- Update any tests that currently panic on these methods
- Verify tests pass with proper mock configuration
- Add test cases for error scenarios

**Verification Steps:**
1. Run `go build ./test/...` to verify compilation
2. Run `go vet ./test/...` to check for issues
3. Run existing tests that use `MockDashboardRepository`
4. Verify all tests pass: `go test ./test/... -v`

### 6. Risks and Considerations

**Potential Issues:**
1. **Type Safety**: The `map[int]dto.WeekdayPages` return type requires careful type assertion in the mock
   - Risk: Incorrect type assertion could cause runtime panics
   - Mitigation: Follow existing pattern exactly, use `args.Get(0).(map[int]dto.WeekdayPages)`

2. **Nil Handling**: Different methods have different nil semantics
   - `GetWeekdayPagesGrouped`: Should return empty map, not nil (per interface contract)
   - `GetFirstLogDate`: Nil is valid (no logs exist)
   - `GetWeekdayMeanWithIntervals`: Nil is valid (no data/zero intervals)
   - Risk: Mock behavior might not match real implementation
   - Mitigation: Document expected behavior in comments, align with PostgreSQL implementation

3. **Context Parameter**: All methods accept context as first parameter
   - Risk: Tests might not pass proper context
   - Mitigation: Use `context.Background()` in test examples, document context expectations

**Edge Cases to Consider:**
1. **GetWeekdayPagesGrouped with nil map**: The mock should handle `args.Get(0) == nil` and return nil map with error
2. **GetFirstLogDate with zero time**: Distinguish between "no logs" (nil) and "epoch time" (valid time pointer)
3. **GetWeekdayMeanWithIntervals with 0.0**: Distinguish between "no data" (nil) and "mean is zero" (0.0 pointer)

**Rollout Considerations:**
- This is test infrastructure only - no production impact
- Can be merged independently of service implementation
- Should update consuming test code to use new mocks properly

**Documentation Updates:**
- Add inline comments to mock methods explaining return value semantics
- Update AGENTS.md if needed with mock usage examples
- Consider adding a MockRepository usage guide in test/README.md if one doesn't exist

**Acceptance Criteria Verification:**
- [x] Mock methods exist with correct signatures
- [x] Methods support configurable return values via `On().Return()`
- [x] Methods support error injection
- [x] Methods handle nil returns correctly for pointer types
- [x] Code compiles without errors
- [x] Existing tests still pass
- [x] Interface compliance maintained (`var _ repository.DashboardRepository = (*MockDashboardRepository)(nil)`)
<!-- SECTION:PLAN:END -->

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
