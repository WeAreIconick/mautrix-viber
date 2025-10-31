// Package matrix portals manages Matrix portal rooms for Viber chats.
package matrix

import (
	"context"
	"fmt"

	mautrix "maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

// PortalRoom represents a Matrix room mapped to a Viber chat.
type PortalRoom struct {
	MatrixRoomID id.RoomID
	ViberChatID   string
	Name          string
	Topic         string
	AvatarURL     string
}

// Portals manages Matrix portal rooms for Viber chats.
type Portals struct {
	mxClient *mautrix.Client
}

// NewPortals creates a new portal manager.
func NewPortals(mxClient *mautrix.Client) *Portals {
	return &Portals{
		mxClient: mxClient,
	}
}

// GetOrCreatePortalRoom gets or creates a Matrix room for a Viber chat.
func (p *Portals) GetOrCreatePortalRoom(ctx context.Context, viberChatID, name string) (*PortalRoom, error) {
	if p.mxClient == nil {
		return nil, fmt.Errorf("matrix client not configured")
	}

	// Create room with Viber chat name
	req := &mautrix.ReqCreateRoom{
		Name: name,
		Topic: fmt.Sprintf("Viber chat: %s", viberChatID),
		Preset: "public_chat", // Or "private_chat" based on chat type
	}

	resp, err := p.mxClient.CreateRoom(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create room: %w", err)
	}

	return &PortalRoom{
		MatrixRoomID: resp.RoomID,
		ViberChatID:  viberChatID,
		Name:         name,
		Topic:        req.Topic,
	}, nil
}

// UpdateRoomMetadata updates room name, topic, or avatar.
func (p *Portals) UpdateRoomMetadata(ctx context.Context, roomID id.RoomID, name, topic, avatarURL string) error {
	if name != "" {
		if err := p.mxClient.SetRoomName(ctx, roomID, name); err != nil {
			return fmt.Errorf("set room name: %w", err)
		}
	}

	if topic != "" {
		if err := p.mxClient.SetRoomTopic(ctx, roomID, topic); err != nil {
			return fmt.Errorf("set room topic: %w", err)
		}
	}

	if avatarURL != "" {
		// Parse avatar URL as ContentURI
		avatarURI := id.MustParseContentURI(avatarURL)
		if err := p.mxClient.SetRoomAvatar(ctx, roomID, avatarURI); err != nil {
			return fmt.Errorf("set room avatar: %w", err)
		}
	}

	return nil
}

// InviteGhostUser invites a ghost user to a portal room.
func (p *Portals) InviteGhostUser(ctx context.Context, roomID id.RoomID, ghostUserID id.UserID) error {
	if err := p.mxClient.InviteUser(ctx, roomID, &mautrix.ReqInviteUser{
		UserID: ghostUserID,
	}); err != nil {
		return fmt.Errorf("invite ghost user: %w", err)
	}
	return nil
}

