package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUserRoleCollection 获取UserRole存储
func GetUserRoleCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return getCollection(ctx, cli, UserRole{})
}

// SchemaUserRole 用户角色
type SchemaUserRole schema.UserRole

// ToUserRole 转换为角色菜单实体
func (a SchemaUserRole) ToUserRole() *UserRole {
	item := new(UserRole)
	util.StructMapToStruct(a, item)
	return item
}

// UserRole 用户角色关联实体
type UserRole struct {
	Model  `bson:",inline"`
	UserID string `bson:"user_id"` // 用户内码
	RoleID string `bson:"role_id"` // 角色内码
}

// CollectionName 集合名
func (a UserRole) CollectionName() string {
	return a.Model.CollectionName("user_role")
}

// CreateIndexes 创建索引
func (a UserRole) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"user_id": 1}},
		{Keys: bson.M{"role_id": 1}},
	})
}

// ToSchemaUserRole 转换为用户角色对象
func (a UserRole) ToSchemaUserRole() *schema.UserRole {
	item := new(schema.UserRole)
	util.StructMapToStruct(a, item)
	return item
}

// UserRoles 用户角色关联列表
type UserRoles []*UserRole

// ToSchemaUserRoles 转换为用户角色对象列表
func (a UserRoles) ToSchemaUserRoles() []*schema.UserRole {
	list := make([]*schema.UserRole, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUserRole()
	}
	return list
}
