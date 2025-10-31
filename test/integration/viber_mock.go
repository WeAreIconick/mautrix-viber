// Package integration provides mock Viber API for testing.
package integration

import (
	"context"
	"fmt"
	"sync"
)

// MockViberAPI provides a mock Viber API server for testing.
type MockViberAPI struct {
	mu           sync.Mutex
	sentMessages []ViberMessage
	webhookURL   string
	errorRate    float64
}

// ViberMessage represents a message sent via mock Viber API.
type ViberMessage struct {
	ReceiverID string
	Message    string
	Type       string
}

// NewMockViberAPI creates a new mock Viber API.
func NewMockViberAPI() *MockViberAPI {
	return &MockViberAPI{
		sentMessages: []ViberMessage{},
	}
}

// SendText sends a text message (mock implementation).
func (m *MockViberAPI) SendText(ctx context.Context, receiverID, text string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.errorRate > 0 {
		// Simulate errors based on error rate
		return fmt.Errorf("mock API error")
	}

	m.sentMessages = append(m.sentMessages, ViberMessage{
		ReceiverID: receiverID,
		Message:    text,
		Type:       "text",
	})
	return nil
}

// GetSentMessages returns all messages sent via the mock API.
func (m *MockViberAPI) GetSentMessages() []ViberMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]ViberMessage{}, m.sentMessages...)
}

// Reset clears all sent messages.
func (m *MockViberAPI) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sentMessages = []ViberMessage{}
}

// SetErrorRate sets the error rate for simulating failures.
func (m *MockViberAPI) SetErrorRate(rate float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorRate = rate
}

// SetWebhookURL sets the webhook URL for the mock API.
func (m *MockViberAPI) SetWebhookURL(url string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.webhookURL = url
}
