package typed

import (
	"time"

	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
)

// Role authority management
type RoleMenu struct {
	ID           string    `gorm:"size:20;primarykey;" json:"id"`
	RoleID       string    `gorm:"size:20;index;" json:"role_id"`
	MenuID       string    `gorm:"size:20;" json:"menu_id"`
	MenuActionID string    `gorm:"size:20;" json:"menu_action_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type RoleMenuQueryParam struct {
	utilx.PaginationParam
	RoleID string `form:"-"`
}

type RoleMenuQueryOptions struct {
	utilx.QueryOptions
}

type RoleMenuQueryResult struct {
	Data       RoleMenus
	PageResult *utilx.PaginationResult
}

type RoleMenus []*RoleMenu
