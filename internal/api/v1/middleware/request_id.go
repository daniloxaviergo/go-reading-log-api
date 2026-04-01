package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	requestIDKey contextKey = "request_id"
)

// RequestIDMiddleware generates a unique UUID for each request
// and propagates it via context for tracing.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a new UUID for each request
		requestID := uuid.New().String()

		// Store request ID in context
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		// Add request ID to response headers
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestIDFromContext retrieves the request ID from context.
func GetRequestIDFromContext(ctx context.Context) string {
	if val := ctx.Value(requestIDKey); val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}
