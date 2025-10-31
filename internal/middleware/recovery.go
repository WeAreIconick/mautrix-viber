// Package middleware provides HTTP middleware including panic recovery.
package middleware

import (
	"net/http"
	"runtime/debug"

	"log/slog"
)

// RecoveryMiddleware recovers from panics and returns a 500 error.
// Logs the panic details for debugging while preventing server crashes.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log panic details
				slog.Error("panic recovered",
					"error", err,
					"path", r.URL.Path,
					"method", r.Method,
					"remote_addr", r.RemoteAddr,
					"stack", string(debug.Stack()),
				)

				// Return 500 error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
