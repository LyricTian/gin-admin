package model

import (
	"context"

	"github.com/LyricTian/gin-admin/v7/internal/app/model/gormx/entity"
	"github.com/LyricTian/gin-admin/v7/internal/app/schema"
	"github.com/LyricTian/gin-admin/v7/pkg/errors"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

// DemoSet 注入Demo
var DemoSet = wire.NewSet(wire.Struct(new(Demo), "*"))

// Demo 示例存储
type Demo struct {
	DB *gorm.DB
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
	opt := a.getQueryOption(opts...)

	db := entity.GetDemoDB(ctx, a.DB)
	if v := params.Code; v != "" {
		db = db.Where("code=?", v)
	}
	if v := params.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("code LIKE ? OR name LIKE ? OR memo LIKE ?", v, v, v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.Demos
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
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
func (a *Demo) Get(ctx context.Context, id string, opts ...schema.DemoQueryOptions) (*schema.Demo, error) {
	db := entity.GetDemoDB(ctx, a.DB).Where("id=?", id)
	var item entity.Demo
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaDemo(), nil
}

// Create 创建数据
func (a *Demo) Create(ctx context.Context, item schema.Demo) error {
	eitem := entity.SchemaDemo(item).ToDemo()
	result := entity.GetDemoDB(ctx, a.DB).Create(eitem)
	return errors.WithStack(result.Error)
}

// Update 更新数据
func (a *Demo) Update(ctx context.Context, id string, item schema.Demo) error {
	eitem := entity.SchemaDemo(item).ToDemo()
	result := entity.GetDemoDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	return errors.WithStack(result.Error)
}

// Delete 删除数据
func (a *Demo) Delete(ctx context.Context, id string) error {
	result := entity.GetDemoDB(ctx, a.DB).Where("id=?", id).Delete(entity.Demo{})
	return errors.WithStack(result.Error)
}

// UpdateStatus 更新状态
func (a *Demo) UpdateStatus(ctx context.Context, id string, status int) error {
	result := entity.GetDemoDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
