package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/config"
	icontext "github.com/LyricTian/gin-admin/internal/app/context"
	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
)

// GetRootUser 获取root用户
func GetRootUser() *schema.User {
	user := config.C.Root
	return &schema.User{
		RecordID: user.UserName,
		UserName: user.UserName,
		RealName: user.RealName,
		Password: util.MD5HashString(user.Password),
	}
}

// CheckIsRootUser 检查是否是root用户
func CheckIsRootUser(ctx context.Context, userID string) bool {
	return GetRootUser().RecordID == userID
}

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
