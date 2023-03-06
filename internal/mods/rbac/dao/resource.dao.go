package dao

import (
	"context"

	"github.com/LyricTian/gin-admin/v10/internal/library/utilx"
	"github.com/LyricTian/gin-admin/v10/internal/mods/rbac/schema"
	"github.com/LyricTian/gin-admin/v10/pkg/errors"
	"gorm.io/gorm"
)

func GetResourceDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return utilx.GetDB(ctx, defDB).Model(new(schema.Resource))
}

type Resource struct {
	DB *gorm.DB
}

func (a *Resource) Query(ctx context.Context, params schema.ResourceQueryParam, opts ...schema.ResourceQueryOptions) (*schema.ResourceQueryResult, error) {
	var opt schema.ResourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetResourceDB(ctx, a.DB)

	if v := params.LikeCode; v != "" {
		db = db.Where("code LIKE ?", v+"%")
	}
	if v := params.Status; v != "" {
		db = db.Where("status=?", v)
	}

	var list schema.Resources
	pr, err := utilx.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.ResourceQueryResult{
		PageResult: pr,
		Data:       list,
	}

	return qr, nil
}

func (a *Resource) Get(ctx context.Context, id string, opts ...schema.ResourceQueryOptions) (*schema.Resource, error) {
	var opt schema.ResourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(schema.Resource)
	ok, err := utilx.FindOne(ctx, GetResourceDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

func (a *Resource) Exists(ctx context.Context, id string) (bool, error) {
	ok, err := utilx.Exists(ctx, GetResourceDB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

func (a *Resource) ExistsCode(ctx context.Context, code string) (bool, error) {
	ok, err := utilx.Exists(ctx, GetResourceDB(ctx, a.DB).Where("code=?", code))
	return ok, errors.WithStack(err)
}

func (a *Resource) Create(ctx context.Context, item *schema.Resource) error {
	result := GetResourceDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *Resource) Update(ctx context.Context, item *schema.Resource) error {
	result := GetResourceDB(ctx, a.DB).Where("id=?", item.ID).Omit("created_at").Updates(item)
	return errors.WithStack(result.Error)
}

func (a *Resource) Delete(ctx context.Context, id string) error {
	result := GetResourceDB(ctx, a.DB).Where("id=?", id).Delete(new(schema.Resource))
	return errors.WithStack(result.Error)
}
