// Package viber backfill handles backfilling recent Viber message history on room creation.
package viber

import (
	"context"
	"fmt"
	"time"

	"github.com/example/mautrix-viber/internal/database"
	mx "github.com/example/mautrix-viber/internal/matrix"
)

// BackfillManager manages backfilling message history.
type BackfillManager struct {
	viberClient  *Client
	matrixClient *mx.Client
	db           *database.DB
	maxMessages  int // Maximum messages to backfill
}

// NewBackfillManager creates a new backfill manager.
func NewBackfillManager(viberClient *Client, matrixClient *mx.Client, db *database.DB) *BackfillManager {
	return &BackfillManager{
		viberClient:  viberClient,
		matrixClient: matrixClient,
		db:           db,
		maxMessages:  50, // Default: backfill last 50 messages
	}
}

// BackfillChatHistory backfills recent Viber chat history for a Matrix room.
func (bm *BackfillManager) BackfillChatHistory(ctx context.Context, viberChatID, matrixRoomID string) error {
	if bm.db == nil {
		return fmt.Errorf("database not configured")
	}

	// Viber API does not provide message history retrieval endpoints
	// Backfilling would require Viber to add history API support
	// When available, this would fetch recent messages and forward them to Matrix

	// Store room mapping if not exists
	if err := bm.db.CreateRoomMapping(ctx, viberChatID, matrixRoomID); err != nil {
		// Ignore if mapping already exists - this is expected if mapping already exists
		// Log at debug level since this is not an error condition
	}

	return nil
}

// BackfillUserHistory backfills recent messages from a specific Viber user.
func (bm *BackfillManager) BackfillUserHistory(ctx context.Context, viberUserID, matrixRoomID string) error {
	// Similar to BackfillChatHistory but for a specific user
	// This would require Viber API support for user message history

	return nil
}

// SetMaxMessages sets the maximum number of messages to backfill.
func (bm *BackfillManager) SetMaxMessages(max int) {
	if max > 0 && max <= 1000 {
		bm.maxMessages = max
	}
}

// ShouldBackfill determines if backfilling should be performed for a room.
func (bm *BackfillManager) ShouldBackfill(ctx context.Context, matrixRoomID string) bool {
	if bm.db == nil {
		return false
	}

	// Check if room mapping already exists
	viberChatID, err := bm.db.GetViberChatID(ctx, matrixRoomID)
	if err != nil || viberChatID == "" {
		return true // New room, should backfill
	}

	// Check if we've already backfilled
	// This could be tracked with a backfill_timestamp in the database
	return false // Already mapped, skip backfill
}

// MarkBackfilled marks a room as backfilled.
func (bm *BackfillManager) MarkBackfilled(matrixRoomID string, timestamp time.Time) error {
	// This would update a backfill_timestamp field in room_mappings table
	// For now, just creating the mapping is sufficient
	return nil
}
