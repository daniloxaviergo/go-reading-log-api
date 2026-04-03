---
id: RDL-017
title: bug build go-api not working
status: To Do
assignee: []
created_date: '2026-04-03 10:43'
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
