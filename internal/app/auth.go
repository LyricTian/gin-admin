package app

import (
	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/pkg/auth"
	"github.com/LyricTian/gin-admin/pkg/auth/jwtauth"
	"github.com/LyricTian/gin-admin/pkg/auth/jwtauth/store/buntdb"
	"github.com/LyricTian/gin-admin/pkg/auth/jwtauth/store/redis"
	jwt "github.com/dgrijalva/jwt-go"
)

// InitAuth 初始化用户认证
func InitAuth() (auth.Auther, error) {
	cfg := config.GetGlobalConfig().JWTAuth

	var opts []jwtauth.Option
	opts = append(opts, jwtauth.SetExpired(cfg.Expired))
	opts = append(opts, jwtauth.SetSigningKey([]byte(cfg.SigningKey)))
	opts = append(opts, jwtauth.SetKeyfunc(func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return []byte(cfg.SigningKey), nil
	}))

	switch cfg.SigningMethod {
	case "HS256":
		opts = append(opts, jwtauth.SetSigningMethod(jwt.SigningMethodHS256))
	case "HS384":
		opts = append(opts, jwtauth.SetSigningMethod(jwt.SigningMethodHS384))
	case "HS512":
		opts = append(opts, jwtauth.SetSigningMethod(jwt.SigningMethodHS512))
	}

	var store jwtauth.Storer
	switch cfg.Store {
	case "file":
		s, err := buntdb.NewStore(cfg.FilePath)
		if err != nil {
			return nil, err
		}
		store = s
	case "redis":
		rcfg := config.GetGlobalConfig().Redis
		store = redis.NewStore(&redis.Config{
			Addr:      rcfg.Addr,
			Password:  rcfg.Password,
			DB:        cfg.RedisDB,
			KeyPrefix: cfg.RedisPrefix,
		})
	}

	return jwtauth.New(store, opts...), nil
}
