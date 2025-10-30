package main

import (
	"fmt"
    "log"
	"net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/example/mautrix-viber/internal/config"
    imatrix "github.com/example/mautrix-viber/internal/matrix"
	"github.com/example/mautrix-viber/internal/viber"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	fmt.Println("mautrix-viber bootstrap")

    env := config.FromEnv()

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
        log.Printf("matrix config incomplete; message relay will be disabled")
    }

    cfg := viber.Config{
        APIToken:      env.APIToken,
        WebhookURL:    env.WebhookURL,
        ListenAddress: env.ListenAddress,
    }

    v := viber.NewClient(cfg, mxClient)
	if err := v.EnsureWebhook(); err != nil {
		log.Fatalf("failed to ensure webhook: %v", err)
	}

    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
    mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
    mux.Handle("/metrics", promhttp.Handler())
    mux.HandleFunc("/viber/webhook", v.WebhookHandler)

    srv := &http.Server{
        Addr:              cfg.ListenAddress,
        Handler:           withServerMiddleware(withRateLimit(mux)),
        ReadTimeout:       10 * time.Second,
        ReadHeaderTimeout: 5 * time.Second,
        WriteTimeout:      15 * time.Second,
        IdleTimeout:       60 * time.Second,
    }

    go func() {
        log.Printf("listening on %s", cfg.ListenAddress)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()

    // Graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
    <-stop
    shutdownCtx, cancel := time.WithTimeout(time.Background(), 15*time.Second)
    defer cancel()
    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Printf("graceful shutdown failed: %v", err)
    }
    log.Printf("shutdown complete")
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

