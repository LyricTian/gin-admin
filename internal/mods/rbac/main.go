package rbac

import (
	"context"
	"path/filepath"

	"github.com/LyricTian/gin-admin/v10/internal/config"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/api"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/pkg/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RBAC struct {
	DB       *gorm.DB
	MenuAPI  *api.Menu
	RoleAPI  *api.Role
	UserAPI  *api.User
	LoginAPI *api.Login
	Casbinx  *Casbinx
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
	if config.C.Storage.DB.AutoMigrate {
		if err := a.AutoMigrate(ctx); err != nil {
			return err
		}
	}

	if err := a.Casbinx.Load(ctx); err != nil {
		return err
	}

	if name := config.C.General.MenuFile; name != "" {
		fullPath := filepath.Join(config.C.General.WorkDir, name)
		if err := a.MenuAPI.MenuBIZ.InitFromFile(ctx, fullPath); err != nil {
			logging.Context(ctx).Error("failed to init menu data", zap.Error(err), zap.String("file", fullPath))
		}
	}

	return nil
}

func (a *RBAC) RegisterV1Routers(ctx context.Context, v1 *gin.RouterGroup) error {
	captcha := v1.Group("captcha")
	{
		captcha.GET("id", a.LoginAPI.GetCaptcha)
		captcha.GET("image", a.LoginAPI.ResponseCaptcha)
	}
	v1.POST("login", a.LoginAPI.Login)

	current := v1.Group("current")
	{
		current.POST("refresh-token", a.LoginAPI.RefreshToken)
		current.GET("user", a.LoginAPI.GetUserInfo)
		current.GET("menus", a.LoginAPI.QueryMenus)
		current.PUT("password", a.LoginAPI.UpdatePassword)
		current.PUT("user", a.LoginAPI.UpdateUser)
		current.POST("logout", a.LoginAPI.Logout)
	}
	menu := v1.Group("menus")
	{
		menu.GET("", a.MenuAPI.Query)
		menu.GET(":id", a.MenuAPI.Get)
		menu.POST("", a.MenuAPI.Create)
		menu.PUT(":id", a.MenuAPI.Update)
		menu.DELETE(":id", a.MenuAPI.Delete)
	}
	role := v1.Group("roles")
	{
		role.GET("", a.RoleAPI.Query)
		role.GET(":id", a.RoleAPI.Get)
		role.POST("", a.RoleAPI.Create)
		role.PUT(":id", a.RoleAPI.Update)
		role.DELETE(":id", a.RoleAPI.Delete)
	}
	user := v1.Group("users")
	{
		user.GET("", a.UserAPI.Query)
		user.GET(":id", a.UserAPI.Get)
		user.POST("", a.UserAPI.Create)
		user.PUT(":id", a.UserAPI.Update)
		user.DELETE(":id", a.UserAPI.Delete)
		user.PATCH(":id/reset-pwd", a.UserAPI.ResetPassword)
	}
	return nil
}

func (a *RBAC) Release(ctx context.Context) error {
	if err := a.Casbinx.Release(ctx); err != nil {
		return err
	}
	return nil
}
