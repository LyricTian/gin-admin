package dao

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"gorm.io/gorm"
)

func GetRoleMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return utilx.GetDB(ctx, defDB).Model(new(typed.RoleMenu))
}

type RoleMenuRepo struct {
	DB *gorm.DB
}

func (a *RoleMenuRepo) Query(ctx context.Context, params typed.RoleMenuQueryParam, opts ...typed.RoleMenuQueryOptions) (*typed.RoleMenuQueryResult, error) {
	var opt typed.RoleMenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetRoleMenuDB(ctx, a.DB)

	if v := params.RoleID; v != "" {
		db = db.Where("role_id=?", v)
	}

	var list typed.RoleMenus
	pr, err := utilx.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &typed.RoleMenuQueryResult{
		PageResult: pr,
		Data:       list,
	}

	return qr, nil
}

func (a *RoleMenuRepo) Get(ctx context.Context, id string, opts ...typed.RoleMenuQueryOptions) (*typed.RoleMenu, error) {
	var opt typed.RoleMenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(typed.RoleMenu)
	ok, err := utilx.FindOne(ctx, GetRoleMenuDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *RoleMenuRepo) Create(ctx context.Context, item *typed.RoleMenu) error {
	result := GetRoleMenuDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *RoleMenuRepo) Update(ctx context.Context, item *typed.RoleMenu) error {
	result := GetRoleMenuDB(ctx, a.DB).Where("id=?", item.ID).Omit("created_at").Updates(item)
	return errors.WithStack(result.Error)
}

func (a *RoleMenuRepo) Delete(ctx context.Context, id string) error {
	result := GetRoleMenuDB(ctx, a.DB).Where("id=?", id).Delete(new(typed.RoleMenu))
	return errors.WithStack(result.Error)
}

func (a *RoleMenuRepo) DeleteByRoleID(ctx context.Context, roleID string) error {
	result := GetRoleMenuDB(ctx, a.DB).Where("role_id=?", roleID).Delete(new(typed.RoleMenu))
	return errors.WithStack(result.Error)
}
