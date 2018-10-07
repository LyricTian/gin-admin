package model

import (
	"gin-admin/src/schema"
)

// IDemo 示例接口
type IDemo interface {
	Query() ([]*schema.Demo, error)
	Get(id int64) (*schema.Demo, error)
	Create(item *schema.Demo) error
	Update(id int64, info map[string]interface{}) error
	Delete(id int64) error
}
