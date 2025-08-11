# Redis Rate Limiter

This document describes the Redis-based rate limiting implementation for the Master Data REST API.

## Overview

The Redis Rate Limiter provides comprehensive rate limiting capabilities using Redis as a backend store. It supports multiple rate limiting strategies, configurable limits, and provides monitoring and management endpoints.

## Features

- **Multiple Rate Limiting Strategies**: IP-based, user-based, API key-based, and tiered rate limiting
- **Configurable Limits**: Different rate limits for different endpoints and user types
- **Redis Backend**: Uses Redis for distributed rate limiting across multiple server instances
- **Graceful Degradation**: Falls back to allowing all requests if Redis is unavailable
- **Monitoring**: Built-in endpoints for monitoring rate limit usage and statistics
- **Management**: Endpoints for resetting rate limits and viewing configuration

## Architecture

### Components

1. **Redis Manager** (`internal/adapters/secondary/redis/connection.go`)
   - Manages Redis connections and health checks
   - Handles connection pooling and configuration

2. **Rate Limiter** (`internal/adapters/secondary/redis/rate_limiter.go`)
   - Core rate limiting logic using Redis sorted sets
   - Lua script for atomic operations
   - Configurable time windows and request limits

3. **Middleware** (`internal/adapters/primary/http/middleware/rate_limiter.go`)
   - Fiber middleware for integrating rate limiting
   - Multiple pre-configured rate limiting strategies
   - Customizable identifier extraction and response handling

4. **HTTP Handler** (`internal/adapters/primary/http/rate_limit_handler.go`)
   - REST endpoints for rate limit management
   - Monitoring and statistics endpoints
   - Configuration and reset functionality

## Configuration

### Environment Variables

```bash
# Redis Configuration
REDIS_ENABLED=true
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=10
REDIS_DIAL_TIMEOUT=5s
REDIS_READ_TIMEOUT=3s
REDIS_WRITE_TIMEOUT=3s
```

### Rate Limiting Configuration

The system supports multiple rate limiting configurations:

#### Global Rate Limiting
- **Default**: 100 requests per minute per IP
- **Key**: `api`

#### Tiered Rate Limiting
- **Geographic Data** (`/api/v1/countries`, `/api/v1/provinces`): 1000 requests/minute
- **Financial Data** (`/api/v1/banks`, `/api/v1/currencies`): 500 requests/minute
- **Language Data** (`/api/v1/languages`): 200 requests/minute

#### Identifier-Based Rate Limiting
- **IP Address**: 1000 requests/minute
- **User ID**: 500 requests/minute
- **API Key**: 2000 requests/minute

## Usage

### Basic Rate Limiting

The rate limiter is automatically applied to all API routes when Redis is enabled:

```go
// In your router setup
if redisManager.IsEnabled() {
    app.Use(middleware.TieredRateLimiter(rateLimiter))
}
```

### Custom Rate Limiting

You can create custom rate limiters for specific use cases:

```go
// IP-based rate limiter
ipLimiter := middleware.IPBasedRateLimiter(rateLimiter, 100, time.Minute)

// User-based rate limiter
userLimiter := middleware.UserBasedRateLimiter(rateLimiter, 50, time.Minute)

// API key-based rate limiter
apiKeyLimiter := middleware.APIKeyBasedRateLimiter(rateLimiter, 200, time.Minute)
```

### Custom Configuration

```go
config := &middleware.RateLimiterConfig{
    Requests:        500,
    Window:          time.Minute,
    Key:            "custom",
    IncludeHeaders: true,
    IdentifierFunc: func(c *fiber.Ctx) string {
        // Custom identifier logic
        return c.Get("X-Custom-ID")
    },
    LimitExceededHandler: func(c *fiber.Ctx, result *redis.RateLimitResult) error {
        // Custom response for rate limit exceeded
        return c.Status(429).JSON(fiber.Map{
            "error": "Custom rate limit message",
            "retry_after": result.RetryAfter.Seconds(),
        })
    },
}

app.Use(middleware.RateLimiter(rateLimiter, config))
```

## API Endpoints

### Rate Limit Information

```
GET /api/v1/rate-limit/info?identifier={identifier}
```

Returns current rate limit status for the specified identifier (defaults to client IP).

**Response:**
```json
{
  "success": true,
  "message": "Rate limit information retrieved successfully",
  "data": {
    "api": {
      "allowed": true,
      "remaining": 95,
      "reset_time": "2024-01-15T10:30:00Z",
      "limit": 100,
      "window": "1m"
    },
    "ip": {
      "allowed": true,
      "remaining": 995,
      "reset_time": "2024-01-15T10:30:00Z",
      "limit": 1000,
      "window": "1m"
    }
  }
}
```

### Rate Limit Statistics

```
GET /api/v1/rate-limit/stats?key={key}
```

Returns comprehensive statistics for rate limiting across all clients.

**Response:**
```json
{
  "success": true,
  "message": "Rate limit statistics retrieved successfully",
  "data": {
    "enabled": true,
    "total_keys": 150,
    "window_size": "1m0s",
    "max_requests": 100,
    "active_limits": {
      "192.168.1.1": {
        "remaining": 95,
        "reset_time": "2024-01-15T10:30:00Z",
        "allowed": true
      }
    }
  }
}
```

### Rate Limit Configuration

```
GET /api/v1/rate-limit/config
```

Returns the current rate limiting configuration.

**Response:**
```json
{
  "success": true,
  "message": "Rate limit configuration retrieved successfully",
  "data": {
    "global": {
      "requests": 100,
      "window": "1m",
      "key": "api"
    },
    "tiered": {
      "geographic": {
        "endpoints": ["/api/v1/countries", "/api/v1/provinces"],
        "requests": 1000,
        "window": "1m",
        "key": "geo"
      }
    }
  }
}
```

### Reset Rate Limit

```
POST /api/v1/rate-limit/reset
```

Resets the rate limit for a specific identifier.

**Request Body:**
```json
{
  "identifier": "192.168.1.1",
  "key": "api"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Rate limit reset successfully",
  "data": {
    "identifier": "192.168.1.1",
    "key": "api",
    "reset_time": "2024-01-15T10:25:00Z"
  }
}
```

## Response Headers

When rate limiting is enabled, the following headers are added to responses:

- `X-RateLimit-Limit`: Maximum requests allowed in the time window
- `X-RateLimit-Remaining`: Remaining requests in the current window
- `X-RateLimit-Reset`: Unix timestamp when the window resets
- `X-RateLimit-RetryAfter`: Seconds to wait before retrying (only when rate limited)

## Rate Limit Exceeded Response

When a rate limit is exceeded, the API returns a `429 Too Many Requests` response:

```json
{
  "success": false,
  "error": "Rate limit exceeded. Try again in 45 seconds",
  "data": {
    "retry_after": 45,
    "reset_time": "2024-01-15T10:30:00Z"
  }
}
```

## Implementation Details

### Redis Data Structure

Rate limits are stored in Redis using sorted sets with the following key pattern:
```
rate_limit:{key}:{identifier}
```

- **Key**: Rate limit type (e.g., `api`, `ip`, `user`)
- **Identifier**: Client identifier (e.g., IP address, user ID, API key)
- **Score**: Timestamp of the request
- **Member**: Timestamp of the request (same as score)

### Lua Script

The rate limiting logic uses a Lua script for atomic operations:

```lua
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
```

### Time Windows

Rate limits use sliding time windows. The window is calculated by truncating the current time to the specified duration:

```go
windowStart := now.Truncate(config.Window)
```

This ensures consistent window boundaries across all requests.

## Testing

### Manual Testing

Use the provided test script to verify rate limiting functionality:

```bash
# Make the script executable
chmod +x scripts/test-rate-limiter.sh

# Run all tests
./scripts/test-rate-limiter.sh test

# Test specific functionality
./scripts/test-rate-limiter.sh limit
./scripts/test-rate-limiter.sh headers
```

### Testing with curl

```bash
# Test rate limiting
for i in {1..15}; do
  curl -H "Authorization: Bearer your-api-key" \
       http://localhost:8080/api/v1/countries
  sleep 0.1
done

# Check rate limit info
curl -H "Authorization: Bearer your-api-key" \
     http://localhost:8080/api/v1/rate-limit/info

# Check rate limit stats
curl -H "Authorization: Bearer your-api-key" \
     http://localhost:8080/api/v1/rate-limit/stats
```

## Monitoring and Alerting

### Metrics

The rate limiter provides several metrics for monitoring:

- **Total active rate limit keys**
- **Rate limit violations per time period**
- **Average response time impact**
- **Redis connection health**

### Health Checks

Redis health is checked on startup and can be monitored via:

```go
err := redisManager.HealthCheck(ctx)
if err != nil {
    log.WithError(err).Error("Redis health check failed")
}
```

### Logging

Rate limiting events are logged with structured logging:

```go
log.WithFields(logrus.Fields{
    "identifier": identifier,
    "key":        config.Key,
    "allowed":    result.Allowed,
    "remaining":  result.Remaining,
}).Info("Rate limit check completed")
```

## Performance Considerations

### Redis Performance

- **Connection Pooling**: Configurable pool size for optimal performance
- **Lua Scripts**: Atomic operations reduce Redis round trips
- **Key Expiration**: Automatic cleanup of expired rate limit data
- **Memory Usage**: Sorted sets provide efficient storage and retrieval

### API Performance

- **Middleware Order**: Rate limiting is applied early in the middleware chain
- **Graceful Degradation**: API continues to function if Redis is unavailable
- **Minimal Overhead**: Rate limiting adds minimal latency to requests
- **Async Operations**: Redis operations are non-blocking

## Security Considerations

### Rate Limit Bypass Prevention

- **Multiple Identifier Types**: Prevents bypassing limits by changing identifiers
- **IP Address Validation**: Ensures IP addresses are properly extracted
- **API Key Validation**: Validates API keys before applying rate limits

### Data Privacy

- **Identifier Hashing**: Consider hashing sensitive identifiers
- **Data Retention**: Rate limit data expires automatically
- **Access Control**: Rate limit endpoints require authentication

## Troubleshooting

### Common Issues

1. **Rate Limiting Not Working**
   - Check if Redis is enabled and running
   - Verify Redis connection in logs
   - Check rate limiter middleware is applied

2. **Unexpected Rate Limiting**
   - Verify identifier extraction logic
   - Check rate limit configuration
   - Review Redis data structure

3. **Performance Issues**
   - Monitor Redis memory usage
   - Check connection pool settings
   - Review Lua script performance

### Debug Mode

Enable debug logging for detailed rate limiting information:

```bash
LOG_LEVEL=debug
```

### Redis CLI Debugging

```bash
# Connect to Redis
redis-cli

# Check rate limit keys
KEYS "rate_limit:*"

# Inspect specific rate limit
ZRANGE rate_limit:api:192.168.1.1 0 -1 WITHSCORES

# Check key expiration
TTL rate_limit:api:192.168.1.1
```

## Future Enhancements

### Planned Features

- **Dynamic Rate Limiting**: Adjust limits based on user tier or subscription
- **Geographic Rate Limiting**: Different limits based on client location
- **Rate Limit Analytics**: Detailed analytics and reporting
- **Webhook Notifications**: Notify external systems of rate limit violations
- **Machine Learning**: Adaptive rate limiting based on usage patterns

### Integration Opportunities

- **Prometheus Metrics**: Export rate limiting metrics
- **Grafana Dashboards**: Visualize rate limiting data
- **Alert Manager**: Alert on rate limit violations
- **External Caching**: Integrate with CDN rate limiting

## Conclusion

The Redis Rate Limiter provides a robust, scalable solution for API rate limiting. It offers flexibility in configuration, comprehensive monitoring capabilities, and graceful degradation when Redis is unavailable. The implementation follows best practices for performance, security, and maintainability.

For questions or issues, please refer to the project documentation or create an issue in the project repository.
