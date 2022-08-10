package typed

import (
	"strings"
	"time"

	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
)

type MenuStatus string

const (
	MenuStatusEnabled  MenuStatus = "enabled"
	MenuStatusDisabled MenuStatus = "disabled"
)

const MenuParentPathDelimiter = "."

// Menu management
type Menu struct {
	ID         string      `gorm:"size:20;primarykey;" json:"id"`
	Name       string      `gorm:"size:64;" json:"name"`
	Sequence   int         `gorm:"index;" json:"sequence"`
	Icon       *string     `gorm:"size:128;" json:"icon"`
	Link       *string     `gorm:"size:255;" json:"link"`
	ParentID   *string     `gorm:"size:20;" json:"parent_id"`
	ParentPath *string     `gorm:"size:512;" json:"parent_path"` // parent path (split by '.')
	Remark     *string     `gorm:"size:1024;" json:"remark"`
	Hide       bool        `json:"hide"`
	Status     MenuStatus  `gorm:"size:20;index;" json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	CreatedBy  string      `gorm:"size:20;" json:"created_by"`
	UpdatedAt  time.Time   `json:"updated_at"`
	UpdatedBy  string      `gorm:"size:20;" json:"updated_by"`
	Actions    MenuActions `gorm:"-" json:"actions"`
	Children   *Menus      `gorm:"-" json:"children"`
}

func (a *Menu) GetParentID() string {
	if a.ParentID == nil {
		return ""
	}
	return *a.ParentID
}

func (a *Menu) GetParentPath() string {
	if a.ParentPath == nil {
		return ""
	}
	return *a.ParentPath
}

func (a *Menu) GetIcon() string {
	if a.Icon == nil {
		return ""
	}
	return *a.Icon
}

func (a *Menu) GetLink() string {
	if a.Link == nil {
		return ""
	}
	return *a.Link
}

func (a *Menu) GetRemark() string {
	if a.Remark == nil {
		return ""
	}
	return *a.Remark
}

func (a *Menu) GenerateParentPath() *string {
	parentPath := a.GetParentPath() + a.ID + MenuParentPathDelimiter
	return &parentPath
}

type MenuQueryParam struct {
	utilx.PaginationParam
	IDList           []string   `form:"-"`
	ParentPathPrefix string     `form:"-"`
	UserID           string     `form:"-"`
	ParentID         *string    `form:"parentID"`
	LikeName         string     `form:"name"`
	Status           MenuStatus `form:"status"`
}

type MenuQueryOptions struct {
	utilx.QueryOptions
}

type MenuQueryResult struct {
	Data       Menus
	PageResult *utilx.PaginationResult
}

type Menus []*Menu

func (a Menus) Len() int {
	return len(a)
}

func (a Menus) Less(i, j int) bool {
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
	idList := make([]string, 0, len(a))
	mIDList := make(map[string]struct{})
	for _, item := range a {
		if _, ok := mIDList[item.ID]; ok {
			continue
		}
		idList = append(idList, item.ID)
		mIDList[item.ID] = struct{}{}

		if item.ParentPath != nil {
			for _, pp := range strings.Split(*item.ParentPath, MenuParentPathDelimiter) {
				if _, ok := mIDList[pp]; ok {
					continue
				}
				idList = append(idList, pp)
				mIDList[pp] = struct{}{}
			}
		}
	}
	return idList
}

func (a Menus) ToTree() Menus {
	var list Menus

	m := a.ToMap()
	for _, item := range a {
		if item.ParentID == nil || *item.ParentPath == "" {
			list = append(list, item)
			continue
		}

		if pitem, ok := m[*item.ParentID]; ok {
			if pitem.Children == nil {
				children := Menus{item}
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		}
	}

	return list
}

func (a Menus) FillActions(m map[string]MenuActions) {
	for _, item := range a {
		item.Actions = m[item.ID]
	}
}

type MenuCreate struct {
	Name     string      `json:"name" binding:"required"`
	Sequence int         `json:"sequence" binding:"required"`
	Icon     string      `json:"icon"`
	Link     string      `json:"link"`
	ParentID string      `json:"parent_id"`
	Remark   string      `json:"remark"`
	Hide     bool        `json:"hide"`
	Actions  MenuActions `json:"actions"`
}

type MenuUpdateStatus struct {
	Status MenuStatus `json:"status" binding:"required"` // enabled/disabled
}
