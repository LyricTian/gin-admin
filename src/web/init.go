package web

import (
	"fmt"

	"github.com/LyricTian/gin-admin/src/service/mysql"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/LyricTian/gin-admin/src/web/router"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Init 初始化所有服务
func Init(db *mysql.DB, enforcer *casbin.Enforcer, ctlCommon *ctl.Common) *gin.Engine {
	gin.SetMode(viper.GetString("run_mode"))
	app := gin.New()

	app.NoMethod(context.WrapContext(func(ctx *context.Context) {
		ctx.ResError(fmt.Errorf("方法不允许"), 405)
	}))

	app.NoRoute(context.WrapContext(func(ctx *context.Context) {
		ctx.ResError(fmt.Errorf("资源不存在"), 404)
	}))

	// 注册中间件
	apiPrefixes := []string{"/api/"}
	app.Use(router.TraceMiddleware(apiPrefixes...))
	app.Use(router.LoggerMiddleware(apiPrefixes, "/api/v1/loggers"))
	app.Use(router.RecoveryMiddleware())
	app.Use(router.SessionMiddleware(db, apiPrefixes...))

	// 注册/api/v1路由
	router.APIV1Handler(app, enforcer, ctlCommon)

	// 加载casbin策略数据
	err := loadCasbinPolicyData(ctlCommon)
	if err != nil {
		panic("加载casbin策略数据发生错误：" + err.Error())
	}

	return app
}

// 加载casbin策略数据，包括角色权限数据、用户角色数据
func loadCasbinPolicyData(ctlCommon *ctl.Common) error {
	err := ctlCommon.RoleAPI.RoleBll.LoadAllPolicy()
	if err != nil {
		return err
	}

	err = ctlCommon.UserAPI.UserBll.LoadAllPolicy()
	if err != nil {
		return err
	}
	return nil
}
