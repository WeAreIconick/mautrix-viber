// Package logger provides structured JSON logging using log/slog.
package logger

import (
	"context"
)

// requestIDKey is the context key for request IDs
type requestIDKey struct{}

// WithRequestID adds a request ID to the context and returns a logger with that ID.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

// GetRequestID extracts the request ID from the context.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return ""
}

// ErrorWithContext logs an error message with request context fields.
func ErrorWithContext(ctx context.Context, msg string, err error, args ...any) {
	attrs := append(args, "error", err)
	if requestID := GetRequestID(ctx); requestID != "" {
		attrs = append(attrs, "request_id", requestID)
	}
	Default.Error(msg, attrs...)
}

// WarnWithContext logs a warning with request context fields.
func WarnWithContext(ctx context.Context, msg string, args ...any) {
	attrs := args
	if requestID := GetRequestID(ctx); requestID != "" {
		attrs = append(attrs, "request_id", requestID)
	}
	Default.Warn(msg, attrs...)
}

// InfoWithContext logs an info message with request context fields.
func InfoWithContext(ctx context.Context, msg string, args ...any) {
	attrs := args
	if requestID := GetRequestID(ctx); requestID != "" {
		attrs = append(attrs, "request_id", requestID)
	}
	Default.Info(msg, attrs...)
}

// DebugWithContext logs a debug message with request context fields.
func DebugWithContext(ctx context.Context, msg string, args ...any) {
	attrs := args
	if requestID := GetRequestID(ctx); requestID != "" {
		attrs = append(attrs, "request_id", requestID)
	}
	Default.Debug(msg, attrs...)
}

