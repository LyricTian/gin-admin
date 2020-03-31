// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package inject

import (
	"github.com/LyricTian/gin-admin/internal/app/api"
	"github.com/LyricTian/gin-admin/internal/app/bll/impl/bll"
	"github.com/LyricTian/gin-admin/internal/app/model/impl/gorm/model"
	"github.com/LyricTian/gin-admin/internal/app/module/adapter"
	"github.com/LyricTian/gin-admin/internal/app/router"
	"github.com/google/wire"
)

//  BuildInjector 生成注入器
func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		InitGormDB,
		InitAuth,
		model.AllSet,
		bll.AllSet,
		api.AllSet,
		router.RouterSet,
		InitHTTPEngine,
		adapter.CasbinAdapterSet,
		InitCasbin,
		MenuSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
