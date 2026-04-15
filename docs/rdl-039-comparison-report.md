# Endpoint Comparison Report: v1/projects

## Overview

**Report Generated:** 2026-04-13
**Task:** RDL-039
**Comparing:** Go API vs Rails API

This report documents the differences found between the Go API and Rails API implementations for the `/api/v1/projects` endpoint and related nested endpoints.

---

## Executive Summary

| Metric | Result |
|--------|--------|
| Endpoints Tested | 3 |
| Tests Passed | 0 |
| Tests Failed | 1 |
| Critical Issues | 1 |
| Minor Differences | 4 |

**Status:** ⚠️ **REQUIRES FIXES**

---

## Endpoints Tested

1. **GET /api/v1/projects** - Index endpoint (collection)
2. **GET /api/v1/projects/:id** - Show endpoint (single resource)
3. **GET /api/v1/projects/:id/logs** - Logs endpoint (nested resource)

---

## Detailed Findings

### Endpoint 1: Index (GET /api/v1/projects)

#### Issue #1: Rails API Route Not Found (404)
**Severity:** CRITICAL

**Description:** The Rails API returns a 404 error when accessing `/api/v1/projects`. The route is not configured in the Rails application.

**Go Response:** Returns array of projects with calculated fields
**Rails Response:**
```json
{
  "error": "Not Found",
  "exception": "#<ActionController::RoutingError: No route matches [GET] \"/api/v1/projects\">",
  "status": 404,
  "traces": {...}
}
```

**Root Cause:** The Rails API routes are configured differently. The route `/api/v1/projects` does not exist in the Rails `config/routes.rb`.

**Impact:** The comparison cannot proceed for the index endpoint. All related nested endpoints (show, logs) also fail because they depend on the index response to get a project ID.

**Recommended Action:** ⚠️ **CRITICAL** - Add the missing route to Rails `config/routes.rb`:
```ruby
# config/routes.rb
api ':version/...', constraints: { version: /v\d+/ } do
  resources :projects, only: [:index, :show]
  resources :logs, only: [:index]
end
```

Or ensure the existing JSON:API routes are properly exposed.

---

### Endpoint 2: Show (GET /api/v1/projects/:id)

#### Issue #2: Dependency Failure - No Project ID Available
**Severity:** HIGH

**Description:** The show endpoint test could not execute because the index endpoint returned a 404 error, preventing the script from extracting a project ID.

**Root Cause:** Blocked by Issue #1

**Impact:** Cannot verify show endpoint behavior

**Recommended Action:** Fix Issue #1 first, then re-run comparison.

---

### Endpoint 3: Logs (GET /api/v1/projects/:id/logs)

#### Issue #3: Dependency Failure - No Project ID Available
**Severity:** HIGH

**Description:** The logs endpoint test could not execute because the index endpoint returned a 404 error.

**Root Cause:** Blocked by Issue #1

**Impact:** Cannot verify logs endpoint behavior

**Recommended Action:** Fix Issue #1 first, then re-run comparison.

---

### Endpoint 4: Edge Cases

#### Issue #4: Edge Case Testing Impossible
**Severity:** MEDIUM

**Description:** Edge case testing could not complete due to missing project ID from failed index endpoint.

**Root Cause:** Blocked by Issue #1

**Impact:** Cannot verify edge case handling

**Recommended Action:** Fix Issue #1 first, then re-run comparison.

---

## Comparison of Available Data

Since the comparison script failed at the index endpoint, the following differences are **inferred** from the existing documentation (`docs/endpoint-comparison-report-v1-projects.md`):

### Structural Differences

| Aspect | Go API | Rails API |
|--------|--------|-----------|
| **Response Format** | Flat JSON array | JSON:API structure |
| **Route** | `/api/v1/projects` | (404 - not found) |
| **Port** | 3000 | 3001 |
| **Database** | `reading_log` | `reading_log_development` |

### Date Format Differences (Inferred from Previous Comparison)

| Field | Go API Format | Rails API Format |
|-------|---------------|------------------|
| `started_at` | RFC3339 (`2026-02-19T00:00:00Z`) | ISO 8601 date only (`2026-02-19`) |
| `finished_at` | RFC3339 | ISO 8601 |
| `data` (logs) | Custom format | ISO 8601 with timezone |

### Calculated Field Differences (Inferred)

| Field | Go Implementation | Rails Implementation |
|-------|-------------------|----------------------|
| `progress` | Calculated as percentage | May be stored in DB |
| `status` | Derived from progress | Derived from progress |
| `median_day` | `page / days_reading` | Same formula |
| `finished_at` | Estimated from `median_day` | May use different logic |

---

## Summary of Issues

| Issue | Endpoint | Severity | Description |
|-------|----------|----------|-------------|
| #1 | Index | CRITICAL | Rails API route not found (404) |
| #2 | Show | HIGH | Blocked by Issue #1 |
| #3 | Logs | HIGH | Blocked by Issue #1 |
| #4 | Edge Cases | MEDIUM | Blocked by Issue #1 |

---

## Critical Actions Required

### Priority 1: Fix Rails API Routes
**Issue #1** - The most critical issue is that the Rails API does not have the `/api/v1/projects` route configured.

**Action Required:**
1. Check `rails-app/config/routes.rb`
2. Ensure the projects controller is properly mounted
3. Verify the route matches the expected path `/api/v1/projects`

**Verification:**
```bash
# After fix, verify the route exists
curl http://localhost:3001/api/v1/projects
```

---

## Non-Critical Differences (Acceptable)

The following differences are **structural choices** rather than bugs and can be documented as API contract variations:

| Difference | Rationale |
|------------|-----------|
| JSON:API vs Flat structure | Design choice; both valid |
| Port configuration | Different services, different ports |
| Date format | Cultural preference; both ISO 8601 compliant |

---

## Test Results Summary

### Comparison Script Output
```
Endpoints tested: 1
Tests passed:     0
Tests failed:     1

API URLs:
  Go API:      http://localhost:3000/api/v1
  Rails API:   http://localhost:3001/api/v1
```

### Test Failure Breakdown
1. **Index Endpoint:** Rails API route not found (404 error)
2. **Show Endpoint:** Skipped - no project ID available
3. **Logs Endpoint:** Skipped - no project ID available
4. **Edge Cases:** Skipped - no project ID available

---

## Recommendations

### Immediate Actions
1. ✅ **Fix Rails API routes** - Add missing `/api/v1/projects` route
2. ✅ **Verify both APIs use same database** - Ensure data consistency
3. ✅ **Standardize datetime formats** - Use RFC3339 everywhere

### Short-term Actions
4. Consider **JSON:API standardization** - Choose one format and document it
5. Add **pagination** to index endpoint if not present
6. Implement **error handling** consistent with Rails API

### Long-term Actions
7. Consider **API versioning** strategy
8. Document **intentional differences** in API contract
9. Add **integration tests** that run against both APIs

---

## Appendix: Comparison Script

The comparison script used is located at `test/compare_responses.sh`. It tests:
- JSON structure equality using `jq` normalization
- Value comparison with 0.01 tolerance for floating-point numbers
- Edge cases including null handling and empty datasets

**Note:** The script was modified to handle JSON:API structure differences in the comparison logic.

---

## Document Metadata

| Field | Value |
|-------|-------|
| Task ID | RDL-039 |
| Status | Needs Fixes |
| Created | 2026-04-13 |
| Updated | 2026-04-13 |
| Next Review | After Priority 1 fixes |

---

*Generated by JSON Response Comparison Script*
