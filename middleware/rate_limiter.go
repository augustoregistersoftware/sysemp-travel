package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
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
) (bool, error) {

	current, err := rl.client.Incr(ctx, key).Result()

	if err != nil {
		return false, err
	}

	// seta expiração só na primeira request
	if current == 1 {

		err = rl.client.Expire(
			ctx,
			key,
			rl.window,
		).Err()

		if err != nil {
			return false, err
		}
	}

	return current <= int64(rl.limit), nil
}

func RateLimiterMiddleware(
	rl *RateLimiter,
) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		clientIP := ctx.ClientIP()

		key := "rate_limit:" + clientIP

		allowed, err := rl.Allow(
			ctx.Request.Context(),
			key,
		)

		if err != nil {

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			ctx.Abort()
			return
		}

		if !allowed {

			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too Many Requests",
			})

			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
