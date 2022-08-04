//go:build wireinject
// +build wireinject

package inject

// The build tag makes sure the stub is not built in the final build.

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/google/wire"
) // end

func BuildInjector(ctx context.Context) (*Injector, func(), error) {
	wire.Build(
		InitAuth,
		InitCache,
		InitDB,
		wire.NewSet(wire.Struct(new(Injector), "*")),
		utilx.TransRepoSet,
		rbac.Set,
	) // end
	return new(Injector), nil, nil
}
