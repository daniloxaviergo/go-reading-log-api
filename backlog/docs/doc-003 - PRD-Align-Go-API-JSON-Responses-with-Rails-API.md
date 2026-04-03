---
id: doc-003
title: 'PRD: Align Go API JSON Responses with Rails API'
type: other
created_date: '2026-04-03 13:54'
---


# PRD: Align Go API JSON Responses with Rails API

# Executive Summary

This PRD addresses the requirement to make the Go and Rails API applications return **identical JSON responses** for the three main endpoints. Both applications connect to the same PostgreSQL database but currently return different response formats due to missing derived calculations in the Go implementation.

**Why necessary**: To enable gradual migration from Rails to Go without breaking client applications that depend on consistent API responses.

**Scope**: Align `GET /v1/projects`, `GET /v1/projects/{project_id}`, and `GET /v1/projects/{project_id}/logs` endpoints to return identical JSON structure, field names, data types, and values.

**Value**: Enables rollback strategy, A/B testing, and smooth deprecation of Rails API.

---

# Key Requirements

| Requirement | Status | Notes |
|-------------|--------|-------|
| Field names match (snake_case) | ✅ Implemented | Both use identical JSON keys |
| Data types match | 🔄 Needs work | Date/time format alignment, null handling |
| Values match (including derived calculations) | ❌ Missing | progress, status, days_unreading, median_day, finished_at, logs_count |
| Same JSON structure (nested objects/arrays) | ✅ Implemented | Same nesting |
| Logs limited to first 4 | ✅ Implemented | Both limit to 4 |
| Eager loading for performance | 🔄 Needs optimization | Single JOIN query instead of N+1 |
| Validation logic shared | 🔄 Needs implementation | Page ≤ total_page, status rules |

---

# Technical Decisions

### Decision 1: Derived Calculations Implementation

**Problem**: Rails uses ActiveModelSerializers with virtual methods for calculated fields like `progress`, `status`, `median_day`, etc. Go DTOs only include raw database fields.

**Decision**: Implement all derived calculations in Go using the same formulas as Rails model methods.

**Rationale**: 
- Ensures identical JSON output
- Enables gradual migration without client changes
- Maintains consistent business logic

**Implementation**:
- `progress`: `(page / total_page) * 100` rounded to 2 decimal places
- `status`: Conditional logic based on logs count, page count, days_unreading
- `days_unreading`: `(Date.today - last_log_or_started_at).days`
- `median_day`: `page / days_reading.round(2)`
- `finished_at`: Future date based on reading rate
- `logs_count`: `logs.size`

### Decision 2: Configuration Values

**Problem**: Rails `status` calculation depends on configuration values (`em_andamento_range`, `dormindo_range`) from `V1::UserConfig`.

**Decision**: Create Go config structure with default values matching Rails defaults (7 days running, 14 days sleeping).

**Rationale**:
- Ensures consistent status determination
- Maintainable configuration approach
- Prevents hardcoded magic numbers

### Decision 3: Date/Time Format Alignment

**Problem**: Date format inconsistencies between Rails and Go implementations.

**Decision**: Format all dates as RFC3339 strings (ISO 8601 with timezone) in JSON output using `time.Now()` matching Rails `Date.today` behavior.

**Rationale**:
- Modern standard for API responses
- Consistent with `logs.data` (TIMESTAMP) and `started_at` (DATE)
- Compatible with `pgx` PostgreSQL driver

### Decision 4: Database Query Optimization

**Problem**: Go implementation uses separate queries for each project (potential N+1).

**Decision**: Use single LEFT OUTER JOIN query with ordering matching Rails (`ORDER BY projects.id, logs.data DESC`).

**Rationale**:
- Matches Rails eager loading behavior
- Better performance under load
- Simpler code maintenance

**SQL**:
```sql
SELECT p.id, p.name, p.total_page, p.started_at, p.page, p.reinicia,
       l.id as log_id, l.data, l.start_page, l.end_page, l.note
FROM projects p
LEFT JOIN logs l ON p.id = l.project_id
ORDER BY p.id, l.data DESC
```

### Decision 5: Shared Validation Logic

**Problem**: Validation logic exists separately in Rails and Go, potentially leading to inconsistencies.

**Decision**: Define shared validation rules as database-level constraints + Go validation package.

**Rules**:
- `page ≤ total_page` (store as constraint)
- `start_page ≤ end_page` (validate in both apps)
- `logs_count ≥ 0` (derive from logs array)
- `status` values: `unstarted`, `finished`, `running`, `sleeping`, `stopped`

---

# Acceptance Criteria

## Functional

| AC | Description | Test |
|----|-------------|------|
| AC1 | `GET /v1/projects` returns identical JSON structure and values | Compare JSON with `jq` |
| AC2 | `GET /v1/projects/{project_id}` returns identical JSON structure and values | Compare JSON with `jq` |
| AC3 | `GET /v1/projects/{project_id}/logs` returns identical JSON structure and values | Compare JSON with `jq` |
| AC4 | Derived fields (progress, status, median_day, finished_at, logs_count, days_unreading) present in Go response | Check JSON output |
| AC5 | Date formats match (RFC3339 for timestamps, ISO date for started_at) | Verify JSON string format |
| AC6 | Null values handled identically (database NULL → JSON null) | Verify JSON output with NULL database values |
| AC7 | Logs limited to first 4 entries in correct order (by data DESC) | Compare array length and ordering |

## Non-Functional

| AC | Description | Threshold |
|----|-------------|-----------|
| NF1 | Performance comparable to Rails (same query pattern) | Max 10% difference |
| NF2 | Memory usage not significantly higher | < 20% increase |
| NF3 | Error responses match format | Same error JSON structure |

---

# Files to Modify

| File | Change Type | Reason |
|------|-------------|--------|
| `internal/domain/dto/project_response.go` | Modify | Add derived calculation fields (status, progress, days_unreading, median_day, finished_at, logs_count) |
| `internal/domain/dto/log_response.go` | Modify | Ensure data field format matches Rails (RFC3339) |
| `internal/adapter/postgres/project_repository.go` | Modify | Update query to calculate derived values or implement in handler |
| `internal/adapter/postgres/log_repository.go` | Modify | Confirm logs ordering matches Rails (data DESC) |
| `internal/repository/project_repository.go` | Modify | Update interface to include derived calculations |
| `internal/api/v1/handlers/projects_handler.go` | Modify | Add derived field calculation logic |
| `internal/api/v1/handlers/logs_handler.go` | Modify | Ensure logs ordering and formatting matches Rails |
| `internal/config/config.go` | Modify | Add app config fields (em_andamento_range, dormindo_range) |
| `internal/domain/models/project.go` | Create | Add calculation methods to match Rails model |

---

# Files Created

| File | Purpose |
|------|---------|
| `internal/validation/validation.go` | Shared validation rules |
| `internal/domain/models/calculations.go` | Derived calculation methods |
| `test/compare_responses.sh` | Script to compare JSON responses using curl + jq |

---

# Validation Rules

| Entity | Field | Rule | Go Implementation |
|--------|-------|------|-------------------|
| Project | page ≤ total_page | Constraint | Database constraint + validation |
| Project | status values | unstarted/finished/running/sleeping/stopped | Conditional logic |
| Project | progress range | 0.00-100.00 | Clamp function |
| Project | days_unreading | ≥ 0 | Validation function |
| Log | start_page ≤ end_page | Constraint | Validation function |
| Log | data timestamp | RFC3339 format | Format function |
| Log | note (optional) | String or null | NULL check |
| Log | text (optional) | String or null | NULL check |

---

# Out of Scope

| Item | Reason |
|------|--------|
| Authentication/Authorization endpoints | Different scope |
| POST/PUT/DELETE endpoints | Read-only comparison |
| Dashboard endpoints | Not part of current requirements |
| Pagination | Both implementations handle similarly |
| Database schema changes | Already defined |
| Rails serialization configuration changes | Use current behavior |

---

# Implementation Checklist

- [ ] **Phase 1: Field Alignment**
  - [ ] Verify field names match (snake_case)
  - [ ] Add null handling for optional dates
  - [ ] Ensure date format consistency (RFC3339)

- [ ] **Phase 2: Derived Calculations**
  - [ ] Implement `progress` calculation
  - [ ] Implement `status` determination
  - [ ] Implement `days_unreading` calculation
  - [ ] Implement `median_day` calculation
  - [ ] Implement `finished_at` calculation
  - [ ] Implement `logs_count` derivation

- [ ] **Phase 3: Query Optimization**
  - [ ] Implement single JOIN query for projects + logs
  - [ ] Add database indexes if needed
  - [ ] Verify query performance

- [ ] **Phase 4: Validation**
  - [ ] Create validation package
  - [ ] Implement shared validation rules
  - [ ] Add tests for edge cases

- [ ] **Phase 5: Testing & Validation**
  - [ ] Compare JSON responses using `jq`
  - [ ] Test edge cases (empty logs, zero total_page, etc.)
  - [ ] Performance comparison test
  - [ ] Database schema verification

---

# Stakeholder Alignment

| Stakeholder | Responsibilities | Verification |
|-------------|------------------|--------------|
| Product Owner | Approve acceptance criteria | Test results |
| Backend Lead | Implement Go changes | Code review |
| QA Engineer | Test JSON alignment | Comparison script |
| DevOps | Ensure same DB | Docker config |
| Client Team | Validate API compatibility | Integration tests |

---

# Traceability Matrix

| Requirement | User Story | AC | Test |
|-------------|------------|----|------|
| Identical JSON structure | As a client developer, I need consistent API responses | AC1, AC2, AC3 | JSON comparison |
| Derived calculations | As a client, I need calculated fields like progress and status | AC4 | JSON field check |
| Null handling | As a client, I expect null for absent values | AC6 | NULL DB values |
| Logs limit | As a client, I expect only first 4 logs | AC7 | Array length |

---

# Validation

## Code Quality
- [ ] All changes follow existing Go project conventions
- [ ] No breaking changes to public interfaces
- [ ] Proper error handling maintained
- [ ] Tests added for all new functionality

## Technical Feasibility
- [ ] All derived calculations implementable in Go
- [ ] Query optimization achievable with current database
- [ ] No performance regressions expected
- [ ] Date/time formatting compatible with PostgreSQL

## User Needs
- [ ] JSON response format preserved
- [ ] All existing features available
- [ ] No breaking changes for clients
- [ ] Consistent behavior between Rails and Go

---

# Ready for Implementation

✅ **Reviewed** - PRD has been refined and all technical questions resolved. Implementation can begin.

- All requirements unambiguous
- Acceptance criteria measurable and testable
- Technical approach validated
- Stakeholders aligned on scope
- Edge cases identified and addressed

**Next step**: Begin Phase 1 implementation (field alignment) in `internal/domain/dto/` packages.