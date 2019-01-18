package bll

import (
	"context"

	gcontext "github.com/LyricTian/gin-admin/src/context"
	"github.com/LyricTian/gin-admin/src/model"
)

// TransFunc 定义事务执行函数
type TransFunc func(context.Context) error

// Common 公共处理
type Common struct {
	TransModel model.ITrans `inject:"ITrans"`
}

// ExecTrans 执行事务
func (a *Common) ExecTrans(ctx context.Context, fn TransFunc) error {
	trans, err := a.TransModel.Begin(ctx)
	if err != nil {
		return err
	}

	err = fn(gcontext.NewTrans(ctx, trans))
	if err != nil {
		_ = a.TransModel.Rollback(ctx, trans)
		return err
	}
	return a.TransModel.Commit(ctx, trans)
}
