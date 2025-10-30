package matrix

import (
    "context"
    "fmt"
    "mime"
    "path/filepath"

    mautrix "maunium.net/go/mautrix"
    "maunium.net/go/mautrix/format"
    "maunium.net/go/mautrix/id"
    "maunium.net/go/mautrix/event"
)

// Client is a minimal Matrix client wrapper for sending messages into a single room.
type Client struct {
    homeserverURL string
    accessToken   string
    defaultRoomID string
    mxClient      *mautrix.Client
}

type Config struct {
    HomeserverURL string
    AccessToken   string
    DefaultRoomID string
}

func NewClient(cfg Config) (*Client, error) {
    mx, err := mautrix.NewClient(cfg.HomeserverURL, "", cfg.AccessToken)
    if err != nil {
        return nil, fmt.Errorf("create matrix client: %w", err)
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
    content := format.RenderMarkdown(text, true, true)
    _, err := c.mxClient.SendMessageEvent(ctx, c.defaultRoomID, mautrix.EventMessage, content)
    if err != nil {
        return fmt.Errorf("send matrix message: %w", err)
    }
    return nil
}

// SendImage uploads bytes to the HS and sends an m.image message.
func (c *Client) SendImage(ctx context.Context, filename string, mimeType string, data []byte, info *mautrix.ImageInfo) error {
    if c.defaultRoomID == "" {
        return fmt.Errorf("default room ID not configured")
    }
    if mimeType == "" {
        // best-effort guess from extension
        mimeType = mime.TypeByExtension(filepath.Ext(filename))
        if mimeType == "" {
            mimeType = "application/octet-stream"
        }
    }
    uploadResp, err := c.mxClient.UploadToContentRepo(data, mimeType, int64(len(data)))
    if err != nil {
        return fmt.Errorf("upload image: %w", err)
    }
    content := mautrix.MessageEventContent{
        MsgType: mautrix.MsgImage,
        Body:    filename,
        URL:     uploadResp.ContentURI.CUString(),
        Info:    info,
    }
    _, err = c.mxClient.SendMessageEvent(ctx, id.RoomID(c.defaultRoomID), mautrix.EventMessage, content)
    if err != nil {
        return fmt.Errorf("send image message: %w", err)
    }
    return nil
}

// EnsureGhostUser attempts to set displayname/avatar for a ghost user.
// Note: Proper puppeting usually requires an appservice. This is a stub for future integration.
func (c *Client) EnsureGhostUser(ctx context.Context, userID id.UserID, displayName string) error {
    // Best-effort: set profile if token has privileges
    if displayName != "" {
        if err := c.mxClient.SetDisplayName(ctx, userID, displayName); err != nil {
            // Ignore errors; may lack permissions
            return nil
        }
    }
    return nil
}

// StartMessageListener starts a background sync and invokes onMessage for message events.
func (c *Client) StartMessageListener(ctx context.Context, onMessage func(ctx context.Context, evt *event.MessageEventContent, roomID id.RoomID, sender id.UserID)) error {
	syncer := c.mxClient.Sync()
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		if evt == nil || evt.Content.Parsed == nil {
			return
		}
		msg, ok := evt.Content.Parsed.(*event.MessageEventContent)
		if !ok {
			return
		}
		onMessage(context.Background(), msg, evt.RoomID, evt.Sender)
	})
	go func() { _ = syncer.SyncWithContext(ctx) }()
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
	
	_, err := c.mxClient.RedactEvent(ctx, roomID, eventID, &mautrix.ReqRedact{
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
	_, err := c.mxClient.SendMessageEvent(ctx, roomID, mautrix.EventMessage, content)
	if err != nil {
		return fmt.Errorf("send matrix message: %w", err)
	}
	return nil
}


