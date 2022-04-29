package inject

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/pkg/cache"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
)

func InitCache(ctx context.Context) (cache.Cacher, func(), error) {
	cfg := config.C.Cache

	switch cfg.Store {
	case "redis":
		rcfg := config.C.Redis
		c := cache.NewRedisCache(&cache.RedisConfig{
			Addr:     rcfg.Addr,
			Password: rcfg.Password,
			DB:       cfg.RedisDB,
		})
		return c, func() {
			err := c.Close(ctx)
			if err != nil {
				logger.WithContext(ctx).Errorf("release redis cache failed: %s", err.Error())
			}
		}, nil
	default:
		c, err := cache.NewBuntdbCache(cfg.Path)
		return c, func() {
			err := c.Close(ctx)
			if err != nil {
				logger.WithContext(ctx).Errorf("release buntdb cache failed: %s", err.Error())
			}
		}, err
	}
}
