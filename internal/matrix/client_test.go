// Package matrix tests - unit tests for Matrix client.
package matrix

import (
	"testing"
)

// TestSendText tests sending text messages.
// Requires mock Matrix client implementation.
func TestSendText(t *testing.T) {
	t.Skip("Test requires mock Matrix client implementation - see test/integration/matrix_mock.go")
}

// TestSendImage tests sending image messages.
// Requires mock Matrix client implementation.
func TestSendImage(t *testing.T) {
	t.Skip("Test requires mock Matrix client implementation - see test/integration/matrix_mock.go")
}

// TestRedactEvent tests redacting events.
// Requires mock Matrix client implementation.
func TestRedactEvent(t *testing.T) {
	t.Skip("Test requires mock Matrix client implementation - see test/integration/matrix_mock.go")
}
