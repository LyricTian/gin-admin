package schema

import (
	"time"

	"github.com/LyricTian/gin-admin/v10/pkg/util"
)

const (
	RoleStatusEnabled  = "enabled"  // Enabled
	RoleStatusDisabled = "disabled" // Disabled
)

// Role management for RBAC
type Role struct {
	ID          string    `json:"id" gorm:"size:20;primarykey;"` // Unique ID
	Code        string    `json:"code" gorm:"size:32;index;"`    // Code of role (unique)
	Name        string    `json:"name" gorm:"size:128;index"`    // Display name of role
	Description string    `json:"description" gorm:"size:1024"`  // Details about role
	Sequence    int       `json:"sequence"`                      // Sequence for sorting
	Status      string    `json:"status" gorm:"size:20;index"`   // Status of role (disabled, enabled)
	CreatedAt   time.Time `json:"created_at" gorm:"index;"`      // Create time
	UpdatedAt   time.Time `json:"updated_at" gorm:"index;"`      // Update time
	Menus       RoleMenus `json:"menus" gorm:"-"`                // Role menu list
}

// Defining the query parameters for the `Role` struct.
type RoleQueryParam struct {
	util.PaginationParam
	LikeName    string     `form:"name"`                                       // Display name of role
	Status      string     `form:"status" binding:"oneof=disabled enabled ''"` // Status of role (disabled, enabled)
	InIDs       []string   `form:"-"`                                          // ID list
	GtUpdatedAt *time.Time `form:"-"`                                          // Update time is greater than
}

// Defining the query options for the `Role` struct.
type RoleQueryOptions struct {
	util.QueryOptions
}

// Defining the query result for the `Role` struct.
type RoleQueryResult struct {
	Data       Roles
	PageResult *util.PaginationResult
}

// Defining the slice of `Role` struct.
type Roles []*Role

// Defining the data structure for creating a `Role` struct.
type RoleForm struct {
	Code        string    `json:"code" binding:"required,max=32"`                   // Code of role (unique)
	Name        string    `json:"name" binding:"required,max=128"`                  // Display name of role
	Description string    `json:"description"`                                      // Details about role
	Sequence    int       `json:"sequence"`                                         // Sequence for sorting
	Status      string    `json:"status" binding:"required,oneof=disabled enabled"` // Status of role (disabled, enabled)
	Menus       RoleMenus `json:"menus"`                                            // Role menu list
}

// A validation function for the `RoleForm` struct.
func (a *RoleForm) Validate() error {
	return nil
}

func (a *RoleForm) FillTo(role *Role) error {
	role.Code = a.Code
	role.Name = a.Name
	role.Description = a.Description
	role.Sequence = a.Sequence
	role.Status = a.Status
	return nil
}
