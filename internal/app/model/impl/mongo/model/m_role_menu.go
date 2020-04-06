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

var _ model.IRoleMenu = (*RoleMenu)(nil)

// RoleMenuSet 注入RoleMenu
var RoleMenuSet = wire.NewSet(wire.Struct(new(RoleMenu), "*"), wire.Bind(new(model.IRoleMenu), new(*RoleMenu)))

// RoleMenu 角色菜单存储
type RoleMenu struct {
	Client *mongo.Client
}

func (a *RoleMenu) getQueryOption(opts ...schema.RoleMenuQueryOptions) schema.RoleMenuQueryOptions {
	var opt schema.RoleMenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *RoleMenu) Query(ctx context.Context, params schema.RoleMenuQueryParam, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenuQueryResult, error) {
	opt := a.getQueryOption(opts...)

	c := entity.GetRoleMenuCollection(ctx, a.Client)
	filter := DefaultFilter(ctx)
	roleIDs := params.RoleIDs
	if v := params.RoleID; v != "" {
		roleIDs = append(roleIDs, v)
	}
	if v := roleIDs; len(v) > 0 {
		filter = append(filter, Filter("role_id", bson.M{"$in": v}))
	}
	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByDESC))

	var list entity.RoleMenus
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RoleMenuQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRoleMenus(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *RoleMenu) Get(ctx context.Context, recordID string, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenu, error) {
	c := entity.GetRoleMenuCollection(ctx, a.Client)
	filter := DefaultFilter(ctx, Filter("_id", recordID))
	var item entity.RoleMenu
	ok, err := FindOne(ctx, c, filter, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaRoleMenu(), nil
}

// Create 创建数据
func (a *RoleMenu) Create(ctx context.Context, item schema.RoleMenu) error {
	eitem := entity.SchemaRoleMenu(item).ToRoleMenu()
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := entity.GetRoleMenuCollection(ctx, a.Client)
	err := Insert(ctx, c, eitem)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *RoleMenu) Update(ctx context.Context, recordID string, item schema.RoleMenu) error {
	eitem := entity.SchemaRoleMenu(item).ToRoleMenu()
	eitem.UpdatedAt = time.Now()
	c := entity.GetRoleMenuCollection(ctx, a.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)), eitem)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *RoleMenu) Delete(ctx context.Context, recordID string) error {
	c := entity.GetRoleMenuCollection(ctx, a.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeleteByRoleID 根据角色ID删除数据
func (a *RoleMenu) DeleteByRoleID(ctx context.Context, roleID string) error {
	c := entity.GetRoleMenuCollection(ctx, a.Client)
	err := DeleteMany(ctx, c, DefaultFilter(ctx, Filter("role_id", roleID)))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
