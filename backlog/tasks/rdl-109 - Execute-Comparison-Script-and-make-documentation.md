---
id: RDL-109
title: Execute Comparison Script and make documentation
status: To Do
assignee:
  - workflow
created_date: '2026-04-27 23:50'
updated_date: '2026-04-27 23:56'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Execute compare to the route v1/dashboard/day.json and make a documentation of differencies

test/compare_responses.sh
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task involves executing the JSON response comparison script specifically for the `/v1/dashboard/day.json` endpoint and documenting any differences found between the Go API and Rails API implementations.

**Technical Approach:**
1. **Preparation Phase**: Ensure both Go and Rails APIs are running with synchronized test data
2. **Execution Phase**: Run the comparison script targeting only the dashboard/day endpoint
3. **Analysis Phase**: Review comparison results and identify specific differences
4. **Documentation Phase**: Create a detailed comparison report following the established format

**Architecture Decisions:**
- Use the existing `test/compare_responses.sh` script with the `--dashboard-only` flag to isolate dashboard endpoint testing
- Follow the documentation format established in `docs/endpoint-comparison-report-v1-projects.md` for consistency
- Document both structural differences (JSON format) and value differences (calculated fields)

**Why This Approach:**
- The comparison script already has the `test_dashboard_day()` function implemented
- Using the existing infrastructure ensures consistency with other comparison tests
- The established documentation format provides a clear structure for reporting findings

---

### 2. Files to Modify

**Files to Read/Analyze:**
- `test/compare_responses.sh` - Existing comparison script (lines 816-850 for `test_dashboard_day` function)
- `internal/api/v1/handlers/dashboard_handler.go` - Go implementation of the Day endpoint (lines 43-115)
- `rails-app/app/controllers/v1/dashboard/days_controller.rb` - Rails implementation (reference)
- `docs/endpoint-comparison-report-v1-projects.md` - Documentation format reference

**Files to Create:**
- `docs/endpoint-comparison-report-dashboard-day.md` - New comparison report documenting differences for `/v1/dashboard/day.json`

**Files to Modify:**
- None (this is a documentation-only task)

---

### 3. Dependencies

**Prerequisites:**
1. **Both APIs must be running:**
   - Go API on port 3000: `make run`
   - Rails API on port 3001: `cd rails-app && rails server -p 3001`

2. **Database must have test data:**
   - Run `make reload-db` to populate test fixtures
   - Ensure both APIs connect to the same database

3. **Required tools:**
   - `curl` - For HTTP requests
   - `jq` (version 1.6+) - For JSON comparison

**Blocking Issues:**
- If Rails API is not available, the comparison cannot be executed
- If test database is empty, the comparison will not have meaningful data

**Setup Steps:**
```bash
# 1. Start Go API
make run

# 2. Start Rails API (in separate terminal)
cd rails-app && rails server -p 3001

# 3. Ensure database has data
make reload-db

# 4. Verify both APIs are accessible
curl http://localhost:3000/healthz
curl http://localhost:3001/healthz
```

---

### 4. Code Patterns

**Documentation Format:**
Follow the structure from `docs/endpoint-comparison-report-v1-projects.md`:

```markdown
# Endpoint Comparison Report: [Endpoint Name]

## Overview
- Report metadata (date, task, comparison target)

## Executive Summary
- Table with metrics (Endpoints Tested, Tests Passed/Failed, Critical Issues)
- Status indicator (✅ PASS / ❌ REQUIRES FIXES)

## Endpoints Tested
- List of endpoints with methods

## Detailed Findings
- Each issue with:
  - Issue number and name
  - Severity level (HIGH/MEDIUM/LOW)
  - Description
  - Go Response example
  - Rails Response example
  - Root cause analysis
  - Impact assessment
  - Recommended action

## Non-Critical Differences
- Acceptable structural choices documented as API contract variations

## Test Results Summary
- Detailed test output from comparison script

## Recommendations
- Priority-ranked list of fixes needed

## Appendix
- Raw comparison script output
```

**Severity Levels:**
- **HIGH**: Critical functional differences that break client compatibility
- **MEDIUM**: Structural differences that require client adaptation
- **LOW**: Minor formatting differences with minimal impact

**Field Comparison Categories:**
1. **Structure**: JSON envelope format, field presence
2. **Data Types**: String vs integer IDs, number precision
3. **Date Formats**: RFC3339 vs ISO date-only
4. **Calculated Fields**: `mean_day`, `per_pages`, `spec_mean_day`, `progress_geral`
5. **Error Handling**: Error response format and status codes

---

### 5. Testing Strategy

**Execution Steps:**

1. **Run the comparison script:**
```bash
./test/compare_responses.sh --dashboard-only \
  -g http://localhost:3000 \
  -r http://localhost:3001
```

2. **Capture output:**
```bash
./test/compare_responses.sh --dashboard-only \
  -g http://localhost:3000 \
  -r http://localhost:3001 \
  > /tmp/dashboard_day_comparison.log 2>&1
```

3. **Manual verification (if needed):**
```bash
# Fetch Go API response
curl http://localhost:3000/v1/dashboard/day.json | jq . > go_response.json

# Fetch Rails API response
curl http://localhost:3001/v1/dashboard/day.json | jq . > rails_response.json

# Compare structures
diff <(jq -S 'keys' go_response.json) <(jq -S 'keys' rails_response.json)

# Compare specific fields
jq '.data.attributes.stats' go_response.json
jq '.data.attributes.stats' rails_response.json
```

**Edge Cases to Document:**
- Empty database scenario (no logs)
- Date parameter handling (`?date=...`)
- Invalid date format errors
- Null/missing field handling
- Floating point precision differences

**Expected Comparison Metrics:**
- Structural match: Same keys and nesting
- Value tolerance rules:
  - `mean_day`: ±0.01 (floating point rounding)
  - `per_pages`: ±0.01 (floating point rounding)
  - `spec_mean_day`: ±0.01 (derived from mean_day)
  - `progress_geral`: ±0.01 (percentage calculation)

**Verification Checklist:**
- [ ] Script executes without errors
- [ ] Both APIs return valid JSON
- [ ] Structure comparison completed
- [ ] Value comparison completed
- [ ] All differences documented
- [ ] Root causes identified
- [ ] Recommendations provided

---

### 6. Risks and Considerations

**Known Risks:**

1. **Database Synchronization:**
   - **Risk**: Go and Rails APIs may connect to different databases
   - **Mitigation**: Verify connection strings in `.env` files before running comparison
   - **Detection**: Check first project ID in responses - should match

2. **Date/Time Calculation Differences:**
   - **Risk**: Timezone differences may cause 1-day variance in date-based calculations
   - **Mitigation**: Document acceptable tolerance (±1 day)
   - **Reference**: See `json-response-comparison.md` tolerance rules

3. **Rails API Availability:**
   - **Risk**: Rails application may not be running or accessible
   - **Mitigation**: Check Rails server status before running script
   - **Fallback**: Document as "Rails API unavailable - comparison skipped"

4. **Test Data Insufficiency:**
   - **Risk**: Empty or minimal test data may not trigger all code paths
   - **Mitigation**: Use `make reload-db` to populate comprehensive fixtures
   - **Verification**: Check that responses contain actual data, not empty arrays

**Acceptable Differences:**
- Floating point rounding (±0.01)
- Date format variations (RFC3339 vs date-only)
- JSON key ordering (normalized by jq)
- Null vs omitted fields (if semantically equivalent)

**Blocking Issues:**
- If comparison script fails to run, document the error and root cause
- If Rails API returns 500 errors, investigate Rails logs
- If Go API returns 500 errors, check `server.log` for stack traces

**Documentation Considerations:**
- Include raw script output in appendix for reproducibility
- Screenshot or copy-paste exact JSON responses for each difference
- Note the exact date/time when comparison was executed (affects date-based calculations)
- Record database state (number of projects, logs) for context

**Follow-up Tasks:**
- Any HIGH severity differences should be logged as separate tasks
- MEDIUM severity differences may be addressed in future alignment phases
- LOW severity differences can be documented as known API variations

---

**Implementation Timeline:**
1. **Preparation** (15 min): Start APIs, verify connectivity, reload database
2. **Execution** (10 min): Run comparison script, capture output
3. **Analysis** (30 min): Review results, identify differences, analyze root causes
4. **Documentation** (45 min): Create comparison report following established format
5. **Review** (15 min): Verify report completeness and accuracy

**Total Estimated Time:** ~2 hours
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
