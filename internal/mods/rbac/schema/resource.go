package schema

import (
	"time"

	"github.com/LyricTian/gin-admin/v10/internal/library/utils"
)

// Defining the `Resource` struct.
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

// Defining the name of the database table that corresponds to the `Resource` struct.
func (a Resource) TableName() string {
	return "resource"
}

// Defining the query parameters for the `Resource` struct.
type ResourceQueryParam struct {
	utils.PaginationParam
	LikeCode string `form:"code"`                                           // Unique code
	Status   string `form:"status" binding:"oneof='enabled' 'disabled' ''"` // Status (enabled/disabled)
}

// Defining the query options for the `Resource` struct.
type ResourceQueryOptions struct {
	utils.QueryOptions
}

// Defining the query result for the `Resource` struct.
type ResourceQueryResult struct {
	Data       Resources
	PageResult *utils.PaginationResult
}

// Defining the slice of `Resource` struct.
type Resources []*Resource

// Defining the data structure for creating a `Resource` struct.
type ResourceSave struct {
	Code        string `json:"code"`                                                 // Unique code (format: module.resource.action)
	Object      string `json:"object" binding:"required"`                            // Resource object
	Action      string `json:"action" binding:"required"`                            // Resource action
	Description string `json:"description"`                                          // Description
	Status      string `json:"status" binding:"required,oneof='enabled' 'disabled'"` // Status (enabled/disabled)
}

// A validation function for the `ResourceSave` struct.
func (a ResourceSave) Validate() error {
	return nil
}
