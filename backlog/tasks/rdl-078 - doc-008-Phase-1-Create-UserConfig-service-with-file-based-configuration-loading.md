---
id: RDL-078
title: >-
  [doc-008 Phase 1] Create UserConfig service with file-based configuration
  loading
status: To Do
assignee:
  - workflow
created_date: '2026-04-21 15:49'
updated_date: '2026-04-21 15:54'
labels:
  - phase-1
  - infrastructure
  - config
dependencies: []
references:
  - REQ-DASH-001
  - AC-DASH-001
  - 'Decision 2: UserConfig Implementation Strategy'
documentation:
  - doc-008
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement internal/service/user_config_service.go to load dashboard configuration from YAML file with hardcoded defaults as fallback. The service must handle missing values gracefully and provide type-safe access to max_faults, prediction_pct, and pages_per_day settings.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Configuration loads from config/dashboard.yaml if present
- [ ] #2 Hardcoded defaults used when config file missing
- [ ] #3 All three settings (max_faults, prediction_pct, pages_per_day) accessible with correct types
- [ ] #4 Unit tests cover both file and default paths
<!-- AC:END -->

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
