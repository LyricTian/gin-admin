package gormrole

import (
	"context"
	"fmt"

	gcontext "github.com/LyricTian/gin-admin/src/context"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/model/gorm/common"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/jinzhu/gorm"
)

// NewModel 实例化角色存储
func NewModel(db *gormplus.DB) *Model {
	db.AutoMigrate(new(Role), new(RoleMenu))
	return &Model{db}
}

// Model 角色存储
type Model struct {
	db *gormplus.DB
}

func (a *Model) getFuncName(name string) string {
	return fmt.Sprintf("gorm.role.%s", name)
}

func (a *Model) getRoleDB(ctx context.Context) *gorm.DB {
	return gormcommon.FromTransDB(ctx, a.db).Model(Role{})
}

func (a *Model) getRoleMenuDB(ctx context.Context) *gorm.DB {
	return gormcommon.FromTransDB(ctx, a.db).Model(RoleMenu{})
}

// Query 查询数据
func (a *Model) Query(ctx context.Context, params schema.RoleQueryParam, pp *schema.PaginationParam) ([]*schema.Role, *schema.PaginationResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := a.getRoleDB(ctx)
	if v := params.Name; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	db = db.Order("id DESC")

	var items []*Role
	pr, err := gormcommon.WrapPageQuery(db, pp, &items)
	if err != nil {
		span.Errorf(err.Error())
		return nil, nil, errors.New("查询数据发生错误")
	}

	sroles := Roles(items).ToSchemaRoles()
	if params.IncludeMenuIDs {
		for i, item := range sroles {
			menuIDs, err := a.QueryMenuIDs(ctx, item.RecordID)
			if err != nil {
				return nil, nil, err
			}
			sroles[i].MenuIDs = menuIDs
		}
	}

	return sroles, pr, nil
}

// Get 查询指定数据
func (a *Model) Get(ctx context.Context, recordID string, includeMenuIDs bool) (*schema.Role, error) {
	span := logger.StartSpan(ctx, "查询指定数据", a.getFuncName("Get"))
	defer span.Finish()

	var item Role
	ok, err := a.db.FindOne(a.getRoleDB(ctx).Where("record_id=?", recordID), &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	srole := item.ToSchemaRole()
	if includeMenuIDs {
		menuIDs, err := a.QueryMenuIDs(ctx, recordID)
		if err != nil {
			return nil, err
		}
		srole.MenuIDs = menuIDs
	}

	return srole, nil
}

// CheckName 检查名称是否存在
func (a *Model) CheckName(ctx context.Context, name string) (bool, error) {
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
func (a *Model) Create(ctx context.Context, item schema.Role) error {
	span := logger.StartSpan(ctx, "创建数据", a.getFuncName("Create"))
	defer span.Finish()

	return gormcommon.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		err := a.CreateRole(ctx, item)
		if err != nil {
			return err
		}

		for _, menuID := range item.MenuIDs {
			err = a.CreateMenu(ctx, RoleMenu{
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
func (a *Model) CreateRole(ctx context.Context, item schema.Role) error {
	span := logger.StartSpan(ctx, "创建角色数据", a.getFuncName("CreateRole"))
	defer span.Finish()

	role := new(Role)
	_ = util.FillStruct(item, role)
	role.Creator, _ = gcontext.FromUserID(ctx)
	result := a.getRoleDB(ctx).Create(role)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("创建角色数据发生错误")
	}
	return nil
}

// Update 更新数据
func (a *Model) Update(ctx context.Context, recordID string, item schema.Role) error {
	span := logger.StartSpan(ctx, "更新数据", a.getFuncName("Update"))
	defer span.Finish()

	return gormcommon.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		err := a.UpdateRole(ctx, recordID, item)
		if err != nil {
			return err
		}

		err = a.DeleteMenu(ctx, recordID)
		if err != nil {
			return err
		}

		for _, menuID := range item.MenuIDs {
			err = a.CreateMenu(ctx, RoleMenu{
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
func (a *Model) UpdateRole(ctx context.Context, recordID string, item schema.Role) error {
	span := logger.StartSpan(ctx, "更新角色数据", a.getFuncName("UpdateRole"))
	defer span.Finish()

	role := new(Role)
	_ = util.FillStruct(item, role)
	result := a.getRoleDB(ctx).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(role)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新角色数据发生错误")
	}
	return nil
}

// Delete 删除数据
func (a *Model) Delete(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除数据", a.getFuncName("Delete"))
	defer span.Finish()

	return gormcommon.ExecTrans(ctx, a.db, func(ctx context.Context) error {
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
func (a *Model) DeleteRole(ctx context.Context, recordID string) error {
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
func (a *Model) UpdateStatus(ctx context.Context, recordID string, status int) error {
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
func (a *Model) QueryMenuIDs(ctx context.Context, recordID string) ([]string, error) {
	span := logger.StartSpan(ctx, "查询菜单ID列表", a.getFuncName("QueryMenuIDs"))
	defer span.Finish()

	var items []*RoleMenu
	result := a.getRoleMenuDB(ctx).Where("role_id=?", recordID).Find(&items)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询菜单ID列表发生错误")
	}

	return RoleMenus(items).ToMenuIDs(), nil
}

// CreateMenu 创建角色菜单
func (a *Model) CreateMenu(ctx context.Context, item RoleMenu) error {
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
func (a *Model) DeleteMenu(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除角色菜单", a.getFuncName("DeleteMenu"))
	defer span.Finish()

	result := a.getRoleMenuDB(ctx).Where("role_id=?", recordID).Delete(RoleMenu{})
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("删除角色菜单发生错误")
	}
	return nil
}
