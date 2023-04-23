package dal

import (
	"context"

	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/internal/utils"
	"github.com/LyricTian/gin-admin/v10/pkg/errors"
	"gorm.io/gorm"
)

// Get user storage instance
func GetUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return utils.GetDB(ctx, defDB).Model(new(schema.User))
}

// User management for RBAC
type User struct {
	DB *gorm.DB
}

// Query users from the database based on the provided parameters and options.
func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetUserDB(ctx, a.DB)
	if v := params.LikeUsername; len(v) > 0 {
		db = db.Where("username LIKE ?", "%"+v+"%")
	}
	if v := params.LikeName; len(v) > 0 {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; len(v) > 0 {
		db = db.Where("status = ?", v)
	}

	var list schema.Users
	pageResult, err := utils.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queryResult := &schema.UserQueryResult{
		PageResult: pageResult,
		Data:       list,
	}
	return queryResult, nil
}

// Get the specified user from the database.
func (a *User) Get(ctx context.Context, id string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(schema.User)
	ok, err := utils.FindOne(ctx, GetUserDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// Exist checks if the specified user exists in the database.
func (a *User) Exists(ctx context.Context, id string) (bool, error) {
	ok, err := utils.Exists(ctx, GetUserDB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

func (a *User) ExistsUsername(ctx context.Context, username string) (bool, error) {
	ok, err := utils.Exists(ctx, GetUserDB(ctx, a.DB).Where("username=?", username))
	return ok, errors.WithStack(err)
}

// Create a new user.
func (a *User) Create(ctx context.Context, item *schema.User) error {
	result := GetUserDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update the specified user in the database.
func (a *User) Update(ctx context.Context, item *schema.User) error {
	result := GetUserDB(ctx, a.DB).Where("id=?", item.ID).Select("*").Omit("created_at").Updates(item)
	return errors.WithStack(result.Error)
}

// Delete the specified user from the database.
func (a *User) Delete(ctx context.Context, id string) error {
	result := GetUserDB(ctx, a.DB).Where("id=?", id).Delete(new(schema.User))
	return errors.WithStack(result.Error)
}
