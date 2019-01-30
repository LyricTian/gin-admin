package model

import (
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/model/gorm/entity"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
)

// InitMenu 初始化菜单存储
func InitMenu(db *gormplus.DB) *Menu {
	db.AutoMigrate(new(entity.Menu))
	return NewMenu(db)
}

// NewMenu 实例化菜单存储
func NewMenu(db *gormplus.DB) *Menu {
	return &Menu{db: db}
}

// Menu 菜单存储
type Menu struct {
	db *gormplus.DB
}

func (a *Menu) getFuncName(name string) string {
	return fmt.Sprintf("gorm.model.Menu.%s", name)
}

func (a *Menu) getMenuDB(ctx context.Context) *gormplus.DB {
	return FromDBWithModel(ctx, a.db, entity.Menu{})
}

func (a *Menu) getQueryOption(opts ...schema.MenuQueryOptions) schema.MenuQueryOptions {
	if len(opts) > 0 {
		return opts[0]
	}
	return schema.MenuQueryOptions{}
}

// Query 查询数据
func (a *Menu) Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (schema.MenuQueryResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := a.getMenuDB(ctx).DB
	if v := params.RecordIDs; len(v) > 0 {
		db = db.Where("record_id IN(?)", v)
	}
	if v := params.Code; v != "" {
		db = db.Where("code LIKE ?", "%"+v+"%")
	}
	if v := params.Name; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.Types; len(v) > 0 {
		db = db.Where("type IN(?)", v)
	}
	if v := params.IsHide; v > 0 {
		db = db.Where("is_hide=?", v)
	}
	if v := params.ParentID; v != nil {
		db = db.Where("parent_id=?", *v)
	}
	if v := params.UserID; v != "" {
		userRoleQuery := FromDBWithModel(ctx, a.db, entity.UserRole{}).Select("role_id").Where("user_id=?", v).SubQuery()
		roleMenuQuery := FromDBWithModel(ctx, a.db, entity.RoleMenu{}).Select("menu_id").Where("role_id IN(?)", userRoleQuery).SubQuery()
		db = db.Where("record_id IN(?)", roleMenuQuery)
	}
	if v := params.ParentPath; v != "" {
		db = db.Where("parent_path LIKE ?", v+"%")
	}
	if v := params.ParentPaths; len(v) > 0 {
		db = db.Where("parent_path IN(?)", v)
	}
	db = db.Order("sequence,id DESC")

	var qr schema.MenuQueryResult
	opt := a.getQueryOption(opts...)
	var items entity.Menus
	pr, err := WrapPageQuery(db, opt.PageParam, &items)
	if err != nil {
		span.Errorf(err.Error())
		return qr, errors.New("查询数据发生错误")
	}
	qr.PageResult = pr
	qr.Data = items.ToSchemaMenus()

	return qr, nil
}

// Get 查询指定数据
func (a *Menu) Get(ctx context.Context, recordID string) (*schema.Menu, error) {
	span := logger.StartSpan(ctx, "查询指定数据", a.getFuncName("Get"))
	defer span.Finish()

	var item entity.Menu
	ok, err := a.db.FindOne(a.getMenuDB(ctx).Where("record_id=?", recordID), &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaMenu(), nil
}

// CheckCode 检查编号是否存在
func (a *Menu) CheckCode(ctx context.Context, code string, parentID string) (bool, error) {
	span := logger.StartSpan(ctx, "检查编号是否存在", a.getFuncName("CheckCode"))
	defer span.Finish()

	db := a.getMenuDB(ctx).Where("code=? AND parent_id=?", code, parentID)
	exists, err := a.db.Check(db)
	if err != nil {
		span.Errorf(err.Error())
		return false, errors.New("检查编号是否存在发生错误")
	}
	return exists, nil
}

// CheckChild 检查子级是否存在
func (a *Menu) CheckChild(ctx context.Context, parentID string) (bool, error) {
	span := logger.StartSpan(ctx, "检查子级是否存在", a.getFuncName("CheckChild"))
	defer span.Finish()

	db := a.getMenuDB(ctx).Where("parent_id=?", parentID)
	exists, err := a.db.Check(db)
	if err != nil {
		span.Errorf(err.Error())
		return false, errors.New("检查子级是否存在发生错误")
	}
	return exists, nil
}

// Create 创建数据
func (a *Menu) Create(ctx context.Context, item schema.Menu) error {
	span := logger.StartSpan(ctx, "创建数据", a.getFuncName("Create"))
	defer span.Finish()

	menu := entity.SchemaMenu(item).ToMenu()
	menu.Creator = FromUserID(ctx)
	result := a.getMenuDB(ctx).Create(menu)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("创建数据发生错误")
	}
	return nil
}

// Update 更新数据
func (a *Menu) Update(ctx context.Context, recordID string, item schema.Menu) error {
	span := logger.StartSpan(ctx, "更新数据", a.getFuncName("Update"))
	defer span.Finish()

	menu := entity.SchemaMenu(item).ToMenu()
	result := a.getMenuDB(ctx).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(menu)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新数据发生错误")
	}
	return nil
}

// UpdateParentPath 更新父级路径
func (a *Menu) UpdateParentPath(ctx context.Context, recordID, parentPath string) error {
	span := logger.StartSpan(ctx, "更新父级路径", a.getFuncName("UpdateParentPath"))
	defer span.Finish()

	result := a.getMenuDB(ctx).Where("record_id=?", recordID).Update("parent_path", parentPath)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新父级路径发生错误")
	}
	return nil
}

// Delete 删除数据
func (a *Menu) Delete(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除数据", a.getFuncName("Delete"))
	defer span.Finish()

	result := a.getMenuDB(ctx).Where("record_id=?", recordID).Delete(entity.Menu{})
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("删除数据发生错误")
	}
	return nil
}
