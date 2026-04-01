---
id: RDL-005
title: '[doc-001 Phase 2] Implement middleware components'
status: To Do
assignee:
  - thomas
created_date: '2026-04-01 00:58'
updated_date: '2026-04-01 10:41'
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

This task focuses on testing and verification of the already-implemented middleware components. The middleware is fully implemented in `internal/api/v1/middleware/` but lacks tests. The PRD requires comprehensive unit tests for each middleware component.

**Key issues to address:**
- Missing unit tests for each middleware component (cors, request_id, recovery, logging)
- The middleware chain in `middleware.go` is implemented but unused in server.go
- Logging middleware has an inefficient responseWriter pattern with unused GetContext() method
- Gorilla/mux usage in routes.go contradicts PRD net/http decision (will be addressed in RDL-006)

**Implementation approach:**
1. Write unit tests for each middleware component
2. Test middleware chaining and order
3. Fix the logging middleware responseWriter pattern
4. Update server.go to use the Chain helper function
5. Add tests for the Chain helper

### 2. Files to Modify

| File | Action | Reason |
|------|--------|--------|
| `internal/api/v1/middleware/cors_test.go` | Create | Unit tests for CORS middleware (verify headers, preflight handling) |
| `internal/api/v1/middleware/request_id_test.go` | Create | Unit tests for request ID generation and context propagation |
| `internal/api/v1/middleware/recovery_test.go` | Create | Unit tests for panic recovery (verify 500 response, logging) |
| `internal/api/v1/middleware/logging_test.go` | Create | Unit tests for logging middleware (verify log output, timing) |
| `internal/api/v1/middleware/middleware_test.go` | Create | Integration tests for middleware chain ordering |
| `internal/api/v1/middleware/logging.go` | Modify | Fix responseWriter to remove unused GetContext() method |
| `cmd/server.go` | Modify | Use middleware.Chain() helper instead of manual chaining |

### 3. Dependencies

**No new dependencies required:**
- Using existing dependencies: `github.com/google/uuid`, `log/slog` (stdlib)

**Prerequisites:**
- All middleware components must exist (they do)
- Logger must be configured before testing logging middleware
- Server must be able to start and respond to requests

**No blocking issues** - All components are self-contained.

### 4. Code Patterns

**Following existing patterns in the codebase:**

**Test helper pattern:**
```go
type testRequest struct {
    method     string
    url        string
    body       io.Reader
    headers    map[string]string
}

func makeRequest(req testRequest) *httptest.ResponseRecorder {
    r := httptest.NewRequest(req.method, req.url, req.body)
    for k, v := range req.headers {
        r.Header.Set(k, v)
    }
    w := httptest.NewRecorder()
    return w
}
```

**Test patterns for each middleware:**
- **CORS**: Verify headers set, OPTIONS returns 204, normal requests pass through
- **Request ID**: Verify UUID format, context value accessible, uniqueness
- **Recovery**: Verify panic caught, 500 returned, error logged
- **Logging**: Verify log output contains required fields, timing calculated

**Naming conventions:**
- Test files: `{middleware}_test.go`
- Test functions: `Test{Middleware}_{Scenario}` (e.g., `TestCORS_PreflightRequest`)
- Helper functions: `new{Middleware}Handler` (e.g., `newRecoveryHandler`)

### 5. Testing Strategy

**Unit tests for each middleware component:**

1. **cors_test.go**:
   - `TestCORSMiddleware_PreflightRequest` - OPTIONS returns 204
   - `TestCORSMiddleware_NormalRequest` - Normal requests pass through
   - `TestCORSMiddleware_HeadersSet` - CORS headers are set correctly

2. **request_id_test.go**:
   - `TestRequestIDMiddleware_GeneratesUniqueIDs` - Each request gets unique ID
   - `TestRequestIDMiddleware_ContextPropagation` - ID accessible in context
   - `TestRequestIDMiddleware_ResponseHeader` - X-Request-ID header set

3. **recovery_test.go**:
   - `TestRecoveryMiddleware_PanicCaught` - Panic returns 500, no crash
   - `TestRecoveryMiddleware_NoPanic` - Normal requests work
   - `TestRecoveryMiddleware_LoggerCalled` - Error logged on panic

4. **logging_test.go**:
   - `TestLoggingMiddleware_LogsRequest` - Request logged with all fields
   - `TestLoggingMiddleware_TimingCalculated` - Duration calculated correctly
   - `TestLoggingMiddleware_StatusCode` - Status captured in log

5. **middleware_test.go**:
   - `TestChain_MiddlewareOrder` - Chain applies middleware in correct order
   - `TestChain_PanicsCaught` - Recovery is outermost in chain
   - `TestChain_ContextPropagation` - Context passed through chain

**Test coverage target: 80%+ on middleware package**

### 6. Risks and Considerations

**Key considerations:**

1. **Middleware order is critical**: The Chain function applies middleware in reverse order to achieve: Recovery → CORS → RequestID → Logging → Handler. This must be verified in tests.

2. **Logging middleware inefficiency**: The responseWriter has a GetContext() method that returns `context.Background()` (useless). This should be removed.

3. **Testing panicRecovery**: Testing panic recovery requires `recover()` in a deferred function, which can be tricky in tests. Need to use subtests or separate functions.

4. **No mocking framework**: Using only stdlib test+httptest. For complex scenarios (testing logger calls), may need to use channels or sync primitives.

5. **Gorilla/mux vs net/http**: routes.go uses gorilla/mux but PRD specifies net/http. This is a separate issue that will be fixed in RDL-006 when implementing handlers. Not blocking for RDL-005 middleware tests.

**Trade-offs:**

1. **Testing approach**: Using httptest.Recorder for state observation vs mocking. Since we need to test actual behavior, using real ResponseRecorder is better.

2. **No table-driven tests for some cases**: Some middleware tests may not benefit from table-driven approach due to different setup requirements.

3. **Test database not required**: Middleware tests are pure HTTP handling, no database needed.

**Deployment considerations:**
- No configuration changes required
- No database migrations required
- Zero-downtime deployment possible
- All changes are additive (new test files)

### Next Steps

1. Create unit tests for each middleware component
2. Run tests and verify coverage > 80%
3. Fix logging.go responseWriter pattern
4. Update server.go to use middleware.Chain() helper
5. Final verification: run `go test ./internal/api/v1/middleware/... -v`
<!-- SECTION:PLAN:END -->
