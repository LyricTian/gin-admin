package inject

import (
	"context"
	"time"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac"
	"github.com/LyricTian/gin-admin/v9/internal/module/sys"
	"github.com/LyricTian/gin-admin/v9/pkg/jwtauth"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
	"github.com/LyricTian/gin-admin/v9/pkg/x/cachex"
	"github.com/LyricTian/gin-admin/v9/pkg/x/gormx"
	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
) // end

// Inject global objects
type Injector struct {
	Auth  jwtauth.Auther
	Cache cachex.Cacher
	DB    *gorm.DB
	RBAC  *rbac.RBAC
	SYS   *sys.SYS
} // end

func InitAuth(ctx context.Context) (jwtauth.Auther, func(), error) {
	cfg := config.C.Middleware.Auth
	var opts []jwtauth.Option
	opts = append(opts, jwtauth.SetExpired(cfg.Expired))
	opts = append(opts, jwtauth.SetSigningKey(cfg.SigningKey, cfg.OldSigningKey))

	var method jwt.SigningMethod
	switch cfg.SigningMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	default:
		method = jwt.SigningMethodHS512
	}
	opts = append(opts, jwtauth.SetSigningMethod(method))

	var cache cachex.Cacher
	switch cfg.Store.Type {
	case "redis":
		cache = cachex.NewRedisCache(cachex.RedisConfig{
			Addr:     cfg.Store.Redis.Addr,
			DB:       cfg.Store.Redis.DB,
			Username: cfg.Store.Redis.Username,
			Password: cfg.Store.Redis.Password,
		}, cachex.WithDelimiter(cfg.Store.Delimiter))
	case "badger":
		cache = cachex.NewBadgerCache(cachex.BadgerConfig{
			Path: cfg.Store.Badger.Path,
		}, cachex.WithDelimiter(cfg.Store.Delimiter))
	default:
		cache = cachex.NewMemoryCache(cachex.MemoryConfig{
			CleanupInterval: time.Second * time.Duration(cfg.Store.Memory.CleanupInterval),
		}, cachex.WithDelimiter(cfg.Store.Delimiter))
	}

	auth := jwtauth.New(jwtauth.NewStoreWithCache(cache), opts...)
	return auth, func() {
		err := auth.Release(ctx)
		if err != nil {
			logger.Context(ctx).Error("Failed to release auth cache", zap.Error(err))
		}
	}, nil
}

func InitCache(ctx context.Context) (cachex.Cacher, func(), error) {
	cfg := config.C.Storage.Cache

	var cache cachex.Cacher
	switch cfg.Type {
	case "redis":
		cache = cachex.NewRedisCache(cachex.RedisConfig{
			Addr:     cfg.Redis.Addr,
			DB:       cfg.Redis.DB,
			Username: cfg.Redis.Username,
			Password: cfg.Redis.Password,
		}, cachex.WithDelimiter(cfg.Delimiter))
	case "badger":
		cache = cachex.NewBadgerCache(cachex.BadgerConfig{
			Path: cfg.Badger.Path,
		}, cachex.WithDelimiter(cfg.Delimiter))
	default:
		cache = cachex.NewMemoryCache(cachex.MemoryConfig{
			CleanupInterval: time.Second * time.Duration(cfg.Memory.CleanupInterval),
		}, cachex.WithDelimiter(cfg.Delimiter))
	}

	return cache, func() {
		err := cache.Close(ctx)
		if err != nil {
			logger.Context(ctx).Error("Failed to close cache", zap.Error(err))
		}
	}, nil
}

func InitDB(ctx context.Context) (*gorm.DB, func(), error) {
	cfg := config.C.Storage.DB
	db, err := gormx.New(gormx.Config{
		Debug:        cfg.Debug,
		DBType:       cfg.Type,
		DSN:          cfg.DSN,
		MaxLifetime:  cfg.MaxLifetime,
		MaxIdleTime:  cfg.MaxIdleTime,
		MaxOpenConns: cfg.MaxOpenConns,
		MaxIdleConns: cfg.MaxIdleConns,
		TablePrefix:  cfg.TablePrefix,
		Replicas: gormx.ReplicasConfig{
			DSNs:         cfg.Replicas.DSNs,
			Tables:       cfg.Replicas.Tables,
			MaxLifetime:  cfg.Replicas.MaxLifetime,
			MaxIdleTime:  cfg.Replicas.MaxIdleTime,
			MaxOpenConns: cfg.Replicas.MaxOpenConns,
			MaxIdleConns: cfg.Replicas.MaxIdleConns,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return db, func() {
		sqlDB, err := db.DB()
		if err == nil {
			err := sqlDB.Close()
			if err != nil {
				logger.Context(context.Background()).Error("Failed to close db", zap.Error(err))
			}
		}
	}, nil
}
