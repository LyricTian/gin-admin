package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetMenuCollection 获取Menu存储
func GetMenuCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return getCollection(ctx, cli, Menu{})
}

// SchemaMenu 菜单对象
type SchemaMenu schema.Menu

// ToMenu 转换为菜单实体
func (a SchemaMenu) ToMenu() *Menu {
	item := new(Menu)
	util.StructMapToStruct(a, item)
	return item
}

// Menu 菜单实体
type Menu struct {
	Model      `bson:",inline"`
	Name       string `bson:"name"`        // 菜单名称
	Sequence   int    `bson:"sequence"`    // 排序值
	Icon       string `bson:"icon"`        // 菜单图标
	Router     string `bson:"router"`      // 访问路由
	ParentID   string `bson:"parent_id"`   // 父级内码
	ParentPath string `bson:"parent_path"` // 父级路径
	ShowStatus int    `bson:"show_status"` // 状态(1:显示 2:隐藏)
	Status     int    `bson:"status"`      // 状态(1:启用 2:禁用)
	Memo       string `bson:"memo"`        // 备注
	Creator    string `bson:"creator"`     // 创建人
}

func (a Menu) String() string {
	return toString(a)
}

// CollectionName 集合名
func (a Menu) CollectionName() string {
	return a.Model.CollectionName("menu")
}

// CreateIndexes 创建索引
func (a Menu) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"name": 1}},
		{Keys: bson.M{"sequence": -1}},
		{Keys: bson.M{"parent_id": 1}},
		{Keys: bson.M{"parent_path": 1}},
		{Keys: bson.M{"show_status": 1}},
		{Keys: bson.M{"status": 1}},
	})
}

// ToSchemaMenu 转换为菜单对象
func (a Menu) ToSchemaMenu() *schema.Menu {
	item := new(schema.Menu)
	util.StructMapToStruct(a, item)
	return item
}

// Menus 菜单实体列表
type Menus []*Menu

// ToSchemaMenus 转换为菜单对象列表
func (a Menus) ToSchemaMenus() []*schema.Menu {
	list := make([]*schema.Menu, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenu()
	}
	return list
}
