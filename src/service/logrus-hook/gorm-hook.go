package logrushook

import (
	"time"

	"github.com/spf13/viper"

	"github.com/LyricTian/gin-admin/src/logger"
	"github.com/LyricTian/gin-admin/src/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// NewGormHook 创建gorm存储的钩子
func NewGormHook(db *gorm.DB, opts ...Option) *Hook {
	db.AutoMigrate(&GormItem{})
	return New(&gormHook{db}, opts...)
}

type gormHook struct {
	db *gorm.DB
}

func (h *gormHook) Exec(entry *logrus.Entry) error {
	item := &GormItem{
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
		buf := util.JSONMarshalToString(data)
		if len(buf) > 4096 {
			buf = string(buf[:4096])
		}
		item.Data = buf
	}

	result := h.db.Create(item)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (h *gormHook) Close() error {
	return nil
}

// GormItem gorm日志项
type GormItem struct {
	ID            uint      `gorm:"column:id;primary_key;auto_increment;"` // id
	Level         string    `gorm:"column:level;size:20;index;"`           // 日志级别
	Message       string    `gorm:"column:message;size:1024;"`             // 消息
	TraceID       string    `gorm:"column:trace_id;size:128;index;"`       // 跟踪ID
	UserID        string    `gorm:"column:user_id;size:36;index;"`         // 用户ID
	SpanID        string    `gorm:"column:span_id;size:128;"`              // 跟踪单元ID
	SpanTitle     string    `gorm:"column:span_title;size:256;"`           // 跟踪单元标题
	SpanFunction  string    `gorm:"column:span_function;size:256;"`        // 跟踪单元函数名
	Data          string    `gorm:"column:data;size:4096;"`                // 日志数据(json)
	TimeConsuming int64     `gorm:"column:time_consuming;index;"`          // 耗时(单位：毫秒)
	Version       string    `gorm:"column:version;size:32;"`               // 服务版本号
	CreatedAt     time.Time `gorm:"column:created_at;index"`               // 创建时间
}

// TableName 表名
func (GormItem) TableName() string {
	return viper.GetString("db_table_prefix") + "logger"
}
