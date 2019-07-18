package internal

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/model"
)

// NewTrans 创建角色管理实例
func NewTrans(trans model.ITrans) *Trans {
	return &Trans{
		TransModel: trans,
	}
}

// Trans 事务管理
type Trans struct {
	TransModel model.ITrans
}

// Exec 执行事务
func (a *Trans) Exec(ctx context.Context, fn func(context.Context) error) error {
	return ExecTrans(ctx, a.TransModel, fn)
}
