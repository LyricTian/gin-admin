package web

import (
	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/web/middleware"
	"github.com/LyricTian/gin-admin/src/web/router"
	"github.com/gin-gonic/gin"
)

// Init 初始化所有服务
func Init(obj *inject.Object) *gin.Engine {
	gin.SetMode(config.GetRunMode())
	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	// 注册中间件
	apiPrefixes := []string{"/api/"}
	if dir := config.GetWWWDir(); dir != "" {
		app.Use(middleware.WWWMiddleware(dir, middleware.AllowPathPrefixSkipper(apiPrefixes...)))
	}

	app.Use(middleware.TraceMiddleware(middleware.NoAllowPathPrefixSkipper(apiPrefixes...)))
	app.Use(middleware.LoggerMiddleware(middleware.NoAllowPathPrefixSkipper(apiPrefixes...)))
	app.Use(middleware.RecoveryMiddleware())

	if config.IsSessionAuth() {
		app.Use(middleware.SessionMiddleware(obj, middleware.NoAllowPathPrefixSkipper(apiPrefixes...)))
	}

	// 注册/api路由
	router.APIHandler(app, obj)

	// 加载casbin策略数据
	// err := loadCasbinPolicyData(obj)
	// if err != nil {
	// 	panic("加载casbin策略数据发生错误：" + err.Error())
	// }

	return app
}

// 加载casbin策略数据，包括角色权限数据、用户角色数据
func loadCasbinPolicyData(obj *inject.Object) error {
	c := obj.CtlCommon

	err := c.RoleAPI.RoleBll.LoadAllPolicy()
	if err != nil {
		return err
	}

	err = c.UserAPI.UserBll.LoadAllPolicy()
	if err != nil {
		return err
	}
	return nil
}
