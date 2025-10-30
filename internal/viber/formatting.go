// Package viber formatting handles rich formatting: replies, threads, reactions, markdown parsing, mentions.
package viber

import (
	"context"
	"fmt"
	"strings"

	mautrix "maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// FormatMessage formats a Matrix message for Viber with rich formatting.
func FormatMessage(msg *event.MessageEventContent) string {
	switch msg.MsgType {
	case event.MsgText, event.MsgNotice:
		text := msg.Body
		// Handle formatted body if present
		if msg.FormattedBody != "" && strings.Contains(msg.Format, "org.matrix.custom.html") {
			text = stripHTML(msg.FormattedBody)
		}
		return text
	case event.MsgImage:
		return fmt.Sprintf("[Image: %s]", msg.Body)
	case event.MsgVideo:
		return fmt.Sprintf("[Video: %s]", msg.Body)
	case event.MsgFile:
		return fmt.Sprintf("[File: %s]", msg.Body)
	case event.MsgAudio:
		return fmt.Sprintf("[Audio: %s]", msg.Body)
	default:
		return fmt.Sprintf("[%s]", msg.MsgType)
	}
}

// stripHTML removes HTML tags from text (simple implementation).
func stripHTML(html string) string {
	text := html
	// Remove common HTML tags
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<br>", "\n")
	text = strings.ReplaceAll(text, "<p>", "")
	text = strings.ReplaceAll(text, "</p>", "\n")
	text = strings.ReplaceAll(text, "<b>", "*")
	text = strings.ReplaceAll(text, "</b>", "*")
	text = strings.ReplaceAll(text, "<i>", "_")
	text = strings.ReplaceAll(text, "</i>", "_")
	text = strings.ReplaceAll(text, "<code>", "`")
	text = strings.ReplaceAll(text, "</code>", "`")
	text = strings.ReplaceAll(text, "<strong>", "*")
	text = strings.ReplaceAll(text, "</strong>", "*")
	text = strings.ReplaceAll(text, "<em>", "_")
	text = strings.ReplaceAll(text, "</em>", "_")
	// Remove any remaining HTML tags (simple regex-like replacement)
	for strings.Contains(text, "<") && strings.Contains(text, ">") {
		start := strings.Index(text, "<")
		end := strings.Index(text[start:], ">")
		if end > 0 {
			text = text[:start] + text[start+end+1:]
		} else {
			break
		}
	}
	return strings.TrimSpace(text)
}

// HandleReply extracts reply information from a Matrix message.
func HandleReply(msg *event.MessageEventContent) (replyToEventID string, replyText string) {
	if msg.RelatesTo == nil {
		return "", ""
	}
	if msg.RelatesTo.Type == event.RelInReplyTo {
		return msg.RelatesTo.EventID.String(), msg.Body
	}
	return "", ""
}

// FormatReply formats a message with reply context for Viber.
func FormatReply(originalText, replyText string) string {
	return fmt.Sprintf("Re: %s\n\n%s", originalText, replyText)
}

// HandleReaction extracts reaction information from a Matrix event.
func HandleReaction(evt *event.Event) (reactedToEventID string, reactionKey string) {
	if evt.Type != event.EventReaction {
		return "", ""
	}
	if evt.Content.RelatesTo == nil {
		return "", ""
	}
	if evt.Content.RelatesTo.Type == event.RelAnnotation {
		return evt.Content.RelatesTo.EventID.String(), evt.Content.RelatesTo.Key
	}
	return "", ""
}

// FormatReaction formats a reaction for Viber (text representation).
func FormatReaction(reactionKey string) string {
	return fmt.Sprintf("Reacted with: %s", reactionKey)
}

// HandleMentions extracts mentions from a Matrix message.
func HandleMentions(msg *event.MessageEventContent) []id.UserID {
	// Parse @mentions from body or formatted body
	// This is a simplified implementation
	var mentions []id.UserID
	text := msg.Body
	if msg.FormattedBody != "" {
		text = msg.FormattedBody
	}
	
	// Simple regex-like parsing for @user:domain mentions
	words := strings.Fields(text)
	for _, word := range words {
		if strings.HasPrefix(word, "@") && strings.Contains(word, ":") {
			// Extract user ID (remove punctuation)
			userID := strings.Trim(word, "@:!?.,")
			if strings.Contains(userID, ":") {
				mentions = append(mentions, id.UserID(userID))
			}
		}
	}
	return mentions
}

// FormatMentions formats mentions for Viber (plain text).
func FormatMentions(text string, mentions []id.UserID) string {
	// Replace @user:domain with plain text names
	result := text
	for _, userID := range mentions {
		// Extract localpart
		parts := strings.Split(string(userID), ":")
		if len(parts) > 0 {
			localpart := strings.TrimPrefix(parts[0], "@")
			result = strings.ReplaceAll(result, string(userID), "@"+localpart)
		}
	}
	return result
}

// FormatForViber formats a complete Matrix message for Viber, including replies, reactions, and mentions.
func FormatForViber(msg *event.MessageEventContent, mxClient *mautrix.Client, roomID id.RoomID, ctx context.Context) string {
	// Get base message text
	text := FormatMessage(msg)
	
	// Handle reply if present
	if replyToID, replyText := HandleReply(msg); replyToID != "" && mxClient != nil {
		// Try to get original message (best-effort)
		if evt, err := mxClient.GetEvent(ctx, roomID, id.EventID(replyToID)); err == nil {
			if origMsg, ok := evt.Content.Parsed.(*event.MessageEventContent); ok {
				originalText := origMsg.Body
				text = FormatReply(originalText, replyText)
			}
		}
	}
	
	// Handle mentions
	mentions := HandleMentions(msg)
	if len(mentions) > 0 {
		text = FormatMentions(text, mentions)
	}
	
	return text
}

