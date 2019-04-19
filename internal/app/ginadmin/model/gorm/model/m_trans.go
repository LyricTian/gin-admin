package model

import (
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/LyricTian/gin-admin/pkg/gormplus"
	"github.com/LyricTian/gin-admin/pkg/logger"
)

// NewTrans 创建事务管理实例
func NewTrans(db *gormplus.DB) *Trans {
	return &Trans{db}
}

// Trans 事务管理
type Trans struct {
	db *gormplus.DB
}

func (a *Trans) getFuncName(name string) string {
	return fmt.Sprintf("gorm.model.Trans.%s", name)
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
