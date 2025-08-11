package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// RateLimiter implements rate limiting using Redis
type RateLimiter struct {
	client *redis.Client
	log    *logrus.Logger
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Requests int           // Number of requests allowed
	Window   time.Duration // Time window for the limit
	Key      string        // Redis key prefix for this limit
}

// RateLimitResult contains the result of a rate limit check
type RateLimitResult struct {
	Allowed    bool          // Whether the request is allowed
	Remaining  int           // Remaining requests in the current window
	ResetTime  time.Time     // When the window resets
	RetryAfter time.Duration // How long to wait before retrying
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(client *redis.Client, log *logrus.Logger) *RateLimiter {
	return &RateLimiter{
		client: client,
		log:    log,
	}
}

// CheckRateLimit checks if a request is allowed based on the rate limit configuration
func (r *RateLimiter) CheckRateLimit(ctx context.Context, identifier string, config RateLimitConfig) (*RateLimitResult, error) {
	if r.client == nil {
		// If Redis is not available, allow all requests
		r.log.Warn("Redis not available, allowing all requests")
		return &RateLimitResult{
			Allowed:   true,
			Remaining: config.Requests,
		}, nil
	}

	// Create a unique key for this identifier and rate limit configuration
	key := fmt.Sprintf("rate_limit:%s:%s", config.Key, identifier)

	// Get current timestamp
	now := time.Now()
	windowStart := now.Truncate(config.Window)

	// Use Lua script for atomic operation
	script := `
		local key = KEYS[1]
		local window_start = tonumber(ARGV[1])
		local window_end = tonumber(ARGV[2])
		local max_requests = tonumber(ARGV[3])
		
		-- Remove expired entries
		redis.call('ZREMRANGEBYSCORE', key, 0, window_start - 1)
		
		-- Count current requests in window
		local current_requests = redis.call('ZCARD', key)
		
		-- Check if limit exceeded
		if current_requests >= max_requests then
			return {0, current_requests, window_end}
		end
		
		-- Add current request
		redis.call('ZADD', key, now(), now())
		redis.call('EXPIRE', key, window_end - now())
		
		return {1, max_requests - current_requests - 1, window_end}
	`

	// Execute Lua script
	result, err := r.client.Eval(ctx, script, []string{key},
		windowStart.Unix(),
		windowStart.Add(config.Window).Unix(),
		config.Requests).Result()

	if err != nil {
		r.log.WithError(err).Error("Failed to execute rate limit script")
		return nil, fmt.Errorf("failed to check rate limit: %w", err)
	}

	// Parse result
	results := result.([]interface{})
	allowed := results[0].(int64) == 1
	remaining := int(results[1].(int64))
	resetTime := time.Unix(results[2].(int64), 0)

	var retryAfter time.Duration
	if !allowed {
		retryAfter = resetTime.Sub(now)
	}

	return &RateLimitResult{
		Allowed:    allowed,
		Remaining:  remaining,
		ResetTime:  resetTime,
		RetryAfter: retryAfter,
	}, nil
}

// GetRateLimitInfo gets information about the current rate limit status
func (r *RateLimiter) GetRateLimitInfo(ctx context.Context, identifier string, config RateLimitConfig) (*RateLimitResult, error) {
	if r.client == nil {
		return &RateLimitResult{
			Allowed:   true,
			Remaining: config.Requests,
		}, nil
	}

	key := fmt.Sprintf("rate_limit:%s:%s", config.Key, identifier)
	now := time.Now()
	windowStart := now.Truncate(config.Window)

	// Count current requests in window
	currentRequests, err := r.client.ZCount(ctx, key,
		"-inf",
		"+inf").Result()

	if err != nil {
		return nil, fmt.Errorf("failed to get rate limit info: %w", err)
	}

	remaining := config.Requests - int(currentRequests)
	if remaining < 0 {
		remaining = 0
	}

	return &RateLimitResult{
		Allowed:   remaining > 0,
		Remaining: remaining,
		ResetTime: windowStart.Add(config.Window),
	}, nil
}

// ResetRateLimit resets the rate limit for a specific identifier
func (r *RateLimiter) ResetRateLimit(ctx context.Context, identifier string, config RateLimitConfig) error {
	if r.client == nil {
		return nil
	}

	key := fmt.Sprintf("rate_limit:%s:%s", config.Key, identifier)
	return r.client.Del(ctx, key).Err()
}

// GetRateLimitStats gets statistics about rate limiting
func (r *RateLimiter) GetRateLimitStats(ctx context.Context, config RateLimitConfig) (map[string]interface{}, error) {
	if r.client == nil {
		return map[string]interface{}{
			"enabled": false,
			"message": "Redis not available",
		}, nil
	}

	pattern := fmt.Sprintf("rate_limit:%s:*", config.Key)
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get rate limit keys: %w", err)
	}

	stats := map[string]interface{}{
		"enabled":       true,
		"total_keys":    len(keys),
		"window_size":   config.Window.String(),
		"max_requests":  config.Requests,
		"active_limits": make(map[string]interface{}),
	}

	// Get details for each active rate limit
	for _, key := range keys {
		identifier := key[len(fmt.Sprintf("rate_limit:%s:", config.Key)):]
		info, err := r.GetRateLimitInfo(ctx, identifier, config)
		if err != nil {
			r.log.WithError(err).WithField("key", key).Warn("Failed to get rate limit info")
			continue
		}

		stats["active_limits"].(map[string]interface{})[identifier] = map[string]interface{}{
			"remaining":  info.Remaining,
			"reset_time": info.ResetTime,
			"allowed":    info.Allowed,
		}
	}

	return stats, nil
}
