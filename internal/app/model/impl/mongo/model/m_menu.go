package model

import (
	"context"
	"fmt"
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

var _ model.IMenu = (*Menu)(nil)

// MenuSet 注入Menu
var MenuSet = wire.NewSet(wire.Struct(new(Menu), "*"), wire.Bind(new(model.IMenu), new(*Menu)))

// Menu 菜单存储
type Menu struct {
	Client *mongo.Client
}

func (a *Menu) getQueryOption(opts ...schema.MenuQueryOptions) schema.MenuQueryOptions {
	var opt schema.MenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *Menu) Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (*schema.MenuQueryResult, error) {
	opt := a.getQueryOption(opts...)

	c := entity.GetMenuCollection(ctx, a.Client)
	filter := DefaultFilter(ctx)
	if v := params.RecordIDs; len(v) > 0 {
		filter = append(filter, Filter("_id", bson.M{"$in": v}))
	}
	if v := params.Name; v != "" {
		filter = append(filter, Filter("name", v))
	}
	if v := params.QueryValue; v != "" {
		filter = append(filter, Filter("$or", bson.A{
			OrRegexFilter("name", v),
			OrRegexFilter("memo", v),
		}))
	}
	if v := params.ParentID; v != nil {
		filter = append(filter, Filter("parent_id", *v))
	}
	if v := params.PrefixParentPath; v != "" {
		filter = append(filter, RegexFilter("parent_path", fmt.Sprintf("^%s.*", v)))
	}
	if v := params.ShowStatus; v != 0 {
		filter = append(filter, Filter("show_status", v))
	}
	if v := params.Status; v != 0 {
		filter = append(filter, Filter("status", v))
	}
	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByDESC))

	var list entity.Menus
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.MenuQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaMenus(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *Menu) Get(ctx context.Context, recordID string, opts ...schema.MenuQueryOptions) (*schema.Menu, error) {
	c := entity.GetMenuCollection(ctx, a.Client)
	filter := DefaultFilter(ctx, Filter("_id", recordID))
	var item entity.Menu
	ok, err := FindOne(ctx, c, filter, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaMenu(), nil
}

// Create 创建数据
func (a *Menu) Create(ctx context.Context, item schema.Menu) error {
	eitem := entity.SchemaMenu(item).ToMenu()
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := entity.GetMenuCollection(ctx, a.Client)
	err := Insert(ctx, c, eitem)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *Menu) Update(ctx context.Context, recordID string, item schema.Menu) error {
	eitem := entity.SchemaMenu(item).ToMenu()
	eitem.UpdatedAt = time.Now()
	c := entity.GetMenuCollection(ctx, a.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)), eitem)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *Menu) Delete(ctx context.Context, recordID string) error {
	c := entity.GetMenuCollection(ctx, a.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Menu) UpdateStatus(ctx context.Context, recordID string, status int) error {
	c := entity.GetMenuCollection(ctx, a.Client)
	err := UpdateFields(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)), bson.M{"status": status})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateParentPath 更新父级路径
func (a *Menu) UpdateParentPath(ctx context.Context, recordID, parentPath string) error {
	c := entity.GetMenuCollection(ctx, a.Client)
	err := UpdateFields(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)), bson.M{"parent_path": parentPath})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
