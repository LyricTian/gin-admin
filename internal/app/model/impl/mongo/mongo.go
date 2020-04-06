package mongo

import (
	"context"
	"time"

	"github.com/LyricTian/gin-admin/internal/app/model/impl/mongo/entity"
	"github.com/LyricTian/gin-admin/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config 配置参数
type Config struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// NewClient 创建mongo客户端实例
func NewClient(cfg *Config) (*mongo.Client, func(), error) {
	var (
		ctx    = context.Background()
		cancel context.CancelFunc
	)

	if t := cfg.Timeout; t > 0 {
		ctx, cancel = context.WithTimeout(ctx, t)
		defer cancel()
	}

	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, nil, err
	}

	cleanFunc := func() {
		err := cli.Disconnect(context.Background())
		if err != nil {
			logger.Errorf(context.Background(), "Mongo disconnect error: %s", err.Error())
		}
	}

	err = cli.Ping(context.Background(), nil)
	if err != nil {
		return nil, cleanFunc, err
	}
	return cli, cleanFunc, nil
}

// CreateIndexes 创建索引
func CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return createIndexes(
		ctx,
		cli,
		new(entity.Demo),
		new(entity.MenuAction),
		new(entity.MenuActionResource),
		new(entity.Menu),
		new(entity.RoleMenu),
		new(entity.Role),
		new(entity.UserRole),
		new(entity.User),
	)
}

type indexer interface {
	CreateIndexes(ctx context.Context, cli *mongo.Client) error
}

func createIndexes(ctx context.Context, cli *mongo.Client, indexes ...indexer) error {
	for _, idx := range indexes {
		err := idx.CreateIndexes(ctx, cli)
		if err != nil {
			return err
		}
	}
	return nil
}
