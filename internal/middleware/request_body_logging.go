// Package middleware provides HTTP middleware for request tracking.
package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/example/mautrix-viber/internal/logger"
)

// RequestBodyLoggingMiddleware logs request and response bodies for debugging.
// Should only be enabled via ENABLE_REQUEST_LOGGING=true flag (disabled by default).
// Warning: Do not enable in production with sensitive data - logs are written unredacted.
func RequestBodyLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only log if explicitly enabled
		// Skip for non-debug endpoints
		if r.URL.Path == "/metrics" || r.URL.Path == "/healthz" || r.URL.Path == "/readyz" {
			next.ServeHTTP(w, r)
			return
		}

		requestID := logger.GetRequestID(r.Context())

		// Read and restore request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("failed to read request body",
				"request_id", requestID,
				"path", r.URL.Path,
				"error", err,
			)
			next.ServeHTTP(w, r)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		// Log request body (only first 1KB to avoid huge logs)
		if len(body) > 0 {
			bodyPreview := string(body)
			if len(bodyPreview) > 1024 {
				bodyPreview = bodyPreview[:1024] + "... (truncated)"
			}
			slog.Debug("request body",
				"request_id", requestID,
				"method", r.Method,
				"path", r.URL.Path,
				"body_preview", bodyPreview,
			)
		}

		// Wrap response writer to capture response
		wrapped := &bodyResponseWriter{
			ResponseWriter: w,
			requestID:      requestID,
			path:           r.URL.Path,
		}

		next.ServeHTTP(wrapped, r)
	})
}

type bodyResponseWriter struct {
	http.ResponseWriter
	requestID string
	path      string
	wroteBody bool
}

func (rw *bodyResponseWriter) Write(b []byte) (int, error) {
	rw.wroteBody = true

	// Log response body preview (only first 512B to avoid huge logs)
	bodyPreview := string(b)
	if len(bodyPreview) > 512 {
		bodyPreview = bodyPreview[:512] + "... (truncated)"
	}

	slog.Debug("response body",
		"request_id", rw.requestID,
		"path", rw.path,
		"size_bytes", len(b),
		"body_preview", bodyPreview,
	)

	return rw.ResponseWriter.Write(b)
}

func (rw *bodyResponseWriter) WriteHeader(code int) {
	rw.ResponseWriter.WriteHeader(code)
}
