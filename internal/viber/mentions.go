// Package viber mentions handles @mentions and converts them between platforms.
package viber

import (
	"context"
	"fmt"
	"strings"

	mautrix "maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	
	"github.com/example/mautrix-viber/internal/database"
	mx "github.com/example/mautrix-viber/internal/matrix"
)

// MentionManager manages @mentions between Viber and Matrix.
type MentionManager struct {
	matrixClient *mx.Client
	db           *database.DB
}

// NewMentionManager creates a new mention manager.
func NewMentionManager(matrixClient *mx.Client, db *database.DB) *MentionManager {
	return &MentionManager{
		matrixClient: matrixClient,
		db:           db,
	}
}

// HandleMention handles a Matrix @mention and converts it to Viber format.
func (mm *MentionManager) HandleMention(ctx context.Context, text string, mentions []id.UserID) string {
	if mm.db == nil {
		return text // No database, return as-is
	}
	
	result := text
	
	// Convert Matrix @user:domain mentions to Viber @username format
	for _, userID := range mentions {
		// Extract localpart
		parts := strings.Split(string(userID), ":")
		if len(parts) > 0 {
			localpart := strings.TrimPrefix(parts[0], "@")
			
			// Try to find Viber user ID for this Matrix user
			// This would require a Matrix user -> Viber user mapping
			// For now, just convert to plain text format
			result = strings.ReplaceAll(result, string(userID), "@"+localpart)
		}
	}
	
	return result
}

// ExtractMentionsFromViber extracts mentions from Viber text.
func (mm *MentionManager) ExtractMentionsFromViber(text string) []string {
	var mentions []string
	
	// Simple extraction of @mentions (username-like patterns)
	words := strings.Fields(text)
	for _, word := range words {
		if strings.HasPrefix(word, "@") && len(word) > 1 {
			// Remove punctuation
			mention := strings.Trim(word, "@:!?.,;")
			if mention != "" {
				mentions = append(mentions, mention)
			}
		}
	}
	
	return mentions
}

// ConvertViberMentionsToMatrix converts Viber mentions to Matrix user IDs.
func (mm *MentionManager) ConvertViberMentionsToMatrix(ctx context.Context, mentions []string) []id.UserID {
	if mm.db == nil {
		return nil
	}
	
	var matrixUserIDs []id.UserID
	
	for _, mention := range mentions {
		// Try to find Matrix user ID for Viber user
		// This would query the database for user mappings
		// For now, create ghost user ID
		ghostID := id.UserID(fmt.Sprintf("@viber_%s:example.com", mention))
		matrixUserIDs = append(matrixUserIDs, ghostID)
	}
	
	return matrixUserIDs
}

// FormatMessageWithMentions formats a message with mentions highlighted.
func (mm *MentionManager) FormatMessageWithMentions(ctx context.Context, text string, mentions []id.UserID) string {
	result := text
	
	// Replace Matrix mentions with plain text versions
	for _, userID := range mentions {
		parts := strings.Split(string(userID), ":")
		if len(parts) > 0 {
			localpart := strings.TrimPrefix(parts[0], "@")
			result = strings.ReplaceAll(result, string(userID), "@"+localpart)
		}
	}
	
	return result
}

