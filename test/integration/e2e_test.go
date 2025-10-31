// Package integration e2e tests - end-to-end integration tests.
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/example/mautrix-viber/internal/database"
)

// TestEndToEndMessageFlow tests the complete message flow from Viber to Matrix.
func TestEndToEndMessageFlow(t *testing.T) {
	t.Skip("Requires full test environment with mock Matrix client")

	// This would test:
	// 1. Viber webhook receives message
	// 2. Signature verification passes
	// 3. Message stored in database
	// 4. Message forwarded to Matrix
	// 5. Message mapping stored
	// 6. Matrix event ID retrieved correctly
}

// TestMatrixToViberFlow tests message flow from Matrix to Viber.
func TestMatrixToViberFlow(t *testing.T) {
	t.Skip("Requires mock Matrix client and Viber API")

	// This would test:
	// 1. Matrix event received
	// 2. Message formatted for Viber
	// 3. Message sent via Viber API
	// 4. Message mapping stored
}

// TestDatabaseConsistency tests database operations maintain consistency.
func TestDatabaseConsistency(t *testing.T) {
	dbPath := "/tmp/test_e2e.db"
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer func() { _ = db.Close() }()

	// Test transaction-like consistency
	viberID := "test_user_e2e"
	matrixRoomID := "!test_room:example.com"
	viberChatID := "test_chat_123"

	// Create user
	ctx := context.Background()
	if err := db.UpsertViberUser(ctx, viberID, "Test User"); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create room mapping
	if err := db.CreateRoomMapping(ctx, viberChatID, matrixRoomID); err != nil {
		t.Fatalf("Failed to create room mapping: %v", err)
	}

	// Verify consistency
	user, err := db.GetViberUser(ctx, viberID)
	if err != nil || user == nil {
		t.Fatal("User not found after creation")
	}

	retrievedRoomID, err := db.GetMatrixRoomID(ctx, viberChatID)
	if err != nil || retrievedRoomID != matrixRoomID {
		t.Fatalf("Room mapping inconsistent: expected %s, got %s", matrixRoomID, retrievedRoomID)
	}
}

// TestConcurrentOperations tests concurrent database operations.
func TestConcurrentOperations(t *testing.T) {
	dbPath := "/tmp/test_concurrent.db"
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Simulate concurrent writes
	ctx := context.Background()
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			viberID := "test_user_" + string(rune(id))
			_ = db.UpsertViberUser(ctx, viberID, "Test User")
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all users were created
	for i := 0; i < 10; i++ {
		viberID := "test_user_" + string(rune(i))
		user, err := db.GetViberUser(ctx, viberID)
		if err != nil || user == nil {
			t.Errorf("User %s not found after concurrent creation", viberID)
		}
	}
}

// TestGracefulShutdown tests graceful shutdown handling.
func TestGracefulShutdown(t *testing.T) {
	t.Skip("Requires actual server instance")

	// This would test:
	// 1. Server starts successfully
	// 2. SIGTERM received
	// 3. Shutdown context respected
	// 4. Connections closed gracefully
	// 5. No hanging goroutines
}

// TestMessageDeduplication tests that duplicate messages are handled correctly.
func TestMessageDeduplication(t *testing.T) {
	dbPath := "/tmp/test_dedup.db"
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	viberMsgID := "msg_123"
	matrixEventID1 := "$event_1"
	matrixEventID2 := "$event_2"
	chatID := "chat_123"
	ctx := context.Background()

	// Create room mapping first
	if err := db.CreateRoomMapping(ctx, chatID, "!room:example.com"); err != nil {
		t.Fatalf("Failed to create room mapping: %v", err)
	}

	// Store first mapping
	if err := db.StoreMessageMapping(ctx, viberMsgID, matrixEventID1, chatID); err != nil {
		t.Fatalf("Failed to store first mapping: %v", err)
	}

	// Try to store duplicate (should update, not fail)
	if err := db.StoreMessageMapping(ctx, viberMsgID, matrixEventID2, chatID); err != nil {
		t.Fatalf("Failed to update mapping: %v", err)
	}

	// Verify latest mapping is stored
	retrieved, err := db.GetMatrixEventID(ctx, viberMsgID)
	if err != nil {
		t.Fatalf("Failed to retrieve mapping: %v", err)
	}

	if retrieved != matrixEventID2 {
		t.Errorf("Expected %s, got %s", matrixEventID2, retrieved)
	}
}

// TestContextCancellation tests that operations respect context cancellation.
func TestContextCancellation(t *testing.T) {
	dbPath := "/tmp/test_context.db"
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Small delay to ensure context is cancelled
	time.Sleep(10 * time.Millisecond)

	// Ping should respect context
	err = db.Ping(ctx)
	if err == nil {
		t.Error("Expected error from cancelled context")
	}
}
