// Package database tests - unit tests for database operations.
package database

import (
	"context"
	"os"
	"testing"
)

func TestUpsertViberUser(t *testing.T) {
	// Create temporary database
	dbPath := "/tmp/test_bridge.db"
	defer os.Remove(dbPath)
	
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	ctx := context.Background()
	
	// Test upsert
	err = db.UpsertViberUser(ctx, "test_user_1", "Test User")
	if err != nil {
		t.Fatalf("Failed to upsert user: %v", err)
	}
	
	// Test retrieval
	user, err := db.GetViberUser(ctx, "test_user_1")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	
	if user == nil {
		t.Fatal("User not found")
	}
	
	if user.ViberName != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", user.ViberName)
	}
	
	// Test update
	err = db.UpsertViberUser(ctx, "test_user_1", "Updated Name")
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	
	user, err = db.GetViberUser(ctx, "test_user_1")
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}
	
	if user.ViberName != "Updated Name" {
		t.Errorf("Expected updated name 'Updated Name', got '%s'", user.ViberName)
	}
}

func TestRoomMapping(t *testing.T) {
	dbPath := "/tmp/test_bridge_rooms.db"
	defer os.Remove(dbPath)
	
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	ctx := context.Background()
	
	// Test create mapping
	err = db.CreateRoomMapping(ctx, "viber_chat_1", "!matrix_room_1:example.com")
	if err != nil {
		t.Fatalf("Failed to create mapping: %v", err)
	}
	
	// Test retrieval
	matrixRoomID, err := db.GetMatrixRoomID(ctx, "viber_chat_1")
	if err != nil {
		t.Fatalf("Failed to get matrix room id: %v", err)
	}
	
	if matrixRoomID != "!matrix_room_1:example.com" {
		t.Errorf("Expected '!matrix_room_1:example.com', got '%s'", matrixRoomID)
	}
	
	// Test reverse lookup
	viberChatID, err := db.GetViberChatID(ctx, "!matrix_room_1:example.com")
	if err != nil {
		t.Fatalf("Failed to get viber chat id: %v", err)
	}
	
	if viberChatID != "viber_chat_1" {
		t.Errorf("Expected 'viber_chat_1', got '%s'", viberChatID)
	}
}

func TestMessageMapping(t *testing.T) {
	dbPath := "/tmp/test_bridge_messages.db"
	defer os.Remove(dbPath)
	
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	ctx := context.Background()
	
	// Create room mapping first
	err = db.CreateRoomMapping(ctx, "viber_chat_1", "!room:example.com")
	if err != nil {
		t.Fatalf("Failed to create room mapping: %v", err)
	}
	
	// Test store message mapping
	err = db.StoreMessageMapping(ctx, "viber_msg_123", "$matrix_event_456", "viber_chat_1")
	if err != nil {
		t.Fatalf("Failed to store message mapping: %v", err)
	}
	
	// Test retrieval
	matrixEventID, err := db.GetMatrixEventID(ctx, "viber_msg_123")
	if err != nil {
		t.Fatalf("Failed to get matrix event id: %v", err)
	}
	
	if matrixEventID != "$matrix_event_456" {
		t.Errorf("Expected '$matrix_event_456', got '%s'", matrixEventID)
	}
}

func TestGroupMembers(t *testing.T) {
	dbPath := "/tmp/test_bridge_groups.db"
	defer os.Remove(dbPath)
	
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	ctx := context.Background()
	
	// Test upsert group member
	err = db.UpsertGroupMember(ctx, "chat_1", "user_1")
	if err != nil {
		t.Fatalf("Failed to upsert group member: %v", err)
	}
	
	err = db.UpsertGroupMember(ctx, "chat_1", "user_2")
	if err != nil {
		t.Fatalf("Failed to upsert second group member: %v", err)
	}
	
	// Test list members
	members, err := db.ListGroupMembers(ctx, "chat_1")
	if err != nil {
		t.Fatalf("Failed to list group members: %v", err)
	}
	
	if len(members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(members))
	}
	
	// Check members are present
	foundUser1 := false
	foundUser2 := false
	for _, member := range members {
		if member == "user_1" {
			foundUser1 = true
		}
		if member == "user_2" {
			foundUser2 = true
		}
	}
	
	if !foundUser1 || !foundUser2 {
		t.Errorf("Expected to find both users, found user_1: %v, user_2: %v", foundUser1, foundUser2)
	}
}

func TestLinkViberUser(t *testing.T) {
	dbPath := "/tmp/test_bridge_link.db"
	defer os.Remove(dbPath)
	
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	ctx := context.Background()
	
	// Create user first
	err = db.UpsertViberUser(ctx, "viber_user_1", "Viber User")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// Link to Matrix user
	err = db.LinkViberUser(ctx, "viber_user_1", "@matrix_user:example.com")
	if err != nil {
		t.Fatalf("Failed to link user: %v", err)
	}
	
	// Retrieve and verify
	user, err := db.GetViberUser(ctx, "viber_user_1")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	
	if user.MatrixUserID == nil {
		t.Fatal("Matrix user ID not linked")
	}
	
	if *user.MatrixUserID != "@matrix_user:example.com" {
		t.Errorf("Expected '@matrix_user:example.com', got '%s'", *user.MatrixUserID)
	}
}
