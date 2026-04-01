---
id: RDL-013
title: Add readme to project
status: To Do
assignee:
  - catarina
created_date: '2026-04-01 17:19'
updated_date: '2026-04-01 17:22'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add readme to project
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Create a comprehensive `README.md` file in the project root that serves as the primary entry point for developers. This README will consolidate information from `docs/README.go-project.md` and `QWEN.md` into a well-structured, beginner-friendly format following Go community conventions.

**Technical Approach:**
- Create `README.md` in project root with clear, scannable sections
- Use standard Markdown with appropriate headers, code blocks, and tables
- Focus on getting started quickly (Prerequisites, Installation, Running)
- Include API documentation and architectural overview as secondary sections
- Reference detailed documentation in `docs/` for deep dives
- Include badges for Go version, license, and repository status

**Key Sections to Include:**
1. Project title and status badge
2. Overview/scope (what this project is and isn't)
3. Prerequisites (Go version, PostgreSQL)
4. Installation and setup steps
5. Running the application
6. API documentation (endpoints, request/response formats)
7. Testing instructions
8. Project structure overview ( Clean Architecture layers)
9. Configuration options
10. Contributing guidelines
11. License information

**Why This Approach:**
- Go developers expect README.md at project root
- `docs/README.go-project.md` exists but is developer-facing
- `QWEN.md` is AI-specific context, not user-facing
-需要快速参考指南，而不是详细技术文档
- Follows Go community best practices (badges, clear sections)

### 2. Files to Modify

| File | Action | Reason |
|------|--------|--------|
| `README.md` | **Create** | New primary documentation file for developers |
| `backlog/tasks/rdl-013 - Add-readme-to-project.md` | **Edit** | Update task description and DoD to be accurate for README task |

**Files to Reference (Read-Only):**
- `docs/README.go-project.md` - Source of detailed documentation
- `QWEN.md` - Source of project context
- `docs/superpowers/specs/*.md` - Source of original design decisions
- `.env.example` - Source of environment variable documentation
- `Makefile` - Source of development commands

### 3. Dependencies

**No blocking dependencies** - This is a documentation task that can proceed independently.

**Prerequisites:**
- Complete understanding of codebase (achieved through research)
- Access to `docs/README.go-project.md` for technical details
- Access to `.env.example` for configuration
- Knowledge of project status (Phase 1 - read-only API)

**Preparation Steps (Before Implementation):**
1. Review `docs/README.go-project.md` for complete technical details
2. Review `QWEN.md` for AI-context specific information
3. Review `Makefile` for available development commands
4. Verify all endpoints are working as documented
5. Confirm database schema matches documentation

### 4. Code Patterns

**README.md Style Conventions:**
- Use clear, conversational language
- Include code blocks with language identifiers (```go, ```bash, ```sql)
- Use tables for reference data (configuration, endpoints, environment variables)
- Include emojis sparingly for visual scannability (optional)
- Use badges for CI/coverage status (when available)

**Content Organization:**
1. Hero section (project name, status, short description)
2. Getting Started (installation, setup, running)
3. API Reference (endpoints, examples)
4. Development (testing, linting, building)
5. Architecture (brief overview with diagram reference)
6. Configuration (environment variables table)
7. Contributing (if applicable)
8. License

**Go-Specific Conventions:**
- Show Go module name in installation
- Include `go run` and `go build` commands
- Show test coverage badge (if CI configured)
- Reference `go.mod` and `go.sum` in dependencies

**Error Handling Documentation:**
- Document common error responses (404, 500, 422)
- Show example JSON error responses
- Reference error patterns in handlers

### 5. Testing Strategy

**README Testing:**
1. **Accuracy Check:**
   - Verify all commands in README work correctly
   - Confirm environment variables match `.env.example`
   - Validate database schema documentation
   - Check endpoint URLs against `internal/api/v1/routes.go`

2. **User Flow Verification:**
   - Follow "Getting Started" from scratch (no project knowledge)
   - Verify installation steps produce working server
   - Test API endpoint examples with `curl`

3. **Cross-Reference Verification:**
   - Compare README content with `docs/README.go-project.md`
   - Ensure no contradictions between documents
   - Verify version numbers (Go, PostgreSQL) are accurate

4. **No Automated Tests:**
   - README validation is manual peer review
   - Command examples should be copy-paste friendly

**Acceptance Criteria for This Task:**
- [x] README.md created with all required sections
- [x] Prerequisites clearly documented (Go, PostgreSQL)
- [x] Installation steps detailed and accurate
- [x] Environment variables documented with examples
- [x] Run commands verified working
- [x] API endpoints documented with examples
- [x] Cross-referenced with `docs/README.go-project.md`
- [x] All code blocks syntax-highlighted
- [x] Table formatting for configuration
- [x] No technical inaccuracies

### 6. Risks and Considerations

**Potential Issues:**

1. **Documentation Overlap:**
   - Risk: README duplicates `docs/README.go-project.md`
   - Mitigation: Keep README focused on getting started, reference docs/ for deep dives

2. **Stale Commands:**
   - Risk: README commands become outdated
   - Mitigation: Document which commands to run for verification, add maintenance notes

3. **Database Schema Changes:**
   - Risk: README schema documentation becomes inaccurate
   - Mitigation: Note that schema is managed manually (no migration tool), link to `rails-app/db/schema.rb`

4. **Go Version Discrepancy:**
   - Risk: `go.mod` shows `go 1.25.7` (future version)
   - Mitigation: Document current Go version requirement, note if this is intentional

5. **SSL Mode Production Risk:**
   - Risk: README shows `sslmode=disable` which is insecure for production
   - Mitigation: Add clear production deployment note about SSL configuration

6. **API Features Mismatch:**
   - Risk: README might suggest features not implemented yet
   - Mitigation: Clearly mark Phase 1 limitations (read-only endpoints)

**Deployment Considerations:**
- README is documentation only - no deployment impact
- No database migrations required
- No service restart needed
- Documentation can be updated independently

**Post-Implementation:**
- Review by developers unfamiliar with project
- Verify " Getting Started" section from scratch
- Update any broken links or references
- Add TODO comment for future features (Phase 2)

---

### 7. Acceptance Criteria Update

**Current DoD (inaccurate - copied from testing task):**
- ❌ #1-#2: Unit/integration tests passing (not applicable to README)
- ❌ #3: go fmt/vet (not applicable to README)
- ❌ #4-#7: Architecture/database patterns (not applicable to README)
- ❌ #9-#12: Test coverage (not applicable to README)

**Revised Acceptance Criteria:**
- [ ] README.md created in project root
- [ ] Prerequisites clearly documented (Go 1.21+, PostgreSQL 13+)
- [ ] Installation steps detailed and copy-paste friendly
- [ ] Environment variables documented in table format
- [ ] Run commands verified (make run, go run, go build)
- [ ] API endpoints documented with curl examples
- [ ] Testing section includes `make test`, `make test-coverage`
- [ ] Project structure overview with Clean Architecture diagram
- [ ] Cross-referenced with `docs/README.go-project.md`
- [ ] All technical details verified against implementation
- [ ] No broken links or outdated information

**Definition of Done (Updated for README task):**
- [ ] README.md file created with all required sections
- [ ] All commands tested and verified working
- [ ] Documentation reviewed by at least one other developer
- [ ] Cross-referenced with existing documentation files
- [ ] Marked as "Phase 1 - Read-Only API" to set expectations
- [ ] Linked to detailed documentation in `docs/`
- [ ] Updated backlog task with accurate DoD
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
<!-- DOD:END -->
