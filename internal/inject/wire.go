//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package inject

import (
	"context"
	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v9/internal/api"
	"github.com/LyricTian/gin-admin/v9/internal/dao"
	"github.com/LyricTian/gin-admin/v9/internal/router"
	"github.com/LyricTian/gin-admin/v9/internal/service"
)

func BuildInjector(ctx context.Context) (*Injector, func(), error) {
	wire.Build(
		InitGormDB,
		InitCache,
		InitJWTAuth,
		InitCasbin,
		InitEngine,
		dao.RepoSet,
		service.ServiceSet,
		api.APISet,
		router.RouterSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
