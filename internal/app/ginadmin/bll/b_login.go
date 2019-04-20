package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/ginadmin/model"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/LyricTian/gin-admin/pkg/util"
)

// 定义错误
var (
	ErrInvalidUserName = errors.NewBadRequestError("无效的用户名")
	ErrInvalidPassword = errors.NewBadRequestError("无效的密码")
	ErrInvalidUser     = errors.NewUnauthorizedError("无效的用户")
	ErrUserDisable     = errors.NewUnauthorizedError("用户被禁用")
	ErrNoPerm          = errors.NewUnauthorizedError("没有权限")
)

// NewLogin 创建登录管理实例
func NewLogin(m *model.Common) *Login {
	return &Login{
		UserModel: m.User,
		RoleModel: m.Role,
		MenuModel: m.Menu,
	}
}

// Login 登录管理
type Login struct {
	UserModel model.IUser
	RoleModel model.IRole
	MenuModel model.IMenu
}

// Verify 登录验证
func (a *Login) Verify(ctx context.Context, userName, password string) (*schema.User, error) {
	// 检查是否是超级用户
	root := GetRootUser()
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
	userID := GetUserID(ctx)
	if isRoot := CheckIsRootUser(ctx, userID); isRoot {
		root := GetRootUser()
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
	userID := GetUserID(ctx)
	isRoot := CheckIsRootUser(ctx, userID)

	// 如果是root用户，则查询所有显示的菜单树
	if isRoot {
		hidden := 0
		result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
			Hidden: &hidden,
		}, schema.MenuQueryOptions{
			IncludeActions: true,
		})
		if err != nil {
			return nil, err
		}
		return result.Data.ToTrees().ToTree(), nil
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
	}, schema.MenuQueryOptions{
		IncludeActions: true,
	})
	if err != nil {
		return nil, err
	} else if len(menuResult.Data) == 0 {
		return nil, ErrNoPerm
	}

	menuActions := roleResult.Data.ToMenuIDActionsMap()
	return menuResult.Data.ToTrees().ForEach(func(item *schema.MenuTree, _ int) {
		// 遍历菜单动作权限
		var actions []*schema.MenuAction
		for _, code := range menuActions[item.RecordID] {
			for _, aitem := range item.Actions {
				if aitem.Code == code {
					actions = append(actions, aitem)
					break
				}
			}
		}
		item.Actions = actions
	}).ToTree(), nil
}

// UpdatePassword 更新当前用户登录密码
func (a *Login) UpdatePassword(ctx context.Context, params schema.UpdatePasswordParam) error {
	userID := GetUserID(ctx)
	if CheckIsRootUser(ctx, userID) {
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
