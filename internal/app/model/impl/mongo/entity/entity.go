package entity

import (
	"context"
	"fmt"
	"time"

	"github.com/LyricTian/gin-admin/internal/app/config"
	"github.com/LyricTian/gin-admin/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Model base model
type Model struct {
	RecordID  string     `bson:"_id"`
	CreatedAt time.Time  `bson:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty"`
}

// CollectionName collection name
func (Model) CollectionName(name string) string {
	return fmt.Sprintf("%s%s", config.C.Mongo.CollectionPrefix, name)
}

// CreateIndexes 创建索引
func (Model) CreateIndexes(ctx context.Context, cli *mongo.Client, m collectioner, indexes []mongo.IndexModel) error {
	models := []mongo.IndexModel{
		{Keys: bson.M{"created_at": 1}},
		{Keys: bson.M{"updated_at": 1}},
		{Keys: bson.M{"deleted_at": 1}},
	}
	if len(indexes) > 0 {
		models = append(models, indexes...)
	}
	_, err := getCollection(ctx, cli, m).Indexes().CreateMany(ctx, models)
	return err
}

func toString(v interface{}) string {
	return util.JSONMarshalToString(v)
}

type collectioner interface {
	CollectionName() string
}

func getCollection(ctx context.Context, cli *mongo.Client, m collectioner) *mongo.Collection {
	return cli.Database(config.C.Mongo.Database).Collection(m.CollectionName())
}
