package retry_test

import (
	"context"
	"fmt"
	"time"

	"github.com/example/mautrix-viber/internal/retry"
)

func ExampleDefaultConfig() {
	// Get default retry configuration
	cfg := retry.DefaultConfig()

	fmt.Printf("Max attempts: %d\n", cfg.MaxAttempts)
	fmt.Printf("Initial delay: %v\n", cfg.InitialDelay)
	fmt.Printf("Max delay: %v\n", cfg.MaxDelay)
	fmt.Printf("Jitter enabled: %v\n", cfg.Jitter)

	// Output:
	// Max attempts: 3
	// Initial delay: 100ms
	// Max delay: 5s
	// Jitter enabled: true
}

func ExampleDo() {
	// Create a context
	ctx := context.Background()

	// Get default retry configuration
	cfg := retry.DefaultConfig()

	// Attempt counter
	attempts := 0

	// Execute a function with retry logic
	err := retry.Do(ctx, cfg, func() error {
		attempts++
		if attempts < 2 {
			return fmt.Errorf("transient error")
		}
		fmt.Println("Operation succeeded")
		return nil
	})

	if err != nil {
		fmt.Printf("Operation failed after retries: %v\n", err)
	}

	// Output:
	// Operation succeeded
}

func ExampleDo_withTimeout() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Get default retry configuration
	cfg := retry.DefaultConfig()

	// Execute a function that will timeout
	err := retry.Do(ctx, cfg, func() error {
		// Simulate a long-running operation
		time.Sleep(50 * time.Millisecond)
		return fmt.Errorf("operation failed")
	})

	if err != nil {
		fmt.Println("Retry stopped due to context timeout")
	}

	// Output:
	// Retry stopped due to context timeout
}
