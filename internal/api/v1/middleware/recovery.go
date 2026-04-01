package middleware

import (
	"log/slog"
	"net/http"
)

// RecoveryMiddleware recovers from panics and returns 500 errors with logging.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace info
				slog.Error("Panic recovered", "error", err)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
