package logger

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

type Logger struct {
	ID        string    `gorm:"size:20;primaryKey;"`
	Level     string    `gorm:"size:20;"`
	TraceID   string    `gorm:"size:64;"`
	UserID    string    `gorm:"size:20;"`
	Tag       string    `gorm:"size:32;"`
	Message   string    `gorm:"size:1024;"`
	Data      string    `gorm:"size:1048576;"`
	CreatedAt time.Time `gorm:"index"`
}

func NewGormHook(db *gorm.DB) *GormHook {
	_ = db.AutoMigrate(new(Logger))

	return &GormHook{
		db: db,
	}
}

// Gorm Logger Hook
type GormHook struct {
	db *gorm.DB
}

func (h *GormHook) Exec(extra map[string]string, b []byte) error {
	msg := &Logger{}
	data := make(map[string]interface{})
	err := jsoniter.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	if v, ok := data["ts"]; ok {
		msg.CreatedAt = time.UnixMilli(int64(v.(float64)))
		delete(data, "ts")
	}
	if v, ok := data["msg"]; ok {
		msg.Message = v.(string)
		delete(data, "msg")
	}
	if v, ok := data["logger"]; ok {
		msg.Tag = v.(string)
		delete(data, "logger")
	}
	if v, ok := data["trace_id"]; ok {
		msg.TraceID = v.(string)
		delete(data, "trace_id")
	}
	if v, ok := data["user_id"]; ok {
		msg.UserID = v.(string)
		delete(data, "user_id")
	}
	if v, ok := data["level"]; ok {
		msg.Level = v.(string)
		delete(data, "level")
	}

	for k, v := range extra {
		data[k] = v
	}

	if len(data) > 0 {
		buf, _ := jsoniter.Marshal(data)
		msg.Data = string(buf)
	}

	return h.db.Create(msg).Error
}

func (h *GormHook) Close() error {
	db, err := h.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
