// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package app

import (
	"github.com/LyricTian/gin-admin/v8/internal/app/api"
	"github.com/LyricTian/gin-admin/v8/internal/app/dao"
	"github.com/LyricTian/gin-admin/v8/internal/app/module/adapter"
	"github.com/LyricTian/gin-admin/v8/internal/app/router"
	"github.com/LyricTian/gin-admin/v8/internal/app/service"
	"github.com/google/wire"
)

// BuildInjector 生成注入器
func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		InitGormDB,
		dao.RepoSet,
		InitAuth,
		InitCasbin,
		InitGinEngine,
		service.ServiceSet,
		api.APISet,
		router.RouterSet,
		adapter.CasbinAdapterSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
