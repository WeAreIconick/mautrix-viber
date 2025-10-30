package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"github.com/example/mautrix-viber/internal/api"
	"github.com/example/mautrix-viber/internal/config"
	"github.com/example/mautrix-viber/internal/database"
	"github.com/example/mautrix-viber/internal/logger"
	imatrix "github.com/example/mautrix-viber/internal/matrix"
	"github.com/example/mautrix-viber/internal/middleware"
	"github.com/example/mautrix-viber/internal/viber"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger.Info("mautrix-viber starting",
		"version", "dev",
	)

	env := config.FromEnv()
	
	// Validate configuration before proceeding
	if err := env.Validate(); err != nil {
		log.Fatalf("configuration validation failed: %v", err)
	}

    // Init Matrix client (optional: only if fully configured)
    var mxClient *imatrix.Client
    if env.MatrixHomeserverURL != "" && env.MatrixAccessToken != "" && env.MatrixDefaultRoomID != "" {
        mc, err := imatrix.NewClient(imatrix.Config{
            HomeserverURL: env.MatrixHomeserverURL,
            AccessToken:   env.MatrixAccessToken,
            DefaultRoomID: env.MatrixDefaultRoomID,
        })
        if err != nil {
            log.Fatalf("failed to initialize matrix client: %v", err)
        }
        mxClient = mc
    } else {
        logger.Info("matrix config incomplete; message relay will be disabled")
    }

    // Open database
    db, err := database.Open(env.DatabasePath)
    if err != nil {
        log.Fatalf("failed to open database: %v", err)
    }
    defer db.Close()

    cfg := viber.Config{
        APIToken:      env.APIToken,
        WebhookURL:    env.WebhookURL,
        ListenAddress: env.ListenAddress,
    }

    v := viber.NewClient(cfg, mxClient, db)
	if err := v.EnsureWebhook(); err != nil {
		log.Fatalf("failed to ensure webhook: %v", err)
	}

    // If Matrix is configured, start listener to forward Matrix -> Viber
    if mxClient != nil && env.ViberDefaultReceiverID != "" {
        if err := mxClient.StartMessageListener(context.Background(), func(ctx context.Context, msg *event.MessageEventContent, roomID id.RoomID, sender id.UserID) {
            // Forward plain text messages to a default Viber receiver for demo purposes
            if msg.MsgType == event.MsgText || msg.MsgType == event.MsgNotice {
                // best-effort: include sender localpart
                text := msg.Body
                if text != "" {
                    _, err := v.SendText(ctx, env.ViberDefaultReceiverID, text)
                    if err != nil {
                        logger.Warn("failed to forward message to Viber",
                            "error", err,
                            "receiver", env.ViberDefaultReceiverID,
                        )
                    }
                }
            }
        }); err != nil {
            logger.Error("matrix listener error",
                "error", err,
            )
        }
    }

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler(db, v))
	mux.HandleFunc("/readyz", readinessHandler(db, v, mxClient))
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/api/info", api.InfoHandler)
	mux.HandleFunc("/viber/webhook", v.WebhookHandler)

	// Build middleware chain: recovery -> request ID -> logging -> rate limit -> body size
	middlewareChain := withServerMiddleware(
		withRateLimit(
			middleware.LoggingMiddleware(
				middleware.RequestIDMiddleware(
					middleware.RecoveryMiddleware(mux),
				),
			),
		),
	)

	srv := &http.Server{
		Addr:              cfg.ListenAddress,
		Handler:           middlewareChain,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB max header size
	}

    go func() {
        logger.Info("http server listening",
            "address", cfg.ListenAddress,
        )
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("server error",
                "error", err,
            )
            log.Fatalf("server error: %v", err)
        }
    }()

    // Graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
    <-stop
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
    if err := srv.Shutdown(shutdownCtx); err != nil {
        logger.Error("graceful shutdown failed",
            "error", err,
        )
    }
    logger.Info("shutdown complete")
}

// withServerMiddleware adds basic defenses like body size limits.
func withServerMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Limit body size to 2MB by default (sufficient for webhook JSON)
        r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
        next.ServeHTTP(w, r)
    })
}

// withRateLimit applies a simple token-bucket rate limiter per client IP.
func withRateLimit(next http.Handler) http.Handler {
    limiter := newIPRateLimiter(5, 10) // 5 req/sec, burst 10
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := clientIP(r)
        if !limiter.Allow(ip) {
            http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}

