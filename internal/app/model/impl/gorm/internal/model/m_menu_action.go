package model

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/errors"
	"github.com/LyricTian/gin-admin/internal/app/model/impl/gorm/internal/entity"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// NewMenuAction 创建菜单动作存储实例
func NewMenuAction(db *gorm.DB) *MenuAction {
	return &MenuAction{db}
}

// MenuAction 菜单动作存储
type MenuAction struct {
	db *gorm.DB
}

func (a *MenuAction) getQueryOption(opts ...schema.MenuActionQueryOptions) schema.MenuActionQueryOptions {
	var opt schema.MenuActionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *MenuAction) Query(ctx context.Context, params schema.MenuActionQueryParam, opts ...schema.MenuActionQueryOptions) (*schema.MenuActionQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetMenuActionDB(ctx, a.db)
	if v := params.MenuID; v != "" {
		db = db.Where("menu_id=?", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.MenuActions
	pr, err := WrapPageQuery(ctx, db, opt.PageParam, &list)
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
func (a *MenuAction) Get(ctx context.Context, recordID string, opts ...schema.MenuActionQueryOptions) (*schema.MenuAction, error) {
	db := entity.GetMenuActionDB(ctx, a.db).Where("record_id=?", recordID)
	var item entity.MenuAction
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaMenuAction(), nil
}

// Create 创建数据
func (a *MenuAction) Create(ctx context.Context, item schema.MenuAction) error {
	eitem := entity.SchemaMenuAction(item).ToMenuAction()
	result := entity.GetMenuActionDB(ctx, a.db).Create(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *MenuAction) Update(ctx context.Context, recordID string, item schema.MenuAction) error {
	eitem := entity.SchemaMenuAction(item).ToMenuAction()
	result := entity.GetMenuActionDB(ctx, a.db).Where("record_id=?", recordID).Omit("record_id").Updates(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *MenuAction) Delete(ctx context.Context, recordID string) error {
	result := entity.GetMenuActionDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.MenuAction{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeleteByMenuID 根据菜单ID删除数据
func (a *MenuAction) DeleteByMenuID(ctx context.Context, menuID string) error {
	result := entity.GetMenuActionDB(ctx, a.db).Where("menu_id=?", menuID).Delete(entity.MenuAction{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
