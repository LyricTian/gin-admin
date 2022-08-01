package dao

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"gorm.io/gorm"
)

func GetMenuActionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return utilx.GetDB(ctx, defDB).Model(new(typed.MenuAction))
}

type MenuActionRepo struct {
	DB *gorm.DB
}

func (a *MenuActionRepo) Query(ctx context.Context, params typed.MenuActionQueryParam, opts ...typed.MenuActionQueryOptions) (*typed.MenuActionQueryResult, error) {
	var opt typed.MenuActionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetMenuActionDB(ctx, a.DB)

	if v := params.MenuID; v != "" {
		db = db.Where("menu_id=?", v)
	}
	if v := params.UserID; v != "" {
		userRoleQuery := GetUserRoleDB(ctx, a.DB).Select("role_id").Where("user_id=?", v)
		roleMenuQuery := GetRoleMenuDB(ctx, a.DB).Distinct("menu_action_id").Where("role_id in (?)", userRoleQuery)
		db = db.Where("id in (?)", roleMenuQuery)
	}

	var list typed.MenuActions
	pr, err := utilx.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &typed.MenuActionQueryResult{
		PageResult: pr,
		Data:       list,
	}

	return qr, nil
}

func (a *MenuActionRepo) Get(ctx context.Context, id string, opts ...typed.MenuActionQueryOptions) (*typed.MenuAction, error) {
	var opt typed.MenuActionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(typed.MenuAction)
	ok, err := utilx.FindOne(ctx, GetMenuActionDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *MenuActionRepo) Create(ctx context.Context, item *typed.MenuAction) error {
	result := GetMenuActionDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *MenuActionRepo) Update(ctx context.Context, item *typed.MenuAction) error {
	result := GetMenuActionDB(ctx, a.DB).Where("id=?", item.ID).Omit("created_at").Updates(item)
	return errors.WithStack(result.Error)
}

func (a *MenuActionRepo) Delete(ctx context.Context, id string) error {
	result := GetMenuActionDB(ctx, a.DB).Where("id=?", id).Delete(new(typed.MenuAction))
	return errors.WithStack(result.Error)
}

func (a *MenuActionRepo) DeleteByMenuID(ctx context.Context, menuID string) error {
	result := GetMenuActionDB(ctx, a.DB).Where("menu_id=?", menuID).Delete(new(typed.MenuAction))
	return errors.WithStack(result.Error)
}
