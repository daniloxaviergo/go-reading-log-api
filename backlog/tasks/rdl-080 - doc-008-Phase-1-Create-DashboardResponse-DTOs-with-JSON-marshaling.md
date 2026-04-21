---
id: RDL-080
title: '[doc-008 Phase 1] Create DashboardResponse DTOs with JSON marshaling'
status: To Do
assignee:
  - catarina
created_date: '2026-04-21 15:49'
updated_date: '2026-04-21 17:14'
labels:
  - phase-1
  - dto
  - api
dependencies: []
references:
  - REQ-DASH-003
  - AC-DASH-003
  - 'Decision 5: Response Format - Chart Configurations'
documentation:
  - doc-008
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement internal/domain/dto/dashboard_response.go defining all response structures including DashboardResponse, EchartConfig, StatsData, and LogEntry. Ensure proper JSON field tags and implement validation methods for each DTO.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 All response DTOs defined with correct JSON tags
- [ ] #2 EchartConfig supports ECharts-style configurations
- [ ] #3 StatsData includes all required aggregate fields
- [ ] #4 Validation methods implemented for each DTO
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task involves **completing the DashboardResponse DTO implementation** with proper JSON marshaling, validation methods, and comprehensive test coverage. The approach extends the existing partial implementation in `internal/domain/dto/dashboard_response.go`.

**Architecture Decision**: Extend the existing dashboard_response.go file to include all response structures for Phase 1 dashboard endpoints. This aligns with Decision 5 from doc-008 which specifies ECharts-style JSON configurations as the response format.

**Key Design Considerations**:
- Use `json` tags for proper serialization matching Rails API conventions (snake_case)
- Implement validation methods for each DTO to ensure data integrity
- Include context embedding for request lifecycle tracking (matching existing pattern in DailyStats, ProjectAggregate, etc.)
- Support all 8 dashboard endpoints: day.json, projects.json, last_days.json, faults.json, speculate_actual.json, faults_week_day.json, mean_progress.json, last_year_total.json

**Response Structure Overview**:
| Endpoint | Main DTO | Additional Data |
|----------|----------|-----------------|
| `/v1/dashboard/day.json` | `DashboardResponse` | `StatsData`, `EchartConfig` |
| `/v1/dashboard/projects.json` | `DashboardResponse` | `StatsData`, `EchartConfig` |
| `/v1/dashboard/last_days.json` | `DashboardResponse` | `StatsData`, `LogEntry[]` |
| `/v1/dashboard/echart/faults.json` | `EchartConfig` | Gauge chart configuration |
| `/v1/dashboard/echart/speculate_actual.json` | `EchartConfig` | Line chart with predicted vs actual |
| `/v1/dashboard/echart/faults_week_day.json` | `EchartConfig` | Radar chart for weekday distribution |
| `/v1/dashboard/echart/mean_progress.json` | `EchartConfig` | Line chart with visual map colors |
| `/v1/dashboard/echart/last_year_total.json` | `EchartConfig` | Weekly trend with average line |

---

### 2. Files to Modify

#### New Files to Create:

```
internal/domain/dto/
└── dashboard_response.go          # Complete response DTOs (extension of existing partial file)
```

#### Modified Files:

```
internal/api/v1/routes.go           # Add dashboard route registrations
internal/api/v1/handlers/dashboard_handler.go  # HTTP handlers (RDL-081)
test/unit/dashboard_response_test.go           # Unit tests for DTOs
```

---

### 3. Dependencies

**Prerequisites for Implementation**:
- [x] Go 1.25.7 environment ready
- [x] Existing `internal/domain/dto/` package structure
- [x] `internal/repository/dashboard_repository.go` interface (RDL-079)
- [x] `internal/adapter/postgres/dashboard_repository.go` implementation (RDL-079)
- [x] `internal/service/user_config_service.go` (RDL-078)
- [x] Existing partial `dashboard_response.go` with DailyStats, ProjectAggregate, FaultStats, WeekdayFaults

**External Dependencies** (already in go.mod):
- No new dependencies required

---

### 4. Code Patterns

**Pattern 1: Extended DTO Definition with JSON Tags**

```go
// internal/domain/dto/dashboard_response.go
package dto

import "context"

// DashboardResponse is the main response structure for all dashboard endpoints
type DashboardResponse struct {
    ctx      context.Context
    Echart   *EchartConfig `json:"echart,omitempty"`
    Stats    *StatsData    `json:"stats,omitempty"`
    Logs     []LogEntry    `json:"logs,omitempty"`
}

// NewDashboardResponse creates a new DashboardResponse with optional fields
func NewDashboardResponse() *DashboardResponse {
    return &DashboardResponse{
        Logs: make([]LogEntry, 0),
    }
}

// SetEchart sets the ECharts configuration
func (d *DashboardResponse) SetEchart(config *EchartConfig) *DashboardResponse {
    d.Echart = config
    return d
}

// SetStats sets the statistics data
func (d *DashboardResponse) SetStats(stats *StatsData) *DashboardResponse {
    d.Stats = stats
    return d
}

// AddLog adds a log entry to the response
func (d *DashboardResponse) AddLog(log LogEntry) *DashboardResponse {
    d.Logs = append(d.Logs, log)
    return d
}
```

**Pattern 2: EchartConfig for Chart Specifications**

```go
// EchartConfig holds ECharts-style chart configurations
type EchartConfig struct {
    ctx       context.Context
    Title     string                 `json:"title,omitempty"`
    Tooltip   map[string]interface{} `json:"tooltip,omitempty"`
    Legend    *Legend                `json:"legend,omitempty"`
    Series    []Series               `json:"series,omitempty"`
    XAxis     *Axis                  `json:"xAxis,omitempty"`
    YAxis     *Axis                  `json:"yAxis,omitempty"`
    Grid      *Grid                  `json:"grid,omitempty"`
    Toolbox   map[string]interface{} `json:"toolbox,omitempty"`
}

// Legend configuration
type Legend struct {
    Show   bool     `json:"show"`
    Data   []string `json:"data,omitempty"`
    Top    string   `json:"top,omitempty"`
    Left   string   `json:"left,omitempty"`
}

// Series represents a chart series
type Series struct {
    Name      string                 `json:"name"`
    Type      string                 `json:"type"`  // line, bar, pie, gauge, radar
    Data      []interface{}          `json:"data"`
    ItemStyle map[string]interface{} `json:"itemStyle,omitempty"`
    LineStyle map[string]interface{} `json:"lineStyle,omitempty"`
    AreaStyle map[string]interface{} `json:"areaStyle,omitempty"`
}

// Axis configuration
type Axis struct {
    Type       string   `json:"type"`
    Name       string   `json:"name,omitempty"`
    BoundaryGap []bool  `json:"boundaryGap,omitempty"`
}

// Grid configuration
type Grid struct {
    Left   string `json:"left,omitempty"`
    Right  string `json:"right,omitempty"`
    Top    string `json:"top,omitempty"`
    Bottom string `json:"bottom,omitempty"`
}
```

**Pattern 3: StatsData for Aggregate Statistics**

```go
// StatsData holds all statistical calculations for dashboard views
type StatsData struct {
    ctx              context.Context
    PreviousWeekPages int     `json:"previous_week_pages"`  // Sum of pages from previous 7 days
    LastWeekPages     int     `json:"last_week_pages"`      // Sum of pages from last 7 days
    PerPages          float64 `json:"per_pages"`            // Ratio: last_week / previous_week
    MeanDay           float64 `json:"mean_day"`             // Average pages per day (current weekday)
    SpecMeanDay       float64 `json:"spec_mean_day"`        // Predicted average for current weekday
    ProgressGeral     float64 `json:"progress_geral"`       // Overall completion percentage
    TotalPages        int     `json:"total_pages"`          // Sum of all project total_page values
    Pages             int     `json:"pages"`                // Sum of all project page values
    CountPages        int     `json:"count_pages"`          // Sum of read_pages in period
    SpeculatePages    int     `json:"speculate_pages"`      // Config-based target for period
}

// NewStatsData creates a new StatsData with zero values
func NewStatsData() *StatsData {
    return &StatsData{
        PerPages:      0.0,
        MeanDay:       0.0,
        SpecMeanDay:   0.0,
        ProgressGeral: 0.0,
    }
}

// RoundToThreeDecimals rounds float64 values to 3 decimal places
func (s *StatsData) RoundToThreeDecimals() {
    s.PerPages = math.Round(s.PerPages*1000) / 1000
    s.MeanDay = math.Round(s.MeanDay*1000) / 1000
    s.SpecMeanDay = math.Round(s.SpecMeanDay*1000) / 1000
    s.ProgressGeral = math.Round(s.ProgressGeral*1000) / 1000
}
```

**Pattern 4: LogEntry for Eager-Loaded Log Data**

```go
// LogEntry represents a log entry with eager-loaded project data
type LogEntry struct {
    ctx        context.Context
    ID         int64     `json:"id"`
    Data       string    `json:"data"`          // RFC3339 formatted timestamp
    StartPage  int       `json:"start_page"`
    EndPage    int       `json:"end_page"`
    Note       *string   `json:"note,omitempty"`
    Project    *Project  `json:"project"`       // Eager-loaded project data
    ReadPages  int       `json:"read_pages"`    // Calculated: end_page - start_page
}

// NewLogEntry creates a new LogEntry with calculated read_pages
func NewLogEntry(id int64, data string, startPage, endPage int, note *string, project *Project) *LogEntry {
    return &LogEntry{
        ID:        id,
        Data:      data,
        StartPage: startPage,
        EndPage:   endPage,
        Note:      note,
        Project:   project,
        ReadPages: endPage - startPage,
    }
}
```

**Pattern 5: Validation Methods**

```go
// Validate validates the DashboardResponse structure
func (d *DashboardResponse) Validate() error {
    if d == nil {
        return fmt.Errorf("dashboard response is nil")
    }
    
    // Validate EchartConfig if present
    if d.Echart != nil {
        if err := d.Echart.Validate(); err != nil {
            return fmt.Errorf("echart config validation failed: %w", err)
        }
    }
    
    // Validate StatsData if present
    if d.Stats != nil {
        if err := d.Stats.Validate(); err != nil {
            return fmt.Errorf("stats data validation failed: %w", err)
        }
    }
    
    // Validate each log entry
    for i, log := range d.Logs {
        if err := log.Validate(); err != nil {
            return fmt.Errorf("log entry %d validation failed: %w", i, err)
        }
    }
    
    return nil
}

// Validate validates the EchartConfig structure
func (e *EchartConfig) Validate() error {
    if e == nil {
        return fmt.Errorf("echart config is nil")
    }
    
    // Validate series
    if len(e.Series) == 0 {
        return fmt.Errorf("series cannot be empty")
    }
    
    for i, series := range e.Series {
        if series.Name == "" {
            return fmt.Errorf("series %d: name is required", i)
        }
        if series.Type == "" {
            return fmt.Errorf("series %d (%s): type is required", i, series.Name)
        }
        if len(series.Data) == 0 {
            return fmt.Errorf("series %d (%s): data cannot be empty", i, series.Name)
        }
    }
    
    return nil
}

// Validate validates the StatsData structure
func (s *StatsData) Validate() error {
    if s == nil {
        return fmt.Errorf("stats data is nil")
    }
    
    // Validate non-negative values
    if s.PreviousWeekPages < 0 {
        return fmt.Errorf("previous_week_pages cannot be negative")
    }
    if s.LastWeekPages < 0 {
        return fmt.Errorf("last_week_pages cannot be negative")
    }
    if s.TotalPages < 0 {
        return fmt.Errorf("total_pages cannot be negative")
    }
    if s.Pages < 0 {
        return fmt.Errorf("pages cannot be negative")
    }
    
    // Validate percentage ranges (0-100)
    if s.PerPages < 0 || s.PerPages > 100 {
        return fmt.Errorf("per_pages must be between 0 and 100")
    }
    if s.ProgressGeral < 0 || s.ProgressGeral > 100 {
        return fmt.Errorf("progress_geral must be between 0 and 100")
    }
    
    // Validate non-negative floats
    if s.MeanDay < 0 {
        return fmt.Errorf("mean_day cannot be negative")
    }
    if s.SpecMeanDay < 0 {
        return fmt.Errorf("spec_mean_day cannot be negative")
    }
    
    return nil
}

// Validate validates the LogEntry structure
func (l *LogEntry) Validate() error {
    if l == nil {
        return fmt.Errorf("log entry is nil")
    }
    
    // Validate required fields
    if l.ID <= 0 {
        return fmt.Errorf("id must be positive")
    }
    if l.StartPage < 0 {
        return fmt.Errorf("start_page cannot be negative")
    }
    if l.EndPage < 0 {
        return fmt.Errorf("end_page cannot be negative")
    }
    if l.ReadPages < 0 {
        return fmt.Errorf("read_pages cannot be negative (end_page must >= start_page)")
    }
    
    // Validate project data if present
    if l.Project != nil {
        if err := l.Project.Validate(); err != nil {
            return fmt.Errorf("project validation failed: %w", err)
        }
    }
    
    return nil
}
```

**Pattern 6: Project Reference in LogEntry**

```go
// Project represents minimal project data for eager loading
type Project struct {
    ctx       context.Context
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    TotalPage int    `json:"total_page"`
    Page      int    `json:"page"`
}

// NewProject creates a new Project instance
func NewProject(id int64, name string, totalPage, page int) *Project {
    return &Project{
        ctx:       context.Background(),
        ID:        id,
        Name:      name,
        TotalPage: totalPage,
        Page:      page,
    }
}

// Validate validates the Project structure
func (p *Project) Validate() error {
    if p == nil {
        return fmt.Errorf("project is nil")
    }
    if p.ID <= 0 {
        return fmt.Errorf("project id must be positive")
    }
    if p.TotalPage <= 0 {
        return fmt.Errorf("total_page must be positive")
    }
    if p.Page < 0 {
        return fmt.Errorf("page cannot be negative")
    }
    if p.Page > p.TotalPage {
        return fmt.Errorf("page cannot exceed total_page")
    }
    return nil
}
```

---

### 5. Testing Strategy

#### Unit Tests Structure:

```go
// test/unit/dashboard_response_test.go
package test

import (
    "context"
    "encoding/json"
    "math"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "go-reading-log-api-next/internal/domain/dto"
)

// TestDashboardResponse_BasicCreation tests basic response creation
func TestDashboardResponse_BasicCreation(t *testing.T) {
    // Arrange
    response := dto.NewDashboardResponse()
    
    // Act
    response.SetEchart(dto.NewEchartConfig()).
        SetStats(dto.NewStatsData()).
        AddLog(*dto.NewLogEntry(1, "2024-01-15T10:30:00Z", 0, 25, nil, nil))
    
    // Assert
    assert.NotNil(t, response)
    assert.NotNil(t, response.Echart)
    assert.NotNil(t, response.Stats)
    assert.Len(t, response.Logs, 1)
}

// TestDashboardResponse_JSONSerialization tests JSON marshaling
func TestDashboardResponse_JSONSerialization(t *testing.T) {
    // Arrange
    response := dto.NewDashboardResponse()
    
    project := dto.NewProject(1, "Test Project", 200, 50)
    log := dto.NewLogEntry(1, "2024-01-15T10:30:00Z", 0, 25, nil, project)
    response.AddLog(*log).
        SetStats(dto.NewStatsData().
            SetPreviousWeekPages(100).
            SetLastWeekPages(150).
            SetPerPages(1.5))
    
    // Act
    data, err := json.Marshal(response)
    
    // Assert
    require.NoError(t, err)
    
    // Verify JSON structure
    var result map[string]interface{}
    err = json.Unmarshal(data, &result)
    require.NoError(t, err)
    
    assert.Equal(t, 1.5, result["stats"].(map[string]interface{})["per_pages"])
    assert.Equal(t, float64(100), result["stats"].(map[string]interface{})["previous_week_pages"])
    assert.Equal(t, float64(150), result["stats"].(map[string]interface{})["last_week_pages"])
}

// TestDashboardResponse_Validation tests validation methods
func TestDashboardResponse_Validation(t *testing.T) {
    testCases := []struct {
        name     string
        response *dto.DashboardResponse
        expectError bool
        errorMsg string
    }{
        {
            name:     "valid response",
            response: dto.NewDashboardResponse().SetStats(dto.NewStatsData()),
            expectError: false,
        },
        {
            name: "negative previous_week_pages",
            response: dto.NewDashboardResponse().
                SetStats(dto.NewStatsData().SetPreviousWeekPages(-10)),
            expectError: true,
            errorMsg: "previous_week_pages cannot be negative",
        },
        {
            name: "invalid percentage",
            response: dto.NewDashboardResponse().
                SetStats(dto.NewStatsData().SetPerPages(150)),
            expectError: true,
            errorMsg: "per_pages must be between 0 and 100",
        },
        {
            name: "empty series",
            response: dto.NewDashboardResponse().
                SetEchart(dto.NewEchartConfigWithSeries([]dto.Series{})),
            expectError: true,
            errorMsg: "series cannot be empty",
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            err := tc.response.Validate()
            
            if tc.expectError {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tc.errorMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

// TestEchartConfig_JSONSerialization tests ECharts config serialization
func TestEchartConfig_JSONSerialization(t *testing.T) {
    // Arrange
    config := dto.NewEchartConfig().
        SetTitle("Reading Progress").
        SetTooltip(map[string]interface{}{
            "trigger": "axis",
            "axisPointer": map[string]interface{}{
                "type": "shadow",
            },
        }).
        AddSeries(dto.Series{
            Name: "Pages Read",
            Type: "line",
            Data: []int{10, 20, 30, 40, 50},
            ItemStyle: map[string]interface{}{
                "color": "#5470C6",
            },
        })
    
    // Act
    data, err := json.Marshal(config)
    
    // Assert
    require.NoError(t, err)
    
    var result map[string]interface{}
    err = json.Unmarshal(data, &result)
    require.NoError(t, err)
    
    assert.Equal(t, "Reading Progress", result["title"])
    assert.Equal(t, "axis", result["tooltip"].(map[string]interface{})["trigger"])
    assert.Len(t, result["series"].([]interface{}), 1)
}

// TestStatsData_Rounding tests decimal rounding
func TestStatsData_Rounding(t *testing.T) {
    // Arrange
    stats := dto.NewStatsData().
        SetPreviousWeekPages(100).
        SetLastWeekPages(33).
        SetPerPages(0.333333).
        SetMeanDay(12.345678).
        SetSpecMeanDay(14.987654)
    
    // Act
    stats.RoundToThreeDecimals()
    
    // Assert
    assert.Equal(t, 0.333, stats.PerPages)
    assert.Equal(t, 12.346, stats.MeanDay)
    assert.Equal(t, 14.988, stats.SpecMeanDay)
}

// TestLogEntry_CalculatedFields tests calculated fields
func TestLogEntry_CalculatedFields(t *testing.T) {
    // Arrange
    project := dto.NewProject(1, "Test Project", 200, 50)
    
    // Act
    log := dto.NewLogEntry(1, "2024-01-15T10:30:00Z", 10, 35, nil, project)
    
    // Assert
    assert.Equal(t, 25, log.ReadPages)  // 35 - 10 = 25
    
    // Test validation with invalid page range
    invalidLog := dto.NewLogEntry(2, "2024-01-15T10:30:00Z", 35, 10, nil, project)
    err := invalidLog.Validate()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "read_pages cannot be negative")
}

// TestDashboardResponse_EmptyValues tests handling of empty/zero values
func TestDashboardResponse_EmptyValues(t *testing.T) {
    // Arrange
    response := dto.NewDashboardResponse()
    
    // Act - Marshal with zero values
    data, err := json.Marshal(response)
    
    // Assert
    require.NoError(t, err)
    
    var result map[string]interface{}
    err = json.Unmarshal(data, &result)
    require.NoError(t, err)
    
    // Empty slices should be marshaled as empty arrays
    assert.Equal(t, []interface{}{}, result["logs"])
    
    // Nil optional fields should be omitted (omitempty)
    assert.NotContains(t, result, "echart")
    assert.NotContains(t, result, "stats")
}

// TestDashboardResponse_ConcurrentAccess tests thread safety
func TestDashboardResponse_ConcurrentAccess(t *testing.T) {
    // Arrange
    response := dto.NewDashboardResponse()
    
    // Act - Concurrent modifications
    done := make(chan bool)
    for i := 0; i < 10; i++ {
        go func(idx int) {
            log := dto.NewLogEntry(int64(idx), "2024-01-15T10:30:00Z", 0, 10, nil, nil)
            response.AddLog(*log)
            done <- true
        }(i)
    }
    
    // Wait for all goroutines
    for i := 0; i < 10; i++ {
        <-done
    }
    
    // Assert
    assert.Len(t, response.Logs, 10)
}

// TestIntegration_DashboardResponseWithRealData tests with simulated real data
func TestIntegration_DashboardResponseWithRealData(t *testing.T) {
    // Arrange - Simulate data from database
    testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
    
    // Create sample logs
    logs := []dto.LogEntry{
        *dto.NewLogEntry(1, testDate, 0, 25, nil, dto.NewProject(1, "Project A", 200, 50)),
        *dto.NewLogEntry(2, testDate, 25, 50, nil, dto.NewProject(1, "Project A", 200, 50)),
        *dto.NewLogEntry(3, testDate, 0, 30, nil, dto.NewProject(2, "Project B", 300, 80)),
    }
    
    // Create stats
    stats := dto.NewStatsData().
        SetPreviousWeekPages(150).
        SetLastWeekPages(200).
        SetPerPages(1.333).
        SetMeanDay(25.5).
        SetSpecMeanDay(29.325)
    
    // Create ECharts config
    echart := dto.NewEchartConfig().
        SetTitle("Daily Progress").
        AddSeries(dto.Series{
            Name: "Pages",
            Type: "line",
            Data: []int{10, 20, 30, 40, 50},
        })
    
    // Build response
    response := dto.NewDashboardResponse().
        SetStats(stats).
        SetEchart(echart)
    
    for _, log := range logs {
        response.AddLog(log)
    }
    
    // Act - Marshal and unmarshal
    data, err := json.Marshal(response)
    require.NoError(t, err)
    
    var result map[string]interface{}
    err = json.Unmarshal(data, &result)
    require.NoError(t, err)
    
    // Assert
    assert.Equal(t, 3, len(result["logs"].([]interface{})))
    assert.Equal(t, 1.333, result["stats"].(map[string]interface{})["per_pages"])
    assert.Equal(t, "Daily Progress", result["echart"].(map[string]interface{})["title"])
}
```

#### Test Coverage Targets:

| Test File | Coverage Areas | Target Lines |
|-----------|---------------|--------------|
| `test/unit/dashboard_response_test.go` | DTO creation, validation, JSON serialization | ~400-500 |
| **Total** | All response structures and methods | ~400-500 |

---

### 6. Risks and Considerations

#### Known Risks:

| Risk | Impact | Mitigation |
|------|--------|------------|
| **JSON tag mismatch with Rails** | Medium - API incompatibility | Verify all JSON tags match Rails conventions (snake_case) |
| **Context propagation** | Low - Memory leaks if not handled | Follow existing pattern: embed context, provide Get/Set methods |
| **Circular references in JSON** | High - Marshal will panic | Ensure Project reference in LogEntry doesn't create cycles |
| **Float precision loss** | Medium - Calculation accuracy | Use `math.Round` for 3-decimal precision as specified |
| **Missing validation on nil pointers** | High - Runtime panics | Add nil checks in all Validate() methods |
| **Time format inconsistency** | Medium - Parsing errors | Use RFC3339 consistently, validate date strings |

#### Design Trade-offs:

1. **Pointer vs Value Types for Optional Fields**: 
   - Chose pointers for optional fields (Echart, Stats, Note) to enable `omitempty` JSON behavior
   - Enables distinguishing between "not set" (nil) and "explicitly zero"

2. **Method Chaining for Builder Pattern**:
   - Returns `*DashboardResponse` from setter methods for fluent API
   - Improves readability: `response.SetStats(...).SetEchart(...)`
   
3. **Embedded Context**:
   - Follows existing pattern in codebase (see `DailyStats`, `ProjectAggregate`)
   - Allows request-scoped values without passing context through every method

4. **Validation at DTO Level**:
   - Validates on construction and before JSON marshaling
   - Fails fast with descriptive errors

5. **Separate Series Type vs Inline**:
   - Created explicit `Series`, `Legend`, `Axis` types for reusability
   - Enables consistent chart configuration across all 8 endpoints

#### Blocking Issues:

1. **Circular Reference Prevention**: Ensure `LogEntry.Project` doesn't create circular references when marshaling
   - Solution: Project should not contain Logs, only minimal identifying data
   
2. **ECharts Configuration Complexity**: Some endpoints require complex chart configurations
   - Solution: Build config incrementally, test each endpoint's requirements separately

#### Deployment Considerations:

- No database migrations required (DTOs are in-memory structures)
- No configuration changes needed
- Backward compatible (new response format extends existing structure)
- Canary deployment safe (old endpoints return same format)

---

### 7. Implementation Checklist

**Phase 1: Core DTO Definitions**
- [ ] Create `internal/domain/dto/dashboard_response.go`
- [ ] Define `DashboardResponse` struct with JSON tags
- [ ] Define `EchartConfig` struct with all ECharts options
- [ ] Define `StatsData` struct with all statistical fields
- [ ] Define `LogEntry` struct with eager-loaded Project
- [ ] Implement constructor functions (`NewDashboardResponse`, etc.)
- [ ] Run `go fmt` and `go vet`

**Phase 2: Validation Methods**
- [ ] Implement `Validate()` for `DashboardResponse`
- [ ] Implement `Validate()` for `EchartConfig`
- [ ] Implement `Validate()` for `StatsData`
- [ ] Implement `Validate()` for `LogEntry`
- [ ] Implement `Validate()` for `Project`
- [ ] Add helper methods: `RoundToThreeDecimals()`, `SetContext()`, etc.

**Phase 3: Unit Tests**
- [ ] Create `test/unit/dashboard_response_test.go`
- [ ] Test basic DTO creation and construction
- [ ] Test JSON marshaling/unmarshaling
- [ ] Test validation with valid data
- [ ] Test validation with invalid data (negative values, empty strings)
- [ ] Test float rounding to 3 decimals
- [ ] Test calculated fields (ReadPages in LogEntry)
- [ ] Test edge cases (nil pointers, zero values)

**Phase 4: Integration Verification**
- [ ] Verify no circular references in JSON output
- [ ] Verify field names match Rails API conventions
- [ ] Verify nil optional fields are omitted from JSON
- [ ] Run full test suite: `go test ./...`
- [ ] Check coverage: `go test -cover ./internal/domain/dto/...`

**Phase 5: Documentation**
- [ ] Add godoc comments to all public types and methods
- [ ] Update QWEN.md with new DTO structure
- [ ] Document JSON response format for each endpoint type

---

### 8. Acceptance Criteria Mapping

| AC Requirement | Implementation Approach |
|----------------|------------------------|
| All response DTOs defined with correct JSON tags | Create `dashboard_response.go` with all structs and `json` tags |
| EchartConfig supports ECharts-style configurations | Implement `EchartConfig`, `Series`, `Legend`, `Axis`, `Grid` types |
| StatsData includes all required aggregate fields | Implement `StatsData` with all 12 statistical fields from PRD |
| Validation methods implemented for each DTO | Add `Validate()` method to all public DTOs with comprehensive checks |

---

### 9. Quality Gates

Before marking task complete:
- [ ] All unit tests pass (target: >95% coverage)
- [ ] `go vet` reports no issues
- [ ] `go fmt` applied consistently
- [ ] No circular import dependencies
- [ ] JSON output matches Rails API format exactly
- [ ] Validation errors are descriptive and actionable
- [ ] Documentation updated in QWEN.md
- [ ] Code follows existing patterns in `internal/domain/dto/`
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
