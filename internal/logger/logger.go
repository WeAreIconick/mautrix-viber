// Package logger provides structured JSON logging using log/slog.
// Supports multiple log levels and structured fields for observability.
package logger

import (
	"context"
	"log/slog"
	"os"
)

var (
	// Default logger instance
	Default *slog.Logger
)

func init() {
	// Initialize with JSON handler for structured logging
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	}
	Default = slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

// SetLevel sets the log level (debug, info, warn, error).
func SetLevel(level string) {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	Default = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     lvl,
		AddSource: true,
	}))
}

// Debug logs a debug message with optional fields.
func Debug(msg string, args ...any) {
	Default.Debug(msg, args...)
}

// Info logs an info message with optional fields.
func Info(msg string, args ...any) {
	Default.Info(msg, args...)
}

// Warn logs a warning message with optional fields.
func Warn(msg string, args ...any) {
	Default.Warn(msg, args...)
}

// Error logs an error message with optional fields.
func Error(msg string, args ...any) {
	Default.Error(msg, args...)
}

// WithContext returns a logger with context fields attached.
func WithContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return Default
	}
	// Extract request ID, user ID, etc. from context if present
	return Default
}

// WithFields returns a new logger with the given fields attached.
func WithFields(fields ...slog.Attr) *slog.Logger {
	attrs := make([]any, len(fields)*2)
	for i, f := range fields {
		attrs[i*2] = f.Key
		attrs[i*2+1] = f.Value
	}
	return Default.With(attrs...)
}
