package model

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/internal/app/model/impl/gorm/entity"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var _ model.IRoleMenu = (*RoleMenu)(nil)

// RoleMenuSet 注入RoleMenu
var RoleMenuSet = wire.NewSet(wire.Struct(new(RoleMenu), "*"), wire.Bind(new(model.IRoleMenu), new(*RoleMenu)))

// RoleMenu 角色菜单存储
type RoleMenu struct {
	DB *gorm.DB
}

func (a *RoleMenu) getQueryOption(opts ...schema.RoleMenuQueryOptions) schema.RoleMenuQueryOptions {
	var opt schema.RoleMenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *RoleMenu) Query(ctx context.Context, params schema.RoleMenuQueryParam, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenuQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetRoleMenuDB(ctx, a.DB)
	if v := params.RoleID; v != "" {
		db = db.Where("role_id=?", v)
	}
	if v := params.RoleIDs; len(v) > 0 {
		db = db.Where("role_id IN(?)", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.RoleMenus
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RoleMenuQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRoleMenus(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *RoleMenu) Get(ctx context.Context, recordID string, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenu, error) {
	db := entity.GetRoleMenuDB(ctx, a.DB).Where("record_id=?", recordID)
	var item entity.RoleMenu
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaRoleMenu(), nil
}

// Create 创建数据
func (a *RoleMenu) Create(ctx context.Context, item schema.RoleMenu) error {
	eitem := entity.SchemaRoleMenu(item).ToRoleMenu()
	result := entity.GetRoleMenuDB(ctx, a.DB).Create(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *RoleMenu) Update(ctx context.Context, recordID string, item schema.RoleMenu) error {
	eitem := entity.SchemaRoleMenu(item).ToRoleMenu()
	result := entity.GetRoleMenuDB(ctx, a.DB).Where("record_id=?", recordID).Updates(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *RoleMenu) Delete(ctx context.Context, recordID string) error {
	result := entity.GetRoleMenuDB(ctx, a.DB).Where("record_id=?", recordID).Delete(entity.RoleMenu{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeleteByRoleID 根据角色ID删除数据
func (a *RoleMenu) DeleteByRoleID(ctx context.Context, roleID string) error {
	result := entity.GetRoleMenuDB(ctx, a.DB).Where("role_id=?", roleID).Delete(entity.RoleMenu{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
