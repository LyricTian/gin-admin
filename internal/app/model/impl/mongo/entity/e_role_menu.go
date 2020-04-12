package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetRoleMenuCollection 获取RoleMenu存储
func GetRoleMenuCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return getCollection(ctx, cli, RoleMenu{})
}

// SchemaRoleMenu 角色菜单
type SchemaRoleMenu schema.RoleMenu

// ToRoleMenu 转换为角色菜单实体
func (a SchemaRoleMenu) ToRoleMenu() *RoleMenu {
	item := new(RoleMenu)
	util.StructMapToStruct(a, item)
	return item
}

// RoleMenu 角色菜单实体
type RoleMenu struct {
	Model    `bson:",inline"`
	RoleID   string `bson:"role_id"`   // 角色ID
	MenuID   string `bson:"menu_id"`   // 菜单ID
	ActionID string `bson:"action_id"` // 动作ID
}

func (a RoleMenu) String() string {
	return toString(a)
}

// CollectionName 集合名
func (a RoleMenu) CollectionName() string {
	return a.Model.CollectionName("role_menu")
}

// CreateIndexes 创建索引
func (a RoleMenu) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"role_id": 1}},
		{Keys: bson.M{"menu_id": 1}},
		{Keys: bson.M{"action_id": 1}},
	})
}

// ToSchemaRoleMenu 转换为角色菜单对象
func (a RoleMenu) ToSchemaRoleMenu() *schema.RoleMenu {
	item := new(schema.RoleMenu)
	util.StructMapToStruct(a, item)
	return item
}

// RoleMenus 角色菜单列表
type RoleMenus []*RoleMenu

// ToSchemaRoleMenus 转换为角色菜单对象列表
func (a RoleMenus) ToSchemaRoleMenus() []*schema.RoleMenu {
	list := make([]*schema.RoleMenu, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRoleMenu()
	}
	return list
}
