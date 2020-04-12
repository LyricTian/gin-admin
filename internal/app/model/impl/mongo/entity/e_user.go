package entity

import (
	"context"

	"github.com/LyricTian/gin-admin/internal/app/schema"
	"github.com/LyricTian/gin-admin/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUserCollection 获取User存储
func GetUserCollection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return getCollection(ctx, cli, User{})
}

// SchemaUser 用户对象
type SchemaUser schema.User

// ToUser 转换为用户实体
func (a SchemaUser) ToUser() *User {
	item := new(User)
	util.StructMapToStruct(a, item)
	return item
}

// User 用户实体
type User struct {
	Model    `bson:",inline"`
	UserName string `bson:"user_name"` // 用户名
	RealName string `bson:"real_name"` // 真实姓名
	Password string `bson:"password"`  // 密码(sha1(md5(明文))加密)
	Email    string `bson:"email"`     // 邮箱
	Phone    string `bson:"phone"`     // 手机号
	Status   int    `bson:"status"`    // 状态(1:启用 2:停用)
	Creator  string `bson:"creator"`   // 创建者
}

func (a User) String() string {
	return toString(a)
}

// CollectionName 集合名
func (a User) CollectionName() string {
	return a.Model.CollectionName("user")
}

// CreateIndexes 创建索引
func (a User) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"user_name": 1}},
		{Keys: bson.M{"real_name": 1}},
		{Keys: bson.M{"status": 1}},
	})
}

// ToSchemaUser 转换为用户对象
func (a User) ToSchemaUser() *schema.User {
	item := new(schema.User)
	util.StructMapToStruct(a, item)
	return item
}

// Users 用户实体列表
type Users []*User

// ToSchemaUsers 转换为用户对象列表
func (a Users) ToSchemaUsers() []*schema.User {
	list := make([]*schema.User, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUser()
	}
	return list
}
