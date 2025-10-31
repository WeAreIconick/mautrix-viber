// Package matrix avatars handles syncing user avatars from Viber to Matrix ghost users.
package matrix

import (
	"context"
	"fmt"

	mautrix "maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

// AvatarManager manages avatar syncing for ghost users.
type AvatarManager struct {
	mxClient *mautrix.Client
	puppeting *Puppeting
}

// NewAvatarManager creates a new avatar manager.
func NewAvatarManager(mxClient *mautrix.Client, puppeting *Puppeting) *AvatarManager {
	return &AvatarManager{
		mxClient:  mxClient,
		puppeting: puppeting,
	}
}

// SyncAvatar syncs a Viber user's avatar to their Matrix ghost user.
func (am *AvatarManager) SyncAvatar(ctx context.Context, viberUserID, avatarURL string) error {
	if am.puppeting == nil {
		return fmt.Errorf("puppeting not configured")
	}
	
	_ = am.puppeting.GetGhostUserID(viberUserID) // For future use
	
	if avatarURL == "" {
		return fmt.Errorf("avatar URL is empty")
	}
	
	// Parse the URL as a Matrix ContentURI (mxc://server/media)
	// For now, this is a placeholder - SetAvatarURL doesn't take userID parameter
	// and requires proper puppeting setup
	_ = avatarURL
	_ = ctx
	return fmt.Errorf("avatar syncing requires appservice puppeting")
}

// SyncAvatarFromViberUser syncs avatar from a Viber UserDetails object.
func (am *AvatarManager) SyncAvatarFromViberUser(ctx context.Context, viberUserID, avatarURL string) error {
	if avatarURL == "" {
		return nil // No avatar to sync
	}
	
	return am.SyncAvatar(ctx, viberUserID, avatarURL)
}

// GetGhostAvatarURL retrieves the current avatar URL for a ghost user.
func (am *AvatarManager) GetGhostAvatarURL(ctx context.Context, viberUserID string) (string, error) {
	if am.puppeting == nil {
		return "", fmt.Errorf("puppeting not configured")
	}
	
	ghostID := am.puppeting.GetGhostUserID(viberUserID)
	
	profile, err := am.mxClient.GetProfile(ctx, ghostID)
	if err != nil {
		return "", fmt.Errorf("get profile: %w", err)
	}
	
	var empty id.ContentURI
	if profile.AvatarURL == empty {
		return "", nil
	}
	return profile.AvatarURL.String(), nil
}

// UpdateAvatarIfChanged updates the avatar only if it's different from current.
func (am *AvatarManager) UpdateAvatarIfChanged(ctx context.Context, viberUserID, newAvatarURL string) error {
	if newAvatarURL == "" {
		return nil
	}
	
	currentURL, err := am.GetGhostAvatarURL(ctx, viberUserID)
	if err != nil {
		// Ghost user may not exist yet, sync anyway
		return am.SyncAvatar(ctx, viberUserID, newAvatarURL)
	}
	
	if currentURL == newAvatarURL {
		return nil // Already synced
	}
	
	return am.SyncAvatar(ctx, viberUserID, newAvatarURL)
}

