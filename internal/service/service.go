package service

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/module/contextx"
	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(
	DemoSet,
) // end

func IsRootUser(ctx context.Context) bool {
	return config.C.IsRootUser(contextx.FromUserID(ctx))
}
