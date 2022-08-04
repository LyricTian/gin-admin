package typed

import (
	"time"

	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
)

type UserStatus string

const (
	UserStatusActivated UserStatus = "activated"
	UserStatusFreezed   UserStatus = "freezed"
)

// User management
type User struct {
	ID        string     `gorm:"size:20;primarykey;" json:"id"`
	Username  string     `gorm:"size:64;index;" json:"username"` // Login username (must be unique)
	Name      string     `gorm:"size:64;" json:"name"`           // Name of the user
	Password  string     `gorm:"size:60;" json:"-"`              // Password of the user
	Email     *string    `gorm:"size:64;" json:"email"`
	Phone     *string    `gorm:"size:20;" json:"phone"`
	Remark    *string    `gorm:"size:1024;" json:"remark"`
	Status    UserStatus `gorm:"size:20;index;" json:"status"`
	CreatedAt time.Time  `gorm:"index;" json:"created_at"`
	CreatedBy string     `gorm:"size:20;" json:"created_by"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdatedBy string     `gorm:"size:20;" json:"updated_by"`
	UserRoles UserRoles  `gorm:"-" json:"user_roles"`
}

type UserQueryParam struct {
	utilx.PaginationParam
	Username     string     `form:"-"`
	Status       UserStatus `form:"status"`
	LikeUsername string     `form:"username"`
	LikeName     string     `form:"name"`
	RoleID       string     `form:"roleID"`
}

type UserQueryOptions struct {
	utilx.QueryOptions
}

type UserQueryResult struct {
	Data       Users
	PageResult *utilx.PaginationResult
}

type Users []*User

func (a Users) ToIDs() []string {
	var idList []string
	for _, item := range a {
		idList = append(idList, item.ID)
	}
	return idList
}

func (a Users) FillUserRoles(m map[string]UserRoles) {
	for _, item := range a {
		item.UserRoles = m[item.ID]
	}
}

type UserCreate struct {
	Username string   `json:"username" binding:"required"`
	Name     string   `json:"name" binding:"required"`
	Password string   `json:"password"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Remark   string   `json:"remark"`
	RoleIDs  []string `json:"role_ids"`
}
