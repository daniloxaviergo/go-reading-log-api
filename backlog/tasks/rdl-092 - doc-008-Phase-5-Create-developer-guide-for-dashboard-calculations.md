---
id: RDL-092
title: '[doc-008 Phase 5] Create developer guide for dashboard calculations'
status: To Do
assignee:
  - catarina
created_date: '2026-04-21 15:52'
updated_date: '2026-04-22 17:10'
labels:
  - phase-5
  - documentation
  - guide
dependencies: []
references:
  - DOC-002
  - Implementation Checklist Phase 5
documentation:
  - doc-008
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create comprehensive developer guide explaining calculation methodologies, configuration options, and troubleshooting procedures. Include explanation of mean_day, progress_geral, fault percentage, and speculative mean calculations.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Calculation methodologies documented
- [ ] #2 Configuration options explained
- [ ] #3 Troubleshooting guide created
- [ ] #4 Developer onboarding information included
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task is a **documentation-only** work item (no code changes required). The goal is to create a comprehensive developer guide explaining dashboard calculation methodologies, configuration options, and troubleshooting procedures.

**Approach:**
- Create a new documentation file `docs/dashboard-developer-guide.md`
- Structure the guide following the existing documentation patterns in `docs/README.go-project.md`
- Include detailed explanations of all 8 dashboard endpoints and their calculations
- Provide code examples for common tasks
- Include troubleshooting section for common issues

**Why this approach:**
- The PRD (doc-008) already defines the API specifications
- Existing documentation in `docs/dashboard-api-reference.md` covers endpoint usage
- This guide should complement existing docs with deeper technical explanations
- No code changes needed - pure documentation effort

---

### 2. Files to Modify

#### New Files to Create

| File Path | Purpose |
|-----------|---------|
| `docs/dashboard-developer-guide.md` | Main developer guide for dashboard calculations |

#### Existing Files to Reference (No Modification Needed)

| File Path | Purpose |
|-----------|---------|
| `docs/dashboard-api-reference.md` | API endpoint reference (existing) |
| `docs/README.go-project.md` | Project structure documentation (existing) |
| `internal/service/dashboard/*.go` | Service implementations (reference for calculations) |
| `internal/domain/dto/dashboard_response.go` | Response DTOs (reference for data structures) |

---

### 3. Dependencies

**Prerequisites:**
- [x] RDL-078 - UserConfig service created (Done)
- [x] RDL-079 - DashboardRepository interface created (Done)
- [x] RDL-080 - DashboardResponse DTOs created (Done)
- [x] RDL-081 - DashboardHandler with all 8 endpoints (Done)
- [x] RDL-082 - DayService implemented (Done)
- [x] RDL-083 - ProjectsService implemented (Done)
- [x] RDL-084 - FaultsService implemented (Done)
- [x] RDL-085 - SpeculateService implemented (Done)
- [x] RDL-086 - WeekdayFaultsService implemented (Done)
- [x] RDL-087 - MeanProgressService implemented (Done)
- [x] RDL-091 - API documentation updated (Done)

**No external dependencies required for this task.**

---

### 4. Code Patterns

Since this is a documentation task, the "patterns" to follow are **documentation conventions**:

#### Documentation Structure Pattern
```
# Document Title

## Table of Contents
- Overview
- Key Concepts
- API Reference
- Calculation Methods
- Configuration
- Examples
- Troubleshooting

## Sections
### Overview
[High-level description]

### Key Concepts
[Fundamental concepts]

### [Topic]
[Detailed explanation with code examples]
```

#### Code Example Pattern
```go
// Include relevant Go code snippets
// Use proper syntax highlighting
// Show both implementation and usage
```

#### Calculation Formula Pattern
```
Formula: name = expression

Where:
- term1: description
- term2: description

Example: value = (a + b) * 2
         value = (10 + 5) * 2 = 30
```

---

### 5. Testing Strategy

**This is a documentation task - no unit or integration tests required.**

However, the documentation should include:
- **Code examples** that can be verified against actual implementation
- **Calculation examples** with known input/output values
- **Configuration examples** that match existing `.env.example` format

**Verification:**
- Review existing service implementations to ensure documentation accuracy
- Verify calculation formulas match `internal/service/dashboard/*.go`
- Check configuration keys match `internal/service/user_config_service.go`
- Validate error handling examples against handler implementations

---

### 6. Risks and Considerations

#### Potential Issues

| Risk | Impact | Mitigation |
|------|--------|------------|
| Documentation drift from implementation | Medium | Reference specific commit/branch; update with code changes |
| Incomplete coverage of all 8 endpoints | Low | Use PRD acceptance criteria as checklist |
| Outdated calculation explanations | Medium | Link to source code for definitive implementation |

#### Design Decisions

1. **Single File Approach**: Create one comprehensive guide rather than multiple files, making it easier for developers to find information.

2. **No Code Generation**: This is pure documentation; no generated code or auto-docs needed.

3. **Location**: Place in `docs/` directory alongside existing API documentation for consistency.

4. **Reference Style**: Link to actual source code for implementation details rather than duplicating code in docs.

#### Acceptance Criteria Mapping

| AC | Status | Verification |
|----|--------|--------------|
| #1 Calculation methodologies documented | To Do | All 8 endpoints have calculation explanations |
| #2 Configuration options explained | To Do | UserConfig service options fully documented |
| #3 Troubleshooting guide created | To Do | Common issues and solutions documented |
| #4 Developer onboarding information included | To Do | Getting started section with examples |

---

## Implementation Steps

### Step 1: Create Documentation Structure
```bash
touch docs/dashboard-developer-guide.md
```

### Step 2: Write Overview and Key Concepts
- Dashboard architecture overview
- Clean Architecture layers for dashboard
- How services interact with repositories
- Request/response flow

### Step 3: Document Calculation Methodologies
For each of the 8 endpoints, document:
- Purpose and use case
- Input parameters
- Calculation formulas
- Output structure
- Example responses

### Step 4: Document Configuration Options
- UserConfig service interface
- All configuration keys
- Default values
- How to customize
- Environment variable examples

### Step 5: Create Troubleshooting Section
- Common errors and solutions
- Debugging tips
- Performance considerations
- Testing strategies

### Step 6: Add Code Examples
- Service usage examples
- Repository query patterns
- Handler implementation patterns
- Test examples

---

## Estimated Effort

| Task | Estimate |
|------|----------|
| Create documentation structure | 15 minutes |
| Write overview and key concepts | 30 minutes |
| Document calculation methodologies (8 endpoints) | 60 minutes |
| Document configuration options | 20 minutes |
| Create troubleshooting guide | 20 minutes |
| Add code examples | 20 minutes |
| Review and refine | 15 minutes |

**Total: ~2.5 hours**

---

## Notes

- This task is **low priority** (Low) and part of Phase 5
- No breaking changes or code modifications needed
- Documentation should be clear enough for new developers to understand the dashboard system
- Link to related documentation (API reference, project README) for cross-referencing
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
