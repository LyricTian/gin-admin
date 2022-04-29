package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// Create redis-based cache
func NewRedisCache(cfg *RedisConfig) Cacher {
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		Password: cfg.Password,
	})

	return &redisCache{
		cli: cli,
	}
}

// Use redis client create cache
func NewRedisCacheWithClient(cli *redis.Client) Cacher {
	return &redisCache{
		cli: cli,
	}
}

// Use redis cluster client create cache
func NewRedisCacheWithClusterClient(cli *redis.ClusterClient) Cacher {
	return &redisCache{
		cli: cli,
	}
}

type redisClienter interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	Exists(keys ...string) *redis.IntCmd
	Del(keys ...string) *redis.IntCmd
	Close() error
}

type redisCache struct {
	cli redisClienter
}

func (a *redisCache) getKey(ns, key string) string {
	return fmt.Sprintf("%s_%s", ns, key)
}

func (a *redisCache) Set(ctx context.Context, namespace, key, value string, expiration ...time.Duration) error {
	var exp time.Duration
	if len(expiration) > 0 {
		exp = expiration[0]
	}

	cmd := a.cli.Set(a.getKey(namespace, key), value, exp)
	return cmd.Err()
}

func (a *redisCache) Get(ctx context.Context, namespace, key string) (string, error) {
	cmd := a.cli.Get(a.getKey(namespace, key))
	if err := cmd.Err(); err != nil && err != redis.Nil {
		return "", err
	}
	return cmd.Val(), nil
}

func (a *redisCache) GetAndDelete(ctx context.Context, namespace, key string) (string, error) {
	value, err := a.Get(ctx, namespace, key)
	if err != nil {
		return "", err
	}
	return value, a.Delete(ctx, namespace, key)
}

func (a *redisCache) Exists(ctx context.Context, namespace, key string) (bool, error) {
	cmd := a.cli.Exists(a.getKey(namespace, key))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (a *redisCache) Delete(ctx context.Context, namespace, key string) error {
	cmd := a.cli.Del(a.getKey(namespace, key))
	if err := cmd.Err(); err != nil && err != redis.Nil {
		return err
	}
	return nil
}

func (a *redisCache) Close(ctx context.Context) error {
	return a.cli.Close()
}
