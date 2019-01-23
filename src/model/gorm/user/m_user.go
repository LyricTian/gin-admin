package gormuser

import (
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/model/gorm/common"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/jinzhu/gorm"
)

// NewModel 实例化用户存储
func NewModel(db *gormplus.DB) *Model {
	db.AutoMigrate(new(User), new(UserRole))
	return &Model{db}
}

// Model 用户存储
type Model struct {
	db *gormplus.DB
}

func (a *Model) getFuncName(name string) string {
	return fmt.Sprintf("gorm.user.%s", name)
}

func (a *Model) getUserDB(ctx context.Context) *gorm.DB {
	return gormcommon.FromTransDB(ctx, a.db).Model(User{})
}

func (a *Model) getUserRoleDB(ctx context.Context) *gorm.DB {
	return gormcommon.FromTransDB(ctx, a.db).Model(UserRole{})
}

// Query 查询数据
func (a *Model) Query(ctx context.Context, params schema.UserQueryParam, pp *schema.PaginationParam) ([]*schema.User, *schema.PaginationResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := a.getUserDB(ctx)
	if v := params.UserName; v != "" {
		db = db.Where("user_name LIKE ?", "%"+v+"%")
	}
	if v := params.RealName; v != "" {
		db = db.Where("real_name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	if v := params.RoleID; v != "" {
		expr := a.getUserRoleDB(ctx).Select("user_id").Where("role_id=?", v).SubQuery()
		db = db.Where("record_id IN(?)", expr)
	}
	db = db.Order("id DESC")

	var items []*User
	pr, err := gormcommon.WrapPageQuery(db, pp, &items)
	if err != nil {
		span.Errorf(err.Error())
		return nil, nil, errors.New("查询数据发生错误")
	}

	sitems := Users(items).ToSchemaUsers(params.IncludePassword)
	if params.IncludeRoleIDs {
		for i, item := range sitems {
			roleIDs, err := a.QueryRoleIDs(ctx, item.RecordID)
			if err != nil {
				return nil, nil, err
			}
			sitems[i].RoleIDs = roleIDs
		}
	}

	return sitems, pr, nil
}

// GetByUserName 根据用户名查询指定数据
func (a *Model) GetByUserName(ctx context.Context, userName string, includePassword, includeRoleIDs bool) (*schema.User, error) {
	span := logger.StartSpan(ctx, "根据用户名查询指定数据", a.getFuncName("GetByUserName"))
	defer span.Finish()

	var item User
	ok, err := a.db.FindOne(a.getUserDB(ctx).Where("user_name=?", userName), &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("根据用户名查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	sitem := item.ToSchemaUser(includePassword)
	if includeRoleIDs {
		roleIDs, err := a.QueryRoleIDs(ctx, item.RecordID)
		if err != nil {
			return nil, err
		}
		sitem.RoleIDs = roleIDs
	}

	return sitem, nil
}

// Get 查询指定数据
func (a *Model) Get(ctx context.Context, recordID string, includePassword, includeRoleIDs bool) (*schema.User, error) {
	span := logger.StartSpan(ctx, "查询指定数据", a.getFuncName("Get"))
	defer span.Finish()

	var item User
	ok, err := a.db.FindOne(a.getUserDB(ctx).Where("record_id=?", recordID), &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	sitem := item.ToSchemaUser(includePassword)
	if includeRoleIDs {
		roleIDs, err := a.QueryRoleIDs(ctx, recordID)
		if err != nil {
			return nil, err
		}
		sitem.RoleIDs = roleIDs
	}

	return sitem, nil
}

// CheckUserName 检查用户名是否存在
func (a *Model) CheckUserName(ctx context.Context, userName string) (bool, error) {
	span := logger.StartSpan(ctx, "检查用户名是否存在", a.getFuncName("CheckUserName"))
	defer span.Finish()

	db := a.getUserDB(ctx).Where("user_name=?", userName)
	exists, err := a.db.Check(db)
	if err != nil {
		span.Errorf(err.Error())
		return false, errors.New("检查用户名是否存在发生错误")
	}
	return exists, nil
}

// Create 创建数据
func (a *Model) Create(ctx context.Context, item schema.User) error {
	span := logger.StartSpan(ctx, "创建数据", a.getFuncName("Create"))
	defer span.Finish()

	return gormcommon.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		err := a.CreateUser(ctx, item)
		if err != nil {
			return err
		}

		for _, roleID := range item.RoleIDs {
			err = a.CreateRole(ctx, UserRole{
				UserID: item.RecordID,
				RoleID: roleID,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// CreateUser 创建用户数据
func (a *Model) CreateUser(ctx context.Context, item schema.User) error {
	span := logger.StartSpan(ctx, "创建用户数据", a.getFuncName("CreateUser"))
	defer span.Finish()

	user := SchemaUser(item).ToUser()
	user.Creator = gormcommon.FromUserID(ctx)
	result := a.getUserDB(ctx).Create(user)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("创建用户数据发生错误")
	}
	return nil
}

// Update 更新数据
func (a *Model) Update(ctx context.Context, recordID string, item schema.User) error {
	span := logger.StartSpan(ctx, "更新数据", a.getFuncName("Update"))
	defer span.Finish()

	return gormcommon.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		err := a.UpdateUser(ctx, recordID, item)
		if err != nil {
			return err
		}

		err = a.DeleteRole(ctx, recordID)
		if err != nil {
			return err
		}

		for _, roleID := range item.RoleIDs {
			err = a.CreateRole(ctx, UserRole{
				UserID: recordID,
				RoleID: roleID,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateUser 更新用户数据
func (a *Model) UpdateUser(ctx context.Context, recordID string, item schema.User) error {
	span := logger.StartSpan(ctx, "更新用户数据", a.getFuncName("UpdateUser"))
	defer span.Finish()

	user := SchemaUser(item).ToUser()
	omits := []string{"record_id", "creator"}
	if user.Password == "" {
		omits = append(omits, "password")
	}
	result := a.getUserDB(ctx).Where("record_id=?", recordID).Omit(omits...).Updates(user)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新用户数据发生错误")
	}
	return nil
}

// Delete 删除数据
func (a *Model) Delete(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除数据", a.getFuncName("Delete"))
	defer span.Finish()

	return gormcommon.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		err := a.DeleteUser(ctx, recordID)
		if err != nil {
			return err
		}

		err = a.DeleteRole(ctx, recordID)
		if err != nil {
			return err
		}

		return nil
	})
}

// DeleteUser 删除用户数据
func (a *Model) DeleteUser(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除用户数据", a.getFuncName("DeleteUser"))
	defer span.Finish()

	result := a.getUserDB(ctx).Where("record_id=?", recordID).Delete(User{})
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("删除用户数据发生错误")
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Model) UpdateStatus(ctx context.Context, recordID string, status int) error {
	span := logger.StartSpan(ctx, "更新状态", a.getFuncName("UpdateStatus"))
	defer span.Finish()

	result := a.getUserDB(ctx).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新状态发生错误")
	}
	return nil
}

// QueryRoleIDs 查询角色ID列表
func (a *Model) QueryRoleIDs(ctx context.Context, recordID string) ([]string, error) {
	span := logger.StartSpan(ctx, "查询角色ID列表", a.getFuncName("QueryRoleIDs"))
	defer span.Finish()

	var items []*UserRole
	result := a.getUserRoleDB(ctx).Where("user_id=?", recordID).Find(&items)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询角色ID列表发生错误")
	}

	return UserRoles(items).ToRoleIDs(), nil
}

// CreateRole 创建用户角色
func (a *Model) CreateRole(ctx context.Context, item UserRole) error {
	span := logger.StartSpan(ctx, "创建用户角色", a.getFuncName("CreateRole"))
	defer span.Finish()

	result := a.getUserRoleDB(ctx).Create(&item)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("创建用户角色发生错误")
	}
	return nil
}

// DeleteRole 删除用户角色
func (a *Model) DeleteRole(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除用户角色", a.getFuncName("DeleteRole"))
	defer span.Finish()

	result := a.getUserRoleDB(ctx).Where("user_id=?", recordID).Delete(UserRole{})
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("删除用户角色发生错误")
	}
	return nil
}
