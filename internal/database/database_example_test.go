package database_test

import (
	"context"
	"fmt"
	"github.com/example/mautrix-viber/internal/database"
	"os"
)

func ExampleOpen() {
	// Open or create a SQLite database
	dbPath := "/tmp/example_bridge.db"
	db, err := database.Open(dbPath)
	if err != nil {
		fmt.Printf("Failed to open database: %v\n", err)
		return
	}
	defer func() { _ = db.Close() }()
	defer func() { _ = os.Remove(dbPath) }() // Cleanup

	// Database is ready to use
	fmt.Println("Database opened successfully")
	// Output:
	// Database opened successfully
}

func ExampleDB_UpsertViberUser() {
	// Open database
	dbPath := "/tmp/example_users.db"
	db, err := database.Open(dbPath)
	if err != nil {
		fmt.Printf("Failed to open database: %v\n", err)
		return
	}
	defer func() { _ = db.Close() }()
	defer func() { _ = os.Remove(dbPath) }() // Cleanup

	// Upsert a Viber user
	ctx := context.Background()
	err = db.UpsertViberUser(ctx, "user123", "Alice")
	if err != nil {
		fmt.Printf("Failed to upsert user: %v\n", err)
		return
	}

	// Retrieve the user
	user, err := db.GetViberUser(ctx, "user123")
	if err != nil {
		fmt.Printf("Failed to get user: %v\n", err)
		return
	}

	if user != nil {
		fmt.Printf("User: %s (%s)\n", user.ViberName, user.ViberID)
	}
	// Output:
	// User: Alice (user123)
}
