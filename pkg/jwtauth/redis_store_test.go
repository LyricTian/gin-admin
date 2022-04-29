package jwtauth

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	addr = "127.0.0.1:55000"
)

func TestRedisStore(t *testing.T) {
	store := NewRedisStore(&RedisConfig{
		Addr:      addr,
		DB:        1,
		KeyPrefix: "test_",
	})

	defer store.Close()

	key := "test"
	ctx := context.Background()
	err := store.Set(ctx, key, 0)

	// If redis connection failed,skip
	if _, ok := err.(*net.OpError); ok {
		return
	}
	assert.Nil(t, err)

	b, err := store.Check(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, b)

	err = store.Delete(ctx, key)
	assert.Nil(t, err)
}
