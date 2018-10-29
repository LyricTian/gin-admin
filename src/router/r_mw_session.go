package router

import (
	"fmt"
	"gin-admin/src/service/mysql"
	"gin-admin/src/util"
	"strings"

	"github.com/go-session/gin-session"

	mysession "github.com/go-session/mysql"

	"github.com/go-session/session"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

// SessionMiddleware session中间件
func SessionMiddleware(db *mysql.DB, prefixes ...string) gin.HandlerFunc {
	sessionConfig := viper.GetStringMap("session")

	var opts []session.Option
	opts = append(opts, session.SetCookieName(util.T(sessionConfig["cookie_name"]).String()))
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

	mode := viper.GetString("run_mode")
	ginConfig := ginsession.DefaultConfig
	ginConfig.Skipper = func(ctx *gin.Context) bool {
		if mode == util.DebugMode {
			return true
		}

		for _, prefix := range prefixes {
			if strings.HasPrefix(ctx.Request.URL.Path, prefix) {
				return false
			}
		}

		return true
	}

	return ginsession.NewWithConfig(ginConfig, opts...)
}
