package common

import (
	"context"

	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
	"github.com/pkg/errors"
)

// GetTrans 获取事务
func GetTrans(trans interface{}) *gormplus.DB {
	return trans.(*gormplus.DB)
}

// NewTrans 事务管理
func NewTrans(db *gormplus.DB) *Trans {
	return &Trans{db}
}

// Trans 事务管理
type Trans struct {
	db *gormplus.DB
}

// Begin 开启事务
func (a *Trans) Begin(ctx context.Context) (interface{}, error) {
	span := logger.StartSpan(ctx, "开启事务", "trans.Begin")
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
	span := logger.StartSpan(ctx, "提交事务", "trans.Commit")
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
	span := logger.StartSpan(ctx, "回滚事务", "trans.Rollback")
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
