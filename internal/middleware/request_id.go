// Package middleware provides HTTP middleware for request tracking.
package middleware

import (
	"context"
	"net/http"

	"github.com/example/mautrix-viber/internal/logger"
	"github.com/google/uuid"
)

// RequestIDMiddleware adds a unique request ID to each request for tracking.
// The request ID is added as a header and to the request context.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get or generate request ID
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add to response header
		w.Header().Set("X-Request-ID", requestID)

		// Add to context using the shared logger package's function
		ctx := logger.WithRequestID(r.Context(), requestID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// GetRequestID extracts the request ID from the context.
// Delegates to the logger package's GetRequestID function.
func GetRequestID(ctx context.Context) string {
	return logger.GetRequestID(ctx)
}
