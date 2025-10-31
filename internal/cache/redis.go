// Package cache provides Redis cache for frequently accessed user/room mappings.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache provides Redis caching functionality.
type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewCache creates a new Redis cache.
func NewCache(redisURL string, ttl time.Duration) (*Cache, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Cache{
		client: client,
		ttl:    ttl,
	}, nil
}

// Get retrieves a value from cache.
func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	if c == nil || c.client == nil {
		return nil, fmt.Errorf("cache not configured")
	}

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Key not found
	}
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}

	return []byte(val), nil
}

// Set stores a value in cache.
func (c *Cache) Set(ctx context.Context, key string, value []byte) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("cache not configured")
	}

	return c.client.Set(ctx, key, value, c.ttl).Err()
}

// Delete removes a key from cache.
func (c *Cache) Delete(ctx context.Context, key string) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("cache not configured")
	}

	return c.client.Del(ctx, key).Err()
}

// GetJSON retrieves and unmarshals JSON from cache.
func (c *Cache) GetJSON(ctx context.Context, key string, v interface{}) error {
	data, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	if data == nil {
		return fmt.Errorf("key not found")
	}

	return json.Unmarshal(data, v)
}

// SetJSON marshals and stores JSON in cache.
func (c *Cache) SetJSON(ctx context.Context, key string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	return c.Set(ctx, key, data)
}

// Close closes the cache connection.
func (c *Cache) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}
