package model

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/internal/app/model/impl/gorm/entity"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var _ model.IUserRole = new(UserRole)

// UserRoleSet 注入UserRole
var UserRoleSet = wire.NewSet(NewUserRole, wire.Bind(new(model.IUserRole), new(*UserRole)))

// NewUserRole 创建用户角色存储实例
func NewUserRole(db *gorm.DB) *UserRole {
	return &UserRole{db}
}

// UserRole 用户角色存储
type UserRole struct {
	db *gorm.DB
}

func (a *UserRole) getQueryOption(opts ...schema.UserRoleQueryOptions) schema.UserRoleQueryOptions {
	var opt schema.UserRoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *UserRole) Query(ctx context.Context, params schema.UserRoleQueryParam, opts ...schema.UserRoleQueryOptions) (*schema.UserRoleQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetUserRoleDB(ctx, a.db)
	if v := params.UserID; v != "" {
		db = db.Where("user_id=?", v)
	}
	if v := params.UserIDs; len(v) > 0 {
		db = db.Where("user_id IN(?)", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.UserRoles
	pr, err := WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.UserRoleQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUserRoles(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *UserRole) Get(ctx context.Context, recordID string, opts ...schema.UserRoleQueryOptions) (*schema.UserRole, error) {
	db := entity.GetUserRoleDB(ctx, a.db).Where("record_id=?", recordID)
	var item entity.UserRole
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaUserRole(), nil
}

// Create 创建数据
func (a *UserRole) Create(ctx context.Context, item schema.UserRole) error {
	eitem := entity.SchemaUserRole(item).ToUserRole()
	result := entity.GetUserRoleDB(ctx, a.db).Create(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *UserRole) Update(ctx context.Context, recordID string, item schema.UserRole) error {
	eitem := entity.SchemaUserRole(item).ToUserRole()
	result := entity.GetUserRoleDB(ctx, a.db).Where("record_id=?", recordID).Omit("record_id").Updates(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *UserRole) Delete(ctx context.Context, recordID string) error {
	result := entity.GetUserRoleDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.UserRole{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeleteByUserID 根据用户ID删除数据
func (a *UserRole) DeleteByUserID(ctx context.Context, userID string) error {
	result := entity.GetUserRoleDB(ctx, a.db).Where("user_id=?", userID).Delete(entity.UserRole{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
