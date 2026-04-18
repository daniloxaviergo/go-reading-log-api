---
id: RDL-063
title: '[doc-005 Phase 2] Ensure median_day field is included in ProjectResponse DTO'
status: To Do
assignee:
  - thomas
created_date: '2026-04-18 11:47'
updated_date: '2026-04-18 13:13'
labels:
  - phase-2
  - median-day
  - dto
dependencies: []
references:
  - 'PRD Section: Key Requirements REQ-003'
  - internal/domain/dto/project_response.go
documentation:
  - doc-005
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update internal/domain/dto/project_response.go to ensure the median_day field is properly exposed in all project API responses. The field should be a float64 pointer that calculates pages per day rounded to 2 decimal places.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 median_day field present in ProjectResponse struct
- [ ] #2 Field serialized correctly to JSON with proper rounding
- [ ] #3 AC-REQ-003.1 verified: Inspect JSON response structure shows median_day
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
## Implementation Plan: Ensure median_day Field in ProjectResponse DTO

### 1. Technical Approach

**Analysis:** The `median_day` field is already defined in the `ProjectResponse` struct and calculated by the domain model, but it's not being populated in the response DTO when projects are fetched from the database.

**Root Cause:** In `internal/adapter/postgres/project_repository.go`, the `GetAllWithLogs()` and `GetWithLogs()` methods calculate derived fields (progress, status, days_unread, finished_at) but **do not call `CalculateMedianDay()`** to populate the `MedianDay` field in the response DTO.

**Solution:** Add the missing median_day calculation in both repository methods that build ProjectResponse objects.

**Architecture Decision:** Follow existing patterns:
- Use `project.CalculateMedianDay()` method (already implemented in domain model)
- Assign result directly to `projectResp.MedianDay`
- The method already handles all edge cases (no started_at, zero/negative days)

---

### 2. Files to Modify

| File | Change | Reason |
|------|--------|--------|
| `internal/adapter/postgres/project_repository.go` | Add `MedianDay` calculation in `GetWithLogs()` | Populate median_day for single project query |
| `internal/adapter/postgres/project_repository.go` | Add `MedianDay` calculation in `GetAllWithLogs()` | Populate median_day for all projects query |
| `internal/domain/dto/project_response_test.go` | Verify median_day serialization | Ensure JSON marshaling works correctly |

---

### 3. Dependencies

**Prerequisites (Already Satisfied):**
- ✅ `CalculateMedianDay()` method implemented in `internal/domain/models/project.go`
- ✅ `MedianDay` field exists in `ProjectResponse` struct with correct JSON tag
- ✅ Timezone configuration support added (RDL-061)
- ✅ Date parsing supports multiple formats (RDL-060)

**No External Dependencies Required**

---

### 4. Code Patterns

**Pattern to Follow:** Consistent with existing derived field calculations

```go
// Existing pattern in project_repository.go (lines ~195-200):
projectResp := dto.NewProjectResponse(
    project.ID,
    project.Name,
    startedAtStr,
    project.TotalPage,
    project.Page,
)
projectResp.LogsCount = logsCount
projectResp.Status = project.CalculateStatus(logsForProject, config.LoadConfig())
projectResp.DaysUnread = daysUnread
projectResp.Progress = project.CalculateProgress()
projectResp.FinishedAt = formatTimePtr(project.CalculateFinishedAt(logsForProject))

// NEW: Add this line
projectResp.MedianDay = project.CalculateMedianDay()
```

**Field Type:** `*float64` (pointer to float64)
- Allows nil serialization when calculation is undefined
- JSON tag: `"median_day,omitempty"` in struct definition

---

### 5. Testing Strategy

**Unit Tests (Verify Calculation Logic):**
1. Test `CalculateMedianDay()` with normal case (page > 0, started_at exists)
2. Test edge cases:
   - No `started_at` → returns 0.0
   - Zero/negative days reading → returns 0.0
   - Page exceeds total → calculates based on actual page

**Integration Tests (Verify Repository):**
1. Create test project with known values
2. Fetch via `GetWithLogs()` and verify `MedianDay` populated
3. Fetch all via `GetAllWithLogs()` and verify `MedianDay` populated
4. Verify JSON serialization includes `median_day`

**Test Commands:**
```bash
# Run unit tests for median day calculation
go test -v ./internal/domain/models/... -run TestProject_CalculateMedianDay

# Run repository integration tests
go test -v ./internal/adapter/postgres/... -run TestProjectRepository

# Run all tests with coverage
go test -cover ./...
```

---

### 6. Implementation Steps

**Step 1: Update `GetWithLogs()` method**
- Location: `internal/adapter/postgres/project_repository.go`, line ~200
- Add: `projectResp.MedianDay = project.CalculateMedianDay()`
- Verify: Test single project endpoint `/v1/projects/{id}.json`

**Step 2: Update `GetAllWithLogs()` method**
- Location: `internal/adapter/postgres/project_repository.go`, line ~400
- Add: `projectResp.MedianDay = project.CalculateMedianDay()`
- Verify: Test all projects endpoint `/v1/projects.json`

**Step 3: Verify JSON Serialization**
- Run tests to ensure `median_day` appears in JSON output
- Check that `omitempty` tag works correctly for nil values

**Step 4: Run Full Test Suite**
```bash
go test -v ./...
go fmt ./...
go vet ./...
```

---

### 7. Risks and Considerations

**Risk:** None identified - this is a straightforward addition following existing patterns.

**Considerations:**
1. **Edge Case Handling:** `CalculateMedianDay()` already handles all edge cases (no started_at, zero days, etc.)
2. **Timezone Consistency:** Method uses context-based timezone from RDL-061
3. **JSON Tag:** Field uses `omitempty` - will be omitted if nil (acceptable behavior)
4. **Performance:** No additional queries needed - calculation is in-memory

---

### 8. Verification Checklist

After implementation, verify:
- [ ] `median_day` field present in `/v1/projects.json` response
- [ ] `median_day` field present in `/v1/projects/{id}.json` response
- [ ] Field value correctly calculated (pages / days_reading, rounded to 2 decimals)
- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] `go fmt` and `go vet` pass with no errors
- [ ] Follows Clean Architecture (domain model calculates, repository populates DTO)

---

### 9. Expected JSON Output

**Before:**
```json
{
  "id": 1,
  "name": "Test Project",
  "total_page": 200,
  "page": 50,
  "started_at": "2024-01-15T10:30:00Z",
  "progress": 25.0,
  "status": "running",
  "logs_count": 4,
  "days_unreading": 5,
  "finished_at": "2024-02-15T00:00:00Z"
  // median_day MISSING
}
```

**After:**
```json
{
  "id": 1,
  "name": "Test Project",
  "total_page": 200,
  "page": 50,
  "started_at": "2024-01-15T10:30:00Z",
  "progress": 25.0,
  "status": "running",
  "logs_count": 4,
  "days_unreading": 5,
  "median_day": 10.0,     // ← NOW PRESENT
  "finished_at": "2024-02-15T00:00:00Z"
}
```

---

### 10. Related Documentation

- PRD Section: REQ-003 - Add median_day field to ProjectResponse DTO
- Acceptance Criteria: AC-REQ-003.1 - Inspect JSON response structure shows median_day
- Implementation matches Rails behavior: `median_day` = `page / days_reading` rounded to 2 decimals
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass use testing-expert subagent for test execution and verification
- [ ] #2 All integration tests pass use testing-expert subagent for test execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Database queries optimized with proper indexes
- [ ] #8 Documentation updated in QWEN.md
- [ ] #9 New code paths include error path tests
- [ ] #10 HTTP handlers test both success and error responses
- [ ] #11 Integration tests verify actual database interactions
- [ ] #12 Tests use testing-expert subagent for test execution and verification
- [ ] #13 Field uses float64 pointer type for optional serialization
- [ ] #14 Rounding to 2 decimal places implemented
<!-- DOD:END -->
