package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/LyricTian/gin-admin/internal/app/ginadmin/config"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/ginplus"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
	"golang.org/x/time/rate"
)

// RateLimiterMiddleware 请求频率限制中间件
func RateLimiterMiddleware(skipper ...SkipperFunc) gin.HandlerFunc {
	cfg := config.GetGlobalConfig().RateLimiter
	if !cfg.Enable {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	rc := config.GetGlobalConfig().Redis
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": rc.Addr,
		},
		Password: rc.Password,
		DB:       cfg.RedisDB,
	})

	limiter := redis_rate.NewLimiter(ring)
	limiter.Fallback = rate.NewLimiter(rate.Inf, 0)

	return func(c *gin.Context) {
		if (len(skipper) > 0 && skipper[0](c)) || limiter == nil {
			c.Next()
			return
		}

		userID := ginplus.GetUserID(c)
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
			ginplus.ResErrorWithStatus(c, errors.New("请求过于频繁"), http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}
