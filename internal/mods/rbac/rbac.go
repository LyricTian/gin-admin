package rbac

import (
	"context"

	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/api"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/dao"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RBAC struct {
	DB          *gorm.DB
	ResourceAPI *api.Resource
}

func (a *RBAC) dbMigrate(ctx context.Context) error {
	return a.DB.AutoMigrate(
		dao.GetResourceDB(ctx, a.DB).Statement.Model,
	)
}

func (a *RBAC) Init(ctx context.Context) error {
	if err := a.dbMigrate(ctx); err != nil {
		return err
	}
	return nil
}

func (a *RBAC) RegisterAPIs(ctx context.Context, gm map[string]*gin.RouterGroup) error {
	apiV1 := gm["/api/v1"]
	{
		gResource := apiV1.Group("resources")
		{
			gResource.GET("", a.ResourceAPI.Query)
			gResource.POST("", a.ResourceAPI.Create)
			gResource.GET(":id", a.ResourceAPI.Get)
			gResource.PUT(":id", a.ResourceAPI.Update)
			gResource.DELETE(":id", a.ResourceAPI.Delete)
		}
	}

	return nil
}
