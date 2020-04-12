package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetDemoCollection 获取demo存储
func GetDemoCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return getCollection(ctx, cli, Demo{})
}

// SchemaDemo demo对象
type SchemaDemo schema.Demo

// ToDemo 转换为demo实体
func (a SchemaDemo) ToDemo() *Demo {
	item := new(Demo)
	util.StructMapToStruct(a, item)
	return item
}

// Demo demo实体
type Demo struct {
	Model   `bson:",inline"`
	Code    string `bson:"code"`    // 编号
	Name    string `bson:"name"`    // 名称
	Memo    string `bson:"memo"`    // 备注
	Status  int    `bson:"status"`  // 状态(1:启用 2:停用)
	Creator string `bson:"creator"` // 创建者
}

func (a Demo) String() string {
	return toString(a)
}

// CollectionName 集合名
func (a Demo) CollectionName() string {
	return a.Model.CollectionName("demo")
}

// CreateIndexes 创建索引
func (a Demo) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"code": 1}},
		{Keys: bson.M{"name": 1}},
		{Keys: bson.M{"status": 1}},
	})
}

// ToSchemaDemo 转换为demo对象
func (a Demo) ToSchemaDemo() *schema.Demo {
	item := new(schema.Demo)
	util.StructMapToStruct(a, item)
	return item
}

// Demos demo列表
type Demos []*Demo

// ToSchemaDemos 转换为demo对象列表
func (a Demos) ToSchemaDemos() []*schema.Demo {
	list := make([]*schema.Demo, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaDemo()
	}
	return list
}
