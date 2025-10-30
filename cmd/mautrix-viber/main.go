package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/example/mautrix-viber/internal/viber"
)

func main() {
	fmt.Println("mautrix-viber bootstrap")

	cfg := viber.Config{
		APIToken:      "",
		WebhookURL:    "",
		ListenAddress: ":8080",
	}

	v := viber.NewClient(cfg)
	if err := v.EnsureWebhook(); err != nil {
		log.Fatalf("failed to ensure webhook: %v", err)
	}

	http.HandleFunc("/viber/webhook", v.WebhookHandler)
	log.Printf("listening on %s", cfg.ListenAddress)
	log.Fatal(http.ListenAndServe(cfg.ListenAddress, nil))
}
