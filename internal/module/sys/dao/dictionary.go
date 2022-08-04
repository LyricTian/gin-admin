package dao

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/module/sys/typed"
	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
	"gorm.io/gorm"
)

func GetDictionaryDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return utilx.GetDB(ctx, defDB).Model(new(typed.Dictionary))
}

type DictionaryRepo struct {
	DB *gorm.DB
}

func (a *DictionaryRepo) Query(ctx context.Context, params typed.DictionaryQueryParam, opts ...typed.DictionaryQueryOptions) (*typed.DictionaryQueryResult, error) {
	var opt typed.DictionaryQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetDictionaryDB(ctx, a.DB)

	if v := params.Ns; len(v) > 0 {
		db = db.Where("ns=?", v)
	}
	if v := params.Key; len(v) > 0 {
		db = db.Where("key=?", v)
	}

	var list typed.Dictionaries
	pr, err := utilx.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &typed.DictionaryQueryResult{
		PageResult: pr,
		Data:       list,
	}

	return qr, nil
}

func (a *DictionaryRepo) Get(ctx context.Context, id string, opts ...typed.DictionaryQueryOptions) (*typed.Dictionary, error) {
	var opt typed.DictionaryQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(typed.Dictionary)
	ok, err := utilx.FindOne(ctx, GetDictionaryDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *DictionaryRepo) Exists(ctx context.Context, id string) (bool, error) {
	exists, err := utilx.Exists(ctx, GetDictionaryDB(ctx, a.DB).Where("id=?", id))
	return exists, errors.WithStack(err)
}

func (a *DictionaryRepo) Create(ctx context.Context, item *typed.Dictionary) error {
	result := GetDictionaryDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *DictionaryRepo) Update(ctx context.Context, item *typed.Dictionary) error {
	result := GetDictionaryDB(ctx, a.DB).Where("id=?", item.ID).Omit("created_at", "created_by").Updates(item)
	return errors.WithStack(result.Error)
}

func (a *DictionaryRepo) Delete(ctx context.Context, id string) error {
	result := GetDictionaryDB(ctx, a.DB).Where("id=?", id).Delete(new(typed.Dictionary))
	return errors.WithStack(result.Error)
}
