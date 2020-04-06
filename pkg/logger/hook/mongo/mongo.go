package bson

import (
	"context"
	"time"

	"github.com/LyricTian/gin-admin/pkg/logger"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config 配置参数
type Config struct {
	URI        string
	Database   string
	Collection string
	Timeout    time.Duration
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// New 创建基于bson的钩子实例(需要指定表名)
func New(cfg *Config) *Hook {
	var (
		ctx    = context.Background()
		cancel context.CancelFunc
	)

	if t := cfg.Timeout; t > 0 {
		ctx, cancel = context.WithTimeout(ctx, t)
		defer cancel()
	}

	cli, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	handleError(err)
	c := cli.Database(cfg.Database).Collection(cfg.Collection)

	return &Hook{
		Client:     cli,
		Collection: c,
	}
}

// Hook bson日志钩子
type Hook struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

// Exec 执行日志写入
func (h *Hook) Exec(entry *logrus.Entry) error {
	item := &LogItem{
		ID:        util.NewRecordID(),
		Level:     entry.Level.String(),
		Message:   entry.Message,
		CreatedAt: entry.Time,
	}

	data := entry.Data
	if v, ok := data[logger.TraceIDKey]; ok {
		item.TraceID, _ = v.(string)
		delete(data, logger.TraceIDKey)
	}
	if v, ok := data[logger.UserIDKey]; ok {
		item.UserID, _ = v.(string)
		delete(data, logger.UserIDKey)
	}
	if v, ok := data[logger.SpanTitleKey]; ok {
		item.SpanTitle, _ = v.(string)
		delete(data, logger.SpanTitleKey)
	}
	if v, ok := data[logger.SpanFunctionKey]; ok {
		item.SpanFunction, _ = v.(string)
		delete(data, logger.SpanFunctionKey)
	}
	if v, ok := data[logger.VersionKey]; ok {
		item.Version, _ = v.(string)
		delete(data, logger.VersionKey)
	}

	if len(data) > 0 {
		item.Data = util.JSONMarshalToString(data)
	}

	_, err := h.Collection.InsertOne(context.Background(), item)
	return err
}

// Close 关闭钩子
func (h *Hook) Close() error {
	return h.Client.Disconnect(context.Background())
}

// LogItem 存储日志项
type LogItem struct {
	ID           string    `bson:"_id"`           // id
	Level        string    `bson:"level"`         // 日志级别
	Message      string    `bson:"message"`       // 消息
	TraceID      string    `bson:"trace_id"`      // 跟踪ID
	UserID       string    `bson:"user_id"`       // 用户ID
	SpanTitle    string    `bson:"span_title"`    // 跟踪单元标题
	SpanFunction string    `bson:"span_function"` // 跟踪单元函数名
	Data         string    `bson:"data"`          // 日志数据(json)
	Version      string    `bson:"version"`       // 服务版本号
	CreatedAt    time.Time `bson:"created_at"`    // 创建时间
}
