package model

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	icontext "github.com/LyricTian/gin-admin/internal/app/context"
	"github.com/LyricTian/gin-admin/pkg/gormplus"
	"github.com/jinzhu/gorm"
)

// ExecTrans 执行事务
func ExecTrans(ctx context.Context, db *gormplus.DB, fn func(context.Context) error) error {
	if _, ok := icontext.FromTrans(ctx); ok {
		return fn(ctx)
	}

	transModel := NewTrans(db)
	trans, err := transModel.Begin(ctx)
	if err != nil {
		return err
	}

	err = fn(icontext.NewTrans(ctx, trans))
	if err != nil {
		_ = transModel.Rollback(ctx, trans)
		return err
	}
	return transModel.Commit(ctx, trans)
}

// WrapPageQuery 包装带有分页的查询
func WrapPageQuery(db *gorm.DB, pp *schema.PaginationParam, out interface{}) (*schema.PaginationResult, error) {
	if pp != nil {
		total, err := gormplus.Wrap(db).FindPage(db, pp.PageIndex, pp.PageSize, out)
		if err != nil {
			return nil, err
		}
		return &schema.PaginationResult{
			Total: total,
		}, nil
	}

	result := db.Find(out)
	return nil, result.Error
}
