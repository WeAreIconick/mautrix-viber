package logger_test

import (
	"github.com/example/mautrix-viber/internal/logger"
)

func ExampleDebug() {
	// Log a debug message
	logger.Debug("Processing webhook", "event_type", "message", "sender_id", "user123")
}

func ExampleInfo() {
	// Log an info message
	logger.Info("Server started", "address", ":8080", "version", "1.0.0")
}

func ExampleWarn() {
	// Log a warning
	logger.Warn("Rate limit approaching", "user_id", "user123", "current_rate", 90, "limit", 100)
}

func ExampleError() {
	// Log an error
	logger.Error("Failed to send message", "error", "connection timeout", "receiver", "user456")
}

func ExampleSetLevel() {
	// Set log level to debug for detailed logging
	logger.SetLevel("debug")

	// Now all debug messages will be logged
	logger.Debug("This debug message will be logged")
}
