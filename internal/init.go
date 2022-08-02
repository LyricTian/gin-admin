package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/LyricTian/captcha"
	"github.com/LyricTian/captcha/store"
	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/inject"
	"github.com/LyricTian/gin-admin/v9/internal/middleware"
	"github.com/LyricTian/gin-admin/v9/internal/x/contextx"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"github.com/LyricTian/gin-admin/v9/pkg/jwtauth"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "net/http/pprof"

	_ "github.com/LyricTian/gin-admin/v9/internal/swagger"
)

func Init(ctx context.Context) (func(), error) {
	initCaptcha()

	injector, cleanInjectFn, err := inject.BuildInjector(ctx)
	if err != nil {
		return nil, err
	}

	// Initialize modules
	{
		if err := injector.RBAC.Init(ctx); err != nil {
			return cleanInjectFn, err
		}
	} // end

	// Initialize http server with gin
	cleanHTTP, err := initHTTPServer(ctx, injector)
	if err != nil {
		return cleanInjectFn, err
	}

	if addr := config.C.General.PprofAddr; addr != "" {
		go func() {
			logger.Context(ctx).Info("Pprof server is listening on " + addr)
			logger.Context(ctx).Error(http.ListenAndServe(addr, nil).Error())
		}()
	}

	return func() {
		if cleanHTTP != nil {
			cleanHTTP()
		}
		if cleanInjectFn != nil {
			cleanInjectFn()
		}
	}, nil
}

func initHTTPServer(ctx context.Context, injector *inject.Injector) (func(), error) {
	gin.SetMode(config.C.General.RunMode)
	app := gin.New()

	app.GET("/health", func(c *gin.Context) {
		utilx.ResOK(c)
	})

	app.NoMethod(func(c *gin.Context) {
		logger.Context(c.Request.Context()).Warn("NoMethod",
			zap.String("method", c.Request.Method),
			zap.String("ip", c.ClientIP()),
			zap.String("remote_addr", c.Request.RemoteAddr),
			zap.String("url", c.Request.URL.String()),
			zap.String("proto", c.Request.Proto),
			zap.String("user_agent", c.GetHeader("User-Agent")),
		)
		utilx.ResError(c, errors.MethodNotAllowed(errors.ErrMethodNotAllowedID, "Method not allowed"))
	})

	app.NoRoute(func(c *gin.Context) {
		logger.Context(c.Request.Context()).Warn("NoMethod",
			zap.String("method", c.Request.Method),
			zap.String("ip", c.ClientIP()),
			zap.String("remote_addr", c.Request.RemoteAddr),
			zap.String("url", c.Request.URL.String()),
			zap.String("proto", c.Request.Proto),
			zap.String("user_agent", c.GetHeader("User-Agent")),
		)
		utilx.ResError(c, errors.NotFound(errors.ErrNotFoundID, "Not found"))
	})

	apiGroup := app.Group("/api")
	{
		initMiddlewares(apiGroup, injector)

		// Register RBAC APIs
		injector.RBAC.RegisterAPI(ctx, apiGroup)
	} // end

	if !config.C.General.DisableSwagger {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	if dir := config.C.Middleware.Static.Dir; dir != "" {
		app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:                dir,
			SkippedPathPrefixes: config.C.Middleware.Static.SkippedPathPrefixes,
		}))
	}

	logger.Context(ctx).Info(fmt.Sprintf("HTTP server is listening on %s", config.C.General.HTTP.Addr))

	srv := &http.Server{
		Addr:         config.C.General.HTTP.Addr,
		Handler:      app,
		ReadTimeout:  time.Second * time.Duration(config.C.General.HTTP.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(config.C.General.HTTP.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(config.C.General.HTTP.IdleTimeout),
	}
	go func() {
		var err error
		if config.C.General.HTTP.CertFile != "" && config.C.General.HTTP.KeyFile != "" {
			srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = srv.ListenAndServeTLS(config.C.General.HTTP.CertFile, config.C.General.HTTP.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(config.C.General.HTTP.ShutdownTimeout))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Context(ctx).Error("Failed to shutdown http server", zap.Error(err))
		}
	}, nil
}

func initMiddlewares(g *gin.RouterGroup, injector *inject.Injector) {
	g.Use(middleware.RecoveryWithConfig(middleware.RecoveryConfig{
		Skip: config.C.Middleware.Recovery.Skip,
	}))

	g.Use(middleware.TraceWithConfig(middleware.TraceConfig{
		SkippedPathPrefixes: config.C.Middleware.Trace.SkippedPathPrefixes,
		RequestHeaderKey:    config.C.Middleware.Trace.RequestHeaderKey,
		ResponseTraceKey:    config.C.Middleware.Trace.ResponseTraceKey,
	}))

	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		SkippedPathPrefixes:      config.C.Middleware.Logger.SkippedPathPrefixes,
		MaxOutputRequestBodyLen:  config.C.Middleware.Logger.MaxOutputRequestBodyLen,
		MaxOutputResponseBodyLen: config.C.Middleware.Logger.MaxOutputResponseBodyLen,
	}))

	g.Use(middleware.CopyBodyWithConfig(middleware.CopyBodyConfig{
		SkippedPathPrefixes: config.C.Middleware.CopyBody.SkippedPathPrefixes,
		MaxContentLen:       config.C.Middleware.CopyBody.MaxContentLen,
	}))

	g.Use(middleware.CORSWithConfig(middleware.CORSConfig{
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

	g.Use(middleware.AuthWithConfig(middleware.AuthConfig{
		Disable:             config.C.Middleware.Auth.Disable,
		SkippedPathPrefixes: config.C.Middleware.Auth.SkippedPathPrefixes,
		DefaultUserID: func(c *gin.Context) string {
			return config.C.Dictionary.RootUser.ID
		},
		ParseUserID: func(c *gin.Context) (string, error) {
			ctx := c.Request.Context()
			sub, err := injector.Auth.ParseSubject(ctx, utilx.GetToken(c))
			if err != nil {
				if err == jwtauth.ErrInvalidToken {
					return "", utilx.ErrInvalidToken
				}
				return "", err
			} else if sub == config.C.Dictionary.RootUser.ID {
				return sub, nil
			}

			joinRole, ok, err := injector.Cache.Get(ctx, utilx.CacheNSForUserRole, sub)
			if err != nil {
				return "", err
			} else if ok {
				c.Request = c.Request.WithContext(contextx.NewRoleIDs(ctx, strings.Split(joinRole, ",")))
				return sub, nil
			}

			roleIDs, err := injector.RBAC.UserBiz.GetRoleIDs(ctx, sub)
			if err != nil {
				return "", err
			}

			err = injector.Cache.Set(ctx, utilx.CacheNSForUserRole, sub, strings.Join(roleIDs, ","))
			if err != nil {
				return "", err
			}
			c.Request = c.Request.WithContext(contextx.NewRoleIDs(ctx, roleIDs))
			return sub, nil
		},
	}))

	g.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Enable:              config.C.Middleware.RateLimiter.Enable,
		SkippedPathPrefixes: config.C.Middleware.RateLimiter.SkippedPathPrefixes,
		Period:              config.C.Middleware.RateLimiter.Period,
		MaxRequestsPerIP:    config.C.Middleware.RateLimiter.MaxRequestsPerIP,
		MaxRequestsPerUser:  config.C.Middleware.RateLimiter.MaxRequestsPerUser,
		StoreType:           config.C.Middleware.RateLimiter.Store.Type,
		MemoryStoreConfig: middleware.RateLimiterMemoryConfig{
			Expiration:      time.Second * time.Duration(config.C.Middleware.RateLimiter.Store.Memory.Expiration),
			CleanupInterval: time.Second * time.Duration(config.C.Middleware.RateLimiter.Store.Memory.CleanupInterval),
		},
		RedisStoreConfig: middleware.RateLimiterRedisConfig{
			Addr:     config.C.Middleware.RateLimiter.Store.Redis.Addr,
			Password: config.C.Middleware.RateLimiter.Store.Redis.Password,
			DB:       config.C.Middleware.RateLimiter.Store.Redis.DB,
			Username: config.C.Middleware.RateLimiter.Store.Redis.Username,
		},
	}))

	g.Use(middleware.CasbinWithConfig(middleware.CasbinConfig{
		SkippedPathPrefixes: config.C.Middleware.Casbin.SkippedPathPrefixes,
		Skipper: func(c *gin.Context) bool {
			if config.C.Middleware.Casbin.Disable ||
				(config.C.Middleware.CORS.Enable && c.Request.Method == http.MethodOptions) ||
				contextx.FromUserID(c.Request.Context()) == config.C.Dictionary.RootUser.ID {
				return true
			}
			return false
		},
		GetEnforcer: func(c *gin.Context) *casbin.Enforcer {
			return injector.RBAC.Casbinx.GetEnforcer()
		},
		GetRoleIDs: func(c *gin.Context) []string {
			return contextx.FromRoleIDs(c.Request.Context())
		},
	}))
}

func initCaptcha() {
	if config.C.Util.Captcha.CacheType == "redis" {
		captcha.SetCustomStore(store.NewRedisStore(
			&redis.Options{
				Addr:     config.C.Util.Captcha.Redis.Addr,
				DB:       config.C.Util.Captcha.Redis.DB,
				Username: config.C.Util.Captcha.Redis.Username,
				Password: config.C.Util.Captcha.Redis.Password,
			},
			captcha.Expiration,
			&logger.PrintLogger{},
			config.C.Util.Captcha.Redis.KeyPrefix,
		))
	}
}
