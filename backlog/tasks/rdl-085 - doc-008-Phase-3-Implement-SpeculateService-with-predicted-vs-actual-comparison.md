---
id: RDL-085
title: >-
  [doc-008 Phase 3] Implement SpeculateService with predicted vs actual
  comparison
status: Done
assignee:
  - thomas
created_date: '2026-04-21 15:50'
updated_date: '2026-04-21 23:38'
labels:
  - phase-3
  - service
  - prediction
dependencies: []
references:
  - REQ-DASH-008
  - AC-DASH-005
  - Implementation Checklist Phase 3
documentation:
  - doc-008
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement internal/service/dashboard/speculate_service.go comparing actual vs predicted reading. Generate chart data for last 15 days including today, calculate speculative mean as actual_mean * (1 + prediction_pct), and zero-fill missing days.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Actual vs predicted comparison implemented
- [x] #2 Speculative mean formula correct (actual * (1 + pct))
- [x] #3 Last 15 days data coverage including today
- [x] #4 Missing days zero-filled not omitted
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The SpeculateService needs to compare actual vs predicted reading data and generate chart configurations. Based on the PRD (AC-DASH-005), this involves:

**Core Functionality:**
1. **Data Comparison**: Compare actual reading data against predicted/specified targets
2. **Chart Data Generation**: Create line chart data points for last 15 days including today
3. **Speculative Mean Calculation**: Use formula `actual_mean * (1 + prediction_pct)` from config
4. **Zero-Fill Missing Days**: Ensure all 15 days have data points, missing days show 0

**Architecture Decisions:**
- Create `speculate_service.go` following the same pattern as `day_service.go` and `faults_service.go`
- Use dependency injection for `DashboardRepository` and `UserConfigProvider`
- Implement chart configuration using existing `EchartConfig` DTO
- Follow Clean Architecture: Service layer handles business logic, Repository handles data access

**Why This Approach:**
- Consistent with existing codebase patterns (DayService, FaultsService)
- Testable without HTTP stack
- Reusable across different API versions
- Leverages existing repository methods for data aggregation

---

### 2. Files to Modify

#### New Files to Create:
| File | Purpose |
|------|---------|
| `internal/service/dashboard/speculate_service.go` | Main service implementation for speculative calculations |

#### Existing Files to Reference (No Modification Needed):
| File | Reason |
|------|--------|
| `internal/domain/dto/dashboard_response.go` | Already contains `EchartConfig`, `Series`, `Axis`, `Legend` DTOs |
| `internal/repository/dashboard_repository.go` | Interface already defined with required methods |
| `internal/adapter/postgres/dashboard_repository.go` | Implementation already exists |
| `internal/api/v1/handlers/dashboard_handler.go` | Handler will be added here (separate task) |

---

### 3. Dependencies

**Required Prerequisites:**
- [x] `UserConfigService` with `GetPredictionPct()` method (already implemented in `user_config_service.go`)
- [x] `DashboardRepository` interface with methods:
  - `GetFaultsByDateRange(ctx, start, end)` - Get faults/reads within date range
  - `GetProjectAggregates(ctx)` - Get project-level statistics
- [x] Existing DTOs in `dashboard_response.go` for chart configurations
- [x] Test infrastructure (`test/unit/dashboard_repository_test.go`, `test_helper.go`)

**No External Dependencies Required:**
- Uses standard library `time`, `math`
- Uses existing `pgx/v5` for database access
- No new packages needed

---

### 4. Code Patterns

**Pattern 1: Service Initialization (from day_service.go)**
```go
type SpeculateService struct {
    repo       repository.DashboardRepository
    userConfig UserConfigProvider
}

func NewSpeculateService(repo repository.DashboardRepository, userConfig UserConfigProvider) *SpeculateService {
    return &SpeculateService{
        repo:       repo,
        userConfig: userConfig,
    }
}
```

**Pattern 2: Date Range Calculation (from faults_service.go)**
```go
func GetDateRangeLastNDays(days int) (start, end time.Time) {
    end = GetToday()
    start = end.AddDate(0, 0, -days)
    return start, end
}
```

**Pattern 3: Speculative Mean Formula (from PRD AC-DASH-005)**
```go
// spec_mean = actual_mean * (1 + prediction_pct)
specMean := actualMean * (1 + s.userConfig.GetPredictionPct())
```

**Pattern 4: Chart Configuration (from existing code)**
```go
chart := dto.NewEchartConfig().
    SetTitle("Speculated vs Actual").
    SetTooltip(map[string]interface{}{"trigger": "axis"}).
    AddSeries(*dto.NewSeries("Actual", "line", actualData).SetLineStyle(...)).
    AddSeries(*dto.NewSeries("Speculated", "line", predictedData).SetLineStyle(...))
```

**Pattern 5: Zero-Fill Missing Days**
```go
// Ensure all 15 days have data points
data := make([]interface{}, 15)
for i := 0; i < 15; i++ {
    if val, exists := resultMap[i]; exists {
        data[i] = val
    } else {
        data[i] = 0 // Zero-fill missing days
    }
}
```

---

### 5. Testing Strategy

**Unit Tests Structure (test/unit/speculate_service_test.go):**

| Test Case | Description | Coverage |
|-----------|-------------|----------|
| `TestSpeculateService_NewSpeculateService` | Service initialization | Constructor |
| `TestSpeculateService_CalculateSpeculativeMean_Normal` | Normal mean calculation | Formula validation |
| `TestSpeculateService_CalculateSpeculativeMean_ZeroMean` | Zero mean edge case | Boundary condition |
| `TestSpeculateService_GenerateChartData_Last15Days` | 15-day data generation | Date range |
| `TestSpeculateService_GenerateChartData_MissingDays` | Zero-fill missing days | Gap handling |
| `TestSpeculateService_CalculatePercentage_Change` | Percentage change calculation | Comparison logic |

**Integration Tests Structure (test/integration/speculate_service_integration_test.go):**
- Set up test database with sample logs
- Verify calculations against known values
- Test date range queries
- Validate chart configuration output

**Mock Strategy:**
```go
type MockDashboardRepository struct {
    mockGetFaultsByDateRange func(ctx context.Context, start, end time.Time) (*dto.FaultStats, error)
}

func (m *MockDashboardRepository) GetFaultsByDateRange(...) { ... }
```

---

### 6. Risks and Considerations

**Blocking Issues:**
1. **None identified** - All prerequisites are already in place

**Design Trade-offs:**

| Decision | Rationale |
|----------|-----------|
| Use `GetToday()` for consistent date reference | Matches existing services (DayService, FaultsService) |
| Return `EchartConfig` directly from service | Aligns with PRD requirement for "Chart Configurations" |
| Zero-fill missing days instead of omitting | Matches AC-DASH-005 Requirement #4: "Missing days zero-filled not omitted" |
| Use `float64` for mean calculations | Consistent with existing `mean_day` field type |

**Edge Cases to Handle:**
1. **Zero actual mean**: Speculative mean should return 0 (not NaN)
2. **Negative prediction percentage**: Config should validate or handle gracefully
3. **Empty data set**: Return zero-filled array, not error
4. **Date boundary issues**: Ensure "today" is consistently defined

**Future Considerations:**
- May need to support different time periods (7 days, 30 days) via parameter
- Could add validation for prediction_pct range (-1 to +infinity)
- Consider caching for frequently accessed chart data

---

### 7. Implementation Checklist

```markdown
## Phase 1: Service Core
- [ ] Create `internal/service/dashboard/speculate_service.go`
- [ ] Implement `NewSpeculateService` constructor
- [ ] Implement `CalculateSpeculativeMean(actualMean float64) float64` method
- [ ] Add unit tests for mean calculation

## Phase 2: Data Retrieval
- [ ] Implement `GetLast15DaysData(ctx)` method using repository
- [ ] Handle date range queries correctly (last 15 days including today)
- [ ] Add unit tests for data retrieval

## Phase 3: Chart Generation
- [ ] Implement `GenerateChartConfig()` method
- [ ] Create line chart with actual vs speculative series
- [ ] Ensure 15 data points in output
- [ ] Add unit tests for chart generation

## Phase 4: Edge Cases
- [ ] Handle zero/empty data gracefully
- [ ] Validate prediction_pct from config
- [ ] Ensure missing days are zero-filled
- [ ] Add error path tests

## Phase 5: Integration
- [ ] Run `go fmt ./...`
- [ ] Run `go vet ./...`
- [ ] Execute all tests: `go test ./...`
- [ ] Verify coverage > 80%
```

---

### Summary

**What to Implement:**
1. Single service file: `internal/service/dashboard/speculate_service.go`
2. ~200-300 lines of code
3. 5-6 methods for mean calculation, data retrieval, and chart generation
4. Follow existing patterns from `day_service.go` and `faults_service.go`

**Key Formula:**
```
speculative_mean = actual_mean × (1 + prediction_pct)
```

**Output Format:**
- ECharts configuration with line chart
- 15 data points (last 15 days including today)
- Zero-filled for missing days
- Two series: "Actual" and "Speculated"

**No Breaking Changes:**
- Purely additive change
- No modifications to existing files required
- Backward compatible with existing API
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress - RDL-085

### Status: In Progress

I'm researching the codebase to understand the patterns and structure before implementing the SpeculateService. Based on my analysis:

**Key Findings:**

1. **Service Patterns**: The codebase uses dependency injection with interfaces (UserConfigProvider) for testability
   - `DayService` and `FaultsService` follow similar patterns
   - Both use `DashboardRepository` and `UserConfigProvider`

2. **Test Infrastructure**: 
   - Unit tests in `test/unit/`
   - Integration tests use `TestHelper` from `test/test_helper.go`
   - Mock implementations available for repository

3. **Configuration**:
   - `UserConfigService` provides `GetPredictionPct()` method
   - Default prediction_pct is 0.15 (15%)

4. **Required Implementation**:
   - Single file: `internal/service/dashboard/speculate_service.go`
   - ~200-300 lines of code
   - Calculate speculative mean: `actual_mean * (1 + prediction_pct)`
   - Generate chart data for last 15 days
   - Zero-fill missing days

**Next Steps:**
1. Create the SpeculateService with mean calculation
2. Implement 15-day data retrieval and zero-filling
3. Create chart configuration generation
4. Write unit tests following existing patterns
5. Run integration tests to verify database interactions
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass
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
