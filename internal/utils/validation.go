// Package utils provides utility functions including input validation.
package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	matrixUserIDRegex = regexp.MustCompile(`^@[a-z0-9._=-]+:[a-z0-9.-]+\.[a-z]{2,}$`)
	matrixRoomIDRegex = regexp.MustCompile(`^![a-zA-Z0-9]+:[a-z0-9.-]+\.[a-z]{2,}$`)
)

// ValidateMatrixUserID validates a Matrix user ID format.
func ValidateMatrixUserID(userID string) error {
	if !matrixUserIDRegex.MatchString(userID) {
		return fmt.Errorf("invalid Matrix user ID format: %s", userID)
	}
	return nil
}

// ValidateMatrixRoomID validates a Matrix room ID format.
func ValidateMatrixRoomID(roomID string) error {
	if !matrixRoomIDRegex.MatchString(roomID) {
		return fmt.Errorf("invalid Matrix room ID format: %s", roomID)
	}
	return nil
}

// ValidateURL validates a URL format.
func ValidateURL(urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme")
	}
	if u.Host == "" {
		return fmt.Errorf("URL must have a host")
	}
	return nil
}

// ValidateHTTPS validates that a URL uses HTTPS.
func ValidateHTTPS(urlStr string) error {
	if err := ValidateURL(urlStr); err != nil {
		return err
	}
	u, _ := url.Parse(urlStr)
	if u.Scheme != "https" {
		return fmt.Errorf("URL must use HTTPS scheme for production")
	}
	return nil
}

// SanitizeInput sanitizes user input to prevent injection attacks.
func SanitizeInput(input string) string {
	// Remove potentially dangerous characters
	input = strings.ReplaceAll(input, "<", "")
	input = strings.ReplaceAll(input, ">", "")
	input = strings.ReplaceAll(input, "&", "")
	input = strings.ReplaceAll(input, "'", "")
	input = strings.ReplaceAll(input, "\"", "")
	input = strings.ReplaceAll(input, ";", "")
	input = strings.ReplaceAll(input, "|", "")
	return strings.TrimSpace(input)
}

// ValidateViberUserID validates a Viber user ID format.
func ValidateViberUserID(userID string) error {
	if userID == "" {
		return fmt.Errorf("Viber user ID cannot be empty")
	}
	if len(userID) > 64 {
		return fmt.Errorf("Viber user ID too long")
	}
	return nil
}

