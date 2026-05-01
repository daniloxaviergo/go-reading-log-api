---
id: RDL-146
title: '[doc-12 Phase 5] Create faults calculation documentation'
status: To Do
assignee:
  - thomas
created_date: '2026-05-01 15:08'
updated_date: '2026-05-01 18:45'
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
Create docs/faults-calculation-explanation.md with detailed explanation of Rails fault calculation logic, SQL query breakdown with CTE structure, and examples with diagrams showing how zero-page days are counted. Include comparison with previous incorrect implementation.
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task involves creating comprehensive technical documentation for the faults calculation logic that was fixed in Phase 4 (RDL-139, RDL-140). The documentation will serve as a reference for developers to understand:

- **What** a "fault" is in the context of the reading log API
- **How** the calculation works at the SQL level
- **Why** the CTE-based approach was chosen
- **Examples** of fault counting with visual diagrams

**Documentation Structure:**
The document `docs/faults-calculation-explanation.md` will be organized into the following sections:

1. **Overview** - High-level explanation of the fault concept
2. **Fault Definition** - Precise definition matching Rails behavior
3. **SQL Query Breakdown** - Detailed explanation of the CTE structure
4. **Visual Examples** - Diagrams showing fault counting scenarios
5. **Edge Cases** - Special handling for NULL values, single-day ranges, empty databases
6. **Rails Comparison** - Side-by-side comparison of Rails vs Go implementation
7. **Previous Incorrect Implementation** - What was wrong and how it was fixed

**Key Concepts to Document:**
- A "fault" = a day with zero pages read (sum of `read_pages` = 0 for that day)
- Days with logs but `sum(end_page - start_page) = 0` ARE counted as faults
- Days with no logs at all ARE counted as faults
- Days with one or more logs where `sum(end_page - start_page) > 0` are NOT faults

### 2. Files to Modify

**Files to Create:**
- `docs/faults-calculation-explanation.md` - Main documentation file

**Files to Read (for reference):**
- `internal/adapter/postgres/dashboard_repository.go` - Current SQL implementation
- `test/faults_rails_comparison_test.go` - Test scenarios and expected values
- `test/fixtures/dashboard/scenarios.go` - Test fixture definitions
- `backlog/tasks/rdl-146 - doc-12-Phase-5-Create-faults-calculation-documentation.md` - Task description
- `AGENTS.md` - Project context and Rails comparison details
- `prd_view doc-012` - PRD with technical decisions

### 3. Dependencies

**Prerequisites:**
- None - This is a documentation-only task
- The faults calculation implementation is already complete (Phase 4)
- Test scenarios are already implemented and passing

**Related Tasks (Already Completed):**
- RDL-139: Update GetFaultsByDateRange SQL query with CTE
- RDL-140: Update GetWeekdayFaults SQL query with CTE
- RDL-141: Verify context timeout and error handling
- RDL-142: Rewrite TestGetWeekdayFaults unit tests
- RDL-143: Update TestDashboardWeekdayFaults integration test
- RDL-144: Add edge case integration tests
- RDL-145: Add Rails comparison test cases

**Related Tasks (Parallel):**
- RDL-147: Update API documentation with fault definition (complementary documentation)

### 4. Code Patterns

**Documentation Patterns to Follow:**
- Use markdown code blocks for SQL queries with syntax highlighting
- Use ASCII diagrams for visual representations
- Follow the existing documentation style in `docs/README.go-project.md`
- Reference specific line numbers in code when applicable
- Include concrete examples with actual dates and page numbers

**Technical Writing Standards:**
- Start each section with a clear heading
- Use bullet points for lists
- Include "Key Takeaway" boxes for important concepts
- Provide both technical and non-technical explanations
- Use consistent terminology (fault vs fault day vs zero-page day)

**Visual Diagram Format:**
```
Date Range: Jan 1-10, 2024 (10 days total)

Day | Date       | Logs | Pages Read | Fault?
----|------------|------|------------|--------
 1  | Jan 1 (Mon)| Yes  | 25         | No
 2  | Jan 2 (Tue)| No   | 0          | YES ⚠️
 3  | Jan 3 (Wed)| Yes  | 30         | No
 4  | Jan 4 (Thu)| Yes  | 0          | YES ⚠️
 5  | Jan 5 (Fri)| No   | 0          | YES ⚠️
...

Total Faults: 3 out of 10 days (30%)
```

### 5. Testing Strategy

**Documentation Validation:**
Since this is a documentation task, "testing" involves:

1. **Accuracy Verification:**
   - Cross-reference all SQL queries with `dashboard_repository.go`
   - Verify all examples match test scenarios in `faults_rails_comparison_test.go`
   - Confirm fault counts match expected values in fixtures

2. **Completeness Check:**
   - Ensure all edge cases from PRD doc-012 are covered
   - Verify both `GetFaultsByDateRange` and `GetWeekdayFaults` are documented
   - Include examples for all test scenarios (30-day, 6-month, leap year)

3. **Review Process:**
   - Self-review for technical accuracy
   - Peer review for clarity and completeness
   - Engineering lead approval (as per Definition of Done)

**Validation Checklist:**
- [ ] SQL queries match current implementation exactly
- [ ] Examples use dates from existing test fixtures
- [ ] Fault counts match expected values in tests
- [ ] Edge cases (NULL, single-day, empty) are documented
- [ ] Rails comparison is accurate and complete
- [ ] Diagrams are clear and easy to understand

### 6. Risks and Considerations

**Known Issues:**
- None - Implementation is complete and tested

**Potential Pitfalls:**
1. **Documentation Drift:** The documentation must stay synchronized with code changes
   - *Mitigation:* Include version reference and last-updated date
   
2. **Complexity:** The CTE-based SQL query may be difficult for some readers to understand
   - *Mitigation:* Break down the query step-by-step with inline comments
   
3. **Terminology Confusion:** "Fault" may be ambiguous to new developers
   - *Mitigation:* Define terms clearly at the beginning and use consistently

**Deployment Considerations:**
- Documentation is version-controlled with the codebase
- No deployment steps required
- Documentation will be accessible via GitHub repository

**Trade-offs:**
- **Depth vs. Breadth:** This document focuses on faults calculation only (in scope per doc-012)
  - Other dashboard calculations are documented separately
- **Technical vs. Non-technical:** Aimed at developers, includes both high-level concepts and low-level SQL details

**Success Criteria:**
- A new developer can understand the faults calculation without reading the code
- A developer can explain why the CTE approach was chosen
- A developer can manually calculate faults for a given dataset and verify against the implementation
- The documentation serves as a reference for future maintenance and enhancements

---

**Implementation Steps Summary:**

1. **Read reference materials** (15 min)
   - Read `dashboard_repository.go` faults methods
   - Read `faults_rails_comparison_test.go`
   - Read `test/fixtures/dashboard/scenarios.go`
   - Review PRD doc-012 technical decisions

2. **Create document structure** (10 min)
   - Set up markdown headings and sections
   - Add table of contents
   - Create placeholder for diagrams

3. **Write content** (60 min)
   - Overview and fault definition
   - SQL query breakdown with CTE explanation
   - Visual examples with diagrams
   - Edge cases documentation
   - Rails comparison section
   - Previous implementation comparison

4. **Review and validate** (15 min)
   - Verify all technical details against code
   - Check examples against test fixtures
   - Proofread for clarity

5. **Submit for review** (5 min)
   - Create pull request
   - Request review from engineering lead

**Total Estimated Time:** ~2 hours
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Starting implementation of faults calculation documentation.

**Progress:**
1. ✅ Read reference materials:
   - `internal/adapter/postgres/dashboard_repository.go` - Current SQL implementation
   - `test/faults_rails_comparison_test.go` - Test scenarios
   - `test/fixtures/dashboard/scenarios.go` - Test fixture definitions
   - `prd_view doc-012` - PRD with technical decisions

2. 🔄 Creating documentation structure for `docs/faults-calculation-explanation.md`

**Next Steps:**
- Create the documentation file with all required sections
- Include SQL query breakdowns with CTE explanations
- Add visual examples and diagrams
- Document edge cases and Rails comparisons
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
