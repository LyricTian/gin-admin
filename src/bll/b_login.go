package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/src/config"
	gcontext "github.com/LyricTian/gin-admin/src/context"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/model"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
)

// 定义错误
var (
	ErrInvalidUserName = errors.NewBadRequestError("无效的用户名")
	ErrInvalidPassword = errors.NewBadRequestError("无效的密码")
	ErrInvalidUser     = errors.NewUnauthorizedError("无效的用户")
	ErrUserDisable     = errors.NewUnauthorizedError("用户被禁用")
	ErrNoPerm          = errors.NewUnauthorizedError("没有权限")
)

// Login 登录管理
type Login struct {
	UserModel model.IUser `inject:"IUser"`
	RoleModel model.IRole `inject:"IRole"`
	MenuModel model.IMenu `inject:"IMenu"`
}

// GetRootUser 获取root用户数据
func (a *Login) GetRootUser() schema.User {
	rootUser := config.GetRootUser()
	return schema.User{
		RecordID: rootUser.UserName,
		UserName: rootUser.UserName,
		RealName: rootUser.RealName,
		Password: util.MD5HashString(rootUser.Password),
	}
}

// CheckIsRoot 检查是否是超级用户
func (a *Login) CheckIsRoot(ctx context.Context, recordID string) bool {
	rootUser := a.GetRootUser()
	if rootUser.RecordID == recordID {
		return true
	}
	return false
}

// Verify 登录验证
func (a *Login) Verify(ctx context.Context, userName, password string) (string, error) {
	// 检查是否是超级用户
	rootUser := a.GetRootUser()
	if userName == rootUser.UserName && rootUser.Password == password {
		return rootUser.RecordID, nil
	}

	user, err := a.UserModel.GetByUserName(ctx, userName, schema.UserQueryOptions{IncludePassword: true})
	if err != nil {
		return "", err
	} else if user == nil {
		return "", ErrInvalidUserName
	} else if user.Password != util.SHA1HashString(password) {
		return "", ErrInvalidPassword
	} else if user.Status != 1 {
		return "", ErrUserDisable
	}

	return user.RecordID, nil
}

// GetUserInfo 获取当前用户登录信息
func (a *Login) GetUserInfo(ctx context.Context) (*schema.UserLoginInfo, error) {
	userID := gcontext.FromUserID(ctx)
	if isRoot := a.CheckIsRoot(ctx, userID); isRoot {
		user := a.GetRootUser()
		loginInfo := &schema.UserLoginInfo{
			UserName: user.UserName,
			RealName: user.RealName,
		}
		return loginInfo, nil
	}

	user, err := a.UserModel.Get(ctx, userID, schema.UserQueryOptions{
		IncludeRoleIDs: true,
	})
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, ErrInvalidUser
	} else if user.Status != 1 {
		return nil, ErrUserDisable
	}

	loginInfo := &schema.UserLoginInfo{
		UserName: user.UserName,
		RealName: user.RealName,
	}

	if len(user.RoleIDs) > 0 {
		roles, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
			RecordIDs: user.RoleIDs,
		})
		if err != nil {
			return nil, err
		}
		loginInfo.RoleNames = roles.Data.ToNames()
	}

	return loginInfo, nil
}

// QueryUserMenuTree 查询当前用户的权限菜单树
func (a *Login) QueryUserMenuTree(ctx context.Context) ([]*schema.MenuTree, error) {
	userID := gcontext.FromUserID(ctx)
	if a.CheckIsRoot(ctx, userID) {
		userID = ""
	}

	// 查询用户的权限
	result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
		IsHide: 2,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	} else if len(result.Data) == 0 {
		return nil, ErrNoPerm
	}

	// 组装权限树
	perms := result.Data
	result, err = a.MenuModel.Query(ctx, schema.MenuQueryParam{
		RecordIDs: result.Data.SplitParentPathToRecordIDs(),
	})
	if err != nil {
		return nil, err
	}
	perms = append(result.Data, perms...)

	return perms.ToTreeList().ToTree(), nil
}

// UpdatePassword 更新当前用户登录密码
func (a *Login) UpdatePassword(ctx context.Context, params schema.UpdatePasswordParam) error {
	userID := gcontext.FromUserID(ctx)
	if a.CheckIsRoot(ctx, userID) {
		return errors.NewBadRequestError("超级用户密码只能通过配置文件修改")
	}

	user, err := a.UserModel.Get(ctx, userID, schema.UserQueryOptions{
		IncludePassword: true,
	})
	if err != nil {
		return err
	} else if user == nil {
		return ErrInvalidUser
	} else if user.Status != 1 {
		return ErrUserDisable
	} else if util.SHA1HashString(params.OldPassword) != user.Password {
		return errors.NewBadRequestError("旧密码不正确")
	}

	return a.UserModel.UpdatePassword(ctx, userID, params.NewPassword)
}
