package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/LyricTian/gin-admin/v10/internal/config"
	"github.com/LyricTian/gin-admin/v10/internal/library/utilx"
	"github.com/LyricTian/gin-admin/v10/internal/library/wirex"
	"github.com/LyricTian/gin-admin/v10/internal/middlewares"
	_ "github.com/LyricTian/gin-admin/v10/internal/swagger"
	"github.com/LyricTian/gin-admin/v10/pkg/errors"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func Start(ctx context.Context) (func(), error) {
	injector, injectorClean, err := wirex.BuildInjector(ctx)
	if err != nil {
		return nil, err
	}

	if err := injector.M.Init(ctx); err != nil {
		return nil, err
	}

	httpServerClean, err := startHTTPServer(ctx, injector, registerAPIs)
	if err != nil {
		return injectorClean, err
	}

	if addr := config.C.General.PprofAddr; addr != "" {
		logging.Context(ctx).Info("Pprof server is listening on " + addr)
		go func() {
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				logging.Context(ctx).Error("Failed to listen pprof server", zap.Error(err))
			}
		}()
	}

	return func() {
		if httpServerClean != nil {
			httpServerClean()
		}

		if injectorClean != nil {
			injectorClean()
		}
	}, nil
}

func registerAPIs(ctx context.Context, e *gin.Engine, injector *wirex.Injector) ([]string, error) {
	apiPrefix := "/api"
	gAPI := e.Group(apiPrefix)

	apiV1Prefix := apiPrefix + "/v1"
	gm := map[string]*gin.RouterGroup{
		apiPrefix:   gAPI,
		apiV1Prefix: gAPI.Group(apiV1Prefix),
	}
	if err := injector.M.RegisterAPIs(ctx, gm); err != nil {
		return nil, err
	}

	return []string{
		apiPrefix + "/",
	}, nil
}

func startHTTPServer(ctx context.Context, injector *wirex.Injector, registerAPIs func(ctx context.Context, e *gin.Engine, injector *wirex.Injector) ([]string, error)) (func(), error) {
	if config.C.General.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()

	e.Use(middlewares.RecoveryWithConfig(middlewares.RecoveryConfig{
		Skip: config.C.Middleware.Recovery.Skip,
	}))

	e.GET("/health", func(c *gin.Context) {
		utilx.ResOK(c)
	})

	e.NoMethod(func(c *gin.Context) {
		utilx.ResError(c, errors.MethodNotAllowed("", "Method not allowed"))
	})

	e.NoRoute(func(c *gin.Context) {
		utilx.ResError(c, errors.NotFound("", "Not found"))
	})

	if !config.C.General.DisableSwagger {
		e.StaticFile("/openapi.json", filepath.Join(config.C.General.ConfigDir, "openapi.json"))
		e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	e.Use(middlewares.CORSWithConfig(middlewares.CORSConfig{
		Enable:                 config.C.Middleware.CORS.Enable,
		AllowAllOrigins:        config.C.Middleware.CORS.AllowAllOrigins,
		AllowOrigins:           config.C.Middleware.CORS.AllowOrigins,
		AllowMethods:           config.C.Middleware.CORS.AllowMethods,
		AllowHeaders:           config.C.Middleware.CORS.AllowHeaders,
		AllowCredentials:       config.C.Middleware.CORS.AllowCredentials,
		ExposeHeaders:          config.C.Middleware.CORS.ExposeHeaders,
		MaxAge:                 config.C.Middleware.CORS.MaxAge,
		AllowWildcard:          config.C.Middleware.CORS.AllowWildcard,
		AllowBrowserExtensions: config.C.Middleware.CORS.AllowBrowserExtensions,
		AllowWebSockets:        config.C.Middleware.CORS.AllowWebSockets,
		AllowFiles:             config.C.Middleware.CORS.AllowFiles,
	}))

	// Register APIs
	allowedPathPrefixes, err := registerAPIs(ctx, e, injector)
	if err != nil {
		return nil, err
	}

	e.Use(middlewares.TraceWithConfig(middlewares.TraceConfig{
		SkippedPathPrefixes: config.C.Middleware.Trace.SkippedPathPrefixes,
		AllowedPathPrefixes: allowedPathPrefixes,
		RequestHeaderKey:    config.C.Middleware.Trace.RequestHeaderKey,
		ResponseTraceKey:    config.C.Middleware.Trace.ResponseTraceKey,
	}))

	e.Use(middlewares.LoggerWithConfig(middlewares.LoggerConfig{
		SkippedPathPrefixes:      config.C.Middleware.Logger.SkippedPathPrefixes,
		AllowedPathPrefixes:      allowedPathPrefixes,
		MaxOutputRequestBodyLen:  config.C.Middleware.Logger.MaxOutputRequestBodyLen,
		MaxOutputResponseBodyLen: config.C.Middleware.Logger.MaxOutputResponseBodyLen,
	}))

	e.Use(middlewares.CopyBodyWithConfig(middlewares.CopyBodyConfig{
		SkippedPathPrefixes: config.C.Middleware.CopyBody.SkippedPathPrefixes,
		AllowedPathPrefixes: allowedPathPrefixes,
		MaxContentLen:       config.C.Middleware.CopyBody.MaxContentLen,
	}))

	e.Use(middlewares.RateLimiterWithConfig(middlewares.RateLimiterConfig{
		Enable:              config.C.Middleware.RateLimiter.Enable,
		SkippedPathPrefixes: config.C.Middleware.RateLimiter.SkippedPathPrefixes,
		Period:              config.C.Middleware.RateLimiter.Period,
		MaxRequestsPerIP:    config.C.Middleware.RateLimiter.MaxRequestsPerIP,
		MaxRequestsPerUser:  config.C.Middleware.RateLimiter.MaxRequestsPerUser,
		StoreType:           config.C.Middleware.RateLimiter.Store.Type,
		MemoryStoreConfig: middlewares.RateLimiterMemoryConfig{
			Expiration:      time.Second * time.Duration(config.C.Middleware.RateLimiter.Store.Memory.Expiration),
			CleanupInterval: time.Second * time.Duration(config.C.Middleware.RateLimiter.Store.Memory.CleanupInterval),
		},
		RedisStoreConfig: middlewares.RateLimiterRedisConfig{
			Addr:     config.C.Middleware.RateLimiter.Store.Redis.Addr,
			Password: config.C.Middleware.RateLimiter.Store.Redis.Password,
			DB:       config.C.Middleware.RateLimiter.Store.Redis.DB,
			Username: config.C.Middleware.RateLimiter.Store.Redis.Username,
		},
	}))

	if dir := config.C.Middleware.Static.Dir; dir != "" {
		e.Use(middlewares.StaticWithConfig(middlewares.StaticConfig{
			Root:                dir,
			SkippedPathPrefixes: allowedPathPrefixes,
		}))
	}

	srv := &http.Server{
		Addr:         config.C.General.HTTP.Addr,
		Handler:      e,
		ReadTimeout:  time.Second * time.Duration(config.C.General.HTTP.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(config.C.General.HTTP.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(config.C.General.HTTP.IdleTimeout),
	}

	logging.Context(ctx).Info(fmt.Sprintf("HTTP server is listening on %s", srv.Addr))
	go func() {
		var err error
		if config.C.General.HTTP.CertFile != "" && config.C.General.HTTP.KeyFile != "" {
			srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = srv.ListenAndServeTLS(config.C.General.HTTP.CertFile, config.C.General.HTTP.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			logging.Context(ctx).Error("Failed to listen http server", zap.Error(err))
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(config.C.General.HTTP.ShutdownTimeout))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logging.Context(ctx).Error("Failed to shutdown http server", zap.Error(err))
		}
	}, nil
}
