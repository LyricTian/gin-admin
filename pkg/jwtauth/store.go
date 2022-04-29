package jwtauth

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/tidwall/buntdb"
)

// JWT token storage interface
type Storer interface {
	Set(ctx context.Context, tokenStr string, expiration time.Duration) error
	Delete(ctx context.Context, tokenStr string) error
	Check(ctx context.Context, tokenStr string) (bool, error)
	Close() error
}

// Implementation Storer interface base buntdb storage
func NewBuntDBStore(path string) (Storer, error) {
	if path != ":memory:" {
		os.MkdirAll(filepath.Dir(path), 0777)
	}

	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}

	return &buntdbStore{
		db: db,
	}, nil
}

type buntdbStore struct {
	db *buntdb.DB
}

func (a *buntdbStore) Set(ctx context.Context, tokenStr string, expiration time.Duration) error {
	return a.db.Update(func(tx *buntdb.Tx) error {
		var opts *buntdb.SetOptions
		if expiration > 0 {
			opts = &buntdb.SetOptions{Expires: true, TTL: expiration}
		}
		_, _, err := tx.Set(tokenStr, "1", opts)
		return err
	})
}

func (a *buntdbStore) Delete(ctx context.Context, tokenStr string) error {
	return a.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(tokenStr)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		return nil
	})
}

func (a *buntdbStore) Check(ctx context.Context, tokenStr string) (bool, error) {
	var exists bool
	err := a.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(tokenStr)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		exists = val == "1"
		return nil
	})
	return exists, err
}

func (a *buntdbStore) Close() error {
	return a.db.Close()
}
