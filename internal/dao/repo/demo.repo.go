package repo

import (
	"context"

	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v9/internal/dao/util"
	"github.com/LyricTian/gin-admin/v9/internal/schema"
	"github.com/LyricTian/gin-admin/v9/pkg/errors"
)

// Injection wire
var DemoSet = wire.NewSet(wire.Struct(new(DemoRepo), "*"))

func GetDemoDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(schema.Demo))
}

type DemoRepo struct {
	DB *gorm.DB
}

func (a *DemoRepo) Query(ctx context.Context, params schema.DemoQueryParam, opts ...schema.DemoQueryOptions) (*schema.DemoQueryResult, error) {
	var opt schema.DemoQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetDemoDB(ctx, a.DB)

	// TODO: Your where condition code here...

	if len(opt.SelectFields) > 0 {
		db = db.Select(opt.SelectFields)
	}

	if len(opt.OrderFields) > 0 {
		db = db.Order(util.ParseOrder(opt.OrderFields))
	}

	var list schema.Demos
	pr, err := util.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.DemoQueryResult{
		PageResult: pr,
		Data:       list,
	}

	return qr, nil
}

func (a *DemoRepo) Get(ctx context.Context, id string) (*schema.Demo, error) {
	item := new(schema.Demo)
	ok, err := util.FindOne(ctx, GetDemoDB(ctx, a.DB).Where("id=?", id), item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *DemoRepo) Create(ctx context.Context, item *schema.Demo) error {
	result := GetDemoDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *DemoRepo) Update(ctx context.Context, item *schema.Demo) error {
	result := GetDemoDB(ctx, a.DB).Where("id=?", item.ID).Updates(item)
	return errors.WithStack(result.Error)
}

func (a *DemoRepo) Delete(ctx context.Context, id string) error {
	result := GetDemoDB(ctx, a.DB).Where("id=?", id).Delete(schema.Demo{})
	return errors.WithStack(result.Error)
}
