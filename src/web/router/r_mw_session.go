package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/LyricTian/gin-admin/src/service/mysql"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/gin-gonic/gin"
	"github.com/go-session/gin-session"
	mysession "github.com/go-session/mysql"
	"github.com/go-session/session"
	"github.com/spf13/viper"
)

// SessionMiddleware session中间件
func SessionMiddleware(db *mysql.DB, allowPrefixes ...string) gin.HandlerFunc {
	sessionConfig := viper.GetStringMap("session")

	var opts []session.Option
	opts = append(opts, session.SetEnableSetCookie(false))
	opts = append(opts, session.SetEnableSIDInURLQuery(false))
	opts = append(opts, session.SetEnableSIDInHTTPHeader(true))
	opts = append(opts, session.SetSessionNameInHTTPHeader(util.T(sessionConfig["header_name"]).String()))
	opts = append(opts, session.SetSign(util.T(sessionConfig["sign"]).Bytes()))
	opts = append(opts, session.SetDomain(util.T(sessionConfig["domain"]).String()))
	opts = append(opts, session.SetCookieLifeTime(util.T(sessionConfig["cookie_life_time"]).Int()))
	opts = append(opts, session.SetExpired(util.T(sessionConfig["expired"]).Int64()))

	if util.T(sessionConfig["store"]).String() == "mysql" {
		tableName := fmt.Sprintf("%s_%s",
			util.T(viper.GetStringMap("mysql")["table_prefix"]).String(),
			util.T(sessionConfig["table"]).String())
		opts = append(opts, session.SetStore(mysession.NewStoreWithDB(db.Db, tableName, 0)))
	}

	ginConfig := ginsession.DefaultConfig
	ginConfig.Skipper = func(c *gin.Context) bool {
		return sessionSkipper(c, allowPrefixes...)
	}
	ginConfig.ErrorHandleFunc = func(c *gin.Context, err error) {
		ctx := context.NewContext(c)
		ctx.ResError(err, http.StatusInternalServerError)
	}

	return ginsession.NewWithConfig(ginConfig, opts...)
}

func sessionSkipper(c *gin.Context, prefixes ...string) bool {
	if viper.GetString("run_mode") == util.DebugMode {
		return true
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(c.Request.URL.Path, prefix) {
			return false
		}
	}

	return true
}

// VerifySessionMiddleware 验证session中间件
func VerifySessionMiddleware(allowPrefixes, skipPrefixes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if sessionSkipper(c, allowPrefixes...) {
			c.Next()
			return
		}

		ctx := context.NewContext(c)
		store := ginsession.FromContext(c)
		userID, ok := store.Get(util.SessionKeyUserID)
		if !ok || userID == nil {

			for _, prefix := range skipPrefixes {
				if strings.HasPrefix(c.Request.URL.Path, prefix) {
					c.Next()
					return
				}
			}

			ctx.ResError(fmt.Errorf("用户未登录"), http.StatusUnauthorized, 9999)
			return
		}
		c.Set(util.SessionKeyUserID, userID)
		c.Next()
	}
}
