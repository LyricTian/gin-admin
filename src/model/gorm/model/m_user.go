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

// InitUser 初始化用户存储
func InitUser(db *gormplus.DB) *User {
	db.AutoMigrate(new(entity.User), new(entity.UserRole))
	return NewUser(db)
}

// NewUser 实例化用户存储
func NewUser(db *gormplus.DB) *User {
	return &User{db: db}
}

// User 用户存储
type User struct {
	db *gormplus.DB
}

func (a *User) getFuncName(name string) string {
	return fmt.Sprintf("gorm.model.User.%s", name)
}

func (a *User) getUserDB(ctx context.Context) *gormplus.DB {
	return FromDBWithModel(ctx, a.db, entity.User{})
}

func (a *User) getUserRoleDB(ctx context.Context) *gormplus.DB {
	return FromDBWithModel(ctx, a.db, entity.UserRole{})
}

func (a *User) getQueryOption(opts ...schema.UserQueryOptions) schema.UserQueryOptions {
	if len(opts) > 0 {
		return opts[0]
	}
	return schema.UserQueryOptions{}
}

// Query 查询数据
func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (schema.UserQueryResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := a.getUserDB(ctx).DB
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

	opt := a.getQueryOption(opts...)
	var qr schema.UserQueryResult
	var items []*entity.User
	pr, err := WrapPageQuery(db, opt.PageParam, &items)
	if err != nil {
		span.Errorf(err.Error())
		return qr, errors.New("查询数据发生错误")
	}
	qr.PageResult = pr

	sitems := make([]*schema.User, len(items))
	for i, item := range items {
		sitems[i], err = a.toSchemaUser(ctx, *item, opts...)
		if err != nil {
			return qr, err
		}
	}
	qr.Data = sitems

	return qr, nil
}

func (a *User) toSchemaUser(ctx context.Context, item entity.User, opts ...schema.UserQueryOptions) (*schema.User, error) {
	opt := a.getQueryOption(opts...)
	sitem := item.ToSchemaUser(opt.IncludePassword)
	if opt.IncludeRoleIDs {
		roleIDs, err := a.QueryRoleIDs(ctx, item.RecordID)
		if err != nil {
			return nil, err
		}
		sitem.RoleIDs = roleIDs
	}

	return sitem, nil
}

// GetByUserName 根据用户名查询指定数据
func (a *User) GetByUserName(ctx context.Context, userName string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	span := logger.StartSpan(ctx, "根据用户名查询指定数据", a.getFuncName("GetByUserName"))
	defer span.Finish()

	var item entity.User
	ok, err := a.db.FindOne(a.getUserDB(ctx).Where("user_name=?", userName), &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("根据用户名查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	return a.toSchemaUser(ctx, item, opts...)
}

// Get 查询指定数据
func (a *User) Get(ctx context.Context, recordID string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	span := logger.StartSpan(ctx, "查询指定数据", a.getFuncName("Get"))
	defer span.Finish()

	var item entity.User
	ok, err := a.db.FindOne(a.getUserDB(ctx).Where("record_id=?", recordID), &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	return a.toSchemaUser(ctx, item, opts...)
}

// CheckUserName 检查用户名是否存在
func (a *User) CheckUserName(ctx context.Context, userName string) (bool, error) {
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
func (a *User) Create(ctx context.Context, item schema.User) error {
	span := logger.StartSpan(ctx, "创建数据", a.getFuncName("Create"))
	defer span.Finish()

	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		err := a.CreateUser(ctx, item)
		if err != nil {
			return err
		}

		for _, roleID := range item.RoleIDs {
			err = a.CreateRole(ctx, entity.UserRole{
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
func (a *User) CreateUser(ctx context.Context, item schema.User) error {
	span := logger.StartSpan(ctx, "创建用户数据", a.getFuncName("CreateUser"))
	defer span.Finish()

	user := entity.SchemaUser(item).ToUser()
	user.Creator = FromUserID(ctx)
	result := a.getUserDB(ctx).Create(user)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("创建用户数据发生错误")
	}
	return nil
}

// Update 更新数据
func (a *User) Update(ctx context.Context, recordID string, item schema.User) error {
	span := logger.StartSpan(ctx, "更新数据", a.getFuncName("Update"))
	defer span.Finish()

	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		err := a.UpdateUser(ctx, recordID, item)
		if err != nil {
			return err
		}

		err = a.DeleteRole(ctx, recordID)
		if err != nil {
			return err
		}

		for _, roleID := range item.RoleIDs {
			err = a.CreateRole(ctx, entity.UserRole{
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
func (a *User) UpdateUser(ctx context.Context, recordID string, item schema.User) error {
	span := logger.StartSpan(ctx, "更新用户数据", a.getFuncName("UpdateUser"))
	defer span.Finish()

	user := entity.SchemaUser(item).ToUser()
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
func (a *User) Delete(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除数据", a.getFuncName("Delete"))
	defer span.Finish()

	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
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
func (a *User) DeleteUser(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除用户数据", a.getFuncName("DeleteUser"))
	defer span.Finish()

	result := a.getUserDB(ctx).Where("record_id=?", recordID).Delete(entity.User{})
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("删除用户数据发生错误")
	}
	return nil
}

// UpdateStatus 更新状态
func (a *User) UpdateStatus(ctx context.Context, recordID string, status int) error {
	span := logger.StartSpan(ctx, "更新状态", a.getFuncName("UpdateStatus"))
	defer span.Finish()

	result := a.getUserDB(ctx).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新状态发生错误")
	}
	return nil
}

// UpdatePassword 更新密码
func (a *User) UpdatePassword(ctx context.Context, recordID, password string) error {
	span := logger.StartSpan(ctx, "更新密码", a.getFuncName("UpdatePassword"))
	defer span.Finish()

	result := a.getUserDB(ctx).Where("record_id=?", recordID).Update("password", password)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新密码发生错误")
	}
	return nil
}

// QueryRoleIDs 查询角色ID列表
func (a *User) QueryRoleIDs(ctx context.Context, recordID string) ([]string, error) {
	span := logger.StartSpan(ctx, "查询角色ID列表", a.getFuncName("QueryRoleIDs"))
	defer span.Finish()

	var items entity.UserRoles
	result := a.getUserRoleDB(ctx).Where("user_id=?", recordID).Find(&items)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询角色ID列表发生错误")
	}

	return items.ToRoleIDs(), nil
}

// CreateRole 创建用户角色
func (a *User) CreateRole(ctx context.Context, item entity.UserRole) error {
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
func (a *User) DeleteRole(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除用户角色", a.getFuncName("DeleteRole"))
	defer span.Finish()

	result := a.getUserRoleDB(ctx).Where("user_id=?", recordID).Delete(entity.UserRole{})
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("删除用户角色发生错误")
	}
	return nil
}
