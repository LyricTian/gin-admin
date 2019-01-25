package model

import (
	"context"

	gcontext "github.com/LyricTian/gin-admin/src/context"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/jinzhu/gorm"
)

// FromUserID 从上下文中获取用户ID
func FromUserID(ctx context.Context) string {
	return gcontext.FromUserID(ctx)
}

// FromDB 从上下文中获取DB，如果不存在则使用默认DB
func FromDB(ctx context.Context, defDB *gormplus.DB) *gormplus.DB {
	trans, ok := gcontext.FromTrans(ctx)
	if ok {
		db, ok := trans.(*gormplus.DB)
		if ok {
			return db
		}
	}
	return defDB
}

// FromDBWithModel 从上下文获取DB，并创建模型
func FromDBWithModel(ctx context.Context, defDB *gormplus.DB, v interface{}) *gormplus.DB {
	return gormplus.Wrap(FromDB(ctx, defDB).Model(v))
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

// ExecTrans 执行事务
func ExecTrans(ctx context.Context, db *gormplus.DB, fn func(context.Context) error) error {
	if _, ok := gcontext.FromTrans(ctx); ok {
		return fn(ctx)
	}

	transModel := NewTrans(db)
	trans, err := transModel.Begin(ctx)
	if err != nil {
		return err
	}

	err = fn(gcontext.NewTrans(ctx, trans))
	if err != nil {
		_ = transModel.Rollback(ctx, trans)
		return err
	}
	return transModel.Commit(ctx, trans)
}
