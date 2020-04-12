package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetMenuActionResourceCollection 获取MenuActionResource存储
func GetMenuActionResourceCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return getCollection(ctx, cli, MenuActionResource{})
}

// SchemaMenuActionResource 菜单动作关联资源
type SchemaMenuActionResource schema.MenuActionResource

// ToMenuActionResource 转换为菜单动作关联资源实体
func (a SchemaMenuActionResource) ToMenuActionResource() *MenuActionResource {
	item := new(MenuActionResource)
	util.StructMapToStruct(a, item)
	return item
}

// MenuActionResource 菜单动作关联资源实体
type MenuActionResource struct {
	Model    `bson:",inline"`
	ActionID string `bson:"action_id"` // 菜单动作ID
	Method   string `bson:"method"`    // 资源请求方式(支持正则)
	Path     string `bson:"path"`      // 资源请求路径（支持/:id匹配）
}

func (a MenuActionResource) String() string {
	return toString(a)
}

// CollectionName 集合名
func (a MenuActionResource) CollectionName() string {
	return a.Model.CollectionName("menu_action_resource")
}

// CreateIndexes 创建索引
func (a MenuActionResource) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"action_id": 1}},
	})
}

// ToSchemaMenuActionResource 转换为菜单动作关联资源对象
func (a MenuActionResource) ToSchemaMenuActionResource() *schema.MenuActionResource {
	item := new(schema.MenuActionResource)
	util.StructMapToStruct(a, item)
	return item
}

// MenuActionResources 菜单动作关联资源列表
type MenuActionResources []*MenuActionResource

// ToSchemaMenuActionResources 转换为菜单动作关联资源对象列表
func (a MenuActionResources) ToSchemaMenuActionResources() []*schema.MenuActionResource {
	list := make([]*schema.MenuActionResource, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuActionResource()
	}
	return list
}
