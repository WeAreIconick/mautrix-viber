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
	
	ghostID := am.puppeting.GetGhostUserID(viberUserID)
	
	if avatarURL == "" {
		return fmt.Errorf("avatar URL is empty")
	}
	
	// Download avatar from Viber URL (if needed, or use direct URL)
	// For now, assume avatarURL is directly usable
	if err := am.mxClient.SetAvatarURL(ctx, ghostID, mautrix.ContentURI(avatarURL)); err != nil {
		return fmt.Errorf("set avatar: %w", err)
	}
	
	return nil
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
	
	return profile.AvatarURL.CUString(), nil
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

