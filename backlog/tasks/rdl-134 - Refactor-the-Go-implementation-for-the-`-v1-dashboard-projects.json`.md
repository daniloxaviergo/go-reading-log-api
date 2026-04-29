---
id: RDL-134
title: Refactor the Go implementation for the `/v1/dashboard/projects.json`
status: To Do
assignee:
  - Thomas
created_date: '2026-04-29 21:42'
updated_date: '2026-04-29 21:45'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Refactor the Go implementation for the `/v1/dashboard/projects.json` endpoint to produce a JSON response that exactly matches the Rails application's JSON API format.

Use curl to check the rails-app
```bash
curl http://0.0.0.0:3001/v1/dashboard/projects.json
```

## Required Changes

### 1. Response Structure
- **Root object**: Change from `{ "projects": [...] }` to `{ "data": [...] }`
- **Each project**: Must follow JSON API format with `id`, `type`, and `attributes` wrapper
- **Remove**: Nested `project` and `logs` objects from each item

### 2. Field Mapping

| Rails Field | Current Go Field | Notes |
|-------------|------------------|-------|
| `id` (string) | `project_id` (int) | Convert to string, use `project_id` |
| `type` | N/A | Always set to `"projects"` |
| `name` | `project_name` | Direct mapping |
| `started-at` | N/A | Calculate from logs or add field |
| `progress` | `progress_geral` (per project) | Calculate percentage |
| `total-page` | `total_pages` | Direct mapping |
| `page` | `pages` | Direct mapping |
| `status` | N/A | Add field (default: `"stopped"`) |
| `logs-count` | `log_count` | Direct mapping |
| `days-unreading` | N/A | Calculate from latest log date |

### 3. Field Naming Convention
- Use **kebab-case** for attribute keys (e.g., `started-at`, `total-page`, `days-unreading`)
- Do NOT use Go's `json:"snake_case"` tags

### 4. Stats Object
Remove unnecessary fields from stats:
- **Keep**: `progress_geral`, `total_pages`, `pages`
- **Remove**: `previous_week_pages`, `last_week_pages`, `mean_day`, `spec_mean_day`, `count_pages`, `speculate_pages`

### 5. Data Type Requirements
- `id`: Must be **string** (e.g., `"446"`)
- `progress`: Must be **float** with decimals (e.g., `96.53`)
- All numeric fields: Use appropriate types (int for counts, float for percentages)

## Implementation Requirements

1. Create a new Go struct that matches the Rails JSON API format exactly
2. Remove the `logs` array from each project response
3. Calculate `progress` as: `(pages / total_page) * 100`
4. Calculate `started-at` from the earliest log date for each project
5. Calculate `days-unreading` from the latest log date to current date
6. Set `status` to `"stopped"` for all projects (or calculate based on activity)
7. Update the stats object to match Rails format exactly
8. Ensure proper JSON serialization with correct field tags
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
The endpoint /v1/dashboard/projects.json currently returns an incorrect JSON structure. Your goal is to modify the Golang code so the response matches the Rails response structure and data format precisely.

To ensure accuracy, follow this Chain-of-Thought reasoning process before writing the code:

Step 1: Analyze the Expected Rails Response Structure.
- Identify the root keys: The Rails response uses "data" (array) and "stats" (object).
- Examine the "data" array items: They follow a JSON:API style with "id" (string), "type" ("projects"), and "attributes" (object).
- List the "attributes" keys: "name", "started-at", "progress", "total-page", "page", "status", "logs-count", "days-unreading".
- Note the "stats" keys: "progress_geral", "total_pages", "pages".

Step 2: Analyze the Current Golang Response Structure.
- Identify the root keys: The Go response uses "projects" (array) and "stats" (object).
- Examine the "projects" array items: They contain a nested "project" object, a "logs" array, and top-level fields like "total_pages".
- Note the field naming differences: "project_name" vs "name", "log_count" vs "logs-count", "started-at" is missing.

Step 3: Identify Discrepancies and Mapping Requirements.
- Root Key Change: Map "projects" to "data".
- Item Structure Flattening: The nested "project" and "logs" structure must be flattened into a single "attributes" object. The "logs" array should be removed from individual items; only "logs-count" should exist in attributes.
- Field Mapping: Map "project_id" to "id" (as string). Map "project_name" to "name". Ensure "started-at", "status", and "days-unreading" are included.
- Stats Alignment: The Go "stats" object has extra fields (previous_week_pages, mean_day, etc.). These must be removed to match the Rails "stats" structure.

Step 4: Data Logic and Consistency.
- Observe that the sample data IDs differ between the Rails and Go examples. Determine if the query logic needs to be adjusted to fetch the same dataset or if this is purely a serialization fix. For this task, prioritize structural alignment to match the Rails schema.
- Ensure data types match: IDs should be strings, progress values as floats.

Step 5: Construct the Solution.
- Write the Go structs or JSON marshaling logic that enforces the structure identified in Step 3.
- Use JSON tags to ensure kebab-case keys where required (e.g., `json:"started-at"`).
- Ensure no extra fields leak into the final JSON output.

Step 6: Verification.
- Review the generated JSON structure against the Rails example provided below.
- Confirm that the root is "data", items have "attributes", and "stats" matches exactly.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress

### Analysis Phase (Complete)
- Reviewed task requirements for `/v1/dashboard/projects.json` endpoint refactoring
- Analyzed current Go implementation in `dashboard_handler.go`
- Reviewed Rails response structure expectations from task description
- Identified key changes needed:
  1. Root object: `{ "projects": [...] }` → `{ "data": [...] }`
  2. JSON:API format with `id`, `type`, `attributes` wrapper
  3. Kebab-case field naming
  4. Remove nested `project` and `logs` objects
  5. Stats object cleanup (remove unnecessary fields)
  6. ID as string type

### Implementation Plan
1. Create new DTO for JSON:API dashboard projects response
2. Update `DashboardHandler.Projects` method to use new structure
3. Update service layer to return flattened project data
4. Update tests to match new response structure
5. Run tests and verify acceptance criteria

### Next Steps
- Create new DTO structures for JSON:API format
- Modify handler to return correct structure
- Update tests
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
