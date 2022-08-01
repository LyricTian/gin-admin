//go:build wireinject
// +build wireinject

package initialize

// The build tag makes sure the stub is not built in the final build.

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"

	"github.com/google/wire"
)

func BuildInjector(ctx context.Context) (*Injector, func(), error) {
	wire.Build(
		InitAuth,
		InitCache,
		InitDB,
		utilx.TransRepoSet,
		rbac.Set,
		wire.NewSet(wire.Struct(new(Injector), "*")),
	)
	return new(Injector), nil, nil
}
