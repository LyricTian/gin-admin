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

// NewUser 创建用户存储实例
func NewUser(db *gormplus.DB) *User {
	return &User{db}
}

// User 用户存储
type User struct {
	db *gormplus.DB
}

func (a *User) getFuncName(name string) string {
	return fmt.Sprintf("gorm.model.User.%s", name)
}

func (a *User) getQueryOption(opts ...schema.UserQueryOptions) schema.UserQueryOptions {
	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := entity.GetUserDB(ctx, a.db).DB
	if v := params.UserName; v != "" {
		db = db.Where("user_name=?", v)
	}
	if v := params.LikeUserName; v != "" {
		db = db.Where("user_name LIKE ?", "%"+v+"%")
	}
	if v := params.LikeRealName; v != "" {
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

	opt := a.getQueryOption(opts...)
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

	err = a.fillSchemaUsers(ctx, qr.Data, opts...)
	if err != nil {
		return nil, err
	}

	return qr, nil
}

func (a *User) fillSchemaUsers(ctx context.Context, items []*schema.User, opts ...schema.UserQueryOptions) error {
	opt := a.getQueryOption(opts...)

	if opt.IncludeRoles {
		userIDs := make([]string, len(items))
		for i, item := range items {
			userIDs[i] = item.RecordID
		}

		var roleList entity.UserRoles
		if opt.IncludeRoles {
			items, err := a.queryRoles(ctx, userIDs...)
			if err != nil {
				return err
			}
			roleList = items
		}

		for i, item := range items {
			if len(roleList) > 0 {
				items[i].Roles = roleList.GetByUserID(item.RecordID)
			}
		}
	}

	return nil
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
	err = a.fillSchemaUsers(ctx, []*schema.User{sitem}, opts...)
	if err != nil {
		return nil, err
	}

	return sitem, nil
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

// 对比并获取需要新增，修改，删除的角色数据
func (a *User) compareUpdateRole(oldList, newList []*entity.UserRole) (clist, dlist, ulist []*entity.UserRole) {
	for _, nitem := range newList {
		exists := false
		for _, oitem := range oldList {
			if oitem.RoleID == nitem.RoleID {
				exists = true
				ulist = append(ulist, nitem)
				break
			}
		}
		if !exists {
			clist = append(clist, nitem)
		}
	}

	for _, oitem := range oldList {
		exists := false
		for _, nitem := range newList {
			if nitem.RoleID == oitem.RoleID {
				exists = true
				break
			}
		}
		if !exists {
			dlist = append(dlist, oitem)
		}
	}

	return
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

		roles, err := a.queryRoles(ctx, recordID)
		if err != nil {
			return err
		}

		clist, dlist, ulist := a.compareUpdateRole(roles, sitem.ToUserRoles())
		for _, item := range clist {
			result := entity.GetUserRoleDB(ctx, a.db).Create(item)
			if err := result.Error; err != nil {
				span.Errorf(err.Error())
				return errors.New("创建用户角色数据发生错误")
			}
		}

		for _, item := range dlist {
			result := entity.GetUserRoleDB(ctx, a.db).Where("user_id=? AND role_id=?", recordID, item.RoleID).Delete(entity.UserRole{})
			if err := result.Error; err != nil {
				span.Errorf(err.Error())
				return errors.New("删除用户角色数据发生错误")
			}
		}

		for _, item := range ulist {
			result := entity.GetUserRoleDB(ctx, a.db).Where("user_id=? AND role_id=?", recordID, item.RoleID).Omit("user_id", "role_id").Updates(item)
			if err := result.Error; err != nil {
				span.Errorf(err.Error())
				return errors.New("更新用户角色数据发生错误")
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

func (a *User) queryRoles(ctx context.Context, userIDs ...string) (entity.UserRoles, error) {
	span := logger.StartSpan(ctx, "查询用户角色数据", a.getFuncName("queryRoles"))
	defer span.Finish()

	var list entity.UserRoles
	result := entity.GetUserRoleDB(ctx, a.db).Where("user_id IN(?)", userIDs).Find(&list)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询用户角色数据发生错误")
	}
	return list, nil
}
