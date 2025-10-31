// Package matrix events handles Matrix event listeners and two-way message bridging.
package matrix

import (
	"context"
	"fmt"
	"strings"

	mautrix "maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"github.com/example/mautrix-viber/internal/logger"
)

// EventHandler handles incoming Matrix events for bridging.
type EventHandler struct {
	mxClient      *mautrix.Client
	viberClient   interface{} // *viber.Client stored as interface{} to avoid circular import
	defaultRoomID string
	onMessage     func(ctx context.Context, evt *event.MessageEventContent, roomID id.RoomID, sender id.UserID)
	ctx           context.Context // Context for event handling (set by Start)
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
// The provided context controls the lifecycle of the event listener.
func (h *EventHandler) Start(ctx context.Context) error {
	h.ctx = ctx // Store context for use in event handlers

	// Access the client's syncer to register event handlers
	if h.mxClient.Syncer == nil {
		return fmt.Errorf("matrix client syncer not configured")
	}

	// Cast to ExtensibleSyncer to access OnEventType method
	extSyncer, ok := h.mxClient.Syncer.(mautrix.ExtensibleSyncer)
	if !ok {
		return fmt.Errorf("syncer does not implement ExtensibleSyncer interface")
	}

	extSyncer.OnEventType(event.EventMessage, h.handleMessage)
	extSyncer.OnEventType(event.EventReaction, h.handleReaction)
	extSyncer.OnEventType(event.EventRedaction, h.handleRedaction)
	extSyncer.OnEventType(event.EphemeralEventTyping, h.handleTyping)
	extSyncer.OnEventType(event.EphemeralEventReceipt, h.handleReceipt)

	go func() {
		if err := h.mxClient.SyncWithContext(ctx); err != nil && err != context.Canceled {
			// Log sync errors - context.Canceled is expected during shutdown
			logger.Error("matrix sync error",
				"error", err,
			)
		}
	}()

	return nil
}

// handleMessage handles Matrix message events.
// Uses the handler's context for message callbacks to allow cancellation.
func (h *EventHandler) handleMessage(ctx context.Context, evt *event.Event) {
	if h.onMessage != nil {
		msgEvt, ok := evt.Content.Parsed.(*event.MessageEventContent)
		if ok {
			// Use handler's context (derived from sync context) for cancellation propagation
			h.onMessage(h.ctx, msgEvt, evt.RoomID, evt.Sender)
		}
	}
	_ = ctx // unused parameter
}

// handleReaction handles Matrix reaction events.
// Viber API does not support reactions, so this is a no-op.
func (h *EventHandler) handleReaction(ctx context.Context, evt *event.Event) {
	// Reactions from Matrix cannot be forwarded to Viber as the Viber API
	// does not support reaction/reply features
	_, _ = ctx, evt // unused parameters
}

// handleRedaction handles Matrix redaction events.
// Redactions (message deletions) are forwarded to Viber via HandleDeletion.
func (h *EventHandler) handleRedaction(ctx context.Context, evt *event.Event) {
	// Redactions are handled separately via the deletion flow
	// Matrix redactions trigger Viber message deletions
	_, _ = ctx, evt // unused parameters
}

// handleTyping handles Matrix typing indicators.
// Typing indicators can be synced to Viber if the API supports it.
func (h *EventHandler) handleTyping(ctx context.Context, evt *event.Event) {
	// Typing indicators require Viber API support for typing events
	// This would use SetTyping if Viber API adds support
	_, _ = ctx, evt // unused parameters
}

// handleReceipt handles Matrix read receipts.
// Read receipts can be synced to Viber if the API supports it.
func (h *EventHandler) handleReceipt(ctx context.Context, evt *event.Event) {
	// Read receipts require Viber API support for read receipt events
	// This would use SendReadReceipt if Viber API adds support
	_, _ = ctx, evt // unused parameters
}

// FormatMatrixMessage formats a Matrix message for Viber, handling rich content.
func FormatMatrixMessage(msg *event.MessageEventContent) string {
	switch msg.MsgType {
	case event.MsgText, event.MsgNotice:
		text := msg.Body
		// Handle formatted body if present
		if msg.FormattedBody != "" && string(msg.Format) == "org.matrix.custom.html" {
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
