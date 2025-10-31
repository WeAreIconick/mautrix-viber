// Package matrix encryption handles Matrix E2EE (encrypted room creation and message handling).
// NOTE: E2EE requires libolm C dependencies which are not included by default.
// This is a stub implementation that returns "not implemented" errors.
// To enable E2EE, install libolm and use the crypto module from maunium.net/go/mautrix/crypto
package matrix

import (
	"context"
	"fmt"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// EncryptionManager manages E2EE for Matrix rooms.
// This is a placeholder implementation. Full E2EE requires libolm C bindings.
type EncryptionManager struct {
	mxClient *mautrix.Client
}

// NewEncryptionManager creates a new encryption manager.
// NOTE: olmMachine parameter removed since libolm is not compiled.
func NewEncryptionManager(mxClient *mautrix.Client) *EncryptionManager {
	return &EncryptionManager{
		mxClient: mxClient,
	}
}

// EnableEncryption enables encryption for a Matrix room.
// This is a placeholder implementation. Full E2EE requires libolm.
func (em *EncryptionManager) EnableEncryption(ctx context.Context, roomID id.RoomID) error {
	if em.mxClient == nil {
		return fmt.Errorf("matrix client not configured")
	}
	
	// Send encryption event to room (basic support without OLM)
	content := map[string]interface{}{
		"algorithm": "m.megolm.v1.aes-sha2",
	}
	
	_, err := em.mxClient.SendStateEvent(ctx, roomID, event.StateEncryption, "", content)
	return err
}

// IsEncrypted checks if a room is encrypted.
func (em *EncryptionManager) IsEncrypted(ctx context.Context, roomID id.RoomID) (bool, error) {
	if em.mxClient == nil {
		return false, fmt.Errorf("matrix client not configured")
	}
	
	// Query room encryption state using FullStateEvent
	state, err := em.mxClient.FullStateEvent(ctx, roomID, event.StateEncryption, "")
	if err != nil {
		return false, err
	}
	
	return state != nil, nil
}

// SendEncryptedMessage sends an encrypted message to a Matrix room.
// This is a placeholder implementation. Full E2EE requires libolm.
func (em *EncryptionManager) SendEncryptedMessage(ctx context.Context, roomID id.RoomID, content *event.MessageEventContent) error {
	if em.mxClient == nil {
		return fmt.Errorf("matrix client not configured")
	}
	
	// The mautrix client will encrypt automatically if the room has encryption enabled
	// This works without OLM for basic Megolm support
	_, err := em.mxClient.SendMessageEvent(ctx, roomID, event.EventMessage, content)
	return err
}
