---
id: RDL-082
title: '[doc-008 Phase 2] Implement DayService with weekly statistics calculations'
status: To Do
assignee:
  - workflow
created_date: '2026-04-21 15:50'
updated_date: '2026-04-21 19:40'
labels:
  - phase-2
  - service
  - calculation
dependencies: []
references:
  - REQ-DASH-005
  - AC-DASH-001
  - Implementation Checklist Phase 2
documentation:
  - doc-008
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement internal/service/dashboard/day_service.go calculating previous_week_pages, last_week_pages, per_pages ratio, mean_day, and spec_mean_day. Use GetToday() for consistent date references and ensure all float values rounded to 3 decimal places.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Previous week page total calculated correctly
- [ ] #2 Last week page total calculated correctly
- [ ] #3 Per pages ratio computed with 3 decimal precision
- [ ] #4 Mean day by weekday calculated accurately
- [ ] #5 Speculative mean derived from actual mean
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The DayService needs to calculate weekly statistics for dashboard endpoints. Based on the PRD and existing codebase, I'll implement:

**Architecture Decision**: Create a dedicated service layer following Clean Architecture principles:
- `internal/service/dashboard/day_service.go` - Business logic for daily/weekly calculations
- Use existing `DashboardRepository` interface for data access
- Leverage existing `UserConfigService` for configuration values
- Follow the established pattern in `user_config_service.go`

**Calculation Logic**:
1. **Previous Week Pages**: Sum of pages from 14-7 days ago (previous week)
2. **Last Week Pages**: Sum of pages from 7 days ago to today (current week so far)
3. **Per Pages Ratio**: `last_week_pages / previous_week_pages * 100` (percentage, 3 decimals)
4. **Mean Day**: Average pages per day for the current weekday across all weeks
5. **Spec Mean Day**: `mean_day * (1 + prediction_pct)` from config

**Date Reference**: Use a `GetToday()` helper function for consistent date calculations, ensuring all calculations use the same reference point.

**Key Trade-offs**:
- Will calculate means across ALL available data (not limited to specific weeks) for accuracy
- Zero division handled explicitly (returns 0.0 instead of NaN/error)
- All float values rounded to 3 decimal places as per AC-DASH-001

---

### 2. Files to Modify

#### New Files to Create:
| File | Purpose |
|------|---------|
| `internal/service/dashboard/day_service.go` | Main service implementing weekly statistics calculations |
| `test/unit/day_service_test.go` | Unit tests for DayService calculations |

#### Existing Files to Modify:
| File | Modification |
|------|--------------|
| `internal/api/v1/handlers/dashboard_handler.go` | Update handlers to use DayService instead of inline calculations |
| `internal/api/v1/routes.go` | Add new route if needed (may be handled by existing routes) |

---

### 3. Dependencies

**Prerequisites**:
- ✅ `internal/service/user_config_service.go` - Already exists, provides config access
- ✅ `internal/repository/dashboard_repository.go` - Already exists, provides data queries
- ✅ `internal/domain/dto/dashboard_response.go` - Already exists, provides response structures

**Blocking Issues**:
- None identified - all dependencies are already in place

**Setup Steps**:
1. Verify database has sufficient log data for meaningful calculations
2. Ensure `GetToday()` helper is available (create if not present)
3. Configure `.env` with `PREDICTION_PCT` value (default 0.15)

---

### 4. Code Patterns

**Follow Existing Patterns**:

1. **Service Layer Pattern** (from `user_config_service.go`):
```go
type DayService struct {
    repo       repository.DashboardRepository
    userConfig *service.UserConfigService
}

func NewDayService(repo repository.DashboardRepository, userConfig *service.UserConfigService) *DayService {
    return &DayService{
        repo:       repo,
        userConfig: userConfig,
    }
}
```

2. **Calculation Method Pattern**:
```go
func (s *DayService) CalculateWeeklyStats(ctx context.Context, referenceDate time.Time) (*dto.WeeklyStats, error) {
    // Use GetToday() for consistent date references
    today := GetToday()
    
    // Calculate date ranges
    prevWeekStart := today.AddDate(0, 0, -14)
    prevWeekEnd := today.AddDate(0, 0, -7)
    lastWeekStart := today.AddDate(0, 0, -7)
    lastWeekEnd := today
    
    // ... calculation logic
}
```

3. **Error Handling Pattern** (from existing handlers):
```go
if err != nil {
    return nil, fmt.Errorf("failed to calculate weekly stats: %w", err)
}
```

4. **Float Rounding Pattern**:
```go
func roundToThreeDecimals(val float64) float64 {
    return math.Round(val*1000) / 1000
}
```

---

### 5. Testing Strategy

**Unit Tests** (`test/unit/day_service_test.go`):
- Test `CalculateWeeklyStats()` with mock repository data
- Verify previous_week_pages calculation
- Verify last_week_pages calculation
- Verify per_pages ratio (including edge case of zero previous)
- Verify mean_day calculation across weekdays
- Verify spec_mean_day derivation
- Test with empty data (should return zeros, not errors)
- Test with partial week data

**Integration Tests**:
- Extend existing `test/integration/projects_integration_test.go` pattern
- Create test fixtures with known page counts
- Verify calculations against expected values
- Test date boundary conditions

**Test Coverage Goals**:
- Unit tests: 100% coverage of calculation logic
- Integration tests: Cover real database queries with known data

---

### 6. Risks and Considerations

**Known Risks**:

1. **Data Sparsity**: If logs are sparse or inconsistent, mean calculations may be skewed
   - *Mitigation*: Document that mean is calculated from all available data, not just complete weeks

2. **Zero Division**: Previous week could have zero pages
   - *Mitigation*: Return 0.0 for ratio when denominator is zero (matches Rails behavior)

3. **Time Zone Sensitivity**: "Today" depends on server timezone
   - *Mitigation*: Use `GetToday()` consistently throughout; document TZ dependency

4. **Performance**: Calculating means across large datasets could be slow
   - *Mitigation*: Repository queries should use proper indexes; consider caching if needed in future

5. **AC-DASH-001 Compliance**: Must match Rails exact behavior
   - *Verification*: Compare output with Rails `/v1/dashboard/day.json` endpoint

**Design Decisions**:

| Decision | Rationale |
|----------|-----------|
| Mean calculated from ALL logs, not just recent | Matches Rails `median_day` behavior which uses all available data |
| Speculative mean uses prediction_pct config | Aligns with existing `UserConfigService` pattern |
| Returns zero values for missing data | Graceful degradation; avoids errors |
| 3 decimal precision for all floats | Explicitly required by AC-DASH-001 |

---

### Implementation Checklist

Before coding, verify:
- [ ] PRD section 5.2 (Technical Decisions) reviewed
- [ ] Existing `user_config_service.go` patterns understood
- [ ] `dashboard_repository.go` interface methods available
- [ ] Test database populated with sample data for verification
- [ ] Rails API `/v1/dashboard/day.json` behavior documented for comparison

Ready to implement when approved.
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
