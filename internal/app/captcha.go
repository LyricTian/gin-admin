package app

import (
	"github.com/LyricTian/captcha"
	"github.com/LyricTian/captcha/store"
	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/pkg/logger"
)

// InitCaptcha 初始化图形验证码
func InitCaptcha() {
	cfg := config.GetGlobalConfig().Captcha
	if cfg.Store == "redis" {
		rc := config.GetGlobalConfig().Redis
		captcha.SetCustomStore(store.NewRedisStore(&store.RedisOptions{
			Addr:     rc.Addr,
			Password: rc.Password,
			DB:       cfg.RedisDB,
		}, captcha.Expiration, logger.StandardLogger(), cfg.RedisPrefix))
	}
}
