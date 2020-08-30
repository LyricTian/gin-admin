// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package app

import (
	"github.com/LyricTian/gin-admin/v7/internal/app/api"
	// "github.com/LyricTian/gin-admin/v7/internal/app/api/mock"
	"github.com/LyricTian/gin-admin/v7/internal/app/bll"
	"github.com/LyricTian/gin-admin/v7/internal/app/module/adapter"
	"github.com/LyricTian/gin-admin/v7/internal/app/router"
	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v7/internal/app/model/gormx/model"
)

// BuildInjector 生成注入器
func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		// mock.MockSet,
		InitGormDB,
		model.ModelSet,
		InitAuth,
		InitCasbin,
		InitGinEngine,
		bll.BllSet,
		api.APISet,
		router.RouterSet,
		adapter.CasbinAdapterSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
