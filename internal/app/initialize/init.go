package initialize

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/LyricTian/captcha"
	"github.com/LyricTian/captcha/store"
	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/pkg/logger"
	"github.com/go-redis/redis"
	"github.com/google/gops/agent"
)

// InitCaptcha 初始化图形验证码
func InitCaptcha() {
	cfg := config.C.Captcha
	if cfg.Store == "redis" {
		rc := config.C.Redis
		captcha.SetCustomStore(store.NewRedisStore(&redis.Options{
			Addr:     rc.Addr,
			Password: rc.Password,
			DB:       cfg.RedisDB,
		}, captcha.Expiration, logger.StandardLogger(), cfg.RedisPrefix))
	}
}

// InitMonitor 初始化服务监控
func InitMonitor(ctx context.Context) {
	if c := config.C.Monitor; c.Enable {
		err := agent.Listen(agent.Options{Addr: c.Addr, ConfigDir: c.ConfigDir, ShutdownCleanup: true})
		if err != nil {
			logger.Errorf(ctx, "Agent monitor error: %s", err.Error())
		}
	}
}

// InitHTTPServer 初始化http服务
func InitHTTPServer(ctx context.Context, handler http.Handler) func() {
	cfg := config.C.HTTP
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		logger.Printf(ctx, "HTTP server is running at %s.", addr)
		var err error
		if cfg.CertFile != "" && cfg.KeyFile != "" {
			srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.ShutdownTimeout))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Errorf(ctx, err.Error())
		}
	}
}
