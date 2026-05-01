---
id: RDL-147
title: '[doc-12 Phase 5] Update API documentation with fault definition'
status: To Do
assignee:
  - thomas
created_date: '2026-05-01 15:08'
updated_date: '2026-05-01 19:07'
labels:
  - documentation
  - phase-5
dependencies: []
documentation:
  - doc-012
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update API documentation to clarify fault definition (days with zero pages read), add calculation examples showing how faults are counted, and document edge cases (NULL handling, single-day ranges, empty database). Ensure documentation matches Rails behavior.
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task involves updating the API documentation to clarify the fault definition and add calculation examples. The approach is:

1. **Add a new Dashboard API section** to the QWEN.md documentation file that includes:
   - Clear definition of what a "fault" is (day with zero pages read)
   - Dashboard endpoints documentation (`/v1/dashboard/echart/faults.json`)
   - Calculation examples showing how faults are counted
   - Edge cases documentation (NULL handling, single-day ranges, empty database)
   - SQL query examples for transparency

2. **Ensure Rails parity documentation** - Document that the Go implementation matches Rails behavior exactly

3. **Reference existing detailed documentation** - Link to `docs/faults-calculation-explanation.md` for deep technical details

**Why this approach:**
- QWEN.md is the primary project context file for AI assistants and developers
- README.md is more high-level and focuses on Phase 1 endpoints
- The detailed technical explanation already exists in `docs/faults-calculation-explanation.md`
- Following the existing documentation pattern in QWEN.md (similar to how Calculated Fields are documented)

### 2. Files to Modify

| File | Change Type | Purpose |
|------|-------------|---------|
| `QWEN.md` | Modify | Add Dashboard API endpoints section with fault calculation documentation |

**Specific sections to add to QWEN.md:**

1. **New section after "Logs Endpoints"** - "Dashboard Endpoints"
   - Health check endpoint (already exists)
   - Faults endpoint documentation
   - Weekday faults endpoint documentation
   - Fault definition and calculation logic

2. **New subsection under "Calculated Fields"** - "Fault Metrics"
   - Definition of fault (day with zero pages)
   - Calculation formula
   - Examples with visual tables
   - Edge cases

**Files to review (no changes needed):**
- `docs/faults-calculation-explanation.md` - Already comprehensive, will be referenced
- `README.md` - High-level overview, no changes needed for Phase 1 focus

### 3. Dependencies

**Prerequisites:**
- RDL-146 (Create faults calculation documentation) - ✅ **COMPLETED**
  - The detailed technical document `docs/faults-calculation-explanation.md` already exists
  - All SQL queries and algorithms are documented there

**Existing Implementation (no code changes needed):**
- `internal/adapter/postgres/dashboard_repository.go` - GetFaultsByDateRange and GetWeekdayFaults methods are implemented
- `internal/api/v1/handlers/dashboard_handler.go` - Handlers are implemented
- All tests are passing (RDL-139 through RDL-146 completed)

**No blocking issues** - This is purely a documentation task.

### 4. Code Patterns

**Documentation Style to Follow:**

1. **End Format** - Match existing endpoint documentation style in QWEN.md:
   ```markdown
   ### Endpoint Name

   | Property | Value |
   |----------|-------|
   | **Method** | GET |
   | **Path** | `/v1/endpoint.json` |
   | **Description** | Description |
   | **Authentication** | None |
   | **Response Code** | 200 OK |

   **Request:**
   ```bash
   curl http://localhost:3000/v1/endpoint.json
   ```

   **Response (200 OK):**
   ```json
   {
     "example": "response"
   }
   ```
   ```

2. **Calculated Fields Format** - Match existing table format:
   ```markdown
   | Field | Type | Description | Formula |
   |-------|------|-------------|---------|
   | `field_name` | type | Description | Formula |
   ```

3. **Edge Cases Format** - Use the pattern from `docs/faults-calculation-explanation.md`:
   ```markdown
   ### Edge Case: Description

   **Scenario:** Brief description

   **Expected Result:** What the API returns

   **Key Takeaway:** Important note
   ```

4. **Examples with Visual Tables** - Use the format from existing documentation:
   ```markdown
   | Day | Date | Logs | Pages Read | Daily Sum | Fault? |
   |-----|------|------|------------|-----------|--------|
   | 1   | Jan 1| Yes  | 25         | 25        | No     |
   | 2   | Jan 2| No   | 0          | NULL      | YES ⚠️ |
   ```

### 5. Testing Strategy

**This is a documentation task - no code tests required.**

**Validation Steps:**
1. **Manual Review** - Verify documentation accuracy against:
   - `docs/faults-calculation-explanation.md` (source of truth)
   - `internal/adapter/postgres/dashboard_repository.go` (actual implementation)
   - Rails application reference (if available)

2. **Consistency Check** - Ensure:
   - Terminology matches (fault = day with zero pages)
   - SQL queries in documentation match actual code
   - Examples are mathematically correct
   - Edge cases align with implementation

3. **Peer Review** - Have team member verify:
   - Documentation is clear and accurate
   - Examples are easy to understand
   - Rails parity is correctly described

### 6. Risks and Considerations

**Known Issues/Considerations:**

1. **Documentation Accuracy Risk**
   - **Mitigation:** Use `docs/faults-calculation-explanation.md` as the source of truth
   - **Mitigation:** Cross-reference with actual SQL queries in `dashboard_repository.go`

2. **Rails Parity Claims**
   - **Consideration:** Document that Go implementation matches Rails behavior
   - **Verification:** RDL-145 added Rails comparison test cases (completed)

3. **Edge Case Coverage**
   - **Must Document:**
     - NULL values in start_page/end_page (treated as 0)
     - Single-day date ranges (start_date == end_date)
     - Empty database (all days are faults)
     - Logs with zero pages (start_page == end_page)
     - Inclusive date range boundaries (BETWEEN $1 AND $2)

4. **Documentation Maintenance**
   - **Consideration:** If faults calculation logic changes, both docs need update
   - **Recommendation:** Add comment in QWEN.md to reference `docs/faults-calculation-explanation.md`

5. **Phase 1 Context**
   - **Note:** Dashboard endpoints are part of the completed Phase 1 work
   - **Clarification:** Document that faults endpoints are read-only (consistent with Phase 1)

**No Blocking Issues Identified**

**Implementation Timeline:**
- Estimated time: 30-45 minutes
- Steps:
  1. Review existing documentation (10 min)
  2. Draft Dashboard endpoints section (15 min)
  3. Add fault calculation examples (10 min)
  4. Review and refine (10 min)

---

## Implementation Checklist

### Phase 1: Content Planning
- [ ] Review `docs/faults-calculation-explanation.md` for accurate information
- [ ] Review existing QWEN.md structure and formatting
- [ ] Identify exact insertion points for new sections

### Phase 2: Documentation Updates
- [ ] Add "Dashboard Endpoints" section with faults endpoint documentation
- [ ] Add "Fault Metrics" subsection under "Calculated Fields"
- [ ] Include visual examples with tables
- [ ] Document all edge cases
- [ ] Add curl examples for testing

### Phase 3: Validation
- [ ] Verify SQL queries match implementation
- [ ] Verify calculation examples are mathematically correct
- [ ] Ensure Rails parity claims are accurate
- [ ] Check formatting consistency with existing docs

### Phase 4: Review
- [ ] Self-review for clarity and completeness
- [ ] Request peer review if needed
- [ ] Update task status to Done

---

*Implementation Plan Version: 1.0*
*Created: 2026-05-01*
*Author: Implementation Planner*
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
### Implementation Progress

**Status:** Completed

**Phase 1: Content Planning** ✅ Completed
- Reviewed `docs/faults-calculation-explanation.md` - comprehensive technical details available
- Reviewed existing QWEN.md structure - identified insertion points
- Identified exact sections to add:
  1. New "Dashboard Endpoints" section after "Logs Endpoints"
  2. New "Fault Metrics" subsection under "Calculated Fields"

**Phase 2: Documentation Updates** ✅ Completed
- ✅ Added "Dashboard Endpoints" section with:
  - Health check endpoint documentation
  - Daily Statistics endpoint (`/v1/dashboard/day.json`)
  - Faults endpoint (`/v1/dashboard/echart/faults.json`)
  - Weekday Faults endpoint (`/v1/dashboard/echart/faults_week_day.json`)
  - Speculate vs Actual endpoint (`/v1/dashboard/echart/speculate_actual.json`)
  - Mean Progress endpoint (`/v1/dashboard/echart/mean_progress.json`)
  - Yearly Total endpoint (`/v1/dashboard/echart/last_year_total.json`)
  - Dashboard Projects endpoint (`/v1/dashboard/projects.json`)
  - Projects With Logs endpoint (`/v1/dashboard/projects_with_logs.json`)
  - Last Days Trend endpoint (`/v1/dashboard/last_days.json`)
  
- ✅ Added "Fault Metrics" subsection with:
  - Clear fault definition (day with zero pages read)
  - Calculation formula and examples
  - Visual table showing 10-day period example
  - SQL query pattern with detailed breakdown
  - Weekday faults documentation
  - All 6 edge cases documented:
    1. Empty Database
    2. Single-Day Range
    3. NULL Values in Logs
    4. Logs with Zero Pages
    5. All Days Have Reading
    6. Date Range Boundaries
  - Rails parity documentation

**Phase 3: Validation** ✅ Completed
- ✅ All unit tests pass
- ✅ All integration tests pass
- ✅ `go fmt` passes with no errors
- ✅ `go vet` passes with no errors
- ✅ Documentation accuracy verified against `docs/faults-calculation-explanation.md`
- ✅ SQL queries in documentation match actual implementation

**Phase 4: Review** ✅ Completed
- ✅ Self-review completed - documentation is clear and accurate
- ✅ Examples are easy to understand
- ✅ Rails parity is correctly described

**Definition of Done Checklist:**
- [x] #1 All unit tests pass
- [x] #2 All integration tests pass execution and verification
- [x] #3 go fmt and go vet pass with no errors
- [x] #4 Clean Architecture layers properly followed (N/A - documentation task)
- [x] #5 Error responses consistent with existing patterns (N/A - documentation task)
- [x] #6 HTTP status codes correct for response type (N/A - documentation task)
- [x] #7 Documentation updated in QWEN.md
- [x] #8 New code paths include error path tests (N/A - documentation task)
- [x] #9 HTTP handlers test both success and error responses (N/A - documentation task)
- [x] #10 Integration tests verify actual database interactions (N/A - documentation task)

**Ready to mark task as Done**
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
