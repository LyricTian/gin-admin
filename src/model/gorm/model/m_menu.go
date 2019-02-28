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

// Menu 菜单存储
type Menu struct {
	db *gormplus.DB
}

// Init 初始化
func (a *Menu) Init(db *gormplus.DB) *Menu {
	db.AutoMigrate(new(entity.Menu))
	a.db = db
	return a
}

func (a *Menu) getFuncName(name string) string {
	return fmt.Sprintf("gorm.model.Menu.%s", name)
}

// Query 查询数据
func (a *Menu) Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (*schema.MenuQueryResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := entity.GetMenuDB(ctx, a.db).DB
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
	if v := params.ParentID; v != nil {
		db = db.Where("parent_id=?", *v)
	}
	if v := params.UserID; v != "" {
		userRoleQuery := entity.GetUserRoleDB(ctx, a.db).Select("role_id").Where("user_id=?", v).SubQuery()
		roleMenuQuery := entity.GetRoleMenuDB(ctx, a.db).Select("menu_id").Where("role_id IN(?)", userRoleQuery).SubQuery()
		db = db.Where("record_id IN(?)", roleMenuQuery)
	}
	if v := params.ParentPath; v != "" {
		db = db.Where("parent_path LIKE ?", v+"%")
	}
	db = db.Order("sequence DESC,id DESC")

	var opt schema.MenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	var list entity.Menus
	pr, err := WrapPageQuery(db, opt.PageParam, &list)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询数据发生错误")
	}
	qr := &schema.MenuQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaMenus(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *Menu) Get(ctx context.Context, recordID string, opts ...schema.MenuQueryOptions) (*schema.Menu, error) {
	span := logger.StartSpan(ctx, "查询指定数据", a.getFuncName("Get"))
	defer span.Finish()

	var item entity.Menu
	ok, err := a.db.FindOne(entity.GetMenuDB(ctx, a.db).Where("record_id=?", recordID), &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaMenu(), nil
}

// CheckCodeWithParentID 检查同一父级下编号是否存在
func (a *Menu) CheckCodeWithParentID(ctx context.Context, code, parentID string) (bool, error) {
	span := logger.StartSpan(ctx, "检查同一父级下编号是否存在", a.getFuncName("CheckCodeWithParentID"))
	defer span.Finish()

	db := entity.GetMenuDB(ctx, a.db).Where("code=? AND parent_id=?", code, parentID)
	exists, err := a.db.Check(db)
	if err != nil {
		span.Errorf(err.Error())
		return false, errors.New("检查同一父级下编号是否存在发生错误")
	}
	return exists, nil
}

// CheckChild 检查子级是否存在
func (a *Menu) CheckChild(ctx context.Context, parentID string) (bool, error) {
	span := logger.StartSpan(ctx, "检查子级是否存在", a.getFuncName("CheckChild"))
	defer span.Finish()

	db := entity.GetMenuDB(ctx, a.db).Where("parent_id=?", parentID)
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
	result := entity.GetMenuDB(ctx, a.db).Create(menu)
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
	result := entity.GetMenuDB(ctx, a.db).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(menu)
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

	result := entity.GetMenuDB(ctx, a.db).Where("record_id=?", recordID).Update("parent_path", parentPath)
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

	result := entity.GetMenuDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.Menu{})
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("删除数据发生错误")
	}
	return nil
}
