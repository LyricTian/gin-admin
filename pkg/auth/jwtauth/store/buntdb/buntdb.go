package buntdb

import (
	"os"
	"path/filepath"
	"time"

	"github.com/tidwall/buntdb"
)

// NewStore 创建基于buntdb的存储
func NewStore(p string) (*Store, error) {
	os.MkdirAll(filepath.Dir(p), 0777)

	db, err := buntdb.Open(p)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

// Store buntdb存储
type Store struct {
	db *buntdb.DB
}

// Set ...
func (a *Store) Set(tokenString string, expiration time.Duration) error {
	return a.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(tokenString, "1", &buntdb.SetOptions{Expires: true, TTL: expiration})
		return err
	})
}

// Check ...
func (a *Store) Check(tokenString string) (bool, error) {
	var exists bool
	err := a.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(tokenString)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		exists = val == "1"
		return nil
	})
	return exists, err
}

// Close ...
func (a *Store) Close() error {
	return a.db.Close()
}
