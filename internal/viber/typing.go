// Package viber typing handles typing indicators and read receipts synchronization.
package viber

import (
	"context"
	"fmt"
)

// SetTyping sends a typing indicator to Matrix for a Viber user.
// Requires Matrix client implementation of SendTyping method.
func (c *Client) SetTyping(ctx context.Context, roomID string, isTyping bool) error {
	if c.matrix == nil {
		return fmt.Errorf("matrix client not configured")
	}

	// Typing indicators require the Matrix client to implement SendTyping
	// This feature requires mautrix client access to typing event API
	return fmt.Errorf("typing indicators require Matrix client SendTyping method implementation")
}

// SendReadReceipt sends a read receipt to Matrix.
// Requires Matrix client implementation of SendReceipt method.
func (c *Client) SendReadReceipt(ctx context.Context, roomID, eventID string) error {
	if c.matrix == nil {
		return fmt.Errorf("matrix client not configured")
	}

	// Read receipts require the Matrix client to implement SendReceipt
	// This feature requires mautrix client access to receipt event API
	return fmt.Errorf("read receipts require Matrix client SendReceipt method implementation")
}
