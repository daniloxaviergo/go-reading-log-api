---
id: RDL-017
title: bug build go-api not working
status: To Do
assignee:
  - workflow
created_date: '2026-04-03 10:43'
updated_date: '2026-04-03 10:45'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
go-api not build
```sh
docker-compose up go-api
WARN[0000] /home/danilo/scripts/github/go-reading-log-api-next/docker-compose.yml: the attribute `version` is obsolete, it will be ignored, please remove it to avoid potential confusion 
WARN[0000] Docker Compose requires buildx plugin to be installed 
Sending build context to Docker daemon  2.923MB
Step 1/14 : FROM golang:1.25.7-alpine AS builder
 ---> 2677ed46f77c
Step 2/14 : WORKDIR /build
 ---> Using cache
 ---> 5a84e6c0f5be
Step 3/14 : RUN apk add --no-cache git
 ---> Using cache
 ---> 60024c3319e6
Step 4/14 : COPY go.mod go.sum ./
[+] up 0/1
 ⠙ Image go-reading-log-api-next-go-api Building                                                   0.4s
COPY failed: file not found in build context or excluded by .dockerignore: stat go.mod: file does not exist
```
<!-- SECTION:DESCRIPTION:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

The Docker build is failing because the `.dockerignore` file explicitly excludes `go.mod` and `go.sum` files, which are required for the Go module to download dependencies during the Docker build process.

The fix involves:
1. Removing the `go.mod` and `go.sum` lines from `.dockerignore`
2. The `go.mod` file must be included for `go mod download` to work
3. The `go.sum` file is also needed for dependency verification

### 2. Files to Modify

- `.dockerignore` - Remove `go.mod` and `go.sum` from the exclusion list
- `docker-compose.yml` - Remove the deprecated `version: "3.8"` line (optional, just a warning)

### 3. Dependencies

- No prerequisites required
- The `go.mod` and `go.sum` files already exist in the project root
- Docker and Docker Compose must be installed locally

### 4. Code Patterns

- Follow existing `.dockerignore` patterns (comments explaining exclusions)
- Keep the Dockerfile pattern of copying `go.mod`/`go.sum` first for layer caching
- Maintain existing Docker build optimization strategy

### 5. Testing Strategy

- Run `docker-compose build go-api` to verify the build completes successfully
- Run `docker-compose up go-api` to verify the container starts
- Check that the application logs appear indicating successful startup
- Verify the server is listening on the configured port (3000)

### 6. Risks and Considerations

- The `.dockerignore` file currently excludes `go.sum` which is needed for `go mod verify` to work
- The `go.mod` exclusion prevents dependency resolution entirely
- This is a straightforward fix with minimal risk
- After the fix, Docker builds should be faster due to proper layer caching with `go.mod`/`go.sum`

### 7. Corrected Definition of Done for This Bug Fix

- [ ] `docker-compose build go-api` completes successfully without errors
- [ ] `docker-compose up go-api` starts the container without the "file does not exist" error
- [ ] Application logs show successful startup (database connection established, server starting)
- [ ] The `go.mod` and `go.sum` files are correctly copied into the Docker build context
- [ ] No regression in existing Docker build optimizations
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
