---
id: RDL-034
title: '[doc-002 Phase 5] Execute JSON response comparison test'
status: To Do
assignee:
  - catarina
created_date: '2026-04-03 14:04'
updated_date: '2026-04-04 04:05'
labels:
  - phase-5
  - json-comparison
  - testing
dependencies: []
references:
  - 'PRD Section: Files Created - compare_responses.sh'
  - 'PRD Section: Acceptance Criteria - AC1'
  - AC2
  - AC3
documentation:
  - doc-002
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create and execute test script comparing Go and Rails API JSON responses for all three endpoints (/v1/projects.json, /v1/projects/{id}.json, /v1/projects/{id}/logs.json). Use curl and jq to extract and compare JSON structures and verify identical field names, values, and formats.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Script automated response comparison for all 3 endpoints
- [ ] #2 JSON structures identical between Go and Rails
- [ ] #3 All field values match within tolerance
- [ ] #4 Edge cases tested (empty logs, null values)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Create a bash script using curl and jq to compare JSON responses between Go and Rails API endpoints. The script will:

1. **Endpoint Testing**: Call all three endpoints:
   - `/api/v1/projects.json`
   - `/api/v1/projects/{id}.json`
   - `/api/v1/projects/{id}/logs.json`

2. **Comparison Strategy**:
   - Use curl to fetch JSON from both APIs
   - Use jq to normalize and compare structures
   - Compare field names, data types, and values
   - Handle date/time format differences (RFC3339 vs Rails formats)

3. **Test Scenarios**:
   - Normal data (project with logs)
   - Empty logs (project without logs)
   - Null values (projects with NULL dates)
   - Multiple projects (order verification)

4. **Output**: Generate human-readable comparison report with pass/fail status and detailed diffs

### 2. Files to Modify

| File | Action | Reason |
|------|--------|--------|
| `test/compare_responses.sh` | Create | Main comparison script using curl + jq |
| `test/compare_responses_test.sh` | Create | Test scenarios and fixtures (optional) |
| `docs/api-comparison-guide.md` | Create | Documentation for manual comparison verification |

### 3. Dependencies

**Prerequisites**:
- Both Rails and Go APIs must be running on configured ports
- Test database must have test data (projects with logs)
- jq version ≥ 1.6 (for slurp and deep comparison)
- curl available in PATH

**Current State**:
- Rails API: Already implements `project_serializer.rb` with derived fields (progress, status, median_day, finished_at, logs_count, days_unreading)
- Go API: Already implements all derived calculations in `internal/domain/models/project.go`:
  - `CalculateProgress()`: (page/total_page)*100 rounded to 2 decimals
  - `CalculateStatus()`: unstarted/finished/running/sleeping/stopped
  - `CalculateDaysUnreading()`: days since last log or started_at
  - `CalculateMedianDay()`: page/days_reading rounded
  - `CalculateFinishedAt()`: estimated finish date
  - `CalculateLogsCount()`: logs.size

**Required Setup**:
```bash
# Start Rails API (if not already running)
cd rails-app && rails server -p 3001

# Start Go API (if not already running)
make run
```

### 4. Code Patterns

**Script Conventions**:
- Use bash best practices (set -euo pipefail)
- Define API URLs as variables (configurable via env vars)
- Use jq for all JSON parsing (not grep/sed)
- Generate normalized JSON for comparison (key order independent)
- Handle edge cases: null values, empty arrays, missing fields

**JSON Normalization**:
```bash
# Sort keys to ensure consistent comparison
jq -S '.' input.json > normalized.json

# Compare two JSON files
diff <(jq -S '.' file1.json) <(jq -S '.' file2.json)
```

**Field Mapping**:
| Rails Field | Go Field | Notes |
|-------------|----------|-------|
| started_at | started_at | May differ in format (date vs datetime) |
| progress | progress | Both calculate same formula |
| status | status | Both use same 5 values |
| logs_count | logs_count | Count of logs array |
| days_unreading | days_unreading | Days since last log |
| median_day | median_day | page/days_reading |
| finished_at | finished_at | Estimated finish date |
| logs | logs | Limited to 4, ordered by data DESC |

**Comparison Tolerance**:
- Dates: Consider same if within 1 second (time zone handling)
- Numbers: Consider same if within 0.01 tolerance (rounding)
- Strings: Exact match required
- Arrays: Same length and element order

### 5. Testing Strategy

**Automated Test Script** (`test/compare_responses.sh`):

1. **Setup Phase**:
   - Validate API endpoints are accessible
   - Fetch test data from Rails API (reference)
   - Fetch corresponding data from Go API (test)

2. **Comparison Tests**:
   ```bash
   test_projects_index() {
     # Compare all projects
     # Verify count, IDs, derived fields
   }
   
   test_project_show() {
     # Compare single project with ID
     # Verify all fields including median_day, finished_at
   }
   
   test_project_logs() {
     # Compare logs for project
     # Verify log count, ordering (data DESC)
   }
   ```

3. **Edge Case Tests**:
   ```bash
   test_empty_logs() {
     # Project with no logs
     # Verify logs_count = 0, logs = []
   }
   
   test_null_dates() {
     # Project with NULL started_at
     # Verify handling of null dates
   }
   ```

4. **Test Output**:
   - Pass/fail summary for each endpoint
   - Detailed diff for failing fields
   - Exit code: 0 (all pass), 1 (any fail)

**Integration Test** (`test/integration/json_comparison_test.go`):
- Write Go test that validates JSON responses using test framework
- Compare against expected values from testdata
- Verify derived field calculations match expected values

### 6. Risks and Considerations

**Known Challenges**:

1. **Date/Time Format Differences**:
   - Rails may use different date formats (ISO 8601 vs custom)
   - Go uses RFC3339 consistently
   - **Mitigation**: Parse both dates and compare as timestamps

2. **Floating Point Precision**:
   - Division results may differ slightly due to rounding
   - **Mitigation**: Compare with tolerance (0.01) for fields like progress, median_day

3. **Null Value Handling**:
   - Rails may omit null fields while Go includes them
   - **Mitigation**: Normalize JSON to include all fields

4. **Array Ordering**:
   - Logs must be in same order (data DESC)
   - **Mitigation**: Compare arrays element-by-element with explicit order assertion

5. **Configuration Differences**:
   - Rails uses em_andamento_range=8, dormindo_range=16
   - Go uses em_andamento_range=7, dormindo_range=14 (different from PRD spec!)
   - **Mitigation**: Document config differences and test with same config

6. **Performance**:
   - Multiple API calls may timeout
   - **Mitigation**: Use reasonable timeouts, parallelize carefully

**Testing Environment Requirements**:
```bash
# Environment variables for script
RAILS_API_URL=http://localhost:3001/api/v1
GO_API_URL=http://localhost:3000/api/v1
```

**Acceptance Criteria Alignment**:
- ✅ AC1: Script automated for all 3 endpoints - Script will have separate test functions
- ✅ AC2: JSON structures identical - jq -S comparison ensures this
- ✅ AC3: Field values match within tolerance - Compare numbers with 0.01 tolerance
- ✅ AC4: Edge cases tested - Separate test functions for empty logs, null dates

**Estimated Duration**: 4-8 hours
- Script development: 2-4 hours
- Testing and debugging: 2-4 hours
- Documentation: 1 hour
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
<!-- DOD:END -->
