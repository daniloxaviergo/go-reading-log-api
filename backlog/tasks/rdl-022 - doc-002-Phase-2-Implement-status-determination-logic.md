---
id: RDL-022
title: '[doc-002 Phase 2] Implement status determination logic'
status: To Do
assignee:
  - workflow
created_date: '2026-04-03 14:02'
updated_date: '2026-04-03 17:22'
labels:
  - phase-2
  - status-logic
  - business-rules
dependencies: []
references:
  - >-
    PRD Section: Technical Decisions - Decision 1: Derived Calculations
    Implementation
  - 'PRD Section: Validation Rules - status values'
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement status determination logic in Go matching Rails ActiveModelSerializer status method. Status depends on days_unreading ranges (configured values) and logs count: unstarted (no logs started), finished (logs count = total_page), running (days_unreading ≤ em_andamento_range), sleeping (days_unreading ≤ dormindo_range), stopped (all other cases).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Unstarted: No logs or no log with data
- [ ] #2 Finished: Logs count equals total_page
- [ ] #3 Running: days_unreading ≤ em_andamento_range
- [ ] #4 Sleeping: days_unreading ≤ dormindo_range
- [ ] #5 Stopped: All other cases
- [ ] #6 Method implemented in Project model or calculations package
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Implement status determination logic in Go that matches the Rails ActiveModelSerializer `status` method. The status is determined by evaluating the following conditions in order:

1. **unstarted**: No logs exist for the project or no log has a data value
2. **finished**: The project's page >= total_page (reading complete)
3. **running**: days_unreading <= em_andamento_range (configurable, default 7 days)
4. **sleeping**: days_unreading <= dormindo_range (configurable, default 14 days)
5. **stopped**: All other cases (days_unreading > dormindo_range but not finished)

**Technical Implementation**:
- Add a `CalculateStatus()` method to the `Project` model in `internal/domain/models/project.go`
- The method will:
  - Accept logs as a parameter (since they're not directly in the model)
  - Check for empty logs → unstarted
  - Check if page >= total_page → finished
  - Calculate days_unreading using the formula: `(current_date - last_log_date_or_started_at)`
  - Compare against configuration ranges to determine running/sleeping/stopped
- Use configuration values from `internal/config/config.go` (em_andamento_range, dormindo_range)

**Why this approach**:
- Matches Rails implementation pattern in `app/models/project.rb`
- Follows Clean Architecture by placing business logic in the domain model
- Uses existing configuration infrastructure
- Returns a string status matching the Rails JSON response format
- Supports multiple status values: unstarted, finished, running, sleeping, stopped

**Key Design Decisions**:
- Calculate status on-demand via method (not stored in model) for accuracy
- Days calculation uses `time.Now()` matching Rails `Date.today` behavior
- Last log date used when available, otherwise started_at date

### 2. Files to Modify

| File | Action | Reason |
|------|--|------|
| `internal/domain/models/project.go` | Modify | Add `CalculateStatus()` method with status determination logic |
| `internal/domain/models/project.go` | Modify | Add `CalculateDaysUnreading()` helper method |
| `internal/domain/models/project_test.go` | Modify | Add unit tests for status determination (5 status types + edge cases) |
| `internal/config/config.go` | Verify | Ensure em_andamento_range (default 7) and dormindo_range (default 14) are set |
| `internal/adapter/postgres/project_repository.go` | Modify | Call `CalculateStatus()` when building ProjectResponse DTOs |
| `internal/api/v1/handlers/projects_handler.go` | Verify | Ensure status field is set from repository responses |
| `test/test_helper.go` | Modify | Update mock repository to support status testing |

### 3. Dependencies

**Prerequisites**:
1. Task RDL-020 completed (progress calculation) - uses similar pattern
2. Configuration values em_andamento_range and dormindo_range must be set (already implemented in config.go)
3. Access to logs through ProjectWithLogs structure in repository layer
4. Date/time handling follows RFC3339 format (RDL-019 in progress)

**Required Existing Infrastructure**:
- `internal/config` package with status range configuration
- `internal/domain/models` package with Project model
- PostgreSQL repository with eager-loaded logs support
- Time package usage consistent with RFC3339

**Missing Components**:
- None - all prerequisites are in place or in progress

### 4. Code Patterns

**Existing conventions to follow**:

1. **Project model pattern** (from `internal/domain/models/project.go`):
```go
type Project struct {
    ctx        context.Context
    ID         int64      `json:"id"`
    Name       string     `json:"name"`
    TotalPage  int        `json:"total_page"`
    StartedAt  *time.Time `json:"started_at"`
    Page       int        `json:"page"`
    Progress   *float64   `json:"progress,omitempty"`
    // ... other fields
}

func (p *Project) CalculateProgress() *float64 {
    // Implementation here
}
```

2. **Status value constants** (to be added):
```go
const (
    StatusUnstarted = "unstarted"
    StatusFinished  = "finished"
    StatusRunning   = "running"
    StatusSleeping  = "sleeping"
    StatusStopped   = "stopped"
)
```

3. **Helper functions pattern**:
```go
func formatDatePtr(t *time.Time) *string {
    if t == nil {
        return nil
    }
    s := t.Format(time.RFC3339)
    return &s
}
```

4. **Calculation timing**: Add methods to Project model, call from repository/DTO layer

**Status determination logic**:
```go
func (p *Project) CalculateStatus(logs []*dto.LogResponse, config *config.Config) *string {
    // 1. Check for unstarted (no logs)
    if len(logs) == 0 {
        return stringPtr(StatusUnstarted)
    }
    
    // 2. Check for finished (page >= total_page)
    if p.Page >= p.TotalPage {
        return stringPtr(StatusFinished)
    }
    
    // 3. Calculate days_unreading
    daysUnreading := p.CalculateDaysUnreading(logs, p.StartedAt)
    
    // 4. Check running (days <= em_andamento_range)
    if daysUnreading <= config.GetEmAndamentoRange() {
        return stringPtr(StatusRunning)
    }
    
    // 5. Check sleeping (days <= dormindo_range)
    if daysUnreading <= config.GetDormindoRange() {
        return stringPtr(StatusSleeping)
    }
    
    // 6. All other cases → stopped
    return stringPtr(StatusStopped)
}
```

### 5. Testing Strategy

**Unit tests to add** (in `internal/domain/models/project_test.go`):

1. `TestProject_CalculateStatus_EmptyLogs` - Verify 'unstarted' status when no logs
2. `TestProject_CalculateStatus_Finished` - Verify 'finished' status when page >= total_page
3. `TestProject_CalculateStatus_Running` - Verify 'running' status when days_unreading <= 7
4. `TestProject_CalculateStatus_Sleeping` - Verify 'sleeping' status when days_unreading <= 14
5. `TestProject_CalculateStatus_Stopped` - Verify 'stopped' status when days_unreading > 14
6. `TestProject_CalculateStatus_EdgeCase_RunningBoundary` - Verify 'running' at exactly 7 days
7. `TestProject_CalculateStatus_EdgeCase_SleepingBoundary` - Verify 'sleeping' at exactly 14 days
8. `TestProject_CalculateStatus_NoStartedAt` - Handle nil started_at gracefully

**Test data setup**:
- Create test projects with varying page/total_page values
- Create log arrays with different last log dates
- Configure test config with status ranges
- Use time manipulation to set specific dates for testing

**Approach**:
- Use `testing.T` for assertions
- Compare returned status string with expected values
- Test all 5 status types and boundary conditions
- Verify configuration values are used correctly

**Integration tests to verify**:
- Full response includes correct status field
- Repository properly passes logs to status calculation
- Handler returns correct JSON response with status

### 6. Risks and Considerations

**Potential pitfalls**:
1. **Time calculation**: Rails uses `Date.today` (date only), Go needs to match this behavior
   - Solution: Use `time.Now().Date()` or extract year/month/day from `time.Now()`

2. **Last log date**: Rails uses `logs.first.data` (most recent log date)
   - Solution: Sort logs by date DESC and use the first one's data field

3. **Starting point for days calculation**: Rails uses `last_read || started_at`
   - Solution: Use last log date if available, otherwise started_at

4. **Null started_at**: Need to handle cases where started_at is nil
   - Solution: Return 'unstarted' or calculate from today (based on logs)

5. **Configuration access**: Need to pass config to the calculation method
   - Solution: Either pass config as parameter or use a constructor pattern

6. **Integer division in days calculation**: Must use float conversion
   - Solution: Use `time.Duration` and convert to days properly

**Design decisions**:
- **Return type**: `*string` matches existing pattern for optional derived fields
- **Status values**: Match Rails exactly (unstarted, finished, running, sleeping, stopped)
- **Configuration passing**: Pass config as parameter for testability and decoupling
- **Days calculation**: Implement as separate method for reusability and testing

**Comparison with Rails**:
- Rails: `days_unreading = (Date.today - base_data).to_i` where base_data = last_read || started_at
- Go: `days = int(time.Now().Sub(baseDate).Hours() / 24)` or use date components

**Edge cases to handle**:
1. No logs, no started_at → return 'unstarted' (cannot calculates days_unreading)
2. Log without data field → treat as no log data
3. Zero total_page → return 'finished' only if page == 0
4. Negative values → return 'unstarted' (invalid data)
5. Far future started_at → days_unreading will be negative → 'unstarted' or 'running'
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
<!-- DOD:END -->
