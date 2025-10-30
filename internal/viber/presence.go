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
	ghostID := id.UserID(fmt.Sprintf("@viber_%s:example.com", viberUserID)) // TODO: Get domain from config
	
	// Set presence
	presence := "offline"
	if isOnline {
		presence = "online"
	}
	
	// Matrix presence API
	_, err := pm.matrixClient.SendPresence(ctx, ghostID, mautrix.Presence(presence))
	if err != nil {
		return fmt.Errorf("set presence: %w", err)
	}
	
	return nil
}

// SyncPresenceFromMatrix syncs presence from Matrix to Viber (if supported).
func (pm *PresenceManager) SyncPresenceFromMatrix(ctx context.Context, matrixUserID id.UserID, presence mautrix.Presence) error {
	// Viber may not support presence updates via API
	// This is a placeholder for future support
	_ = matrixUserID
	_ = presence
	return nil
}

// GetViberPresence gets current presence for a Viber user (if available).
func (pm *PresenceManager) GetViberPresence(ctx context.Context, viberUserID string) (bool, error) {
	// Viber API may not expose presence directly
	// This would need to be inferred from activity or webhook events
	// Placeholder for future implementation
	return false, nil
}

// HandlePresenceUpdate handles a presence update from Viber webhook.
func (pm *PresenceManager) HandlePresenceUpdate(ctx context.Context, viberUserID string, isOnline bool) error {
	return pm.SyncPresenceFromViber(ctx, viberUserID, isOnline)
}

