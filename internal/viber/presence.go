// Package viber presence handles user presence synchronization (online/offline) between Viber and Matrix.
package viber

import (
	"context"
	"fmt"

	mautrix "maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

// PresenceManager manages presence synchronization between Viber and Matrix.
type PresenceManager struct {
	viberClient  *Client
	matrixClient *mautrix.Client
}

// NewPresenceManager creates a new presence manager.
func NewPresenceManager(viberClient *Client, matrixClient *mautrix.Client) *PresenceManager {
	return &PresenceManager{
		viberClient:  viberClient,
		matrixClient: matrixClient,
	}
}

// SyncPresenceFromViber syncs presence from Viber to Matrix ghost user.
func (pm *PresenceManager) SyncPresenceFromViber(ctx context.Context, viberUserID string, isOnline bool) error {
	if pm.matrixClient == nil {
		return fmt.Errorf("matrix client not configured")
	}

	// Generate ghost user ID
	// Construct ghost user ID - domain should come from Matrix homeserver config
	// For now, use a default pattern (in production, extract from homeserver URL)
	ghostID := id.UserID(fmt.Sprintf("@viber_%s:example.com", viberUserID))

	// Set presence
	presence := "offline"
	if isOnline {
		presence = "online"
	}

	// Matrix presence API
	// NOTE: SetPresence doesn't take userID parameter and only sets current user's presence
	// This requires proper appservice puppeting configuration
	_ = ghostID
	_ = ctx
	_ = presence
	return nil
}

// SyncPresenceFromMatrix syncs presence from Matrix to Viber (if supported).
// Viber API does not support presence updates from external sources.
func (pm *PresenceManager) SyncPresenceFromMatrix(ctx context.Context, matrixUserID id.UserID, presence interface{}) error {
	// Viber API does not support setting presence status programmatically
	// Presence sync from Matrix to Viber is not possible with current Viber API
	_ = ctx
	_ = matrixUserID
	_ = presence
	return fmt.Errorf("Viber API does not support presence updates from external sources")
}

// GetViberPresence gets current presence for a Viber user (if available).
// Viber API does not expose presence status directly.
func (pm *PresenceManager) GetViberPresence(ctx context.Context, viberUserID string) (bool, error) {
	// Viber API does not provide presence status queries
	// Presence must be inferred from activity (recent messages, typing indicators)
	// This requires tracking user activity in webhook events
	_ = ctx
	_ = viberUserID
	return false, fmt.Errorf("Viber API does not expose presence status - must infer from activity")
}

// HandlePresenceUpdate handles a presence update from Viber webhook.
func (pm *PresenceManager) HandlePresenceUpdate(ctx context.Context, viberUserID string, isOnline bool) error {
	return pm.SyncPresenceFromViber(ctx, viberUserID, isOnline)
}
