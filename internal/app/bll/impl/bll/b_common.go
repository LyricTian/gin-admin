package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/v6/internal/app/icontext"
	"github.com/LyricTian/gin-admin/v6/internal/app/model"
)

// TransFunc 定义事务执行函数
type TransFunc func(context.Context) error

// ExecTrans 执行事务
func ExecTrans(ctx context.Context, transModel model.ITrans, fn TransFunc) error {
	return transModel.Exec(ctx, fn)
}

// ExecTransWithLock 执行事务（加锁）
func ExecTransWithLock(ctx context.Context, transModel model.ITrans, fn TransFunc) error {
	if !icontext.FromTransLock(ctx) {
		ctx = icontext.NewTransLock(ctx)
	}
	return ExecTrans(ctx, transModel, fn)
}

// NewNoTrans 不使用事务执行
func NewNoTrans(ctx context.Context) context.Context {
	if !icontext.FromNoTrans(ctx) {
		return icontext.NewNoTrans(ctx)
	}
	return ctx
}
