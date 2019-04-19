package gorm

import (
	"time"

	"github.com/LyricTian/gin-admin/pkg/logger"
	"github.com/LyricTian/gin-admin/pkg/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var tableName string

// Config 配置参数
type Config struct {
	DBType       string
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	TableName    string
}

// New 创建基于gorm的钩子实例(需要指定表名)
func New(c *Config) *Hook {
	tableName = c.TableName

	db, err := gorm.Open(c.DBType, c.DSN)
	if err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)

	db.AutoMigrate(new(LogItem))
	return &Hook{
		db: db,
	}
}

// Hook gorm日志钩子
type Hook struct {
	db *gorm.DB
}

// Exec 执行日志写入
func (h *Hook) Exec(entry *logrus.Entry) error {
	item := &LogItem{
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
	if v, ok := data[logger.SpanIDKey]; ok {
		item.SpanID, _ = v.(string)
		delete(data, logger.SpanIDKey)
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
	if v, ok := data[logger.TimeConsumingKey]; ok {
		item.TimeConsuming, _ = v.(int64)
		delete(data, logger.TimeConsumingKey)
	}

	if len(data) > 0 {
		item.Data = util.JSONMarshalToString(data)
	}

	result := h.db.Create(item)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

// Close 关闭钩子
func (h *Hook) Close() error {
	return h.db.Close()
}

// LogItem 存储日志项
type LogItem struct {
	ID            uint      `gorm:"column:id;primary_key;auto_increment;"` // id
	Level         string    `gorm:"column:level;size:20;index;"`           // 日志级别
	Message       string    `gorm:"column:message;size:1024;"`             // 消息
	TraceID       string    `gorm:"column:trace_id;size:128;index;"`       // 跟踪ID
	UserID        string    `gorm:"column:user_id;size:36;index;"`         // 用户ID
	SpanID        string    `gorm:"column:span_id;size:128;"`              // 跟踪单元ID
	SpanTitle     string    `gorm:"column:span_title;size:256;"`           // 跟踪单元标题
	SpanFunction  string    `gorm:"column:span_function;size:256;"`        // 跟踪单元函数名
	Data          string    `gorm:"column:data;type:text;"`                // 日志数据(json)
	TimeConsuming int64     `gorm:"column:time_consuming;index;"`          // 耗时(单位：微妙)
	Version       string    `gorm:"column:version;index;size:32;"`         // 服务版本号
	CreatedAt     time.Time `gorm:"column:created_at;index"`               // 创建时间
}

// TableName 表名
func (LogItem) TableName() string {
	return tableName
}
