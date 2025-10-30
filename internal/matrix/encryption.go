// Package matrix encryption handles Matrix E2EE (encrypted room creation and message handling).
package matrix

import (
	"context"
	"fmt"

	mautrix "maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/id"
)

// EncryptionManager manages E2EE for Matrix rooms.
type EncryptionManager struct {
	mxClient      *mautrix.Client
	olmMachine    *crypto.OlmMachine
}

// NewEncryptionManager creates a new encryption manager.
func NewEncryptionManager(mxClient *mautrix.Client, olmMachine *crypto.OlmMachine) *EncryptionManager {
	return &EncryptionManager{
		mxClient:   mxClient,
		olmMachine: olmMachine,
	}
}

// EnableEncryption enables encryption for a Matrix room.
func (em *EncryptionManager) EnableEncryption(ctx context.Context, roomID id.RoomID) error {
	if em.mxClient == nil {
		return fmt.Errorf("matrix client not configured")
	}
	
	// Send encryption event to room
	content := map[string]interface{}{
		"algorithm": "m.megolm.v1.aes-sha2",
	}
	
	_, err := em.mxClient.SendStateEvent(ctx, roomID, "m.room.encryption", "", content)
	return err
}

// IsEncrypted checks if a room is encrypted.
func (em *EncryptionManager) IsEncrypted(ctx context.Context, roomID id.RoomID) (bool, error) {
	if em.mxClient == nil {
		return false, fmt.Errorf("matrix client not configured")
	}
	
	// Query room encryption state
	state, err := em.mxClient.GetStateEvent(ctx, roomID, "m.room.encryption", "")
	if err != nil {
		return false, err
	}
	
	return state != nil, nil
}

// SendEncryptedMessage sends an encrypted message to a Matrix room.
func (em *EncryptionManager) SendEncryptedMessage(ctx context.Context, roomID id.RoomID, content *mautrix.MessageEventContent) error {
	if em.olmMachine == nil {
		return fmt.Errorf("olm machine not configured")
	}
	
	// OLM machine handles encryption automatically when room is encrypted
	// The mautrix client will encrypt if the room has encryption enabled
	_, err := em.mxClient.SendMessageEvent(ctx, roomID, mautrix.EventMessage, content)
	return err
}

