package gormmodel

import (
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/model/gorm/entity"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
)

// InitRole 初始化角色存储
func InitRole(db *gormplus.DB) *Role {
	db.AutoMigrate(new(gormentity.Role), new(gormentity.RoleMenu))
	return NewRole(db)
}

// NewRole 实例化角色存储
func NewRole(db *gormplus.DB) *Role {
	return &Role{db: db}
}

// Role 角色存储
type Role struct {
	db *gormplus.DB
}

func (a *Role) getFuncName(name string) string {
	return fmt.Sprintf("gorm.role.%s", name)
}

func (a *Role) getRoleDB(ctx context.Context) *gormplus.DB {
	return FromTransDBWithModel(ctx, a.db, gormentity.Role{})
}

func (a *Role) getRoleMenuDB(ctx context.Context) *gormplus.DB {
	return FromTransDBWithModel(ctx, a.db, gormentity.RoleMenu{})
}

func (a *Role) getQueryOption(opts ...schema.RoleQueryOptions) schema.RoleQueryOptions {
	if len(opts) > 0 {
		return opts[0]
	}
	return schema.RoleQueryOptions{}
}

// Query 查询数据
func (a *Role) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (schema.RoleQueryResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := a.getRoleDB(ctx).DB
	if v := params.Name; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	if v := params.RecordIDs; len(v) > 0 {
		db = db.Where("record_id IN(?)", v)
	}
	db = db.Order("id DESC")

	opt := a.getQueryOption(opts...)
	var qr schema.RoleQueryResult
	var items gormentity.Roles
	pr, err := WrapPageQuery(db, opt.PageParam, &items)
	if err != nil {
		span.Errorf(err.Error())
		return qr, errors.New("查询数据发生错误")
	}
	qr.PageResult = pr

	sitems := make([]*schema.Role, len(items))
	for i, item := range items {
		sitems[i], err = a.toSchemaRole(ctx, *item, opts...)
		if err != nil {
			return qr, err
		}
	}
	qr.Data = sitems

	return qr, nil
}

func (a *Role) toSchemaRole(ctx context.Context, item gormentity.Role, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	opt := a.getQueryOption(opts...)
	sitem := item.ToSchemaRole()
	if opt.IncludeMenuIDs {
		menuIDs, err := a.QueryMenuIDs(ctx, item.RecordID)
		if err != nil {
			return nil, err
		}
		sitem.MenuIDs = menuIDs
	}

	return sitem, nil
}

// Get 查询指定数据
func (a *Role) Get(ctx context.Context, recordID string, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	span := logger.StartSpan(ctx, "查询指定数据", a.getFuncName("Get"))
	defer span.Finish()

	var item gormentity.Role
	ok, err := a.db.FindOne(a.getRoleDB(ctx).Where("record_id=?", recordID), &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	return a.toSchemaRole(ctx, item, opts...)
}

// CheckName 检查名称是否存在
func (a *Role) CheckName(ctx context.Context, name string) (bool, error) {
	span := logger.StartSpan(ctx, "检查名称是否存在", a.getFuncName("CheckName"))
	defer span.Finish()

	db := a.getRoleDB(ctx).Where("name=?", name)
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
		err := a.CreateRole(ctx, item)
		if err != nil {
			return err
		}

		for _, menuID := range item.MenuIDs {
			err = a.CreateMenu(ctx, gormentity.RoleMenu{
				RoleID: item.RecordID,
				MenuID: menuID,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// CreateRole 创建角色数据
func (a *Role) CreateRole(ctx context.Context, item schema.Role) error {
	span := logger.StartSpan(ctx, "创建角色数据", a.getFuncName("CreateRole"))
	defer span.Finish()

	role := gormentity.SchemaRole(item).ToRole()
	role.Creator = FromUserID(ctx)
	result := a.getRoleDB(ctx).Create(role)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("创建角色数据发生错误")
	}
	return nil
}

// Update 更新数据
func (a *Role) Update(ctx context.Context, recordID string, item schema.Role) error {
	span := logger.StartSpan(ctx, "更新数据", a.getFuncName("Update"))
	defer span.Finish()

	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		err := a.UpdateRole(ctx, recordID, item)
		if err != nil {
			return err
		}

		err = a.DeleteMenu(ctx, recordID)
		if err != nil {
			return err
		}

		for _, menuID := range item.MenuIDs {
			err = a.CreateMenu(ctx, gormentity.RoleMenu{
				RoleID: recordID,
				MenuID: menuID,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateRole 更新角色数据
func (a *Role) UpdateRole(ctx context.Context, recordID string, item schema.Role) error {
	span := logger.StartSpan(ctx, "更新角色数据", a.getFuncName("UpdateRole"))
	defer span.Finish()

	role := gormentity.SchemaRole(item).ToRole()
	result := a.getRoleDB(ctx).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(role)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新角色数据发生错误")
	}
	return nil
}

// Delete 删除数据
func (a *Role) Delete(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除数据", a.getFuncName("Delete"))
	defer span.Finish()

	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		err := a.DeleteRole(ctx, recordID)
		if err != nil {
			return err
		}

		err = a.DeleteMenu(ctx, recordID)
		if err != nil {
			return err
		}

		return nil
	})
}

// DeleteRole 删除角色数据
func (a *Role) DeleteRole(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除角色数据", a.getFuncName("DeleteRole"))
	defer span.Finish()

	result := a.getRoleDB(ctx).Where("record_id=?", recordID).Delete(Role{})
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("删除角色数据发生错误")
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Role) UpdateStatus(ctx context.Context, recordID string, status int) error {
	span := logger.StartSpan(ctx, "更新状态", a.getFuncName("UpdateStatus"))
	defer span.Finish()

	result := a.getRoleDB(ctx).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新状态发生错误")
	}
	return nil
}

// QueryMenuIDs 查询菜单ID列表
func (a *Role) QueryMenuIDs(ctx context.Context, recordID string) ([]string, error) {
	span := logger.StartSpan(ctx, "查询菜单ID列表", a.getFuncName("QueryMenuIDs"))
	defer span.Finish()

	var items []*gormentity.RoleMenu
	result := a.getRoleMenuDB(ctx).Where("role_id=?", recordID).Find(&items)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询菜单ID列表发生错误")
	}

	return gormentity.RoleMenus(items).ToMenuIDs(), nil
}

// CreateMenu 创建角色菜单
func (a *Role) CreateMenu(ctx context.Context, item gormentity.RoleMenu) error {
	span := logger.StartSpan(ctx, "创建角色菜单", a.getFuncName("CreateMenu"))
	defer span.Finish()

	result := a.getRoleMenuDB(ctx).Create(&item)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("创建角色菜单发生错误")
	}
	return nil
}

// DeleteMenu 删除角色菜单
func (a *Role) DeleteMenu(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除角色菜单", a.getFuncName("DeleteMenu"))
	defer span.Finish()

	result := a.getRoleMenuDB(ctx).Where("role_id=?", recordID).Delete(gormentity.RoleMenu{})
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("删除角色菜单发生错误")
	}
	return nil
}
