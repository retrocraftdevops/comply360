package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/comply360/shared/errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	// Default rate limit: 1000 requests per minute per tenant
	defaultRequestsPerMinute = 1000
	defaultWindowDuration    = time.Minute
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int
	WindowDuration    time.Duration
}

// RateLimiter middleware implements rate limiting per tenant using Redis
func RateLimiter(redisClient *redis.Client) gin.HandlerFunc {
	return RateLimiterWithConfig(redisClient, &RateLimitConfig{
		RequestsPerMinute: defaultRequestsPerMinute,
		WindowDuration:    defaultWindowDuration,
	})
}

// RateLimiterWithConfig allows custom rate limiting configuration
func RateLimiterWithConfig(redisClient *redis.Client, config *RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip rate limiting for health check
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

	// Get tenant ID from context (set by TenantMiddleware)
	tenantIDVal, exists := c.Get("tenant_id")
	if !exists {
		// No tenant context, allow request (auth endpoints, etc.)
		c.Next()
		return
	}

	tenantID := fmt.Sprintf("%v", tenantIDVal)
	ctx := context.Background()

		// Create rate limit key
		rateLimitKey := fmt.Sprintf("rate_limit:tenant:%s", tenantID)

		// Get current count
		currentCount, err := redisClient.Get(ctx, rateLimitKey).Int()
		if err != nil && err != redis.Nil {
			// Redis error, allow request but log error
			c.Next()
			return
		}

		// Check if rate limit exceeded
		if currentCount >= config.RequestsPerMinute {
			// Get TTL to inform client when to retry
			ttl, _ := redisClient.TTL(ctx, rateLimitKey).Result()

			c.Header("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerMinute))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))
			c.Header("Retry-After", strconv.FormatInt(int64(ttl.Seconds()), 10))

			c.JSON(http.StatusTooManyRequests, errors.NewAPIErrorWithDetails(
				errors.ErrRateLimitExceeded,
				fmt.Sprintf("Rate limit exceeded. Maximum %d requests per minute.", config.RequestsPerMinute),
				map[string]interface{}{
					"limit":     config.RequestsPerMinute,
					"remaining": 0,
					"reset_at":  time.Now().Add(ttl).Unix(),
				},
			))
			c.Abort()
			return
		}

		// Increment counter
		pipe := redisClient.Pipeline()
		incrCmd := pipe.Incr(ctx, rateLimitKey)

		// Set expiration on first request
		if currentCount == 0 {
			pipe.Expire(ctx, rateLimitKey, config.WindowDuration)
		}

		_, err = pipe.Exec(ctx)
		if err != nil {
			// Redis error, allow request but log error
			c.Next()
			return
		}

		// Get new count
		newCount := int(incrCmd.Val())
		remaining := config.RequestsPerMinute - newCount
		if remaining < 0 {
			remaining = 0
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerMinute))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.WindowDuration).Unix(), 10))

		c.Next()
	}
}
