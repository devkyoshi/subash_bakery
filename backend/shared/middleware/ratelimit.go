package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/yourusername/erp-system/shared/utils"
)

type RateLimiter struct {
	redisClient *redis.Client
	limit       int
	window      time.Duration
}

func NewRateLimiter(redisClient *redis.Client, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		limit:       limit,
		window:      window,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get identifier (user ID or IP address)
		identifier := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			identifier = fmt.Sprintf("user:%s", userID)
		}

		key := fmt.Sprintf("ratelimit:%s", identifier)
		ctx := context.Background()

		// Increment counter
		count, err := rl.redisClient.Incr(ctx, key).Result()
		if err != nil {
			// If Redis fails, allow the request
			c.Next()
			return
		}

		// Set expiry on first request
		if count == 1 {
			rl.redisClient.Expire(ctx, key, rl.window)
		}

		// Check if limit exceeded
		if count > int64(rl.limit) {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Too many requests", nil)
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.limit-int(count)))

		c.Next()
	}
}
