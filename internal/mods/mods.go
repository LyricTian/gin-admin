package mods

import (
	"context"

	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Collection of wire providers
var Set = wire.NewSet(
	wire.Struct(new(Mods), "*"),
	rbac.Set,
) // end

type Mods struct {
	RBAC *rbac.RBAC
}

func (a *Mods) Init(ctx context.Context) error {
	if err := a.RBAC.Init(ctx); err != nil {
		return err
	}

	return nil
}

func (a *Mods) RegisterAPIs(ctx context.Context, gm map[string]*gin.RouterGroup) error {
	if err := a.RBAC.RegisterAPIs(ctx, gm); err != nil {
		return err
	}

	return nil
}
