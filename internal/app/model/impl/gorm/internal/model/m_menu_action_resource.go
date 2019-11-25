package model

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/errors"
	"github.com/LyricTian/gin-admin/internal/app/model/impl/gorm/internal/entity"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// NewMenuActionResource 创建菜单动作关联资源存储实例
func NewMenuActionResource(db *gorm.DB) *MenuActionResource {
	return &MenuActionResource{db}
}

// MenuActionResource 菜单动作关联资源存储
type MenuActionResource struct {
	db *gorm.DB
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

	db := entity.GetMenuActionResourceDB(ctx, a.db)
	if v := params.MenuID; v != "" {
		subQuery := entity.GetMenuActionDB(ctx, a.db).Where("menu_id=?", v).Select("record_id").SubQuery()
		db = db.Where("action_id IN(?)", subQuery)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.MenuActionResources
	pr, err := WrapPageQuery(ctx, db, opt.PageParam, &list)
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
func (a *MenuActionResource) Get(ctx context.Context, recordID string, opts ...schema.MenuActionResourceQueryOptions) (*schema.MenuActionResource, error) {
	db := entity.GetMenuActionResourceDB(ctx, a.db).Where("record_id=?", recordID)
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
	result := entity.GetMenuActionResourceDB(ctx, a.db).Create(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *MenuActionResource) Update(ctx context.Context, recordID string, item schema.MenuActionResource) error {
	eitem := entity.SchemaMenuActionResource(item).ToMenuActionResource()
	result := entity.GetMenuActionResourceDB(ctx, a.db).Where("record_id=?", recordID).Omit("record_id").Updates(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *MenuActionResource) Delete(ctx context.Context, recordID string) error {
	result := entity.GetMenuActionResourceDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.MenuActionResource{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeleteByActionID 根据动作ID删除数据
func (a *MenuActionResource) DeleteByActionID(ctx context.Context, actionID string) error {
	result := entity.GetMenuActionResourceDB(ctx, a.db).Where("action_id IN(?)", actionID).Delete(entity.MenuAction{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeleteByMenuID 根据菜单ID删除数据
func (a *MenuActionResource) DeleteByMenuID(ctx context.Context, menuID string) error {
	subQuery := entity.GetMenuActionDB(ctx, a.db).Where("menu_id=?", menuID).Select("record_id").SubQuery()
	result := entity.GetMenuActionResourceDB(ctx, a.db).Where("action_id IN(?)", subQuery).Delete(entity.MenuAction{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
