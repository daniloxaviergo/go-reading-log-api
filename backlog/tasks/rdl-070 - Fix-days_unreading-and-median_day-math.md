---
id: RDL-070
title: Fix days_unreading and median_day math
status: To Do
assignee:
  - catarina
created_date: '2026-04-21 10:15'
updated_date: '2026-04-21 10:22'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
In route the fields days_unreading and median_day is diferent
http://0.0.0.0:3001/v1/projects/450.json -> Rails-Api
http://0.0.0.0:3000/v1/projects/450.json -> Go-Api

| Field | Go-Api | Rails-Api |
| days-unreading | 19 | 15 |
| median-day | 11.33 | 12.12 |

Dont change the rails-app
Look the code rais-app to check the math and fix
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan: Fix days_unreading and median_day Math

## 1. Technical Approach

### Root Cause Analysis

After comparing Go implementation with Rails API code, I've identified the **actual bug**:

**Rails Code:**
```ruby
def days_unreading
  base_data = last_read && last_read[:data].to_date || started_at
  @days_unreading = (Date.today - base_data).to_i
end

def median_day
  (page.to_f / days_reading.to_f).round(2)
end

def days_reading
  (Date.today - started_at).to_i
end
```

**Key:** Rails `Date.today` returns the current date **in the application timezone** (America/Sao_Paulo/BRT).

**Current Go Code (BUGGY):**
```go
now := time.Now()  // Gets current time in SYSTEM timezone (likely UTC)
nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tzLocation)
```

**The Problem:**
- `time.Now()` returns time in system timezone (UTC)
- Extracting year/month/day from UTC time
- Then applying BRT timezone to those extracted values

**Example of the bug:**
- Current UTC time: `2026-04-21 02:00:00`
- Extracted date parts: Year=2026, Month=4, Day=21
- After applying BRT (-3 hours): Date becomes `2026-04-20`!
- This causes a **1-day discrepancy**

**The Fix:**
Convert to the target timezone FIRST, THEN extract date parts:
```go
now := time.Now().In(tzLocation)  // Convert to target timezone first
nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tzLocation)
```

### Files to Modify

| File | Changes |
|------|---------|
| `internal/domain/models/project.go` | Fix `CalculateDaysUnreading()` - use `.In(tzLocation)` before extracting date |
| `internal/domain/models/project.go` | Fix `CalculateMedianDay()` - use `.In(tzLocation)` before extracting date |
| `internal/domain/models/project.go` | Fix `CalculateFinishedAt()` - use `.In(tzLocation)` before extracting date |

### Code Changes

**Before (buggy):**
```go
func (p *Project) CalculateDaysUnreading(logs []*dto.LogResponse) *int {
    // ...
    now := time.Now()  // ❌ Gets UTC time
    
    ctx := p.GetContext()
    tzLocation := getTimezoneFromContext(ctx)
    
    nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tzLocation)
    // ...
}
```

**After (fixed):**
```go
func (p *Project) CalculateDaysUnreading(logs []*dto.LogResponse) *int {
    // ...
    ctx := p.GetContext()
    tzLocation := getTimezoneFromContext(ctx)
    
    now := time.Now().In(tzLocation)  // ✅ Convert to target timezone first
    nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tzLocation)
    // ...
}
```

### Testing Strategy

**Unit Tests:**
```bash
# Test days_unreading calculation with timezone
go test -v ./internal/domain/models/... -run TestProject_CalculateDaysUnreading_Timezone

# Test median_day calculation with timezone
go test -v ./internal/domain/models/... -run TestProject_CalculateMedianDay_Timezone

# Run all project model tests
go test -v ./internal/domain/models/...
```

**Integration Tests:**
```bash
# Run all integration tests
go test -v ./test/integration/...

# Run repository tests
go test -v ./internal/adapter/postgres/...
```

**Manual Verification:**
1. Start the server with BRT timezone configured
2. Query `/v1/projects/450.json`
3. Compare `days_unreading` value with Rails API
4. Values should match (within 1 day tolerance)

### Expected Results

After fix:
- `days_unreading`: Should match Rails API (15 days)
- `median_day`: Should match Rails API (12.12)
- `finished_at`: Should be calculated correctly

### Risks and Considerations

**Risk:** Minimal - this is a straightforward timezone fix following established patterns.

**Consideration:** The fix ensures all date calculations use the application's configured timezone consistently, matching Rails `Date.today` behavior exactly.

### Acceptance Criteria Alignment

| AC | Status | Verification |
|----|--------|--------------|
| All unit tests pass | To Do | Run `go test ./...` |
| All integration tests pass | To Do | Run `go test -v ./test/integration/...` |
| go fmt and go vet pass | To Do | Run linters |
| days_unreading matches Rails | To Do | Verify value matches Rails API |
| median_day calculation correct | To Do | Verify value matches Rails API |

### Summary

This task fixes a **timezone conversion bug** where Go was extracting date parts from UTC time before applying the target timezone, causing date shifts. The fix ensures `time.Now()` is converted to the target timezone FIRST, matching Rails `Date.today` behavior exactly.
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 go fmt and go vet pass with no errors
- [ ] #2 Clean Architecture layers properly followed
- [ ] #3 Error responses consistent with existing patterns
- [ ] #4 HTTP status codes correct for response type
- [ ] #5 Documentation updated in QWEN.md
- [ ] #6 New code paths include error path tests
- [ ] #7 HTTP handlers test both success and error responses
- [ ] #8 Integration tests verify actual database interactions
- [ ] #9 Tests use testing-expert subagent for test execution and verification
- [ ] #10 All unit tests pass
- [ ] #11 All integration tests pass
<!-- DOD:END -->
