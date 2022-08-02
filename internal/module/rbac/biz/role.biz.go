package biz

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/dao"
	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/contextx"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"github.com/LyricTian/gin-admin/v9/pkg/logger"
	"github.com/LyricTian/gin-admin/v9/pkg/util/xid"
	"github.com/LyricTian/gin-admin/v9/pkg/x/cachex"
	"go.uber.org/zap"
)

type RoleBiz struct {
	TransRepo    utilx.TransRepo
	RoleRepo     dao.RoleRepo
	RoleMenuRepo dao.RoleMenuRepo
	UserRoleRepo dao.UserRoleRepo
	Cache        cachex.Cacher
}

func (a *RoleBiz) Query(ctx context.Context, params typed.RoleQueryParam) (*typed.RoleQueryResult, error) {
	params.Pagination = true
	queryOpts := utilx.QueryOptions{
		OrderFields: []utilx.OrderByParam{
			{Field: "sequence", Direction: utilx.DESC},
			{Field: "updated_at", Direction: utilx.DESC},
		},
	}

	if params.Result == "select" {
		params.Pagination = false
		queryOpts.SelectFields = []string{"id", "name"}
	}

	result, err := a.RoleRepo.Query(ctx, params, typed.RoleQueryOptions{
		QueryOptions: queryOpts,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *RoleBiz) Get(ctx context.Context, id string) (*typed.Role, error) {
	role, err := a.RoleRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if role == nil {
		return nil, errors.NotFound(errors.ErrNotFoundID, "Role not found")
	}

	roleMenuResult, err := a.RoleMenuRepo.Query(ctx, typed.RoleMenuQueryParam{
		RoleID: id,
	})
	if err != nil {
		return nil, err
	}
	role.RoleMenus = roleMenuResult.Data

	return role, nil
}

func (a *RoleBiz) Create(ctx context.Context, createItem typed.RoleCreate) (*typed.Role, error) {
	role := &typed.Role{
		ID:        xid.NewID(),
		Name:      createItem.Name,
		Sequence:  createItem.Sequence,
		Remark:    &createItem.Remark,
		Status:    typed.RoleStatusEnabled,
		CreatedBy: contextx.FromUserID(ctx),
	}

	err := a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.RoleRepo.Create(ctx, role); err != nil {
			return err
		}

		for _, roleMenu := range createItem.RoleMenus {
			err := a.RoleMenuRepo.Create(ctx, &typed.RoleMenu{
				ID:           xid.NewID(),
				RoleID:       role.ID,
				MenuID:       roleMenu.MenuID,
				MenuActionID: roleMenu.MenuActionID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (a *RoleBiz) Update(ctx context.Context, id string, createItem typed.RoleCreate) error {
	oldRole, err := a.RoleRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldRole == nil {
		return errors.NotFound(errors.ErrNotFoundID, "Role not found")
	}
	oldRole.Name = createItem.Name
	oldRole.Sequence = createItem.Sequence
	oldRole.Remark = &createItem.Remark
	oldRole.UpdatedBy = contextx.FromUserID(ctx)

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.RoleRepo.Update(ctx, oldRole); err != nil {
			return err
		}

		if err := a.RoleMenuRepo.DeleteByRoleID(ctx, id); err != nil {
			return err
		}

		for _, roleMenu := range createItem.RoleMenus {
			err := a.RoleMenuRepo.Create(ctx, &typed.RoleMenu{
				ID:           xid.NewID(),
				RoleID:       oldRole.ID,
				MenuID:       roleMenu.MenuID,
				MenuActionID: roleMenu.MenuActionID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (a *RoleBiz) clearUserCacheByID(ctx context.Context, id string) error {
	userRoleResult, err := a.UserRoleRepo.Query(ctx, typed.UserRoleQueryParam{
		RoleID: id,
	}, typed.UserRoleQueryOptions{
		QueryOptions: utilx.QueryOptions{
			SelectFields: []string{"user_id"},
		},
	})
	if err != nil {
		return err
	}

	for _, userRole := range userRoleResult.Data {
		err = a.Cache.Delete(ctx, utilx.CacheNSForUserRole, userRole.UserID)
		if err != nil {
			logger.Context(ctx).Error("Failed to delete user role cache", zap.String("user_id", userRole.UserID), zap.Error(err))
		}
	}

	return nil
}

func (a *RoleBiz) Delete(ctx context.Context, id string) error {
	exists, err := a.RoleRepo.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound(errors.ErrNotFoundID, "Role not found")
	}

	err = a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		if err := a.RoleRepo.Delete(ctx, id); err != nil {
			return err
		}

		if err := a.RoleMenuRepo.DeleteByRoleID(ctx, id); err != nil {
			return err
		}

		if err := a.UserRoleRepo.DeleteByRoleID(ctx, id); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Set deleted role flag to cache for casbin
	err = a.Cache.Set(ctx, utilx.CacheNSForDeletedRole, utilx.CacheKeyForDeletedRole, "1")
	if err != nil {
		logger.Context(ctx).Error("Failed to set deleted role cache", zap.Error(err))
	}

	// Clear user cache by role id
	return a.clearUserCacheByID(ctx, id)
}

func (a *RoleBiz) UpdateStatus(ctx context.Context, id string, status typed.RoleStatus) error {
	exists, err := a.RoleRepo.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound(errors.ErrNotFoundID, "Role not found")
	}

	err = a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.RoleRepo.UpdateStatus(ctx, id, status)
	})
	if err != nil {
		return err
	}

	// Clear user cache by role id
	return a.clearUserCacheByID(ctx, id)
}
