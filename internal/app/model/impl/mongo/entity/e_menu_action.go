package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetMenuActionCollection 获取MenuAction存储
func GetMenuActionCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return getCollection(ctx, cli, MenuAction{})
}

// SchemaMenuAction 菜单动作
type SchemaMenuAction schema.MenuAction

// ToMenuAction 转换为菜单动作实体
func (a SchemaMenuAction) ToMenuAction() *MenuAction {
	item := new(MenuAction)
	util.StructMapToStruct(a, item)
	return item
}

// MenuAction 菜单动作实体
type MenuAction struct {
	Model  `bson:",inline"`
	MenuID string `bson:"menu_id"` // 菜单ID
	Code   string `bson:"code"`    // 动作编号
	Name   string `bson:"name"`    // 动作名称
}

func (a MenuAction) String() string {
	return toString(a)
}

// CollectionName 集合名
func (a MenuAction) CollectionName() string {
	return a.Model.CollectionName("menu_action")
}

// CreateIndexes 创建索引
func (a MenuAction) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"menu_id": 1}},
	})
}

// ToSchemaMenuAction 转换为菜单动作对象
func (a MenuAction) ToSchemaMenuAction() *schema.MenuAction {
	item := new(schema.MenuAction)
	util.StructMapToStruct(a, item)
	return item
}

// MenuActions 菜单动作列表
type MenuActions []*MenuAction

// ToSchemaMenuActions 转换为菜单动作对象列表
func (a MenuActions) ToSchemaMenuActions() []*schema.MenuAction {
	list := make([]*schema.MenuAction, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuAction()
	}
	return list
}
