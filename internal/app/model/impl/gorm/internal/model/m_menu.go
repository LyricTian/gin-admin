package model

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/errors"
	"github.com/LyricTian/gin-admin/internal/app/model/impl/gorm/internal/entity"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/gormplus"
)

// NewMenu 创建菜单存储实例
func NewMenu(db *gormplus.DB) *Menu {
	return &Menu{db}
}

// Menu 菜单存储
type Menu struct {
	db *gormplus.DB
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
		return nil, errors.WithStack(err)
	}
	qr := &schema.MenuQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaMenus(),
	}

	err = a.fillSchemaMenus(ctx, qr.Data, opts...)
	if err != nil {
		return nil, err
	}

	return qr, nil
}

// 填充菜单对象数据
func (a *Menu) fillSchemaMenus(ctx context.Context, items []*schema.Menu, opts ...schema.MenuQueryOptions) error {
	opt := a.getQueryOption(opts...)

	if opt.IncludeActions || opt.IncludeResources {

		menuIDs := make([]string, len(items))
		for i, item := range items {
			menuIDs[i] = item.RecordID
		}

		var actionList entity.MenuActions
		var resourceList entity.MenuResources
		if opt.IncludeActions {
			items, err := a.queryActions(ctx, menuIDs...)
			if err != nil {
				return err
			}
			actionList = items
		}

		if opt.IncludeResources {
			items, err := a.queryResources(ctx, menuIDs...)
			if err != nil {
				return err
			}
			resourceList = items
		}

		for i, item := range items {
			if len(actionList) > 0 {
				items[i].Actions = actionList.GetByMenuID(item.RecordID)
			}
			if len(resourceList) > 0 {
				items[i].Resources = resourceList.GetByMenuID(item.RecordID)
			}
		}
	}

	return nil
}

// Get 查询指定数据
func (a *Menu) Get(ctx context.Context, recordID string, opts ...schema.MenuQueryOptions) (*schema.Menu, error) {
	var item entity.Menu
	ok, err := a.db.FindOne(entity.GetMenuDB(ctx, a.db).Where("record_id=?", recordID), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	sitem := item.ToSchemaMenu()
	err = a.fillSchemaMenus(ctx, []*schema.Menu{sitem}, opts...)
	if err != nil {
		return nil, err
	}

	return sitem, nil
}

// Create 创建数据
func (a *Menu) Create(ctx context.Context, item schema.Menu) error {
	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaMenu(item)
		result := entity.GetMenuDB(ctx, a.db).Create(sitem.ToMenu())
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		for _, item := range sitem.ToMenuActions() {
			result := entity.GetMenuActionDB(ctx, a.db).Create(item)
			if err := result.Error; err != nil {
				return errors.WithStack(err)
			}
		}

		for _, item := range sitem.ToMenuResources() {
			result := entity.GetMenuResourceDB(ctx, a.db).Create(item)
			if err := result.Error; err != nil {
				return errors.WithStack(err)
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
func (a *Menu) updateActions(ctx context.Context, menuID string, items entity.MenuActions) error {
	list, err := a.queryActions(ctx, menuID)
	if err != nil {
		return err
	}

	clist, dlist, ulist := a.compareUpdateAction(list, items)
	for _, item := range clist {
		result := entity.GetMenuActionDB(ctx, a.db).Create(item)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
	}

	for _, item := range dlist {
		result := entity.GetMenuActionDB(ctx, a.db).Where("menu_id=? AND code=?", menuID, item.Code).Delete(entity.MenuAction{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
	}

	for _, item := range ulist {
		result := entity.GetMenuActionDB(ctx, a.db).Where("menu_id=? AND code=?", menuID, item.Code).Omit("menu_id", "code").Updates(item)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
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
func (a *Menu) updateResources(ctx context.Context, menuID string, items entity.MenuResources) error {
	list, err := a.queryResources(ctx, menuID)
	if err != nil {
		return err
	}

	clist, dlist, ulist := a.compareUpdateResource(list, items)
	for _, item := range clist {
		result := entity.GetMenuResourceDB(ctx, a.db).Create(item)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
	}

	for _, item := range dlist {
		result := entity.GetMenuResourceDB(ctx, a.db).Where("menu_id=? AND code=?", menuID, item.Code).Delete(entity.MenuResource{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
	}

	for _, item := range ulist {
		result := entity.GetMenuResourceDB(ctx, a.db).Where("menu_id=? AND code=?", menuID, item.Code).Omit("menu_id", "code").Updates(item)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Update 更新数据
func (a *Menu) Update(ctx context.Context, recordID string, item schema.Menu) error {
	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaMenu(item)
		result := entity.GetMenuDB(ctx, a.db).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(sitem.ToMenu())
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		err := a.updateActions(ctx, recordID, sitem.ToMenuActions())
		if err != nil {
			return err
		}

		err = a.updateResources(ctx, recordID, sitem.ToMenuResources())
		if err != nil {
			return err
		}

		return nil
	})
}

// UpdateParentPath 更新父级路径
func (a *Menu) UpdateParentPath(ctx context.Context, recordID, parentPath string) error {
	result := entity.GetMenuDB(ctx, a.db).Where("record_id=?", recordID).Update("parent_path", parentPath)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *Menu) Delete(ctx context.Context, recordID string) error {
	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		result := entity.GetMenuDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.Menu{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		result = entity.GetMenuActionDB(ctx, a.db).Where("menu_id=?", recordID).Delete(entity.MenuAction{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		result = entity.GetMenuResourceDB(ctx, a.db).Where("menu_id=?", recordID).Delete(entity.MenuResource{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

func (a *Menu) queryActions(ctx context.Context, menuIDs ...string) (entity.MenuActions, error) {
	var list entity.MenuActions
	result := entity.GetMenuActionDB(ctx, a.db).Where("menu_id IN(?)", menuIDs).Find(&list)
	if err := result.Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return list, nil
}

func (a *Menu) queryResources(ctx context.Context, menuIDs ...string) (entity.MenuResources, error) {
	var list entity.MenuResources
	result := entity.GetMenuResourceDB(ctx, a.db).Where("menu_id IN(?)", menuIDs).Find(&list)
	if err := result.Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return list, nil
}
