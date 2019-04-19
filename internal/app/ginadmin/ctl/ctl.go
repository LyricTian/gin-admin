package ctl

import (
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/bll"
	"github.com/LyricTian/gin-admin/pkg/auth"
)

// Common 提供统一的控制器管理
type Common struct {
	Demo  *Demo
	Login *Login
	Menu  *Menu
	Role  *Role
	User  *User
}

// NewCommon 创建统一的控制器管理
func NewCommon(b *bll.Common, a auth.Auther) *Common {
	return &Common{
		Demo:  NewDemo(b),
		Login: NewLogin(b, a),
		Menu:  NewMenu(b),
		Role:  NewRole(b),
		User:  NewUser(b),
	}
}
