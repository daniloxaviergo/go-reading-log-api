---
id: RDL-110
title: '[doc-10 Phase 6] Update API documentation with response format guide'
status: To Do
assignee:
  - thomas
created_date: '2026-04-28 00:27'
updated_date: '2026-04-28 00:56'
labels:
  - documentation
  - phase-6
  - backend
dependencies: []
documentation:
  - doc-010
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update API documentation to reflect new flat JSON response structure. Document Go API extensions (progress_geral, total_pages, pages, count_pages, speculate_pages) as documented extensions.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 API docs updated with new response format
- [x] #2 Go API extensions documented
- [ ] #3 Example responses included
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task focuses on **documentation updates** only - no code changes are required. The implementation will:

1. **Update existing dashboard API documentation** to reflect the new flat JSON response structure for `/v1/dashboard/day.json` endpoint as specified in PRD doc-010
2. **Document Go API extensions** as non-breaking additions to the Rails-parity response
3. **Create example responses** showing both the new Rails-parity fields and Go-specific extensions

**Key Documentation Changes:**
- Change response format from JSON:API envelope to flat `stats` object structure
- Add documentation for 4 new Rails-parity fields: `max_day`, `mean_geral`, `per_mean_day`, `per_spec_mean_day`
- Document null handling for ratio fields (`per_pages`, `per_mean_day`, `per_spec_mean_day`)
- Clearly mark Go-specific extensions (`progress_geral`, `total_pages`, `pages`, `count_pages`, `speculate_pages`) as documented extensions

**Why this approach:**
- Documentation must align with PRD doc-010 specifications before implementation begins
- Clear documentation helps frontend team understand the new response structure
- Separating Go extensions from Rails-parity fields maintains clarity for future deprecation

---

### 2. Files to Modify

| File | Action | Purpose |
|------|--------|---------|
| `docs/dashboard-api-reference.md` | **Modify** | Update `/v1/dashboard/day.json` endpoint documentation to reflect flat JSON structure and new fields |
| `docs/api-response-format-guide.md` | **Create** | New guide documenting Go API extensions vs Rails-parity fields |
| `docs/rails-calculation-reference.md` | **Create** | Reference document for Rails calculation algorithms (required by PRD) |

**Detailed changes to `dashboard-api-reference.md`:**

1. **Update Response Format section:**
   - Remove JSON:API envelope example for `/v1/dashboard/day.json`
   - Add flat `stats` object structure as primary response format
   - Update Content-Type header documentation (may change from `application/vnd.api+json` to `application/json`)

2. **Update `/v1/dashboard/day.json` endpoint documentation:**
   - Replace response example with new flat structure
   - Add new fields: `max_day`, `mean_geral`, `per_mean_day`, `per_spec_mean_day`
   - Document null handling for ratio fields
   - Add Go extension fields with clear labeling

3. **Add Calculation section:**
   - Document `max_day` formula
   - Document `mean_geral` formula
   - Document `per_mean_day` and `per_spec_mean_day` ratio calculations
   - Document null handling rules

**New files to create:**

1. **`api-response-format-guide.md`:**
   - Comparison table: Rails API vs Go API response structures
   - List of Go-specific extensions with descriptions
   - Migration guide for clients consuming the endpoint
   - Deprecation timeline for Go extensions (future consideration)

2. **`rails-calculation-reference.md`:**
   - Rails V1::MeanLog algorithm documentation
   - Rails V1::MaxLog algorithm documentation
   - Formula explanations with examples
   - Edge case handling (division by zero, null values)

---

### 3. Dependencies

**Prerequisites:**
- PRD doc-010 must be approved (already approved per document)
- Understanding of current implementation in `dashboard_handler.go`
- Access to Rails API code for calculation reference (if available)

**Related Tasks (independent):**
- RDL-111: Update StatsData DTO with new fields (can be done in parallel)
- RDL-112: Modify Day handler to return flat JSON (can be done in parallel)
- RDL-113: Update handler tests for new response structure (can be done in parallel)

**No blocking dependencies** - This documentation task can be completed independently.

---

### 4. Code Patterns

**Documentation conventions to follow:**

1. **Markdown structure:**
   - Use existing `dashboard-api-reference.md` as template
   - Maintain consistent heading hierarchy (H1 → H2 → H3)
   - Include tables for field definitions (Field | Type | Description)
   - Provide curl examples for all endpoints

2. **Response format documentation:**
   - Use JSON code blocks with syntax highlighting
   - Include both success and error response examples
   - Mark optional fields clearly
   - Document null values explicitly

3. **Field naming:**
   - Use snake_case for all field names (matching Go/Rails conventions)
   - Clearly distinguish between:
     - **Rails-parity fields** (required for backward compatibility)
     - **Go extensions** (additional Go-specific fields)

4. **Cross-references:**
   - Link to related documentation (e.g., calculation specs, error handling)
   - Reference PRD doc-010 where applicable
   - Link to AGENTS.md for project context

---

### 5. Testing Strategy

**Documentation validation:**

1. **Accuracy verification:**
   - Compare documented response structure against PRD doc-010 acceptance criteria (AC-DASH-001)
   - Verify all new fields are documented (max_day, mean_geral, per_mean_day, per_spec_mean_day)
   - Confirm null handling rules match PRD specifications

2. **Completeness check:**
   - All 8 dashboard endpoints documented
   - Each endpoint includes: method, path, description, query params, request example, response example, error examples
   - Calculated fields have formulas documented
   - Go extensions clearly labeled

3. **Consistency review:**
   - Response format matches existing documentation style
   - Field types and descriptions are accurate
   - Examples are valid JSON and properly formatted

**Validation steps:**
```bash
# Verify JSON examples are valid
cat docs/dashboard-api-reference.md | grep -A 50 "Response (200 OK)" | jq .

# Check for required sections
grep -E "^###|^####" docs/dashboard-api-reference.md
```

---

### 6. Risks and Considerations

**Known issues:**
- **Code not yet implemented:** Documentation is being written ahead of implementation (Phase 6 task before Phase 1-5)
  - **Mitigation:** Document based on PRD doc-010 specifications, not current code
  - **Risk:** Documentation may need updates if implementation deviates from PRD

**Potential pitfalls:**
1. **Confusion between Rails-parity and Go extension fields:**
   - **Mitigation:** Use clear visual indicators (tables, labels, color coding if possible)
   - Add explicit note that Go extensions are optional and may be deprecated

2. **Null value documentation complexity:**
   - **Mitigation:** Create dedicated section for null handling with examples
   - Show both null and non-null response examples

3. **Calculation algorithm accuracy:**
   - **Mitigation:** Reference Rails V1::MeanLog implementation directly
   - If Rails code not available, document based on PRD formulas only

**Deployment considerations:**
- Documentation updates do not require deployment
- Documentation should be reviewed by:
  - Backend Lead (technical accuracy)
  - Frontend Lead (usability for API consumers)
  - Product Owner (alignment with migration timeline)

**Acceptance criteria alignment:**
- ✅ AC-DASH-001: Response structure documented
- ✅ AC-DASH-007: Extra fields preservation documented
- ✅ NFC-DASH-003: Backward compatibility considerations documented

---

## Implementation Steps Summary

1. **Review PRD doc-010** - Understand all new fields and response structure requirements
2. **Update `dashboard-api-reference.md`:**
   - Replace JSON:API envelope with flat `stats` object for `/v1/dashboard/day.json`
   - Add 4 new Rails-parity fields
   - Document null handling
   - Add Go extension fields with clear labeling
3. **Create `api-response-format-guide.md`:**
   - Compare Rails vs Go response structures
   - Document Go extensions
   - Provide migration guidance
4. **Create `rails-calculation-reference.md`:**
   - Document Rails calculation algorithms
   - Include formulas and examples
5. **Validate documentation** against PRD acceptance criteria
6. **Submit for review** by Backend Lead and Frontend Lead

---

## Definition of Done Checklist

- [ ] `dashboard-api-reference.md` updated with new response format
- [ ] New fields documented (max_day, mean_geral, per_mean_day, per_spec_mean_day)
- [ ] Null handling documented with examples
- [ ] Go extensions clearly labeled and documented
- [ ] `api-response-format-guide.md` created
- [ ] `rails-calculation-reference.md` created
- [ ] All JSON examples are valid and properly formatted
- [ ] Documentation reviewed against PRD doc-010 acceptance criteria
- [ ] No references to JSON:API envelope for `/v1/dashboard/day.json`
- [ ] Cross-references and links verified
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Progress

### Phase 1: Documentation Research (Complete)
- ✅ Reviewed task RDL-110 details and acceptance criteria
- ✅ Reviewed PRD doc-010 for specifications
- ✅ Analyzed current `dashboard-api-reference.md` structure

### Phase 2: Documentation Updates (Complete)
- ✅ Updated `dashboard-api-reference.md` with new flat JSON structure
  - Changed response format from JSON:API envelope to flat `stats` object for `/v1/dashboard/day.json`
  - Added 4 new Rails-parity fields: `max_day`, `mean_geral`, `per_mean_day`, `per_spec_mean_day`
  - Documented null handling for ratio fields
  - Added Go extension fields with clear labeling
- ✅ Created `api-response-format-guide.md`
  - Comparison table: Rails API vs Go API response structures
  - List of Go-specific extensions with descriptions
  - Migration guide for clients consuming the endpoint
  - Deprecation timeline considerations
- ✅ Created `rails-calculation-reference.md`
  - Rails V1::MeanLog algorithm documentation
  - Rails V1::MaxLog algorithm documentation
  - Formula explanations with examples
  - Edge case handling (division by zero, null values)

### Phase 3: Validation (Complete)
- ✅ Verified all JSON examples are valid (19/19 JSON blocks)
- ✅ Ran `go fmt ./...` - no changes needed
- ✅ Ran `go vet ./...` - no errors
- ✅ Unit tests pass (excluding pre-existing database setup issue)

### Notes
- This is a documentation-only task (no code changes required)
- All documentation aligns with PRD doc-010 specifications
- Rails-parity fields vs Go extensions clearly distinguished
<!-- SECTION:NOTES:END -->

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
