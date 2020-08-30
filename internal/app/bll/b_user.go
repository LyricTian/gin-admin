package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/v7/internal/app/model/gormx/model"
	"github.com/LyricTian/gin-admin/v7/internal/app/schema"
	"github.com/LyricTian/gin-admin/v7/pkg/errors"
	"github.com/LyricTian/gin-admin/v7/pkg/util/hash"
	"github.com/LyricTian/gin-admin/v7/pkg/util/uuid"
	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
)

// UserSet 注入User
var UserSet = wire.NewSet(wire.Struct(new(User), "*"))

// User 用户管理
type User struct {
	Enforcer      *casbin.SyncedEnforcer
	TransModel    *model.Trans
	UserModel     *model.User
	UserRoleModel *model.UserRole
	RoleModel     *model.Role
}

// Query 查询数据
func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	return a.UserModel.Query(ctx, params, opts...)
}

// QueryShow 查询显示项数据
func (a *User) QueryShow(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserShowQueryResult, error) {
	result, err := a.UserModel.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	} else if result == nil {
		return nil, nil
	}

	userRoleResult, err := a.UserRoleModel.Query(ctx, schema.UserRoleQueryParam{
		UserIDs: result.Data.ToIDs(),
	})
	if err != nil {
		return nil, err
	}

	roleResult, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
		IDs: userRoleResult.Data.ToRoleIDs(),
	})
	if err != nil {
		return nil, err
	}

	return result.ToShowResult(userRoleResult.Data.ToUserIDMap(), roleResult.Data.ToMap()), nil
}

// Get 查询指定数据
func (a *User) Get(ctx context.Context, id string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	item, err := a.UserModel.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	userRoleResult, err := a.UserRoleModel.Query(ctx, schema.UserRoleQueryParam{
		UserID: id,
	})
	if err != nil {
		return nil, err
	}
	item.UserRoles = userRoleResult.Data

	return item, nil
}

// Create 创建数据
func (a *User) Create(ctx context.Context, item schema.User) (*schema.IDResult, error) {
	err := a.checkUserName(ctx, item)
	if err != nil {
		return nil, err
	}

	item.Password = hash.SHA1String(item.Password)
	item.ID = uuid.MustString()
	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		for _, urItem := range item.UserRoles {
			urItem.ID = uuid.MustString()
			urItem.UserID = item.ID
			err := a.UserRoleModel.Create(ctx, *urItem)
			if err != nil {
				return err
			}
		}

		return a.UserModel.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return schema.NewIDResult(item.ID), nil
}

func (a *User) checkUserName(ctx context.Context, item schema.User) error {
	if item.UserName == schema.GetRootUser().UserName {
		return errors.New400Response("用户名不合法")
	}

	result, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		UserName:        item.UserName,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("用户名已经存在")
	}
	return nil
}

// Update 更新数据
func (a *User) Update(ctx context.Context, id string, item schema.User) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.UserName != item.UserName {
		err := a.checkUserName(ctx, item)
		if err != nil {
			return err
		}
	}

	if item.Password != "" {
		item.Password = hash.SHA1String(item.Password)
	} else {
		item.Password = oldItem.Password
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt
	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		addUserRoles, delUserRoles := a.compareUserRoles(ctx, oldItem.UserRoles, item.UserRoles)
		for _, rmitem := range addUserRoles {
			rmitem.ID = uuid.MustString()
			rmitem.UserID = id
			err := a.UserRoleModel.Create(ctx, *rmitem)
			if err != nil {
				return err
			}
		}

		for _, rmitem := range delUserRoles {
			err := a.UserRoleModel.Delete(ctx, rmitem.ID)
			if err != nil {
				return err
			}
		}

		return a.UserModel.Update(ctx, id, item)
	})
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}

func (a *User) compareUserRoles(ctx context.Context, oldUserRoles, newUserRoles schema.UserRoles) (addList, delList schema.UserRoles) {
	mOldUserRoles := oldUserRoles.ToMap()
	mNewUserRoles := newUserRoles.ToMap()

	for k, item := range mNewUserRoles {
		if _, ok := mOldUserRoles[k]; ok {
			delete(mOldUserRoles, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldUserRoles {
		delList = append(delList, item)
	}
	return
}

// Delete 删除数据
func (a *User) Delete(ctx context.Context, id string) error {
	oldItem, err := a.UserModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	err = a.TransModel.Exec(ctx, func(ctx context.Context) error {
		err := a.UserRoleModel.DeleteByUserID(ctx, id)
		if err != nil {
			return err
		}

		return a.UserModel.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}

// UpdateStatus 更新状态
func (a *User) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.UserModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}
	oldItem.Status = status

	err = a.UserModel.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}
