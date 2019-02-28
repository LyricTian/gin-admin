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

// User 用户存储
type User struct {
	db *gormplus.DB
}

// Init 初始化
func (a *User) Init(db *gormplus.DB) *User {
	db.AutoMigrate(new(entity.User), new(entity.UserRole))
	a.db = db
	return a
}

func (a *User) getFuncName(name string) string {
	return fmt.Sprintf("gorm.model.User.%s", name)
}

// Query 查询数据
func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := entity.GetUserDB(ctx, a.db).DB
	if v := params.UserName; v != "" {
		db = db.Where("user_name LIKE ?", "%"+v+"%")
	}
	if v := params.RealName; v != "" {
		db = db.Where("real_name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	if v := params.RoleIDs; len(v) > 0 {
		subQuery := entity.GetUserRoleDB(ctx, a.db).Select("user_id").Where("role_id IN(?)", v).SubQuery()
		db = db.Where("record_id IN(?)", subQuery)
	}
	db = db.Order("id DESC")

	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	var list entity.Users
	pr, err := WrapPageQuery(db, opt.PageParam, &list)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询数据发生错误")
	}

	qr := &schema.UserQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUsers(),
	}

	for i, item := range qr.Data {
		err := a.fillSchemaUser(ctx, item, opts...)
		if err != nil {
			return nil, err
		}
		qr.Data[i] = item
	}

	return qr, nil
}

func (a *User) fillSchemaUser(ctx context.Context, item *schema.User, opts ...schema.UserQueryOptions) error {
	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.IncludeRoleIDs {
		roleIDs, err := a.QueryRoleIDs(ctx, item.RecordID)
		if err != nil {
			return err
		}
		item.RoleIDs = roleIDs
	}

	return nil
}

// GetByUserName 根据用户名查询指定数据
func (a *User) GetByUserName(ctx context.Context, userName string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	span := logger.StartSpan(ctx, "根据用户名查询指定数据", a.getFuncName("GetByUserName"))
	defer span.Finish()

	var item entity.User
	ok, err := a.db.FindOne(entity.GetUserDB(ctx, a.db).Where("user_name=?", userName), &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("根据用户名查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	sitem := item.ToSchemaUser()
	err = a.fillSchemaUser(ctx, sitem, opts...)
	if err != nil {
		return nil, err
	}

	return sitem, nil
}

// Get 查询指定数据
func (a *User) Get(ctx context.Context, recordID string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	span := logger.StartSpan(ctx, "查询指定数据", a.getFuncName("Get"))
	defer span.Finish()

	var item entity.User
	ok, err := a.db.FindOne(entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID), &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	sitem := item.ToSchemaUser()
	err = a.fillSchemaUser(ctx, sitem, opts...)
	if err != nil {
		return nil, err
	}

	return sitem, nil
}

// CheckUserName 检查用户名是否存在
func (a *User) CheckUserName(ctx context.Context, userName string) (bool, error) {
	span := logger.StartSpan(ctx, "检查用户名是否存在", a.getFuncName("CheckUserName"))
	defer span.Finish()

	db := entity.GetUserDB(ctx, a.db).Where("user_name=?", userName)
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
		sitem := entity.SchemaUser(item)
		result := entity.GetUserDB(ctx, a.db).Create(sitem.ToUser())
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("创建用户数据发生错误")
		}

		for _, eitem := range sitem.ToUserRoles() {
			result := entity.GetUserRoleDB(ctx, a.db).Create(eitem)
			if err := result.Error; err != nil {
				span.Errorf(err.Error())
				return errors.New("创建用户角色发生错误")
			}
		}
		return nil
	})
}

// Update 更新数据
func (a *User) Update(ctx context.Context, recordID string, item schema.User) error {
	span := logger.StartSpan(ctx, "更新数据", a.getFuncName("Update"))
	defer span.Finish()

	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaUser(item)
		omits := []string{"record_id", "creator"}
		if sitem.Password == "" {
			omits = append(omits, "password")
		}

		result := entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID).Omit(omits...).Updates(sitem.ToUser())
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("更新用户数据发生错误")
		}

		result = entity.GetUserRoleDB(ctx, a.db).Where("user_id=?", recordID).Delete(entity.UserRole{})
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("删除用户角色发生错误")
		}

		for _, eitem := range sitem.ToUserRoles() {
			result := entity.GetUserRoleDB(ctx, a.db).Create(eitem)
			if err := result.Error; err != nil {
				span.Errorf(err.Error())
				return errors.New("创建用户角色发生错误")
			}
		}
		return nil
	})
}

// Delete 删除数据
func (a *User) Delete(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除数据", a.getFuncName("Delete"))
	defer span.Finish()

	return ExecTrans(ctx, a.db, func(ctx context.Context) error {
		result := entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.User{})
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("删除用户数据发生错误")
		}

		result = entity.GetUserRoleDB(ctx, a.db).Where("user_id=?", recordID).Delete(entity.UserRole{})
		if err := result.Error; err != nil {
			span.Errorf(err.Error())
			return errors.New("删除用户角色发生错误")
		}

		return nil
	})
}

// UpdateStatus 更新状态
func (a *User) UpdateStatus(ctx context.Context, recordID string, status int) error {
	span := logger.StartSpan(ctx, "更新状态", a.getFuncName("UpdateStatus"))
	defer span.Finish()

	result := entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID).Update("status", status)
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

	result := entity.GetUserDB(ctx, a.db).Where("record_id=?", recordID).Update("password", password)
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

	var list entity.UserRoles
	result := entity.GetUserRoleDB(ctx, a.db).Where("user_id=?", recordID).Find(&list)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询角色ID列表发生错误")
	}

	return list.ToRoleIDs(), nil
}
