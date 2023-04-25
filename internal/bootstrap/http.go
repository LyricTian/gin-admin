package bootstrap

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/LyricTian/gin-admin/v10/internal/config"
	"github.com/LyricTian/gin-admin/v10/internal/consts"
	"github.com/LyricTian/gin-admin/v10/internal/middlewares"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/internal/utils"
	"github.com/LyricTian/gin-admin/v10/internal/wirex"
	"github.com/LyricTian/gin-admin/v10/pkg/errors"
	"github.com/LyricTian/gin-admin/v10/pkg/jwtx"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func startHTTPServer(ctx context.Context, injector *wirex.Injector) (func(), error) {
	if config.C.General.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()

	e.GET("/health", func(c *gin.Context) {
		utils.ResOK(c)
	})

	e.Use(middlewares.RecoveryWithConfig(middlewares.RecoveryConfig{
		Skip: config.C.Middleware.Recovery.Skip,
	}))

	e.NoMethod(func(c *gin.Context) {
		utils.ResError(c, errors.MethodNotAllowed("", "Method not allowed"))
	})

	e.NoRoute(func(c *gin.Context) {
		utils.ResError(c, errors.NotFound("", "Not found"))
	})

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

	allowedPathPrefixes := injector.M.RouterPrefixes()
	e.Use(middlewares.TraceWithConfig(middlewares.TraceConfig{
		AllowedPathPrefixes: allowedPathPrefixes,
		SkippedPathPrefixes: config.C.Middleware.Trace.SkippedPathPrefixes,
		RequestHeaderKey:    config.C.Middleware.Trace.RequestHeaderKey,
		ResponseTraceKey:    config.C.Middleware.Trace.ResponseTraceKey,
	}))

	e.Use(middlewares.LoggerWithConfig(middlewares.LoggerConfig{
		AllowedPathPrefixes:      allowedPathPrefixes,
		SkippedPathPrefixes:      config.C.Middleware.Logger.SkippedPathPrefixes,
		MaxOutputRequestBodyLen:  config.C.Middleware.Logger.MaxOutputRequestBodyLen,
		MaxOutputResponseBodyLen: config.C.Middleware.Logger.MaxOutputResponseBodyLen,
	}))

	e.Use(middlewares.CopyBodyWithConfig(middlewares.CopyBodyConfig{
		AllowedPathPrefixes: allowedPathPrefixes,
		SkippedPathPrefixes: config.C.Middleware.CopyBody.SkippedPathPrefixes,
		MaxContentLen:       config.C.Middleware.CopyBody.MaxContentLen,
	}))

	e.Use(middlewares.AuthWithConfig(middlewares.AuthConfig{
		AllowedPathPrefixes: allowedPathPrefixes,
		SkippedPathPrefixes: config.C.Middleware.Auth.SkippedPathPrefixes,
		ParseUserID: func(c *gin.Context) (string, error) {
			rootID := config.C.General.Root.ID
			if config.C.Middleware.Auth.Disable {
				return rootID, nil
			}

			errInvalidToken := errors.Unauthorized("com.invalid.token", "Invalid access token")
			token := utils.GetToken(c)
			if token == "" {
				return "", errInvalidToken
			}

			ctx := c.Request.Context()
			ctx = utils.NewUserToken(ctx, token)

			userID, err := injector.Auth.ParseSubject(ctx, token)
			if err != nil {
				if err == jwtx.ErrInvalidToken {
					return "", errInvalidToken
				}
				return "", err
			} else if userID == rootID {
				c.Request = c.Request.WithContext(utils.NewIsRootUser(ctx))
				return userID, nil
			}

			userCacheVal, ok, err := injector.Cache.Get(ctx, consts.CacheNSForUser, userID)
			if err != nil {
				return "", err
			} else if ok {
				userCache := utils.ParseUserCache(userCacheVal)
				c.Request = c.Request.WithContext(utils.NewUserCache(ctx, userCache))
				return userID, nil
			}

			// Check user status, if not activated, force to logout
			user, err := injector.M.RBAC.UserAPI.UserBIZ.UserDAL.Get(ctx, userID, schema.UserQueryOptions{
				QueryOptions: utils.QueryOptions{SelectFields: []string{"status"}},
			})
			if err != nil {
				return "", err
			} else if user == nil || user.Status != schema.UserStatusActivated {
				return "", errInvalidToken
			}

			roleIDs, err := injector.M.RBAC.UserAPI.UserBIZ.GetRoleIDs(ctx, userID)
			if err != nil {
				return "", err
			}

			userCache := utils.UserCache{
				RoleIDs: roleIDs,
			}
			err = injector.Cache.Set(ctx, consts.CacheNSForUser, userID, userCache.String())
			if err != nil {
				return "", err
			}

			c.Request = c.Request.WithContext(utils.NewUserCache(ctx, userCache))
			return userID, nil
		},
	}))

	e.Use(middlewares.RateLimiterWithConfig(middlewares.RateLimiterConfig{
		Enable:              config.C.Middleware.RateLimiter.Enable,
		AllowedPathPrefixes: allowedPathPrefixes,
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

	e.Use(middlewares.CasbinWithConfig(middlewares.CasbinConfig{
		AllowedPathPrefixes: allowedPathPrefixes,
		SkippedPathPrefixes: config.C.Middleware.Casbin.SkippedPathPrefixes,
		Skipper: func(c *gin.Context) bool {
			if config.C.Middleware.Casbin.Disable ||
				utils.FromIsRootUser(c.Request.Context()) {
				return true
			}
			return false
		},
		GetEnforcer: func(c *gin.Context) *casbin.Enforcer {
			return injector.M.RBAC.Casbinx.GetEnforcer()
		},
		GetSubjects: func(c *gin.Context) []string {
			return utils.FromUserCache(c.Request.Context()).RoleIDs
		},
	}))

	// register routers
	if err := injector.M.RegisterRouters(ctx, e); err != nil {
		return nil, err
	}

	if !config.C.General.DisableSwagger {
		e.StaticFile("/openapi.json", filepath.Join(config.C.General.ConfigDir, "openapi.json"))
		e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

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
