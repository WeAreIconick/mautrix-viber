// Package main provides a basic example of using the mautrix-viber bridge.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/example/mautrix-viber/internal/config"
	"github.com/example/mautrix-viber/internal/database"
	"github.com/example/mautrix-viber/internal/matrix"
	"github.com/example/mautrix-viber/internal/viber"
)

func main() {
	// Load configuration
	cfg := config.FromEnv()
	
	// Open database
	db, err := database.Open(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	// Initialize Matrix client
	mxClient, err := matrix.NewClient(matrix.Config{
		HomeserverURL: cfg.MatrixHomeserverURL,
		AccessToken:   cfg.MatrixAccessToken,
		DefaultRoomID: cfg.MatrixDefaultRoomID,
	})
	if err != nil {
		log.Fatalf("Failed to initialize Matrix client: %v", err)
	}
	
	// Initialize Viber client
	viberClient := viber.NewClient(viber.Config{
		APIToken:   cfg.APIToken,
		WebhookURL: cfg.WebhookURL,
	}, mxClient, db)
	
	// Ensure webhook is registered
	if err := viberClient.EnsureWebhook(); err != nil {
		log.Fatalf("Failed to register webhook: %v", err)
	}
	
	fmt.Println("Bridge initialized successfully!")
	fmt.Println("Waiting for messages...")
	
	// Keep running
	ctx := context.Background()
	<-ctx.Done()
	time.Sleep(time.Second)
}

