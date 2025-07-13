package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

// RateLimiterPolicy defines one policy.
type RateLimiterPolicy struct {
	Name     string
	Requests int
	Window   time.Duration
}

// RedisRateLimiter is a reusable factory for multiple policies.
type RedisRateLimiter struct {
	RedisClient *redis.Client
	Policies    map[string]RateLimiterPolicy
}

// NewRedisRateLimiter sets up reusable limiter with predefined policies.
func NewRateLimiter(redisClient *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{
		RedisClient: redisClient,
		Policies: map[string]RateLimiterPolicy{
			"aggressive": {
				Name:     "aggressive",
				Requests: 5,
				Window:   30 * time.Second,
			},
			"moderate": {
				Name:     "moderate",
				Requests: 20,
				Window:   1 * time.Minute,
			},
			"relaxed": {
				Name:     "relaxed",
				Requests: 100,
				Window:   5 * time.Minute,
			},
		},
	}
}

// Handler returns a gin.HandlerFunc for the given policy name.
func (rl *RedisRateLimiter) Handler(policyName string) gin.HandlerFunc {
	policy, ok := rl.Policies[policyName]
	if !ok {
		panic("Unknown rate limiter policy: " + policyName)
	}

	return func(c *gin.Context) {
		ctx := context.Background()
		clientIP := c.ClientIP()
		key := "ratelimit:" + policy.Name + ":" + clientIP

		count, err := rl.RedisClient.Incr(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiter error"})
			c.Abort()
			return
		}

		if count == 1 {
			rl.RedisClient.Expire(ctx, key, policy.Window)
		}

		if count > int64(policy.Requests) {
			ttl, _ := rl.RedisClient.TTL(ctx, key).Result()
			c.Header("Retry-After", strconv.Itoa(int(ttl.Seconds())))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"retry_after": ttl.Seconds(),
				"policy":      policy.Name,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *RedisRateLimiter) Aggressive() gin.HandlerFunc {
	return rl.Handler("aggressive")
}

func (rl *RedisRateLimiter) Moderate() gin.HandlerFunc {
	return rl.Handler("moderate")
}

func (rl *RedisRateLimiter) Relaxed() gin.HandlerFunc {
	return rl.Handler("relaxed")
}
