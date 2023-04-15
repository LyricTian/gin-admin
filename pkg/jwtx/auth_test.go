package jwtx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	cache := NewMemoryCache(MemoryConfig{CleanupInterval: time.Second})

	store := NewStoreWithCache(cache)
	ctx := context.Background()
	jwtAuth := New(store)

	userID := "test"
	token, err := jwtAuth.GenerateToken(ctx, userID)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	id, err := jwtAuth.ParseSubject(ctx, token.GetAccessToken())
	assert.Nil(t, err)
	assert.Equal(t, userID, id)

	err = jwtAuth.DestroyToken(ctx, token.GetAccessToken())
	assert.Nil(t, err)

	id, err = jwtAuth.ParseSubject(ctx, token.GetAccessToken())
	assert.NotNil(t, err)
	assert.EqualError(t, err, ErrInvalidToken.Error())
	assert.Empty(t, id)

	err = jwtAuth.Release(ctx)
	assert.Nil(t, err)
}
