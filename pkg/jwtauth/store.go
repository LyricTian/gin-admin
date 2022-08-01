package jwtauth

import (
	"context"
	"time"

	"github.com/LyricTian/gin-admin/v9/pkg/x/cachex"
)

// Storer is the interface that storage the token.
type Storer interface {
	Set(ctx context.Context, tokenStr string, expiration time.Duration) error
	Delete(ctx context.Context, tokenStr string) error
	Check(ctx context.Context, tokenStr string) (bool, error)
	Close(ctx context.Context) error
}

type storeOptions struct {
	CacheNS string // default "jwt"
}

type StoreOption func(*storeOptions)

func WithCacheNS(ns string) StoreOption {
	return func(o *storeOptions) {
		o.CacheNS = ns
	}
}

func NewStoreWithCache(cache cachex.Cacher, opts ...StoreOption) Storer {
	s := &storeImpl{
		cache: cache,
		opts: &storeOptions{
			CacheNS: "jwt",
		},
	}
	for _, opt := range opts {
		opt(s.opts)
	}
	return s
}

type storeImpl struct {
	opts  *storeOptions
	cache cachex.Cacher
}

func (s *storeImpl) Set(ctx context.Context, tokenStr string, expiration time.Duration) error {
	return s.cache.Set(ctx, s.opts.CacheNS, tokenStr, "", expiration)
}

func (s *storeImpl) Delete(ctx context.Context, tokenStr string) error {
	return s.cache.Delete(ctx, s.opts.CacheNS, tokenStr)
}

func (s *storeImpl) Check(ctx context.Context, tokenStr string) (bool, error) {
	_, found, err := s.cache.Get(ctx, s.opts.CacheNS, tokenStr)
	return found, err
}

func (s *storeImpl) Close(ctx context.Context) error {
	return s.cache.Close(ctx)
}
