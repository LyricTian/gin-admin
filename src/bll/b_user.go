package bll

import (
	"context"

	"github.com/LyricTian/gin-admin/src/config"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/model"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/casbin/casbin"
)

// User 用户管理
type User struct {
	UserModel model.IUser      `inject:"IUser"`
	RoleModel model.IRole      `inject:"IRole"`
	Enforcer  *casbin.Enforcer `inject:""`
	CommonBll *Common          `inject:""`
}

// GetRoot 获取root用户数据
func (a *User) GetRoot() schema.User {
	user := config.GetRoot()
	return schema.User{
		RecordID: user.UserName,
		UserName: user.UserName,
		RealName: user.RealName,
		Password: util.MD5HashString(user.Password),
	}
}

// CheckIsRoot 检查是否是root
func (a *User) CheckIsRoot(ctx context.Context, recordID string) bool {
	return a.GetRoot().RecordID == recordID
}

// QueryPage 查询分页数据
func (a *User) QueryPage(ctx context.Context, params schema.UserQueryParam, pp *schema.PaginationParam) ([]*schema.UserPageShow, *schema.PaginationResult, error) {
	result, err := a.UserModel.Query(ctx, params, schema.UserQueryOptions{
		PageParam:      pp,
		IncludeRoleIDs: true,
	})
	if err != nil {
		return nil, nil, err
	} else if len(result.Data) == 0 {
		return nil, nil, nil
	}

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
		IncludeRoleIDs: true,
	})
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}
	item.Password = ""

	return item, nil
}

func (a *User) check(ctx context.Context, item schema.User, oldItem *schema.User) error {
	if oldItem == nil || oldItem.UserName != item.UserName {
		exists, err := a.UserModel.CheckUserName(ctx, item.UserName)
		if err != nil {
			return err
		} else if exists {
			return errors.NewBadRequestError("用户名已经存在，请重新填写")
		}
	}

	roleResult, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{RecordIDs: item.RoleIDs})
	if err != nil {
		return err
	} else if len(roleResult.Data) == 0 {
		return errors.NewBadRequestError("请选择有效的角色")
	}
	return nil
}

// Create 创建数据
func (a *User) Create(ctx context.Context, item schema.User) (*schema.User, error) {
	err := a.check(ctx, item, nil)
	if err != nil {
		return nil, err
	}

	item.Password = util.SHA1HashString(item.Password)
	item.RecordID = util.MustUUID()
	item.Creator = a.CommonBll.GetUserID(ctx)
	err = a.UserModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	err = a.LoadPolicy(ctx, item.RecordID)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Update 更新数据
func (a *User) Update(ctx context.Context, recordID string, item schema.User) error {
	oldItem, err := a.UserModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	err = a.check(ctx, item, oldItem)
	if err != nil {
		return err
	}

	if item.Password != "" {
		item.Password = util.SHA1HashString(item.Password)
	}

	err = a.UserModel.Update(ctx, recordID, item)
	if err != nil {
		return err
	}

	return a.LoadPolicy(ctx, recordID)
}

// Delete 删除数据
func (a *User) Delete(ctx context.Context, recordIDs ...string) error {
	err := a.CommonBll.ExecTrans(ctx, func(ctx context.Context) error {
		for _, recordID := range recordIDs {
			err := a.UserModel.Delete(ctx, recordID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, recordID := range recordIDs {
		a.Enforcer.DeleteUser(recordID)
	}

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
		err = a.LoadPolicy(ctx, recordID)
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
	}, schema.UserQueryOptions{IncludeRoleIDs: true})
	if err != nil {
		return err
	}

	for _, item := range result.Data {
		for _, roleID := range item.RoleIDs {
			a.Enforcer.AddRoleForUser(item.RecordID, roleID)
		}
	}

	return nil
}

// LoadPolicy 加载用户权限策略
func (a *User) LoadPolicy(ctx context.Context, recordID string) error {
	item, err := a.UserModel.Get(ctx, recordID, schema.UserQueryOptions{
		IncludeRoleIDs: true,
	})
	if err != nil {
		return err
	}

	a.Enforcer.DeleteRolesForUser(recordID)
	for _, roleID := range item.RoleIDs {
		a.Enforcer.AddRoleForUser(recordID, roleID)
	}
	return nil
}
