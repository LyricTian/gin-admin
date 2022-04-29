package jwtauth

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type RedisConfig struct {
	Addr      string
	DB        int
	Password  string
	KeyPrefix string
}

func NewRedisStore(cfg *RedisConfig) Storer {
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		Password: cfg.Password,
	})

	return &redisStore{
		cli:    cli,
		prefix: cfg.KeyPrefix,
	}
}

func NewRedisStoreWithClient(cli *redis.Client, keyPrefix string) Storer {
	return &redisStore{
		cli:    cli,
		prefix: keyPrefix,
	}
}

func NewRedisStoreWithClusterClient(cli *redis.ClusterClient, keyPrefix string) Storer {
	return &redisStore{
		cli:    cli,
		prefix: keyPrefix,
	}
}

type redisClienter interface {
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Exists(keys ...string) *redis.IntCmd
	Del(keys ...string) *redis.IntCmd
	Close() error
}

type redisStore struct {
	cli    redisClienter
	prefix string
}

func (a *redisStore) wrapperKey(key string) string {
	return fmt.Sprintf("%s%s", a.prefix, key)
}

func (a *redisStore) Set(ctx context.Context, tokenStr string, expiration time.Duration) error {
	cmd := a.cli.Set(a.wrapperKey(tokenStr), "1", expiration)
	return cmd.Err()
}

func (a *redisStore) Delete(ctx context.Context, tokenStr string) error {
	cmd := a.cli.Del(a.wrapperKey(tokenStr))
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

func (a *redisStore) Check(ctx context.Context, tokenStr string) (bool, error) {
	cmd := a.cli.Exists(a.wrapperKey(tokenStr))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (a *redisStore) Close() error {
	return a.cli.Close()
}
