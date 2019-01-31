package inject

import (
	"context"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/errors"
	mgorm "github.com/LyricTian/gin-admin/src/model/gorm"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/casbin/casbin"
	"github.com/facebookgo/inject"
	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
	"golang.org/x/time/rate"
)

// Object 注入对象
type Object struct {
	GormDB      *gormplus.DB
	Enforcer    *casbin.Enforcer
	CtlCommon   *ctl.Common
	RateLimiter *redis_rate.Limiter
}

// Init 初始化依赖注入
func Init(ctx context.Context) (*Object, error) {
	g := new(inject.Graph)
	obj := &Object{
		RateLimiter: getRateLimiter(),
	}

	// 注入存储层
	switch {
	case config.IsGormDB():
		db, err := mgorm.Init(ctx, g)
		if err != nil {
			return nil, err
		}
		obj.GormDB = db
	default:
		return nil, errors.New("unknown model")
	}

	// 注入casbin
	enforcer := casbin.NewEnforcer(config.GetCasbinModelConf(), false)
	g.Provide(&inject.Object{Value: enforcer})
	obj.Enforcer = enforcer

	// 注入控制器
	ctlCommon := new(ctl.Common)
	g.Provide(&inject.Object{Value: ctlCommon})
	obj.CtlCommon = ctlCommon

	if err := g.Populate(); err != nil {
		return nil, err
	}

	return obj, nil
}

func getRateLimiter() *redis_rate.Limiter {
	rateConfig := config.GetRateLimiterConfig()
	if !rateConfig.Enable {
		return nil
	}

	redisConfig := config.GetRedisConfig()
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": redisConfig.Addr,
		},
		Password: redisConfig.Password,
		DB:       rateConfig.RedisDB,
	})
	limiter := redis_rate.NewLimiter(ring)
	limiter.Fallback = rate.NewLimiter(rate.Inf, 0)
	return limiter
}
