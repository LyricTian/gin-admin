package user

import (
	"context"

	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao/util"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/errors"
)

// UserSet 注入User
var UserSet = wire.NewSet(wire.Struct(new(UserRepo), "*"))

// UserRepo 用户存储
type UserRepo struct {
	DB *gorm.DB
}

func (a *UserRepo) getQueryOption(opts ...schema.UserQueryOptions) schema.UserQueryOptions {
	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *UserRepo) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := GetUserDB(ctx, a.DB)
	if v := params.UserName; v != "" {
		db = db.Where("user_name=?", v)
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	if v := params.RoleIDs; len(v) > 0 {
		subQuery := GetUserRoleDB(ctx, a.DB).
			Select("user_id").
			Where("role_id IN (?)", v)
		db = db.Where("id IN (?)", subQuery)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("user_name LIKE ? OR real_name LIKE ?", v, v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(util.ParseOrder(opt.OrderFields))

	var list Users
	pr, err := util.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.UserQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUsers(),
	}
	return qr, nil
}

// Get 查询指定数据
func (a *UserRepo) Get(ctx context.Context, id uint64, opts ...schema.UserQueryOptions) (*schema.User, error) {
	var item User
	ok, err := util.FindOne(ctx, GetUserDB(ctx, a.DB).Where("id=?", id), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaUser(), nil
}

// Create 创建数据
func (a *UserRepo) Create(ctx context.Context, item schema.User) error {
	sitem := SchemaUser(item)
	result := GetUserDB(ctx, a.DB).Create(sitem.ToUser())
	return errors.WithStack(result.Error)
}

// Update 更新数据
func (a *UserRepo) Update(ctx context.Context, id uint64, item schema.User) error {
	eitem := SchemaUser(item).ToUser()
	result := GetUserDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// Delete 删除数据
func (a *UserRepo) Delete(ctx context.Context, id uint64) error {
	result := GetUserDB(ctx, a.DB).Where("id=?", id).Delete(User{})
	return errors.WithStack(result.Error)
}

// UpdateStatus 更新状态
func (a *UserRepo) UpdateStatus(ctx context.Context, id uint64, status int) error {
	result := GetUserDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}

// UpdatePassword 更新密码
func (a *UserRepo) UpdatePassword(ctx context.Context, id uint64, password string) error {
	result := GetUserDB(ctx, a.DB).Where("id=?", id).Update("password", password)
	return errors.WithStack(result.Error)
}
