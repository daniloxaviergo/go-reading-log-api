---
id: RDL-068
title: >-
  [doc-005 Phase 4] Update .env.example with timezone configuration and create
  test data artifacts
status: To Do
assignee:
  - catarina
created_date: '2026-04-18 11:48'
updated_date: '2026-04-18 15:55'
labels:
  - phase-4
  - configuration
  - test-data
dependencies: []
references:
  - 'PRD Section: Configuration Files'
  - .env.example
  - docker-compose.yml
documentation:
  - doc-005
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update .env.example with TZ_LOCATION configuration example, create test data files (project-450-go.json, project-450-rails.json), and ensure docker-compose.yml has consistent timezone across containers.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 TZ_LOCATION documented in .env.example
- [ ] #2 Test data artifacts created for project 450
- [ ] #3 docker-compose.yml ensures consistent timezone
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task focuses on **documentation and test data preparation** for Phase 4 of the API response alignment project (doc-005). The implementation involves three main components:

**A. Timezone Configuration Documentation (.env.example)**
- Add `TZ_LOCATION` environment variable with clear examples
- Document how timezone affects date calculations (matching Rails `Date.today` behavior)
- Provide Brazil timezone as default since this is a Brazilian reading log application
- Include comments explaining the impact on `days_unreading` and `finished_at` calculations

**B. Test Data Artifacts for Project 450**
- Create `test/data/` directory structure
- Capture current API responses for project 450 from both Go and Rails APIs
- Store as JSON files for regression testing reference
- Include both show endpoint and logs endpoint responses

**C. Docker Compose Timezone Consistency**
- Add timezone configuration to all service containers
- Ensure PostgreSQL, Go API, and Rails API use consistent timezone settings
- Document the approach in comments

**Architecture Decision:** This is a documentation/configuration task, not a code change. The goal is to establish the configuration baseline before Phase 4 implementation begins.

---

### 2. Files to Modify

#### Modified Files:

| File | Change Type | Description |
|------|-------------|-------------|
| `.env.example` | Modify | Add `TZ_LOCATION` configuration with examples and documentation |
| `docker-compose.yml` | Modify | Add timezone environment variables to all services |

#### New Files to Create:

| File | Purpose |
|------|---------|
| `test/data/project-450-go.json` | Recorded Go API response for project 450 (show endpoint) |
| `test/data/project-450-rails.json` | Recorded Rails API response for project 450 (show endpoint) |
| `test/data/project-450-go-logs.json` | Recorded Go API response for project 450 logs |
| `test/data/project-450-rails-logs.json` | Recorded Rails API response for project 450 logs |

#### Existing Files to Reference:

| File | Use Case |
|------|----------|
| `test/compare_responses.sh` | Will use these test data files for regression testing |
| `.env.example` (current) | Base for adding timezone configuration |

---

### 3. Dependencies

**Prerequisites:**
- [x] Go API running and accessible on port 3000
- [x] Rails API running and accessible on port 3001  
- [x] Database populated with project ID 450
- [x] `jq` installed for JSON formatting

**Blocking Issues:** None - this is a documentation/preparation task that enables future implementation.

**Setup Steps Required:**
```bash
# Ensure APIs are running
make docker-up

# Verify project 450 exists in database
curl http://localhost:3000/v1/projects/450.json
curl http://localhost:3001/v1/projects/450.json
```

---

### 4. Code Patterns

**Configuration File Pattern (`.env.example`):**
```bash
# Comment block explaining the configuration
# Variable name in ALL_CAPS
# Example values with comments
# Default values when applicable
```

**Docker Compose Pattern:**
```yaml
services:
  service_name:
    environment:
      - ENV_VAR_NAME=value
      - ANOTHER_VAR=${FROM_ENV}
```

**Test Data Format:**
- JSON format matching API response structure
- Pretty-printed for readability
- Include complete response with all fields

---

### 5. Testing Strategy

**Unit Tests (for configuration parsing):**
- Verify `.env.example` variables can be parsed correctly
- Test timezone default values
- Validate environment variable precedence

**Integration Tests:**
- Use `test/compare_responses.sh` to compare captured test data against live API
- Verify project 450 responses match expected artifacts
- Check timezone-aware date calculations

**Manual Verification Steps:**
1. Start services with new configuration
2. Verify no startup errors related to TZ_LOCATION
3. Test endpoint `/v1/projects/450.json` returns consistent data
4. Compare captured test data against current API response

---

### 6. Risks and Considerations

**Risk 1: Timezone Calculation Discrepancies**
- **Mitigation:** Document that `TZ_LOCATION` must match Rails application timezone
- **Impact:** If mismatched, `days_unreading` may differ by hours/days
- **Solution:** Use same timezone identifier in both Go and Rails configs

**Risk 2: Test Data Becomes Stale**
- **Mitigation:** Include timestamp in test data file comments
- **Impact:** Regression tests may fail if API changes
- **Solution:** Document when test data was captured; update periodically

**Risk 3: Docker Compose Environment Variable Conflicts**
- **Mitigation:** Use consistent variable naming across all services
- **Impact:** Services may use different timezones causing inconsistent calculations
- **Solution:** Define TZ_LOCATION once in docker-compose.yml and reference via `${TZ_LOCATION}`

**Design Decision:** This task does NOT implement timezone support in the Go code - it only prepares the configuration infrastructure. The actual implementation will happen in Phase 4 after this documentation is in place.

---

### Implementation Checklist

#### Phase 1: Documentation Updates (30 minutes)
- [ ] Update `.env.example` with `TZ_LOCATION` configuration
- [ ] Add explanatory comments for timezone impact
- [ ] Include multiple timezone examples (America/Sao_Paulo, Europe/London, Asia/Tokyo)

#### Phase 2: Docker Compose Updates (15 minutes)
- [ ] Add timezone environment variable to PostgreSQL service
- [ ] Add timezone environment variable to Go API service
- [ ] Add timezone environment variable to Rails API service
- [ ] Verify all services reference same timezone config

#### Phase 3: Test Data Creation (30 minutes)
- [ ] Create `test/data/` directory
- [ ] Capture project 450 response from Go API
- [ ] Capture project 450 response from Rails API
- [ ] Capture logs for project 450 from both APIs
- [ ] Format JSON with `jq` for readability
- [ ] Add metadata comments (timestamp, API version)

#### Phase 4: Verification (30 minutes)
- [ ] Run `make docker-up` to verify configuration loads
- [ ] Test `/v1/projects/450.json` endpoint
- [ ] Compare captured data with live response
- [ ] Document any discrepancies

---

### Expected Outcomes

1. **`.env.example`** now documents timezone configuration with clear examples
2. **`docker-compose.yml`** ensures consistent timezone across all containers
3. **Test data artifacts** provide baseline for regression testing
4. **Documentation** enables future Phase 4 implementation

---

### Approval Requirements

Before proceeding with implementation:
- [ ] Review PRD Section: Configuration Files requirements
- [ ] Confirm test data structure matches `test/compare_responses.sh` expectations
- [ ] Verify timezone approach aligns with Rails `Date.today` behavior
- [ ] Ensure no breaking changes to existing configurations
<!-- SECTION:PLAN:END -->

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
- [ ] #13 Configuration validated with docker-compose
- [ ] #14 Test data matches actual API responses
<!-- DOD:END -->
