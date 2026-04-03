---
id: RDL-021
title: '[doc-002 Phase 2] Create app configuration for status ranges'
status: To Do
assignee:
  - thomas
created_date: '2026-04-03 14:02'
updated_date: '2026-04-03 16:39'
labels:
  - phase-2
  - configuration
  - setup
dependencies: []
references:
  - 'PRD Section: Technical Decisions - Decision 2: Configuration Values'
  - 'PRD Section: Files to Modify - config.go'
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create a Go configuration structure in `internal/config/config.go` with `em_andamento_range` (7 days default) and `dormindo_range` (14 days default) values matching Rails configuration. Add methods to access these values from status calculation logic.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 em_andamento_range = 7 days
- [x] #2 dormindo_range = 14 days
- [x] #3 Access methods for configuration values
- [x] #4 Configuration loads from environment variables or defaults
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task adds application configuration for status range values to the Go codebase, matching the Rails `V1::UserConfig` behavior.

**Implementation Strategy:**
- Add two new integer fields to the `Config` struct: `EmAndamentoRange` (7 days default) and `DormindoRange` (14 days default)
- Provide environment variable overrides: `EM_ANDAMENTO_RANGE` and `DORMINDO_RANGE`
- Add getter methods `GetEmAndamentoRange()` and `GetDormindoRange()` for accessing values
- Keep default values consistent with Rails configuration (note: Rails uses 8 and 16 days, but task specifies 7 and 14)
- Follow existing config pattern using `getEnvAsInt` helper

**Architecture Alignment:**
- Maintains clean architecture: config is in `internal/config` package
- No breaking changes: new fields are additive only
- Config loaded at application start in `cmd/main.go` and passed to handlers/repositories as needed

**Why This Approach:**
- Simple, follow existing patterns in config.go
- Environment variable support for flexibility
- Accessor methods provide encapsulation for future validation
- No dependency on external configuration management libraries

### 2. Files to Modify

**Existing Files:**
- `internal/config/config.go` - Add configuration fields and getter methods
- `internal/config/config_test.go` - Add tests for new configuration fields

**No New Files Required:**
- Configuration is added to existing structure

### 3. Dependencies

**Prerequisites:**
- None - this is a foundational task for phase 2
- No other tasks must be completed first
- Task RDL-022 (status determination logic) depends on this task

**Related Tasks:**
- RDL-022 - Status determination logic will use these config values
- RDL-023 - days_unreading calculation may reference ranges

**Environment Setup:**
- `.env` file can optionally include `EM_ANDAMENTO_RANGE` and `DORMINDO_RANGE`
- `.env.example` should document these new environment variables

### 4. Code Patterns

**Follow Existing Patterns:**
- Use `getEnvAsInt` helper for integer environment variable parsing
- Default values match Go standard: lowercase with underscores in comments
- Field names: `EmAndamentoRange` (PascalCase for exported fields)
- Getter methods: `GetEmAndamentoRange()` and `GetDormindoRange()`

**Naming Conventions:**
- Config fields: PascalCase (`EmAndamentoRange`, `DormindoRange`)
- Environment variables: UPPER_SNAKE_CASE (`EM_ANDAMENTO_RANGE`, `DORMINDO_RANGE`)
- Getter methods: `Get<Fieldname>()` pattern consistent with Go conventions

**Error Handling:**
- Invalid values fall back to defaults (like existing `DB_PORT`)
- No error returned: config loading is non-blocking for production resilience

### 5. Testing Strategy

**Unit Tests to Add:**
1. `TestLoadConfigDefaultValues` - Verify default values (7 and 14 days)
2. `TestLoadConfigEnvironmentVariables` - Verify env var overrides
3. `TestLoadConfigInvalidRangeValues` - Test invalid values fall back to defaults
4. `TestGetEmAndamentoRange` - Verify getter method returns correct value
5. `TestGetDormindoRange` - Verify getter method returns correct value

**Test Coverage:**
- All existing tests continue to pass
- New tests for configuration fields
- Edge cases: negative numbers, zero, very large values

**Test Execution:**
- Run with `go test ./internal/config/...`
- Verify coverage with `go test -cover ./internal/config/...`
- Run with `make test` for project-wide test suite

### 6. Risks and Considerations

**Configuration Value Discrepancy:**
- ⚠️ Rails uses `em_andamento_range: 8` and `dormindo_range: 16`
- ⚠️ Task specification requests `em_andamento_range: 7` and `dormindo_range: 14`
- **Decision**: Follow task specification (7 and 14) as primary requirement
- **Note**: Document discrepancy in code comments

**Future-Proofing Considerations:**
- Consider making range values configurable at runtime if needed later
- Consider adding validation in getter methods (e.g., positive values only)
- Consider adding a method `IsValidRangeValue(val int) bool` for status calculation logic

**Deployment Considerations:**
- Default values are backward compatible (no database migration needed)
- Environment variables are optional
- No breaking changes to existing API responses at config level
- Status determination (RDL-022) must use these configs to affect behavior

**Documentation:**
- Create/update config documentation in `QWEN.md` or separate config doc
- Document environment variables in `.env.example`
- Add inline comments explaining purpose of each range

### Implementation Checklist (After Approval)

```markdown
- [ ] Add `EmAndamentoRange int` field to Config struct with 7-day default
- [ ] Add `DormindoRange int` field to Config struct with 14-day default
- [ ] Add `getEnvAsInt` calls in `LoadConfig()` for both new config values
- [ ] Add `GetEmAndamentoRange()` method to Config
- [ ] Add `GetDormindoRange()` method to Config
- [ ] Add tests for default values (7 and 14)
- [ ] Add tests for environment variable overrides
- [ ] Add tests for getter methods
- [ ] Update `.env.example` with new environment variables
- [ ] Run `go fmt`, `go vet`, and all tests
- [ ] Update documentation (QWEN.md or config docs)
- [ ] Verify no breaking changes to existing functionality
```
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
2026-04-03: Added EmAndamentoRange (7 days default) and DormindoRange (14 days default) fields to Config struct

2026-04-03: Implemented getEnvAsInt validation to reject negative values - status ranges must be positive

2026-04-03: Added 5 new test functions for status range configuration

2026-04-03: All 130 tests pass across 11 packages

2026-04-03: Build succeeds with no errors or warnings
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [x] #2 All integration tests pass use testing-expert subagent for test execution and verification
- [x] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Database queries optimized with proper indexes
- [ ] #8 Documentation updated in QWEN.md
- [ ] #9 New code paths include error path tests
- [ ] #10 HTTP handlers test both success and error responses
- [ ] #11 Integration tests verify actual database interactions
- [ ] #12 Tests use testing-expert subagent for test execution and verification
<!-- DOD:END -->
