---
id: RDL-075
title: '[doc-007 Phase 3] Update integration tests for JSON:API schema compliance'
status: To Do
assignee:
  - catarina
created_date: '2026-04-21 12:11'
updated_date: '2026-04-21 13:35'
labels:
  - testing
  - backend
dependencies: []
references:
  - AC-FUNC-01
  - AC-NFUNC-02
documentation:
  - doc-007
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update the test/integration/logs_endpoint_test.go file to validate the new JSON:API response structure, including checks for RFC3339 date format, relationship references, and payload size reduction.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Tests validate JSON:API schema
- [ ] #2 Date format checked for RFC3339
- [ ] #3 Payload size verified
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The goal is to update integration tests for the logs endpoint to validate JSON:API schema compliance. The Go API has already been updated to return JSON:API formatted responses (based on code review of `log_response.go`, `jsonapi_response.go`, and `logs_handler.go`).

**Key Observations from Code Review:**
1.  **DTO Structure:** `LogResponse` already uses `time.Time` for the `Data` field (RFC3339 compliant) and has a `Relationships` struct with `Project` reference.
2.  **Handler Logic:** `logs_handler.go` constructs JSON:API envelopes using `dto.JSONAPIData`, `NewJSONAPIEnvelopeWithIncluded`, and populates `relationships.project.data` with string IDs.
3.  **Current Test Gap:** The existing test file (`logs_integration_test.go`) parses responses but lacks specific assertions for the JSON:API schema elements like `relationships`, `included`, and strict RFC3339 format validation.

**Approach:**
1.  Enhance `test_context.go` to include robust JSON:API parsing helpers that can extract `relationships` and `included` data.
2.  Update `logs_integration_test.go` to add specific test cases for:
    *   Schema structure validation (presence of `data`, `relationships`, `included`).
    *   RFC3339 date format verification using regex or time parsing.
    *   Payload size verification (comparing against expected baseline).
    *   Relationship reference correctness (ensuring project is a reference, not embedded in attributes).
3.  Ensure tests handle the new envelope structure (`{"data": {...}, "included": [...]}`) correctly.

**Why this approach:** The backend logic is already implemented to output JSON:API. The test suite needs to catch up to validate that the output matches the specification (RDL-071/Doc-007). This ensures the Phase 3 goal of "JSON:API schema compliance" is met before deployment.

---

### 2. Files to Modify

| File | Action | Reason |
| :--- | :--- | :--- |
| `test/integration/test_context.go` | **Modify** | Add helper methods to parse JSON:API envelopes, extract relationships, and verify included resources. Update `ParseLogResponseArray` to handle the new envelope structure properly. |
| `test/integration/logs_integration_test.go` | **Modify** | Add new test cases specifically for JSON:API schema compliance (relationships, date format, payload size). Refactor existing tests to use the new parsing helpers. |
| `docs/api-changes/logs-endpoint-refinement.md` | **Create/Verify** | Ensure documentation exists for the changes being tested (part of PRD requirements). |

---

### 3. Dependencies

*   **Existing Implementation:** Requires the JSON:API response structure to be fully implemented in `internal/domain/dto/` and `internal/api/v1/handlers/`. (Verified as complete based on code review).
*   **Test Database:** Requires a running PostgreSQL instance with the `reading_log_test` schema.
*   **Go Version:** Go 1.25.7+ for native time.RFC3339 support.

---

### 4. Code Patterns

*   **JSON:API Schema Validation:** Use `encoding/json` to unmarshal into `dto.JSONAPIEnvelope`. Verify `data` is an array of objects containing `type`, `id`, `attributes`, and `relationships`.
*   **RFC3339 Verification:** Parse the `data` string field using `time.Parse(time.RFC3339, ...)`. If it parses successfully, the format is correct.
*   **Payload Size:** Calculate response body length in bytes. Compare against a calculated baseline (e.g., previous embedded object size).
*   **Relationship Checking:** Assert that `attributes.project` does NOT exist (it should be in `relationships.project.data`).
*   **Consistent Error Handling:** Use `t.Fatalf` for fatal setup errors and `t.Errorf` for assertion failures.

---

### 5. Testing Strategy

*   **Unit Tests:** Not required for this task (focus is integration).
*   **Integration Tests:**
    *   **Test 1 (Schema Structure):** Verify response contains `data`, `included`, and `relationships.project.data`.
    *   **Test 2 (RFC3339 Date):** Parse the `data` field string to ensure it matches RFC3339 format.
    *   **Test 3 (Payload Size):** Measure response size and assert it is below a certain threshold (or reduced by X%).
    *   **Test 4 (Relationship Reference):** Confirm project data is in `included` array, not embedded in log attributes.
*   **Edge Cases:** Empty logs list, single log, multiple logs (limit check).

---

### 6. Risks and Considerations

*   **Breaking Changes:** The response format changes from flat JSON to JSON:API envelope. Clients consuming the old format will break. This is intentional per PRD.
*   **Parsing Complexity:** The test helpers need to be robust enough to handle both single-object and array responses within the envelope.
*   **Performance:** Adding relationship resolution in tests might slow down the suite slightly, but it's necessary for accuracy.
*   **Documentation Sync:** Ensure `docs/api-changes/logs-endpoint-refinement.md` is updated to reflect the exact schema tested here.
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
- [ ] #11 100% coverage for modified files
<!-- DOD:END -->
