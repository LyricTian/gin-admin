package cachex

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

// Create redis-based cache
func NewRedisCache(cfg RedisConfig, opts ...Option) Cacher {
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return newRedisCache(cli, opts...)
}

// Use redis client create cache
func NewRedisCacheWithClient(cli *redis.Client, opts ...Option) Cacher {
	return newRedisCache(cli, opts...)
}

// Use redis cluster client create cache
func NewRedisCacheWithClusterClient(cli *redis.ClusterClient, opts ...Option) Cacher {
	return newRedisCache(cli, opts...)
}

func newRedisCache(cli redisClienter, opts ...Option) Cacher {
	defaultOpts := &options{
		Delimiter: defaultDelimiter,
	}

	for _, o := range opts {
		o(defaultOpts)
	}

	return &redisCache{
		opts: defaultOpts,
		cli:  cli,
	}
}

type redisClienter interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
	Close() error
}

type redisCache struct {
	opts *options
	cli  redisClienter
}

func (a *redisCache) getKey(ns, key string) string {
	return fmt.Sprintf("%s%s%s", ns, a.opts.Delimiter, key)
}

func (a *redisCache) Set(ctx context.Context, ns, key, value string, expiration ...time.Duration) error {
	var exp time.Duration
	if len(expiration) > 0 {
		exp = expiration[0]
	}

	cmd := a.cli.Set(ctx, a.getKey(ns, key), value, exp)
	return cmd.Err()
}

func (a *redisCache) Get(ctx context.Context, ns, key string) (string, bool, error) {
	cmd := a.cli.Get(ctx, a.getKey(ns, key))
	if err := cmd.Err(); err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, err
	}
	return cmd.Val(), true, nil
}

func (a *redisCache) Exists(ctx context.Context, ns, key string) (bool, error) {
	cmd := a.cli.Exists(ctx, a.getKey(ns, key))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (a *redisCache) Delete(ctx context.Context, ns, key string) error {
	b, err := a.Exists(ctx, ns, key)
	if err != nil {
		return err
	} else if !b {
		return nil
	}

	cmd := a.cli.Del(ctx, a.getKey(ns, key))
	if err := cmd.Err(); err != nil && err != redis.Nil {
		return err
	}
	return nil
}

func (a *redisCache) GetAndDelete(ctx context.Context, ns, key string) (string, bool, error) {
	value, ok, err := a.Get(ctx, ns, key)
	if err != nil {
		return "", false, err
	} else if !ok {
		return "", false, nil
	}

	cmd := a.cli.Del(ctx, a.getKey(ns, key))
	if err := cmd.Err(); err != nil && err != redis.Nil {
		return "", false, err
	}
	return value, true, nil
}

func (a *redisCache) Iterator(ctx context.Context, ns string, fn func(ctx context.Context, key, value string) bool) error {
	var cursor uint64 = 0

LB_LOOP:
	for {
		cmd := a.cli.Scan(ctx, cursor, a.getKey(ns, "*"), 100)
		if err := cmd.Err(); err != nil {
			return err
		}

		keys, c, err := cmd.Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			cmd := a.cli.Get(ctx, key)
			if err := cmd.Err(); err != nil {
				if err == redis.Nil {
					continue
				}
				return err
			}
			if next := fn(ctx, strings.TrimPrefix(key, a.getKey(ns, "")), cmd.Val()); !next {
				break LB_LOOP
			}
		}

		if c == 0 {
			break
		}
		cursor = c
	}

	return nil
}

func (a *redisCache) Close(ctx context.Context) error {
	return a.cli.Close()
}
