// Package viber location handles location sharing with map preview in Matrix.
package viber

import (
	"context"
	"fmt"
)

// HandleLocation handles a Viber location message and forwards it to Matrix with map preview.
func (c *Client) HandleLocation(ctx context.Context, lat, lon float64, label string) error {
	if c.matrix == nil {
		return fmt.Errorf("matrix client not configured")
	}
	
	// Format location message
	locationText := fmt.Sprintf("📍 Location: %s\nLatitude: %.6f\nLongitude: %.6f", label, lat, lon)
	if label == "" {
		locationText = fmt.Sprintf("📍 Location\nLatitude: %.6f\nLongitude: %.6f", lat, lon)
	}
	
	// Generate map preview URL (e.g., OpenStreetMap)
	mapURL := fmt.Sprintf("https://www.openstreetmap.org/?mlat=%.6f&mlon=%.6f&zoom=15", lat, lon)
	
	// Send as text with map URL
	text := fmt.Sprintf("%s\n\nMap: %s", locationText, mapURL)
	
	return c.matrix.SendText(ctx, text)
}

// ForwardLocation forwards a location from Viber location message.
func (c *Client) ForwardLocation(ctx context.Context, locationURL string) error {
	if c.matrix == nil {
		return fmt.Errorf("matrix client not configured")
	}
	
	// Parse location data from URL or data
	// This is a placeholder - actual implementation would parse Viber location format
	text := fmt.Sprintf("[Location] %s", locationURL)
	return c.matrix.SendText(ctx, text)
}

// SendLocationToViber sends a location message to Viber from Matrix.
func (c *Client) SendLocationToViber(ctx context.Context, receiver string, lat, lon float64) (*SendMessageResponse, error) {
	return c.SendLocation(ctx, receiver, lat, lon)
}

