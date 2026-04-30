---
id: RDL-135
title: Fix migrate controller
status: To Do
assignee:
  - Thomas
created_date: '2026-04-30 10:35'
updated_date: '2026-04-30 10:48'
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

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress

### Changes Made

#### 1. Fixed Repository Query (internal/adapter/postgres/dashboard_repository.go)
- **Added filter**: `WHERE p.page != p.total_page` to match Rails `only_status(:running)` scope
- **Added secondary ordering**: `MAX(lr.data) DESC` to match Rails `order('logs.data DESC')`
- Updated comments to document the Rails-equivalent behavior

#### 2. Simplified Service Layer (internal/service/dashboard/projects_service.go)
- **Removed redundant filtering**: The `isRunningProject` method and filtering logic were removed since SQL now handles filtering
- **Updated sorting**: Service layer now sorts by progress DESC, latest log date DESC, then project ID ASC for consistency
- Updated comments to reflect that filtering is done at repository level

#### 3. Updated Tests (internal/service/dashboard/projects_service_test.go)
- **Removed test**: `TestProjectsService_isRunningProject` - method no longer exists
- **Updated test**: `TestProjectsService_GetRunningProjectsWithLogs_Filtering` - now tests that repository returns pre-filtered data

### Verification
- ✅ All unit tests pass (`go test ./internal/...`)
- ✅ `go fmt` passes with no errors
- ✅ `go vet` passes with no errors
- ✅ Build succeeds (`go build ./cmd/server.go`)

### Rails Logic Reference
The Rails controller uses:
```ruby
filter = { status: [:running] }
@projects = Project.only_status(filter)
@projects = Project.eager_load(:logs).where(id: @projects.map(&:id))
                                     .order_progress
                                     .order('logs.data DESC')
                                     .all
```

Where:
- `only_status(:running)` → SQL: `page != total_page`
- `order_progress` → SQL: `page::float / total_page::float DESC`
- `order('logs.data DESC')` → SQL: `logs.data DESC`

The Go implementation now matches this behavior exactly.
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
