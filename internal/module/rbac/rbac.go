package rbac

import (
	"context"
	"path/filepath"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/api"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/biz"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/dao"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Collection of RBAC wire providers
var Set = wire.NewSet(
	wire.Struct(new(RBAC), "*"),
	wire.Struct(new(Casbinx), "*"),
	wire.Struct(new(dao.MenuRepo), "*"),
	wire.Struct(new(dao.MenuActionRepo), "*"),
	wire.Struct(new(dao.MenuActionResourceRepo), "*"),
	wire.Struct(new(biz.MenuBiz), "*"),
	wire.Struct(new(api.MenuAPI), "*"),
	wire.Struct(new(dao.RoleRepo), "*"),
	wire.Struct(new(dao.RoleMenuRepo), "*"),
	wire.Struct(new(biz.RoleBiz), "*"),
	wire.Struct(new(api.RoleAPI), "*"),
	wire.Struct(new(dao.UserRepo), "*"),
	wire.Struct(new(dao.UserRoleRepo), "*"),
	wire.Struct(new(biz.UserBiz), "*"),
	wire.Struct(new(api.UserAPI), "*"),
	wire.Struct(new(biz.LoginBiz), "*"),
	wire.Struct(new(api.LoginAPI), "*"),
) // end

// RBAC module is a RBAC service
type RBAC struct {
	DB       *gorm.DB
	Casbinx  *Casbinx
	MenuBiz  *biz.MenuBiz
	UserBiz  *biz.UserBiz
	MenuAPI  *api.MenuAPI
	RoleAPI  *api.RoleAPI
	UserAPI  *api.UserAPI
	LoginAPI *api.LoginAPI
} // end

func (a *RBAC) Init(ctx context.Context) error {
	// Auto migrate tables for RBAC
	if err := a.autoMigrate(ctx); err != nil {
		return err
	}

	if !config.C.General.DisableInitMenu {
		// Initialize menu data from json file
		err := a.MenuBiz.InitFromJSON(ctx, filepath.Join(config.C.General.ConfigDir, "menu.json"))
		if err != nil {
			logger.Context(ctx).Error("Failed to init menu from json file", zap.Error(err))
			return err
		}
	}

	if !config.C.Middleware.Casbin.Disable {
		// Initialize casbin
		if err := a.Casbinx.Load(ctx); err != nil {
			return err
		}
		go a.Casbinx.AutoLoad(ctx)
	}

	return nil
}

func (a *RBAC) autoMigrate(ctx context.Context) error {
	return a.DB.AutoMigrate(
		&typed.Menu{},
		&typed.MenuAction{},
		&typed.MenuActionResource{},
		&typed.Role{},
		&typed.RoleMenu{},
		&typed.User{},
		&typed.UserRole{},
	) // end
}

func (a *RBAC) RegisterAPI(ctx context.Context, group *gin.RouterGroup) {
	r := group.Group("rbac")
	v1 := r.Group("v1")
	{
		gLogin := v1.Group("login")
		{
			gLogin.POST("", a.LoginAPI.Login)
			gLogin.GET("captchaid", a.LoginAPI.GetCaptchaID)
			gLogin.GET("captcha", a.LoginAPI.WriteCaptchaImage)
		}

		gCurrent := v1.Group("current")
		{
			gCurrent.POST("logout", a.LoginAPI.Logout)
			gCurrent.POST("refreshtoken", a.LoginAPI.RefreshToken)
			gCurrent.PUT("password", a.LoginAPI.UpdatePassword)
			gCurrent.GET("user", a.LoginAPI.GetCurrentUser)
			gCurrent.GET("menus", a.LoginAPI.QueryPrivilegeMenus)
		}

		gMenu := v1.Group("menus")
		{
			gMenu.GET("", a.MenuAPI.Query)
			gMenu.GET(":id", a.MenuAPI.Get)
			gMenu.POST("", a.MenuAPI.Create)
			gMenu.PUT(":id", a.MenuAPI.Update)
			gMenu.DELETE(":id", a.MenuAPI.Delete)
			gMenu.PUT(":id/status", a.MenuAPI.UpdateStatus)
		}

		gRole := v1.Group("roles")
		{
			gRole.GET("", a.RoleAPI.Query)
			gRole.GET(":id", a.RoleAPI.Get)
			gRole.POST("", a.RoleAPI.Create)
			gRole.PUT(":id", a.RoleAPI.Update)
			gRole.DELETE(":id", a.RoleAPI.Delete)
			gRole.PUT(":id/status", a.RoleAPI.UpdateStatus)
		}

		gUser := v1.Group("users")
		{
			gUser.GET("", a.UserAPI.Query)
			gUser.GET(":id", a.UserAPI.Get)
			gUser.POST("", a.UserAPI.Create)
			gUser.PUT(":id", a.UserAPI.Update)
			gUser.DELETE(":id", a.UserAPI.Delete)
			gUser.PUT(":id/status", a.UserAPI.UpdateStatus)
		}
	} // end
}
