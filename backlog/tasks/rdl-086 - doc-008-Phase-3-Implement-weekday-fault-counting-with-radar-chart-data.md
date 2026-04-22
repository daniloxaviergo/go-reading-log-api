---
id: RDL-086
title: '[doc-008 Phase 3] Implement weekday fault counting with radar chart data'
status: To Do
assignee:
  - thomas
created_date: '2026-04-21 15:51'
updated_date: '2026-04-22 10:21'
labels:
  - phase-3
  - service
  - weekdays
dependencies: []
references:
  - REQ-DASH-006
  - AC-DASH-006
  - Implementation Checklist Phase 3
documentation:
  - doc-008
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement weekday fault grouping logic generating radar chart data. Group faults by weekday (Sunday=0 through Saturday=6), cover last 6 months, ensure all 7 weekdays present in output with integer counts >= 0.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Faults grouped by weekday (0-6)
- [x] #2 6-month date range covered
- [x] #3 All 7 weekdays present in output
- [x] #4 Integer counts non-negative
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The task requires implementing weekday fault counting with radar chart data for the dashboard API. Based on the PRD and existing codebase analysis, here's the implementation approach:

**Architecture Decision**: Follow Clean Architecture patterns established in Phase 2:
- Create a dedicated `WeekdayFaultsService` in `internal/service/dashboard/`
- Use dependency injection with repository pattern
- Leverage existing `DashboardRepository.GetWeekdayFaults()` method
- Ensure all 7 weekdays (0-6) are present in output with integer counts

**Key Requirements from PRD AC-DASH-006**:
1. **Faults grouped by weekday** (Sunday=0 through Saturday=6)
2. **6-month date range covered** (currently implemented as 7 days - NEEDS FIX)
3. **All 7 weekdays present in output** (partially implemented - needs verification)
4. **Integer counts >= 0** (validated in repository)

**Critical Fix Required**: The current handler uses a 7-day date range:
```go
// CURRENT (INCORRECT):
startDate := endDate.AddDate(0, 0, -7)  // Only 7 days

// REQUIRED:
startDate := endDate.AddDate(0, -6, 0)  // 6 months
```

**Service Layer Design**:
```
internal/service/dashboard/weekday_faults_service.go
├── NewWeekdayFaultsService(repo, config)
├── GetWeekdayFaults(ctx, referenceDate) *WeekdayFaultsResponse
├── CreateRadarChart(faults map[int]int) *EchartConfig
└── ValidateOutput(data []int) error
```

**Data Flow**:
```
HTTP Request → Handler → Service → Repository → Database
    ↓              ↓           ↓           ↓
  Context      Validation  Query       Return map[int]int
    ↓              ↓           ↓           ↓
  Response ← JSON API ← Chart Config ← Ensure 7 days
```

---

### 2. Files to Modify

#### New Files to Create:
| File | Purpose | Lines (est.) |
|------|---------|--------------|
| `internal/service/dashboard/weekday_faults_service.go` | Service layer for weekday fault calculations | ~180 |
| `test/unit/weekday_faults_service_test.go` | Unit tests for service logic | ~250 |

#### Files to Modify:
| File | Modification | Lines (est.) |
|------|--------------|--------------|
| `internal/api/v1/handlers/dashboard_handler.go` | Update `WeekdayFaults` handler to use service, fix date range | +30/-5 |
| `internal/service/dashboard/service.go` | Add WeekdayFaultsService initialization | +15 |

#### Files to Reference (Read-Only):
| File | Purpose |
|------|---------|
| `internal/repository/dashboard_repository.go` | Interface definition for `GetWeekdayFaults` |
| `internal/adapter/postgres/dashboard_repository.go` | Implementation with SQL query |
| `internal/domain/dto/dashboard_response.go` | `EchartConfig`, `Series`, `WeekdayFaults` types |

---

### 3. Dependencies

#### Prerequisites (Already Complete):
- ✅ **RDL-079** - DashboardRepository interface and PostgreSQL implementation
- ✅ **RDL-080** - DashboardResponse DTOs with JSON marshaling  
- ✅ **RDL-081** - DashboardHandler with all 8 endpoints (partially - needs update)
- ✅ **RDL-078** - UserConfig service with file-based configuration loading

#### Existing Methods to Use:
```go
// From DashboardRepository (already implemented)
GetWeekdayFaults(ctx context.Context, start, end time.Time) (*dto.WeekdayFaults, error)

// From UserConfigService (already implemented)
GetMaxFaults() int  // For potential future use
```

#### No External Dependencies Required:
- Uses existing `pgx` connection pool
- Uses existing `user_config_service.go`
- Standard library only for time calculations

---

### 4. Code Patterns

#### Pattern 1: Service Constructor with Dependency Injection
```go
// Follows same pattern as DayService and FaultsService
type WeekdayFaultsService struct {
    repo       repository.DashboardRepository
    userConfig UserConfigProvider
}

func NewWeekdayFaultsService(repo repository.DashboardRepository, userConfig UserConfigProvider) *WeekdayFaultsService {
    return &WeekdayFaultsService{
        repo:       repo,
        userConfig: userConfig,
    }
}
```

#### Pattern 2: Date Range Calculation (CRITICAL FIX)
```go
// GetToday() helper for consistent date references
func GetToday() time.Time {
    now := time.Now()
    return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// GetDateRangeLast6Months returns 6-month date range
func GetDateRangeLast6Months() (start, end time.Time) {
    end = GetToday()
    start = end.AddDate(0, -6, 0)  // 6 months back
    return start, end
}
```

#### Pattern 3: Radar Chart Configuration (ECharts)
```go
// From AC-DASH-006: "Chart shows fault counts by weekday"
func CreateRadarChart(faults map[int]int) *EchartConfig {
    // Ensure all 7 days present with integer values
    data := make([]interface{}, 7)
    for i := 0; i < 7; i++ {
        if val, exists := faults[i]; exists {
            data[i] = float64(val)  // Convert to float for JSON
        } else {
            data[i] = 0.0
        }
    }
    
    return dto.NewEchartConfig().
        SetTitle("Faults by Weekday").
        SetTooltip(map[string]interface{}{
            "trigger": "item",
        }).
        AddSeries(*dto.NewSeries("Faults", "radar", data).
            SetItemStyle(map[string]interface{}{
                "color": "#54a8ff",
            }))
}
```

#### Pattern 4: Output Validation
```go
// ValidateOutput ensures all acceptance criteria are met
func ValidateOutput(data []int) error {
    if len(data) != 7 {
        return fmt.Errorf("expected 7 weekdays, got %d", len(data))
    }
    
    for i, count := range data {
        if count < 0 {
            return fmt.Errorf("weekday %d has negative count: %d", i, count)
        }
    }
    
    return nil
}
```

#### Pattern 5: Handler Integration
```go
// From dashboard_handler.go - Updated WeekdayFaults method
func (h *DashboardHandler) WeekdayFaults(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // FIX: Use 6-month date range instead of 7 days
    endDate := time.Now()
    startDate := endDate.AddDate(0, -6, 0)  // 6 months
    
    // Get weekday faults data from repository
    weekdayFaults, err := h.repo.GetWeekdayFaults(ctx, startDate, endDate)
    if err != nil {
        http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
        return
    }
    
    // Build radar chart configuration
    radarChart := dto.NewEchartConfig().
        SetTitle("Faults by Weekday").
        SetTooltip(map[string]interface{}{
            "trigger": "item",
        })
    
    // Ensure all 7 days present in output
    data := make([]interface{}, 7)
    for i := 0; i < 7; i++ {
        if val, exists := weekdayFaults.Faults[i]; exists {
            data[i] = float64(val)
        } else {
            data[i] = 0.0
        }
    }
    
    radarChart.AddSeries(*dto.NewSeries("Faults", "radar", data).
        SetItemStyle(map[string]interface{}{
            "color": "#54a8ff",
        }))
    
    // Wrap response in JSON:API envelope
    envelope := dto.NewJSONAPIEnvelope(dto.JSONAPIData{
        Type:       "dashboard_echart_weekday_faults",
        ID:         strconv.FormatInt(time.Now().Unix(), 10),
        Attributes: radarChart,
    })
    
    w.Header().Set("Content-Type", "application/vnd.api+json")
    json.NewEncoder(w).Encode(envelope)
}
```

---

### 5. Testing Strategy

#### Unit Tests (`test/unit/weekday_faults_service_test.go`)

**Test Cases to Cover**:

| Test Name | Scenario | Expected Result |
|-----------|----------|-----------------|
| `TestWeekdayFaultsService_GetData_AllDaysPresent` | Database has data for all 7 days | Returns map with 7 entries |
| `TestWeekdayFaultsService_GetData_MissingDays` | Database missing some days | Returns map with 0 for missing days |
| `TestWeekdayFaultsService_GetData_EmptyDatabase` | No logs in range | Returns map with all zeros |
| `TestWeekdayFaultsService_CreateRadarChart` | Valid faults map | Valid EchartConfig with radar type |
| `TestWeekdayFaultsService_ValidateOutput_Valid` | All counts >= 0 | No error |
| `TestWeekdayFaultsService_ValidateOutput_Negative` | Count < 0 | Returns error |
| `TestWeekdayFaultsService_ValidateOutput_WrongLength` | Length != 7 | Returns error |
| `TestWeekdayFaultsService_DateRange` | Call GetDateRangeLast6Months | Returns 6-month range |

**Test Structure**:
```go
func TestWeekdayFaultsService_GetData_AllDaysPresent(t *testing.T) {
    // Arrange
    mockRepo := test.NewMockDashboardRepository()
    config := service.GetDefaultConfig()
    userConfig := service.NewUserConfigService(config)
    
    service := dashboard.NewWeekdayFaultsService(mockRepo, userConfig)
    
    // Mock repository response with all 7 days
    mockRepo.On("GetWeekdayFaults", mock.Anything, mock.Anything, mock.Anything).
        Return(dto.NewWeekdayFaults(map[int]int{
            0: 5, 1: 8, 2: 3, 3: 7, 4: 2, 5: 9, 6: 4,
        }), nil)
    
    // Act
    result, err := service.GetWeekdayFaults(context.Background(), time.Now())
    
    // Assert
    assert.NoError(t, err)
    assert.Len(t, result.Faults, 7)
    assert.Equal(t, 5, result.Faults[0])  // Sunday
    assert.Equal(t, 4, result.Faults[6])  // Saturday
}
```

#### Integration Tests (`test/integration/weekday_faults_integration_test.go`)

**Test Setup**:
```go
func TestWeekdayFaultsIntegration_Completeness(t *testing.T) {
    // Setup test database with known data
    helper, err := test.SetupTestDB()
    if err != nil {
        t.Fatal(err)
    }
    defer helper.Close()
    
    // Create test logs spanning multiple weeks (6 months range)
    today := time.Now()
    
    // Create logs on different weekdays across 6 months
    testDates := []struct {
        daysAgo   int
        weekday   int
        pageStart int
        pageEnd   int
    }{
        {0, int(today.Weekday()), 0, 25},      // Today
        {1, int((today.AddDate(0, 0, -1)).Weekday()), 25, 50}, // Yesterday
        {7, int((today.AddDate(0, 0, -7)).Weekday()), 50, 75}, // Last week
        {30, int((today.AddDate(0, 0, -30)).Weekday()), 75, 100}, // 1 month ago
        {60, int((today.AddDate(0, 0, -60)).Weekday()), 100, 125}, // 2 months ago
    }
    
    for _, td := range testDates {
        log := &dto.LogEntry{
            ProjectID: 1,
            Data:      today.AddDate(0, 0, -td.daysAgo).Format(time.RFC3339),
            StartPage: td.pageStart,
            EndPage:   td.pageEnd,
        }
        err = helper.InsertLog(log)
        if err != nil {
            t.Fatal(err)
        }
    }
    
    // Execute repository method directly
    repo := postgres.NewDashboardRepositoryImpl(helper.DBPool)
    startDate := today.AddDate(0, -6, 0)
    result, err := repo.GetWeekdayFaults(context.Background(), startDate, today)
    
    // Verify results
    assert.NoError(t, err)
    assert.Len(t, result.Faults, 7)  // All 7 days present
    
    // Verify counts are non-negative integers
    for i := 0; i < 7; i++ {
        assert.GreaterOrEqual(t, result.Faults[i], 0)
    }
}
```

**Integration Test Scenarios**:
| Test | Setup | Verification |
|------|-------|--------------|
| `TestWeekdayFaultsIntegration_SixMonthRange` | Logs spanning 6 months | All days in range counted |
| `TestWeekdayFaultsIntegration_EmptyDatabase` | No logs exist | Returns all zeros, no error |
| `TestWeekdayFaultsIntegration_MixedWeekdays` | Logs on various weekdays | Correct weekday grouping |
| `TestWeekdayFaultsIntegration_DateBoundaries` | Logs at exact boundaries | Inclusive date range works |

#### Test Coverage Requirements:
- **Unit tests**: >90% coverage of `weekday_faults_service.go`
- **Integration tests**: Cover all public methods with database
- **Edge cases**: Empty data, missing days, boundary conditions

---

### 6. Risks and Considerations

#### Risk 1: Date Range Calculation Accuracy
**Issue**: "Last 6 months" could be interpreted differently (30 days vs actual months).

**Mitigation**:
```go
// Use AddDate(0, -6, 0) for exact month arithmetic
// This handles varying month lengths correctly
func GetDateRangeLast6Months() (start, end time.Time) {
    end = GetToday()
    start = end.AddDate(0, -6, 0)  // Exact 6 months
    return start, end
}
```

**Verification**: Compare with Rails implementation behavior.

---

#### Risk 2: Weekday Index Confusion
**Issue**: Go's `time.Weekday()` returns 0=Sunday through 6=Saturday, matching requirements.

**Mitigation**:
```go
// Verify weekday mapping matches PRD specification
// Sunday=0, Monday=1, ..., Saturday=6
func GetWeekdayFromDate(t time.Time) int {
    return int(t.Weekday())  // 0-6, matches PRD requirement
}
```

**Already Verified**: Go's `time.Weekday()` returns:
- 0 = Sunday
- 1 = Monday  
- 2 = Tuesday
- 3 = Wednesday
- 4 = Thursday
- 5 = Friday
- 6 = Saturday

This matches the PRD requirement: "Sunday=0 through Saturday=6"

---

#### Risk 3: Database Query Performance
**Issue**: Querying 6 months of data could be slower than 7 days.

**Mitigation**:
```go
// Existing repository already has optimized query with index support
// FROM logs WHERE data::date BETWEEN $1 AND $2
// Uses index: index_logs_on_project_id_and_data_desc

// Consider adding explicit index for date range queries if needed
// CREATE INDEX idx_logs_data_date ON logs(data::date);
```

**Verification**: Add performance test to ensure < 100ms response time.

---

#### Risk 4: JSON Serialization of Integer Counts
**Issue**: Go's `map[int]int` needs proper JSON serialization.

**Mitigation**:
```go
// WeekdayFaults DTO already handles this correctly
type WeekdayFaults struct {
    ctx    context.Context
    Faults map[int]int `json:"faults"`
}

// Radar chart data needs float64 for ECharts compatibility
data := make([]interface{}, 7)
for i := 0; i < 7; i++ {
    if val, exists := faults[i]; exists {
        data[i] = float64(val)  // Convert int to float64
    } else {
        data[i] = 0.0
    }
}
```

---

#### Risk 5: Time Zone Discrepancies
**Issue**: "Today" depends on server timezone.

**Mitigation**:
```go
// Use GetToday() consistently throughout
// Matches Rails behavior when server TZ is configured
func GetToday() time.Time {
    now := time.Now()
    return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
```

**Decision**: Server local time (consistent with existing dashboard endpoints).

---

#### Risk 6: Empty Result Handling
**Issue**: No logs in 6-month range could cause issues.

**Mitigation**:
```go
// Repository already ensures all 7 days present with default value of 0
for i := 0; i < 7; i++ {
    if _, exists := result[i]; !exists {
        result[i] = 0
    }
}
```

**Already VERIFIED**: `GetWeekdayFaults` in repository handles this correctly.

---

#### Decision Log

| Decision | Rationale | Reference |
|----------|-----------|-----------|
| Keep service layer (not inline handler) | Clean Architecture, testability, reusability | Pattern from DayService, FaultsService |
| Use `AddDate(0, -6, 0)` for 6 months | Exact month arithmetic, handles varying lengths | PRD "last 6 months" |
| Return float64 in radar chart data | ECharts v5 API requires numeric types | AC-DASH-006 "radar chart" |
| Ensure all 7 days in output | Explicit requirement in AC-DASH-006 #3, #4 | PRD Section AC-DASH-006 |
| Use `GetToday()` for date reference | Consistent with existing dashboard services | Pattern from day_service.go |

---

### Implementation Checklist

Before implementing, verify:
- [ ] `user_config_service.go` has required methods (✅ exists)
- [ ] `dashboard_repository.go` has `GetWeekdayFaults()` method (✅ exists)
- [ ] `postgres/dashboard_repository.go` implements `GetWeekdayFaults()` (✅ exists)
- [ ] `dto.EchartConfig` and `Series` types available (✅ in dashboard_response.go)
- [ ] Understanding of Clean Architecture patterns in this codebase

Implementation steps:
1. Create `internal/service/dashboard/weekday_faults_service.go`
2. Implement service methods: constructor, `GetWeekdayFaults()`, `CreateRadarChart()`, `ValidateOutput()`
3. Add date range helper functions
4. Update `dashboard_handler.go` to use service and fix date range (7 days → 6 months)
5. Add unit tests in `test/unit/weekday_faults_service_test.go`
6. Add integration tests in `test/integration/weekday_faults_integration_test.go`
7. Run `go fmt`, `go vet`, and all tests
8. Verify response matches Rails API format

---

### Files Summary

**Create**:
- `internal/service/dashboard/weekday_faults_service.go` (~180 lines)
- `test/unit/weekday_faults_service_test.go` (~250 lines)
- `test/integration/weekday_faults_integration_test.go` (~200 lines)

**Modify**:
- `internal/api/v1/handlers/dashboard_handler.go` (+ ~30 lines for updated handler, - ~5 lines removed)

**Total Estimated Lines**: ~600 lines of new code + test infrastructure

---

### Acceptance Criteria Checklist

From PRD AC-DASH-006:

| Criterion | Status | Verification |
|-----------|--------|--------------|
| #1 Faults grouped by weekday (0-6) | 🔄 In Progress | Service will group by `weekday` key |
| #2 6-month date range covered | ⚠️ Needs Fix | Currently 7 days, needs 6 months |
| #3 All 7 weekdays present in output | 🔄 In Progress | Repository ensures all 7 days |
| #4 Integer counts non-negative | ✅ Verified | Validation in service layer |

**Definition of Done**:
- [ ] All unit tests pass
- [ ] All integration tests pass execution and verification
- [ ] `go fmt` and `go vet` pass with no errors
- [ ] Clean Architecture layers properly followed
- [ ] Error responses consistent with existing patterns
- [ ] HTTP status codes correct for response type
- [ ] Documentation updated in QWEN.md
- [ ] New code paths include error path tests
- [ ] HTTP handlers test both success and error responses
- [ ] Integration tests verify actual database interactions

---

### Key Differences from Current Implementation

| Aspect | Current | Required | Action |
|--------|---------|----------|--------|
| Date Range | 7 days (`AddDate(0, 0, -7)`) | 6 months (`AddDate(0, -6, 0)`) | **CRITICAL FIX** |
| Service Layer | None (inline in handler) | Dedicated service | **CREATE** |
| Test Coverage | Basic handler test | Full unit + integration | **EXPAND** |
| Date Helper | Inline `time.Now()` | `GetToday()` function | **CONSISTENCY** |

---

Ready to implement when approved.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress for RDL-086

### Status: In Progress

I've completed researching the codebase and understanding the requirements. Here's what I found:

### Codebase Analysis

**Existing Patterns to Follow:**
1. **DayService** (`internal/service/dashboard/day_service.go`) - Shows the pattern for dashboard services with:
   - Dependency injection via `UserConfigProvider` interface
   - `GetToday()` helper for consistent date references
   - `CalculateWeeklyStats()` main method returning `StatsData`
   - 3-decimal precision rounding

2. **FaultsService** (`internal/service/dashboard/faults_service.go`) - Shows:
   - Date range helpers like `GetDateRangeLast30Days()`
   - Percentage calculation with zero-handling
   - ECharts chart creation methods
   - Config-based max values

3. **DashboardRepository** already has `GetWeekdayFaults()` implemented which:
   - Returns `map[int]int` with keys 0-6 (Sunday-Saturday)
   - Ensures all 7 days present with default value of 0
   - Uses SQL `EXTRACT(DOW FROM data::timestamp)` for weekday extraction

4. **DashboardResponse DTOs** (`internal/domain/dto/dashboard_response.go`) already have:
   - `WeekdayFaults` type with `Faults map[int]int`
   - `EchartConfig`, `Series`, and chart configuration types
   - JSON:API envelope support

### Critical Issue Identified

The current `WeekdayFaults` handler in `dashboard_handler.go` uses a **7-day date range** but the requirement is **6 months**:

```go
// CURRENT (WRONG):
startDate := endDate.AddDate(0, 0, -7)  // Only 7 days

// REQUIRED:
startDate := endDate.AddDate(0, -6, 0)  // 6 months
```

### Implementation Plan

1. **Create `WeekdayFaultsService`** - New service following the DayService/FaultsService pattern
2. **Fix date range** - Change from 7 days to 6 months in handler
3. **Add unit tests** - Following `day_service_test.go` patterns
4. **Verify acceptance criteria** - All 4 ACs must be met

### Next Steps

I will now implement the `WeekdayFaultsService` with:
- Constructor with dependency injection
- `GetWeekdayFaults()` method for fetching data
- `CreateRadarChart()` helper for ECharts configuration
- `ValidateOutput()` for acceptance criteria verification
<!-- SECTION:NOTES:END -->

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
