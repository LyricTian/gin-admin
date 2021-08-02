package menu

import (
	"context"

	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao/util"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/errors"
)

// MenuSet 注入Menu
var MenuSet = wire.NewSet(wire.Struct(new(MenuRepo), "*"))

// MenuRepo 菜单存储
type MenuRepo struct {
	DB *gorm.DB
}

func (a *MenuRepo) getQueryOption(opts ...schema.MenuQueryOptions) schema.MenuQueryOptions {
	var opt schema.MenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *MenuRepo) Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (*schema.MenuQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := GetMenuDB(ctx, a.DB)
	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}
	if v := params.Name; v != "" {
		db = db.Where("name=?", v)
	}
	if v := params.ParentID; v != nil {
		db = db.Where("parent_id=?", *v)
	}
	if v := params.PrefixParentPath; v != "" {
		db = db.Where("parent_path LIKE ?", v+"%")
	}
	if v := params.IsShow; v != 0 {
		db = db.Where("show_status=?", v)
	}
	if v := params.Status; v != 0 {
		db = db.Where("status=?", v)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ?", v)
	}

	if len(opt.SelectFields) > 0 {
		db = db.Select(opt.SelectFields)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(util.ParseOrder(opt.OrderFields))

	var list Menus
	pr, err := util.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.MenuQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaMenus(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *MenuRepo) Get(ctx context.Context, id uint64, opts ...schema.MenuQueryOptions) (*schema.Menu, error) {
	var item Menu
	ok, err := util.FindOne(ctx, GetMenuDB(ctx, a.DB).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaMenu(), nil
}

// Create 创建数据
func (a *MenuRepo) Create(ctx context.Context, item schema.Menu) error {
	eitem := SchemaMenu(item).ToMenu()
	result := GetMenuDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

// Update 更新数据
func (a *MenuRepo) Update(ctx context.Context, id uint64, item schema.Menu) error {
	eitem := SchemaMenu(item).ToMenu()
	result := GetMenuDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// UpdateParentPath 更新父级路径
func (a *MenuRepo) UpdateParentPath(ctx context.Context, id uint64, parentPath string) error {
	result := GetMenuDB(ctx, a.DB).Where("id=?", id).Update("parent_path", parentPath)
	return errors.WithStack(result.Error)
}

// Delete 删除数据
func (a *MenuRepo) Delete(ctx context.Context, id uint64) error {
	result := GetMenuDB(ctx, a.DB).Where("id=?", id).Delete(Menu{})
	return errors.WithStack(result.Error)
}

// UpdateStatus 更新状态
func (a *MenuRepo) UpdateStatus(ctx context.Context, id uint64, status int) error {
	result := GetMenuDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
