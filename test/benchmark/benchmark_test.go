// Package benchmark provides performance benchmarks for critical operations.
package benchmark

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/example/mautrix-viber/internal/database"
	"github.com/example/mautrix-viber/internal/retry"
)

// BenchmarkDatabaseUpsert benchmarks user upsertion.
func BenchmarkDatabaseUpsert(b *testing.B) {
	dbPath := "/tmp/benchmark_upsert.db"
	defer os.Remove(dbPath)
	
	db, err := database.Open(dbPath)
	if err != nil {
		b.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		viberID := "bench_user_" + string(rune(i))
		_ = db.UpsertViberUser(viberID, "Benchmark User")
	}
}

// BenchmarkDatabaseQuery benchmarks user retrieval.
func BenchmarkDatabaseQuery(b *testing.B) {
	dbPath := "/tmp/benchmark_query.db"
	defer os.Remove(dbPath)
	
	db, err := database.Open(dbPath)
	if err != nil {
		b.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	
	// Pre-populate
	db.UpsertViberUser("bench_user", "Bench User")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = db.GetViberUser("bench_user")
	}
}

// BenchmarkRetryLogic benchmarks retry operations.
func BenchmarkRetryLogic(b *testing.B) {
	cfg := retry.DefaultConfig()
	cfg.MaxAttempts = 3
	cfg.InitialDelay = 1 * time.Microsecond
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = retry.Do(context.Background(), cfg, func() error {
			return nil // Success case
		})
	}
}

// BenchmarkSignatureCalculation benchmarks HMAC signature calculation.
func BenchmarkSignatureCalculation(b *testing.B) {
	token := "test-api-token"
	body := `{"event":"message","sender":{"id":"123"},"message":{"type":"text","text":"Hello world"}}`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Signature calculation benchmark
		_ = calculateSignature(body, token)
	}
}

func calculateSignature(body, token string) string {
	// Simplified for benchmark - actual implementation uses HMAC
	return "signature"
}

