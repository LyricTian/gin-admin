package typed

import (
	"time"

	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
)

// Dictionary management for key/value pairs
type Dictionary struct {
	ID        string    `gorm:"size:20;primarykey;" json:"id"`
	Namespace string    `gorm:"size:64;index;" json:"namespace"` // Namespace of the dictionary
	Key       string    `gorm:"size:128;index;" json:"key"`      // Key of the dictionary
	Value     *string   `gorm:"size:4096;" json:"value"`         // Value of the key
	Remark    *string   `gorm:"size:1024;" json:"remark"`        // Remark of the key
	CreatedAt time.Time `gorm:"index;" json:"created_at"`
	CreatedBy string    `gorm:"size:20;" json:"created_by"`
	UpdatedAt time.Time `gorm:"index;" json:"updated_at"`
	UpdatedBy string    `gorm:"size:20;" json:"updated_by"`
}

type DictionaryQueryParam struct {
	utilx.PaginationParam
	Namespace string `form:"namespace"`
	Key       string `form:"key"`
}

type DictionaryQueryOptions struct {
	utilx.QueryOptions
}

type DictionaryQueryResult struct {
	Data       Dictionaries
	PageResult *utilx.PaginationResult
}

type Dictionaries []*Dictionary

type DictionaryCreate struct {
	Namespace string `json:"namespace" binding:"required"`
	Key       string `json:"key" binding:"required"`
	Value     string `json:"value"`
	Remark    string `json:"remark"`
}
