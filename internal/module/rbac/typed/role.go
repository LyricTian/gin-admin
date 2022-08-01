package typed

import (
	"time"

	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
)

type RoleStatus string

const (
	RoleStatusEnabled  RoleStatus = "enabled"
	RoleStatusDisabled RoleStatus = "disabled"
)

// Role management
type Role struct {
	ID        string     `gorm:"size:20;primarykey;" json:"id"`
	Name      string     `gorm:"size:64;" json:"name"`
	Sequence  int        `gorm:"index;" json:"sequence"`
	Remark    *string    `gorm:"size:1024;" json:"remark"`
	Status    RoleStatus `gorm:"size:20;index;" json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy string     `gorm:"size:20;" json:"created_by"`
	UpdatedAt time.Time  `gorm:"index;" json:"updated_at"`
	UpdatedBy string     `gorm:"size:20;" json:"updated_by"`
	RoleMenus RoleMenus  `gorm:"-" json:"role_menus"`
}

type RoleQueryParam struct {
	utilx.PaginationParam
	IDList      []string   `form:"-"`
	GtUpdatedAt *time.Time `form:"-"`
	Status      RoleStatus `form:"status"`
	LikeName    string     `form:"name"`
	Result      string     `form:"result"` // select/paginate
}

type RoleQueryOptions struct {
	utilx.QueryOptions
}

type RoleQueryResult struct {
	Data       Roles
	PageResult *utilx.PaginationResult
}

type Roles []*Role

func (a Roles) ToMap() map[string]*Role {
	var m = make(map[string]*Role)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

type RoleCreate struct {
	Name      string    `json:"name" binding:"required"`
	Sequence  int       `json:"sequence" binding:"required"`
	Remark    string    `json:"remark"`
	RoleMenus RoleMenus `json:"menus"`
}
