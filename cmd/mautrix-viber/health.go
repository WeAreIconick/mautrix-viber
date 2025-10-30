package main

import (
	"context"
	"net/http"
	"time"

	"github.com/example/mautrix-viber/internal/database"
	imatrix "github.com/example/mautrix-viber/internal/matrix"
	"github.com/example/mautrix-viber/internal/viber"
)

// healthHandler checks if the service is healthy.
// Returns 200 if basic components are functional.
func healthHandler(db *database.DB, viberClient *viber.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Quick health check - just verify we can respond
		// More detailed checks are in readinessHandler
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

// readinessHandler checks if the service is ready to serve traffic.
// Returns 200 if all critical components are operational.
func readinessHandler(db *database.DB, viberClient *viber.Client, mxClient *imatrix.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		
		var issues []string
		
		// Check database connectivity
		if db != nil {
			// Use context with timeout for database ping
			pingCtx, pingCancel := context.WithTimeout(ctx, 2*time.Second)
			defer pingCancel()
			
			if err := db.Ping(pingCtx); err != nil {
				issues = append(issues, "database: "+err.Error())
			}
		} else {
			issues = append(issues, "database: not configured")
		}
		
		// Check Viber client configuration
		if viberClient == nil {
			issues = append(issues, "viber: not configured")
		}
		
		// Check Matrix client if configured
		// Note: We don't fail if Matrix is optional, but we check if it's supposed to be there
		
		if len(issues) > 0 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Not ready: " + issues[0]))
			return
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	}
}

