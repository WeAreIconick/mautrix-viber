// Package viber power_levels syncs admin/moderator permissions between Viber groups and Matrix rooms.
package viber

import (
	"context"
	"fmt"

	"maunium.net/go/mautrix/id"
	
	"github.com/example/mautrix-viber/internal/database"
	mx "github.com/example/mautrix-viber/internal/matrix"
)

// PowerLevel represents user power level in a room.
type PowerLevel struct {
	UserID     id.UserID
	PowerLevel int64
	Role       string // "admin", "moderator", "member"
}

// PowerLevelManager manages power levels between Viber groups and Matrix rooms.
type PowerLevelManager struct {
	matrixClient *mx.Client
	db           *database.DB
}

// NewPowerLevelManager creates a new power level manager.
func NewPowerLevelManager(matrixClient *mx.Client, db *database.DB) *PowerLevelManager {
	return &PowerLevelManager{
		matrixClient: matrixClient,
		db:           db,
	}
}

// SyncPowerLevels syncs power levels from Viber group to Matrix room.
func (plm *PowerLevelManager) SyncPowerLevels(ctx context.Context, viberChatID string, admins []string, moderators []string) error {
	if plm.db == nil {
		return fmt.Errorf("database not configured")
	}
	
	matrixRoomID, err := plm.db.GetMatrixRoomID(ctx, viberChatID)
	if err != nil || matrixRoomID == "" {
		return fmt.Errorf("matrix room id not found for chat %s", viberChatID)
	}
	
	// Set admin power level (50 in Matrix)
	for _, viberUserID := range admins {
		ghostID := id.UserID(fmt.Sprintf("@viber_%s:example.com", viberUserID))
		if err := plm.setUserPowerLevel(ctx, id.RoomID(matrixRoomID), ghostID, 50); err != nil {
			return fmt.Errorf("set admin power level: %w", err)
		}
	}
	
	// Set moderator power level (50 in Matrix, or separate if needed)
	for _, viberUserID := range moderators {
		ghostID := id.UserID(fmt.Sprintf("@viber_%s:example.com", viberUserID))
		if err := plm.setUserPowerLevel(ctx, id.RoomID(matrixRoomID), ghostID, 50); err != nil {
			return fmt.Errorf("set moderator power level: %w", err)
		}
	}
	
	return nil
}

// setUserPowerLevel sets a user's power level in a Matrix room.
// Requires Matrix client implementation of SetPowerLevel method.
func (plm *PowerLevelManager) setUserPowerLevel(ctx context.Context, roomID id.RoomID, userID id.UserID, level int64) error {
	if plm.matrixClient == nil {
		return fmt.Errorf("matrix client not configured")
	}
	
	// Power level setting requires Matrix client to implement SetPowerLevel
	// This would use PUT /_matrix/client/r0/rooms/{roomId}/state/m.room.power_levels
	// Requires the mautrix client to expose power level modification methods
	_ = ctx
	_ = roomID
	_ = userID
	_ = level
	return fmt.Errorf("power level setting requires Matrix client SetPowerLevel method implementation")
}

// GetPowerLevel gets a user's power level in a Matrix room.
// Returns the user's current power level by querying Matrix room state.
func (plm *PowerLevelManager) GetPowerLevel(ctx context.Context, roomID id.RoomID, userID id.UserID) (int64, error) {
	if plm.matrixClient == nil {
		return 0, fmt.Errorf("matrix client not configured")
	}
	
	// Power level retrieval requires querying Matrix room state:
	// GET /_matrix/client/r0/rooms/{roomId}/state/m.room.power_levels
	// Requires the mautrix client to expose room state query methods
	_ = ctx
	_ = roomID
	_ = userID
	return 0, fmt.Errorf("power level retrieval requires Matrix client room state query methods")
}

// IsAdmin checks if a user is an admin in a room.
func (plm *PowerLevelManager) IsAdmin(ctx context.Context, roomID id.RoomID, userID id.UserID) (bool, error) {
	level, err := plm.GetPowerLevel(ctx, roomID, userID)
	if err != nil {
		return false, err
	}
	
	// Admin power level is typically 50 or higher
	return level >= 50, nil
}

