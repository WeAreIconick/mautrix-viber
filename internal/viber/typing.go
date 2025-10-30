// Package viber typing handles typing indicators and read receipts synchronization.
package viber

import (
	"context"
	"fmt"
)

// SetTyping sends a typing indicator to Matrix for a Viber user.
func (c *Client) SetTyping(ctx context.Context, roomID string, isTyping bool) error {
	if c.matrix == nil {
		return fmt.Errorf("matrix client not configured")
	}

	// TODO: Implement SendTyping in matrix client
	// For now, this is a placeholder
	_ = roomID
	_ = isTyping
	return nil
}

// SendReadReceipt sends a read receipt to Matrix.
func (c *Client) SendReadReceipt(ctx context.Context, roomID, eventID string) error {
	if c.matrix == nil {
		return fmt.Errorf("matrix client not configured")
	}

	// TODO: Implement SendReceipt in matrix client
	_ = roomID
	_ = eventID
	return nil
}

