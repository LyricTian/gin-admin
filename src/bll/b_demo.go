package bll

import (
	"gin-admin/src/model"
	"gin-admin/src/schema"
)

// Demo 示例
type Demo struct {
	DemoModel model.IDemo `inject:"IDemo"`
}

// Query 查询示例数据
func (a *Demo) Query() ([]*schema.Demo, error) {
	return a.DemoModel.Query()
}

// Get 获取单条示例数据
func (a *Demo) Get(id int64) (*schema.Demo, error) {
	return a.DemoModel.Get(id)
}

// Create 增加示例数据
func (a *Demo) Create(item *schema.Demo) error {
	return a.DemoModel.Create(item)
}

// Update 更新示例数据
func (a *Demo) Update(id int64, item *schema.Demo) error {
	info := map[string]interface{}{
		"code": item.Code,
		"name": item.Name,
	}

	return a.DemoModel.Update(id, info)
}

// Delete 删除示例数据
func (a *Demo) Delete(id int64) error {
	return a.DemoModel.Delete(id)
}
