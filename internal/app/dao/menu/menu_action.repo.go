package menu

import (
	"context"

	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao/util"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/errors"
)

// MenuActionSet 注入MenuAction
var MenuActionSet = wire.NewSet(wire.Struct(new(MenuActionRepo), "*"))

// MenuActionRepo 菜单动作存储
type MenuActionRepo struct {
	DB *gorm.DB
}

func (a *MenuActionRepo) getQueryOption(opts ...schema.MenuActionQueryOptions) schema.MenuActionQueryOptions {
	var opt schema.MenuActionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *MenuActionRepo) Query(ctx context.Context, params schema.MenuActionQueryParam, opts ...schema.MenuActionQueryOptions) (*schema.MenuActionQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := GetMenuActionDB(ctx, a.DB)
	if v := params.MenuID; v > 0 {
		db = db.Where("menu_id=?", v)
	}
	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByASC))
	db = db.Order(util.ParseOrder(opt.OrderFields))

	var list MenuActions
	pr, err := util.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.MenuActionQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaMenuActions(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *MenuActionRepo) Get(ctx context.Context, id uint64, opts ...schema.MenuActionQueryOptions) (*schema.MenuAction, error) {
	db := GetMenuActionDB(ctx, a.DB).Where("id=?", id)
	var item MenuAction
	ok, err := util.FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaMenuAction(), nil
}

// Create 创建数据
func (a *MenuActionRepo) Create(ctx context.Context, item schema.MenuAction) error {
	eitem := SchemaMenuAction(item).ToMenuAction()
	result := GetMenuActionDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

// Update 更新数据
func (a *MenuActionRepo) Update(ctx context.Context, id uint64, item schema.MenuAction) error {
	eitem := SchemaMenuAction(item).ToMenuAction()
	result := GetMenuActionDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// Delete 删除数据
func (a *MenuActionRepo) Delete(ctx context.Context, id uint64) error {
	result := GetMenuActionDB(ctx, a.DB).Where("id=?", id).Delete(MenuAction{})
	return errors.WithStack(result.Error)
}

// DeleteByMenuID 根据菜单ID删除数据
func (a *MenuActionRepo) DeleteByMenuID(ctx context.Context, menuID uint64) error {
	result := GetMenuActionDB(ctx, a.DB).Where("menu_id=?", menuID).Delete(MenuAction{})
	return errors.WithStack(result.Error)
}
