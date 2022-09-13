package typed

import (
	"strings"
	"time"

	"github.com/LyricTian/gin-admin/v9/internal/x/utilx"
)

const (
	DictionaryPathDelimiter = "."
)

// Dictionary management for key/value pairs
type Dictionary struct {
	ID         string        `gorm:"size:20;primarykey;" json:"id"`
	Key        string        `gorm:"size:128;index;" json:"key"`         // Key of the dictionary (Unique key for the same parent)
	Value      *string       `gorm:"size:4096;" json:"value"`            // Value of the key
	Remark     *string       `gorm:"size:1024;" json:"remark"`           // Remark of the key
	ParentID   *string       `gorm:"size:20;index;" json:"parent_id"`    // Parent ID
	ParentPath *string       `gorm:"size:255;index;" json:"parent_path"` // Parent path (split by .)(max depth: 10)
	CreatedAt  time.Time     `gorm:"index;" json:"created_at"`
	CreatedBy  string        `gorm:"size:20;" json:"created_by"`
	UpdatedAt  time.Time     `gorm:"index;" json:"updated_at"`
	UpdatedBy  string        `gorm:"size:20;" json:"updated_by"`
	Children   *Dictionaries `gorm:"-" json:"children"`
}

type DictionaryQueryParam struct {
	utilx.PaginationParam
	Key                string   `form:"key"`
	QueryValue         string   `form:"queryValue"`
	ParentID           *string  `form:"parentID"`
	LikeLeftParentPath string   `form:"-"`
	IDs                []string `form:"-"`
}

type DictionaryQueryOptions struct {
	utilx.QueryOptions
}

type DictionaryQueryResult struct {
	Data       Dictionaries
	PageResult *utilx.PaginationResult
}

type Dictionaries []*Dictionary

func (a Dictionaries) SplitParentIDs() []string {
	parentIDList := make([]string, 0, len(a))
	mIDList := make(map[string]struct{})
	for _, item := range a {
		if _, ok := mIDList[item.ID]; ok {
			continue
		}
		mIDList[item.ID] = struct{}{}

		if item.ParentPath != nil {
			for _, pp := range strings.Split(*item.ParentPath, DictionaryPathDelimiter) {
				if pp == "" {
					continue
				}
				if _, ok := mIDList[pp]; ok {
					continue
				}
				parentIDList = append(parentIDList, pp)
				mIDList[pp] = struct{}{}
			}
		}
	}
	return parentIDList
}

func (a Dictionaries) ToMap() map[string]*Dictionary {
	m := make(map[string]*Dictionary)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

func (a Dictionaries) ToTree() Dictionaries {
	var list Dictionaries

	m := a.ToMap()
	for _, item := range a {
		if item.ParentID == nil || *item.ParentPath == "" {
			list = append(list, item)
			continue
		}

		if pitem, ok := m[*item.ParentID]; ok {
			if pitem.Children == nil {
				children := Dictionaries{item}
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		}
	}

	return list
}

type DictionaryCreate struct {
	Key      string `json:"key" binding:"required"`
	Value    string `json:"value"`
	Remark   string `json:"remark"`
	ParentID string `json:"parent_id"`
}

type DictionaryUpdate struct {
	Key    string `json:"key" binding:"required"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}
