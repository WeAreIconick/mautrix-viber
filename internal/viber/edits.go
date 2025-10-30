// Package viber edits handles message edits and deletions (Viber deletions → Matrix redactions).
package viber

import (
	"context"
	"fmt"

	"maunium.net/go/mautrix/id"
)

// HandleDeletion handles a Viber message deletion and redacts the corresponding Matrix event.
func (c *Client) HandleDeletion(ctx context.Context, viberMsgID string) error {
	if c.matrix == nil {
		return fmt.Errorf("matrix client not configured")
	}
	
	if c.db == nil {
		return fmt.Errorf("database not configured")
	}
	
	// Get Matrix event ID from database
	matrixEventID, err := c.db.GetMatrixEventID(viberMsgID)
	if err != nil {
		return fmt.Errorf("get matrix event id: %w", err)
	}
	
	if matrixEventID == "" {
		return fmt.Errorf("matrix event id not found for viber message %s", viberMsgID)
	}
	
	// Get room ID from database (we track this in message_mappings)
	// For now, use default room if available
	roomID := c.matrix.GetDefaultRoomID()
	if roomID == "" {
		return fmt.Errorf("room id not available")
	}
	
	// Redact the Matrix event
	if err := c.matrix.RedactEvent(ctx, id.RoomID(roomID), id.EventID(matrixEventID)); err != nil {
		return fmt.Errorf("redact matrix event: %w", err)
	}
	
	return nil
}

// HandleEdit handles a Viber message edit and updates the corresponding Matrix event.
func (c *Client) HandleEdit(ctx context.Context, viberMsgID, newText string) error {
	if c.matrix == nil {
		return fmt.Errorf("matrix client not configured")
	}
	
	// Matrix doesn't support editing messages directly
	// Instead, we'd redact the old and send a new message with edit indication
	// For now, this is a placeholder
	
	if err := c.HandleDeletion(ctx, viberMsgID); err != nil {
		return fmt.Errorf("delete old message: %w", err)
	}
	
	// Send new message with "(edited)" indicator
	editedText := fmt.Sprintf("%s (edited)", newText)
	return c.matrix.SendText(ctx, editedText)
	
	// TODO: Proper Matrix message editing support would require sending a new message
	// with m.new_content indicating it's an edit of the original
}

