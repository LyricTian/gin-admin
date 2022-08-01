package dao

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"gorm.io/gorm"
)

func GetMenuActionResourceDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return utilx.GetDB(ctx, defDB).Model(new(typed.MenuActionResource))
}

type MenuActionResourceRepo struct {
	DB *gorm.DB
}

func (a *MenuActionResourceRepo) Query(ctx context.Context, params typed.MenuActionResourceQueryParam, opts ...typed.MenuActionResourceQueryOptions) (*typed.MenuActionResourceQueryResult, error) {
	var opt typed.MenuActionResourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetMenuActionResourceDB(ctx, a.DB)

	if v := params.MenuID; v != "" {
		db = db.Where("menu_id=?", v)
	}
	if v := params.RoleID; v != "" {
		roleMenuQuery := GetRoleMenuDB(ctx, a.DB).Distinct("menu_action_id").Where("role_id=?", v)
		db = db.Where("action_id in (?)", roleMenuQuery)
	}

	var list typed.MenuActionResources
	pr, err := utilx.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &typed.MenuActionResourceQueryResult{
		PageResult: pr,
		Data:       list,
	}

	return qr, nil
}

func (a *MenuActionResourceRepo) Get(ctx context.Context, id string, opts ...typed.MenuActionResourceQueryOptions) (*typed.MenuActionResource, error) {
	var opt typed.MenuActionResourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(typed.MenuActionResource)
	ok, err := utilx.FindOne(ctx, GetMenuActionResourceDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *MenuActionResourceRepo) Create(ctx context.Context, item *typed.MenuActionResource) error {
	result := GetMenuActionResourceDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *MenuActionResourceRepo) Update(ctx context.Context, item *typed.MenuActionResource) error {
	result := GetMenuActionResourceDB(ctx, a.DB).Where("id=?", item.ID).Omit("created_at").Updates(item)
	return errors.WithStack(result.Error)
}

func (a *MenuActionResourceRepo) Delete(ctx context.Context, id string) error {
	result := GetMenuActionResourceDB(ctx, a.DB).Where("id=?", id).Delete(new(typed.MenuActionResource))
	return errors.WithStack(result.Error)
}

func (a *MenuActionResourceRepo) DeleteByActionID(ctx context.Context, actionID string) error {
	result := GetMenuActionResourceDB(ctx, a.DB).Where("action_id=?", actionID).Delete(new(typed.MenuActionResource))
	return errors.WithStack(result.Error)
}

func (a *MenuActionResourceRepo) DeleteByMenuID(ctx context.Context, menuID string) error {
	result := GetMenuActionResourceDB(ctx, a.DB).Where("menu_id=?", menuID).Delete(new(typed.MenuActionResource))
	return errors.WithStack(result.Error)
}
