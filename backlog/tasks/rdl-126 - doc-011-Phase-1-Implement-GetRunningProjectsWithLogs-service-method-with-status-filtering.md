---
id: RDL-126
title: >-
  [doc-011 Phase 1] Implement GetRunningProjectsWithLogs service method with
  status filtering
status: To Do
assignee:
  - thomas
created_date: '2026-04-28 11:16'
updated_date: '2026-04-28 12:08'
labels:
  - feature
  - backend
  - phase-1
dependencies: []
documentation:
  - doc-011
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create internal/service/dashboard/projects_service.go with GetRunningProjectsWithLogs() method that filters projects by calculated 'running' status using 7-day threshold. Implement progress calculation (page/total_page*100) and ordering logic (progress DESC, id ASC).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 GetRunningProjectsWithLogs returns only projects with status='running'
- [ ] #2 Status calculation uses 7-day threshold for running status
- [ ] #3 Progress calculated as (page/total_page)*100
- [ ] #4 Projects ordered by progress DESC, then id ASC
- [ ] #5 Division by zero handled returning 0.0
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The implementation will follow Clean Architecture principles, creating a service layer method that filters projects by their calculated "running" status using a 7-day threshold. The approach:

**Architecture Flow:**
1. **Repository Layer**: Add `GetRunningProjectsWithLogs()` to `DashboardRepository` interface with SQL query that joins projects and logs
2. **Service Layer**: Implement `GetRunningProjectsWithLogs()` in `ProjectsService` that:
   - Calls repository to fetch running projects with eager-loaded logs
   - Calculates progress as `(page / total_page) * 100`
   - Orders results by progress DESC, then id ASC
   - Handles division by zero (returns 0.0)
3. **Domain Layer**: Leverage existing `Project.CalculateStatus()` method which uses config's `EmAndamentoRange` (7 days) for running status determination

**Why this approach:**
- Separation of concerns: Repository handles SQL, Service handles business logic
- Reuses existing status calculation logic from `Project` model
- Follows existing patterns from `DayService` and other dashboard services
- Enables unit testing with mock repositories

**Key Decisions:**
- Status filtering will be done in Go (not SQL) to leverage existing `CalculateStatus()` logic
- Progress calculation matches `Project.CalculateProgress()` formula
- Ordering implemented in Go after fetching data to ensure consistency with calculated fields

---

### 2. Files to Modify

**New Files to Create:**
1. `internal/service/dashboard/projects_service_test.go` - Unit tests for the new service method

**Existing Files to Modify:**
1. `internal/repository/dashboard_repository.go`
   - Add `GetRunningProjectsWithLogs(ctx context.Context) ([]*dto.ProjectWithLogs, error)` to interface

2. `internal/adapter/postgres/dashboard_repository.go`
   - Implement `GetRunningProjectsWithLogs()` with SQL query:
     - JOIN projects and logs
     - Filter by calculated running status (requires fetching all projects and filtering in Go, or using CASE expression)
     - Eager-load first 4 logs per project ordered by data DESC
     - Order by progress DESC, id ASC

3. `internal/service/dashboard/projects_service.go`
   - Add `GetRunningProjectsWithLogs(ctx context.Context) ([]*ProjectWithLogs, error)` method
   - Implement progress calculation: `(page / total_page) * 100`
   - Implement ordering logic (progress DESC, id ASC)
   - Handle edge cases (zero total_page, nil projects)

4. `internal/domain/dto/dashboard_response.go`
   - Verify `ProjectWithLogs` struct has all required fields for the response
   - Add any missing fields for running projects filter

5. `internal/api/v1/handlers/dashboard_handler.go` (Future task RDL-129)
   - Will be updated to use the new service method

6. `internal/api/v1/routes.go` (Future task RDL-125)
   - Route registration handled separately

---

### 3. Dependencies

**Prerequisites:**
1. ✅ Existing `DashboardRepository` interface in place
2. ✅ `ProjectsService` struct already exists with `NewProjectsService()` constructor
3. ✅ `Project` model with `CalculateStatus()` and `CalculateDaysUnreading()` methods
4. ✅ `Config` with `EmAndamentoRange` (7 days) and `DormindoRange` (14 days)
5. ✅ `ProjectAggregateResponse` DTO exists
6. ✅ `ProjectWithLogs` struct exists in service layer

**Blocking Issues:**
- None - all dependencies are in place

**Setup Steps:**
1. Review existing `ProjectsService.GetAll()` method for pattern reference
2. Review `Project.CalculateStatus()` implementation for status logic
3. Review existing repository SQL query patterns in `dashboard_repository.go`

---

### 4. Code Patterns

**Follow Existing Patterns:**

1. **Service Method Pattern** (from `DayService`):
```go
func (s *ProjectsService) GetRunningProjectsWithLogs(ctx context.Context) ([]*ProjectWithLogs, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    // Call repository
    projects, err := s.repo.GetRunningProjectsWithLogs(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get running projects: %w", err)
    }
    
    // Business logic
    // ...
    
    return results, nil
}
```

2. **Error Handling Pattern**:
```go
if err != nil {
    return nil, fmt.Errorf("failed to get running projects: %w", err)
}
```

3. **Progress Calculation Pattern** (from `Project.CalculateProgress()`):
```go
func calculateProgress(page, totalPage int) float64 {
    if totalPage <= 0 {
        return 0.0
    }
    if page <= 0 {
        return 0.0
    }
    progress := (float64(page) / float64(totalPage)) * 100.0
    return math.Round(progress*1000) / 1000 // Round to 3 decimals
}
```

4. **Ordering Pattern** (from `ProjectsService.GetAll()`):
```go
sort.Slice(results, func(i, j int) bool {
    if results[i].Progress != results[j].Progress {
        return results[i].Progress > results[j].Progress // DESC
    }
    return results[i].Project.ProjectID < results[j].Project.ProjectID // ASC
})
```

5. **Context Timeout Pattern**:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

**Naming Conventions:**
- Method names: `GetRunningProjectsWithLogs` (camelCase)
- Variables: `runningProjects`, `progress`, `totalPages` (camelCase)
- Error messages: "failed to get running projects: %w"

---

### 5. Testing Strategy

**Unit Tests** (`internal/service/dashboard/projects_service_test.go`):

1. **Test GetRunningProjectsWithLogs - Normal Case**:
```go
func TestProjectsService_GetRunningProjectsWithLogs(t *testing.T) {
    mockRepo := &MockDashboardRepository{}
    service := NewProjectsService(mockRepo, nil)
    
    // Setup: Mock returns 2 running projects
    mockRepo.mockGetRunningProjectsWithLogs = func(ctx context.Context) ([]*dto.ProjectWithLogs, error) {
        return []*dto.ProjectWithLogs{
            {Project: &dto.ProjectAggregateResponse{ProjectID: 1, Progress: 50.0}},
            {Project: &dto.ProjectAggregateResponse{ProjectID: 2, Progress: 25.0}},
        }, nil
    }
    
    // Execute
    results, err := service.GetRunningProjectsWithLogs(ctx)
    
    // Verify: Ordered by progress DESC
    assert.NoError(t, err)
    assert.Len(t, results, 2)
    assert.Equal(t, 50.0, results[0].Project.Progress)
    assert.Equal(t, 25.0, results[1].Project.Progress)
}
```

2. **Test Progress Calculation - Division by Zero**:
```go
t.Run("division by zero handling", func(t *testing.T) {
    // Setup: Project with total_page = 0
    // Verify: Progress returns 0.0
})
```

3. **Test Ordering - Equal Progress**:
```go
t.Run("equal progress ordering by id", func(t *testing.T) {
    // Setup: Two projects with same progress
    // Verify: Ordered by id ASC
})
```

4. **Test Empty Results**:
```go
t.Run("no running projects", func(t *testing.T) {
    // Setup: Mock returns empty slice
    // Verify: Returns empty slice, no error
})
```

5. **Test Repository Error**:
```go
t.Run("repository error", func(t *testing.T) {
    // Setup: Mock returns error
    // Verify: Returns error with proper wrapping
})
```

**Integration Tests** (`test/integration/dashboard_projects_test.go` - Future task):
- Create test database with fixtures
- Insert projects with different statuses (running, finished, stopped)
- Verify only running projects returned
- Verify progress calculation matches formula
- Verify ordering is correct

**Edge Cases to Cover:**
- Projects with `total_page = 0` (division by zero → progress = 0.0)
- Projects with `page = 0` (progress = 0.0)
- Projects with `page > total_page` (progress clamped to 100.0)
- Empty project list
- Single project
- Float rounding to 3 decimals

---

### 6. Risks and Considerations

**Known Issues:**
1. **Status Filtering Complexity**: The status is calculated dynamically based on `days_unreading`, which requires:
   - Fetching all projects (not just running ones) from database
   - Calculating status for each project in Go
   - Filtering running projects after calculation
   
   **Mitigation**: Implement in service layer where `Project.CalculateStatus()` is available

2. **Performance**: Fetching all projects and filtering in Go may be less efficient than SQL filtering
   
   **Mitigation**: For Phase 1, correctness is prioritized. SQL optimization can be added later if needed.

3. **Float Precision**: Progress calculation must round to 3 decimals consistently
   
   **Mitigation**: Use `math.Round(progress*1000) / 1000` pattern from existing code

**Trade-offs:**
- ✅ **Correctness**: Using Go's `CalculateStatus()` ensures consistency with existing logic
- ⚠️ **Performance**: Filtering in Go vs SQL - acceptable for Phase 1, can optimize later
- ✅ **Testability**: Service layer approach enables easy unit testing with mocks

**Deployment Considerations:**
- No database schema changes required
- No configuration changes required
- Backward compatible - new method doesn't affect existing endpoints

**Rollout Plan:**
1. Implement service method (this task RDL-126)
2. Add repository method (RDL-128)
3. Create unit tests (RDL-130)
4. Update handler (RDL-129)
5. Register route (RDL-125)
6. Integration tests (future)

**Acceptance Criteria Mapping:**
- ✅ #1 GetRunningProjectsWithLogs returns only projects with status='running' → Service filters by `CalculateStatus() == StatusRunning`
- ✅ #2 Status calculation uses 7-day threshold → Uses `config.EmAndamentoRange` (default 7)
- ✅ #3 Progress calculated as (page/total_page)*100 → Matches `Project.CalculateProgress()`
- ✅ #4 Projects ordered by progress DESC, then id ASC → Sort in service layer
- ✅ #5 Division by zero handled returning 0.0 → Check `totalPage <= 0` before division
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
