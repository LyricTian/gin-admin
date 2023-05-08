package schema

import (
	"strings"
	"time"

	"github.com/LyricTian/gin-admin/v10/pkg/util"
)

const (
	MenuStatusDisabled = "disabled"
	MenuStatusEnabled  = "enabled"
)

var (
	MenusOrderParams = []util.OrderByParam{
		{Field: "sequence", Direction: util.DESC},
		{Field: "created_at", Direction: util.DESC},
	}
)

// Menu management for RBAC
type Menu struct {
	ID          string        `json:"id" gorm:"size:20;primarykey;"`      // Unique ID
	Code        string        `json:"code" gorm:"size:32;index;"`         // Code of menu (unique for each level)
	Name        string        `json:"name" gorm:"size:128;index"`         // Display name of menu
	Description string        `json:"description" gorm:"size:1024"`       // Details about menu
	Sequence    int           `json:"sequence" gorm:"index;"`             // Sequence for sorting
	Type        string        `json:"type" gorm:"size:20;index"`          // Type of menu (group, page, button)
	Path        string        `json:"path" gorm:"size:255;"`              // Access path of menu
	Properties  string        `json:"properties" gorm:"type:text;"`       // Properties of menu (JSON)
	Status      string        `json:"status" gorm:"size:20;index"`        // Status of menu (disabled, enabled)
	ParentID    string        `json:"parent_id" gorm:"size:20;index;"`    // Parent ID (From Menu.ID)
	ParentPath  string        `json:"parent_path" gorm:"size:255;index;"` // Parent path (split by .)
	Children    *Menus        `json:"children" gorm:"-"`                  // Child menus
	CreatedAt   time.Time     `json:"created_at" gorm:"index;"`           // Create time
	UpdatedAt   time.Time     `json:"updated_at" gorm:"index;"`           // Update time
	Resources   MenuResources `json:"resources" gorm:"-"`                 // Resources of menu
}

func (a Menu) TableName() string {
	return "menu"
}

// Defining the query parameters for the `Menu` struct.
type MenuQueryParam struct {
	util.PaginationParam
	InIDs            []string `form:"-"`
	LikeName         string   `form:"name"` // Display name of menu
	Status           string   `form:"-"`    // Status of menu (disabled, enabled)
	ParentID         string   `form:"-"`    // Parent ID (From Menu.ID)
	ParentPathPrefix string   `form:"-"`    // Parent path (split by .)
	UserID           string   `form:"-"`    // User ID
	RoleID           string   `form:"-"`    // Role ID
}

// Defining the query options for the `Menu` struct.
type MenuQueryOptions struct {
	util.QueryOptions
}

// Defining the query result for the `Menu` struct.
type MenuQueryResult struct {
	Data       Menus
	PageResult *util.PaginationResult
}

// Defining the slice of `Menu` struct.
type Menus []*Menu

func (a Menus) Len() int {
	return len(a)
}

func (a Menus) Less(i, j int) bool {
	if a[i].Sequence == a[j].Sequence {
		return a[i].CreatedAt.Unix() > a[j].CreatedAt.Unix()
	}
	return a[i].Sequence > a[j].Sequence
}

func (a Menus) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Menus) ToMap() map[string]*Menu {
	m := make(map[string]*Menu)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

func (a Menus) SplitParentIDs() []string {
	parentIDs := make([]string, 0, len(a))
	idMapper := make(map[string]struct{})
	for _, item := range a {
		if _, ok := idMapper[item.ID]; ok {
			continue
		}
		idMapper[item.ID] = struct{}{}
		if pp := item.ParentPath; pp != "" {
			for _, pid := range strings.Split(pp, util.TreePathDelimiter) {
				if pid == "" {
					continue
				}
				if _, ok := idMapper[pid]; ok {
					continue
				}
				parentIDs = append(parentIDs, pid)
				idMapper[pid] = struct{}{}
			}
		}
	}
	return parentIDs
}

func (a Menus) ToTree() Menus {
	var list Menus
	m := a.ToMap()
	for _, item := range a {
		if item.ParentPath == "" {
			list = append(list, item)
			continue
		}
		if parent, ok := m[item.ParentPath]; ok {
			if parent.Children == nil {
				children := Menus{item}
				parent.Children = &children
				continue
			}
			*parent.Children = append(*parent.Children, item)
		}
	}
	return list
}

// Defining the data structure for creating a `Menu` struct.
type MenuForm struct {
	Code        string        `json:"code" binding:"required,max=32"`                   // Code of menu (unique for each level)
	Name        string        `json:"name" binding:"required,max=128"`                  // Display name of menu
	Description string        `json:"description"`                                      // Details about menu
	Sequence    int           `json:"sequence"`                                         // Sequence for sorting
	Type        string        `json:"type" binding:"required,oneof=group menu button"`  // Type of menu (group, page, button)
	Path        string        `json:"path"`                                             // Access path of menu
	Properties  string        `json:"properties"`                                       // Properties of menu (JSON)
	Status      string        `json:"status" binding:"required,oneof=disabled enabled"` // Status of menu (disabled, enabled)
	ParentID    string        `json:"parent_id"`                                        // Parent ID (From Menu.ID)
	Resources   MenuResources `json:"resources"`                                        // Resources of menu
}

// A validation function for the `MenuForm` struct.
func (a *MenuForm) Validate() error {
	return nil
}

func (a *MenuForm) FillTo(menu *Menu) error {
	menu.Code = a.Code
	menu.Name = a.Name
	menu.Description = a.Description
	menu.Sequence = a.Sequence
	menu.Type = a.Type
	menu.Path = a.Path
	menu.Properties = a.Properties
	menu.Status = a.Status
	menu.ParentID = a.ParentID
	return nil
}
