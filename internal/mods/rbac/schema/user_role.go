package schema

import (
	"time"

	"github.com/LyricTian/gin-admin/v10/internal/utils"
)

// User roles for RBAC
type UserRole struct {
	ID        string    `json:"id" gorm:"size:20;primarykey"` // Unique ID
	UserID    string    `json:"user_id" gorm:"size:20;index"` // From User.ID
	RoleID    string    `json:"role_id" gorm:"size:20;index"` // From Role.ID
	CreatedAt time.Time `json:"created_at" gorm:"index;"`     // Create time
	UpdatedAt time.Time `json:"updated_at" gorm:"index;"`     // Update time
	RoleName  string    `json:"role_name" gorm:"<-:false"`    // From Role.Name
}

func (a UserRole) TableName() string {
	return "user_role"
}

// Defining the query parameters for the `UserRole` struct.
type UserRoleQueryParam struct {
	utils.PaginationParam
	InUserIDs []string `form:"-"` // From User.ID
	UserID    string   `form:"-"` // From User.ID
	RoleID    string   `form:"-"` // From Role.ID
}

// Defining the query options for the `UserRole` struct.
type UserRoleQueryOptions struct {
	utils.QueryOptions
	JoinRole bool // Join role table
}

// Defining the query result for the `UserRole` struct.
type UserRoleQueryResult struct {
	Data       UserRoles
	PageResult *utils.PaginationResult
}

// Defining the slice of `UserRole` struct.
type UserRoles []*UserRole

func (a UserRoles) ToUserIDMap() map[string]UserRoles {
	m := make(map[string]UserRoles)
	for _, userRole := range a {
		m[userRole.UserID] = append(m[userRole.UserID], userRole)
	}
	return m
}

// Defining the data structure for creating a `UserRole` struct.
type UserRoleForm struct {
}

// A validation function for the `UserRoleForm` struct.
func (a *UserRoleForm) Validate() error {
	return nil
}

func (a *UserRoleForm) FillTo(userRole *UserRole) error {
	return nil
}
