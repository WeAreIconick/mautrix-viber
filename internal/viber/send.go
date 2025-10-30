// Package viber send implements Viber API message sending functions.
// Supports text, image, video, file, location, contact, and URL messages.
package viber

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/example/mautrix-viber/internal/metrics"
)

// SendMessageRequest represents a Viber send message API request.
type SendMessageRequest struct {
	Receiver      string            `json:"receiver"`
	Type          string            `json:"type"`
	Text          string            `json:"text,omitempty"`
	Media         string            `json:"media,omitempty"`
	Thumbnail     string            `json:"thumbnail,omitempty"`
	Duration      int               `json:"duration,omitempty"`
	Size          int64             `json:"size,omitempty"`
	FileName      string            `json:"file_name,omitempty"`
	Location      *Location         `json:"location,omitempty"`
	Contact       *Contact          `json:"contact,omitempty"`
	TrackingData  string            `json:"tracking_data,omitempty"`
	Keyboard      *Keyboard         `json:"keyboard,omitempty"`
	RichMedia     *RichMedia        `json:"rich_media,omitempty"`
	MinAPIVersion int               `json:"min_api_version,omitempty"`
}

// Location represents a location for Viber location messages.
type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Contact represents a contact card for Viber contact messages.
type Contact struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Avatar      string `json:"avatar,omitempty"`
}

// Keyboard represents a Viber keyboard layout.
type Keyboard struct {
	Type    string   `json:"Type"`
	Buttons []Button `json:"Buttons"`
}

// Button represents a button in a Viber keyboard.
type Button struct {
	ActionType string `json:"ActionType"`
	ActionBody string `json:"ActionBody"`
	Text       string `json:"Text"`
}

// RichMedia represents rich media content for Viber messages.
type RichMedia struct {
	Type    string   `json:"Type"`
	Buttons []Button `json:"Buttons"`
}

// SendMessageResponse represents the response from Viber send message API.
type SendMessageResponse struct {
	Status        int    `json:"status"`
	StatusMessage string `json:"status_message"`
	MessageToken  int64  `json:"message_token,omitempty"`
	ChatHostname  string `json:"chat_hostname,omitempty"`
}

// SendText sends a text message to a Viber user.
func (c *Client) SendText(ctx context.Context, receiver, text string) (*SendMessageResponse, error) {
	return c.SendMessage(ctx, SendMessageRequest{
		Receiver: receiver,
		Type:     "text",
		Text:     text,
	})
}

// SendImage sends an image message to a Viber user.
// mediaURL should be a publicly accessible image URL.
func (c *Client) SendImage(ctx context.Context, receiver, mediaURL, thumbnailURL string) (*SendMessageResponse, error) {
	return c.SendMessage(ctx, SendMessageRequest{
		Receiver:  receiver,
		Type:      "picture",
		Media:     mediaURL,
		Thumbnail: thumbnailURL,
	})
}

// SendVideo sends a video message to a Viber user.
func (c *Client) SendVideo(ctx context.Context, receiver, mediaURL string, size int64, duration int) (*SendMessageResponse, error) {
	return c.SendMessage(ctx, SendMessageRequest{
		Receiver: receiver,
		Type:     "video",
		Media:    mediaURL,
		Size:     size,
		Duration: duration,
	})
}

// SendFile sends a file message to a Viber user.
func (c *Client) SendFile(ctx context.Context, receiver, mediaURL string, size int64, filename string) (*SendMessageResponse, error) {
	return c.SendMessage(ctx, SendMessageRequest{
		Receiver: receiver,
		Type:     "file",
		Media:    mediaURL,
		Size:     size,
		FileName: filename,
	})
}

// SendLocation sends a location message to a Viber user.
func (c *Client) SendLocation(ctx context.Context, receiver string, lat, lon float64) (*SendMessageResponse, error) {
	return c.SendMessage(ctx, SendMessageRequest{
		Receiver: receiver,
		Type:     "location",
		Location: &Location{Lat: lat, Lon: lon},
	})
}

// SendContact sends a contact card message to a Viber user.
func (c *Client) SendContact(ctx context.Context, receiver string, contact Contact) (*SendMessageResponse, error) {
	return c.SendMessage(ctx, SendMessageRequest{
		Receiver: receiver,
		Type:     "contact",
		Contact:  &contact,
	})
}

// SendURL sends a URL message to a Viber user.
func (c *Client) SendURL(ctx context.Context, receiver, urlStr string) (*SendMessageResponse, error) {
	return c.SendMessage(ctx, SendMessageRequest{
		Receiver: receiver,
		Type:     "url",
		Media:    urlStr,
	})
}

// SendMessage sends a generic message to a Viber user using the send_message API.
func (c *Client) SendMessage(ctx context.Context, req SendMessageRequest) (*SendMessageResponse, error) {
	if c.config.APIToken == "" {
		return nil, fmt.Errorf("api token not configured")
	}

	start := time.Now()
	defer func() {
		metrics.RecordOperationDuration("viber_send_message", time.Since(start))
	}()

	apiBaseURL := c.config.ViberAPIBaseURL
	if apiBaseURL == "" {
		apiBaseURL = "https://chatapi.viber.com"
	}
	apiURL := apiBaseURL + "/pa/send_message"
	body, err := json.Marshal(req)
	if err != nil {
		metrics.RecordError("viber_marshal_failure", "send")
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		metrics.RecordError("viber_request_failure", "send")
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Viber-Auth-Token", c.config.APIToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		metrics.RecordError("viber_send_failure", "send")
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		metrics.RecordError("viber_read_failure", "send")
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		metrics.RecordError("viber_api_error", "send")
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	var sendResp SendMessageResponse
	if err := json.Unmarshal(respBody, &sendResp); err != nil {
		metrics.RecordError("viber_decode_failure", "send")
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if sendResp.Status != 0 {
		metrics.RecordError("viber_api_error", "send")
		return nil, fmt.Errorf("viber api error %d: %s", sendResp.Status, sendResp.StatusMessage)
	}

	return &sendResp, nil
}

// GetUserDetails retrieves user information from Viber API.
func (c *Client) GetUserDetails(ctx context.Context, userID string) (*UserDetails, error) {
	if c.config.APIToken == "" {
		return nil, fmt.Errorf("api token not configured")
	}

	apiBaseURL := c.config.ViberAPIBaseURL
	if apiBaseURL == "" {
		apiBaseURL = "https://chatapi.viber.com"
	}
	apiURL := fmt.Sprintf("%s/pa/get_user_details?id=%s", apiBaseURL, url.QueryEscape(userID))
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("X-Viber-Auth-Token", c.config.APIToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	var userResp struct {
		Status int         `json:"status"`
		User   UserDetails `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if userResp.Status != 0 {
		return nil, fmt.Errorf("viber api error %d", userResp.Status)
	}

	return &userResp.User, nil
}

// UserDetails represents user information from Viber API.
type UserDetails struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar,omitempty"`
	Language    string `json:"language,omitempty"`
	Country     string `json:"country,omitempty"`
	APIVersion int    `json:"api_version,omitempty"`
}

