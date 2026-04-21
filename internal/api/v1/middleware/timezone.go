package middleware

import (
	"context"
	"net/http"

	"go-reading-log-api-next/internal/config"
)

// TimezoneMiddleware sets the timezone in the request context based on configuration.
// This ensures all date calculations use the application's configured timezone (BRT by default),
// matching Rails Date.today behavior exactly.
func TimezoneMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get timezone from config
			tzLocation := cfg.TZLocation

			// Create a new context with the timezone
			ctx := r.Context()
			ctx = context.WithValue(ctx, "timezone", tzLocation)

			// Continue with the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
