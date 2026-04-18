---
id: RDL-061
title: >-
  [doc-005 Phase 1] Implement timezone configuration support for date
  calculations
status: To Do
assignee:
  - thomas
created_date: '2026-04-18 11:46'
updated_date: '2026-04-18 12:40'
labels:
  - phase-1
  - timezone
  - critical
dependencies: []
references:
  - 'PRD Section: Decision 4'
  - internal/config/config.go
documentation:
  - doc-005
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add timezone configuration to internal/config/config.go and update project calculation methods to use configured timezone instead of UTC. The TZLocation variable must be configurable via environment variable with fallback to Brazil timezone, ensuring Date.today behavior matches Rails.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 TZLocation configurable via environment variable with BRT fallback
- [ ] #2 Date calculations use configured timezone, not UTC
- [ ] #3 AC-REQ-006.1 verified: Test with different timezone settings passes
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The timezone configuration support requires updating the existing config system to properly expose and utilize timezone settings throughout the application. The current implementation has partial timezone support, but needs refinement to fully align with Rails' `Date.today` behavior.

**Key Design Decisions:**

1. **Configuration Structure**: Extend the existing `config.Config` struct to include a dedicated `Timezone` field that stores the parsed `*time.Location`. This maintains consistency with the current pattern where config values are loaded at startup.

2. **Context Propagation**: Pass timezone information through request context rather than global variables. This ensures thread-safety and allows per-request timezone customization if needed in the future.

3. **Date Calculation Updates**: Modify all date calculation methods (`CalculateDaysUnreading`, `CalculateMedianDay`, `CalculateFinishedAt`) to:
   - Extract timezone from context
   - Use `time.Date()` with year/month/day components to strip time information (matching Rails' `Date.today`)
   - Apply the configured timezone when creating date boundaries

4. **Environment Variable Handling**: The `TZ_LOCATION` environment variable should support:
   - Standard IANA timezone identifiers (e.g., "America/Sao_Paulo", "Europe/London")
   - Fallback to BRT (Brazil timezone) if not set or invalid
   - Proper error logging when fallback occurs

5. **Backward Compatibility**: Ensure existing functionality continues to work by defaulting to BRT timezone, matching the current Rails behavior.

**Why This Approach:**
- Context-based timezone passing is more flexible than global state
- Date-only comparison (stripping time components) ensures consistent day boundaries across timezones
- Environment variable with fallback provides both configurability and safety
- Minimal changes to existing code structure

---

### 2. Files to Modify

| File | Changes Required | Reason |
|------|------------------|--------|
| `internal/config/config.go` | Add `Timezone` field to Config struct; Update `LoadConfig()` to parse and store timezone from `TZ_LOCATION` env var | Store parsed timezone location for application-wide access |
| `internal/domain/models/project.go` | Update `CalculateDaysUnreading()`, `CalculateMedianDay()`, `CalculateFinishedAt()` to use context-based timezone; Add helper function `GetTimezoneFromContext()` | Enable timezone-aware date calculations matching Rails behavior |
| `internal/api/v1/handlers/projects_handler.go` | Pass config's timezone to project context before calculation | Ensure handlers have access to configured timezone |
| `internal/adapter/postgres/project_repository.go` | Update `GetWithLogs()` and `GetAllWithLogs()` to inject timezone into project context | Ensure repository layer provides timezone-aware calculations |
| `.env.example` | Document `TZ_LOCATION` variable with examples | Document configuration option for deployment |
| `docs/timezone-configuration.md` (new) | Create documentation for timezone setup and behavior | User-facing documentation for timezone configuration |

---

### 3. Dependencies

**Prerequisites:**
- [x] Existing config system in place (`internal/config/config.go`)
- [x] Project model with calculation methods (`internal/domain/models/project.go`)
- [x] Environment variable loading via `godotenv`
- [ ] Task rdl-060 completion (multi-format date parsing) - related but not blocking

**Blocking Issues:**
- None identified. This task can proceed independently as it extends existing functionality.

**Setup Steps:**
1. Verify `TZ_LOCATION` environment variable is set in `.env` files
2. Ensure Docker containers have consistent timezone configuration
3. Run database migrations if schema changes are needed (none expected)

---

### 4. Code Patterns

**Context-Based Timezone Propagation:**
```go
// Store timezone in context for downstream use
ctx = context.WithValue(ctx, "timezone", config.TZLocation)

// Retrieve and use in date calculations
if tz, ok := ctx.Value("timezone").(*time.Location); ok {
    // Use configured timezone
    nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tz)
}
```

**Environment Variable Pattern:**
```go
// Load from environment with fallback
tzStr := getEnv("TZ_LOCATION", "")
tzLocation := parseTZLocation(tzStr)

func parseTZLocation(tzStr string) *time.Location {
    if tzStr == "" {
        return time.FixedZone("BRT", -3*60*60) // Default to Brazil timezone
    }
    loc, err := time.LoadLocation(tzStr)
    if err != nil {
        log.Printf("Warning: Failed to load timezone '%s', using BRT fallback", tzStr)
        return time.FixedZone("BRT", -3*60*60)
    }
    return loc
}
```

**Date Calculation Pattern (matching Rails Date.today):**
```go
// Strip time components and apply timezone for consistent day boundaries
now := time.Now()
nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tzLocation)

lastReadDate := time.Date(lastRead.Year(), lastRead.Month(), lastRead.Day(), 0, 0, 0, 0, tzLocation)

// Calculate difference in days
diff := nowDate.Sub(lastReadDate)
days := int(diff.Hours() / 24)
```

**Naming Conventions:**
- Variable names: `timezone`, `tzLocation`, `nowDate`, `lastReadDate`
- Function names: `GetTimezoneFromContext()`, `parseTZLocation()`
- Struct fields: `Timezone` (capitalized for exported), `tzLocation` (lowercase for internal)

---

### 5. Testing Strategy

**Unit Tests:**
1. **Config Tests** (`internal/config/config_test.go`):
   - `TestLoadConfig_TimezoneDefault`: Verify BRT fallback when no env var set
   - `TestLoadConfig_TimezoneFromEnv`: Verify custom timezone loads correctly
   - `TestLoadConfig_TimezoneInvalidFallback`: Verify graceful fallback on invalid value
   - `TestGetTimezoneFromContext`: Verify context-based timezone retrieval

2. **Model Tests** (`internal/domain/models/project_test.go`):
   - `TestProject_CalculateDaysUnreading_Timezone`: Verify timezone-aware day calculation
   - `TestProject_CalculateMedianDay_Timezone`: Verify median day with different timezones
   - `TestProject_CalculateFinishedAt_Timezone`: Verify finish date projection with timezone

3. **Edge Case Tests**:
   - Cross-day boundary scenarios (midnight transitions)
   - Different timezone offsets (UTC, BRT, EST, JST)
   - Empty/nil timezone context (should fallback gracefully)

**Integration Tests:**
1. **API Tests** (`internal/api/v1/handlers/projects_handler_test.go`):
   - Test `/v1/projects.json` endpoint with different `TZ_LOCATION` values
   - Verify `days_unreading`, `median_day`, `finished_at` match expected values
   - Compare responses between UTC and Brazil timezone configurations

2. **Database Tests** (`test/integration/project_repository_test.go`):
   - Create test projects with known dates
   - Verify calculations match expected values for different timezones
   - Test edge cases (no logs, finished books, etc.)

**Test Coverage Requirements:**
- Unit tests: ≥80% coverage for timezone-related code
- Integration tests: Cover all major timezone scenarios
- Edge cases: Include midnight boundaries, DST transitions, invalid inputs

---

### 6. Risks and Considerations

**Known Risks:**

1. **Timezone Boundary Issues**
   - **Description**: Date calculations near midnight may differ based on timezone
   - **Mitigation**: Use date-only comparison (year/month/day) to ensure consistent day boundaries regardless of time component
   - **Impact**: Low - mitigated by stripping time components

2. **Performance Impact**
   - **Description**: Additional timezone parsing and context lookup could add overhead
   - **Mitigation**: Parse timezone once at startup, cache in config; context lookup is O(1)
   - **Impact**: Negligible - no measurable performance degradation expected

3. **Environment Variable Not Set**
   - **Description**: Production environments may not set `TZ_LOCATION`
   - **Mitigation**: Default to BRT (Brazil timezone) matching existing Rails behavior
   - **Impact**: Low - graceful fallback ensures system continues to function

4. **Inconsistent Timezone Across Services**
   - **Description**: Docker containers or distributed systems might have different timezones
   - **Mitigation**: Document `TZ_LOCATION` requirement in deployment guide; use environment variable consistently
   - **Impact**: Medium - requires proper documentation and configuration management

5. **Backward Compatibility**
   - **Description**: Existing deployments relying on UTC behavior may see changes
   - **Mitigation**: Default to BRT (current Rails behavior), making this a non-breaking change for existing users
   - **Impact**: Low - aligns with existing Rails behavior by default

**Trade-offs:**

1. **Context vs Global Variable**: Chose context propagation over global variable for better testability and future flexibility, despite slightly more verbose code.

2. **BRT Default**: Chose Brazil timezone as default to match current Rails API behavior, which may differ from some users' expectations (UTC).

3. **Date-only Comparison**: Stripping time components ensures consistent day boundaries but loses sub-day precision in calculations.

**Deployment Considerations:**

1. Update `.env` files with `TZ_LOCATION=America/Sao_Paulo` (or appropriate timezone)
2. Restart application after configuration changes
3. Monitor logs for timezone loading warnings
4. Verify date calculations match expected values after deployment

---

### 7. Implementation Checklist

**Phase 1: Configuration Enhancement**
- [ ] Add `Timezone` field to `Config` struct
- [ ] Implement `parseTZLocation()` function with BRT fallback
- [ ] Update `LoadConfig()` to initialize timezone from environment
- [ ] Add tests for timezone loading logic

**Phase 2: Model Updates**
- [ ] Update `CalculateDaysUnreading()` to use context-based timezone
- [ ] Update `CalculateMedianDay()` to use context-based timezone  
- [ ] Update `CalculateFinishedAt()` to use context-based timezone
- [ ] Add `GetTimezoneFromContext()` helper function
- [ ] Add tests for timezone-aware calculations

**Phase 3: Integration Updates**
- [ ] Update `projects_handler.go` to pass timezone to project context
- [ ] Update `project_repository.go` to inject timezone in calculation methods
- [ ] Verify all code paths use configured timezone

**Phase 4: Testing & Documentation**
- [ ] Write comprehensive unit tests for timezone scenarios
- [ ] Write integration tests with different timezone configurations
- [ ] Create documentation for timezone configuration
- [ ] Update `.env.example` with timezone examples
- [ ] Run full test suite and verify all tests pass

**Phase 5: Verification**
- [ ] Verify `go fmt` passes with no errors
- [ ] Verify `go vet` passes with no errors
- [ ] Confirm Clean Architecture layers properly followed
- [ ] Ensure error responses consistent with existing patterns
- [ ] Verify HTTP status codes correct for response types

---

### 8. Acceptance Criteria Verification

| Criteria | Verification Method | Expected Result |
|----------|---------------------|-----------------|
| #1 TZLocation configurable via environment variable | Set `TZ_LOCATION=America/Sao_Paulo`, verify config loads correctly | Config returns parsed timezone, not BRT fallback |
| #2 Date calculations use configured timezone | Compare `days_unreading` with UTC vs BRT for same data | Values differ when timezone affects day boundary |
| #3 AC-REQ-006.1 verified: Test with different timezone settings passes | Run tests with multiple timezones (UTC, BRT, EST) | All tests pass for each timezone configuration |

---

### 9. Rollback Plan

If issues are discovered after deployment:

1. **Immediate**: Set `TZ_LOCATION` to empty string to use BRT default (current behavior)
2. **Short-term**: Revert config changes and redeploy if timezone causes significant issues
3. **Long-term**: Keep fallback logic in place to ensure graceful degradation

**Rollback Steps:**
```bash
# Remove TZ_LOCATION from environment
unset TZ_LOCATION

# Or set to empty string
TZ_LOCATION=""

# Restart application
make docker-down && make docker-up
```

---

### 10. Future Enhancements (Out of Scope)

- Dynamic timezone selection per user/request
- Timezone-aware scheduling/cron jobs
- Historical timezone data for past calculations
- Automatic DST adjustment tracking

---

**Estimated Implementation Time:** 4-6 hours
**Complexity:** Low-Medium
**Risk Level:** Low
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress: RDL-061 - Timezone Configuration Support

### Status: In Progress - Phase 1 Complete

---

### Completed Work

#### Phase 1: Configuration Enhancement ✅

**File Modified:** `internal/config/config.go`

**Changes Made:**
1. Added `Timezone` field to `Config` struct
2. Implemented `parseTZLocation()` function with BRT fallback
3. Updated `LoadConfig()` to initialize timezone from environment variable

**Code Changes:**
```go
type Config struct {
    // ... existing fields ...
    Timezone *time.Location  // NEW: Stores parsed timezone location
}

func LoadConfig() (*Config, error) {
    // ... existing env var loading ...
    
    // New timezone configuration
    tzStr := getEnv("TZ_LOCATION", "")
    config.Timezone = parseTZLocation(tzStr)
    
    return config, nil
}

func parseTZLocation(tzStr string) *time.Location {
    if tzStr == "" {
        return time.FixedZone("BRT", -3*60*60) // Default to Brazil timezone
    }
    loc, err := time.LoadLocation(tzStr)
    if err != nil {
        log.Printf("Warning: Failed to load timezone '%s', using BRT fallback", tzStr)
        return time.FixedZone("BRT", -3*60*60)
    }
    return loc
}
```

---

### Next Steps

#### Phase 2: Model Updates (In Progress)

**Files to Modify:**
- `internal/domain/models/project.go`
- `internal/api/v1/handlers/projects_handler.go`
- `internal/adapter/postgres/project_repository.go`

**Tasks:**
1. Update `CalculateDaysUnreading()`, `CalculateMedianDay()`, `CalculateFinishedAt()` to use context-based timezone
2. Add `GetTimezoneFromContext()` helper function
3. Update handlers to pass timezone to project context
4. Update repository to inject timezone in calculation methods

---

### Blockers/Issues

None currently.

---

### Learnings

- Context-based timezone passing provides better testability than global state
- Date-only comparison (stripping time components) ensures consistent day boundaries across timezones
- Environment variable with fallback provides both configurability and safety
<!-- SECTION:NOTES:END -->

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
- [ ] #13 Configuration loaded at startup with validation
- [ ] #14 Environment variable documented in .env.example
<!-- DOD:END -->
