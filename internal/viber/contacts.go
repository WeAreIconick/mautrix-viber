// Package viber contacts handles contact card sharing (vCard) messages.
package viber

import (
	"context"
	"fmt"
)

// HandleContact handles a Viber contact card message and forwards it to Matrix.
func (c *Client) HandleContact(ctx context.Context, contactName, phoneNumber, avatarURL string) error {
	if c.matrix == nil {
		return fmt.Errorf("matrix client not configured")
	}
	
	// Format contact card as text
	contactText := fmt.Sprintf("📇 Contact Card\nName: %s\nPhone: %s", contactName, phoneNumber)
	if avatarURL != "" {
		contactText += fmt.Sprintf("\nAvatar: %s", avatarURL)
	}
	
	return c.matrix.SendText(ctx, contactText)
}

// ForwardContact forwards a contact from Viber contact message.
func (c *Client) ForwardContact(ctx context.Context, contactData string) error {
	if c.matrix == nil {
		return fmt.Errorf("matrix client not configured")
	}
	
	// Parse vCard data if needed
	text := fmt.Sprintf("[Contact Card] %s", contactData)
	return c.matrix.SendText(ctx, text)
}

// SendContactToViber sends a contact card message to Viber from Matrix.
func (c *Client) SendContactToViber(ctx context.Context, receiver string, contactName, phoneNumber, avatarURL string) error {
	contact := Contact{
		Name:        contactName,
		PhoneNumber: phoneNumber,
		Avatar:      avatarURL,
	}
	_, err := c.SendContact(ctx, receiver, contact)
	return err
}

