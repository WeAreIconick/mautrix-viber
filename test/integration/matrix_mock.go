// Package integration provides mock Matrix client for testing.
package integration

import (
	"context"
	"fmt"
)

// MockMatrixClient provides a mock Matrix client for testing.
type MockMatrixClient struct {
	SentMessages []SentMessage
	ErrorRate    float64 // Percentage of operations that should fail
}

// SentMessage represents a message sent via the mock client.
type SentMessage struct {
	RoomID  string // Matrix room ID as string
	Content string
}

// SendText sends a text message (mock implementation).
func (m *MockMatrixClient) SendText(ctx context.Context, text string) error {
	if m.ErrorRate > 0 && rand() < m.ErrorRate {
		return fmt.Errorf("mock error")
	}
	
	m.SentMessages = append(m.SentMessages, SentMessage{
		Content: text,
	})
	return nil
}

// GetSentMessages returns all messages sent via the mock client.
func (m *MockMatrixClient) GetSentMessages() []SentMessage {
	return m.SentMessages
}

// Reset clears all sent messages.
func (m *MockMatrixClient) Reset() {
	m.SentMessages = []SentMessage{}
}

func rand() float64 {
	return 0.5 // Simplified - use proper random in real implementation
}

