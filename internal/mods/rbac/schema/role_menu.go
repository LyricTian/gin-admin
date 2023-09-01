package schema

import (
	"time"

	"github.com/LyricTian/gin-admin/v10/pkg/util"
)

// Role permissions for RBAC
type RoleMenu struct {
	ID        string    `json:"id" gorm:"size:20;primarykey"` // Unique ID
	RoleID    string    `json:"role_id" gorm:"size:20;index"` // From Role.ID
	MenuID    string    `json:"menu_id" gorm:"size:20;index"` // From Menu.ID
	CreatedAt time.Time `json:"created_at" gorm:"index;"`     // Create time
	UpdatedAt time.Time `json:"updated_at" gorm:"index;"`     // Update time
}

// Defining the query parameters for the `RoleMenu` struct.
type RoleMenuQueryParam struct {
	util.PaginationParam
	RoleID string `form:"-"` // From Role.ID
}

// Defining the query options for the `RoleMenu` struct.
type RoleMenuQueryOptions struct {
	util.QueryOptions
}

// Defining the query result for the `RoleMenu` struct.
type RoleMenuQueryResult struct {
	Data       RoleMenus
	PageResult *util.PaginationResult
}

// Defining the slice of `RoleMenu` struct.
type RoleMenus []*RoleMenu

// Defining the data structure for creating a `RoleMenu` struct.
type RoleMenuForm struct {
}

// A validation function for the `RoleMenuForm` struct.
func (a *RoleMenuForm) Validate() error {
	return nil
}

func (a *RoleMenuForm) FillTo(roleMenu *RoleMenu) error {
	return nil
}
