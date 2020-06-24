package model

import (
	"context"
	"time"

	"github.com/LyricTian/gin-admin/v6/internal/app/model"
	"github.com/LyricTian/gin-admin/v6/internal/app/model/impl/mongo/entity"
	"github.com/LyricTian/gin-admin/v6/internal/app/schema"
	"github.com/LyricTian/gin-admin/v6/pkg/errors"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ model.IMenuAction = (*MenuAction)(nil)

// MenuActionSet 注入MenuAction
var MenuActionSet = wire.NewSet(wire.Struct(new(MenuAction), "*"), wire.Bind(new(model.IMenuAction), new(*MenuAction)))

// MenuAction 菜单动作存储
type MenuAction struct {
	Client *mongo.Client
}

func (a *MenuAction) getQueryOption(opts ...schema.MenuActionQueryOptions) schema.MenuActionQueryOptions {
	var opt schema.MenuActionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *MenuAction) Query(ctx context.Context, params schema.MenuActionQueryParam, opts ...schema.MenuActionQueryOptions) (*schema.MenuActionQueryResult, error) {
	opt := a.getQueryOption(opts...)

	c := entity.GetMenuActionCollection(ctx, a.Client)
	filter := DefaultFilter(ctx)
	if v := params.MenuID; v != "" {
		filter = append(filter, Filter("menu_id", v))
	}
	if v := params.IDs; len(v) > 0 {
		filter = append(filter, Filter("_id", bson.M{"$in": v}))
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByASC))

	var list entity.MenuActions
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.MenuActionQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaMenuActions(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *MenuAction) Get(ctx context.Context, id string, opts ...schema.MenuActionQueryOptions) (*schema.MenuAction, error) {
	c := entity.GetMenuActionCollection(ctx, a.Client)
	filter := DefaultFilter(ctx, Filter("_id", id))
	var item entity.MenuAction
	ok, err := FindOne(ctx, c, filter, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaMenuAction(), nil
}

// Create 创建数据
func (a *MenuAction) Create(ctx context.Context, item schema.MenuAction) error {
	eitem := entity.SchemaMenuAction(item).ToMenuAction()
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := entity.GetMenuActionCollection(ctx, a.Client)
	err := Insert(ctx, c, eitem)
	return errors.WithStack(err)
}

// Update 更新数据
func (a *MenuAction) Update(ctx context.Context, id string, item schema.MenuAction) error {
	eitem := entity.SchemaMenuAction(item).ToMenuAction()
	eitem.UpdatedAt = time.Now()
	c := entity.GetMenuActionCollection(ctx, a.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", id)), eitem)
	return errors.WithStack(err)
}

// Delete 删除数据
func (a *MenuAction) Delete(ctx context.Context, id string) error {
	c := entity.GetMenuActionCollection(ctx, a.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", id)))
	return errors.WithStack(err)
}

// DeleteByMenuID 根据菜单ID删除数据
func (a *MenuAction) DeleteByMenuID(ctx context.Context, menuID string) error {
	c := entity.GetMenuActionCollection(ctx, a.Client)
	err := DeleteMany(ctx, c, DefaultFilter(ctx, Filter("menu_id", menuID)))
	return errors.WithStack(err)
}
