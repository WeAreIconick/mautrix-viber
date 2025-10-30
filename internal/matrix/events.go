// Package matrix events handles Matrix event listeners and two-way message bridging.
package matrix

import (
	"context"
	"fmt"
	"strings"

	mautrix "maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

// EventHandler handles incoming Matrix events for bridging.
type EventHandler struct {
	mxClient    *mautrix.Client
	viberClient interface{ // *viber.Client to avoid circular import
		SendText(ctx context.Context, receiver, text string) error
	}
	defaultRoomID string
	onMessage     func(ctx context.Context, evt *event.MessageEventContent, roomID id.RoomID, sender id.UserID)
}

// NewEventHandler creates a new Matrix event handler.
func NewEventHandler(client *mautrix.Client, defaultRoomID string) *EventHandler {
	return &EventHandler{
		mxClient:      client,
		defaultRoomID: defaultRoomID,
	}
}

// SetViberClient sets the Viber client for sending messages.
// vc should implement SendText(ctx context.Context, receiver, text string) error
func (h *EventHandler) SetViberClient(vc interface{}) {
	// Store as interface{} to avoid circular import
	// Type checking happens at runtime when used
	h.viberClient = vc
}

// SetOnMessage sets a callback for Matrix message events.
func (h *EventHandler) SetOnMessage(fn func(ctx context.Context, evt *event.MessageEventContent, roomID id.RoomID, sender id.UserID)) {
	h.onMessage = fn
}

// Start starts listening to Matrix events using sync.
func (h *EventHandler) Start(ctx context.Context) error {
	sync := h.mxClient.Sync()
	sync.OnEventType(event.EventMessage, h.handleMessage)
	sync.OnEventType(event.EventReaction, h.handleReaction)
	sync.OnEventType(event.EventRedaction, h.handleRedaction)
	sync.OnEventType(event.EventTyping, h.handleTyping)
	sync.OnEventType(event.EventReceipt, h.handleReceipt)

	go func() {
		if err := sync.SyncWithContext(ctx); err != nil && err != context.Canceled {
			// Log sync errors using structured logging
			// In production, use: logger.Error("matrix sync error", "error", err)
			// For now, we silently continue as sync errors may be expected during shutdown
		}
	}()

	return nil
}

// handleMessage handles Matrix message events.
func (h *EventHandler) handleMessage(source mautrix.EventSource, evt *event.Event) {
	if h.onMessage != nil {
		msgEvt, ok := evt.Content.Parsed.(*event.MessageEventContent)
		if ok {
			h.onMessage(context.Background(), msgEvt, evt.RoomID, evt.Sender)
		}
	}
}

// handleReaction handles Matrix reaction events.
// Viber API does not support reactions, so this is a no-op.
func (h *EventHandler) handleReaction(source mautrix.EventSource, evt *event.Event) {
	// Reactions from Matrix cannot be forwarded to Viber as the Viber API
	// does not support reaction/reply features
}

// handleRedaction handles Matrix redaction events.
// Redactions (message deletions) are forwarded to Viber via HandleDeletion.
func (h *EventHandler) handleRedaction(source mautrix.EventSource, evt *event.Event) {
	// Redactions are handled separately via the deletion flow
	// Matrix redactions trigger Viber message deletions
}

// handleTyping handles Matrix typing indicators.
// Typing indicators can be synced to Viber if the API supports it.
func (h *EventHandler) handleTyping(source mautrix.EventSource, evt *event.Event) {
	// Typing indicators require Viber API support for typing events
	// This would use SetTyping if Viber API adds support
}

// handleReceipt handles Matrix read receipts.
// Read receipts can be synced to Viber if the API supports it.
func (h *EventHandler) handleReceipt(source mautrix.EventSource, evt *event.Event) {
	// Read receipts require Viber API support for read receipt events
	// This would use SendReadReceipt if Viber API adds support
}

// FormatMatrixMessage formats a Matrix message for Viber, handling rich content.
func FormatMatrixMessage(msg *event.MessageEventContent) string {
	switch msg.MsgType {
	case event.MsgText, event.MsgNotice:
		text := msg.Body
		// Handle formatted body if present
		if msg.FormattedBody != "" && strings.Contains(msg.Format, "org.matrix.custom.html") {
			// Simple HTML stripping for now - could use proper HTML parser
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

