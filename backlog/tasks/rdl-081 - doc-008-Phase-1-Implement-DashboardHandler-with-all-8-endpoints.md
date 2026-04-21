---
id: RDL-081
title: '[doc-008 Phase 1] Implement DashboardHandler with all 8 endpoints'
status: To Do
assignee:
  - thomas
created_date: '2026-04-21 15:49'
updated_date: '2026-04-21 18:09'
labels:
  - phase-1
  - handler
  - api
dependencies: []
references:
  - REQ-DASH-004
  - AC-DASH-004
  - Implementation Checklist Phase 1
documentation:
  - doc-008
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create internal/api/v1/handlers/dashboard_handler.go implementing HTTP handlers for /v1/dashboard/day.json, /v1/dashboard/projects.json, /v1/dashboard/last_days.json, and ECharts endpoints. Include proper error handling, response formatting, and integration with service layer.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 All 8 dashboard endpoints implemented as HTTP handlers
- [ ] #2 Error handling consistent with existing patterns
- [ ] #3 Response formatting matches API conventions
- [ ] #4 Unit tests cover both success and error scenarios
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The DashboardHandler implementation requires creating HTTP handlers for 8 dashboard endpoints that aggregate data from the existing repository layer. The approach follows Clean Architecture principles:

**Architecture Strategy:**
- Create `dashboard_handler.go` in `internal/api/v1/handlers/` following the established pattern from `projects_handler.go`
- Leverage existing `DashboardRepository` interface and `UserConfigService` already created
- Use `DashboardResponse` DTOs from `internal/domain/dto/` for consistent response formatting
- Implement proper error handling matching existing patterns (validation errors return 400, not found returns 404, server errors return 500)
- Register routes in `internal/api/v1/routes.go` following the established pattern

**Endpoint Implementation:**
| Endpoint | Method | Service Call | Response |
|----------|--------|--------------|----------|
| `/v1/dashboard/day.json` | GET | DayService | DailyStats with weekday breakdown |
| `/v1/dashboard/projects.json` | GET | ProjectsService | Project aggregates with progress |
| `/v1/dashboard/last_days.json` | GET | LastDaysService | Trend data with type parameter |
| `/v1/dashboard/echart/faults.json` | GET | FaultsService | Gauge chart config |
| `/v1/dashboard/echart/speculate_actual.json` | GET | SpeculateService | Line chart config |
| `/v1/dashboard/echart/faults_week_day.json` | GET | WeekdayFaultsService | Radar chart config |
| `/v1/dashboard/echart/mean_progress.json` | GET | MeanProgressService | Line chart config |
| `/v1/dashboard/echart/last_year_total.json` | GET | YearlyTotalService | Yearly trend chart |

**Key Design Decisions:**
- Use context with timeout (15 seconds) for all database operations
- Return zero values instead of errors for empty data (graceful degradation)
- Round float values to 3 decimal places for consistency
- Handle nil config gracefully with defaults from `UserConfigService`
- Follow JSON:API envelope pattern for response structure

---

### 2. Files to Modify

#### New Files to Create:

| File Path | Purpose |
|-----------|---------|
| `internal/api/v1/handlers/dashboard_handler.go` | HTTP handlers for all 8 dashboard endpoints |
| `internal/service/dashboard/day_service.go` | Daily statistics calculation logic |
| `internal/service/dashboard/projects_service.go` | Project aggregate calculation logic |
| `internal/service/dashboard/last_days_service.go` | Last days trend data logic |
| `internal/service/dashboard/faults_service.go` | Fault counting and percentage logic |
| `internal/service/dashboard/speculate_service.go` | Speculated vs actual reading comparison |
| `internal/service/dashboard/weekday_faults_service.go` | Weekday fault distribution |
| `internal/service/dashboard/mean_progress_service.go` | Mean progress calculation |
| `internal/service/dashboard/yearly_total_service.go` | Yearly trend data |

#### Files to Modify:

| File Path | Modification |
|-----------|--------------|
| `internal/api/v1/routes.go` | Add dashboard route registrations in `SetupRoutes()` |
| `go.mod` | May need new dependencies (check imports) |

#### Test Files to Create:

| File Path | Purpose |
|-----------|---------|
| `test/unit/dashboard_handler_test.go` | Unit tests for all 8 handlers |
| `test/integration/dashboard_integration_test.go` | Integration tests with database |
| `test/service/dashboard_service_test.go` | Service layer unit tests |

---

### 3. Dependencies

**Existing Dependencies (No New Setup Required):**
- ✅ `internal/repository/dashboard_repository.go` - Interface already exists
- ✅ `internal/adapter/postgres/dashboard_repository.go` - Implementation already exists
- ✅ `internal/domain/dto/dashboard_response.go` - Response DTOs already exist
- ✅ `internal/service/user_config_service.go` - Config service already exists
- ✅ `internal/config/config.go` - Main config with timezone support

**Prerequisites to Verify:**
1. Database must be accessible and contain test data
2. Run `go mod tidy` after creating new files to resolve imports
3. Ensure `internal/api/v1/handlers/` directory exists
4. Verify `internal/service/dashboard/` directory structure

**Blocking Issues:**
- None identified - all foundational pieces are in place from Phase 1

---

### 4. Code Patterns

**Pattern 1: Handler Structure (from projects_handler.go)**
```go
type DashboardHandler struct {
    repo        repository.DashboardRepository
    userConfig  *service.UserConfigService
}

func NewDashboardHandler(repo repository.DashboardRepository, userConfig *service.UserConfigService) *DashboardHandler {
    return &DashboardHandler{
        repo:       repo,
        userConfig: userConfig,
    }
}
```

**Pattern 2: Error Handling (from projects_handler.go)**
```go
// Success response with JSON:API envelope
envelope := dto.NewJSONAPIEnvelope(dto.JSONAPIData{
    Type:       "dashboard_day",
    ID:         strconv.FormatInt(time.Now().Unix(), 10),
    Attributes: stats,
})
w.Header().Set("Content-Type", "application/vnd.api+json")
json.NewEncoder(w).Encode(envelope)

// Error response
http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
```

**Pattern 3: Service Method Structure**
```go
func (s *DayService) CalculateDailyStats(ctx context.Context, date time.Time) (*dto.DailyStats, error) {
    stats, err := s.repo.GetDailyStats(ctx, date)
    if err != nil {
        return nil, fmt.Errorf("failed to get daily stats: %w", err)
    }
    // Apply business logic calculations
    return stats, nil
}
```

**Pattern 4: Response Formatting**
```go
// Round to 3 decimal places for consistency
func roundToThreeDecimals(val float64) float64 {
    return math.Round(val*1000) / 1000
}

// Handle zero division gracefully
func safeDivide(numerator, denominator int) float64 {
    if denominator == 0 {
        return 0.0
    }
    return float64(numerator) / float64(denominator)
}
```

---

### 5. Testing Strategy

**Unit Tests Pattern (from projects_handler_test.go):**
```go
func TestDashboardHandler_DayEndpoint(t *testing.T) {
    // Setup mock repository
    mockRepo := &mocks.DashboardRepository{}
    mockRepo.On("GetDailyStats", mock.Anything, mock.Anything).
        Return(dto.NewDailyStats(100, 5), nil)
    
    // Setup user config service
    userConfig := service.NewUserConfigService(service.GetDefaultConfig())
    
    // Create handler
    handler := handlers.NewDashboardHandler(mockRepo, userConfig)
    
    // Make request
    req := httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json", nil)
    recorder := httptest.NewRecorder()
    handler.Day(recorder, req)
    
    // Verify response
    assert.Equal(t, http.StatusOK, recorder.Code)
    // ... additional assertions
}
```

**Integration Tests Pattern (from projects_integration_test.go):**
```go
func TestDashboardIntegration_DayEndpoint(t *testing.T) {
    ctx := integration.Setup(t)
    defer ctx.Teardown(t)
    
    // Create test data
    projectID := ctx.CreateTestProject(t)
    ctx.CreateTestLog(t, projectID)
    
    // Make request
    recorder := ctx.MakeRequest(t, httptest.NewRequest(http.MethodGet, "/v1/dashboard/day.json", nil))
    
    // Verify response
    assert.Equal(t, http.StatusOK, recorder.Code)
    // ... verify calculated values
}
```

**Test Coverage Requirements:**
- Each handler tests success path with mock data
- Each handler tests error paths (nil repo, invalid config)
- Integration tests verify database queries work correctly
- Edge cases: empty data, zero values, null timestamps
- Concurrent request testing for performance validation

---

### 6. Risks and Considerations

**Known Risks:**

1. **Database Query Performance**
   - Risk: Complex aggregations may be slow on large datasets
   - Mitigation: Add database indexes (already done per RDL-028), implement caching if needed in Phase 3
   - Impact: Low - Current implementation uses COALESCE and simple aggregations

2. **Timezone Handling**
   - Risk: Date calculations may differ from Rails if timezone not configured
   - Mitigation: Use `config.Config.TZLocation` consistently, document default behavior
   - Impact: Medium - Users must configure TZ_LOCATION for exact parity

3. **Fault Counting Ambiguity**
   - Risk: PRD says "count ALL faults" but actual fault criteria unclear
   - Mitigation: Current implementation counts all logs as faults; clarify with stakeholder if needed
   - Impact: Medium - May require query adjustment after clarification

4. **Empty Data Handling**
   - Risk: Dashboard returns errors when no data exists
   - Mitigation: Return zero values wrapped in successful response (already implemented in repository)
   - Impact: Low - Graceful degradation already designed

5. **Route Registration Order**
   - Risk: Dashboard routes may conflict with existing patterns
   - Mitigation: Use `/v1/dashboard/` prefix, register before specific project/log routes
   - Impact: Low - Clear namespace separation

**Design Trade-offs:**

| Decision | Rationale |
|----------|-----------|
| No dedicated service layer for Phase 1 | Reduces complexity; can refactor to `internal/service/dashboard/` if scope grows |
| Direct repository calls in handlers | Matches existing `projects_handler.go` pattern; simpler for Phase 1 |
| JSON:API envelope for all responses | Consistent with existing API contract |
| Default config fallback | Graceful degradation when config file missing |

**Deployment Considerations:**
- No database migrations required (uses existing schema)
- Configuration file (`dashboard_config.yaml`) must be created or defaults will apply
- Restart required after code changes (no hot reload)
- Monitor logs for dashboard-specific errors

---

### Implementation Checklist

Before coding, verify:
- [ ] `internal/api/v1/handlers/` directory exists
- [ ] `internal/service/dashboard/` directory created
- [ ] Run `go mod tidy` to ensure all imports resolve
- [ ] Database test data available for integration tests
- [ ] `go fmt` and `go vet` pass on existing code

After coding, verify:
- [ ] All 8 endpoints respond with 200 OK (with test data)
- [ ] Error responses return correct status codes (400, 404, 500)
- [ ] Response JSON matches expected format in PRD
- [ ] `go test ./...` passes with >85% coverage
- [ ] `go fmt` and `go vet` pass without errors
- [ ] Documentation updated in QWEN.md
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
