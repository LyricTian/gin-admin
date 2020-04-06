package model

import (
	"context"

	icontext "github.com/LyricTian/gin-admin/internal/app/context"
	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ model.ITrans = new(Trans)

// TransSet 注入Trans
var TransSet = wire.NewSet(wire.Struct(new(Trans), "*"), wire.Bind(new(model.ITrans), new(*Trans)))

// Trans 事务管理
type Trans struct {
	Client *mongo.Client
}

// Exec 执行事务
func (a *Trans) Exec(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := icontext.FromTrans(ctx); ok {
		return fn(ctx)
	}

	session, err := a.Client.StartSession()
	if err != nil {
		return errors.WithStack(err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		err := fn(icontext.NewTrans(sessCtx, true))
		return nil, err
	})

	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
