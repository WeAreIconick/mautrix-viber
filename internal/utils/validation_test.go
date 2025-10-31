// Package utils tests - unit tests for validation utilities.
package utils

import (
	"testing"
)

func TestValidateMatrixUserID(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{"valid user id", "@alice:matrix.example.com", false},
		{"invalid format", "alice", true},
		{"missing @", "alice:matrix.example.com", true},
		{"missing domain", "@alice", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMatrixUserID(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMatrixUserID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMatrixRoomID(t *testing.T) {
	tests := []struct {
		name    string
		roomID  string
		wantErr bool
	}{
		{"valid room id", "!room123:matrix.example.com", false},
		{"invalid format", "room123", true},
		{"missing !", "room123:matrix.example.com", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMatrixRoomID(tt.roomID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMatrixRoomID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid https", "https://example.com", false},
		{"valid http", "http://example.com", false},
		{"no scheme", "example.com", true},
		{"invalid scheme", "ftp://example.com", true},
		{"no host", "https://", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateHTTPS(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid https", "https://example.com", false},
		{"http not allowed", "http://example.com", true},
		{"invalid", "not-a-url", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHTTPS(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHTTPS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal text", "hello world", "hello world"},
		{"with html", "<script>alert('xss')</script>", "scriptalertxssscript"},
		{"with quotes", "test'\"text", "testtext"},
		{"with semicolon", "test;DROP TABLE", "testDROP TABLE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeInput(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeInput() = %v, want %v", result, tt.expected)
			}
		})
	}
}
