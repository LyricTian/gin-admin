package gorm

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/LyricTian/gin-admin/v8/pkg/logger"
)

// Create logger hook from gorm
func New(db *gorm.DB) *Hook {
	db.AutoMigrate(new(Logger))

	return &Hook{
		db: db,
	}
}

// Grom Logger Hook
type Hook struct {
	db *gorm.DB
}

func (h *Hook) Exec(entry *logrus.Entry) error {
	item := &Logger{
		Level:     entry.Level.String(),
		Message:   entry.Message,
		CreatedAt: entry.Time,
	}

	ctx := entry.Context
	item.TraceID = logger.FromTraceIDContext(ctx)
	item.UserID = logger.FromUserIDContext(ctx)
	item.UserName = logger.FromUserNameContext(ctx)
	item.Tag = logger.FromTagContext(ctx)

	data := entry.Data
	if v := logger.FromStackContext(ctx); v != nil {
		data["stack"] = fmt.Sprintf("%+v", v)
	}

	if len(data) > 0 {
		b, _ := json.Marshal(data)
		item.Data = string(b)
	}

	return h.db.Create(item).Error
}

func (h *Hook) Close() error {
	db, err := h.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

type Logger struct {
	ID        uint      `gorm:"primaryKey;"`     // id
	Level     string    `gorm:"size:20;index;"`  // 日志级别
	TraceID   string    `gorm:"size:128;index;"` // 跟踪ID
	UserID    uint64    `gorm:"index;"`          // 用户ID
	UserName  string    `gorm:"size:64;index;"`  // 用户名
	Tag       string    `gorm:"size:128;index;"` // Tag
	Message   string    `gorm:"size:1024;"`      // 消息
	Data      string    `gorm:"type:text;"`      // 日志数据(json)
	CreatedAt time.Time `gorm:"index"`           // 创建时间
}
