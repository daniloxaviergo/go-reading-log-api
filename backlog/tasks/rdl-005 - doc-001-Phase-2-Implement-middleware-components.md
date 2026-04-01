---
id: RDL-005
title: '[doc-001 Phase 2] Implement middleware components'
status: To Do
assignee:
  - workflow
created_date: '2026-04-01 00:58'
updated_date: '2026-04-01 02:22'
labels: []
dependencies: []
references:
  - 'PRD Section: Technical Decisions'
  - 'Implementation Checklist: Core Components'
  - 'PRD Section: Middleware support'
documentation:
  - doc-001
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create middleware handlers in internal/api/v1/middleware/ for CORS, request ID generation, panic recovery, and request logging.

Implement middleware chain that propagates context with timeout and passes through all request layers.

Ensure CORS allows all origins to match Rails app behavior.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 CORS middleware implemented allowing all origins
- [ ] #2 Request ID middleware generates unique IDs for each request
- [ ] #3 Recovery middleware prevents panic propagation
- [ ] #4 Context propagation with timeout working correctly
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
### 1. Technical Approach

Implement four middleware components for the `net/http` stdlib router:
1. **CORS middleware** - Sets CORS headers (AllowedOrigins: *) to match Rails app behavior
2. **Request ID middleware** - Generates unique UUIDs per request and propagates via context
3. **Panic recovery middleware** - Recovers from panics and returns 500 errors with logging
4. **Request logging middleware** - Logs request details (method, path, status, duration)

Create a middleware chain function that wraps handlers in the correct order:
`Recovery -> CORS -> RequestID -> Logging -> Handler`

Use Go's context package for request-scoped values (request ID propagation).
Implement handlers as `http.HandlerFunc` to work with `net/http` router.
Follow the pattern already established in the codebase (slog logging, context propagation, error wrapping).

### 2. Files to Modify

| File | Action | Reason |
|------|--------|--------|
| `internal/api/v1/middleware/cors.go` | Create | CORS middleware implementation with AllowAllOrigins |
| `internal/api/v1/middleware/request_id.go` | Create | Request ID generation and propagation middleware |
| `internal/api/v1/middleware/recovery.go` | Create | Panic recovery middleware with error logging |
| `internal/api/v1/middleware/logging.go` | Create | Request logging middleware with timing information |
| `internal/api/v1/middleware/middleware.go` | Create | Middleware chain helper to apply all middleware |
| `cmd/server.go` | Modify | Wire up middleware chain and register handlers |

### 3. Dependencies

**Existing dependencies (no new dependencies required):**
- `net/http` (stdlib) - HTTP handling
- `context` (stdlib) - Context propagation
- `log/slog` (stdlib) - Structured logging
- `github.com/google/uuid` - UUID generation for request IDs

**Prerequisites:**
- Task RDL-004 (Configuration management) should be complete - middleware needs config for log level
- Domain layer must be initialized before middleware can use context patterns
- Base logger must be configured before logging middleware can be used

**No blocking issues** - All components are standalone and can be implemented in any order.

### 4. Code Patterns

**Follow existing patterns:**
- Use `context.WithTimeout` for request-level timeouts (already used in repository layer)
- Use slog for logging (already used in config and logger packages)
- Return wrapped errors with `fmt.Errorf("%w", err)` for error chain traversal
- Use pointer types for optional fields (already used in DTOs)
- File structure mirrors repository pattern in `internal/adapter/postgres/`

**Naming conventions:**
- Middleware functions: `MiddlewareNameMiddleware` (e.g., `CORSMiddleware`)
- Error constants: `ErrMiddlewareName` (e.g., `ErrRecovery`)
- Variables: snake_case for context keys (e.g., `requestIDKey`)

**Context key pattern:**
```go
type contextKey string
const requestIDKey contextKey = "request_id"
```

**Error handling pattern:**
```go
if err != nil {
    return nil, fmt.Errorf("operation failed: %w", err)
}
```

### 5. Testing Strategy

**Unit tests for each middleware:**
- `internal/api/v1/middleware/cors_test.go`
- `internal/api/v1/middleware/request_id_test.go`
- `internal/api/v1/middleware/recovery_test.go`
- `internal/api/v1/middleware/logging_test.go`
- `internal/api/v1/middleware/middleware_test.go`

**Test coverage:**
- CORS: Verify headers are set (Access-Control-Allow-Origin, etc.)
- Request ID: Verify unique IDs generated, context propagation works
- Recovery: Verify panic is caught, error logged, 500 response returned
- Logging: Verify log output format includes all required fields

**Test approach:**
- Use `httptest.ResponseRecorder` to capture response
- Use `httptest.NewRequest` to create test requests
- Test middleware in isolation (unit tests)
- Test middleware chain (integration test)
- Expected coverage: >80% on middleware package

### 6. Risks and Considerations

**Key considerations:**

1. **Middleware order matters**: Recovery must be outermost (catches all panics), CORS must be early (sets headers before response), request ID before logging (so log has ID), logging before handler (captures handler execution)

2. **No external dependencies**: Using `net/http` stdlib means middleware must be implemented from scratch (no chi/mux with built-in middleware)

3. **Context propagation**: Request ID stored in context must be accessible in handlers and repository layer for tracing

4. **Error response format**: Should match existing error format: `{"error": "<message>"}`

5. **Performance**: Minimal overhead - context operations are cheap, logging should be async or non-blocking where possible

6. **No breaking changes**: Must work with existing handlers in `cmd/server.go`

**Trade-offs:**
- Not using a router框架 (chi/mux) means manual middleware chaining
- Manual context propagation required (no automatic injection)
- Error responses must be manually formatted in recovery middleware

**Deployment considerations:**
- No configuration changes required (all middleware work with existing config)
- Zero-downtime deployment possible (middleware is additive)
- No database migration required
<!-- SECTION:PLAN:END -->
