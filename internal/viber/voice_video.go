// Package viber voice_video handles voice messages and video messages with transcoding if needed.
package viber

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// HandleVoiceMessage handles a Viber voice message and forwards it to Matrix.
func (c *Client) HandleVoiceMessage(ctx context.Context, mediaURL string, duration int, size int64) error {
	if c.matrix == nil {
		return fmt.Errorf("matrix client not configured")
	}
	
	// Download voice message
	resp, err := c.httpClient.Get(mediaURL)
	if err != nil {
		return fmt.Errorf("download voice: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: status %d", resp.StatusCode)
	}
	
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read voice data: %w", err)
	}
	
	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "audio/ogg" // Default for voice messages
	}
	
	// Upload to Matrix and send as m.audio
	// TODO: Implement SendAudio in matrix client
	// For now, send as file
	return c.matrix.SendImage(ctx, fmt.Sprintf("voice_%d.ogg", duration), mimeType, data, nil)
}

// HandleVideoMessage handles a Viber video message and forwards it to Matrix.
func (c *Client) HandleVideoMessage(ctx context.Context, mediaURL, thumbnailURL string, duration int, size int64) error {
	if c.matrix == nil {
		return fmt.Errorf("matrix client not configured")
	}
	
	// Download video
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
		mimeType = "video/mp4" // Default
	}
	
	// Upload to Matrix and send as m.video
	// TODO: Implement SendVideo in matrix client
	// For now, send as file
	return c.matrix.SendImage(ctx, fmt.Sprintf("video_%d.mp4", duration), mimeType, data, nil)
}

// TranscodeIfNeeded transcodes media if needed for Matrix compatibility.
func (c *Client) TranscodeIfNeeded(ctx context.Context, inputData []byte, inputMime, outputMime string) ([]byte, error) {
	// Placeholder for transcoding logic
	// This would require ffmpeg or similar tool
	// For now, return original data if mime types match
	if inputMime == outputMime {
		return inputData, nil
	}
	
	// TODO: Implement actual transcoding
	// Would need external tool like ffmpeg
	return inputData, fmt.Errorf("transcoding not implemented: %s -> %s", inputMime, outputMime)
}

