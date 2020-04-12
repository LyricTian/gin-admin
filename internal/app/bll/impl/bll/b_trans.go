package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/bll"
	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/google/wire"
)

var _ bll.ITrans = new(Trans)

// TransSet 注入Trans
var TransSet = wire.NewSet(wire.Struct(new(Trans), "*"), wire.Bind(new(bll.ITrans), new(*Trans)))

// Trans 事务管理
type Trans struct {
	TransModel model.ITrans
}

// Exec 执行事务
func (a *Trans) Exec(ctx context.Context, fn func(context.Context) error) error {
	return ExecTrans(ctx, a.TransModel, fn)
}
