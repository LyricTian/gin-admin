package bll

import (
	"context"
	"gin-admin/src/model"
	"gin-admin/src/schema"
	"gin-admin/src/util"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// Menu 菜单管理
type Menu struct {
	MenuModel model.IMenu `inject:"IMenu"`
	lock      sync.RWMutex
}

// QueryPage 查询分页数据
func (a *Menu) QueryPage(ctx context.Context, params schema.MenuQueryParam, pageIndex, pageSize uint) (int64, []*schema.MenuQueryResult, error) {
	return a.MenuModel.QueryPage(ctx, params, pageIndex, pageSize)
}

// QueryTree 查询菜单树
func (a *Menu) QueryTree(ctx context.Context, params schema.MenuSelectQueryParam) ([]map[string]interface{}, error) {
	items, err := a.MenuModel.QuerySelect(ctx, params)
	if err != nil {
		return nil, err
	}

	treeData := util.Slice2Tree(util.StructsToMapSlice(items), "record_id", "parent_id")
	return treeData, nil
}

// Get 查询指定数据
func (a *Menu) Get(ctx context.Context, recordID string) (*schema.Menu, error) {
	return a.MenuModel.Get(ctx, recordID)
}

// Create 创建数据
func (a *Menu) Create(ctx context.Context, item *schema.Menu) error {
	if item.Code != "" {
		exists, err := a.MenuModel.CheckCode(ctx, item.Code, item.Type, item.ParentID)
		if err != nil {
			return err
		} else if exists {
			return errors.New("菜单编号已经存在")
		}
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	levelCodes, err := a.MenuModel.QueryLevelCodesByParentID(item.ParentID)
	if err != nil {
		return err
	}

	levelCode := util.GetLevelCode(levelCodes)
	if len(levelCode) == 0 {
		return errors.New("无效的分级码")
	}

	item.LevelCode = levelCode
	item.ID = 0
	item.RecordID = util.NewUUID()
	item.Created = time.Now().Unix()
	item.Deleted = 0
	return a.MenuModel.Create(ctx, item)
}

// Update 更新数据
func (a *Menu) Update(ctx context.Context, recordID string, item *schema.Menu) error {
	oldItem, err := a.MenuModel.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.New("无效的参数")
	} else if item.Code != "" && item.Code != oldItem.Code {
		exists, err := a.MenuModel.CheckCode(ctx, item.Code, item.Type, item.ParentID)
		if err != nil {
			return err
		} else if exists {
			return errors.New("菜单编号已经存在")
		}
	}

	info := util.StructToMap(item)
	delete(info, "id")
	delete(info, "record_id")
	delete(info, "level_code")
	delete(info, "creator")
	delete(info, "created")
	delete(info, "deleted")

	if item.ParentID != oldItem.ParentID {
		a.lock.Lock()
		defer a.lock.Unlock()

		levelCodes, err := a.MenuModel.QueryLevelCodesByParentID(item.ParentID)
		if err != nil {
			return err
		}

		levelCode := util.GetLevelCode(levelCodes)
		if len(levelCode) == 0 {
			return errors.New("无效的分级码")
		}

		return a.MenuModel.UpdateWithLevelCode(ctx, recordID, info, oldItem.LevelCode, levelCode)
	}

	return a.MenuModel.Update(ctx, recordID, info)
}

// Delete 删除数据
func (a *Menu) Delete(ctx context.Context, recordID string) error {
	exists, err := a.MenuModel.CheckChild(ctx, recordID)
	if err != nil {
		return err
	} else if exists {
		return errors.New("含有子级菜单，不能删除")
	}

	return a.MenuModel.Delete(ctx, recordID)
}

// UpdateStatus 更新状态
func (a *Menu) UpdateStatus(ctx context.Context, recordID string, status int) error {
	info := map[string]interface{}{
		"status": status,
	}
	return a.MenuModel.Update(ctx, recordID, info)
}
