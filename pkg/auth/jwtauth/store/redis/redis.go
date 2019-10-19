package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// Config redis配置参数
type Config struct {
	Addr      string // 地址(IP:Port)
	DB        int    // 数据库
	Password  string // 密码
	KeyPrefix string // 存储key的前缀
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

// NewStoreWithClient 使用redis客户端创建存储实例
func NewStoreWithClient(cli *redis.Client, keyPrefix string) *Store {
	return &Store{
		cli:    cli,
		prefix: keyPrefix,
	}
}

// NewStoreWithClusterClient 使用redis集群客户端创建存储实例
func NewStoreWithClusterClient(cli *redis.ClusterClient, keyPrefix string) *Store {
	return &Store{
		cli:    cli,
		prefix: keyPrefix,
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

func (s *Store) wrapperKey(key string) string {
	return fmt.Sprintf("%s%s", s.prefix, key)
}

// Set ...
func (s *Store) Set(ctx context.Context, tokenString string, expiration time.Duration) error {
	cmd := s.cli.Set(s.wrapperKey(tokenString), "1", expiration)
	return cmd.Err()
}

// Delete ...
func (s *Store) Delete(ctx context.Context, tokenString string) error {
	cmd := s.cli.Del(tokenString)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

// Check ...
func (s *Store) Check(ctx context.Context, tokenString string) (bool, error) {
	cmd := s.cli.Exists(s.wrapperKey(tokenString))
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

// Close ...
func (s *Store) Close() error {
	return s.cli.Close()
}
