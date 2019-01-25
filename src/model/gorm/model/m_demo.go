package gormmodel

import (
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin/src/errors"
	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/model/gorm/entity"
	"github.com/LyricTian/gin-admin/src/schema"
	"github.com/LyricTian/gin-admin/src/service/gormplus"
)

// InitDemo 初始化demo存储
func InitDemo(db *gormplus.DB) *Demo {
	db.AutoMigrate(new(gormentity.Demo))
	return NewDemo(db)
}

// NewDemo 实例化demo存储
func NewDemo(db *gormplus.DB) *Demo {
	return &Demo{db: db}
}

// Demo demo存储
type Demo struct {
	db *gormplus.DB
}

func (a *Demo) getFuncName(name string) string {
	return fmt.Sprintf("gorm.demo.%s", name)
}

func (a *Demo) getDemoDB(ctx context.Context) *gormplus.DB {
	return FromTransDBWithModel(ctx, a.db, gormentity.Demo{})
}

func (a *Demo) getQueryOption(opts ...schema.DemoQueryOptions) schema.DemoQueryOptions {
	if len(opts) > 0 {
		return opts[0]
	}
	return schema.DemoQueryOptions{}
}

// Query 查询数据
func (a *Demo) Query(ctx context.Context, params schema.DemoQueryParam, opts ...schema.DemoQueryOptions) (schema.DemoQueryResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := a.getDemoDB(ctx).DB
	if v := params.Code; v != "" {
		db = db.Where("code LIKE ?", "%"+v+"%")
	}
	if v := params.Name; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	db = db.Order("id DESC")

	var qr schema.DemoQueryResult
	opt := a.getQueryOption(opts...)
	var items gormentity.Demos
	pr, err := WrapPageQuery(db, opt.PageParam, &items)
	if err != nil {
		span.Errorf(err.Error())
		return qr, errors.New("查询数据发生错误")
	}
	qr.PageResult = pr
	qr.Data = items.ToSchemaDemos()

	return qr, nil
}

// Get 查询指定数据
func (a *Demo) Get(ctx context.Context, recordID string) (*schema.Demo, error) {
	span := logger.StartSpan(ctx, "查询指定数据", a.getFuncName("Get"))
	defer span.Finish()

	db := a.getDemoDB(ctx).Where("record_id=?", recordID)
	var item gormentity.Demo
	ok, err := a.db.FindOne(db, &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaDemo(), nil
}

// CheckCode 检查编号是否存在
func (a *Demo) CheckCode(ctx context.Context, code string) (bool, error) {
	span := logger.StartSpan(ctx, "检查编号是否存在", a.getFuncName("CheckCode"))
	defer span.Finish()

	db := a.getDemoDB(ctx).Where("code=?", code)
	exists, err := a.db.Check(db)
	if err != nil {
		span.Errorf(err.Error())
		return false, errors.New("检查编号是否存在发生错误")
	}
	return exists, nil
}

// Create 创建数据
func (a *Demo) Create(ctx context.Context, item schema.Demo) error {
	span := logger.StartSpan(ctx, "创建数据", a.getFuncName("Create"))
	defer span.Finish()

	demo := gormentity.SchemaDemo(item).ToDemo()
	demo.Creator = FromUserID(ctx)
	result := a.getDemoDB(ctx).Create(demo)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("创建数据发生错误")
	}
	return nil
}

// Update 更新数据
func (a *Demo) Update(ctx context.Context, recordID string, item schema.Demo) error {
	span := logger.StartSpan(ctx, "更新数据", a.getFuncName("Update"))
	defer span.Finish()

	demo := gormentity.SchemaDemo(item).ToDemo()
	result := a.getDemoDB(ctx).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(demo)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新数据发生错误")
	}
	return nil
}

// Delete 删除数据
func (a *Demo) Delete(ctx context.Context, recordID string) error {
	span := logger.StartSpan(ctx, "删除数据", a.getFuncName("Delete"))
	defer span.Finish()

	result := a.getDemoDB(ctx).Where("record_id=?", recordID).Delete(gormentity.Demo{})
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("删除数据发生错误")
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Demo) UpdateStatus(ctx context.Context, recordID string, status int) error {
	span := logger.StartSpan(ctx, "更新状态", a.getFuncName("UpdateStatus"))
	defer span.Finish()

	result := a.getDemoDB(ctx).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新状态发生错误")
	}
	return nil
}
