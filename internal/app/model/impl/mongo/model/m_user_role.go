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

var _ model.IUserRole = (*UserRole)(nil)

// UserRoleSet 注入UserRole
var UserRoleSet = wire.NewSet(wire.Struct(new(UserRole), "*"), wire.Bind(new(model.IUserRole), new(*UserRole)))

// UserRole 用户角色存储
type UserRole struct {
	Client *mongo.Client
}

func (a *UserRole) getQueryOption(opts ...schema.UserRoleQueryOptions) schema.UserRoleQueryOptions {
	var opt schema.UserRoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *UserRole) Query(ctx context.Context, params schema.UserRoleQueryParam, opts ...schema.UserRoleQueryOptions) (*schema.UserRoleQueryResult, error) {
	opt := a.getQueryOption(opts...)

	c := entity.GetUserRoleCollection(ctx, a.Client)
	filter := DefaultFilter(ctx)
	userIDs := params.UserIDs
	if v := params.UserID; v != "" {
		userIDs = append(userIDs, v)
	}
	if v := userIDs; len(v) > 0 {
		filter = append(filter, Filter("user_id", bson.M{"$in": v}))
	}
	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByDESC))

	var list entity.UserRoles
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.UserRoleQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUserRoles(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *UserRole) Get(ctx context.Context, recordID string, opts ...schema.UserRoleQueryOptions) (*schema.UserRole, error) {
	c := entity.GetUserRoleCollection(ctx, a.Client)
	filter := DefaultFilter(ctx, Filter("_id", recordID))
	var item entity.UserRole
	ok, err := FindOne(ctx, c, filter, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaUserRole(), nil
}

// Create 创建数据
func (a *UserRole) Create(ctx context.Context, item schema.UserRole) error {
	eitem := entity.SchemaUserRole(item).ToUserRole()
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := entity.GetUserRoleCollection(ctx, a.Client)
	err := Insert(ctx, c, eitem)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *UserRole) Update(ctx context.Context, recordID string, item schema.UserRole) error {
	eitem := entity.SchemaUserRole(item).ToUserRole()
	eitem.UpdatedAt = time.Now()
	c := entity.GetUserRoleCollection(ctx, a.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)), eitem)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *UserRole) Delete(ctx context.Context, recordID string) error {
	c := entity.GetUserRoleCollection(ctx, a.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", recordID)))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeleteByUserID 根据用户ID删除数据
func (a *UserRole) DeleteByUserID(ctx context.Context, userID string) error {
	c := entity.GetUserRoleCollection(ctx, a.Client)
	err := DeleteMany(ctx, c, DefaultFilter(ctx, Filter("user_id", userID)))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
