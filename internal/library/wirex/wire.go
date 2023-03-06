//go:build wireinject
// +build wireinject

package wirex

// The build tag makes sure the stub is not built in the final build.

import (
	"context"

	"github.com/google/wire"

	"github.com/LyricTian/gin-admin/v10/internal/library/utilx"
	"github.com/LyricTian/gin-admin/v10/internal/mods"
) // end

func BuildInjector(ctx context.Context) (*Injector, func(), error) {
	wire.Build(
		InitCacher,
		InitDB,
		wire.NewSet(wire.Struct(new(utilx.Trans), "*")),
		wire.NewSet(wire.Struct(new(Injector), "*")),
		mods.Set,
	) // end
	return new(Injector), nil, nil
}
