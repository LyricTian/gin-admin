package dao

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"

	"gorm.io/gorm"
)

func GetMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return utilx.GetDB(ctx, defDB).Model(new(typed.Menu))
}

type MenuRepo struct {
	DB *gorm.DB
}

func (a *MenuRepo) Query(ctx context.Context, params typed.MenuQueryParam, opts ...typed.MenuQueryOptions) (*typed.MenuQueryResult, error) {
	var opt typed.MenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetMenuDB(ctx, a.DB)

	if v := len(params.IDList); v > 0 {
		db = db.Where("id in (?)", v)
	}
	if v := params.LikeName; v != "" {
		db = db.Where("name like ?", "%"+v+"%")
	}
	if v := params.ParentID; v != nil {
		db = db.Where("parent_id=?", *v)
	}
	if v := params.ParentPathPrefix; v != "" {
		db = db.Where("parent_path like ?", v+"%")
	}
	if v := params.Status; v != "" {
		db = db.Where("status=?", v)
	}
	if v := params.UserID; v != "" {
		userRoleQuery := GetUserRoleDB(ctx, a.DB).Select("role_id").Where("user_id=?", v)
		roleMenuQuery := GetRoleMenuDB(ctx, a.DB).Distinct("menu_id").Where("role_id in (?)", userRoleQuery)
		db = db.Where("menu_id in (?)", roleMenuQuery)
	}

	var list typed.Menus
	pr, err := utilx.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &typed.MenuQueryResult{
		PageResult: pr,
		Data:       list,
	}

	return qr, nil
}

func (a *MenuRepo) Get(ctx context.Context, id string, opts ...typed.MenuQueryOptions) (*typed.Menu, error) {
	var opt typed.MenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(typed.Menu)
	ok, err := utilx.FindOne(ctx, GetMenuDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *MenuRepo) Create(ctx context.Context, item *typed.Menu) error {
	result := GetMenuDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *MenuRepo) Update(ctx context.Context, item *typed.Menu) error {
	result := GetMenuDB(ctx, a.DB).Where("id=?", item.ID).Omit("created_at", "created_by").Updates(item)
	return errors.WithStack(result.Error)
}

func (a *MenuRepo) Delete(ctx context.Context, id string) error {
	result := GetMenuDB(ctx, a.DB).Where("id=?", id).Delete(new(typed.Menu))
	return errors.WithStack(result.Error)
}

func (a *MenuRepo) UpdateParentPath(ctx context.Context, id string, parentPath string) error {
	result := GetMenuDB(ctx, a.DB).Where("id=?", id).Update("parent_path", parentPath)
	return errors.WithStack(result.Error)
}

func (a *MenuRepo) UpdateStatus(ctx context.Context, id string, status typed.MenuStatus) error {
	result := GetMenuDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
