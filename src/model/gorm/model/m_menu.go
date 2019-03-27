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
	db.AutoMigrate(
		new(entity.Menu),
		new(entity.MenuAction),
		new(entity.MenuResource),
	)
	a.db = db
	return a
}

func (a *Menu) getFuncName(name string) string {
	return fmt.Sprintf("gorm.model.Menu.%s", name)
}

func (a *Menu) getQueryOption(opts ...schema.MenuQueryOptions) schema.MenuQueryOptions {
	var opt schema.MenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *Menu) Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (*schema.MenuQueryResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := entity.GetMenuDB(ctx, a.db).DB
	if v := params.RecordIDs; len(v) > 0 {
		db = db.Where("record_id IN(?)", v)
	}
	if v := params.LikeName; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.ParentID; v != nil {
		db = db.Where("parent_id=?", *v)
	}
	if v := params.PrefixParentPath; v != "" {
		db = db.Where("parent_path LIKE ?", v+"%")
	}
	if v := params.Hidden; v != nil {
		db = db.Where("hidden=?", *v)
	}
	db = db.Order("sequence DESC,id DESC")

	opt := a.getQueryOption(opts...)
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

	for _, item := range qr.Data {
		err := a.fillSchemaMenu(ctx, item, opts...)
		if err != nil {
			return nil, err
		}
	}

	return qr, nil
}

// 填充菜单对象数据
func (a *Menu) fillSchemaMenu(ctx context.Context, item *schema.Menu, opts ...schema.MenuQueryOptions) error {
	opt := a.getQueryOption(opts...)

	if opt.IncludeActions {
		list, err := a.queryActions(ctx, item.RecordID)
		if err != nil {
			return err
		}
		item.Actions = list.ToSchemaMenuActions()
	}

	if opt.IncludeResources {
		list, err := a.queryResources(ctx, item.RecordID)
		if err != nil {
			return err
		}
		item.Resources = list.ToSchemaMenuResources()
	}
	return nil
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

	sitem := item.ToSchemaMenu()
	err = a.fillSchemaMenu(ctx, sitem, opts...)
	if err != nil {
		return nil, err
	}

	return sitem, nil
}

// Create 创建数据
func (a *Menu) Create(ctx context.Context, item schema.Menu) error {
	span := logger.StartSpan(ctx, "创建数据", a.getFuncName("Create"))
	defer span.Finish()

	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaMenu(item)
		result := entity.GetMenuDB(ctx, a.db).Create(sitem.ToMenu())
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("创建菜单数据发生错误")
		}

		for _, item := range sitem.ToMenuActions() {
			result := entity.GetMenuActionDB(ctx, a.db).Create(item)
			if err := result.Error; err != nil {
				span.Errorf(err.Error())
				return errors.New("创建菜单动作数据发生错误")
			}
		}

		for _, item := range sitem.ToMenuResources() {
			result := entity.GetMenuResourceDB(ctx, a.db).Create(item)
			if err := result.Error; err != nil {
				span.Errorf(err.Error())
				return errors.New("创建菜单资源数据发生错误")
			}
		}

		return nil
	})
}

// 对比并获取需要新增，修改，删除的动作项
func (a *Menu) compareUpdateAction(oldList, newList entity.MenuActions) (clist, dlist, ulist []*entity.MenuAction) {
	oldMap, newMap := oldList.ToMap(), newList.ToMap()

	for _, nitem := range newList {
		if _, ok := oldMap[nitem.Code]; ok {
			ulist = append(ulist, nitem)
			continue
		}
		clist = append(clist, nitem)
	}

	for _, oitem := range oldList {
		if _, ok := newMap[oitem.Code]; !ok {
			dlist = append(dlist, oitem)
		}
	}
	return
}

// 更新动作数据
func (a *Menu) updateActions(ctx context.Context, span *logger.Entry, menuID string, items entity.MenuActions) error {
	list, err := a.queryActions(ctx, menuID)
	if err != nil {
		return err
	}

	clist, dlist, ulist := a.compareUpdateAction(list, items)
	for _, item := range clist {
		result := entity.GetMenuActionDB(ctx, a.db).Create(item)
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("创建菜单动作数据发生错误")
		}
	}

	for _, item := range dlist {
		result := entity.GetMenuActionDB(ctx, a.db).Where("menu_id=? AND code=?", menuID, item.Code).Delete(entity.MenuAction{})
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("删除菜单动作数据发生错误")
		}
	}

	for _, item := range ulist {
		result := entity.GetMenuActionDB(ctx, a.db).Where("menu_id=? AND code=?", menuID, item.Code).Omit("menu_id", "code").Updates(item)
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("更新菜单动作数据发生错误")
		}
	}
	return nil
}

// 对比并获取需要新增，修改，删除的资源项
func (a *Menu) compareUpdateResource(oldList, newList entity.MenuResources) (clist, dlist, ulist []*entity.MenuResource) {
	oldMap, newMap := oldList.ToMap(), newList.ToMap()

	for _, nitem := range newList {
		if _, ok := oldMap[nitem.Code]; ok {
			ulist = append(ulist, nitem)
			continue
		}
		clist = append(clist, nitem)
	}

	for _, oitem := range oldList {
		if _, ok := newMap[oitem.Code]; !ok {
			dlist = append(dlist, oitem)
		}
	}
	return
}

// 更新资源数据
func (a *Menu) updateResources(ctx context.Context, span *logger.Entry, menuID string, items entity.MenuResources) error {
	list, err := a.queryResources(ctx, menuID)
	if err != nil {
		return err
	}

	clist, dlist, ulist := a.compareUpdateResource(list, items)
	for _, item := range clist {
		result := entity.GetMenuResourceDB(ctx, a.db).Create(item)
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("创建菜单资源数据发生错误")
		}
	}

	for _, item := range dlist {
		result := entity.GetMenuResourceDB(ctx, a.db).Where("menu_id=? AND code=?", menuID, item.Code).Delete(entity.MenuResource{})
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("删除菜单资源数据发生错误")
		}
	}

	for _, item := range ulist {
		result := entity.GetMenuResourceDB(ctx, a.db).Where("menu_id=? AND code=?", menuID, item.Code).Omit("menu_id", "code").Updates(item)
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("更新菜单资源数据发生错误")
		}
	}
	return nil
}

// Update 更新数据
func (a *Menu) Update(ctx context.Context, recordID string, item schema.Menu) error {
	span := logger.StartSpan(ctx, "更新数据", a.getFuncName("Update"))
	defer span.Finish()

	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaMenu(item)
		result := entity.GetMenuDB(ctx, a.db).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(sitem.ToMenu())
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("更新数据发生错误")
		}

		err := a.updateActions(ctx, span, recordID, sitem.ToMenuActions())
		if err != nil {
			return err
		}

		err = a.updateResources(ctx, span, recordID, sitem.ToMenuResources())
		if err != nil {
			return err
		}

		return nil
	})
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

	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		result := entity.GetMenuDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.Menu{})
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("删除数据发生错误")
		}

		result = entity.GetMenuActionDB(ctx, a.db).Where("menu_id=?", recordID).Delete(entity.MenuAction{})
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("删除菜单动作数据发生错误")
		}

		result = entity.GetMenuResourceDB(ctx, a.db).Where("menu_id=?", recordID).Delete(entity.MenuResource{})
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("删除菜单资源数据发生错误")
		}
		return nil
	})
}

func (a *Menu) queryActions(ctx context.Context, menuID string) (entity.MenuActions, error) {
	span := logger.StartSpan(ctx, "查询菜单动作数据", a.getFuncName("queryActions"))
	defer span.Finish()

	var list entity.MenuActions
	result := entity.GetMenuActionDB(ctx, a.db).Where("menu_id=?", menuID).Find(&list)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询菜单动作数据发生错误")
	}

	return list, nil
}

func (a *Menu) queryResources(ctx context.Context, menuID string) (entity.MenuResources, error) {
	span := logger.StartSpan(ctx, "查询菜单资源数据", a.getFuncName("queryResources"))
	defer span.Finish()

	var list entity.MenuResources
	result := entity.GetMenuResourceDB(ctx, a.db).Where("menu_id=?", menuID).Find(&list)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询菜单资源数据发生错误")
	}

	return list, nil
}
