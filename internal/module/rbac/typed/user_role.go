package typed

import (
	"time"

	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
)

// User authority management
type UserRole struct {
	ID        string    `gorm:"size:20;primarykey;" json:"id"`
	UserID    string    `gorm:"size:20;index;" json:"user_id"`
	RoleID    string    `gorm:"size:20;index;" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	Role      *Role     `gorm:"-" json:"role"`
}

type UserRoleQueryParam struct {
	utilx.PaginationParam
	RoleID     string   `json:"-"`
	UserID     string   `form:"-"`
	UserIDList []string `form:"-"`
}

type UserRoleQueryOptions struct {
	utilx.QueryOptions
}

type UserRoleQueryResult struct {
	Data       UserRoles
	PageResult *utilx.PaginationResult
}

type UserRoles []*UserRole

func (a UserRoles) ToRoleIDs() []string {
	var ids []string
	m := make(map[string]bool)
	for _, v := range a {
		if _, ok := m[v.RoleID]; ok {
			continue
		}
		ids = append(ids, v.RoleID)
		m[v.RoleID] = true
	}
	return ids
}

func (a UserRoles) FillRole(roleMap map[string]*Role) {
	for _, v := range a {
		v.Role = roleMap[v.RoleID]
	}
}

func (a UserRoles) ToUserIDMap() map[string]UserRoles {
	m := make(map[string]UserRoles)
	for _, v := range a {
		m[v.UserID] = append(m[v.UserID], v)
	}
	return m
}
