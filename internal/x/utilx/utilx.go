package utilx

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/config"
	"github.com/LyricTian/gin-admin/v9/internal/x/contextx"
)

func IsRootUser(ctx context.Context) bool {
	return contextx.FromUserID(ctx) == config.C.Dictionary.RootUser.ID
}
