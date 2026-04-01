---
id: RDL-006
title: '[doc-001 Phase 3] Implement API handlers for projects and logs endpoints'
status: To Do
assignee:
  - catarina
created_date: '2026-04-01 00:58'
updated_date: '2026-04-01 02:27'
labels: []
dependencies: []
references:
  - 'PRD Section: Acceptance Criteria'
  - 'Implementation Checklist: API Layer'
  - 'PRD Section: Key Requirements'
documentation:
  - doc-001
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement handlers in internal/api/v1/handlers/ for all required endpoints: GET /api/v1/projects, GET /api/v1/projects/:id, GET /api/v1/projects/:project_id/logs, and GET /healthz.

Each handler should use repository interfaces for data access and return proper JSON responses matching Rails API behavior.

Implement error handling to return {"error": "<resource> not found"} for missing records and {"error": "Internal server error"} for unexpected errors.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 GET /api/v1/projects returns array of projects with eager-loaded logs ordered by logs descending
- [ ] #2 GET /api/v1/projects/:id returns single project with eager-loaded logs
- [ ] #3 GET /api/v1/projects/:project_id/logs returns first 4 logs for project with project eager-loaded
- [ ] #4 Error responses match Rails API format
<!-- AC:END -->
