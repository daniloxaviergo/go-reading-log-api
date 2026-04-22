---
id: RDL-087
title: '[doc-008 Phase 3] Implement mean progress calculation with visual map colors'
status: To Do
assignee:
  - workflow
created_date: '2026-04-21 15:51'
updated_date: '2026-04-22 10:41'
labels:
  - phase-3
  - service
  - progress
dependencies: []
references:
  - REQ-DASH-007
  - AC-DASH-007
  - Implementation Checklist Phase 3
documentation:
  - doc-008
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement mean progress calculation logic. Calculate daily progress as (daily_pages / mean_pages) * 100 - 100, apply visual map color ranges (gray 0-10%, cyan 10-20%, blue 20-50%, green >50%, red negative), cover last 30 days.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Daily progress percentage calculated correctly
- [ ] #2 Color ranges applied per specification
- [ ] #3 Last 30 days data coverage
- [ ] #4 Visual map configuration generated
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

**Understanding the Task:**
The task requires implementing a mean progress calculation with visual map colors for the dashboard. Based on the PRD (AC-DASH-007), this involves:

1. **Mean Progress Calculation**: Calculate daily progress as `(daily_pages / mean_pages) * 100 - 100`
2. **Visual Map Colors**: Apply color coding based on progress ranges:
   - Gray: 0% to 10%
   - Cyan: 10% to 20%
   - Blue: 20% to 50%
   - Green: >50%
   - Red: Negative values
3. **Data Coverage**: Last 30 days of data
4. **Chart Type**: Line chart with color-coded data points

**Architecture Decision:**
Create a new `MeanProgressService` following the established pattern in `internal/service/dashboard/`. This service will:
- Use dependency injection for testability
- Implement calculations separate from HTTP handling
- Return ECharts configuration objects

**Why This Approach:**
- Consistent with existing services (DayService, FaultsService, SpeculateService)
- Testable without HTTP stack
- Reusable across different API versions
- Clear separation of concerns

---

### 2. Files to Modify

#### New Files to Create:

| File Path | Purpose |
|-----------|---------|
| `internal/service/dashboard/mean_progress_service.go` | Core calculation logic for mean progress with color mapping |

#### Existing Files to Modify:

| File Path | Modification |
|-----------|--------------|
| `internal/api/v1/handlers/dashboard_handler.go` | Replace placeholder `MeanProgress` handler with full implementation using service |
| `internal/domain/dto/dashboard_response.go` | Add visual map color configuration support to Series/Chart |
| `internal/repository/dashboard_repository.go` | Verify `GetLogsByDateRange` method exists (already exists) |

---

### 3. Dependencies

**Prerequisites:**
- [x] Phase 1 dashboard endpoints implemented (day.json, projects.json)
- [x] UserConfigService available for configuration access
- [x] DashboardRepository interface with all required methods
- [x] ECharts DTOs in place (EchartConfig, Series, etc.)

**External Dependencies:**
- Go 1.25.7
- PostgreSQL with existing schema
- No new external packages required

---

### 4. Code Patterns

**Pattern 1: Service Layer Calculation**
```go
// Calculate mean progress for a single day
func (s *MeanProgressService) CalculateDailyProgress(dailyPages, meanPages float64) float64 {
    if meanPages == 0 {
        return 0.0
    }
    return math.Round(((dailyPages/meanPages)*100-100)*1000)/1000
}
```

**Pattern 2: Color Mapping**
```go
// Determine color based on progress percentage
func (s *MeanProgressService) GetColorForProgress(progress float64) string {
    switch {
    case progress < 0:
        return "#ff4d4f" // Red
    case progress >= 0 && progress < 10:
        return "#959595" // Gray
    case progress >= 10 && progress < 20:
        return "#1890ff" // Cyan
    case progress >= 20 && progress < 50:
        return "#67c23a" // Green
    default: // >= 50
        return "#108ee9" // Blue
    }
}
```

**Pattern 3: ECharts Visual Map**
```go
// Configure visual map for continuous color encoding
visualMap := map[string]interface{}{
    "show": true,
    "min": -100,  // Minimum expected progress
    "max": 100,   // Maximum expected progress
    "inRange": map[string]interface{}{
        "color": []string{
            "#ff4d4f",  // Red (negative)
            "#959595",  // Gray (0-10%)
            "#1890ff",  // Cyan (10-20%)
            "#67c23a",  // Green (20-50%)
            "#108ee9",  // Blue (>50%)
        },
    },
}
```

---

### 5. Testing Strategy

**Unit Tests (`test/mean_progress_service_test.go`):**

| Test Case | Description |
|-----------|-------------|
| `TestMeanProgress_CalculateDailyProgress_ZeroMean` | Handle division by zero |
| `TestMeanProgress_CalculateDailyProgress_Negative` | Negative progress calculation |
| `TestMeanProgress_GetColorForProgress_Gray` | Gray color for 0-10% |
| `TestMeanProgress_GetColorForProgress_Cyan` | Cyan color for 10-20% |
| `TestMeanProgress_GetColorForProgress_Blue` | Blue color for 20-50% |
| `TestMeanProgress_GetColorForProgress_Green` | Green color for >50% |
| `TestMeanProgress_GetColorForProgress_Red` | Red color for negative |
| `TestMeanProgress_GenerateChartConfig_Last30Days` | Verify 30-day data coverage |

**Integration Tests (`test/dashboard_integration_test.go`):**

| Test Case | Description |
|-----------|-------------|
| `TestMeanProgressEndpoint_200OK` | Successful response with chart config |
| `TestMeanProgressEndpoint_DataCoverage` | Verify last 30 days included |
| `TestMeanProgressEndpoint_ColorRanges` | Verify colors match specification |
| `TestMeanProgressEndpoint_EmptyDatabase` | Handle empty logs table |

---

### 6. Risks and Considerations

**Risk 1: Mean Pages Calculation**
- **Issue**: `mean_pages` needs clarification - is it average per day, per project, or overall?
- **Mitigation**: Based on PRD context, use average pages per day across all projects for the date range
- **Decision**: Calculate mean as `total_pages / days_in_range`

**Risk 2: Color Range Boundaries**
- **Issue**: Overlapping boundaries (e.g., does 10% belong to gray or cyan?)
- **Mitigation**: Use half-open intervals `[0, 10)`, `[10, 20)`, etc.
- **Decision**: 10% exactly will be Cyan (first range that includes it)

**Risk 3: Zero Division**
- **Issue**: Mean pages could be zero if no logs exist
- **Mitigation**: Return 0% progress with gray color when mean is zero
- **Decision**: Handle gracefully without returning error

**Risk 4: Date Range Calculation**
- **Issue**: "Last 30 days" could mean different things (calendar vs rolling)
- **Mitigation**: Use rolling 30-day window ending on today
- **Decision**: Today - 29 days to Today (inclusive = 30 days)

**Risk 5: ECharts Configuration**
- **Issue**: Visual map configuration may need tuning for optimal display
- **Mitigation**: Start with simple color ranges, allow customization via config
- **Decision**: Implement basic visual map first, extend if needed

---

### 7. Implementation Steps

**Step 1: Create MeanProgressService**
```bash
touch internal/service/dashboard/mean_progress_service.go
```
- Implement `CalculateDailyProgress` method
- Implement `GetColorForProgress` method  
- Implement `GenerateChartConfig` method
- Add unit tests

**Step 2: Update DashboardHandler**
- Replace placeholder `MeanProgress` handler with full implementation
- Wire up service dependency
- Return proper ECharts configuration

**Step 3: Extend DTOs (if needed)**
- Ensure Series supports visual map through ItemStyle
- Add VisualMap configuration to EchartConfig if not present

**Step 4: Integration Testing**
- Run against test database
- Verify data coverage (exactly 30 days)
- Verify color assignments match specification

**Step 5: Documentation**
- Update QWEN.md with new endpoint details
- Document calculation formulas
- Add example responses
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
