package middleware

import (
	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/gin-gonic/gin"
	"github.com/go-session/gin-session"
	"github.com/go-session/gorm"
	"github.com/go-session/session"
)

// SessionMiddleware session中间件
func SessionMiddleware(obj *inject.Object, skipper SkipperFunc) gin.HandlerFunc {
	options := config.GetSessionConfig()

	var opts []session.Option
	opts = append(opts, session.SetEnableSetCookie(false))
	opts = append(opts, session.SetEnableSIDInURLQuery(false))
	opts = append(opts, session.SetEnableSIDInHTTPHeader(true))
	opts = append(opts, session.SetSessionNameInHTTPHeader(options.HeaderName))
	opts = append(opts, session.SetSign([]byte(options.Sign)))
	opts = append(opts, session.SetExpired(options.Expired))

	if options.EnableStore {
		if config.IsGormDB() && obj.GormDB != nil {
			tableName := config.GetGormTablePrefix() + "session"
			store := gormsession.NewStoreWithDB(obj.GormDB.DB, tableName, 0)
			opts = append(opts, session.SetStore(store))
		}
	}

	ginConfig := ginsession.DefaultConfig
	ginConfig.Skipper = skipper
	ginConfig.ErrorHandleFunc = func(c *gin.Context, err error) {
		ctx := context.New(c)
		if err == session.ErrInvalidSessionID {
			ctx.ResError(errors.NewBadRequestError("无效的会话"))
			return
		}
		logger.StartSpan(ctx.CContext(), "session中间件", "SessionMiddleware").Errorf("服务器会话发生错误: %s", err.Error())
		ctx.ResError(errors.NewInternalServerError())
	}

	return ginsession.NewWithConfig(ginConfig, opts...)
}

// VerifySessionMiddleware 验证session中间件
func VerifySessionMiddleware(skipper SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if skipper != nil && skipper(c) {
			c.Next()
			return
		}

		ctx := context.New(c)
		store := ginsession.FromContext(c)
		userID, ok := store.Get(context.ContextKeyUserID)
		if !ok || userID == nil {
			if !config.IsDebugMode() {
				ctx.ResError(errors.NewUnauthorizedError("用户未登录"), 9999)
				return
			}
			// 调试模式下，如果用户未登录，则使用root用户
			userID = config.GetRootUser().UserName
		}

		c.Set(context.ContextKeyUserID, userID)
		c.Next()
	}
}
