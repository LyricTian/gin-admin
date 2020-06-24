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

var _ model.IUser = (*User)(nil)

// UserSet 注入User
var UserSet = wire.NewSet(wire.Struct(new(User), "*"), wire.Bind(new(model.IUser), new(*User)))

// User 用户存储
type User struct {
	Client *mongo.Client
}

func (a *User) getQueryOption(opts ...schema.UserQueryOptions) schema.UserQueryOptions {
	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	opt := a.getQueryOption(opts...)

	c := entity.GetUserCollection(ctx, a.Client)
	filter := DefaultFilter(ctx)
	if v := params.UserName; v != "" {
		filter = append(filter, Filter("user_name", v))
	}
	if v := params.QueryValue; v != "" {
		filter = append(filter, Filter("$or", bson.A{
			OrRegexFilter("user_name", v),
			OrRegexFilter("real_name", v),
			OrRegexFilter("phone", v),
			OrRegexFilter("email", v),
		}))
	}
	if v := params.RoleIDs; len(v) > 0 {
		result, err := entity.GetUserRoleCollection(ctx, a.Client).Distinct(ctx, "user_id", DefaultFilter(ctx, Filter("role_id", bson.M{"$in": v})))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		filter = append(filter, Filter("_id", bson.M{"$in": result}))
	}
	if v := params.Status; v > 0 {
		filter = append(filter, Filter("status", v))
	}
	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByDESC))

	var list entity.Users
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.UserQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUsers(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *User) Get(ctx context.Context, id string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	c := entity.GetUserCollection(ctx, a.Client)
	filter := DefaultFilter(ctx, Filter("_id", id))
	var item entity.User
	ok, err := FindOne(ctx, c, filter, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaUser(), nil
}

// Create 创建数据
func (a *User) Create(ctx context.Context, item schema.User) error {
	eitem := entity.SchemaUser(item).ToUser()
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := entity.GetUserCollection(ctx, a.Client)
	err := Insert(ctx, c, eitem)
	return errors.WithStack(err)
}

// Update 更新数据
func (a *User) Update(ctx context.Context, id string, item schema.User) error {
	eitem := entity.SchemaUser(item).ToUser()
	eitem.UpdatedAt = time.Now()
	c := entity.GetUserCollection(ctx, a.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", id)), eitem)
	return errors.WithStack(err)
}

// Delete 删除数据
func (a *User) Delete(ctx context.Context, id string) error {
	c := entity.GetUserCollection(ctx, a.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", id)))
	return errors.WithStack(err)
}

// UpdateStatus 更新状态
func (a *User) UpdateStatus(ctx context.Context, id string, status int) error {
	c := entity.GetUserCollection(ctx, a.Client)
	err := UpdateFields(ctx, c, DefaultFilter(ctx, Filter("_id", id)), bson.M{"status": status})
	return errors.WithStack(err)
}

// UpdatePassword 更新密码
func (a *User) UpdatePassword(ctx context.Context, id, password string) error {
	c := entity.GetUserCollection(ctx, a.Client)
	err := UpdateFields(ctx, c, DefaultFilter(ctx, Filter("_id", id)), bson.M{"password": password})
	return errors.WithStack(err)
}
