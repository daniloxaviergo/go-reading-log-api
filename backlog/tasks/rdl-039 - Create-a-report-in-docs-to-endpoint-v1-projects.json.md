---
id: RDL-039
title: Create a report in docs to endpoint v1/projects.json
status: To Do
assignee:
  - workflow
created_date: '2026-04-12 20:40'
updated_date: '2026-04-12 20:45'
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

This task requires executing the existing JSON response comparison script and documenting any differences found between the Go and Rails API implementations.

**Approach:**
- Execute the `test/compare_responses.sh` script which compares JSON responses from both APIs
- Analyze any structural or value differences found
- Create a comprehensive report in the `docs/` directory documenting:
  - All differences found between the two implementations
  - Root cause analysis for each difference
  - Whether differences are acceptable or need fixing

**Architecture decision:** The report should be a Markdown document that captures:
- Comparison methodology used
- All observed differences organized by category
- Impact assessment of each difference
- Recommended actions

This approach ensures we don't just report that differences exist, but provide actionable insights into what needs to be addressed.

### 2. Files to Modify

**Files to Create:**
- `docs/endpoint-comparison-report-v1-projects.md` - Main comparison report

**Files to Reference (read-only):**
- `test/compare_responses.sh` - The comparison script to execute
- `internal/api/v1/handlers/projects.go` - Go API handler for `/api/v1/projects`
- `app/controllers/api/v1/projects_controller.rb` - Rails API controller for `/api/v1/projects`

**Files to Potentially Modify (if differences require fixes):**
- `internal/api/v1/handlers/projects.go` - Go handler corrections
- `internal/domain/dto/` - DTO structure corrections

### 3. Dependencies

**Prerequisites:**
- Both Go and Rails APIs must be running simultaneously
- Go API running on `http://localhost:3000`
- Rails API running on `http://localhost:3001`
- Required tools: `curl`, `jq` (version 1.6+)
- Database must contain test data (at least one project with logs)

**Setup Steps:**
```bash
# Start both APIs (if not already running)
make docker-up  # Or manual startup

# Verify both APIs are accessible
curl http://localhost:3000/api/v1/projects
curl http://localhost:3001/api/v1/projects
```

### 4. Code Patterns

**Comparison Script Analysis:**
The script tests three endpoints:
1. `GET /api/v1/projects` - Index endpoint (collection)
2. `GET /api/v1/projects/:id` - Show endpoint (single resource)
3. `GET /api/v1/projects/:id/logs` - Logs endpoint (nested resource)

**Key Comparison Points:**
- JSON structure and key names (snake_case vs camelCase)
- Null handling (null vs omitted field)
- Number formatting (int vs float)
- Date/time formats (RFC3339 compliance)
- Calculated fields consistency (progress, status, logs_count, etc.)

**Reporting Convention:**
Use this structure for each difference found:
```markdown
### Difference: [Field/Structure Name]
- **Location:** [Endpoint + context]
- **Go Value:** [actual value]
- **Rails Value:** [actual value]
- **Impact:** [Low/Medium/High]
- **Recommended Action:** [Fix/Document as intentional/Investigate]
```

### 5. Testing Strategy

**Script Execution:**
```bash
cd /home/danilo/scripts/github/go-reading-log-api-next
./test/compare_responses.sh
```

**Manual Verification Steps:**
1. Run the comparison script and capture all output
2. Manually inspect each reported difference
3. Verify if differences are due to:
   - Implementation bugs
   - Intentional design decisions
   - Version/encoding differences

**Report Quality Check:**
- All differences must be documented, not just "fixed" ones
- Include sample JSON payloads for significant differences
- Provide clear guidance on whether each difference requires action

### 6. Risks and Considerations

**Known Risks:**
- The comparison script may fail if APIs are not properly running
- Floating point comparison tolerance (0.01) may mask real differences
- Date/time differences due to timezone handling
- Ordering of array elements may cause false positives

**Design Considerations:**
- Some differences may be acceptable (e.g., Go using float64 vs Rails integer)
- Field ordering in JSON objects may differ (standard is to sort keys)
- Null vs omitted field is a common source of differences
- The report should distinguish between "blocking" and "non-blocking" differences

**Blocking Issues:**
- None identified - this is a documentation/reporting task
- However, the report may uncover blocking issues that require follow-up tasks

**Output Location:**
The final report should be saved as `docs/endpoint-comparison-report-v1-projects.md` to match the project's documentation structure.
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
