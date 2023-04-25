package rbac

import (
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/api"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/biz"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/dal"
	"github.com/google/wire"
)

// Collection of wire providers
var Set = wire.NewSet(
	wire.Struct(new(RBAC), "*"),
	wire.Struct(new(Casbinx), "*"),
	wire.Struct(new(dal.Menu), "*"),
	wire.Struct(new(biz.Menu), "*"),
	wire.Struct(new(api.Menu), "*"),
	wire.Struct(new(dal.MenuResource), "*"),
	wire.Struct(new(dal.Role), "*"),
	wire.Struct(new(biz.Role), "*"),
	wire.Struct(new(api.Role), "*"),
	wire.Struct(new(dal.RoleMenu), "*"),
	wire.Struct(new(dal.User), "*"),
	wire.Struct(new(biz.User), "*"),
	wire.Struct(new(api.User), "*"),
	wire.Struct(new(dal.UserRole), "*"),
	wire.Struct(new(biz.Login), "*"),
	wire.Struct(new(api.Login), "*"),
)
