package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetRoleCollection 获取Role存储
func GetRoleCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return getCollection(ctx, cli, Role{})
}

// SchemaRole 角色对象
type SchemaRole schema.Role

// ToRole 转换为角色实体
func (a SchemaRole) ToRole() *Role {
	item := new(Role)
	util.StructMapToStruct(a, item)
	return item
}

// Role 角色实体
type Role struct {
	Model    `bson:",inline"`
	Name     string `bson:"name"`     // 角色名称
	Sequence int    `bson:"sequence"` // 排序值
	Memo     string `bson:"memo"`     // 备注
	Status   int    `bson:"status"`   // 状态(1:启用 2:禁用)
	Creator  string `bson:"creator"`  // 创建者
}

func (a Role) String() string {
	return toString(a)
}

// CollectionName 集合名
func (a Role) CollectionName() string {
	return a.Model.CollectionName("role")
}

// CreateIndexes 创建索引
func (a Role) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"name": 1}},
		{Keys: bson.M{"sequence": -1}},
		{Keys: bson.M{"status": 1}},
	})
}

// ToSchemaRole 转换为角色对象
func (a Role) ToSchemaRole() *schema.Role {
	item := new(schema.Role)
	util.StructMapToStruct(a, item)
	return item
}

// Roles 角色实体列表
type Roles []*Role

// ToSchemaRoles 转换为角色对象列表
func (a Roles) ToSchemaRoles() []*schema.Role {
	list := make([]*schema.Role, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRole()
	}
	return list
}
