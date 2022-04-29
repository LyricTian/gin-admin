package hook

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Logger struct {
	ID        string    `gorm:"size:20;primaryKey;"`
	Level     string    `gorm:"size:20;"`
	TraceID   string    `gorm:"size:128;"`
	UserID    string    `gorm:"size:20;index;"`
	Tag       string    `gorm:"size:32;"`
	Message   string    `gorm:"size:1024;"`
	Data      string    `gorm:"size:4096;"`
	CreatedAt time.Time `gorm:"index"`
}

// Create logger hook from gorm
func NewGormHook(db *gorm.DB) *GormHook {
	db.AutoMigrate(new(Logger))

	return &GormHook{
		db: db,
	}
}

// Gorm Logger Hook
type GormHook struct {
	db *gorm.DB
}

func (h *GormHook) Exec(entry *logrus.Entry) error {
	item := &Logger{
		ID:        xid.New().String(),
		Level:     entry.Level.String(),
		Message:   entry.Message,
		CreatedAt: entry.Time,
	}

	data := entry.Data
	if v, ok := data["trace_id"]; ok {
		item.TraceID = fmt.Sprintf("%v", v)
		delete(data, "trace_id")
	}
	if v, ok := data["user_id"]; ok {
		item.UserID = fmt.Sprintf("%v", v)
		delete(data, "user_id")
	}
	if v, ok := data["tag"]; ok {
		item.Tag = fmt.Sprintf("%v", v)
		delete(data, "tag")
	}

	if len(data) > 0 {
		b, _ := json.Marshal(data)
		item.Data = string(b)
	}

	return h.db.Create(item).Error
}

func (h *GormHook) Close() error {
	db, err := h.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
