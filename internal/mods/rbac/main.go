package rbac

import (
	"context"

	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/api"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RBAC struct {
	DB *gorm.DB

	MenuAPI *api.Menu
	RoleAPI *api.Role
	UserAPI *api.User
}

func (a *RBAC) AutoMigrate(ctx context.Context) error {
	return a.DB.AutoMigrate(

		new(schema.Menu),
		new(schema.MenuResource),
		new(schema.Role),
		new(schema.RoleMenu),
		new(schema.User),
		new(schema.UserRole),
	)
}

func (a *RBAC) Init(ctx context.Context) error {
	if err := a.AutoMigrate(ctx); err != nil {
		return err
	}
	return nil
}

func (a *RBAC) RegisterV1Routers(ctx context.Context, v1 *gin.RouterGroup) error {

	menu := v1.Group("/menus")
	{
		menu.GET("", a.MenuAPI.Query)
		menu.GET("/:id", a.MenuAPI.Get)
		menu.POST("", a.MenuAPI.Create)
		menu.PUT("", a.MenuAPI.Update)
		menu.DELETE("/:id", a.MenuAPI.Delete)
	}
	role := v1.Group("/roles")
	{
		role.GET("", a.RoleAPI.Query)
		role.GET("/:id", a.RoleAPI.Get)
		role.POST("", a.RoleAPI.Create)
		role.PUT("", a.RoleAPI.Update)
		role.DELETE("/:id", a.RoleAPI.Delete)
	}
	user := v1.Group("/users")
	{
		user.GET("", a.UserAPI.Query)
		user.GET("/:id", a.UserAPI.Get)
		user.POST("", a.UserAPI.Create)
		user.PUT("", a.UserAPI.Update)
		user.DELETE("/:id", a.UserAPI.Delete)
	}
	return nil
}
