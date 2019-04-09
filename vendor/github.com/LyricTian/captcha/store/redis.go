package store

import (
	"encoding/hex"
	"time"

	"github.com/go-redis/redis"
)

// NewRedisStore create an instance of a redis store
func NewRedisStore(opts *RedisOptions, expiration time.Duration, out Logger, prefix ...string) Store {
	if opts == nil {
		panic("options cannot be nil")
	}
	return NewRedisStoreWithCli(
		redis.NewClient(opts.redisOptions()),
		expiration,
		out,
		prefix...,
	)
}

// NewRedisStoreWithCli create an instance of a redis store
func NewRedisStoreWithCli(cli *redis.Client, expiration time.Duration, out Logger, prefix ...string) Store {
	store := &redisStore{
		cli:        cli,
		expiration: expiration,
		out:        out,
	}
	if len(prefix) > 0 {
		store.prefix = prefix[0]
	}
	return store
}

// NewRedisClusterStore create an instance of a redis cluster store
func NewRedisClusterStore(opts *RedisClusterOptions, expiration time.Duration, out Logger, prefix ...string) Store {
	if opts == nil {
		panic("options cannot be nil")
	}
	return NewRedisClusterStoreWithCli(
		redis.NewClusterClient(opts.redisClusterOptions()),
		expiration,
		out,
		prefix...,
	)
}

// NewRedisClusterStoreWithCli create an instance of a redis cluster store
func NewRedisClusterStoreWithCli(cli *redis.ClusterClient, expiration time.Duration, out Logger, prefix ...string) Store {
	store := &redisStore{
		cli:        cli,
		expiration: expiration,
		out:        out,
	}
	if len(prefix) > 0 {
		store.prefix = prefix[0]
	}
	return store
}

type clienter interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(keys ...string) *redis.IntCmd
}

type redisStore struct {
	cli        clienter
	prefix     string
	out        Logger
	expiration time.Duration
}

func (s *redisStore) getKey(id string) string {
	return s.prefix + id
}

func (s *redisStore) printf(format string, args ...interface{}) {
	if s.out != nil {
		s.out.Printf(format, args...)
	}
}

func (s *redisStore) Set(id string, digits []byte) {
	cmd := s.cli.Set(s.getKey(id), hex.EncodeToString(digits), s.expiration)
	if err := cmd.Err(); err != nil {
		s.printf("redis execution set command error: %s", err.Error())
	}
	return
}

func (s *redisStore) Get(id string, clear bool) []byte {
	key := s.getKey(id)
	cmd := s.cli.Get(key)
	if err := cmd.Err(); err != nil {
		if err == redis.Nil {
			return nil
		}
		s.printf("redis execution get command error: %s", err.Error())
		return nil
	}

	b, err := hex.DecodeString(cmd.Val())
	if err != nil {
		s.printf("hex decoding error: %s", err.Error())
		return nil
	}

	if clear {
		cmd := s.cli.Del(key)
		if err := cmd.Err(); err != nil {
			s.printf("redis execution del command error: %s", err.Error())
			return nil
		}
	}

	return b
}
