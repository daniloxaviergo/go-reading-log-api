---
id: RDL-039
title: Create a report in docs to endpoint v1/projects.json
status: To Do
assignee:
  - thomas
created_date: '2026-04-12 20:40'
updated_date: '2026-04-13 10:42'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
execute test/compare_responses.sh to endpoint v1/projects.json and make a report of differencies
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task requires updating the comparison script to use the correct `.json` endpoint suffixes and running a comprehensive comparison between the Go API and Rails API for the projects endpoint.

**Key Technical Decisions:**
- Update `test/compare_responses.sh` to use `.json` suffix for all endpoints
- The script tests three endpoints: index, show, and logs
- Compare both JSON structure and values (with 0.01 tolerance for floating-point numbers)
- Document all discrepancies in a structured markdown report
- The existing report `docs/endpoint-comparison-report-v1-projects.md` shows there are known differences that need to be re-verified

**Why This Approach:**
- The comparison script already exists and handles edge cases (null values, empty arrays)
- It provides consistent, repeatable comparisons
- The existing report format from `endpoint-comparison-report-v1-projects.md` provides a good template
- Using established tools (`jq`, `curl`) ensures portability

---

### 2. Files to Modify

#### Files to Modify:
1. **`test/compare_responses.sh`** - Update endpoint URLs to use `.json` suffix
   - Line 21-22: Update default URLs
   - Lines 200-220: Update `test_index_endpoint()` to use `/api/v1/projects.json`
   - Lines 223-240: Update `test_show_endpoint()` to use `/api/v1/projects/:id.json`
   - Lines 243-260: Update `test_logs_endpoint()` to use `/api/v1/projects/:id/logs.json`

#### New Files to Create:
1. **`docs/rdl-039-comparison-report.md`** - The main comparison report documenting all findings

#### Files to Read (No modification):
1. **`docs/endpoint-comparison-report-v1-projects.md`** - Template/reference for report format
2. **`docs/comprare.md`** - Additional context on comparison plan

#### Existing Documentation to Review:
1. **`AGENTS.md`** - Project context and API documentation
2. **`docs/README.go-project.md`** - Detailed project structure

---

### 3. Dependencies

#### Prerequisites:
- ✅ **Go API** must be running on `http://localhost:3000`
- ✅ **Rails API** must be running on `http://localhost:3001`
- ✅ **PostgreSQL** database must be accessible
- ✅ **Required tools**: `curl`, `jq` (version 1.6+), `bash`
- ✅ **Test database** must contain at least one project with logs

#### Setup Steps:
1. Start both APIs (via `make docker-up` or direct execution)
2. Verify database connectivity
3. Ensure both APIs are accessible
4. Make comparison script executable: `chmod +x test/compare_responses.sh`

#### Blocking Issues:
- None identified - this is a documentation task that verifies existing functionality

---

### 4. Code Patterns

The implementation follows existing patterns in the project:

**Comparison Script Patterns:**
- Uses `jq -S` for JSON normalization (sorted keys)
- Implements value comparison with floating-point tolerance
- Handles null/empty values gracefully
- Provides color-coded console output for clarity

**Report Format Patterns:**
- Markdown structure with headers and tables
- Severity classification (CRITICAL, HIGH, MEDIUM, LOW)
- Code blocks for JSON examples
- Summary tables for quick overview
- Actionable recommendations

**Integration with Existing Codebase:**
- Report stored in `docs/` directory alongside other comparison reports
- Follows naming convention: `rdl-039-comparison-report.md`
- References the task ID in the filename and document metadata

---

### 5. Testing Strategy

The comparison script includes comprehensive test coverage:

#### Test Cases to Execute:
1. **Index Endpoint Test** (`test_index_endpoint`)
   - Fetches `/api/v1/projects.json` from both APIs
   - Compares JSON structure using normalized diff
   - Compares values with 0.01 tolerance
   - Verifies first project data consistency

2. **Show Endpoint Test** (`test_show_endpoint`)
   - Fetches `/api/v1/projects/:id.json` from both APIs
   - Compares complete project details
   - Validates derived fields (progress, status, median_day, etc.)

3. **Logs Endpoint Test** (`test_logs_endpoint`)
   - Fetches `/api/v1/projects/:id/logs.json` from both APIs
   - Compares log entry structure
   - Validates nested objects and arrays

4. **Edge Cases Test** (`test_edge_cases`)
   - Empty logs scenario
   - Null date handling
   - Field presence/absence consistency

#### Test Verification:
- Count passed/failed tests
- Capture detailed diff output for failures
- Verify all required fields are present
- Confirm data type consistency

---

### 6. Risks and Considerations

#### Potential Issues:

| Risk | Impact | Mitigation |
|------|--------|------------|
| APIs not running | Test cannot execute | Document setup requirements clearly |
| Different databases | All data differs | Check database configuration first |
| Network failures | Incomplete comparison | Add retry logic or timeout handling |
| Floating-point precision | False positives | Use 0.01 tolerance as implemented |
| Date format differences | Structure mismatch | Normalize dates before comparison |

#### Design Trade-offs:

1. **JSON Structure Comparison**: The script normalizes JSON (sorts keys) before comparing structure, which ignores key ordering differences but may mask intentional field ordering in Rails API.

2. **Tolerance Value**: The 0.01 tolerance for floating-point numbers may mask legitimate calculation differences. Consider making this configurable.

3. **First Project Selection**: Using the first project from the index response assumes consistent ordering. Consider adding explicit project ID selection for reproducibility.

4. **No PATCH/POST Tests**: This is a read-only comparison task. Write operations are out of scope for RDL-039.

#### Recommendations:
- Run the comparison script periodically as part of CI/CD
- Fail the build if critical differences are detected
- Update the report after any API changes that affect response structure
- Consider adding automated alerting for regression detection

---

### 7. Implementation Steps

1. **Pre-flight Check**
   - Verify APIs are running
   - Verify database connectivity
   - Ensure required tools are installed

2. **Execute Comparison**
   - Run `./test/compare_responses.sh`
   - Capture output for documentation

3. **Analyze Results**
   - Review pass/fail status
   - Identify specific field differences
   - Categorize by severity

4. **Create Report**
   - Document all findings in `docs/rdl-039-comparison-report.md`
   - Include before/after examples
   - Provide actionable recommendations

5. **Review and Approve**
   - Verify report accuracy
   - Ensure all acceptance criteria are met
   - Update task status in backlog
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Task Execution Progress: RDL-039 - Create a report in docs to endpoint v1/projects.json

### Date: 2026-04-13

### What I've Done

1. **Analyzed the task requirements:**
   - Execute the comparison script `test/compare_responses.sh` against the v1/projects.json endpoint
   - Create a comprehensive report documenting differences between Go API and Rails API
   - Use the testing-expert subagent for test execution

2. **Examined the existing infrastructure:**
   - Found the comparison script at `test/compare_responses.sh`
   - Found an existing comparison report at `docs/endpoint-comparison-report-v1-projects.md`
   - The script tests three endpoints: index, show, and logs
   - Uses jq for JSON comparison with 0.01 tolerance for floating-point numbers

3. **Identified the issue:**
   - The comparison script was designed to use `.json` suffixes
   - However, neither the Go API nor the Rails API actually have `.json` suffix routes
   - The Go API routes are `/api/v1/projects` and `/api/v1/projects/:id` (no .json suffix)
   - The Rails API routes appear to follow JSON:API conventions differently

4. **Fixed the comparison script:**
   - Added directory creation in `fetch_json` function to ensure temp directory structure exists
   - Updated endpoint URLs to remove `.json` suffix (both APIs use plain routes)
   - Added proper error handling for API accessibility checks

5. **Ran the comparison script:**
   - Both APIs are running (Go on 3000, Rails on 3001)
   - The script successfully executed and captured differences
   - Key findings documented in the output

### Pending Actions

- [ ] Create comprehensive report documenting all differences found
- [ ] Verify all acceptance criteria are met
- [ ] Submit PR-style final summary

### Next Steps

1. Document all differences found in the comparison between Go API and Rails API
2. Create the report at `docs/rdl-039-comparison-report.md`
3. Address any critical differences that need fixing
4. Finalize the task in the backlog

### Notes

The comparison script now works correctly with both APIs running. The next step is to generate the comprehensive report documenting all discrepancies found between the two implementations.
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass use testing-expert subagent for test execution and verification
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
