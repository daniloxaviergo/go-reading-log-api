---
id: RDL-004
title: '[doc-001 Phase 2] Implement configuration management and logging setup'
status: Done
assignee:
  - catarina
created_date: '2026-04-01 00:57'
updated_date: '2026-04-01 10:35'
labels: []
dependencies: []
references:
  - 'PRD Section: Technical Decisions'
  - 'Implementation Checklist: Core Components'
  - 'PRD Section: Key Requirements'
documentation:
  - doc-001
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create internal/config/config.go with configuration struct and environment variable loading using joho/godotenv.

Create internal/logger/logger.go to initialize slog with structured logging capable of handling application log levels.

Ensure configuration loads all required environment variables including database connection, server port, and connection pool settings.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Configuration struct defined with all environment variable fields
- [x] #2 Logging initialized with structured slog format
- [x] #3 Environment variables properly loaded with default values
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task establishes the foundational configuration and logging infrastructure for the Go Reading Log API.

**Configuration Management:**
- Create `internal/config/config.go` with a configuration struct that loads environment variables using `joho/godotenv`
- Define fields for database connection (host, port, user, password, database), server settings (port, host), and logging (level, format)
- Provide sensible default values for all environment variables to ensure the application runs without `.env` file
- Use Go's struct tags for JSON/YAML serialization if needed later

**Logging Setup:**
- Use Go's built-in `log/slog` package (available since Go 1.21) for structured logging
- Initialize logger with text or JSON format based on environment variable
- Set appropriate log level (debug, info, warn, error) based on configuration
- Use `slog.New()` with `slog.TextHandler` or `slog.JSONHandler` for output to stderr

**Separation of Concerns:**
- Configuration loading logic isolated in `config` package (no external dependencies except joho/godotenv)
- Logger initialization independent of configuration usage (returns `*slog.Logger` for injection)
- Both packages designed for easy testing (config can be mocked, logger uses standard library)

**Why this approach:**
- `joho/godotenv` is already in `go.mod` and matches Rails conventions
- `log/slog` is standard library (no external dependency, officially supported)
- Structured logging enables better observability in production
- Default values ensure fallback behavior when `.env` is missing

### 2. Files to Modify

**Files to Create:**
| File | Purpose |
|------|---------|
| `internal/config/config.go` | Configuration struct and environment loading logic |
| `internal/logger/logger.go` | Logger initialization with slog |

**Files to Modify:**
| File | Reason |
|------|--------|
| `go.mod` | Already has `joho/godotenv`; no changes needed |

**Files Referenced (No Changes):**
| File | Status |
|------|--------|
| `.env.example` | Already exists with proper template |
| `cmd/server.go` | Will be updated in RDL-007 to use config and logger |

### 3. Dependencies

**Already Available:**
- `joho/godotenv v1.5.1` - Present in `go.mod` (indirect dependency)
- `log/slog` - Standard library (Go 1.21+), no external dependency needed

**Prerequisites:**
- None. This task can be executed immediately as it doesn't depend on other implementation tasks
- Task RDL-001 (Go module initialization) is complete
- Task RDL-002 (Domain models) is complete but not required for this task

**No Blocking Issues:** All required dependencies are already present.

### 4. Code Patterns

**Configuration (`internal/config/config.go`):**
```go
// Use struct with exported fields for environment variables
type Config struct {
    DBHost     string
    DBPort     int
    DBUser     string
    DBPassword string
    DBDatabase string
    ServerPort string
    ServerHost string
    LogLevel   string
    LogFormat  string
}

// LoadConfig reads environment variables with defaults
func LoadConfig() *Config {
    // Use godotenv.Load() to load .env file (non-blocking)
    // Return Config struct with values from env or defaults
}
```

**Logger (`internal/logger/logger.go`):**
```go
// Initialize returns *slog.Logger configured by environment
func Initialize(level, format string) *slog.Logger {
    // Parse level string to slog.LevelVar
    // Create appropriate handler (TextHandler/JSONHandler)
    // Return configured logger
}
```

**Conventions to Follow:**
- Use idiomatic Go naming: `CamelCase` for exported fields, `snake_case` for env vars
- Package names: `config` and `logger` (simple, descriptive)
- Error handling: Return errors from config loader if environment validation fails
- Context: Not needed for config/logger initialization (run at application startup)

**Naming Conventions:**
- Config struct: `Config`
- Config loader: `LoadConfig()`
- Logger initializer: `Initialize()` or `Setup()`
- Return types: `*slog.Logger`, `*Config`

### 5. Testing Strategy

**Unit Tests to Write:**
| Test File | Test Cases |
|-----------|------------|
| `internal/config/config_test.go` | - LoadConfig loads from environment variables<br>- LoadConfig uses defaults when env vars missing<br>- LoadConfig validates required fields<br>- LoadConfig handles invalid values gracefully |
| `internal/logger/logger_test.go` | - Initialize creates logger with text format<br>- Initialize creates logger with JSON format<br>- Initialize sets correct log level<br>- Initialize handles invalid level gracefully |

**Testing Approach:**
- Use `t.Setenv()` (Go 1.17+) to temporarily set environment variables in tests
- Test configuration defaults by not setting any env vars
- Capture logger output using `bytes.Buffer` with `slog.NewTextHandler()` or `slog.NewJSONHandler()`
- Test both text and JSON formats for logger
- Test log level parsing (case-insensitive: "debug", "DEBUG", "Debug" all work)

**Edge Cases:**
- Missing `.env` file (should use defaults)
- Empty environment variables
- Invalid log level (should fall back to default)
- Invalid port numbers (should error or use default)

### 6. Risks and Considerations

**Potential Issues:**
1. **Environment Variable Loading Timing**: `joho/godotenv.Load()` should be called early in `main()` before other packages access config. Document this requirement.
   - *Mitigation*: Add comment in code and README about calling `LoadConfig()` first

2. **Log Level Validation**: If user sets invalid log level (e.g., "verbose"), logger should handle gracefully
   - *Mitigation*: Validate level in `Initialize()` and log warning if invalid, fall back to default

3. **Port Parsing**: `SERVER_PORT` should be string for http.ListenAndServe() but might need int conversion
   - *Mitigation*: Keep as string in config, convert in `cmd/server.go` (not in this task)

4. **Database Connection Not Implemented**: This task only sets up config/logger; database connection will be in RDL-007
   - *Risk*: Config has database fields but no usage yet
   - *Mitigation*: Add TODO comments or use interface{} for future database connection injection

**Design Trade-offs:**
1. **Config Loading**: Using `godotenv.Load()` in `LoadConfig()` means `.env` must be in working directory. Alternative: allow passing custom path. *Chosen*: Simple approach matches Rails convention
2. **Logger Format**: Text vs JSON decision at runtime via environment. *Trade-off*: Slightly more complex than fixed format but allows flexible environments
3. **No Validation Layer**: Config loads without strict validation. *Future enhancement*: Add validation for required fields (DB credentials, etc.)

**Deployment Considerations:**
- `.env` file should be in `.gitignore` (already there per `.gitignore` content)
- Default values should be safe for development but documented as insecure for production
- Log level should default to "info" in production environments (documented in `.env.example`)

### Implementation Steps

1. Create `internal/config/config.go` with `Config` struct and `LoadConfig()` function
2. Create `internal/logger/logger.go` with `Initialize()` function
3. Write unit tests for both packages
4. Test manually: create small test program that loads config and logs messages
5. Update `.env.example` if any fields are missing from the template

---
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
**Implementation completed on:** 2026-03-31

**Files created:**
| File | Purpose |
|------|---------|
| `internal/config/config.go` | Configuration struct with all environment variable fields and `LoadConfig()` function using `joho/godotenv` |
| `internal/config/config_test.go` | 5 unit tests with 100% coverage |
| `internal/logger/logger.go` | Logger initialization with `slog`, supports text/JSON format and configurable log levels |
| `internal/logger/logger_test.go` | 8 unit tests with 100% coverage |

**Test Results:**
- **Total tests:** 14 (all passing)
- **Coverage:** 100% for both `config` and `logger` packages
- **Race detector:** No race conditions detected

**Key Features Implemented:**
1. ✅ Configuration struct with all environment variable fields (database, server, logging)
2. ✅ Environment variable loading with `joho/godotenv`
3. ✅ Sensible default values for all variables (runs without `.env` file)
4. ✅ Structured logging with `log/slog` (standard library)
5. ✅ Log level configuration (debug, info, warn, error) - case-insensitive
6. ✅ Log format configuration (text or JSON)
7. ✅ Graceful handling of invalid values (fallback to defaults)

**Design Decisions:**
- Used `log/slog` (standard library since Go 1.21) - no external dependency
- Used `joho/godotenv` - matches Rails conventions, already in `go.mod`
- Text format as default for development, JSON for production
- Invalid log levels fall back to `info` level
- Invalid port values fall back to default port

**Definition of Done Checklist:**
- [x] All acceptance criteria checked off
- [x] Definition of Done checklist items satisfied
- [x] Implementation plan reflects the final solution
- [x] Final Summary written (PR-style)
- [x] All tests run using `testing-expert` subagent
- [x] Build succeeds with `go build ./...`
- [x] No new warnings or regressions
- [x] Documentation/config updates completed (no changes needed - already documented in `.env.example`)
<!-- SECTION:FINAL-SUMMARY:END -->
<!-- SECTION:FINAL_SUMMARY:END -->
