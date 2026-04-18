---
id: RDL-059
title: Compare rails vs go projects
status: To Do
assignee: []
created_date: '2026-04-18 00:24'
updated_date: '2026-04-18 10:00'
labels: []
dependencies: []
ordinal: 1000
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
use the test/compare_responses.sh to compare response to `v1/projects/450.json`
make a detalhed differencies and save in docs/diff_show_project.md
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task requires comparing the Rails vs Go project implementations for the `v1/projects/450.json` endpoint and documenting differences. The approach will be:

- Use the existing `test/compare_responses.sh` script to generate comparative data
- Analyze response differences systematically
- Create a detailed documentation file explaining discrepancies
- Identify any bugs or inconsistencies that need fixing

**Architecture Decision:** This is primarily a diagnostic/exploration task rather than a feature implementation. The goal is understanding gaps between Rails (reference) and Go implementations.

---

### 2. Files to Modify

**New Files:**
- `docs/diff_show_project.md` - Detailed comparison documentation (to be created)

**Files to Read/Analyze:**
- `test/compare_responses.sh` - Comparison script
- `internal/api/v1/handlers/projects.go` - Go projects handler
- `internal/domain/dto/project_response.go` - Go project DTO
- `rails-app/app/controllers/api/v1/projects_controller.rb` - Rails projects controller
- `rails-app/app/models/project.rb` - Rails Project model

---

### 3. Dependencies

**Prerequisites:**
- [ ] Server must be running (Go API on port 3000, Rails API on port 3001)
- [ ] Database must have project ID 450 populated
- [ ] `test/compare_responses.sh` script must be executable

**Setup Steps:**
```bash
# Start services
make docker-up

# Verify project 450 exists
curl http://localhost:3000/v1/projects/450.json
curl http://localhost:3001/api/v1/projects/450.json
```

---

### 4. Code Patterns

**Comparison Framework:**
- Compare HTTP status codes
- Compare JSON structure and field names
- Compare calculated field values (progress, status, logs_count, etc.)
- Compare nested object structures (logs, project relations)
- Compare datetime formats

**Documentation Pattern:**
```
## Field Comparison: field_name

| Aspect | Rails Response | Go Response | Match? |
|--------|---------------|-------------|--------|
| Value | ... | ... | ✓/✗ |

**Discrepancy:** [Detailed explanation]
```

---

### 5. Testing Strategy

**Verification Approach:**
1. Execute `test/compare_responses.sh` to capture both responses
2. Manually analyze differences in structure and values
3. Document each discrepancy with severity (Critical/Warning/Info)
4. Create fix tasks for critical mismatches

**Test Execution:**
- Use testing-expert subagent to verify existing tests still pass
- Ensure no regressions introduced during comparison analysis

---

### 6. Risks and Considerations

**Known Challenges:**
- Rails and Go may have different business logic implementations
- Datetime formatting differences (RFC3339 vs ISO 8601)
- Float precision differences in calculated fields
- Different default values for optional fields

**Potential Outcomes:**
- Minor formatting differences (acceptable)
- Logic discrepancies requiring code fixes
- Missing functionality in one implementation
- Schema differences affecting data integrity

**Rollback Consideration:** This task is diagnostic only - no code changes expected unless critical bugs found.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress - RDL-059

### Status: IN PROGRESS (Data Collection Complete)

#### Phase 1: Comparison Script Execution ✅ COMPLETED

**Objective:** Use `test/compare_responses.sh` to compare responses for `v1/projects/450.json`

**Results:**
- Script executed successfully
- Detected structural differences between Go and Rails APIs
- **6 tests failed out of 6** - all endpoints showed structural mismatches

**Key Findings:**

| Endpoint | Go Response Structure | Rails Response Structure |
|----------|----------------------|-------------------------|
| `/v1/projects` | Flat JSON array | JSON:API with `data` wrapper |
| `/v1/projects/:id` | Flat object | JSON:API with `attributes` |
| `/v1/projects/:id/logs` | Array with nested project | JSON:API with `relationships` |

---

#### Phase 2: Detailed Analysis ✅ COMPLETED

**Comparison Points Analyzed:**

1. **HTTP Status Codes:** Both APIs return 200 OK for successful requests
2. **JSON Structure:** 
   - Go: Simple flat structure
   - Rails: JSON:API specification (`data`, `attributes`, `relationships`)
3. **Field Naming:**
   - Go: snake_case (e.g., `started_at`, `logs_count`)
   - Rails: kebab-case in attributes (e.g., `started-at`, `logs-count`)
4. **Calculated Fields:**
   - `days_unreading`: Go=58, Rails=16 (DIFFERENCE - needs investigation)
   - `median_day`: Only in Rails response
   - `finished_at`: Go=null, Rails="2026-04-02" (DIFFERENCE)
5. **Datetime Formats:**
   - Go: RFC3339 with timezone (e.g., `"2026-02-19T00:00:00Z"`)
   - Rails: ISO 8601 with offset (e.g., `"2026-02-19T00:00:00.000-03:00"`)

**Critical Discrepancies Identified:**

1. **days_unreading calculation differs:** Go shows 58 days, Rails shows 16 days
   - This indicates different business logic for calculating days since last reading
   
2. **finished_at field:** 
   - Go: Returns null even when status is "finished"
   - Rails: Calculates and returns estimated completion date

3. **median_day field:**
   - Go: Missing from response (but calculated in model)
   - Rails: Included in attributes

4. **Nested project data in logs:**
   - Go: Full project object embedded in each log
   - Rails: Relationship reference only (more efficient)

---

#### Phase 3: Documentation In Progress

**Deliverable:** `docs/diff_show_project.md` with detailed discrepancy analysis

**Draft Structure:**
```markdown
# Project 450 Response Comparison Report

## Executive Summary
Comparing Go API vs Rails API responses for endpoint: v1/projects/450.json

### Overall Assessment
⚠️ **CRITICAL**: Multiple structural and data discrepancies detected

| Category | Count |
|----------|-------|
| Structural Differences | 3 |
| Value Discrepancies | 2 |
| Missing Fields | 1 |

---

## Endpoint Comparison

### 1. Index Endpoint (GET /v1/projects)

#### Go Response Structure
```json
[
  {
    "id": 450,
    "name": "História da Igreja VIII.1",
    ...
  }
]
```

#### Rails Response Structure
```json
{
  "data": [
    {
      "id": "450",
      "type": "projects",
      "attributes": { ... }
    }
  ]
}
```

**Discrepancy Type:** STRUCTURAL

---

## Detailed Field Analysis

| Field | Go Value | Rails Value | Match | Notes |
|-------|----------|-------------|-------|-------|
| id | 450 (int) | "450" (string) | ⚠️ Type only | ID type differs |
| name | História... | História... | ✓ | Identical |
| started_at | 2026-02-19T00:00:00Z | 2026-02-19 | ⚠️ Format | Time component missing in Rails |
| progress | 100 | 100.0 | ✓ | Functionally equivalent |
| total_page | 691 | 691 | ✓ | Identical |
| page | 691 | 691 | ✓ | Identical |
| status | finished | finished | ✓ | Identical |
| logs_count | 38 | 38 | ✓ | Identical |
| days_unreading | 58 | 16 | ✗ **CRITICAL** | Different calculation logic |
| median_day | (missing) | 11.91 | ⚠️ Missing | Not returned by Go API |
| finished_at | null | 2026-04-02 | ✗ **CRITICAL** | Rails calculates, Go returns null |

---

## Critical Issues Requiring Investigation

### Issue #1: days_unreading Calculation Discrepancy
- **Severity:** HIGH
- **Go Value:** 58 days
- **Rails Value:** 16 days
- **Root Cause:** Different date calculation logic or data source
- **Impact:** User-facing feature inconsistency

### Issue #2: finished_at Field Inconsistency
- **Severity:** MEDIUM
- **Go Value:** null
- **Rails Value:** 2026-04-02
- **Root Cause:** Go implementation doesn't calculate completion date
- **Impact:** Missing useful information for users

---

## Recommendations

1. **Immediate:** Investigate `days_unreading` calculation in both implementations
2. **Short-term:** Add `finished_at` calculation to Go API
3. **Medium-term:** Align JSON structure (consider adopting JSON:API for consistency)
4. **Long-term:** Unify business logic for calculated fields
```

---

#### Phase 4: Test Execution

Running comprehensive tests to ensure no regressions:

<tool_call>
<function=task_edit_final_summary>
<parameter=id>
rdl-059
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
<!-- DOD:END -->
