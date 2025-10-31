// Package queue provides message queue for reliable message delivery with retry logic.
package queue

import (
	"context"
	"sync"
	"time"
)

// MessageJob represents a message delivery job.
type MessageJob struct {
	ID        string
	Type      string // "viber_to_matrix", "matrix_to_viber"
	Payload   interface{}
	Retries   int
	MaxRetries int
	CreatedAt time.Time
}

// Queue represents a message queue.
type Queue struct {
	mu      sync.Mutex
	jobs    []MessageJob
	workers int
	handler func(ctx context.Context, job MessageJob) error
}

// NewQueue creates a new message queue.
func NewQueue(workers int, handler func(ctx context.Context, job MessageJob) error) *Queue {
	q := &Queue{
		workers: workers,
		handler: handler,
		jobs:    make([]MessageJob, 0),
	}
	
	// Start workers
	for i := 0; i < workers; i++ {
		go q.worker(i)
	}
	
	return q
}

// Enqueue adds a job to the queue.
func (q *Queue) Enqueue(job MessageJob) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs = append(q.jobs, job)
}

// worker processes jobs from the queue.
func (q *Queue) worker(id int) {
	for {
		q.mu.Lock()
		if len(q.jobs) == 0 {
			q.mu.Unlock()
			time.Sleep(100 * time.Millisecond)
			continue
		}
		
		job := q.jobs[0]
		q.jobs = q.jobs[1:]
		q.mu.Unlock()
		
		// Process job
		ctx := context.Background()
		err := q.handler(ctx, job)
		
		if err != nil && job.Retries < job.MaxRetries {
			// Retry with exponential backoff
			job.Retries++
			delay := time.Duration(job.Retries) * 100 * time.Millisecond
			time.Sleep(delay)
			q.Enqueue(job)
		}
	}
}

// Length returns the current queue length.
func (q *Queue) Length() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.jobs)
}

