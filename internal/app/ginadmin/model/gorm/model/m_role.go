package model

import (
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin/internal/app/ginadmin/model/gorm/entity"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/LyricTian/gin-admin/pkg/gormplus"
	"github.com/LyricTian/gin-admin/pkg/logger"
)

// NewRole 创建角色存储实例
func NewRole(db *gormplus.DB) *Role {
	return &Role{db}
}

// Role 角色存储
type Role struct {
	db *gormplus.DB
}

func (a *Role) getFuncName(name string) string {
	return fmt.Sprintf("gorm.model.Role.%s", name)
}

func (a *Role) getQueryOption(opts ...schema.RoleQueryOptions) schema.RoleQueryOptions {
	var opt schema.RoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *Role) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := entity.GetRoleDB(ctx, a.db).DB
	if v := params.RecordIDs; len(v) > 0 {
		db = db.Where("record_id IN(?)", v)
	}
	if v := params.Name; v != "" {
		db = db.Where("name=?", v)
	}
	if v := params.LikeName; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.UserID; v != "" {
		subQuery := entity.GetUserRoleDB(ctx, a.db).Where("user_id=?", v).Select("role_id").SubQuery()
		db = db.Where("record_id IN(?)", subQuery)
	}
	db = db.Order("sequence DESC,id DESC")

	opt := a.getQueryOption(opts...)
	var list entity.Roles
	pr, err := WrapPageQuery(db, opt.PageParam, &list)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询数据发生错误")
	}
	qr := &schema.RoleQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRoles(),
	}

	err = a.fillSchameRoles(ctx, qr.Data, opts...)
	if err != nil {
		return nil, err
	}

	return qr, nil
}

// 填充角色对象
func (a *Role) fillSchameRoles(ctx context.Context, items []*schema.Role, opts ...schema.RoleQueryOptions) error {
	opt := a.getQueryOption(opts...)

	if opt.IncludeMenus {

		roleIDs := make([]string, len(items))
		for i, item := range items {
			roleIDs[i] = item.RecordID
		}

		var menuList entity.RoleMenus
		if opt.IncludeMenus {
			items, err := a.queryMenus(ctx, roleIDs...)
			if err != nil {
				return err
			}
			menuList = items
		}

		for i, item := range items {
			if len(menuList) > 0 {
				items[i].Menus = menuList.GetByRoleID(item.RecordID)
			}
		}
	}
	return nil
}

// Get 查询指定数据
func (a *Role) Get(ctx context.Context, recordID string, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	span := logger.StartSpan(ctx, "查询指定数据", a.getFuncName("Get"))
	defer span.Finish()

	var role entity.Role
	ok, err := a.db.FindOne(entity.GetRoleDB(ctx, a.db).Where("record_id=?", recordID), &role)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	sitem := role.ToSchemaRole()
	err = a.fillSchameRoles(ctx, []*schema.Role{sitem}, opts...)
	if err != nil {
		return nil, err
	}

	return sitem, nil
}

// Create 创建数据
func (a *Role) Create(ctx context.Context, item schema.Role) error {
	span := logger.StartSpan(ctx, "创建数据", a.getFuncName("Create"))
	defer span.Finish()

	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaRole(item)
		result := entity.GetRoleDB(ctx, a.db).Create(sitem.ToRole())
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("创建角色数据发生错误")
		}

		for _, item := range sitem.ToRoleMenus() {
			result := entity.GetRoleMenuDB(ctx, a.db).Create(item)
			if err := result.Error; err != nil {
				span.Errorf(err.Error())
				return errors.New("创建角色菜单数据发生错误")
			}
		}

		return nil
	})
}

// 对比并获取需要新增，修改，删除的菜单项
func (a *Role) compareUpdateMenu(oldList, newList entity.RoleMenus) (clist, dlist, ulist entity.RoleMenus) {
	oldMap, newMap := oldList.ToMap(), newList.ToMap()

	for _, nitem := range newList {
		if _, ok := oldMap[nitem.MenuID]; ok {
			ulist = append(ulist, nitem)
			continue
		}
		clist = append(clist, nitem)
	}

	for _, oitem := range oldList {
		if _, ok := newMap[oitem.MenuID]; !ok {
			dlist = append(dlist, oitem)
		}
	}
	return
}

// 更新菜单数据
func (a *Role) updateMenus(ctx context.Context, span *logger.Entry, roleID string, items entity.RoleMenus) error {
	list, err := a.queryMenus(ctx, roleID)
	if err != nil {
		return err
	}

	clist, dlist, ulist := a.compareUpdateMenu(list, items)
	for _, item := range clist {
		result := entity.GetRoleMenuDB(ctx, a.db).Create(item)
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("创建角色菜单数据发生错误")
		}
	}

	for _, item := range dlist {
		result := entity.GetRoleMenuDB(ctx, a.db).Where("role_id=? AND menu_id=?", roleID, item.MenuID).Delete(entity.RoleMenu{})
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("删除角色菜单数据发生错误")
		}
	}

	for _, item := range ulist {
		result := entity.GetRoleMenuDB(ctx, a.db).Where("role_id=? AND menu_id=?", roleID, item.MenuID).Omit("role_id", "menu_id").Updates(item)
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("更新角色菜单数据发生错误")
		}
	}
	return nil
}

// Update 更新数据
func (a *Role) Update(ctx context.Context, recordID string, item schema.Role) error {
	span := logger.StartSpan(ctx, "更新数据", a.getFuncName("Update"))
	defer span.Finish()

	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaRole(item)
		result := entity.GetRoleDB(ctx, a.db).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(sitem.ToRole())
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("更新角色数据发生错误")
		}

		err := a.updateMenus(ctx, span, recordID, sitem.ToRoleMenus())
		if err != nil {
			return err
		}

		return nil
	})
}

// Delete 删除数据
func (a *Role) Delete(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除数据", a.getFuncName("Delete"))
	defer span.Finish()

	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		result := entity.GetRoleDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.Role{})
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("删除角色数据发生错误")
		}

		result = entity.GetRoleMenuDB(ctx, a.db).Where("role_id=?", recordID).Delete(entity.RoleMenu{})
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("删除角色菜单数据发生错误")
		}

		return nil
	})
}

func (a *Role) queryMenus(ctx context.Context, roleIDs ...string) (entity.RoleMenus, error) {
	span := logger.StartSpan(ctx, "查询角色菜单数据", a.getFuncName("queryMenus"))
	defer span.Finish()

	var list entity.RoleMenus
	result := entity.GetRoleMenuDB(ctx, a.db).Where("role_id IN(?)", roleIDs).Find(&list)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询角色菜单数据发生错误")
	}

	return list, nil
}
