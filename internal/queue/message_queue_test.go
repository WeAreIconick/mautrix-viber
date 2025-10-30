// Package queue tests - unit tests for message queue.
package queue

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestQueue_Enqueue(t *testing.T) {
	q := NewQueue(1, func(ctx context.Context, job MessageJob) error {
		return nil
	})
	
	job := MessageJob{
		ID:        "test-job-1",
		Type:      "test",
		MaxRetries: 3,
		CreatedAt: time.Now(),
	}
	
	q.Enqueue(job)
	
	// Give worker time to process
	time.Sleep(50 * time.Millisecond)
	
	// Queue should be empty or processing
	length := q.Length()
	if length > 1 {
		t.Errorf("Expected queue length <= 1, got %d", length)
	}
}

func TestQueue_Retry(t *testing.T) {
	attempts := 0
	q := NewQueue(1, func(ctx context.Context, job MessageJob) error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary error")
		}
		return nil
	})
	
	job := MessageJob{
		ID:        "test-job-2",
		Type:      "test",
		MaxRetries: 3,
		CreatedAt: time.Now(),
	}
	
	q.Enqueue(job)
	
	// Wait for processing and retry
	time.Sleep(200 * time.Millisecond)
	
	if attempts < 2 {
		t.Errorf("Expected at least 2 attempts, got %d", attempts)
	}
}

func TestQueue_Length(t *testing.T) {
	q := NewQueue(1, func(ctx context.Context, job MessageJob) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	
	// Enqueue multiple jobs
	for i := 0; i < 5; i++ {
		q.Enqueue(MessageJob{
			ID:        "job-" + string(rune(i)),
			Type:      "test",
			MaxRetries: 1,
			CreatedAt: time.Now(),
		})
	}
	
	// Check initial length
	length := q.Length()
	if length != 5 {
		t.Errorf("Expected queue length 5, got %d", length)
	}
	
	// Wait for processing
	time.Sleep(100 * time.Millisecond)
	
	// Should be processing
	finalLength := q.Length()
	if finalLength > 5 {
		t.Errorf("Queue length should decrease, got %d", finalLength)
	}
}

