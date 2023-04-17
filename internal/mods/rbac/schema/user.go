package schema

import (
	"time"

	"github.com/LyricTian/gin-admin/v10/internal/utils"
)

// User management for RBAC
type User struct {
	ID        string    `json:"id" gorm:"size:20;primarykey;"` // Unique ID
	Username  string    `json:"username" gorm:"size:64;index"` // Username for login
	Name      string    `json:"name" gorm:"size:64;index"`     // Name of user
	Password  string    `json:"password" gorm:"size:64;"`      // Password for login (encrypted)
	Phone     string    `json:"phone" gorm:"size:32;"`         // Phone number of user
	Email     string    `json:"email" gorm:"size:128;"`        // Email of user
	Remark    string    `json:"remark" gorm:"size:1024;"`      // Remark of user
	Status    string    `json:"status" gorm:"size:20;index"`   // Status of user (activated, freezed)
	CreatedAt time.Time `json:"created_at" gorm:"index;"`      // Create time
	UpdatedAt time.Time `json:"updated_at" gorm:"index;"`      // Update time
}

// Defining the query parameters for the `User` struct.
type UserQueryParam struct {
	utils.PaginationParam
	LikeUsername string `form:"username"` // Username for login
	LikeName     string `form:"name"`     // Name of user
	Status       string `form:"status"`   // Status of user (activated, freezed)
}

// Defining the query options for the `User` struct.
type UserQueryOptions struct {
	utils.QueryOptions
}

// Defining the query result for the `User` struct.
type UserQueryResult struct {
	Data       Users
	PageResult *utils.PaginationResult
}

// Defining the slice of `User` struct.
type Users []*User

// Defining the data structure for creating a `User` struct.
type UserForm struct {
	Username string `json:"username" binding:"required,max=64"`                // Username for login
	Name     string `json:"name" binding:"required,max=64"`                    // Name of user
	Password string `json:"password"`                                          // Password for login (encrypted)
	Phone    string `json:"phone"`                                             // Phone number of user
	Email    string `json:"email"`                                             // Email of user
	Remark   string `json:"remark"`                                            // Remark of user
	Status   string `json:"status" binding:"required,oneof=activated freezed"` // Status of user (activated, freezed)
}

// A validation function for the `UserForm` struct.
func (a *UserForm) Validate() error {
	return nil
}

func (a *UserForm) FillTo(user *User) *User {
	user.Username = a.Username
	user.Name = a.Name
	user.Password = a.Password
	user.Phone = a.Phone
	user.Email = a.Email
	user.Remark = a.Remark
	user.Status = a.Status
	return user
}
