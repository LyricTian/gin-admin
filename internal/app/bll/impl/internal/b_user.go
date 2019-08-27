package internal

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/errors"
	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/casbin/casbin"
)

// NewUser 创建菜单管理实例
func NewUser(
	e *casbin.Enforcer,
	mUser model.IUser,
	mRole model.IRole,
) *User {
	return &User{
		Enforcer:  e,
		UserModel: mUser,
		RoleModel: mRole,
	}
}

// User 用户管理
type User struct {
	Enforcer  *casbin.Enforcer
	UserModel model.IUser
	RoleModel model.IRole
}

// Query 查询数据
func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	return a.UserModel.Query(ctx, params, opts...)
}

// QueryShow 查询显示项数据
func (a *User) QueryShow(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserShowQueryResult, error) {
	userResult, err := a.UserModel.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	} else if userResult == nil {
		return nil, nil
	}

	result := &schema.UserShowQueryResult{
		PageResult: userResult.PageResult,
	}
	if len(userResult.Data) == 0 {
		return result, nil
	}

	roleResult, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
		RecordIDs: userResult.Data.ToRoleIDs(),
	})
	if err != nil {
		return nil, err
	}
	result.Data = userResult.Data.ToUserShows(roleResult.Data.ToMap())
	return result, nil
}

// Get 查询指定数据
func (a *User) Get(ctx context.Context, recordID string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	item, err := a.UserModel.Get(ctx, recordID, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}
	return item, nil
}

func (a *User) checkUserName(ctx context.Context, userName string) error {
	if userName == GetRootUser().UserName {
		return errors.ErrResourceExists
	}

	result, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		UserName: userName,
	}, schema.UserQueryOptions{
		PageParam: &schema.PaginationParam{PageSize: -1},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.ErrResourceExists
	}
	return nil
}

func (a *User) getUpdate(ctx context.Context, recordID string) (*schema.User, error) {
	nitem, err := a.Get(ctx, recordID, schema.UserQueryOptions{
		IncludeRoles: true,
	})
	if err != nil {
		return nil, err
	}

	err = a.LoadPolicy(ctx, *nitem)
	if err != nil {
		return nil, err
	}
	return nitem, nil
}

// Create 创建数据
func (a *User) Create(ctx context.Context, item schema.User) (*schema.User, error) {
	if item.Password == "" {
		return nil, errors.ErrUserNotEmptyPwd
	}

	err := a.checkUserName(ctx, item.UserName)
	if err != nil {
		return nil, err
	}

	item.Password = util.SHA1HashString(item.Password)
	item.RecordID = util.MustUUID()
	err = a.UserModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return a.getUpdate(ctx, item.RecordID)
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

	return a.getUpdate(ctx, recordID)
}

// Delete 删除数据
func (a *User) Delete(ctx context.Context, recordID string) error {
	oldItem, err := a.UserModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	err = a.UserModel.Delete(ctx, recordID)
	if err != nil {
		return err
	}
	a.Enforcer.DeleteUser(recordID)
	return nil
}

// UpdateStatus 更新状态
func (a *User) UpdateStatus(ctx context.Context, recordID string, status int) error {
	oldItem, err := a.UserModel.Get(ctx, recordID, schema.UserQueryOptions{
		IncludeRoles: true,
	})
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	err = a.UserModel.UpdateStatus(ctx, recordID, status)
	if err != nil {
		return err
	}

	if status == 2 {
		a.Enforcer.DeleteUser(recordID)
	} else {
		err = a.LoadPolicy(ctx, *oldItem)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadPolicy 加载用户权限策略
func (a *User) LoadPolicy(ctx context.Context, item schema.User) error {
	a.Enforcer.DeleteRolesForUser(item.RecordID)
	for _, roleID := range item.Roles.ToRoleIDs() {
		a.Enforcer.AddRoleForUser(item.RecordID, roleID)
	}
	return nil
}
