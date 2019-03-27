package inject

import (
	"context"

	"github.com/LyricTian/gin-admin/src/auth"
	"github.com/LyricTian/gin-admin/src/config"
	mgorm "github.com/LyricTian/gin-admin/src/model/gorm"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/casbin/casbin"
	"github.com/dgrijalva/jwt-go"
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
	Auth        *auth.Auth
}

// Init 初始化依赖注入
func Init(ctx context.Context) (*Object, error) {
	g := new(inject.Graph)
	obj := &Object{
		RateLimiter: rateLimiter(),
		Auth:        userAuth(),
		Enforcer:    casbin.NewEnforcer(config.GetCasbinModelConf(), false),
		CtlCommon:   new(ctl.Common),
	}

	// 注入存储层
	if s := config.GetStore(); s == "gorm" {
		db, err := mgorm.Init(ctx, g)
		if err != nil {
			return nil, err
		}
		obj.GormDB = db
	}

	// 注入auth
	g.Provide(&inject.Object{Value: obj.Auth})

	// 注入casbin
	g.Provide(&inject.Object{Value: obj.Enforcer})

	// 注入控制器
	g.Provide(&inject.Object{Value: obj.CtlCommon})

	if err := g.Populate(); err != nil {
		return nil, err
	}

	return obj, nil
}

// 初始化用户认证
func userAuth() *auth.Auth {
	var opts []auth.Option
	ac := config.GetAuth()
	opts = append(opts, auth.SetExpired(ac.Expired))
	opts = append(opts, auth.SetSigningKey([]byte(ac.SigningKey)))
	opts = append(opts, auth.SetKeyfunc(func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return []byte(ac.SigningKey), nil
	}))

	switch ac.SigningMethod {
	case "HS256":
		opts = append(opts, auth.SetSigningMethod(jwt.SigningMethodHS256))
	case "HS384":
		opts = append(opts, auth.SetSigningMethod(jwt.SigningMethodHS384))
	case "HS512":
		opts = append(opts, auth.SetSigningMethod(jwt.SigningMethodHS512))
	}

	switch ac.Store {
	case "file":
		opts = append(opts, auth.SetBlackStore(auth.NewFileBlackStore(ac.FilePath)))
	case "redis":
		opts = append(opts, auth.SetBlackStore(auth.NewRedisBlackStore(&auth.RedisConfig{
			Addr:     config.GetRedis().Addr,
			Password: config.GetRedis().Password,
			DB:       ac.RedisDB,
		}, ac.RedisPrefix)))
	}

	return auth.New(opts...)
}

// 初始化基于redis的访问频率限制(如果redis服务不通，则自动降级为内存模式)
func rateLimiter() *redis_rate.Limiter {
	rlc := config.GetRateLimiter()
	if !rlc.Enable {
		return nil
	}

	rc := config.GetRedis()
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": rc.Addr,
		},
		Password: rc.Password,
		DB:       rlc.RedisDB,
	})

	l := redis_rate.NewLimiter(ring)
	l.Fallback = rate.NewLimiter(rate.Inf, 0)
	return l
}
