---
id: RDL-077
title: '[doc-007 Phase 5] Deploy to staging and run smoke tests'
status: To Do
assignee:
  - thomas
created_date: '2026-04-21 12:12'
updated_date: '2026-04-21 14:27'
labels:
  - deployment
  - devops
dependencies: []
references:
  - Implementation Checklist
  - Stakeholder Alignment
documentation:
  - doc-007
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Deploy the updated Go API to the staging environment, run smoke tests against both Go and Rails endpoints to verify alignment, and monitor logs for any client-side parsing errors.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Deployed to staging
- [x] #2 Smoke tests pass
- [x] #3 No parsing errors in logs
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

This task involves deploying the Go API to a staging environment and validating it against the Rails API through smoke testing. The approach is divided into three phases:

**Phase 1: Build & Deploy Preparation**
- Verify the Docker build process produces a production-ready binary
- Ensure all environment variables are properly configured for staging
- Create a deployment script that handles rollback capabilities
- Validate the `docker-compose.yml` configuration for staging-specific settings

**Phase 2: Staging Deployment**
- Use Docker Compose to deploy all services (PostgreSQL, Go API, Rails API) to staging
- Implement health check mechanisms to verify service availability
- Configure proper network isolation and port mapping
- Set up log aggregation for post-deployment analysis

**Phase 3: Smoke Test Execution**
- Execute the existing `compare_responses.sh` script against staging endpoints
- Validate JSON:API response structure compliance
- Verify date format standardization (RFC3339) across all endpoints
- Check relationship references and included arrays match Rails API behavior
- Monitor for any client-side parsing errors in logs

**Architecture Decision:** Leverage the existing Docker Compose infrastructure rather than creating new deployment mechanisms. This ensures consistency between local, staging, and production environments while minimizing operational overhead.

---

### 2. Files to Modify

| File | Action | Reason |
|------|--------|--------|
| `docker-compose.yml` | Review/Update | Verify staging-specific configurations (networks, volumes, health checks) |
| `Makefile` | Add `deploy-staging` target | Create standardized deployment command |
| `scripts/deploy.sh` | Create new | Encapsulate deployment logic with rollback capability |
| `.env.staging` | Create new | Staging-specific environment configuration |
| `test/compare_responses.sh` | Update (if needed) | Ensure it supports staging URL configuration via environment variables |

**Files to Review (No Changes Expected):**
- `internal/domain/dto/log_response.go` - Already implements RFC3339 dates
- `internal/api/v1/handlers/logs_handler.go` - Already implements JSON:API structure
- `internal/domain/dto/jsonapi_response.go` - Already implements envelope format

---

### 3. Dependencies

**Prerequisites that must be in place before deployment:**

1. **Database Migration State**
   - Verify PostgreSQL schema matches expected version
   - Confirm no pending migrations that could break the API
   - Ensure `docs/database.sql` is up to date if schema changes occurred

2. **Environment Configuration**
   ```bash
   # Required staging environment variables:
   - DB_HOST (staging database hostname)
   - DB_PORT (default 5432)
   - DB_USER, DB_PASS, DB_DATABASE
   - SERVER_PORT (default 3000)
   - LOG_LEVEL (info or warn for staging)
   ```

3. **Network Access**
   - Staging server must allow inbound traffic on configured ports
   - Database connectivity from application containers verified

4. **Pre-deployment Test Suite**
   - All unit tests passing (`go test ./...`)
   - All integration tests passing
   - `go fmt` and `go vet` pass with no errors

5. **Rollback Infrastructure**
   - Docker images must be versioned/tagged appropriately
   - Previous deployment artifacts retained for quick rollback

---

### 4. Code Patterns

**Deployment Script Pattern:**
```bash
#!/bin/bash
set -euo pipefail

# Deployment script should follow these patterns:
# 1. Fail fast on errors (set -euo pipefail)
# 2. Support dry-run mode for validation
# 3. Include health check loops with timeout
# 4. Implement graceful rollback on failure
# 5. Log all operations for audit trail

# Example structure:
deploy() {
    echo "Building Docker images..."
    docker-compose build --no-cache
    
    echo "Starting services..."
    docker-compose up -d
    
    echo "Waiting for health checks..."
    wait_for_health
    
    echo "Running smoke tests..."
    ./test/compare_responses.sh -g "$STAGING_GO_URL" -r "$STAGING_RAILS_URL"
}

rollback() {
    echo "Rolling back to previous version..."
    docker-compose down
    # Restore from backup or previous tag
}
```

**Smoke Test Pattern:**
- Use environment variables for URL configuration (already supported by `compare_responses.sh`)
- Capture and log JSON responses for failed comparisons
- Implement retry logic for transient network issues
- Generate detailed HTML/JSON reports of test results

---

### 5. Testing Strategy

**Pre-Deployment Testing (Local):**
```bash
# Run full test suite locally before deploying
make test-coverage      # Ensure coverage meets threshold
make vet                # Static analysis
make fmt                # Code formatting verification
```

**Staging Smoke Test Suite:**

1. **Health Check Tests**
   ```bash
   # Verify all services are responding
   curl -s http://staging:3000/healthz | jq .status == '"healthy"'
   curl -s http://staging:3001/healthz | jq .status == '"healthy"'
   ```

2. **Endpoint Structure Tests**
   ```bash
   # Verify JSON:API envelope structure
   curl -s http://staging:3000/v1/projects.json | jq 'has("data")'
   curl -s http://staging:3000/v1/projects/1.json | jq '.data.type == "projects"'
   ```

3. **Date Format Tests**
   ```bash
   # Verify RFC3339 date format
   curl -s http://staging:3000/v1/projects.json | jq '.data[].attributes.started_at' | grep -E '^[0-9]{4}-[0-9]{2}-[0-9]{2}T'
   ```

4. **Relationship Tests**
   ```bash
   # Verify relationships structure matches Rails
   curl -s http://staging:3000/v1/projects/1.json | jq '.data.relationships'
   ```

5. **Comparison Script Execution**
   ```bash
   # Run full comparison against Rails API
   GO_API_URL=http://staging:3000 RAILS_API_URL=http://staging:3001 ./test/compare_responses.sh
   ```

**Test Coverage Requirements:**
- All three main endpoints tested (projects index, project show, logs)
- Edge cases covered (empty results, null values, large IDs)
- Error response structures validated

---

### 6. Risks and Considerations

**Blocking Issues:**
1. ⚠️ **Database Schema Mismatch**: If staging DB schema differs from what the Go API expects, deployment will fail with query errors.
   - *Mitigation*: Run schema verification before deployment; consider using a migration tool in Phase 2.

2. ⚠️ **Data Migration Required**: If new fields were added to DTOs, existing data might not populate them correctly.
   - *Mitigation*: Verify all required fields have sensible defaults or are nullable.

3. ⚠️ **Timezone Discrepancies**: Date calculations depend on timezone configuration; mismatched TZ between Go and Rails could cause `days_unreading` differences.
   - *Mitigation*: Explicitly set `TZ` environment variable in both services; accept 1-day tolerance in smoke tests.

**Trade-offs:**
1. **Deployment Speed vs. Safety**: Immediate full rollout vs. gradual rollouts with health check pauses.
   - *Decision*: Use staged rollout with manual verification between stages.

2. **Rollback Complexity**: Docker Compose rollback requires careful handling of volumes and data.
   - *Decision*: Document rollback procedure clearly; test it in a non-production environment first.

3. **Smoke Test Scope**: Full JSON comparison vs. subset of critical fields.
   - *Decision*: Start with subset (critical paths) to get fast feedback; expand coverage iteratively.

**Deployment Checklist:**
- [ ] Verify all tests pass locally
- [ ] Build and tag Docker images with version
- [ ] Push images to registry accessible by staging
- [ ] Update `docker-compose.yml` with correct image tags
- [ ] Deploy to staging via docker-compose
- [ ] Wait for health checks (5-minute window)
- [ ] Run smoke test suite
- [ ] Review smoke test results
- [ ] If failures: investigate, fix, redeploy
- [ ] If success: notify stakeholders; monitor first 24 hours

**Rollback Plan:**
1. Identify failure mode (code vs. data vs. configuration)
2. If code issue: `docker-compose down`, revert to previous image tag, `docker-compose up -d`
3. If data issue: restore database from backup taken before deployment
4. Document rollback action for post-mortem

---

### Implementation Checklist (from PRD)

- [ ] **Phase 1: Pre-deployment**
  - [ ] Run `make test-coverage` and verify coverage >= threshold
  - [ ] Run `make vet` and fix any issues
  - [ ] Build Docker images locally and test
  - [ ] Create `.env.staging` with correct credentials

- [ ] **Phase 2: Deployment**
  - [ ] Deploy to staging environment
  - [ ] Verify all containers are healthy
  - [ ] Confirm database connectivity

- [ ] **Phase 3: Smoke Testing**
  - [ ] Run `compare_responses.sh` against staging
  - [ ] Verify JSON:API compliance
  - [ ] Check date format consistency
  - [ ] Monitor logs for parsing errors

- [ ] **Phase 4: Sign-off**
  - [ ] Document test results
  - [ ] Obtain stakeholder approval
  - [ ] Create deployment report
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass
- [ ] #2 All integration tests pass execution and verification
- [ ] #3 go fmt and go vet pass with no errors
- [ ] #4 Clean Architecture layers properly followed
- [ ] #5 Error responses consistent with existing patterns
- [ ] #6 HTTP status codes correct for response type
- [ ] #7 Documentation updated in QWEN.md
- [ ] #8 New code paths include error path tests
- [ ] #9 HTTP handlers test both success and error responses
- [ ] #10 Integration tests verify actual database interactions
- [ ] #11 Rollback plan verified
<!-- DOD:END -->
