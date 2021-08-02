package user

import (
	"context"

	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao/util"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/errors"
)

// UserRoleSet 注入UserRole
var UserRoleSet = wire.NewSet(wire.Struct(new(UserRoleRepo), "*"))

// UserRoleRepo 用户角色存储
type UserRoleRepo struct {
	DB *gorm.DB
}

func (a *UserRoleRepo) getQueryOption(opts ...schema.UserRoleQueryOptions) schema.UserRoleQueryOptions {
	var opt schema.UserRoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *UserRoleRepo) Query(ctx context.Context, params schema.UserRoleQueryParam, opts ...schema.UserRoleQueryOptions) (*schema.UserRoleQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := GetUserRoleDB(ctx, a.DB)
	if v := params.UserID; v > 0 {
		db = db.Where("user_id=?", v)
	}
	if v := params.UserIDs; len(v) > 0 {
		db = db.Where("user_id IN (?)", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(util.ParseOrder(opt.OrderFields))

	var list UserRoles
	pr, err := util.WrapPageQuery(ctx, db, params.PaginationParam, &list)
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
func (a *UserRoleRepo) Get(ctx context.Context, id uint64, opts ...schema.UserRoleQueryOptions) (*schema.UserRole, error) {
	db := GetUserRoleDB(ctx, a.DB).Where("id=?", id)
	var item UserRole
	ok, err := util.FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaUserRole(), nil
}

// Create 创建数据
func (a *UserRoleRepo) Create(ctx context.Context, item schema.UserRole) error {
	eitem := SchemaUserRole(item).ToUserRole()
	result := GetUserRoleDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

// Update 更新数据
func (a *UserRoleRepo) Update(ctx context.Context, id uint64, item schema.UserRole) error {
	eitem := SchemaUserRole(item).ToUserRole()
	result := GetUserRoleDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// Delete 删除数据
func (a *UserRoleRepo) Delete(ctx context.Context, id uint64) error {
	result := GetUserRoleDB(ctx, a.DB).Where("id=?", id).Delete(UserRole{})
	return errors.WithStack(result.Error)
}

// DeleteByUserID 根据用户ID删除数据
func (a *UserRoleRepo) DeleteByUserID(ctx context.Context, userID uint64) error {
	result := GetUserRoleDB(ctx, a.DB).Where("user_id=?", userID).Delete(UserRole{})
	return errors.WithStack(result.Error)
}
