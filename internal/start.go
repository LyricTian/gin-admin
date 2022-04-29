package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/inject"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
	loggerhook "github.com/LyricTian/gin-admin/v9/pkg/logger/hook"

	"github.com/LyricTian/captcha"
	"github.com/LyricTian/captcha/store"
	"github.com/go-redis/redis"
	"github.com/google/gops/agent"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type options struct {
	ConfigDir string
	WWWDir    string
	Version   string
}

type Option func(*options)

func SetConfigDir(dir string) Option {
	return func(o *options) {
		o.ConfigDir = dir
	}
}

func SetWWWDir(dir string) Option {
	return func(o *options) {
		o.WWWDir = dir
	}
}

func SetVersion(s string) Option {
	return func(o *options) {
		o.Version = s
	}
}

func Start(ctx context.Context, opts ...Option) (func(), error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	cfgFile := filepath.Join(o.ConfigDir, "config.toml")
	logger.WithContext(ctx).Printf("Load config file: %s", cfgFile)
	config.MustLoad(cfgFile)

	if v := o.WWWDir; v != "" {
		config.C.WWW = v
	}

	config.C.ConfigDir = o.ConfigDir
	config.Print()

	logger.WithContext(ctx).Printf("Start server, #run_mode %s, #version %s, #pid %d", config.C.RunMode, o.Version, os.Getpid())

	loggerCleanFunc, err := InitLogger(ctx)
	if err != nil {
		return nil, err
	}
	monitorCleanFunc := InitMonitor(ctx)
	InitCaptcha()

	injector, injectorCleanFunc, err := inject.BuildInjector(ctx)
	if err != nil {
		return nil, err
	}

	httpServerCleanFunc := InitHTTPServer(ctx, injector.Engine)

	return func() {
		httpServerCleanFunc()
		injectorCleanFunc()
		monitorCleanFunc()
		loggerCleanFunc()
	}, nil
}

func InitLogger(ctx context.Context) (func(), error) {
	c := config.C.Log
	logger.SetLevel(logger.Level(c.Level))
	logger.SetFormatter(c.Format)

	var file *rotatelogs.RotateLogs
	if c.Output != "" {
		switch c.Output {
		case "stdout":
			logger.SetOutput(os.Stdout)
		case "stderr":
			logger.SetOutput(os.Stderr)
		case "file":
			if name := c.OutputFile; name != "" {
				_ = os.MkdirAll(filepath.Dir(name), 0755)
				logf, err := rotatelogs.New(
					name+".%Y-%m-%d",
					rotatelogs.WithLinkName(name),
					rotatelogs.WithMaxAge(time.Duration(c.RotationMaxAge)*time.Hour*24),
					rotatelogs.WithRotationCount(c.RotationCount),
					rotatelogs.WithRotationSize(int64(c.RotationSize)*1024*1024),
				)
				if err != nil {
					return nil, err
				}
				logger.SetOutput(logf)
				file = logf
			}
		}
	}

	var hooks []*loggerhook.Hook
	for _, h := range c.Hooks {
		var levels []logger.Level
		for _, lvl := range h.Levels {
			plvl, err := logger.ParseLevel(lvl)
			if err != nil {
				return nil, err
			}
			levels = append(levels, plvl)
		}

		extra := map[string]interface{}{
			"app_name": config.C.AppName,
		}

		switch h.Type {
		case "gorm":
			db, err := inject.NewGormDB(ctx)
			if err != nil {
				return nil, err
			}

			h := loggerhook.New(loggerhook.NewGormHook(db),
				loggerhook.SetMaxWorkers(h.MaxThread),
				loggerhook.SetMaxJobs(h.MaxBuffer),
				loggerhook.SetLevels(levels...),
				loggerhook.SetExtra(extra),
			)
			logger.AddHook(h)
			hooks = append(hooks, h)
		}
	}

	return func() {
		if file != nil {
			file.Close()
		}

		for _, h := range hooks {
			h.Flush()
		}
	}, nil
}

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

func InitMonitor(ctx context.Context) func() {
	if c := config.C.Monitor; c.Enable {
		// ShutdownCleanup set false to prevent automatically closes on os.Interrupt
		// and close agent manually before service shutting down
		err := agent.Listen(agent.Options{Addr: c.Addr, ConfigDir: c.ConfigDir, ShutdownCleanup: false})
		if err != nil {
			logger.WithContext(ctx).Errorf("Agent monitor error: %s", err.Error())
		}
		return func() {
			agent.Close()
		}
	}
	return func() {}
}

func InitHTTPServer(ctx context.Context, handler http.Handler) func() {
	cfg := config.C.HTTP
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		logger.WithContext(ctx).Printf("HTTP server start, #addr %s", addr)

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
			logger.WithContext(ctx).Errorf(err.Error())
		}
	}
}
