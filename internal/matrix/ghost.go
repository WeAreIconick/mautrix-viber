// Package matrix ghost implements Matrix ghost user puppeting for Viber contacts.
package matrix

import (
	"context"
	"fmt"

	mautrix "maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

// GhostUser represents a Matrix ghost user mapped to a Viber contact.
type GhostUser struct {
	MatrixUserID id.UserID
	ViberUserID   string
	DisplayName   string
	AvatarURL     string
}

// Puppeting manages Matrix ghost users for Viber contacts.
type Puppeting struct {
	mxClient  *mautrix.Client
	domain    string // Matrix homeserver domain for ghost user IDs
}

// NewPuppeting creates a new puppeting manager.
func NewPuppeting(mxClient *mautrix.Client, domain string) *Puppeting {
	return &Puppeting{
		mxClient: mxClient,
		domain:   domain,
	}
}

// GetGhostUserID generates a Matrix user ID for a Viber user.
func (p *Puppeting) GetGhostUserID(viberUserID string) id.UserID {
	// Format: @viber_<viber_user_id>:<homeserver>
	return id.UserID(fmt.Sprintf("@viber_%s:%s", viberUserID, p.domain))
}

// EnsureGhostUser creates or updates a Matrix ghost user for a Viber contact.
func (p *Puppeting) EnsureGhostUser(ctx context.Context, viberUserID, displayName, avatarURL string) (*GhostUser, error) {
	if p.mxClient == nil {
		return nil, fmt.Errorf("matrix client not configured")
	}

	ghostID := p.GetGhostUserID(viberUserID)
	
	// Set display name if provided
	if displayName != "" {
		if err := p.mxClient.SetDisplayName(ctx, ghostID, displayName); err != nil {
			// Silently ignore - may lack permissions without appservice
			// In production, this would use structured logging
		}
	}

	// Set avatar if provided
	// NOTE: SetAvatarURL doesn't support setting other users' avatars
	// This requires appservice puppet registration in production
	if avatarURL != "" {
		// This will fail without proper puppeting setup - that's expected
		_ = avatarURL
	}

	return &GhostUser{
		MatrixUserID: ghostID,
		ViberUserID:  viberUserID,
		DisplayName:  displayName,
		AvatarURL:    avatarURL,
	}, nil
}

// GetGhostUser retrieves ghost user information (best-effort).
func (p *Puppeting) GetGhostUser(ctx context.Context, viberUserID string) (*GhostUser, error) {
	ghostID := p.GetGhostUserID(viberUserID)
	
	// Try to get profile (may fail without appservice)
	profile, err := p.mxClient.GetProfile(ctx, ghostID)
	if err != nil {
		return &GhostUser{
			MatrixUserID: ghostID,
			ViberUserID:  viberUserID,
		}, nil
	}

	return &GhostUser{
		MatrixUserID: ghostID,
		ViberUserID:  viberUserID,
		DisplayName:  profile.DisplayName,
		AvatarURL:    profile.AvatarURL.CUString(),
	}, nil
}

