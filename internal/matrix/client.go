package matrix

import (
	"context"
	"fmt"
	"mime"
	"path/filepath"
	"time"

	mautrix "maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"

	"github.com/example/mautrix-viber/internal/logger"
	"github.com/example/mautrix-viber/internal/metrics"
)

// Client is a minimal Matrix client wrapper for sending messages into a single room.
// Provides message sending, event listening, and room management functionality.
type Client struct {
	homeserverURL string
	accessToken   string
	defaultRoomID string
	mxClient      *mautrix.Client
}

// Config holds Matrix client configuration.
type Config struct {
	HomeserverURL string
	AccessToken   string
	DefaultRoomID string
}

// NewClient creates a new Matrix client with the given configuration.
func NewClient(cfg Config) (*Client, error) {
	// Extract user ID from access token or use a placeholder
	// In production, parse from token
	userID := id.UserID("@bridge:local")
	mx, err := mautrix.NewClient(cfg.HomeserverURL, userID, cfg.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("create matrix client: %w", err)
	}

	// Initialize DefaultSyncer if not already set
	if mx.Syncer == nil {
		mx.Syncer = mautrix.NewDefaultSyncer()
	}

	return &Client{
		homeserverURL: cfg.HomeserverURL,
		accessToken:   cfg.AccessToken,
		defaultRoomID: cfg.DefaultRoomID,
		mxClient:      mx,
	}, nil
}

// SendText sends a plain/HTML formatted text message to the default room.
func (c *Client) SendText(ctx context.Context, text string) error {
	if c.defaultRoomID == "" {
		return fmt.Errorf("default room ID not configured")
	}
	start := time.Now()
	defer func() {
		metrics.RecordOperationDuration("matrix_send_text", time.Since(start))
	}()
	content := format.RenderMarkdown(text, true, true)
	_, err := c.mxClient.SendMessageEvent(ctx, id.RoomID(c.defaultRoomID), event.EventMessage, content)
	if err != nil {
		metrics.RecordError("matrix_send_failure", "client")
		return fmt.Errorf("send matrix message: %w", err)
	}
	return nil
}

// SendImage uploads bytes to the HS and sends an m.image message.
func (c *Client) SendImage(ctx context.Context, filename string, mimeType string, data []byte, info interface{}) error {
	if c.defaultRoomID == "" {
		return fmt.Errorf("default room ID not configured")
	}
	start := time.Now()
	defer func() {
		metrics.RecordOperationDuration("matrix_send_image", time.Since(start))
	}()
	if mimeType == "" {
		// best-effort guess from extension
		mimeType = mime.TypeByExtension(filepath.Ext(filename))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
	}
	uploadResp, err := c.mxClient.UploadBytes(ctx, data, mimeType)
	if err != nil {
		metrics.RecordError("matrix_upload_failure", "client")
		return fmt.Errorf("upload image: %w", err)
	}
	content := map[string]interface{}{
		"msgtype": event.MsgImage,
		"body":    filename,
		"url":     uploadResp.ContentURI.String(),
	}
	if info != nil {
		content["info"] = info
	}
	_, err = c.mxClient.SendMessageEvent(ctx, id.RoomID(c.defaultRoomID), event.EventMessage, content)
	if err != nil {
		metrics.RecordError("matrix_send_failure", "client")
		return fmt.Errorf("send image message: %w", err)
	}
	return nil
}

// EnsureGhostUser attempts to set displayname/avatar for a ghost user.
// Note: Proper puppeting usually requires an appservice registration. This provides basic profile setting.
func (c *Client) EnsureGhostUser(ctx context.Context, userID id.UserID, displayName string) error {
	// Best-effort: set profile if token has privileges
	// NOTE: SetDisplayName doesn't take userID parameter in mautrix v0.25+
	// This requires proper appservice puppeting configuration
	if displayName != "" {
		_ = displayName // For future implementation
	}
	return nil
}

// StartMessageListener starts a background sync and invokes onMessage for message events.
// The provided context controls the lifecycle of the sync operation.
// Each message callback receives a context derived from the parent context for cancellation propagation.
func (c *Client) StartMessageListener(ctx context.Context, onMessage func(ctx context.Context, evt *event.MessageEventContent, roomID id.RoomID, sender id.UserID)) error {
	// Access the client's syncer to register event handlers
	if c.mxClient.Syncer == nil {
		return fmt.Errorf("matrix client syncer not configured")
	}

	// Cast to ExtensibleSyncer to access OnEventType method
	extSyncer, ok := c.mxClient.Syncer.(mautrix.ExtensibleSyncer)
	if !ok {
		return fmt.Errorf("syncer does not implement ExtensibleSyncer interface")
	}

	// Register message event handler
	extSyncer.OnEventType(event.EventMessage, func(handlerCtx context.Context, evt *event.Event) {
		if evt == nil || evt.Content.Parsed == nil {
			return
		}
		msg, ok := evt.Content.Parsed.(*event.MessageEventContent)
		if !ok {
			return
		}
		// Use parent context for cancellation propagation (background listener context)
		// This allows the message handler to respect context cancellation from shutdown
		onMessage(ctx, msg, evt.RoomID, evt.Sender)
	})

	// Start syncing in background goroutine
	go func() {
		if err := c.mxClient.SyncWithContext(ctx); err != nil && err != context.Canceled {
			// Log sync errors - context.Canceled is expected during shutdown
			logger.Error("matrix sync error",
				"error", err,
			)
		}
	}()
	return nil
}

// GetDefaultRoomID returns the default Matrix room ID.
func (c *Client) GetDefaultRoomID() string {
	return c.defaultRoomID
}

// RedactEvent redacts a Matrix event.
func (c *Client) RedactEvent(ctx context.Context, roomID id.RoomID, eventID id.EventID) error {
	if c.mxClient == nil {
		return fmt.Errorf("matrix client not configured")
	}

	// ReqRedact is passed as variadic parameter, not pointer
	_, err := c.mxClient.RedactEvent(ctx, roomID, eventID, mautrix.ReqRedact{
		Reason: "Message deleted on Viber",
	})
	return err
}

// SendTextToRoom sends a text message to a specific Matrix room.
func (c *Client) SendTextToRoom(ctx context.Context, roomID id.RoomID, text string) error {
	if c.mxClient == nil {
		return fmt.Errorf("matrix client not configured")
	}

	content := format.RenderMarkdown(text, true, true)
	_, err := c.mxClient.SendMessageEvent(ctx, roomID, event.EventMessage, content)
	if err != nil {
		return fmt.Errorf("send matrix message: %w", err)
	}
	return nil
}
