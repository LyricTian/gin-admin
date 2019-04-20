package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/ginadmin/model"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/casbin/casbin"
)

// NewUser 创建菜单管理实例
func NewUser(m *model.Common, e *casbin.Enforcer) *User {
	return &User{
		UserModel: m.User,
		RoleModel: m.Role,
		Enforcer:  e,
	}
}

// User 用户管理
type User struct {
	UserModel model.IUser
	RoleModel model.IRole
	Enforcer  *casbin.Enforcer
}

// QueryPage 查询分页数据
func (a *User) QueryPage(ctx context.Context, params schema.UserQueryParam, pp *schema.PaginationParam) ([]*schema.UserPageShow, *schema.PaginationResult, error) {
	result, err := a.UserModel.Query(ctx, params, schema.UserQueryOptions{
		PageParam:    pp,
		IncludeRoles: true,
	})
	if err != nil {
		return nil, nil, err
	} else if len(result.Data) == 0 {
		return nil, result.PageResult, nil
	}

	// 填充角色数据
	roleResult, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
		RecordIDs: result.Data.ToRoleIDs(),
	})
	if err != nil {
		return nil, nil, err
	}

	pageResult := result.Data.ToPageShows(roleResult.Data.ToMap())
	return pageResult, result.PageResult, nil
}

// Get 查询指定数据
func (a *User) Get(ctx context.Context, recordID string) (*schema.User, error) {
	item, err := a.UserModel.Get(ctx, recordID, schema.UserQueryOptions{
		IncludeRoles: true,
	})
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}
	// 不输出用户密码
	item.Password = ""

	return item, nil
}

func (a *User) checkUserName(ctx context.Context, userName string) error {
	if userName == GetRootUser().UserName {
		return errors.NewBadRequestError("用户名不合法")
	}

	result, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		UserName: userName,
	}, schema.UserQueryOptions{
		PageParam: &schema.PaginationParam{PageSize: -1},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.NewBadRequestError("用户名已经存在")
	}
	return nil
}

// Create 创建数据
func (a *User) Create(ctx context.Context, item schema.User) (*schema.User, error) {
	if item.Password == "" {
		return nil, errors.NewBadRequestError("用户密码不允许为空")
	}

	err := a.checkUserName(ctx, item.UserName)
	if err != nil {
		return nil, err
	}

	item.Password = util.SHA1HashString(item.Password)
	item.RecordID = util.MustUUID()
	item.Creator = GetUserID(ctx)
	err = a.UserModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	nitem, err := a.Get(ctx, item.RecordID)
	if err != nil {
		return nil, err
	}

	err = a.LoadPolicy(ctx, nitem)
	if err != nil {
		return nil, err
	}
	return nitem, nil
}

// Update 更新数据
func (a *User) Update(ctx context.Context, recordID string, item schema.User) (*schema.User, error) {
	oldItem, err := a.UserModel.Get(ctx, recordID)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	} else if oldItem.UserName != item.UserName {
		err := a.checkUserName(ctx, item.UserName)
		if err != nil {
			return nil, err
		}
	}

	if item.Password != "" {
		item.Password = util.SHA1HashString(item.Password)
	}

	err = a.UserModel.Update(ctx, recordID, item)
	if err != nil {
		return nil, err
	}

	nitem, err := a.Get(ctx, item.RecordID)
	if err != nil {
		return nil, err
	}

	err = a.LoadPolicy(ctx, nitem)
	if err != nil {
		return nil, err
	}
	return nitem, nil
}

// Delete 删除数据
func (a *User) Delete(ctx context.Context, recordID string) error {
	err := a.UserModel.Delete(ctx, recordID)
	if err != nil {
		return err
	}
	a.Enforcer.DeleteUser(recordID)
	return nil
}

// UpdateStatus 更新状态
func (a *User) UpdateStatus(ctx context.Context, recordID string, status int) error {
	err := a.UserModel.UpdateStatus(ctx, recordID, status)
	if err != nil {
		return err
	}

	if status == 2 {
		a.Enforcer.DeleteUser(recordID)
	} else {
		err = a.LoadPolicyWithRecordID(ctx, recordID)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadAllPolicy 加载所有的用户策略
func (a *User) LoadAllPolicy(ctx context.Context) error {
	result, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		Status: 1,
	}, schema.UserQueryOptions{IncludeRoles: true})
	if err != nil {
		return err
	}

	for _, item := range result.Data {
		err := a.LoadPolicy(ctx, item)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadPolicyWithRecordID 加载用户权限策略
func (a *User) LoadPolicyWithRecordID(ctx context.Context, recordID string) error {
	item, err := a.UserModel.Get(ctx, recordID, schema.UserQueryOptions{
		IncludeRoles: true,
	})
	if err != nil {
		return err
	}
	return a.LoadPolicy(ctx, item)
}

// LoadPolicy 加载用户权限策略
func (a *User) LoadPolicy(ctx context.Context, item *schema.User) error {
	a.Enforcer.DeleteRolesForUser(item.RecordID)
	for _, roleID := range item.Roles.ToRoleIDs() {
		a.Enforcer.AddRoleForUser(item.RecordID, roleID)
	}
	return nil
}
