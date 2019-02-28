package router

import (
	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/web/auth"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/LyricTian/gin-admin/src/web/middleware"
	"github.com/gin-gonic/gin"
)

// APIHandler /api路由
func APIHandler(app *gin.Engine, obj *inject.Object) {
	g := app.Group("/api")

	// 用户授权(session/jwt)
	g.Use(auth.Entry(auth.SkipperFunc(middleware.NoAllowPathPrefixSkipper("/api"))))

	// 用户身份授权
	g.Use(middleware.UserAuthMiddleware(
		middleware.AllowMethodAndPathPrefixSkipper(
			context.JoinRouter("GET", "/api/v1/login"),
			context.JoinRouter("POST", "/api/v1/login"),
		),
	))

	// casbin权限校验中间件
	g.Use(middleware.CasbinMiddleware(obj.Enforcer,
		middleware.AllowMethodAndPathPrefixSkipper(
			context.JoinRouter("GET", "/api/v1/login"),
			context.JoinRouter("POST", "/api/v1/login"),
			context.JoinRouter("PUT", "/api/v1/current/password"),
			context.JoinRouter("GET", "/api/v1/current/user"),
			context.JoinRouter("GET", "/api/v1/current/menutree"),
		),
	))

	// 请求频率限制中间件
	if config.GetRateLimiterConfig().Enable {
		g.Use(middleware.RateLimiterMiddleware(obj.RateLimiter))
	}

	c := obj.CtlCommon

	// 注册路由

	// 注册/demos路由
	GET(g, "/v1/demos", c.DemoCtl.Query, SetName("查询示例数据"))
	GET(g, "/v1/demos/:id", c.DemoCtl.Get, SetName("查询指定示例数据"))
	POST(g, "/v1/demos", c.DemoCtl.Create, SetName("创建示例数据"))
	PUT(g, "/v1/demos/:id", c.DemoCtl.Update, SetName("更新示例数据"))
	DELETE(g, "/v1/demos/:id", c.DemoCtl.Delete, SetName("删除示例数据"))
	DELETE(g, "/v1/demos", c.DemoCtl.DeleteMany, SetName("删除多条示例数据"))
	PATCH(g, "/v1/demos/:id/enable", c.DemoCtl.Enable, SetName("启用示例数据"))
	PATCH(g, "/v1/demos/:id/disable", c.DemoCtl.Disable, SetName("禁用示例数据"))

	// 注册登录模块路由
	GET(g, "/v1/login/captchaid", c.LoginCtl.GetCaptchaID, SetName("获取验证码ID"))
	GET(g, "/v1/login/captcha", c.LoginCtl.GetCaptcha, SetName("获取图形验证码"))
	POST(g, "/v1/login", c.LoginCtl.Login, SetName("用户登录"))
	POST(g, "/v1/login/exit", c.LoginCtl.Logout, SetName("用户登出"))
	PUT(g, "/v1/current/password", c.LoginCtl.UpdatePassword, SetName("更新个人密码"))
	GET(g, "/v1/current/user", c.LoginCtl.GetUserInfo, SetName("获取当前用户信息"))
	GET(g, "/v1/current/menutree", c.LoginCtl.QueryUserMenuTree, SetName("查询当前用户菜单树"))

	// 注册/menus路由
	GET(g, "/v1/menus", c.MenuCtl.Query, SetName("查询菜单数据"))
	GET(g, "/v1/menus/:id", c.MenuCtl.Get, SetName("查询指定菜单数据"))
	POST(g, "/v1/menus", c.MenuCtl.Create, SetName("创建菜单数据"))
	PUT(g, "/v1/menus/:id", c.MenuCtl.Update, SetName("更新菜单数据"))
	DELETE(g, "/v1/menus/:id", c.MenuCtl.Delete, SetName("删除菜单数据"))
	DELETE(g, "/v1/menus", c.MenuCtl.DeleteMany, SetName("删除多条菜单数据"))

	// 注册/resources路由
	GET(g, "/v1/resources", c.ResourceCtl.Query, SetName("查询资源数据"))
	GET(g, "/v1/resources/:id", c.ResourceCtl.Get, SetName("查询指定资源数据"))
	POST(g, "/v1/resources", c.ResourceCtl.Create, SetName("创建资源数据"))
	PUT(g, "/v1/resources/:id", c.ResourceCtl.Update, SetName("更新资源数据"))
	DELETE(g, "/v1/resources/:id", c.ResourceCtl.Delete, SetName("删除资源数据"))
	DELETE(g, "/v1/resources", c.ResourceCtl.DeleteMany, SetName("删除多条资源数据"))

	// 注册/roles路由
	GET(g, "/v1/roles", c.RoleCtl.Query, SetName("查询角色数据"))
	GET(g, "/v1/roles/:id", c.RoleCtl.Get, SetName("查询指定角色数据"))
	POST(g, "/v1/roles", c.RoleCtl.Create, SetName("创建角色数据"))
	PUT(g, "/v1/roles/:id", c.RoleCtl.Update, SetName("更新角色数据"))
	DELETE(g, "/v1/roles/:id", c.RoleCtl.Delete, SetName("删除角色数据"))
	DELETE(g, "/v1/roles", c.RoleCtl.DeleteMany, SetName("删除多条角色数据"))

	// 注册/users路由
	GET(g, "/v1/users", c.UserCtl.Query, SetName("查询用户数据"))
	GET(g, "/v1/users/:id", c.UserCtl.Get, SetName("查询指定用户数据"))
	POST(g, "/v1/users", c.UserCtl.Create, SetName("创建用户数据"))
	PUT(g, "/v1/users/:id", c.UserCtl.Update, SetName("更新用户数据"))
	DELETE(g, "/v1/users/:id", c.UserCtl.Delete, SetName("删除用户数据"))
	DELETE(g, "/v1/users", c.UserCtl.DeleteMany, SetName("删除多条用户数据"))
	PATCH(g, "/v1/users/:id/enable", c.UserCtl.Enable, SetName("启用用户数据"))
	PATCH(g, "/v1/users/:id/disable", c.UserCtl.Disable, SetName("禁用用户数据"))
}
