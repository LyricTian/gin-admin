package internal

import (
	"context"
	"net/http"

	"github.com/LyricTian/captcha"
	"github.com/LyricTian/gin-admin/internal/app/errors"
	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/auth"
	"github.com/LyricTian/gin-admin/pkg/util"
)

// NewLogin 创建登录管理实例
func NewLogin(
	a auth.Auther,
	mUser model.IUser,
	mRole model.IRole,
	mMenu model.IMenu,
) *Login {
	return &Login{
		Auth:      a,
		UserModel: mUser,
		RoleModel: mRole,
		MenuModel: mMenu,
	}
}

// Login 登录管理
type Login struct {
	UserModel model.IUser
	RoleModel model.IRole
	MenuModel model.IMenu
	Auth      auth.Auther
}

// GetCaptcha 获取图形验证码信息
func (a *Login) GetCaptcha(ctx context.Context, length int) (*schema.LoginCaptcha, error) {
	captchaID := captcha.NewLen(length)
	item := &schema.LoginCaptcha{
		CaptchaID: captchaID,
	}
	return item, nil
}

// ResCaptcha 生成并响应图形验证码
func (a *Login) ResCaptcha(ctx context.Context, w http.ResponseWriter, captchaID string, width, height int) error {
	err := captcha.WriteImage(w, captchaID, width, height)
	if err != nil {
		if err == captcha.ErrNotFound {
			return errors.ErrNotFound
		}
		return errors.WithStack(err)
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")
	return nil
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
		return nil, errors.ErrInvalidUserName
	}

	item := result.Data[0]
	if item.Password != util.SHA1HashString(password) {
		return nil, errors.ErrInvalidPassword
	} else if item.Status != 1 {
		return nil, errors.ErrUserDisable
	}

	return item, nil
}

// GenerateToken 生成令牌
func (a *Login) GenerateToken(ctx context.Context, userID string) (*schema.LoginTokenInfo, error) {
	tokenInfo, err := a.Auth.GenerateToken(userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	item := &schema.LoginTokenInfo{
		AccessToken: tokenInfo.GetAccessToken(),
		TokenType:   tokenInfo.GetTokenType(),
		ExpiresAt:   tokenInfo.GetExpiresAt(),
	}
	return item, nil
}

// DestroyToken 销毁令牌
func (a *Login) DestroyToken(ctx context.Context, tokenString string) error {
	err := a.Auth.DestroyToken(tokenString)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *Login) getAndCheckUser(ctx context.Context, userID string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	user, err := a.UserModel.Get(ctx, userID, opts...)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.ErrInvalidUser
	} else if user.Status != 1 {
		return nil, errors.ErrUserDisable
	}
	return user, nil
}

// GetLoginInfo 获取当前用户登录信息
func (a *Login) GetLoginInfo(ctx context.Context, userID string) (*schema.UserLoginInfo, error) {
	if isRoot := CheckIsRootUser(ctx, userID); isRoot {
		root := GetRootUser()
		loginInfo := &schema.UserLoginInfo{
			UserName: root.UserName,
			RealName: root.RealName,
		}
		return loginInfo, nil
	}

	user, err := a.getAndCheckUser(ctx, userID, schema.UserQueryOptions{
		IncludeRoles: true,
	})
	if err != nil {
		return nil, err
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
func (a *Login) QueryUserMenuTree(ctx context.Context, userID string) ([]*schema.MenuTree, error) {
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
		return nil, errors.ErrNoPerm
	}

	// 查询角色权限菜单列表
	menuResult, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
		RecordIDs: roleResult.Data.ToMenuIDs(),
	})
	if err != nil {
		return nil, err
	} else if len(menuResult.Data) == 0 {
		return nil, errors.ErrNoPerm
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
		return nil, errors.ErrNoPerm
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
func (a *Login) UpdatePassword(ctx context.Context, userID string, params schema.UpdatePasswordParam) error {
	if CheckIsRootUser(ctx, userID) {
		return errors.ErrLoginNotAllowModifyPwd
	}

	user, err := a.getAndCheckUser(ctx, userID)
	if err != nil {
		return err
	} else if util.SHA1HashString(params.OldPassword) != user.Password {
		return errors.ErrLoginInvalidOldPwd
	}

	params.NewPassword = util.SHA1HashString(params.NewPassword)
	return a.UserModel.UpdatePassword(ctx, userID, params.NewPassword)
}
