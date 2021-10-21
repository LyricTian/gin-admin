package role

import (
	"context"

	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/internal/app/dao/util"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/LyricTian/gin-admin/v8/pkg/errors"
)

var RoleSet = wire.NewSet(wire.Struct(new(RoleRepo), "*"))

type RoleRepo struct {
	DB *gorm.DB
}

func (a *RoleRepo) getQueryOption(opts ...schema.RoleQueryOptions) schema.RoleQueryOptions {
	var opt schema.RoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *RoleRepo) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := GetRoleDB(ctx, a.DB)
	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}
	if v := params.Name; v != "" {
		db = db.Where("name=?", v)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ?", v)
	}

	if len(opt.SelectFields) > 0 {
		db = db.Select(opt.SelectFields)
	}

	if len(opt.OrderFields) > 0 {
		db = db.Order(util.ParseOrder(opt.OrderFields))
	}

	var list Roles
	pr, err := util.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RoleQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRoles(),
	}

	return qr, nil
}

func (a *RoleRepo) Get(ctx context.Context, id uint64, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	var role Role
	ok, err := util.FindOne(ctx, GetRoleDB(ctx, a.DB).Where("id=?", id), &role)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return role.ToSchemaRole(), nil
}

func (a *RoleRepo) Create(ctx context.Context, item schema.Role) error {
	eitem := SchemaRole(item).ToRole()
	result := GetRoleDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

func (a *RoleRepo) Update(ctx context.Context, id uint64, item schema.Role) error {
	eitem := SchemaRole(item).ToRole()
	result := GetRoleDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

func (a *RoleRepo) Delete(ctx context.Context, id uint64) error {
	result := GetRoleDB(ctx, a.DB).Where("id=?", id).Delete(Role{})
	return errors.WithStack(result.Error)
}

func (a *RoleRepo) UpdateStatus(ctx context.Context, id uint64, status int) error {
	result := GetRoleDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
