package cachex

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/dgraph-io/badger/v3"
)

type BadgerConfig struct {
	Path string
}

// Create badger-based cache
func NewBadgerCache(cfg BadgerConfig, opts ...Option) Cacher {
	defaultOpts := &options{
		Delimiter: defaultDelimiter,
	}

	for _, o := range opts {
		o(defaultOpts)
	}

	badgerOpts := badger.DefaultOptions(cfg.Path)
	badgerOpts = badgerOpts.WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(badgerOpts)
	if err != nil {
		panic(err)
	}

	return &badgerCache{
		opts: defaultOpts,
		db:   db,
	}
}

type badgerCache struct {
	opts *options
	db   *badger.DB
}

func (a *badgerCache) getKey(ns, key string) string {
	return fmt.Sprintf("%s%s%s", ns, a.opts.Delimiter, key)
}

func (a *badgerCache) strToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func (a *badgerCache) bytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func (a *badgerCache) Set(ctx context.Context, ns, key, value string, expiration ...time.Duration) error {
	return a.db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(a.strToBytes(a.getKey(ns, key)), a.strToBytes(value))
		if len(expiration) > 0 {
			entry = entry.WithTTL(expiration[0])
		}
		return txn.SetEntry(entry)
	})
}

func (a *badgerCache) Get(ctx context.Context, ns, key string) (string, bool, error) {
	value := ""
	ok := false
	err := a.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(a.strToBytes(a.getKey(ns, key)))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			return err
		}
		ok = true
		val, err := item.ValueCopy(nil)
		value = a.bytesToStr(val)
		return err
	})
	if err != nil {
		return "", false, err
	}
	return value, ok, nil
}

func (a *badgerCache) Exists(ctx context.Context, ns, key string) (bool, error) {
	exists := false
	err := a.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(a.strToBytes(a.getKey(ns, key)))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			return err
		}
		exists = true
		return nil
	})
	return exists, err
}

func (a *badgerCache) Delete(ctx context.Context, ns, key string) error {
	b, err := a.Exists(ctx, ns, key)
	if err != nil {
		return err
	} else if !b {
		return nil
	}

	return a.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(a.strToBytes(a.getKey(ns, key)))
	})
}

func (a *badgerCache) GetAndDelete(ctx context.Context, ns, key string) (string, bool, error) {
	value, ok, err := a.Get(ctx, ns, key)
	if err != nil {
		return "", false, err
	} else if !ok {
		return "", false, nil
	}

	err = a.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(a.strToBytes(a.getKey(ns, key)))
	})
	if err != nil {
		return "", false, err
	}

	return value, true, nil
}

func (a *badgerCache) Iterator(ctx context.Context, ns string, fn func(ctx context.Context, key, value string) bool) error {
	return a.db.View(func(txn *badger.Txn) error {
		iterOpts := badger.DefaultIteratorOptions
		iterOpts.Prefix = a.strToBytes(a.getKey(ns, ""))
		it := txn.NewIterator(iterOpts)
		defer it.Close()

		it.Rewind()
		for it.Valid() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			key := a.bytesToStr(item.Key())
			if !fn(ctx, strings.TrimPrefix(key, a.getKey(ns, "")), a.bytesToStr(val)) {
				break
			}
			it.Next()
		}
		return nil
	})
}

func (a *badgerCache) Close(ctx context.Context) error {
	return a.db.Close()
}
