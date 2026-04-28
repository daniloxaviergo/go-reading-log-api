---
id: RDL-114
title: '[doc-10 Phase 2] Implement MaxDay field and repository method'
status: To Do
assignee:
  - book
created_date: '2026-04-28 00:29'
updated_date: '2026-04-28 02:06'
labels:
  - repository
  - phase-2
  - backend
dependencies: []
documentation:
  - doc-010
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add GetMaxByWeekday() repository method in dashboard_repository.go to query maximum pages read for a specific weekday. Implementation: max(pages_read_on_each_occurrence_of_weekday).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 GetMaxByWeekday() method implemented in adapter
- [ ] #2 Interface method added to repository contract
- [ ] #3 Returns maximum pages for target weekday
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
