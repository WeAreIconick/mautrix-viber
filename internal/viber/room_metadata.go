// Package viber room_metadata syncs room names, topics, and avatars between platforms.
package viber

import (
	"context"
	"fmt"

	"maunium.net/go/mautrix/id"

	"github.com/example/mautrix-viber/internal/database"
	mx "github.com/example/mautrix-viber/internal/matrix"
)

// RoomMetadataManager manages room metadata synchronization.
type RoomMetadataManager struct {
	matrixClient *mx.Client
	db           *database.DB
	portals      *mx.Portals
}

// NewRoomMetadataManager creates a new room metadata manager.
func NewRoomMetadataManager(matrixClient *mx.Client, db *database.DB, portals *mx.Portals) *RoomMetadataManager {
	return &RoomMetadataManager{
		matrixClient: matrixClient,
		db:           db,
		portals:      portals,
	}
}

// SyncRoomName syncs room name from Viber to Matrix.
func (rmm *RoomMetadataManager) SyncRoomName(ctx context.Context, viberChatID, name string) error {
	if rmm.db == nil {
		return fmt.Errorf("database not configured")
	}

	matrixRoomID, err := rmm.db.GetMatrixRoomID(ctx, viberChatID)
	if err != nil || matrixRoomID == "" {
		return fmt.Errorf("matrix room id not found for chat %s", viberChatID)
	}

	if rmm.portals == nil {
		return fmt.Errorf("portals not configured")
	}

	// Update room name
	if err := rmm.portals.UpdateRoomMetadata(ctx, id.RoomID(matrixRoomID), name, "", ""); err != nil {
		return fmt.Errorf("update room name: %w", err)
	}

	return nil
}

// SyncRoomTopic syncs room topic from Viber to Matrix.
func (rmm *RoomMetadataManager) SyncRoomTopic(ctx context.Context, viberChatID, topic string) error {
	if rmm.db == nil {
		return fmt.Errorf("database not configured")
	}

	matrixRoomID, err := rmm.db.GetMatrixRoomID(ctx, viberChatID)
	if err != nil || matrixRoomID == "" {
		return fmt.Errorf("matrix room id not found for chat %s", viberChatID)
	}

	if rmm.portals == nil {
		return fmt.Errorf("portals not configured")
	}

	// Update room topic
	if err := rmm.portals.UpdateRoomMetadata(ctx, id.RoomID(matrixRoomID), "", topic, ""); err != nil {
		return fmt.Errorf("update room topic: %w", err)
	}

	return nil
}

// SyncRoomAvatar syncs room avatar from Viber to Matrix.
func (rmm *RoomMetadataManager) SyncRoomAvatar(ctx context.Context, viberChatID, avatarURL string) error {
	if rmm.db == nil {
		return fmt.Errorf("database not configured")
	}

	matrixRoomID, err := rmm.db.GetMatrixRoomID(ctx, viberChatID)
	if err != nil || matrixRoomID == "" {
		return fmt.Errorf("matrix room id not found for chat %s", viberChatID)
	}

	if rmm.portals == nil {
		return fmt.Errorf("portals not configured")
	}

	// Update room avatar
	if err := rmm.portals.UpdateRoomMetadata(ctx, id.RoomID(matrixRoomID), "", "", avatarURL); err != nil {
		return fmt.Errorf("update room avatar: %w", err)
	}

	return nil
}

// SyncAllMetadata syncs all room metadata from Viber to Matrix.
func (rmm *RoomMetadataManager) SyncAllMetadata(ctx context.Context, viberChatID, name, topic, avatarURL string) error {
	if rmm.db == nil {
		return fmt.Errorf("database not configured")
	}

	matrixRoomID, err := rmm.db.GetMatrixRoomID(ctx, viberChatID)
	if err != nil || matrixRoomID == "" {
		return fmt.Errorf("matrix room id not found for chat %s", viberChatID)
	}

	if rmm.portals == nil {
		return fmt.Errorf("portals not configured")
	}

	// Update all metadata at once
	if err := rmm.portals.UpdateRoomMetadata(ctx, id.RoomID(matrixRoomID), name, topic, avatarURL); err != nil {
		return fmt.Errorf("update room metadata: %w", err)
	}

	return nil
}
