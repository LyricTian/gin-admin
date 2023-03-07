package schema

import (
	"time"

	"github.com/LyricTian/gin-admin/v10/internal/library/utilx"
	"github.com/go-playground/validator/v10"
)

type Resource struct {
	ID          string    `gorm:"size:20;primarykey;" json:"id"` // Unique ID
	Code        string    `gorm:"size:128;index;" json:"code"`   // Unique code (format: module.resource.action)
	Object      string    `gorm:"size:128;" json:"object"`       // Resource object
	Action      string    `gorm:"size:128;" json:"action"`       // Resource action
	Description string    `gorm:"size:256;" json:"description"`  // Description
	Status      string    `gorm:"size:20;index;" json:"status"`  // Status (enabled/disabled)
	CreatedAt   time.Time `gorm:"index;" json:"created_at"`      // Create time
	UpdatedAt   time.Time `gorm:"index;" json:"updated_at"`      // Update time
}

type ResourceQueryParam struct {
	utilx.PaginationParam
	LikeCode string `form:"code"`                                           // Unique code
	Status   string `form:"status" binding:"oneof='enabled' 'disabled' ''"` // Status (enabled/disabled)
}

type ResourceQueryOptions struct {
	utilx.QueryOptions
}

type ResourceQueryResult struct {
	Data       Resources
	PageResult *utilx.PaginationResult
}

type Resources []*Resource

type ResourceCreate struct {
	Code        string `json:"code"`                                                 // Unique code (format: module.resource.action)
	Object      string `json:"object" binding:"required"`                            // Resource object
	Action      string `json:"action" binding:"required"`                            // Resource action
	Description string `json:"description"`                                          // Description
	Status      string `json:"status" binding:"required,oneof='enabled' 'disabled'"` // Status (enabled/disabled)
}

func (a ResourceCreate) Validate() error {
	v := validator.New()
	v.SetTagName("binding")
	return v.Struct(a)
}
