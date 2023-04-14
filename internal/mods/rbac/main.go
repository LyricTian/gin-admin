package rbac

import (
	"context"

	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/api"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RBAC struct {
	DB          *gorm.DB
	ResourceAPI *api.Resource
}

func (a *RBAC) AutoMigrate(ctx context.Context) error {
	return a.DB.AutoMigrate(
		new(schema.Resource),
	)
}

func (a *RBAC) Init(ctx context.Context) error {
	if err := a.AutoMigrate(ctx); err != nil {
		return err
	}
	return nil
}

func (a *RBAC) RegisterV1Routers(ctx context.Context, v1 *gin.RouterGroup) error {
	resource := v1.Group("resources")
	{
		resource.GET("", a.ResourceAPI.Query)
		resource.POST("", a.ResourceAPI.Create)
		resource.GET(":id", a.ResourceAPI.Get)
		resource.PUT(":id", a.ResourceAPI.Update)
		resource.DELETE(":id", a.ResourceAPI.Delete)
	}
	return nil
}
