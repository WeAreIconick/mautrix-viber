// Package viber file_upload handles downloading files from Viber and uploading to Matrix content repo with progress.
package viber

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FileUploadProgress represents upload progress information.
type FileUploadProgress struct {
	BytesUploaded   int64
	TotalBytes      int64
	PercentComplete float64
	Status          string // "downloading", "uploading", "complete", "error"
}

// UploadFile downloads a file from Viber and uploads it to Matrix with progress tracking.
func (c *Client) UploadFile(ctx context.Context, fileURL, filename, mimeType string, progressChan chan<- FileUploadProgress) error {
	if c.matrix == nil {
		return fmt.Errorf("matrix client not configured")
	}

	defer close(progressChan)

	// Step 1: Download from Viber
	if progressChan != nil {
		progressChan <- FileUploadProgress{Status: "downloading"}
	}

	resp, err := c.httpClient.Get(fileURL)
	if err != nil {
		if progressChan != nil {
			progressChan <- FileUploadProgress{Status: "error"}
		}
		return fmt.Errorf("download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if progressChan != nil {
			progressChan <- FileUploadProgress{Status: "error"}
		}
		return fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	totalBytes := resp.ContentLength
	if mimeType == "" {
		mimeType = resp.Header.Get("Content-Type")
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
	}

	// Read file data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		if progressChan != nil {
			progressChan <- FileUploadProgress{Status: "error"}
		}
		return fmt.Errorf("read file data: %w", err)
	}

	if progressChan != nil {
		progressChan <- FileUploadProgress{
			BytesUploaded:   int64(len(data)),
			TotalBytes:      totalBytes,
			PercentComplete: 50.0,
			Status:          "uploading",
		}
	}

	// Step 2: Upload to Matrix
	if err := c.matrix.SendImage(ctx, filename, mimeType, data, nil); err != nil {
		if progressChan != nil {
			progressChan <- FileUploadProgress{Status: "error"}
		}
		return fmt.Errorf("upload to matrix: %w", err)
	}

	if progressChan != nil {
		progressChan <- FileUploadProgress{
			BytesUploaded:   int64(len(data)),
			TotalBytes:      totalBytes,
			PercentComplete: 100.0,
			Status:          "complete",
		}
	}

	return nil
}

// UploadFileWithRetry uploads a file with retry logic.
func (c *Client) UploadFileWithRetry(ctx context.Context, fileURL, filename, mimeType string, maxRetries int) error {
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			delay := time.Duration(attempt) * time.Second
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		if err := c.UploadFile(ctx, fileURL, filename, mimeType, nil); err != nil {
			lastErr = err
			continue
		}

		return nil
	}

	return fmt.Errorf("max retries (%d) exceeded: %w", maxRetries, lastErr)
}
