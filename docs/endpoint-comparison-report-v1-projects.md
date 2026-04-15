# Endpoint Comparison Report: v1/projects.json

## Overview

**Report Generated:** 2026-04-12
**Task:** RDL-039
**Comparing:** Go API vs Rails API

This report documents the differences found between the Go API and Rails API implementations for the `/api/v1/projects.json` endpoint and related nested endpoints.

---

## Executive Summary

| Metric | Result |
|--------|--------|
| Endpoints Tested | 3 |
| Tests Passed | 0 |
| Tests Failed | 3 |
| Critical Issues | 2 |
| Minor Differences | 4 |

**Status:** ❌ **REQUIRES FIXES**

---

## Endpoints Tested

1. **GET /api/v1/projects** - Index endpoint (collection)
2. **GET /api/v1/projects/:id** - Show endpoint (single resource)
3. **GET /api/v1/projects/:id/logs** - Logs endpoint (nested resource)

---

## Detailed Findings

### Endpoint 1: Index (GET /api/v1/projects)

#### Issue #1: Different Project Data Returned
**Severity:** HIGH

**Description:** The Go API and Rails API return completely different projects from the same database query. The Go API returns project ID 1 ("Filocalia"), while the Rails API returns project ID 450 ("História da Igreja VIII.1").

**Go Response (First Project):**
```json
{
  "days_unreading": 3354,
  "finished_at": null,
  "id": 1,
  "logs_count": 50,
  "name": "Filocalia",
  "page": 655,
  "progress": null,
  "started_at": "2017-02-04T00:00:00Z",
  "status": "stopped",
  "total_page": 1267
}
```

**Rails Response (First Project):**
```json
{
  "days_unreading": 10,
  "finished_at": null,
  "id": 450,
  "logs_count": 38,
  "name": "História da Igreja VIII.1",
  "page": 691,
  "progress": 100.0,
  "started_at": "2026-02-19",
  "status": "finished",
  "total_page": 691
}
```

**Root Cause:** Both APIs are querying different databases:
- Go API: Using `reading_log` database
- Rails API: Using `reading_log_development` database

**Impact:** All calculated fields differ because they're operating on different data.

**Recommended Action:** ⚠️ **CRITICAL** - Ensure both APIs connect to the same database instance.

---

#### Issue #2: Different JSON Structure (JSON:API vs Flat)
**Severity:** MEDIUM

**Description:** The Rails API returns JSON:API format while the Go API returns a flat array structure.

**Go Structure:**
```json
[
  {
    "id": 1,
    "name": "Filocalia",
    "total_page": 1267,
    ...
  }
]
```

**Rails Structure (JSON:API):**
```json
{
  "data": [
    {
      "id": "450",
      "type": "projects",
      "attributes": {
        "name": "História da Igreja VIII.1",
        "total-page": 691,
        ...
      }
    }
  ]
}
```

**Root Cause:** Rails API uses Active Model Serializers with JSON:API adapter; Go API returns plain JSON.

**Impact:** Client code must handle two different response formats.

**Recommended Action:** Consider standardizing on one format or documenting the difference explicitly.

---

### Endpoint 2: Show (GET /api/v1/projects/:id)

#### Issue #3: Date Format Inconsistency
**Severity:** HIGH

**Description:** Timestamp fields use different formats between APIs.

**Go API (`started_at`):**
```json
"started_at": "2026-02-19T00:00:00Z"
```
Format: RFC3339 with 'Z' suffix (UTC)

**Rails API (`started_at`):**
```json
"started_at": "2026-02-19"
```
Format: ISO 8601 date only (no time component)

**Root Cause:** 
- Go uses `time.Time` marshaling which defaults to RFC3339
- Rails uses date-only serialization for this field

**Impact:** Date parsing differs; time information is lost in Rails response.

**Recommended Action:** Align date formats. Consider:
- Using `2026-02-19T00:00:00Z` format in Rails
- Or using `2026-02-19` format in Go for consistency

---

#### Issue #4: `finished_at` Field Differences
**Severity:** MEDIUM

**Description:** The `finished_at` field shows different behavior.

**Go API:**
```json
"finished_at": null
```

**Rails API:**
```json
"finished_at": "2026-04-02"
```

**Root Cause:** Different calculation logic or data state for the finished_at field.

**Impact:** Inconsistent status tracking for completed projects.

**Recommended Action:** Investigate and unify the `finished_at` calculation logic.

---

#### Issue #5: `progress` Field Type Difference
**Severity:** LOW

**Description:** The `progress` field type differs between APIs.

**Go API:**
```json
"progress": null
```

**Rails API:**
```json
"progress": 100.0
```

**Root Cause:** Go implementation may not calculate progress for all cases; Rails always calculates.

**Impact:** Clients expecting progress percentage may receive null in Go API.

**Recommended Action:** Ensure consistent progress calculation in Go API.

---

### Endpoint 3: Logs (GET /api/v1/projects/:id/logs)

#### Issue #6: Nested Project Object in Logs
**Severity:** MEDIUM

**Description:** The Go API includes a full project object within each log entry, while the Rails API does not.

**Go API Log Structure:**
```json
{
  "data": "2026-04-02 21:21:53",
  "end_page": 691,
  "id": 9092,
  "note": null,
  "project": {
    "days_unreading": null,
    "finished_at": null,
    "id": 450,
    "logs_count": null,
    "name": "História da Igreja VIII.1",
    "page": 691,
    "progress": null,
    "started_at": "2026-02-19T00:00:00Z",
    "status": null,
    "total_page": 691
  },
  "start_page": 665
}
```

**Rails API Log Structure:**
```json
{
  "data": "2026-04-02T18:21:53.000-03:00",
  "end_page": 691,
  "id": 9092,
  "note": null,
  "start_page": 665
}
```

**Root Cause:** Go API includes associated project data via JOIN query; Rails API follows JSON:API relationships pattern.

**Impact:** Response size differs; clients expecting project data need different handling.

**Recommended Action:** Consider adding a query parameter to control include depth (e.g., `?include=project`).

---

#### Issue #7: Date Format in Logs
**Severity:** HIGH

**Description:** Log `data` field uses different formats.

**Go API:**
```json
"data": "2026-04-02 21:21:53"
```
Format: Custom datetime string (not RFC3339 compliant)

**Rails API:**
```json
"data": "2026-04-02T18:21:53.000-03:00"
```
Format: ISO 8601 with timezone offset

**Root Cause:** Go API uses custom string formatting; Rails uses standard JSON:API date format.

**Impact:** Timezone handling differs; parsing requires different logic.

**Recommended Action:** ⚠️ **CRITICAL** - Standardize on RFC3339 format for all datetime fields.

---

## Summary of Issues

| Issue | Endpoint | Severity | Description |
|-------|----------|----------|-------------|
| #1 | Index | HIGH | Different databases being queried |
| #2 | Index | MEDIUM | JSON:API vs Flat structure |
| #3 | Show | HIGH | Date format inconsistency |
| #4 | Show | MEDIUM | `finished_at` field differences |
| #5 | Show | LOW | `progress` field type difference |
| #6 | Logs | MEDIUM | Nested project object in logs |
| #7 | Logs | HIGH | Date format in logs field |

---

## Critical Actions Required

### Priority 1: Database Configuration
**Issue #1** - The most critical issue is that the Go and Rails APIs are connected to different databases. This means:
- All comparison data is invalid
- Users get different data depending on which API they use
- Data consistency cannot be ensured

**Action:** Configure both APIs to use the same PostgreSQL database (`reading_log`).

### Priority 2: Datetime Format Standardization
**Issues #3, #7** - Inconsistent datetime formats make it difficult for clients to parse timestamps.

**Action:** 
- Apply RFC3339 format consistently across all datetime fields
- Ensure timezone information is preserved
- Consider adding `Z` suffix for UTC times

---

## Non-Critical Differences (Acceptable)

The following differences are structural choices rather than bugs:

| Difference | Rationale |
|------------|-----------|
| JSON:API vs Flat structure | Design choice; both valid |
| `logs_count` field presence | Can be derived from array length |
| Field ordering | JSON keys are unordered by spec |

---

## Test Results

### Comparison Script Output
```
Endpoints tested: 3
Tests passed: 0
Tests failed: 3

API URLs:
  Go API:      http://localhost:3000/api/v1
  Rails API:   http://localhost:3001/v1
```

### Test Failures Breakdown
1. **Index Endpoint:** Structure mismatch (different projects, different JSON format)
2. **Show Endpoint:** Structure mismatch (date formats, field differences)
3. **Logs Endpoint:** Structure mismatch (nested object, date formats)

---

## Recommendations

### Immediate Actions
1. ✅ **Fix database connection** - Point Go API to the correct database
2. ✅ **Standardize datetime formats** - Use RFC3339 everywhere
3. ✅ **Align `finished_at` calculation** - Ensure consistent logic

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
| Created | 2026-04-12 |
| Updated | 2026-04-12 |
| Next Review | After Priority 1 fixes |

---

*Generated by JSON Response Comparison Script*
