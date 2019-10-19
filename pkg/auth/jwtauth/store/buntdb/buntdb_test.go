package buntdb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	store, err := NewStore(":memory:")
	assert.Nil(t, err)

	defer store.Close()

	key := "test"
	ctx := context.Background()
	err = store.Set(ctx, key, 0)
	assert.Nil(t, err)

	b, err := store.Check(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, true, b)

	err = store.Delete(ctx, key)
	assert.Nil(t, err)
}
