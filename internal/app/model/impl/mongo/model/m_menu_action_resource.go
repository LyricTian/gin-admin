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

var _ model.IMenuActionResource = (*MenuActionResource)(nil)

// MenuActionResourceSet 注入MenuActionResource
var MenuActionResourceSet = wire.NewSet(wire.Struct(new(MenuActionResource), "*"), wire.Bind(new(model.IMenuActionResource), new(*MenuActionResource)))

// MenuActionResource 菜单动作关联资源存储
type MenuActionResource struct {
	Client *mongo.Client
}

func (a *MenuActionResource) getQueryOption(opts ...schema.MenuActionResourceQueryOptions) schema.MenuActionResourceQueryOptions {
	var opt schema.MenuActionResourceQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *MenuActionResource) Query(ctx context.Context, params schema.MenuActionResourceQueryParam, opts ...schema.MenuActionResourceQueryOptions) (*schema.MenuActionResourceQueryResult, error) {
	opt := a.getQueryOption(opts...)

	c := entity.GetMenuActionResourceCollection(ctx, a.Client)
	filter := DefaultFilter(ctx)
	menuIDs := params.MenuIDs
	if v := params.MenuID; v != "" {
		menuIDs = append(menuIDs, v)
	}
	if v := menuIDs; len(v) > 0 {
		actionIDs, err := a.queryActionIDs(ctx, v...)
		if err != nil {
			return nil, err
		}
		filter = append(filter, Filter("action_id", bson.M{"$in": actionIDs}))
	}
	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByASC))

	var list entity.MenuActionResources
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.MenuActionResourceQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaMenuActionResources(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *MenuActionResource) Get(ctx context.Context, recordID string, opts ...schema.MenuActionResourceQueryOptions) (*schema.MenuActionResource, error) {
	c := entity.GetMenuActionResourceCollection(ctx, a.Client)
	filter := DefaultFilter(ctx, Filter("_id", recordID))
	var item entity.MenuActionResource
	ok, err := FindOne(ctx, c, filter, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaMenuActionResource(), nil
}

// Create 创建数据
func (a *MenuActionResource) Create(ctx context.Context, item schema.MenuActionResource) error {
	eitem := entity.SchemaMenuActionResource(item).ToMenuActionResource()
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := entity.GetMenuActionResourceCollection(ctx, a.Client)
	err := Insert(ctx, c, eitem)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *MenuActionResource) Update(ctx context.Context, recordID string, item schema.MenuActionResource) error {
	eitem := entity.SchemaMenuActionResource(item).ToMenuActionResource()
	eitem.UpdatedAt = time.Now()
	c := entity.GetMenuActionResourceCollection(ctx, a.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)), eitem)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *MenuActionResource) Delete(ctx context.Context, recordID string) error {
	c := entity.GetMenuActionResourceCollection(ctx, a.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeleteByActionID 根据动作ID删除数据
func (a *MenuActionResource) DeleteByActionID(ctx context.Context, actionID string) error {
	c := entity.GetMenuActionResourceCollection(ctx, a.Client)
	err := DeleteMany(ctx, c, DefaultFilter(ctx, Filter("action_id", actionID)))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeleteByMenuID 根据菜单ID删除数据
func (a *MenuActionResource) DeleteByMenuID(ctx context.Context, menuID string) error {
	actionIDs, err := a.queryActionIDs(ctx, menuID)
	if err != nil {
		return err
	}

	c := entity.GetMenuActionResourceCollection(ctx, a.Client)
	err = DeleteMany(ctx, c, DefaultFilter(ctx, Filter("action_id", bson.M{"$in": actionIDs})))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *MenuActionResource) queryActionIDs(ctx context.Context, menuIDs ...string) ([]interface{}, error) {
	result, err := entity.GetMenuActionCollection(ctx, a.Client).Distinct(ctx, "_id", DefaultFilter(ctx, Filter("menu_id", bson.M{"$in": menuIDs})))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}
