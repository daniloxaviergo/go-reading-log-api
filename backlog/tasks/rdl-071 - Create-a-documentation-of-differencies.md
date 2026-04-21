---
id: RDL-071
title: Create a documentation of differencies
status: To Do
assignee:
  - catarina
created_date: '2026-04-21 10:35'
updated_date: '2026-04-21 10:58'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create a documentation of differencies of jsons and write in docs/diff_show_logs.md
http://0.0.0.0:3001/v1/projects/450/logs.json -> Rails-Api
http://0.0.0.0:3000/v1/projects/450/logs.json -> Go-Api

Dont change rails-app
The fix should be in golang
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task focuses on **documentation** of differences between the Go API and Rails API responses for the logs endpoint (`/v1/projects/{project_id}/logs.json`).

**Approach:**
1. **Analyze existing code**: Review both Go and Rails implementation to identify structural and data differences
2. **Document field-by-field comparisons**: Create detailed mapping of all response fields
3. **Document calculation differences**: Capture any discrepancies in derived fields
4. **Create comparison examples**: Provide side-by-side JSON examples for clarity
5. **Write migration guidance**: Help clients understand how to handle differences

**Why this approach:**
- The task explicitly asks for documentation, not code changes
- Comparing implementations ensures accurate documentation
- Providing migration guidance adds practical value for API consumers

---

### 2. Files to Modify

#### New Files to Create:
| File | Purpose | Lines |
|------|---------|-------|
| `docs/diff_show_logs.md` | Main documentation of logs endpoint differences | ~500-800 |

#### Files to Reference (Read-Only):
| File | Purpose |
|------|---------|
| `internal/api/v1/handlers/logs_handler.go` | Go API logs implementation |
| `rails-app/app/controllers/v1/logs_controller.rb` | Rails API logs implementation |
| `docs/diff_show_project.md` | Existing project comparison (format reference) |
| `docs/api-response-alignment.md` | General API alignment documentation |

#### Files to Check for Context:
| File | Purpose |
|------|---------|
| `internal/domain/dto/project_response.go` | Go DTO definitions |
| `internal/domain/dto/log_response.go` | Log response DTO |
| `rails-app/app/serializers/*` | Rails serializer configurations |

---

### 3. Dependencies

**Prerequisites:**
- ✅ Go API logs endpoint implementation (RDL-047, RDL-057)
- ✅ Rails API logs endpoint exists
- ✅ Existing comparison documentation patterns (RDL-039, RDL-059)

**No blocking issues** - This is a documentation task that can proceed independently.

---

### 4. Code Patterns

**Documentation Style:**
```markdown
# Logs Endpoint Comparison Report

## Overview
Comparing Go API vs Rails API responses for endpoint: `v1/projects/{id}/logs.json`

## Field Comparison Table
| Field | Go Value | Rails Value | Match | Notes |
|-------|----------|-------------|-------|-------|
| id | 9092 (int) | "9092" (string) | ⚠️ Type | ID format differs |

## Code Examples
Provide Go and Rails code snippets showing how each API constructs the response.

## Migration Guide
Highlight what clients need to change when migrating from Rails to Go API.
```

**Key Patterns to Document:**
1. **ID Type**: Integer vs String (JSON:API spec)
2. **Date Format**: RFC3339 vs ISO 8601 vs Custom datetime
3. **Nested Objects**: Go embeds project in logs; Rails does not
4. **Field Naming**: snake_case vs kebab-case

---

### 5. Testing Strategy

**Documentation Verification:**
- [ ] Compare actual API responses from both endpoints
- [ ] Verify field mappings are accurate
- [ ] Validate example JSON is correct
- [ ] Ensure migration guidance is actionable

**No unit/integration tests required** - This is a documentation task.

---

### 6. Risks and Considerations

**Known Issues to Document:**

1. **ID Type Difference (CRITICAL)**
   - Go: `9092` (integer)
   - Rails: `"9092"` (string, per JSON:API spec)
   - Impact: Client code must handle string parsing

2. **Date Format Inconsistency (HIGH)**
   - Go logs: `2026-04-18 21:21:53` (custom format)
   - Rails logs: `2026-04-02T18:21:53.000-03:00` (ISO 8601)
   - Impact: Date parsing requires multi-format support

3. **Embedded Project Object (MEDIUM)**
   - Go: Includes full project object in each log
   - Rails: No embedded project (follows JSON:API relationships)
   - Impact: Response size and data structure differs

4. **Field Naming (LOW)**
   - Go: `start_page`, `end_page`
   - Rails: `start-page`, `end-page`
   - Impact: Client field access requires mapping

**Recommendations:**
- Update Go API to use JSON:API compliant date format (RFC3339)
- Consider removing embedded project object from logs (use relationship references)
- Document field naming convention clearly for clients
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
<!-- DOD:END -->
