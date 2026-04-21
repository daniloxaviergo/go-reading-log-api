---
id: doc-007
title: Logs Endpoint Alignment PRD - RDL-071
type: other
created_date: '2026-04-21 12:03'
---


# Logs Endpoint Alignment PRD

## Executive Summary

**Why necessary:** Provides a 10-second summary for stakeholders (managers, product owners) to quickly understand the feature, scope, and value without reading the full document.

This PRD addresses critical structural discrepancies identified between the Go API and Rails API for the `v1/projects/{id}/logs.json` endpoint. The primary goal is to align the Go API implementation with the JSON:API specification followed by the Rails API to ensure consistency, improve performance, and reduce client-side parsing complexity.

**Scope:** Refinement of existing Go API response structure for logs endpoint.
**Key Changes:** Date format standardization (RFC3339), replacement of embedded project objects with relationship references, and ID type alignment.
**Impact:** Medium-High. Affects all downstream clients consuming the logs endpoint.

---

## Key Requirements

**Why necessary:** Captures the core features in a scannable format with status tracking. Shows progress and what's been decided vs. what's pending.

| ID | Requirement | Description | Priority | Status |
|:---:|-------------|-------------|:--------:|:------:|
| REQ-01 | Date Format Standardization | Convert all datetime fields (`data`) to RFC3339 (ISO 8601) format including timezone. | Critical | To Do |
| REQ-02 | Relationship Reference Implementation | Replace embedded `project` object with JSON:API relationship reference (`relationships.project.data`). | Critical | To Do |
| REQ-03 | ID Type Alignment | Change top-level and nested IDs from integer to string type for JSON:API compliance. | High | To Do |
| REQ-04 | Response Size Optimization | Reduce payload size by removing data duplication (denormalization) in log entries. | Medium | To Do |
| REQ-05 | Field Naming Consistency | Maintain snake_case for Go API fields (internal standard) while ensuring JSON serialization matches spec. | Low | To Do |

---

## Technical Decisions

**Why necessary:** Documents the "why" behind architecture choices. Crucial for onboarding developers, avoiding repeated design debates, and understanding trade-offs made.

### Decision 1: Adoption of JSON:API Specification
**Context:** The Rails API implements a strict JSON:API 1.0 specification, while the Go API uses a custom structure.
**Choice:** The Go API will be refactored to comply with JSON:API 1.0 specifications for the logs endpoint.
**Rationale:** 
- Ensures interoperability between the two services.
- Provides a standardized contract for frontend clients.
- Reduces maintenance overhead from maintaining two different formats.
**Validation:** Confirmed via `subagent "palha"` technical review.

### Decision 2: RFC3339 Datetime Format
**Context:** Go API uses `2026-04-02 21:21:53` (custom string), Rails uses `2026-04-02T18:21:53.000-03:00` (RFC3339).
**Choice:** All `data` fields in the Go API will be changed from `string` to `time.Time` type, which marshals to RFC3339 by default.
**Rationale:** 
- Timezone information is critical for distributed systems.
- Standard libraries handle parsing effortlessly across languages.
- Eliminates custom parsing logic required by clients.
**Trade-off:** Slight increase in payload size due to timezone info, offset by network efficiency gains.

### Decision 3: Relationship References vs. Embedded Objects
**Context:** Go API embeds full project objects in every log entry (~350 bytes), Rails uses references (~180 bytes + included array).
**Choice:** Implement relationship references (`relationships.project.data`) and utilize the `included` array for project data, similar to Rails implementation.
**Rationale:** 
- **Performance:** ~50% reduction in response size (estimated 170 bytes saved per log entry).
- **Data Integrity:** Prevents stale data if project details change.
- **Spec Compliance:** JSON:API spec discourages denormalization.
**Trade-off:** Slightly more complex client-side logic to resolve relationships (though modern clients handle this easily).

### Decision 4: ID Type Conversion
**Context:** Go API uses `int` (9092), Rails API uses `string` ("9092").
**Choice:** IDs will be serialized as strings in the JSON response to match JSON:API spec, though internal Go structs may retain `int64`.
**Rationale:** 
- JSON:API specification mandates string IDs.
- Prevents issues with large integers exceeding JavaScript MAX_SAFE_INTEGER.
- Ensures consistency across the API ecosystem.

---

## Acceptance Criteria

**Why necessary:** Defines objective, testable conditions that mark completion. Separated into Functional and Non-Functional.

### Functional Acceptance Criteria
1.  **AC-FUNC-01:** A GET request to `/v1/projects/{id}/logs.json` returns a valid JSON:API document structure with `data` and `included` arrays.
2.  **AC-FUNC-02:** The `data` field within log attributes is serialized as an RFC3339 string (e.g., `2026-04-02T18:21:53.000Z`).
3.  **AC-FUNC-03:** Log entries do not contain an embedded `project` object in `attributes`; instead, they contain a `relationships.project.data` object with `id` and `type`.
4.  **AC-FUNC-04:** The top-level `id` field in log entries is serialized as a string.
5.  **AC-FUNC-05:** Project data appears in the `included` array when requested or by default, matching Rails API behavior.

### Non-Functional Acceptance Criteria
1.  **AC-NFUNC-01:** Response time for the logs endpoint does not degrade by more than 10% compared to the current implementation.
2.  **AC-NFUNC-02:** The payload size is reduced by at least 40% compared to the current embedded implementation.
3.  **AC-NFUNC-03:** All existing tests pass without modification (backward compatibility layer maintained if needed, though full migration is preferred).
4.  **AC-NFUNC-04:** Error responses follow the JSON:API error format specification.

---

## Files to Modify

**Why necessary:** Actionable implementation checklist for developers. Shows exactly where changes are needed and the rationale for each file.

| File Path | Change Description | Rationale |
|:----------|:-------------------|:----------|
| `internal/domain/dto/log_response.go` | Update `LogResponse` struct to use `time.Time` for `Data` field and add `Relationships` struct. Remove `Project` field from attributes. | Core DTO definition for API response structure. |
| `internal/domain/dto/log_response.go` | Add `Included` slice to hold related project resources. | Required for JSON:API `included` array. |
| `internal/api/v1/handlers/logs_handler.go` | Update `GetProjectLogs` handler to populate `relationships` and `included` arrays correctly. | Logic layer implementing the new structure. |
| `internal/adapter/postgres/log_repository.go` | Ensure queries fetch project IDs efficiently for relationship building. | Performance optimization for relationship assembly. |
| `internal/api/v1/routes.go` | Verify route configuration supports the new response format (no changes expected, but verification needed). | Safety check. |

---

## Files Created

**Why necessary:** Documents artifacts generated for this PRD (user-facing docs vs. technical docs).

| File Path | Purpose |
|:----------|:--------|
| `docs/api-changes/logs-endpoint-refinement.md` | Detailed changelog for API consumers describing the breaking changes and migration steps. |
| `test/integration/logs_endpoint_test.go` | Updated integration tests covering the new JSON:API structure. |

---

## Validation Rules

**Why necessary:** Ensures data integrity and consistent user experience. Shows that validation logic is well-defined, shared between TUI and CLI (DRY principle), and consistent with error messages.

### Input Validation (for POST/PUT if extended later)
| Field | Rule | Error Message |
|:------|:-----|:--------------|
| `data` | Must be valid RFC3339 string | "Invalid date format. Expected RFC3339 (ISO 8601)." |
| `start_page` | Must be integer >= 0 | "Start page must be a non-negative integer." |
| `end_page` | Must be integer > start_page | "End page must be greater than start page." |

### Output Validation
| Field | Rule | Check Method |
|:------|:-----|:-------------|
| `data` | RFC3339 format | Regex: `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(.\d+)?(Z|[+-]\d{2}:\d{2})` |
| `id` | String type | Type assertion in test suite |
| `relationships.project.data.id` | Present and non-empty | Schema validation |
| `included` | Array of resources | Schema validation |

### Error Response Format
All errors must adhere to the JSON:API error object structure:
```json
{
  "errors": [
    {
      "status": "400",
      "title": "Bad Request",
      "detail": "Validation failed for field 'data'",
      "source": { "pointer": "/data/attributes/data" }
    }
  ]
}
```

---

## Out of Scope

**Why necessary:** Explicitly prevents scope creep. Manages expectations about what the MVP does not include.

- **Database Schema Changes:** No migration scripts or DDL changes are required for this PRD. The existing schema supports the required data extraction.
- **Rails API Modifications:** This PRD focuses on aligning the Go API to match Rails, not vice versa.
- **Authentication/Authorization:** The logs endpoint remains public/unauthenticated as per current specs.
- **Pagination Implementation:** While related to JSON:API, full pagination links (`links.next`, `links.prev`) are out of scope for this specific refinement but noted for future work.
- **Filtering/Sorting:** Advanced query parameters (e.g., `?sort=data`) are not included.
- **Webhook Subscriptions:** Not part of the current endpoint functionality.

---

## Implementation Checklist

**Why necessary:** Tracks implementation progress. Helps developers know what steps remain.

- [ ] **Phase 1: Refactoring**
    - [ ] Update `LogResponse` DTO to use `time.Time` for `Data`.
    - [ ] Add `Relationships` struct to `LogResponse`.
    - [ ] Remove `Project` field from `LogAttributes`.
    - [ ] Update `GetProjectLogs` query logic to fetch project IDs.
- [ ] **Phase 2: Serialization**
    - [ ] Implement JSON marshaling for relationships (`relationships.project.data`).
    - [ ] Implement `included` array population logic.
    - [ ] Ensure ID serialization converts `int64` to `string`.
- [ ] **Phase 3: Testing**
    - [ ] Write unit tests for new DTO structure.
    - [ ] Update integration tests to validate JSON:API schema compliance.
    - [ ] Run performance benchmarks to verify response size reduction.
- [ ] **Phase 4: Documentation**
    - [ ] Update API documentation (`docs/api-changes/...`).
    - [ ] Create migration guide for clients (JavaScript/Python examples).
    - [ ] Update Postman collections if applicable.
- [ ] **Phase 5: Deployment**
    - [ ] Deploy to staging environment.
    - [ ] Run smoke tests against both Go and Rails endpoints.
    - [ ] Monitor logs for parsing errors on client side (if possible).

---

## Stakeholder Alignment

**Why necessary:** Ensures all parties understand their responsibilities.

| Stakeholder | Responsibility | Sign-off Required |
|:------------|:---------------|:------------------|
| **Backend Team (Go)** | Implementation of JSON:API structure, performance optimization. | Yes |
| **Backend Team (Rails)** | Review alignment with existing Rails implementation. | Yes |
| **Frontend Team** | Verify compatibility with existing clients; provide feedback on `included` array usage. | Yes |
| **Product Owner** | Prioritization of this task vs. other backlog items. | Yes |
| **DevOps** | Deployment strategy and rollback plan verification. | Yes |

### Acceptance Verification
- **Owner:** Backend Team Lead
- **Reviewer:** Frontend Architect
- **Criteria:** All Acceptance Criteria (Section 4) must be met.

---

## Traceability Matrix

**Why necessary:** Links requirements → user stories → acceptance criteria → tests. Enables impact analysis for changes, test coverage verification, and requirement completeness validation.

| Req ID | Source Issue | User Story | AC ID | Test File |
|:-------|:-------------|:-----------|:-----:|:----------|
| REQ-01 | Issue #1 (Date) | As a client developer, I want consistent date formats so I can parse them easily. | AC-FUNC-02, AC-NFUNC-02 | `logs_endpoint_test.go` |
| REQ-02 | Issue #2 (Embedding) | As a mobile user, I want smaller payloads to save bandwidth. | AC-FUNC-03, AC-NFUNC-02 | `logs_endpoint_test.go` |
| REQ-03 | Issue #3 (ID Type) | As a frontend dev, I want standard ID types to avoid casting errors. | AC-FUNC-04 | `logs_endpoint_test.go` |
| REQ-04 | Issue #2 (Performance) | As a system admin, I want faster response times. | AC-NFUNC-01 | `benchmark_test.go` |

### Impact Analysis
- **High Impact:** Changes affect all clients consuming the logs endpoint.
- **Dependency:** None. This is an internal structural change (assuming JSON:API was planned anyway).
- **Risk:** Medium. Breaking change for clients relying on embedded objects.

---

## Validation

**Why necessary:** Proof the PRD is sound before implementation begins. Confirms code quality standards, technical feasibility, and alignment with user needs.

### Code Quality Standards Checklist
- [ ] Follows Go 1.25.7 best practices (`go fmt`, `go vet`).
- [ ] No circular dependencies introduced.
- [ ] Context timeouts properly handled (5s limit maintained).
- [ ] Error wrapping uses `%w` format.

### Technical Feasibility Check
- [x] **Date RFC3339:** Supported natively by `time.Time` in Go 1.25.
- [x] **JSON:API Relationships:** Standard library `encoding/json` supports custom marshaling via `MarshalJSON()`.
- [x] **Response Size:** Measurable and verifiable via benchmarking.

### User Need Alignment
- [x] Resolves critical discrepancy reported in RDL-071.
- [x] Aligns with existing Rails API contract (interoperability).
- [x] Addresses performance concerns raised in comparison report.

---

## Ready for Implementation

**Why necessary:** Final gatekeeping statement. Indicates the PRD has been reviewed and is unambiguous enough for developers to begin work.

✅ **APPROVED FOR DEVELOPMENT**

This PRD is ready for implementation pending stakeholder sign-off on Section 8 (Stakeholder Alignment). The technical constraints are well-understood, and the required changes are localized to the `internal/domain/dto` and `internal/api/v1/handlers` packages.

**Next Step:** Assign developer resources to Phase 1 of the Implementation Checklist.