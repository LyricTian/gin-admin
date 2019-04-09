package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate"
)

// RateLimiterMiddleware 请求频率限制中间件
func RateLimiterMiddleware(limiter *redis_rate.Limiter, skipper ...SkipperFunc) gin.HandlerFunc {
	cfg := config.GetRateLimiter()
	return func(c *gin.Context) {
		if (len(skipper) > 0 && skipper[0](c)) || limiter == nil {
			c.Next()
			return
		}

		ctx := context.New(c)
		userID := ctx.GetUserID()
		if userID == "" {
			c.Next()
			return
		}

		limit := cfg.Count
		rate, delay, allowed := limiter.AllowMinute(userID, limit)
		if !allowed {
			h := c.Writer.Header()
			h.Set("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
			h.Set("X-RateLimit-Remaining", strconv.FormatInt(limit-rate, 10))
			delaySec := int64(delay / time.Second)
			h.Set("X-RateLimit-Delay", strconv.FormatInt(delaySec, 10))
			ctx.ResErrorWithStatus(errors.New("请求过于频繁"), http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}
