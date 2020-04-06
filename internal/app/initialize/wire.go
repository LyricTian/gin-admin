// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package initialize

import (
	"github.com/LyricTian/gin-admin/internal/app/api"
	"github.com/LyricTian/gin-admin/internal/app/api/mock"
	"github.com/LyricTian/gin-admin/internal/app/bll/impl/bll"
	gormModel "github.com/LyricTian/gin-admin/internal/app/model/impl/gorm/model"
	mongoModel "github.com/LyricTian/gin-admin/internal/app/model/impl/mongo/model"
	"github.com/LyricTian/gin-admin/internal/app/module/adapter"
	"github.com/LyricTian/gin-admin/internal/app/router"
	"github.com/google/wire"
)

var buildContainer = []interface{}{}

//  BuildGormInjector 生成基于gorm实现的存储注入器
func BuildGormInjector() (*Injector, func(), error) {
	wire.Build(
		InitAuth,
		bll.BllSet,
		api.APISet,
		mock.MockSet,
		router.RouterSet,
		InitGinEngine,
		adapter.CasbinAdapterSet,
		InitCasbin,
		MenuSet,
		InjectorSet,
		InitGormDB,
		gormModel.ModelSet,
	)
	return new(Injector), nil, nil
}

// BuildMongoInjector 生成基于mongo实现的存储注入器
func BuildMongoInjector() (*Injector, func(), error) {
	wire.Build(
		InitAuth,
		bll.BllSet,
		api.APISet,
		mock.MockSet,
		router.RouterSet,
		InitGinEngine,
		adapter.CasbinAdapterSet,
		InitCasbin,
		MenuSet,
		InjectorSet,
		InitMongo,
		mongoModel.ModelSet,
	)
	return new(Injector), nil, nil
}
