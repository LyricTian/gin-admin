package bll

import (
	"context"
	"gin-admin/src/model"
	"gin-admin/src/schema"
	"gin-admin/src/util"

	"github.com/pkg/errors"
)

// 定义错误
var (
	ErrInvalidUserName = errors.New("无效的用户名")
	ErrInvalidPassword = errors.New("无效的密码")
	ErrUserDisable     = errors.New("用户被禁用")
)

// Login 登录管理
type Login struct {
	UserModel model.IUser `inject:"IUser"`
	MenuModel model.IMenu `inject:"IMenu"`
}

// Verify 登录验证
func (a *Login) Verify(ctx context.Context, userName, password string) (*schema.User, error) {
	user, err := a.UserModel.GetByUserName(ctx, userName, false)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, ErrInvalidUserName
	} else if user.Status != 1 {
		return nil, ErrUserDisable
	} else if user.Password != util.SHA1HashString(password) {
		return nil, ErrInvalidPassword
	}

	return user, nil
}
