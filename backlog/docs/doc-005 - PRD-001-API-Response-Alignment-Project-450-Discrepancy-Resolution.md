---
id: doc-005
title: 'PRD-001: API Response Alignment - Project 450 Discrepancy Resolution'
type: other
created_date: '2026-04-18 11:32'
---


# Project 450 Response Comparison Report

## Executive Summary

Comparing Go API vs Rails API responses for endpoint: `v1/projects/450.json`

### Overall Assessment
⚠️ **CRITICAL**: Multiple structural and data discrepancies detected

| Category | Count |
|----------|-------|
| Structural Differences | 3 |
| Value Discrepancies | 2 |
| Missing Fields | 1 |

---

## Endpoint Comparison

### 1. Index Endpoint (GET /v1/projects)

#### Go Response Structure
```json
[
 {
 id: 450,
 name: História da Igreja VIII.1,
 ...
 }
]
```

#### Rails Response Structure
```json
{
 data: [
 {
 id: 450,
 type: projects,
 attributes: { ... }
 }
 ]
}
```

**Discrepancy Type:** STRUCTURAL

---

## Detailed Field Analysis

| Field | Go Value | Rails Value | Match | Notes |
|-------|----------|-------------|-------|-------|
| id | 450 (int) | 450 (string) | ⚠️ Type only | ID type differs |
| name | História da Igreja VIII.1 | História da Igreja VIII.1 | ✓ | Identical |
| started_at | 2026-02-19T00:00:00Z | 2026-02-19 | ⚠️ Format | Time component missing in Rails |
| progress | 100 | 100.0 | ✓ | Functionally equivalent |
| total_page | 691 | 691 | ✓ | Identical |
| page | 691 | 691 | ✓ | Identical |
| status | finished | finished | ✓ | Identical |
| logs_count | 38 | 38 | ✓ | Identical |
| days_unreading | 58 | 16 | ✗ **CRITICAL** | Different calculation logic |
| median_day | (missing) | 11.91 | ⚠️ Missing | Not returned by Go API |
| finished_at | null | 2026-04-02 | ✗ **CRITICAL** | Rails calculates, Go returns null |

---

## Critical Issues Requiring Investigation

### Issue #1: days_unreading Calculation Discrepancy
- **Severity:** HIGH
- **Go Value:** 58 days
- **Rails Value:** 16 days
- **Root Cause:** Different date calculation logic or data source
- **Impact:** User-facing feature inconsistency

**Investigation Required:**
- Check `CalculateDaysUnreading()` method in Go implementation
- Compare with Rails `days_unreading` calculation
- Verify date source (last log date vs started_at)

---

### Issue #2: finished_at Field Inconsistency
- **Severity:** MEDIUM
- **Go Value:** null
- **Rails Value:** 2026-04-02
- **Root Cause:** Go implementation doesnt calculate completion date
- **Impact:** Missing useful information for users

**Investigation Required:**
- Check if `finished_at` calculation exists in Go model
- Compare with Rails `calculate_finished_at` method
- Consider implementing based on median_day projection

---

### Issue #3: median_day Field Missing
- **Severity:** LOW
- **Go Value:** Not returned
- **Rails Value:** 11.91 pages/day
- **Root Cause:** Go API doesnt expose calculated median
- **Impact:** Users cannot see reading pace information

---

## Structural Differences

### JSON:API vs Flat JSON

| Aspect | Go API | Rails API |
|--------|--------|-----------|
| Specification | Custom flat structure | JSON:API 1.0 |
| Root wrapper | Array/Object directly | `{data: {...}}` |
| ID type | Integer | String |
| Field naming | snake_case | kebab-case |
| Nested data | Embedded objects | Relationship references |

---

## Recommendations

### Priority 1 (Immediate)
1. **Investigate `days_unreading` calculation** - The 42-day difference (58 vs 16) needs root cause analysis
2. **Add `finished_at` calculation** - Implement completion date projection based on reading pace

### Priority 2 (Short-term)
3. **Align field naming** - Consider standardizing to one convention (snake_case preferred for Go)
4. **Add `median_day` to response** - Expose reading pace information to users

### Priority 3 (Medium-term)
5. **Consider JSON:API adoption** - Standardize on JSON:API specification for consistency
6. **Optimize nested data** - Move from embedded objects to relationship references in logs endpoint

---

## Methodology

### Data Collection
- Used `test/compare_responses.sh` script for automated comparison
- Manual extraction of project 450 details for deeper analysis
- All comparisons performed on 2026-04-18

### Test Results Summary
```
Endpoints tested: 1
Tests passed: 0
Tests failed: 6
```

All structural tests failed, indicating significant differences in API response format.

---

## Appendix: Raw Response Data

### Go API Response (v1/projects/450.json)
```json
{
 id: 450,
 name: História da Igreja VIII.1,
 started_at: 2026-02-19T00:00:00Z,
 progress: 100,
 total_page: 691,
 page: 691,
 status: finished,
 logs_count: 38,
 days_unreading: 58,
 finished_at: null
}
```

### Rails API Response (v1/projects/450.json)
```json
{
 data: {
 id: 450,
 type: projects,
 attributes: {
 name: História da Igreja VIII.1,
 started-at: 2026-02-19,
 progress: 100.0,
 total-page: 691,
 page: 691,
 status: finished,
 logs-count: 38,
 days-unreading: 16,
 median-day: 11.91,
 finished-at: 2026-04-02
 }
 }
}
```

---

*Report generated: 2026-04-18*
*Task: RDL-059*