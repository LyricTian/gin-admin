package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/tidwall/buntdb"
)

// A cache module with namespace
type Cacher interface {
	Set(ctx context.Context, namespace, key, value string, expiration ...time.Duration) error
	Get(ctx context.Context, namespace, key string) (string, error)
	GetAndDelete(ctx context.Context, namespace, key string) (string, error)
	Exists(ctx context.Context, namespace, key string) (bool, error)
	Delete(ctx context.Context, namespace, key string) error
	Close(ctx context.Context) error
}

// Create buntdb-based cache
func NewBuntdbCache(path string) (Cacher, error) {
	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}

	return &buntdbCache{
		db: db,
	}, nil
}

type buntdbCache struct {
	db *buntdb.DB
}

func (a *buntdbCache) getKey(ns, key string) string {
	return fmt.Sprintf("%s_%s", ns, key)
}

func (a *buntdbCache) Set(ctx context.Context, namespace, key, value string, expiration ...time.Duration) error {
	return a.db.Update(func(tx *buntdb.Tx) error {
		var opts *buntdb.SetOptions
		if len(expiration) > 0 {
			opts = &buntdb.SetOptions{Expires: true, TTL: expiration[0]}
		}
		_, _, err := tx.Set(a.getKey(namespace, key), value, opts)
		return err
	})
}

func (a *buntdbCache) Get(ctx context.Context, namespace, key string) (string, error) {
	var result string
	err := a.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(a.getKey(namespace, key))
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		result = val
		return nil
	})
	return result, err
}

func (a *buntdbCache) GetAndDelete(ctx context.Context, namespace, key string) (string, error) {
	value, err := a.Get(ctx, namespace, key)
	if err != nil {
		return "", err
	}
	return value, a.Delete(ctx, namespace, key)
}

func (a *buntdbCache) Exists(ctx context.Context, namespace, key string) (bool, error) {
	var exists bool
	err := a.db.View(func(tx *buntdb.Tx) error {
		_, err := tx.Get(a.getKey(namespace, key))
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		exists = err != buntdb.ErrNotFound
		return nil
	})
	return exists, err
}

func (a *buntdbCache) Delete(ctx context.Context, namespace, key string) error {
	return a.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(a.getKey(namespace, key))
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		return nil
	})
}

func (a *buntdbCache) Close(ctx context.Context) error {
	return a.db.Close()
}
