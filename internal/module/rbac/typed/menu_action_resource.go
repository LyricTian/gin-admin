package typed

import (
	"time"

	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
)

// Menu action resource management (router permissions for casbin)
type MenuActionResource struct {
	ID        string    `gorm:"size:20;primarykey;" json:"id"`
	MenuID    string    `gorm:"size:20;index;" json:"menu_id"`
	ActionID  string    `gorm:"size:20;index;" json:"action_id"`
	Method    string    `gorm:"size:20;" json:"method"`
	Path      string    `gorm:"size:512;" json:"path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MenuActionResourceQueryParam struct {
	utilx.PaginationParam
	MenuID string `form:"-"`
	RoleID string `form:"-"`
}

type MenuActionResourceQueryOptions struct {
	utilx.QueryOptions
}

type MenuActionResourceQueryResult struct {
	Data       MenuActionResources
	PageResult *utilx.PaginationResult
}

type MenuActionResources []*MenuActionResource

func (a MenuActionResources) ToMap() map[string]*MenuActionResource {
	m := make(map[string]*MenuActionResource)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

func (a MenuActionResources) ToActionIDMap() map[string]MenuActionResources {
	m := make(map[string]MenuActionResources)
	for _, v := range a {
		m[v.ActionID] = append(m[v.ActionID], v)
	}
	return m
}
