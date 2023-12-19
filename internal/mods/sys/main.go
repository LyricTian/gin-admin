package sys

import (
	"context"

	"github.com/LyricTian/gin-admin/v10/internal/mods/sys/api"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SYS struct {
	DB        *gorm.DB
	LoggerAPI *api.Logger
}

func (a *SYS) AutoMigrate(ctx context.Context) error {
	return a.DB.AutoMigrate()
}

func (a *SYS) Init(ctx context.Context) error {
	if err := a.AutoMigrate(ctx); err != nil {
		return err
	}
	return nil
}

func (a *SYS) RegisterV1Routers(ctx context.Context, v1 *gin.RouterGroup) error {
	logger := v1.Group("loggers")
	{
		logger.GET("", a.LoggerAPI.Query)
	}
	return nil
}

func (a *SYS) Release(ctx context.Context) error {
	return nil
}
