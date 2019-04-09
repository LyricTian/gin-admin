package auth

import (
	"time"

	"github.com/go-redis/redis"
)

type redisClienter interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Expire(key string, expiration time.Duration) *redis.BoolCmd
	Exists(keys ...string) *redis.IntCmd
	TxPipeline() redis.Pipeliner
	Del(keys ...string) *redis.IntCmd
	Close() error
}

// RedisConfig redis配置参数
type RedisConfig struct {
	Addr     string
	DB       int
	Password string
}

// NewRedisBlackStore 创建基于redis的黑名单存储实例
func NewRedisBlackStore(cfg *RedisConfig, keyPrefix string) BlackStorer {
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		Password: cfg.Password,
	})
	return NewRedisBlackStoreWithCli(cli, keyPrefix)
}

// NewRedisBlackStoreWithCli 创建基于redis的黑名单存储实例
func NewRedisBlackStoreWithCli(cli *redis.Client, keyPrefix string) BlackStorer {
	return &redisBlackStore{
		cli:    cli,
		prefix: keyPrefix,
	}
}

type redisBlackStore struct {
	cli    redisClienter
	prefix string
}

func (a *redisBlackStore) Set(tokenString string, expiration time.Duration) error {
	cmd := a.cli.Set(a.prefix+tokenString, "1", expiration)
	return cmd.Err()
}

func (a *redisBlackStore) Check(tokenString string) (bool, error) {
	cmd := a.cli.Exists(a.prefix + tokenString)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (a *redisBlackStore) Close() error {
	return a.cli.Close()
}
