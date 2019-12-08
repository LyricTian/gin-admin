package redis_rate

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"golang.org/x/time/rate"
)

const redisPrefix = "rate"

type rediser interface {
	Del(...string) *redis.IntCmd
	Pipelined(func(pipe redis.Pipeliner) error) ([]redis.Cmder, error)
}

// Limiter controls how frequently events are allowed to happen.
type Limiter struct {
	redis rediser

	// Optional fallback limiter used when Redis is unavailable.
	Fallback *rate.Limiter
}

func NewLimiter(redis rediser) *Limiter {
	return &Limiter{
		redis: redis,
	}
}

// Reset resets the rate limit for the name in the given rate limit window.
func (l *Limiter) Reset(name string, dur time.Duration) error {
	udur := int64(dur / time.Second)
	slot := time.Now().Unix() / udur

	name = allowName(name, slot)
	return l.redis.Del(name).Err()
}

// Reset resets the rate limit for the name and limit.
func (l *Limiter) ResetRate(name string, rateLimit rate.Limit) error {
	if rateLimit == 0 {
		return nil
	}
	if rateLimit == rate.Inf {
		return nil
	}

	dur := time.Second
	limit := int64(rateLimit)
	if limit == 0 {
		limit = 1
		dur *= time.Duration(1 / rateLimit)
	}
	slot := time.Now().UnixNano() / dur.Nanoseconds()

	name = allowRateName(name, dur, slot)
	return l.redis.Del(name).Err()
}

// AllowN reports whether an event with given name may happen at time now.
// It allows up to maxn events within duration dur, with each interaction
// incrementing the limit by n.
func (l *Limiter) AllowN(
	name string, maxn int64, dur time.Duration, n int64,
) (count int64, delay time.Duration, allow bool) {
	udur := int64(dur / time.Second)
	utime := time.Now().Unix()
	slot := utime / udur
	delay = time.Duration((slot+1)*udur-utime) * time.Second

	if l.Fallback != nil {
		allow = l.Fallback.Allow()
	}

	name = allowName(name, slot)
	count, err := l.incr(name, dur, n)
	if err == nil {
		allow = count <= maxn
	}

	return count, delay, allow
}

// Allow is shorthand for AllowN(name, max, dur, 1).
func (l *Limiter) Allow(name string, maxn int64, dur time.Duration) (count int64, delay time.Duration, allow bool) {
	return l.AllowN(name, maxn, dur, 1)
}

// AllowMinute is shorthand for Allow(name, maxn, time.Minute).
func (l *Limiter) AllowMinute(name string, maxn int64) (count int64, delay time.Duration, allow bool) {
	return l.Allow(name, maxn, time.Minute)
}

// AllowHour is shorthand for Allow(name, maxn, time.Hour).
func (l *Limiter) AllowHour(name string, maxn int64) (count int64, delay time.Duration, allow bool) {
	return l.Allow(name, maxn, time.Hour)
}

// AllowRate reports whether an event may happen at time now.
// It allows up to rateLimit events each second.
func (l *Limiter) AllowRate(name string, rateLimit rate.Limit) (delay time.Duration, allow bool) {
	if rateLimit == 0 {
		return 0, false
	}
	if rateLimit == rate.Inf {
		return 0, true
	}

	dur := time.Second
	limit := int64(rateLimit)
	if limit == 0 {
		limit = 1
		dur *= time.Duration(1 / rateLimit)
	}
	now := time.Now()
	slot := now.UnixNano() / dur.Nanoseconds()

	if l.Fallback != nil {
		allow = l.Fallback.Allow()
	}

	name = allowRateName(name, dur, slot)
	count, err := l.incr(name, dur, 1)
	if err == nil {
		allow = count <= limit
	}

	if !allow {
		delay = time.Duration(slot+1)*dur - time.Duration(now.UnixNano())
	}

	return delay, allow
}

func (l *Limiter) incr(name string, dur time.Duration, n int64) (int64, error) {
	var incr *redis.IntCmd
	_, err := l.redis.Pipelined(func(pipe redis.Pipeliner) error {
		incr = pipe.IncrBy(name, n)
		pipe.Expire(name, dur)
		return nil
	})

	rate, _ := incr.Result()
	return rate, err
}

func allowName(name string, slot int64) string {
	return fmt.Sprintf("%s:%s-%d", redisPrefix, name, slot)
}

func allowRateName(name string, dur time.Duration, slot int64) string {
	return fmt.Sprintf("%s:%s-%d-%d", redisPrefix, name, dur, slot)
}
