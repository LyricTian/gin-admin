package dao

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"gorm.io/gorm"
)

func GetUserRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return utilx.GetDB(ctx, defDB).Model(new(typed.UserRole))
}

type UserRoleRepo struct {
	DB *gorm.DB
}

func (a *UserRoleRepo) Query(ctx context.Context, params typed.UserRoleQueryParam, opts ...typed.UserRoleQueryOptions) (*typed.UserRoleQueryResult, error) {
	var opt typed.UserRoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetUserRoleDB(ctx, a.DB)

	if v := params.RoleID; v != "" {
		db = db.Where("role_id=?", v)
	}
	if v := params.UserID; v != "" {
		db = db.Where("user_id=?", v)
	}
	if v := params.UserIDList; len(v) > 0 {
		db = db.Where("user_id IN (?)", v)
	}

	var list typed.UserRoles
	pr, err := utilx.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &typed.UserRoleQueryResult{
		PageResult: pr,
		Data:       list,
	}

	return qr, nil
}

func (a *UserRoleRepo) Get(ctx context.Context, id string, opts ...typed.UserRoleQueryOptions) (*typed.UserRole, error) {
	var opt typed.UserRoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(typed.UserRole)
	ok, err := utilx.FindOne(ctx, GetUserRoleDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *UserRoleRepo) GetRoleIDsByUserID(ctx context.Context, userID string) ([]string, error) {
	var roleIDs []string
	result := GetUserRoleDB(ctx, a.DB).Where("user_id=?", userID).Pluck("role_id", &roleIDs)
	return roleIDs, errors.WithStack(result.Error)
}

func (a *UserRoleRepo) Create(ctx context.Context, item *typed.UserRole) error {
	result := GetUserRoleDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *UserRoleRepo) Update(ctx context.Context, item *typed.UserRole) error {
	result := GetUserRoleDB(ctx, a.DB).Where("id=?", item.ID).Omit("created_at").Updates(item)
	return errors.WithStack(result.Error)
}

func (a *UserRoleRepo) Delete(ctx context.Context, id string) error {
	result := GetUserRoleDB(ctx, a.DB).Where("id=?", id).Delete(new(typed.UserRole))
	return errors.WithStack(result.Error)
}

func (a *UserRoleRepo) DeleteByUserID(ctx context.Context, userID string) error {
	result := GetUserRoleDB(ctx, a.DB).Where("user_id=?", userID).Delete(new(typed.UserRole))
	return errors.WithStack(result.Error)
}

func (a *UserRoleRepo) DeleteByRoleID(ctx context.Context, roleID string) error {
	result := GetUserRoleDB(ctx, a.DB).Where("role_id=?", roleID).Delete(new(typed.UserRole))
	return errors.WithStack(result.Error)
}
