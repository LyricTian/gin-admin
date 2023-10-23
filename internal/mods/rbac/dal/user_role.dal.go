package dal

import (
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/pkg/errors"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
	"gorm.io/gorm"
)

// Get user role storage instance
func GetUserRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.UserRole))
}

// User roles for RBAC
type UserRole struct {
	DB *gorm.DB
}

// Query user roles from the database based on the provided parameters and options.
func (a *UserRole) Query(ctx context.Context, params schema.UserRoleQueryParam, opts ...schema.UserRoleQueryOptions) (*schema.UserRoleQueryResult, error) {
	var opt schema.UserRoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := a.DB.Table(fmt.Sprintf("%s AS a", new(schema.UserRole).TableName()))
	if opt.JoinRole {
		db = db.Joins(fmt.Sprintf("left join %s b on a.role_id=b.id", new(schema.Role).TableName()))
		db = db.Select("a.*,b.name as role_name")
	}

	if v := params.InUserIDs; len(v) > 0 {
		db = db.Where("a.user_id IN (?)", v)
	}
	if v := params.UserID; len(v) > 0 {
		db = db.Where("a.user_id = ?", v)
	}
	if v := params.RoleID; len(v) > 0 {
		db = db.Where("a.role_id = ?", v)
	}

	var list schema.UserRoles
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queryResult := &schema.UserRoleQueryResult{
		PageResult: pageResult,
		Data:       list,
	}
	return queryResult, nil
}

// Get the specified user role from the database.
func (a *UserRole) Get(ctx context.Context, id string, opts ...schema.UserRoleQueryOptions) (*schema.UserRole, error) {
	var opt schema.UserRoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(schema.UserRole)
	ok, err := util.FindOne(ctx, GetUserRoleDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// Exist checks if the specified user role exists in the database.
func (a *UserRole) Exists(ctx context.Context, id string) (bool, error) {
	ok, err := util.Exists(ctx, GetUserRoleDB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

// Create a new user role.
func (a *UserRole) Create(ctx context.Context, item *schema.UserRole) error {
	result := GetUserRoleDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update the specified user role in the database.
func (a *UserRole) Update(ctx context.Context, item *schema.UserRole) error {
	result := GetUserRoleDB(ctx, a.DB).Where("id=?", item.ID).Select("*").Omit("created_at").Updates(item)
	return errors.WithStack(result.Error)
}

// Delete the specified user role from the database.
func (a *UserRole) Delete(ctx context.Context, id string) error {
	result := GetUserRoleDB(ctx, a.DB).Where("id=?", id).Delete(new(schema.UserRole))
	return errors.WithStack(result.Error)
}

func (a *UserRole) DeleteByUserID(ctx context.Context, userID string) error {
	result := GetUserRoleDB(ctx, a.DB).Where("user_id=?", userID).Delete(new(schema.UserRole))
	return errors.WithStack(result.Error)
}

func (a *UserRole) DeleteByRoleID(ctx context.Context, roleID string) error {
	result := GetUserRoleDB(ctx, a.DB).Where("role_id=?", roleID).Delete(new(schema.UserRole))
	return errors.WithStack(result.Error)
}
