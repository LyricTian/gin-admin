package gormcommon

import (
	"context"
	"fmt"

	gcontext "github.com/LyricTian/gin-admin/src/context"
	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
)

// FromTransDB 从上下文中获取事务DB，如果不存在则使用默认DB
func FromTransDB(ctx context.Context, defaultDB *gormplus.DB) *gormplus.DB {
	trans, ok := gcontext.FromTrans(ctx)
	if ok {
		db, ok := trans.(*gormplus.DB)
		if ok {
			return db
		}
	}
	return defaultDB
}

// ExecTrans 执行事务
func ExecTrans(ctx context.Context, db *gormplus.DB, fn func(context.Context) error) error {
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

// NewTrans 事务管理
func NewTrans(db *gormplus.DB) *Trans {
	return &Trans{db}
}

// Trans 事务管理
type Trans struct {
	db *gormplus.DB
}

func (a *Trans) getFuncName(name string) string {
	return fmt.Sprintf("trans.%s", name)
}

// Begin 开启事务
func (a *Trans) Begin(ctx context.Context) (interface{}, error) {
	span := logger.StartSpan(ctx, "开启事务", a.getFuncName("Begin"))
	defer span.Finish()

	result := a.db.Begin()
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("开启事务发生错误")
	}
	return gormplus.Wrap(result), nil
}

// Commit 提交事务
func (a *Trans) Commit(ctx context.Context, trans interface{}) error {
	span := logger.StartSpan(ctx, "提交事务", a.getFuncName("Commit"))
	defer span.Finish()

	db, ok := trans.(*gormplus.DB)
	if !ok {
		span.Warnf("未知的事务类型")
		return errors.New("未知的事务类型")
	}

	result := db.Commit()
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("提交事务发生错误")
	}
	return nil
}

// Rollback 回滚事务
func (a *Trans) Rollback(ctx context.Context, trans interface{}) error {
	span := logger.StartSpan(ctx, "回滚事务", a.getFuncName("Rollback"))
	defer span.Finish()

	db, ok := trans.(*gormplus.DB)
	if !ok {
		span.Warnf("未知的事务类型")
		return errors.New("未知的事务类型")
	}

	result := db.Rollback()
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("回滚事务发生错误")
	}
	return nil
}
