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
func (h *EventHandler) SetViberClient(vc interface{}) {
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
			fmt.Printf("sync error: %v\n", err)
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
func (h *EventHandler) handleReaction(source mautrix.EventSource, evt *event.Event) {
	// TODO: Forward reactions to Viber when supported
}

// handleRedaction handles Matrix redaction events.
func (h *EventHandler) handleRedaction(source mautrix.EventSource, evt *event.Event) {
	// TODO: Forward redactions to Viber (delete messages)
}

// handleTyping handles Matrix typing indicators.
func (h *EventHandler) handleTyping(source mautrix.EventSource, evt *event.Event) {
	// TODO: Sync typing indicators to Viber
}

// handleReceipt handles Matrix read receipts.
func (h *EventHandler) handleReceipt(source mautrix.EventSource, evt *event.Event) {
	// TODO: Sync read receipts to Viber
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

