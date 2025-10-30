package matrix

import (
    "context"
    "fmt"
    "mime"
    "path/filepath"

    mautrix "maunium.net/go/mautrix"
    "maunium.net/go/mautrix/format"
    "maunium.net/go/mautrix/id"
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


