package model

import (
	"context"
	"fmt"
	"strings"

	icontext "github.com/LyricTian/gin-admin/internal/app/context"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/jinzhu/gorm"
)

// TransFunc 定义事务执行函数
type TransFunc func(context.Context) error

// ExecTrans 执行事务
func ExecTrans(ctx context.Context, db *gorm.DB, fn TransFunc) error {
	transModel := &Trans{DB: db}
	return transModel.Exec(ctx, fn)
}

// ExecTransWithLock 执行事务（加锁）
func ExecTransWithLock(ctx context.Context, db *gorm.DB, fn TransFunc) error {
	if !icontext.FromTransLock(ctx) {
		ctx = icontext.NewTransLock(ctx)
	}
	return ExecTrans(ctx, db, fn)
}

// WrapPageQuery 包装带有分页的查询
func WrapPageQuery(ctx context.Context, db *gorm.DB, pp schema.PaginationParam, out interface{}) (*schema.PaginationResult, error) {
	if pp.OnlyCount {
		var count int
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
func FindPage(ctx context.Context, db *gorm.DB, pp schema.PaginationParam, out interface{}) (int, error) {
	var count int
	err := db.Count(&count).Error
	if err != nil {
		return 0, err
	} else if count == 0 {
		return count, nil
	}

	current, pageSize := pp.GetCurrent(), pp.GetPageSize()
	db = db.Offset((current - 1) * pageSize).Limit(pageSize)
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
	var count int
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
