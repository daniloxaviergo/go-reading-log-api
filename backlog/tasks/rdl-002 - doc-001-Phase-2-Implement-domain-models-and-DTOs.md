---
id: RDL-002
title: '[doc-001 Phase 2] Implement domain models and DTOs'
status: To Do
assignee:
  - workflow
created_date: '2026-04-01 00:57'
updated_date: '2026-04-01 01:25'
labels: []
dependencies: []
references:
  - 'PRD Section: Key Requirements'
  - 'Implementation Checklist: Core Components'
  - 'PRD Section: Files to Modify'
documentation:
  - doc-001
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement domain models for Project and Log entities in internal/domain/project.go and internal/domain/log.go.

Create response DTOs in internal/domain/dto/ for JSON serialization: project_response.go, log_response.go, and health_check_response.go.

Ensure all structs have appropriate JSON tags and embed context for data flow throughout the application.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Project and Log domain models implemented with all required fields
- [ ] #2 Response DTOs created with correct JSON tags for API compatibility
- [ ] #3 Context properly embedded in all models for request lifecycle
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Implement domain models and DTOs following Clean Architecture principles:

**Domain Models** (`internal/domain/models/`):
- Create `project.go` and `log.go` with Go structs mirroring Rails database schema
- Include all fields from Rails migrations
- Use appropriate Go types (int64 for IDs, time.Time for dates, pointers for nullable fields)
- Add JSON tags for API serialization compatibility

**Response DTOs** (`internal/domain/dto/`):
- Create `project_response.go` with fields matching Rails serializer
- Create `log_response.go` matching Rails log serializer
- Create `health_check_response.go` for health endpoints
- All structs must have JSON tags for exact API compatibility

**Context Embedding**:
- Embed `context.Context` in models for request lifecycle
- Follow existing patterns in the codebase

### 2. Files to Modify

**Create new files:**
- `internal/domain/models/project.go` - Project domain model
- `internal/domain/models/log.go` - Log domain model
- `internal/domain/dto/project_response.go` - Project JSON response DTO
- `internal/domain/dto/log_response.go` - Log JSON response DTO
- `internal/domain/dto/health_check_response.go` - Health check response DTO

**Reference files:**
- `rails-app/db/schema.rb` - Database schema for field definitions
- `rails-app/app/serializers/project_serializer.rb` - Project response structure
- `rails-app/app/serializers/log_serializer.rb` - Log response structure
- `go.mod` - Existing dependencies for context

### 3. Dependencies

**Prerequisites:**
- Go 1.25.7 (already in go.mod)
- `context` package (stdlib, no dependency needed)
- No new dependencies required for this task

**Blocking issues:** None - this is foundational domain work

### 4. Code Patterns

**Naming Conventions:**
- Struct names: PascalCase (Project, Log)
- Field names: PascalCase (matches Rails serializer output)
- File names: lowercase with underscore (project.go, log.go)

**Struct Field Mapping:**
- DB field | Go Type | Notes
- `id` | `int64` | Primary key
- `project_id` | `int64` | Foreign key
- `date` | `time.Time` | DATETIME field
- `integer` | `int` | Integer field
- `boolean` | `bool` | Boolean field
- `text` | `*string` | Nullable text (pointer)
- `date` (DB date) | `*string` | DATE type, use pointer for nullable

**Context Embedding:**
- Embed context for request lifecycle propagation
- Example: `ctx context.Context` as a field

**JSON Tags:**
- Use exact field names matching Rails serializer
- Example: `json:"total_page"` not `json:"totalPage"`

### 5. Testing Strategy

**Verification Steps:**
1. Verify structs compile and have correct JSON serialization
2. Test JSON marshaling/unmarshaling for each DTO
3. Verify field names match Rails serializer output exactly
4. Test null/nullable field handling with pointers

**Test Cases:**
- Project model with all fields populated
- Project model with nullable fields as nil
- Log model with various combinations
- JSON serialization matches Rails output format

**Tools:**
- Built-in `encoding/json` package
- No external test framework needed for basic verification

### 6. Risks and Considerations

**Field Type Decisions:**
- `date` column in Rails is DateTime - use `time.Time`
- `integer` columns maps to Go `int`
- All fields are nullable in DB except `created_at`/`updated_at` - use pointers
- `reinicia` is boolean in Rails - Go `bool`

**API Compatibility:**
- Must match Rails serializer output exactly
- `total_page` not `total_page` (underscore matches DB)
- `page` field exists in both Rails and Go model

**Timing:**
- This is Phase 2 work - depends on Phase 1 (setup) being done
- No blocking on other tasks - can run in parallel with repository implementations
<!-- SECTION:PLAN:END -->
