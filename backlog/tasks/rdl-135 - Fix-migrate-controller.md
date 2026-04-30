---
id: RDL-135
title: Fix migrate controller
status: To Do
assignee: []
created_date: '2026-04-30 10:35'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Your objective is to migrate the controller located at rails-app/app/controllers/v1/dashboard/projects_controller.rb to Golang while ensuring filtering and ordering logic works correctly

The Golang migration is partially complete, but the logic for filtering and ordering (ORDER BY) is incorrect compared to the Rails source.

**Requirements:**
1. **Source of Truth:** Analyze the provided Rails controller code to understand the exact filtering and sorting logic.
2. **Fix Logic:** Update the Golang code to replicate the Rails behavior accurately, specifically addressing the broken filter and order by functionality.
3. **Struct Constraint:** Do **not** modify the existing Go struct definitions. The JSON response structure must remain unchanged.
4. **Validation:** Use the following command to verify the expected response from the Rails backend:
   ```bash
   curl http://0.0.0.0:3001/v1/dashboard/projects.json
   ```
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Analyze the original Rails controller code to extract all query parameters, scope filters, and ordering clauses.
2. Examine the current Golang implementation or design the new handler to identify why filters are not running and why ordering is incorrect.
3. Translate the Rails ActiveRecord query logic into equivalent Golang database queries, paying close attention to how conditions and order clauses are constructed.
4. Implement the corrected Golang code, ensuring that the filter logic matches the Rails behavior exactly.
5. Verify the ordering logic by comparing the expected SQL or query execution order with the Rails implementation.
6. Cross-reference the output using the provided curl command curl http://0.0.0.0:3001/v1/dashboard/projects.json to validate that the response data and structure are identical.
7. Finalize the Golang code, ensuring no struct changes are made and all logical steps are preserved.
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
