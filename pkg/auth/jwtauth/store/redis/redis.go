package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// Config redis配置参数
type Config struct {
	Addr      string
	DB        int
	Password  string
	KeyPrefix string
}

// NewStore 创建基于redis存储实例
func NewStore(cfg *Config) *Store {
	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		Password: cfg.Password,
	})
	return &Store{
		cli:    cli,
		prefix: cfg.KeyPrefix,
	}
}

type redisClienter interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Expire(key string, expiration time.Duration) *redis.BoolCmd
	Exists(keys ...string) *redis.IntCmd
	TxPipeline() redis.Pipeliner
	Del(keys ...string) *redis.IntCmd
	Close() error
}

// Store redis存储
type Store struct {
	cli    redisClienter
	prefix string
}

func (a *Store) wrapperKey(key string) string {
	return fmt.Sprintf("%s%s", a.prefix, key)
}

// Set ...
func (a *Store) Set(tokenString string, expiration time.Duration) error {
	cmd := a.cli.Set(a.wrapperKey(tokenString), "1", expiration)
	return cmd.Err()
}

// Check ...
func (a *Store) Check(tokenString string) (bool, error) {
	cmd := a.cli.Exists(a.wrapperKey(tokenString))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

// Close ...
func (a *Store) Close() error {
	return a.cli.Close()
}
