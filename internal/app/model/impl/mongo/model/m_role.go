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

var _ model.IRole = (*Role)(nil)

// RoleSet 注入Role
var RoleSet = wire.NewSet(wire.Struct(new(Role), "*"), wire.Bind(new(model.IRole), new(*Role)))

// Role 角色存储
type Role struct {
	Client *mongo.Client
}

func (a *Role) getQueryOption(opts ...schema.RoleQueryOptions) schema.RoleQueryOptions {
	var opt schema.RoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *Role) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	opt := a.getQueryOption(opts...)

	c := entity.GetRoleCollection(ctx, a.Client)
	filter := DefaultFilter(ctx)

	if v := params.RecordIDs; len(v) > 0 {
		filter = append(filter, Filter("_id", bson.M{"$in": v}))
	}
	if v := params.Name; v != "" {
		filter = append(filter, Filter("name", v))
	}
	if v := params.UserID; v != "" {
		result, err := entity.GetUserRoleCollection(ctx, a.Client).Distinct(ctx, "role_id", DefaultFilter(ctx, Filter("user_id", v)))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		filter = append(filter, Filter("_id", bson.M{"$in": result}))
	}
	if v := params.QueryValue; v != "" {
		filter = append(filter, Filter("$or", bson.A{
			OrRegexFilter("name", v),
			OrRegexFilter("memo", v),
		}))
	}
	if v := params.Status; v > 0 {
		filter = append(filter, Filter("status", v))
	}
	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByDESC))

	var list entity.Roles
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RoleQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRoles(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *Role) Get(ctx context.Context, recordID string, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	c := entity.GetRoleCollection(ctx, a.Client)
	filter := DefaultFilter(ctx, Filter("_id", recordID))
	var item entity.Role
	ok, err := FindOne(ctx, c, filter, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaRole(), nil
}

// Create 创建数据
func (a *Role) Create(ctx context.Context, item schema.Role) error {
	eitem := entity.SchemaRole(item).ToRole()
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := entity.GetRoleCollection(ctx, a.Client)
	err := Insert(ctx, c, eitem)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *Role) Update(ctx context.Context, recordID string, item schema.Role) error {
	eitem := entity.SchemaRole(item).ToRole()
	eitem.UpdatedAt = time.Now()
	c := entity.GetRoleCollection(ctx, a.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)), eitem)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *Role) Delete(ctx context.Context, recordID string) error {
	c := entity.GetRoleCollection(ctx, a.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Role) UpdateStatus(ctx context.Context, recordID string, status int) error {
	c := entity.GetRoleCollection(ctx, a.Client)
	err := UpdateFields(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)), bson.M{"status": status})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
