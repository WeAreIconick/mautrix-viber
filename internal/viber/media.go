// Package viber media handles full media support: video, audio, files, stickers, locations, contacts.
package viber

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ForwardMedia forwards various media types from Viber to Matrix.
func (c *Client) ForwardMedia(ctx context.Context, msgType, mediaURL, filename, thumbnail string) error {
	if c.matrix == nil {
		return fmt.Errorf("matrix client not configured")
	}

	switch strings.ToLower(msgType) {
	case "video":
		return c.forwardVideo(ctx, mediaURL, filename)
	case "audio", "file":
		return c.forwardFile(ctx, mediaURL, filename)
	case "sticker":
		return c.forwardSticker(ctx, mediaURL, thumbnail)
	case "location":
		return c.forwardLocation(ctx, mediaURL)
	case "contact":
		return c.forwardContact(ctx, mediaURL)
	default:
		return fmt.Errorf("unsupported media type: %s", msgType)
	}
}

// forwardVideo downloads and forwards a video to Matrix.
func (c *Client) forwardVideo(ctx context.Context, mediaURL, filename string) error {
	resp, err := c.httpClient.Get(mediaURL)
	if err != nil {
		return fmt.Errorf("download video: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read video data: %w", err)
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "video/mp4"
	}

	// Upload to Matrix and send as m.video
	// TODO: Implement SendVideo in matrix client
	return c.matrix.SendImage(ctx, filename, mimeType, data, nil) // Placeholder
}

// forwardFile downloads and forwards a file to Matrix.
func (c *Client) forwardFile(ctx context.Context, mediaURL, filename string) error {
	resp, err := c.httpClient.Get(mediaURL)
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read file data: %w", err)
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Upload to Matrix and send as m.file
	// TODO: Implement SendFile in matrix client
	return c.matrix.SendImage(ctx, filename, mimeType, data, nil) // Placeholder
}

// forwardSticker forwards a sticker (as image for now, could be enhanced to Matrix stickers).
func (c *Client) forwardSticker(ctx context.Context, mediaURL, thumbnail string) error {
	// Use thumbnail if available, otherwise media URL
	url := thumbnail
	if url == "" {
		url = mediaURL
	}
	if url == "" {
		return fmt.Errorf("no sticker URL available")
	}

	// Forward as image (Matrix sticker support would require additional implementation)
	return c.forwardImage(ctx, url, "sticker.png")
}

// HandleSticker handles a Viber sticker and forwards it to Matrix.
func (c *Client) HandleSticker(ctx context.Context, stickerURL, thumbnailURL string) error {
	return c.forwardSticker(ctx, stickerURL, thumbnailURL)
}

// forwardImage downloads and forwards an image.
func (c *Client) forwardImage(ctx context.Context, mediaURL, filename string) error {
	resp, err := c.httpClient.Get(mediaURL)
	if err != nil {
		return fmt.Errorf("download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read image data: %w", err)
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/png"
	}

	return c.matrix.SendImage(ctx, filename, mimeType, data, nil)
}

// forwardLocation forwards a location message (text representation).
func (c *Client) forwardLocation(ctx context.Context, locationData string) error {
	// Parse location data and send as text message
	// Format: "Location: lat,lon" or URL
	text := fmt.Sprintf("[Location] %s", locationData)
	return c.matrix.SendText(ctx, text)
}

// forwardContact forwards a contact card (vCard representation).
func (c *Client) forwardContact(ctx context.Context, contactData string) error {
	// Parse contact data and send as text message
	text := fmt.Sprintf("[Contact] %s", contactData)
	return c.matrix.SendText(ctx, text)
}

