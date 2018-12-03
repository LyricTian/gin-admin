package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/src/model"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// 定义错误
var (
	ErrInvalidUser     = errors.New("无效的用户")
	ErrInvalidUserName = errors.New("无效的用户名")
	ErrInvalidPassword = errors.New("无效的密码")
	ErrUserDisable     = errors.New("用户被禁用")
)

// Login 登录管理
type Login struct {
	UserModel model.IUser `inject:"IUser"`
	RoleModel model.IRole `inject:"IRole"`
	MenuModel model.IMenu `inject:"IMenu"`
}

func (a *Login) getRootUser() schema.User {
	rootUser := viper.GetStringSlice("system_root_user")
	if len(rootUser) == 2 {
		return schema.User{
			RecordID: rootUser[0],
			UserName: rootUser[0],
			RealName: "超级用户",
			Password: rootUser[1],
		}
	}
	return schema.User{}
}

// CheckIsRoot 检查是否是超级用户
func (a *Login) CheckIsRoot(ctx context.Context, recordID string) bool {
	rootUser := a.getRootUser()
	if rootUser.RecordID == recordID {
		return true
	}
	return false
}

// Verify 登录验证
func (a *Login) Verify(ctx context.Context, userName, password string) (*schema.User, error) {
	rootUser := a.getRootUser()
	if userName == rootUser.UserName &&
		util.MD5HashString(rootUser.Password) == password {
		return &rootUser, nil
	}

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

// GetCurrentUserInfo 获取当前用户信息
func (a *Login) GetCurrentUserInfo(ctx context.Context, userID string) (*schema.LoginInfo, error) {
	if isRoot := a.CheckIsRoot(ctx, userID); isRoot {
		rootUser := a.getRootUser()
		return &schema.LoginInfo{
			UserName: rootUser.UserName,
			RealName: rootUser.RealName,
		}, nil
	}

	user, err := a.UserModel.Get(ctx, userID, true)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, ErrInvalidUser
	} else if user.Status != 1 {
		return nil, ErrUserDisable
	}

	info := &schema.LoginInfo{
		UserName: user.UserName,
		RealName: user.RealName,
	}

	// 查询用户角色
	if len(user.RoleIDs) > 0 {
		roleItems, err := a.RoleModel.QuerySelect(ctx, schema.RoleSelectQueryParam{RecordIDs: user.RoleIDs})
		if err == nil && len(roleItems) > 0 {
			roleNames := make([]string, len(roleItems))
			for i, item := range roleItems {
				roleNames[i] = item.Name
			}
			info.RoleNames = roleNames
		}
	}

	return info, nil
}

// QueryCurrentUserMenus 查询当前用户菜单
func (a *Login) QueryCurrentUserMenus(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	params := schema.MenuSelectQueryParam{
		Status:     1,
		SystemCode: viper.GetString("system_code"),
		IsHide:     2,
		Types:      []int{20, 30},
	}

	if isRoot := a.CheckIsRoot(ctx, userID); !isRoot {
		params.UserID = userID
	}

	items, err := a.MenuModel.QuerySelect(ctx, params)
	if err != nil {
		return nil, err
	}

	treeData := util.Slice2Tree(util.StructsToMapSlice(items), "record_id", "parent_id")
	return treeData, nil
}
