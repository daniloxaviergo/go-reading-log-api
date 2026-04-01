package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware logs request details (method, path, status, duration).
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get request ID from context for tracing
		requestID := GetRequestIDFromContext(r.Context())

		// Create a response wrapper to capture status code
		wrapper := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)

		// Log request details
		slog.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapper.status,
			"duration", duration.String(),
			"request_id", requestID,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
