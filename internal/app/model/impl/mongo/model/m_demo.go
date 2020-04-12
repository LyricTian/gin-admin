package model

import (
	"context"
	"time"

	"github.com/LyricTian/gin-admin/internal/app/model"
	"github.com/LyricTian/gin-admin/internal/app/model/impl/mongo/entity"
	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ model.IDemo = (*Demo)(nil)

// DemoSet 注入Demo
var DemoSet = wire.NewSet(wire.Struct(new(Demo), "*"), wire.Bind(new(model.IDemo), new(*Demo)))

// Demo 示例存储
type Demo struct {
	Client *mongo.Client
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

	c := entity.GetDemoCollection(ctx, a.Client)
	filter := DefaultFilter(ctx)
	if v := params.Code; v != "" {
		filter = append(filter, Filter("code", v))
	}
	if v := params.QueryValue; v != "" {
		filter = append(filter, Filter("$or", bson.A{
			OrRegexFilter("code", v),
			OrRegexFilter("name", v),
			OrRegexFilter("memo", v),
		}))
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByDESC))

	var list entity.Demos
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
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
	c := entity.GetDemoCollection(ctx, a.Client)
	filter := DefaultFilter(ctx, Filter("_id", recordID))
	var item entity.Demo
	ok, err := FindOne(ctx, c, filter, &item)
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
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := entity.GetDemoCollection(ctx, a.Client)
	err := Insert(ctx, c, eitem)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *Demo) Update(ctx context.Context, recordID string, item schema.Demo) error {
	eitem := entity.SchemaDemo(item).ToDemo()
	eitem.UpdatedAt = time.Now()
	c := entity.GetDemoCollection(ctx, a.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)), eitem)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *Demo) Delete(ctx context.Context, recordID string) error {
	c := entity.GetDemoCollection(ctx, a.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Demo) UpdateStatus(ctx context.Context, recordID string, status int) error {
	c := entity.GetDemoCollection(ctx, a.Client)
	err := UpdateFields(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)), bson.M{"status": status})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
