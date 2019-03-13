package bll

import (
	"context"

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
	UserBll   *User       `inject:""`
	UserModel model.IUser `inject:"IUser"`
	RoleModel model.IRole `inject:"IRole"`
	MenuModel model.IMenu `inject:"IMenu"`
}

// Verify 登录验证
func (a *Login) Verify(ctx context.Context, userName, password string) (*schema.User, error) {
	// 检查是否是超级用户
	root := a.UserBll.GetRoot()
	if userName == root.UserName && root.Password == password {
		return root, nil
	}

	result, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		UserName: userName,
	})
	if err != nil {
		return nil, err
	} else if len(result.Data) == 0 {
		return nil, ErrInvalidUserName
	}

	item := result.Data[0]
	if item.Password != util.SHA1HashString(password) {
		return nil, ErrInvalidPassword
	} else if item.Status != 1 {
		return nil, ErrUserDisable
	}

	return item, nil
}

// GetUserInfo 获取当前用户登录信息
func (a *Login) GetUserInfo(ctx context.Context) (*schema.UserLoginInfo, error) {
	userID := gcontext.FromUserID(ctx)
	if isRoot := a.UserBll.CheckIsRoot(ctx, userID); isRoot {
		root := a.UserBll.GetRoot()
		loginInfo := &schema.UserLoginInfo{
			UserName: root.UserName,
			RealName: root.RealName,
		}
		return loginInfo, nil
	}

	user, err := a.UserModel.Get(ctx, userID, schema.UserQueryOptions{
		IncludeRoles: true,
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

	if roleIDs := user.Roles.ToRoleIDs(); len(roleIDs) > 0 {
		roles, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
			RecordIDs: roleIDs,
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
	isRoot := a.UserBll.CheckIsRoot(ctx, userID)

	// 如果是root用户，则查询所有显示的菜单树，并设定所有的菜单操作权限为读写权限
	if isRoot {
		hidden := 0
		result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
			Hidden: &hidden,
		})
		if err != nil {
			return nil, err
		}
		return result.Data.ToTrees().ForEach(func(item *schema.MenuTree, _ int) {
			item.OperPerm = "rw"
		}).ToTree(), nil
	}

	roleResult, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
		UserID: userID,
	}, schema.RoleQueryOptions{
		IncludeMenus: true,
	})
	if err != nil {
		return nil, err
	} else if len(roleResult.Data) == 0 {
		return nil, ErrNoPerm
	}

	// 查询角色权限菜单列表
	menuResult, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
		RecordIDs: roleResult.Data.ToMenuIDs(),
	})
	if err != nil {
		return nil, err
	} else if len(menuResult.Data) == 0 {
		return nil, ErrNoPerm
	}

	// 拆分并查询菜单树
	menuResult, err = a.MenuModel.Query(ctx, schema.MenuQueryParam{
		RecordIDs: menuResult.Data.SplitAndGetAllRecordIDs(),
	})
	if err != nil {
		return nil, err
	} else if len(menuResult.Data) == 0 {
		return nil, ErrNoPerm
	}

	pm := roleResult.Data.ToMenuIDPermMap()
	return menuResult.Data.ToTrees().ForEach(func(item *schema.MenuTree, _ int) {
		item.OperPerm = pm[item.RecordID]
	}).ToTree(), nil
}

// UpdatePassword 更新当前用户登录密码
func (a *Login) UpdatePassword(ctx context.Context, params schema.UpdatePasswordParam) error {
	userID := gcontext.FromUserID(ctx)
	if a.UserBll.CheckIsRoot(ctx, userID) {
		return errors.NewBadRequestError("超级管理员密码只能通过配置文件修改")
	}

	user, err := a.UserModel.Get(ctx, userID)
	if err != nil {
		return err
	} else if user == nil {
		return ErrInvalidUser
	} else if user.Status != 1 {
		return ErrUserDisable
	} else if util.SHA1HashString(params.OldPassword) != user.Password {
		return errors.NewBadRequestError("旧密码不正确")
	}

	params.NewPassword = util.SHA1HashString(params.NewPassword)
	return a.UserModel.UpdatePassword(ctx, userID, params.NewPassword)
}
