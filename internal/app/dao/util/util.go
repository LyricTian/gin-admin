package util

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/LyricTian/gin-admin/v8/internal/app/contextx"
	"github.com/LyricTian/gin-admin/v8/internal/app/schema"
)

// Define base model
type Model struct {
	ID        uint64 `gorm:"primaryKey;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Get gorm.DB from context
func GetDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	trans, ok := contextx.FromTrans(ctx)
	if ok && !contextx.FromNoTrans(ctx) {
		db, ok := trans.(*gorm.DB)
		if ok {
			if contextx.FromTransLock(ctx) {
				db = db.Clauses(clause.Locking{Strength: "UPDATE"})
			}
			return db
		}
	}

	return defDB
}

// Get gorm.DB.Model from context
func GetDBWithModel(ctx context.Context, defDB *gorm.DB, m interface{}) *gorm.DB {
	return GetDB(ctx, defDB).Model(m)
}

// Define transaction execute function
type TransFunc func(context.Context) error

func ExecTrans(ctx context.Context, db *gorm.DB, fn TransFunc) error {
	transModel := &Trans{DB: db}
	return transModel.Exec(ctx, fn)
}

func ExecTransWithLock(ctx context.Context, db *gorm.DB, fn TransFunc) error {
	if !contextx.FromTransLock(ctx) {
		ctx = contextx.NewTransLock(ctx)
	}
	return ExecTrans(ctx, db, fn)
}

func WrapPageQuery(ctx context.Context, db *gorm.DB, pp schema.PaginationParam, out interface{}) (*schema.PaginationResult, error) {
	if pp.OnlyCount {
		var count int64
		err := db.Count(&count).Error
		if err != nil {
			return nil, err
		}
		return &schema.PaginationResult{Total: count}, nil
	} else if !pp.Pagination {
		err := db.Find(out).Error
		return nil, err
	}

	total, err := FindPage(ctx, db, pp, out)
	if err != nil {
		return nil, err
	}

	return &schema.PaginationResult{
		Total:    total,
		Current:  pp.GetCurrent(),
		PageSize: pp.GetPageSize(),
	}, nil
}

func FindPage(ctx context.Context, db *gorm.DB, pp schema.PaginationParam, out interface{}) (int64, error) {
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

	err = db.Find(out).Error
	return count, err
}

func FindOne(ctx context.Context, db *gorm.DB, out interface{}) (bool, error) {
	result := db.First(out)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func Check(ctx context.Context, db *gorm.DB) (bool, error) {
	var count int64
	result := db.Count(&count)
	if err := result.Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Define order fields convert function
type OrderFieldFunc func(string) string

func ParseOrder(items []*schema.OrderField, handle ...OrderFieldFunc) string {
	orders := make([]string, len(items))

	for i, item := range items {
		key := item.Key
		if len(handle) > 0 {
			key = handle[0](key)
		}

		direction := "ASC"
		if item.Direction == schema.OrderByDESC {
			direction = "DESC"
		}
		orders[i] = fmt.Sprintf("%s %s", key, direction)
	}

	return strings.Join(orders, ",")
}
