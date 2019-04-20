package bll

import (
	"context"
	"sync"

	"github.com/LyricTian/gin-admin/internal/app/ginadmin/config"
	icontext "github.com/LyricTian/gin-admin/internal/app/ginadmin/context"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/model"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
	"github.com/LyricTian/gin-admin/pkg/auth"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/casbin/casbin"
)

// GetUserID 获取用户ID
func GetUserID(ctx context.Context) string {
	userID, ok := icontext.FromUserID(ctx)
	if ok {
		return userID
	}
	return ""
}

// TransFunc 定义事务执行函数
type TransFunc func(context.Context) error

// ExecTrans 执行事务
func ExecTrans(ctx context.Context, transModel model.ITrans, fn TransFunc) error {
	if _, ok := icontext.FromTrans(ctx); ok {
		return fn(ctx)
	}
	trans, err := transModel.Begin(ctx)
	if err != nil {
		return err
	}

	err = fn(icontext.NewTrans(ctx, trans))
	if err != nil {
		_ = transModel.Rollback(ctx, trans)
		return err
	}
	return transModel.Commit(ctx, trans)
}

var (
	rootUser     *schema.User
	rootUserOnce sync.Once
)

// GetRootUser 获取root用户
func GetRootUser() *schema.User {
	rootUserOnce.Do(func() {
		user := config.GetGlobalConfig().Root
		rootUser = &schema.User{
			RecordID: user.UserName,
			UserName: user.UserName,
			RealName: user.RealName,
			Password: util.MD5HashString(user.Password),
		}
	})
	return rootUser
}

// CheckIsRootUser 检查是否是root用户
func CheckIsRootUser(ctx context.Context, recordID string) bool {
	return GetRootUser().RecordID == recordID
}

// Common 提供统一的业务逻辑处理
type Common struct {
	Demo  *Demo
	Login *Login
	Menu  *Menu
	Role  *Role
	User  *User
}

// NewCommon 创建统一的业务逻辑处理
func NewCommon(m *model.Common, a auth.Auther, e *casbin.Enforcer) *Common {
	return &Common{
		Demo:  NewDemo(m),
		Login: NewLogin(m, a),
		Menu:  NewMenu(m),
		Role:  NewRole(m, e),
		User:  NewUser(m, e),
	}
}
