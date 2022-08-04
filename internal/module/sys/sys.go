package sys

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/module/sys/api"
	"github.com/LyricTian/gin-admin/v9/internal/module/sys/biz"
	"github.com/LyricTian/gin-admin/v9/internal/module/sys/dao"
	"github.com/LyricTian/gin-admin/v9/internal/module/sys/typed"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// Collection of SYS wire providers
var Set = wire.NewSet(
	wire.Struct(new(SYS), "*"),
	wire.Struct(new(dao.DictionaryRepo), "*"),
	wire.Struct(new(biz.DictionaryBiz), "*"),
	wire.Struct(new(api.DictionaryAPI), "*"),
) // end

// SYS module is a SYS service
type SYS struct {
	DB            *gorm.DB
	DictionaryAPI *api.DictionaryAPI
} // end

func (a *SYS) Init(ctx context.Context) error {
	// Auto migrate tables for SYS
	if err := a.autoMigrate(ctx); err != nil {
		return err
	}

	return nil
}

func (a *SYS) autoMigrate(ctx context.Context) error {
	return a.DB.AutoMigrate(
		&typed.Dictionary{},
	) // end
}

func (a *SYS) RegisterAPI(ctx context.Context, group *gin.RouterGroup) {
	r := group.Group("sys")
	v1 := r.Group("v1")
	{
		gDictionary := v1.Group("dictionaries")
		{
			gDictionary.GET("", a.DictionaryAPI.Query)
			gDictionary.GET(":id", a.DictionaryAPI.Get)
			gDictionary.POST("", a.DictionaryAPI.Create)
			gDictionary.PUT(":id", a.DictionaryAPI.Update)
			gDictionary.DELETE(":id", a.DictionaryAPI.Delete)
		}
	} // end
}
