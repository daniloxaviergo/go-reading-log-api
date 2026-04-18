---
id: RDL-059
title: Compare rails vs go projects
status: To Do
assignee: []
created_date: '2026-04-18 00:24'
updated_date: '2026-04-18 09:32'
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
