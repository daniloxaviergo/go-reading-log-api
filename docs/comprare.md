# JSON Response Comparison Report Plan

## Overview

This document outlines the plan for creating a JSON response comparison report using `test/compare_responses.sh` to verify that the Go API and Rails API return identical JSON responses.

## Current Status

### Script Location
- **compare_responses.sh**: `/home/danilo/scripts/github/go-reading-log-api-next/test/compare_responses.sh`

### Script Features
- Compares JSON responses between Go API (port 3000) and Rails API (port 3001)
- Tests 3 endpoints:
  - `GET /api/v1/projects`
  - `GET /api/v1/projects/{id}`
  - `GET /api/v1/projects/{id}/logs`
- Uses `curl` + `jq` for JSON comparison
- Supports tolerance for floating-point numbers (0.01)
- Tests edge cases (empty logs, null dates)

### Issue Identified
The routes in `compare_responses.sh` currently do NOT have `.json` suffix. The plan is to update them to use `.json` suffix.

## Plan: Update Routes to Use .json Suffix

### Routes to Update

| Current Route | Updated Route |
|---------------|---------------|
| `/api/v1/projects` | `/api/v1/projects.json` |
| `/api/v1/projects/:id` | `/api/v1/projects/:id.json` |
| `/api/v1/projects/:id/logs` | `/api/v1/projects/:id/logs.json` |

### Files to Modify

1. **test/compare_responses.sh**
   - Update `GO_API_URL` default to include `.json` or add suffix per endpoint
   - Update `RAILS_API_URL` default similarly
   - Update `test_index_endpoint()` - fetch from `/api/v1/projects.json`
   - Update `test_show_endpoint()` - fetch from `/api/v1/projects/{id}.json`
   - Update `test_logs_endpoint()` - fetch from `/api/v1/projects/{id}/logs.json`
   - Update help text to reflect `.json` routes

2. **docs/comprare.md** (this file)
   - Document the comparison plan and results

## Implementation Steps

### Step 1: Update compare_responses.sh
Modify the script to use `.json` suffix for all API endpoints:

```bash
# Current (line 21-22):
GO_API_URL="${GO_API_URL:-http://localhost:3000/api/v1}"
RAILS_API_URL="${RAILS_API_URL:-http://localhost:3001/api/v1}"

# Updated:
GO_API_URL="${GO_API_URL:-http://localhost:3000/api/v1}"
RAILS_API_URL="${RAILS_API_URL:-http://localhost:3001/api/v1}"

# Then add .json suffix in fetch calls:
fetch_json "$GO_API_URL/projects.json" "$go_file"
fetch_json "$RAILS_API_URL/projects.json" "$rails_file"
```

### Step 2: Run the Comparison Script
Execute the updated script:

```bash
# Make script executable if needed
chmod +x test/compare_responses.sh

# Run the comparison
./test/compare_responses.sh
```

### Step 3: Generate Report
Create `docs/comprare.md` with:
- Test results summary (pass/fail per endpoint)
- Field-by-field differences (if any)
- Performance metrics
- Recommendations

## Expected Outcomes

### Success Criteria
- All 3 endpoints return identical JSON structure
- All field values match within tolerance
- Edge cases handled correctly

### Report Contents
1. **Executive Summary**
   - Pass/fail status
   - Number of tests passed/failed

2. **Per-Endpoint Results**
   - Index endpoint (/api/v1/projects.json)
   - Show endpoint (/api/v1/projects/:id.json)
   - Logs endpoint (/api/v1/projects/:id/logs.json)

3. **Detailed Analysis**
   - Fields that differ
   - Values that differ
   - Data type mismatches

4. **Edge Cases**
   - Empty logs handling
   - Null date handling

5. **Recommendations**
   - Fixes needed
   - Further testing required

## Prerequisites

### Required Tools
- `curl` - HTTP client for fetching JSON
- `jq` - JSON parser for comparison (version 1.6+ recommended)
- `bash` - Shell for script execution

### Required Services
- Go API running on `http://localhost:3000`
- Rails API running on `http://localhost:3001`

### Required Data
- Test database must have at least one project with logs

## Notes

- The script already handles most edge cases (null values, empty arrays)
- Floating-point comparisons use 0.01 tolerance
- JSON structures are normalized (sorted keys) before comparison
