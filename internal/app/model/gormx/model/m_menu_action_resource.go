package model

import (
	"context"

	"github.com/LyricTian/gin-admin/v7/internal/app/model/gormx/entity"
	"github.com/LyricTian/gin-admin/v7/internal/app/schema"
	"github.com/LyricTian/gin-admin/v7/pkg/errors"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

// MenuActionResourceSet 注入MenuActionResource
var MenuActionResourceSet = wire.NewSet(wire.Struct(new(MenuActionResource), "*"))

// MenuActionResource 菜单动作关联资源存储
type MenuActionResource struct {
	DB *gorm.DB
}

func (a *MenuActionResource) getQueryOption(opts ...schema.MenuActionResourceQueryOptions) schema.MenuActionResourceQueryOptions {
	var opt schema.MenuActionResourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *MenuActionResource) Query(ctx context.Context, params schema.MenuActionResourceQueryParam, opts ...schema.MenuActionResourceQueryOptions) (*schema.MenuActionResourceQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetMenuActionResourceDB(ctx, a.DB)
	if v := params.MenuID; v != "" {
		subQuery := entity.GetMenuActionDB(ctx, a.DB).
			Where("menu_id=?", v).
			Select("id").SubQuery()
		db = db.Where("action_id IN ?", subQuery)
	}
	if v := params.MenuIDs; len(v) > 0 {
		subQuery := entity.GetMenuActionDB(ctx, a.DB).Where("menu_id IN (?)", v).Select("id").SubQuery()
		db = db.Where("action_id IN ?", subQuery)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByASC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.MenuActionResources
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.MenuActionResourceQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaMenuActionResources(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *MenuActionResource) Get(ctx context.Context, id string, opts ...schema.MenuActionResourceQueryOptions) (*schema.MenuActionResource, error) {
	db := entity.GetMenuActionResourceDB(ctx, a.DB).Where("id=?", id)
	var item entity.MenuActionResource
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaMenuActionResource(), nil
}

// Create 创建数据
func (a *MenuActionResource) Create(ctx context.Context, item schema.MenuActionResource) error {
	eitem := entity.SchemaMenuActionResource(item).ToMenuActionResource()
	result := entity.GetMenuActionResourceDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

// Update 更新数据
func (a *MenuActionResource) Update(ctx context.Context, id string, item schema.MenuActionResource) error {
	eitem := entity.SchemaMenuActionResource(item).ToMenuActionResource()
	result := entity.GetMenuActionResourceDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// Delete 删除数据
func (a *MenuActionResource) Delete(ctx context.Context, id string) error {
	result := entity.GetMenuActionResourceDB(ctx, a.DB).Where("id=?", id).Delete(entity.MenuActionResource{})
	return errors.WithStack(result.Error)
}

// DeleteByActionID 根据动作ID删除数据
func (a *MenuActionResource) DeleteByActionID(ctx context.Context, actionID string) error {
	result := entity.GetMenuActionResourceDB(ctx, a.DB).Where("action_id =?", actionID).Delete(entity.MenuActionResource{})
	return errors.WithStack(result.Error)
}

// DeleteByMenuID 根据菜单ID删除数据
func (a *MenuActionResource) DeleteByMenuID(ctx context.Context, menuID string) error {
	subQuery := entity.GetMenuActionDB(ctx, a.DB).Where("menu_id=?", menuID).Select("id").SubQuery()
	result := entity.GetMenuActionResourceDB(ctx, a.DB).Where("action_id IN ?", subQuery).Delete(entity.MenuActionResource{})
	return errors.WithStack(result.Error)
}
