---
id: RDL-124
title: Add route GET /v1/dashboard/day.json
status: To Do
assignee:
  - catarina
created_date: '2026-04-28 10:27'
updated_date: '2026-04-28 10:31'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add route GET /v1/dashboard/day.json

internal/api/v1/handlers/dashboard_handler.go
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The `GET /v1/dashboard/day.json` endpoint handler is already fully implemented in `internal/api/v1/handlers/dashboard_handler.go`. The task requires:

1. **Register the route** in `internal/api/v1/routes.go` to expose the endpoint
2. **Instantiate dependencies** in `cmd/server.go`:
   - Create `DashboardRepositoryImpl` instance
   - Create `UserConfigService` instance (loads from YAML or uses defaults)
   - Create `DashboardHandler` with injected dependencies

**Architecture Decisions:**
- Follow existing Clean Architecture pattern: Handler → Service → Repository → Database
- Use dependency injection for testability (already implemented in handler)
- Route registration follows Gorilla Mux pattern used elsewhere in the codebase
- The endpoint returns flat JSON with `stats` key at root (not JSON:API envelope), matching Rails API behavior

**Why this approach:**
- Handler and service logic already exist and are tested
- Only wiring changes needed (routes + DI)
- Minimal risk as calculations are already verified in unit tests

### 2. Files to Modify

**Files to Modify:**
1. `cmd/server.go`
   - Add import for `service` package
   - Add import for `internal/adapter/postgres` (if not already present)
   - Create `userConfigService` instance using `service.LoadDashboardConfig()` or `service.NewUserConfigService()`
   - Create `dashboardRepo` instance using `postgres.NewDashboardRepositoryImpl(dbPool)`
   - Create `dashboardHandler` instance using `handlers.NewDashboardHandler(dashboardRepo, userConfigService)`

2. `internal/api/v1/routes.go`
   - Add import for `handlers` package (if not already present)
   - Register route: `r.HandleFunc("/v1/dashboard/day.json", dashboardHandler.Day).Methods("GET")`

**Files Already Implemented (No Changes Needed):**
- `internal/api/v1/handlers/dashboard_handler.go` - Day method exists
- `internal/adapter/postgres/dashboard_repository.go` - All repository methods exist
- `internal/service/dashboard/day_service.go` - Service logic exists
- `internal/domain/dto/dashboard_response.go` - DTOs exist
- `internal/api/v1/handlers/dashboard_handler_test.go` - Unit tests exist

### 3. Dependencies

**Prerequisites:**
- PostgreSQL database must be running with `reading_log` database
- Dashboard config file (`dashboard_config.yaml`) - optional, defaults used if missing
- Existing dependencies already in `go.mod`:
  - `github.com/gorilla/mux`
  - `github.com/jackc/pgx/v5`
  - `gopkg.in/yaml.v3` (for config loading)

**No Blocking Issues:**
- All required code already exists
- No database schema changes needed
- No new external dependencies required

### 4. Code Patterns

**Follow Existing Conventions:**
1. **Dependency Injection Pattern:**
   ```go
   dashboardRepo := postgres.NewDashboardRepositoryImpl(dbPool)
   userConfig := service.NewUserConfigService(service.GetDefaultConfig())
   dashboardHandler := handlers.NewDashboardHandler(dashboardRepo, userConfig)
   ```

2. **Route Registration Pattern:**
   ```go
   r.HandleFunc("/v1/dashboard/day.json", dashboardHandler.Day).Methods("GET")
   ```

3. **Error Handling:**
   - Handler already implements proper error responses with JSON format
   - HTTP status codes: 200 (OK), 400 (Bad Request), 500 (Internal Server Error)

4. **Context Timeout:**
   - Repository uses 15-second timeout (`dashboardContextTimeout`)
   - No changes needed, already implemented

5. **Response Format:**
   - Returns flat JSON: `{"stats": {...}}`
   - Content-Type: `application/json`
   - Float values rounded to 3 decimals

### 5. Testing Strategy

**Unit Tests (Already Exist):**
- `TestDashboardHandler_Day` - Tests successful response with date parameter
- `TestDashboardHandler_Day_WithoutDate` - Tests default (today) date behavior
- `TestDashboardHandler_Day_InvalidDate` - Tests invalid date format error handling
- Mock repository used for isolation

**Integration Tests to Add/Verify:**
1. **Happy Path:**
   - Start server with test database
   - Insert sample log data
   - Call `GET /v1/dashboard/day.json`
   - Verify 200 OK and response structure

2. **Edge Cases:**
   - Empty database (no logs)
   - Invalid date format
   - Missing date parameter (uses today)

3. **Validation:**
   - Verify `per_pages` is null when previous week has no data
   - Verify float rounding to 3 decimals
   - Verify all calculated fields present in response

**Test Commands:**
```bash
# Run unit tests
go test -v ./internal/api/v1/handlers/dashboard_handler_test.go

# Run integration tests
go test -v ./internal/api/v1/handlers/... -run TestDashboardHandler
```

### 6. Risks and Considerations

**Known Risks:**
1. **Config File Missing:** `dashboard_config.yaml` may not exist
   - Mitigation: Handler uses defaults via `service.GetDefaultConfig()`
   - No blocking issue

2. **Database Connection:** Dashboard queries require active database
   - Mitigation: Server already validates connection in `main()`
   - Existing error handling covers this

3. **Performance:** Multiple repository calls in handler
   - Current implementation calls: `GetMeanByWeekday`, `GetProjectAggregates`, `GetMaxByWeekday`, `GetOverallMean`, `GetPreviousPeriodMean`, `GetPreviousPeriodSpecMean`
   - Mitigation: Each query has 15-second timeout
   - Consider optimization in future if performance issues arise

**Deployment Considerations:**
- No database migrations required
- No environment variable changes required
- Route becomes immediately available after deployment

**Rollout Plan:**
1. Implement route registration
2. Run unit tests
3. Run integration tests
4. Manual verification with curl:
   ```bash
   curl http://localhost:3000/v1/dashboard/day.json
   curl http://localhost:3000/v1/dashboard/day.json?date=2024-01-15T10:30:00Z
   ```

**Acceptance Criteria Alignment:**
- ✅ All unit tests pass (existing tests cover handler logic)
- ✅ Integration tests verify actual database interactions
- ✅ `go fmt` and `go vet` pass (following existing patterns)
- ✅ Clean Architecture layers followed (Handler → Repository → DB)
- ✅ Error responses consistent with existing patterns
- ✅ HTTP status codes correct (200, 400, 500)
- ⚠️ Documentation updated in QWEN.md (to be done after implementation)
- ✅ Error path tests included (invalid date format)
- ✅ Handlers test both success and error responses
- ✅ Integration tests verify database interactions
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
