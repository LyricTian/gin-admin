package typed

import (
	"time"

	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
)

// Menu action management (button permissions)
type MenuAction struct {
	ID        string              `gorm:"size:20;primarykey;" json:"id"`
	MenuID    string              `gorm:"size:20;index;" json:"menu_id"`
	Code      string              `gorm:"size:64;" json:"code"`
	Name      string              `gorm:"size:64;" json:"name"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Resources MenuActionResources `gorm:"-" json:"resources"`
}

type MenuActionQueryParam struct {
	utilx.PaginationParam
	MenuID string `form:"-"`
	UserID string `form:"-"`
}

type MenuActionQueryOptions struct {
	utilx.QueryOptions
}

type MenuActionQueryResult struct {
	Data       MenuActions
	PageResult *utilx.PaginationResult
}

type MenuActions []*MenuAction

func (a MenuActions) ToMap() map[string]*MenuAction {
	m := make(map[string]*MenuAction)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

func (a MenuActions) FillResources(m map[string]MenuActionResources) MenuActions {
	for _, v := range a {
		v.Resources = m[v.ID]
	}
	return a
}

func (a MenuActions) ToMenuIDMap() map[string]MenuActions {
	m := make(map[string]MenuActions)
	for _, v := range a {
		m[v.MenuID] = append(m[v.MenuID], v)
	}
	return m
}
