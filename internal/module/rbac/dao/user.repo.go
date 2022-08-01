package dao

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/module/rbac/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"gorm.io/gorm"
)

func GetUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return utilx.GetDB(ctx, defDB).Model(new(typed.User))
}

type UserRepo struct {
	DB *gorm.DB
}

func (a *UserRepo) Query(ctx context.Context, params typed.UserQueryParam, opts ...typed.UserQueryOptions) (*typed.UserQueryResult, error) {
	var opt typed.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetUserDB(ctx, a.DB)

	if v := params.Username; v != "" {
		db = db.Where("username=?", v)
	}
	if v := params.Status; v != "" {
		db = db.Where("status=?", v)
	}
	if v := params.LikeName; v != "" {
		db = db.Where("name like ?", "%"+v+"%")
	}
	if v := params.LikeUsername; v != "" {
		db = db.Where("username like ?", "%"+v+"%")
	}
	if v := params.RoleID; v != "" {
		db = db.Where("id in (?)", GetUserRoleDB(ctx, a.DB).Select("user_id").Where("role_id=?", v))
	}

	var list typed.Users
	pr, err := utilx.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &typed.UserQueryResult{
		PageResult: pr,
		Data:       list,
	}

	return qr, nil
}

func (a *UserRepo) Get(ctx context.Context, id string, opts ...typed.UserQueryOptions) (*typed.User, error) {
	var opt typed.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(typed.User)
	ok, err := utilx.FindOne(ctx, GetUserDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *UserRepo) Exists(ctx context.Context, id string) (bool, error) {
	exists, err := utilx.Exists(ctx, GetUserDB(ctx, a.DB).Where("id=?", id))
	return exists, errors.WithStack(err)
}

func (a *UserRepo) ExistsUsername(ctx context.Context, username string) (bool, error) {
	exists, err := utilx.Exists(ctx, GetUserDB(ctx, a.DB).Where("username=?", username))
	return exists, errors.WithStack(err)
}

func (a *UserRepo) Create(ctx context.Context, item *typed.User) error {
	result := GetUserDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *UserRepo) Update(ctx context.Context, item *typed.User) error {
	result := GetUserDB(ctx, a.DB).Where("id=?", item.ID).Omit("created_at", "created_by").Updates(item)
	return errors.WithStack(result.Error)
}

func (a *UserRepo) Delete(ctx context.Context, id string) error {
	result := GetUserDB(ctx, a.DB).Where("id=?", id).Delete(new(typed.User))
	return errors.WithStack(result.Error)
}

func (a *UserRepo) UpdateStatus(ctx context.Context, id string, status typed.UserStatus) error {
	result := GetUserDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}

func (a *UserRepo) UpdatePassword(ctx context.Context, id string, password string) error {
	result := GetUserDB(ctx, a.DB).Where("id=?", id).Update("password", password)
	return errors.WithStack(result.Error)
}
