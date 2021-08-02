package dao

import (
	"strings"

	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/internal/app/config"
	"github.com/LyricTian/gin-admin/v8/internal/app/dao/menu"
	"github.com/LyricTian/gin-admin/v8/internal/app/dao/role"
	"github.com/LyricTian/gin-admin/v8/internal/app/dao/user"
	"github.com/LyricTian/gin-admin/v8/internal/app/dao/util"
) // end

// RepoSet repo injection
var RepoSet = wire.NewSet(
	util.TransSet,
	menu.MenuActionResourceSet,
	menu.MenuActionSet,
	menu.MenuSet,
	role.RoleMenuSet,
	role.RoleSet,
	user.UserRoleSet,
	user.UserSet,
) // end

// Define repo type alias
type (
	TransRepo              = util.Trans
	MenuActionResourceRepo = menu.MenuActionResourceRepo
	MenuActionRepo         = menu.MenuActionRepo
	MenuRepo               = menu.MenuRepo
	RoleMenuRepo           = role.RoleMenuRepo
	RoleRepo               = role.RoleRepo
	UserRoleRepo           = user.UserRoleRepo
	UserRepo               = user.UserRepo
) // end

// Auto migration for given models
func AutoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	return db.AutoMigrate(
		new(menu.MenuActionResource),
		new(menu.MenuAction),
		new(menu.Menu),
		new(role.RoleMenu),
		new(role.Role),
		new(user.UserRole),
		new(user.User),
	) // end
}
