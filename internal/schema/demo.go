package schema

import "time"

type Demo struct {
	ID          string    `gorm:"size:20;primarykey;" json:"id"`
	Code        string    `gorm:"size:50;index;" json:"code" binding:"required"`
	Name        string    `gorm:"size:50;index;" json:"name" binding:"required"`
	Description string    `gorm:"size:1024;" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Query parameters for db
type DemoQueryParam struct {
	PaginationParam
}

// Query options for db (order or select fields)
type DemoQueryOptions struct {
	OrderFields  []*OrderField
	SelectFields []string
}

// Query result from db
type DemoQueryResult struct {
	Data       Demos
	PageResult *PaginationResult
}

// Object List
type Demos []*Demo
