---
id: RDL-130
title: '[doc-011 Phase 3] Create unit tests for projects service layer'
status: To Do
assignee:
  - thomas
created_date: '2026-04-28 11:17'
updated_date: '2026-04-28 14:20'
labels:
  - testing
  - backend
  - phase-3
dependencies: []
documentation:
  - doc-011
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create internal/service/dashboard/projects_service_test.go with unit tests for GetRunningProjectsWithLogs and CalculateStats methods. Test status filtering logic, stats calculation, progress ordering, and edge cases (zero projects, division by zero). Use mock repository for isolation.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Test status filtering returns only running projects
- [ ] #2 Test stats calculation with known input values
- [ ] #3 Test progress ordering (DESC by progress, ASC by id)
- [ ] #4 Test edge case: zero projects returns empty array
- [ ] #5 Test edge case: division by zero returns 0.0
- [ ] #6 Test coverage > 85% for service layer
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task involves creating comprehensive unit tests for the ProjectsService layer, specifically for the `GetRunningProjectsWithLogs` and `CalculateStats` methods. The tests will use mock repositories to isolate the service layer from database dependencies.

**Current State Analysis:**
- `internal/service/dashboard/projects_service.go` already exists with the required methods
- `internal/service/dashboard/projects_service_test.go` already exists with substantial test coverage
- The existing tests cover most acceptance criteria but need verification for completeness

**Technical Strategy:**
1. **Verify existing test coverage** - Run tests and measure coverage to identify gaps
2. **Enhance test completeness** - Add missing edge cases if coverage is below 85%
3. **Validate test quality** - Ensure tests follow project patterns and best practices
4. **Document findings** - Update task with coverage report and any additional recommendations

**Key Implementation Decisions:**
- Use existing mock infrastructure (`MockDashboardRepositoryForProjects`, `MockPgxPoolForProjects`)
- Follow existing test naming conventions (`TestProjectsService_MethodName_Subject`)
- Use `assert` and `require` from testify package consistently
- Test both happy paths and error scenarios
- Verify float rounding to 3 decimal places for all calculated fields

### 2. Files to Modify

**Files to Review/Verify:**
- `internal/service/dashboard/projects_service.go` - Read to understand implementation details
- `internal/service/dashboard/projects_service_test.go` - Review existing tests for completeness
- `internal/domain/dto/dashboard_response.go` - Understand DTO structures used in tests
- `internal/repository/dashboard_repository.go` - Understand repository interface

**Files to Modify (if needed):**
- `internal/service/dashboard/projects_service_test.go` - Add missing test cases if coverage < 85%
  - Additional edge cases for `GetRunningProjectsWithLogs`
  - Additional edge cases for `CalculateStats`
  - Error path tests for database operations

**Files to Create (if needed):**
- None expected - existing test infrastructure is sufficient

### 3. Dependencies

**Prerequisites:**
- [x] ProjectsService implementation exists (`internal/service/dashboard/projects_service.go`)
- [x] DashboardRepository interface exists (`internal/repository/dashboard_repository.go`)
- [x] DTOs are defined (`internal/domain/dto/dashboard_response.go`)
- [x] Test helper infrastructure exists (`test/test_helper.go`)
- [x] Mock implementations exist in test file

**Blocking Issues:**
- None identified - all dependencies are in place

**Setup Steps:**
1. Verify Go environment is set up (`go version` should show 1.25.7)
2. Ensure test database is accessible (for any integration verification)
3. Run existing tests to establish baseline: `go test -v ./internal/service/dashboard/...`
4. Generate coverage report: `go test -cover ./internal/service/dashboard/...`

### 4. Code Patterns

**Testing Patterns to Follow:**
1. **Test Structure:**
   ```go
   func TestProjectsService_MethodName(t *testing.T) {
       mockRepo := &MockDashboardRepositoryForProjects{}
       mockPool := &MockPgxPoolForProjects{}
       service := NewProjectsService(mockRepo, mockPool)
       ctx := context.Background()

       t.Run("test_case_name", func(t *testing.T) {
           // Setup mock expectations
           // Execute method under test
           // Verify results
       })
   }
   ```

2. **Assertion Patterns:**
   - Use `require.NoError(t, err)` for error checks when error is expected to be nil
   - Use `assert.NoError(t, err)` when continuing after error check
   - Use `assert.Len(t, slice, expected)` for slice length verification
   - Use `assert.InDelta(t, expected, actual, delta)` for float comparisons (delta: 0.001)

3. **Mock Setup Pattern:**
   ```go
   mockRepo.mockGetRunningProjectsWithLogs = func(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
       return []*dto.ProjectWithLogs{/* test data */}, nil
   }
   ```

4. **Naming Conventions:**
   - Test functions: `Test<ServiceName>_<MethodName>_<Subject>`
   - Sub-tests: Descriptive names explaining the scenario
   - Mock types: `Mock<InterfaceName>For<ServiceName>`

5. **Context Usage:**
   - Always use `context.Background()` for unit tests
   - No timeout needed for mock-based tests

### 5. Testing Strategy

**Test Coverage Goals:**
- Target: > 85% line coverage for `projects_service.go`
- Focus areas: All public methods, error paths, edge cases

**Test Categories:**

**A. GetRunningProjectsWithLogs Tests:**
1. **Happy Path Tests:**
   - Normal case: Multiple running projects with valid data
   - Single project case
   - Projects with varying progress levels (verify ordering)

2. **Filtering Tests:**
   - Filter out finished projects (pages >= total_pages)
   - Filter out projects without logs
   - Filter out projects with pages exceeding total

3. **Ordering Tests:**
   - Order by progress DESC
   - Order by id ASC when progress is equal
   - Mixed scenario with both criteria

4. **Edge Case Tests:**
   - Zero projects returns empty array
   - Division by zero (total_pages = 0) returns 0.0 progress
   - Zero pages with valid total returns 0.0 progress
   - Float rounding to 3 decimal places (e.g., 1/3 = 33.333)

5. **Error Path Tests:**
   - Repository returns error
   - Error message includes context ("failed to get running projects")

**B. CalculateStats Tests:**
1. **Happy Path Tests:**
   - Normal case: Multiple projects with known values
   - Single project case
   - Verify stats.total_pages = sum of all total_page values
   - Verify stats.pages = sum of all page values
   - Verify stats.progress_geral = round((pages/total_pages)*100, 3)

2. **Edge Case Tests:**
   - Zero projects returns stats with all values at 0
   - Division by zero (total_pages = 0) returns 0.0 for progress_geral
   - Zero pages with valid total
   - Float rounding to 3 decimal places

3. **Error Path Tests:**
   - Repository returns error
   - Database query for page value returns error
   - Project not found (pgx.ErrNoRows) - should skip and continue

**C. Helper Method Tests:**
1. **calculateProgress Tests:**
   - Normal progress calculation
   - Zero pages
   - Zero total pages (division by zero)
   - Both zero
   - Full progress (100%)
   - Exceeds total (not clamped)
   - Decimal rounding

2. **isRunningProject Tests:**
   - Running project (has logs, not finished)
   - Finished project (pages >= total_pages)
   - Project exceeded (pages > total_pages)
   - No logs (empty slice)
   - Nil logs slice

**Test Execution Commands:**
```bash
# Run all tests with verbose output
go test -v ./internal/service/dashboard/projects_service_test.go

# Run tests with coverage
go test -cover ./internal/service/dashboard/projects_service_test.go

# Run tests with coverage and HTML report
go test -coverprofile=coverage.out ./internal/service/dashboard/...
go tool cover -html=coverage.out

# Check specific coverage percentage
go test -cover ./internal/service/dashboard/projects_service_test.go 2>&1 | grep "coverage:"
```

**Edge Cases to Cover:**
- Empty input (zero projects, empty logs)
- Division by zero scenarios
- Nil pointer handling
- Error propagation
- Float precision and rounding
- Boundary values (0, 1, max int)

### 6. Risks and Considerations

**Known Risks:**
1. **Coverage Gap Risk:** Existing tests may not reach 85% coverage
   - Mitigation: Identify uncovered lines and add targeted test cases
   - Use `go test -coverprofile` to identify gaps

2. **Mock Completeness Risk:** Mock implementations may not cover all scenarios
   - Mitigation: Review mock methods against actual repository interface
   - Add additional mock methods if needed

3. **Float Comparison Risk:** Direct float comparison can fail due to precision
   - Mitigation: Always use `assert.InDelta(t, expected, actual, 0.001)` for floats
   - Document rounding behavior in test comments

**Trade-offs:**
- **Mock vs Integration:** Using mocks provides isolation but may miss integration issues
  - Decision: Unit tests focus on business logic; integration tests (RDL-131, RDL-132) cover database interactions
- **Test Complexity vs Coverage:** Adding too many edge cases can make tests hard to maintain
  - Decision: Focus on meaningful edge cases that could cause real bugs

**Deployment Considerations:**
- Tests must pass before code merge
- Coverage report should be reviewed in PR
- No breaking changes to existing test infrastructure

**Definition of Done for This Task:**
- [ ] All unit tests pass (`go test ./internal/service/dashboard/...`)
- [ ] Test coverage > 85% for `projects_service.go`
- [ ] `go fmt` and `go vet` pass with no errors
- [ ] All acceptance criteria from task description are met
- [ ] Test code follows existing project patterns
- [ ] Error handling tests verify proper error wrapping
- [ ] Edge cases are documented in test comments

**Acceptance Criteria Mapping:**
| AC | Description | Test Coverage |
|----|-------------|---------------|
| #1 | Test status filtering returns only running projects | `TestProjectsService_GetRunningProjectsWithLogs_Filtering` |
| #2 | Test stats calculation with known input values | `TestProjectsService_CalculateStats` |
| #3 | Test progress ordering (DESC by progress, ASC by id) | `TestProjectsService_GetRunningProjectsWithLogs` |
| #4 | Test edge case: zero projects returns empty array | `TestProjectsService_GetRunningProjectsWithLogs` |
| #5 | Test edge case: division by zero returns 0.0 | Both methods have dedicated tests |
| #6 | Test coverage > 85% for service layer | To be verified with coverage report |
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
