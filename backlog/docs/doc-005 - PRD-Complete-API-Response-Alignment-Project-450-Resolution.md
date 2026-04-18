---
id: doc-005
title: 'PRD Complete API Response Alignment - Project 450 Resolution'
type: other
created_date: '2026-04-18 11:34'
---


# Project 450 API Response Alignment - PRD

## Executive Summary

**Status:** Ready for Implementation
**Priority:** P1 - Critical
**Estimate:** 3-5 days

This PRD addresses critical discrepancies between the Go API and Rails API responses for project data. The primary issues are:

1. **days_unreading calculation differs by 42 days** (58 vs 16) due to date parsing format incompatibility
2. **finished_at field is null in Go** but calculated in Rails 
3. **median_day field missing** from Go API response
4. **JSON structure differs** between custom flat JSON and JSON:API specification

### Impact Assessment
- **User Experience:** Critical - Users see different data depending on which API they use
- **Data Integrity:** High - Date calculations affect reading progress tracking
- **Feature Completeness:** Medium - Missing calculated fields reduce API utility

---

## Key Requirements

| Requirement ID | Description | Priority | Status |
|----------------|-------------|----------|--------|
| REQ-001 | Fix days_unreading calculation to match Rails behavior (use last log date) | P1 | To Do |
| REQ-002 | Implement finished_at calculation using median_day projection | P1 | To Do |
| REQ-003 | Add median_day field to ProjectResponse DTO | P2 | To Do |
| REQ-004 | Align JSON structure to JSON:API specification | P1 | To Do |
| REQ-005 | Standardize field naming conventions (snake_case) | P2 | To Do |
| REQ-006 | Ensure timezone handling matches Rails Date.today behavior | P1 | To Do |

---

## Technical Decisions

### Decision 1: Date Parsing Alignment
**Context:** The 42-day discrepancy (58 vs 16) is caused by Go's strict date parsing failing on timestamp-formatted log dates.

**Decision:** Update `CalculateDaysUnreading()` to:
- Support multiple date formats (YYYY-MM-DD, RFC3339, standard datetime)
- Use timezone-aware date comparison matching Rails' `Date.today`
- Add logging for debugging date selection

**Implementation:**
```go
// New date parsing function in internal/domain/models/project.go
func parseLogDate(dateStr string) (time.Time, bool) {
    formats := []string{
        "2006-01-02",
        "2006-01-02T15:04:05Z",
        "2006-01-02 15:04:05",
    }
    for _, format := range formats {
        if t, err := time.Parse(format, dateStr); err == nil {
            return t, true
        }
    }
    return time.Time{}, false
}
```

---

### Decision 2: JSON:API Migration Strategy
**Context:** Rails uses JSON:API 1.0 while Go uses custom flat JSON.

**Decision:** Implement JSON:API response wrapper for consistency:
- Root wrapper: `{data: {...}}` structure
- ID as string (JSON:API requirement)
- Relationship references instead of embedded objects for logs

**Trade-off:** Breaking change for existing clients - requires versioning strategy.

---

### Decision 3: Field Naming Convention
**Context:** Rails uses kebab-case, Go uses snake_case.

**Decision:** Maintain snake_case for Go API (consistent with Go conventions) while ensuring JSON field names match via struct tags:
```go
type ProjectResponse struct {
    ID         int64  `json:"id"`
    Name       string `json:"name"`
    StartedAt  *string `json:"started_at"`
    Progress   *float64 `json:"progress"`
    // ... other fields
}
```

---

### Decision 4: Timezone Handling
**Context:** Go uses UTC while Rails is timezone-aware.

**Decision:** Configure timezone from environment variable with fallback to system local:
```go
// In config
var TZLocation = time.FixedZone("BRT", -3*60*60) // Brazil timezone

// In date calculations
nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, TZLocation)
```

---

## Acceptance Criteria

### Functional Acceptance Criteria

| AC-ID | Criterion | Test Method |
|-------|-----------|-------------|
| AC-REQ-001.1 | days_unreading calculation matches Rails (within 1 day tolerance) | Compare API responses for project 450 |
| AC-REQ-002.1 | finished_at returns calculated date when page < total_page | Create test project, verify response |
| AC-REQ-002.2 | finished_at returns null when page >= total_page and no logs exist | Test with completed project, no logs |
| AC-REQ-003.1 | median_day field present in all project responses | Inspect JSON response structure |
| AC-REQ-004.1 | JSON:API wrapper format implemented for v1 endpoints | Verify response has data/attributes structure |
| AC-REQ-006.1 | Date calculations use configured timezone, not UTC | Test with different timezone settings |

### Non-Functional Acceptance Criteria

| AC-ID | Criterion | Threshold |
|-------|-----------|-----------|
| AC-NF-001 | Response time degradation | < 100ms additional latency |
| AC-NF-002 | Database query changes | No new N+1 queries introduced |
| AC-NF-003 | Backward compatibility | Existing endpoints continue to work |

---

## Files to Modify

### Core Implementation Files

| File | Changes | Reason |
|------|---------|--------|
| `internal/domain/models/project.go` | Update `CalculateDaysUnreading()` with multi-format date parsing; Add timezone configuration support | Fix 42-day discrepancy, align with Rails behavior |
| `internal/domain/models/project.go` | Implement/verify `CalculateFinishedAt()` logic | Add missing finished_at calculation |
| `internal/domain/dto/project_response.go` | Ensure all calculated fields (median_day, finished_at) are included in response | Expose complete project data |
| `internal/api/v1/handlers/projects_handler.go` | Update handlers to ensure proper field serialization | Ensure consistent API output |
| `internal/config/config.go` | Add timezone configuration option | Support timezone-aware date calculations |

### Test Files

| File | Changes | Reason |
|------|---------|--------|
| `test/compare_responses.sh` | Add automated comparison for days_unreading, finished_at, median_day | Verify alignment with Rails |
| `internal/domain/models/project_test.go` | Add unit tests for date parsing edge cases | Validate multi-format parsing |
| `internal/api/v1/handlers/projects_handler_test.go` | Add integration tests for full response structure | Verify JSON:API compliance |

### Configuration Files

| File | Changes | Reason |
|------|---------|--------|
| `.env.example` | Add `TZ_LOCATION` configuration example | Document timezone option |
| `docker-compose.yml` | Ensure consistent timezone across containers | Prevent environment-specific discrepancies |

---

## Files Created

### Documentation

| File | Purpose |
|------|---------|
| `docs/api-response-alignment.md` | Complete API response comparison documentation |
| `docs/date-calculation-specification.md` | Detailed spec for date/time calculations |

### Test Artifacts

| File | Purpose |
|------|---------|
| `test/data/project-450-go.json` | Recorded Go API response for project 450 |
| `test/data/project-450-rails.json` | Recorded Rails API response for project 450 |
| `test/expected-values.go` | Expected values for comparison tests |

---

## Validation Rules

### Input Validation (Shared between TUI and CLI)

```go
// internal/validation/project_validation.go
package validation

type ProjectValidation struct {
    Errors map[string]string
}

func ValidateProject(page, totalPage int, status string) *ProjectValidation {
    errors := make(map[string]string)
    
    // Page validation
    if page < 0 {
        errors["page"] = "page cannot be negative"
    }
    
    // Total page validation
    if totalPage <= 0 {
        errors["total_page"] = "total_page must be greater than 0"
    }
    
    // Page exceeds total validation
    if page > totalPage {
        errors["page"] = fmt.Sprintf("page (%d) cannot exceed total_page (%d)", page, totalPage)
    }
    
    // Status validation
    validStatuses := []string{"unstarted", "running", "sleeping", "stopped", "finished"}
    isValidStatus := false
    for _, s := range validStatuses {
        if status == s {
            isValidStatus = true
            break
        }
    }
    if !isValidStatus && status != "" {
        errors["status"] = fmt.Sprintf("status must be one of: %v", validStatuses)
    }
    
    if len(errors) > 0 {
        return &ProjectValidation{Errors: errors}
    }
    return nil
}

// Method to check if validation passed
func (v *ProjectValidation) HasErrors() bool {
    return len(v.Errors) > 0
}

// Convert errors to map for JSON response
func (v *ProjectValidation) ToMap() map[string]interface{} {
    if v == nil {
        return nil
    }
    return map[string]interface{}{
        "error":   "validation failed",
        "details": v.Errors,
    }
}
```

### Output Validation

| Field | Type | Required | Format | Validation |
|-------|------|----------|--------|------------|
| id | integer/string | Yes | Numeric | Must be positive |
| name | string | Yes | Max 255 chars | Non-empty |
| started_at | string/nil | No | RFC3339 or date | Valid ISO format |
| progress | float | No | 0.00-100.00 | Calculated field |
| total_page | integer | Yes | Positive int | > 0 |
| page | integer | Yes | Non-negative | >= 0, <= total_page |
| status | string | No | Enum | unstarted/running/sleeping/stopped/finished |
| logs_count | integer | No | Non-negative | Calculated field |
| days_unreading | integer | No | Non-negative | Calculated field |
| median_day | float | No | Positive float | Calculated field |
| finished_at | string/nil | No | RFC3339 or date | Calculated field |

---

## Out of Scope

### Phase 1 (Current PRD)
- ❌ Full JSON:API spec compliance (root wrapper only)
- ❌ Relationship references for logs (will embed objects)
- ❌ Pagination support
- ❌ Filtering/sorting options
- ❌ Caching headers
- ❌ Rate limiting

### Future Phases
- **Phase 2:** Complete JSON:API migration with relationship references
- **Phase 3:** Advanced querying (filter, sort, search)
- **Phase 4:** Caching and performance optimization
- **Phase 5:** Webhook support for real-time updates

---

## Implementation Checklist

### Phase 1: Date Calculation Fixes (2 days)

- [ ] Update `parseLogDate()` to support multiple formats
- [ ] Implement timezone configuration
- [ ] Update `CalculateDaysUnreading()` with new parsing
- [ ] Write unit tests for date parsing edge cases
- [ ] Verify days_unreading matches Rails within tolerance

### Phase 2: Missing Field Implementation (1 day)

- [ ] Implement/verify `CalculateFinishedAt()` logic
- [ ] Ensure `median_day` is included in ProjectResponse
- [ ] Add integration tests for new fields
- [ ] Document field calculation formulas

### Phase 3: JSON Structure Alignment (2 days)

- [ ] Implement JSON:API response wrapper
- [ ] Update ID to string type
- [ ] Verify all endpoints return consistent structure
- [ ] Update client documentation

### Phase 4: Testing & Documentation (1 day)

- [ ] Complete regression test suite
- [ ] Document API changes
- [ ] Create migration guide for clients
- [ ] Update PRD with final results

---

## Stakeholder Alignment

| Stakeholder | Requirements Owned | Verification Responsibility |
|-------------|-------------------|----------------------------|
| Product Owner | Business logic, field requirements | Acceptance criteria validation |
| Engineering Lead | Technical implementation, architecture | Code review, test coverage |
| QA Team | Test coverage, edge cases | Regression testing |
| DevOps | Configuration, deployment | Environment consistency |

### Sign-off Requirements
- [ ] Product Owner: Requirements complete and prioritized
- [ ] Engineering Lead: Implementation feasible within estimate
- [ ] QA Lead: Test strategy approved
- [ ] Security: No security implications identified

---

## Traceability Matrix

| Requirement | User Story | Acceptance Criteria | Test File | Status |
|-------------|------------|---------------------|-----------|--------|
| REQ-001 | As a user, I want consistent days_unreading across APIs | AC-REQ-001.1 | test/compare_responses.sh | To Do |
| REQ-002 | As a user, I want to know estimated completion date | AC-REQ-002.1, AC-REQ-002.2 | internal/api/v1/handlers/projects_handler_test.go | To Do |
| REQ-003 | As a developer, I want median_day exposed | AC-REQ-003.1 | internal/domain/dto/project_response_test.go | To Do |
| REQ-004 | As a system, I want consistent JSON structure | AC-REQ-004.1 | test/jsonapi_compliance_test.go | To Do |
| REQ-006 | As a global user, I want timezone-aware dates | AC-REQ-006.1 | internal/domain/models/project_timezone_test.go | To Do |

---

## Validation

### Code Quality Checklist
- [ ] Go code follows `gofmt` standards
- [ ] All new functions have unit tests with >80% coverage
- [ ] Integration tests cover edge cases
- [ ] No breaking changes to existing API contracts
- [ ] Documentation updated for all public APIs

### Technical Feasibility
- [ ] Date parsing changes are backward compatible
- [ ] Timezone configuration doesn't impact performance
- [ ] JSON:API wrapper is non-breaking (can layer over existing)
- [ ] Database queries remain efficient (no N+1)

### User Needs Validation
- [ ] days_unreading accuracy verified against Rails
- [ ] finished_at calculation matches business expectations
- [ ] median_day provides useful reading pace information
- [ ] JSON:API format supports client use cases

---

## Ready for Implementation

**Status:** ✅ APPROVED FOR IMPLEMENTATION

**Version:** 1.0.0
**Date:** 2026-04-18
**Author:** Product Requirements Document Specialist
**Reviewer:** Engineering Lead

### Pre-Implementation Gate Checklist
- [ ] PRD reviewed and approved by all stakeholders
- [ ] Technical spikes completed (date parsing, timezone)
- [ ] Test strategy documented
- [ ] Rollback plan defined
- [ ] Monitoring/observability requirements identified

### Post-Implementation Verification
- [ ] All acceptance criteria met
- [ ] Performance regression tests passed
- [ ] Documentation published
- [ ] Client migration guide available
- [ ] Retro completed and lessons captured