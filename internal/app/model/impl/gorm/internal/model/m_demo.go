package model

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/errors"
	"github.com/LyricTian/gin-admin/internal/app/model/impl/gorm/internal/entity"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/gormplus"
)

// NewDemo 创建demo存储实例
func NewDemo(db *gormplus.DB) *Demo {
	return &Demo{db}
}

// Demo demo存储
type Demo struct {
	db *gormplus.DB
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
		return nil, errors.WithStack(err)
	}
	qr := &schema.DemoQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaDemos(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *Demo) Get(ctx context.Context, recordID string, opts ...schema.DemoQueryOptions) (*schema.Demo, error) {
	db := entity.GetDemoDB(ctx, a.db).Where("record_id=?", recordID)
	var item entity.Demo
	ok, err := a.db.FindOne(db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaDemo(), nil
}

// Create 创建数据
func (a *Demo) Create(ctx context.Context, item schema.Demo) error {
	demo := entity.SchemaDemo(item).ToDemo()
	result := entity.GetDemoDB(ctx, a.db).Create(demo)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *Demo) Update(ctx context.Context, recordID string, item schema.Demo) error {
	demo := entity.SchemaDemo(item).ToDemo()
	result := entity.GetDemoDB(ctx, a.db).Where("record_id=?", recordID).Omit("record_id", "creator").Updates(demo)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *Demo) Delete(ctx context.Context, recordID string) error {
	result := entity.GetDemoDB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.Demo{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Demo) UpdateStatus(ctx context.Context, recordID string, status int) error {
	result := entity.GetDemoDB(ctx, a.db).Where("record_id=?", recordID).Update("status", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
