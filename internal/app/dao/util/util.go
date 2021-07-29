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

// Model Define basic
type Model struct {
	ID        uint64 `gorm:"primaryKey;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetDB Get gorm.DB from context
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

// GetDBWithModel Get gorm.DB.Model from context
func GetDBWithModel(ctx context.Context, defDB *gorm.DB, m interface{}) *gorm.DB {
	return GetDB(ctx, defDB).Model(m)
}

// TransFunc 定义事务执行函数
type TransFunc func(context.Context) error

// ExecTrans 执行事务
func ExecTrans(ctx context.Context, db *gorm.DB, fn TransFunc) error {
	transModel := &Trans{DB: db}
	return transModel.Exec(ctx, fn)
}

// ExecTransWithLock 执行事务（加锁）
func ExecTransWithLock(ctx context.Context, db *gorm.DB, fn TransFunc) error {
	if !contextx.FromTransLock(ctx) {
		ctx = contextx.NewTransLock(ctx)
	}
	return ExecTrans(ctx, db, fn)
}

// WrapPageQuery 包装带有分页的查询
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

// FindPage 查询分页数据
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

// FindOne 查询单条数据
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

// Check 检查数据是否存在
func Check(ctx context.Context, db *gorm.DB) (bool, error) {
	var count int64
	result := db.Count(&count)
	if err := result.Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// OrderFieldFunc 排序字段转换函数
type OrderFieldFunc func(string) string

// ParseOrder 解析排序字段
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
