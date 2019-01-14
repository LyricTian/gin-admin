package middleware

import (
	"fmt"

	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/go-session/gorm"

	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/gin-gonic/gin"
	"github.com/go-session/gin-session"
	"github.com/go-session/session"
	"github.com/spf13/viper"
)

// SessionMiddleware session中间件
func SessionMiddleware(obj *inject.Object) gin.HandlerFunc {
	var config struct {
		HeaderName  string `toml:"header_name" yaml:"header_name" json:"header_name"`
		Sign        string `toml:"sign" yaml:"sign" json:"sign"`
		Expired     int64  `toml:"expired" yaml:"expired" json:"expired"`
		EnableStore bool   `toml:"enable_store" yaml:"enable_store" json:"enable_store"`
	}

	err := viper.UnmarshalKey("session", &config)
	if err != nil {
		panic(err.Error())
	}

	var opts []session.Option
	opts = append(opts, session.SetEnableSetCookie(false))
	opts = append(opts, session.SetEnableSIDInURLQuery(false))
	opts = append(opts, session.SetEnableSIDInHTTPHeader(true))
	opts = append(opts, session.SetSessionNameInHTTPHeader(config.HeaderName))
	opts = append(opts, session.SetSign([]byte(config.Sign)))
	opts = append(opts, session.SetExpired(config.Expired))

	if config.EnableStore {
		if mode := viper.GetString("db_mode"); mode == "gorm" && obj.GormDB != nil {
			opts = append(opts, session.SetStore(gormsession.NewDefaultStore(obj.GormDB)))
		}
	}

	ginConfig := ginsession.DefaultConfig
	ginConfig.ErrorHandleFunc = func(c *gin.Context, err error) {
		ctx := context.New(c)
		if err == session.ErrInvalidSessionID {
			ctx.ResError(util.NewBadRequestError("无效的会话"))
			return
		}
		logger.Start(ctx.CContext()).Errorf("服务器会话发生错误: %s", err.Error())
		ctx.ResError(util.NewInternalServerError("服务器会话发生错误"))
	}

	return ginsession.NewWithConfig(ginConfig, opts...)
}

// VerifySessionMiddleware 验证session中间件
func VerifySessionMiddleware(skipPrefixes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.New(c)
		store := ginsession.FromContext(c)
		userID, ok := store.Get(util.SessionKeyUserID)

		// 调试模式使用root用户
		if viper.GetString("run_mode") == util.DebugMode {
			if !ok || userID == nil {
				userID = viper.GetString("system_root_user")
			}
			c.Set(util.ContextKeyUserID, userID)
			c.Next()
			return
		}

		if !ok || userID == nil {
			p := fmt.Sprintf("%s%s", c.Request.Method, c.Request.URL.Path)
			if util.CheckPrefix(p, skipPrefixes...) {
				c.Next()
				return
			}
			ctx.ResError(util.NewUnauthorizedError("用户未登录"), 9999)
			return
		}
		c.Set(util.ContextKeyUserID, userID)
		c.Next()
	}
}
