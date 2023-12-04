package schema

import (
	"time"

	"github.com/LyricTian/gin-admin/v10/internal/config"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
)

// Logger management
type Logger struct {
	ID        string    `gorm:"size:20;primaryKey;" json:"id"`           // Unique ID
	Level     string    `gorm:"size:20;index;" json:"level"`             // Log level
	TraceID   string    `gorm:"size:64;index;" json:"trace_id"`          // Trace ID
	UserID    string    `gorm:"size:20;index;" json:"user_id"`           // User ID
	Tag       string    `gorm:"size:32;index;" json:"tag"`               // Log tag
	Message   string    `gorm:"size:1024;" json:"message"`               // Log message
	Stack     string    `gorm:"type:text;" json:"stack"`                 // Error stack
	Data      string    `gorm:"type:text;" json:"data"`                  // Log data
	CreatedAt time.Time `gorm:"index;" json:"created_at"`                // Create time
	LoginName string    `json:"login_name" gorm:"<-:false;-:migration;"` // From User.Username
	UserName  string    `json:"user_name" gorm:"<-:false;-:migration;"`  // From User.Name
}

func (a *Logger) TableName() string {
	return config.C.FormatTableName("logger")
}

// Defining the query parameters for the `Logger` struct.
type LoggerQueryParam struct {
	util.PaginationParam
	Level        string `form:"level"`     // Log level
	TraceID      string `form:"traceID"`   // Trace ID
	LikeUserName string `form:"userName"`  // User Name
	Tag          string `form:"tag"`       // Log tag
	LikeMessage  string `form:"message"`   // Log message
	StartTime    string `form:"startTime"` // Start time
	EndTime      string `form:"endTime"`   // End time
}

// Defining the query options for the `Logger` struct.
type LoggerQueryOptions struct {
	util.QueryOptions
}

// Defining the query result for the `Logger` struct.
type LoggerQueryResult struct {
	Data       Loggers
	PageResult *util.PaginationResult
}

// Defining the slice of `Logger` struct.
type Loggers []*Logger
