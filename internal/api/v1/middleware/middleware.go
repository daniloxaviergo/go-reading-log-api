package middleware

import (
	"context"
	"net/http"
	"time"
)

// Chain wraps handlers with middleware in the specified order.
// Order: Recovery -> CORS -> RequestID -> Logging -> Handler
func Chain(next http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	// Apply middleware in reverse order (innermost first)
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}
	return next
}

// ContextWithTimeout creates a context with a timeout for request-level operations.
func ContextWithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

// DefaultTimeout is the default request timeout duration.
const DefaultTimeout = 30 * time.Second
