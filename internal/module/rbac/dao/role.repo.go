package dao

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"gorm.io/gorm"
)

func GetRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return utilx.GetDB(ctx, defDB).Model(new(typed.Role))
}

type RoleRepo struct {
	DB *gorm.DB
}

func (a *RoleRepo) Query(ctx context.Context, params typed.RoleQueryParam, opts ...typed.RoleQueryOptions) (*typed.RoleQueryResult, error) {
	var opt typed.RoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetRoleDB(ctx, a.DB)

	if v := len(params.IDList); v > 0 {
		db = db.Where("id in (?)", v)
	}
	if v := params.LikeName; v != "" {
		db = db.Where("name like ?", "%"+v+"%")
	}
	if v := params.Status; v != "" {
		db = db.Where("status=?", v)
	}
	if v := params.GtUpdatedAt; v != nil {
		db = db.Where("updated_at > ?", *v)
	}

	var list typed.Roles
	pr, err := utilx.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &typed.RoleQueryResult{
		PageResult: pr,
		Data:       list,
	}

	return qr, nil
}

func (a *RoleRepo) Get(ctx context.Context, id string, opts ...typed.RoleQueryOptions) (*typed.Role, error) {
	var opt typed.RoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(typed.Role)
	ok, err := utilx.FindOne(ctx, GetRoleDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *RoleRepo) Exists(ctx context.Context, id string) (bool, error) {
	exists, err := utilx.Exists(ctx, GetRoleDB(ctx, a.DB).Where("id=?", id))
	return exists, errors.WithStack(err)
}

func (a *RoleRepo) Create(ctx context.Context, item *typed.Role) error {
	result := GetRoleDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *RoleRepo) Update(ctx context.Context, item *typed.Role) error {
	result := GetRoleDB(ctx, a.DB).Where("id=?", item.ID).Omit("created_at", "created_by").Updates(item)
	return errors.WithStack(result.Error)
}

func (a *RoleRepo) Delete(ctx context.Context, id string) error {
	result := GetRoleDB(ctx, a.DB).Where("id=?", id).Delete(new(typed.Role))
	return errors.WithStack(result.Error)
}

func (a *RoleRepo) UpdateStatus(ctx context.Context, id string, status typed.RoleStatus) error {
	result := GetRoleDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
