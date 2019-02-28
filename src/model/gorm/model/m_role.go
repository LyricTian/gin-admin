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

// Role 角色存储
type Role struct {
	db *gormplus.DB
}

// Init 初始化
func (a *Role) Init(db *gormplus.DB) *Role {
	db.AutoMigrate(new(entity.Role), new(entity.RoleMenu))
	a.db = db
	return a
}

func (a *Role) getFuncName(name string) string {
	return fmt.Sprintf("gorm.model.Role.%s", name)
}

// Query 查询数据
func (a *Role) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := entity.GetRoleDB(ctx, a.db).DB
	if v := params.Name; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.RecordIDs; len(v) > 0 {
		db = db.Where("record_id IN(?)", v)
	}
	db = db.Order("sequence DESC,id DESC")

	var opt schema.RoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

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

	for i, item := range qr.Data {
		err := a.fillSchameRole(ctx, item, opts...)
		if err != nil {
			return nil, err
		}
		qr.Data[i] = item
	}

	return qr, nil
}

// 填充角色对象
func (a *Role) fillSchameRole(ctx context.Context, item *schema.Role, opts ...schema.RoleQueryOptions) error {
	var opt schema.RoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.IncludeMenuIDs {
		menuIDs, err := a.QueryMenuIDs(ctx, item.RecordID)
		if err != nil {
			return err
		}
		item.MenuIDs = menuIDs
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
	err = a.fillSchameRole(ctx, sitem, opts...)
	if err != nil {
		return nil, err
	}

	return sitem, nil
}

// CheckName 检查名称是否存在
func (a *Role) CheckName(ctx context.Context, name string) (bool, error) {
	span := logger.StartSpan(ctx, "检查名称是否存在", a.getFuncName("CheckName"))
	defer span.Finish()

	db := entity.GetRoleDB(ctx, a.db).Where("name=?", name)
	exists, err := a.db.Check(db)
	if err != nil {
		span.Errorf(err.Error())
		return false, errors.New("检查名称是否存在发生错误")
	}
	return exists, nil
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

		for _, eitem := range sitem.ToRoleMenus() {
			result := entity.GetRoleMenuDB(ctx, a.db).Create(eitem)
			if err := result.Error; err != nil {
				span.Errorf(err.Error())
				return errors.New("创建角色菜单数据发生错误")
			}
		}

		return nil
	})
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

		result = entity.GetRoleMenuDB(ctx, a.db).Where("role_id=?", recordID).Delete(entity.RoleMenu{})
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("删除角色菜单发生错误")
		}

		for _, eitem := range sitem.ToRoleMenus() {
			result := entity.GetRoleMenuDB(ctx, a.db).Create(eitem)
			if err := result.Error; err != nil {
				span.Errorf(err.Error())
				return errors.New("创建角色菜单数据发生错误")
			}
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

// QueryMenuIDs 查询角色菜单ID列表
func (a *Role) QueryMenuIDs(ctx context.Context, roleID string) ([]string, error) {
	span := logger.StartSpan(ctx, "查询角色菜单ID列表", a.getFuncName("QueryMenuIDs"))
	defer span.Finish()

	var list entity.RoleMenus
	result := entity.GetRoleMenuDB(ctx, a.db).Where("role_id=?", roleID).Find(&list)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询角色菜单ID列表发生错误")
	}

	return list.ToMenuIDs(), nil
}
