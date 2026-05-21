package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

type RateLimitResult struct {
	Allowed    bool
	Limit      int
	Remaining  int64
	ResetAfter time.Duration
}

func NewRateLimiter(
	client *redis.Client,
	limit int,
	window time.Duration,
) *RateLimiter {
	return &RateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Allow(
	ctx context.Context,
	key string,
) (RateLimitResult, error) {
	current, err := rl.client.Incr(ctx, key).Result()
	if err != nil {
		return RateLimitResult{}, err
	}

	// Seta expiracao so na primeira request da janela.
	if current == 1 {
		err = rl.client.Expire(
			ctx,
			key,
			rl.window,
		).Err()

		if err != nil {
			return RateLimitResult{}, err
		}
	}

	ttl, err := rl.client.TTL(ctx, key).Result()
	if err != nil {
		return RateLimitResult{}, err
	}

	if ttl < 0 {
		ttl = rl.window
	}

	remaining := int64(rl.limit) - current
	if remaining < 0 {
		remaining = 0
	}

	return RateLimitResult{
		Allowed:    current <= int64(rl.limit),
		Limit:      rl.limit,
		Remaining:  remaining,
		ResetAfter: ttl,
	}, nil
}

func RateLimiterMiddleware(
	rl *RateLimiter,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := rl.Allow(
			ctx.Request.Context(),
			rateLimitKey(ctx),
		)

		if err != nil {
			ctx.Header("X-RateLimit-Status", "unavailable")
			ctx.Next()
			return
		}

		ctx.Header("X-RateLimit-Limit", strconv.Itoa(result.Limit))
		ctx.Header("X-RateLimit-Remaining", strconv.FormatInt(result.Remaining, 10))
		ctx.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(result.ResetAfter).Unix(), 10))

		if !result.Allowed {
			retryAfter := secondsUntil(result.ResetAfter)

			ctx.Header("Retry-After", retryAfter)
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too Many Requests",
				"retry_after": retryAfter,
			})

			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func rateLimitKey(ctx *gin.Context) string {
	route := ctx.FullPath()
	if route == "" {
		route = ctx.Request.URL.Path
	}

	return "rate_limit:" + ctx.ClientIP() + ":" + ctx.Request.Method + ":" + route
}

func secondsUntil(duration time.Duration) string {
	if duration <= 0 {
		return "1"
	}

	seconds := int64(duration / time.Second)
	if duration%time.Second != 0 {
		seconds++
	}
	if seconds < 1 {
		seconds = 1
	}

	return strconv.FormatInt(seconds, 10)
}
