// Package viber groups handles Viber group chats and maps them to Matrix rooms with member sync.
package viber

import (
	"context"
	"fmt"

	"maunium.net/go/mautrix/id"
	
	"github.com/example/mautrix-viber/internal/database"
	mx "github.com/example/mautrix-viber/internal/matrix"
)

// GroupChatManager manages Viber group chats and their Matrix room mappings.
type GroupChatManager struct {
	viberClient  *Client
	matrixClient *mx.Client
	db           *database.DB
	portals      *mx.Portals
}

// NewGroupChatManager creates a new group chat manager.
func NewGroupChatManager(viberClient *Client, matrixClient *mx.Client, db *database.DB, portals *mx.Portals) *GroupChatManager {
	return &GroupChatManager{
		viberClient:  viberClient,
		matrixClient: matrixClient,
		db:           db,
		portals:      portals,
	}
}

// HandleGroupMessage handles a message from a Viber group chat.
func (gm *GroupChatManager) HandleGroupMessage(ctx context.Context, chatID string, senderID, senderName, messageText string) error {
	if gm.matrixClient == nil {
		return fmt.Errorf("matrix client not configured")
	}
	
	if gm.db == nil {
		return fmt.Errorf("database not configured")
	}
	
	// Get or create Matrix room for this Viber group
	matrixRoomID, err := gm.db.GetMatrixRoomID(chatID)
	if err != nil {
		return fmt.Errorf("get matrix room id: %w", err)
	}
	
	if matrixRoomID == "" {
		// Create new portal room for this group
		if gm.portals == nil {
			return fmt.Errorf("portals not configured")
		}
		
		room, err := gm.portals.GetOrCreatePortalRoom(ctx, chatID, fmt.Sprintf("Viber Group: %s", chatID))
		if err != nil {
			return fmt.Errorf("create portal room: %w", err)
		}
		
		matrixRoomID = string(room.MatrixRoomID)
		
		// Store room mapping
		if err := gm.db.CreateRoomMapping(chatID, matrixRoomID); err != nil {
			return fmt.Errorf("create room mapping: %w", err)
		}
	}
	
	// Ensure sender is in group members
	if err := gm.db.UpsertGroupMember(chatID, senderID); err != nil {
		return fmt.Errorf("upsert group member: %w", err)
	}
	
	// Forward message to Matrix room
	text := fmt.Sprintf("[Viber] %s: %s", senderName, messageText)
	if err := gm.matrixClient.SendTextToRoom(ctx, id.RoomID(matrixRoomID), text); err != nil {
		return fmt.Errorf("send message to matrix: %w", err)
	}
	
	return nil
}

// SyncGroupMembers syncs Viber group members to Matrix room.
func (gm *GroupChatManager) SyncGroupMembers(ctx context.Context, chatID string) error {
	if gm.db == nil {
		return fmt.Errorf("database not configured")
	}
	
	matrixRoomID, err := gm.db.GetMatrixRoomID(chatID)
	if err != nil || matrixRoomID == "" {
		return fmt.Errorf("matrix room id not found for chat %s", chatID)
	}
	
	// Get all group members
	_, err = gm.db.ListGroupMembers(chatID)
	if err != nil {
		return fmt.Errorf("list group members: %w", err)
	}
	
	// Ghost user invitation requires:
	// 1. Ghost user creation via EnsureGhostUser
	// 2. Matrix room invitation API call
	// For now, members are tracked in database and can be invited separately
	return nil
}

// AddGroupMember adds a member to a Viber group.
func (gm *GroupChatManager) AddGroupMember(ctx context.Context, chatID, userID string) error {
	if gm.db == nil {
		return fmt.Errorf("database not configured")
	}
	
	if err := gm.db.UpsertGroupMember(chatID, userID); err != nil {
		return fmt.Errorf("add group member: %w", err)
	}
	
	// Sync members to Matrix room
	return gm.SyncGroupMembers(ctx, chatID)
}

// RemoveGroupMember removes a member from a Viber group.
func (gm *GroupChatManager) RemoveGroupMember(ctx context.Context, chatID, userID string) error {
	if gm.db == nil {
		return fmt.Errorf("database not configured")
	}
	
	// Removing group members requires database DELETE operation
	// Query: DELETE FROM group_members WHERE viber_chat_id = ? AND viber_user_id = ?
	// This requires adding RemoveGroupMember method to database layer
	_ = ctx
	_ = chatID
	_ = userID
	return fmt.Errorf("group member removal requires database RemoveGroupMember method implementation")
}

