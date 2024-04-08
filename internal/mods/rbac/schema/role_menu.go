package schema

import (
	"time"

	"github.com/LyricTian/gin-admin/v10/internal/config"
	"github.com/LyricTian/gin-admin/v10/pkg/util"
)

// RoleMenu Role permissions for RBAC
type RoleMenu struct {
	ID        string    `json:"id" gorm:"size:20;primarykey"` // Unique ID
	RoleID    string    `json:"role_id" gorm:"size:20;index"` // From Role.ID
	MenuID    string    `json:"menu_id" gorm:"size:20;index"` // From Menu.ID
	CreatedAt time.Time `json:"created_at" gorm:"index;"`     // Create time
	UpdatedAt time.Time `json:"updated_at" gorm:"index;"`     // Update time
}

func (a *RoleMenu) TableName() string {
	return config.C.FormatTableName("role_menu")
}

// RoleMenuQueryParam Defining the query parameters for the `RoleMenu` struct.
type RoleMenuQueryParam struct {
	util.PaginationParam
	RoleID string `form:"-"` // From Role.ID
}

// RoleMenuQueryOptions Defining the query options for the `RoleMenu` struct.
type RoleMenuQueryOptions struct {
	util.QueryOptions
}

// RoleMenuQueryResult Defining the query result for the `RoleMenu` struct.
type RoleMenuQueryResult struct {
	Data       RoleMenus
	PageResult *util.PaginationResult
}

// RoleMenus Defining the slice of `RoleMenu` struct.
type RoleMenus []*RoleMenu

// RoleMenuForm Defining the data structure for creating a `RoleMenu` struct.
type RoleMenuForm struct {
}

// Validate A validation function for the `RoleMenuForm` struct.
func (a *RoleMenuForm) Validate() error {
	return nil
}

func (a *RoleMenuForm) FillTo(roleMenu *RoleMenu) error {
	return nil
}
