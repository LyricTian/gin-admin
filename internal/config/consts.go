package config

import "github.com/LyricTian/gin-admin/v10/pkg/errors"

const (
	CacheNSForUser = "user"
	CacheNSForRole = "role"
)

const (
	CacheKeyForSyncToCasbin = "sync:casbin"
)

var (
	ErrInvalidToken = errors.Unauthorized("com.invalid.token", "Invalid access token")
)
