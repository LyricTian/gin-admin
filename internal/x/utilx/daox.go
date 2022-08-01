package utilx

import (
	"context"

	"github.com/LyricTian/gin-admin/v9/internal/x/contextx"

	"github.com/google/wire"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var TransRepoSet = wire.NewSet(wire.Struct(new(TransRepo), "*"))

type TransRepo struct {
	DB *gorm.DB
}

type TransFunc func(context.Context) error

func (a *TransRepo) Exec(ctx context.Context, fn TransFunc) error {
	if _, ok := contextx.FromTrans(ctx); ok {
		return fn(ctx)
	}

	return a.DB.Transaction(func(db *gorm.DB) error {
		return fn(contextx.NewTrans(ctx, db))
	})
}

func GetDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	db := defDB
	if tdb, ok := contextx.FromTrans(ctx); ok {
		db = tdb
	}
	if contextx.FromRowLock(ctx) {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	return db
}

func wrapQueryOptions(db *gorm.DB, opts QueryOptions) *gorm.DB {
	if len(opts.SelectFields) > 0 {
		db = db.Select(opts.SelectFields)
	}
	if len(opts.OmitFields) > 0 {
		db = db.Omit(opts.OmitFields...)
	}
	if len(opts.OrderFields) > 0 {
		db = db.Order(opts.OrderFields.ToSQL())
	}
	return db
}

func WrapPageQuery(ctx context.Context, db *gorm.DB, pp PaginationParam, opts QueryOptions, out interface{}) (*PaginationResult, error) {
	if pp.OnlyCount {
		var count int64
		err := db.Count(&count).Error
		if err != nil {
			return nil, err
		}
		return &PaginationResult{Total: count}, nil
	} else if !pp.Pagination {
		pageSize := pp.PageSize
		if pageSize > 0 {
			db = db.Limit(pageSize)
		}

		db = wrapQueryOptions(db, opts)
		err := db.Find(out).Error
		return nil, err
	}

	total, err := FindPage(ctx, db, pp, opts, out)
	if err != nil {
		return nil, err
	}

	return &PaginationResult{
		Total:    total,
		Current:  pp.GetCurrent(),
		PageSize: pp.GetPageSize(),
	}, nil
}

func FindPage(ctx context.Context, db *gorm.DB, pp PaginationParam, opts QueryOptions, out interface{}) (int64, error) {
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return 0, err
	} else if count == 0 {
		return count, nil
	}

	current, pageSize := pp.GetCurrent(), pp.GetPageSize()
	if current > 0 && pageSize > 0 {
		db = db.Offset((current - 1) * pageSize).Limit(pageSize)
	} else if pageSize > 0 {
		db = db.Limit(pageSize)
	}

	db = wrapQueryOptions(db, opts)
	err = db.Find(out).Error
	return count, err
}

func FindOne(ctx context.Context, db *gorm.DB, opts QueryOptions, out interface{}) (bool, error) {
	db = wrapQueryOptions(db, opts)
	result := db.First(out)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func Exists(ctx context.Context, db *gorm.DB) (bool, error) {
	var count int64
	result := db.Count(&count)
	if err := result.Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
