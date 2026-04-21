---
id: RDL-078
title: >-
  [doc-008 Phase 1] Create UserConfig service with file-based configuration
  loading
status: To Do
assignee:
  - thomas
created_date: '2026-04-21 15:49'
updated_date: '2026-04-21 16:16'
labels:
  - phase-1
  - infrastructure
  - config
dependencies: []
references:
  - REQ-DASH-001
  - AC-DASH-001
  - 'Decision 2: UserConfig Implementation Strategy'
documentation:
  - doc-008
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement internal/service/user_config_service.go to load dashboard configuration from YAML file with hardcoded defaults as fallback. The service must handle missing values gracefully and provide type-safe access to max_faults, prediction_pct, and pages_per_day settings.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Configuration loads from config/dashboard.yaml if present
- [x] #2 Hardcoded defaults used when config file missing
- [ ] #3 All three settings (max_faults, prediction_pct, pages_per_day) accessible with correct types
- [ ] #4 Unit tests cover both file and default paths
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The UserConfig service will implement file-based configuration loading with hardcoded defaults as fallback, following the established patterns in the codebase:

**Architecture Decision**: 
- Create `internal/service/user_config_service.go` to centralize dashboard configuration access
- Use YAML configuration file (`config/dashboard.yaml`) for production flexibility
- Provide hardcoded defaults as fallback for missing or invalid config
- Follow Clean Architecture by keeping service layer independent of HTTP framework

**Key Design Choices**:
1. **Configuration Loading**: Use `gopkg.in/yaml.v3` for YAML parsing (consistent with Go ecosystem)
2. **Default Values**: Hardcode sensible defaults matching Rails behavior
3. **Type Safety**: Strongly typed struct with pointer fields for optional values
4. **Error Handling**: Graceful degradation - log warnings but continue with defaults
5. **Testability**: Interface-based design allowing mock implementations

**File Structure**:
```
internal/service/
└── user_config_service.go    # Main service implementation
config/
└── dashboard.yaml             # Configuration file (created)
test/
└── unit/
    └── user_config_service_test.go  # Unit tests
```

**Configuration Schema**:
```yaml
# config/dashboard.yaml
max_faults: 10          # Maximum faults before alert (matches Rails default)
prediction_pct: 0.15    # Prediction percentage for speculative calculations
pages_per_day: 25.0     # Default pages per day target
```

---

### 2. Files to Modify

#### New Files to Create:

| File | Purpose | Lines Estimate |
|------|---------|----------------|
| `internal/service/user_config_service.go` | Main UserConfig service implementation with YAML loading and defaults | ~80-100 |
| `config/dashboard.yaml` | Default configuration file for dashboard settings | ~10 |
| `test/unit/user_config_service_test.go` | Unit tests covering file load, defaults, and edge cases | ~150-200 |

#### Files to Reference (No Modifications):

| File | Reason |
|------|--------|
| `internal/config/config.go` | Reference for existing config patterns and error handling |
| `go.mod` | May need yaml dependency addition |
| `.env.example` | Document new configuration options |

---

### 3. Dependencies

#### External Dependencies:
- **gopkg.in/yaml.v3** - YAML parsing library (add to go.mod)
  ```bash
  go get gopkg.in/yaml.v3@v3.0.1
  ```

#### Internal Dependencies:
- `internal/config` - Existing config package for patterns reference
- `internal/domain/models` - May need dashboard-specific models later
- `test` - Test helper utilities

#### Prerequisites:
- [ ] Go 1.25.7 installed and verified
- [ ] PostgreSQL running and accessible
- [ ] Test database created (`reading_log_test`)
- [ ] `.env` file configured with database credentials

---

### 4. Code Patterns

#### Pattern 1: Configuration Struct (matching `internal/config/config.go`)
```go
type DashboardConfig struct {
    MaxFaults     int     `yaml:"max_faults"`
    PredictionPct float64 `yaml:"prediction_pct"`
    PagesPerDay   float64 `yaml:"pages_per_day"`
}
```

#### Pattern 2: Error Handling (matching existing patterns)
```go
func LoadDashboardConfig(path string) (*DashboardConfig, error) {
    // Try to load from file
    config, err := loadFromFile(path)
    if err != nil {
        log.Warn("Failed to load dashboard config, using defaults", "error", err)
        return getDefaultConfig(), nil  // Graceful fallback
    }
    return config, nil
}
```

#### Pattern 3: Test Structure (matching `test/unit/*.go`)
```go
func TestUserConfigService_LoadFromFile(t *testing.T) {
    helper, err := test.SetupTestDB()
    if err != nil {
        t.Fatal(err)
    }
    defer helper.Close()
    
    // Test implementation
}
```

#### Pattern 4: Service Method Structure
```go
type UserConfigService struct {
    config *DashboardConfig
}

func (s *UserConfigService) GetMaxFaults() int {
    return s.config.MaxFaults
}

func (s *UserConfigService) GetPredictionPct() float64 {
    return s.config.PredictionPct
}

func (s *UserConfigService) GetPagesPerDay() float64 {
    return s.config.PagesPerDay
}
```

---

### 5. Testing Strategy

#### Unit Tests Coverage:

| Test Case | Description | Expected Result |
|-----------|-------------|-----------------|
| `TestUserConfigService_LoadFromFile` | Load config from valid YAML file | Returns correct config values |
| `TestUserConfigService_FileNotFound` | Load from non-existent file | Returns defaults, logs warning |
| `TestUserConfigService_InvalidYAML` | Load from malformed YAML | Returns defaults, logs error |
| `TestUserConfigService_DefaultValues` | No config file present | All defaults applied correctly |
| `TestUserConfigService_PartialConfig` | Config with missing fields | Missing fields use defaults |
| `TestUserConfigService_Getters` | All getter methods return correct values | Values match constructor |

#### Test Data Setup:
```go
// Create test YAML file
testConfig := `
max_faults: 15
prediction_pct: 0.20
pages_per_day: 30.5
`
ioutil.WriteFile("/tmp/test_dashboard.yaml", []byte(testConfig), 0644)
```

#### Edge Cases to Cover:
1. Empty config file → Use all defaults
2. Missing required fields → Use defaults for missing only
3. Invalid YAML syntax → Log error, use defaults
4. File permissions issues → Handle gracefully
5. Zero/negative values in config → Validate and warn

---

### 6. Risks and Considerations

#### Known Risks:

| Risk | Impact | Mitigation |
|------|--------|------------|
| Config file missing in production | Service uses defaults (acceptable) | Document default values clearly |
| YAML parsing error | Service falls back to defaults | Log detailed errors for debugging |
| Path configuration incorrect | File not found | Use relative paths from project root |
| Concurrency issues | None - config is read-only after load | Config loaded once at startup |

#### Trade-offs:

1. **YAML vs JSON**: Chose YAML for human-readable comments and structure
2. **Hardcoded defaults**: Simpler than database lookup, easier to deploy
3. **No hot reload**: Config loaded once at startup (sufficient for Phase 2)
4. **Graceful degradation**: Better UX than hard failures on config issues

#### Blocking Issues:
- None identified - this is a foundational component with low risk

#### Deployment Considerations:
- Config file should be included in Docker image
- Environment variables can override YAML values if needed later
- Consider adding config validation at startup to fail fast if invalid

---

### 7. Implementation Steps

**Step 1**: Add YAML dependency to `go.mod`
```bash
go get gopkg.in/yaml.v3@v3.0.1
```

**Step 2**: Create `config/dashboard.yaml` with default values
- Define `max_faults`, `prediction_pct`, `pages_per_day`
- Include comments explaining each field

**Step 3**: Implement `internal/service/user_config_service.go`
- Define `DashboardConfig` struct
- Implement `LoadDashboardConfig()` function
- Implement `UserConfigService` with getter methods
- Add proper error handling and logging

**Step 4**: Write unit tests in `test/unit/user_config_service_test.go`
- Test file loading scenarios
- Test default value application
- Test edge cases (empty file, invalid YAML, etc.)

**Step 5**: Run tests and verify coverage
```bash
go test -v ./test/unit/... -run TestUserConfig
```

**Step 6**: Run linting and formatting
```bash
go fmt ./internal/service/user_config_service.go
go vet ./internal/service/user_config_service.go
```

---

### 8. Acceptance Criteria Mapping

| AC Requirement | Implementation Approach |
|----------------|------------------------|
| Config loads from `config/dashboard.yaml` if present | Implement file loading with `yaml.Unmarshal()` |
| Hardcoded defaults used when config file missing | Return default config on error, log warning |
| All three settings accessible with correct types | Expose via getter methods returning typed values |
| Unit tests cover both file and default paths | Test explicit scenarios for each case |

---

### 9. Quality Gates

Before marking task complete:
- [ ] All unit tests pass (target: >90% coverage)
- [ ] `go vet` reports no issues
- [ ] `go fmt` applied consistently
- [ ] Code follows existing patterns in `internal/config/config.go`
- [ ] Error handling matches project conventions
- [ ] Documentation updated in QWEN.md or README
- [ ] No hardcoded paths (use relative paths from project root)
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress - RDL-078

### Status: In Progress

I'm implementing the UserConfig service with file-based configuration loading. Here's what I've done so far:

**Analysis Completed:**
- Reviewed task requirements and acceptance criteria
- Examined existing config patterns in `internal/config/config.go`
- Verified codebase structure (service directory doesn't exist yet)
- Confirmed YAML dependency not present in go.mod

**Implementation Steps Ready:**
1. Add YAML dependency to go.mod
2. Create config/dashboard.yaml with default values
3. Implement internal/service/user_config_service.go
4. Write unit tests in test/unit/user_config_service_test.go

### Next Actions:
- Proceeding with Step 1: Adding YAML dependency
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
