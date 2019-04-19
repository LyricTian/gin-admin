package model

import (
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin/internal/app/ginadmin/model/gorm/entity"
	"github.com/LyricTian/gin-admin/internal/app/ginadmin/schema"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/LyricTian/gin-admin/pkg/gormplus"
	"github.com/LyricTian/gin-admin/pkg/logger"
)

// NewDemo 创建demo存储实例
func NewDemo(db *gormplus.DB) *Demo {
	return &Demo{db}
}

// Demo demo存储
type Demo struct {
	db *gormplus.DB
}

func (a *Demo) getFuncName(name string) string {
	return fmt.Sprintf("gorm.model.Demo.%s", name)
}

func (a *Demo) getQueryOption(opts ...schema.DemoQueryOptions) schema.DemoQueryOptions {
	var opt schema.DemoQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *Demo) Query(ctx context.Context, params schema.DemoQueryParam, opts ...schema.DemoQueryOptions) (*schema.DemoQueryResult, error) {
	span := logger.StartSpan(ctx, "查询数据", a.getFuncName("Query"))
	defer span.Finish()

	db := entity.GetDemoDB(ctx, a.db).DB
	if v := params.Code; v != "" {
		db = db.Where("code=?", v)
	}
	if v := params.LikeCode; v != "" {
		db = db.Where("code LIKE ?", "%"+v+"%")
	}
	if v := params.LikeName; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	db = db.Order("id DESC")

	opt := a.getQueryOption(opts...)
	var list entity.Demos
	pr, err := WrapPageQuery(db, opt.PageParam, &list)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询数据发生错误")
	}
	qr := &schema.DemoQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaDemos(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *Demo) Get(ctx context.Context, recordID string, opts ...schema.DemoQueryOptions) (*schema.Demo, error) {
	span := logger.StartSpan(ctx, "查询指定数据", a.getFuncName("Get"))
	defer span.Finish()

	db := entity.GetDemoDB(ctx, a.db).Where("record_id=?", recordID)
	var item entity.Demo
	ok, err := a.db.FindOne(db, &item)
	if err != nil {
		span.Errorf(err.Error())
		return nil, errors.New("查询指定数据发生错误")
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaDemo(), nil
}

// Create 创建数据
func (a *Demo) Create(ctx context.Context, item schema.Demo) error {
	span := logger.StartSpan(ctx, "创建数据", a.getFuncName("Create"))
	defer span.Finish()

	demo := entity.SchemaDemo(item).ToDemo()
	result := entity.GetDemoDB(ctx, a.db).Create(demo)
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

	demo := entity.SchemaDemo(item).ToDemo()
	result := entity.GetDemoDB(ctx, a.db).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(demo)
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

	result := entity.GetDemoDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.Demo{})
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

	result := entity.GetDemoDB(ctx, a.db).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		span.Errorf(err.Error())
		return errors.New("更新状态发生错误")
	}
	return nil
}
