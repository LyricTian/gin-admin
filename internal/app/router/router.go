package router

import (
	"github.com/LyricTian/gin-admin/internal/app/api"
	"github.com/LyricTian/gin-admin/pkg/auth"
	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
)

// RouterSet 注入router
var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"))

// Router 路由管理器
type Router struct {
	Auth           auth.Auther
	CasbinEnforcer *casbin.SyncedEnforcer
	DemoAPI        *api.Demo
	LoginAPI       *api.Login
	MenuAPI        *api.Menu
	RoleAPI        *api.Role
	UserAPI        *api.User
}
